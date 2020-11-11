package decisionlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

const mainAuthorizationPrefix = "decisionlogs"

type service struct {
	dao              daos.Dao
	validator        *validator.Validate
	partitionSvc     PartitionService
	authorizationSvc authorization.Service
}

func (s *service) MigrateDB(systemLogger log.Logger) error {
	systemLogger.Debug("Migrate database for Decision Logs")

	return s.dao.MigrateDB()
}

func (s *service) FindByIDOrDecisionID(ctx context.Context, id, did *string, projection *models.Projection) (*models.DecisionLog, error) {
	// Check if we are in the find by id case
	if id != nil {
		return s.findByID(ctx, *id, projection)
	}

	// Other case
	return s.findByDecisionID(ctx, *did, projection)
}

func (s *service) findByDecisionID(ctx context.Context, did string, projection *models.Projection) (*models.DecisionLog, error) {
	// Check authorization
	err := s.authorizationSvc.CheckAuthorized(
		ctx,
		fmt.Sprintf("%s:FindByDecisionID", mainAuthorizationPrefix),
		fmt.Sprintf("%s:%s", mainAuthorizationPrefix, did),
	)
	// Check error
	if err != nil {
		return nil, err
	}

	return s.dao.FindOneByDecisionID(did, projection)
}

func (s *service) findByID(ctx context.Context, id string, projection *models.Projection) (*models.DecisionLog, error) {
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

func (s *service) UnsecureCreate(partitionID string, inp []map[string]interface{}) error {
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

	// Loop over inp
	for i := 0; i < len(inp); i++ {
		data := inp[i]

		bb, err := json.Marshal(data)
		// Check error
		if err != nil {
			return err
		}

		// Create decision logs object
		dl := &models.DecisionLog{
			OriginalMessage: string(bb),
			PartitionID:     partition.ID,
		}
		// Add decision id
		if data["decision_id"] != nil {
			dl.DecisionID = data["decision_id"].(string)
		}
		// Add path
		if data["path"] != nil {
			dl.Path = data["path"].(string)
		}
		// Add requested by
		if data["requested_by"] != nil {
			dl.RequestedBy = data["requested_by"].(string)
		}
		// Add timestamp
		if data["timestamp"] != nil {
			tiStr := data["timestamp"].(string)

			ti, err2 := time.Parse(time.RFC3339, tiStr)
			// Check error
			if err2 != nil {
				return err2
			}

			dl.Timestamp = ti
		}

		// Validate input
		err = s.validator.Struct(dl)
		// Check error
		if err != nil {
			return err
		}

		// Check if decision logs already exists in db
		item, err := s.dao.FindOneByDecisionID(dl.DecisionID, nil)
		// Check error
		if err != nil {
			return err
		}
		// Check if item exists
		if item != nil {
			// Skip
			continue
		}

		// Save decision log object
		err = s.dao.Save(dl)
		// Check error
		if err != nil {
			return err
		}
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
) ([]*models.DecisionLog, *pagination.PageOutput, error) {
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
