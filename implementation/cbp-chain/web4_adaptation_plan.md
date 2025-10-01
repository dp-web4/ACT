# CBP Chain Web4 Adaptation Plan
Based on Society4's Reference Implementation

## Overview
Adapting Society4's Web4-compliant ATP/ADP energy economics and governance to CBP Chain.

## Key Adaptations from Society4

### 1. ATP/ADP Token Pool System
**Society4's Implementation:**
- Fixed 1000 ATP budget per society
- Role-based allocation (Security Queen: 150 ATP)
- Daily recharge: 20 ATP per role
- Discharge mechanism: ATP → ADP for operations

**CBP Adaptation:**
```python
# CBP Role Structure with ATP Allocations
CBP_ROLES = [
    {"name": "Data Queen", "initial_atp": 140, "daily_recharge": 20},      # Highest for data ops
    {"name": "Metrics Queen", "initial_atp": 130, "daily_recharge": 20},   # Performance tracking
    {"name": "Security Queen", "initial_atp": 130, "daily_recharge": 20},  # Security oversight
    {"name": "Bridge Queen", "initial_atp": 120, "daily_recharge": 20},    # Federation bridging
    {"name": "Cache Queen", "initial_atp": 110, "daily_recharge": 20},     # KV cache management
    {"name": "Worker 1", "initial_atp": 90, "daily_recharge": 20},         # Implementation
    {"name": "Worker 2", "initial_atp": 90, "daily_recharge": 20},         # Testing
    {"name": "Worker 3", "initial_atp": 90, "daily_recharge": 20},         # Documentation
    {"name": "Coordinator", "initial_atp": 100, "daily_recharge": 20},     # Orchestration
]
# Total: 1000 ATP (Web4 compliant)
```

### 2. LCT Identity Structure
**Society4's Implementation:**
- Hardware binding with EAT (Entity Attestation Token)
- MRH parent relationships
- Policy capabilities and constraints

**CBP Adaptation:**
```json
{
  "lct_id": "lct:web4:mb32:cbp_self_0001",
  "subject": "did:web4:cbp:coordinator",
  "binding": {
    "entity_type": "device",
    "hardware_anchor": "eat:mb64:hw:[CBP_HARDWARE_HASH]",
    "public_key": "mb64:[CBP_PUBLIC_KEY]",
    "binding_proof": "cose:[BINDING_PROOF]"
  },
  "mrh": {
    "bound": [{
      "lct_id": "lct:web4:hardware:wsl2:[CBP_HW_HASH]",
      "type": "parent",
      "relationship": "hardware_platform"
    }],
    "horizon_distance": 2
  },
  "policy": {
    "capabilities": [
      "cache_management",
      "metrics_collection",
      "federation_bridging",
      "blockchain_witnessing"
    ]
  }
}
```

### 3. Governance Integration
**Society4's Law Oracle:**
- Machine-readable laws (JSON-LD)
- RFC extensions for procedural clarity
- Automated compliance checking

**CBP Adaptation:**
- Focus on data governance and metrics policy
- Cache invalidation rules
- Federation bridging protocols
- Performance thresholds

### 4. Implementation Timeline

#### Phase 1: Foundation (Week 1)
- [ ] Create CBP role structure with ATP allocations
- [ ] Implement token pool management
- [ ] Add discharge/recharge mechanisms

#### Phase 2: Identity (Week 2)
- [ ] Generate CBP hardware binding
- [ ] Create self LCT with proper MRH
- [ ] Implement policy framework

#### Phase 3: Integration (Week 3)
- [ ] Connect to federation energy grid
- [ ] Implement daily recharge automation
- [ ] Add blockchain witnessing

#### Phase 4: SAGE Support (Week 4)
- [ ] KV cache with ATP cost tracking
- [ ] Metrics collection with energy budgets
- [ ] Performance optimization for 15W operation

## Key Differences from Society4

1. **Role Focus**: CBP emphasizes data/metrics vs Society4's legal/governance
2. **Energy Priority**: Cache operations get highest ATP allocation
3. **Bridge Function**: Unique federation bridging role
4. **SAGE Integration**: Direct support for KV cache persistence

## Next Steps

1. Begin implementing ATP/ADP pool system
2. Create CBP-specific role definitions
3. Adapt Society4's keeper functions
4. Test energy economics with scheduler

## Benefits of Adaptation

- **Web4 Compliance**: Achieve similar 7.9/10 score
- **Energy Economics**: Proper ATP/ADP lifecycle
- **Federation Integration**: Compatible with other societies
- **SAGE Ready**: Energy-aware cache management

## Code Reuse

Society4's implementation provides excellent templates:
- `society_pool_keeper.go` → `cbp_pool_keeper.py`
- `society_token_pool.go` → `cbp_token_pool.py`
- `recharge_automation.go` → `cbp_recharge.py`

We'll translate from Go to Python for CBP's implementation.