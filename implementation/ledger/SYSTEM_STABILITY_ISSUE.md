# System Stability Issue with Ignite CLI

## Problem
Running `ignite chain serve` causes system crashes/instability on WSL2. The system hard crashes requiring restart.

## Symptoms
- Process gets killed with "signal: killed"
- System becomes unresponsive
- Claude Code app crashes with "process exited with code 1"
- Requires machine restart

## Possible Causes
1. **Memory pressure in WSL2** - Ignite operations may be consuming all available memory
2. **Resource limits** - WSL2 has default memory limits that may be too low
3. **Go 1.24 incompatibility** - Despite fixes, Go 1.24 may have deeper issues with Ignite

## Workarounds Attempted
1. ✅ Installed all protoc plugins manually
2. ✅ Proto generation works with `ignite generate proto-go`
3. ✅ Building works with `ignite chain build --skip-proto`
4. ❌ Full `ignite chain serve` still crashes

## Potential Solutions to Try

### 1. Increase WSL2 Memory
Create or edit `%USERPROFILE%\.wslconfig`:
```ini
[wsl2]
memory=8GB
processors=4
swap=4GB
```

### 2. Use Direct Binary Instead of Ignite
```bash
# Build
ignite chain build --skip-proto

# Then start directly
racecar-webd start --api.enable --grpc.enable
```

### 3. Downgrade to Go 1.23
Despite sonic warnings, Go 1.23 might be more stable with Ignite.

### 4. Run in Docker Container
Isolate the blockchain in a container with controlled resources.

## Current Status
- Proto generation: ✅ WORKING
- Binary build: ✅ WORKING
- Direct start: ❌ VALIDATOR ISSUES
- Ignite serve: ❌ CAUSES CRASHES

## Next Steps
1. Try increasing WSL2 memory allocation
2. Fix validator configuration for direct start
3. Consider downgrading Go version if instability persists