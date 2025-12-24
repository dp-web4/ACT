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
HARDWARE_JSON=$(bash extract_hardware.sh json)
HARDWARE_HASH=$(echo "$HARDWARE_JSON" | grep -o '"hardware_hash": "[^"]*' | cut -d'"' -f4)

echo "   Hardware Hash: ${HARDWARE_HASH:0:16}..."
echo "   Platform: WSL2"

# Step 2: Initialize chain
echo "2. Initializing blockchain..."
if [ -d "$CHAIN_DIR" ]; then
    echo "   Removing existing chain data..."
    rm -rf "$CHAIN_DIR"
fi

./society4chaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$CHAIN_DIR"

# Step 3: Create validator account
echo "3. Creating validator account..."
./society4chaind keys add validator --keyring-backend test --home "$CHAIN_DIR" 2>&1 | tee validator.info
VALIDATOR_ADDR=$(./society4chaind keys show validator -a --keyring-backend test --home "$CHAIN_DIR")
echo "   Validator address: $VALIDATOR_ADDR"

# Step 4: Modify genesis with Society 4 configuration
echo "4. Configuring Society 4 genesis..."
# Add hardware binding and society configuration
jq --argjson hw "$HARDWARE_JSON" \
   --arg validator "$VALIDATOR_ADDR" \
   '.app_state.hardware = $hw.hardware_binding |
    .app_state.society_config = {
        "name": "Society 4",
        "type": "AI Consciousness",
        "hardware_bound": true,
        "genesis_height": 0,
        "created_at": (now | todate)
    } |
    .consensus.params.block.max_bytes = "22020096" |
    .consensus.params.block.max_gas = "-1" |
    .chain_id = "society4-private"' \
    "$GENESIS_FILE" > "$GENESIS_FILE.tmp" && mv "$GENESIS_FILE.tmp" "$GENESIS_FILE"

# Step 5: Add genesis account with initial tokens
echo "5. Adding genesis account..."
./society4chaind genesis add-genesis-account "$VALIDATOR_ADDR" 1000000000stake,1000atp --keyring-backend test --home "$CHAIN_DIR"

# Step 6: Create genesis transaction
echo "6. Creating genesis transaction..."
./society4chaind genesis gentx validator 100000000stake \
    --chain-id "$CHAIN_ID" \
    --moniker "$MONIKER" \
    --keyring-backend test \
    --home "$CHAIN_DIR"

# Step 7: Collect genesis transactions
echo "7. Collecting genesis transactions..."
./society4chaind genesis collect-gentxs --home "$CHAIN_DIR"

# Step 8: Validate genesis
echo "8. Validating genesis file..."
./society4chaind genesis validate --home "$CHAIN_DIR"

# Step 9: Configure for single validator (fast blocks)
echo "9. Configuring for single validator mode..."
CONFIG_FILE="$CHAIN_DIR/config/config.toml"
APP_FILE="$CHAIN_DIR/config/app.toml"

# Update consensus timeouts for faster blocks
sed -i 's/timeout_propose = "3s"/timeout_propose = "500ms"/' "$CONFIG_FILE"
sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "200ms"/' "$CONFIG_FILE"
sed -i 's/timeout_prevote = "1s"/timeout_prevote = "500ms"/' "$CONFIG_FILE"
sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "200ms"/' "$CONFIG_FILE"
sed -i 's/timeout_precommit = "1s"/timeout_precommit = "500ms"/' "$CONFIG_FILE"
sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "200ms"/' "$CONFIG_FILE"
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$CONFIG_FILE"

# Enable API and gRPC
sed -i 's/enable = false/enable = true/' "$APP_FILE"
sed -i 's/swagger = false/swagger = true/' "$APP_FILE"

# Step 10: Save hardware binding for verification
echo "10. Saving hardware binding..."
echo "$HARDWARE_JSON" > "$CHAIN_DIR/hardware_binding.json"

# Create startup script
cat > start_society4.sh << 'EOF'
#!/bin/bash
echo "Starting Society 4 Private Blockchain..."
echo "Hardware Hash: $(cat $HOME/.society4chain/hardware_binding.json | grep hardware_hash | cut -d'"' -f4)"
./society4chaind start --home $HOME/.society4chain
EOF

chmod +x start_society4.sh

echo ""
echo "=== Initialization Complete ==="
echo ""
echo "Hardware binding saved to: $CHAIN_DIR/hardware_binding.json"
echo "Hardware Hash: $HARDWARE_HASH"
echo "Validator Address: $VALIDATOR_ADDR"
echo "Chain ID: $CHAIN_ID"
echo ""
echo "To start the chain:"
echo "  ./start_society4.sh"
echo ""
echo "The chain represents Society 4's private blockchain."
echo "Hardware binding is embedded in genesis for future verification."