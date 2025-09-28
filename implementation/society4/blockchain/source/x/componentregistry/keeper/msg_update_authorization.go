package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateAuthorization(ctx context.Context, msg *types.MsgUpdateAuthorization) (*types.MsgUpdateAuthorizationResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.ComponentId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_id cannot be empty")
	}
	if msg.AuthRules == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentType, "auth_rules cannot be empty")
	}

	// Validate auth_rules is valid JSON
	var authRulesMap map[string]interface{}
	if err := json.Unmarshal([]byte(msg.AuthRules), &authRulesMap); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentType, "auth_rules must be valid JSON")
	}

	// Check if component exists
	component, err := k.Components.Get(ctx, msg.ComponentId)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "component not found")
	}

	// Extract authorization parameters from auth rules
	authLevel, ok := authRulesMap["authorization_level"].(string)
	if !ok {
		authLevel = AuthLevelBasic // Default to basic
	}

	lctID, ok := authRulesMap["lct_id"].(string)
	if !ok {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentType, "lct_id is required in auth_rules")
	}

	context, ok := authRulesMap["context"].(string)
	if !ok {
		context = "component_authorization" // Default context
	}

	// Create new pairing authorization using trust-based logic
	authorization, err := k.Keeper.CreatePairingAuthorization(ctx, msg.ComponentId, lctID, authLevel, context)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, fmt.Sprintf("failed to create authorization: %v", err))
	}

	// Update component capabilities with authorization info
	component.Capabilities = msg.AuthRules

	// Store the updated component
	if err := k.Components.Set(ctx, msg.ComponentId, component); err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to update component")
	}

	// Emit authorization updated event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventAuthorizationUpdated{
		ComponentId: msg.ComponentId,
		Updater:     msg.Creator,
	})

	// Log authorization creation (using the authorization variable)
	_ = authorization // Use the variable to avoid unused variable error

	return &types.MsgUpdateAuthorizationResponse{}, nil
}
