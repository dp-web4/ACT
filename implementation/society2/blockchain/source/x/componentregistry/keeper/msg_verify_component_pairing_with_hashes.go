package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) VerifyComponentPairingWithHashes(ctx context.Context, msg *types.MsgVerifyComponentPairingWithHashes) (*types.MsgVerifyComponentPairingWithHashesResponse, error) {
	// Validate input
	if msg.Verifier == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "verifier cannot be empty")
	}
	if msg.ComponentHashA == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_hash_a cannot be empty")
	}
	if msg.ComponentHashB == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_hash_b cannot be empty")
	}

	// Verify component pairing using hashes using the keeper
	canPair, reason, err := k.Keeper.VerifyComponentPairingWithHashes(ctx, msg.ComponentHashA, msg.ComponentHashB)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to verify component pairing with hashes")
	}

	// Get trust score if components can pair
	trustScore := "0.0"
	if canPair {
		// Get anonymous metadata for trust score calculation using the keeper
		metadataA, err := k.Keeper.GetAnonymousComponentMetadata(ctx, msg.ComponentHashA)
		if err == nil {
			if score, exists := metadataA["trust_score"]; exists {
				trustScore = score.(string)
			}
		}
	}

	// Emit pairing verification event (anonymous)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitTypedEvent(&types.EventComponentVerified{
		ComponentId: msg.ComponentHashA, // Use hash as component ID
		Status:      "pairing_verified",
		Verifier:    msg.Verifier,
	})

	return &types.MsgVerifyComponentPairingWithHashesResponse{
		CanPair:    canPair,
		Reason:     reason,
		TrustScore: trustScore,
	}, nil
}
