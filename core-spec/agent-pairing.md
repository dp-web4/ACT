# Agent Pairing Specification

## Overview

Agent pairing establishes the cryptographic relationship between a human's root LCT and their delegated agent LCTs. This enables humans to safely delegate specific authorities to AI agents while maintaining control and accountability.

## Core Concepts

### Agent Types
```yaml
Personal Assistant:
  - General purpose AI helper
  - Broad but limited permissions
  - Direct human supervision

Specialized Agent:
  - Domain-specific expertise
  - Narrow deep permissions
  - Task-focused operation

Autonomous Agent:
  - Self-directed operation
  - Trigger-based activation
  - Result-oriented goals
```

### Pairing Mechanics

```python
class AgentPairing:
    def __init__(self, human_lct, agent_config):
        self.human_lct = human_lct
        self.agent_config = agent_config
        self.shared_secret = None
        self.permissions = {}
        self.constraints = {}
    
    def establish_pairing(self):
        # 1. Generate agent keypair
        agent_keys = generate_x25519_keypair()
        
        # 2. Perform Diffie-Hellman exchange
        self.shared_secret = dh_exchange(
            human_lct.private_key,
            agent_keys.public_key
        )
        
        # 3. Derive session keys
        keys = derive_keys(self.shared_secret)
        
        # 4. Create pairing certificate
        certificate = {
            "human_lct": self.human_lct.id,
            "agent_public_key": agent_keys.public_key,
            "permissions": self.permissions,
            "constraints": self.constraints,
            "created_at": timestamp(),
            "expires_at": self.calculate_expiry()
        }
        
        # 5. Sign with human's root key
        signature = self.human_lct.sign(certificate)
        
        return PairedAgent(certificate, signature, keys)
```

## Permission Model

### Permission Scopes
```yaml
Financial:
  Read:
    - View balances
    - Transaction history
    - Market data
  Write:
    - Transfer limits
    - Approved recipients
    - Time windows

Communication:
  Read:
    - Email access
    - Message history
    - Contact lists
  Write:
    - Send on behalf
    - Draft creation
    - Response automation

Computation:
  Execute:
    - Resource limits
    - API access
    - Tool usage
  Delegate:
    - Sub-agent creation
    - Task distribution
    - Result aggregation
```

### Constraint Framework
```python
class Constraints:
    def __init__(self):
        self.time_limits = {}
        self.value_limits = {}
        self.rate_limits = {}
        self.geo_limits = {}
    
    def add_time_constraint(self, operation, window):
        """Restrict operation to time window"""
        self.time_limits[operation] = {
            "start": window.start,
            "end": window.end,
            "timezone": window.timezone
        }
    
    def add_value_constraint(self, operation, limit):
        """Limit value per operation"""
        self.value_limits[operation] = {
            "max_per_transaction": limit.per_tx,
            "max_per_day": limit.per_day,
            "max_total": limit.total
        }
    
    def add_rate_limit(self, operation, rate):
        """Limit operation frequency"""
        self.rate_limits[operation] = {
            "max_per_minute": rate.per_minute,
            "max_per_hour": rate.per_hour,
            "burst_limit": rate.burst
        }
    
    def validate(self, operation, context):
        """Check all constraints for operation"""
        for constraint_type in [self.time_limits, 
                               self.value_limits, 
                               self.rate_limits]:
            if not self._check_constraint(
                constraint_type, 
                operation, 
                context
            ):
                return False
        return True
```

## Delegation Patterns

### Direct Delegation
```yaml
Pattern: Human → Agent
Use Case: Personal assistant
Properties:
  - Single level
  - Direct accountability
  - Immediate revocation
  
Example:
  Human LCT → Email Agent
  Permissions: Read/send emails
  Constraints: No financial, 100 msgs/day
```

### Hierarchical Delegation
```yaml
Pattern: Human → Agent → Sub-agents
Use Case: Complex task management
Properties:
  - Multi-level
  - Cascading permissions
  - Tracked lineage

Example:
  Human LCT → Project Manager Agent
    → Code Review Agent
    → Testing Agent
    → Deployment Agent
```

### Cooperative Delegation
```yaml
Pattern: Human → Multiple Agents (parallel)
Use Case: Distributed processing
Properties:
  - Parallel operation
  - Shared context
  - Result synthesis

Example:
  Human LCT → Research Agent A
           → Research Agent B
           → Research Agent C
           → Synthesis Agent
```

## Revocation Mechanisms

### Immediate Revocation
```python
def revoke_agent_immediately(human_lct, agent_lct):
    # 1. Create revocation certificate
    revocation = {
        "agent_lct": agent_lct.id,
        "revoked_at": timestamp(),
        "reason": "immediate_revocation"
    }
    
    # 2. Sign with human's root key
    signature = human_lct.sign(revocation)
    
    # 3. Broadcast to network
    broadcast_revocation(revocation, signature)
    
    # 4. Update local state
    agent_lct.status = "revoked"
    
    return revocation
```

### Scheduled Revocation
```python
def schedule_revocation(human_lct, agent_lct, when):
    # Create future-dated revocation
    scheduled = {
        "agent_lct": agent_lct.id,
        "revoke_after": when,
        "type": "scheduled"
    }
    
    # Sign and store
    signature = human_lct.sign(scheduled)
    store_scheduled_revocation(scheduled, signature)
    
    return scheduled
```

### Conditional Revocation
```python
def conditional_revocation(human_lct, agent_lct, conditions):
    # Revoke if conditions met
    conditional = {
        "agent_lct": agent_lct.id,
        "conditions": conditions,
        "type": "conditional"
    }
    
    # Examples of conditions:
    # - Value threshold exceeded
    # - Geographic boundary crossed
    # - Suspicious activity detected
    # - Trust score drops below threshold
    
    return conditional
```

## Security Considerations

### Key Management
- Agent keys never leave secure enclave
- Regular key rotation schedule
- Quantum-resistant algorithms ready
- Hardware security module integration

### Audit Trail
```python
class AuditLog:
    def log_pairing(self, human_lct, agent_lct, permissions):
        entry = {
            "event": "agent_pairing",
            "human": human_lct.id,
            "agent": agent_lct.id,
            "permissions": permissions,
            "timestamp": timestamp(),
            "witness": gather_witnesses()
        }
        self.append_to_ledger(entry)
    
    def log_action(self, agent_lct, action, result):
        entry = {
            "event": "agent_action",
            "agent": agent_lct.id,
            "action": action,
            "result": result,
            "timestamp": timestamp(),
            "proof_of_agency": generate_proof()
        }
        self.append_to_ledger(entry)
```

### Attack Mitigation
- Replay attack prevention via nonces
- Man-in-the-middle protection via TLS
- Privilege escalation detection
- Anomaly detection on agent behavior

## Implementation Checklist

### Client Requirements
- [ ] Secure key generation
- [ ] Pairing certificate creation
- [ ] Permission interface
- [ ] Constraint configuration
- [ ] Revocation mechanisms
- [ ] Audit log viewer

### Agent Requirements
- [ ] Key storage
- [ ] Permission checking
- [ ] Constraint enforcement
- [ ] Proof-of-agency generation
- [ ] Result reporting
- [ ] Graceful revocation handling

### Network Requirements
- [ ] Certificate distribution
- [ ] Revocation propagation
- [ ] Witness network integration
- [ ] Trust tensor updates
- [ ] Audit log aggregation

## Standard Compliance

Complies with:
- Web4 LCT specification
- Agency Delegation (AGY) framework
- Agentic Context Protocol (ACP)
- Trust tensor (T3) calculations

---

*"Agent pairing transforms AI from a tool into a trusted delegate, acting with your authority but within your boundaries."*