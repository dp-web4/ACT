# Federation Resilience Architecture

## Current Federation Status
- **Society-1** (10.0.0.72) - Genesis society, this machine
- **act-society-claude** (10.0.0.147) - Connected
- **Sprout** (10.0.0.36) - Joined previously
- **Society4** - Ready to join

## Key Resilience Features for Society TODO System

### 1. Machine-Contextual Configuration

Each society maintains its own context file:
```yaml
# society_context.yaml
society:
  id: "society_001"
  machine_id: "dp-desktop-72"
  network_interface: "eth0"
  local_ip: "10.0.0.72"
  public_endpoints:
    - "tcp://10.0.0.72:26656"  # P2P
    - "http://10.0.0.72:26657"  # RPC
    - "http://10.0.0.72:1317"   # REST
  
federation:
  known_peers:
    - id: "act-society-claude"
      endpoints: ["tcp://10.0.0.147:26656"]
      last_seen: "2025-01-22T20:15:00Z"
      trust_level: 0.85
    - id: "sprout"
      endpoints: ["tcp://10.0.0.36:26656"]
      last_seen: "2025-01-21T15:30:00Z"
      trust_level: 0.75
  
  reconnect_policy:
    max_retries: 10
    backoff_seconds: [1, 2, 4, 8, 16, 32, 64, 128, 256, 512]
    give_up_after_seconds: 3600
```

### 2. Communication Dropout Recovery

#### Automatic Peer Recovery
```go
func (k Keeper) MonitorPeerHealth(ctx sdk.Context) {
    for _, peer := range k.GetKnownPeers(ctx) {
        if !k.IsPeerHealthy(peer) {
            k.InitiateReconnect(ctx, peer)
        }
    }
}
```

#### TODO State Synchronization
When a society reconnects after dropout:
1. Exchange TODO state hashes
2. Identify divergence points
3. Reconcile using consensus voting
4. Update local state

### 3. Portable Task Distribution

#### Task Metadata Structure
```go
type PortableTodo struct {
    // Universal identifiers
    GlobalID      string    // UUID v4
    OriginSociety string    // Creating society
    
    // Portable content
    Title         string
    Description   string
    Priority      TodoPriority
    
    // Machine-agnostic requirements
    RequiredCapabilities []string
    MinTrustLevel       sdk.Dec
    
    // Federation tracking
    AcceptedBy    []string  // Society IDs
    RejectedBy    []string
    InProgressAt  string    // Current executor society
    
    // Resilience features
    Checkpoints   []Checkpoint
    LastHeartbeat time.Time
    FailoverPlan  FailoverStrategy
}
```

### 4. Consensus During Partitions

#### Degraded Mode Operations
When fewer than 50% of federation members are available:
```
STATE: DEGRADED_CONSENSUS
- Accept only CRITICAL priority todos
- Increase ATP costs by 2x
- Queue non-critical for later consensus
- Log all decisions for reconciliation
```

#### Split-Brain Prevention
```go
func (k Keeper) CanProcessTodoInDegraded(ctx sdk.Context, todo TodoItem) bool {
    activeNodes := k.CountActivePeers(ctx)
    totalNodes := k.CountTotalPeers(ctx)
    quorum := float64(activeNodes) / float64(totalNodes)
    
    if quorum < 0.5 {
        // Only critical with witness
        return todo.Priority == CRITICAL && 
               len(todo.Witnesses) >= 2
    }
    return true
}
```

### 5. Cross-Society TODO Migration

When a society goes offline while executing a TODO:
```
1. Heartbeat timeout (30 seconds)
2. Grace period (2 minutes)  
3. Automatic reassignment:
   - Find next eligible society
   - Transfer task state
   - Checkpoint progress
   - Update ATP allocations
```

### 6. Federation Event Log

Distributed event log for audit trail:
```json
{
  "event_type": "todo_migration",
  "timestamp": "2025-01-22T20:30:45Z",
  "from_society": "society_002",
  "to_society": "society_001",
  "todo_id": "todo_12345",
  "reason": "heartbeat_timeout",
  "checkpoint": "step_3_of_7",
  "witnesses": ["society_003", "society_004"]
}
```

## Implementation Strategy

### Phase 1: Local Resilience (Current)
- [x] Basic TODO system
- [x] Wake/sleep cycles
- [ ] Local state persistence
- [ ] Checkpoint system

### Phase 2: Federation Awareness
- [ ] Peer health monitoring
- [ ] State synchronization protocol
- [ ] Degraded mode operations
- [ ] TODO migration on dropout

### Phase 3: Full Distribution
- [ ] Distributed consensus for TODOs
- [ ] Cross-society ATP pooling
- [ ] Global reputation system
- [ ] Automatic load balancing

## Testing Scenarios

### Scenario 1: Rolling Dropouts
1. Start all 4 societies
2. Create shared TODO
3. Drop Society-2 during execution
4. Verify automatic migration to Society-3
5. Bring Society-2 back online
6. Verify state reconciliation

### Scenario 2: Network Partition
1. Split federation: (S1,S2) | (S3,S4)
2. Submit todos to both partitions
3. Verify degraded mode activation
4. Reunite network
5. Verify consensus reconciliation

### Scenario 3: Machine Migration
1. Save Society-1 state
2. Move to new machine (different IP)
3. Update peer configurations
4. Restore and reconnect
5. Verify TODO continuity

## Configuration Files

### Per-Society Configuration
```bash
# Each society maintains:
/society/config/
  ├── node_key.json        # Identity
  ├── priv_validator_key.json
  ├── society_context.yaml # Machine context
  └── federation_peers.json # Known peers
```

### Portable TODO State
```bash
/society/data/
  ├── todos/
  │   ├── active/         # Current todos
  │   ├── checkpoints/    # Progress snapshots
  │   └── completed/      # History
  └── federation/
      ├── peer_health.db  # Peer status
      └── event_log.db    # Audit trail
```

## Key Principles

1. **No Single Point of Failure**: Any society can fail without stopping the federation
2. **Portable State**: TODOs can migrate between societies seamlessly
3. **Contextual Awareness**: Each society knows its machine capabilities
4. **Graceful Degradation**: Reduced functionality better than no functionality
5. **Eventual Consistency**: Reconciliation when federation reunites
6. **Trust-Based Recovery**: Higher trust societies have priority in conflicts

## Monitoring Commands

```bash
# Check federation health
racecarwebd query societytodo federation-status

# View peer connections
racecarwebd query societytodo peer-health

# Check TODO migrations
racecarwebd query societytodo migration-log

# View degraded mode status
racecarwebd query societytodo consensus-state
```

## Emergency Procedures

### Force TODO Migration
```bash
racecarwebd tx societytodo force-migrate \
  --todo-id=todo_123 \
  --to-society=society_backup \
  --reason="manual_intervention"
```

### Reset Peer Connection
```bash
racecarwebd tx societytodo reset-peer \
  --peer-id=society_002 \
  --new-endpoint="tcp://10.0.0.146:26656"
```

### Enter Maintenance Mode
```bash
racecarwebd tx societytodo maintenance-mode \
  --duration=3600 \
  --accept-critical-only=true
```

---

*"The federation is not a rigid structure but a living network that adapts, recovers, and evolves."*