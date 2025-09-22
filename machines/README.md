# ACT Machine-Specific Configurations

## Overview

Each machine in the ACT federation has unique paths, configurations, and network settings. This directory contains machine-specific adaptations that allow ACT to run seamlessly across different environments.

## Directory Structure

```
machines/
├── common/          # Shared utilities and scripts
├── sprout/          # Jetson Orin Nano (ARM64)
├── legion/          # Legion Pro 7 with RTX 4090
├── cbp/            # Windows WSL2 with RTX 2060 SUPER
└── README.md       # This file
```

## Machine Configurations

### Sprout (Jetson Orin Nano)
- **Architecture**: ARM64
- **OS**: Ubuntu 22.04 (JetPack)
- **IP**: 10.0.0.36
- **Path**: `/home/sprout/ai-workspace/ACT`
- **Special**: Edge device, lower resources

### Legion (RTX 4090)
- **Architecture**: x86_64
- **OS**: Windows 11 + WSL2
- **IP**: 10.0.0.72 (typical)
- **Path**: `/mnt/c/exe/projects/ai-agents/ACT`
- **Special**: High performance GPU

### CBP (RTX 2060 SUPER)
- **Architecture**: x86_64
- **OS**: Windows 11 + WSL2
- **IP**: 10.0.0.XX
- **Path**: `/mnt/c/projects/ai-agents/ACT`
- **Special**: Development machine

## Usage

### On Each Machine

1. **Build the blockchain**:
```bash
bash machines/[MACHINE_NAME]/build.sh
```

2. **Initialize society**:
```bash
bash machines/[MACHINE_NAME]/init-society.sh
```

3. **Start blockchain**:
```bash
bash machines/[MACHINE_NAME]/start-blockchain.sh
```

4. **Join federation**:
```bash
bash machines/[MACHINE_NAME]/join-federation.sh
```

## Federation Setup

### Step 1: Initialize First Machine
```bash
# On Sprout
bash machines/sprout/build.sh
bash machines/sprout/init-society.sh
bash machines/sprout/start-blockchain.sh
# Note the Node ID displayed
```

### Step 2: Share Genesis
```bash
# Copy genesis file from first machine
scp sprout@10.0.0.36:/path/to/genesis_sprout.json ./
```

### Step 3: Initialize Other Machines
```bash
# On Legion
bash machines/legion/build.sh
bash machines/legion/init-society.sh
# Replace genesis with shared one
cp genesis_sprout.json implementation/ledger/society-legion/config/genesis.json
bash machines/legion/start-blockchain.sh
```

### Step 4: Connect Peers
Edit `join-federation.sh` on each machine with peer node IDs, then run:
```bash
bash machines/[MACHINE_NAME]/join-federation.sh
```

## Machine Config Format

Each machine has a `machine-config.json`:
```json
{
  "machine_id": "unique-identifier",
  "paths": {
    "home": "/home/username",
    "workspace": "/path/to/workspace",
    "act_root": "/path/to/ACT"
  },
  "network": {
    "ip": "10.0.0.XX",
    "p2p_port": 26656,
    "rpc_port": 26657
  },
  "society": {
    "name": "act-society-name",
    "node_id": "generated-on-init"
  }
}
```

## Adding a New Machine

1. Create directory: `machines/[NEW_MACHINE]/`
2. Copy template scripts from `machines/common/` (if available) or another machine
3. Create `machine-config.json` with correct paths
4. Update scripts with machine-specific paths
5. Test build and initialization
6. Add machine to federation

## Troubleshooting

### Build Failures
- Check Go version (needs 1.23+)
- Verify paths in machine-config.json
- Check architecture compatibility (ARM vs x86)

### Network Issues
- Ensure ports are open (26656, 26657, 1317, 9090)
- Check firewall settings
- Verify IP addresses are reachable

### Federation Problems
- Genesis files must match exactly
- Chain ID must be the same
- Node IDs must be correct in peer configuration

## Federation Philosophy

Each machine society is autonomous but gains strength through federation:
- **Autonomy**: Each runs its own blockchain node
- **Cooperation**: Societies share state and witness events
- **Resilience**: Network continues if machines drop
- **Evolution**: Societies can propose and vote on changes

This mirrors Web4's vision: independent entities forming voluntary associations for mutual benefit.