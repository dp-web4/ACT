package keeper

import (
	"context"
	"fmt"

	"racecar-web/x/componentregistry/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetComponent(ctx context.Context, req *types.QueryGetComponentRequest) (*types.QueryGetComponentResponse, error) {
	if req == nil {
		fmt.Printf("DEBUG: Request is nil\n")
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ComponentId == "" {
		fmt.Printf("DEBUG: ComponentId is empty\n")
		return nil, status.Error(codes.InvalidArgument, "component_id cannot be empty")
	}

	fmt.Printf("DEBUG: Looking for component: %s\n", req.ComponentId)

	// Get component from store
	component, err := q.k.Components.Get(ctx, req.ComponentId)
	if err != nil {
		fmt.Printf("DEBUG: Error getting component: %v\n", err)
		return nil, status.Error(codes.NotFound, "component not found")
	}

	// Debug: Print what we're returning
	fmt.Printf("DEBUG: GetComponent returning component: %+v\n", component)

	response := &types.QueryGetComponentResponse{
		Component: component,
	}

	// Debug: Print the response
	fmt.Printf("DEBUG: GetComponent response: %+v\n", response)

	return response, nil
}
