package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetComponentVerification(ctx context.Context, req *types.QueryGetComponentVerificationRequest) (*types.QueryGetComponentVerificationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ComponentId == "" {
		return nil, status.Error(codes.InvalidArgument, "component_id cannot be empty")
	}

	// Get component from store to check if it exists
	component, err := q.k.Components.Get(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "component not found")
	}

	// Create verification object based on component data
	verification := types.ComponentVerification{
		ComponentId:          req.ComponentId,
		Status:               "verified", // Default status for registered components
		VerifiedAt:           component.LastVerifiedAt,
		VerificationMethod:   "registration",
		VerificationEvidence: "component_registered",
		Notes:                "Component verified during registration",
	}

	return &types.QueryGetComponentVerificationResponse{
		Verification: verification,
	}, nil
}
