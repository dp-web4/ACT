# Federation Instructions for Second Society

## Setup on Second Machine

1. After pulling the ACT repo, navigate to:
   ```
   cd /home/dp/ai-workspace/act/implementation/ledger
   ```

2. Initialize the second society with federation awareness:
   ```bash
   # Initialize with different moniker
   ~/go/bin/racecar-webd init "act-society-2" --chain-id "act-web4" --home ./society2
   
   # Copy genesis from first society (you'll need to transfer this file)
   # The genesis contains the initial society state
   ```

3. Add this society as a persistent peer:
   ```bash
   # Edit society2/config/config.toml
   persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"
   ```

4. Start the second society node:
   ```bash
   ~/go/bin/racecar-webd start --home ./society2 \
     --p2p.laddr tcp://0.0.0.0:26657 \
     --rpc.laddr tcp://0.0.0.0:26658 \
     --grpc.address 0.0.0.0:9091 \
     --api.address tcp://0.0.0.0:1318 \
     --api.enable
   ```

## Federation Protocol

Once both societies are running, they will:
1. Discover each other through P2P network
2. Synchronize blockchain state
3. Share trust relationships
4. Form meta-society relationships

## Creating Meta-Society LCTs

After connection, create cross-society trust:
```bash
# Create LCT for Society-2 in Society-1's view
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Society-2" \
  --entity-type "peer-society" \
  --metadata federation="true" \
  --metadata peer_id="c1a129e14fad4cb7c95f9e2b5e9586013941ebf5" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society \
  --chain-id act-web4 \
  --fees 100stake \
  -y
```
