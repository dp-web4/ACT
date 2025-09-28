# Pairing Queue Module

## Table of Contents

1. [Overview](#overview)
2. [Key Concepts](#key-concepts)
3. [State Management](#state-management)
4. [Messages (Transactions)](#messages-transactions)
5. [Queries](#queries)
6. [Events](#events)
7. [Parameters](#parameters)
8. [Queue Processing](#queue-processing)
9. [Integration Guide](#integration-guide)

## Overview

The Pairing Queue module handles offline pairing scenarios where one or both components are not immediately available for pairing. It maintains queues of pending operations and processes them when components come online, ensuring reliable pairing even in disconnected environments common in battery management systems.

### Purpose
- Queue pairing requests when target components are offline
- Process queued operations when components become available
- Support proxy-based pairing for permanently offline devices
- Maintain request status and retry mechanisms

### Dependencies
- **Component Registry**: Verifies component existence and status
- **Pairing**: Executes actual pairing operations

### Module Store Key
`pairingqueue`

## Key Concepts

### Pairing Request
A queued request to establish pairing between components:
- **Request ID**: Unique identifier for the request
- **Target Component**: Component that needs to come online
- **Initiator**: Component requesting the pairing
- **Request Type**: STANDARD, PROXY, EMERGENCY
- **Priority**: Processing priority level
- **Expiry**: Maximum time to keep in queue

### Offline Operation
Metadata about operations pending for offline components:
- **Component ID**: Target component
- **Operation Count**: Number of pending operations
- **Last Seen**: When component was last online
- **Expected Return**: Estimated return time (if known)

### Queue Types
- **Standard Queue**: Normal priority operations
- **Priority Queue**: High priority operations (e.g., safety-critical)
- **Proxy Queue**: Operations that can be completed via proxy

## State Management

### Stored Types

#### 1. PairingRequest
```protobuf
message PairingRequest {
  string request_id = 1;
  string initiator_component = 2;
  string target_component = 3;
  string request_type = 4;  // STANDARD, PROXY, EMERGENCY
  string priority = 5;      // LOW, NORMAL, HIGH, CRITICAL
  string status = 6;        // QUEUED, PROCESSING, COMPLETED, FAILED, EXPIRED
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp expires_at = 8;
  uint32 retry_count = 9;
  string proxy_component = 10;  // For proxy requests
  map<string, string> metadata = 11;
}
```

#### 2. OfflineOperation
```protobuf
message OfflineOperation {
  string component_id = 1;
  repeated string pending_request_ids = 2;
  uint32 operation_count = 3;
  google.protobuf.Timestamp last_seen = 4;
  google.protobuf.Timestamp expected_return = 5;
  string offline_reason = 6;
  bool accepts_proxy = 7;
}
```

#### 3. QueueStatus
```protobuf
message QueueStatus {
  string component_id = 1;
  uint32 standard_queue_size = 2;
  uint32 priority_queue_size = 3;
  uint32 proxy_queue_size = 4;
  google.protobuf.Timestamp oldest_request = 5;
  google.protobuf.Timestamp last_processed = 6;
}
```

### Store Layout
```
pairingqueue/
├── requests/
│   └── {request_id} → PairingRequest
├── offline/
│   └── {component_id} → OfflineOperation
├── queues/
│   └── {component_id}/
│       ├── standard/
│       │   └── {priority}/{timestamp} → request_id
│       ├── priority/
│       │   └── {priority}/{timestamp} → request_id
│       └── proxy/
│           └── {proxy_id} → request_id
└── params → Params
```

## Messages (Transactions)

### 1. QueuePairingRequest
Queues a pairing request for an offline component.

**Input**:
```protobuf
message MsgQueuePairingRequest {
  string creator = 1;
  string initiator_component = 2;
  string target_component = 3;
  string request_type = 4;
  string priority = 5;
  uint32 expiry_hours = 6;
  string proxy_component = 7;  // Optional
  map<string, string> metadata = 8;
}
```

**Validation**:
- Initiator component must exist and be online
- Target component must exist
- If proxy specified, proxy component must exist and be authorized
- Expiry within allowed range
- Creator must own initiator component

**Effects**:
- Creates PairingRequest record
- Adds to appropriate queue based on type and priority
- Updates OfflineOperation for target component
- Emits `request_queued` event

**Example**:
```bash
racecar-webd tx pairingqueue queue-pairing-request \
  --initiator="comp_pack1" \
  --target="comp_module_offline" \
  --request-type="STANDARD" \
  --priority="NORMAL" \
  --expiry-hours=72 \
  --from mykey
```

### 2. ProcessOfflineQueue
Processes all queued requests for a component when it comes online.

**Input**:
```protobuf
message MsgProcessOfflineQueue {
  string creator = 1;
  string component_id = 2;
  bool process_proxy_queue = 3;
  uint32 max_operations = 4;  // Optional limit
}
```

**Validation**:
- Component must exist
- Creator must have processing rights (component owner or validator)
- Component should have pending operations

**Effects**:
- Processes queued requests in priority order
- Initiates actual pairing operations via Pairing module
- Updates request statuses
- Removes completed/failed requests from queue
- Emits `queue_processed` event

**Processing Order**:
1. Critical priority requests
2. High priority requests  
3. Normal priority requests
4. Low priority requests
5. Proxy requests (if enabled)

### 3. CancelRequest
Cancels a queued pairing request.

**Input**:
```protobuf
message MsgCancelRequest {
  string creator = 1;
  string request_id = 2;
  string cancellation_reason = 3;
}
```

**Validation**:
- Request must exist and be in QUEUED status
- Creator must own the initiator component or be admin
- Must provide cancellation reason

**Effects**:
- Updates request status to CANCELLED
- Removes from processing queue
- Updates queue statistics
- Emits `request_cancelled` event

## Queries

### 1. GetQueuedRequests
Gets all queued requests for a component.

**Request**:
```protobuf
message QueryGetQueuedRequestsRequest {
  string component_id = 1;
  string queue_type = 2;    // STANDARD, PRIORITY, PROXY, ALL
  string status_filter = 3; // Optional: QUEUED, PROCESSING, etc.
  cosmos.base.query.v1beta1.PageRequest pagination = 4;
}
```

**Response**:
```protobuf
message QueryGetQueuedRequestsResponse {
  repeated PairingRequest requests = 1;
  QueueStatus queue_status = 2;
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

### 2. GetRequestStatus
Gets the status of a specific pairing request.

**Request**:
```protobuf
message QueryGetRequestStatusRequest {
  string request_id = 1;
}
```

**Response**:
```protobuf
message QueryGetRequestStatusResponse {
  PairingRequest request = 1;
  uint32 queue_position = 2;
  google.protobuf.Timestamp estimated_processing = 3;
}
```

### 3. ListProxyQueue
Lists requests that can be processed via proxy.

**Request**:
```protobuf
message QueryListProxyQueueRequest {
  string proxy_component = 1;
  cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
```

**Response**:
```protobuf
message QueryListProxyQueueResponse {
  repeated PairingRequest proxy_requests = 1;
  repeated string authorized_targets = 2;
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

## Events

### request_queued
Emitted when a pairing request is queued.
```json
{
  "type": "request_queued",
  "attributes": [
    {"key": "request_id", "value": "req_abc123"},
    {"key": "initiator", "value": "comp_pack1"},
    {"key": "target", "value": "comp_module_offline"},
    {"key": "priority", "value": "NORMAL"},
    {"key": "expires_at", "value": "2024-01-18T10:30:00Z"}
  ]
}
```

### queue_processed
Emitted when offline queue is processed.
```json
{
  "type": "queue_processed",
  "attributes": [
    {"key": "component_id", "value": "comp_module_offline"},
    {"key": "requests_processed", "value": "5"},
    {"key": "successful", "value": "4"},
    {"key": "failed", "value": "1"},
    {"key": "processing_time_ms", "value": "1250"}
  ]
}
```

### request_cancelled
Emitted when a request is cancelled.
```json
{
  "type": "request_cancelled",
  "attributes": [
    {"key": "request_id", "value": "req_abc123"},
    {"key": "cancelled_by", "value": "cosmos1..."},
    {"key": "reason", "value": "Target component decommissioned"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Maximum time to keep requests in queue (hours)
  uint32 max_queue_time_hours = 1;
  
  // Maximum requests per component in queue
  uint32 max_requests_per_component = 2;
  
  // Maximum operations to process in single batch
  uint32 max_batch_processing_size = 3;
  
  // Enable proxy-based processing
  bool enable_proxy_processing = 4;
  
  // Auto-cancel expired requests
  bool auto_cancel_expired = 5;
  
  // Retry failed requests
  bool enable_retry = 6;
  
  // Maximum retry attempts
  uint32 max_retry_attempts = 7;
}
```

**Default Values**:
```json
{
  "max_queue_time_hours": 168,    // 1 week
  "max_requests_per_component": 100,
  "max_batch_processing_size": 10,
  "enable_proxy_processing": true,
  "auto_cancel_expired": true,
  "enable_retry": true,
  "max_retry_attempts": 3
}
```

## Queue Processing

### Processing Algorithm

```go
func ProcessQueue(componentID string) error {
    // Get all queued requests for component
    requests := getQueuedRequests(componentID)
    
    // Sort by priority and timestamp
    sort.Slice(requests, func(i, j int) bool {
        if requests[i].Priority != requests[j].Priority {
            return priorityValue(requests[i].Priority) > priorityValue(requests[j].Priority)
        }
        return requests[i].CreatedAt.Before(requests[j].CreatedAt)
    })
    
    // Process each request
    for _, req := range requests {
        if req.Status != "QUEUED" {
            continue
        }
        
        // Update status to processing
        updateRequestStatus(req.RequestId, "PROCESSING")
        
        // Attempt pairing
        err := pairingKeeper.InitiatePairing(req.InitiatorComponent, req.TargetComponent)
        if err != nil {
            handleFailure(req, err)
        } else {
            updateRequestStatus(req.RequestId, "COMPLETED")
        }
    }
    
    return nil
}
```

### Priority System

| Priority | Value | Use Case |
|----------|-------|----------|
| CRITICAL | 4 | Safety-critical pairings |
| HIGH | 3 | Operational requirements |
| NORMAL | 2 | Standard pairing requests |
| LOW | 1 | Maintenance or optional pairings |

### Retry Logic

```go
func handleFailure(req PairingRequest, err error) {
    if req.RetryCount < maxRetryAttempts {
        // Exponential backoff
        delay := time.Duration(math.Pow(2, float64(req.RetryCount))) * time.Minute
        scheduleRetry(req.RequestId, time.Now().Add(delay))
        
        req.RetryCount++
        updateRequest(req)
    } else {
        updateRequestStatus(req.RequestId, "FAILED")
        emitEvent("request_failed", req.RequestId, err.Error())
    }
}
```

## Integration Guide

### For Field Service Technicians

1. **Queue Operations for Offline Devices**:
```bash
# Queue pairing for module that's currently offline
racecar-webd tx pairingqueue queue-pairing-request \
  --initiator="comp_pack1" \
  --target="comp_module_123" \
  --priority="HIGH" \
  --expiry-hours=24 \
  --from service-key

# Check queue status
racecar-webd query pairingqueue get-queued-requests comp_module_123
```

2. **Process Queue When Device Returns**:
```bash
# When device comes back online
racecar-webd tx pairingqueue process-offline-queue \
  --component-id="comp_module_123" \
  --max-operations=5 \
  --from service-key
```

### For System Administrators

1. **Queue Monitoring**:
```go
// Monitor queue sizes across all components
func monitorQueues(ctx context.Context) {
    // Get all components with pending operations
    components := getAllComponentsWithQueues()
    
    for _, comp := range components {
        queueResp, err := queryClient.GetQueuedRequests(ctx,
            &types.QueryGetQueuedRequestsRequest{
                ComponentId: comp.ID,
                QueueType: "ALL",
            })
        
        if err != nil {
            continue
        }
        
        // Alert if queue too large
        if queueResp.QueueStatus.StandardQueueSize > 50 {
            alertLargeQueue(comp.ID, queueResp.QueueStatus)
        }
        
        // Alert if oldest request too old
        age := time.Since(queueResp.QueueStatus.OldestRequest)
        if age > 24*time.Hour {
            alertStaleQueue(comp.ID, age)
        }
    }
}
```

2. **Proxy Management**:
```go
// Setup proxy for permanently offline devices
func setupProxy(targetID, proxyID string) error {
    // Verify proxy authorization
    authResp, err := componentRegistry.CheckAuthorization(proxyID, targetID, "PROXY")
    if err != nil || !authResp.Authorized {
        return errors.New("proxy not authorized")
    }
    
    // Queue with proxy
    msg := &types.MsgQueuePairingRequest{
        InitiatorComponent: initiatorID,
        TargetComponent: targetID,
        RequestType: "PROXY",
        ProxyComponent: proxyID,
        Priority: "NORMAL",
    }
    
    return submitTransaction(msg)
}
```

### For Module Developers

The Pairing Queue module exposes the following keeper interface:

```go
type PairingQueueKeeper interface {
    // Request operations
    QueueRequest(ctx sdk.Context, req types.PairingRequest) error
    GetRequest(ctx sdk.Context, requestID string) (types.PairingRequest, bool)
    UpdateRequestStatus(ctx sdk.Context, requestID, status string) error
    
    // Queue operations
    GetQueuedRequests(ctx sdk.Context, componentID string) []types.PairingRequest
    ProcessQueue(ctx sdk.Context, componentID string) error
    
    // Offline operations
    MarkComponentOffline(ctx sdk.Context, componentID, reason string)
    MarkComponentOnline(ctx sdk.Context, componentID string) error
    
    // Utility methods
    GetQueueSize(ctx sdk.Context, componentID string) (int, error)
    CleanupExpiredRequests(ctx sdk.Context)
}
```

Example integration in Component Registry:

```go
func (k Keeper) HandleComponentStatusChange(ctx sdk.Context, componentID, newStatus string) {
    if newStatus == "ONLINE" {
        // Component came online, process its queue
        if err := k.pairingQueue.ProcessQueue(ctx, componentID); err != nil {
            k.Logger(ctx).Error("failed to process queue", "component", componentID, "error", err)
        }
    } else if newStatus == "OFFLINE" {
        // Mark component as offline
        k.pairingQueue.MarkComponentOffline(ctx, componentID, "status_change")
    }
}
```

### Best Practices

1. **Queue Management**:
   - Set appropriate expiry times based on expected offline duration
   - Use priority levels effectively for critical operations
   - Monitor queue sizes to prevent unbounded growth

2. **Proxy Usage**:
   - Only use trusted components as proxies
   - Verify proxy authorization before delegating
   - Log all proxy operations for audit

3. **Error Handling**:
   - Implement proper retry logic with exponential backoff
   - Handle permanent failures gracefully
   - Provide meaningful error messages to operators

4. **Performance**:
   - Process queues in batches to avoid blockchain congestion
   - Clean up expired and completed requests regularly
   - Use pagination for large queue queries