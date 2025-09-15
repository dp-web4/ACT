# Human-LCT Binding Specification

## Overview

Human-LCT binding establishes the cryptographic connection between a human's biological identity and their Linked Context Token in Web4. This binding is permanent, unforgeable, and forms the foundation of human participation in the trust-native internet.

## Binding Process

### 1. Initial Creation
```yaml
Human Identity Verification:
  - Biometric capture (optional)
  - Government ID verification (optional)
  - Social attestation (minimum 3 witnesses)
  - Physical presence proof

LCT Generation:
  - Ed25519 keypair creation
  - Root LCT minting
  - Initial MRH graph creation
  - Society registration
```

### 2. Binding Ceremony
```python
def bind_human_to_lct(human_identity, lct_params):
    # 1. Generate root LCT
    root_lct = LCT.generate(
        entity_type="human",
        binding_type="root",
        identity_proof=human_identity
    )
    
    # 2. Create witness records
    witnesses = gather_witnesses(minimum=3)
    for witness in witnesses:
        witness.sign_attestation(root_lct, human_identity)
    
    # 3. Register with society
    society.register_human_lct(
        lct=root_lct,
        witness_attestations=witnesses,
        birth_certificate=create_birth_cert()
    )
    
    # 4. Establish initial trust
    root_lct.trust_tensor = T3.initialize(
        competence=0.5,  # Neutral start
        reliability=0.5,
        transparency=1.0  # Full for root binding
    )
    
    return root_lct
```

### 3. Recovery Mechanisms
```yaml
Recovery Methods:
  Social Recovery:
    - Threshold signatures from trusted contacts
    - Minimum 3 of 5 guardians
    - Time-locked for security
  
  Biometric Recovery:
    - Encrypted biometric template
    - Distributed storage
    - Multi-factor required
  
  Legal Recovery:
    - Court order process
    - Society authority override
    - Full audit trail
```

## Security Properties

### Unforgeability
- **Witness requirement**: Multiple independent attestations
- **Society registration**: Immutable ledger record
- **Biometric binding**: Optional hardware-level security
- **Time-based strengthening**: Trust accumulates over time

### Privacy Protection
- **Selective disclosure**: Reveal only necessary attributes
- **Zero-knowledge proofs**: Prove properties without revealing data
- **Pseudonymous operation**: Multiple agent LCTs for different contexts
- **Right to be forgotten**: Deactivation without deletion

## Agent Pairing

### Creating Agent LCTs
```python
def create_agent_lct(human_root_lct, purpose):
    # 1. Generate agent keypair
    agent_keys = generate_keypair()
    
    # 2. Create pairing proof
    pairing = {
        "root_lct": human_root_lct.id,
        "agent_public_key": agent_keys.public,
        "purpose": purpose,
        "permissions": define_permissions(purpose),
        "expiry": calculate_expiry()
    }
    
    # 3. Sign with root LCT
    signature = human_root_lct.sign(pairing)
    
    # 4. Create agent LCT
    agent_lct = LCT.generate(
        entity_type="agent",
        binding_type="paired",
        parent_lct=human_root_lct,
        pairing_proof=signature
    )
    
    return agent_lct
```

### Permission Model
```yaml
Permission Levels:
  Read-Only:
    - View public data
    - Query services
    - No state changes
  
  Limited Write:
    - Specific actions only
    - Value limits
    - Time restrictions
  
  Full Delegation:
    - Act as human
    - Sign on behalf
    - Revocable anytime
```

## Implementation Requirements

### Client Application
- Secure key storage (hardware security module preferred)
- Biometric capture capability
- Network connectivity for witnessing
- Backup and recovery UI

### Society Integration
- Birth certificate issuance
- Witness network access
- Trust tensor tracking
- Recovery support

### Security Measures
- Multi-factor authentication
- Regular key rotation for agents
- Audit logging of all operations
- Breach detection and response

## User Experience

### Onboarding Flow
1. **Welcome**: Explain Web4 and LCTs
2. **Identity Verification**: Choose method (bio/gov/social)
3. **Witness Gathering**: Invite or find witnesses
4. **Binding Ceremony**: Generate keys with witnesses
5. **Backup Setup**: Configure recovery methods
6. **First Agent**: Create initial agent LCT

### Daily Operations
- **Single sign-on**: LCT as universal identity
- **Agent management**: Create, monitor, revoke
- **Trust dashboard**: View reputation and relationships
- **Transaction history**: Complete audit trail

## Compliance Considerations

### GDPR/Privacy
- Data minimization
- Purpose limitation
- Right to erasure (within constraints)
- Data portability

### KYC/AML
- Optional enhanced verification
- Transaction monitoring
- Suspicious activity reporting
- Regulatory reporting APIs

## Migration Path

### From Web2
```yaml
Gradual Migration:
  Phase 1: Create LCT alongside existing accounts
  Phase 2: Link LCT to existing identities
  Phase 3: Use LCT for authentication
  Phase 4: Migrate data and relationships
  Phase 5: Deprecate old accounts
```

### Interoperability
- OAuth/OIDC bridge for legacy systems
- SAML assertions from LCT
- API gateway for translation
- Progressive enhancement approach

## Standard Compliance

This specification complies with:
- Web4 core specification
- ACP framework requirements
- Society-Authority-Law governance
- ATP/ADP economic model

## Security Audit Checklist

- [ ] Key generation entropy verified
- [ ] Witness attestation signatures valid
- [ ] Society registration confirmed
- [ ] Recovery mechanisms tested
- [ ] Permission boundaries enforced
- [ ] Audit logs comprehensive
- [ ] Privacy controls functioning
- [ ] Breach detection active

---

*"Your LCT is not just your identity—it's your presence, your reputation, and your agency in the digital realm."*