package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"racecar-web/x/pairing/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (ms msgServer) CompletePairing(ctx context.Context, msg *types.MsgCompletePairing) (*types.MsgCompletePairingResponse, error) {
	if msg.Creator == "" || msg.ChallengeId == "" || msg.ComponentAAuth == "" || msg.ComponentBAuth == "" {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "missing required fields")
	}

	// Validate creator address format (we'll add authorization logic later)
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid creator address: %s", err)
	}

	// Get the pairing session
	session, err := ms.PairingSessions.Get(ctx, msg.ChallengeId)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "pairing session not found: %s", msg.ChallengeId)
	}

	// Check if session has expired
	if session.ExpiresAt < time.Now().Unix() {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "pairing session has expired")
	}

	// Validate component A authentication
	expectedHashA := sha256.Sum256([]byte(msg.ChallengeId + "component_a"))
	if hex.EncodeToString(expectedHashA[:]) != msg.ComponentAAuth {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "component A authentication failed")
	}

	// Validate component B authentication
	expectedHashB := sha256.Sum256([]byte(msg.ChallengeId + "component_b"))
	if hex.EncodeToString(expectedHashB[:]) != msg.ComponentBAuth {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "component B authentication failed")
	}

	// Update session status to completed
	session.Status = "completed"
	session.EstablishedAt = time.Now().Unix()
	session.SessionKeys = fmt.Sprintf("session_keys_%s_%d", msg.ChallengeId, time.Now().Unix())

	// Store updated session
	if err := ms.PairingSessions.Set(ctx, msg.ChallengeId, session); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "failed to update pairing session: %s", err)
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("pairing_completed",
			sdk.NewAttribute("challenge_id", msg.ChallengeId),
			sdk.NewAttribute("lct_id", session.LctId),
			sdk.NewAttribute("status", "completed"),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("established_at", fmt.Sprintf("%d", session.EstablishedAt)),
		),
	)

	return &types.MsgCompletePairingResponse{
		LctId:        session.LctId,
		SessionKeys:  session.SessionKeys,
		TrustSummary: fmt.Sprintf("trust_score:0.85,context:pairing_completed"),
	}, nil
}
