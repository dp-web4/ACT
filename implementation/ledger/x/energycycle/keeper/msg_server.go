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
