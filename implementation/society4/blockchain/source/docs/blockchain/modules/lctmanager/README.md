# LCT Manager Module

## Table of Contents

1. [Overview](#overview)
2. [Key Concepts](#key-concepts)
3. [State Management](#state-management)
4. [Messages (Transactions)](#messages-transactions)
5. [Queries](#queries)
6. [Events](#events)
7. [Parameters](#parameters)
8. [Integration Guide](#integration-guide)

## Overview

The LCT Manager module is responsible for creating and managing Linked Context Tokens (LCTs), which represent relationships between components in the Web4-ModBatt system. LCTs are the fundamental building blocks that enable secure interactions, energy transfers, and trust building between battery components.

### Purpose
- Create and manage component relationships through LCTs
- Validate access permissions for LCT operations
- Track relationship lifecycle and status changes
- Maintain component relationship mappings

### Dependencies
- **Component Registry**: Verifies component identities and authorizations

### Module Store Key
`lctmanager`

## Key Concepts

### Linked Context Token (LCT)
An LCT is a blockchain-native token representing a relationship between two or more components:
- **Unique Identifier**: Each LCT has a globally unique ID
- **Participants**: Components involved in the relationship
- **Relationship Type**: Defines the nature of the relationship
- **Permissions**: Operations allowed within the relationship
- **Status**: Current state of the relationship

### Component Relationship
A mapping that tracks all LCTs associated with a specific component:
- **Component ID**: The component being tracked
- **LCT List**: All LCTs where this component participates
- **Role**: The component's role in each relationship

### Relationship Types
- **PEER**: Equal partnership between components
- **PARENT_CHILD**: Hierarchical relationship (e.g., pack-module)
- **CONTROLLER_DEVICE**: Control relationship
- **ENERGY_SHARING**: Dedicated energy exchange relationship

## State Management

### Stored Types

#### 1. LinkedContextToken
```protobuf
message LinkedContextToken {
  string lct_id = 1;
  repeated string participant_components = 2;
  string relationship_type = 3;
  repeated string permissions = 4;
  string status = 5;  // PENDING, ACTIVE, SUSPENDED, TERMINATED
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
  map<string, string> metadata = 8;
  string created_by = 9;
}
```

#### 2. ComponentRelationship
```protobuf
message ComponentRelationship {
  string component_id = 1;
  repeated LCTReference lct_refs = 2;
  uint32 active_relationship_count = 3;
  google.protobuf.Timestamp last_activity = 4;
}

message LCTReference {
  string lct_id = 1;
  string role = 2;  // INITIATOR, PARTICIPANT, OBSERVER
  string relationship_type = 3;
  string status = 4;
}
```

### Store Layout
```
lctmanager/
├── lcts/
│   └── {lct_id} → LinkedContextToken
├── relationships/
│   └── {component_id}/
│       └── index → ComponentRelationship
└── params → Params
```

## Messages (Transactions)

### 1. CreateLCTRelationship
Creates a new LCT between components.

**Input**:
```protobuf
message MsgCreateLCTRelationship {
  string creator = 1;
  repeated string participant_components = 2;
  string relationship_type = 3;
  repeated string permissions = 4;
  map<string, string> metadata = 5;
}
```

**Validation**:
- All components must exist in Component Registry
- Components must have mutual authorization
- Relationship type must be valid
- Creator must own at least one participant component
- No duplicate active relationships of same type between components

**Effects**:
- Creates new LinkedContextToken
- Updates ComponentRelationship for all participants
- Emits `lct_created` event

**Example**:
```bash
racecar-webd tx lctmanager create-lct-relationship \
  --participants="comp_module1,comp_pack1" \
  --relationship-type="PARENT_CHILD" \
  --permissions="READ,WRITE,ENERGY_TRANSFER" \
  --from mykey
```

### 2. UpdateLCTStatus
Updates the status of an existing LCT.

**Input**:
```protobuf
message MsgUpdateLCTStatus {
  string creator = 1;
  string lct_id = 2;
  string new_status = 3;  // ACTIVE, SUSPENDED, TERMINATED
  string reason = 4;
}
```

**Validation**:
- LCT must exist
- Creator must be participant in LCT
- Status transition must be valid
- Cannot reactivate terminated LCTs

**Effects**:
- Updates LCT status
- Updates timestamp
- Updates component relationship counts
- Emits `lct_status_updated` event

**Status Transitions**:
```
PENDING → ACTIVE → SUSPENDED → TERMINATED
   ↓        ↓         ↓
   └────────┴─────────┴────────→ TERMINATED
```

### 3. TerminateLCTRelationship
Terminates an LCT relationship.

**Input**:
```protobuf
message MsgTerminateLCTRelationship {
  string creator = 1;
  string lct_id = 2;
  string termination_reason = 3;
}
```

**Validation**:
- LCT must exist and be active/suspended
- Creator must be participant or have admin rights
- Must provide termination reason

**Effects**:
- Sets LCT status to TERMINATED
- Updates component relationships
- Prevents future operations on this LCT
- Emits `lct_terminated` event

## Queries

### 1. GetLCT
Retrieves a specific LCT by ID.

**Request**:
```protobuf
message QueryGetLCTRequest {
  string lct_id = 1;
}
```

**Response**:
```protobuf
message QueryGetLCTResponse {
  LinkedContextToken lct = 1;
  repeated ComponentInfo participants = 2;
}
```

**Example**:
```bash
racecar-webd query lctmanager get-lct lct_abc123def456
```

### 2. GetComponentRelationships
Gets all relationships for a specific component.

**Request**:
```protobuf
message QueryGetComponentRelationshipsRequest {
  string component_id = 1;
  string status_filter = 2;  // Optional: ACTIVE, SUSPENDED, etc.
  cosmos.base.query.v1beta1.PageRequest pagination = 3;
}
```

**Response**:
```protobuf
message QueryGetComponentRelationshipsResponse {
  repeated LinkedContextToken lcts = 1;
  ComponentRelationship relationship_summary = 2;
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

### 3. ValidateLCTAccess
Validates if a component has access to perform operations on an LCT.

**Request**:
```protobuf
message QueryValidateLCTAccessRequest {
  string lct_id = 1;
  string component_id = 2;
  string required_permission = 3;
}
```

**Response**:
```protobuf
message QueryValidateLCTAccessResponse {
  bool has_access = 1;
  string role = 2;
  repeated string granted_permissions = 3;
  string lct_status = 4;
}
```

**Example**:
```bash
racecar-webd query lctmanager validate-lct-access \
  lct_abc123 comp_module1 ENERGY_TRANSFER
```

## Events

### lct_created
Emitted when a new LCT is created.
```json
{
  "type": "lct_created",
  "attributes": [
    {"key": "lct_id", "value": "lct_abc123def456"},
    {"key": "relationship_type", "value": "PARENT_CHILD"},
    {"key": "participants", "value": "comp_module1,comp_pack1"},
    {"key": "creator", "value": "cosmos1..."},
    {"key": "permissions", "value": "READ,WRITE,ENERGY_TRANSFER"}
  ]
}
```

### lct_status_updated
Emitted when LCT status changes.
```json
{
  "type": "lct_status_updated",
  "attributes": [
    {"key": "lct_id", "value": "lct_abc123def456"},
    {"key": "old_status", "value": "ACTIVE"},
    {"key": "new_status", "value": "SUSPENDED"},
    {"key": "reason", "value": "Maintenance mode"},
    {"key": "updated_by", "value": "cosmos1..."}
  ]
}
```

### lct_terminated
Emitted when an LCT is terminated.
```json
{
  "type": "lct_terminated",
  "attributes": [
    {"key": "lct_id", "value": "lct_abc123def456"},
    {"key": "termination_reason", "value": "Component decommissioned"},
    {"key": "terminated_by", "value": "cosmos1..."},
    {"key": "active_duration_hours", "value": "8760"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Maximum LCTs per component
  uint32 max_lcts_per_component = 1;
  
  // Maximum participants in a single LCT
  uint32 max_participants_per_lct = 2;
  
  // Allowed relationship types
  repeated string allowed_relationship_types = 3;
  
  // Default permissions for new LCTs
  repeated string default_permissions = 4;
  
  // Require mutual authorization for LCT creation
  bool require_mutual_authorization = 5;
  
  // Auto-terminate inactive LCTs after days (0 = disabled)
  uint32 auto_terminate_days = 6;
}
```

**Default Values**:
```json
{
  "max_lcts_per_component": 100,
  "max_participants_per_lct": 10,
  "allowed_relationship_types": [
    "PEER", 
    "PARENT_CHILD", 
    "CONTROLLER_DEVICE", 
    "ENERGY_SHARING"
  ],
  "default_permissions": ["READ"],
  "require_mutual_authorization": true,
  "auto_terminate_days": 0
}
```

## Integration Guide

### For Application Developers

1. **Creating Relationships Between Components**:
```go
// First ensure components are paired and authorized
// Then create LCT
msg := &types.MsgCreateLCTRelationship{
    Creator: controllerAddr,
    ParticipantComponents: []string{moduleID, packID},
    RelationshipType: "PARENT_CHILD",
    Permissions: []string{"READ", "WRITE", "ENERGY_TRANSFER"},
    Metadata: map[string]string{
        "pack_position": "A1",
        "connection_type": "series",
        "max_current": "100A",
    },
}
```

2. **Monitoring Relationship Health**:
```go
// Query all active relationships for a component
resp, err := queryClient.GetComponentRelationships(ctx, 
    &types.QueryGetComponentRelationshipsRequest{
        ComponentId: moduleID,
        StatusFilter: "ACTIVE",
    })

// Check each LCT's last activity
for _, lct := range resp.Lcts {
    if time.Since(lct.UpdatedAt) > 24*time.Hour {
        // Consider suspending inactive relationships
    }
}
```

### For Energy System Integrators

1. **Validate Before Energy Operations**:
```go
// Before initiating energy transfer
accessResp, err := queryClient.ValidateLCTAccess(ctx,
    &types.QueryValidateLCTAccessRequest{
        LctId: lctID,
        ComponentId: sourceModule,
        RequiredPermission: "ENERGY_TRANSFER",
    })

if !accessResp.HasAccess {
    return fmt.Errorf("no energy transfer permission on LCT %s", lctID)
}

if accessResp.LctStatus != "ACTIVE" {
    return fmt.Errorf("LCT not active: %s", accessResp.LctStatus)
}
```

2. **Handle Relationship Lifecycle**:
```go
// Suspend during maintenance
suspendMsg := &types.MsgUpdateLCTStatus{
    Creator: maintainerAddr,
    LctId: lctID,
    NewStatus: "SUSPENDED",
    Reason: "Scheduled maintenance",
}

// Reactivate after maintenance
activateMsg := &types.MsgUpdateLCTStatus{
    Creator: maintainerAddr,
    LctId: lctID,
    NewStatus: "ACTIVE",
    Reason: "Maintenance completed",
}
```

### For Module Developers

The LCT Manager exposes the following keeper interface:

```go
type LCTManagerKeeper interface {
    // LCT operations
    GetLCT(ctx sdk.Context, lctID string) (types.LinkedContextToken, bool)
    SetLCT(ctx sdk.Context, lct types.LinkedContextToken)
    
    // Relationship operations
    GetComponentRelationships(ctx sdk.Context, componentID string) (types.ComponentRelationship, bool)
    UpdateComponentRelationship(ctx sdk.Context, componentID string, lctRef types.LCTReference)
    
    // Validation methods
    ValidateLCTAccess(ctx sdk.Context, lctID, componentID, permission string) bool
    IsLCTActive(ctx sdk.Context, lctID string) bool
    
    // Utility methods
    GetAllActiveLCTs(ctx sdk.Context) []types.LinkedContextToken
    GetLCTsByType(ctx sdk.Context, relType string) []types.LinkedContextToken
}
```

Example integration in Energy Cycle module:

```go
func (k Keeper) ValidateEnergyTransfer(ctx sdk.Context, lctID string, source, target string) error {
    // Get LCT details
    lct, found := k.lctManager.GetLCT(ctx, lctID)
    if !found {
        return errors.New("LCT not found")
    }
    
    // Verify both components are participants
    isSourceParticipant := false
    isTargetParticipant := false
    
    for _, participant := range lct.ParticipantComponents {
        if participant == source {
            isSourceParticipant = true
        }
        if participant == target {
            isTargetParticipant = true
        }
    }
    
    if !isSourceParticipant || !isTargetParticipant {
        return errors.New("components not part of LCT relationship")
    }
    
    // Check LCT has energy transfer permission
    hasPermission := false
    for _, perm := range lct.Permissions {
        if perm == "ENERGY_TRANSFER" {
            hasPermission = true
            break
        }
    }
    
    if !hasPermission {
        return errors.New("LCT does not allow energy transfers")
    }
    
    // Verify LCT is active
    if !k.lctManager.IsLCTActive(ctx, lctID) {
        return errors.New("LCT is not active")
    }
    
    return nil
}
```

### Best Practices

1. **Relationship Design**:
   - Use specific relationship types for different use cases
   - Include relevant metadata for relationship context
   - Set appropriate permissions based on relationship type

2. **Lifecycle Management**:
   - Regularly check relationship activity
   - Suspend inactive relationships to save resources
   - Terminate relationships when components are decommissioned

3. **Permission Management**:
   - Follow principle of least privilege
   - Only grant necessary permissions
   - Regularly audit permission usage

4. **Error Handling**:
   - Always check LCT status before operations
   - Handle all validation errors gracefully
   - Provide meaningful error messages