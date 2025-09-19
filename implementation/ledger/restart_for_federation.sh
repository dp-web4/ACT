#!/bin/bash
# Restart society with network-accessible endpoints for federation

echo "=== Restarting Society for Federation ==="
echo ""

# Get current PID
if [ -f society.pid ]; then
    PID=$(cat society.pid)
    echo "Stopping current society (PID: $PID)..."
    kill $PID 2>/dev/null
    sleep 5
fi

echo "Starting society with network-accessible endpoints..."

# Start with all interfaces accessible
~/go/bin/racecar-webd start --home ./society \
  --rpc.laddr tcp://0.0.0.0:26657 \
  --api.address tcp://0.0.0.0:1317 \
  --grpc.address 0.0.0.0:9090 \
  --api.enable \
  --grpc.enable &

NEW_PID=$!
echo $NEW_PID > society.pid

sleep 5

echo ""
echo "Society restarted with federation-ready configuration:"
echo "  P2P: 0.0.0.0:26656 (already was accessible)"
echo "  RPC: 0.0.0.0:26657 (now accessible from network)"
echo "  API: 0.0.0.0:1317 (now accessible from network)"
echo "  gRPC: 0.0.0.0:9090"
echo "  PID: $NEW_PID"
echo ""
echo "Other machines on the network can now connect to:"
echo "  10.0.0.72:26656 (P2P)"
echo "  http://10.0.0.72:26657 (RPC)"
echo "  http://10.0.0.72:1317 (API)"
echo ""
echo "Society is ready for federation!"