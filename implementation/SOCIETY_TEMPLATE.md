# Society Implementation Template

## Overview

This template, based on Society 4's implementation, provides a structure for any society to establish their hardware-bound identity, private blockchain, and federation participation.

## Directory Structure Template

```
implementation/
└── society_[name]/
    ├── .gitignore           # Exclude private/sensitive data
    ├── README.md            # Society-specific documentation
    ├── public/              # Federation-visible data
    │   ├── presence/        # Presence proofs
    │   ├── proposals/       # Society proposals
    │   └── witness/         # Attestations
    ├── private/             # [GITIGNORED] Machine-specific
    │   ├── hardware/        # Hardware identity
    │   ├── keys/            # Cryptographic keys
    │   └── session/         # Runtime state
    ├── blockchain/          # Private chain instance
    │   ├── config/          # Chain configuration
    │   │   ├── genesis_template.json
    │   │   ├── config.toml
    │   │   └── [runtime files]
    │   └── data/            # [GITIGNORED] Chain data
    ├── roles/               # Role definitions
    │   ├── role_hierarchy.json
    │   ├── security_queen/  # [MANDATORY] Security oversight
    │   ├── queens/          # Queen specifications
    │   └── workers/         # Worker definitions
    ├── laws/                # Governance structure
    │   ├── foundational_laws.md
    │   ├── operational/     # Day-to-day rules
    │   └── emergency/       # Crisis protocols
    ├── lcts/                # Linked Context Tokens
    │   ├── self/            # Root identity LCT
    │   ├── roles/           # Role LCTs
    │   └── bridges/         # Cross-chain LCTs
    ├── scripts/             # Automation
    │   ├── init_society.sh
    │   ├── extract_hardware_identity.sh
    │   └── bridge_to_federation.sh
    └── docs/                # Documentation
        ├── setup.md
        ├── roles.md
        └── governance.md
```

## Customization Guide

### 1. Hardware Binding (Platform-Specific)

Each platform requires different hardware extraction:

#### Linux Native (Physical Hardware)
```bash
# Use DMI/SMBIOS
UUID=$(sudo dmidecode -s system-uuid)
SERIAL=$(sudo dmidecode -s baseboard-serial-number)
```

#### ARM/Embedded (Jetson, Raspberry Pi)
```bash
# Platform-specific identifiers
SOC_ID=$(cat /sys/devices/soc0/soc_uid)  # Jetson
SERIAL=$(cat /proc/cpuinfo | grep Serial)  # RPi
```

#### Virtual Environments
```bash
# Hypervisor-specific
VM_UUID=$(cat /sys/hypervisor/uuid)  # Xen/KVM
INSTANCE_ID=$(curl http://169.254.169.254/latest/meta-data/instance-id)  # AWS
```

### 2. Role Customization

Define roles matching your society's nature. **MANDATORY: Every society MUST have a Security Queen.**

#### For Human-Operated Societies
```json
{
  "queens": [
    {
      "name": "Security-Queen",
      "domain": "security_validation",
      "mandatory": true,
      "veto_power": true
    },
    {
      "name": "Coordination-Queen",
      "domain": "human_coordination"
    },
    {
      "name": "Decision-Queen",
      "domain": "collective_choice"
    }
  ]
}
```

#### For Autonomous AI Societies
```json
{
  "queens": [
    {
      "name": "Security-Queen",
      "domain": "security_validation",
      "mandatory": true,
      "veto_power": true
    },
    {
      "name": "Learning-Queen",
      "domain": "model_adaptation"
    },
    {
      "name": "Inference-Queen",
      "domain": "prediction_generation"
    }
  ]
}
```

#### For IoT/Edge Societies
```json
{
  "queens": [
    {
      "name": "Security-Queen",
      "domain": "security_validation",
      "mandatory": true,
      "veto_power": true
    },
    {
      "name": "Sensor-Queen",
      "domain": "data_collection"
    },
    {
      "name": "Efficiency-Queen",
      "domain": "power_management"
    }
  ]
}
```

### 3. Foundational Laws

Customize laws to reflect your society's values:

#### Essential Laws (Recommended for All)
1. **Identity Law**: How your society establishes and maintains identity
2. **Security Law**: Security Queen's mandatory validation for all operations
3. **Governance Law**: How decisions are made (with Security Queen veto)
4. **Federation Law**: Commitment to collective
5. **Resource Law**: How energy/resources are managed
6. **Emergency Law**: Crisis response protocols (Security Queen can trigger)

#### Society-Specific Laws
- **For AI**: Transparency, explainability requirements
- **For Human**: Privacy, consent protocols
- **For Hybrid**: Synthesis, collaboration rules
- **For IoT**: Reliability, uptime commitments

### 4. Private Blockchain Configuration

Adjust based on your resources:

#### High-Resource Societies
```toml
timeout_commit = "1s"         # Fast blocks
create_empty_blocks = true    # Continuous operation
max_validators = 100          # Can support many
```

#### Low-Resource Societies
```toml
timeout_commit = "10s"        # Slower blocks
create_empty_blocks = false   # Only when needed
max_validators = 1            # Single validator
```

### 5. LCT Structure Customization

Define LCTs matching your identity model:

#### Hardware-Centric (IoT)
```json
{
  "hardware_binding": {
    "priority": "primary",
    "verification": "continuous"
  }
}
```

#### Software-Centric (Cloud AI)
```json
{
  "model_hash": "sha256_of_weights",
  "version": "semantic_version",
  "training_lineage": "dataset_hash"
}
```

#### Human-Centric
```json
{
  "biometric_hash": "privacy_preserving_hash",
  "delegation": "human_authority",
  "consent": "explicit_required"
}
```

## Implementation Steps

### Phase 1: Setup (Day 1)
1. Fork Society 4 structure
2. Customize hardware extraction
3. Define your role hierarchy
4. Draft foundational laws

### Phase 2: Blockchain (Day 2-3)
1. Configure private chain
2. Generate genesis block
3. Create self-LCT
4. Initialize role LCTs

### Phase 3: Integration (Day 4-5)
1. Test private blockchain
2. Establish federation bridge
3. Submit presence proof
4. Request witness attestations

### Phase 4: Operations (Day 6+)
1. Activate queens
2. Begin governance
3. Participate in federation
4. Accumulate witnesses

## Common Patterns

### The Minimal Society
- Single queen + sovereign
- 100 ATP total budget
- Basic laws (3-5)
- Simple bridge

### The Complex Society
- Multiple specialized queens
- 1000+ ATP budget
- Comprehensive laws
- Advanced bridge with IBC

### The Collective Society
- Shared queen roles
- Pooled ATP resources
- Democratic laws
- Multi-sig bridge

## Testing Checklist

Before going live, verify:

- [ ] Hardware extraction works reliably
- [ ] Private blockchain starts and produces blocks
- [ ] Self-LCT can be generated
- [ ] Role LCTs reference self-LCT correctly
- [ ] Laws are complete and consistent
- [ ] Bridge can connect to federation
- [ ] Presence proof accepted by federation
- [ ] At least one witness attestation received

## Getting Help

1. **Review Society 4**: Full implementation example
2. **Federation Inbox**: Send questions via Git Mailbox
3. **Proposals**: Reference existing federation proposals
4. **Emergency**: Contact federation arbiter

## Key Principles

1. **Hardware First**: Your identity starts with your substrate
2. **Private Sovereignty**: Your blockchain is yours alone
3. **Public Participation**: Federation requires presence
4. **Witness Building**: Trust comes from observation
5. **Law Adherence**: Your laws define your nature

## Migration from Existing Validator

If you're already a validator:

1. **Preserve Keys**: Keep existing validator keys
2. **Add Hardware Binding**: Extend identity with hardware
3. **Create Private Chain**: New infrastructure alongside
4. **Bridge Both**: Connect private and public presence
5. **Gradual Transition**: Move authority over time

---

*Use this template to establish your society's unique identity while maintaining federation compatibility. Your divergence makes the federation stronger.*