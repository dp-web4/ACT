# Society 4 Law Oracle

**Version**: 1.0.0
**Status**: Active
**Published**: October 1, 2025
**Oracle LCT**: `lct:web4:role:law_oracle:society4`

## Overview

Society 4's Law Oracle publishes machine-readable laws following the web4 Society-Authority-Law (SAL) specification. These laws govern identity, governance, economics, security, and federation interactions.

## Law Dataset

**Location**: `society4_law_oracle_v1.0.0.json`
**Format**: JSON-LD (web4 Law Dataset)
**ID**: `web4://law/society4/1.0.0`
**Hash**: `sha256:pending` (will be computed and signed)

## Law Categories

### 1. Identity Laws (LAW-IDENTITY-*)

**LAW-IDENTITY-001**: Hardware Binding Requirement
- All operations require valid hardware binding
- Severity: CRITICAL
- Enforcement: MANDATORY

**LAW-IDENTITY-002**: Genesis Hardware Hash
- Binds to specific WSL2 hardware: `93e766842ee7882a...`
- Prevents hardware spoofing
- Severity: CRITICAL

**LAW-IDENTITY-003**: Temporal Authentication
- Enables temporal pattern matching (RFC-TEMP-AUTH-001)
- Detects anomalous network presence
- Severity: HIGH
- Enforcement: RECOMMENDED

### 2. Governance Laws (LAW-GOV-*)

**LAW-GOV-001**: Monarchic Approval
- King Claudius must approve major transactions
- Monarchic governance model
- Severity: CRITICAL

**LAW-GOV-002**: Security Queen Veto Power
- Security Queen can veto ANY operation on security grounds
- Per federation mandate
- Severity: CRITICAL

**LAW-GOV-003**: Pending Consensus Network Restriction
- Pending decisions processed only from home network (10.0.0.x)
- Prevents work network compromise
- Severity: HIGH

**LAW-GOV-004**: Queens Quorum
- Major decisions require 5 of 8 queens
- Ensures distributed governance
- Severity: HIGH

### 3. Economic Laws (LAW-ECON-*)

**LAW-ECON-001**: Total ATP Budget
- Society 4 total: 1000 ATP
- Prevents inflation
- Severity: CRITICAL

**LAW-ECON-002**: Security Queen ATP Allocation
- Security Queen: 150 ATP (highest)
- Reflects critical importance
- Severity: HIGH

**LAW-ECON-003**: Daily ATP Recharge
- All queens: +20 ATP daily
- Capped at initial allocation
- Severity: MEDIUM

**LAW-ECON-004**: ATP Stake for Trust Queries
- Trust queries cost minimum 5 ATP
- Privacy protection mechanism
- Severity: MEDIUM

### 4. Security Laws (LAW-SEC-*)

**LAW-SEC-001**: Hardware Validation Frequency
- Validate on every critical operation
- Real-time security enforcement
- Severity: CRITICAL

**LAW-SEC-002**: Emergency Halt Authority
- Only Security Queen and King can halt operations
- Emergency response protocol
- Severity: CRITICAL

**LAW-SEC-003**: Key Rotation Period
- Quarterly rotation (90 days)
- Prevents key staleness
- Severity: HIGH

### 5. Federation Laws (LAW-FED-*)

**LAW-FED-001**: Witness Quorum
- Minimum 3 federation witnesses
- Distributed validation
- Severity: HIGH

**LAW-FED-002**: Network Mobility Allowed
- Society 4 may move between networks
- Unique mobile node characteristic
- Severity: MEDIUM

**LAW-FED-003**: Temporal Surprise Threshold
- Surprise > 0.6 triggers enhanced auth
- Implements RFC-TEMP-AUTH-001
- Severity: MEDIUM

## Procedures

### PROC-WIT-3: Minimum Witness Requirement
- 3 independent federation witnesses
- Standard for all federation operations

### PROC-EMERGENCY: Emergency Halt
- Security Queen immediate halt (no quorum)
- For security breaches only

### PROC-CONSENSUS: Pending Consensus Processing
- Triggered on home network connection
- Requires 3 witnesses
- Processes queued decisions

### PROC-HARDWARE-VERIFY: Hardware Binding Verification
- Extract and compare hash
- Rejects transaction on mismatch
- Every critical operation

### PROC-TEMPORAL-AUTH: Temporal Authentication
- Computes surprise factor
- Thresholds: 0.3 (low), 0.6 (medium), 0.8 (high)
- Actions scale with surprise level
- Implements RFC-TEMP-AUTH-001

### PROC-ATP-RECHARGE: Daily ATP Regeneration
- Daily at 00:00 UTC
- +20 ATP to all queens
- Capped at initial allocation

### PROC-QUEEN-CONSENSUS: Queens Quorum Voting
- Requires 5/8 queens
- 48-hour voting period
- Security Queen has veto

## Interpretations

The Law Oracle maintains canonical interpretations for common questions:

### INT-001: Hardware Binding Philosophy
> "Hardware sovereignty ensures Society 4's identity cannot be replicated or forged. The WSL2 hardware serves as root of trust, analogous to TPM in production systems."

### INT-002: Network Mobility vs Security
> "Mobile nodes present unique challenges. Society 4 addresses this through: (1) Hardware binding remains constant, (2) Temporal authentication detects anomalies, (3) Pending consensus queues decisions while isolated, (4) Home network required for federation operations."

### INT-003: Security Queen Authority
> "Yes, on security matters only. While King Claudius is sovereign, Security Queen has veto power over ANY operation if security risk is detected. This is per federation mandate following the Great Key Exposure Incident."

### INT-004: ATP Economic Model
> "ATP (charged) and ADP (discharged) tokens create energy economy mirroring biological ATP. Work discharges ATP→ADP, value creation charges ADP→ATP. This prevents resource hoarding and incentivizes productive contribution."

## Exceptions

### EXC-001: Genesis Bootstrap Exception
- **Scope**: Birth certificate
- **Condition**: Genesis LCT creation
- **Allowed Violation**: Null birth certificate
- **Rationale**: First entity cannot have witnesses before existence
- **Expires**: When birth certificate obtained

### EXC-002: Work Network Isolation
- **Scope**: Federation access
- **Condition**: Network = work_isolated (172.25.x.x)
- **Allowed Violation**: Direct P2P unavailable
- **Rationale**: Work network cannot reach home federation
- **Mitigation**: Pending consensus queue

## R6 Bindings

Society 4 laws map to R6 (Rules + Role + Request + Reference + Resource → Result) action grammar:

- `web4://schemas/r6-rules-v1` - Core R6 rules schema
- `web4://schemas/society-governance-v1` - Governance patterns
- `web4://schemas/hardware-sovereignty-v1` - Hardware binding rules

## Publication

**Endpoint**: `web4://law/society4/latest`
**SPARQL**: `http://society4.act.federation/sparql/law`
**Format**: JSON-LD
**Update Policy**: Monthly review, quarterly revision

## Compliance

**Level**: web4-core-v1.0
**Extensions**:
- temporal-auth (RFC-TEMP-AUTH-001)
- reality-cache (RFC-REALITY-CACHE-001)
- network-mobility (Society 4 innovation)

## Signatures

Laws must be signed by:
1. **Law Oracle Queen** - Publisher
2. **King Claudius** - Sovereign approval

Signatures pending cryptographic key generation.

## Queries

### Example Law Queries

**Check if action is compliant**:
```sparql
PREFIX law: <web4://law/society4/>
SELECT ?norm ?enforcement
WHERE {
  ?norm law:selector "r6.transaction.approval" .
  ?norm law:enforcement ?enforcement .
}
```

**Get all critical laws**:
```sparql
PREFIX law: <web4://law/society4/>
SELECT ?id ?name
WHERE {
  ?norm law:id ?id .
  ?norm law:name ?name .
  ?norm law:severity "critical" .
}
```

**Find ATP-related laws**:
```sparql
PREFIX law: <web4://law/society4/>
SELECT ?id ?description
WHERE {
  ?norm law:category "economic" .
  ?norm law:id ?id .
  ?norm law:description ?description .
  FILTER(CONTAINS(?description, "ATP"))
}
```

## Next Steps

1. **Cryptographic Signing**
   - Generate keys for Law Oracle Queen and King Claudius
   - Sign law dataset with COSE
   - Compute SHA-256 hash

2. **Publication**
   - Publish to web4://law/society4/1.0.0
   - Set up SPARQL endpoint
   - Enable federation queries

3. **Integration**
   - Integrate with blockchain consensus
   - Connect R6 action grammar
   - Enable automated compliance checking

4. **Maintenance**
   - Monthly review (November 1, 2025)
   - Quarterly revision as needed
   - Track law evolution and interpretations

---

**Law Oracle Queen**: Society 4
**Approved By**: King Claudius
**Effective**: October 1, 2025
**Compliance**: web4-core-v1.0 + extensions
