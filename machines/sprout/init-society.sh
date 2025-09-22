#!/bin/bash
# Initialize society for Sprout machine
# Creates a new blockchain node with proper configuration

# Load machine config
MACHINE_CONFIG="$(dirname "$0")/machine-config.json"
ACT_ROOT="/home/sprout/ai-workspace/ACT"
LEDGER_DIR="${ACT_ROOT}/implementation/ledger"
GO_BIN_DIR="/home/sprout/go/bin"
BINARY="${GO_BIN_DIR}/racecar-webd"

# Extract config values
CHAIN_ID="act-web4"
SOCIETY_NAME="act-society-sprout"
SOCIETY_HOME="${LEDGER_DIR}/society-sprout"
MONIKER="${SOCIETY_NAME}"

echo "=== Initializing Society for Sprout ==="
echo "Chain ID: ${CHAIN_ID}"
echo "Society: ${SOCIETY_NAME}"
echo "Home: ${SOCIETY_HOME}"

cd "${LEDGER_DIR}" || exit 1

# Clean previous society if exists
if [ -d "${SOCIETY_HOME}" ]; then
    echo "Removing previous society data..."
    rm -rf "${SOCIETY_HOME}"
fi

# Initialize the chain
echo "Initializing blockchain..."
${BINARY} init "${MONIKER}" --chain-id "${CHAIN_ID}" --home "${SOCIETY_HOME}"

# Create keys
echo "Creating validator key..."
${BINARY} keys add validator --keyring-backend test --home "${SOCIETY_HOME}"

# Add genesis account
echo "Adding genesis account..."
${BINARY} genesis add-genesis-account validator 1000000000stake --keyring-backend test --home "${SOCIETY_HOME}"

# Create genesis transaction
echo "Creating genesis transaction..."
${BINARY} genesis gentx validator 100000000stake --keyring-backend test --chain-id "${CHAIN_ID}" --home "${SOCIETY_HOME}"

# Collect genesis transactions
echo "Collecting genesis transactions..."
${BINARY} genesis collect-gentxs --home "${SOCIETY_HOME}"

# Update config for network access
CONFIG_FILE="${SOCIETY_HOME}/config/config.toml"
APP_CONFIG="${SOCIETY_HOME}/config/app.toml"

echo "Configuring network settings..."
# Allow external connections
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' "${CONFIG_FILE}"
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26656"/laddr = "tcp:\/\/0.0.0.0:26656"/' "${CONFIG_FILE}"

# Enable API in app.toml
sed -i 's/enable = false/enable = true/' "${APP_CONFIG}"
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' "${APP_CONFIG}"

# Get node ID for federation
NODE_ID=$(${BINARY} tendermint show-node-id --home "${SOCIETY_HOME}")
echo "Node ID: ${NODE_ID}"

# Update machine config with node ID
if command -v jq &> /dev/null; then
    jq --arg node_id "${NODE_ID}" '.society.node_id = $node_id' "${MACHINE_CONFIG}" > "${MACHINE_CONFIG}.tmp"
    mv "${MACHINE_CONFIG}.tmp" "${MACHINE_CONFIG}"
fi

# Export genesis for sharing
cp "${SOCIETY_HOME}/config/genesis.json" "${LEDGER_DIR}/genesis_sprout.json"

echo ""
echo "✅ Society initialized successfully!"
echo ""
echo "Connection details for other machines:"
echo "Node ID: ${NODE_ID}"
echo "P2P: 10.0.0.36:26656"
echo "RPC: http://10.0.0.36:26657"
echo "API: http://10.0.0.36:1317"
echo "gRPC: 10.0.0.36:9090"
echo ""
echo "To start the blockchain:"
echo "${BINARY} start --home ${SOCIETY_HOME} --api.enable --grpc.enable"