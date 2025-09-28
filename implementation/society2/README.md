# Society2 - Democratic Web4 Society

## Overview

Society2 is an alternative Web4 society implementation designed to explore interoperability between dissimilar private chains. While Society4 uses hierarchical governance, Society2 implements a flat democratic model, enabling us to test cross-chain federation between different governance paradigms.

## Key Differences from Society4

| Aspect | Society4 | Society2 |
|--------|----------|----------|
| **Governance** | Hierarchical (Queen/roles) | Democratic (1-citizen-1-vote) |
| **Consensus** | Proof-of-Stake | Proof-of-Contribution |
| **Energy Model** | Traditional ATP | Renewable ATP sources |
| **Federation** | Selective | Open by default |
| **Decision Making** | Role-based authority | Majority consensus |
| **Trust Model** | Reputation-based | Contribution-based |

## Architecture

```
society2/
├── blockchain/         # Independent blockchain source
│   └── source/        # Modified from society4 for democratic governance
├── laws/              # Democratic constitution
│   └── society2_constitution.md
├── keys/              # Society-specific keys
│   └── society2_identity.json
├── config/            # Genesis and network config
├── federation_config.json  # Inter-chain setup
├── init_society2.sh   # Initialize society2
└── README.md         # This file
```

## Network Configuration

Society2 runs on different ports to allow simultaneous operation with Society4:

- **P2P**: 26556 (Society4: 26656)
- **RPC**: 26557 (Society4: 26657)
- **API**: 1217 (Society4: 1317)
- **gRPC**: 9091 (Society4: 9090)

## Setup Instructions

### 1. Initialize Society2

```bash
cd /mnt/c/exe/projects/ai-agents/ACT/implementation/society2
bash init_society2.sh
```

### 2. Start Society2 Blockchain

```bash
cd blockchain/source
./society2d start --home=./data
```

### 3. Verify Society2 is Running

```bash
# Check status
curl -s http://localhost:26557/status | jq .result.node_info.network

# Should return: "web4-society2-001"
```

## Inter-Chain Communication

Society2 is designed to federate with Society4, demonstrating:

### 1. Governance Bridge
- Democratic proposals from Society2 can be considered by Society4's hierarchy
- Hierarchical decisions from Society4 require democratic ratification in Society2

### 2. Energy Exchange
- Different ATP generation models (renewable vs traditional)
- Exchange rate mechanisms for cross-chain energy transfers
- Proof-of-contribution rewards vs stake-based rewards

### 3. Trust Tensor Sharing
- Cross-chain reputation calculations
- Federated trust scores
- Inter-society relationship tracking

### 4. LCT Pairing
- Citizens can maintain identity across both societies
- Hardware binding remains society-specific
- Cross-chain component attestation

## Federation Setup

To connect Society2 with Society4:

### 1. Start Both Chains
```bash
# Terminal 1 - Society4
cd ../society4/blockchain/source
./society4d start --home=./data

# Terminal 2 - Society2
cd ../society2/blockchain/source
./society2d start --home=./data
```

### 2. Create IBC Connection
```bash
# Create client on Society2 pointing to Society4
./society2d tx ibc connection create-client \
  --chain-id web4-society4-001 \
  --node tcp://localhost:26657

# Create client on Society4 pointing to Society2
./society4d tx ibc connection create-client \
  --chain-id web4-society2-001 \
  --node tcp://localhost:26557
```

### 3. Establish Channels
```bash
# Token transfer channel
./society2d tx ibc channel open \
  --port transfer \
  --counterparty-port transfer

# Trust tensor channel
./society2d tx ibc channel open \
  --port trusttensor \
  --counterparty-port trusttensor
```

## Use Cases

### 1. Democratic Proposal System
Citizens submit proposals that require majority approval rather than hierarchical sign-off.

### 2. Contribution-Based Mining
Citizens earn ATP through verified contributions rather than staking.

### 3. Open Federation
Any compatible society can join without permission, following democratic vote.

### 4. Renewable Energy
ATP generation tied to sustainable practices and renewable sources.

## Testing Interoperability

### Test Scenarios

1. **Cross-Chain Token Transfer**
   - Transfer ATP from Society2 to Society4
   - Verify exchange rate application

2. **Federated Governance**
   - Submit proposal in Society2
   - Require acknowledgment from Society4

3. **Trust Tensor Synchronization**
   - Update trust score in Society2
   - Verify propagation to Society4

4. **LCT Migration**
   - Create LCT in Society2
   - Pair with equivalent in Society4

## Development Notes

### Customization Points

The blockchain source in `blockchain/source/` has been modified from Society4:

1. **Governance Module**: Changed from role-based to vote-based
2. **Consensus**: Modified validator selection for contribution tracking
3. **Energy Cycle**: Added renewable source tracking
4. **Federation**: Enabled by default in genesis

### Building from Source

```bash
cd blockchain/source
go mod download
go build -tags society2 -o society2d ./app
```

### Configuration Files

- `keys/society2_identity.json` - Hardware-bound society identity
- `laws/society2_constitution.md` - Democratic governance rules
- `federation_config.json` - Inter-chain communication setup
- `config/genesis_template.json` - Genesis configuration

## Monitoring

### Check Society2 Status
```bash
curl http://localhost:26557/status
```

### Query Governance Proposals
```bash
curl http://localhost:1217/cosmos/gov/v1/proposals
```

### Monitor Federation
```bash
curl http://localhost:1217/ibc/core/channel/v1/channels
```

## Troubleshooting

### Port Conflicts
Ensure Society4 and Society2 use different ports as configured.

### Build Failures
Check Go version (1.21+) and ensure all dependencies are downloaded.

### Federation Issues
Verify both chains are running and IBC relayer is active.

## Future Enhancements

1. **Quadratic Voting**: Implement more sophisticated democratic mechanisms
2. **Liquid Democracy**: Allow vote delegation
3. **DAO Treasury**: Automated fund management
4. **Cross-Chain Governance**: Unified proposal system
5. **Energy Markets**: Automated ATP/ADP trading

## License

Part of the ACT/Web4 implementation - AGPL-3.0