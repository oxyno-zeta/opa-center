package common

import (
	"fmt"
	"io"
	"strconv"
)

// SortOrderEnum is used to sort data.
// This must be used as a pointer in other structures to be used automatically in sort.
// Moreover, a tag containing the database field must be declared.
// Example:
// type Sort struct {
// 	Field1 *SortOrderEnum `dbfield:"field_1"`
// }
// .
type SortOrderEnum string

var (
	SortOrderEnumAsc  SortOrderEnum = "ASC"
	SortOrderEnumDesc SortOrderEnum = "DESC"
)

var AllSortOrderEnum = []SortOrderEnum{
	SortOrderEnumAsc,
	SortOrderEnumDesc,
}

func (e SortOrderEnum) IsValid() bool {
	switch e {
	case SortOrderEnumAsc, SortOrderEnumDesc:
		return true
	}

	return false
}

func (e SortOrderEnum) String() string {
	return string(e)
}

func (e *SortOrderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrderEnum", str)
	}

	return nil
}

func (e SortOrderEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
