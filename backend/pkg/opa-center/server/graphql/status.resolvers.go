package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	models1 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	models2 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/generated"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/mappers"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/utils"
)

func (r *statusResolver) ID(ctx context.Context, obj *models2.Status) (string, error) {
	return utils.ToIDRelay(mappers.StatusIDPrefix, obj.ID), nil
}

func (r *statusResolver) CreatedAt(ctx context.Context, obj *models2.Status) (string, error) {
	return utils.FormatTime(obj.CreatedAt), nil
}

func (r *statusResolver) UpdatedAt(ctx context.Context, obj *models2.Status) (string, error) {
	return utils.FormatTime(obj.UpdatedAt), nil
}

func (r *statusResolver) Partition(ctx context.Context, obj *models2.Status) (*models1.Partition, error) {
	// Create projection object
	projection := models1.Projection{}
	// Get projection
	err := utils.ManageSimpleProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}

	// Call business
	return r.BusiServices.PartitionsSvc.FindByID(ctx, obj.PartitionID, &projection)
}

// Status returns generated.StatusResolver implementation.
func (r *Resolver) Status() generated.StatusResolver { return &statusResolver{r} }

type statusResolver struct{ *Resolver }
