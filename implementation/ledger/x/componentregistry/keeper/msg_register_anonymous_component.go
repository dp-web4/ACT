package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RegisterAnonymousComponent(ctx context.Context, msg *types.MsgRegisterAnonymousComponent) (*types.MsgRegisterAnonymousComponentResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.RealComponentId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "real_component_id cannot be empty")
	}
	if msg.ManufacturerId == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "manufacturer_id cannot be empty")
	}
	if msg.ComponentType == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_type cannot be empty")
	}

	// Register the component anonymously using the keeper
	component, err := k.Keeper.RegisterAnonymousComponent(ctx, msg.RealComponentId, msg.ManufacturerId, msg.ComponentType)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to register anonymous component")
	}

	// Emit anonymous component registered event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventAnonymousComponentRegistered{
		ComponentHash:    component.ComponentId,
		CategoryHash:     component.CategoryHash,
		ManufacturerHash: component.ManufacturerHash,
		Creator:          msg.Creator,
	})

	return &types.MsgRegisterAnonymousComponentResponse{
		ComponentHash:    component.ComponentId,
		ManufacturerHash: component.ManufacturerHash,
		CategoryHash:     component.CategoryHash,
		Status:           component.Status,
		TrustAnchor:      component.TrustAnchor,
	}, nil
}
