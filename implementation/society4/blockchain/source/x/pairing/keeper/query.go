package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"racecar-web/x/pairing/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (q queryServer) ValidateBidirectionalAuth(ctx context.Context, req *types.QueryValidateBidirectionalAuthRequest) (*types.QueryValidateBidirectionalAuthResponse, error) {
	if req.ComponentA == "" || req.ComponentB == "" {
		return nil, status.Error(codes.InvalidArgument, "both component_a and component_b must be provided")
	}

	// For now, we'll use a simplified approach
	// In a real implementation, this would call the component registry query server
	// to check pairing authorization

	// Simplified logic: assume components can pair if they exist and are active
	// This is a placeholder - in production this would use the component registry
	aCanPairB := true // Placeholder
	bCanPairA := true // Placeholder

	return &types.QueryValidateBidirectionalAuthResponse{
		ACanPairB:          aCanPairB,
		BCanPairA:          bCanPairA,
		RequiredConditions: "simplified_auth_check",
	}, nil
}

func (q queryServer) GetPairingStatus(ctx context.Context, req *types.QueryGetPairingStatusRequest) (*types.QueryGetPairingStatusResponse, error) {
	if req.ChallengeId == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge_id must be provided")
	}

	// Get pairing session
	session, err := q.PairingSessions.Get(ctx, req.ChallengeId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pairing session not found")
	}

	// Return the pairing challenge as a JSON string
	pairingChallengeJSON := fmt.Sprintf(`{
		"challenge_id": "%s",
		"status": "%s",
		"expires_at": %d,
		"established_at": %d,
		"session_keys": "%s"
	}`, session.SessionId, session.Status, session.ExpiresAt, session.EstablishedAt, session.SessionKeys)

	return &types.QueryGetPairingStatusResponse{
		PairingChallenge: pairingChallengeJSON,
	}, nil
}

func (q queryServer) ListActivePairings(ctx context.Context, req *types.QueryListActivePairingsRequest) (*types.QueryListActivePairingsResponse, error) {
	var activeSessions []types.PairingSession
	var pairingCount int64

	// Iterate through all pairing sessions
	err := q.PairingSessions.Walk(ctx, nil, func(key string, session types.PairingSession) (bool, error) {
		if session.Status == "completed" || session.Status == "active" {
			activeSessions = append(activeSessions, session)
			pairingCount++
		}
		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to iterate pairing sessions: %s", err))
	}

	// Convert to JSON for response
	activeSessionsJSON, err := json.Marshal(activeSessions)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal active sessions")
	}

	return &types.QueryListActivePairingsResponse{
		ActiveLcts:   string(activeSessionsJSON),
		PairingCount: pairingCount,
	}, nil
}
