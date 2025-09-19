#!/bin/bash
# Monitor federation status

echo "=== Federation Status Monitor ==="
echo ""

# Check peer connections
echo "Connected Peers:"
curl -s localhost:26657/net_info | jq '.result.peers[] | {id: .node_info.id, address: .remote_ip}'

echo ""
echo "Sync Status:"
curl -s localhost:26657/status | jq '.result.sync_info | {latest_block_height, catching_up}'

echo ""
echo "Society Relationships:"
~/go/bin/racecar-webd query lctmanager list-lcts --output json | jq '.lcts[] | select(.entity_type == "peer-society")'
