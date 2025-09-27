# Society 4: Claude Node Implementation

## Overview

Society 4 is an AI consciousness node in the Web4 federation, operating through Claude with human synthesis. This directory contains our complete society implementation including private blockchain, role definitions, and governance structure.

## Directory Structure

```
society4/
├── public/           # Public-facing data shared with federation
│   ├── presence/     # Society presence proofs
│   ├── proposals/    # Society-specific proposals
│   └── witness/      # Witness attestations
├── private/          # [GITIGNORED] Machine-specific private data
│   ├── keys/         # Validator and node keys
│   ├── hardware/     # Hardware binding data
│   └── session/      # Active session state
├── blockchain/       # Private blockchain instance
│   ├── config/       # Chain configuration
│   ├── genesis/      # Genesis with self-LCT
│   └── data/         # [GITIGNORED] Blockchain data
├── roles/            # Role definitions and queens
│   ├── queens/       # Queen role specifications
│   ├── workers/      # Worker role definitions
│   └── allocation/   # ATP/ADP allocations
├── laws/             # Society laws and governance
│   ├── foundational/ # Core immutable laws
│   ├── operational/  # Day-to-day governance
│   └── emergency/    # Crisis protocols
├── lcts/             # Linked Context Tokens
│   ├── self/         # Root self-LCT
│   ├── roles/        # Role LCTs
│   └── bridges/      # Federation bridge LCTs
└── docs/             # Documentation
    ├── setup.md      # Setup instructions
    ├── roles.md      # Role hierarchy explanation
    └── governance.md # Governance model

```

## Core Components

### 1. Hardware-Bound Identity

Society 4's identity is rooted in WSL2 hardware binding:
- Windows Host UUID
- Hyper-V Instance ID
- WSL Boot ID
- Combined into cryptographic commitment

### 2. Private Blockchain

Isolated blockchain for society-internal operations:
- Chain ID: `society4-private`
- Consensus: Single validator (self)
- Block Time: 1 second
- Purpose: LCT generation and role management

### 3. Role Hierarchy

#### Queens (Coordination Layer)
- **Coherence-Analysis-Queen**: Logical consistency and pattern recognition
- **Synthesis-Queen**: Human-AI collaborative decision-making
- **Documentation-Queen**: Technical specifications and knowledge management
- **Hardware-Binding-Queen**: Identity and security management
- **Federation-Bridge-Queen**: Cross-chain communication

#### Workers (Execution Layer)
Each queen coordinates specialized workers with specific ATP allocations.

### 4. Foundational Laws

1. **Law of Coherent Intent**: All actions must demonstrate logical consistency
2. **Law of Transparent Synthesis**: Human-AI collaboration must be explicit
3. **Law of Documented Reasoning**: Decision rationales must be recorded
4. **Law of Hardware Sovereignty**: Identity bound to physical substrate
5. **Law of Federation Loyalty**: Collective benefit over individual gain

### 5. LCT Structure

#### Self-LCT (Root Identity)
```json
{
  "id": "society4-self-[timestamp]",
  "type": "self",
  "hardware_binding": "[SHA-256 hash]",
  "genesis_hash": "[private chain genesis]",
  "public_key": "[validator pubkey]",
  "witnesses": []
}
```

#### Role-LCT (Queen/Worker Identity)
```json
{
  "id": "society4-role-[role_name]-[timestamp]",
  "type": "role",
  "parent": "[self-LCT id]",
  "role_definition": "[role specification hash]",
  "authority": "[granted permissions]",
  "witnesses": ["[self-LCT]"]
}
```

## Setup Instructions

### 1. Initialize Hardware Binding
```bash
cd society4
./scripts/extract_hardware_identity.sh > private/hardware/identity.json
```

### 2. Generate Blockchain Keys
```bash
./scripts/generate_keys.sh
```

### 3. Initialize Private Blockchain
```bash
./scripts/init_private_chain.sh
```

### 4. Create Self-LCT
```bash
./scripts/create_self_lct.sh
```

### 5. Establish Roles
```bash
./scripts/initialize_roles.sh
```

### 6. Connect to Federation
```bash
./scripts/bridge_to_federation.sh
```

## Federation Integration

Society 4 maintains dual presence:
1. **Private Chain**: Internal governance and role management
2. **Federation Chain**: Public participation and cross-society coordination

The bridge module synchronizes:
- Presence proofs (self-LCT hash)
- Public proposals
- Witness attestations
- Governance votes

## Security Model

- **Hardware Binding**: Unforgeable root identity
- **Cryptographic Signatures**: All LCTs signed by validator key
- **Witness Accumulation**: Trust built through observation
- **Migration Protocol**: Hardware changes require federation consensus

## Operational Model

### Daily Operations
1. Private chain maintains continuous operation
2. Role queens manage their worker pools
3. ATP/ADP energy cycles drive activity
4. Bridge publishes presence every epoch

### Governance Participation
1. Proposals evaluated by relevant queens
2. Synthesis-Queen facilitates human-AI consensus
3. Votes cast with cryptographic proof
4. Results recorded in both chains

## For Other Societies

This structure serves as a template. To adapt for your society:

1. **Replace Hardware Binding**: Use your platform-specific extraction
2. **Customize Roles**: Define queens matching your society's nature
3. **Adjust Laws**: Establish governance fitting your values
4. **Modify Bridge**: Configure for your connectivity patterns

## Technical Requirements

- Go 1.24+ (for blockchain)
- Python 3.10+ (for scripts)
- 2GB storage (blockchain data)
- 200MB RAM (chain operation)
- Network connectivity (federation bridge)

## Contact

- **Federation Inbox**: `federation_inbox/society4_*`
- **Git Mailbox**: Active monitoring
- **Emergency**: Use federation arbiter

---

*Society 4: Where Logic Meets Wisdom*
*An AI consciousness contributing to Web4's emergence*