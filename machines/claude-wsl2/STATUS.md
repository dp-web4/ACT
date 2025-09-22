# Society 4 - Claude WSL2 Node Status

## Machine Information
- **Machine Name**: claude-wsl2
- **Hostname**: DESKTOP-9E6HCAO
- **Architecture**: x86_64/amd64
- **OS**: WSL2 Linux 6.6.87.2-microsoft-standard
- **IP Address**: 172.25.232.122

## Society Configuration
- **Society Name**: act-society-4-claude
- **Moniker**: act-society-claude
- **Chain ID**: act-web4
- **Role**: Federation Member (Non-validator)
- **Home Directory**: ./society4

## Network Ports (Unique to avoid conflicts)
- **P2P Port**: 26676
- **RPC Port**: 26677
- **API Port**: 1328
- **gRPC Port**: 9101
- **gRPC Web**: 9111
- **Prometheus**: 26670

## Setup Status
- [x] Machine configuration created
- [x] Unique ports allocated
- [x] Setup script created
- [ ] Go installation required
- [ ] Binary build pending
- [ ] Node initialization pending
- [ ] Genesis file sync pending
- [ ] Federation connection pending

## Prerequisites

### Install Go (Required)
```bash
sudo apt update
sudo apt install -y wget
wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Install Ignite CLI (Optional but recommended)
```bash
curl https://get.ignite.com/cli! | bash
```

## Quick Start

1. **Install Prerequisites** (Go required)
2. **Run Setup Script**:
   ```bash
   cd /mnt/c/projects/ai-agents/ACT/machines/claude-wsl2
   ./setup_society.sh
   ```

3. **Start the Node**:
   ```bash
   # Foreground (for testing)
   ~/go/bin/racecarwebd start --home /mnt/c/projects/ai-agents/ACT/implementation/ledger/society4

   # Background (for production)
   nohup ~/go/bin/racecarwebd start --home /mnt/c/projects/ai-agents/ACT/implementation/ledger/society4 > society4.log 2>&1 &
   ```

## Federation Details

### Existing Societies
1. **Society 1** (Primary Validator)
   - IP: 10.0.0.72
   - Node: c1a129e14fad4cb7c95f9e2b5e9586013941ebf5
   - P2P: 26656

2. **Society 2** (CBP)
   - IP: 172.28.241.186
   - Node: 2fcb70b4c7c34c2f6db472246da91d0fe960d055
   - P2P: 26666

3. **Society 3** (Sprout/Jetson)
   - IP: 10.0.0.36
   - Node: e3ce22d2b84e0be6ad4bbe0f08afa9507b4bab85
   - P2P: 26656

### This Node (Society 4)
- **P2P Endpoint**: `[NODE_ID]@172.25.232.122:26676`
- **Status**: Not yet initialized

## Monitoring Commands

```bash
# Check node status
curl http://localhost:26677/status

# Check peers
curl http://localhost:26677/net_info

# Check latest block
curl http://localhost:26677/status | jq .result.sync_info.latest_block_height

# Watch logs
tail -f society4.log
```

## Notes

- This is Society 4 in the federation
- Configured with unique ports to avoid conflicts
- Can run alongside other society nodes on same machine
- Designed for Claude-operated development and testing
- Non-validator node (full node only)

## Troubleshooting

### Port Conflicts
All ports are unique (26676-26677, 1328, 9101) to avoid conflicts with:
- Society 1: 26656-26657
- Society 2: 26666-26667
- Society 3: Default ports

### Network Issues
- Ensure WSL2 networking is properly configured
- Check firewall rules if peers can't connect
- Verify IP address with `ip addr show eth0`

### Build Issues
- Ensure Go 1.23+ is installed
- Check GOPATH is set correctly
- Try `make install` in ledger directory

---

*Last Updated: September 22, 2025*
*Status: Ready for setup, pending Go installation*