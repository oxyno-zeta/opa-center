package common

import (
	"fmt"
	"reflect"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"gorm.io/gorm"
)

// Supported enum type for testing purpose.
var supportedEnumType = reflect.TypeOf(new(SortOrderEnum))

func ManageSortOrder(sort interface{}, db *gorm.DB) (*gorm.DB, error) {
	// Create result
	res := db
	// Get reflect value of sort object
	rVal := reflect.ValueOf(sort)
	// Get kind of sort
	rKind := rVal.Kind()
	// Check nil
	if rKind == reflect.Invalid || (rKind == reflect.Ptr && rVal.IsNil()) {
		// Stop here
		return manageDefaultSort(res), nil
	}
	// Check if kind is supported
	if rKind != reflect.Struct && rKind != reflect.Ptr {
		return nil, errors.NewInvalidInputError("sort must be an object")
	}

	// Indirect value
	indirect := reflect.Indirect(rVal)
	// Get indirect data
	indData := indirect.Interface()
	// Get type of indirect value
	typeOfIndi := reflect.TypeOf(indData)
	// Variable to know at the end if one sort was applied
	sortApplied := false

	// Loop over all num fields
	for i := 0; i < indirect.NumField(); i++ {
		// Get field type
		fType := typeOfIndi.Field(i)
		// Get tag on field
		tagVal := fType.Tag.Get(dbColTagName)
		// Check that field have a tag set and correct
		if tagVal == "" || tagVal == "-" {
			// Skip this value
			continue
		}
		// Check that type is supported
		if fType.Type != supportedEnumType {
			return nil, errors.NewInvalidInputError(fmt.Sprintf("field %s with sort tag must be a *SortOrderEnum", fType.Name))
		}
		// Get field value
		fVal := indirect.Field(i)
		// Test if field is nil
		if fVal.IsNil() {
			// Skip field because of nil
			continue
		}
		// Get value from field
		val := fVal.Interface()
		// Cast value to Sort Order Enum
		enu := val.(*SortOrderEnum)
		// Apply order
		res = res.Order(fmt.Sprintf("%s %s", tagVal, enu.String()))
		// Store sort applied
		sortApplied = true
	}

	// Check if one sort was applied or not in order to put the default one
	if !sortApplied {
		res = manageDefaultSort(res)
	}

	return res, nil
}

func manageDefaultSort(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}
