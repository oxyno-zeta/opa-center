package daos

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

// Dao represent a partition object service.
type Dao interface {
	// MigrateDB will migrate database
	MigrateDB() error
	// Get data paginated
	GetAllPaginated(
		page *pagination.PageInput,
		sort *models.SortOrder,
		filter *models.Filter,
		projection *models.Projection,
	) ([]*models.Partition, *pagination.PageOutput, error)
	// Save will save partition object
	Save(ins *models.Partition) (*models.Partition, error)
	// Find by name
	FindByName(name string, projection *models.Projection) (*models.Partition, error)
	// Find by id
	FindByID(id string, projection *models.Projection) (*models.Partition, error)
}

func NewDao(db database.DB) Dao {
	return &service{
		db: db,
	}
}
