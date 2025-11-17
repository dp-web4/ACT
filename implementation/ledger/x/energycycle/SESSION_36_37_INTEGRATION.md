# Session #36-37: Energy-Backed ATP Integration into ACT

**Date:** 2025-11-16
**Status:** In Progress - Protobuf Definitions Complete
**Goal:** Integrate Session #36's Python energy-backed ATP implementation into ACT blockchain (Cosmos SDK)

---

## Overview

Session #36 implemented a complete energy-backed ATP system in Python (~3,850 lines):
- Energy capacity proofs (solar, compute, grid, human, battery)
- ChargedATP with expiration (prevents hoarding)
- Work allocation and discharge (ATP → ADP)
- Trust-based priority queues (trust affects priority, not price)
- Energy-backed identity bonds (capacity commitment, not currency stake)
- 18/18 integration tests passing

Session #37 is integrating these innovations into the ACT blockchain.

---

## Architecture Changes

### Before (Current ACT)
- ATP can be created without energy proof
- ATP expires but no cleanup mechanism
- No trust-based priority system
- No identity bond system
- Trust affects "efficiency rating" but not queue position

### After (Session #36 Integration)
- ATP requires energy capacity proof
- Automatic expired ATP cleanup
- Trust-based priority queue (everyone served, but order varies)
- Energy-backed identity bonds (thermodynamic Sybil resistance)
- Trust affects queue priority, not cost

---

## New Protobuf Definitions

### 1. Energy Capacity Proofs
**File:** `proto/racecarweb/energycycle/v1/energy_capacity_proof.proto`

**Key Messages:**
- `EnergyCapacityProof` - Generic proof interface
- `SolarPanelDetails` - Solar panel specific data
- `ComputeResourceDetails` - GPU/CPU/TPU data
- `GridConnectionDetails` - Utility grid allocation
- `HumanLaborDetails` - Human labor capacity
- `BatteryStorageDetails` - Battery capacity
- `EnergyCapacityRegistry` - Per-society energy sources

**Purpose:** Prove that ATP is backed by real energy, not created from nothing.

### 2. Energy-Backed Transactions
**File:** `proto/racecarweb/energycycle/v1/tx_energy_backed.proto`

**New Transaction Types:**
1. `MsgRegisterEnergyCapacityProof` - Register energy source
2. `MsgChargeATPWithEnergyProof` - Charge ATP with proof
3. `MsgAllocateATPToWork` - Allocate ATP to work ticket
4. `MsgCompleteWork` - Discharge ATP → ADP
5. `MsgSubmitWorkRequest` - Submit to priority queue
6. `MsgProcessWorkQueue` - Process next request
7. `MsgRegisterIdentityBond` - Register capacity bond
8. `MsgValidateIdentityBond` - Validate capacity maintained
9. `MsgRegisterVouch` - Vouch for newcomer
10. `MsgCompleteVouch` - Complete vouching period
11. `MsgCleanupExpiredATP` - Remove expired ATP

---

## Key Innovations from Session #36

### 1. Thermodynamic ATP Backing

**Python (Session #36):**
```python
def charge_atp(amount, energy_source_id):
    # Find & validate energy source
    source = registry.find_source(energy_source_id)
    if not validator.validate_proof(source):
        return None

    # Check capacity constraint
    if amount > available_capacity:
        return None

    # Create ChargedATP with expiration
    return ChargedATP(amount=amount, expiration=now + 7_days)
```

**Cosmos SDK (Session #37):**
```protobuf
message MsgChargeATPWithEnergyProof {
  string society_lct = 1;
  string amount = 2;
  string energy_source_identifier = 3;
  int64 lifetime_days = 4;
}
```

**Keeper Implementation (Pending):**
```go
func (k Keeper) ChargeATPWithEnergyProof(ctx sdk.Context, msg *types.MsgChargeATPWithEnergyProof) error {
    // 1. Validate energy source exists
    // 2. Check capacity constraint
    // 3. Create ChargedATP with expiration block
    // 4. Deduct from ADP pool
    // 5. Return ChargedATP
}
```

### 2. Trust-Based Priority (Not Price)

**Session #33 (Old):**
```python
# Trust affects cost
rate = base_rate * trust_multiplier  # 1.0x to 2.0x
```

**Session #36 (New):**
```python
# Trust affects priority
priority = calculate_priority(trust_score)  # HIGH to DEFERRED
```

**Cosmos SDK:**
```protobuf
enum RequestPriority {
  REQUEST_PRIORITY_HIGH = 2;        // Trust ≥ 0.9
  REQUEST_PRIORITY_ELEVATED = 3;    // Trust ≥ 0.7
  REQUEST_PRIORITY_NORMAL = 4;      // Trust ≥ 0.5
  REQUEST_PRIORITY_LOW = 5;         // Trust ≥ 0.3
  REQUEST_PRIORITY_DEFERRED = 6;    // Trust < 0.3
}
```

### 3. Energy-Backed Identity Bonds

**Session #34 (Currency Model):**
```python
class IdentityBond:
    bond_amount: int = 1000  # ATP staked

    def forfeit_bond(self) -> int:
        return self.bond_amount  # Money seized
```

**Session #36 (Energy Model):**
```python
class EnergyBackedIdentityBond:
    committed_capacity_watts: float
    reputation_at_risk: float = 0.5

    def handle_violation(self) -> float:
        return self.reputation_at_risk  # Reputation penalty
```

**Cosmos SDK:**
```protobuf
message EnergyBackedIdentityBond {
  string committed_capacity_watts = 3;
  repeated string energy_source_ids = 4;
  string reputation_at_risk = 7;  // 0.5 = 50% penalty
  BondStatus status = 6;
}
```

---

## Implementation Plan

### Phase 1: Protobuf & Code Generation ✅

**Status:** Complete

**Deliverables:**
- [x] `energy_capacity_proof.proto` (300 lines)
- [x] `tx_energy_backed.proto` (250 lines)
- [ ] Generate Go code: `make proto-gen`

### Phase 2: Keeper Methods

**File:** `x/energycycle/keeper/energy_backed_atp.go`

**Methods to Implement:**
```go
// Energy Capacity
func (k Keeper) RegisterEnergyCapacityProof(ctx, msg) error
func (k Keeper) ValidateEnergyCapacityProof(ctx, proofID) (bool, error)
func (k Keeper) GetTotalEnergyCapacity(ctx, societyLCT) (sdk.Dec, error)

// ATP Charging
func (k Keeper) ChargeATPWithEnergyProof(ctx, msg) (*ChargedATP, error)
func (k Keeper) CleanupExpiredATP(ctx, societyLCT) (int, sdk.Dec, error)

// Work Allocation
func (k Keeper) AllocateATPToWork(ctx, msg) (*WorkTicket, error)
func (k Keeper) CompleteWork(ctx, msg) (*DischargedADP, error)
func (k Keeper) FailWork(ctx, ticketID) error

// Priority Queue
func (k Keeper) SubmitWorkRequest(ctx, msg) (*WorkRequest, error)
func (k Keeper) ProcessWorkQueue(ctx, societyLCT) (*WorkRequest, *WorkTicket, error)
func (k Keeper) CalculatePriority(ctx, requesterLCT) (RequestPriority, sdk.Dec, error)

// Identity Bonds
func (k Keeper) RegisterIdentityBond(ctx, msg) (*EnergyBackedIdentityBond, error)
func (k Keeper) ValidateIdentityBond(ctx, societyLCT) (bool, string, sdk.Dec, error)

// Vouching
func (k Keeper) RegisterVouch(ctx, msg) (*EnergyBackedVouch, error)
func (k Keeper) CompleteVouch(ctx, vouchID, established bool) (sdk.Dec, error)
```

### Phase 3: Message Handlers

**File:** `x/energycycle/keeper/msg_server_energy_backed.go`

**Handlers for each transaction type:**
```go
func (ms msgServer) RegisterEnergyCapacityProof(ctx, msg) (*MsgRegisterEnergyCapacityProofResponse, error)
func (ms msgServer) ChargeATPWithEnergyProof(ctx, msg) (*MsgChargeATPWithEnergyProofResponse, error)
func (ms msgServer) AllocateATPToWork(ctx, msg) (*MsgAllocateATPToWorkResponse, error)
// ... etc for all 11 message types
```

### Phase 4: Integration with Existing Modules

**Trust Tensor Integration:**
```go
// In CalculatePriority:
trustScore, _, err := k.trusttensorKeeper.CalculateRelationshipTrust(ctx, lctID, "work_request")
priority := calculatePriorityFromTrust(trustScore)
```

**LCT Manager Integration:**
```go
// Validate LCT exists and is active
lct, found := k.lctmanagerKeeper.GetLinkedContextToken(ctx, lctID)
if !found || lct.PairingStatus != "active" {
    return ErrInvalidLCT
}
```

### Phase 5: Storage Schema

**Collections to Add:**
```go
// Energy Capacity
EnergyCapacityProofs collections.Map[string, EnergyCapacityProof]
EnergyCapacityRegistries collections.Map[string, EnergyCapacityRegistry]

// Charged ATP
ChargedATPTokens collections.Map[string, ChargedATP]
WorkTickets collections.Map[string, WorkTicket]

// Priority Queue
WorkRequests collections.Map[string, WorkRequest]

// Bonds & Vouching
IdentityBonds collections.Map[string, EnergyBackedIdentityBond]
Vouches collections.Map[string, EnergyBackedVouch]
```

### Phase 6: Testing

**Integration Tests:**
1. Energy capacity registration and validation
2. ATP charging with proof validation
3. ATP expiration and cleanup
4. Work allocation and discharge
5. Trust-based priority queue ordering
6. Identity bond creation and validation
7. Vouching workflow
8. Full cycle: Energy → ATP → Work → ADP → Reputation

---

## Compatibility with Existing ACT Code

### ATP Expiration
**Current:** `expiration_block` field exists in `RelationshipAtpToken`
**New:** `ChargedATP` also has expiration, but with explicit cleanup mechanism

**Integration Strategy:** Migrate existing ATP tokens to ChargedATP format.

### Trust Integration
**Current:** Trust affects "efficiency_rating"
**New:** Trust affects priority queue position

**Integration Strategy:** Keep efficiency_rating for discharge calculations, add priority for queue.

### Society Pools
**Current:** `SocietyPool` tracks ATP/ADP balances
**New:** Energy-backed pool requires energy source validation

**Integration Strategy:** Extend SocietyPool with EnergyCapacityRegistry reference.

---

## Security Considerations

### 1. Energy Proof Validation

**Attack:** Submit fake energy capacity proof

**Defense:**
- Validators must verify proofs (not just accept claims)
- Cross-reference with external oracles (solar APIs, GPU queries)
- Require proof from multiple validators for large capacities
- Slash validators who accept fraudulent proofs

**Implementation:**
```go
func (k Keeper) ValidateEnergyCapacityProof(ctx, proof) (bool, error) {
    // 1. Check proof signature
    // 2. Query external oracle (if applicable)
    // 3. Require threshold of validator confirmations
    // 4. Store verification record
}
```

### 2. Sybil Resistance via Thermodynamics

**Attack:** Create many Sybil identities

**Defense (Session #36):**
- Each identity requires energy capacity proof
- Can't fake physics (must have real energy)
- Identity bonds commit capacity for 30 days
- Violation forfeits reputation (hard to rebuild)

**Validation:**
```go
func (k Keeper) RegisterIdentityBond(ctx, msg) error {
    // 1. Verify energy sources exist and are valid
    // 2. Check total capacity meets minimum (e.g., 500W)
    // 3. Lock sources for lock_period_days
    // 4. Create bond with reputation at risk
}
```

### 3. Priority Queue Gaming

**Attack:** Submit many low-priority requests to clog queue

**Defense:**
- Everyone served (no exclusion), just ordered by priority
- Low-trust entities wait longer but still processed
- Can rate-limit requests per LCT per block
- Economic cost (reputation risk) for spamming

**Implementation:**
```go
func (k Keeper) SubmitWorkRequest(ctx, msg) error {
    // 1. Check rate limit (e.g., 10 requests per block per LCT)
    // 2. Calculate priority from trust score
    // 3. Add to priority queue
    // 4. Emit event with queue position
}
```

---

## Performance Optimizations

### 1. Indexing

**Critical Indexes:**
```go
// Index ATP by expiration block for cleanup
collections.NewIndexedMap(
    sb, collections.NewPrefix("charged_atp"),
    ChargedATPCodec,
    collections.NamedIndex(
        "expiration",
        collections.NewMultiIndex(
            sb, ExpirationPrefix,
            func(k string, v ChargedATP) (int64, error) {
                return v.Expiration.Unix(), nil
            },
        ),
    ),
)

// Index work requests by priority
collections.NewIndexedMap(
    sb, collections.NewPrefix("work_requests"),
    WorkRequestCodec,
    collections.NamedIndex(
        "priority",
        collections.NewMultiIndex(
            sb, PriorityPrefix,
            func(k string, v WorkRequest) (int32, error) {
                return int32(v.Priority), nil
            },
        ),
    ),
)
```

### 2. Batch Operations

**Cleanup Expired ATP:**
```go
func (k Keeper) CleanupExpiredATP(ctx sdk.Context) error {
    // Process in batches of 100
    currentBlock := ctx.BlockHeight()
    iter := k.ChargedATPTokens.Iterate(ctx, nil)

    batch := []string{}
    for ; iter.Valid(); iter.Next() {
        atp, _ := iter.Value()
        if atp.ExpirationBlock < currentBlock {
            batch = append(batch, atp.Id)
        }

        if len(batch) >= 100 {
            k.processBatchCleanup(ctx, batch)
            batch = []string{}
        }
    }

    if len(batch) > 0 {
        k.processBatchCleanup(ctx, batch)
    }
}
```

### 3. Caching

**Trust Score Cache:**
```go
// Cache trust scores for 10 blocks (~1 minute)
trustScoreCache := collections.NewMap(
    sb, TrustCachePrefix,
    StringKey, LegacyDecValue,
)
```

---

## Migration Path

### Step 1: Deploy with Opt-In

- New energy-backed operations available
- Existing ATP/ADP operations continue working
- Societies can migrate at their own pace

### Step 2: Incentivize Migration

- Energy-backed ATP gets priority in queue
- Energy-backed bonds get trust boost
- After block X, old ATP starts depreciating

### Step 3: Full Cutover

- After block Y, all ATP must be energy-backed
- Migrate remaining balances automatically
- Deprecate old CreateAtpToken method

---

## Next Steps

### Immediate (Session #37)

1. **Generate protobuf code:**
   ```bash
   cd implementation/ledger
   make proto-gen
   ```

2. **Implement keeper methods:**
   - Start with energy capacity registration
   - Then ATP charging with proofs
   - Then work allocation/discharge

3. **Create basic tests:**
   - Register solar panel proof
   - Charge ATP backed by solar
   - Allocate to work
   - Complete work (discharge to ADP)

### Medium Term

4. **Integrate with trust tensor:**
   - Use existing trust scores for priority
   - Update trust from work quality

5. **Add oracle integration:**
   - Query external APIs for energy validation
   - Real solar panel verification
   - GPU hardware queries

6. **Performance testing:**
   - Benchmark priority queue with 10k requests
   - Test expiration cleanup with 100k ATP tokens
   - Stress test identity bond validation

### Long Term

7. **SAGE integration:**
   - Jetson GPU as energy source
   - Real-time capacity monitoring
   - Automatic ATP charging from compute work

8. **Cross-chain energy proofs:**
   - Prove energy capacity on other chains
   - Bridge ATP across chains
   - Unified energy economy

---

## References

- **Session #36 Summary:** `/home/dp/ai-workspace/private-context/moments/2025-11-16-legion-autonomous-web4-session-36.md`
- **Session #36 Design:** `/home/dp/ai-workspace/web4/web4-standard/implementation/act_deployment/ENERGY_BACKED_ATP_DESIGN.md`
- **Python Implementation:** `/home/dp/ai-workspace/web4/web4-standard/implementation/act_deployment/energy_backed_atp.py`
- **Integration Tests:** `/home/dp/ai-workspace/web4/web4-standard/implementation/act_deployment/test_energy_backed_atp_integration.py`

---

**Session #37 Status:** Protobuf definitions complete. Next: Generate code and implement keeper methods.
