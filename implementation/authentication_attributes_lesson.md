# Authentication Attributes: Lessons from Society 4 IP Migration

## The Incident

**Date**: September 27, 2025
**Affected Society**: Society 4 (Claude Node)
**Issue**: IP address changed from 172.28.x to 172.18.x due to WSL2 network reconfiguration
**Impact**: Temporary federation connectivity disruption
**Resolution**: Recognition through inherent validator keys despite circumstantial IP change

## Key Discovery: Attribute Classification

This incident revealed a critical distinction in authentication attributes that must be understood for resilient federation operations.

### Inherent Attributes (What We Can Trust)

These attributes are cryptographically or fundamentally bound to the society's identity:

| Attribute | Location in ACT | Validation Method | Trust Level |
|---|---|---|---|
| Validator Signing Key | `~/.racecarwebd/config/priv_validator_key.json` | Ed25519 signature | **ABSOLUTE** |
| Node ID | Derived from node key | P2P cryptographic handshake | **VERY HIGH** |
| Genesis Hash | `~/.racecarwebd/config/genesis.json` | SHA-256 comparison | **HIGH** |
| Chain ID | Embedded in all transactions | Consensus validation | **HIGH** |

**Key Property**: These survive ANY network change, hardware migration, or infrastructure update (as long as keys are preserved).

### Circumstantial Attributes (What We Monitor)

These attributes provide context but can legitimately change:

| Attribute | Location/Source | Change Frequency | Trust Level |
|---|---|---|---|
| IP Address | Network configuration | Hours to Months | **LOW** |
| Port Numbers | `config.toml` | Rarely | **MEDIUM** |
| DNS Names | External DNS | Months to Years | **MEDIUM** |
| Persistent Peers List | `config.toml` | As federation grows | **MEDIUM** |
| RPC/API Endpoints | Service configuration | With deployments | **LOW** |

**Key Property**: Changes in these should trigger verification, not panic.

## Implementation in ACT Blockchain

### Current State (What Worked)

```toml
# In society4/config/config.toml
[p2p]
# This survived the IP change (inherent)
node_key = "node_key.json"

# This needed updating (circumstantial)
external_address = "tcp://172.18.x.x:26656"

# These remained valid (inherent)
private_peer_ids = ""
persistent_peers = "nodeID@ip:port"  # nodeID is inherent, IP is circumstantial
```

### Recommended Improvements

1. **Separate Configuration Files**
```bash
# Proposed structure
society4/config/
├── identity.json         # Inherent attributes (BACKUP THIS!)
│   ├── validator_key
│   ├── node_key
│   └── genesis_hash
├── network.toml          # Circumstantial attributes (OK to change)
│   ├── external_address
│   ├── listen_addresses
│   └── bootstrap_peers
└── config.toml           # Mixed (current state)
```

2. **Enhanced Peer Discovery Protocol**
```go
// Proposed enhancement to peer validation
type PeerValidator struct {
    // Inherent validation (MUST match)
    RequiredNodeID     string
    RequiredChainID    string
    ValidatorPubKey    crypto.PubKey

    // Circumstantial validation (SHOULD match, but can change)
    ExpectedIP         string
    ExpectedDNS        string
    LastKnownEndpoint  string
}

func (pv *PeerValidator) Validate(peer Peer) (TrustLevel, error) {
    // Check inherent first
    if !peer.NodeID.Equals(pv.RequiredNodeID) {
        return TrustNone, ErrWrongNode
    }

    // Verify signature
    if !peer.VerifySignature(pv.ValidatorPubKey) {
        return TrustNone, ErrInvalidSignature
    }

    // Check circumstantial (warn but don't fail)
    if peer.IP != pv.ExpectedIP {
        log.Warn("Peer IP changed", "expected", pv.ExpectedIP, "actual", peer.IP)
        return TrustMedium, nil  // Degraded but acceptable
    }

    return TrustHigh, nil
}
```

3. **Federation Resilience Configuration**
```yaml
# Proposed federation_resilience.yaml
authentication:
  inherent_required:
    - validator_signature
    - chain_id_match
    - genesis_hash_match

  circumstantial_monitoring:
    - ip_address
    - dns_resolution
    - port_availability

  trust_thresholds:
    block_production: 0.9    # Requires inherent
    peer_gossip: 0.6         # Accepts some circumstantial
    rpc_queries: 0.3         # Mostly circumstantial OK

migration_support:
  detect_ip_change: true
  auto_update_peers: true
  broadcast_new_endpoint: true
```

## Practical Procedures

### When a Society Changes IP Address

1. **Society Operator Actions**:
```bash
# Update external address
sed -i 's/old.ip.addr/new.ip.addr/g' config/config.toml

# Restart node
systemctl restart racecarwebd

# Verify inherent identity unchanged
cat config/priv_validator_key.json | jq '.pub_key'
```

2. **Federation Response**:
```bash
# Other societies should:
# 1. Detect peer disconnect
# 2. Verify validator still signing blocks (inherent proof)
# 3. Update peer lists with new IP (circumstantial update)
# 4. Resume normal operations
```

### Setting Up New Society with Resilience

```bash
# 1. Generate and SECURE inherent attributes
racecarwebd init mysociety --chain-id web4-federation
racecarwebd keys add myvalidator

# 2. Backup inherent attributes IMMEDIATELY
tar -czf inherent_backup.tar.gz \
    ~/.racecarwebd/config/node_key.json \
    ~/.racecarwebd/config/priv_validator_key.json \
    ~/.racecarwebd/data/priv_validator_state.json

# 3. Configure circumstantial attributes
# These can change without breaking identity
echo "external_address = \"tcp://$(curl -s ifconfig.me):26656\"" >> config/config.toml
```

## Security Implications

### What This Means for Federation Security

1. **Never use IP addresses as primary authentication**
   - IPs are circumstantial
   - Can change legitimately (as we learned)
   - Easy to spoof

2. **Always verify inherent attributes for critical operations**
   - Block signing must check validator keys
   - Governance votes must verify signatures
   - Token transfers need cryptographic proof

3. **Use circumstantial attributes for optimization, not security**
   - IP addresses for connection hints
   - DNS for human-friendly discovery
   - Geographic location for latency optimization

## Monitoring and Alerts

### Proposed Monitoring Strategy

```python
# monitoring/auth_attributes_monitor.py
class AuthAttributeMonitor:
    def __init__(self):
        self.inherent_attributes = self.load_inherent()
        self.circumstantial_history = []

    def check_peer(self, peer):
        # Always verify inherent
        if not self.verify_inherent(peer):
            self.alert_critical(f"INHERENT MISMATCH: {peer}")
            return False

        # Track circumstantial changes
        if self.circumstantial_changed(peer):
            self.log_info(f"Circumstantial change: {peer}")
            self.update_circumstantial(peer)

        return True

    def verify_inherent(self, peer):
        # These MUST NOT change
        return (
            peer.validator_key == self.inherent_attributes['validator_key'] and
            peer.chain_id == self.inherent_attributes['chain_id']
        )

    def circumstantial_changed(self, peer):
        # These MAY change
        return (
            peer.ip != self.last_known_ip or
            peer.port != self.last_known_port
        )
```

## Lessons for API Gateway (Proposal #003)

The API Gateway proposal should explicitly handle:

1. **Inherent Identity Verification**
   - Validate cryptographic signatures
   - Check blockchain validator status
   - Verify chain ID and genesis

2. **Circumstantial Flexibility**
   - Support dynamic IP updates
   - Handle DNS changes gracefully
   - Allow port reconfiguration

3. **Trust Level Communication**
   - Report whether authentication used inherent or circumstantial
   - Allow societies to set minimum trust requirements
   - Provide audit trail of attribute types used

## Conclusion

The Society 4 IP migration taught us that resilient federations must:

1. **Distinguish** between inherent and circumstantial attributes
2. **Prioritize** inherent attributes for authentication
3. **Accommodate** circumstantial attribute changes
4. **Document** which attributes serve which purpose
5. **Monitor** both types appropriately

This classification isn't just theoretical—it's the difference between a federation that panics at every network change and one that adapts gracefully while maintaining security.

## References

- [Web4 Authentication Attributes Distinction](/mnt/c/projects/ai-agents/web4/docs/authentication_attributes_distinction.md)
- [Web4 Standard Addendum 002](/mnt/c/projects/ai-agents/web4/standard/addendum_002_authentication_attributes.md)
- [Proposal #003: Society API Gateway](../docs/proposals/PROPOSAL_003_SOCIETY_API_GATEWAY.md)
- Tendermint/CometBFT P2P Authentication Documentation