# Web4-ModBatt Blockchain Architecture

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Module Architecture](#module-architecture)
3. [Data Flow](#data-flow)
4. [State Management](#state-management)
5. [Transaction Flow](#transaction-flow)
6. [Query System](#query-system)
7. [Inter-Module Communication](#inter-module-communication)

## System Architecture

### High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│  (Web UI, CLI, Mobile Apps, Embedded Devices)                   │
└────────────────┬─────────────────────────┬──────────────────────┘
                 │                         │
                 ▼                         ▼
┌────────────────────────┐       ┌────────────────────────┐
│      REST API          │       │      gRPC Gateway      │
│  (OpenAPI/Swagger)     │       │   (Protocol Buffers)   │
└────────────┬───────────┘       └───────────┬────────────┘
             │                               │
             ▼                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Application Layer (app/)                      │
│  ┌────────────┬─────────────┬─────────────┬─────────────────┐  │
│  │   Router   │   Keeper    │   Handler   │   Query Server  │  │
│  └────────────┴─────────────┴─────────────┴─────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────┐
│                         Module Layer (x/)                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Component  │    LCT     │  Pairing  │ Pairing Queue   │   │
│  │  Registry   │  Manager   │           │                  │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │   Energy    │   Trust    │           │                  │   │
│  │   Cycle     │   Tensor   │           │                  │   │
│  └─────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────┐
│                      Cosmos SDK Core                             │
│  ┌─────────────┬──────────────┬──────────────┬──────────────┐  │
│  │    Bank     │    Staking   │     Gov      │     Auth     │  │
│  └─────────────┴──────────────┴──────────────┴──────────────┘  │
└──────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────┐
│                   Consensus Layer (CometBFT)                     │
│  ┌─────────────┬──────────────┬──────────────┬──────────────┐  │
│  │  Consensus  │  Networking  │   Mempool    │     P2P      │  │
│  └─────────────┴──────────────┴──────────────┴──────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Network Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Validator  │◄────►  Validator  │◄────►  Validator  │
│    Node 1   │     │    Node 2   │     │    Node 3   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌──────────────────────────────────────────────────────┐
│              Tendermint P2P Network                  │
└──────────────────────────────────────────────────────┘
       ▲                   ▲                   ▲
       │                   │                   │
┌──────┴──────┐     ┌──────┴──────┐     ┌──────┴──────┐
│  Full Node  │     │  Full Node  │     │  Light      │
│      1      │     │      2      │     │  Client     │
└─────────────┘     └─────────────┘     └─────────────┘
```

## Module Architecture

### Standard Module Structure

Each module follows the Cosmos SDK module pattern:

```
x/<module_name>/
├── keeper/              # State management and business logic
│   ├── keeper.go       # Keeper struct and constructor
│   ├── msg_server.go   # Message handler implementations
│   ├── query.go        # Query handler implementations
│   └── genesis.go      # Genesis state handling
├── types/              # Types and interfaces
│   ├── codec.go        # Codec registration
│   ├── errors.go       # Module-specific errors
│   ├── keys.go         # Store keys and prefixes
│   ├── msgs.go         # Message type implementations
│   └── *.pb.go         # Generated protobuf code
└── module/             # Module integration
    ├── module.go       # Module implementation
    ├── autocli.go      # CLI command registration
    └── depinject.go    # Dependency injection
```

### Keeper Pattern

Each module uses a Keeper for state management:

```go
type Keeper struct {
    cdc          codec.BinaryCodec
    storeKey     storetypes.StoreKey
    memKey       storetypes.StoreKey
    paramstore   paramtypes.Subspace
    
    // Inter-module dependencies
    bankKeeper   types.BankKeeper
    accountKeeper types.AccountKeeper
    // Module-specific keepers...
}
```

## Data Flow

### Component Registration Flow

```
Client Request
     │
     ▼
REST API/gRPC
     │
     ▼
Tx Handler ──────► Component Registry Module
     │                      │
     │                      ▼
     │              Validate Component
     │                      │
     │                      ▼
     │              Store Component Identity
     │                      │
     │                      ▼
     └──────────────  Emit Events
                           │
                           ▼
                    Update State Tree
```

### Energy Transfer Flow

```
Initiate Transfer
     │
     ▼
Validate LCT ◄──── LCT Manager
     │
     ▼
Check Trust ◄──── Trust Tensor
     │
     ▼
Execute Transfer ◄─── Energy Cycle
     │
     ▼
Update Balances
     │
     ▼
Record History
```

## State Management

### Store Structure

Each module maintains its state in a dedicated store:

```
Module Store Keys:
├── ComponentRegistry
│   ├── components/<component_id>
│   └── authorizations/<component_id>/<partner_id>
├── LCTManager
│   ├── lcts/<lct_id>
│   └── relationships/<component_id>/<index>
├── Pairing
│   ├── sessions/<session_id>
│   └── pairings/<component_id>/<partner_id>
├── PairingQueue
│   ├── requests/<request_id>
│   └── queues/<component_id>/<index>
├── EnergyCycle
│   ├── operations/<operation_id>
│   ├── balances/<lct_id>
│   └── history/<lct_id>/<index>
└── TrustTensor
    ├── tensors/<relationship_id>
    └── witnesses/<tensor_id>/<witness_id>
```

### State Transitions

State changes follow ACID properties through CometBFT consensus:

1. **Atomicity**: All state changes in a transaction succeed or fail together
2. **Consistency**: Validation ensures state remains consistent
3. **Isolation**: Transactions are processed sequentially
4. **Durability**: Committed state is persisted to disk

## Transaction Flow

### Message Processing Pipeline

```
1. Client constructs and signs transaction
                │
                ▼
2. Transaction submitted to node
                │
                ▼
3. Basic validation (signatures, fees)
                │
                ▼
4. Added to mempool
                │
                ▼
5. Included in block proposal
                │
                ▼
6. Consensus on block
                │
                ▼
7. Execute transaction
                │
                ▼
8. Update state
                │
                ▼
9. Emit events
                │
                ▼
10. Return result to client
```

## Query System

### Query Types

1. **State Queries**: Direct state lookups
   - Get component by ID
   - Get LCT details
   - Get trust score

2. **List Queries**: Paginated lists
   - List authorized partners
   - List active pairings
   - List energy history

3. **Computed Queries**: Calculated results
   - Calculate relationship trust
   - Calculate energy balance
   - Validate permissions

### Query Processing

```
Query Request
     │
     ▼
gRPC Handler
     │
     ▼
Query Server ───► Module Keeper
     │                 │
     │                 ▼
     │           Read State
     │                 │
     │                 ▼
     └───────── Return Result
```

## Inter-Module Communication

### Dependency Graph

```
┌─────────────────┐
│ Component       │
│ Registry        │◄─────────────────┐
└────────┬────────┘                  │
         │                           │
         ▼                           │
┌─────────────────┐         ┌────────────────┐
│ LCT Manager     │◄────────┤ Pairing Queue  │
└────────┬────────┘         └────────────────┘
         │                           ▲
         ▼                           │
┌─────────────────┐         ┌────────────────┐
│ Energy Cycle    │         │ Pairing        │
└────────┬────────┘         └────────────────┘
         │
         ▼
┌─────────────────┐
│ Trust Tensor    │
└─────────────────┘
```

### Interface Contracts

Modules communicate through well-defined keeper interfaces:

```go
// Example: LCTManager expects from ComponentRegistry
type ComponentRegistryKeeper interface {
    GetComponent(ctx sdk.Context, componentID string) (Component, bool)
    IsAuthorized(ctx sdk.Context, componentID, partnerID string) bool
}

// Example: Energyycle expects from LCTManager
type LCTManagerKeeper interface {
    GetLCT(ctx sdk.Context, lctID string) (LCT, bool)
    ValidateLCTAccess(ctx sdk.Context, lctID, componentID string) bool
}
```

### Event System

Modules emit events for asynchronous communication:

```go
ctx.EventManager().EmitEvent(
    sdk.NewEvent(
        "component_registered",
        sdk.NewAttribute("component_id", componentID),
        sdk.NewAttribute("component_type", componentType),
    ),
)
```

Other modules or external systems can subscribe to these events for reactive processing.