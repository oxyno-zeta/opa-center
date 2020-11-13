package daos

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

// Dao represent a decision logs access object service.
type Dao interface {
	// MigrateDB will migrate database
	MigrateDB() error
	// Save will save object in database
	Save(ins *models.Status) error
	// FindByID will find by id
	FindByID(di string, projection *models.Projection) (*models.Status, error)
	// Get data paginated
	GetAllPaginated(
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.Status, *pagination.PageOutput, error)
}

func NewDao(db database.DB) Dao {
	return &service{
		db: db,
	}
}
