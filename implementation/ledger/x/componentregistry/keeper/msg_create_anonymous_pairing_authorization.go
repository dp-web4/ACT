package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAnonymousPairingAuthorization(ctx context.Context, msg *types.MsgCreateAnonymousPairingAuthorization) (*types.MsgCreateAnonymousPairingAuthorizationResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.ComponentHashA == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_hash_a cannot be empty")
	}
	if msg.ComponentHashB == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_hash_b cannot be empty")
	}
	if msg.RuleHash == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "rule_hash cannot be empty")
	}

	// Create anonymous pairing authorization using the keeper
	auth, err := k.Keeper.CreateAnonymousPairingAuthorization(ctx, msg.ComponentHashA, msg.ComponentHashB, msg.RuleHash)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to create anonymous pairing authorization")
	}

	// Emit anonymous pairing authorized event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventAnonymousPairingAuthorized{
		AuthId:         auth.AuthId,
		ComponentHashA: auth.ComponentHashA,
		ComponentHashB: auth.ComponentHashB,
		Creator:        msg.Creator,
	})

	return &types.MsgCreateAnonymousPairingAuthorizationResponse{
		AuthId:    auth.AuthId,
		Status:    auth.Status,
		ExpiresAt: auth.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
