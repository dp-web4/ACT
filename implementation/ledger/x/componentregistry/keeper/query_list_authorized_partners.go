package keeper

import (
	"context"
	"encoding/json"

	"racecar-web/x/componentregistry/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListAuthorizedPartners(ctx context.Context, req *types.QueryListAuthorizedPartnersRequest) (*types.QueryListAuthorizedPartnersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ComponentId == "" {
		return nil, status.Error(codes.InvalidArgument, "component_id cannot be empty")
	}

	// Get the component
	component, err := q.k.Components.Get(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "component not found")
	}

	// Parse capabilities (using as authorization rules)
	var authRules map[string]interface{}
	if component.Capabilities != "" {
		if err := json.Unmarshal([]byte(component.Capabilities), &authRules); err != nil {
			return &types.QueryListAuthorizedPartnersResponse{
				AuthorizedComponents: "[]",
			}, nil
		}
	}

	// Extract authorized partner types from rules
	var authorizedPartners []string
	if authRules != nil {
		if pairingRules, ok := authRules["pairing"].(map[string]interface{}); ok {
			// Check outbound rules (components this component can pair with)
			if outboundRules, ok := pairingRules["outbound"].(map[string]interface{}); ok {
				if allowedTypes, ok := outboundRules["allowed_types"].([]interface{}); ok {
					for _, allowedType := range allowedTypes {
						if typeStr, ok := allowedType.(string); ok {
							authorizedPartners = append(authorizedPartners, typeStr)
						}
					}
				}
			}

			// Check general pairing rules
			if allowedTypes, ok := pairingRules["allowed_types"].([]interface{}); ok {
				for _, allowedType := range allowedTypes {
					if typeStr, ok := allowedType.(string); ok {
						// Avoid duplicates
						found := false
						for _, existing := range authorizedPartners {
							if existing == typeStr {
								found = true
								break
							}
						}
						if !found {
							authorizedPartners = append(authorizedPartners, typeStr)
						}
					}
				}
			}
		}
	}

	// Convert to JSON string
	authorizedPartnersJSON, err := json.Marshal(authorizedPartners)
	if err != nil {
		return &types.QueryListAuthorizedPartnersResponse{
			AuthorizedComponents: "[]",
		}, nil
	}

	return &types.QueryListAuthorizedPartnersResponse{
		AuthorizedComponents: string(authorizedPartnersJSON),
	}, nil
}
