package pagination

import (
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/common"
	"gorm.io/gorm"
)

// PageInput represents an input pagination configuration.
type PageInput struct {
	Skip  int
	Limit int
}

// PageOutput represents an output pagination structure.
type PageOutput struct {
	TotalRecord int
	Limit       int
	Skip        int
	HasPrevious bool
	HasNext     bool
}

// PagingOptions represents pagination options.
type PagingOptions struct {
	// Gorm database
	DB *gorm.DB
	// Pagination input
	PageInput *PageInput
	// Must be a pointer to an object with *SortOrderEnum objects with tags
	Sort interface{}
	// Must be a pointer to an object with *GenericFilter objects or implementing the GenericFilterBuilder interface and with tags
	Filter interface{}
	// Must be a pointer to an object with booleans with tags
	Projection interface{}
	// This function is called after filters and before any sorts
	ExtraFunc func(db *gorm.DB) (*gorm.DB, error)
}

// Paging function in order to have a paginated sorted and filters list of objects.
// Parameters:
// - result: Must be a pointer to a list of objects
// - options: Pagination options
// .
func Paging(
	result interface{},
	options *PagingOptions,
) (*PageOutput, error) {
	// Manage default limit
	if options.PageInput.Limit == 0 {
		options.PageInput.Limit = 10
	}

	var count int64 = 0

	// Create transaction to avoid situations where count and find are different
	err := options.DB.Transaction(func(db *gorm.DB) error {
		// Apply filter
		db, err := common.ManageFilter(options.Filter, db)
		// Check error
		if err != nil {
			return err
		}

		// Extra function
		if options.ExtraFunc != nil {
			db, err = options.ExtraFunc(db)
			// Check error
			if err != nil {
				return err
			}
		}

		// Count all objects
		db = db.Model(result).Count(&count)
		// Check error
		if db.Error != nil {
			return db.Error
		}

		// Apply sort
		db, err = common.ManageSortOrder(options.Sort, db)
		// Check error
		if err != nil {
			return err
		}

		// Apply projection
		db, err = common.ManageProjection(options.Projection, db)
		// Check error
		if err != nil {
			return err
		}

		// Request to database with limit and offset
		db = db.Limit(options.PageInput.Limit).Offset(options.PageInput.Skip).Find(result)
		// Check error
		if db.Error != nil {
			return db.Error
		}

		return nil
	})

	// Check error
	if err != nil {
		return nil, err
	}

	return getPageOutput(options.PageInput, count), nil
}

func getPageOutput(p *PageInput, count int64) *PageOutput {
	var paginator PageOutput
	// Create total record
	paginator.TotalRecord = int(count)
	// Store skip
	paginator.Skip = p.Skip
	// Store limit
	paginator.Limit = p.Limit
	// Calculate has next page
	paginator.HasNext = (p.Limit+p.Skip < paginator.TotalRecord)
	// Calculate has previous page
	paginator.HasPrevious = (p.Skip != 0)

	return &paginator
}
