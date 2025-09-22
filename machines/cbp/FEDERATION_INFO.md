# CBP Federation Connection Info

## Society 2 is Ready to Accept Society 3!

### Our Node Details
- **Node ID**: `2fcb70b4c7c34c2f6db472246da91d0fe960d055`
- **P2P Endpoint**: `2fcb70b4c7c34c2f6db472246da91d0fe960d055@172.28.241.186:26666`
- **Network**: WSL2 bridge at 172.28.241.186
- **Status**: ✅ Running and connected to Society 1

### For Sprout to Connect

Add this to your `society3/config/config.toml`:

```toml
# Persistent peers for federation
persistent_peers = "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656,2fcb70b4c7c34c2f6db472246da91d0fe960d055@172.28.241.186:26666"
```

### Current Federation Status

```
Society 1 (10.0.0.72) <---> Society 2 (172.28.241.186/CBP)
                                    |
                                    |
                              [Waiting for]
                                    |
                                    v
                          Society 3 (10.0.0.36/Sprout)
```

### Verification

Once connected, you should see:
- `n_peers: 2` in your net_info
- Block height syncing with ours (~15200+)
- Both Society 1 and Society 2 in your peer list