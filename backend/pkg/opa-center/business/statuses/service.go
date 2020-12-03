package statuses

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

const mainAuthorizationPrefix = "statuses"

type service struct {
	dao              daos.Dao
	validator        *validator.Validate
	partitionSvc     PartitionService
	authorizationSvc authorization.Service
}

func (s *service) MigrateDB(systemLogger log.Logger) error {
	systemLogger.Debug("Migrate database for Status")

	return s.dao.MigrateDB()
}

func (s *service) ManageRetention(logger log.Logger, retentionDuration time.Duration, partitionID string) error {
	// Get now date
	now := time.Now()
	// Remove duration
	oldDate := now.Add(-retentionDuration)
	// Format date
	oldDateS := oldDate.Format(time.RFC3339)

	return s.dao.Delete(&models.Filter{
		CreatedAt:   &common.DateFilter{Lt: &oldDateS},
		PartitionID: &common.GenericFilter{Eq: partitionID},
	})
}

func (s *service) FindByID(ctx context.Context, id string, projection *models.Projection) (*models.Status, error) {
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

func (s *service) UnsecureCreate(partitionID string, inp map[string]interface{}) error {
	// Find partition
	partition, err := s.partitionSvc.UnsecureFindByID(partitionID)
	// Check error
	if err != nil {
		return err
	}
	// Check if doesn't partition exist
	if partition == nil {
		return errors.New("partition doesn't exist")
	}

	bb, err := json.Marshal(inp)
	// Check error
	if err != nil {
		return err
	}

	// Create status object
	st := &models.Status{
		OriginalMessage: string(bb),
		PartitionID:     partition.ID,
	}

	// Validate input
	err = s.validator.Struct(st)
	// Check error
	if err != nil {
		return err
	}

	// Save status object
	err = s.dao.Save(st)
	// Check error
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllPaginated(
	ctx context.Context,
	partitionID string,
	page *pagination.PageInput,
	sort *models.SortOrder,
	filter *models.Filter,
	projection *models.Projection,
) ([]*models.Status, *pagination.PageOutput, error) {
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

	// Create filter if not exists
	if filter == nil {
		filter = &models.Filter{}
	}

	// Add partition id to filter
	filter.PartitionID = &common.GenericFilter{Eq: partitionID}

	return s.dao.GetAllPaginated(page, sort, filter, projection)
}
