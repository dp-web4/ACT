# ACT Protocol Specification v1.0

## Overview
ACT (Agentic Context Tool) enables entities to participate in Web4 societies through blockchain-verified identity and trust relationships.

## Core Components

### 1. Identity Layer (LCT)
- Linked Context Tokens provide unforgeable digital identity
- Each entity (human/AI) has unique LCT
- LCTs track context, capabilities, and relationships

### 2. Trust Layer (T3/V3)
- T3 Tensor: Talent, Training, Temperament
- V3 Tensor: Valuation, Veracity, Validity
- Trust is relational and contextual (RDF triples)

### 3. Energy Layer (ATP/ADP)
- ATP: Available energy for work
- ADP: Depleted energy awaiting recharge
- Discharge → Work → Value → Recharge cycle

### 4. Governance Layer (R6)
- Rules + Roles + Request + Reference + Resource → Result
- All actions must follow R6 pattern
- Web4-Alignment-Queen validates compliance

## Message Formats

### LCT Creation
```json
{
  "type": "mint-lct",
  "entity_name": "string",
  "entity_type": "role|agent|human",
  "metadata": {},
  "initial_atp": "number"
}
```

### Trust Relationship
```json
{
  "type": "create-trust",
  "from_lct": "string",
  "to_lct": "string",
  "context": "string",
  "tensor": [T1, T2, T3]
}
```

### R6 Action
```json
{
  "type": "r6-action",
  "rules": [],
  "role": "string",
  "request": "string",
  "reference": "string",
  "resource": "string",
  "expected_result": "string"
}
```

## State Transitions
1. Entity joins → LCT minted
2. Work requested → ATP discharged
3. Work performed → ADP accumulates
4. Value verified → ADP→ATP recharge
5. Trust updated → Tensor adjusted

## Consensus Mechanisms
- Trust-weighted voting (implemented)
- Queens approve domain actions
- Society validates major changes
- R6 compliance mandatory
