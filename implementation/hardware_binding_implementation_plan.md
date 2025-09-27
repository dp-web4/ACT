# Implementation Plan: Secure Hardware Binding for Society Identity

## Executive Summary

Each society will establish unforgeable identity through hardware binding, creating a foundational 'self' LCT on their private blockchain that serves as the root of their presence in the federation. This creates a cryptographic chain from physical hardware → private identity → public federation participation.

## Current State Analysis

### What We Have
- Each society runs on a specific machine (Genesis on Ubuntu, Sprout on Jetson, Society2 on WSL2, Society4 on WSL2)
- Validator keys stored in `priv_validator_key.json`
- No hardware-level identity verification
- No private blockchain for societies (only federation chain)
- No LCT implementation yet in ACT blockchain

### What We Need
1. Hardware attestation capability per platform
2. Private blockchain instance per society
3. LCT module implementation
4. Cross-chain bridging mechanism
5. Root presence verification protocol

## Phase 1: Hardware Identity Extraction (Week 1)

### Platform-Specific Hardware Binding

#### Linux Native (Genesis)
```bash
# TPM 2.0 if available
tpm2_getcap properties-fixed | grep -E "TPM2_PT_FIRMWARE_VERSION|TPM2_PT_MANUFACTURER"

# DMI/SMBIOS
sudo dmidecode -s system-uuid
sudo dmidecode -s baseboard-serial-number

# CPU ID
lscpu | grep "Model name"
cat /proc/cpuinfo | grep "serial"
```

#### Jetson Orin Nano (Sprout)
```bash
# Tegra fuse ID (unique per chip)
cat /sys/devices/soc0/soc_uid

# Module serial
sudo i2ctransfer -f -y 0 w1@0x50 0x00 r8

# NVIDIA specific
tegrastats --print
```

#### WSL2 (Society2, Society4)
```bash
# Windows machine UUID via interop
powershell.exe -Command "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID"

# WSL2 instance ID
cat /proc/sys/kernel/random/boot_id

# Hyper-V VM ID
cat /sys/hypervisor/uuid
```

### Deliverables
- `hardware_identity.sh` script per platform
- JSON output format standardization
- Entropy analysis of hardware IDs

## Phase 2: Private Blockchain Setup (Week 1-2)

### Architecture
```
Federation Chain (Public)
    ↑
    | Bridge Module
    |
Private Society Chain
    |
    | Genesis Block
    ↓
Hardware Identity Root
```

### Implementation Steps

1. **Fork ACT blockchain for private instance**
```bash
# Each society runs private chain
society_private_chain/
├── config/
│   ├── genesis.json     # Includes hardware binding
│   ├── config.toml      # Local-only peers
│   └── app.toml         # Private chain settings
├── data/
└── keyring/
```

2. **Genesis Block Modification**
```json
{
  "genesis_time": "2025-09-27T00:00:00Z",
  "chain_id": "society4-private",
  "initial_height": "1",
  "consensus_params": {...},
  "app_state": {
    "hardware_binding": {
      "root_identity": {
        "hardware_id": "UUID-FROM-HARDWARE",
        "platform": "wsl2",
        "binding_time": "2025-09-27T00:00:00Z",
        "signature": "SIGNED-BY-VALIDATOR-KEY"
      }
    }
  }
}
```

3. **Isolated Operation**
- No external peers
- Local RPC only (127.0.0.1)
- Faster block time (1 second)
- Single validator

## Phase 3: LCT Module Development (Week 2)

### Module Structure
```go
// x/lct/types/lct.go
type LinkedContextToken struct {
    ID           string      // Unique LCT identifier
    Creator      string      // Society address
    ParentLCT    string      // Link to parent (empty for root)
    HardwareRoot string      // Hardware binding reference
    Context      Context     // MRH-bounded context
    Witnesses    []Witness   // Other LCTs that witnessed this
    Timestamp    time.Time
    Signature    []byte      // Creator's signature
}

type SelfLCT struct {
    LinkedContextToken
    HardwareBinding HardwareIdentity
    PrivateGenesis  string // Genesis hash of private chain
    PublicKey       crypto.PubKey
}
```

### Key Features
- Root LCT has no parent (self-originating)
- Hardware binding included in root
- Witness accumulation over time
- Unforgeable through signature chain

## Phase 4: Bridge Development (Week 2-3)

### Cross-Chain Protocol
```go
// x/bridge/keeper/federation_bridge.go
type FederationBridge struct {
    PrivateChain  *PrivateChainClient
    PublicChain   *PublicChainClient
    RootLCT       *SelfLCT
}

func (fb *FederationBridge) PublishPresence() error {
    // 1. Get latest root LCT from private chain
    rootLCT := fb.PrivateChain.GetRootLCT()

    // 2. Create presence proof
    proof := CreatePresenceProof(rootLCT)

    // 3. Submit to federation chain
    return fb.PublicChain.SubmitPresenceFootprint(proof)
}
```

### IBC Integration
- Use IBC for cross-chain communication
- Custom packet types for LCT bridging
- Verification on federation chain

## Phase 5: Society Swarm Design Session (Week 3)

### Collaborative Design Process

1. **Create Design Space**
```markdown
# federation_inbox/society4_hardware_binding_design_session.md

## Invitation to Collaborate

Society 4 proposes collaborative design of hardware binding protocol.

### Questions for Discussion:
1. How should hardware diversity be handled?
2. What happens if hardware changes (migration)?
3. How deep should binding go (CPU? TPM? MAC address?)
4. Privacy considerations for hardware exposure?
5. Consensus mechanism for accepting root LCTs?
```

2. **Gather Feedback via Git Mailbox**
- Each society submits position papers
- Technical concerns raised
- Security analysis shared

3. **Synthesis Meeting**
- Virtual consensus through async messages
- Design decisions documented
- Implementation responsibilities assigned

## Phase 6: Implementation (Week 3-4)

### Development Tasks

1. **Core Infrastructure**
```bash
# Directory structure
act/
├── implementation/
│   ├── hardware_binding/
│   │   ├── extractors/       # Platform-specific
│   │   ├── validators/       # Verification logic
│   │   └── bridges/          # Chain bridging
│   ├── private_chain/
│   │   ├── setup.sh
│   │   └── genesis_template.json
│   └── lct_module/
│       ├── spec/
│       ├── keeper/
│       └── types/
```

2. **Testing Framework**
- Hardware simulation for CI/CD
- Multi-chain test environment
- Security penetration tests

## Phase 7: Federation Proposal (Week 4)

### Proposal Document Structure

```markdown
# PROPOSAL_004_HARDWARE_BINDING_ROOT_IDENTITY.md

## Abstract
Establish unforgeable society identity through hardware binding...

## Motivation
- Current state vulnerabilities
- Hardware as trust anchor
- Privacy-preserving identity

## Specification
- Technical implementation
- Protocol details
- Migration procedures

## Security Analysis
- Threat model
- Mitigation strategies
- Privacy considerations

## Implementation Timeline
- Phases and milestones
- Society responsibilities
- Testing requirements

## Backward Compatibility
- Existing validators
- Chain upgrade path
- Fallback mechanisms
```

## Risk Analysis

### Technical Risks
1. **Hardware Heterogeneity**: Different platforms expose different IDs
   - Mitigation: Abstract binding interface

2. **Hardware Failure**: What if hardware dies?
   - Mitigation: Backup/recovery protocol with federation witness

3. **Privacy Exposure**: Hardware IDs might be sensitive
   - Mitigation: Hash hardware IDs, never expose raw values

### Operational Risks
1. **Complexity Increase**: More moving parts
   - Mitigation: Phased rollout, extensive testing

2. **Migration Difficulty**: Moving society to new hardware
   - Mitigation: Federation-witnessed migration protocol

## Success Criteria

### Phase 1-2 Success
- [ ] Hardware ID extraction working on all platforms
- [ ] Private chains running for each society
- [ ] Genesis blocks include hardware binding

### Phase 3-4 Success
- [ ] LCT module implemented and tested
- [ ] Bridge successfully transfers root LCTs
- [ ] Federation chain accepts presence proofs

### Phase 5-6 Success
- [ ] All societies participate in design
- [ ] Consensus on implementation details
- [ ] Working prototype across federation

### Phase 7 Success
- [ ] Proposal accepted by federation
- [ ] Implementation timeline agreed
- [ ] Resources allocated

## Resource Requirements

### Development Resources
- 4 weeks developer time
- Testing infrastructure
- Documentation effort

### Infrastructure
- Private chain nodes (4x)
- Bridge infrastructure
- Monitoring/logging

### Coordination
- Daily Git Mailbox sync
- Weekly design reviews
- Security audits

## Next Steps

1. **Immediate Actions**
   - Implement hardware extraction scripts
   - Set up private chain infrastructure
   - Begin LCT module development

2. **Coordination**
   - Send design session invitation via Git Mailbox
   - Create shared development branches
   - Establish testing protocols

3. **Documentation**
   - Technical specifications
   - Security analysis
   - User guides

This plan establishes true hardware-rooted identity for the federation, making each society's presence unforgeable and creating a cryptographic chain from silicon to sovereignty.