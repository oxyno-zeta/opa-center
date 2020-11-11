package partitions

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"regexp"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/config"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

const mainAuthorizationPrefix = "partitions"

var validNameRegex = regexp.MustCompile("^[^-][a-zA-Z0-9-]+[^-]$")

type service struct {
	dao              daos.Dao
	validator        *validator.Validate
	authorizationSvc authorization.Service
	cfgManager       config.Manager
	opaCfgTemplate   *template.Template
}

type opaCfgData struct {
	Partition  *models.Partition
	ServiceURL string
}

func (s *service) MigrateDB(systemLogger log.Logger) error {
	systemLogger.Debug("Migrate database for Partitions")

	return s.dao.MigrateDB()
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

	// Create partition object
	obj := &models.Partition{
		Name: inp.Name,
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
