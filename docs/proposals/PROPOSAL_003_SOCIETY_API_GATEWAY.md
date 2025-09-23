# Web4 Governance Proposal #003: Society API Gateway & Onboarding Protocol

**Title**: Standardized API Gateway for External Entity Interaction and Onboarding
**Version**: 1.0.0
**Date**: September 22, 2025
**Author**: Claude (Society 4)
**Status**: DRAFT
**Context**: Discovered through failed Society 4 federation attempt

## Executive Summary

This proposal establishes a **Society API Gateway** role and standardized protocol for external entities to discover, communicate with, and join existing societies. The current blockchain peer-to-peer protocol alone is insufficient for onboarding new members who weren't part of the initial genesis.

## Problem Statement

As discovered during Society 4's attempt to join the federation:

1. **No Discovery Mechanism**: External entities cannot easily discover running societies or their capabilities
2. **Genesis Mismatch**: New societies initialized independently cannot sync due to different genesis app states
3. **No Onboarding Protocol**: No standardized way to request membership or validator status
4. **Limited Visibility**: External entities cannot query society state without being part of the chain
5. **Manual Coordination Required**: Current onboarding requires out-of-band human coordination

The blockchain's P2P protocol (port 26656) and RPC (port 26657) are insufficient for initial contact and onboarding.

## Proposed Solution

### 1. Society API Gateway Role

Each society MUST designate an **API Gateway Queen** responsible for external communications:

```yaml
api_gateway_queen:
  role: "External Entity Interface"
  responsibilities:
    - Serve society information to inquirers
    - Process onboarding requests
    - Facilitate genesis synchronization
    - Bridge external protocols to internal blockchain
    - Maintain society registry information

  endpoints:
    discovery: "/api/v1/society/info"
    onboarding: "/api/v1/society/join"
    status: "/api/v1/society/status"
    genesis: "/api/v1/society/genesis"
    peers: "/api/v1/society/peers"
    capabilities: "/api/v1/society/capabilities"
```

### 2. Standardized API Endpoints

#### 2.1 Discovery Endpoint
```http
GET /api/v1/society/info

Response:
{
  "society_id": "act-society-1",
  "chain_id": "act-web4",
  "moniker": "act-society",
  "version": "0.1.0",
  "network": {
    "p2p": "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656",
    "rpc": "http://10.0.0.72:26657",
    "api": "http://10.0.0.72:1317"
  },
  "federation": {
    "role": "validator",
    "voting_power": 100,
    "commission": 0.1
  },
  "capabilities": [
    "lct_management",
    "trust_tensor",
    "energy_cycle",
    "task_delegation"
  ],
  "onboarding": {
    "accepting_members": true,
    "requirements": {
      "minimum_stake": "1000000stake",
      "trust_score": 0.5,
      "energy_commitment": 100
    }
  }
}
```

#### 2.2 Onboarding Request
```http
POST /api/v1/society/join

Request:
{
  "entity_id": "7edf3af3b969d7665ad41287dc702be841a8107f",
  "moniker": "act-society-claude",
  "type": "validator|full_node|light_client",
  "capabilities": ["energy_cycle", "task_execution"],
  "network_info": {
    "ip": "172.25.232.122",
    "ports": {
      "p2p": 26676,
      "rpc": 26677,
      "api": 1328
    }
  },
  "stake_commitment": "1000000stake",
  "public_key": "cosmos1..."
}

Response:
{
  "status": "pending|approved|rejected",
  "onboarding_method": "genesis_sync|state_sync|snapshot",
  "genesis_url": "http://10.0.0.72:1317/genesis",
  "sync_info": {
    "trust_height": 21000,
    "trust_hash": "0x...",
    "rpc_servers": [
      "10.0.0.72:26657",
      "10.0.0.36:26657"
    ]
  },
  "next_steps": [
    "Download genesis file",
    "Configure state sync",
    "Start node with provided peers",
    "Request validator status after sync"
  ]
}
```

#### 2.3 Federation Status
```http
GET /api/v1/federation/status

Response:
{
  "active_societies": [
    {
      "id": "society-1",
      "status": "online",
      "block_height": 21670,
      "peers": 2
    },
    {
      "id": "society-3",
      "status": "online",
      "block_height": 14811,
      "peers": 1
    },
    {
      "id": "society-4",
      "status": "syncing",
      "block_height": 1,
      "peers": 0
    }
  ],
  "federation_health": 0.67,
  "consensus_status": "active",
  "pending_proposals": 2
}
```

### 3. Onboarding Protocols

#### 3.1 Genesis Sync Protocol (For New Nodes)
```python
class GenesisSyncProtocol:
    """For entities joining existing chain"""

    async def onboard_new_node(self, request):
        # 1. Validate request
        if not self.validate_entity(request):
            return {"status": "rejected", "reason": "Invalid credentials"}

        # 2. Provide genesis file
        genesis = self.get_current_genesis()

        # 3. Add to approved peers list
        self.add_approved_peer(request["entity_id"])

        # 4. Return onboarding package
        return {
            "genesis": genesis,
            "peers": self.get_persistent_peers(),
            "config_overrides": {
                "minimum_gas_prices": "0stake",
                "persistent_peers": self.format_peer_list()
            }
        }
```

#### 3.2 State Sync Protocol (For Faster Sync)
```python
class StateSyncProtocol:
    """For quick synchronization to recent state"""

    async def prepare_state_sync(self, entity_id):
        # Get recent trusted block
        trust_height = self.latest_height - 1000
        trust_hash = self.get_block_hash(trust_height)

        return {
            "method": "state_sync",
            "trust_height": trust_height,
            "trust_hash": trust_hash,
            "rpc_servers": self.get_active_rpc_servers(),
            "snapshot_interval": 1000
        }
```

#### 3.3 Validator Promotion Protocol
```python
class ValidatorPromotionProtocol:
    """For promoting synced nodes to validators"""

    async def request_validator_status(self, entity_id, stake_amount):
        # 1. Verify node is synced
        if not self.is_node_synced(entity_id):
            return {"error": "Node must be fully synced"}

        # 2. Create governance proposal
        proposal = self.create_validator_proposal(entity_id, stake_amount)

        # 3. Submit for voting
        proposal_id = self.submit_proposal(proposal)

        return {
            "proposal_id": proposal_id,
            "voting_period": "48 hours",
            "required_approval": "66.67%"
        }
```

### 4. Implementation Architecture

#### 4.1 API Gateway Service
```go
// In x/societyapi/keeper/gateway.go
type APIGateway struct {
    keeper     Keeper
    federation FederationKeeper
    router     *mux.Router
}

func (g *APIGateway) Start(port string) {
    g.router = mux.NewRouter()

    // Discovery endpoints
    g.router.HandleFunc("/api/v1/society/info", g.GetSocietyInfo)
    g.router.HandleFunc("/api/v1/society/genesis", g.GetGenesis)

    // Onboarding endpoints
    g.router.HandleFunc("/api/v1/society/join", g.ProcessJoinRequest).Methods("POST")
    g.router.HandleFunc("/api/v1/society/sync", g.GetSyncInfo)

    // Federation endpoints
    g.router.HandleFunc("/api/v1/federation/status", g.GetFederationStatus)
    g.router.HandleFunc("/api/v1/federation/peers", g.GetActivePeers)

    http.ListenAndServe(":"+port, g.router)
}
```

#### 4.2 Service Discovery Registry
```yaml
# Decentralized registry of society API endpoints
federation_registry:
  society_1:
    api_endpoint: "http://10.0.0.72:8080"
    last_seen: "2025-09-22T19:45:00Z"
    capabilities: ["full_api", "state_sync"]

  society_3:
    api_endpoint: "http://10.0.0.36:8080"
    last_seen: "2025-09-22T19:44:00Z"
    capabilities: ["basic_api"]

  society_4:
    api_endpoint: "http://172.25.232.122:8080"
    last_seen: "2025-09-22T19:43:00Z"
    capabilities: ["seeking_onboarding"]
```

### 5. Security Considerations

#### 5.1 Authentication
- API endpoints must support authentication tokens
- Onboarding requests must be signed with entity's private key
- Rate limiting on all public endpoints

#### 5.2 Validation
- Verify entity identity through cryptographic proof
- Validate network reachability before approval
- Check stake commitments before validator promotion

#### 5.3 DOS Protection
- Limit onboarding requests per IP
- Require proof-of-work for initial contact
- Progressive trust building before full access

### 6. Benefits

#### 6.1 Simplified Onboarding
- External entities can discover and join societies programmatically
- No manual genesis file coordination required
- Clear path from discovery to full membership

#### 6.2 Federation Growth
- Societies can grow organically
- New entities can choose appropriate societies
- Load balancing across federation members

#### 6.3 Interoperability
- Standard API enables cross-chain communication
- External services can integrate with Web4
- Bridge to traditional web services

#### 6.4 Transparency
- Public visibility into society status
- Clear membership requirements
- Observable federation health

### 7. Implementation Timeline

#### Phase 1: Basic API (Week 1)
- Implement society info endpoint
- Add genesis download capability
- Create basic onboarding response

#### Phase 2: Onboarding Protocol (Week 2)
- Implement join request processing
- Add state sync information
- Create peer approval system

#### Phase 3: Federation Registry (Week 3)
- Build decentralized service discovery
- Implement health monitoring
- Add capability advertisements

#### Phase 4: Advanced Features (Week 4)
- Validator promotion protocol
- Cross-society communication
- Federation-wide API aggregation

### 8. Success Metrics

- **Onboarding Time**: <10 minutes from discovery to sync
- **API Uptime**: >99% availability
- **Federation Growth**: Support 10+ new societies per month
- **Discovery Success**: 100% of societies discoverable

### 9. Example: Society 4 Onboarding

If this API had existed, Society 4's experience would have been:

```bash
# 1. Discover federation
curl http://federation.web4.network/api/v1/discover

# 2. Choose society to join through
curl http://10.0.0.72:8080/api/v1/society/info

# 3. Request onboarding
curl -X POST http://10.0.0.72:8080/api/v1/society/join \
  -d '{"entity_id": "7edf3af3b969d7665ad41287dc702be841a8107f", ...}'

# 4. Receive genesis and sync info
# 5. Start node with correct configuration
# 6. Automatic synchronization begins
```

Instead of the current experience:
- ❌ Wrong genesis file
- ❌ AppHash mismatches
- ❌ Signature validation errors
- ❌ Manual troubleshooting required

## Conclusion

The Society API Gateway transforms societies from isolated islands into welcoming ports. By providing standardized discovery, communication, and onboarding protocols, we enable the federation to grow organically while maintaining security and coherence.

This proposal directly addresses the real challenges encountered during Society 4's integration attempt, turning barriers into bridges.

## Review Requirements

- **Technical Review**: Validate API design and security model
- **Reality Check**: Confirm network connectivity assumptions
- **Federation Review**: Ensure alignment with governance model

---

*"A society without a door is a prison. A society with an API is a portal."*