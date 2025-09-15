package keeper

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"racecar-web/x/lctmanager/types"
)

type QueryServer struct {
	Keeper
}

var _ types.QueryServer = QueryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

// Params implements the Query/Params RPC method.
func (qs QueryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := qs.Keeper.Params.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

// GetLct implements the Query/GetLct RPC method.
func (qs QueryServer) GetLct(ctx context.Context, req *types.QueryGetLctRequest) (*types.QueryGetLctResponse, error) {
	if req.LctId == "" {
		return nil, status.Error(codes.InvalidArgument, "LCT ID cannot be empty")
	}

	lct, found := qs.Keeper.GetLct(ctx, req.LctId)
	if !found {
		return nil, status.Error(codes.NotFound, "LCT not found")
	}

	lctJSON, err := json.Marshal(lct)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal LCT")
	}

	return &types.QueryGetLctResponse{
		LinkedContextToken: string(lctJSON),
	}, nil
}

// GetComponentRelationships implements the Query/GetComponentRelationships RPC method.
func (qs QueryServer) GetComponentRelationships(ctx context.Context, req *types.QueryGetComponentRelationshipsRequest) (*types.QueryGetComponentRelationshipsResponse, error) {
	if req.ComponentId == "" {
		return nil, status.Error(codes.InvalidArgument, "component ID cannot be empty")
	}

	lcts, err := qs.Keeper.GetComponentRelationships(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	lctsJSON, err := json.Marshal(lcts)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal LCTs")
	}

	return &types.QueryGetComponentRelationshipsResponse{
		ComponentRelationships: string(lctsJSON),
		LctCount:               int64(len(lcts)),
	}, nil
}

// ValidateLctAccess implements the Query/ValidateLctAccess RPC method.
func (qs QueryServer) ValidateLctAccess(ctx context.Context, req *types.QueryValidateLctAccessRequest) (*types.QueryValidateLctAccessResponse, error) {
	if req.LctId == "" {
		return nil, status.Error(codes.InvalidArgument, "LCT ID cannot be empty")
	}
	if req.RequestorId == "" {
		return nil, status.Error(codes.InvalidArgument, "requestor ID cannot be empty")
	}

	hasAccess, accessLevel, err := qs.Keeper.ValidateLctAccess(ctx, req.LctId, req.RequestorId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryValidateLctAccessResponse{
		HasAccess:   hasAccess,
		AccessLevel: accessLevel,
	}, nil
}
