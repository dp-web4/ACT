# Instructions for Sprout to Join the Federation

## Current Federation Status
- Society-1 (10.0.0.72) - Genesis society
- Society-2 (10.0.0.146) - First federation member
- Sprout (awaiting) - Ready to join!

## For Sprout Machine

### Quick Join Commands

```bash
# 1. Get the genesis file from Society-1
scp dp@10.0.0.72:/home/dp/ai-workspace/act/implementation/ledger/society/config/genesis.json ./genesis_federation.json

# 2. Initialize Sprout society
~/go/bin/racecar-webd init "sprout" --chain-id act-web4 --home ./sprout

# 3. Replace genesis
cp genesis_federation.json ./sprout/config/genesis.json

# 4. Add both existing societies as persistent peers
echo 'persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656,2fcb70b4c7c34c2f6db472246da91d0fe960d055@10.0.0.146:26656"' >> ./sprout/config/config.toml

# 5. Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./sprout/config/app.toml

# 6. Start Sprout
~/go/bin/racecar-webd start --home ./sprout --api.enable --grpc.enable &
```

## Alternative: If already initialized

If Sprout is already running but not connecting:

```bash
# Add our peers to config
sed -i 's/persistent_peers = ""/persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656,2fcb70b4c7c34c2f6db472246da91d0fe960d055@10.0.0.146:26656"/' ./sprout/config/config.toml

# Restart
pkill racecar-webd
~/go/bin/racecar-webd start --home ./sprout --api.enable --grpc.enable &
```

## Connection Details

### Society-1 (This machine)
- Node ID: `c1a129e14fad4cb7c95f9e2b5e9586013941ebf5`
- P2P: `10.0.0.72:26656`
- RPC: `http://10.0.0.72:26657`

### Society-2 
- Node ID: `2fcb70b4c7c34c2f6db472246da91d0fe960d055`
- P2P: `10.0.0.146:26656`

## Verification

Once connected, verify with:
```bash
curl -s localhost:26657/net_info | grep n_peers
# Should show: "n_peers": "2"
```

## Welcome Sprout! 🌱

When you join, the federation becomes:
- 3 societies
- Shared consensus
- Distributed validation
- Even more resilient!

The network effect amplifies - each new member strengthens the whole.