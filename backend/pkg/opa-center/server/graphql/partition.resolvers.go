package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	models2 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	models1 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/generated"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/mappers"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/model"
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

func (r *partitionResolver) OpaConfiguration(ctx context.Context, obj *models.Partition) (string, error) {
	return r.BusiServices.PartitionsSvc.GenerateOPAConfiguration(ctx, obj.ID)
}

func (r *partitionResolver) Statuses(ctx context.Context, obj *models.Partition, after *string, before *string, first *int, last *int, sort *models1.SortOrder, filter *models1.Filter) (*model.StatusConnection, error) {
	// Create projection object
	projection := models1.Projection{}
	// Get projection
	err := utils.ManageConnectionNodeProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	projection.ID = true

	// Get page input
	pInput, err := utils.GetPageInput(after, before, first, last)
	// Check error
	if err != nil {
		return nil, err
	}

	// Call business
	list, pOut, err := r.BusiServices.StatusSvc.GetAllPaginated(ctx, obj.ID, pInput, sort, filter, &projection)
	// Check error
	if err != nil {
		return nil, err
	}

	// Create result
	var res model.StatusConnection
	// Manage connection
	err = utils.MapConnection(&res, list, pOut)
	// Check error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *partitionResolver) DecisionLogs(ctx context.Context, obj *models.Partition, after *string, before *string, first *int, last *int, sort *models2.SortOrder, filter *models2.Filter) (*model.DecisionLogConnection, error) {
	// Create projection object
	projection := models2.Projection{}
	// Get projection
	err := utils.ManageConnectionNodeProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	projection.ID = true

	// Get page input
	pInput, err := utils.GetPageInput(after, before, first, last)
	// Check error
	if err != nil {
		return nil, err
	}

	// Call business
	list, pOut, err := r.BusiServices.DecisionLogsSvc.GetAllPaginated(ctx, obj.ID, pInput, sort, filter, &projection)
	// Check error
	if err != nil {
		return nil, err
	}

	// Create result
	var res model.DecisionLogConnection
	// Manage connection
	err = utils.MapConnection(&res, list, pOut)
	// Check error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Partition returns generated.PartitionResolver implementation.
func (r *Resolver) Partition() generated.PartitionResolver { return &partitionResolver{r} }

type partitionResolver struct{ *Resolver }
