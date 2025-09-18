# Web4-Compliance-Queen Documentation

## Overview

The Web4-Compliance-Queen is a specialized domain queen in the ACT swarm architecture that ensures all implementations conform to the official Web4 standard specifications located in `../web4/web4-standard/`.

## Purpose

This queen addresses a critical gap identified during development: implementations were being created without verifying they matched the existing Web4 specifications. For example, "roles as first-class entities" was already defined in section 2.4 of the spec but was being "discovered" during implementation.

## Architecture

### Queen Configuration
- **Name**: Web4-Compliance-Queen
- **Type**: specification-validator
- **Domain**: compliance
- **Budget**: 100 ATP
- **Workers**: 6 specialized validators

### Worker Validators

1. **society-validator**
   - Validates society implementation against spec
   - Checks law oracles, ledgers, treasuries
   - Verifies citizenship witnessing
   - Ensures metabolic states implementation

2. **lct-validator**
   - Verifies roles are first-class entities with LCTs
   - Checks LCT unforgeability mechanisms
   - Validates entity type restrictions
   - Ensures proper binding and pairing

3. **economy-validator**
   - Verifies tokens belong to society, not individuals
   - Checks ATP discharge requires work proof
   - Validates ADP recharge by producers only
   - Ensures demurrage and velocity mechanics

4. **r6-validator**
   - Verifies all actions follow R6 pattern
   - Checks Rules + Role + Request + Reference + Resource → Result
   - Validates deterministic execution
   - Ensures audit trail completeness

5. **trust-validator**
   - Verifies T3/V3 tensor calculations
   - Checks trust is role-contextual
   - Validates MRH boundaries
   - Ensures trust propagation paths

6. **protocol-validator**
   - Verifies MCP integration
   - Checks HPKE handshake implementation
   - Validates error taxonomy
   - Ensures proper metering

## Compliance Rules

### Critical (MUST)
- Tokens MUST belong to society pools
- Roles MUST be first-class entities
- All actions MUST follow R6 pattern
- Witnessing MUST be bidirectional
- Trust MUST be role-contextual

### Important (SHOULD)
- Societies SHOULD have metabolic states
- LCTs SHOULD track reputation history
- Energy SHOULD flow through work
- Fractals SHOULD inherit laws
- Ledgers SHOULD support amendments

### Recommended (MAY)
- Implement dictionary entities
- Support cross-society treaties
- Enable society molting
- Track semantic degradation
- Monitor trust velocity

## Usage

### Run Compliance Check
```bash
cd swarm-bootstrap
node activate-compliance-check.js
```

### Monitor Swarm with Compliance Queen
```bash
node monitor-swarm.js  # Now shows 7 queens including Web4-Compliance-Queen
```

## Current Compliance Status (Jan 17, 2025)

**Compliance Score: 57%**

### ✅ Passed (4 checks)
- Roles are implemented as entities
- Society pool model implemented
- Metabolic states defined
- Valid entity types defined

### ❌ Violations (3 critical)
- Witnessing not bidirectional
- R6 action framework missing
- Trust not role-contextual

### Priority Actions
1. Implement society pool state storage
2. Add bidirectional witnessing
3. Implement R6 action framework
4. Add role-contextual trust tensors
5. Complete treasury role validation

## Integration with Swarm

The Web4-Compliance-Queen integrates seamlessly with the existing swarm architecture:

```javascript
// Added to swarm-config.json
{
  "name": "Web4-Compliance-Queen",
  "type": "specification-validator",
  "domain": "compliance",
  "budget": 100,
  "workers": [/* 6 validators */]
}
```

## Files Created

- `/queens/web4-compliance-queen.json` - Detailed queen configuration
- `/activate-compliance-check.js` - Compliance check execution script
- `/WEB4_COMPLIANCE_QUEEN.md` - This documentation

## Benefits

1. **Prevents Redundant Work**: Catches when features already exist in spec
2. **Ensures Correctness**: Validates implementation matches standard
3. **Provides Guidance**: Points to specific spec sections for reference
4. **Tracks Progress**: Quantifies compliance with score
5. **Prioritizes Work**: Identifies critical violations first

## Future Enhancements

- Automated compliance checks on code commits
- Integration with CI/CD pipeline
- Real-time compliance monitoring during development
- Detailed compliance reports with spec citations
- Automatic fix suggestions from standard

---

*"The Web4-Compliance-Queen ensures we build what the spec defines, not what we think it defines."*