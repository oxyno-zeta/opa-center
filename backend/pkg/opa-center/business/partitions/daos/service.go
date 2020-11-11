package daos

import (
	"errors"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
	"gorm.io/gorm"
)

type service struct {
	db database.DB
}

func (s *service) MigrateDB() error {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Migrate
	err := gdb.AutoMigrate(&models.Partition{})

	return err
}

func (s *service) Save(ins *models.Partition) (*models.Partition, error) {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Save
	res := gdb.Save(ins)
	// Check error
	if res.Error != nil {
		return nil, res.Error
	}
	// Return result
	return ins, nil
}

func (s *service) GetAllPaginated(
	page *pagination.PageInput,
	sort *models.SortOrder,
	filter *models.Filter,
	projection *models.Projection,
) ([]*models.Partition, *pagination.PageOutput, error) {
	// Get gorm db
	db := s.db.GetGormDB()
	// result
	res := make([]*models.Partition, 0)
	// Find todos
	pageOut, err := pagination.Paging(&res, &pagination.PagingOptions{
		DB:         db,
		Filter:     filter,
		PageInput:  page,
		Projection: projection,
		Sort:       sort,
	})
	// Check error
	if err != nil {
		return nil, nil, err
	}

	return res, pageOut, nil
}

func (s *service) FindByName(name string, projection *models.Projection) (*models.Partition, error) {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Apply projection
	gdb, err := common.ManageProjection(projection, gdb)
	// Check error
	if err != nil {
		return nil, err
	}
	// Create result
	var res models.Partition
	// Request database
	dbres := gdb.Where("name = ?", name).First(&res)
	// Check error
	if dbres.Error != nil {
		// Check if error is a not found error
		if errors.Is(dbres.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// Error
		return nil, dbres.Error
	}
	// Return result
	return &res, nil
}

func (s *service) FindByID(id string, projection *models.Projection) (*models.Partition, error) {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Apply projection
	gdb, err := common.ManageProjection(projection, gdb)
	// Check error
	if err != nil {
		return nil, err
	}
	// Create result
	var res models.Partition
	// Request database
	dbres := gdb.Where("id = ?", id).First(&res)
	// Check error
	if dbres.Error != nil {
		// Check if error is a not found error
		if errors.Is(dbres.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// Error
		return nil, dbres.Error
	}
	// Return result
	return &res, nil
}
