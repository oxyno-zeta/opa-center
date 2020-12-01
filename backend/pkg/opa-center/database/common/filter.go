package common

import (
	"fmt"
	"reflect"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"gorm.io/gorm"
)

// AND field.
const andFieldName = "AND"

// OR field.
const orFieldName = "OR"

func ManageFilter(filter interface{}, db *gorm.DB) (*gorm.DB, error) {
	return manageFilter(filter, db, db, false)
}

func manageFilter(filter interface{}, db, originalDB *gorm.DB, skipInputNotObject bool) (*gorm.DB, error) { //nolint:unparam // Because seems to be bugged...
	// Create result
	res := db
	// Get reflect value of filter object
	rVal := reflect.ValueOf(filter)
	// Get kind of filter
	rKind := rVal.Kind()
	// Check if filter isn't nil
	if rKind == reflect.Invalid || (rKind == reflect.Ptr && rVal.IsNil()) {
		// Stop here
		return res, nil
	}
	// Check if kind is supported
	if rKind != reflect.Struct && rKind != reflect.Ptr {
		// Check if skip input not an object is enabled
		// This is used in recursive calls in order to avoid errors when OR or AND cases aren't an object supported
		if skipInputNotObject {
			return db, nil
		}

		// No skip => Error
		return nil, errors.NewInvalidInputError("filter must be an object")
	}

	// Indirect value
	indirect := reflect.Indirect(rVal)
	// Get indirect data
	indData := indirect.Interface()
	// Get type of indirect value
	typeOfIndi := reflect.TypeOf(indData)

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
		// Check if value is a pointer or not
		if fVal.Kind() != reflect.Ptr {
			return nil, errors.NewInvalidInputError(
				fmt.Sprintf("field %s with filter tag must be a *GenericFilter or implement GenericFilterBuilder interface", fType.Name),
			)
		}
		// Test if field is nil
		if fVal.IsNil() {
			// Skip field because of nil
			continue
		}
		// Get value from field
		val := fVal.Interface()

		// Try to cast it as GenericFilterBuilder
		v1, castGFB := val.(GenericFilterBuilder)
		// Check that type is supported
		if !castGFB {
			return nil, errors.NewInvalidInputError(
				fmt.Sprintf("field %s with filter tag must be a *GenericFilter or implement GenericFilterBuilder interface", fType.Name),
			)
		}

		// Cast value
		v, err := v1.GetGenericFilter()
		// Check error
		if err != nil {
			return nil, err
		}
		// Manage filter request
		res2, err := manageFilterRequest(tagVal, v, res)
		// Check error
		if err != nil {
			return nil, err
		}

		res = res2
	}

	// Manage AND cases
	// Check in type that AND key exists
	_, exists := typeOfIndi.FieldByName(andFieldName)
	// Check if it exists
	if exists {
		// AND field is detected
		// Get field AND
		andRVal := indirect.FieldByName(andFieldName)
		// Check that type is a slice
		if andRVal.Kind() == reflect.Slice {
			// Loop over items in array
			for i := 0; i < andRVal.Len(); i++ {
				// Get element at index
				andElementRVal := andRVal.Index(i)
				// Get value behind
				andElement := andElementRVal.Interface()
				// Call manage filter
				res2, err := manageFilter(andElement, originalDB, originalDB, true)
				// Check error
				if err != nil {
					return nil, err
				}
				// Save result
				res = res.Where(res2)
			}
		}
	}

	// Manage OR cases
	// Check in type that OR key exists
	_, exists = typeOfIndi.FieldByName(orFieldName)
	// Check if it exists
	if exists {
		// OR field is detected
		// Get field OR
		orRVal := indirect.FieldByName(orFieldName)

		// Check that type is a slice
		if orRVal.Kind() == reflect.Slice {
			// Get array length
			lgt := orRVal.Len()
			// Check length in order to ignore it it is 0
			if lgt != 0 {
				// Array is populated
				// Loop over elements
				for i := 0; i < lgt; i++ {
					// Get element
					elemRVal := orRVal.Index(i)
					// Get data behind
					elem := elemRVal.Interface()
					// Call manage filter WITH the original db in order to create a pure subquery
					// See here: https://gorm.io/docs/advanced_query.html#Group-Conditions
					res2, err := manageFilter(elem, originalDB, originalDB, true)
					// Check error
					if err != nil {
						return nil, err
					}
					// Manage result
					// Check if it is the first element
					if i == 0 {
						// First element must be managed as a where
						res = res.Where(res2)
					} else {
						// This is the other cases
						res = res.Or(res2)
					}
				}
			}
		}
	}

	// Return
	return res, nil
}

func manageFilterRequest(dbCol string, v *GenericFilter, db *gorm.DB) (*gorm.DB, error) {
	// Create result
	dbRes := db
	// Check Equal case
	if v.Eq != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s = ?", dbCol), v.Eq)
	}
	// Check not equal case
	if v.NotEq != nil {
		dbRes = dbRes.Not(fmt.Sprintf("%s = ?", dbCol), v.NotEq)
	}
	// Check greater and equal than case
	if v.Gte != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s >= ?", dbCol), v.Gte)
	}
	// Check not greater and equal than case
	if v.NotGte != nil {
		dbRes = dbRes.Not(fmt.Sprintf("%s >= ?", dbCol), v.NotGte)
	}
	// Check greater than case
	if v.Gt != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s > ?", dbCol), v.Gt)
	}
	// Check not greater than case
	if v.NotGt != nil {
		dbRes = dbRes.Not(fmt.Sprintf("%s > ?", dbCol), v.NotGt)
	}
	// Check less and equal than case
	if v.Lte != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s <= ?", dbCol), v.Lte)
	}
	// Check not less and equal than case
	if v.NotLte != nil {
		dbRes = dbRes.Not(fmt.Sprintf("%s <= ?", dbCol), v.NotLte)
	}
	// Check less than case
	if v.Lt != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s < ?", dbCol), v.Lt)
	}
	// Check not less than case
	if v.NotLt != nil {
		dbRes = dbRes.Not(fmt.Sprintf("%s < ?", dbCol), v.NotLt)
	}
	// Check contains case
	if v.Contains != nil {
		// Get string value
		s, err := getStringValue(v.Contains)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("contains " + err.Error())
		}

		dbRes = dbRes.Where(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%%%s%%", s))
	}
	// Check not contains case
	if v.NotContains != nil {
		// Get string value
		s, err := getStringValue(v.NotContains)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("notContains " + err.Error())
		}

		dbRes = dbRes.Not(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%%%s%%", s))
	}
	// Check starts with case
	if v.StartsWith != nil {
		// Get string value
		s, err := getStringValue(v.StartsWith)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("startsWith " + err.Error())
		}

		dbRes = dbRes.Where(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%s%%", s))
	}
	// Check not starts with case
	if v.NotStartsWith != nil {
		// Get string value
		s, err := getStringValue(v.NotStartsWith)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("notStartsWith " + err.Error())
		}

		dbRes = dbRes.Not(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%s%%", s))
	}
	// Check ends with case
	if v.EndsWith != nil {
		// Get string value
		s, err := getStringValue(v.EndsWith)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("endsWith " + err.Error())
		}

		dbRes = dbRes.Where(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%%%s", s))
	}
	// Check not ends with case
	if v.NotEndsWith != nil {
		// Get string value
		s, err := getStringValue(v.NotEndsWith)
		// Check error
		if err != nil {
			return nil, errors.NewInvalidInputError("notEndsWith " + err.Error())
		}

		dbRes = dbRes.Not(fmt.Sprintf("%s LIKE ?", dbCol), fmt.Sprintf("%%%s", s))
	}
	// Check in case
	if v.In != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s IN (?)", dbCol), v.In)
	}
	// Check not in case
	if v.NotIn != nil {
		dbRes = dbRes.Where(fmt.Sprintf("%s NOT IN (?)", dbCol), v.NotIn)
	}
	// Check is null case
	if v.IsNull {
		dbRes = dbRes.Where(fmt.Sprintf("%s IS NULL", dbCol))
	}
	// Check is not null case
	if v.IsNotNull {
		dbRes = dbRes.Where(fmt.Sprintf("%s IS NOT NULL", dbCol))
	}

	// Return
	return dbRes, nil
}

func getStringValue(x interface{}) (string, error) {
	// Get reflect value
	val := reflect.ValueOf(x)
	// Check if val is a pointer
	if val.Kind() == reflect.Ptr {
		// Indirect for pointers
		val = reflect.Indirect(val)
	}
	// Check if type is acceptable
	if val.Kind() != reflect.String {
		return "", errors.NewInvalidInputError("value must be a string or *string")
	}

	return val.String(), nil
}
