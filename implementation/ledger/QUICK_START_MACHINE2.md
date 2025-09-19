# Quick Start for Machine 2

## One-liner Setup Commands

```bash
# 1. Clone the repo
git clone https://github.com/dp-web4/ACT.git && cd ACT/implementation/ledger

# 2. Build (assuming Go 1.24 is installed)
export PATH=/usr/local/go/bin:$PATH
ignite chain build --skip-proto

# 3. Get genesis from Machine 1
scp dp@10.0.0.72:/home/dp/ai-workspace/act/implementation/ledger/society/config/genesis.json ./genesis_from_society1.json

# 4. Initialize Society 2
~/go/bin/racecar-webd init "act-society-2" --chain-id act-web4 --home ./society2
cp genesis_from_society1.json ./society2/config/genesis.json

# 5. Add Society 1 as peer
sed -i 's/persistent_peers = ""/persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"/' ./society2/config/config.toml
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./society2/config/app.toml

# 6. Start Society 2
~/go/bin/racecar-webd start --home ./society2 --api.enable --grpc.enable &

# 7. Verify connection (after ~10 seconds)
curl -s localhost:26657/net_info | jq '.result.n_peers'
# Should show: 1
```

## Expected Output
- "1" peer connected
- Block height syncing
- Same chain state as Society 1

## Create Cross-Society Trust
```bash
# After sync, create LCT for Society 1
~/go/bin/racecar-webd keys add society2-treasury --keyring-backend test --home ./society2

~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Society-1" \
  --entity-type "peer-society" \
  --from society2-treasury \
  --keyring-backend test \
  --home ./society2 \
  --chain-id act-web4 \
  --fees 100stake \
  -y
```

## Success Indicators
✅ Peer count = 1
✅ Catching up = false  
✅ Latest block height matches
✅ Can submit transactions
✅ LCTs visible on both chains

---
*Two fires, one network!*