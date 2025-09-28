package keeper

import (
	"context"

	"racecar-web/x/componentregistry/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) GetAnonymousComponentMetadata(ctx context.Context, msg *types.MsgGetAnonymousComponentMetadata) (*types.MsgGetAnonymousComponentMetadataResponse, error) {
	// Validate input
	if msg.Requester == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "requester cannot be empty")
	}
	if msg.ComponentHash == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidComponentID, "component_hash cannot be empty")
	}

	// Get anonymous component metadata using the keeper
	metadata, err := k.Keeper.GetAnonymousComponentMetadata(ctx, msg.ComponentHash)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get anonymous component metadata")
	}

	// Extract non-sensitive metadata
	componentType := "unknown"
	if t, exists := metadata["type"]; exists {
		componentType = t.(string)
	}

	status := "unknown"
	if s, exists := metadata["status"]; exists {
		status = s.(string)
	}

	trustAnchor := "unknown"
	if ta, exists := metadata["trust_anchor"]; exists {
		trustAnchor = ta.(string)
	}

	lastVerified := "unknown"
	if lv, exists := metadata["last_verified"]; exists {
		lastVerified = lv.(string)
	}

	return &types.MsgGetAnonymousComponentMetadataResponse{
		ComponentHash: msg.ComponentHash,
		Type:          componentType,
		Status:        status,
		TrustAnchor:   trustAnchor,
		LastVerified:  lastVerified,
	}, nil
}
