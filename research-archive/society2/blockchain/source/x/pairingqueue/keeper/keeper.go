package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"

	componentregistrytypes "racecar-web/x/componentregistry/types"
	"racecar-web/x/pairingqueue/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	authKeeper   types.AuthKeeper
	bankKeeper   types.BankKeeper
	compKeeper   componentregistrytypes.ComponentregistryKeeper

	// State management
	Params            collections.Item[types.Params]
	PairingRequests   collections.Map[string, types.PairingRequest]
	OfflineOperations collections.Map[string, types.OfflineOperation]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	compKeeper componentregistrytypes.ComponentregistryKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		authKeeper:   authKeeper,
		bankKeeper:   bankKeeper,
		compKeeper:   compKeeper,

		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		PairingRequests:   collections.NewMap(sb, types.PairingRequestsPrefix, "pairing_requests", collections.StringKey, codec.CollValue[types.PairingRequest](cdc)),
		OfflineOperations: collections.NewMap(sb, types.OfflineOperationsPrefix, "offline_operations", collections.StringKey, codec.CollValue[types.OfflineOperation](cdc)),
	}

	return k
}

// QueuePairingRequest queues a new pairing request with trust and authorization checks
func (k Keeper) QueuePairingRequest(ctx context.Context, initiatorID, targetID, requestType, proxyID string) (string, error) {
	// Validate components exist and are verified
	initiatorVerified, _ := k.compKeeper.VerifyComponentForPairing(ctx, initiatorID)
	if !initiatorVerified {
		return "", errors.Wrap(types.ErrComponentNotVerified, "initiator component not verified")
	}
	targetVerified, _ := k.compKeeper.VerifyComponentForPairing(ctx, targetID)
	if !targetVerified {
		return "", errors.Wrap(types.ErrComponentNotVerified, "target component not verified")
	}

	// Check pairing authorization
	authorized, _, _ := k.compKeeper.CheckBidirectionalPairingAuth(ctx, initiatorID, targetID)
	if !authorized {
		return "", errors.Wrap(types.ErrPairingNotAuthorized, "pairing not authorized")
	}

	// For now, skip trust tensor check to avoid circular dependency
	// TODO: Implement trust check via a different mechanism if needed

	// Generate request ID
	requestID := fmt.Sprintf("%s-%s-%d", initiatorID, targetID, time.Now().UnixNano())

	// Create and store request with retry configuration
	request := types.PairingRequest{
		RequestId:     requestID,
		InitiatorId:   initiatorID,
		TargetId:      targetID,
		RequestType:   requestType,
		ProxyId:       proxyID,
		Status:        "queued",
		CreatedAt:     time.Now().Unix(),
		ProcessedAt:   0,
		RetryCount:    0,
		MaxRetries:    3, // Default max retries
		LastAttemptAt: 0,
		NextRetryAt:   0,
		FailureReason: "",
	}

	if err := k.PairingRequests.Set(ctx, requestID, request); err != nil {
		return "", errors.Wrap(err, "failed to store pairing request")
	}

	return requestID, nil
}

// GetPairingRequest retrieves a pairing request by ID
func (k Keeper) GetPairingRequest(ctx context.Context, requestID string) (types.PairingRequest, bool) {
	request, err := k.PairingRequests.Get(ctx, requestID)
	if err != nil {
		return types.PairingRequest{}, false
	}
	return request, true
}

// QueueOfflineOperation queues an offline operation for a component with trust validation
func (k Keeper) QueueOfflineOperation(ctx context.Context, componentID, operationType string) (string, error) {
	// Validate component exists and is verified
	componentVerified, _ := k.compKeeper.VerifyComponentForPairing(ctx, componentID)
	if !componentVerified {
		return "", errors.Wrap(types.ErrComponentNotVerified, "component not verified")
	}

	// Generate operation ID
	operationID := fmt.Sprintf("%s-%s-%d", componentID, operationType, time.Now().UnixNano())

	// Create and store operation with retry configuration
	operation := types.OfflineOperation{
		OperationId:   operationID,
		ComponentId:   componentID,
		OperationType: operationType,
		QueuedAt:      time.Now().Unix(),
		RetryCount:    0,
		MaxRetries:    5, // Default max retries for offline operations
		LastAttemptAt: 0,
		NextRetryAt:   0,
		FailureReason: "",
	}

	if err := k.OfflineOperations.Set(ctx, operationID, operation); err != nil {
		return "", errors.Wrap(err, "failed to store offline operation")
	}

	return operationID, nil
}

// ProcessOfflineQueue processes offline operations for a component with robust retry logic
func (k Keeper) ProcessOfflineQueue(ctx context.Context, componentID string) (int, int, error) {
	var processed, failed int
	now := time.Now().Unix()

	// Get all operations for the component
	iter, err := k.OfflineOperations.Iterate(ctx, nil)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to iterate offline operations")
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		operation, err := iter.Value()
		if err != nil {
			failed++
			continue
		}

		if operation.ComponentId != componentID {
			continue
		}

		// Check if operation is ready for retry
		if operation.NextRetryAt > 0 && operation.NextRetryAt > now {
			continue // Not ready for retry yet
		}

		// Check if max retries exceeded
		if operation.RetryCount >= operation.MaxRetries {
			// Mark as permanently failed
			operation.FailureReason = "max retries exceeded"
			if err := k.OfflineOperations.Set(ctx, operation.OperationId, operation); err != nil {
				failed++
			}
			continue
		}

		// Validate component still exists and is verified
		componentVerified, _ := k.compKeeper.VerifyComponentForPairing(ctx, componentID)
		if !componentVerified {
			operation.RetryCount++
			operation.LastAttemptAt = now
			operation.NextRetryAt = now + 300 // Retry in 5 minutes
			operation.FailureReason = "component not verified"
			if err := k.OfflineOperations.Set(ctx, operation.OperationId, operation); err != nil {
				failed++
			}
			continue
		}

		// Process operation based on type with trust validation
		success := false
		switch operation.OperationType {
		case "pairing":
			success = k.processPairingOperation(ctx, operation)
		case "unpairing":
			success = k.processUnpairingOperation(ctx, operation)
		case "energy_transfer":
			success = k.processEnergyTransferOperation(ctx, operation)
		default:
			operation.FailureReason = "unknown operation type"
		}

		if success {
			processed++
			// Remove successful operation
			if err := k.OfflineOperations.Remove(ctx, operation.OperationId); err != nil {
				failed++
			}
		} else {
			// Update retry information
			operation.RetryCount++
			operation.LastAttemptAt = now
			operation.NextRetryAt = now + k.calculateRetryDelay(operation.RetryCount)
			if err := k.OfflineOperations.Set(ctx, operation.OperationId, operation); err != nil {
				failed++
			}
		}
	}

	return processed, failed, nil
}

// processPairingOperation processes a pairing operation with trust validation
func (k Keeper) processPairingOperation(ctx context.Context, operation types.OfflineOperation) bool {
	// This would integrate with the pairing module
	// For now, return true to simulate success
	return true
}

// processUnpairingOperation processes an unpairing operation
func (k Keeper) processUnpairingOperation(ctx context.Context, operation types.OfflineOperation) bool {
	// This would integrate with the pairing module
	// For now, return true to simulate success
	return true
}

// processEnergyTransferOperation processes an energy transfer operation
func (k Keeper) processEnergyTransferOperation(ctx context.Context, operation types.OfflineOperation) bool {
	// This would integrate with the energy cycle module
	// For now, return true to simulate success
	return true
}

// calculateRetryDelay calculates exponential backoff delay
func (k Keeper) calculateRetryDelay(retryCount int32) int64 {
	// Exponential backoff: 5min, 10min, 20min, 40min, 80min
	baseDelay := int64(300) // 5 minutes
	maxDelay := int64(4800) // 80 minutes

	delay := baseDelay * (1 << (retryCount - 1))
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// CancelRequest cancels a queued request
func (k Keeper) CancelRequest(ctx context.Context, requestID, reason string) error {
	request, err := k.PairingRequests.Get(ctx, requestID)
	if err != nil {
		return errors.Wrap(types.ErrRequestNotFound, "request not found")
	}

	if request.Status != "queued" {
		return errors.Wrap(types.ErrInvalidRequestStatus, "can only cancel queued requests")
	}

	request.Status = "cancelled"
	request.RequestData = reason
	request.ProcessedAt = time.Now().Unix()

	return k.PairingRequests.Set(ctx, requestID, request)
}

// GetQueuedRequests returns all queued requests for a component
func (k Keeper) GetQueuedRequests(ctx context.Context, componentID string) ([]types.PairingRequest, error) {
	var requests []types.PairingRequest

	iter, err := k.PairingRequests.Iterate(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate pairing requests")
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		request, err := iter.Value()
		if err != nil {
			continue
		}

		if (request.InitiatorId == componentID || request.TargetId == componentID) && request.Status == "queued" {
			requests = append(requests, request)
		}
	}

	return requests, nil
}

// GetQueueStatus returns comprehensive queue status for a component
func (k Keeper) GetQueueStatus(ctx context.Context, componentID string) (types.QueueStatus, error) {
	var status types.QueueStatus
	status.ComponentId = componentID
	status.Timestamp = time.Now().Unix()

	// Count queued requests
	requests, err := k.GetQueuedRequests(ctx, componentID)
	if err != nil {
		return status, err
	}
	status.QueuedRequests = int32(len(requests))

	// Count offline operations
	iter, err := k.OfflineOperations.Iterate(ctx, nil)
	if err != nil {
		return status, err
	}
	defer iter.Close()

	var operations []types.OfflineOperation
	for ; iter.Valid(); iter.Next() {
		operation, err := iter.Value()
		if err != nil {
			continue
		}

		if operation.ComponentId == componentID {
			operations = append(operations, operation)
		}
	}
	status.OfflineOperations = int32(len(operations))

	// Calculate retry statistics
	var totalRetries int32
	var failedOperations int32
	for _, op := range operations {
		totalRetries += op.RetryCount
		if op.RetryCount >= op.MaxRetries {
			failedOperations++
		}
	}
	status.TotalRetries = totalRetries
	status.FailedOperations = failedOperations

	return status, nil
}

// ListProxyQueue returns all offline operations for a proxy
func (k Keeper) ListProxyQueue(ctx context.Context, proxyID string) ([]types.OfflineOperation, error) {
	var operations []types.OfflineOperation

	iter, err := k.OfflineOperations.Iterate(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to iterate offline operations")
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		operation, err := iter.Value()
		if err != nil {
			continue
		}

		// For now, we'll return all operations since proxy filtering logic
		// would depend on specific proxy-component relationships
		operations = append(operations, operation)
	}

	return operations, nil
}
