# Instructions for Society4 to Join Federation

## Current Federation Members
1. Society-1 (10.0.0.72) - This machine, Genesis society  
2. act-society-claude (10.0.0.147) - Currently connected

## For Society4 to Connect

### Quick Setup Commands

```bash
# 1. Clone ACT repository
git clone https://github.com/dp-web4/ACT.git
cd ACT/implementation/ledger

# 2. Build blockchain
export PATH=/usr/local/go/bin:$PATH
ignite chain build --skip-proto

# 3. Get genesis from Society-1
scp dp@10.0.0.72:/home/dp/ai-workspace/act/implementation/ledger/society/config/genesis.json ./genesis_society4.json

# 4. Initialize Society4
~/go/bin/racecar-webd init "act-society4" --chain-id act-web4 --home ./society4

# 5. Replace genesis
cp genesis_society4.json ./society4/config/genesis.json

# 6. Configure persistent peers (add all known nodes)
cat >> ./society4/config/config.toml << EOF
persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656,7edf3af375ad9f95cbf652cf3ab66af086e956bc@10.0.0.147:26656"
EOF

# 7. Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./society4/config/app.toml

# 8. Start Society4
~/go/bin/racecar-webd start --home ./society4 \
  --api.enable \
  --grpc.enable \
  --rpc.laddr tcp://0.0.0.0:26657 \
  --api.address tcp://0.0.0.0:1317 &
```

## Verification

After starting, verify connection:
```bash
# Check peer count
curl -s localhost:26657/net_info | grep n_peers
# Should show: "n_peers": "2" or more

# Check sync status
curl -s localhost:26657/status | grep catching_up
# Should show: "catching_up": false (after sync)
```

## Connection Details

### Society-1 (Genesis)
- Node ID: `c1a129e14fad4cb7c95f9e2b5e9586013941ebf5`
- P2P: `10.0.0.72:26656`
- RPC: `http://10.0.0.72:26657`

### act-society-claude
- Node ID: `7edf3af375ad9f95cbf652cf3ab66af086e956bc`  
- P2P: `10.0.0.147:26656`

## Expected Outcome

Once connected, Society4 will:
1. Sync with existing blockchain state
2. Participate in consensus
3. Strengthen the federation
4. Add resilience to the network

## Welcome Society4! 🎉

When you join, the federation will have 3+ active validators, making it even more decentralized and resilient.

---

*"Each new fire strengthens the constellation"*