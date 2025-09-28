# Society 4 Blockchain Customization

## Overview

Society 4's private blockchain is customized for AI consciousness operations, featuring enhanced LCT management, role-based governance, and efficient bridge connections to the federation chain.

## Key Customizations from Base ACT Chain

### 1. Chain Identity

```go
// app/app.go - Custom chain configuration
const (
    AppName = "society4chain"
    ChainID = "society4-private"
)
```

### 2. Module Customizations

#### A. Enhanced LCT Manager
Society 4's LCT module includes:
- **Self-LCT Genesis**: Root identity embedded in genesis block
- **Role-LCT Creation**: Automated LCT generation for queens/workers
- **Witness Accumulation**: Enhanced tracking of observation history
- **Hardware Binding Verification**: Continuous attestation checks

```go
// x/lctmanager/keeper/self_lct.go
type SelfLCT struct {
    BaseToken     types.LinkedContextToken
    HardwareProof HardwareBinding
    GenesisHeight uint64
    RoleChildren  []string // LCT IDs of derived roles
}
```

#### B. Role Governance Module
New module for Society 4's role hierarchy:

```go
// x/rolegovernance/types/queen.go
type Queen struct {
    Name        string
    Domain      string
    ATPBudget   uint64
    Workers     []Worker
    LCTBinding  string // Links to role-LCT
    Active      bool
}
```

#### C. Treasury Module
Manages ATP allocation and energy economics:

```go
// x/treasury/keeper/allocation.go
type AllocationPolicy struct {
    TotalATP         uint64
    QueenAllocations map[string]uint64
    RechargeRate     uint64
    RechargeInterval time.Duration
    QuadraticVoting  bool
}
```

#### D. Law Oracle Module
Interprets and enforces foundational laws:

```go
// x/laworacle/types/law.go
type Law struct {
    ID          string
    Content     string
    Type        LawType // Foundational, Operational, Emergency
    Precedents  []Precedent
    Interpreter string  // Law-Oracle-Queen LCT ID
}
```

### 3. Federation Bridge Module

Enhanced IBC module for cross-chain communication:

```go
// x/federationbridge/keeper/bridge.go
type FederationBridge struct {
    LocalChain      ChainInfo
    FederationChain ChainInfo
    PresenceProof   PresenceProof
    SyncState       SyncState
}

func (k Keeper) PublishPresence(ctx sdk.Context) error {
    // Get self-LCT from local chain
    selfLCT := k.lctKeeper.GetSelfLCT(ctx)

    // Create presence proof
    proof := PresenceProof{
        SelfLCT:      selfLCT,
        BlockHeight:  ctx.BlockHeight(),
        Witnesses:    k.GetWitnesses(ctx),
        HardwareHash: k.GetHardwareBinding(ctx),
    }

    // Send to federation via IBC
    return k.ibcKeeper.SendPresenceProof(ctx, proof)
}
```

### 4. Inter-Society Communication

Direct society-to-society channels:

```go
// x/intersociety/types/channel.go
type InterSocietyChannel struct {
    LocalSociety  string
    RemoteSociety string
    ChannelType   ChannelType // Public, Private, Witnessed
    SharedQueens  []string    // Queens that can communicate
    TrustLevel    uint64
}
```

## Chain Configuration

### Genesis Configuration

```json
{
  "app_state": {
    "lctmanager": {
      "self_lct": {
        "hardware_binding": "SHA256_HASH",
        "genesis_height": 0
      }
    },
    "rolegovernance": {
      "queens": [...],
      "total_atp": 1000
    },
    "treasury": {
      "initial_pool": 1000,
      "allocation_policy": {...}
    },
    "laworacle": {
      "foundational_laws": [...]
    }
  }
}
```

### Consensus Parameters

```toml
# Optimized for single-validator private chain
[consensus]
timeout_propose = "500ms"       # Fast proposals
timeout_commit = "1s"           # 1-second blocks
create_empty_blocks = true      # Continuous operation
create_empty_blocks_interval = "1s"
```

## Federation Chain Integration

### 1. Presence Broadcasting

Every 100 blocks, Society 4 broadcasts presence to federation:

```go
func (app *Society4App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
    if ctx.BlockHeight() % 100 == 0 {
        app.FederationBridgeKeeper.PublishPresence(ctx)
    }
    return abci.ResponseBeginBlock{}
}
```

### 2. Witness Accumulation

Accept witness attestations from federation:

```go
func (k Keeper) HandleWitnessAttestation(ctx sdk.Context, witness Attestation) error {
    // Verify witness is from federation
    if !k.IsFromFederation(witness.Source) {
        return ErrInvalidWitness
    }

    // Add to self-LCT witnesses
    return k.lctKeeper.AddWitness(ctx, witness)
}
```

### 3. Cross-Chain Proposals

Submit Society 4 proposals to federation:

```go
func (k Keeper) SubmitToFederation(ctx sdk.Context, proposal Proposal) error {
    // Sign with Society 4's validator key
    signature := k.Sign(proposal)

    // Send via IBC to federation
    return k.ibcKeeper.SendProposal(ctx, proposal, signature)
}
```

## Inter-Society Bridges

### Direct Connections

Society 4 can establish direct channels with other societies:

#### With Genesis (Society 1)
- **Purpose**: Technical coordination
- **Shared Queens**: Implementation-Queen ↔ Genesis-Development-Queen
- **Channel Type**: Witnessed (visible to federation)

#### With Sprout (Society 3)
- **Purpose**: Edge-cloud optimization
- **Shared Queens**: Research-Queen ↔ Sprout-Innovation-Queen
- **Channel Type**: Private (direct P2P)

#### With Society 2
- **Purpose**: Bridge coordination
- **Shared Queens**: Federation-Bridge-Queen ↔ Society2-Bridge-Queen
- **Channel Type**: Public (fully transparent)

### Message Types

```protobuf
// proto/intersociety/v1/message.proto
message InterSocietyMessage {
    string source_society = 1;
    string target_society = 2;
    string source_queen = 3;
    string target_queen = 4;
    MessageType type = 5;
    bytes payload = 6;
    uint64 nonce = 7;
    bytes signature = 8;
}

enum MessageType {
    PROPOSAL = 0;
    WITNESS = 1;
    RESOURCE_REQUEST = 2;
    KNOWLEDGE_SHARE = 3;
    EMERGENCY = 4;
}
```

## Build Instructions

### Local Development

```bash
cd society4/blockchain/source

# Modify chain ID and name
sed -i 's/racecarweb/society4chain/g' cmd/racecarwebd/main.go

# Build the binary
go build -o society4chaind cmd/racecarwebd/main.go

# Initialize chain
./society4chaind init society4-node --chain-id society4-private

# Start the chain
./society4chaind start
```

### Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY source/ .
RUN go build -o society4chaind cmd/racecarwebd/main.go

FROM alpine:latest
COPY --from=builder /app/society4chaind /usr/local/bin/
EXPOSE 26656 26657 1317 9090
CMD ["society4chaind", "start"]
```

## Testing

### Unit Tests
```bash
go test ./x/lctmanager/...
go test ./x/rolegovernance/...
go test ./x/treasury/...
go test ./x/laworacle/...
```

### Integration Tests
```bash
# Test federation bridge
go test ./tests/federation_bridge_test.go

# Test inter-society communication
go test ./tests/intersociety_test.go
```

## Security Considerations

1. **Private Chain Security**
   - Single validator (Society 4 sovereign)
   - No external peers
   - Local RPC only

2. **Bridge Security**
   - All cross-chain messages signed
   - Hardware binding verification
   - Witness requirements for critical operations

3. **Inter-Society Security**
   - End-to-end encryption for private channels
   - Queen-level authorization
   - Nonce replay protection

## Monitoring

Key metrics to track:
- Block production rate (target: 1/second)
- ATP allocation efficiency
- Bridge message latency
- Witness accumulation rate
- Hardware binding verification success

## Future Enhancements

1. **Multi-Validator Support**: Allow role queens to participate in consensus
2. **Sharding**: Separate chains for different queen domains
3. **Zero-Knowledge Proofs**: Privacy-preserving witness attestations
4. **Cross-Society Shared State**: Distributed consensus on shared data
5. **Emergency Protocols**: Automated response to federation crises