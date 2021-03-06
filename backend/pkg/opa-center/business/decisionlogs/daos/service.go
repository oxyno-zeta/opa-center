package daos

import (
	"errors"

	daosmodels "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/daos/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
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
	err := gdb.AutoMigrate(&daosmodels.DecisionLog{})

	return err
}

func (s *service) Delete(filter *models.Filter) error {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Apply filter
	db, err := common.ManageFilter(filter, gdb)
	// Check error
	if err != nil {
		return err
	}

	return db.Unscoped().Delete(&models.DecisionLog{}).Error
}

func (s *service) FindByID(id string, projection *models.Projection) (*models.DecisionLog, error) {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Create result
	var res daosmodels.DecisionLog

	// Manage projection
	gdb, err := common.ManageProjection(projection, gdb)
	// Check error
	if err != nil {
		return nil, err
	}

	// Find in db
	dbres := gdb.Where("id = ?", id).First(&res)

	// check error
	if dbres.Error != nil {
		// Check if error is a not found error
		if errors.Is(dbres.Error, gorm.ErrRecordNotFound) {
			// Return nil as answer
			return nil, nil
		}
		// Another error
		return nil, dbres.Error
	}

	// Map result
	mres, err := fromDao(&res)
	// Check error
	if err != nil {
		return nil, err
	}

	return mres, nil
}

func (s *service) FindOneByDecisionID(did string, projection *models.Projection) (*models.DecisionLog, error) {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Create result
	var res daosmodels.DecisionLog

	// Manage projection
	gdb, err := common.ManageProjection(projection, gdb)
	// Check error
	if err != nil {
		return nil, err
	}

	// Find in db
	dbres := gdb.Where("decision_id = ?", did).First(&res)

	// check error
	if dbres.Error != nil {
		// Check if error is a not found error
		if errors.Is(dbres.Error, gorm.ErrRecordNotFound) {
			// Return nil as answer
			return nil, nil
		}
		// Another error
		return nil, dbres.Error
	}

	// Map result
	mres, err := fromDao(&res)
	// Check error
	if err != nil {
		return nil, err
	}

	return mres, nil
}

func (s *service) Save(ins *models.DecisionLog) error {
	// Get gorm database
	gdb := s.db.GetGormDB()
	// Transform object
	input := toDao(ins)
	// Save
	res := gdb.Save(input)

	return res.Error
}

func (s *service) GetAllPaginated(
	page *pagination.PageInput,
	sort *models.SortOrder,
	filter *models.Filter,
	projection *models.Projection,
) ([]*models.DecisionLog, *pagination.PageOutput, error) {
	// Get gorm db
	db := s.db.GetGormDB()
	// result
	dres := make([]*daosmodels.DecisionLog, 0)
	// Find todos
	pageOut, err := pagination.Paging(&dres, &pagination.PagingOptions{
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

	// Result
	res := make([]*models.DecisionLog, 0)
	// Loop over list
	for i := 0; i < len(dres); i++ {
		it := dres[i]
		// Map
		r, err := fromDao(it)
		// Check error
		if err != nil {
			return nil, nil, err
		}
		// Append
		res = append(res, r)
	}

	return res, pageOut, nil
}
