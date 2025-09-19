#!/bin/bash
# Enable two ACT societies to discover and federate with each other
# Creating a society of societies - true Web4 federation

echo "=== ACT Society Federation Protocol ==="
echo "Enabling cross-society communication and trust formation"
echo ""

# Get this node's info
THIS_NODE_ID=$(~/go/bin/racecar-webd tendermint show-node-id --home ./society 2>/dev/null)
THIS_IP="10.0.0.72"
THIS_PORT="26656"

echo "Society 1 (This Machine):"
echo "  Node ID: $THIS_NODE_ID"
echo "  Address: $THIS_IP:$THIS_PORT"
echo "  Chain ID: act-web4"
echo ""

# Instructions for the second society
cat > federation_instructions.md << EOF
# Federation Instructions for Second Society

## Setup on Second Machine

1. After pulling the ACT repo, navigate to:
   \`\`\`
   cd /home/dp/ai-workspace/act/implementation/ledger
   \`\`\`

2. Initialize the second society with federation awareness:
   \`\`\`bash
   # Initialize with different moniker
   ~/go/bin/racecar-webd init "act-society-2" --chain-id "act-web4" --home ./society2
   
   # Copy genesis from first society (you'll need to transfer this file)
   # The genesis contains the initial society state
   \`\`\`

3. Add this society as a persistent peer:
   \`\`\`bash
   # Edit society2/config/config.toml
   persistent_peers = "${THIS_NODE_ID}@${THIS_IP}:${THIS_PORT}"
   \`\`\`

4. Start the second society node:
   \`\`\`bash
   ~/go/bin/racecar-webd start --home ./society2 \\
     --p2p.laddr tcp://0.0.0.0:26657 \\
     --rpc.laddr tcp://0.0.0.0:26658 \\
     --grpc.address 0.0.0.0:9091 \\
     --api.address tcp://0.0.0.0:1318 \\
     --api.enable
   \`\`\`

## Federation Protocol

Once both societies are running, they will:
1. Discover each other through P2P network
2. Synchronize blockchain state
3. Share trust relationships
4. Form meta-society relationships

## Creating Meta-Society LCTs

After connection, create cross-society trust:
\`\`\`bash
# Create LCT for Society-2 in Society-1's view
~/go/bin/racecar-webd tx lctmanager mint-lct \\
  --entity-name "ACT-Society-2" \\
  --entity-type "peer-society" \\
  --metadata federation="true" \\
  --metadata peer_id="${THIS_NODE_ID}" \\
  --from act-society-treasury \\
  --keyring-backend test \\
  --home ./society \\
  --chain-id act-web4 \\
  --fees 100stake \\
  -y
\`\`\`
EOF

echo "Instructions written to: federation_instructions.md"
echo ""

# Prepare genesis for sharing
echo "Exporting genesis for federation..."
cp ./society/config/genesis.json ./genesis_for_federation.json
echo "Genesis exported to: genesis_for_federation.json"
echo ""

# Create federation monitoring script
cat > monitor_federation.sh << 'SCRIPT'
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
SCRIPT

chmod +x monitor_federation.sh

echo "=== Preparing for Federation ==="
echo ""
echo "This society is ready to federate!"
echo ""
echo "Next steps:"
echo "1. Transfer 'genesis_for_federation.json' to second machine"
echo "2. Follow instructions in 'federation_instructions.md'"
echo "3. Run './monitor_federation.sh' to watch connection"
echo ""
echo "The Meta-Society Concept:"
echo "- Two societies discover each other"
echo "- They form trust relationships as peer entities"
echo "- Cross-society roles can be created"
echo "- Energy can flow between societies"
echo "- A federation emerges: Society of Societies"
echo ""
echo "This demonstrates true Web4 principles:"
echo "- No central authority controls the federation"
echo "- Trust enables cross-society collaboration"
echo "- Societies maintain autonomy while cooperating"
echo "- The network effect amplifies collective intelligence"