package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAnonymousRevocationEvent(ctx context.Context, msg *types.MsgCreateAnonymousRevocationEvent) (*types.MsgCreateAnonymousRevocationEventResponse, error) {
	// Validate input
	if msg.Creator == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator cannot be empty")
	}
	if msg.TargetHash == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "target_hash cannot be empty")
	}
	if msg.RevocationType == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "revocation_type cannot be empty")
	}
	if msg.UrgencyLevel == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "urgency_level cannot be empty")
	}
	if msg.ReasonCategory == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "reason_category cannot be empty")
	}

	// Generate initiator hash from creator address
	initiatorHash := k.Keeper.generateHash(msg.Creator)

	// Create anonymous revocation event using the keeper
	revocation, err := k.Keeper.CreateAnonymousRevocationEvent(ctx, msg.TargetHash, msg.RevocationType, msg.UrgencyLevel, msg.ReasonCategory, initiatorHash)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to create anonymous revocation event")
	}

	// Emit anonymous revocation created event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventAnonymousRevocationCreated{
		RevocationId:   revocation.RevocationId,
		TargetHash:     revocation.TargetHash,
		RevocationType: revocation.RevocationType,
		UrgencyLevel:   revocation.UrgencyLevel,
		Creator:        msg.Creator,
	})

	return &types.MsgCreateAnonymousRevocationEventResponse{
		RevocationId: revocation.RevocationId,
		Status:       "revoked",
		EffectiveAt:  revocation.EffectiveAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
