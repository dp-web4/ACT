# ACT Society Trust Graph Architecture

## Web4-Aligned Role-Contextual Trust

In Web4, trust and value are not properties of entities but relationships between them. This document describes the RDF-style trust graph implementation for the ACT society.

## Core Principles

### 1. Trust is Relational, Not Absolute
- No entity has an inherent "trust score"
- Trust exists only in relationships between entities
- Expressed as RDF triples: `<Subject> <Predicate> <Object> <Value>`

### 2. Trust is Contextual
- Trust for coding ≠ Trust for decision-making
- Each relationship has a specific context
- Context determines which aspect of trust applies

### 3. Trust is Directional
- A trusts B does not imply B trusts A
- Asymmetric relationships are the norm
- Bidirectional trust must be explicitly established

### 4. Trust Evolves Through Witnessing
- Initial trust is provisional
- Witnessed performance updates trust tensors
- Trust degrades without reinforcement

## T3 Tensor (Trust Tensor) - Role-Contextual

### Structure
```rdf
<FromRole> <trusts-for-{context}> <ToRole> [Talent, Training, Temperament]
```

### Examples
```rdf
<ACT-Build-Society> <trusts-for-orchestration> <ACT-Genesis-Queen> [0.9, 0.8, 0.9]
<ACT-Genesis-Queen> <trusts-for-alignment> <Web4-Alignment-Queen> [0.95, 0.9, 1.0]
<Web4-Alignment-Queen> <trusts-for-patterns> <Pattern-Validator> [0.8, 0.7, 0.9]
```

### Dimensions
- **Talent**: Natural aptitude for the role in this context
- **Training**: Acquired skills relevant to this context
- **Temperament**: Behavioral reliability in this context

## V3 Tensor (Value Tensor) - Role-Contextual

### Structure
```rdf
<FromRole> <values-{type}-from> <ToRole> [Valuation, Veracity, Validity]
```

### Examples
```rdf
<ACT-Build-Society> <values-creation-from> <LCT-Coder> [0.8, 0.7, 0.9]
<ACT-Genesis-Queen> <values-coordination-from> <Web4-Alignment-Queen> [0.9, 0.95, 0.9]
<Society> <values-storage-from> <Treasury> [0.1, 0.9, 1.0]
```

### Dimensions
- **Valuation**: Subjective worth of the value provided
- **Veracity**: Objective accuracy of value claims
- **Validity**: Confirmed delivery of promised value

## Trust Graph for ACT Society

### Hierarchical Trust Flow
```
Society
    └── trusts-for-orchestration → Genesis-Queen
            ├── trusts-for-alignment → Web4-Alignment-Queen
            │       └── trusts-for-patterns → Pattern-Validator
            ├── trusts-for-reality → Reality-Alignment-Queen
            │       └── trusts-for-impossibility → Impossibility-Detector
            ├── trusts-for-identity → LCT-Infrastructure-Queen
            │       └── trusts-for-coding → LCT-Coder
            ├── trusts-for-protocol → ACP-Protocol-Queen
            ├── trusts-for-governance → Demo-Society-Queen
            ├── trusts-for-integration → MCP-Bridge-Queen
            ├── trusts-for-interface → Client-Interface-Queen
            └── trusts-for-economy → ATP-Economy-Queen
```

### Value Flow Network
```
Creation (Workers) → Exchange (Queens) → Storage (Society)
```

## Implementation Details

### Blockchain Storage
Trust relationships are stored as:
- **Subject LCT**: The entity giving trust
- **Object LCT**: The entity receiving trust
- **Context**: The specific domain of trust
- **Tensor Values**: [T1, T2, T3] or [V1, V2, V3]
- **Timestamp**: When relationship was established
- **Witness Marks**: Confirmations of the relationship

### Query Patterns
```sql
-- Get all entities Genesis Queen trusts
SELECT * FROM trust_tensors WHERE from_lct = 'ACT-Genesis-Queen'

-- Get trust for specific context
SELECT * FROM trust_tensors 
WHERE from_lct = 'ACT-Genesis-Queen' 
AND to_lct = 'Web4-Alignment-Queen'
AND context = 'alignment'

-- Calculate transitive trust
WITH RECURSIVE trust_path AS (
  SELECT from_lct, to_lct, trust_value, 1 as depth
  FROM trust_tensors
  WHERE from_lct = 'Society'
  UNION ALL
  SELECT t.from_lct, t.to_lct, t.trust_value * tp.trust_value, tp.depth + 1
  FROM trust_tensors t
  JOIN trust_path tp ON t.from_lct = tp.to_lct
  WHERE tp.depth < 3
)
SELECT * FROM trust_path
```

### Energy Mechanics
- **ATP Discharge**: When trust is exercised (role performs work)
- **ADP Accumulation**: As work is performed
- **ADP→ATP Recharge**: When value is recognized by trust source
- **Trust Update**: Successful value delivery increases trust tensor

## Current Society State

### Entities (Role LCTs)
1. **ACT-Build-Society** - The society entity (owns all tokens)
2. **ACT-Genesis-Queen** - Meta-orchestrator role
3. **Web4-Alignment-Queen** - Pattern harmony validator
4. **Reality-Alignment-Queen** - Assumption checker
5. **LCT-Infrastructure-Queen** - Identity architect
6. **ACP-Protocol-Queen** - Protocol designer
7. **Demo-Society-Queen** - Governance builder
8. **MCP-Bridge-Queen** - Integration coordinator
9. **Client-Interface-Queen** - UI architect
10. **ATP-Economy-Queen** - Energy modeler
11. **Pattern-Validator** - Worker (Web4 alignment)
12. **LCT-Coder** - Worker (Infrastructure)

### Agent Assignments
All roles currently filled by: Claude instances
Future: Other agents can be invited to fill roles

### Trust Network Status
- Genesis block height: 14+
- Society treasury: 10M ATP, 10M ADP
- Trust relationships: Being established
- Energy flow: Ready to begin with work

## Usage

### Create Trust Relationship
```bash
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "RoleA" \
  --to-lct "RoleB" \
  --tensor-type "T3" \
  --context "specific-domain" \
  --talent "0.8" \
  --training "0.9" \
  --temperament "0.85" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y
```

### Query Trust Relationship
```bash
~/go/bin/racecar-webd query trusttensor relationship \
  --from "RoleA" \
  --to "RoleB" \
  --context "specific-domain"
```

## Philosophical Alignment

This implementation embodies Web4's core philosophy:
- **No central authority**: Trust emerges from relationships
- **Context-aware**: Trust is domain-specific
- **Evolution through use**: Performance updates trust
- **Energy follows trust**: ATP flows along trust edges
- **Witnessed presence**: All relationships are observable

The society becomes a living trust network where work flows along edges of established trust, value accumulates at nodes, and the entire system self-organizes through witnessed interactions.

---

*"In Web4, you don't have trust - you participate in trust relationships. You don't possess value - you create and exchange it. The network is the entity, the relationships are the substance."*