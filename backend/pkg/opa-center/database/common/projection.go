package common

import (
	"fmt"
	"reflect"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"gorm.io/gorm"
)

func ManageProjection(projection interface{}, db *gorm.DB) (*gorm.DB, error) {
	// Create result
	res := db
	// Get reflect value of projection object
	rVal := reflect.ValueOf(projection)
	// Get kind of projection
	rKind := rVal.Kind()
	// Check if projection isn't nil
	if rKind == reflect.Invalid || (rKind == reflect.Ptr && rVal.IsNil()) {
		// Stop here
		return res, nil
	}
	// Check if kind is supported
	if rKind != reflect.Struct && rKind != reflect.Ptr {
		// No skip => Error
		return nil, errors.NewInvalidInputError("projection must be an object")
	}

	// Indirect value
	indirect := reflect.Indirect(rVal)
	// Get indirect data
	indData := indirect.Interface()
	// Get type of indirect value
	typeOfIndi := reflect.TypeOf(indData)

	// Create array of projection results
	selectArray := make([]string, 0)
	// Store boolean to avoid checking length of select array which consum
	selectArrayFilled := false

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
		// Get field value
		fVal := indirect.Field(i)
		// Check if value is a boolean or not
		if fVal.Kind() != reflect.Bool {
			return nil, errors.NewInvalidInputError(
				fmt.Sprintf("field %s with projection tag must be a boolean", fType.Name),
			)
		}
		// Get value from field
		val := fVal.Interface()
		// Cast it to boolean
		v := val.(bool)

		// Manage projection if enabled
		if v {
			selectArray = append(selectArray, tagVal)
			selectArrayFilled = true
		}
	}

	// Check if projection array is filled or not
	if selectArrayFilled {
		res = res.Select(selectArray)
	}

	// Default case
	return res, nil
}
