package decisionlogs

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	pmodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

type Service interface {
	// Migrate database
	MigrateDB(systemLogger log.Logger) error
	// Create decision log used internally only
	UnsecureCreate(partitionName string, inp []map[string]interface{}) error
	// Get data paginated
	GetAllPaginated(
		ctx context.Context,
		partitionID string,
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.DecisionLog, *pagination.PageOutput, error)
	// Find by id or decision id
	FindByIDOrDecisionID(ctx context.Context, id, did *string, projection *models.Projection) (*models.DecisionLog, error)
	// Manage retention data
	ManageRetention(logger log.Logger, retentionDuration time.Duration, partitionID string) error
}

type PartitionService interface {
	UnsecureFindByID(id string) (*pmodels.Partition, error)
}

func NewService(db database.DB, authoSvc authorization.Service, partitionSvc PartitionService) Service {
	// Create dao
	dao := daos.NewDao(db)

	return &service{dao: dao, validator: validator.New(), partitionSvc: partitionSvc, authorizationSvc: authoSvc}
}
