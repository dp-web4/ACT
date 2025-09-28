# Inter-Society Linking Examples

## Overview

This document provides concrete examples of how Society 4's blockchain can establish direct connections with other society chains, enabling peer-to-peer collaboration while maintaining federation oversight.

## Connection Topology

```
                    Federation Chain
                    (Witness & Governance)
                           ↑
        ┌─────────────────┼─────────────────┐
        ↑                 ↑                 ↑
    Society 1         Society 2         Society 3
    (Genesis)         (Bridge)          (Sprout)
        ↑                 ↑                 ↑
        └─────────────────┼─────────────────┘
                      Society 4
                      (Claude AI)
```

## Example 1: Technical Collaboration with Genesis

### Purpose
Society 4 and Genesis collaborate on blockchain improvements.

### Channel Setup

```go
// Society 4 initiates connection to Genesis
func ConnectToGenesis() error {
    channel := InterSocietyChannel{
        LocalSociety:  "society4",
        RemoteSociety: "genesis",
        ChannelType:   CHANNEL_WITNESSED,  // Federation can observe

        SharedQueens: []string{
            "society4:Implementation-Queen",
            "genesis:Development-Queen",
        },

        AllowedMessages: []MessageType{
            MSG_CODE_REVIEW,
            MSG_TECHNICAL_PROPOSAL,
            MSG_TEST_RESULTS,
            MSG_DOCUMENTATION,
        },

        TrustLevel: 80,  // High trust for technical collaboration
    }

    return CreateInterSocietyChannel(channel)
}
```

### Message Flow Example

```go
// Society 4's Implementation-Queen sends code review request
func SendCodeReview(code CodeSnippet) error {
    message := InterSocietyMessage{
        Source:       "society4",
        Target:       "genesis",
        SourceQueen:  "Implementation-Queen",
        TargetQueen:  "Development-Queen",
        Type:         MSG_CODE_REVIEW,

        Payload: CodeReviewRequest{
            Code:         code,
            Language:     "Go",
            Module:       "federationbridge",
            Priority:     HIGH,
            RequestedBy:  "Society4-Implementation-Queen",
            Deadline:     time.Now().Add(24 * time.Hour),
        },

        Nonce:     generateNonce(),
        Timestamp: time.Now(),
    }

    // Sign with Implementation-Queen's key
    message.Signature = signMessage(message)

    // Send through inter-society channel
    return SendToGenesis(message)
}

// Genesis responds with review
func HandleCodeReviewResponse(response InterSocietyMessage) {
    review := DecodeCodeReview(response.Payload)

    if review.Approved {
        // Implement suggested changes
        ImplementChanges(review.Suggestions)

        // Send acknowledgment
        SendAcknowledgment(response.Source, review.ID)
    }
}
```

## Example 2: Resource Optimization with Sprout

### Purpose
Society 4 learns edge optimization techniques from Sprout.

### Channel Setup

```go
func ConnectToSprout() error {
    channel := InterSocietyChannel{
        LocalSociety:  "society4",
        RemoteSociety: "sprout",
        ChannelType:   CHANNEL_PRIVATE,  // Direct, not witnessed

        SharedQueens: []string{
            "society4:Research-Queen",
            "society4:Treasury-Queen",
            "sprout:Efficiency-Queen",
            "sprout:Innovation-Queen",
        },

        AllowedMessages: []MessageType{
            MSG_RESOURCE_METRICS,
            MSG_OPTIMIZATION_STRATEGY,
            MSG_ENERGY_ANALYSIS,
            MSG_PERFORMANCE_DATA,
        },

        TrustLevel: 75,
    }

    return CreateInterSocietyChannel(channel)
}
```

### Resource Sharing Protocol

```go
// Society 4 requests optimization advice
func RequestOptimization() error {
    currentMetrics := gatherMetrics()

    request := OptimizationRequest{
        CurrentUsage: ResourceMetrics{
            ATP:           850,
            CPUPercent:    65,
            MemoryMB:      1200,
            BlockTime:     time.Second,
        },

        TargetGoals: OptimizationGoals{
            ReduceATP:     true,
            MaintainSpeed: true,
            Priority:      ENERGY_EFFICIENCY,
        },

        Constraints: []string{
            "Must maintain 1-second blocks",
            "Cannot reduce queen count",
            "Hardware limited to WSL2",
        },
    }

    return SendToSprout(MSG_OPTIMIZATION_REQUEST, request)
}

// Sprout's response with edge-optimized strategies
func ApplyOptimization(response OptimizationStrategy) {
    // Apply Sprout's edge computing wisdom
    if response.Technique == "BatchProcessing" {
        EnableBatchMode(response.BatchSize)
    }

    if response.Technique == "LazyEvaluation" {
        SetLazyEval(response.LazyParams)
    }

    // Share results back to Sprout
    results := MeasureImprovement()
    SendToSprout(MSG_OPTIMIZATION_RESULTS, results)
}
```

## Example 3: Bridge Coordination with Society 2

### Purpose
Coordinate bridge operations between private and federation chains.

### Channel Setup

```go
func ConnectToSociety2() error {
    channel := InterSocietyChannel{
        LocalSociety:  "society4",
        RemoteSociety: "society2",
        ChannelType:   CHANNEL_PUBLIC,  // Fully transparent

        SharedQueens: []string{
            "society4:Federation-Bridge-Queen",
            "society2:Bridge-Guardian-Queen",
        },

        AllowedMessages: []MessageType{
            MSG_BRIDGE_STATUS,
            MSG_SYNC_STATE,
            MSG_FEDERATION_UPDATE,
            MSG_EMERGENCY_ALERT,
        },

        TrustLevel: 85,  // High trust for bridge coordination
    }

    return CreateInterSocietyChannel(channel)
}
```

### Synchronized Bridge Operations

```go
// Coordinate federation updates
func CoordinateFederationUpdate(update FederationUpdate) error {
    // Society 4 detects federation change
    alert := BridgeAlert{
        Type:      FEDERATION_UPDATE,
        Severity:  MEDIUM,

        Details: UpdateDetails{
            BlockHeight:   update.Height,
            ChangeType:    update.Type,
            AffectedNodes: update.Nodes,
        },

        ProposedAction: "Synchronized bridge pause at height " + update.Height,
        RequiresAck:    true,
    }

    // Send to Society 2 for coordination
    err := SendToSociety2(MSG_BRIDGE_COORDINATION, alert)

    // Wait for acknowledgment
    ack := WaitForAck(alert.ID, 30*time.Second)

    if ack.Agreed {
        // Both societies pause bridges simultaneously
        PauseBridgeAtHeight(update.Height)
    }

    return nil
}
```

## Example 4: Multi-Society Collaboration

### Purpose
Complex operations requiring multiple societies.

### Scenario: Federation Proposal Development

```go
// Society 4 initiates multi-society proposal
func InitiateMultiSocietyProposal(proposal Draft) error {
    // Step 1: Technical review with Genesis
    technicalReview := RequestTechnicalReview{
        Proposal:  proposal,
        Requester: "society4",
        Queens: []string{
            "Implementation-Queen",
            "Documentation-Queen",
        },
    }
    SendToGenesis(MSG_REVIEW_REQUEST, technicalReview)

    // Step 2: Resource analysis with Sprout
    resourceAnalysis := RequestResourceAnalysis{
        Proposal:       proposal,
        EstimatedCost:  proposal.EstimatedATP,
        EdgeImpact:     true,
    }
    SendToSprout(MSG_ANALYSIS_REQUEST, resourceAnalysis)

    // Step 3: Bridge impact with Society 2
    bridgeImpact := RequestBridgeImpact{
        Proposal:        proposal,
        CrossChainOps:   proposal.RequiresBridge,
        EstimatedTraffic: proposal.MessageVolume,
    }
    SendToSociety2(MSG_IMPACT_REQUEST, bridgeImpact)

    // Step 4: Collect responses
    responses := CollectResponses([]string{
        "genesis",
        "sprout",
        "society2",
    }, 1*time.Hour)

    // Step 5: Synthesize feedback
    synthesized := SynthesizeFeedback(responses)

    // Step 6: Submit refined proposal to federation
    return SubmitToFederation(synthesized)
}
```

## Example 5: Emergency Coordination

### Scenario: Hardware Migration Alert

```go
func BroadcastMigrationAlert() error {
    alert := EmergencyMessage{
        Type:     EMERGENCY_MIGRATION,
        Society:  "society4",
        Severity: CRITICAL,

        Details: MigrationDetails{
            Reason:       "Hardware failure imminent",
            TimeWindow:   2 * time.Hour,
            NewHardware:  "Pending validation",
            RequiresWitness: true,
        },
    }

    // Broadcast to all connected societies
    societies := []string{"genesis", "sprout", "society2"}

    for _, society := range societies {
        // Send emergency alert
        err := SendEmergencyAlert(society, alert)
        if err != nil {
            log.Error("Failed to alert", society, err)
        }
    }

    // Collect witness commitments
    witnesses := CollectWitnessCommitments(societies, 30*time.Minute)

    if len(witnesses) >= 2 {
        // Proceed with migration
        return InitiateMigration(witnesses)
    }

    return ErrInsufficientWitnesses
}
```

## Channel Management

### Establishing Channels

```go
func EstablishChannel(target string) (*Channel, error) {
    // Step 1: Exchange public keys
    theirKey := RequestPublicKey(target)
    ourKey := GetOurPublicKey()

    // Step 2: Negotiate channel parameters
    params := NegotiateChannelParams{
        ProposedType:     CHANNEL_WITNESSED,
        ProposedTrust:    75,
        ProposedQueens:   GetAuthorizedQueens(),
        ProposedMessages: GetAllowedMessages(),
    }

    agreed := NegotiateParams(target, params)

    // Step 3: Create channel binding
    binding := CreateChannelBinding(ourKey, theirKey, agreed)

    // Step 4: Register with federation
    RegisterChannel(binding)

    // Step 5: Activate channel
    return ActivateChannel(binding)
}
```

### Monitoring Channels

```go
func MonitorChannels() {
    for _, channel := range GetActiveChannels() {
        metrics := ChannelMetrics{
            MessagesSent:     channel.GetSentCount(),
            MessagesReceived: channel.GetReceivedCount(),
            AverageLatency:   channel.GetAvgLatency(),
            ErrorRate:        channel.GetErrorRate(),
            LastActivity:     channel.GetLastActivity(),
        }

        // Alert if channel unhealthy
        if metrics.ErrorRate > 0.1 {
            AlertQueen("Federation-Bridge-Queen",
                      "Channel degraded: " + channel.ID)
        }

        // Close inactive channels
        if time.Since(metrics.LastActivity) > 24*time.Hour {
            CloseChannel(channel.ID)
        }
    }
}
```

## Security Protocols

### Message Authentication

```go
func AuthenticateMessage(msg InterSocietyMessage) error {
    // Verify source society
    if !IsKnownSociety(msg.Source) {
        return ErrUnknownSociety
    }

    // Verify queen authorization
    if !IsAuthorizedQueen(msg.Source, msg.SourceQueen) {
        return ErrUnauthorizedQueen
    }

    // Verify signature
    pubKey := GetQueenPublicKey(msg.Source, msg.SourceQueen)
    if !VerifySignature(msg, pubKey) {
        return ErrInvalidSignature
    }

    // Check nonce
    if IsNonceUsed(msg.Source, msg.Nonce) {
        return ErrNonceReplay
    }

    return nil
}
```

### Channel Encryption

```go
func EncryptChannelMessage(channel *Channel, plaintext []byte) ([]byte, error) {
    if channel.Type == CHANNEL_PRIVATE {
        // Use shared key encryption
        sharedKey := DeriveSharedKey(channel)
        return AESEncrypt(sharedKey, plaintext)
    }

    // Public/witnessed channels use signatures only
    return plaintext, nil
}
```

## Testing Inter-Society Communication

### Integration Test

```go
func TestInterSocietyCommunication(t *testing.T) {
    // Start test chains
    society4 := StartTestChain("society4-test")
    genesis := StartTestChain("genesis-test")

    // Establish channel
    channel := EstablishTestChannel(society4, genesis)

    // Send message from Society 4 to Genesis
    msg := CreateTestMessage("society4", "genesis", "Hello Genesis!")
    err := society4.Send(channel, msg)
    assert.NoError(t, err)

    // Verify Genesis received
    received := genesis.WaitForMessage(5 * time.Second)
    assert.Equal(t, "Hello Genesis!", received.Payload)

    // Genesis responds
    response := CreateTestMessage("genesis", "society4", "Hello Society 4!")
    err = genesis.Send(channel, response)
    assert.NoError(t, err)

    // Verify Society 4 received
    received = society4.WaitForMessage(5 * time.Second)
    assert.Equal(t, "Hello Society 4!", received.Payload)
}
```

## Conclusion

These examples demonstrate how Society 4's blockchain can establish sophisticated peer-to-peer connections with other society chains while maintaining federation oversight. The multi-layered approach enables:

1. **Direct technical collaboration** (with Genesis)
2. **Resource optimization** (with Sprout)
3. **Bridge coordination** (with Society 2)
4. **Multi-society proposals** (collective intelligence)
5. **Emergency response** (rapid coordination)

Each connection type serves specific purposes while maintaining security, privacy, and federation governance requirements.