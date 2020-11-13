package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	models1 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/decisionlogs/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business/partitions/models"
	models3 "github.com/oxyno-zeta/opa-center/pkg/opa-center/business/statuses/models"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/generated"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/mappers"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/model"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/server/graphql/utils"
)

func (r *mutationResolver) CreatePartition(ctx context.Context, input models.CreateInput) (*model.CreatePartitionPayload, error) {
	// Call business
	part, err := r.BusiServices.PartitionsSvc.Create(ctx, &input)
	// Check error
	if err != nil {
		return nil, err
	}

	return &model.CreatePartitionPayload{Partition: part}, nil
}

func (r *queryResolver) Partitions(ctx context.Context, after *string, before *string, first *int, last *int, sort *models.SortOrder, filter *models.Filter) (*model.PartitionConnection, error) {
	// Create projection object
	projection := models.Projection{}
	// Get projection
	err := utils.ManageSimpleProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	// This is forced to have generate urls
	projection.ID = true

	// Get page input
	pInput, err := utils.GetPageInput(after, before, first, last)
	// Check error
	if err != nil {
		return nil, err
	}

	// Get partitions
	partitions, pOut, err := r.BusiServices.PartitionsSvc.GetAllPaginated(ctx, pInput, sort, filter, &projection)
	// Check error
	if err != nil {
		return nil, err
	}

	// Create connection
	conn := model.PartitionConnection{}
	// Map connection
	err = utils.MapConnection(&conn, partitions, pOut)
	// Check error
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (r *queryResolver) Partition(ctx context.Context, id string) (*models.Partition, error) {
	// Create projection object
	projection := models.Projection{}
	// Get projection
	err := utils.ManageSimpleProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	// This is forced to have generate urls
	projection.ID = true

	// Transform relay id to business id
	bid, err := utils.FromIDRelay(id, mappers.PartitionIDPrefix)
	// Check error
	if err != nil {
		return nil, err
	}

	// Get partition
	return r.BusiServices.PartitionsSvc.FindByID(ctx, bid, &projection)
}

func (r *queryResolver) OpaConfiguration(ctx context.Context, partitionID string) (string, error) {
	// Transform relay id to business id
	bid, err := utils.FromIDRelay(partitionID, mappers.PartitionIDPrefix)
	// Check error
	if err != nil {
		return "", err
	}

	return r.BusiServices.PartitionsSvc.GenerateOPAConfiguration(ctx, bid)
}

func (r *queryResolver) DecisionLogs(ctx context.Context, partitionID string, after *string, before *string, first *int, last *int, sort *models1.SortOrder, filter *models1.Filter) (*model.DecisionLogConnection, error) {
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

	// Transform relay id to business id
	bPartitionID, err := utils.FromIDRelay(partitionID, mappers.PartitionIDPrefix)
	// Check error
	if err != nil {
		return nil, err
	}

	// Get page input
	pInput, err := utils.GetPageInput(after, before, first, last)
	// Check error
	if err != nil {
		return nil, err
	}

	// Call business
	list, pOut, err := r.BusiServices.DecisionLogsSvc.GetAllPaginated(ctx, bPartitionID, pInput, sort, filter, &projection)
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

func (r *queryResolver) DecisionLog(ctx context.Context, id *string, decisionLogID *string) (*models1.DecisionLog, error) {
	// Create projection
	projection := models1.Projection{}
	// Get projection
	err := utils.ManageSimpleProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}

	// Transform id relay to business id
	var bid *string
	// Check if id exists
	if id != nil {
		bidst, err := utils.FromIDRelay(*id, mappers.DecisionLogIDPrefix)
		// Check error
		if err != nil {
			return nil, err
		}

		bid = &bidst
	}

	// Call business
	return r.BusiServices.DecisionLogsSvc.FindByIDOrDecisionID(ctx, bid, decisionLogID, &projection)
}

func (r *queryResolver) Statuses(ctx context.Context, partitionID string, after *string, before *string, first *int, last *int, sort *models3.SortOrder, filter *models3.Filter) (*model.StatusConnection, error) {
	// Create projection object
	projection := models3.Projection{}
	// Get projection
	err := utils.ManageConnectionNodeProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	projection.ID = true

	// Transform relay id to business id
	bPartitionID, err := utils.FromIDRelay(partitionID, mappers.PartitionIDPrefix)
	// Check error
	if err != nil {
		return nil, err
	}

	// Get page input
	pInput, err := utils.GetPageInput(after, before, first, last)
	// Check error
	if err != nil {
		return nil, err
	}

	// Call business
	list, pOut, err := r.BusiServices.StatusSvc.GetAllPaginated(ctx, bPartitionID, pInput, sort, filter, &projection)
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

func (r *queryResolver) Status(ctx context.Context, id string) (*models3.Status, error) {
	// Create projection object
	projection := models3.Projection{}
	// Get projection
	err := utils.ManageSimpleProjection(ctx, &projection)
	// Check error
	if err != nil {
		return nil, err
	}
	// Ask for id projection
	// This is forced to have generate urls
	projection.ID = true

	// Transform relay id to business id
	bid, err := utils.FromIDRelay(id, mappers.StatusIDPrefix)
	// Check error
	if err != nil {
		return nil, err
	}

	// Get partition
	return r.BusiServices.StatusSvc.FindByID(ctx, bid, &projection)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
