package daos

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

// Dao represent a decision logs access object service.
type Dao interface {
	// MigrateDB will migrate database
	MigrateDB() error
	// Save will save object in database
	Save(ins *models.DecisionLog) error
	// FindOneByDecisionID will find one decision log by decision id
	FindOneByDecisionID(did string, projection *models.Projection) (*models.DecisionLog, error)
	// FindByID will find by id
	FindByID(di string, projection *models.Projection) (*models.DecisionLog, error)
	// Get data paginated
	GetAllPaginated(
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.DecisionLog, *pagination.PageOutput, error)
	// Delete permanently with filter
	Delete(filter *models.Filter) error
}

func NewDao(db database.DB) Dao {
	return &service{
		db: db,
	}
}
