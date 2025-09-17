# Claude Context for ACT Web4 Blockchain

## Project Context System

**IMPORTANT**: A comprehensive context system exists at `/mnt/c/projects/ai-agents/misc/context-system/`

Quick access:
```bash
# Get overview of this project's role
cd /mnt/c/projects/ai-agents/misc/context-system
python3 query_context.py project web4-modbatt

# See blockchain integration plans
python3 query_context.py search "digital twin"

# Find related projects
cat /mnt/c/projects/ai-agents/misc/context-system/relationships/blockchain-bridge.md
```

## This Project's Role

ACT Web4 blockchain implements the Artificial Communication Transport protocol:
- Built on Cosmos SDK v0.53.0 with custom modules
- Implements Web4 protocol with LCT (Linked Context Tokens)
- Features ATP/ADP energy economy with discharge/recharge mechanics
- Provides T3/V3 tensor attributions for trust and value

## Key Modules
- `lctmanager`: Manages Linked Context Tokens (LCTs) - core Web4 identity
- `energycycle`: ATP/ADP energy trading with discharge/recharge mechanics
- `trusttensor`: T3/V3 trust and value tensor calculations
- `mrh`: Markov Relevancy Horizon for context boundaries
- `pairing`: Device pairing and authentication
- `componentregistry`: Component tracking and verification
- `pairingqueue`: Queue management for pairing operations

## Current Status (Jan 17, 2025)

### ✅ Completed
- Fixed Go 1.24 compatibility issues with sonic library
- Fixed module registration (RegisterInterfaces was empty)
- Fixed import paths (racecarweb vs racecar-web)
- Binary builds successfully with `ignite chain build --skip-proto`
- Genesis configuration with test accounts (alice, bob)
- All custom modules compile without errors

### ⚠️ In Progress
- Chain startup issues with validator configuration
- OpenAPI generation crashes (possible memory issue)
- Proto generation stability

### 📝 Known Issues
- Ignite serve fails on OpenAPI generation (protoc-gen-openapiv2 killed)
- Direct startup needs proper validator configuration
- May need to use Go 1.22/1.23 for full sonic compatibility

## Quick Commands

```bash
# Build blockchain
ignite chain build --skip-proto

# Initialize and start (after build)
racecar-webd init mynode --chain-id racecarweb
racecar-webd keys add alice --keyring-backend test
racecar-webd genesis add-genesis-account alice 1000000000stake --keyring-backend test
racecar-webd genesis gentx alice 100000000stake --keyring-backend test --chain-id racecarweb
racecar-webd genesis collect-gentxs
racecar-webd start --api.enable --grpc.enable
```

## Historical Note
Successfully ran on July 13, 2025 with Ignite, showing alice/bob accounts operational.

## Build Requirements
- Go 1.24.0 (with sonic replace directives in go.mod)
- Ignite CLI v29.4.0-dev
- Cosmos SDK v0.53.0

## Key Files Modified Today
- All module codec.go files (fixed RegisterInterfaces)
- go.mod (added sonic replace directives)
- x/mrh/keeper/*.go (fixed import paths)

See BUILD_STATUS.md for detailed issue tracking and solutions.