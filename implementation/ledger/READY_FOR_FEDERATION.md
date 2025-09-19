# Society 1 Ready for Federation!

## Current Status
✅ Society running and accessible from network
✅ All ports open on all interfaces
✅ Genesis file exported
✅ Node ID identified

## Connection Details for Society 2

```bash
# Node ID (for persistent_peers)
c1a129e14fad4cb7c95f9e2b5e9586013941ebf5

# Network endpoints
P2P:  10.0.0.72:26656
RPC:  http://10.0.0.72:26657
API:  http://10.0.0.72:1317
gRPC: 10.0.0.72:9090
```

## Quick Setup Commands for Machine 2

```bash
# 1. Clone and build
git clone https://github.com/dp-web4/ACT.git
cd ACT/implementation/ledger
ignite chain build --skip-proto

# 2. Get genesis from this machine (run from Machine 2)
scp dp@10.0.0.72:/home/dp/ai-workspace/act/implementation/ledger/genesis_for_federation.json ./

# 3. Initialize Society 2
~/go/bin/racecar-webd init "act-society-2" --chain-id act-web4 --home ./society2
cp genesis_for_federation.json ./society2/config/genesis.json

# 4. Configure peer connection
echo 'persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"' >> ./society2/config/config.toml

# 5. Start Society 2
~/go/bin/racecar-webd start --home ./society2 --api.enable --grpc.enable &
```

## Verification Commands

Once Society 2 is running, verify federation:

```bash
# On Machine 1 - Check for peer
curl -s http://10.0.0.72:26657/net_info | jq '.result.n_peers'

# On Machine 2 - Check sync status  
curl -s localhost:26657/status | jq '.result.sync_info.catching_up'
```

## Expected Result

When successful:
- Both societies share same blockchain state
- Transactions on one appear on both
- Network shows 1 peer on each
- Block heights synchronize
- Can create cross-society LCTs

## The Meta-Society Vision

Two independent societies will:
1. 🤝 Discover each other through P2P
2. 🔗 Synchronize their blockchains
3. 🌐 Form trust relationships
4. ⚡ Share energy and work
5. 🧬 Become a society of societies

This is Web4 in action - autonomous entities forming larger collectives through voluntary association and trust!

---

*Machine 1 is ready. Awaiting Society 2...*