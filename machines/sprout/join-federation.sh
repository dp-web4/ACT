#!/bin/bash
# Join federation script for Sprout machine
# Connects to other society nodes to form a Web4 federation

# Configuration
ACT_ROOT="/home/sprout/ai-workspace/ACT"
LEDGER_DIR="${ACT_ROOT}/implementation/ledger"
GO_BIN_DIR="/home/sprout/go/bin"
BINARY="${GO_BIN_DIR}/racecar-webd"
SOCIETY_HOME="${LEDGER_DIR}/society-sprout"
CONFIG_FILE="${SOCIETY_HOME}/config/config.toml"

echo "=== Joining Web4 Federation from Sprout ==="

# Function to add peer
add_peer() {
    local PEER_ID=$1
    local PEER_ADDR=$2
    local PEER_NAME=$3
    
    echo "Adding peer: ${PEER_NAME}"
    echo "  ID: ${PEER_ID}"
    echo "  Address: ${PEER_ADDR}"
    
    # Add to persistent peers
    PEER_STRING="${PEER_ID}@${PEER_ADDR}"
    
    # Check if peer already exists
    if grep -q "${PEER_ID}" "${CONFIG_FILE}"; then
        echo "  Peer already configured"
    else
        # Add to persistent_peers
        sed -i "/persistent_peers =/s/\"\"/\"${PEER_STRING}\"/" "${CONFIG_FILE}"
        sed -i "/persistent_peers =/s/\"\(.*\)\"/\"\1,${PEER_STRING}\"/" "${CONFIG_FILE}"
        echo "  ✅ Peer added"
    fi
}

# Known federation members (update as machines join)
echo "Configuring known federation members..."

# Legion (if available)
# add_peer "LEGION_NODE_ID" "10.0.0.72:26656" "Legion-RTX4090"

# CBP (if available)  
# add_peer "CBP_NODE_ID" "10.0.0.XX:26656" "CBP-RTX2060"

# Get our own node info
OUR_NODE_ID=$(${BINARY} tendermint show-node-id --home "${SOCIETY_HOME}")
echo ""
echo "Our Node Info:"
echo "  ID: ${OUR_NODE_ID}"
echo "  P2P: 10.0.0.36:26656"
echo ""

# Check current peers
echo "Checking network status..."
if curl -s http://localhost:26657/net_info > /dev/null 2>&1; then
    PEER_COUNT=$(curl -s http://localhost:26657/net_info | jq -r '.result.n_peers')
    echo "Current peers: ${PEER_COUNT}"
    
    if [ "${PEER_COUNT}" -gt 0 ]; then
        echo "Connected peers:"
        curl -s http://localhost:26657/net_info | jq -r '.result.peers[].node_info.moniker'
    fi
else
    echo "Blockchain not running. Start it first with:"
    echo "${ACT_ROOT}/machines/sprout/start-blockchain.sh"
    exit 1
fi

echo ""
echo "To complete federation:"
echo "1. Share your node ID with other machines:"
echo "   ${OUR_NODE_ID}@10.0.0.36:26656"
echo ""
echo "2. Add other machines' peers by editing this script"
echo ""
echo "3. Restart blockchain to connect:"
echo "   kill $(cat ${LEDGER_DIR}/society-sprout.pid)"
echo "   ${ACT_ROOT}/machines/sprout/start-blockchain.sh"