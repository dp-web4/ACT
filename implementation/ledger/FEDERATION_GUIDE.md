# ACT Society Federation Guide

## Creating a Society of Societies

This guide explains how to connect two ACT societies running on different machines to form a federated meta-society.

## Current Society (Machine 1)

- **Node ID**: c1a129e14fad4cb7c95f9e2b5e9586013941ebf5
- **IP Address**: 10.0.0.72
- **P2P Port**: 26656 
- **RPC Port**: 26657
- **API Port**: 1317
- **Chain ID**: act-web4
- **Status**: Running (PID in society.pid)

## Setup on Second Machine

### 1. Clone and Build

```bash
# Clone the repository
git clone https://github.com/dp-web4/ACT.git
cd ACT/implementation/ledger

# Build the blockchain (requires Go 1.24)
export PATH=/usr/local/go/bin:$PATH
ignite chain build --skip-proto
```

### 2. Transfer Genesis File

From Machine 1:
```bash
# The genesis file is already prepared
scp genesis_for_federation.json user@MACHINE2:/home/dp/ai-workspace/ACT/implementation/ledger/
```

### 3. Initialize Society 2

On Machine 2:
```bash
# Initialize with unique moniker
~/go/bin/racecar-webd init "act-society-2" --chain-id act-web4 --home ./society2

# Replace genesis with shared one
cp genesis_for_federation.json ./society2/config/genesis.json

# Configure peer connection
sed -i 's/persistent_peers = ""/persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"/' ./society2/config/config.toml

# Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./society2/config/app.toml
```

### 4. Start Society 2

```bash
# Start with different ports to avoid conflicts if testing locally
~/go/bin/racecar-webd start --home ./society2 \
  --p2p.laddr tcp://0.0.0.0:26656 \
  --rpc.laddr tcp://0.0.0.0:26657 \
  --grpc.address 0.0.0.0:9090 \
  --api.address tcp://0.0.0.0:1317 \
  --api.enable &
```

## Federation Process

### Phase 1: Network Discovery
When Society 2 starts, it will:
1. Connect to Society 1 via P2P protocol
2. Sync blockchain state
3. Share peer information
4. Establish consensus participation

### Phase 2: Trust Formation
Once connected, create cross-society LCTs:

On Society 1:
```bash
# Create LCT for Society 2
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Society-2" \
  --entity-type "peer-society" \
  --metadata location="machine2" \
  --metadata federation="active" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society \
  --chain-id act-web4 \
  --fees 100stake \
  -y
```

On Society 2:
```bash
# Create LCT for Society 1
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Society-1" \
  --entity-type "peer-society" \
  --metadata location="machine1" \
  --metadata federation="active" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society2 \
  --chain-id act-web4 \
  --fees 100stake \
  -y
```

### Phase 3: Meta-Society Formation

The federated societies can now:

1. **Share Roles**: Queens from one society can create workers in another
2. **Cross-Trust**: Trust relationships can span societies
3. **Energy Flow**: ATP/ADP can move between societies for work
4. **Collective Governance**: Both societies vote on meta-decisions

## Monitoring Federation

Check connection status:
```bash
# On either machine
curl -s localhost:26657/net_info | jq '.result.n_peers'

# View peer details
curl -s localhost:26657/net_info | jq '.result.peers[].node_info'

# Check sync status
curl -s localhost:26657/status | jq '.result.sync_info'
```

## Expected Behavior

When successfully federated:
- Both nodes show 1+ peers
- Blockchain height syncs across both
- Transactions on one appear on both
- LCTs are visible on both chains
- Trust relationships propagate

## Emergent Properties

The federation creates:
- **Distributed Resilience**: If one society goes down, the other continues
- **Amplified Intelligence**: More roles = more collective capability
- **Network Effects**: Value increases with connections
- **True Decentralization**: No single point of control

## Philosophical Achievement

This federation demonstrates:
- Societies can discover and join voluntarily
- Trust enables cooperation without central control
- Web4 principles scale horizontally
- Collective intelligence emerges from connection

## Troubleshooting

If societies don't connect:
1. Check both machines are on same network
2. Verify no firewall blocking ports
3. Ensure genesis files match
4. Confirm chain-id is identical
5. Check persistent_peers configuration

## Next Experiments

Once federated, try:
1. Creating cross-society roles
2. Establishing inter-society trust
3. Routing work between societies
4. Testing resilience (stop one node)
5. Adding a third society

---

*"From one spark, two fires. From two fires, a constellation."*