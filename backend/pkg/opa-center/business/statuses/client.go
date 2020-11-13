package statuses

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/authx/authorization"
	pmodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/daos"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

type Service interface {
	// Migrate database
	MigrateDB(systemLogger log.Logger) error
	// Create decision log used internally only
	UnsecureCreate(partitionName string, inp map[string]interface{}) error
	// Get data paginated
	GetAllPaginated(
		ctx context.Context,
		partitionID string,
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.Status, *pagination.PageOutput, error)
	// Find by id
	FindByID(ctx context.Context, id string, projection *models.Projection) (*models.Status, error)
}

type PartitionService interface {
	UnsecureFindByID(id string) (*pmodels.Partition, error)
}

func NewService(db database.DB, authoSvc authorization.Service, partitionSvc PartitionService) Service {
	// Create dao
	dao := daos.NewDao(db)

	return &service{dao: dao, validator: validator.New(), partitionSvc: partitionSvc, authorizationSvc: authoSvc}
}
