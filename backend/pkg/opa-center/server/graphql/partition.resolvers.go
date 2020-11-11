package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/generated"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/mappers"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/utils"
)

func (r *partitionResolver) ID(ctx context.Context, obj *models.Partition) (string, error) {
	return utils.ToIDRelay(mappers.PartitionIDPrefix, obj.ID), nil
}

func (r *partitionResolver) CreatedAt(ctx context.Context, obj *models.Partition) (string, error) {
	return utils.FormatTime(obj.CreatedAt), nil
}

func (r *partitionResolver) UpdatedAt(ctx context.Context, obj *models.Partition) (string, error) {
	return utils.FormatTime(obj.UpdatedAt), nil
}

// Partition returns generated.PartitionResolver implementation.
func (r *Resolver) Partition() generated.PartitionResolver { return &partitionResolver{r} }

type partitionResolver struct{ *Resolver }
