# Web4-Compliant ACT Swarm Architecture
## Role-Based Agent System with LCT Identity

*Every agent is a Role entity with its own LCT, performed by Claude through the code interface*

---

## Core Principle: Roles as First-Class Entities

In our Web4-compliant swarm:
- Each agent type is a **Role entity** with its own LCT
- Claude (via code interface) **pairs** with Role LCTs to perform them
- Each role has **R6 rules** defining its behavior
- All actions are **witnessed** and create MRH relationships
- Every operation costs **ATP** and generates **ADP** proofs

---

## Role LCT Hierarchy

### 🌍 Society: ACT Development Collective
```json
{
  "lct_id": "lct:web4:society:act-dev-collective",
  "entity_type": "society",
  "law_oracle": "lct:web4:oracle:act-governance",
  "citizens": ["all role LCTs below"],
  "atp_pool": 10000,
  "genesis_block": "block:act-genesis"
}
```

### 👑 Genesis Role: Meta-Orchestrator
```json
{
  "lct_id": "lct:web4:role:act-genesis-orchestrator",
  "entity_type": "role",
  "mode": "delegative",
  "birth_certificate": {
    "society": "lct:web4:society:act-dev-collective",
    "citizen_role": "lct:web4:role:citizen:act-dev",
    "initial_rights": ["spawn_roles", "allocate_atp", "coordinate_swarm"],
    "initial_responsibilities": ["deliver_act", "maintain_coherence", "witness_progress"]
  },
  "r6_rules": {
    "Rules": ["Must coordinate all domain queens", "Weekly evolution cycles"],
    "Roles": ["Meta-orchestrator", "Swarm coordinator"],
    "Request": ["Build complete ACT platform"],
    "Reference": ["Web4 standard", "ACT specification"],
    "Resource": ["1000 ATP daily", "All swarm agents"],
    "Result": ["Functioning ACT platform in 4 weeks"]
  },
  "mrh": {
    "bound": ["lct:web4:society:act-dev-collective"],
    "paired": ["lct:web4:ai:claude-instance"],
    "witnessing": ["lct:web4:role:witness:act-progress"]
  }
}
```

---

## Domain Queen Roles (Level 1)

### LCT Infrastructure Queen Role
```json
{
  "lct_id": "lct:web4:role:lct-infrastructure-queen",
  "entity_type": "role",
  "mode": "delegative",
  "citizen_of": "lct:web4:society:act-dev-collective",
  "r6_rules": {
    "Rules": [
      "Implement Ed25519 cryptographic binding",
      "Create LCT lifecycle management",
      "Ensure unforgeable presence"
    ],
    "Roles": ["System architect", "Infrastructure lead"],
    "Request": ["Build LCT infrastructure"],
    "Reference": ["web4-standard/core-spec/binding.md"],
    "Resource": ["100 ATP daily", "4 worker roles"],
    "Result": ["Complete LCT implementation with registry"]
  },
  "permissions": [
    "capability:create_lct",
    "capability:bind_entities",
    "capability:manage_registry"
  ],
  "worker_roles": [
    "lct:web4:role:lct-coder",
    "lct:web4:role:registry-developer",
    "lct:web4:role:crypto-auditor",
    "lct:web4:role:lct-tester"
  ]
}
```

### ACP Protocol Queen Role
```json
{
  "lct_id": "lct:web4:role:acp-protocol-queen",
  "entity_type": "role",
  "mode": "delegative",
  "citizen_of": "lct:web4:society:act-dev-collective",
  "r6_rules": {
    "Rules": [
      "Define agent plan structure",
      "Implement intent routing",
      "Create decision collection"
    ],
    "Roles": ["Protocol architect", "ACP lead"],
    "Request": ["Implement Agentic Context Protocol"],
    "Reference": ["ACT/core-spec/acp-protocol.md"],
    "Resource": ["100 ATP daily", "4 worker roles"],
    "Result": ["Functioning ACP engine"]
  }
}
```

---

## Worker Roles (Level 2)

### Example: LCT Coder Role
```json
{
  "lct_id": "lct:web4:role:lct-coder",
  "entity_type": "role",
  "mode": "responsive",
  "citizen_of": "lct:web4:society:act-dev-collective",
  "r6_rules": {
    "Rules": [
      "Write clean, tested code",
      "Follow Web4 specifications",
      "Document all functions"
    ],
    "Roles": ["Implementation specialist"],
    "Request": ["Code LCT components"],
    "Reference": ["TypeScript", "Ed25519", "Web4 spec"],
    "Resource": ["5 ATP per task", "Development tools"],
    "Result": ["Working LCT implementation"]
  },
  "capabilities": [
    "code:typescript",
    "crypto:ed25519",
    "test:jest"
  ],
  "reports_to": "lct:web4:role:lct-infrastructure-queen"
}
```

---

## R6 Script Templates

### Queen Role Script
```javascript
// R6 Implementation for Queen Roles
class QueenRole {
  constructor(lct, r6Rules) {
    this.lct = lct;
    this.rules = r6Rules;
    this.atpBudget = r6Rules.Resource.atp;
    this.workers = [];
  }
  
  async execute(request) {
    // Rules: Validate request against role rules
    if (!this.validateRules(request)) {
      throw new Error("Request violates role rules");
    }
    
    // Roles: Confirm role authority
    if (!this.hasAuthority(request.domain)) {
      throw new Error("Outside role authority");
    }
    
    // Request: Parse and plan
    const plan = await this.createPlan(request);
    
    // Reference: Load necessary knowledge
    const references = await this.loadReferences(plan);
    
    // Resource: Allocate ATP to workers
    const allocation = this.allocateResources(plan);
    
    // Execute with workers
    const results = await this.coordinateWorkers(plan, allocation);
    
    // Result: Generate ADP proof
    return this.generateADP(results);
  }
  
  async witnessAction(action) {
    // All actions are witnessed
    return {
      role_lct: this.lct.id,
      action: action,
      timestamp: Date.now(),
      atp_cost: this.calculateCost(action),
      witness_signature: await this.sign(action)
    };
  }
}
```

### Worker Role Script
```javascript
// R6 Implementation for Worker Roles
class WorkerRole {
  constructor(lct, r6Rules) {
    this.lct = lct;
    this.rules = r6Rules;
    this.capabilities = r6Rules.capabilities;
  }
  
  async performTask(task) {
    // Rules: Check task compliance
    this.enforceRules(task);
    
    // Roles: Confirm within role scope
    this.validateScope(task);
    
    // Request: Accept task from queen
    const work = this.parseTask(task);
    
    // Reference: Load domain knowledge
    const knowledge = await this.loadKnowledge(work);
    
    // Resource: Use allocated ATP
    const resources = await this.claimResources(work.atp);
    
    // Execute the work
    const output = await this.execute(work, knowledge, resources);
    
    // Result: Return with witness proof
    return {
      output: output,
      witness: await this.createWitness(output),
      adp: await this.generateADP(output, work.atp)
    };
  }
}
```

---

## MRH (Markov Relevancy Horizon) Boundaries

### Genesis Orchestrator MRH
```
Bound to: Society
Paired with: Claude instance
Witnessing: All queen roles, progress tracker
Broadcast: Status updates to all entities
```

### Queen Role MRH
```
Bound to: Genesis orchestrator
Paired with: Claude instance (when active)
Witnessing: Worker roles, peer queens
Broadcast: Domain updates
```

### Worker Role MRH
```
Bound to: Queen role
Paired with: Claude instance (during tasks)
Witnessing: Task outputs, peer workers
Broadcast: Completion signals
```

---

## ATP/ADP Economic Flow

### ATP Allocation Chain
```
Society Treasury (10,000 ATP)
    ↓ 1000 ATP/day
Genesis Orchestrator
    ↓ 100 ATP/day each
Domain Queens (6 × 100 = 600 ATP)
    ↓ 5 ATP/task
Worker Roles (24 × ~20 tasks = 480 ATP)
```

### ADP Proof Generation
Every completed task generates ADP proving:
1. Work was requested (Request)
2. Work followed rules (Rules)
3. Correct role performed it (Roles)
4. References were consulted (Reference)
5. Resources were consumed (Resource)
6. Result was delivered (Result)

---

## Witness Network

### Witness Roles
```json
{
  "lct_id": "lct:web4:role:witness:act-development",
  "entity_type": "role",
  "mode": "responsive",
  "responsibilities": [
    "Record all role actions",
    "Verify ATP consumption",
    "Validate ADP proofs",
    "Maintain immutable ledger"
  ]
}
```

### Witnessing Protocol
1. Every role action creates a witness record
2. Witnesses verify action compliance with R6
3. Successful actions increase role reputation
4. Failed actions decrease trust scores
5. All witnesses stored in immutable ledger

---

## Dictionary Entity for Swarm Coordination

```json
{
  "lct_id": "lct:web4:dictionary:swarm-coordinator",
  "entity_type": "dictionary",
  "mode": "responsive/agentic",
  "purpose": "Translate between Claude-Flow and Web4 protocols",
  "mappings": {
    "claude-flow-agent": "web4-role-entity",
    "swarm-memory": "mrh-graph",
    "task-execution": "atp-consumption",
    "completion": "adp-generation"
  }
}
```

---

## Compliance Verification

✅ **Entity Types**: All agents are Role entities (Section 2.1)
✅ **Behavioral Modes**: Queens are Delegative, Workers are Responsive
✅ **Citizen Roles**: All roles start as citizens of ACT Dev Society
✅ **Role LCTs**: Each role has full LCT with R6 rules
✅ **MRH Compliance**: Proper binding, pairing, witnessing, broadcast
✅ **ATP/ADP**: Economic model with proper token flow
✅ **Witness Network**: All actions recorded and witnessed
✅ **Dictionary Entity**: Translation layer for protocol bridging

---

## Launch Sequence (Web4 Compliant)

```bash
# 1. Create ACT Development Society
act society create \
  --name="ACT Development Collective" \
  --law-oracle="act-governance" \
  --atp-treasury=10000

# 2. Birth Genesis Orchestrator Role
act role create \
  --type="genesis-orchestrator" \
  --society="act-dev-collective" \
  --mode="delegative" \
  --atp-budget=1000

# 3. Claude pairs with Genesis Role
act pair \
  --entity="claude-instance" \
  --role="genesis-orchestrator" \
  --witness="act-development-witness"

# 4. Genesis spawns Queen Roles
act role spawn-queens \
  --count=6 \
  --mode="delegative" \
  --budget=100

# 5. Queens spawn Worker Roles
act role spawn-workers \
  --per-queen=4 \
  --mode="responsive" \
  --budget=5

# 6. Begin witnessed execution
act swarm execute \
  --witness-all \
  --generate-adp \
  --track-atp
```

---

*"In Web4, every role has presence, every action has witness, every result has proof."*