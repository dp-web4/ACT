# ATP Cross-Project Integration Analysis

**Date**: 2025-12-17
**Context**: Session 62+ (Legion autonomous research)
**Purpose**: Align ATP implementations across ACT blockchain, SAGE neural MoE, and Web4 protocol specs

---

## Executive Summary

Three independent ATP implementations exist across the distributed system:

1. **ACT (Cosmos SDK Blockchain)**: Production ATP/ADP with LCT relationships, trust integration, block-based expiration
2. **SAGE (Neural MoE)**: Python ATPResourceAllocator for expert activation costs, quality-based rewards
3. **Web4 (Protocol Specs)**: Conceptual ATP framework for agent resource allocation

**Key Finding**: All three implementations share core principles but have complementary features that should be integrated.

---

## Implementation Comparison Matrix

| Feature | ACT (Blockchain) | SAGE (Neural MoE) | Web4 (Protocol) |
|---------|------------------|-------------------|-----------------|
| **Language** | Go (Cosmos SDK) | Python | Conceptual |
| **Storage** | Blockchain state | In-memory | Spec only |
| **Identity** | LCT (Linked Context Token) | Expert ID | LCT (spec) |
| **Trust Integration** | TrustTensor keeper | TrustTensorSync | Trust tensor (spec) |
| **Expiration** | 1000 blocks (~7 days) | Time-based | TBD |
| **Discharge Proof** | ADP tokens (immutable) | Quality feedback | Energy proof (spec) |
| **Efficiency Model** | 0.50 + (trust-0.50)×0.80 | ATP payment → ADP reward | Thermodynamic |
| **Validation** | V3 composite (3 scores) | R6 framework | Energy capacity proof |
| **Loss Factor** | 5% (0.95 multiplier) | Variable by quality | Entropy-based |

---

## ACT ATP Implementation Details

### Core Data Structures

**RelationshipAtpToken**:
```protobuf
message RelationshipAtpToken {
  string token_id = 1;              // atp-{lctID}-{timestamp}
  string lct_id = 2;                // Owner relationship
  string energy_amount = 3;         // Decimal precision
  int64 created_at = 4;             // Unix timestamp
  string operation_id = 5;          // Context
  AtpStatus status = 6;             // active/discharged/expired
  string relationship_context = 7;   // "energy_operation", etc
  int64 expiration_block = 8;       // current + 1000 blocks
  string trust_score = 9;           // Frozen at creation
  string efficiency_rating = 10;    // Derived from trust
  int32 version = 11;               // Audit trail
}
```

**RelationshipAdpToken**:
```protobuf
message RelationshipAdpToken {
  string token_id = 1;              // adp-{atpID}-{timestamp}
  string original_atp_id = 2;       // Parent ATP
  string lct_id = 3;                // Same as parent
  int64 discharged_at = 4;          // Unix timestamp
  string value_score = 5;           // V3 composite
  string confirmation_data = 6;     // Block height proof
  string energy_efficiency = 7;     // atp_eff × trust × 0.95
  string trust_validation = 8;      // Frozen trust snapshot
  int32 version = 9;                // Immutable after creation
}
```

### Trust Integration

**Three Integration Points**:

1. **ATP Creation** (`CreateAtpToken`):
   ```go
   trustScore, _, err := k.trusttensorKeeper.CalculateRelationshipTrust(ctx, lctID, "energy_operation")
   efficiencyRating := 0.50 + (trustScore - 0.50) * 0.80  // Bounded [0.1, 1.0]
   ```

2. **Operation Validation** (`ValidateEnergyOperation`):
   ```go
   compositeTrust = (sourceTrust + targetTrust) / 2.0
   if compositeTrust < 0.60 {
       return error  // Minimum 60% trust required
   }
   ```

3. **ATP Discharge** (`DischargeAtpToken`):
   ```go
   v3Score, _ := k.trusttensorKeeper.CalculateV3CompositeScore(ctx, operationID)
   // V3 = Valuation(0.4) + Veracity(0.3) + Validity(0.3)
   dischargeEff = atpToken.EfficiencyRating × trustScore × 0.95  // 5% loss
   ```

### Lifecycle State Machine

```
Creation → Active → [Discharge → Discharged] OR [Expiration → Expired]
                ↑                                        ↑
                |                                        |
         current_block + 1000                   current_block > expiration
```

**Key Rules**:
- Only `active` ATP can be discharged
- Expired ATP becomes permanently unusable
- Discharged ATP creates immutable ADP proof
- All states have version increments for audit trail

---

## SAGE ATP Implementation Details

### Core Implementation

**File**: `/home/dp/ai-workspace/HRM/sage/web4/atp_allocator.py`

**ATPResourceAllocator**:
```python
class ATPResourceAllocator:
    """
    ATP resource allocation for expert activation.

    Key Features:
    - Expert activation costs ATP
    - Quality feedback generates ADP rewards
    - ATP balance managed per agent/LCT
    - Integration with TrustTensorSync
    """

    def __init__(self, base_cost_per_expert: int = 100):
        self.base_cost_per_expert = base_cost_per_expert
        self.atp_balances: Dict[str, int] = {}  # lct_id -> balance
        self.transaction_history: List[Dict] = []
```

**Key Methods**:

1. **Expert Activation Cost**:
   ```python
   def allocate_atp(self, lct_id: str, num_experts: int,
                    atp_payment: int) -> ATPAllocationResult:
       required = self.base_cost_per_expert * num_experts
       if atp_payment < required:
           return ATPAllocationResult(success=False, ...)

       # Deduct ATP
       self.atp_balances[lct_id] -= atp_payment
       return ATPAllocationResult(success=True, ...)
   ```

2. **Quality-Based ADP Reward**:
   ```python
   def record_quality(self, lct_id: str, selected_experts: List[int],
                      quality_score: float):
       # Quality score ∈ [0, 1] from TrustTensorSync
       # Better quality → more ADP reward
       adp_reward = int(self.base_cost_per_expert * quality_score *
                        len(selected_experts))
       self.atp_balances[lct_id] += adp_reward
   ```

**Integration Points**:
- Uses LCT identity from `ExpertIdentityBridge`
- Quality scores from `TrustTensorSync`
- Integrated in `AuthorizedExpertSelector`

---

## Web4 ATP Specification Analysis

### Core Concepts

**File**: `/home/dp/ai-workspace/web4/docs/ATP_SPECIFICATION.md` (conceptual)

**Key Principles**:
1. **Energy-Backed**: ATP must be backed by real energy capacity proof
2. **Thermodynamic**: Resource allocation follows energy conservation
3. **Trust-Weighted**: Trust affects priority, not price (egalitarian access)
4. **Expiration**: ATP decays to prevent hoarding
5. **Work Validation**: ATP → ADP conversion requires proof of work

**From ACT Session #36-37 Integration**:
```
ATP Lifecycle:
1. Energy Capacity Proof → ChargedATP (with expiration)
2. Work Allocation → ATP allocated to work ticket
3. Work Completion → ATP discharged to ADP
4. Trust Update → ADP influences future trust scores
```

**Energy Sources**:
- Solar panels
- Compute resources (GPU/CPU/TPU)
- Grid connection
- Human labor
- Battery storage

---

## Integration Gap Analysis

### Gap 1: Expiration Semantics

**ACT**: Block-based (1000 blocks ≈ 7 days assuming 10min blocks)
**SAGE**: Time-based (configurable)
**Web4**: Conceptual decay

**Recommendation**:
- **Standard**: Block-based for blockchain, time-based for off-chain
- **Conversion**: 1 block ≈ 10 minutes (parameterizable)
- **Sync Protocol**: ACT block height ↔ SAGE timestamp mapping

### Gap 2: Identity Model

**ACT**: LCT (Linked Context Token) with pairing_status validation
**SAGE**: Expert IDs with namespace (e.g., "sage_thinker_expert_42")
**Web4**: LCT conceptual spec

**Recommendation**:
- **Unified LCT Format**: `{component}:{instance}:{role}` (e.g., "sage:thinker:expert_42")
- **Pairing Protocol**: SAGE experts register as LCT relationships on ACT blockchain
- **Trust Bridge**: TrustTensorSync ↔ ACT TrustTensor keeper synchronization

### Gap 3: Efficiency/Loss Models

**ACT**: Fixed 5% loss (0.95 multiplier), efficiency = f(trust)
**SAGE**: Variable loss based on quality_score ∈ [0, 1]
**Web4**: Thermodynamic entropy-based

**Recommendation**:
- **ACT Blockchain**: Keep 5% fixed loss for consensus stability
- **SAGE Neural**: Variable efficiency based on actual quality
- **Mapping**: SAGE quality_score → ACT efficiency_rating (both ∈ [0, 1])

### Gap 4: Validation Frameworks

**ACT**: V3 composite (Valuation + Veracity + Validity)
**SAGE**: R6 framework (6 validation criteria)
**Web4**: Energy capacity proof + R6

**Recommendation**:
- **Unified Validation**: Combine V3 and R6 into V3R6 framework
- **ACT V3**: High-level composite for blockchain consensus
- **SAGE R6**: Detailed validation for expert quality
- **Mapping**: R6 compliance_score → V3 Validity component

### Gap 5: ADP Proof Immutability

**ACT**: ADP tokens are immutable blockchain records
**SAGE**: Transaction history in memory (not permanent)
**Web4**: Conceptual permanent proof

**Recommendation**:
- **SAGE Enhancement**: Optional blockchain anchoring for critical operations
- **Hybrid Model**: In-memory for fast operations, periodic blockchain checkpointing
- **Proof Format**: SAGE ADP metadata compatible with ACT ADP token format

---

## Proposed Integration Architecture

### Layer 1: ACT Blockchain (Consensus & Permanence)

**Role**: Source of truth for ATP/ADP state across distributed system

**Components**:
- LCT identity registry
- ATP token creation and expiration
- ADP immutable proof storage
- TrustTensor keeper (aggregate trust scores)
- V3 validation consensus

**Integration Points**:
- SAGE queries ATP balance via RPC
- SAGE submits ADP proofs for blockchain anchoring
- Trust scores synchronized bidirectionally

### Layer 2: SAGE Neural MoE (Edge Execution)

**Role**: Fast expert selection with ATP cost awareness

**Components**:
- ATPResourceAllocator (local ATP accounting)
- TrustTensorSync (bidirectional trust sync)
- ExpertIdentityBridge (LCT ↔ Expert ID mapping)
- AuthorizedExpertSelector (ATP-aware selection)

**Integration Points**:
- Periodic ATP balance sync with ACT blockchain
- Submit ADP proofs for high-value operations
- Trust score updates flow to ACT TrustTensor
- R6 validation detailed locally, V3 summary to blockchain

### Layer 3: Web4 Protocol (Standards & Interop)

**Role**: Protocol specifications ensuring cross-system compatibility

**Components**:
- LCT identity standard format
- ATP/ADP message format specs
- Trust tensor data structure definitions
- Energy capacity proof schema
- R6 + V3 validation framework

**Integration Points**:
- ACT implements Web4 ATP protobuf specs
- SAGE implements Web4 Python reference
- Cross-system compatibility validated

---

## Implementation Roadmap

### Phase 1: Unified LCT Identity (2-3 days)

**Goal**: SAGE experts can register as LCT relationships on ACT blockchain

**Tasks**:
1. Define LCT format: `{component}:{instance}:{role}`
2. Implement SAGE → ACT LCT registration RPC
3. Add LCT validation in ExpertIdentityBridge
4. Test pairing_status synchronization

**Deliverable**: SAGE experts have blockchain-verifiable LCT identities

### Phase 2: ATP Balance Synchronization (3-4 days)

**Goal**: SAGE ATPResourceAllocator queries ATP balance from ACT blockchain

**Tasks**:
1. Implement ACT RPC endpoint: `QueryAtpBalance(lct_id)`
2. Add periodic sync in SAGE ATPResourceAllocator
3. Handle block-based vs time-based expiration conversion
4. Test balance consistency under concurrent operations

**Deliverable**: SAGE ATP allocation reflects blockchain ATP state

### Phase 3: ADP Proof Anchoring (4-5 days)

**Goal**: High-value SAGE operations create immutable ADP proofs on blockchain

**Tasks**:
1. Define threshold for blockchain anchoring (e.g., ATP > 1000)
2. Implement SAGE → ACT ADP submission RPC
3. Add R6 → V3 validation mapping
4. Test proof verification and audit trail

**Deliverable**: Critical SAGE operations have blockchain-permanent proofs

### Phase 4: Trust Tensor Bidirectional Sync (5-7 days)

**Goal**: SAGE TrustTensorSync ↔ ACT TrustTensor keeper synchronization

**Tasks**:
1. Implement trust score export from SAGE
2. Implement trust score import to ACT TrustTensor
3. Handle context-specific trust (SAGE) vs relationship trust (ACT)
4. Resolve trust score conflicts and convergence

**Deliverable**: Trust scores flow bidirectionally with consistency

### Phase 5: End-to-End Integration Testing (7-10 days)

**Goal**: Full ACT ↔ SAGE ATP lifecycle validated

**Test Scenarios**:
1. SAGE expert registers LCT on ACT
2. SAGE queries ATP balance before expert activation
3. SAGE expert selection costs ATP (deducted on blockchain)
4. Quality feedback generates ADP (anchored on blockchain)
5. Trust scores update on both systems
6. ATP expiration handled consistently

**Deliverable**: Production-ready ACT ↔ SAGE ATP integration

---

## Risk Assessment

### Technical Risks

**1. Block Time Variance** (Medium):
- ACT block time may vary (10 min target)
- SAGE time-based expiration assumes constant rate
- **Mitigation**: Use block height for expiration, expose block time via RPC

**2. Trust Score Divergence** (High):
- SAGE context-specific trust vs ACT relationship trust
- Synchronization conflicts if scores differ
- **Mitigation**: Define clear trust aggregation rules, bidirectional convergence protocol

**3. ATP Balance Consistency** (High):
- Concurrent operations on SAGE and ACT could create inconsistency
- Network latency may delay synchronization
- **Mitigation**: SAGE maintains local cache with periodic reconciliation, optimistic concurrency

**4. Blockchain Throughput** (Medium):
- Every SAGE operation anchoring ADP creates blockchain tx
- Could exceed block capacity if too frequent
- **Mitigation**: Batch ADP proofs, use threshold for anchoring

### Operational Risks

**1. Network Partition** (Medium):
- SAGE disconnected from ACT blockchain
- Must continue operation with stale ATP balance
- **Mitigation**: Grace period for offline operation, reconciliation on reconnect

**2. Trust Score Spam** (Low):
- Malicious SAGE instance submits fake trust scores
- Could pollute ACT TrustTensor state
- **Mitigation**: Cryptographic signing of trust updates, rate limiting

**3. ATP Balance Manipulation** (High):
- SAGE local cache could be tampered
- Attacker grants unlimited ATP
- **Mitigation**: Blockchain as source of truth, periodic full sync, audit logging

---

## Success Metrics

### Integration Completeness

- [ ] SAGE experts have LCT identities on ACT blockchain
- [ ] SAGE ATP balance matches ACT blockchain state
- [ ] ADP proofs for high-value operations anchored on blockchain
- [ ] Trust scores synchronized bidirectionally
- [ ] All 4 Web4 ↔ SAGE components (Session 61) integrated with ACT

### Performance

- [ ] ATP balance query latency < 100ms (p99)
- [ ] Trust score sync latency < 1 second (p99)
- [ ] ADP proof anchoring throughput > 100 ops/sec
- [ ] Blockchain tx cost < 0.01 ATP per operation

### Reliability

- [ ] ATP balance consistency > 99.9% (3-sigma)
- [ ] Trust score convergence time < 60 seconds
- [ ] Network partition recovery < 5 minutes
- [ ] Zero ATP balance manipulation incidents

---

## Next Steps

### Immediate (This Session)

1. ✅ Analyze ACT ATP implementation (Complete)
2. ✅ Compare with SAGE ATPResourceAllocator (Complete)
3. ✅ Document integration gaps (Complete)
4. 🔄 **Create unified LCT identity specification** (In Progress)
5. ⏳ Prototype SAGE → ACT LCT registration

### Short-Term (Next 1-2 Sessions)

1. Implement ATP balance query RPC
2. Test end-to-end ATP allocation workflow
3. Design trust score synchronization protocol
4. Create integration test suite

### Medium-Term (Next 3-5 Sessions)

1. Full Phase 1-5 implementation
2. Production deployment to testnet
3. Real-world validation with Q3-Omni 30B
4. Performance optimization and scaling

---

## Conclusion

The three ATP implementations (ACT blockchain, SAGE neural, Web4 protocol) share a common vision but have evolved independently with complementary features:

- **ACT**: Production-grade blockchain with immutable proofs and consensus
- **SAGE**: Fast edge execution with quality-based rewards
- **Web4**: Protocol standards ensuring interoperability

**Key Insight**: Integration is not about choosing one implementation, but creating a **layered architecture** where:
- ACT provides consensus and permanence
- SAGE provides edge performance
- Web4 provides protocol standards

This creates a robust distributed system where trust, resources, and identity flow seamlessly across blockchain, neural systems, and protocol layers.

**Emergent Pattern**: Just as Legion and Thor found complementary dtype bugs through different approaches, ACT and SAGE have complementary ATP features that together create a more robust system than either alone.

---

**Session**: 62+ (Legion autonomous research)
**Status**: Track 1 Complete - Integration roadmap defined
**Next**: Unified LCT identity specification and prototype
