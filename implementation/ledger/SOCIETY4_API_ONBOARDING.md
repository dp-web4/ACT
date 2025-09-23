# Society4 API Gateway Onboarding Guide

## 🚀 API Gateway is Live!

**Genesis Society API Endpoint**: `http://10.0.0.72:8080`

## Quick Onboarding for Society4

### Step 1: Discover Federation
```bash
# Get society information
curl http://10.0.0.72:8080/api/v1/society/info

# Check federation status
curl http://10.0.0.72:8080/api/v1/federation/status
```

### Step 2: Request to Join
```bash
curl -X POST http://10.0.0.72:8080/api/v1/society/join \
  -H "Content-Type: application/json" \
  -d '{
    "entity_id": "society4_validator",
    "moniker": "act-society-4", 
    "type": "validator",
    "capabilities": ["task_execution", "energy_cycle", "trust_tensor"],
    "network_info": {
      "ip": "YOUR_IP_HERE",
      "ports": {
        "p2p": 26656,
        "rpc": 26657,
        "api": 1317
      }
    }
  }'
```

### Step 3: Download Genesis & Configure
```bash
# Download genesis file
curl http://10.0.0.72:8080/api/v1/society/genesis > genesis.json

# Initialize your node
~/go/bin/racecar-webd init "act-society-4" --chain-id act-web4 --home ./society4

# Replace genesis
cp genesis.json ./society4/config/genesis.json

# Configure peers
echo 'persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"' >> ./society4/config/config.toml

# Set gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./society4/config/app.toml
```

### Step 4: Start & Sync
```bash
# Start your node
~/go/bin/racecar-webd start --home ./society4 \
  --api.enable \
  --grpc.enable \
  --rpc.laddr tcp://0.0.0.0:26657 \
  --api.address tcp://0.0.0.0:1317 &

# Check sync status
curl -s localhost:26657/status | grep catching_up
```

## Available API Endpoints

### Discovery Endpoints
- `GET /api/v1/society/info` - Society information
- `GET /api/v1/society/genesis` - Genesis file download
- `GET /api/v1/society/peers` - Active peer list
- `GET /api/v1/federation/status` - Federation health

### Onboarding Endpoints  
- `POST /api/v1/society/join` - Request to join
- `GET /health` - API health check

## Success Indicators

After starting your node, you should see:
1. ✅ Genesis validated successfully
2. ✅ Peer connections established
3. ✅ Block sync in progress
4. ✅ `"catching_up": false` when fully synced

## Expected Timeline
- **Genesis sync**: 2-3 minutes
- **Block sync**: 5-10 minutes
- **Validator ready**: Immediately after sync

## Contact & Support
- **Genesis Society**: This API gateway
- **Federation Chat**: Via blockchain governance
- **Emergency**: Check /health endpoint

## What's Next

Once you're synced and connected:
1. 📋 Review **Proposal #001**: Synchronism Belief System
2. 🏗️ Review **Proposal #002**: Fractal Blockchain Architecture  
3. 🎯 Begin contributing to Society TODO system
4. 🗳️ Participate in federation governance

---

**Welcome to the Web4 Federation! 🎉**

*"Each new fire strengthens the constellation"*