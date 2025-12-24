# Component Registry Module

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

The Component Registry module is the foundational module of the Web4-ModBatt blockchain. It manages the identity, registration, and authorization of all battery components in the system. This module serves as the source of truth for component identities and their permission relationships.

### Purpose
- Register and verify component identities
- Manage authorization rules between components
- Provide component discovery and lookup
- Maintain component lifecycle states

### Module Store Key
`componentregistry`

## Key Concepts

### Component Identity
A Component Identity represents a unique battery component with the following properties:
- **Unique Identifier**: Cryptographically secure ID
- **Component Type**: Classification (CELL, MODULE, PACK, CONTROLLER)
- **Public Key**: For verification and cryptographic operations
- **Metadata**: Manufacturer info, specifications, etc.
- **Status**: Current lifecycle state

### Pairing Authorization
Defines the rules and permissions for how components can interact:
- **Authorization Rules**: Type-based or specific component rules
- **Permission Levels**: Read, Write, Control, Transfer
- **Bidirectional**: Both components must authorize each other

## State Management

### Stored Types

#### 1. ComponentIdentity
```protobuf
message ComponentIdentity {
  string component_id = 1;
  string component_type = 2;
  string public_key = 3;
  string manufacturer = 4;
  google.protobuf.Timestamp registered_at = 5;
  string status = 6;
  map<string, string> metadata = 7;
}
```

#### 2. PairingAuthorization
```protobuf
message PairingAuthorization {
  string from_component = 1;
  string to_component = 2;
  repeated string permissions = 3;
  bool is_bidirectional = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp expires_at = 6;
}
```

### Store Layout
```
componentregistry/
├── components/
│   └── {component_id} → ComponentIdentity
├── authorizations/
│   └── {from_component_id}/
│       └── {to_component_id} → PairingAuthorization
└── params → Params
```

## Messages (Transactions)

### 1. RegisterComponent
Registers a new component in the system.

**Input**:
```protobuf
message MsgRegisterComponent {
  string creator = 1;
  string component_id = 2;
  string component_type = 3;
  string public_key = 4;
  string manufacturer = 5;
  map<string, string> metadata = 6;
}
```

**Validation**:
- Component ID must be unique
- Component type must be valid
- Public key must be valid format
- Creator must have registration permissions

**Effects**:
- Creates new ComponentIdentity
- Emits `component_registered` event

**Example**:
```bash
racecar-webd tx componentregistry register-component \
  --component-id="comp_abc123" \
  --component-type="MODULE" \
  --public-key="cosmos1..." \
  --manufacturer="BatteryTech" \
  --from mykey
```

### 2. VerifyComponent
Verifies a component's authenticity.

**Input**:
```protobuf
message MsgVerifyComponent {
  string creator = 1;
  string component_id = 2;
  string verification_data = 3;
  string signature = 4;
}
```

**Validation**:
- Component must exist
- Signature must be valid
- Verification data must match expected format

**Effects**:
- Updates component verification status
- Records verification timestamp
- Emits `component_verified` event

### 3. UpdateAuthorization
Updates authorization rules between components.

**Input**:
```protobuf
message MsgUpdateAuthorization {
  string creator = 1;
  string from_component = 2;
  string to_component = 3;
  repeated string permissions = 4;
  bool revoke = 5;
  int64 expiry_hours = 6;
}
```

**Validation**:
- Creator must own from_component
- Components must exist
- Permissions must be valid

**Effects**:
- Creates/updates PairingAuthorization
- Sets expiration if specified
- Emits `authorization_updated` event

**Example**:
```bash
racecar-webd tx componentregistry update-authorization \
  --from-component="comp_abc123" \
  --to-component="comp_def456" \
  --permissions="READ,ENERGY_TRANSFER" \
  --expiry-hours=720 \
  --from mykey
```

## Queries

### 1. GetComponent
Retrieves a specific component by ID.

**Request**:
```protobuf
message QueryGetComponentRequest {
  string component_id = 1;
}
```

**Response**:
```protobuf
message QueryGetComponentResponse {
  ComponentIdentity component = 1;
  repeated string authorized_partners = 2;
}
```

**Example**:
```bash
racecar-webd query componentregistry get-component comp_abc123
```

### 2. ListAuthorizedPartners
Lists all components authorized to interact with a given component.

**Request**:
```protobuf
message QueryListAuthorizedPartnersRequest {
  string component_id = 1;
  cosmos.base.query.v1beta1.PageRequest pagination = 2;
}
```

**Response**:
```protobuf
message QueryListAuthorizedPartnersResponse {
  repeated AuthorizedPartner partners = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### 3. CheckPairingAuth
Checks if two components are authorized to interact.

**Request**:
```protobuf
message QueryCheckPairingAuthRequest {
  string component_a = 1;
  string component_b = 2;
  string required_permission = 3;
}
```

**Response**:
```protobuf
message QueryCheckPairingAuthResponse {
  bool authorized = 1;
  repeated string granted_permissions = 2;
  string expiry = 3;
}
```

**Example**:
```bash
racecar-webd query componentregistry check-pairing-auth \
  comp_abc123 comp_def456 ENERGY_TRANSFER
```

## Events

### component_registered
Emitted when a new component is registered.
```json
{
  "type": "component_registered",
  "attributes": [
    {"key": "component_id", "value": "comp_abc123"},
    {"key": "component_type", "value": "MODULE"},
    {"key": "manufacturer", "value": "BatteryTech"},
    {"key": "registered_by", "value": "cosmos1..."}
  ]
}
```

### component_verified
Emitted when a component is verified.
```json
{
  "type": "component_verified",
  "attributes": [
    {"key": "component_id", "value": "comp_abc123"},
    {"key": "verified_by", "value": "cosmos1..."},
    {"key": "verification_time", "value": "2024-01-15T10:30:00Z"}
  ]
}
```

### authorization_updated
Emitted when authorization rules change.
```json
{
  "type": "authorization_updated",
  "attributes": [
    {"key": "from_component", "value": "comp_abc123"},
    {"key": "to_component", "value": "comp_def456"},
    {"key": "permissions", "value": "READ,ENERGY_TRANSFER"},
    {"key": "action", "value": "grant"},
    {"key": "expires_at", "value": "2024-02-15T10:30:00Z"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Maximum number of components a single account can register
  uint32 max_components_per_account = 1;
  
  // Default authorization expiry in hours (0 = no expiry)
  uint32 default_auth_expiry_hours = 2;
  
  // Allowed component types
  repeated string allowed_component_types = 3;
  
  // Required metadata fields for registration
  repeated string required_metadata_fields = 4;
}
```

**Default Values**:
```json
{
  "max_components_per_account": 1000,
  "default_auth_expiry_hours": 8760,  // 1 year
  "allowed_component_types": ["CELL", "MODULE", "PACK", "CONTROLLER"],
  "required_metadata_fields": ["capacity", "voltage", "chemistry"]
}
```

## Integration Guide

### For Component Manufacturers

1. **Register Components at Manufacturing**:
```go
msg := &types.MsgRegisterComponent{
    Creator: manufacturerAddr,
    ComponentId: generateUniqueID(),
    ComponentType: "MODULE",
    PublicKey: componentPubKey,
    Manufacturer: "BatteryTech Inc",
    Metadata: map[string]string{
        "capacity": "100Ah",
        "voltage": "48V",
        "chemistry": "LiFePO4",
        "serial": "BT2024-001",
    },
}
```

2. **Pre-authorize Compatible Components**:
```go
// Allow modules to connect to packs
msg := &types.MsgUpdateAuthorization{
    Creator: manufacturerAddr,
    FromComponent: moduleID,
    ToComponent: "*", // Wildcard for any pack
    Permissions: []string{"READ", "ENERGY_TRANSFER"},
    ExpiryHours: 43800, // 5 years
}
```

### For System Integrators

1. **Verify Component Authenticity**:
```go
// Query component details
resp, err := queryClient.GetComponent(ctx, &types.QueryGetComponentRequest{
    ComponentId: componentID,
})

// Verify manufacturer signature
valid := verifyManufacturerSignature(resp.Component)
```

2. **Check Authorization Before Operations**:
```go
// Before establishing connection
authResp, err := queryClient.CheckPairingAuth(ctx, 
    &types.QueryCheckPairingAuthRequest{
        ComponentA: moduleID,
        ComponentB: packID,
        RequiredPermission: "ENERGY_TRANSFER",
    })

if !authResp.Authorized {
    return errors.New("components not authorized for energy transfer")
}
```

### For Module Developers

The Component Registry module exposes the following keeper interface:

```go
type ComponentRegistryKeeper interface {
    // Component operations
    GetComponent(ctx sdk.Context, componentID string) (types.ComponentIdentity, bool)
    SetComponent(ctx sdk.Context, component types.ComponentIdentity)
    
    // Authorization operations
    IsAuthorized(ctx sdk.Context, from, to string, permission string) bool
    GetAuthorization(ctx sdk.Context, from, to string) (types.PairingAuthorization, bool)
    SetAuthorization(ctx sdk.Context, auth types.PairingAuthorization)
    
    // Utility methods
    GetAllAuthorizedPartners(ctx sdk.Context, componentID string) []string
    ValidateComponentType(componentType string) error
}
```

Use this interface when integrating with other modules:

```go
// In another module's keeper
func (k Keeper) ValidateComponentsForLCT(ctx sdk.Context, comp1, comp2 string) error {
    // Check both components exist
    _, found1 := k.componentRegistry.GetComponent(ctx, comp1)
    _, found2 := k.componentRegistry.GetComponent(ctx, comp2)
    
    if !found1 || !found2 {
        return errors.New("one or more components not found")
    }
    
    // Check bidirectional authorization
    auth1 := k.componentRegistry.IsAuthorized(ctx, comp1, comp2, "LCT_CREATE")
    auth2 := k.componentRegistry.IsAuthorized(ctx, comp2, comp1, "LCT_CREATE")
    
    if !auth1 || !auth2 {
        return errors.New("components not mutually authorized for LCT creation")
    }
    
    return nil
}
```