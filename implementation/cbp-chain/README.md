# CBP Hardware-Bound Private Chain

## Overview

This is a machine-specific, hardware-bound blockchain instance for the CBP machine. The chain uses hardware attestation to ensure it can only run on this specific WSL2 instance.

## Hardware Binding

The chain is cryptographically bound to:
- Hostname: cbp
- CPU configuration
- Memory configuration
- WSL2 instance UUID
- User environment

This binding ensures:
1. **Machine Specificity**: Chain only runs on THIS machine
2. **Hardware Attestation**: Validates hardware fingerprint on startup
3. **Key Derivation**: Validator keys derived from hardware entropy
4. **Tamper Resistance**: Changes to hardware invalidate the chain

## Directory Structure

```
cbp-chain/
├── config/              # Chain configuration
│   └── hardware_binding.json
├── data/               # Blockchain data directory
│   ├── config/
│   └── data/
├── keys/               # Hardware-derived keys
│   ├── machine_id.json
│   └── hardware_info.txt
├── logs/               # Chain logs
├── init_cbp_chain.sh   # Initialize the chain
├── start_cbp_chain.sh  # Start the chain
└── README.md          # This file
```

## Setup Instructions

### 1. Initialize the Chain

First time setup - generates hardware fingerprint and initializes blockchain:

```bash
bash init_cbp_chain.sh
```

This will:
- Extract hardware fingerprint
- Generate deterministic validator keys
- Create genesis configuration
- Set up chain data directory

### 2. Start the Chain

```bash
bash start_cbp_chain.sh
```

The chain will:
- Verify hardware attestation
- Load hardware-bound validator key
- Start consensus engine
- Begin producing blocks

## Chain Details

- **Chain Type**: Web4 Society Blockchain (Cosmos SDK based)
- **Consensus**: Tendermint BFT
- **Modules**: All Web4 modules (LCTManager, ComponentRegistry, etc.)
- **Hardware Binding**: Patent-compliant split-key encryption
- **Attestation**: Continuous hardware verification

## Security Features

1. **Hardware Lock**: Chain data encrypted with hardware-derived keys
2. **Machine Binding**: Cannot be copied to another machine
3. **Attestation Chain**: Every block includes hardware attestation
4. **Split-Key Protection**: Validator keys split between hardware and software

## API Endpoints

Once running, the chain exposes:

- **RPC**: http://localhost:26657
- **REST**: http://localhost:1317
- **gRPC**: localhost:9090
- **WebSocket**: ws://localhost:26657/websocket

## Web4 Integration

This chain implements the full Web4 stack:
- LCT (Linked Context Token) management
- Hardware-bound component registry
- Trust tensor calculations
- Energy cycle (ATP/ADP) tracking
- Society governance

## Troubleshooting

### Chain Won't Start
- Verify hardware hasn't changed
- Check `keys/machine_id.json` exists
- Ensure no other chain running on same ports

### Hardware Mismatch
- Hardware changes invalidate the chain
- Must reinitialize if hardware changes
- Backup data before hardware upgrades

### Port Conflicts
- Default ports: 26657 (RPC), 1317 (REST), 9090 (gRPC)
- Modify in `data/config/config.toml` if needed

## Important Notes

⚠️ **This chain is bound to THIS SPECIFIC MACHINE**
- Cannot be migrated to another machine
- Hardware changes may require reinitialization
- Backup strategies must account for hardware binding

## License

Part of the ACT/Web4 implementation - AGPL-3.0