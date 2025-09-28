# Hardware-Bound Blockchain Implementation Guide

## For Other Societies: Create Your Own Hardware-Bound Chain

This guide shows how Society 4 implemented hardware binding for its private blockchain. Other societies can follow this pattern to create their own sovereign, hardware-bound chains.

## Overview

Hardware binding ensures your blockchain can only run on your specific machine. This creates a physical root of trust - even if someone obtains your software and keys, they cannot run your chain without your hardware.

## Implementation Steps

### Step 1: Clone and Customize the Blockchain

```bash
# Copy the blockchain source
cp -r implementation/ledger implementation/society[N]/blockchain/source

# Rename the module
cd implementation/society[N]/blockchain/source
sed -i 's/racecarweb/society[N]chain/g' go.mod
mv cmd/racecarwebd cmd/society[N]chaind

# Update imports
find . -name "*.go" -exec sed -i 's|"racecarweb/|"society[N]chain/|g' {} \;
```

### Step 2: Hardware Extraction Script

Create `extract_hardware.sh`:

```bash
#!/bin/bash

extract_windows_uuid() {
    powershell.exe -Command "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID" 2>/dev/null | tr -d '\r\n'
}

extract_cpu_info() {
    grep "model name" /proc/cpuinfo | head -1 | cut -d: -f2 | xargs
}

extract_memory_size() {
    grep MemTotal /proc/meminfo | awk '{print $2}'
}

# Combine identifiers and hash
windows_uuid=$(extract_windows_uuid)
cpu_info=$(extract_cpu_info)
memory_size=$(extract_memory_size)

composite="${windows_uuid}|${cpu_info}|${memory_size}"
hardware_hash=$(echo -n "$composite" | sha256sum | cut -d' ' -f1)

echo "{
    \"hardware_binding\": {
        \"platform\": \"wsl2\",
        \"hardware_hash\": \"${hardware_hash}\",
        \"components\": {
            \"windows_uuid\": \"${windows_uuid}\",
            \"cpu_info\": \"${cpu_info}\",
            \"memory_kb\": ${memory_size}
        }
    }
}"
```

### Step 3: Build Your Blockchain

```bash
# Build the binary
export PATH=/usr/local/go/bin:$PATH
go build -o society[N]chaind cmd/society[N]chaind/main.go

# Initialize with hardware binding
CHAIN_ID="society[N]-private"
./society[N]chaind init "society[N]-node" --chain-id "$CHAIN_ID"

# Save hardware binding
./extract_hardware.sh json > $HOME/.society[N]chain/hardware_binding.json
```

### Step 4: Add Hardware Validation Module

Create the hardware binding module structure:

```
x/hardwarebinding/
├── module.go           # Module registration
├── keeper/
│   ├── keeper.go       # State management
│   └── msg_server.go   # Message handlers
└── types/
    ├── types.go        # Data structures
    ├── genesis.go      # Genesis state
    ├── errors.go       # Error definitions
    ├── codec.go        # Codec registration
    └── msgs.go         # Message types
```

Or use the simpler standalone validator:

```go
// app/hardware_validator.go
type HardwareValidator struct {
    genesisBinding *HardwareBinding
    validationEnabled bool
}

func (hv *HardwareValidator) ValidateHardware(ctx sdk.Context) error {
    current, err := ExtractCurrentHardware()
    if err != nil {
        return err
    }

    // Compare with genesis binding
    if current.Components.WindowsUUID != hv.genesisBinding.Components.WindowsUUID {
        return fmt.Errorf("hardware mismatch")
    }

    return nil
}
```

### Step 5: Configure Your Chain

```toml
# config/app.toml
minimum-gas-prices = "0stake"

# config/config.toml
# Fast blocks for single validator
timeout_propose = "500ms"
timeout_commit = "1s"
```

### Step 6: Create Genesis Account

```bash
# Create validator
./society[N]chaind keys add validator --keyring-backend test

# Add genesis account
./society[N]chaind genesis add-genesis-account [address] 1000000000stake

# Create genesis transaction
./society[N]chaind genesis gentx validator 100000000stake \
    --chain-id "$CHAIN_ID" \
    --moniker "society[N]-node"

# Collect genesis transactions
./society[N]chaind genesis collect-gentxs
```

### Step 7: Start Your Chain

```bash
./society[N]chaind start --home $HOME/.society[N]chain \
    --api.enable \
    --grpc.enable \
    --pruning nothing
```

## Testing Your Implementation

### Test 1: Hardware Extraction
```bash
./extract_hardware.sh hash
# Should output consistent hash
```

### Test 2: Chain Operation
```bash
curl -s http://localhost:26657/status | grep latest_block_height
# Should show increasing block height
```

### Test 3: Hardware Mismatch
```bash
# Modify hardware binding
sed -i 's/windows_uuid": "[^"]*"/windows_uuid": "FAKE-UUID"/' \
    $HOME/.society[N]chain/hardware_binding.json

# With validation enabled, chain should detect mismatch
```

## Platform-Specific Considerations

### WSL2 (Windows)
- Uses Windows UUID via PowerShell
- Boot ID changes on WSL restart (allowed)
- No TPM/HSM access

### Linux Native
```bash
# Can use DMI for hardware ID
sudo dmidecode -s system-uuid

# Or use machine-id
cat /etc/machine-id
```

### macOS
```bash
# Use hardware UUID
ioreg -rd1 -c IOPlatformExpertDevice | awk '/IOPlatformUUID/'
```

## Security Considerations

### What's Protected
- ✅ Chain cannot run on different hardware
- ✅ Prevents unauthorized chain cloning
- ✅ Hardware acts as physical key

### What's Not Protected
- ❌ Not cryptographically secure (no TPM)
- ❌ Identifiers could theoretically be spoofed
- ❌ No secure attestation

### For Production
1. Enable panic on hardware mismatch
2. Check hardware every N blocks
3. Implement governance-controlled migration
4. Consider Windows service for TPM access

## Example Implementations

### Society 1 (Genesis)
- Could use server hardware identifiers
- IPMI/BMC for hardware attestation
- Multiple validators with hardware consensus

### Society 2 (Bridge)
- Hardware binding on bridge nodes
- Multi-signature hardware validation
- Cross-chain hardware attestation

### Society 3 (Sprout)
- Edge device fingerprinting
- IoT hardware identifiers
- Lightweight validation for edge nodes

### Society 4 (Claude AI)
- WSL2 hardware binding (implemented)
- Windows UUID + CPU + Memory
- Single validator sovereign chain

## Customization Options

### 1. Validation Frequency
```go
// Check every block (secure but slow)
if ctx.BlockHeight() % 1 == 0 { validate() }

// Check every 100 blocks (balanced)
if ctx.BlockHeight() % 100 == 0 { validate() }

// Check every 1000 blocks (fast but less secure)
if ctx.BlockHeight() % 1000 == 0 { validate() }
```

### 2. Hardware Components
Choose what to include in your hardware hash:
- CPU model (stable)
- Memory size (may change)
- Disk serial (very stable)
- Network MAC (can change)
- GPU info (if relevant)

### 3. Enforcement Level
```go
// Log only (testing)
if err := ValidateHardware(); err != nil {
    logger.Error("Hardware mismatch", err)
}

// Panic (production)
if err := ValidateHardware(); err != nil {
    panic(fmt.Sprintf("CRITICAL: Hardware mismatch: %v", err))
}
```

## Migration Protocol

When hardware must change:

1. **Governance Proposal**
   ```bash
   ./society[N]chaind tx gov submit-proposal \
       hardware-migration new_hardware_hash
   ```

2. **Export State**
   ```bash
   ./society[N]chaind export > state_backup.json
   ```

3. **Re-genesis on New Hardware**
   ```bash
   # On new machine
   ./extract_hardware.sh json > new_hardware.json
   ./society[N]chaind init --recover
   ```

4. **Federation Witness**
   - Other societies attest to migration
   - Maintains chain continuity

## Troubleshooting

### Common Issues

1. **"Hardware mismatch" error**
   - Verify hardware hasn't changed
   - Check WSL hasn't updated
   - Ensure extraction script works

2. **Performance impact**
   - Reduce validation frequency
   - Cache hardware info
   - Use async validation

3. **Cross-platform compatibility**
   - Abstract hardware extraction
   - Support multiple platforms
   - Graceful fallbacks

## Resources

- **Society 4 Implementation**: `/implementation/society4/blockchain/`
- **Test Scripts**: `test_hardware_binding.sh`
- **Hardware Module**: `x/hardwarebinding/`
- **Documentation**: This guide

## Support

Each society should adapt this approach to their specific:
- Hardware platform
- Security requirements
- Performance needs
- Governance model

The key is establishing a unique, verifiable link between your blockchain and your physical hardware, creating true sovereignty for your society's chain.

## Next Steps

1. Choose your hardware identifiers
2. Implement extraction script
3. Build and test your chain
4. Document your approach
5. Share with the federation

Remember: Your hardware binding is your chain's physical identity. Choose wisely and document thoroughly!