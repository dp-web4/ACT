package keeper

import (
	"context"
	"encoding/json"

	"racecar-web/x/componentregistry/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) CheckPairingAuth(ctx context.Context, req *types.QueryCheckPairingAuthRequest) (*types.QueryCheckPairingAuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ComponentA == "" || req.ComponentB == "" {
		return nil, status.Error(codes.InvalidArgument, "both component_a and component_b must be provided")
	}

	// Get component A
	componentA, err := q.k.Components.Get(ctx, req.ComponentA)
	if err != nil {
		return &types.QueryCheckPairingAuthResponse{
			ACanPairB: false,
			BCanPairA: false,
			Reason:    "component_a not found",
		}, nil
	}

	// Get component B
	componentB, err := q.k.Components.Get(ctx, req.ComponentB)
	if err != nil {
		return &types.QueryCheckPairingAuthResponse{
			ACanPairB: false,
			BCanPairA: false,
			Reason:    "component_b not found",
		}, nil
	}

	// Check if components are active
	if componentA.Status != "active" || componentB.Status != "active" {
		return &types.QueryCheckPairingAuthResponse{
			ACanPairB: false,
			BCanPairA: false,
			Reason:    "one or both components are not active",
		}, nil
	}

	// Parse capabilities for component A (using as authorization rules)
	var authRulesA map[string]interface{}
	if componentA.Capabilities != "" {
		if err := json.Unmarshal([]byte(componentA.Capabilities), &authRulesA); err != nil {
			return &types.QueryCheckPairingAuthResponse{
				ACanPairB: false,
				BCanPairA: false,
				Reason:    "invalid capabilities for component_a",
			}, nil
		}
	}

	// Parse capabilities for component B (using as authorization rules)
	var authRulesB map[string]interface{}
	if componentB.Capabilities != "" {
		if err := json.Unmarshal([]byte(componentB.Capabilities), &authRulesB); err != nil {
			return &types.QueryCheckPairingAuthResponse{
				ACanPairB: false,
				BCanPairA: false,
				Reason:    "invalid capabilities for component_b",
			}, nil
		}
	}

	// Check if A can pair with B
	aCanPairB := checkPairingPermission(authRulesA, componentB.CategoryHash, "outbound")

	// Check if B can pair with A
	bCanPairA := checkPairingPermission(authRulesB, componentA.CategoryHash, "inbound")

	reason := "authorization check completed"
	if !aCanPairB && !bCanPairA {
		reason = "neither component has pairing authorization"
	} else if !aCanPairB {
		reason = "component_a cannot pair with component_b"
	} else if !bCanPairA {
		reason = "component_b cannot pair with component_a"
	}

	return &types.QueryCheckPairingAuthResponse{
		ACanPairB: aCanPairB,
		BCanPairA: bCanPairA,
		Reason:    reason,
	}, nil
}

// checkPairingPermission checks if a component can pair with another component type
func checkPairingPermission(authRules map[string]interface{}, targetCategoryHash, direction string) bool {
	if authRules == nil {
		// Default to allowing pairing if no rules are specified
		return true
	}

	// Check for explicit pairing rules
	if pairingRules, ok := authRules["pairing"].(map[string]interface{}); ok {
		// Check direction-specific rules
		if directionRules, ok := pairingRules[direction].(map[string]interface{}); ok {
			// Check if target component type is allowed
			if allowedTypes, ok := directionRules["allowed_types"].([]interface{}); ok {
				for _, allowedType := range allowedTypes {
					if allowedType == targetCategoryHash {
						return true
					}
				}
				return false
			}
		}

		// Check general pairing rules
		if allowedTypes, ok := pairingRules["allowed_types"].([]interface{}); ok {
			for _, allowedType := range allowedTypes {
				if allowedType == targetCategoryHash {
					return true
				}
			}
			return false
		}
	}

	// Default to allowing pairing if no specific rules are found
	return true
}
