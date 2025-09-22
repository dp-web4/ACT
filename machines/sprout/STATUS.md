# Sprout Machine Status - ACT Blockchain Setup

## Current State (Sep 21, 2025)

### ✅ Completed
1. **Machine-adaptable structure created**
   - Created `machines/` directory with subdirectories for sprout, legion, cbp
   - Full configuration in `machines/sprout/machine-config.json`
   - Build and deployment scripts created

2. **Go 1.24.0 installed**
   - Downloaded and installed at `/usr/local/go`
   - Required for Web4 blockchain compilation

3. **Blockchain binary built**
   - Binary exists at `/home/sprout/go/bin/racecar-webd` (155MB)
   - Built successfully but missing some commands (init, keys, genesis)

### ⚠️ Issues Encountered

1. **Binary Command Structure**
   - The built binary lacks `init`, `keys`, and `genesis` commands
   - Can only run `start`, `export`, `rollback`, etc.
   - This prevents proper blockchain initialization

2. **Module Name Discrepancy**
   - Original: `racecar-web`
   - Current: `racecarweb`
   - This may be causing build issues

3. **CMD Directory Structure**
   - Ignite expects `cmd/racecar-webd/`
   - We have `cmd/racecarwebd/`
   - Created symlink but still having issues

### 🔍 Key Findings

1. **Working on Other Machines**
   - The blockchain is running successfully on Legion and CBP machines
   - They must have the proper cmd structure that we're missing

2. **Ignite CLI Issues**
   - Getting "sonic only supports go1.17~1.23" warnings with Go 1.24
   - Some Ignite commands failing with cancel reader errors on ARM64

3. **Files Modified**
   ```
   modified:   go.mod (changed module name)
   modified:   go.sum
   created:    cmd/racecarwebd/ (incomplete structure)
   ```

### 📝 Next Steps

1. **Option A: Get proper cmd structure from working machines**
   - Pull the actual `cmd/` directory from Legion or CBP
   - These have working init/keys/genesis commands

2. **Option B: Use Ignite chain serve**
   - Let Ignite handle the initialization
   - May need to work around the cancel reader issue

3. **Option C: Manual initialization**
   - Create config files manually
   - Use the binary's limited commands

### 🔧 Technical Details

- **Architecture**: ARM64 (aarch64)
- **OS**: Ubuntu 22.04 (JetPack 6.2.1)
- **Go Version**: 1.24.0
- **Binary Path**: `/home/sprout/go/bin/racecar-webd`
- **Society Home**: `/home/sprout/ai-workspace/ACT/implementation/ledger/society-sprout`
- **Network Config**:
  - IP: 10.0.0.36
  - P2P: 26656
  - RPC: 26657
  - API: 1317
  - gRPC: 9090

### 💭 Hypothesis

The other machines likely:
1. Had Ignite properly scaffold the full cmd structure
2. OR manually created the proper cmd files with all commands
3. AND committed these to a branch we haven't pulled

We need to either:
- Get the working cmd structure from another machine
- Fix Ignite to work properly on ARM64
- Manually add the missing command implementations