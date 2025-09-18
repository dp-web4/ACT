package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"racecar-web/x/energycycle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams implements the Msg/UpdateParams message type.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if err := ms.Params.Set(ctx, msg.Params); err != nil {
		return nil, errorsmod.Wrap(err, "failed to set params")
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// CreateRelationshipEnergyOperation implements the Msg/CreateRelationshipEnergyOperation message type.
func (k msgServer) CreateRelationshipEnergyOperation(ctx context.Context, msg *types.MsgCreateRelationshipEnergyOperation) (*types.MsgCreateRelationshipEnergyOperationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	invalidInputErr := errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid input")

	// Validate input
	if msg.SourceLct == "" || msg.TargetLct == "" {
		return nil, errorsmod.Wrap(invalidInputErr, "both source and target LCTs must be provided")
	}

	if msg.EnergyAmount == "" {
		return nil, errorsmod.Wrap(invalidInputErr, "energy amount must be provided")
	}

	// Validate operation type
	switch msg.OperationType {
	case types.OperationTypeDischarge, types.OperationTypeCharge, types.OperationTypeTransfer, types.OperationTypeBalance:
		// Valid operation types
	default:
		return nil, errorsmod.Wrap(invalidInputErr, "invalid operation type")
	}

	// Get current block height and timestamp
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	timestamp := time.Now().Unix()

	// Generate a unique operation ID
	operationId := fmt.Sprintf("op-%s-%s-%s-%d", msg.SourceLct, msg.TargetLct, msg.OperationType, timestamp)

	// Validate the energy operation using ATP/ADP logic
	isValid, validationMsg, err := k.ValidateEnergyOperation(ctx, operationId, msg.SourceLct, msg.TargetLct, msg.EnergyAmount, msg.OperationType)
	if err != nil {
		return nil, errorsmod.Wrap(err, "energy operation validation failed")
	}

	if !isValid {
		return nil, errorsmod.Wrap(invalidInputErr, validationMsg)
	}

	// Create ATP token for energy storage
	var atpTokenID string
	if msg.OperationType == types.OperationTypeCharge {
		atpToken, err := k.CreateAtpToken(ctx, msg.SourceLct, msg.EnergyAmount, operationId, "energy_operation", blockHeight)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to create ATP token")
		}
		atpTokenID = atpToken.TokenId
	}

	// Get trust score
	trustScore := "0.5" // Default
	if k.trusttensorKeeper != nil {
		trustScore, _, err = k.trusttensorKeeper.CalculateRelationshipTrust(ctx, msg.SourceLct, "energy_operation")
		if err != nil {
			// Use default trust score
		}
	}

	// Create energy operation
	operation := &types.EnergyOperation{
		OperationId:      operationId,
		SourceLct:        msg.SourceLct,
		TargetLct:        msg.TargetLct,
		EnergyAmount:     msg.EnergyAmount,
		OperationType:    msg.OperationType,
		Status:           types.StatusCreated,
		Timestamp:        timestamp,
		BlockHeight:      blockHeight,
		TrustScore:       trustScore,
		AtpTokenId:       atpTokenID,
		EnergyEfficiency: "0.8", // Default efficiency
		ValidationData:   validationMsg,
		Version:          1,
	}

	// Store the operation
	err = k.EnergyOperations.Set(ctx, operationId, *operation)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to store energy operation")
	}

	// Return response
	return &types.MsgCreateRelationshipEnergyOperationResponse{
		OperationId:    operationId,
		TrustValidated: trustScore != "0.5",
	}, nil
}

// ExecuteEnergyTransfer implements the Msg/ExecuteEnergyTransfer message type.
func (k msgServer) ExecuteEnergyTransfer(ctx context.Context, msg *types.MsgExecuteEnergyTransfer) (*types.MsgExecuteEnergyTransferResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.OperationId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "operation ID must be provided")
	}

	// Get the energy operation
	operation, err := k.EnergyOperations.Get(ctx, msg.OperationId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "energy operation not found")
	}

	// Check if operation is in correct state
	if operation.Status != types.StatusCreated {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "operation is not in created state")
	}

	// Get current block height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// Execute based on operation type
	switch operation.OperationType {
	case types.OperationTypeDischarge:
		// Discharge ATP token to create ADP token
		if operation.AtpTokenId == "" {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no ATP token associated with discharge operation")
		}

		adpToken, err := k.DischargeAtpToken(ctx, operation.AtpTokenId, operation.OperationId, blockHeight)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to discharge ATP token")
		}

		// Update operation with ADP token ID
		operation.AdpTokenId = adpToken.TokenId
		operation.Status = types.StatusCompleted
		operation.EnergyEfficiency = adpToken.EnergyEfficiency
		operation.Version++

	case types.OperationTypeTransfer:
		// For transfers, we need to validate energy balance and create new ATP tokens
		// This is a simplified implementation
		operation.Status = types.StatusCompleted
		operation.Version++

	default:
		// For other operation types, just mark as completed
		operation.Status = types.StatusCompleted
		operation.Version++
	}

	// Update the operation
	err = k.EnergyOperations.Set(ctx, operation.OperationId, operation)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to update energy operation")
	}

	return &types.MsgExecuteEnergyTransferResponse{}, nil
}

// ValidateRelationshipValue implements the Msg/ValidateRelationshipValue message type.
func (k msgServer) ValidateRelationshipValue(ctx context.Context, msg *types.MsgValidateRelationshipValue) (*types.MsgValidateRelationshipValueResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	invalidInputErr := errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid input")

	if msg.OperationId == "" {
		return nil, errorsmod.Wrap(invalidInputErr, "operation ID must be provided")
	}

	// Get the energy operation
	operation, err := k.EnergyOperations.Get(ctx, msg.OperationId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "energy operation not found")
	}

	// Calculate V3 composite score using trust tensor
	v3Score, err := k.trusttensorKeeper.CalculateV3CompositeScore(ctx, msg.OperationId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to calculate V3 score")
	}

	// Get current block height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// If this is a discharge operation and we have an ATP token, create ADP token
	var adpTokenIDs []string
	if operation.OperationType == types.OperationTypeDischarge && operation.AtpTokenId != "" {
		adpToken, err := k.DischargeAtpToken(ctx, operation.AtpTokenId, msg.OperationId, blockHeight)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to create ADP token")
		}
		adpTokenIDs = append(adpTokenIDs, adpToken.TokenId)

		// Update operation with ADP token ID
		operation.AdpTokenId = adpToken.TokenId
		operation.Status = types.StatusValidated
		operation.Version++
		k.EnergyOperations.Set(ctx, operation.OperationId, operation)
	}

	// Convert ADP token IDs to string
	adpTokensStr := ""
	if len(adpTokenIDs) > 0 {
		adpTokensStr = adpTokenIDs[0] // For simplicity, just use the first one
	}

	return &types.MsgValidateRelationshipValueResponse{
		V_3Score:  v3Score.String(),
		AdpTokens: adpTokensStr,
	}, nil
}

// DischargeATP implements the Msg/DischargeATP RPC method for converting ATP to ADP during work.
func (ms msgServer) DischargeATP(ctx context.Context, msg *types.MsgDischargeATP) (*types.MsgDischargeATPResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.LctId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "worker LCT ID cannot be empty")
	}

	if msg.Amount == "" || msg.Amount == "0" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "ATP amount must be greater than 0")
	}

	if msg.WorkDescription == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "work description required for discharge")
	}

	// Get current block height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// Generate work ID for tracking
	workId := fmt.Sprintf("work-%s-%d", msg.LctId, blockHeight)

	// TODO: Verify society has sufficient ATP tokens in pool
	// For now, we assume ATP is available

	// Create ADP token representing discharged energy
	adpToken := &types.RelationshipAdpToken{
		TokenId:          fmt.Sprintf("adp-%s", workId),
		OriginalAtpId:    fmt.Sprintf("atp-pool-%s", msg.Amount), // Reference to pool ATP
		LctId:            msg.LctId,
		DischargedAt:     blockHeight,
		ValueScore:       "0", // Will be set by V3 tensor validation
		ConfirmationData: msg.WorkDescription,
		EnergyEfficiency: "1.0", // Will be calculated from actual vs planned
		TrustValidation:  "pending",
		ValidationBlock:  blockHeight + 10, // Validation window
		OperationContext: fmt.Sprintf("work:%s:target:%s", msg.WorkDescription, msg.TargetLct),
		Version:          1,
	}

	// Store the ADP token
	if err := ms.RelationshipAdpTokens.Set(ctx, adpToken.TokenId, *adpToken); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store ADP token")
	}

	// Calculate energy released for work
	energyReleased := msg.Amount // Energy made available for work
	adpCreated := msg.Amount     // Same amount of ADP created (state change only)

	// If trusttensor keeper is available, initiate V3 tracking
	if ms.trusttensorKeeper != nil && msg.TargetLct != "" {
		// Track work relationship between LCTs for V3 tensor calculation
		_, _, err = ms.trusttensorKeeper.CalculateRelationshipTrust(ctx, msg.LctId, workId)
		if err != nil {
			// Log but don't fail - V3 tracking is async
			// TODO: Add logger when available
		}
	}

	// Emit event for tracking energy expenditure
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("atp_discharged_for_work",
			sdk.NewAttribute("worker_lct", msg.LctId),
			sdk.NewAttribute("work_id", workId),
			sdk.NewAttribute("work_description", msg.WorkDescription),
			sdk.NewAttribute("target_lct", msg.TargetLct),
			sdk.NewAttribute("atp_consumed", msg.Amount),
			sdk.NewAttribute("adp_created", adpCreated),
			sdk.NewAttribute("adp_token_id", adpToken.TokenId),
			sdk.NewAttribute("energy_released", energyReleased),
			sdk.NewAttribute("validation_window", fmt.Sprintf("%d", adpToken.ValidationBlock)),
			sdk.NewAttribute("creator", creator.String()),
		),
	)

	// TODO: Update society token pool balances
	// Society ATP balance -= msg.Amount
	// Society ADP balance += msg.Amount
	remainingAtp := "0" // Placeholder - should query society pool

	return &types.MsgDischargeATPResponse{
		EnergyReleased: energyReleased,
		AdpCreated:     adpCreated,
		WorkId:         workId,
		RemainingAtp:   remainingAtp,
	}, nil
}

// RechargeADP implements the Msg/RechargeADP RPC method for energy producers to charge ADP tokens to ATP.
func (ms msgServer) RechargeADP(ctx context.Context, msg *types.MsgRechargeADP) (*types.MsgRechargeADPResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.LctId == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "producer LCT ID cannot be empty")
	}

	if msg.Amount == "" || msg.Amount == "0" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than 0")
	}

	// Validate energy source (must be a producer entity type)
	validSources := map[string]bool{
		"solar": true, "wind": true, "wave": true, "nuclear": true,
		"geothermal": true, "grid": true, "battery": true,
	}
	if !validSources[msg.EnergySource] {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid energy source: %s", msg.EnergySource)
	}

	// Get current block height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()

	// Generate recharge operation ID
	rechargeId := fmt.Sprintf("charge-%s-%s-%d", msg.LctId, msg.EnergySource, blockHeight)

	// In Web4, we're converting society's ADP tokens to ATP by adding real energy
	// This requires validation proof that actual energy was harvested/generated
	if msg.ValidationProof == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "energy generation validation proof required")
	}

	// Create ATP token representing the charged energy
	atpToken := &types.RelationshipAtpToken{
		TokenId:             fmt.Sprintf("atp-%s", rechargeId),
		LctId:               msg.LctId,
		EnergyAmount:        msg.Amount,
		CreatedAt:           blockHeight,
		OperationId:         rechargeId,
		Status:              "charged",
		RelationshipContext: fmt.Sprintf("producer:%s:source:%s", msg.LctId, msg.EnergySource),
		ExpirationBlock:     blockHeight + 100000, // ATP tokens expire if not used
		TrustScore:          "1.0", // Producer entities have high trust
		EfficiencyRating:    msg.ValidationProof, // Store validation in efficiency field
		Version:             1,
	}

	// Store the ATP token
	if err := ms.RelationshipAtpTokens.Set(ctx, atpToken.TokenId, *atpToken); err != nil {
		return nil, errorsmod.Wrap(err, "failed to store ATP token")
	}

	// Track energy consumed from environment (cost of charging)
	energyConsumed := msg.Amount // Real energy harvested

	// Emit event for society-level tracking
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("adp_charged_to_atp",
			sdk.NewAttribute("producer_lct", msg.LctId),
			sdk.NewAttribute("energy_source", msg.EnergySource),
			sdk.NewAttribute("adp_consumed", msg.Amount),
			sdk.NewAttribute("atp_created", atpToken.EnergyAmount),
			sdk.NewAttribute("atp_token_id", atpToken.TokenId),
			sdk.NewAttribute("energy_harvested", energyConsumed),
			sdk.NewAttribute("validation_proof", msg.ValidationProof),
			sdk.NewAttribute("creator", creator.String()),
		),
	)

	// TODO: Update society token pool balances
	// Society ADP balance -= msg.Amount
	// Society ATP balance += msg.Amount
	remainingAdp := "0"  // Placeholder - should query society pool
	newAtpBalance := msg.Amount // Placeholder - should query society pool

	return &types.MsgRechargeADPResponse{
		AtpCreated:     atpToken.EnergyAmount,
		EnergyConsumed: energyConsumed,
		RemainingAdp:   remainingAdp,
		NewAtpBalance:  newAtpBalance,
	}, nil
}

// MintADP implements the Msg/MintADP RPC method for treasury role to mint new ADP tokens.
func (ms msgServer) MintADP(ctx context.Context, msg *types.MsgMintADP) (*types.MsgMintADPResponse, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Validate input
	if msg.Amount == "" || msg.Amount == "0" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than 0")
	}

	if msg.SocietyLct == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "society LCT ID cannot be empty")
	}

	if msg.RoleLct == "" {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "treasury role LCT ID cannot be empty")
	}

	// TODO: Verify that role_lct is actually a treasury role
	// TODO: Verify that role_lct belongs to the society_lct
	// For now, we'll trust the input during genesis

	// Get current block height
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockHeight := sdkCtx.BlockHeight()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Generate mint operation ID
	mintId := fmt.Sprintf("mint-adp-%s-%d", msg.SocietyLct, blockHeight)

	// Create ADP tokens for the society treasury
	// In Web4, these represent discharged energy available to be charged
	// The treasury role has the authority to mint initial ADP allocation

	// Store the mint operation for audit trail
	// TODO: Create proper mint record storage

	// For now, we track the minted amount conceptually
	// In production, this would update the society's ADP balance in state

	// Emit event for mint operation
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("adp_minted",
			sdk.NewAttribute("society_lct", msg.SocietyLct),
			sdk.NewAttribute("treasury_role", msg.RoleLct),
			sdk.NewAttribute("amount", msg.Amount),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("mint_id", mintId),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", blockHeight)),
			sdk.NewAttribute("timestamp", timestamp),
			sdk.NewAttribute("creator", creator.String()),
		),
	)

	// TODO: Actually update society ADP balance in state
	// For now, we'll return success with the minted amount
	societyBalance := msg.Amount // This should be queried from state after update

	return &types.MsgMintADPResponse{
		MintedAmount:   msg.Amount,
		SocietyBalance: societyBalance,
		MintId:         mintId,
		Timestamp:      timestamp,
	}, nil
}
