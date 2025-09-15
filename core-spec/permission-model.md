# Permission Model Specification

## Overview

The ACT permission model defines how authority flows from human root LCTs through agent LCTs to actions in the Web4 ecosystem. It implements principle of least privilege with granular, revocable, and auditable permissions.

## Core Principles

### 1. Least Privilege
Every agent receives only the minimum permissions required for its intended function.

### 2. Explicit Grant
No implicit permissions - everything must be explicitly granted by the human.

### 3. Temporal Bounds
All permissions have time limits, either explicit or inherited from the agent pairing.

### 4. Value Limits
Financial and resource permissions always include value constraints.

### 5. Full Auditability
Every permission grant, use, and revocation is recorded immutably.

## Permission Taxonomy

### Categories
```yaml
Identity:
  profile:
    read: View profile information
    write: Modify profile data
    share: Share with third parties
  
  credentials:
    use: Authenticate as human
    delegate: Create sub-agents
    attest: Sign attestations

Financial:
  balance:
    read: View ATP/ADP balances
    query: Check transaction history
  
  transfer:
    send: Send tokens
    receive: Accept tokens
    exchange: Trade currencies
  
  limits:
    per_transaction: Maximum per operation
    per_day: Daily spending limit
    total: Absolute maximum

Communication:
  messaging:
    read: Access messages
    send: Send messages
    delete: Remove messages
  
  channels:
    email: Email access
    social: Social media
    mcp: MCP server communication

Computation:
  resources:
    cpu: Processing allocation
    memory: RAM limits
    storage: Disk quotas
  
  tools:
    read: File system access
    write: File creation/modification
    execute: Run programs
    network: API calls

Governance:
  voting:
    view: See proposals
    vote: Cast votes
    propose: Create proposals
  
  witnessing:
    observe: Act as witness
    attest: Sign attestations
    validate: Verify claims
```

## Permission Specification Language

### Syntax
```python
class Permission:
    def __init__(self, scope, action, constraints=None):
        self.scope = scope      # e.g., "financial.transfer"
        self.action = action    # e.g., "send"
        self.constraints = constraints or {}
    
    def to_dict(self):
        return {
            "scope": self.scope,
            "action": self.action,
            "constraints": self.constraints,
            "granted_at": timestamp(),
            "expires_at": self.calculate_expiry()
        }

# Example permission specifications
permissions = [
    Permission(
        scope="financial.transfer",
        action="send",
        constraints={
            "max_amount": 100,
            "currency": "ATP",
            "recipients": ["whitelist_only"],
            "daily_limit": 1000
        }
    ),
    Permission(
        scope="communication.messaging",
        action="send",
        constraints={
            "channels": ["email"],
            "rate_limit": "10/hour",
            "require_approval": True
        }
    ),
    Permission(
        scope="computation.tools",
        action="read",
        constraints={
            "paths": ["/shared/data/*"],
            "exclude": ["*.key", "*.secret"]
        }
    )
]
```

## Constraint Types

### Value Constraints
```python
class ValueConstraint:
    def __init__(self, max_per_tx=None, daily_limit=None, 
                 total_limit=None, currency="ATP"):
        self.max_per_tx = max_per_tx
        self.daily_limit = daily_limit
        self.total_limit = total_limit
        self.currency = currency
        self.spent_today = 0
        self.spent_total = 0
    
    def validate(self, amount):
        if self.max_per_tx and amount > self.max_per_tx:
            return False, "Exceeds per-transaction limit"
        
        if self.daily_limit and self.spent_today + amount > self.daily_limit:
            return False, "Exceeds daily limit"
        
        if self.total_limit and self.spent_total + amount > self.total_limit:
            return False, "Exceeds total limit"
        
        return True, "OK"
    
    def record_spend(self, amount):
        self.spent_today += amount
        self.spent_total += amount
```

### Time Constraints
```python
class TimeConstraint:
    def __init__(self, valid_from=None, valid_until=None, 
                 time_windows=None, timezone="UTC"):
        self.valid_from = valid_from or now()
        self.valid_until = valid_until
        self.time_windows = time_windows or []
        self.timezone = timezone
    
    def validate(self, when=None):
        when = when or now()
        
        # Check validity period
        if when < self.valid_from:
            return False, "Permission not yet valid"
        
        if self.valid_until and when > self.valid_until:
            return False, "Permission expired"
        
        # Check time windows (e.g., business hours only)
        if self.time_windows:
            in_window = any(
                window.contains(when) 
                for window in self.time_windows
            )
            if not in_window:
                return False, "Outside allowed time window"
        
        return True, "OK"
```

### Rate Constraints
```python
class RateConstraint:
    def __init__(self, per_second=None, per_minute=None, 
                 per_hour=None, burst_size=None):
        self.limits = {
            "second": per_second,
            "minute": per_minute,
            "hour": per_hour
        }
        self.burst_size = burst_size
        self.history = []
    
    def validate(self, when=None):
        when = when or now()
        self.cleanup_old_entries(when)
        
        for period, limit in self.limits.items():
            if limit:
                window = self.get_window(period, when)
                count = self.count_in_window(window)
                if count >= limit:
                    return False, f"Rate limit exceeded for {period}"
        
        return True, "OK"
    
    def record_use(self, when=None):
        when = when or now()
        self.history.append(when)
        self.cleanup_old_entries(when)
```

## Permission Inheritance

### Delegation Rules
```yaml
Inheritance Model:
  Root → Agent: Agent gets subset of human's permissions
  Agent → Sub-agent: Sub-agent gets subset of agent's permissions
  
  Rules:
  - Cannot grant permissions you don't have
  - Cannot extend time beyond your own expiry
  - Cannot increase value limits
  - Cannot remove constraints
```

### Implementation
```python
class PermissionInheritance:
    def delegate(self, from_entity, to_entity, requested_permissions):
        granted = []
        
        for requested in requested_permissions:
            # Check if delegator has permission
            if not from_entity.has_permission(requested):
                continue  # Skip permissions delegator doesn't have
            
            # Get delegator's permission
            delegator_perm = from_entity.get_permission(requested)
            
            # Create restricted version
            delegated = self.restrict_permission(
                delegator_perm, 
                requested
            )
            
            granted.append(delegated)
        
        return granted
    
    def restrict_permission(self, parent, requested):
        # New permission is intersection of parent and requested
        restricted = Permission(
            scope=requested.scope,
            action=requested.action
        )
        
        # Apply most restrictive constraints
        restricted.constraints = self.merge_constraints(
            parent.constraints,
            requested.constraints,
            mode="most_restrictive"
        )
        
        return restricted
```

## Approval Workflows

### Human Approval
```python
class HumanApproval:
    def __init__(self, human_lct):
        self.human_lct = human_lct
        self.pending_approvals = []
    
    def request_approval(self, agent_lct, action, context):
        request = {
            "id": generate_id(),
            "agent": agent_lct.id,
            "action": action,
            "context": context,
            "requested_at": timestamp(),
            "status": "pending"
        }
        
        self.pending_approvals.append(request)
        self.notify_human(request)
        
        return request["id"]
    
    def approve(self, request_id, constraints=None):
        request = self.find_request(request_id)
        request["status"] = "approved"
        request["approved_at"] = timestamp()
        request["constraints"] = constraints
        
        # Sign approval
        signature = self.human_lct.sign(request)
        request["signature"] = signature
        
        return request
    
    def deny(self, request_id, reason=None):
        request = self.find_request(request_id)
        request["status"] = "denied"
        request["denied_at"] = timestamp()
        request["reason"] = reason
        
        return request
```

### Automated Approval
```python
class AutomatedApproval:
    def __init__(self, rules):
        self.rules = rules
    
    def evaluate(self, request):
        for rule in self.rules:
            if rule.matches(request):
                if rule.should_approve(request):
                    return "approved", rule.constraints
                else:
                    return "denied", rule.reason
        
        # No matching rule - require human approval
        return "pending", None
```

## Security Enforcement

### Permission Checking
```python
class PermissionChecker:
    def check(self, entity, action, context):
        # 1. Find applicable permissions
        permissions = entity.get_permissions_for_action(action)
        
        if not permissions:
            return False, "No permission for action"
        
        # 2. Check each permission
        for permission in permissions:
            # Check temporal constraints
            if not permission.is_valid_now():
                continue
            
            # Check value constraints
            if not permission.check_value_limit(context.value):
                continue
            
            # Check rate limits
            if not permission.check_rate_limit():
                continue
            
            # Check other constraints
            if not permission.check_constraints(context):
                continue
            
            # Permission is valid
            return True, permission
        
        return False, "No valid permission found"
```

### Audit Logging
```python
class PermissionAudit:
    def log_grant(self, grantor, grantee, permissions):
        self.append({
            "event": "permission_grant",
            "grantor": grantor.id,
            "grantee": grantee.id,
            "permissions": [p.to_dict() for p in permissions],
            "timestamp": timestamp(),
            "witness": gather_witnesses()
        })
    
    def log_use(self, entity, permission, action, result):
        self.append({
            "event": "permission_use",
            "entity": entity.id,
            "permission": permission.to_dict(),
            "action": action,
            "result": result,
            "timestamp": timestamp()
        })
    
    def log_revoke(self, revoker, entity, permissions):
        self.append({
            "event": "permission_revoke",
            "revoker": revoker.id,
            "entity": entity.id,
            "permissions": [p.to_dict() for p in permissions],
            "timestamp": timestamp()
        })
```

## Compliance Features

### Regulatory Compliance
- KYC/AML thresholds
- Geographic restrictions
- Sanctions checking
- Reporting requirements

### Privacy Compliance
- Data minimization
- Purpose limitation
- Consent management
- Audit trails

## Implementation Checklist

- [ ] Permission parser and validator
- [ ] Constraint enforcement engine
- [ ] Inheritance calculator
- [ ] Approval workflow system
- [ ] Audit logger
- [ ] Permission UI components
- [ ] Rate limiter
- [ ] Value tracker
- [ ] Time window validator
- [ ] Revocation propagator

---

*"Permissions in ACT are not just access controls—they're the explicit, auditable expression of delegated human intent."*