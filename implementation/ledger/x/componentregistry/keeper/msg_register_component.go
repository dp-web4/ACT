package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RegisterComponent(ctx context.Context, msg *types.MsgRegisterComponent) (*types.MsgRegisterComponentResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.ComponentId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_id cannot be empty")
	}
	if msg.ComponentType == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentType, "component_type cannot be empty")
	}

	// Validate component ID format (should be alphanumeric with optional hyphens/underscores)
	if !isValidComponentID(msg.ComponentId) {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_id must be alphanumeric with optional hyphens/underscores")
	}

	// Validate component type (should be one of the standard types)
	if !isValidComponentType(msg.ComponentType) {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentType, "component_type must be one of: module, pack, host_ecu, sensor, actuator")
	}

	// Check if component already exists
	_, err := k.Components.Get(ctx, msg.ComponentId)
	if err == nil {
		return nil, errorsmod.Wrap(types.ErrComponentExists, fmt.Sprintf("component %s already registered", msg.ComponentId))
	}

	// Extract manufacturer ID from manufacturer data
	manufacturerId := extractManufacturerID(msg.ManufacturerData)
	if manufacturerId == "" {
		manufacturerId = "unknown" // Default if not found
	}

	// Create component identity
	componentIdentity := fmt.Sprintf("comp_%s_%s", msg.ComponentType, msg.ComponentId)

	// Generate LCT ID for this component
	lctId := fmt.Sprintf("lct-%s-%s-%d", msg.ComponentType, msg.ComponentId, time.Now().Unix())

	// Create the component with LCT integration
	component := types.Component{
		ComponentId:            msg.ComponentId,
		ManufacturerId:         manufacturerId,
		ComponentType:          msg.ComponentType,
		HardwareSpecs:          msg.ManufacturerData,
		QualityScore:           "0.8", // Default quality score
		Status:                 "active",
		CreatedAt:              time.Now(),
		LastVerifiedAt:         time.Now(),
		Capabilities:           "{}",
		RelationshipIds:        []string{},
		LctId:                  lctId,
		EncryptedDeviceKeyHalf: []byte{}, // Will be set when device pairs with LCT
		TrustAnchor:            msg.Creator,
		AuthorizationRules:     "{\"allowed_pairing_types\": [\"module_to_pack\", \"pack_to_host\"], \"max_connections\": 2, \"trust_threshold\": 0.75}",
	}

	// Store the component
	if err := k.Components.Set(ctx, msg.ComponentId, component); err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to store component")
	}

	// Update manufacturer index
	if err := k.ManufacturerComponents.Set(ctx, manufacturerId, component); err != nil {
		return nil, errorsmod.Wrap(types.ErrComponentNotFound, "failed to update manufacturer index")
	}

	// Emit event with LCT information
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventComponentRegistered{
		ComponentId:    msg.ComponentId,
		ComponentType:  msg.ComponentType,
		ManufacturerId: manufacturerId,
		Creator:        msg.Creator,
	})

	// Emit additional event for LCT creation
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("component_lct_created",
			sdk.NewAttribute("component_id", msg.ComponentId),
			sdk.NewAttribute("lct_id", lctId),
			sdk.NewAttribute("trust_anchor", msg.Creator),
			sdk.NewAttribute("status", "pending_device_pairing"),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	return &types.MsgRegisterComponentResponse{
		ComponentIdentity: componentIdentity,
		LctId:             lctId,
		Status:            "registered_with_lct",
	}, nil
}

// Helper functions for validation
func isValidComponentID(componentId string) bool {
	if len(componentId) < 3 || len(componentId) > 50 {
		return false
	}

	// Check if contains only alphanumeric characters, hyphens, and underscores
	for _, char := range componentId {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}

func isValidComponentType(componentType string) bool {
	validTypes := []string{"module", "pack", "host_ecu", "sensor", "actuator"}
	for _, validType := range validTypes {
		if componentType == validType {
			return true
		}
	}
	return false
}

func extractManufacturerID(manufacturerData string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(manufacturerData), &data); err == nil {
		if id, ok := data["manufacturer_id"].(string); ok && id != "" {
			return id
		}
		if id, ok := data["mfg_id"].(string); ok && id != "" {
			return id
		}
	}
	return "unknown"
}
