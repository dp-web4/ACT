# Web4 Compliance Work Plan for ACT Ledger

## Executive Summary

This work plan outlines the modifications needed to adapt the existing ModBatt blockchain implementation to serve as the ledger for the ACT prototype society, fully compliant with the Web4 standard.

## Current State Analysis

### Existing Implementation
- **Base**: Cosmos SDK blockchain with Tendermint consensus
- **Origin**: ModBatt battery management system (racecarweb)
- **Modules**: Component Registry, Energy Cycle, Trust Tensor, LCT Manager, Pairing Queue
- **Architecture**: Well-structured with proper separation of concerns
- **Testing**: Comprehensive real blockchain testing infrastructure

### Key Strengths
1. **Solid Foundation**: Cosmos SDK provides robust blockchain infrastructure
2. **Modular Design**: Clean module separation aligns with Web4 components
3. **Trust Concepts**: Already implements trust tensors (T3-like)
4. **Energy Tokens**: Has ATP/ADP token concepts (needs alignment)
5. **Component Registry**: Can be adapted for LCT registry

### Critical Gaps for Web4 Compliance

#### 1. LCT Implementation (HIGH PRIORITY)
**Current**: Basic LCT struct with limited fields
**Required**: Full Web4 LCT specification
- Missing Ed25519/X25519 cryptographic identity
- No MRH (Markov Relevancy Horizon) graphs
- No fractal context boundaries
- No proper entity type system

#### 2. MRH Integration (HIGH PRIORITY)
**Current**: No MRH implementation
**Required**: Complete MRH system
- RDF graph storage and traversal
- Context horizon calculations
- Fractal boundary management
- Witness relationship tracking

#### 3. Trust Tensor Alignment (MEDIUM PRIORITY)
**Current**: Uses Talent/Training/Temperament (T3)
**Required**: Web4 T3 (Competence/Reliability/Transparency)
- Rename tensor dimensions
- Add V3 (Value tensor) implementation
- Implement trust-as-gravity calculations
- Add trust degradation over distance

#### 4. ATP/ADP Cycle (MEDIUM PRIORITY)
**Current**: Basic energy tokens
**Required**: Full Web4 ATP/ADP specification
- Semifungible token implementation
- Society-managed pools
- Anti-hoarding mechanisms
- Proper charge/discharge cycle

#### 5. SAL Framework (HIGH PRIORITY)
**Current**: No governance framework
**Required**: Society-Authority-Law implementation
- Birth certificate issuance
- Law oracle for rule enforcement
- Authority delegation mechanism
- Immutable governance records

#### 6. Agency Delegation (HIGH PRIORITY)
**Current**: Basic pairing mechanism
**Required**: Full AGY framework
- Proof-of-agency generation
- Delegation certificates
- Revocation mechanisms
- Audit trail

#### 7. ACP Implementation (HIGH PRIORITY)
**Current**: No autonomous agent support
**Required**: Agentic Context Protocol
- Agent Plans with triggers
- Intent generation
- Decision collection
- Execution records

#### 8. Dictionary Entities (LOW PRIORITY - Phase 2)
**Current**: None
**Required**: Living dictionary entities
- Semantic bridging
- Compression-trust management
- Cross-domain translation

#### 9. R6 Framework (MEDIUM PRIORITY)
**Current**: No formal action framework
**Required**: R6 transaction validation
- Rules + Role + Request + Reference + Resource → Result
- All transactions must follow R6 pattern

## Implementation Work Plan

### Phase 1: Core Web4 Identity (Week 1-2)

#### 1.1 LCT Module Overhaul
```go
// New LCT structure
type LCT struct {
    // Core Identity
    ID           string    // lct:web4:act:entity:uuid
    EntityType   string    // human|agent|dictionary|society
    PublicKey    []byte    // Ed25519 public key
    
    // MRH Context
    MRHGraph     string    // IPFS hash of RDF graph
    ContextDepth uint32    // Fractal depth
    
    // Trust Tensors
    T3Competence    float64
    T3Reliability   float64
    T3Transparency  float64
    V3Value         float64
    
    // Relationships
    ParentLCT    string    // For agents
    Witnesses    []string  // LCT IDs
    
    // Lifecycle
    BirthCert    string    // Society-issued
    CreatedAt    time.Time
    RevokedAt    *time.Time
    Status       string    // active|revoked|suspended
}
```

#### 1.2 MRH Implementation
- Add RDF graph storage (use IPFS for actual graphs)
- Implement context boundary calculations
- Add witness relationship management
- Create MRH traversal algorithms

#### 1.3 Cryptographic Updates
- Implement Ed25519 key generation
- Add X25519 for Diffie-Hellman
- Zero on-chain key storage (already good)
- Add key derivation for agent LCTs

### Phase 2: Trust and Value (Week 2-3)

#### 2.1 Trust Tensor Alignment
```go
// Rename existing tensor dimensions
type TrustTensor struct {
    TensorID     string
    LctID        string
    
    // Old -> New mapping
    Competence   float64  // was TalentScore
    Reliability  float64  // was TrainingScore  
    Transparency float64  // was TemperamentScore
    
    // Add value tensor
    ValueTensor  float64
    
    // Trust gravity
    TrustMass    float64  // Calculated from T3
    TrustRadius  float64  // Influence distance
}
```

#### 2.2 ATP/ADP Implementation
```go
type ATPToken struct {
    TokenID      string
    State        string    // charged|discharged
    Amount       uint64
    
    // Society management
    SocietyPool  string    // Society that manages it
    
    // Anti-hoarding
    LastUsed     time.Time
    DecayRate    float64
    
    // Provenance
    ChargedBy    string    // Producer LCT
    ChargeProof  string    // Proof of value creation
}
```

### Phase 3: Governance Framework (Week 3-4)

#### 3.1 SAL Implementation
```go
type Society struct {
    SocietyID    string
    Name         string
    
    // Governance
    LawOracle    string    // Smart contract address
    Authority    []string  // LCT IDs with authority
    
    // Economics
    ATPPool      uint64    // Total ATP in circulation
    ADPPool      uint64    // Total ADP available
    
    // Members
    Citizens     []string  // LCT IDs
}

type BirthCertificate struct {
    CertID       string
    LctID        string
    SocietyID    string
    IssuedAt     time.Time
    Witnesses    []string
    Metadata     map[string]string
}
```

#### 3.2 Law Oracle
- Implement rule engine for society laws
- Add permission validation
- Create dispute resolution mechanism
- Implement governance voting

### Phase 4: Agency and Autonomy (Week 4-5)

#### 4.1 Agency Delegation (AGY)
```go
type AgencyDelegation struct {
    DelegationID string
    FromLCT      string    // Delegator
    ToLCT        string    // Agent
    
    // Permissions
    Scope        []string  // What can be done
    Constraints  map[string]interface{}
    
    // Temporal
    ValidFrom    time.Time
    ValidUntil   time.Time
    
    // Proof
    ProofOfAgency string   // Cryptographic proof
    Signature     []byte
}
```

#### 4.2 Agentic Context Protocol (ACP)
```go
type AgentPlan struct {
    PlanID       string
    AgentLCT     string
    
    // Triggers
    Triggers     []Trigger
    
    // Actions
    Actions      []Action
    
    // Constraints
    RequiresApproval bool
    ATPBudget       uint64
}

type Intent struct {
    IntentID     string
    PlanID       string
    Context      map[string]interface{}
    GeneratedAt  time.Time
}

type ExecutionRecord struct {
    RecordID     string
    IntentID     string
    Result       interface{}
    ProofOfWork  string
    Witnesses    []string
}
```

### Phase 5: Integration and Testing (Week 5-6)

#### 5.1 Module Integration
- Wire all new modules together
- Update existing modules for Web4 compliance
- Implement cross-module interactions
- Add comprehensive event emission

#### 5.2 Demo Society Setup
- Create genesis society configuration
- Implement bootstrap witnesses
- Set up initial ATP/ADP pools
- Configure law oracle with basic rules

#### 5.3 Testing Infrastructure
- Adapt existing test framework
- Add Web4-specific test scenarios
- Implement compliance validation tests
- Performance benchmarking

### Phase 6: ACT-Specific Features (Week 6-7)

#### 6.1 Human-LCT Binding
- Implement secure binding ceremony
- Add witness attestation flow
- Create recovery mechanisms
- Implement social recovery

#### 6.2 MCP Bridge Integration
- Add MCP server connectors
- Implement trust-aware routing
- Add ATP metering for MCP calls
- Create result attestation

#### 6.3 UI Backend Support
- Create REST API endpoints
- Add WebSocket for real-time updates
- Implement query optimizations
- Add caching layer

## Technical Modifications

### Module Renaming
```bash
# Rename modules to align with Web4
x/componentregistry -> x/lctregistry
x/lctmanager -> x/lctlifecycle  
x/trusttensor -> x/trustvalue
x/energycycle -> x/atpcycle
x/pairingqueue -> x/agency
```

### New Modules to Create
```bash
x/society      # SAL framework
x/mrh          # Markov Relevancy Horizon
x/acp          # Agentic Context Protocol
x/dictionary   # Dictionary entities (Phase 2)
```

### Proto Updates
- Update all proto definitions for Web4 types
- Add new message types for Web4 operations
- Implement proper versioning
- Add backward compatibility where needed

## Migration Strategy

### Data Migration
1. Map existing components to LCTs
2. Convert trust scores to T3 format
3. Transform energy tokens to ATP/ADP
4. Create initial society structure

### Backward Compatibility
- Keep legacy fields temporarily
- Provide migration transactions
- Support gradual transition
- Document breaking changes

## Testing Requirements

### Unit Tests
- [ ] LCT creation and management
- [ ] MRH graph operations
- [ ] Trust tensor calculations
- [ ] ATP/ADP lifecycle
- [ ] Agency delegation
- [ ] ACP execution

### Integration Tests
- [ ] Human-LCT binding flow
- [ ] Agent pairing process
- [ ] Trust propagation
- [ ] Value flow through system
- [ ] Governance operations

### Compliance Tests
- [ ] Web4 standard validation
- [ ] Security requirements
- [ ] Performance benchmarks
- [ ] Scalability tests

## Risk Mitigation

### Technical Risks
- **Complexity**: Incremental implementation with testing
- **Performance**: Optimize critical paths, use caching
- **Storage**: Use IPFS for large data (MRH graphs)

### Security Risks
- **Key Management**: Zero on-chain storage (already good)
- **Authority Abuse**: Multi-signature requirements
- **Sybil Attacks**: Witness requirements and trust scores

## Success Metrics

### Phase 1 Complete
- [ ] LCTs fully Web4 compliant
- [ ] MRH graphs functional
- [ ] Cryptographic identity working

### Phase 2 Complete
- [ ] Trust tensors aligned with Web4
- [ ] ATP/ADP cycle operational
- [ ] Value tracking functional

### Phase 3 Complete
- [ ] Society governance active
- [ ] Law oracle enforcing rules
- [ ] Birth certificates issued

### Phase 4 Complete
- [ ] Agency delegation working
- [ ] ACP plans executing
- [ ] Autonomous agents operational

### Phase 5 Complete
- [ ] All modules integrated
- [ ] Demo society running
- [ ] Tests passing

### Phase 6 Complete
- [ ] Human-LCT binding functional
- [ ] MCP bridges operational
- [ ] UI backend ready

## Resource Requirements

### Development
- 6-7 weeks of focused development
- Cosmos SDK expertise
- Web4 standard knowledge
- Testing infrastructure

### Infrastructure
- Development blockchain node
- IPFS node for MRH graphs
- PostgreSQL for indexing
- Redis for caching

## Conclusion

The existing ModBatt blockchain provides a solid foundation for the ACT ledger. With systematic modifications following this work plan, we can achieve full Web4 compliance while maintaining the robustness of the Cosmos SDK infrastructure. The modular approach allows for incremental implementation and testing, reducing risk and ensuring quality.

The transformation will create a production-ready ledger for the ACT prototype society, demonstrating Web4's vision of trust-native computing where humans and AI agents interact as peers in a cryptographically-secured, economically-aligned ecosystem.

---

*"From battery management to trust management—the same principles of energy flow apply to value and agency in Web4."*