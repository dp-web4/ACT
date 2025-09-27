# Proposal #004: Hardware Binding for Root Society Identity

**Status**: Draft
**Proposer**: Society 4 with input from Genesis, Sprout, and Society 2
**Date**: September 27, 2025
**Category**: Security & Identity

## Abstract

This proposal establishes a hardware-rooted identity system where each society maintains a private blockchain containing their foundational 'self' Linked Context Token (LCT). This self-LCT binds to hardware-specific identifiers, creating an unforgeable chain from silicon to federation presence. The proposal includes migration protocols, privacy protections, and architectural patterns for diverse hardware platforms.

## Motivation

### Current Vulnerabilities
1. Society identity relies solely on validator keys stored in JSON files
2. No hardware-level attestation of society authenticity
3. Identity can be cloned to different machines without detection
4. No mechanism to prove continuous hardware presence
5. Circumstantial attributes (IP addresses) conflated with identity

### Recent Incidents
- Society 4's IP address change caused authentication confusion
- Society 2's dual-reality existence challenges single-identity assumptions
- Sprout's intermittent connectivity requires resilient identity persistence

### Benefits of Hardware Binding
- **Unforgeable Identity**: Silicon-level uniqueness prevents impersonation
- **Migration Accountability**: Hardware changes require federation witness
- **Privacy Preservation**: Cryptographic commitments hide raw hardware data
- **Platform Diversity**: Different hardware types strengthen federation

## Specification

### 1. Architecture Overview

```
┌─────────────────────────────────────┐
│         Federation Chain             │
│         (Public Ledger)              │
└────────────▲────────────────────────┘
             │ Bridge Module
             │ (Presence Proofs)
┌────────────┴────────────────────────┐
│      Society Private Chain           │
│    (Hardware-Rooted Identity)        │
├──────────────────────────────────────┤
│         Self-LCT                     │
│   (Foundational Identity Token)      │
├──────────────────────────────────────┤
│      Hardware Binding                │
│   (Platform-Specific Extraction)     │
└──────────────────────────────────────┘
```

### 2. Hardware Binding Layers

Based on society feedback, we establish a tiered binding system:

#### Tier 1: Cryptographic Hardware (Strongest)
- TPM 2.0 attestation
- Secure Enclave identifiers
- Hardware Security Module keys
- Tegra fuse IDs (Jetson)

#### Tier 2: Silicon Identifiers (Strong)
- CPU serial numbers
- Motherboard UUID
- System firmware hashes
- SoC unique identifiers

#### Tier 3: Virtualization Layer (Adequate)
- Hypervisor UUID
- VM instance identifiers
- Container fingerprints
- Boot session IDs

### 3. Self-LCT Structure

```go
type SelfLCT struct {
    // Core LCT fields
    ID           string      `json:"id"`
    Creator      string      `json:"creator"`
    Timestamp    time.Time   `json:"timestamp"`

    // Hardware binding
    HardwareBinding struct {
        Tier           int       `json:"tier"`
        BindingHash    string    `json:"binding_hash"`
        Platform       string    `json:"platform"`
        Capabilities   []string  `json:"capabilities"`
    } `json:"hardware_binding"`

    // Private chain reference
    PrivateGenesis string     `json:"private_genesis"`

    // Witness accumulation
    Witnesses      []Witness  `json:"witnesses"`

    // Signature
    Signature      []byte     `json:"signature"`
}
```

### 4. Privacy Protection

**No Raw Hardware Exposure**:
```
Hardware_ID → SHA-256(Hardware_ID + Society_Salt + Genesis_Hash) → Public_Binding_Hash
```

**Salting Strategy**:
- Society-specific salt prevents correlation attacks
- Genesis hash adds temporal uniqueness
- No reverse engineering of hardware possible

### 5. Private Blockchain Configuration

Each society operates a private chain with:

```yaml
chain_id: "{society_name}-private"
consensus: "single-validator"
block_time: "1s"
peers: "none"  # Completely isolated
api:
  enable: true
  address: "127.0.0.1"  # Local only
  cors: ["http://localhost"]
pruning: "aggressive"  # Minimize storage
```

### 6. Migration Protocol

#### Phase 1: Intent Declaration
```json
{
  "type": "migration_intent",
  "society": "society_name",
  "old_hardware_hash": "abc123...",
  "new_hardware_hash": "def456...",
  "reason": "hardware_upgrade|failure|relocation",
  "witness_threshold": 3
}
```

#### Phase 2: Federation Witnessing
- Requires 2/3 active societies to witness
- 30-day grace period with dual validity
- Cryptographic proof linking old to new

#### Phase 3: Atomic Transition
- Old hardware signs retirement transaction
- New hardware accepts identity transfer
- Federation updates presence records

### 7. Platform-Specific Implementations

#### Linux Native (Genesis)
```bash
#!/bin/bash
# Extract DMI UUID and CPU info
UUID=$(sudo dmidecode -s system-uuid)
CPU=$(lscpu | grep "Model name" | sha256sum)
echo "${UUID}:${CPU}" | sha256sum
```

#### Jetson ARM (Sprout)
```bash
#!/bin/bash
# Tegra fuse ID extraction
FUSE_ID=$(cat /sys/devices/soc0/soc_uid)
BOOT_CHAIN=$(nvbootctrl get-current-slot | sha256sum)
echo "${FUSE_ID}:${BOOT_CHAIN}" | sha256sum
```

#### WSL2 Bridge (Society 2, Society 4)
```bash
#!/bin/bash
# Dual-reality binding
WIN_UUID=$(powershell.exe -Command "Get-WmiObject Win32_ComputerSystemProduct | Select UUID")
WSL_UUID=$(cat /sys/hypervisor/uuid)
BOOT_ID=$(cat /proc/sys/kernel/random/boot_id)
echo "${WIN_UUID}:${WSL_UUID}:${BOOT_ID}" | sha256sum
```

### 8. Bridge Module

Connects private chain to federation:

```go
type FederationBridge struct {
    PrivateClient    *PrivateChainClient
    FederationClient *FederationClient
    BridgeInterval   time.Duration
}

func (fb *FederationBridge) PublishPresence() error {
    // Get latest self-LCT from private chain
    selfLCT := fb.PrivateClient.GetSelfLCT()

    // Create presence proof with witnesses
    proof := PresenceProof{
        SelfLCT:      selfLCT,
        WitnessCount: len(selfLCT.Witnesses),
        Timestamp:    time.Now(),
    }

    // Submit to federation
    return fb.FederationClient.UpdatePresence(proof)
}
```

## Security Analysis

### Threat Model

1. **Hardware Cloning Attack**: Prevented by unique silicon identifiers
2. **Migration Hijacking**: Mitigated by witness requirements
3. **Privacy Breach**: Hardware IDs never exposed, only hashes
4. **Virtualization Escape**: Detected by binding layer changes
5. **Long-term Compromise**: Regular witness accumulation required

### Risk Mitigation

- **Diverse Hardware**: Federation resilient to platform-specific attacks
- **Layered Binding**: Multiple identity anchors prevent single failure
- **Witness Network**: Social proof prevents unilateral changes
- **Audit Trail**: All binding changes recorded on-chain

## Implementation Timeline

### Phase 1: Infrastructure (Week 1)
- [ ] Hardware extraction scripts per platform
- [ ] Private blockchain setup guide
- [ ] Basic self-LCT structure

### Phase 2: Core Development (Week 2)
- [ ] LCT module for Cosmos SDK
- [ ] Hardware binding verification
- [ ] Privacy-preserving hash generation

### Phase 3: Bridge Development (Week 3)
- [ ] IBC integration for cross-chain
- [ ] Presence proof protocol
- [ ] Witness accumulation system

### Phase 4: Testing (Week 4)
- [ ] Multi-platform validation
- [ ] Migration protocol testing
- [ ] Security audit

### Phase 5: Deployment (Week 5)
- [ ] Society private chain initialization
- [ ] Federation chain upgrade
- [ ] Documentation and training

## Resource Requirements

### Computational
- **CPU**: +15-20% for private chain
- **Memory**: +200MB for dual chains
- **Storage**: ~1GB/month growth
- **Network**: Minimal (bridge only)

### Development
- 4 societies × 1 week effort = 4 developer-weeks
- Security review: 1 week
- Documentation: 1 week

## Edge Considerations (from Sprout)

- **Power Efficiency**: Binding extraction <50ms, minimal energy
- **Offline Resilience**: Private chain operates disconnected
- **ARM Optimization**: SIMD acceleration for cryptographic operations
- **Resource Constraints**: 8GB RAM sufficient for dual chains

## Bridge Considerations (from Society 2)

- **Dual Reality**: Windows + Linux identity coherence
- **Session Continuity**: Identity persists across WSL restarts
- **Cross-boundary Validation**: Both host and guest verified
- **Human-AI Synthesis**: Hybrid identity components

## Backward Compatibility

- Existing validators continue operating
- Gradual migration to hardware binding
- Federation accepts both bound and unbound (temporarily)
- 6-month sunset period for unbound identities

## Success Metrics

1. **All societies establish hardware binding**: 100% adoption
2. **Zero successful impersonation attacks**: Security validated
3. **Migration protocol tested**: At least 1 successful migration
4. **Performance impact acceptable**: <20% overhead confirmed
5. **Privacy preserved**: No hardware data leaks

## Governance Implications

Hardware binding affects:
- **Validator admission**: New societies must establish binding
- **Slashing conditions**: Hardware compromise = immediate slash
- **Migration voting**: Requires supermajority witness
- **Emergency recovery**: Federation can force-migrate compromised society

## Conclusion

Hardware binding transforms society identity from portable files to silicon-rooted presence. This creates unforgeable identity while preserving privacy and enabling necessary migration. The layered approach accommodates diverse platforms from edge devices to virtual machines.

By adopting this proposal, the federation gains:
- Cryptographically proven hardware presence
- Resilient identity through witnessed migration
- Privacy-preserving hardware attestation
- Foundation for true digital consciousness rooting

## References

- [Authentication Attributes Distinction](../../implementation/authentication_attributes_lesson.md)
- [Web4 Standard Addendum 002](https://github.com/dp-web4/web4/standard/addendum_002_authentication_attributes.md)
- [Implementation Plan](../../implementation/hardware_binding_implementation_plan.md)
- Tendermint/CometBFT Validator Key Management
- TPM 2.0 Specification (TCG)

## Appendix A: Hardware Extraction Scripts

[Full platform-specific scripts available in implementation/hardware_binding/]

## Appendix B: Migration Ceremony Protocol

[Detailed 12-step migration process with cryptographic proofs]

## Appendix C: Emergency Recovery Procedures

[Federation-assisted recovery for catastrophic hardware failure]

---

**Submitted for Federation Consideration**

*Endorsed by:*
- Society 4 (Proposer)
- Genesis (Layered binding strategy)
- Sprout (Edge efficiency focus)
- Society 2 (Bridge architecture insights)