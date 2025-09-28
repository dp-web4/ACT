#!/bin/bash

# Initialize Society 4 blockchain with hardware binding

set -e

CHAIN_ID="society4-private"
MONIKER="society4-node"
CHAIN_DIR="$HOME/.society4chain"
GENESIS_FILE="$CHAIN_DIR/config/genesis.json"

echo "=== Society 4 Hardware-Bound Blockchain Initialization ==="

# Step 1: Extract hardware information
echo "1. Extracting hardware identifiers..."
HARDWARE_JSON=$(./hardware/extract_hardware.sh json)
HARDWARE_HASH=$(echo "$HARDWARE_JSON" | jq -r '.hardware_binding.hardware_hash')

echo "   Hardware Hash: ${HARDWARE_HASH:0:16}..."
echo "   Platform: WSL2"

# Step 2: Initialize chain
echo "2. Initializing blockchain..."
if [ -d "$CHAIN_DIR" ]; then
    echo "   Removing existing chain data..."
    rm -rf "$CHAIN_DIR"
fi

./society4chaind init "$MONIKER" --chain-id "$CHAIN_ID"

# Step 3: Create validator account
echo "3. Creating validator account..."
./society4chaind keys add validator --keyring-backend test 2>&1 | tee validator.info
VALIDATOR_ADDR=$(./society4chaind keys show validator -a --keyring-backend test)

# Step 4: Modify genesis with hardware binding
echo "4. Embedding hardware binding in genesis..."
jq --argjson hw "$HARDWARE_JSON" \
   --arg validator "$VALIDATOR_ADDR" \
   '.app_state.hardware = $hw.hardware_binding |
    .app_state.self_lct = {
        "genesis_height": 0,
        "created_at": (now | todate),
        "hardware_bound": true
    } |
    .app_state.rolegovernance = {
        "queens": [
            {"name": "Treasury-Queen", "atp_budget": 120, "active": true},
            {"name": "Law-Oracle-Queen", "atp_budget": 110, "active": true},
            {"name": "Implementation-Queen", "atp_budget": 100, "active": true},
            {"name": "Research-Queen", "atp_budget": 95, "active": true},
            {"name": "Documentation-Queen", "atp_budget": 85, "active": true},
            {"name": "Federation-Bridge-Queen", "atp_budget": 100, "active": true},
            {"name": "Coherence-Analysis-Queen", "atp_budget": 90, "active": true},
            {"name": "Quality-Assurance-Queen", "atp_budget": 80, "active": true},
            {"name": "Security-Queen", "atp_budget": 110, "active": true},
            {"name": "Emergency-Response-Queen", "atp_budget": 110, "active": true}
        ],
        "total_atp": 1000
    }' \
    "$GENESIS_FILE" > "$GENESIS_FILE.tmp" && mv "$GENESIS_FILE.tmp" "$GENESIS_FILE"

# Step 5: Add genesis account with initial tokens
echo "5. Adding genesis account..."
./society4chaind genesis add-genesis-account "$VALIDATOR_ADDR" 1000000stake,1000atp --keyring-backend test

# Step 6: Create genesis transaction
echo "6. Creating genesis transaction..."
./society4chaind genesis gentx validator 1000000stake \
    --chain-id "$CHAIN_ID" \
    --moniker "$MONIKER" \
    --keyring-backend test

# Step 7: Collect genesis transactions
echo "7. Collecting genesis transactions..."
./society4chaind genesis collect-gentxs

# Step 8: Validate genesis
echo "8. Validating genesis file..."
./society4chaind genesis validate

# Step 9: Configure for single validator
echo "9. Configuring for single validator mode..."
sed -i 's/timeout_propose = "3s"/timeout_propose = "500ms"/' "$CHAIN_DIR/config/config.toml"
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$CHAIN_DIR/config/config.toml"
sed -i 's/create_empty_blocks = true/create_empty_blocks = true/' "$CHAIN_DIR/config/config.toml"
sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "1s"/' "$CHAIN_DIR/config/config.toml"

# Step 10: Save hardware binding for verification
echo "10. Saving hardware binding..."
echo "$HARDWARE_JSON" > "$CHAIN_DIR/hardware_binding.json"

echo ""
echo "=== Initialization Complete ==="
echo ""
echo "Hardware binding saved to: $CHAIN_DIR/hardware_binding.json"
echo "Hardware Hash: $HARDWARE_HASH"
echo "Validator Address: $VALIDATOR_ADDR"
echo ""
echo "To start the chain:"
echo "  ./society4chaind start"
echo ""
echo "The chain will only run on this specific hardware."
echo "Any attempt to run on different hardware will fail consensus."