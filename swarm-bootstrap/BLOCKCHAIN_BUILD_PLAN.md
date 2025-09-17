# ACT Blockchain Build Plan
## Mission: Get Web4 Blockchain Running with LCT and ATP/ADP

### Primary Deliverables
1. **Mint and bind LCTs** - Create unforgeable digital identities
2. **ATP/ADP semifungible tokens** - Energy-based economy  
3. **Discharge/Recharge mechanisms** - ATP ↔ ADP conversion
4. **T3/V3 attributions** - Trust tensor updates on LCTs

### Phase 1: Proto Generation & Compilation (2 hours)
**Queen**: Protocol Queen
**Workers**: Proto Generator, Type Validator, Import Fixer

Tasks:
1. Run `make proto-gen` to generate Go types from protos
2. Fix import paths in generated code
3. Resolve type mismatches between modules
4. Ensure all proto types compile

Expected Issues:
- Missing gogoproto imports
- Incorrect package paths  
- Type registration conflicts

### Phase 2: Module Wiring (3 hours)
**Queen**: Infrastructure Queen
**Workers**: Module Wirer, Codec Registrar, Genesis Builder

Tasks:
1. Update app/app.go to include all 5 Web4 modules
2. Register message types and codecs
3. Add store keys for each module
4. Configure module manager
5. Set up module routing

Key Integration Points:
```go
// app.go additions needed
- lctmanagerKeeper
- trusttensorKeeper  
- energycycleKeeper
- componentregistryKeeper
- pairingqueueKeeper
```

### Phase 3: LCT Implementation (4 hours)
**Queen**: Identity Queen  
**Workers**: Crypto Specialist, Binding Agent, Genesis Minter

Core Functions:
```go
// Mint new LCT
MintLCT(ctx, entityType, ed25519PubKey) (lctID, error)

// Bind LCT to entity (permanent)
BindLCT(ctx, lctID, entityID, bindingProof) error

// Get LCT with full MRH graph
GetLCTWithMRH(ctx, lctID) (*LCT, *MRH, error)
```

Birth Certificate Structure:
- Genesis timestamp
- Parent LCT (if any)
- Entity type declaration
- Initial T3/V3 tensors

### Phase 4: ATP/ADP Token System (5 hours)
**Queen**: Economy Queen
**Workers**: Token Minter, Pool Manager, Velocity Tracker

Core Mechanics:
```go
// Mint initial ATP/ADP pairs
MintATPADP(ctx, societyPool, amount) error

// Discharge ATP to ADP (R6 action)
DischargeATP(ctx, fromLCT, amount, r6Action) (*ADP, error)

// Recharge ADP to ATP (productive work)
RechargeADP(ctx, toLCT, adpToken, workProof) (*ATP, error)

// Track velocity and apply demurrage
ApplyDemurrage(ctx, pool) error
```

Token Properties:
- Semifungible (tracked by charge state)
- Pool-managed (no individual wallets)
- Velocity requirements (must flow)
- Demurrage on idle tokens

### Phase 5: T3/V3 Attribution System (3 hours)
**Queen**: Trust Queen
**Workers**: Tensor Calculator, Attribution Engine, Reputation Updater

Trust Updates:
```go
// Update T3 tensor after action
UpdateT3(ctx, lctID, role, dimension, delta) error

// Update V3 tensor from outcome
UpdateV3(ctx, lctID, context, outcome, witnesses) error

// Calculate trust distance
GetTrustDistance(ctx, fromLCT, toLCT, role) (float64, error)

// Apply trust gravity (attraction/repulsion)
ApplyTrustGravity(ctx, lctID) error
```

Attribution Rules:
- Positive outcomes increase relevant tensor dimensions
- Witnessed actions have higher impact
- Trust decays with distance
- Role-specific reputation

### Phase 6: Genesis Configuration (2 hours)
**Queen**: Genesis Queen
**Workers**: State Builder, Validator Setup, Chain Initializer

Genesis State:
```json
{
  "lctmanager": {
    "params": {},
    "genesis_lcts": [
      {
        "id": "genesis-orchestrator",
        "entity_type": "GENESIS",
        "t3_tensor": [1.0, 1.0, 1.0],
        "v3_tensor": [1.0, 1.0, 1.0]
      }
    ]
  },
  "energycycle": {
    "atp_pool": 10000,
    "adp_pool": 0,
    "velocity_requirement": 0.1,
    "demurrage_rate": 0.001
  }
}
```

### Phase 7: Transaction Testing (2 hours)
**Queen**: Testing Queen
**Workers**: Transaction Builder, Result Validator, Witness Logger

Test Scenarios:
1. **Mint LCT**: Create new digital identity
2. **Bind LCT**: Permanently attach to entity
3. **Transfer ATP**: Discharge for action
4. **Recharge ADP**: Convert through work
5. **Update Trust**: Modify T3/V3 tensors

Success Criteria:
- Chain starts without errors
- Transactions execute and are witnessed
- State changes persist across blocks
- Query endpoints return correct data

### Swarm Configuration

**ATP Budget**: 2000 ATP (significant but justified)
**Timeline**: 21 hours of focused work
**Parallelization**: Queens can work simultaneously after Phase 2
**Witness Mode**: Full activity logging for audit

### Critical Success Factors
1. Proto generation must succeed first
2. Module wiring enables everything else
3. Genesis state must be valid
4. Each module must expose gRPC query endpoints
5. CLI commands must work for testing

### Risk Mitigation
- Start with minimal features, add complexity later
- Use existing Cosmos SDK patterns
- Test each module in isolation first
- Keep witness logs for debugging

### Expected Outcome
A running blockchain that can:
- Create and manage LCT identities
- Process ATP/ADP energy transactions
- Update trust tensors based on actions
- Provide a foundation for Web4 society

This is the blueprint for making Web4 real - from concept to running chain!