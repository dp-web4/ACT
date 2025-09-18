# WSL2 Machine Context - ACT Blockchain

## Last Updated: January 17, 2025

### System Configuration
- **WSL2 Memory**: 12GB RAM + 32GB Swap
- **Platform**: Linux 6.6.87.2-microsoft-standard-WSL2
- **Go Version**: 1.24.0 (with sonic library replace directives)
- **Ignite CLI**: v29.4.0-dev
- **Cosmos SDK**: v0.53.0

### Critical Fixes Applied

#### Memory Management
- WSL2 was exhausting memory causing system crashes
- Fixed with `.wslconfig`:
  ```
  [wsl2]
  memory=12GB
  swap=32GB
  ```
- User feedback: "disabling seems like a bullet to the head solution" (about disabling OpenAPI)

#### Go 1.24 Compatibility
- bytedance/sonic library incompatibility
- Fixed with replace directives in go.mod
- All protoc plugins installed and working

#### Module Registration
- Fixed "type_url has not been registered yet" panics
- Implemented RegisterInterfaces in all codec.go files
- Fixed import paths (racecarweb → racecar-web)

### Current Status: ✅ FULLY OPERATIONAL

- Blockchain builds successfully
- ATP/ADP energy cycle implemented
- Society pools working
- 95% Web4 compliant

### Key Implementation Files

#### Society Pool Infrastructure
- `x/energycycle/keeper/society_pool.go` - Complete pool management
- `proto/racecarweb/energycycle/v1/society_pool.proto` - Pool structure
- `x/energycycle/module/autocli.go` - Fixed CLI parameters

#### Energy Operations
- MintADP: Treasury role mints to society pool
- DischargeATP: Workers convert ATP→ADP for work
- RechargeADP: Producers convert ADP→ATP with energy

### Build Commands
```bash
# Generate protobuf
ignite generate proto-go

# Build blockchain
ignite chain build --skip-proto

# Initialize and start
racecar-webd init mynode --chain-id racecarweb
racecar-webd keys add alice --keyring-backend test
racecar-webd genesis add-genesis-account alice 1000000000stake --keyring-backend test
racecar-webd genesis gentx alice 100000000stake --keyring-backend test --chain-id racecarweb
racecar-webd genesis collect-gentxs
racecar-webd start --api.enable --grpc.enable
```

### Test Accounts
- **Alice**: cosmos1pmnw5epy2zflns3lrmzxgcm86zsly7nr24jqrp (100M stake + 20k tokens)
- **Bob**: cosmos1m8jcll5nn036hpgktn7pcndljqe5jhc8ujznjk (100M stake + 10k tokens)

### Running Services
- **Tendermint RPC**: http://0.0.0.0:26657
- **REST API**: http://0.0.0.0:1317
- **Token Faucet**: http://0.0.0.0:4500

### Blockchain Performance
- "Idle isn't" - uses 18% CPU while doing "nothing"
- This insight led to metabolic states for digital organisms
- Maintenance costs are real for living systems

### Context Insights from Implementation

#### The Entrepreneur's Paradox
- Created elaborate swarm architecture
- Did all work manually anyway
- This was OPTIMAL, not a failure
- Full context privilege > distributed execution
- Delegation cost often > execution cost

#### Swarm Value
- Perspective and validation > execution
- Web4-Compliance-Queen caught 95% compliance (not 42%)
- Multiple viewpoints prevent tunnel vision

#### What We Built (8 hours)
- Complete society-owned token pools
- Energy conservation (ATP ⟷ ADP)
- Role-based operations
- Metabolic states framework
- 95% Web4 compliant

### Historical Note
Successfully ran on July 13, 2025 with Ignite, showing alice/bob accounts operational.
Rebuilt January 17, 2025 with complete ATP/ADP energy economy.

### Files to Preserve
- `/test-energy-cycle.sh` - Integration test script
- `/ENERGY_CYCLE_IMPLEMENTATION.md` - Complete documentation
- `/WEB4_COMPLIANCE_REPORT.md` - 95% compliance validation
- `/CONTEXT_INSIGHTS.md` - Delegation paradox insights