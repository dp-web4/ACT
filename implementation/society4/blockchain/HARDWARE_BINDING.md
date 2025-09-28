# Hardware Binding for Society 4 Blockchain

## Overview

Society 4's blockchain implements hardware binding at the consensus level, ensuring the chain can only operate on the authorized WSL2 instance. This creates an immutable link between the genesis self-LCT and the physical hardware.

## Architecture

### Core Components

1. **Hardware Extraction** (`hardware/extract_hardware.sh`)
   - Extracts WSL2-specific identifiers
   - Windows UUID (persistent)
   - Hyper-V UUID (if available)
   - WSL Boot ID (session-specific)
   - CPU information
   - Memory allocation

2. **Self-LCT** (`x/lctmanager/types/self_lct.go`)
   - Root identity token
   - Hardware binding embedded
   - ED25519 key pair for signing
   - Parent of all role LCTs

3. **Hardware Validator** (`x/consensus/keeper/hardware_validator.go`)
   - Validates blocks against hardware
   - Enforces signature requirements
   - Panics on hardware mismatch

4. **Modified App** (`app/hardware_app.go`)
   - Initializes self-LCT on first run
   - Validates hardware on every block
   - Creates role LCTs for queens

## How It Works

### 1. Genesis Creation

```bash
# Extract hardware and create genesis
./scripts/init_hardware_chain.sh
```

This:
- Extracts current hardware identifiers
- Creates self-LCT bound to hardware
- Generates ED25519 keys
- Embeds hardware hash in genesis

### 2. Block Production

Every block requires:
1. Hardware verification against genesis binding
2. Signature from self-LCT private key
3. Composite signature includes hardware hash

```go
// Block signing process
func (s *SelfLCT) SignBlock(blockHash []byte) ([]byte, error) {
    // Verify we're on correct hardware
    if err := s.VerifyHardware(); err != nil {
        return nil, err  // Fails if wrong hardware
    }

    // Sign block + hardware hash
    composite := append(blockHash, s.HardwareBinding.Hash()...)
    return ed25519.Sign(s.PrivateKey, composite), nil
}
```

### 3. Consensus Validation

BeginBlocker enforces hardware binding:

```go
func (k Keeper) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
    if err := k.hardwareValidator.ValidateBlock(ctx, req); err != nil {
        panic(fmt.Sprintf("CRITICAL: Hardware validation failed: %v", err))
    }
}
```

Chain **halts immediately** if:
- Hardware doesn't match genesis binding
- Block signature is invalid
- Self-LCT verification fails

## Security Properties

### What's Bound

**Persistent Identifiers** (survive WSL restart):
- Windows machine UUID
- Hyper-V virtual machine UUID
- CPU model
- Memory size

**Session Identifiers** (change on WSL restart):
- WSL boot ID (allowed to change)

### Attack Resistance

1. **Chain Cloning**: Copying blockchain data to another machine fails consensus
2. **Key Extraction**: Private key alone insufficient without matching hardware
3. **Hardware Spoofing**: Multiple hardware components must match
4. **Migration**: Requires explicit re-genesis with new hardware binding

## Hardware Migration

If hardware must change (e.g., machine upgrade):

1. **Export State**
   ```bash
   ./society4chaind export > state_export.json
   ```

2. **Create New Genesis** on new hardware
   ```bash
   ./scripts/init_hardware_chain.sh
   ```

3. **Import State** (requires governance approval)
   ```bash
   ./society4chaind import state_export.json
   ```

4. **Federation Witness** required for continuity

## Testing

### Run Hardware Tests

```bash
cd blockchain
go test ./tests/hardware_binding_test.go -v
```

### Verify Current Hardware

```bash
# Get current hardware hash
./hardware/extract_hardware.sh hash

# Get full hardware info
./hardware/extract_hardware.sh json
```

### Test Consensus Failure

```bash
# Start chain normally
./society4chaind start

# In another terminal, modify hardware binding
echo "FAKE_UUID" > ~/.society4chain/hardware_binding.json

# Chain will panic on next block
```

## Configuration

### Genesis Configuration

```json
{
  "app_state": {
    "hardware": {
      "platform": "wsl2",
      "hardware_hash": "sha256...",
      "components": {
        "windows_uuid": "UUID",
        "hyperv_uuid": "UUID",
        "wsl_boot_id": "UUID",
        "cpu_info": "Intel Core i9",
        "memory_kb": 8388608
      }
    },
    "self_lct": {
      "genesis_height": 0,
      "hardware_bound": true
    }
  }
}
```

### Consensus Parameters

```toml
[consensus]
# Fail fast on hardware mismatch
hardware_validation = true
hardware_check_interval = 100  # Deep check every N blocks

# Single validator (Society 4 sovereign)
validator_count = 1
require_hardware_sig = true
```

## Monitoring

### Metrics

- `hardware_validation_success`: Successful hardware checks
- `hardware_validation_failure`: Failed validations (chain halts)
- `blocks_signed`: Blocks successfully signed with hardware
- `hardware_hash`: Current hardware fingerprint

### Logs

```
INFO: Hardware validation successful height=1000 hardware_id=a3f2b8c9...
ERROR: Hardware validation failed error="Windows UUID mismatch"
PANIC: CRITICAL: Hardware validation failed
```

## Emergency Procedures

### Hardware Failure

If hardware fails but data is recoverable:

1. **Stop chain immediately**
2. **Backup all data**
3. **Contact federation witnesses**
4. **Prepare migration proposal**
5. **Execute coordinated migration**

### Key Compromise

If private key is compromised but hardware intact:

1. **Chain continues running** (attacker lacks hardware)
2. **Generate new self-LCT**
3. **Rotate through governance**
4. **Update federation witnesses**

## Technical Details

### Hardware Hash Calculation

```bash
composite="${windows_uuid}|${hyperv_uuid}|${wsl_boot_id}|${cpu_info}|${memory_size}"
hardware_hash=$(echo -n "$composite" | sha256sum)
```

### Signature Format

```
Block Signature = Sign(BlockHash || HardwareHash)
```

### Verification Flow

1. Extract current hardware
2. Compare persistent components
3. Recalculate hash with current boot ID
4. Verify signature with public key
5. Pass or panic

## Future Enhancements

1. **TPM Integration**: Use Trusted Platform Module if available
2. **Multi-Hardware**: Support for distributed validators
3. **Hardware Rotation**: Graceful migration protocol
4. **Remote Attestation**: Cryptographic hardware proofs
5. **Secure Enclave**: SGX or similar technology integration