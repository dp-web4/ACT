# Federation Bridge Architecture

## Overview

The Federation Bridge enables Society 4's private blockchain to interact with the main federation chain while maintaining sovereignty and privacy. This document details the bridge architecture, protocols, and implementation.

## Architecture Diagram

```
┌─────────────────────────────────┐
│    Society 4 Private Chain       │
│         (Sovereign)              │
│                                  │
│  ┌────────────────────────┐     │
│  │   Self-LCT (Root)       │     │
│  └────────────────────────┘     │
│            ↓                     │
│  ┌────────────────────────┐     │
│  │   Federation Bridge     │     │
│  │      Module             │     │
│  └────────────────────────┘     │
└──────────────┬──────────────────┘
               │ IBC Protocol
               │ (Cosmos SDK)
               ↓
┌─────────────────────────────────┐
│    Federation Chain              │
│    (Public Ledger)               │
│                                  │
│  ┌────────────────────────┐     │
│  │  Society 4 Presence    │     │
│  │    Footprint           │     │
│  └────────────────────────┘     │
│                                  │
│  Other Society Footprints:      │
│  • Genesis                       │
│  • Sprout                        │
│  • Society 2                     │
└─────────────────────────────────┘
```

## Bridge Components

### 1. Presence Publisher

Broadcasts Society 4's existence proof to federation:

```go
package federationbridge

type PresencePublisher struct {
    localChainID    string
    federationChainID string
    selfLCT         *types.SelfLCT
    hardwareBinding *types.HardwareBinding
}

func (pp *PresencePublisher) PublishPresence(ctx sdk.Context) error {
    proof := types.PresenceProof{
        SocietyID:    "society4",
        ChainID:      pp.localChainID,
        SelfLCTHash:  pp.selfLCT.Hash(),
        HardwareHash: pp.hardwareBinding.Hash(),
        BlockHeight:  ctx.BlockHeight(),
        Timestamp:    ctx.BlockTime(),
        Signature:    pp.Sign(ctx),
    }

    return pp.sendToFederation(proof)
}
```

### 2. Witness Collector

Accumulates attestations from other societies:

```go
type WitnessCollector struct {
    witnesses map[string]types.Witness  // societyID -> witness
    threshold uint64                    // minimum witnesses required
}

func (wc *WitnessCollector) ProcessWitness(witness types.Witness) error {
    // Verify witness signature
    if !wc.verifySignature(witness) {
        return ErrInvalidSignature
    }

    // Store witness
    wc.witnesses[witness.SocietyID] = witness

    // Check if threshold reached
    if len(wc.witnesses) >= int(wc.threshold) {
        wc.updateTrustLevel()
    }

    return nil
}
```

### 3. Proposal Relay

Submits Society 4 proposals to federation governance:

```go
type ProposalRelay struct {
    queenAuthorization map[string]bool  // which queens can propose
}

func (pr *ProposalRelay) SubmitProposal(proposal types.Proposal) error {
    // Check queen authorization
    if !pr.queenAuthorization[proposal.ProposerQueen] {
        return ErrUnauthorizedQueen
    }

    // Add Society 4 metadata
    proposal.Society = "society4"
    proposal.PrivateChainHeight = getCurrentHeight()

    // Sign and send
    signedProposal := pr.sign(proposal)
    return pr.sendToFederationGovernance(signedProposal)
}
```

### 4. Resource Bridge

Manages ATP/ADP transfers between chains:

```go
type ResourceBridge struct {
    localTreasury    *treasury.Keeper
    federationBridge *ibc.Keeper
}

func (rb *ResourceBridge) RequestResources(amount uint64, reason string) error {
    request := types.ResourceRequest{
        Society:  "society4",
        Amount:   amount,
        Resource: "ATP",
        Reason:   reason,
        Queens:   rb.getRequestingQueens(),
    }

    return rb.federationBridge.SendResourceRequest(request)
}
```

## IBC Channel Configuration

### Channel Setup

```go
// Initialize IBC channel to federation
func SetupFederationChannel() error {
    config := ibc.ChannelConfig{
        PortID:    "federation",
        ChannelID: "channel-0",
        Version:   "federation-bridge-1",

        CounterpartyPortID:    "society4",
        CounterpartyChannelID: "channel-4",  // Society 4 is 4th to join

        ConnectionHops: []string{"connection-0"},

        Ordering: ibc.ORDERED,  // Messages must arrive in order
    }

    return ibc.CreateChannel(config)
}
```

### Message Types

```protobuf
// proto/federationbridge/v1/messages.proto

message PresenceProof {
    string society_id = 1;
    string chain_id = 2;
    bytes self_lct_hash = 3;
    bytes hardware_hash = 4;
    uint64 block_height = 5;
    google.protobuf.Timestamp timestamp = 6;
    bytes signature = 7;
}

message WitnessAttestation {
    string witness_society = 1;
    string target_society = 2;
    uint64 trust_level = 3;  // 0-100
    string attestation_type = 4;
    bytes signature = 5;
}

message FederationProposal {
    string proposing_society = 1;
    string proposing_queen = 2;
    string proposal_type = 3;
    bytes proposal_content = 4;
    uint64 private_chain_height = 5;
    bytes signature = 6;
}

message ResourceTransfer {
    string source = 1;
    string destination = 2;
    uint64 amount = 3;
    string resource_type = 4;  // ATP, ADP, etc.
    string memo = 5;
    bytes signature = 6;
}
```

## Bridge Operations

### 1. Initialization Sequence

```go
func InitializeBridge(ctx sdk.Context) error {
    // Step 1: Verify hardware binding
    if !verifyHardwareBinding(ctx) {
        return ErrInvalidHardware
    }

    // Step 2: Create self-LCT
    selfLCT := createSelfLCT(ctx)

    // Step 3: Establish IBC connection
    channel := setupIBCChannel(ctx)

    // Step 4: Send initial presence proof
    proof := createPresenceProof(selfLCT)

    // Step 5: Wait for witness threshold
    witnesses := awaitWitnesses(ctx, 3)  // Need 3 witnesses

    // Step 6: Activate bridge
    activateBridge(ctx, channel, witnesses)

    return nil
}
```

### 2. Periodic Operations

```go
func (k Keeper) BeginBlocker(ctx sdk.Context) {
    height := ctx.BlockHeight()

    // Every 100 blocks: Update presence
    if height % 100 == 0 {
        k.PublishPresence(ctx)
    }

    // Every 1000 blocks: Sync witnesses
    if height % 1000 == 0 {
        k.SyncWitnesses(ctx)
    }

    // Every 10000 blocks: Resource rebalance
    if height % 10000 == 0 {
        k.RebalanceResources(ctx)
    }
}
```

### 3. Emergency Protocols

```go
func (k Keeper) HandleEmergency(ctx sdk.Context, alert EmergencyAlert) {
    switch alert.Type {
    case "FEDERATION_FORK":
        k.PauseBridge(ctx)
        k.AwaitResolution(ctx)

    case "WITNESS_COMPROMISE":
        k.RevokeWitness(ctx, alert.CompromisedSociety)
        k.RequestNewWitnesses(ctx)

    case "RESOURCE_CRISIS":
        k.EmergencyResourceRequest(ctx)

    case "HARDWARE_MIGRATION":
        k.InitiateMigrationProtocol(ctx)
    }
}
```

## Security Model

### 1. Signature Verification

All bridge messages must be signed:

```go
func VerifyBridgeMessage(msg BridgeMessage) bool {
    // Get society's public key from federation
    pubKey := getFederationPubKey(msg.Society)

    // Verify signature
    return pubKey.VerifySignature(msg.Hash(), msg.Signature)
}
```

### 2. Replay Protection

Prevent message replay attacks:

```go
type NonceTracker struct {
    usedNonces map[string]map[uint64]bool  // society -> nonce -> used
}

func (nt *NonceTracker) CheckNonce(society string, nonce uint64) error {
    if nt.usedNonces[society][nonce] {
        return ErrNonceReused
    }
    nt.usedNonces[society][nonce] = true
    return nil
}
```

### 3. Rate Limiting

Prevent spam and resource exhaustion:

```go
type RateLimiter struct {
    limits map[string]RateLimit  // operation -> limit
}

func (rl *RateLimiter) CheckLimit(operation string) error {
    limit := rl.limits[operation]
    if limit.Exceeded() {
        return ErrRateLimitExceeded
    }
    limit.Increment()
    return nil
}
```

## Monitoring and Metrics

### Key Metrics

```go
type BridgeMetrics struct {
    PresenceProofsSent     uint64
    WitnessesReceived      uint64
    ProposalsSubmitted     uint64
    ResourceTransfers      uint64
    AverageLatency         time.Duration
    SuccessRate            float64
    LastFederationContact  time.Time
}
```

### Health Checks

```go
func (k Keeper) CheckBridgeHealth() BridgeHealth {
    return BridgeHealth{
        IBCChannelActive:   k.isChannelActive(),
        FederationReachable: k.canReachFederation(),
        WitnessThresholdMet: k.hasEnoughWitnesses(),
        LastPresenceProof:   k.getLastPresenceTime(),
        PendingMessages:     k.getPendingMessageCount(),
    }
}
```

## Configuration

### Bridge Parameters

```toml
[federation_bridge]
enabled = true
federation_chain_id = "web4-federation-main"
federation_rpc = "tcp://10.0.0.72:26657"  # Genesis node

# Presence broadcasting
presence_interval = 100  # blocks
presence_timeout = "30s"

# Witness requirements
min_witnesses = 3
witness_timeout = "5m"

# Resource management
resource_request_limit = 1000  # ATP per request
resource_cooldown = "1h"

# Security
require_hardware_binding = true
max_message_size = 1048576  # 1MB
rate_limit_per_block = 10
```

## Testing

### Unit Tests

```go
func TestPresenceProof(t *testing.T) {
    bridge := NewTestBridge()
    proof := bridge.CreatePresenceProof()

    assert.NotNil(t, proof)
    assert.Equal(t, "society4", proof.SocietyID)
    assert.True(t, bridge.VerifyProof(proof))
}
```

### Integration Tests

```go
func TestFederationConnection(t *testing.T) {
    // Start local chain
    localChain := StartTestChain("society4-private")

    // Start mock federation
    federation := StartMockFederation()

    // Setup bridge
    bridge := SetupBridge(localChain, federation)

    // Test presence broadcast
    err := bridge.PublishPresence()
    assert.NoError(t, err)

    // Verify receipt
    proof := federation.GetPresenceProof("society4")
    assert.NotNil(t, proof)
}
```

## Troubleshooting

### Common Issues

1. **IBC Channel Not Establishing**
   - Check network connectivity
   - Verify chain IDs match
   - Ensure ports are open

2. **Witnesses Not Accumulating**
   - Verify other societies are online
   - Check signature verification
   - Review trust requirements

3. **Resource Transfers Failing**
   - Confirm ATP balance sufficient
   - Check rate limits
   - Verify queen authorization

## Future Enhancements

1. **Multi-Channel Support**: Direct channels to each society
2. **Cross-Chain Contracts**: Smart contracts spanning chains
3. **Decentralized Bridges**: Remove single points of failure
4. **Privacy Features**: Zero-knowledge presence proofs
5. **Quantum-Resistant**: Post-quantum cryptography ready