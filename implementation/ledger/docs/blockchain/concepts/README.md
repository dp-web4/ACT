# Web4-ModBatt Core Concepts

## Table of Contents

1. [Overview](#overview)
2. [Component Identity](#component-identity)
3. [Linked Context Tokens (LCTs)](#linked-context-tokens-lcts)
4. [Trust Tensors](#trust-tensors)
5. [Energy Operations](#energy-operations)
6. [Pairing Mechanisms](#pairing-mechanisms)
7. [Authorization Framework](#authorization-framework)
8. [Attention Tokens (ATP/ADP)](#attention-tokens-atpadp)

## Overview

The Web4-ModBatt blockchain introduces several novel concepts designed specifically for managing modular battery systems in a decentralized environment. These concepts form the foundation of how components interact, build trust, and exchange energy.

## Component Identity

### Definition
A Component Identity represents a unique battery component (cell, module, or pack) registered on the blockchain. Each component has a cryptographically secure identity that cannot be forged or duplicated.

### Key Attributes
- **Component ID**: Unique identifier (typically a hash or UUID)
- **Component Type**: Classification (CELL, MODULE, PACK, CONTROLLER)
- **Manufacturer Data**: Origin information and specifications
- **Public Key**: For cryptographic operations and verification
- **Registration Timestamp**: When the component joined the network
- **Status**: Active, Inactive, Decommissioned

### Example Structure
```json
{
  "component_id": "comp_abc123def456",
  "component_type": "MODULE",
  "manufacturer": "BatteryTech Inc",
  "public_key": "cosmos1abc...",
  "registered_at": "2024-01-15T10:30:00Z",
  "status": "ACTIVE",
  "specifications": {
    "capacity": "100Ah",
    "voltage": "48V",
    "chemistry": "LiFePO4"
  }
}
```

## Linked Context Tokens (LCTs)

### Definition
LCTs are blockchain-native tokens that represent relationships between components. They serve as the primary mechanism for establishing, managing, and tracking component interactions.

### Purpose
- **Relationship Representation**: Each LCT embodies a specific relationship between two or more components
- **Permission Management**: LCTs define what operations are allowed between related components
- **State Tracking**: Maintain the current state of component relationships
- **History Recording**: Immutable log of all relationship changes

### LCT Lifecycle
```
Creation → Active → Updated → Terminated
    │         │        │          │
    └─────────┴────────┴──────────┴──→ Historical Record
```

### Key Properties
- **LCT ID**: Unique identifier for the relationship
- **Participant Components**: List of components in the relationship
- **Relationship Type**: PEER, PARENT_CHILD, CONTROLLER_DEVICE
- **Permissions**: Allowed operations (READ, WRITE, ENERGY_TRANSFER)
- **Status**: PENDING, ACTIVE, SUSPENDED, TERMINATED
- **Creation Context**: Why and how the relationship was established

## Trust Tensors

### Definition
Trust Tensors are multi-dimensional mathematical structures that quantify and track trust relationships between components over time. They use tensor mathematics to capture complex trust dynamics.

### Components of Trust Tensors
1. **Base Trust Score**: Initial trust value (0.0 to 1.0)
2. **Dimensional Factors**:
   - Reliability: Historical performance
   - Authenticity: Verification status
   - Behavior: Operational patterns
   - Longevity: Time in network
3. **Witness Confirmations**: Third-party validations
4. **Decay Function**: Trust degradation over time without interaction

### Trust Calculation
```
Trust = Σ(dimension_weight × dimension_score × witness_factor) × decay_factor
```

### Trust Evolution
Trust scores evolve based on:
- Successful operations (+)
- Failed operations (-)
- Time without interaction (decay)
- Witness testimonials (±)
- Anomalous behavior (-)

## Energy Operations

### Definition
Energy Operations represent the transfer, allocation, or consumption of energy between components. These operations are tracked on-chain for transparency and accountability.

### Operation Types
1. **Direct Transfer**: Energy moves from one component to another
2. **Pooled Distribution**: Energy distributed among multiple components
3. **Reserved Allocation**: Energy reserved for future use
4. **Emergency Discharge**: Rapid energy release for critical needs

### Energy Flow Model
```
Source Component → [Validation] → [Trust Check] → [Transfer] → Target Component
                        │              │                │
                        ▼              ▼                ▼
                  Authorization    Trust Score    Update Balances
```

### Key Attributes
- **Operation ID**: Unique transaction identifier
- **Source LCT**: Relationship enabling the transfer
- **Energy Amount**: Quantity in watt-hours (Wh)
- **Transfer Rate**: Power in watts (W)
- **Validation Status**: Pending, Approved, Completed, Failed
- **Trust Requirement**: Minimum trust score needed

## Pairing Mechanisms

### Definition
Pairing is the process by which components establish secure, authenticated relationships. The system supports multiple pairing modes for different scenarios.

### Pairing Types

#### 1. Bidirectional Pairing
Both components actively participate in authentication:
```
Component A ←─[Challenge]─→ Component B
     │                           │
     └──[Response + Verify]──────┘
```

#### 2. Offline Pairing
For components not currently online:
- Request queued on blockchain
- Processed when target comes online
- Supports proxy authorization

#### 3. Factory Pairing
Pre-authorized at manufacturing:
- Components ship with pairing credentials
- Automatic pairing on first connection
- Reduced setup complexity

### Pairing Security
- **Challenge-Response**: Cryptographic proof of identity
- **Time-Bound Sessions**: Pairing attempts expire
- **Mutual Authentication**: Both parties verify each other
- **Revocation Support**: Pairings can be terminated

## Authorization Framework

### Definition
The Authorization Framework defines rules and permissions governing component interactions. It ensures only authorized operations occur between components.

### Authorization Levels

1. **Read-Only**: Access component data and status
2. **Operational**: Control component behavior
3. **Energy Transfer**: Move energy between components
4. **Administrative**: Update configurations and permissions
5. **Emergency**: Override normal restrictions in critical situations

### Rule Structure
```json
{
  "rule_id": "auth_rule_001",
  "from_component_type": "CONTROLLER",
  "to_component_type": "MODULE",
  "permissions": ["READ", "OPERATIONAL", "ENERGY_TRANSFER"],
  "conditions": {
    "min_trust_score": 0.7,
    "requires_pairing": true,
    "time_restrictions": "business_hours"
  }
}
```

### Dynamic Authorization
Permissions can change based on:
- Trust scores
- Operational context
- Time of day
- Emergency status
- Network conditions

## Attention Tokens (ATP/ADP)

### Definition
Attention Tokens represent the computational and operational "attention" that components pay to relationships and operations.

### Token Types

#### ATP (Attention Transfer Protocol)
- **Purpose**: Direct attention allocation between components
- **Use Cases**:
  - Priority processing requests
  - Dedicated monitoring
  - Exclusive operations

#### ADP (Attention Distribution Protocol)
- **Purpose**: Distributed attention across multiple relationships
- **Use Cases**:
  - Load balancing
  - Fair resource allocation
  - Network-wide operations

### Attention Economy
```
Component Resources = Base Capacity - Σ(Allocated Attention)

Where:
- Base Capacity = Component's total attention budget
- Allocated Attention = ATP + ADP commitments
```

### Key Properties
- **Fungibility**: Attention tokens can be exchanged
- **Time-Decay**: Unused attention expires
- **Priority Levels**: Higher attention = faster processing
- **Resource Protection**: Prevents attention exhaustion

### Example Usage
```json
{
  "attention_allocation": {
    "total_capacity": 1000,
    "allocated": {
      "atp_transfers": 300,
      "adp_distributions": 200,
      "system_reserved": 100
    },
    "available": 400
  }
}
```

## Concept Interactions

These concepts work together to create a comprehensive battery management system:

1. **Components** establish **Pairings**
2. **Pairings** create **LCTs**
3. **LCTs** enable **Energy Operations**
4. **Operations** build **Trust**
5. **Trust** influences **Authorization**
6. **Authorization** governs future **Operations**
7. **Attention Tokens** prioritize **Processing**

This creates a self-reinforcing system where good behavior increases capabilities while protecting against malicious actors.