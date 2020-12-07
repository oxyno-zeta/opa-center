package partitions

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
	"github.com/robfig/cron/v3"
)

const mainAuthorizationPrefix = "partitions"

var validNameRegex = regexp.MustCompile("^[^-][a-zA-Z0-9-]+[^-]$")

type service struct {
	dao                daos.Dao
	validator          *validator.Validate
	authorizationSvc   authorization.Service
	cfgManager         config.Manager
	opaCfgTemplate     *template.Template
	retentionScheduler *cron.Cron
	decisionLogsSvc    RetentionService
	statusesSvc        RetentionService
	logger             log.Logger
}

type opaCfgData struct {
	Partition  *models.Partition
	ServiceURL string
}

func (s *service) AddServices(decisionLogsSvc, statusesSvc RetentionService) {
	s.decisionLogsSvc = decisionLogsSvc
	s.statusesSvc = statusesSvc
}

func (s *service) Initialize() error {
	return s.initializeCron(true)
}

func (s *service) initializeCron(isStartup bool) error {
	// Get configuration
	cfg := s.cfgManager.GetConfig().Center

	// Create new cron
	c := cron.New(
		cron.WithLogger(s.logger.GetCronLogger()),
	)
	// Create task
	task := &RetentionCleanTask{s: s, logger: s.logger.WithField("task", "retention-clean-process")}
	// Add task
	_, err := c.AddJob(cfg.CronRetentionProcess, task)
	// Check error
	if err != nil {
		return err
	}

	// Start cron
	c.Start()

	// Check if a startup run is asked
	if isStartup && !cfg.SkipRetentionProcessAtStartup {
		// Start go routine to start task
		go task.Run()
	}

	// Store scheduler
	s.retentionScheduler = c

	return nil
}

func (s *service) Reload() error {
	// Stop scheduler
	s.retentionScheduler.Stop()

	// Restart initialize cron
	return s.initializeCron(false)
}

func (s *service) MigrateDB(systemLogger log.Logger) error {
	systemLogger.Debug("Migrate database for Partitions")

	return s.dao.MigrateDB()
}

func (s *service) CheckAuthenticated(ctx context.Context, partitionID, authorizationHeader string) error {
	// Get logger
	logger := log.GetLoggerFromContext(ctx)

	// Check that authentication header isn't empty
	if authorizationHeader == "" {
		logger.Error("No authorization header content detected")

		return errors.NewUnauthorizedError("unauthorized")
	}

	// Parse authentication header
	authorizationSplit := strings.Split(authorizationHeader, " ")
	// Check that first part is "Token"
	if authorizationSplit[0] != "Token" {
		logger.Error("Authorization header don't start with Token prefix")

		return errors.NewUnauthorizedError("unauthorized")
	}

	// Find partition with given id
	partition, err := s.dao.FindByID(partitionID, &models.Projection{ID: true, AuthorizationToken: true})
	// Check error
	if err != nil {
		return err
	}
	// Check if partition is found
	if partition == nil {
		logger.Errorf("Partition with id %d cannot be found in database", partitionID)

		return errors.NewUnauthorizedError("unauthorized")
	}

	// Check if upload is authenticated
	if partition.AuthorizationToken != authorizationSplit[1] {
		logger.Error("Authorization token not equal to selected partition one")

		return errors.NewUnauthorizedError("unauthorized")
	}

	return nil
}

func (s *service) GetAllPaginated(
	ctx context.Context,
	page *pagination.PageInput,
	sort *models.SortOrder,
	filter *models.Filter,
	projection *models.Projection,
) ([]*models.Partition, *pagination.PageOutput, error) {
	// Check authorization
	err := s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:List", mainAuthorizationPrefix),
		"",
	)
	// Check error
	if err != nil {
		return nil, nil, err
	}

	return s.dao.GetAllPaginated(page, sort, filter, projection)
}

func (s *service) validateCreateInput(inp *models.CreateInput) error {
	// Validate input
	err := s.validator.Struct(inp)
	if err != nil {
		return err
	}

	// Validate name
	if !validNameRegex.MatchString(inp.Name) {
		return errors.NewInvalidInputError(fmt.Sprintf(errInvalidNameTemplate, validNameRegex.String()))
	}

	// Validate decision log retention duration
	if inp.DecisionLogRetention != "" {
		// Try to parse duration
		_, err := time.ParseDuration(inp.DecisionLogRetention)
		// Check error
		if err != nil {
			return errors.NewInvalidInputErrorWithError(err)
		}
	}

	// Validate status retention duration
	if inp.StatusDataRetention != "" {
		// Try to parse duration
		_, err := time.ParseDuration(inp.StatusDataRetention)
		// Check error
		if err != nil {
			return errors.NewInvalidInputErrorWithError(err)
		}
	}

	return nil
}

func (s *service) Create(ctx context.Context, inp *models.CreateInput) (*models.Partition, error) {
	// Get logger from context
	logger := log.GetLoggerFromContext(ctx)

	// Validate input
	err := s.validateCreateInput(inp)
	if err != nil {
		return nil, err
	}

	// Check authorization
	err = s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:Create", mainAuthorizationPrefix),
		fmt.Sprintf("%s:%s", mainAuthorizationPrefix, inp.Name),
	)
	// Check error
	if err != nil {
		return nil, err
	}

	// Create new uuid as "authorization token"
	uuid, err := uuid.NewV4()
	// Check error
	if err != nil {
		return nil, err
	}

	// Create partition object
	obj := &models.Partition{
		Name:                 inp.Name,
		DecisionLogRetention: inp.DecisionLogRetention,
		StatusDataRetention:  inp.StatusDataRetention,
		AuthorizationToken:   uuid.String(),
	}

	// Search if it already exists
	logger.Infof("Check if partition %s already exists in database", obj.Name)
	dbE, err := s.dao.FindByName(obj.Name, &models.Projection{ID: true})
	// Check error
	if err != nil {
		return nil, err
	}
	// Check if item already exists in database
	if dbE != nil {
		return nil, errors.NewConflictError(fmt.Sprintf("partition with name %s already exists", obj.Name))
	}

	// Call dao
	res, err := s.dao.Save(obj)
	// Check error
	if err != nil {
		return nil, err
	}

	// Log
	logger.Infof("Partition %s successfully created", res.Name)

	return res, nil
}

func (s *service) validateUpdateInput(inp *models.UpdateInput) error {
	// Validate input
	err := s.validator.Struct(inp)
	if err != nil {
		return err
	}

	// Validate decision log retention duration
	if inp.DecisionLogRetention != nil && *inp.DecisionLogRetention != "" {
		// Try to parse duration
		_, err := time.ParseDuration(*inp.DecisionLogRetention)
		// Check error
		if err != nil {
			return errors.NewInvalidInputErrorWithError(err)
		}
	}

	// Validate status retention duration
	if inp.StatusDataRetention != nil && *inp.StatusDataRetention != "" {
		// Try to parse duration
		_, err := time.ParseDuration(*inp.StatusDataRetention)
		// Check error
		if err != nil {
			return errors.NewInvalidInputErrorWithError(err)
		}
	}

	return nil
}

func (s *service) Update(ctx context.Context, inp *models.UpdateInput) (*models.Partition, error) {
	// Validate input
	err := s.validateUpdateInput(inp)
	// Check error
	if err != nil {
		return nil, err
	}

	// Find partition
	res, err := s.dao.FindByID(inp.ID, nil)
	// Check error
	if err != nil {
		return nil, err
	}
	// Check if exists
	if res == nil {
		return nil, errors.NewNotFoundError("partition not found")
	}

	// Check authorization
	err = s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:Update", mainAuthorizationPrefix),
		fmt.Sprintf("%s:%s", mainAuthorizationPrefix, res.Name),
	)
	// Check error
	if err != nil {
		return nil, err
	}

	// Update only necessary fields

	// Store if something changed
	edited := false

	// Check if decision log retention is set
	if inp.DecisionLogRetention != nil {
		res.DecisionLogRetention = *inp.DecisionLogRetention
		edited = true
	}

	// Check if status retention is set
	if inp.StatusDataRetention != nil {
		res.StatusDataRetention = *inp.StatusDataRetention
		edited = true
	}

	// Check if nothing was edited
	if !edited {
		return res, nil
	}

	// Save
	res, err = s.dao.Save(res)
	// Check error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UnsecureFindByID(id string) (*models.Partition, error) {
	return s.dao.FindByID(id, nil)
}

func (s *service) FindByID(ctx context.Context, id string, projection *models.Projection) (*models.Partition, error) {
	// TODO Change this to a better solution
	// Check authorization
	err := s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:FindByID", mainAuthorizationPrefix),
		fmt.Sprintf("%s:%s", mainAuthorizationPrefix, id),
	)
	// Check error
	if err != nil {
		return nil, err
	}

	return s.dao.FindByID(id, projection)
}

func (s *service) GenerateOPAConfiguration(ctx context.Context, id string) (string, error) {
	// Check authorization
	err := s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:GenerateOPAConfiguration", mainAuthorizationPrefix),
		fmt.Sprintf("%s:%s", mainAuthorizationPrefix, id),
	)
	// Check error
	if err != nil {
		return "", err
	}

	// Get partition from id
	partition, err := s.dao.FindByID(id, nil)
	// Check error
	if err != nil {
		return "", err
	}

	// Get configuration
	cfg := s.cfgManager.GetConfig()
	// Create service url
	serviceURL := path.Join(cfg.Center.BaseURL, "/api/")

	// Create template data
	data := opaCfgData{
		Partition:  partition,
		ServiceURL: serviceURL,
	}

	// Generate configuration
	// Generate template in buffer
	buf := &bytes.Buffer{}
	// Execute template
	err = s.opaCfgTemplate.Execute(buf, data)
	// Check error
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
