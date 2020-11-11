package utils

import (
	"reflect"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

const edgesFieldName = "Edges"
const nodeFieldName = "Node"
const pageInfoFieldName = "PageInfo"
const cursorFieldName = "Cursor"

// Store supported type of page info.
var pageInfoSupportedType = reflect.TypeOf(new(PageInfo))
var cursorSupportedType = reflect.TypeOf("")

func MapConnection(connectionResult interface{}, list interface{}, pageOut *pagination.PageOutput) error {
	// Validate that connection result isn't nil
	if connectionResult == nil {
		return errors.NewInvalidInputError("connection result argument mustn't be nil")
	}
	// Validate that input isn't nil
	if list == nil {
		return errors.NewInvalidInputError("list argument mustn't be nil")
	}

	// Create page info cursors
	startCursor := ""
	endCursor := ""

	// Reflect input
	listVal := reflect.ValueOf(list)
	// Check that list is a slice
	if listVal.Kind() != reflect.Slice {
		return errors.NewInvalidInputError("list argument must be a slice")
	}
	// Get list length
	listLen := listVal.Len()
	// Get last index of list
	lastIndex := listLen - 1

	// Create reflect object from input value
	connectionResultVal := reflect.ValueOf(connectionResult)
	// Check that top is a pointer
	if connectionResultVal.Kind() != reflect.Ptr {
		return errors.NewInvalidInputError("connection result argument must be a pointer to a connection object")
	}

	// Get connection pointer type
	resType := reflect.TypeOf(connectionResult)
	// Get real type of connection and not a pointer to as type
	connectionType := resType.Elem()
	// Get edges type
	edgesType, exists := connectionType.FieldByName(edgesFieldName)
	// Check if field Edges exists
	if !exists {
		return errors.NewInvalidInputError("field Edges not found in connection object")
	}
	// Check that edges is a slice
	if edgesType.Type.Kind() != reflect.Slice {
		return errors.NewInvalidInputError("field Edges must be a slice")
	}
	// Get page info type
	pageInfoType, exists := connectionType.FieldByName(pageInfoFieldName)
	// Check if field PageInfo exists
	if !exists {
		return errors.NewInvalidInputError("field PageInfo not found in connection object")
	}
	// Check that page info from struct is equal to page info
	if pageInfoType.Type != pageInfoSupportedType {
		return errors.NewInvalidInputError("field PageInfo isn't with the type *PageInfo")
	}

	// Get edge pointer type
	edgeTypePtr := edgesType.Type.Elem()
	// Check that Edge type is a pointer in slice
	if edgeTypePtr.Kind() != reflect.Ptr {
		return errors.NewInvalidInputError("field Edges must be a slice of struct pointers ([]*Edge)")
	}
	// Get edge type
	edgeTypeStruct := edgeTypePtr.Elem()
	// Check that Edge is really a struct after pointer
	if edgeTypeStruct.Kind() != reflect.Struct {
		return errors.NewInvalidInputError("field Edges must be a slice of struct pointers ([]*Edge)")
	}

	// Check that Edge have a Cursor key
	cursorType, exists := edgeTypeStruct.FieldByName(cursorFieldName)
	// Check if exists
	if !exists {
		return errors.NewInvalidInputError("field Cursor not found in Edge object")
	}
	// Check that cursor type is string
	if cursorType.Type != cursorSupportedType {
		return errors.NewInvalidInputError("field Cursor from Edge object must be a string")
	}

	// Get list type
	listType := reflect.TypeOf(list)
	// Get item type from list
	itemListType := listType.Elem()

	// Check that Edge have a Node key
	nodeType, exists := edgeTypeStruct.FieldByName(nodeFieldName)
	// Check if exists
	if !exists {
		return errors.NewInvalidInputError("field Node not found in Edge object")
	}
	// Check that node value is same type as item in list
	if !nodeType.Type.ConvertibleTo(itemListType) {
		return errors.NewInvalidInputError("field Node must have the same type of items in the list argument")
	}

	// Indirect top value
	topValInd := reflect.Indirect(connectionResultVal)
	// Get edges values
	edgesVal := topValInd.FieldByName(edgesFieldName)

	// Loop over all items in list
	for i := 0; i < listLen; i++ {
		// Create cursor for element
		cursor := GetPaginateCursor(i, pageOut.Skip)

		// Store start cursor if it is the first element
		if i == 0 {
			startCursor = cursor
		}
		// Store end cursor if it is the last element
		if i == lastIndex {
			endCursor = cursor
		}

		// Create new edge object
		edgePtr := reflect.New(edgeTypeStruct)
		// Indirect edge object to get value of pointer
		edgeStruct := reflect.Indirect(edgePtr)

		// Get element in list
		listEl := listVal.Index(i)

		// Set node value
		edgeStruct.FieldByName(nodeFieldName).Set(listEl)
		// Set cursor value
		edgeStruct.FieldByName(cursorFieldName).Set(reflect.ValueOf(cursor))
		// Append in the edges list of value
		edgesVal.Set(reflect.Append(edgesVal, edgePtr))
	}

	// Create page info object
	pageInfo := GetPageInfo(startCursor, endCursor, pageOut)
	// Add page info structure in connection
	topValInd.FieldByName(pageInfoFieldName).Set(reflect.ValueOf(pageInfo))

	return nil
}
