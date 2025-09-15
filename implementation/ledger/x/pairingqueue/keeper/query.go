package keeper

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"racecar-web/x/pairingqueue/types"
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

// GetQueuedRequests implements the Query/GetQueuedRequests RPC method.
func (qs QueryServer) GetQueuedRequests(ctx context.Context, req *types.QueryGetQueuedRequestsRequest) (*types.QueryGetQueuedRequestsResponse, error) {
	if req.ComponentId == "" {
		return nil, status.Error(codes.InvalidArgument, "component ID cannot be empty")
	}

	requests, err := qs.Keeper.GetQueuedRequests(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	requestsJSON, err := json.Marshal(requests)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal requests")
	}

	return &types.QueryGetQueuedRequestsResponse{
		PairingRequests: string(requestsJSON),
	}, nil
}

// GetRequestStatus implements the Query/GetRequestStatus RPC method.
func (qs QueryServer) GetRequestStatus(ctx context.Context, req *types.QueryGetRequestStatusRequest) (*types.QueryGetRequestStatusResponse, error) {
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request ID cannot be empty")
	}

	request, found := qs.Keeper.GetPairingRequest(ctx, req.RequestId)
	if !found {
		return nil, status.Error(codes.NotFound, "request not found")
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal request")
	}

	return &types.QueryGetRequestStatusResponse{
		PairingRequest: string(requestJSON),
	}, nil
}

// ListProxyQueue implements the Query/ListProxyQueue RPC method.
func (qs QueryServer) ListProxyQueue(ctx context.Context, req *types.QueryListProxyQueueRequest) (*types.QueryListProxyQueueResponse, error) {
	if req.ProxyId == "" {
		return nil, status.Error(codes.InvalidArgument, "proxy ID cannot be empty")
	}

	operations, err := qs.Keeper.ListProxyQueue(ctx, req.ProxyId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	operationsJSON, err := json.Marshal(operations)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal operations")
	}

	return &types.QueryListProxyQueueResponse{
		OfflineOperations: string(operationsJSON),
	}, nil
}
