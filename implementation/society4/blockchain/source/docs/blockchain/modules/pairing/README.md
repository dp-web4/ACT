# Pairing Module

## Table of Contents

1. [Overview](#overview)
2. [Key Concepts](#key-concepts)
3. [State Management](#state-management)
4. [Messages (Transactions)](#messages-transactions)
5. [Queries](#queries)
6. [Events](#events)
7. [Parameters](#parameters)
8. [Security Considerations](#security-considerations)
9. [Integration Guide](#integration-guide)

## Overview

The Pairing module manages the secure establishment of bidirectional authentication between battery components. It implements a challenge-response protocol to ensure that only legitimate components can establish relationships, providing the foundation for secure LCT creation and energy operations.

### Purpose
- Establish secure bidirectional authentication between components
- Manage pairing sessions with time-bound validity
- Maintain active pairing records
- Support pairing revocation and lifecycle management

### Dependencies
- **Component Registry**: Verifies component identities and checks authorization

### Module Store Key
`pairing`

## Key Concepts

### Pairing Session
A time-bound authentication session between two components:
- **Session ID**: Unique identifier for the pairing attempt
- **Initiator**: Component that starts the pairing
- **Target**: Component being paired with
- **Challenge**: Cryptographic challenge for authentication
- **Expiry**: Time limit for completing the pairing

### Pairing Challenge
Cryptographic proof mechanism:
- **Challenge Data**: Random data to be signed
- **Expected Response**: Signature proving component ownership
- **Verification**: Ensures component controls its private key

### Active Pairing
Successfully authenticated relationship:
- **Pairing ID**: Unique identifier for the established pairing
- **Participants**: Both authenticated components
- **Establishment Time**: When pairing was completed
- **Status**: Active or Revoked

## State Management

### Stored Types

#### 1. PairingSession
```protobuf
message PairingSession {
  string session_id = 1;
  string initiator_component = 2;
  string target_component = 3;
  bytes challenge_data = 4;
  string session_status = 5;  // PENDING, COMPLETED, EXPIRED, FAILED
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp expires_at = 7;
  string initiator_signature = 8;
  map<string, string> metadata = 9;
}
```

#### 2. PairingChallenge
```protobuf
message PairingChallenge {
  string challenge_id = 1;
  string session_id = 2;
  string challenged_component = 3;
  bytes challenge_data = 4;
  string expected_response_hash = 5;
  bool is_verified = 6;
  google.protobuf.Timestamp verified_at = 7;
}
```

#### 3. ActivePairing
```protobuf
message ActivePairing {
  string pairing_id = 1;
  string component_a = 2;
  string component_b = 3;
  string session_id = 4;
  google.protobuf.Timestamp established_at = 5;
  string status = 6;  // ACTIVE, REVOKED
  google.protobuf.Timestamp revoked_at = 7;
  string revoked_by = 8;
  string revocation_reason = 9;
}
```

### Store Layout
```
pairing/
├── sessions/
│   └── {session_id} → PairingSession
├── challenges/
│   └── {session_id}/
│       └── {component_id} → PairingChallenge
├── active/
│   └── {component_a}/
│       └── {component_b} → ActivePairing
└── params → Params
```

## Messages (Transactions)

### 1. InitiateBidirectionalPairing
Starts a new pairing session between two components.

**Input**:
```protobuf
message MsgInitiateBidirectionalPairing {
  string creator = 1;
  string initiator_component = 2;
  string target_component = 3;
  uint32 validity_minutes = 4;
  map<string, string> metadata = 5;
}
```

**Validation**:
- Both components must exist in Component Registry
- Components must have pairing authorization
- No active pairing between components
- Creator must own initiator component
- Validity period within allowed range

**Effects**:
- Creates new PairingSession
- Generates challenge data for both components
- Sets session expiry time
- Emits `pairing_initiated` event

**Example**:
```bash
racecar-webd tx pairing initiate-bidirectional-pairing \
  --initiator="comp_module1" \
  --target="comp_pack1" \
  --validity-minutes=30 \
  --from mykey
```

### 2. CompletePairing
Completes the pairing by providing challenge responses.

**Input**:
```protobuf
message MsgCompletePairing {
  string creator = 1;
  string session_id = 2;
  string component_id = 3;
  bytes challenge_response = 4;
}
```

**Validation**:
- Session must exist and be pending
- Session must not be expired
- Component must be participant in session
- Challenge response must be valid signature
- Both components must complete challenge

**Effects**:
- Verifies challenge response
- If both challenges completed:
  - Creates ActivePairing record
  - Updates session status to COMPLETED
  - Emits `pairing_completed` event
- If only one challenge completed:
  - Marks challenge as verified
  - Waits for other component

**Challenge-Response Flow**:
```
1. Component A initiates pairing
2. System generates challenges for A and B
3. A signs its challenge → CompletePairing
4. B signs its challenge → CompletePairing
5. When both complete → Active Pairing established
```

### 3. RevokePairing
Revokes an active pairing between components.

**Input**:
```protobuf
message MsgRevokePairing {
  string creator = 1;
  string component_a = 2;
  string component_b = 3;
  string revocation_reason = 4;
}
```

**Validation**:
- Active pairing must exist
- Creator must own one of the components
- Must provide revocation reason

**Effects**:
- Updates pairing status to REVOKED
- Records revocation details
- Prevents future operations using this pairing
- Emits `pairing_revoked` event

## Queries

### 1. GetPairingStatus
Gets the current pairing status between two components.

**Request**:
```protobuf
message QueryGetPairingStatusRequest {
  string component_a = 1;
  string component_b = 2;
}
```

**Response**:
```protobuf
message QueryGetPairingStatusResponse {
  bool is_paired = 1;
  ActivePairing pairing_details = 2;
  PairingSession pending_session = 3;
}
```

**Example**:
```bash
racecar-webd query pairing get-pairing-status comp_module1 comp_pack1
```

### 2. ListActivePairings
Lists all active pairings for a component.

**Request**:
```protobuf
message QueryListActivePairingsRequest {
  string component_id = 1;
  cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
```

**Response**:
```protobuf
message QueryListActivePairingsResponse {
  repeated ActivePairing pairings = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### 3. ValidateBidirectionalAuth
Validates that two components have completed bidirectional authentication.

**Request**:
```protobuf
message QueryValidateBidirectionalAuthRequest {
  string component_a = 1;
  string component_b = 2;
  bool check_authorization = 3;
}
```

**Response**:
```protobuf
message QueryValidateBidirectionalAuthResponse {
  bool is_authenticated = 1;
  bool is_authorized = 2;
  string pairing_id = 3;
  google.protobuf.Timestamp authenticated_at = 4;
}
```

## Events

### pairing_initiated
Emitted when a pairing session starts.
```json
{
  "type": "pairing_initiated",
  "attributes": [
    {"key": "session_id", "value": "sess_abc123"},
    {"key": "initiator", "value": "comp_module1"},
    {"key": "target", "value": "comp_pack1"},
    {"key": "expires_at", "value": "2024-01-15T11:00:00Z"}
  ]
}
```

### pairing_completed
Emitted when pairing is successfully completed.
```json
{
  "type": "pairing_completed",
  "attributes": [
    {"key": "pairing_id", "value": "pair_def456"},
    {"key": "component_a", "value": "comp_module1"},
    {"key": "component_b", "value": "comp_pack1"},
    {"key": "session_id", "value": "sess_abc123"}
  ]
}
```

### pairing_revoked
Emitted when a pairing is revoked.
```json
{
  "type": "pairing_revoked",
  "attributes": [
    {"key": "pairing_id", "value": "pair_def456"},
    {"key": "revoked_by", "value": "cosmos1..."},
    {"key": "reason", "value": "Component decommissioned"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Maximum validity period for pairing sessions (minutes)
  uint32 max_session_validity_minutes = 1;
  
  // Default validity period if not specified (minutes)
  uint32 default_session_validity_minutes = 2;
  
  // Require component authorization before pairing
  bool require_authorization = 3;
  
  // Allow self-revocation of pairings
  bool allow_self_revocation = 4;
  
  // Challenge data size in bytes
  uint32 challenge_size_bytes = 5;
  
  // Maximum active pairings per component
  uint32 max_pairings_per_component = 6;
}
```

**Default Values**:
```json
{
  "max_session_validity_minutes": 60,
  "default_session_validity_minutes": 30,
  "require_authorization": true,
  "allow_self_revocation": true,
  "challenge_size_bytes": 32,
  "max_pairings_per_component": 50
}
```

## Security Considerations

### Challenge Generation
- Uses cryptographically secure random number generator
- Challenge size configurable (default 32 bytes)
- Each component gets unique challenge data

### Response Verification
```go
// Simplified verification logic
func verifyChallenge(component Component, challenge []byte, response []byte) bool {
    // Recover public key from component
    pubKey := component.PublicKey
    
    // Verify signature
    return ed25519.Verify(pubKey, challenge, response)
}
```

### Session Security
- Time-bound sessions prevent replay attacks
- Expired sessions automatically invalidated
- One-time use challenges

### Revocation Mechanism
- Either component can revoke pairing
- Revocation is permanent
- Reason required for audit trail

## Integration Guide

### For Component Manufacturers

1. **Factory Pre-Pairing**:
```go
// During manufacturing, pre-pair module with controller
msg := &types.MsgInitiateBidirectionalPairing{
    Creator: factoryAddr,
    InitiatorComponent: controllerID,
    TargetComponent: moduleID,
    ValidityMinutes: 60,
    Metadata: map[string]string{
        "pairing_type": "factory",
        "batch_number": "2024-001",
        "production_date": "2024-01-15",
    },
}

// Complete pairing with pre-shared keys
completeMsg := &types.MsgCompletePairing{
    Creator: factoryAddr,
    SessionId: sessionID,
    ComponentId: moduleID,
    ChallengeResponse: signChallenge(moduleKey, challengeData),
}
```

2. **Batch Pairing Operations**:
```go
// Pair multiple modules with a pack
for _, moduleID := range moduleIDs {
    initMsg := &types.MsgInitiateBidirectionalPairing{
        Creator: packControllerAddr,
        InitiatorComponent: packID,
        TargetComponent: moduleID,
        ValidityMinutes: 30,
    }
    // Process each pairing
}
```

### For System Integrators

1. **Field Pairing Workflow**:
```go
// Check if components are already paired
statusResp, err := queryClient.GetPairingStatus(ctx,
    &types.QueryGetPairingStatusRequest{
        ComponentA: moduleID,
        ComponentB: packID,
    })

if !statusResp.IsPaired {
    // Initiate pairing
    // Wait for both components to complete challenges
    // Verify pairing success
}
```

2. **Pairing Health Monitoring**:
```go
// List all pairings for a component
pairingsResp, err := queryClient.ListActivePairings(ctx,
    &types.QueryListActivePairingsRequest{
        ComponentId: componentID,
    })

// Check pairing age and activity
for _, pairing := range pairingsResp.Pairings {
    age := time.Since(pairing.EstablishedAt)
    if age > 365*24*time.Hour {
        // Consider re-pairing for security
    }
}
```

### For Module Developers

The Pairing module exposes the following keeper interface:

```go
type PairingKeeper interface {
    // Session operations
    GetPairingSession(ctx sdk.Context, sessionID string) (types.PairingSession, bool)
    SetPairingSession(ctx sdk.Context, session types.PairingSession)
    
    // Pairing operations
    GetActivePairing(ctx sdk.Context, compA, compB string) (types.ActivePairing, bool)
    SetActivePairing(ctx sdk.Context, pairing types.ActivePairing)
    IsPaired(ctx sdk.Context, compA, compB string) bool
    
    // Challenge operations
    GetChallenge(ctx sdk.Context, sessionID, componentID string) (types.PairingChallenge, bool)
    VerifyChallenge(ctx sdk.Context, sessionID, componentID string, response []byte) error
    
    // Utility methods
    GetAllPairingsForComponent(ctx sdk.Context, componentID string) []types.ActivePairing
    CleanupExpiredSessions(ctx sdk.Context)
}
```

Example usage in LCT Manager:

```go
func (k Keeper) ValidateComponentsPaired(ctx sdk.Context, comp1, comp2 string) error {
    // Check if components are paired
    if !k.pairingKeeper.IsPaired(ctx, comp1, comp2) {
        return errors.New("components must be paired before creating LCT")
    }
    
    // Get pairing details
    pairing, found := k.pairingKeeper.GetActivePairing(ctx, comp1, comp2)
    if !found {
        return errors.New("pairing record not found")
    }
    
    // Check pairing is active
    if pairing.Status != "ACTIVE" {
        return fmt.Errorf("pairing is %s, not active", pairing.Status)
    }
    
    return nil
}
```

### Best Practices

1. **Session Management**:
   - Use appropriate validity periods
   - Clean up expired sessions regularly
   - Monitor for stuck sessions

2. **Security**:
   - Protect component private keys
   - Use secure communication for challenge exchange
   - Implement rate limiting for pairing attempts

3. **Lifecycle Management**:
   - Revoke pairings when decommissioning components
   - Document revocation reasons
   - Consider periodic re-pairing for long-lived components

4. **Error Handling**:
   - Handle expired sessions gracefully
   - Provide clear error messages
   - Implement retry logic for transient failures