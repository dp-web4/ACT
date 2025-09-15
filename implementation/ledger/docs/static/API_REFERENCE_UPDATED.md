# Web4 Race Car Battery Management System - Complete API Reference

**Version**: 3.0  
**Chain ID**: `racecarweb`  
**REST Base URL**: `http://localhost:8080`  
**gRPC Base URL**: `localhost:9092`  
**Documentation Format**: Comprehensive API Guide

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication & Security](#authentication--security)
3. [API Bridge Interface](#api-bridge-interface)
4. [Account Management](#account-management)
5. [Component Registry API](#component-registry-api)
6. [LCT Manager API](#lct-manager-api)
7. [Pairing API](#pairing-api)
8. [Trust Tensor API](#trust-tensor-api)
9. [Energy Operations API](#energy-operations-api)
10. [Queue Management API](#queue-management-api)
11. [Authorization Management API](#authorization-management-api)
12. [Enhanced Trust Tensor API](#enhanced-trust-tensor-api)
13. [Real-time Streaming](#real-time-streaming)
14. [Error Handling](#error-handling)
15. [Performance Considerations](#performance-considerations)
16. [Business Logic Integration](#business-logic-integration)

---

## Overview

The Web4 Race Car Battery Management System provides a comprehensive blockchain-based solution for secure device pairing, trust management, energy operations, and intelligent queue processing. The system supports both REST and gRPC interfaces for maximum flexibility.

### Key Features

- **Dual Interface Support**: REST and gRPC APIs
- **Secure Pairing**: Split-key generation for device authentication
- **Trust Management**: Advanced tensor-based trust scoring with learning algorithms
- **Energy Operations**: Real-time energy transfer and monitoring with ATP/ADP token system
- **Queue Management**: Intelligent offline queue processing with retry logic
- **Authorization System**: Trust-based component authorization with configurable rules
- **Real-time Streaming**: Live battery status updates via gRPC
- **Business Logic Integration**: Comprehensive trust tensor calculations and energy cycle management

### Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   C++ Client    │    │   REST Client   │    │   gRPC Client   │
│   (RAD Studio)  │    │   (Web/HTTP)    │    │   (High Perf)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      API Bridge           │
                    │   (Go Application)        │
                    │  REST: :8080              │
                    │  gRPC: :9092              │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │   Blockchain (Cosmos)     │
                    │   Chain ID: racecarweb    │
                    │                           │
                    │  ┌─────────────────────┐  │
                    │  │  Trust Tensor T3/V3 │  │
                    │  │  Calculations       │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │  Energy Cycle       │  │
                    │  │  ATP/ADP Logic      │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │  Authorization      │  │
                    │  │  Management         │  │
                    │  └─────────────────────┘  │
                    │  ┌─────────────────────┐  │
                    │  │  Queue Processing   │  │
                    │  │  & Retry Logic      │  │
                    │  └─────────────────────┘  │
                    └───────────────────────────┘
```

---

## Authentication & Security

### Account Management

The API Bridge provides simplified account management that abstracts the complexity of Cosmos SDK authentication.

#### Get Accounts

**REST Endpoint**: `GET /accounts`

**gRPC Service**: `APIBridgeService.GetAccounts`

**Description**: Retrieve all available accounts for transaction signing

**REST Response**:
```json
{
  "accounts": [
    {
      "name": "demo-user",
      "address": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
      "key_type": "secp256k1"
    }
  ]
}
```

**gRPC Response**:
```protobuf
message GetAccountsResponse {
  repeated Account accounts = 1;
}

message Account {
  string name = 1;
  string address = 2;
  string key_type = 3;
}
```

#### Create Account

**REST Endpoint**: `POST /accounts`

**gRPC Service**: `APIBridgeService.CreateAccount`

**Request Body**:
```json
{
  "name": "new-user"
}
```

**Response**:
```json
{
  "name": "new-user",
  "address": "cosmos1abc123...",
  "key_type": "secp256k1"
}
```

### Security Features

- **SSL/TLS**: Required for production environments
- **Transaction Signing**: Automatic signing via Ignite CLI
- **Input Validation**: Comprehensive parameter checking
- **Rate Limiting**: API usage monitoring

---

## API Bridge Interface

The API Bridge provides a unified interface that abstracts blockchain complexity and provides both REST and gRPC access.

### Endpoint Configuration

```json
{
  "rest": {
    "endpoint": "http://localhost:8080",
    "timeout": 30,
    "retry_attempts": 3
  },
  "grpc": {
    "endpoint": "localhost:9092",
    "timeout": 30,
    "max_message_size": 1048576
  }
}
```

### Common Response Format

All API responses follow a consistent format:

```json
{
  "status": "success",
  "data": {
    // Response data
  },
  "tx_hash": "ABC123DEF456...",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

---

## Component Registry API

### Register Component

**REST Endpoint**: `POST /components/register`

**gRPC Service**: `APIBridgeService.RegisterComponent`

**Description**: Register a new component in the blockchain registry

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_data": "ModBatt Pro Series MB-4680-LFP",
  "context": "race-car-battery"
}
```

**Response**:
```json
{
  "component_id": "comp_abc123def456",
  "component_identity": "comp_identity_hash",
  "component_data": "ModBatt Pro Series MB-4680-LFP",
  "context": "race-car-battery",
  "creator": "demo-user",
  "lct_id": "lct_abc123def456",
  "status": "registered",
  "tx_hash": "ABC123DEF456..."
}
```

### Get Component

**REST Endpoint**: `GET /components/{component_id}`

**gRPC Service**: `APIBridgeService.GetComponent`

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "component_id": "comp_abc123def456",
  "component_identity": "comp_identity_hash",
  "component_data": "ModBatt Pro Series MB-4680-LFP",
  "context": "race-car-battery",
  "creator": "demo-user",
  "lct_id": "lct_abc123def456",
  "status": "active",
  "tx_hash": "ABC123DEF456..."
}
```

### Get Component Identity

**REST Endpoint**: `GET /components/{component_id}/identity`

**gRPC Service**: `APIBridgeService.GetComponentIdentity`

**Description**: Retrieve component identity information

### Verify Component

**REST Endpoint**: `POST /components/verify`

**gRPC Service**: `APIBridgeService.VerifyComponent`

**Request Body**:
```json
{
  "verifier": "demo-user",
  "component_id": "comp_abc123def456",
  "context": "verification-context"
}
```

---

## LCT Manager API

### Create LCT (Linked Context Token)

**REST Endpoint**: `POST /lct/create`

**gRPC Service**: `APIBridgeService.CreateLCT`

**Description**: Create a Linked Context Token for secure device relationships

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "context": "battery-motor-pairing",
  "proxy_id": "proxy_001"
}
```

**Response**:
```json
{
  "lct_id": "lct_abc123def456",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "context": "battery-motor-pairing",
  "proxy_id": "proxy_001",
  "status": "active",
  "created_at": 1704067200,
  "creator": "demo-user",
  "tx_hash": "ABC123DEF456...",
  "lct_key_half": "encrypted_lct_key_half",
  "device_key_half": "encrypted_device_key_half"
}
```

### Get LCT

**REST Endpoint**: `GET /lct/{lct_id}`

**gRPC Service**: `APIBridgeService.GetLCT`

**Path Parameters**:
- `lct_id`: LCT identifier

### Update LCT Status

**REST Endpoint**: `PUT /lct/{lct_id}/status`

**gRPC Service**: `APIBridgeService.UpdateLCTStatus`

**Request Body**:
```json
{
  "creator": "demo-user",
  "lct_id": "lct_abc123def456",
  "status": "maintenance",
  "context": "scheduled-maintenance"
}
```

---

## Pairing API

### Initiate Pairing

**REST Endpoint**: `POST /pairing/initiate`

**gRPC Service**: `APIBridgeService.InitiatePairing`

**Description**: Initiate a bidirectional pairing process between components

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "operational_context": "race-car-operation",
  "proxy_id": "proxy_001",
  "force_immediate": false
}
```

**Response**:
```json
{
  "challenge_id": "challenge_abc123def456",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "operational_context": "race-car-operation",
  "proxy_id": "proxy_001",
  "force_immediate": false,
  "status": "pending",
  "created_at": 1704067200,
  "creator": "demo-user",
  "tx_hash": "ABC123DEF456..."
}
```

### Complete Pairing

**REST Endpoint**: `POST /pairing/complete`

**gRPC Service**: `APIBridgeService.CompletePairing`

**Description**: Complete pairing with split-key generation

**Request Body**:
```json
{
  "creator": "demo-user",
  "challenge_id": "challenge_abc123def456",
  "component_a_auth": "battery-authentication",
  "component_b_auth": "motor-authentication",
  "session_context": "race-session-001"
}
```

**Response**:
```json
{
  "lct_id": "lct_abc123def456",
  "session_keys": "encrypted_session_keys",
  "trust_summary": "trust_summary_data",
  "tx_hash": "ABC123DEF456...",
  "split_key_a": "32_byte_hex_string_a",
  "split_key_b": "32_byte_hex_string_b"
}
```

**Important**: The split keys are generated cryptographically and each half is 32 bytes (64 hex characters). Combined, they form a 64-byte key for strong cryptographic security.

### Revoke Pairing

**REST Endpoint**: `POST /pairing/revoke`

**gRPC Service**: `APIBridgeService.RevokePairing`

**Request Body**:
```json
{
  "creator": "demo-user",
  "lct_id": "lct_abc123def456",
  "reason": "component-replacement",
  "notify_offline": true
}
```

### Get Pairing Status

**REST Endpoint**: `GET /pairing/status/{challenge_id}`

**gRPC Service**: `APIBridgeService.GetPairingStatus`

**Path Parameters**:
- `challenge_id`: Pairing challenge identifier

---

## Trust Tensor API

### Create Trust Tensor

**REST Endpoint**: `POST /trust/create`

**gRPC Service**: `APIBridgeService.CreateTrustTensor`

**Description**: Create a trust relationship tensor

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "context": "trust-evaluation",
  "initial_score": 0.8
}
```

**Response**:
```json
{
  "tensor_id": "tensor_abc123def456",
  "score": 0.8,
  "status": "active",
  "tx_hash": "ABC123DEF456..."
}
```

### Get Trust Tensor

**REST Endpoint**: `GET /trust/{tensor_id}`

**gRPC Service**: `APIBridgeService.GetTrustTensor`

### Update Trust Score

**REST Endpoint**: `PUT /trust/{tensor_id}/score`

**gRPC Service**: `APIBridgeService.UpdateTrustScore`

**Request Body**:
```json
{
  "creator": "demo-user",
  "tensor_id": "tensor_abc123def456",
  "score": 0.9,
  "context": "performance-improvement"
}
```

---

## Energy Operations API

### Create Energy Operation

**REST Endpoint**: `POST /energy/create`

**gRPC Service**: `APIBridgeService.CreateEnergyOperation`

**Description**: Create an energy operation between components

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_a": "comp_abc123def456",
  "component_b": "comp_def456ghi789",
  "operation_type": "transfer",
  "amount": 100.5,
  "context": "energy-transfer"
}
```

**Response**:
```json
{
  "operation_id": "op_abc123def456",
  "operation_type": "transfer",
  "amount": 100.5,
  "status": "pending",
  "tx_hash": "ABC123DEF456..."
}
```

### Execute Energy Transfer

**REST Endpoint**: `POST /energy/transfer`

**gRPC Service**: `APIBridgeService.ExecuteEnergyTransfer`

**Request Body**:
```json
{
  "creator": "demo-user",
  "operation_id": "op_abc123def456",
  "amount": 100.5,
  "context": "transfer-execution"
}
```

### Get Energy Balance

**REST Endpoint**: `GET /energy/balance/{component_id}`

**gRPC Service**: `APIBridgeService.GetEnergyBalance`

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "balance": 1500.75
}
```

---

## Queue Management API

The Queue Management API provides intelligent offline queue processing with retry logic and failure handling for components that may be temporarily unavailable.

### Queue Pairing Request

**REST Endpoint**: `POST /api/v1/queue/pairing-request`

**gRPC Service**: `APIBridgeService.QueuePairingRequest`

**Description**: Queue a pairing request for offline processing

**Request Body**:
```json
{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "proxy_id": "proxy_001"
}
```

**Response**:
```json
{
  "request_id": "queue_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "proxy_id": "proxy_001",
  "status": "queued",
  "created_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Get Queue Status

**REST Endpoint**: `GET /api/v1/queue/status/{component_id}`

**gRPC Service**: `APIBridgeService.GetQueueStatus`

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "component_id": "battery_001",
  "status": "active",
  "pending_requests": 5,
  "processed_requests": 12,
  "last_processed": 1704067200
}
```

### Process Offline Queue

**REST Endpoint**: `POST /api/v1/queue/process-offline/{component_id}`

**gRPC Service**: `APIBridgeService.ProcessOfflineQueue`

**Description**: Process all pending operations for a component

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "component_id": "battery_001",
  "status": "processed",
  "processed_requests": 3,
  "failed_requests": 0,
  "processed_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Cancel Request

**REST Endpoint**: `DELETE /api/v1/queue/cancel/{request_id}`

**gRPC Service**: `APIBridgeService.CancelRequest`

**Description**: Cancel a queued request

**Path Parameters**:
- `request_id`: Request identifier

**Request Body**:
```json
{
  "reason": "component_offline"
}
```

**Response**:
```json
{
  "request_id": "queue_battery_001_battery_002_1704067200",
  "status": "cancelled",
  "reason": "component_offline",
  "cancelled_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Get Queued Requests

**REST Endpoint**: `GET /api/v1/queue/requests/{component_id}`

**gRPC Service**: `APIBridgeService.GetQueuedRequests`

**Description**: Get all queued requests for a component

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "requests": [
    {
      "request_id": "req_001",
      "component_a": "battery_001",
      "component_b": "battery_002",
      "operational_context": "energy_transfer",
      "status": "pending",
      "created_at": 1704066900
    },
    {
      "request_id": "req_002",
      "component_a": "battery_001",
      "component_b": "battery_003",
      "operational_context": "pairing",
      "status": "queued",
      "created_at": 1704066600
    }
  ],
  "count": 2,
  "component_id": "battery_001"
}
```

### List Proxy Queue

**REST Endpoint**: `GET /api/v1/queue/proxy/{proxy_id}`

**gRPC Service**: `APIBridgeService.ListProxyQueue`

**Description**: List all operations for a proxy

**Path Parameters**:
- `proxy_id`: Proxy identifier

**Response**:
```json
{
  "operations": [
    {
      "operation_id": "op_001",
      "component_a": "battery_001",
      "component_b": "battery_002",
      "operation_type": "energy_transfer",
      "status": "pending",
      "created_at": 1704066900
    },
    {
      "operation_id": "op_002",
      "component_a": "battery_003",
      "component_b": "battery_004",
      "operation_type": "pairing",
      "status": "completed",
      "created_at": 1704066600
    }
  ],
  "count": 2,
  "proxy_id": "proxy_001"
}
```

---

## Authorization Management API

The Authorization Management API provides trust-based component authorization with configurable rules and comprehensive lifecycle management.

### Create Pairing Authorization

**REST Endpoint**: `POST /api/v1/authorization/pairing`

**gRPC Service**: `APIBridgeService.CreatePairingAuthorization`

**Description**: Create a pairing authorization between components

**Request Body**:
```json
{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "authorization_rules": "trust_score > 0.7"
}
```

**Response**:
```json
{
  "authorization_id": "auth_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "authorization_rules": "trust_score > 0.7",
  "status": "active",
  "created_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Get Component Authorizations

**REST Endpoint**: `GET /api/v1/authorization/component/{component_id}`

**gRPC Service**: `APIBridgeService.GetComponentAuthorizations`

**Description**: Get all authorizations for a component

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "authorizations": [
    {
      "authorization_id": "auth_001",
      "component_a": "battery_001",
      "component_b": "battery_002",
      "operational_context": "energy_transfer",
      "authorization_rules": "trust_score > 0.7",
      "status": "active",
      "created_at": 1704063600
    },
    {
      "authorization_id": "auth_002",
      "component_a": "battery_001",
      "component_b": "battery_003",
      "operational_context": "pairing",
      "authorization_rules": "verified_component",
      "status": "active",
      "created_at": 1704060000
    }
  ],
  "count": 2,
  "component_id": "battery_001"
}
```

### Update Authorization

**REST Endpoint**: `PUT /api/v1/authorization/{authorization_id}`

**gRPC Service**: `APIBridgeService.UpdateAuthorization`

**Description**: Update authorization rules or context

**Path Parameters**:
- `authorization_id`: Authorization identifier

**Request Body**:
```json
{
  "operational_context": "high_priority_energy_transfer",
  "authorization_rules": "trust_score > 0.8",
  "status": "active"
}
```

**Response**:
```json
{
  "authorization_id": "auth_001",
  "operational_context": "high_priority_energy_transfer",
  "authorization_rules": "trust_score > 0.8",
  "status": "active",
  "updated_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Revoke Authorization

**REST Endpoint**: `DELETE /api/v1/authorization/{authorization_id}`

**gRPC Service**: `APIBridgeService.RevokeAuthorization`

**Description**: Revoke an authorization

**Path Parameters**:
- `authorization_id`: Authorization identifier

**Request Body**:
```json
{
  "reason": "security_concern"
}
```

**Response**:
```json
{
  "authorization_id": "auth_001",
  "status": "revoked",
  "reason": "security_concern",
  "revoked_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Check Pairing Authorization

**REST Endpoint**: `GET /api/v1/authorization/check`

**gRPC Service**: `APIBridgeService.CheckPairingAuthorization`

**Description**: Check if pairing is authorized between components

**Query Parameters**:
- `component_a`: First component identifier
- `component_b`: Second component identifier
- `operational_context`: Operational context

**Response**:
```json
{
  "authorized": true,
  "authorization_id": "auth_001",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "trust_score": 0.85,
  "verification_status": "verified",
  "checked_at": 1704067200
}
```

---

## Enhanced Trust Tensor API

The Enhanced Trust Tensor API provides advanced trust calculations with learning algorithms, evidence weighting, and composite scoring for sophisticated trust management.

### Calculate Relationship Trust

**REST Endpoint**: `POST /api/v1/trust-enhanced/calculate`

**gRPC Service**: `APIBridgeService.CalculateRelationshipTrust`

**Description**: Calculate trust score for a relationship using advanced algorithms

**Request Body**:
```json
{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer"
}
```

**Response**:
```json
{
  "tensor_id": "tensor_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "trust_score": 0.82,
  "status": "calculated",
  "calculated_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Get Relationship Tensor

**REST Endpoint**: `GET /api/v1/trust-enhanced/relationship`

**gRPC Service**: `APIBridgeService.GetRelationshipTensor`

**Description**: Get detailed trust tensor information for a relationship

**Query Parameters**:
- `component_a`: First component identifier
- `component_b`: Second component identifier

**Response**:
```json
{
  "tensor_id": "tensor_battery_001_battery_002",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "score": 0.82,
  "status": "active",
  "created_at": 1704063600,
  "last_updated": 1704067200,
  "evidence_count": 15,
  "learning_rate": 0.1
}
```

### Update Tensor Score

**REST Endpoint**: `PUT /api/v1/trust-enhanced/score`

**gRPC Service**: `APIBridgeService.UpdateTensorScore`

**Description**: Update trust score with evidence weighting and learning algorithms

**Request Body**:
```json
{
  "creator": "demo-user",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "score": 0.85,
  "context": "successful_energy_transfer"
}
```

**Response**:
```json
{
  "tensor_id": "tensor_battery_001_battery_002",
  "creator": "demo-user",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "score": 0.85,
  "context": "successful_energy_transfer",
  "updated_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Trust Tensor Features

The Enhanced Trust Tensor API includes advanced features:

- **Composite Scoring**: Combines multiple trust factors
- **Learning Rate**: Adaptive trust score adjustments
- **Evidence Weighting**: Recent interactions weighted more heavily
- **Decay Factor**: Trust scores decay over time
- **Context Awareness**: Different operational contexts have different trust requirements

---

## Real-time Streaming

### Battery Status Streaming (gRPC Only)

**gRPC Service**: `APIBridgeService.StreamBatteryStatus`

**Description**: Stream real-time battery status updates

**Request**:
```protobuf
message StreamBatteryStatusRequest {
  string component_id = 1;
  int32 update_interval_seconds = 2;
}
```

**Response Stream**:
```protobuf
message BatteryStatusUpdate {
  string component_id = 1;
  double voltage = 2;
  double current = 3;
  double temperature = 4;
  double state_of_charge = 5;
  string status = 6;
  int64 timestamp = 7;
}
```

**Example Usage (C++)**:
```cpp
// Initialize gRPC client
GRPCClient client("localhost:9092");

// Stream battery status
client.streamBatteryStatus("battery-001", 2, [](const BatteryStatusUpdate& update) {
    std::cout << "Voltage: " << update.voltage << "V" << std::endl;
    std::cout << "Current: " << update.current << "A" << std::endl;
    std::cout << "Temperature: " << update.temperature << "°C" << std::endl;
    std::cout << "SOC: " << update.stateOfCharge << "%" << std::endl;
    std::cout << "Status: " << update.status << std::endl;
});
```

---

## Error Handling

### Standard Error Format

```json
{
  "error": {
    "code": "INVALID_COMPONENT_ID",
    "message": "Component ID format is invalid",
    "details": {
      "component_id": "invalid-format",
      "expected_format": "comp_[a-z0-9]{12}"
    }
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "request_id": "req_abc123def456"
}
```

### Common Error Codes

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `INVALID_COMPONENT_ID` | Component ID format error | 400 |
| `COMPONENT_NOT_FOUND` | Component doesn't exist | 404 |
| `PAIRING_FAILED` | Pairing operation failed | 400 |
| `INSUFFICIENT_PERMISSIONS` | Authorization error | 403 |
| `BLOCKCHAIN_ERROR` | Blockchain transaction failed | 500 |
| `RATE_LIMIT_EXCEEDED` | Too many requests | 429 |

### Retry Logic

The API Bridge implements automatic retry logic with exponential backoff:

```json
{
  "retry_config": {
    "max_attempts": 3,
    "initial_delay": 1000,
    "max_delay": 10000,
    "backoff_multiplier": 2
  }
}
```

---

## Performance Considerations

### REST vs gRPC Comparison

| Metric | REST | gRPC | Improvement |
|--------|------|------|-------------|
| Account List | 45ms | 12ms | 73% faster |
| Component Register | 120ms | 35ms | 71% faster |
| LCT Create | 180ms | 52ms | 71% faster |
| Pairing Initiate | 150ms | 42ms | 72% faster |
| Pairing Complete | 200ms | 58ms | 71% faster |

### Optimization Recommendations

1. **Use gRPC for high-frequency operations**
2. **Implement connection pooling**
3. **Batch operations when possible**
4. **Use streaming for real-time data**
5. **Cache frequently accessed data**

### Rate Limits

- **REST API**: 100 requests/minute per IP
- **gRPC API**: 500 requests/minute per connection
- **Streaming**: 10 concurrent streams per client

---

## Business Logic Integration

The Web4 Race Car Battery Management System integrates sophisticated business logic across all modules to ensure secure, efficient, and intelligent operations.

### Trust Tensor T3/V3 Calculations

The system implements advanced trust tensor calculations with the following features:

#### Composite Scoring Algorithm
```go
// Trust score calculation combines multiple factors
composite_score = (base_score * 0.4) + (interaction_score * 0.3) + (time_score * 0.3)
```

#### Learning Rate Adaptation
- **Initial Learning Rate**: 0.1 (10% weight for new evidence)
- **Adaptive Adjustment**: Learning rate decreases as evidence accumulates
- **Evidence Weighting**: Recent interactions weighted more heavily than older ones

#### Decay Factor Implementation
- **Time-based Decay**: Trust scores decay over time if no new interactions occur
- **Context-specific Decay**: Different operational contexts have different decay rates
- **Minimum Threshold**: Trust scores cannot decay below a minimum threshold

### Energy Cycle ATP/ADP Logic

The system manages energy cycles using ATP (Adenosine Triphosphate) and ADP (Adenosine Diphosphate) tokens:

#### ATP Token Creation
```go
// ATP tokens represent available energy
atp_tokens = energy_balance * conversion_factor
```

#### ADP Token Discharge
```go
// ADP tokens represent discharged energy
adp_tokens = discharged_energy * conversion_factor
```

#### Energy Balance Calculation
```go
// Real-time energy balance tracking
energy_balance = atp_tokens - adp_tokens
```

#### Operation Validation
- **Trust Score Check**: Operations require minimum trust score
- **Energy Availability**: Sufficient ATP tokens required
- **Authorization Verification**: Component authorization validated

### Component Registry Authorization

The authorization system provides multi-level security:

#### Authorization Levels
1. **Level 1**: Basic component verification
2. **Level 2**: Trust score validation
3. **Level 3**: Operational context verification
4. **Level 4**: Real-time authorization check

#### Pairing Validation
```go
// Authorization check for pairing operations
authorized = is_verified(component_a) && 
            is_verified(component_b) && 
            check_trust_score(component_a, component_b) >= min_trust_score &&
            validate_operational_context(context)
```

### Pairing Queue Offline Processing

The queue system provides robust offline processing capabilities:

#### Retry Logic
- **Exponential Backoff**: Retry delays increase exponentially
- **Maximum Retries**: Configurable maximum retry attempts
- **Failure Handling**: Failed operations logged and reported

#### Queue Management
```go
// Queue processing with trust and authorization checks
for operation in queue {
    if check_authorization(operation) && check_trust_score(operation) {
        process_operation(operation)
        update_trust_score(operation)
    } else {
        mark_for_retry(operation)
    }
}
```

#### Offline Operation Handling
- **Persistent Storage**: Operations stored on-chain
- **Status Tracking**: Real-time status updates
- **Proxy Management**: Proxy components handle offline operations

### Integration Benefits

#### Security
- **Multi-factor Authentication**: Trust scores, authorization, and verification
- **Real-time Validation**: Continuous security checks
- **Audit Trail**: Complete transaction history

#### Efficiency
- **Intelligent Routing**: Operations routed based on trust and availability
- **Load Balancing**: Distributed processing across components
- **Resource Optimization**: Energy and computational resource management

#### Reliability
- **Fault Tolerance**: Graceful handling of component failures
- **Data Consistency**: Atomic operations ensure data integrity
- **Recovery Mechanisms**: Automatic recovery from failures

### Performance Metrics

The system tracks comprehensive performance metrics:

#### Trust Metrics
- **Trust Score Distribution**: Statistical analysis of trust scores
- **Learning Rate Effectiveness**: Impact of learning rate on trust accuracy
- **Decay Factor Impact**: Effect of decay on trust score stability

#### Energy Metrics
- **ATP/ADP Balance**: Real-time energy token tracking
- **Transfer Efficiency**: Energy transfer success rates
- **Token Utilization**: ATP/ADP token usage patterns

#### Queue Metrics
- **Processing Latency**: Time from queue to completion
- **Success Rates**: Operation success percentages
- **Retry Patterns**: Retry frequency and patterns

---

## Client Examples

### C++ Client (RAD Studio Compatible)

```cpp
#include "RESTClient.h"
#include "GRPCClient.h"

// REST API usage
RESTClient restClient("http://localhost:8080");
auto result = restClient.registerComponent("user1", "battery-data", "context");

// gRPC API usage
GRPCClient grpcClient("localhost:9092");
auto lct = grpcClient.createLCT("user1", "comp_a", "comp_b", "context", "proxy");

// Real-time streaming
grpcClient.streamBatteryStatus("battery-001", 2, [](const BatteryStatusUpdate& update) {
    // Handle battery status updates
});
```

### JavaScript Client

```javascript
// REST API
const response = await fetch('http://localhost:8080/components/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    creator: 'user1',
    component_data: 'battery-data',
    context: 'context'
  })
});

// gRPC API (using grpc-web)
const client = new APIBridgeServiceClient('http://localhost:9092');
const result = await client.registerComponent({
  creator: 'user1',
  component_data: 'battery-data',
  context: 'context'
});
```

### Python Client

```python
import requests
import grpc

# REST API
response = requests.post('http://localhost:8080/components/register', json={
    'creator': 'user1',
    'component_data': 'battery-data',
    'context': 'context'
})

# gRPC API
import apibridge_pb2
import apibridge_pb2_grpc

channel = grpc.insecure_channel('localhost:9092')
stub = apibridge_pb2_grpc.APIBridgeServiceStub(channel)
result = stub.RegisterComponent(apibridge_pb2.RegisterComponentRequest(
    creator='user1',
    component_data='battery-data',
    context='context'
))
```

---

## Migration Guide

### From Direct Blockchain API to API Bridge

**Before (Direct Cosmos SDK)**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.componentregistry.v1.MsgRegisterComponent",
        "creator": "cosmos1...",
        "component_id": "MODBATT-MOD-RC001-001",
        "component_type": "module",
        "manufacturer_data": "ModBatt Pro Series"
      }]
    }
  }
}
```

**After (API Bridge)**:
```json
{
  "creator": "demo-user",
  "component_data": "ModBatt Pro Series MB-4680-LFP",
  "context": "race-car-battery"
}
```

### Benefits of API Bridge

1. **Simplified Interface**: No need to understand Cosmos SDK internals
2. **Automatic Transaction Signing**: Handled by the bridge
3. **Error Handling**: Comprehensive error management
4. **Performance**: Optimized for high-throughput operations
5. **Dual Interface**: Choose REST or gRPC based on needs

---

## Conclusion

The Web4 Race Car Battery Management System provides a comprehensive, secure, and high-performance solution for blockchain-based battery management. The dual REST/gRPC interface ensures compatibility with all development environments while maintaining optimal performance for real-time operations.

For implementation examples and detailed guides, refer to the C++ demo application in the `api-bridge/cpp-demo/` directory. 