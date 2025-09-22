# Fix Instructions for Sprout (Jetson) Build Issues

## Problem Identified
The binary on sprout lacks `init`, `keys`, and `genesis` commands because of module naming inconsistencies.

## Working Configuration on CBP

### What CBP Has:
1. **Two cmd directories** (both exist):
   - `/cmd/racecar-webd/` (with hyphen) - OLD, uses module "racecar-web"
   - `/cmd/racecarwebd/` (no hyphen) - NEW, uses module "racecarweb"

2. **Module name in go.mod**: `racecarweb` (no hyphen)

3. **Working binary has ALL commands**:
   - init, keys, genesis, start, export, tx, query, etc.

## The Fix for Sprout

### Step 1: Pull Latest Changes
```bash
cd /home/sprout/ai-workspace/ACT
git pull
```

### Step 2: Check Module Name
```bash
head -5 implementation/ledger/go.mod
# Should show: module racecarweb
```

### Step 3: Use the CORRECT cmd Directory
The new `cmd/racecarwebd/` (no hyphen) structure that was just pulled has:
- `main.go` with correct imports: `"racecarweb/app"` and `"racecarweb/cmd/racecarwebd/cmd"`
- Full command structure in `cmd/` subdirectory:
  - `commands.go` - All command definitions
  - `config.go` - Configuration commands
  - `root.go` - Root command with all subcommands
  - `testnet.go` - Testnet commands
  - `testnet_multi_node.go` - Multi-node setup

### Step 4: Build with Correct Path
```bash
cd implementation/ledger

# Clean any old builds
rm -f ~/go/bin/racecar-webd
rm -f ~/.go/bin/racecar-webd

# Build from the NEW cmd directory (no hyphen)
go build -o ~/go/bin/racecarwebd ./cmd/racecarwebd

# Create symlink for compatibility
ln -s ~/go/bin/racecarwebd ~/go/bin/racecar-webd
```

### Step 5: Verify Commands
```bash
racecarwebd 2>&1 | grep -A 20 "Available Commands:"
```

Should show:
- init
- keys
- genesis
- start
- tx
- query
- etc.

## Key Differences Found

| Aspect | Sprout (Broken) | CBP (Working) |
|--------|----------------|---------------|
| Module name | Confused (racecar-web vs racecarweb) | `racecarweb` |
| CMD directory | Missing proper structure | Both old and new exist |
| Imports in main.go | Wrong module path | Correct: `racecarweb/...` |
| Binary commands | Only basic (start, export) | Full suite (init, keys, genesis, etc.) |

## The Root Cause

The issue is that sprout was trying to build from an incomplete or incorrectly named cmd structure. The latest pull includes the proper `cmd/racecarwebd/` directory with all necessary command implementations.

## Quick Test After Fix

```bash
# Initialize a test node to verify all commands work
racecarwebd init test-node --chain-id test

# Check if keys command works
racecarwebd keys list

# Check genesis commands
racecarwebd genesis --help
```

If these work, sprout is ready to join the federation!