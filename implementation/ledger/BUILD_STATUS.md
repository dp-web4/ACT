# ACT Blockchain Build Status

## Date: January 17, 2025

## Overview
Working to restore the ACT Web4 blockchain to a buildable and runnable state. The blockchain successfully ran on July 13, 2025 using Ignite CLI, as evidenced by historical logs showing alice and bob accounts operating on chain ID "racecarweb".

## Issues Encountered

### 1. Go 1.24 Compatibility Issue
**Problem:** Go 1.24 introduces breaking changes that cause linking errors with bytedance/sonic library:
```
link: github.com/bytedance/sonic/loader: invalid reference to runtime.lastmoduledatap
```

**Solution Applied:** Added replace directives in go.mod:
```go
replace (
    github.com/bytedance/sonic => github.com/bytedance/sonic v1.12.1
    github.com/bytedance/sonic/loader => github.com/bytedance/sonic/loader v0.2.0
)
```
**Status:** ✅ RESOLVED - Binary builds successfully

### 2. Missing Module Registration
**Problem:** All custom modules (trusttensor, energycycle, lctmanager, etc.) had empty RegisterInterfaces functions, causing panic on startup:
```
panic: type_url /racecarweb.trusttensor.v1.MsgUpdateParams has not been registered yet
```

**Solution Applied:** Fixed codec.go for all modules to properly register message types:
```go
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
  registry.RegisterImplementations((*sdk.Msg)(nil),
    &MsgUpdateParams{},
  )
  msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
```
**Status:** ✅ RESOLVED - Modules register correctly

### 3. Import Path Issues
**Problem:** MRH module files incorrectly imported "racecarweb" instead of "racecar-web"

**Solution Applied:** Fixed all import statements in:
- x/mrh/keeper/keeper.go
- x/mrh/keeper/witness.go
- tests/unit/web4/phase1_test.go

**Status:** ✅ RESOLVED

### 4. Genesis Validator Configuration
**Problem:** Chain fails to start with "validator set is empty after InitGenesis"

**Partial Solution:**
- Created test accounts (alice, bob)
- Added genesis accounts with stake
- Created gentx for alice as validator
- Collected gentxs into genesis

**Status:** ⚠️ PARTIALLY RESOLVED - Genesis configured but chain startup still failing

### 5. OpenAPI Generation Failure
**Problem:** Ignite serve fails during proto generation with:
```
failed to generate openapi spec
plugin go tool github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2: signal: killed
```

**Solution Progress:**
- Installed missing protoc plugins:
  - `protoc-gen-openapiv2@latest`
  - `protoc-gen-grpc-gateway@latest`
  - `protoc-gen-go@latest`
  - `protoc-gen-go-grpc@latest`
- Proto generation (`ignite generate proto-go`) now works successfully

**Status:** ⚠️ PARTIALLY RESOLVED - Proto generation works but full serve still causes system instability

## Current Build Status

### What Works:
- ✅ Go 1.24 environment configured
- ✅ Binary builds successfully with `ignite chain build --skip-proto`
- ✅ All modules compile without errors
- ✅ Genesis initialization works
- ✅ Test accounts created
- ✅ Proto generation works with `ignite generate proto-go`
- ✅ All protoc plugins installed and working

### What Doesn't Work:
- ❌ Full `ignite chain serve` causes system instability/crashes
- ❌ Direct `racecar-webd start` fails with validator issues
- ❌ System becomes unstable during heavy Ignite operations (possible WSL2 memory issue)

## Next Steps Required

1. **Fix OpenAPI Generation:**
   - Try increasing memory limits
   - Consider disabling OpenAPI generation temporarily
   - May need to downgrade to Go 1.22/1.23 as sonic suggests

2. **Alternative Startup Methods:**
   - Try `ignite chain serve --skip-proto`
   - Manually configure validator and start without Ignite
   - Use minimal configuration first

3. **Verification Needed:**
   - Ensure all proto files are properly formatted
   - Check if buf.work.yaml conflicts need resolution
   - Validate genesis.json structure

## Commands That Work

```bash
# Install protoc plugins (required first)
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate proto files
ignite generate proto-go

# Build the binary
ignite chain build --skip-proto

# Initialize node
racecar-webd init mynode --chain-id racecarweb

# Add accounts
racecar-webd keys add alice --keyring-backend test
racecar-webd keys add bob --keyring-backend test

# Add genesis accounts
racecar-webd genesis add-genesis-account alice 1000000000stake --keyring-backend test
racecar-webd genesis add-genesis-account bob 1000000000stake --keyring-backend test

# Create validator
racecar-webd genesis gentx alice 100000000stake --keyring-backend test --chain-id racecarweb
racecar-webd genesis collect-gentxs

# Set minimum gas price
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ~/.racecar-web/config/app.toml
```

## Environment Details
- Go Version: 1.24.0
- Ignite CLI: v29.4.0-dev
- Cosmos SDK: v0.53.0
- Platform: WSL2 Linux
- Working Directory: /mnt/c/exe/projects/ai-agents/ACT/implementation/ledger