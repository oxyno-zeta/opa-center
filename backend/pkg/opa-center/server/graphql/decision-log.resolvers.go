package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	models1 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/generated"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/mappers"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/utils"
)

func (r *decisionLogResolver) ID(ctx context.Context, obj *models.DecisionLog) (string, error) {
	return utils.ToIDRelay(mappers.DecisionLogIDPrefix, obj.ID), nil
}

func (r *decisionLogResolver) CreatedAt(ctx context.Context, obj *models.DecisionLog) (string, error) {
	return utils.FormatTime(obj.CreatedAt), nil
}

func (r *decisionLogResolver) UpdatedAt(ctx context.Context, obj *models.DecisionLog) (string, error) {
	return utils.FormatTime(obj.UpdatedAt), nil
}

func (r *decisionLogResolver) Timestamp(ctx context.Context, obj *models.DecisionLog) (string, error) {
	return utils.FormatTime(obj.Timestamp), nil
}

func (r *decisionLogResolver) Partition(ctx context.Context, obj *models.DecisionLog) (*models1.Partition, error) {
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

// DecisionLog returns generated.DecisionLogResolver implementation.
func (r *Resolver) DecisionLog() generated.DecisionLogResolver { return &decisionLogResolver{r} }

type decisionLogResolver struct{ *Resolver }
