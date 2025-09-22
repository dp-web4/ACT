#!/bin/bash
# Start blockchain for Sprout machine

# Load configuration
ACT_ROOT="/home/sprout/ai-workspace/ACT"
LEDGER_DIR="${ACT_ROOT}/implementation/ledger"
GO_BIN_DIR="/home/sprout/go/bin"
BINARY="${GO_BIN_DIR}/racecar-webd"
SOCIETY_HOME="${LEDGER_DIR}/society-sprout"
PID_FILE="${LEDGER_DIR}/society-sprout.pid"

# Check if already running
if [ -f "${PID_FILE}" ]; then
    OLD_PID=$(cat "${PID_FILE}")
    if kill -0 "${OLD_PID}" 2>/dev/null; then
        echo "Blockchain already running with PID ${OLD_PID}"
        exit 1
    else
        echo "Removing stale PID file"
        rm "${PID_FILE}"
    fi
fi

echo "=== Starting ACT Blockchain on Sprout ==="
echo "Society home: ${SOCIETY_HOME}"

# Start the blockchain in background
cd "${LEDGER_DIR}" || exit 1

echo "Launching blockchain..."
${BINARY} start \
    --home "${SOCIETY_HOME}" \
    --api.enable \
    --api.address "tcp://0.0.0.0:1317" \
    --grpc.enable \
    --grpc.address "0.0.0.0:9090" \
    --rpc.laddr "tcp://0.0.0.0:26657" \
    --p2p.laddr "tcp://0.0.0.0:26656" \
    --minimum-gas-prices "0stake" \
    > "${LEDGER_DIR}/blockchain.log" 2>&1 &

# Save PID
BLOCKCHAIN_PID=$!
echo ${BLOCKCHAIN_PID} > "${PID_FILE}"

echo "Blockchain started with PID: ${BLOCKCHAIN_PID}"
echo "Waiting for blockchain to initialize..."

# Wait for RPC to be available
MAX_TRIES=30
TRIES=0
while [ ${TRIES} -lt ${MAX_TRIES} ]; do
    if curl -s http://localhost:26657/status > /dev/null 2>&1; then
        echo "✅ Blockchain is running!"
        break
    fi
    sleep 1
    TRIES=$((TRIES + 1))
    echo -n "."
done

if [ ${TRIES} -eq ${MAX_TRIES} ]; then
    echo ""
    echo "❌ Blockchain failed to start. Check logs at:"
    echo "${LEDGER_DIR}/blockchain.log"
    exit 1
fi

echo ""
echo "Society is active!"
echo ""
echo "Access points:"
echo "  RPC:  http://10.0.0.36:26657"
echo "  API:  http://10.0.0.36:1317"
echo "  gRPC: 10.0.0.36:9090"
echo "  P2P:  10.0.0.36:26656"
echo ""
echo "Logs: tail -f ${LEDGER_DIR}/blockchain.log"
echo "Stop: kill $(cat ${PID_FILE})"