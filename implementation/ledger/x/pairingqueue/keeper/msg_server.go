package keeper

import (
	"context"

	"cosmossdk.io/errors"

	"racecar-web/x/pairingqueue/types"
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

// UpdateParams implements the Msg/UpdateParams message type.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if err := ms.Params.Set(ctx, msg.Params); err != nil {
		return nil, errors.Wrap(err, "failed to set params")
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// QueuePairingRequest implements the Msg/QueuePairingRequest message type.
func (ms msgServer) QueuePairingRequest(ctx context.Context, msg *types.MsgQueuePairingRequest) (*types.MsgQueuePairingRequestResponse, error) {
	requestID, err := ms.Keeper.QueuePairingRequest(ctx, msg.InitiatorId, msg.TargetId, msg.RequestType, msg.ProxyId)
	if err != nil {
		return nil, err
	}

	return &types.MsgQueuePairingRequestResponse{
		RequestId: requestID,
		Status:    "queued",
	}, nil
}

// ProcessOfflineQueue implements the Msg/ProcessOfflineQueue message type.
func (ms msgServer) ProcessOfflineQueue(ctx context.Context, msg *types.MsgProcessOfflineQueue) (*types.MsgProcessOfflineQueueResponse, error) {
	processed, failed, err := ms.Keeper.ProcessOfflineQueue(ctx, msg.ComponentId)
	if err != nil {
		return nil, err
	}

	return &types.MsgProcessOfflineQueueResponse{
		ProcessedCount: int64(processed),
		FailedCount:    int64(failed),
	}, nil
}

// CancelRequest implements the Msg/CancelRequest message type.
func (ms msgServer) CancelRequest(ctx context.Context, msg *types.MsgCancelRequest) (*types.MsgCancelRequestResponse, error) {
	if err := ms.Keeper.CancelRequest(ctx, msg.RequestId, msg.Reason); err != nil {
		return nil, err
	}

	return &types.MsgCancelRequestResponse{}, nil
}
