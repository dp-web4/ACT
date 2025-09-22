#!/bin/bash

# Start script for CBP (Society 2) blockchain
# This machine runs Society 2 in the ACT federation

echo "========================================="
echo "Starting ACT Blockchain - Society 2 (CBP)"
echo "========================================="

# Configuration
MACHINE_NAME="cbp"
SOCIETY_NAME="act-society-2"
CHAIN_ID="act-web4"
HOME_DIR="./society2"
BINARY="racecar-webd"

# Check if already running
if [ -f society2.pid ]; then
    PID=$(cat society2.pid)
    if ps -p $PID > /dev/null 2>&1; then
        echo "❌ Blockchain already running with PID $PID"
        echo "To restart, run: kill $PID && rm society2.pid"
        exit 1
    fi
fi

# Navigate to ledger directory
cd /mnt/c/exe/projects/ai-agents/ACT/implementation/ledger || {
    echo "❌ Failed to navigate to ledger directory"
    exit 1
}

# Start the blockchain
echo "🚀 Starting blockchain..."
echo "   Network: Society 2 (Federation Member)"
echo "   Peer: Society 1 @ 10.0.0.72"
echo "   Ports: P2P=26666, RPC=26667, API=1318, gRPC=9091"

$BINARY start \
    --home $HOME_DIR \
    --p2p.laddr tcp://0.0.0.0:26666 \
    --rpc.laddr tcp://0.0.0.0:26667 \
    --grpc.address 0.0.0.0:9091 \
    --api.address tcp://0.0.0.0:1318 \
    --api.enable \
    > society2.log 2>&1 &

# Save PID
echo $! > society2.pid
PID=$!

# Wait a moment for startup
sleep 5

# Check if running
if ps -p $PID > /dev/null 2>&1; then
    echo "✅ Blockchain started with PID $PID"
    echo "   Log: tail -f society2.log"
    echo ""

    # Check federation status
    sleep 3
    PEERS=$(curl -s http://localhost:26667/net_info 2>/dev/null | grep -o '"n_peers":"[^"]*' | cut -d'"' -f4)
    if [ "$PEERS" = "1" ]; then
        echo "✅ Federation established! Connected to Society 1"
    else
        echo "⚠️  Waiting for federation connection..."
    fi

    # Show block height
    HEIGHT=$(curl -s http://localhost:26667/status 2>/dev/null | grep -o '"latest_block_height":"[^"]*' | cut -d'"' -f4)
    echo "📊 Current block height: $HEIGHT"
else
    echo "❌ Failed to start blockchain"
    echo "Check society2.log for errors"
    exit 1
fi

echo ""
echo "========================================="
echo "Society 2 is ready for federation!"
echo "========================================="