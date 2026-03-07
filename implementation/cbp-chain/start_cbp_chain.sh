#!/bin/bash

# Start CBP Hardware-Bound Chain (ACT Blockchain)
# Chain data lives at ~/.racecarweb (Cosmos SDK default home)
# Endpoints: RPC :26657, REST :1317, gRPC :9090 (all on 0.0.0.0)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

if [ ! -f "keys/machine_id.json" ]; then
    echo "Error: Machine ID not found. Run init_cbp_chain.sh first"
    exit 1
fi

# Get chain info
CHAIN_ID=$(grep chain_id keys/machine_id.json | cut -d'"' -f4)
HARDWARE_ID=$(grep machine_id keys/machine_id.json | cut -d'"' -f4)
CBP_IP=$(hostname -I 2>/dev/null | awk '{print $1}')

echo "Starting CBP ACT Chain"
echo "  Chain ID:    ${CHAIN_ID}"
echo "  Hardware ID: ${HARDWARE_ID:0:16}..."
echo "  RPC:         http://${CBP_IP}:26657"
echo "  REST:        http://${CBP_IP}:1317"
echo "  gRPC:        ${CBP_IP}:9090"
echo ""

# Ensure validator state file exists
mkdir -p "$HOME/.racecarweb/data"
if [ ! -f "$HOME/.racecarweb/data/priv_validator_state.json" ]; then
    echo '{"height":"0","round":0,"step":0}' > "$HOME/.racecarweb/data/priv_validator_state.json"
fi

# Start the chain (uses default home ~/.racecarweb)
exec ./racecarwebd start