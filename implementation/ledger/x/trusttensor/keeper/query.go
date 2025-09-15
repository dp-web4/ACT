package keeper

import (
	"context"

	"racecar-web/x/trusttensor/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	Keeper
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := q.Keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) GetRelationshipTensor(ctx context.Context, req *types.QueryGetRelationshipTensorRequest) (*types.QueryGetRelationshipTensorResponse, error) {
	// TODO: Implement actual logic
	return &types.QueryGetRelationshipTensorResponse{RelationshipTrustTensor: ""}, nil
}

func (q queryServer) CalculateRelationshipTrust(ctx context.Context, req *types.QueryCalculateRelationshipTrustRequest) (*types.QueryCalculateRelationshipTrustResponse, error) {
	// TODO: Implement actual logic
	return &types.QueryCalculateRelationshipTrustResponse{TrustScore: "", Factors: ""}, nil
}

func (q queryServer) GetTensorHistory(ctx context.Context, req *types.QueryGetTensorHistoryRequest) (*types.QueryGetTensorHistoryResponse, error) {
	// TODO: Implement actual logic
	return &types.QueryGetTensorHistoryResponse{TensorEntries: ""}, nil
}
