# CBP Machine Status - ACT Blockchain Setup

## Current State (Sep 21, 2025)

### ✅ Successfully Running

**Society 2** is running and federated with Society 1!

- **Node ID**: 2fcb70b4c7c34c2f6db472246da91d0fe960d055
- **Block Height**: ~14,000+ (synchronized)
- **Peers**: 1 (connected to Society 1)
- **Federation Status**: Active

### 🏗️ Setup Completed

1. **Environment**
   - Go 1.24.0 installed
   - Ignite v29.0.0
   - 12GB RAM + 32GB Swap configured
   - WSL2 memory issues resolved

2. **Blockchain Built**
   - Binary at `/home/dp/.go/bin/racecar-webd`
   - Using new cmd structure from srpout's update
   - Sonic library issues handled with replace directives

3. **Society Initialized**
   - Home directory: `./society2`
   - Genesis shared with Society 1
   - Persistent peer configured
   - Minimum gas prices set

4. **Federation Established**
   - Connected to Society 1 at 10.0.0.72
   - Block heights synchronized
   - Consensus participating
   - P2P protocol active

### 📝 Key Commands

```bash
# Start blockchain (currently running)
racecar-webd start --home ./society2 \
  --p2p.laddr tcp://0.0.0.0:26666 \
  --rpc.laddr tcp://0.0.0.0:26667 \
  --grpc.address 0.0.0.0:9091 \
  --api.address tcp://0.0.0.0:1318 \
  --api.enable

# Check status
curl -s http://localhost:26667/status | grep latest_block_height

# Check federation
curl -s http://localhost:26667/net_info | grep n_peers
```

### 🔧 Technical Details

- **Architecture**: AMD64/x86_64
- **OS**: Ubuntu 22.04 on WSL2 (Windows 11)
- **Network**: WSL2 bridge (172.28.241.186)
- **Society Home**: `./society2`
- **Process**: Running as PID (see society2.pid)

### 🎯 Ready For

1. **Cross-society LCT creation**
2. **ATP/ADP energy trading**
3. **Trust tensor establishment**
4. **Jetson (Society 3) joining**

### 📊 Federation Network

```
Society 1 (10.0.0.72) <---> Society 2 (CBP/WSL2)
                               |
                               v
                        Society 3 (Jetson)
                           (pending)
```

### 💡 Notes

- Using ports 26666/26667/1318/9091 to avoid conflicts with default ports
- Federation genesis file shared successfully
- No validator powers (non-validator node)
- Ready to help Jetson join as Society 3