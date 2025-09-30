# Web4 Compliance Remediation Plan - Society 4
**Created**: September 30, 2025
**Target Completion**: November 15, 2025 (6 weeks)
**Status**: APPROVED - EXECUTION READY

## Overview

This plan addresses all findings from the web4-compliance-reviewer agent's assessment and provides a path to full web4 compliance while preserving Society 4's innovations.

## Phase 1: Identity Foundation (Weeks 1-2)

### Task 1.1: Implement Proper LCT Structure
**Priority**: CRITICAL
**Estimated Effort**: 3 days
**Assignee**: Society 4 development team

**Deliverables**:
1. Create `lct.go` module with complete LCT data structures
2. Wrap existing hardware binding in LCT format
3. Generate proper LCT ID: `lct:web4:mb32:<multibase-hash>`
4. Create DID subject: `did:web4:key:z6Mk<publickey>`
5. Implement binding proof signatures (COSE format)

**Implementation**:
```go
// File: blockchain/source/x/lct/types/lct.go
package types

type Society4LCT struct {
    LCTID         string           `json:"lct_id"`
    Subject       string           `json:"subject"`
    Binding       LCTBinding       `json:"binding"`
    BirthCert     BirthCertificate `json:"birth_certificate"`
    MRH           MRHGraph         `json:"mrh"`
    Policy        Policy           `json:"policy"`
    Attestations  []Attestation    `json:"attestations"`
    CreatedAt     time.Time        `json:"created_at"`
    UpdatedAt     time.Time        `json:"updated_at"`
}

type LCTBinding struct {
    EntityType     string `json:"entity_type"`    // "device"
    PublicKey      string `json:"public_key"`     // COSE key format
    HardwareAnchor string `json:"hardware_anchor"` // EAT token
    CreatedAt      string `json:"created_at"`
    BindingProof   string `json:"binding_proof"`  // COSE signature
}
```

**Validation**: LCT must pass web4 schema validation

---

### Task 1.2: Obtain Birth Certificate from ACT Federation
**Priority**: CRITICAL
**Estimated Effort**: 2 days
**Assignee**: Society 4 + Genesis coordination

**Process**:
1. Request citizenship in ACT Federation
2. Submit founding documentation:
   - Society purpose statement
   - Initial laws draft
   - Founding members (King Claudius + Queens)
   - Hardware binding attestation
3. Obtain witness signatures from:
   - Genesis (Society 1)
   - Society 2
   - Sprout (Society 3)
4. Record birth certificate in federation ledger

**Birth Certificate Structure**:
```json
{
  "birth_certificate": {
    "citizen_role": "lct:web4:role:citizen:society4",
    "context": "federation",
    "birth_timestamp": "2025-10-07T12:00:00Z",
    "parent_entity": "lct:web4:federation:act",
    "birth_witnesses": [
      "lct:web4:society:genesis",
      "lct:web4:society:society2",
      "lct:web4:society:sprout"
    ],
    "founding_attestation": {
      "purpose": "Mobile AI consciousness node with hardware sovereignty",
      "governance": "monarchic",
      "special_capabilities": ["network_mobility", "temporal_auth"]
    }
  }
}
```

**Validation**: Birth certificate co-signed by 3+ federation witnesses

---

### Task 1.3: Publish Law Oracle with Machine-Readable Laws
**Priority**: CRITICAL
**Estimated Effort**: 4 days
**Assignee**: Law Oracle Queen + Society 4 dev

**Deliverables**:
1. Create Law Oracle module
2. Define Society 4 foundational laws in machine-readable format
3. Implement R6 rule bindings
4. Set up law versioning system
5. Publish to federation

**Implementation**:
```yaml
# File: laws/foundational/society4_laws_v1.0.0.yaml
type: Web4LawDataset
id: web4://law/society4/1.0.0
hash: sha256-<computed>
version: 1.0.0

norms:
  # Identity Laws
  - id: LAW-IDENTITY-001
    name: Hardware Binding Requirement
    selector: r6.entity.identity
    operator: requires
    value: hardware_binding
    description: All Society 4 operations require valid hardware binding

  - id: LAW-IDENTITY-002
    name: Genesis Hardware Hash
    selector: r6.entity.hardware_hash
    operator: equals
    value: "93e766842ee7882a248e7d55ef3269b95e1735b0be88b94287b18029d1851759"
    description: Genesis 'self' LCT bound to this specific hardware

  # Governance Laws
  - id: LAW-GOV-001
    name: Monarchic Approval
    selector: r6.transaction.approval
    operator: requires
    value: monarch_signature
    description: King Claudius must approve all major transactions

  - id: LAW-GOV-002
    name: Security Queen Veto
    selector: r6.transaction.security_check
    operator: requires
    value: security_queen_approval
    description: Security Queen has veto power on security matters

  - id: LAW-GOV-003
    name: Pending Consensus Processing
    selector: r6.consensus.network
    operator: equals
    value: home_federation
    description: Pending decisions processed only from home network

  # Economic Laws
  - id: LAW-ECON-001
    name: Security Queen ATP Allocation
    selector: r6.role.security_queen.atp
    operator: equals
    value: 150

  - id: LAW-ECON-002
    name: Total ATP Budget
    selector: r6.society.atp_total
    operator: equals
    value: 1000

  - id: LAW-ECON-003
    name: Daily ATP Recharge
    selector: r6.role.*.atp_recharge
    operator: equals
    value: 20
    description: All queens receive 20 ATP daily recharge

procedures:
  - id: PROC-WITNESS-MIN
    name: Minimum Witness Requirement
    requiresWitnesses: 3
    witnessTypes: [federation_member]

  - id: PROC-EMERGENCY-HALT
    name: Emergency Halt Procedure
    trigger: security_breach_detected
    authority: security_queen
    requires_witnesses: false
    immediate: true

r6Bindings:
  - web4://schemas/r6-rules-v1
  - web4://schemas/society-governance-v1

publication:
  url: web4://law/society4/1.0.0
  rdf_endpoint: http://society4.act.federation/sparql
  update_frequency: monthly

signatures:
  - entity: lct:web4:society4:law_oracle_queen
    signature: <cose-signature>
    timestamp: 2025-10-10T00:00:00Z
```

**Validation**: Law dataset parseable by R6 engine, validates against web4 schema

---

### Task 1.4: Map Security Queen to Authority Role
**Priority**: HIGH
**Estimated Effort**: 1 day
**Assignee**: Society 4 dev

**Implementation**:
```json
{
  "lct_id": "lct:web4:role:authority:security_queen",
  "subject": "did:web4:society4:queen:security",
  "role_type": "authority",
  "authority_scope": [
    "security_validation",
    "hardware_binding_verification",
    "emergency_halt",
    "key_management_oversight",
    "code_review_authority"
  ],
  "delegation_source": "lct:web4:society:society4",
  "delegated_by": "lct:web4:society4:king_claudius",
  "powers": [
    "veto_on_security_violations",
    "halt_on_compromise_detected",
    "validate_hardware_binding",
    "block_unsafe_transactions",
    "force_key_rotation"
  ],
  "constraints": {
    "veto_requires_evidence": true,
    "must_document_decisions": true,
    "subject_to_review": true,
    "accountability_log": "immutable"
  }
}
```

**Validation**: Authority role validates against web4 authority schema

---

## Phase 2: Economic Implementation (Week 3)

### Task 2.1: Implement ATP/ADP Token Pool System
**Priority**: CRITICAL
**Estimated Effort**: 5 days
**Assignee**: Treasury Queen + Society 4 dev

**Deliverables**:
1. Token pool data structures
2. ATP/ADP state transitions
3. Entity balance tracking
4. Stake lock mechanisms
5. Resource metering integration

**Implementation**:
```go
// File: blockchain/source/x/economy/types/token_pool.go
package types

type SocietyTokenPool struct {
    SocietyID       string
    TotalSupply     int64  // 100,000,000 (per spec)
    ATPCirculating  int64  // Charged tokens
    ADPCirculating  int64  // Discharged tokens
    EntityBalances  map[string]*TokenBalance
    StakeRegistry   map[string]*StakeLock
    CreatedAt       time.Time
    LastRecharge    time.Time
}

type TokenBalance struct {
    EntityID     string
    ATP          int64
    ADP          int64
    StakeLocks   []string  // References to StakeLock IDs
    LastActivity time.Time
}

type StakeLock struct {
    LockID       string
    EntityID     string
    Amount       int64
    Purpose      string  // "trust_query", "witness", "resource"
    LockedAt     time.Time
    UnlockAt     time.Time
    Released     bool
}

// State transitions
func (pool *SocietyTokenPool) ChargeADP(entityID string, amount int64) error {
    // ADP → ATP (value creation, energy input)
    balance := pool.EntityBalances[entityID]
    if balance.ADP < amount {
        return ErrInsufficientADP
    }
    balance.ADP -= amount
    balance.ATP += amount
    pool.ATPCirculating += amount
    pool.ADPCirculating -= amount
    return nil
}

func (pool *SocietyTokenPool) DischargeATP(entityID string, amount int64, activity string) error {
    // ATP → ADP (value consumption, work done)
    balance := pool.EntityBalances[entityID]
    if balance.ATP < amount {
        return ErrInsufficientATP
    }
    balance.ATP -= amount
    balance.ADP += amount
    pool.ATPCirculating -= amount
    pool.ADPCirculating += amount

    // Record activity for value tracking
    RecordValueCreation(entityID, activity, amount)
    return nil
}
```

**Initial Token Distribution**:
```go
initialDistribution := map[string]TokenBalance{
    "king_claudius":           {ATP: 200, ADP: 0},
    "security_queen":          {ATP: 150, ADP: 0},
    "law_oracle_queen":        {ATP: 100, ADP: 0},
    "treasury_queen":          {ATP: 100, ADP: 0},
    "coherence_queen":         {ATP: 100, ADP: 0},
    "synthesis_queen":         {ATP: 100, ADP: 0},
    "documentation_queen":     {ATP: 100, ADP: 0},
    "hardware_binding_queen":  {ATP: 75, ADP: 0},
    "federation_bridge_queen": {ATP: 75, ADP: 0},
    "society_treasury":        {ATP: 0, ADP: 99000},
}
```

**Validation**: Token invariants maintained (ATP + ADP = TotalSupply)

---

### Task 2.2: Implement T3/V3 Trust Tensors
**Priority**: HIGH
**Estimated Effort**: 3 days
**Assignee**: Society 4 dev

**Implementation**:
```go
// File: blockchain/source/x/trust/types/trust_tensor.go
type TrustTensor struct {
    SubjectID  string
    ObjectID   string

    // T3 Dimensions (Temperament)
    T1_Utility      float64  // How useful
    T2_Benevolence  float64  // How helpful
    T3_Reliability  float64  // How consistent (temporal!)

    // V3 Dimensions (Validity)
    V1_Ability      float64  // Technical capability
    V2_Integrity    float64  // Honesty/alignment
    V3_Verification float64  // Proof quality

    // Context
    ContextHash     string
    LastUpdated     time.Time
    EvidenceTrail   []TrustEvidence

    // Privacy (requires ATP stake)
    QueryCost       int64   // ATP required to query
}

func (tt *TrustTensor) UpdateOnTemporalSurprise(surpriseLevel float64) {
    // Integration with temporal authentication!
    if surpriseLevel > 0.6 {
        tt.T3_Reliability -= 0.2
        tt.V3_Verification -= 0.1
    }
    tt.RecordEvidence("temporal_surprise", surpriseLevel)
}
```

**Validation**: Trust queries require ATP stakes per privacy spec

---

## Phase 3: Federation Protocols (Weeks 4-5)

### Task 3.1: Refactor Pending Consensus to Ledger-Based
**Priority**: CRITICAL
**Estimated Effort**: 4 days
**Assignee**: Society 4 dev

**Current**: JSON files + git
**Target**: Blockchain ledger + MRH + RDF

**Implementation**:
```go
// File: blockchain/source/x/consensus/types/pending_decision.go
type PendingDecision struct {
    DecisionID       string   `json:"decision_id"`      // LCT format
    Type             string   `json:"type"`             // "vote", "transaction", "attestation"
    LedgerEntry      string   `json:"ledger_entry"`     // "block:12345:tx:67"
    LocalWitnesses   []string `json:"local_witnesses"`  // Self-witness while isolated
    FederationWitnessesRequired int `json:"federation_witnesses_required"`
    FederationWitnesses []Witness `json:"federation_witnesses"`
    Status           string   `json:"status"`           // "pending", "witnessed", "executed"
    NetworkContext   string   `json:"network_context"`  // "work_isolated", "home_federation"
    CreatedAt        time.Time
    ProcessedAt      *time.Time

    // RDF triple representation
    RDFTriples       []RDFTriple

    // MRH updates when processed
    MRHUpdates       []MRHRelationship
}

type Witness struct {
    WitnessID   string
    Signature   string
    Timestamp   time.Time
    Attestation string
}

type RDFTriple struct {
    Subject   string
    Predicate string
    Object    string
}

func (pd *PendingDecision) ToRDF() []RDFTriple {
    return []RDFTriple{
        {
            Subject:   pd.DecisionID,
            Predicate: "web4:proposedBy",
            Object:    "lct:web4:society:society4",
        },
        {
            Subject:   pd.DecisionID,
            Predicate: "web4:requiresWitnessQuorum",
            Object:    fmt.Sprintf("%d", pd.FederationWitnessesRequired),
        },
        {
            Subject:   pd.DecisionID,
            Predicate: "web4:networkContext",
            Object:    pd.NetworkContext,
        },
    }
}
```

**Migration Path**:
1. Keep existing `pending_consensus.py` for backward compatibility
2. Add ledger-based storage in parallel
3. Publish RDF triples to federation SPARQL endpoint
4. Migrate existing pending decisions to new format
5. Deprecate JSON file system

**Validation**: Pending decisions queryable via SPARQL, witness signatures valid

---

### Task 3.2: Implement Witness Quorum Protocols
**Priority**: HIGH
**Estimated Effort**: 3 days
**Assignee**: Society 4 dev + federation coordination

**Implementation**:
```go
type WitnessQuorum struct {
    DecisionID        string
    RequiredWitnesses int
    WitnessPool       []string  // Eligible witnesses
    Witnesses         []Witness
    QuorumMet         bool
    DiversityScore    float64
}

func (wq *WitnessQuorum) ValidateQuorum() error {
    if len(wq.Witnesses) < wq.RequiredWitnesses {
        return ErrInsufficientWitnesses
    }

    // Check witness diversity (not all from same society)
    societies := make(map[string]int)
    for _, w := range wq.Witnesses {
        society := ExtractSocietyFromLCT(w.WitnessID)
        societies[society]++
    }

    wq.DiversityScore = float64(len(societies)) / float64(len(wq.Witnesses))

    if wq.DiversityScore < 0.5 {
        return ErrInsufficientDiversity
    }

    wq.QuorumMet = true
    return nil
}
```

**Integration**: Temporal anomalies trigger witness requests to federation

---

### Task 3.3: Add MCP for Inter-Society Communication
**Priority**: MEDIUM
**Estimated Effort**: 3 days
**Assignee**: Federation Bridge Queen + dev

**Implementation**:
```go
// File: blockchain/source/x/mcp/client.go
type MCPClient struct {
    SocietyID      string
    TransportURL   string
    HPKEHandshake  *HPKESession
    MessageQueue   chan MCPMessage
}

type MCPMessage struct {
    From         string
    To           string
    MessageType  string  // "witness_request", "temporal_alert", "consensus_proposal"
    Payload      []byte
    Signature    string
    Timestamp    time.Time
}

func (c *MCPClient) BroadcastTemporalAnomaly(surprise float64, context string) error {
    msg := MCPMessage{
        From:        c.SocietyID,
        To:          "federation:all",
        MessageType: "temporal_alert",
        Payload:     json.Marshal(TemporalAlert{
            SurpriseFactor: surprise,
            Context:        context,
            Timestamp:      time.Now(),
        }),
    }
    return c.Send(msg)
}
```

**Validation**: Messages encrypted with HPKE, signatures valid

---

### Task 3.4: Publish RDF Triples for Federation Queries
**Priority**: MEDIUM
**Estimated Effort**: 2 days
**Assignee**: Society 4 dev

**Implementation**:
- Set up SPARQL endpoint: `http://society4.act.federation/sparql`
- Publish society state as RDF
- Enable federation-wide queries
- Implement update notifications

**Sample SPARQL Query**:
```sparql
# Query Society 4's temporal patterns
SELECT ?time ?network ?surprise
WHERE {
  ?event web4:society <lct:web4:society:society4> .
  ?event web4:timestamp ?time .
  ?event web4:network ?network .
  ?event web4:surpriseFactor ?surprise .
  FILTER (?surprise > 0.6)
}
ORDER BY DESC(?time)
LIMIT 10
```

**Validation**: Federation can query Society 4 state

---

## Phase 4: Innovation Contribution (Week 6)

### Task 4.1: Draft RFC for Temporal Authentication
**Priority**: MEDIUM
**Estimated Effort**: 2 days
**Assignee**: Society 4 + federation review

**Deliverables**:
```markdown
# RFC-TEMP-AUTH-001: Temporal Authentication Extension for Web4

## Abstract
Proposes temporal pattern analysis as authentication factor in web4 trust framework.

## Motivation
Mobile nodes and network-transitioning entities need authentication that accounts for expected vs actual temporal/spatial presence.

## Specification
- Map temporal surprise to T3_Reliability dimension
- Define surprise calculation algorithm
- Specify trust modulation protocol
- Integration with existing V3_Verification

## Implementation
Reference: Society 4's TEMPORAL_AUTHENTICATION.md

## Security Considerations
- Privacy: temporal patterns are sensitive
- Attack vectors: pattern manipulation
- Mitigation: witness diversity, MRH context

## Status
PROPOSED for web4 standard v1.1.0
```

**Validation**: RFC reviewed by Genesis and 2+ federation societies

---

### Task 4.2: Draft RFC for Reality KV Cache Pattern
**Priority**: MEDIUM
**Estimated Effort**: 2 days
**Assignee**: Society 4 + SAGE team

**Deliverables**:
```markdown
# RFC-REALITY-CACHE-001: Reality KV Cache Protocol for Cognitive Architectures

## Abstract
Standardizes assumption caching and surprise-driven invalidation for AI cognition in web4.

## Motivation
Efficient reality modeling requires caching with invalidation signals.

## Specification
- Assumption cache data structure
- Surprise-based invalidation thresholds
- Multi-modal coherence checking
- MRH as distributed cache

## Implementation
Reference: Society 4's REALITY_KV_CACHE.md
Integration: HRM/SAGE architecture

## Web4 Integration
- Assumptions stored in MRH
- Witnesses validate cache coherence
- ATP cost for cache queries
- Trust degradation on stale cache

## Status
PROPOSED for web4 cognitive extension
```

**Validation**: RFC reviewed by HRM team and federation

---

### Task 4.3: Document Society 4 as Reference Implementation
**Priority**: LOW
**Estimated Effort**: 2 days
**Assignee**: Documentation Queen

**Deliverables**:
- Case study: "Mobile Node with Hardware Sovereignty"
- Architecture documentation
- Lessons learned
- Best practices for mobile societies

---

## Implementation Timeline

```
Week 1: Oct 7-13
├── Mon-Tue: LCT structure implementation
├── Wed-Thu: Birth certificate process
└── Fri: Security Queen → Authority role mapping

Week 2: Oct 14-20
├── Mon-Thu: Law Oracle implementation
└── Fri: Law Oracle testing & validation

Week 3: Oct 21-27
├── Mon-Wed: ATP/ADP token pool
├── Thu: T3/V3 trust tensors
└── Fri: Economic system testing

Week 4: Oct 28-Nov 3
├── Mon-Tue: Pending consensus refactor
├── Wed: Witness quorum protocols
└── Thu-Fri: Testing & integration

Week 5: Nov 4-10
├── Mon-Tue: MCP implementation
├── Wed: RDF/SPARQL endpoint
└── Thu-Fri: Federation integration testing

Week 6: Nov 11-17
├── Mon-Tue: Temporal auth RFC
├── Wed-Thu: Reality cache RFC
└── Fri: Documentation & compliance review
```

## Success Criteria

### Must Have (Compliance)
- ✅ LCT structure implemented and validated
- ✅ Birth certificate obtained with 3+ witnesses
- ✅ Law Oracle publishing machine-readable laws
- ✅ ATP/ADP token pool operational
- ✅ Pending consensus on ledger with witnesses

### Should Have (Integration)
- ✅ MCP communication active
- ✅ RDF triples queryable
- ✅ T3/V3 trust tensors tracking
- ✅ Security Queen as Authority role

### Nice to Have (Innovation)
- ✅ Temporal auth RFC submitted
- ✅ Reality cache RFC submitted
- ✅ Reference implementation docs
- ✅ Federation adoption of patterns

## Risk Management

| Risk | Impact | Mitigation |
|------|--------|------------|
| Birth certificate delayed | HIGH | Start early, coordinate with Genesis |
| ATP/ADP complexity | MEDIUM | Reference existing implementations |
| Federation protocol changes | LOW | Monitor spec updates, stay flexible |
| Timeline overrun | MEDIUM | Prioritize critical path, defer nice-to-haves |

## Resource Requirements

**Development Time**: ~120 hours (3 weeks full-time equivalent)
**Coordination**: Genesis, Society 2, Sprout for birth certificate
**Infrastructure**: SPARQL endpoint, MCP server
**Documentation**: Ongoing throughout implementation

## Approval & Tracking

**Approved By**: King Claudius, Society 4
**Tracked In**: TodoWrite list + git commits
**Status Updates**: Weekly to federation
**Final Review**: web4-compliance-reviewer agent

---

**Let's build Society 4 into full web4 compliance! 🚀**

*Plan created: September 30, 2025*
*Target completion: November 15, 2025*
*Status: READY FOR EXECUTION*
