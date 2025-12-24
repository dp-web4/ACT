package keeper

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"racecar-web/x/pairing/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (ms msgServer) InitiateBidirectionalPairing(ctx context.Context, msg *types.MsgInitiateBidirectionalPairing) (*types.MsgInitiateBidirectionalPairingResponse, error) {
	if msg.Creator == "" || msg.ComponentA == "" || msg.ComponentB == "" {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "missing required fields")
	}

	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid creator address: %s", err)
	}

	// Check if components exist and are active using component registry
	_, found := ms.componentregistryKeeper.GetComponentIdentity(ctx, msg.ComponentA)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "component A not found: %s", msg.ComponentA)
	}

	_, found = ms.componentregistryKeeper.GetComponentIdentity(ctx, msg.ComponentB)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "component B not found: %s", msg.ComponentB)
	}

	// Check pairing authorization using component registry
	aCanPairB, bCanPairA, reason := ms.componentregistryKeeper.CheckBidirectionalPairingAuth(ctx, msg.ComponentA, msg.ComponentB)
	if !aCanPairB || !bCanPairA {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "bidirectional pairing not authorized: %s", reason)
	}

	// Generate unique challenge ID
	challengeId := fmt.Sprintf("challenge-%s-%s-%d", msg.ComponentA, msg.ComponentB, time.Now().Unix())

	// Generate random challenge data (32 bytes)
	challengeData := make([]byte, 32)
	if _, err := rand.Read(challengeData); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to generate challenge data: %s", err)
	}

	// Create pairing session with correct fields
	session := types.PairingSession{
		SessionId:     challengeId,
		LctId:         fmt.Sprintf("lct-%s-%s-%d", msg.ComponentA, msg.ComponentB, time.Now().Unix()),
		SessionKeys:   "",                                     // Will be set when pairing completes
		EstablishedAt: 0,                                      // Will be set when pairing completes
		ExpiresAt:     time.Now().Add(5 * time.Minute).Unix(), // 5 minute timeout
		Status:        "pending",
	}

	// Store pairing session
	if err := ms.PairingSessions.Set(ctx, challengeId, session); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to store pairing session: %s", err)
	}

	// Create LCT relationship using LCT manager
	lctId, _, err := ms.lctmanagerKeeper.CreateLCTRelationship(ctx, msg.ComponentA, msg.ComponentB, msg.OperationalContext, msg.ProxyId)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to create LCT relationship: %s", err)
	}

	// Check if components are offline and need queueing
	queueId := ""
	if msg.ForceImmediate {
		// For now, we'll skip queueing if force immediate is set
		// In a real implementation, this would check component online status
	} else {
		// Queue the pairing request for offline processing
		// Note: We'll need to implement this properly with the pairing queue module
		// For now, we'll skip queueing to avoid interface issues
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("bidirectional_pairing_initiated",
			sdk.NewAttribute("challenge_id", challengeId),
			sdk.NewAttribute("component_a", msg.ComponentA),
			sdk.NewAttribute("component_b", msg.ComponentB),
			sdk.NewAttribute("lct_id", lctId),
			sdk.NewAttribute("status", "pending"),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("operational_context", msg.OperationalContext),
		),
	)

	return &types.MsgInitiateBidirectionalPairingResponse{
		ChallengeId: challengeId,
		LctId:       lctId,
		Status:      "pending",
		QueueId:     queueId,
	}, nil
}
