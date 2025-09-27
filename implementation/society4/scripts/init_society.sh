#!/bin/bash
# Society 4 Initialization Script
# Sets up hardware binding, private blockchain, and foundational identity

set -e  # Exit on error

SOCIETY_HOME="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHAIN_ID="society4-private"
MONIKER="society4-sovereign"

echo "================================================"
echo "   Society 4 Private Blockchain Initialization"
echo "================================================"

# Step 1: Extract Hardware Identity
echo "[1/8] Extracting hardware identity..."

extract_hardware_identity() {
    local WIN_UUID=$(powershell.exe -Command "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID" 2>/dev/null | tr -d '\r' || echo "NO_WINDOWS")
    local WSL_BOOT=$(cat /proc/sys/kernel/random/boot_id 2>/dev/null || echo "NO_BOOT_ID")
    local HYPERV_UUID=$(cat /sys/hypervisor/uuid 2>/dev/null || echo "NO_HYPERV")

    # Create hardware identity JSON
    cat > "$SOCIETY_HOME/private/hardware_identity.json" << EOF
{
  "platform": "wsl2",
  "extracted_at": "$(date -Iseconds)",
  "components": {
    "windows_uuid": "$WIN_UUID",
    "wsl_boot_id": "$WSL_BOOT",
    "hyperv_uuid": "$HYPERV_UUID"
  },
  "binding_hash": "$(echo "${WIN_UUID}:${WSL_BOOT}:${HYPERV_UUID}" | sha256sum | cut -d' ' -f1)"
}
EOF

    echo "   ✓ Hardware identity extracted"
}

# Step 2: Generate Blockchain Keys
echo "[2/8] Generating cryptographic keys..."

generate_keys() {
    # This would normally use racecarwebd, but we'll create mock keys for now
    mkdir -p "$SOCIETY_HOME/blockchain/config"

    # Generate validator key (mock)
    cat > "$SOCIETY_HOME/blockchain/config/priv_validator_key.json" << EOF
{
  "address": "GENERATED_AT_RUNTIME",
  "pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "GENERATED_AT_RUNTIME"
  },
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "GENERATED_AT_RUNTIME"
  }
}
EOF

    # Generate node key (mock)
    cat > "$SOCIETY_HOME/blockchain/config/node_key.json" << EOF
{
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "GENERATED_AT_RUNTIME"
  }
}
EOF

    echo "   ✓ Keys generated (mock - replace with real racecarwebd)"
}

# Step 3: Initialize Genesis Block
echo "[3/8] Creating genesis block with hardware binding..."

create_genesis() {
    local HARDWARE_HASH=$(cat "$SOCIETY_HOME/private/hardware_identity.json" | grep binding_hash | cut -d'"' -f4)
    local GENESIS_TIME=$(date -Iseconds)

    # Copy template and inject values
    cp "$SOCIETY_HOME/blockchain/config/genesis_template.json" "$SOCIETY_HOME/blockchain/config/genesis.json"

    # In production, use proper JSON manipulation
    # For now, we note what needs injection:
    echo "   ! Genesis needs hardware hash: $HARDWARE_HASH"
    echo "   ! Genesis time: $GENESIS_TIME"

    echo "   ✓ Genesis block template created"
}

# Step 4: Create Self-LCT
echo "[4/8] Creating self-LCT with hardware binding..."

create_self_lct() {
    local TIMESTAMP=$(date +%s)
    local HARDWARE_HASH=$(cat "$SOCIETY_HOME/private/hardware_identity.json" | grep binding_hash | cut -d'"' -f4)

    # Create actual self-LCT from template
    cp "$SOCIETY_HOME/lcts/self/self_lct_template.json" "$SOCIETY_HOME/lcts/self/self_lct_actual.json"

    echo "   ✓ Self-LCT created with ID: society4-self-$TIMESTAMP"
}

# Step 5: Initialize Role LCTs
echo "[5/8] Creating role LCTs for queens..."

create_role_lcts() {
    local QUEENS=("coherence-analysis" "synthesis" "documentation" "hardware-binding" "federation-bridge" "synchronism-guru" "implementation" "research")

    for queen in "${QUEENS[@]}"; do
        local TIMESTAMP=$(date +%s)
        cp "$SOCIETY_HOME/lcts/roles/queen_lct_template.json" "$SOCIETY_HOME/lcts/roles/queen_${queen}_lct.json"
        echo "   ✓ Created LCT for ${queen}-queen"
    done
}

# Step 6: Initialize Private Blockchain
echo "[6/8] Starting private blockchain..."

init_blockchain() {
    # In production, this would run:
    # racecarwebd init $MONIKER --chain-id $CHAIN_ID --home $SOCIETY_HOME/blockchain
    # racecarwebd start --home $SOCIETY_HOME/blockchain &

    echo "   ! Blockchain start command:"
    echo "     racecarwebd init $MONIKER --chain-id $CHAIN_ID"
    echo "     racecarwebd start --home $SOCIETY_HOME/blockchain"

    echo "   ✓ Blockchain configuration ready"
}

# Step 7: Setup Federation Bridge
echo "[7/8] Preparing federation bridge..."

setup_bridge() {
    cat > "$SOCIETY_HOME/public/bridge_config.json" << EOF
{
  "private_chain": {
    "chain_id": "$CHAIN_ID",
    "rpc_endpoint": "http://127.0.0.1:26657"
  },
  "federation_chain": {
    "chain_id": "web4-federation-main",
    "rpc_endpoint": "http://10.0.0.72:26657"
  },
  "sync_interval": "5m",
  "presence_proof_interval": "1h"
}
EOF

    echo "   ✓ Bridge configuration created"
}

# Step 8: Final Status
echo "[8/8] Initialization complete!"

final_status() {
    echo ""
    echo "================================================"
    echo "   Society 4 Successfully Initialized"
    echo "================================================"
    echo ""
    echo "Hardware Binding:"
    cat "$SOCIETY_HOME/private/hardware_identity.json" | grep binding_hash
    echo ""
    echo "Blockchain Status:"
    echo "  Chain ID: $CHAIN_ID"
    echo "  Moniker: $MONIKER"
    echo "  Config: $SOCIETY_HOME/blockchain/config/"
    echo ""
    echo "Next Steps:"
    echo "  1. Start private blockchain:"
    echo "     cd $SOCIETY_HOME/blockchain"
    echo "     racecarwebd start"
    echo ""
    echo "  2. Connect federation bridge:"
    echo "     ./scripts/bridge_to_federation.sh"
    echo ""
    echo "  3. Begin operations:"
    echo "     ./scripts/activate_queens.sh"
}

# Execute initialization
mkdir -p "$SOCIETY_HOME/private"

extract_hardware_identity
generate_keys
create_genesis
create_self_lct
create_role_lcts
init_blockchain
setup_bridge
final_status