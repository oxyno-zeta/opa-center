package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/thoas/go-funk"
)

const graphqlFieldTagKey = "graphqlfield"

func ManageConnectionNodeProjection(
	ctx context.Context,
	projectionOut interface{},
) error {
	// Validate projection out
	err := validateProjectionOut(projectionOut)
	// Check error
	if err != nil {
		return err
	}

	// Get operation context
	octx := graphql.GetOperationContext(ctx)
	// Get graphql fields
	fields := graphql.CollectFieldsCtx(ctx, nil)

	// Find edges in fields
	fEdgesInt := funk.Find(fields, func(f graphql.CollectedField) bool {
		return f.Name == "edges"
	})
	// Check if edges field exists in graphql fields
	if fEdgesInt == nil {
		// Not found in graphql projection => stop
		return nil
	}
	// Cast field
	fEdges := fEdgesInt.(graphql.CollectedField)

	// Get graphql fields projection in edges
	inFEdges := graphql.CollectFields(octx, fEdges.Selections, nil)

	// Find node in edges graphql fields
	fNodeInt := funk.Find(inFEdges, func(f graphql.CollectedField) bool {
		return f.Name == "node"
	})
	// Check if node field exists in graphql fields
	if fNodeInt == nil {
		// Not found in graphql projection => stop
		return nil
	}
	// Cast field
	fNode := fNodeInt.(graphql.CollectedField)

	// Start projection on this path
	err = manageGraphqlProjection(
		graphql.CollectFields(octx, fNode.Selections, nil),
		projectionOut,
	)
	// Check error
	if err != nil {
		return err
	}

	// Default
	return nil
}

func ManageSimpleProjection(
	ctx context.Context,
	projectionOut interface{},
) error {
	// Validate projection out
	err := validateProjectionOut(projectionOut)
	// Check error
	if err != nil {
		return err
	}

	// Manage graphql projection
	err = manageGraphqlProjection(
		graphql.CollectFieldsCtx(ctx, nil),
		projectionOut,
	)
	// Check error
	if err != nil {
		return err
	}

	// Default
	return nil
}

func validateProjectionOut(projectionOut interface{}) error {
	// Check if input is nil
	if projectionOut == nil {
		return errors.NewInvalidInputError("projection output cannot be nil")
	}

	// Get projection type
	projOutType := reflect.TypeOf(projectionOut)
	// Check if projection out is a pointer
	if projOutType.Kind() != reflect.Ptr {
		return errors.NewInvalidInputError("projection output must be a pointer to an object")
	}
	// Get projection out value
	projOutVal := reflect.ValueOf(projectionOut)
	// Indirect value
	indVal := reflect.Indirect(projOutVal)
	// Check indirect value type
	if indVal.Kind() != reflect.Struct {
		return errors.NewInvalidInputError("projection output must be a pointer to an object")
	}

	// Default
	return nil
}

func manageGraphqlProjection(
	gfields []graphql.CollectedField,
	projectionOut interface{},
) error {
	// Get reflect ptr value
	pOutPtrVal := reflect.ValueOf(projectionOut)
	// Get reflect value
	pOutVal := reflect.Indirect(pOutPtrVal)
	// Get reflect ptr type
	pOutPtrType := reflect.TypeOf(projectionOut)
	// Get reflect type
	pOutType := pOutPtrType.Elem()

	// Loop over projection struct fields
	for i := 0; i < pOutType.NumField(); i++ {
		// Get field
		fieldType := pOutType.Field(i)
		// Get tag value for graphql field
		tagValue := fieldType.Tag.Get(graphqlFieldTagKey)
		// Check if tag exists or ignored
		if tagValue == "" || tagValue == "-" {
			// Continue to next field
			continue
		}

		// Check if field is asked in graphql
		gfieldInt := funk.Find(gfields, func(gfield graphql.CollectedField) bool {
			return gfield.Name == tagValue
		})
		// Check if field isn't found
		if gfieldInt == nil {
			// Field isn't found => continue to next field
			continue
		}

		// Check if field is a boolean
		if fieldType.Type.Kind() == reflect.Bool {
			pOutVal.Field(i).SetBool(true)
			// Stop here
			continue
		}

		// Field is found but type isn't supported
		return errors.NewInvalidInputError(fmt.Sprintf("field %s must be a boolean", fieldType.Name))
	}

	// Default
	return nil
}
