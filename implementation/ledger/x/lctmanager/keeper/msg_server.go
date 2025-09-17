package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"racecar-web/x/lctmanager/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// UpdateParams implements the Msg/UpdateParams RPC method.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Convert authority to string for comparison
	authorityStr, err := ms.addressCodec.BytesToString(ms.authority)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "invalid authority encoding: %s", err)
	}

	if authorityStr != msg.Authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "invalid authority; expected %s, got %s", authorityStr, msg.Authority)
	}

	if err := ms.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// MintLCT implements the Msg/MintLCT RPC method for creating new Linked Context Tokens.
func (ms msgServer) MintLCT(ctx context.Context, msg *types.MsgMintLCT) (*types.MsgMintLCTResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate entity type
	validTypes := map[string]bool{
		"agent":   true,
		"human":   true,
		"device":  true,
		"service": true,
		"swarm":   true,
	}
	if !validTypes[msg.EntityType] {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "invalid entity type: %s", msg.EntityType)
	}

	// Generate unique LCT ID
	lctId := fmt.Sprintf("lct-%s-%s-%d", msg.EntityType, msg.EntityName, time.Now().UnixNano())

	// Prevent duplicate LCTs
	if _, found := ms.Keeper.GetLinkedContextToken(ctx, lctId); found {
		return nil, errors.Wrapf(types.ErrLctExists, "LCT already exists: %s", lctId)
	}

	// Create entity address (derived from LCT ID for uniqueness)
	hash := sha256.Sum256([]byte(lctId))
	entityAddress := sdk.AccAddress(hash[:20]).String()

	// Parse initial ADP amount (default to "1000" if empty)
	adpAmount := msg.InitialAdpAmount
	if adpAmount == "" {
		adpAmount = "1000"
	}

	// Create LCT struct for the entity
	lct := types.LinkedContextToken{
		LctId:              lctId,
		ComponentAId:       msg.EntityName,
		ComponentBId:       msg.EntityType,
		PairingStatus:      "active",
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		TrustAnchor:        creator.String(),
		OperationalContext: fmt.Sprintf("%s:%s", msg.EntityType, msg.EntityName),
		ProxyComponentId:   "", // No proxy for direct entities
		LctKeyHalf:         "", // Key management handled off-chain
		LastContactAt:      time.Now().Unix(),
		AuthorizationRules: "{}", // Default empty rules
	}

	// Store metadata if provided
	if len(msg.Metadata) > 0 {
		metadataStr := "{"
		for k, v := range msg.Metadata {
			metadataStr += fmt.Sprintf(`"%s":"%s",`, k, v)
		}
		metadataStr = metadataStr[:len(metadataStr)-1] + "}"
		lct.AuthorizationRules = metadataStr
	}

	// Store the LCT
	if err := ms.Keeper.SetLinkedContextToken(ctx, lct); err != nil {
		return nil, err
	}

	// TODO: Initialize T3/V3 tensors when trusttensor module is ready
	// TODO: Allocate initial ADP tokens when energycycle module is ready

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_minted",
			sdk.NewAttribute("lct_id", lctId),
			sdk.NewAttribute("entity_name", msg.EntityName),
			sdk.NewAttribute("entity_type", msg.EntityType),
			sdk.NewAttribute("entity_address", entityAddress),
			sdk.NewAttribute("adp_balance", adpAmount),
			sdk.NewAttribute("creator", creator.String()),
		),
	)

	return &types.MsgMintLCTResponse{
		LctId:         lctId,
		EntityAddress: entityAddress,
		AdpBalance:    adpAmount,
		Status:        "minted",
	}, nil
}

// CreateLctRelationship implements the Msg/CreateLctRelationship RPC method.
func (ms msgServer) CreateLctRelationship(ctx context.Context, msg *types.MsgCreateLctRelationship) (*types.MsgCreateLctRelationshipResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}
	if msg.ComponentA == "" || msg.ComponentB == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "component IDs cannot be empty")
	}

	// Generate unique LCT ID
	lctId := fmt.Sprintf("lct-%s-%s-%d", msg.ComponentA, msg.ComponentB, time.Now().Unix())

	// Prevent duplicate LCTs
	if _, found := ms.Keeper.GetLinkedContextToken(ctx, lctId); found {
		return nil, errors.Wrapf(types.ErrLctExists, "LCT already exists: %s", lctId)
	}

	// Create LCT struct (no key halves stored)
	lct := types.LinkedContextToken{
		LctId:              lctId,
		ComponentAId:       msg.ComponentA,
		ComponentBId:       msg.ComponentB,
		PairingStatus:      "pending",
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		TrustAnchor:        creator.String(),
		OperationalContext: msg.Context,
		ProxyComponentId:   msg.ProxyId,
		LctKeyHalf:         "", // No key half stored on-chain
		LastContactAt:      time.Now().Unix(),
		AuthorizationRules: "{}", // Default empty rules
	}

	if err := ms.Keeper.SetLinkedContextToken(ctx, lct); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_relationship_created",
			sdk.NewAttribute("lct_id", lctId),
			sdk.NewAttribute("component_a", msg.ComponentA),
			sdk.NewAttribute("component_b", msg.ComponentB),
			sdk.NewAttribute("status", "pending"),
		),
	)

	return &types.MsgCreateLctRelationshipResponse{
		LctId:       lctId,
		KeyExchange: "", // Key exchange handled off-chain or via ephemeral event, not stored
		Status:      "pending",
	}, nil
}

// UpdateLctStatus implements the Msg/UpdateLctStatus RPC method.
func (ms msgServer) UpdateLctStatus(ctx context.Context, msg *types.MsgUpdateLctStatus) (*types.MsgUpdateLctStatusResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Get LCT
	lct, found := ms.Keeper.GetLinkedContextToken(ctx, msg.LctId)
	if !found {
		return nil, types.ErrLctNotFound
	}

	// Verify creator is authorized
	if creator.String() != lct.TrustAnchor {
		return nil, types.ErrInvalidSigner
	}

	// Validate the new status
	if !types.IsValidLCTStatus(msg.NewStatus) {
		return nil, types.ErrInvalidLctStatus
	}

	// Update status
	if err := ms.Keeper.UpdateLctStatus(ctx, msg.LctId, msg.NewStatus, msg.Reason); err != nil {
		return nil, err
	}

	return &types.MsgUpdateLctStatusResponse{}, nil
}

// TerminateLctRelationship implements the Msg/TerminateLctRelationship RPC method.
func (ms msgServer) TerminateLctRelationship(ctx context.Context, msg *types.MsgTerminateLctRelationship) (*types.MsgTerminateLctRelationshipResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Get LCT
	lct, found := ms.Keeper.GetLinkedContextToken(ctx, msg.LctId)
	if !found {
		return nil, types.ErrLctNotFound
	}

	// Verify creator is authorized
	if creator.String() != lct.TrustAnchor {
		return nil, types.ErrInvalidSigner
	}

	// Update status to terminated
	if err := ms.Keeper.UpdateLctStatus(ctx, msg.LctId, types.StatusTerminated, msg.Reason); err != nil {
		return nil, err
	}

	// Handle offline notification if requested
	if msg.NotifyOffline && lct.ProxyComponentId != "" && ms.pairingqueueKeeper != nil {
		// Use the pairingqueue keeper to queue offline operation
		// This is a simplified approach - in a real implementation you'd want more robust offline handling
		_, err := ms.pairingqueueKeeper.QueueOfflineOperation(ctx, lct.ProxyComponentId, "termination_notification")
		if err != nil {
			// Log the error but don't fail the termination
			// The proxy component can still be notified when it comes online
			ms.logger.Error("failed to queue offline operation",
				"error", err,
				"lct_id", msg.LctId,
				"proxy_id", lct.ProxyComponentId,
			)
		}
	}

	return &types.MsgTerminateLctRelationshipResponse{}, nil
}

// InitiateLCTMediatedPairing implements the Msg/InitiateLCTMediatedPairing RPC method.
func (ms msgServer) InitiateLCTMediatedPairing(ctx context.Context, msg *types.MsgInitiateLCTMediatedPairing) (*types.MsgInitiateLCTMediatedPairingResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.InitiatorLctId == "" || msg.TargetLctId == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "initiator and target LCT IDs cannot be empty")
	}

	// Get initiator LCT
	initiatorLct, found := ms.Keeper.GetLinkedContextToken(ctx, msg.InitiatorLctId)
	if !found {
		return nil, errors.Wrapf(types.ErrLctNotFound, "initiator LCT not found: %s", msg.InitiatorLctId)
	}

	// Get target LCT
	targetLct, found := ms.Keeper.GetLinkedContextToken(ctx, msg.TargetLctId)
	if !found {
		return nil, errors.Wrapf(types.ErrLctNotFound, "target LCT not found: %s", msg.TargetLctId)
	}

	// Verify creator is authorized (must be trust anchor of either LCT)
	if creator.String() != initiatorLct.TrustAnchor && creator.String() != targetLct.TrustAnchor {
		return nil, types.ErrInvalidSigner
	}

	// Generate unique pairing ID
	pairingId := fmt.Sprintf("pairing-%s-%s-%d", msg.InitiatorLctId, msg.TargetLctId, time.Now().Unix())

	// Create LCT relationship ID
	lctRelationshipId := fmt.Sprintf("lct-rel-%s-%s-%d", initiatorLct.ComponentAId, targetLct.ComponentAId, time.Now().Unix())

	// Generate cryptographically secure challenge data using our crypto primitives
	challengeKeyShare, err := GenerateKeyShare()
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to generate challenge key: %s", err)
	}

	// Create challenge data using the key share
	challengeData := append(challengeKeyShare[:], []byte(fmt.Sprintf("-%s-%s-%d", pairingId, lctRelationshipId, time.Now().Unix()))...)

	// Hash the challenge for expected response
	expectedResponse := sha256.Sum256(challengeData)

	// Securely zero the challenge key share
	ZeroKey(&challengeKeyShare)

	// Create pairing challenge
	challenge := types.PairingChallenge{
		ChallengeId:      fmt.Sprintf("challenge-%s", pairingId),
		PairingId:        pairingId,
		ChallengeData:    challengeData,
		ExpectedResponse: expectedResponse[:],
		ExpiresAt:        msg.ExpiresAt,
		Status:           "pending",
	}

	// Store pairing challenge (in real implementation, this would be stored in keeper)
	// For now, we'll just use it for the response

	// Emit event for off-chain systems
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_mediated_pairing_initiated",
			sdk.NewAttribute("pairing_id", pairingId),
			sdk.NewAttribute("initiator_lct_id", msg.InitiatorLctId),
			sdk.NewAttribute("target_lct_id", msg.TargetLctId),
			sdk.NewAttribute("context", msg.Context),
			sdk.NewAttribute("lct_relationship_id", lctRelationshipId),
			sdk.NewAttribute("status", "initiated"),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	return &types.MsgInitiateLCTMediatedPairingResponse{
		PairingId:         pairingId,
		Status:            "initiated",
		ChallengeId:       challenge.ChallengeId,
		ChallengeData:     challengeData,
		LctRelationshipId: lctRelationshipId,
	}, nil
}

// CompleteLCTMediatedPairing implements the Msg/CompleteLCTMediatedPairing RPC method.
func (ms msgServer) CompleteLCTMediatedPairing(ctx context.Context, msg *types.MsgCompleteLCTMediatedPairing) (*types.MsgCompleteLCTMediatedPairingResponse, error) {
	// Validate creator address format
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.PairingId == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "pairing ID cannot be empty")
	}

	if len(msg.SessionKeyData) == 0 {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "session key data must be provided")
	}

	// Generate split-key pair for this pairing session
	keyShareA, keyShareB, err := ms.generateSplitKeyPair()
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to generate split keys: %s", err)
	}

	// Encrypt session key data using split-key cryptography
	encryptedSessionKeyInitiator, err := ms.encryptWithSplitKey(keyShareA, keyShareB, msg.SessionKeyData)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to encrypt initiator session key: %s", err)
	}

	// For the target, we use a different key combination (in real implementation,
	// each LCT would have its own key share)
	encryptedSessionKeyTarget, err := ms.encryptWithSplitKey(keyShareB, keyShareA, msg.SessionKeyData)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to encrypt target session key: %s", err)
	}

	// Hash the combined session key for audit (never store the actual key)
	hashedCombinedSessionKey := sha256.Sum256(msg.SessionKeyData)

	// Store the SplitKey metadata (no cryptographic secrets)
	splitKey := types.SplitKey{
		LctId:       fmt.Sprintf("split-%s", msg.PairingId),
		Status:      "active",
		CreatedAt:   time.Now().Unix(),
		ActivatedAt: time.Now().Unix(),
	}

	// Store the split key metadata
	if err := ms.SplitKeys.Set(ctx, splitKey.LctId, splitKey); err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to store split key metadata: %s", err)
	}

	// Calculate trust score (simplified - in real implementation, this would be calculated from T3/V3 tensors)
	trustScore := "0.85"

	// Emit event for off-chain audit
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_mediated_pairing_completed",
			sdk.NewAttribute("pairing_id", msg.PairingId),
			sdk.NewAttribute("split_key_id", splitKey.LctId),
			sdk.NewAttribute("status", "completed"),
			sdk.NewAttribute("hashed_combined_session_key", hex.EncodeToString(hashedCombinedSessionKey[:])),
			sdk.NewAttribute("trust_score", trustScore),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	// Securely zero the key shares from memory
	ZeroKey(&keyShareA)
	ZeroKey(&keyShareB)

	return &types.MsgCompleteLCTMediatedPairingResponse{
		PairingId:                    msg.PairingId,
		Status:                       "completed",
		LctRelationshipId:            fmt.Sprintf("lct-rel-%s", msg.PairingId),
		EncryptedSessionKeyInitiator: encryptedSessionKeyInitiator,
		EncryptedSessionKeyTarget:    encryptedSessionKeyTarget,
		HashedCombinedSessionKey:     hashedCombinedSessionKey[:],
		TrustScore:                   trustScore,
	}, nil
}

// EncryptLCTMessage implements the Msg/EncryptLCTMessage RPC method.
func (ms msgServer) EncryptLCTMessage(ctx context.Context, msg *types.MsgEncryptLCTMessage) (*types.MsgEncryptLCTMessageResponse, error) {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.LctId == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "LCT ID cannot be empty")
	}

	if len(msg.Message) == 0 {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "message cannot be empty")
	}

	// Encrypt the message using our cryptographic operations
	encryptedMessage, err := ms.EncryptMessageForLCT(ctx, msg.LctId, msg.Message)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to encrypt message: %s", err)
	}

	// Emit event for audit
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_message_encrypted",
			sdk.NewAttribute("lct_id", msg.LctId),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("message_length", fmt.Sprintf("%d", len(msg.Message))),
			sdk.NewAttribute("encrypted_length", fmt.Sprintf("%d", len(encryptedMessage))),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	return &types.MsgEncryptLCTMessageResponse{
		LctId:            msg.LctId,
		EncryptedMessage: encryptedMessage,
		Status:           "encrypted",
	}, nil
}

// GenerateLCTChallenge implements the Msg/GenerateLCTChallenge RPC method.
func (ms msgServer) GenerateLCTChallenge(ctx context.Context, msg *types.MsgGenerateLCTChallenge) (*types.MsgGenerateLCTChallengeResponse, error) {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.LctId == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "LCT ID cannot be empty")
	}

	// Generate cryptographic challenge
	challenge, err := ms.Keeper.GenerateLCTChallenge(ctx, msg.LctId)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to generate challenge: %s", err)
	}

	// Emit event for audit
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_challenge_generated",
			sdk.NewAttribute("lct_id", msg.LctId),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("challenge_length", fmt.Sprintf("%d", len(challenge))),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	return &types.MsgGenerateLCTChallengeResponse{
		LctId:     msg.LctId,
		Challenge: challenge,
		Status:    "generated",
	}, nil
}

// VerifyLCTChallenge implements the Msg/VerifyLCTChallenge RPC method.
func (ms msgServer) VerifyLCTChallenge(ctx context.Context, msg *types.MsgVerifyLCTChallenge) (*types.MsgVerifyLCTChallengeResponse, error) {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.LctId == "" {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "LCT ID cannot be empty")
	}

	if len(msg.Challenge) == 0 {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "challenge cannot be empty")
	}

	if len(msg.Response) == 0 {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "response cannot be empty")
	}

	// Verify the challenge response
	verified, err := ms.VerifyLCTChallengeResponse(ctx, msg.LctId, msg.Challenge, msg.Response)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "failed to verify challenge: %s", err)
	}

	// Emit event for audit
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("lct_challenge_verified",
			sdk.NewAttribute("lct_id", msg.LctId),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("verified", fmt.Sprintf("%t", verified)),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", time.Now().Unix())),
		),
	)

	return &types.MsgVerifyLCTChallengeResponse{
		LctId:    msg.LctId,
		Verified: verified,
		Status:   "verified",
	}, nil
}
