package common

import (
	"time"
)

// GenericFilter is a structure that will handle filters other than Date.
// This must be used as a pointer in other structures to be used automatically in filters.
// Moreover, a tag containing the database field must be declared.
// Example:
// type Filter struct {
// 	Field1 *GenericFilter `dbfield:"field_1"`
// }
// .
type GenericFilter struct {
	// Allow to test equality to
	Eq interface{}
	// Allow to test non equality to
	NotEq interface{}
	// Allow to test greater or equal than
	Gte interface{}
	// Allow to test not greater or equal than
	NotGte interface{}
	// Allow to test greater than
	Gt interface{}
	// Allow to test not greater than
	NotGt interface{}
	// Allow to test less or equal than
	Lte interface{}
	// Allow to test not less or equal than
	NotLte interface{}
	// Allow to test less than
	Lt interface{}
	// Allow to test not less than
	NotLt interface{}
	// Allow to test if a string contains another string.
	// Contains must be a string
	Contains interface{}
	// Allow to test if a string isn't containing another string.
	// NotContains must be a string
	NotContains interface{}
	// Allow to test if a string starts with another string.
	// StartsWith with must be a string
	StartsWith interface{}
	// Allow to test if a string isn't starting with another string.
	// NotStartsWith must be a string
	NotStartsWith interface{}
	// Allow to test if a string ends with another string.
	// EndsWith with must be a string
	EndsWith interface{}
	// Allow to test if a string isn't ending with another string.
	// NotEndsWith must be a string
	NotEndsWith interface{}
	// Allow to test if value is in array
	In interface{}
	// Allow to test if value isn't in array
	NotIn interface{}
}

// DateFilter is a structure that will handle filters for dates.
// This must be used as a pointer in other structures to be used automatically in filters.
// Moreover, a tag containing the database field must be declared.
// Example:
// type Filter struct {
// 	Field1 *DateFilter `dbfield:"field_1"`
// }
// .
type DateFilter struct {
	// Allow to test equality to
	Eq *string
	// Allow to test non equality to
	NotEq *string
	// Allow to test greater or equal than
	Gte *string
	// Allow to test not greater or equal than
	NotGte *string
	// Allow to test greater than
	Gt *string
	// Allow to test not greater than
	NotGt *string
	// Allow to test less or equal than
	Lte *string
	// Allow to test not less or equal than
	NotLte *string
	// Allow to test less than
	Lt *string
	// Allow to test not less than
	NotLt *string
	// Allow to test if value is in array
	In []string
	// Allow to test if value isn't in array
	NotIn []string
}

// GenericFilterBuilder is an interface that must be implemented in order to work automatic filter.
// This is done like this in order to add more fields in GenericFilter without the need of upgrading
// all code in other to be compatible.
type GenericFilterBuilder interface {
	GetGenericFilter() (*GenericFilter, error)
}

func (g *GenericFilter) GetGenericFilter() (*GenericFilter, error) { return g, nil }

func (d *DateFilter) GetGenericFilter() (*GenericFilter, error) {
	// Create result
	res := &GenericFilter{}

	// Eq case
	if d.Eq != nil {
		// Parse time
		t, err := parseTime(*d.Eq)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Eq = t
	}

	// Not Eq case
	if d.NotEq != nil {
		// Parse time
		t, err := parseTime(*d.NotEq)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotEq = t
	}

	// Gte case
	if d.Gte != nil {
		// Parse time
		t, err := parseTime(*d.Gte)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Gte = t
	}

	// Not Gte case
	if d.NotGte != nil {
		// Parse time
		t, err := parseTime(*d.NotGte)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotGte = t
	}

	// Gt case
	if d.Gt != nil {
		// Parse time
		t, err := parseTime(*d.Gt)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Gt = t
	}

	// Not Gt case
	if d.NotGt != nil {
		// Parse time
		t, err := parseTime(*d.NotGt)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotGt = t
	}

	// Lte case
	if d.Lte != nil {
		// Parse time
		t, err := parseTime(*d.Lte)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Lte = t
	}

	// Not Lte case
	if d.NotLte != nil {
		// Parse time
		t, err := parseTime(*d.NotLte)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotLte = t
	}

	// Lt case
	if d.Lt != nil {
		// Parse time
		t, err := parseTime(*d.Lt)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Lt = t
	}

	// Not Lt case
	if d.NotLt != nil {
		// Parse time
		t, err := parseTime(*d.NotLt)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotLt = t
	}

	// In case
	if d.In != nil {
		// Parse time
		t, err := parseTimes(d.In)
		// Check error
		if err != nil {
			return nil, err
		}

		res.In = t
	}

	// Not In case
	if d.NotIn != nil {
		// Parse time
		t, err := parseTimes(d.NotIn)
		// Check error
		if err != nil {
			return nil, err
		}

		res.NotIn = t
	}

	return res, nil
}

func parseTime(x string) (*time.Time, error) {
	// Parse date
	t, err := time.Parse(time.RFC3339, x)
	// Check error
	if err != nil {
		return nil, err
	}

	// Force utc
	t = t.UTC()

	return &t, nil
}

func parseTimes(x []string) ([]*time.Time, error) {
	// Prepare result
	res := make([]*time.Time, 0)

	// Loop over all values
	for _, v := range x {
		// Parse time
		t, err := parseTime(v)
		// Check error
		if err != nil {
			return nil, err
		}
		// Append
		res = append(res, t)
	}

	return res, nil
}
