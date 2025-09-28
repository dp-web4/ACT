package keeper

import (
	"context"
	"encoding/json"
	"time"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) VerifyComponent(ctx context.Context, msg *types.MsgVerifyComponent) (*types.MsgVerifyComponentResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.ComponentId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_id cannot be empty")
	}

	// Check if component exists
	component, err := k.Components.Get(ctx, msg.ComponentId)
	if err != nil {
		return &types.MsgVerifyComponentResponse{
			IsValid:       false,
			ComponentData: "",
		}, nil
	}

	// Check if component is active
	if component.Status != "active" {
		// Emit verification failed event
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitTypedEvent(&types.EventComponentVerified{
			ComponentId: msg.ComponentId,
			Status:      "failed_inactive",
			Verifier:    msg.Creator,
		})

		return &types.MsgVerifyComponentResponse{
			IsValid:       false,
			ComponentData: "",
		}, nil
	}

	// Create verification record
	verification := types.ComponentVerification{
		ComponentId:          msg.ComponentId,
		Status:               "verified",
		VerifiedAt:           time.Now(),
		VerificationMethod:   "manual",
		VerificationEvidence: "blockchain_verification",
		Notes:                "Verified via blockchain query",
	}

	// Store verification
	if err := k.ComponentVerifications.Set(ctx, msg.ComponentId, verification); err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to store verification")
	}

	// Update component last verified timestamp
	component.LastVerifiedAt = verification.VerifiedAt
	if err := k.Components.Set(ctx, msg.ComponentId, component); err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to update component")
	}

	// Create component data for response (privacy-focused)
	componentData := map[string]interface{}{
		"component_id":             component.ComponentId,
		"category_hash":            component.CategoryHash,
		"manufacturer_hash":        component.ManufacturerHash,
		"authorization_rules_hash": component.AuthorizationRulesHash,
		"status":                   component.Status,
		"registered_at":            component.RegisteredAt,
		"last_verified_at":         component.LastVerifiedAt,
		"verification_metadata":    component.VerificationMetadata,
		"relationship_hashes":      component.RelationshipHashes,
		"trust_anchor":             component.TrustAnchor,
	}

	// Convert to JSON string
	componentDataJSON, err := json.Marshal(componentData)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to marshal component data")
	}

	// Emit verification success event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventComponentVerified{
		ComponentId: msg.ComponentId,
		Status:      "verified",
		Verifier:    msg.Creator,
	})

	return &types.MsgVerifyComponentResponse{
		IsValid:       true,
		ComponentData: string(componentDataJSON),
	}, nil
}
