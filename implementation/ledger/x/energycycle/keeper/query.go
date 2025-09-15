package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"racecar-web/x/energycycle/types"
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

// GetRelationshipEnergyBalance implements the Query/GetRelationshipEnergyBalance RPC method.
func (qs QueryServer) GetRelationshipEnergyBalance(ctx context.Context, req *types.QueryGetRelationshipEnergyBalanceRequest) (*types.QueryGetRelationshipEnergyBalanceResponse, error) {
	if req.LctId == "" {
		return nil, status.Error(codes.InvalidArgument, "LCT ID cannot be empty")
	}

	// Calculate actual energy balance using ATP tokens
	atpBalance, err := qs.Keeper.CalculateEnergyBalance(ctx, req.LctId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// For now, return simplified response with ATP balance
	// In a full implementation, you'd also calculate ADP balance and trust-weighted balance
	return &types.QueryGetRelationshipEnergyBalanceResponse{
		AtpBalance:           atpBalance.String(),
		AdpBalance:           "0.000", // TODO: Calculate ADP balance
		TotalEnergy:          atpBalance.String(),
		TrustWeightedBalance: atpBalance.String(), // TODO: Apply trust weighting
	}, nil
}

// CalculateRelationshipV3 implements the Query/CalculateRelationshipV3 RPC method.
func (qs QueryServer) CalculateRelationshipV3(ctx context.Context, req *types.QueryCalculateRelationshipV3Request) (*types.QueryCalculateRelationshipV3Response, error) {
	if req.OperationId == "" {
		return nil, status.Error(codes.InvalidArgument, "operation ID cannot be empty")
	}

	// Calculate actual V3 composite score using trust tensor
	v3Score, err := qs.Keeper.trusttensorKeeper.CalculateV3CompositeScore(ctx, req.OperationId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCalculateRelationshipV3Response{
		V_3Tensor: v3Score.String(),
	}, nil
}

// GetEnergyFlowHistory implements the Query/GetEnergyFlowHistory RPC method.
func (qs QueryServer) GetEnergyFlowHistory(ctx context.Context, req *types.QueryGetEnergyFlowHistoryRequest) (*types.QueryGetEnergyFlowHistoryResponse, error) {
	if req.LctId == "" {
		return nil, status.Error(codes.InvalidArgument, "LCT ID cannot be empty")
	}

	// For now, return mock energy flow history
	// In a full implementation, this would query the actual energy flow history
	return &types.QueryGetEnergyFlowHistoryResponse{
		EnergyOperations: `[{"operation_id":"energy-op-discharge-12345","source_lct":"lct-MODBATT-PACK-A-HOST-1704067200","target_lct":"lct-MODBATT-MOD-001-PACK-A-1704067200","energy_amount":"50.0","operation_type":"discharge","status":"validated","timestamp":1704067200,"trust_score":"0.85"}]`,
	}, nil
}
