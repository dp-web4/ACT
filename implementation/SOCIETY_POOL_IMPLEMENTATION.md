# Society Pool Implementation

## Overview

This document describes the implementation of Web4 Society pools in ACT, based on the foundational Web4 Society concept documented at [github.com/dp-web4/web4](https://github.com/dp-web4/web4).

## Core Principle: Tokens Belong to Society

**CRITICAL**: ATP and ADP tokens belong to the SOCIETY, not to individual entities. This is a fundamental design principle that prevents hoarding and ensures value flows to active contributors.

## Implementation in ACT

### 1. Society Treasury

```go
// Society owns all tokens
type SocietyTreasury struct {
    SocietyLCT string
    ATP        sdk.Coin  // Charged energy tokens
    ADP        sdk.Coin  // Discharged energy tokens
}
```

### 2. Token Flow Model

```
Society Treasury (Pool)
    ├── ATP (charged energy)
    │    └── Allocated to citizens for work
    │         └── Discharged to ADP through work
    └── ADP (discharged energy)
         └── Can be recharged to ATP by producers
```

### 3. Citizen Rights vs Ownership

Citizens have **RIGHTS** to request ATP allocation, not ownership:

```go
// Citizens request, don't own
type ATPAllocationRequest struct {
    CitizenLCT      string
    RequestedAmount sdk.Coin
    Purpose         string
    WorkDescription string
}
```

### 4. Energy Cycle

```go
// DischargeATP - Work converts ATP to ADP
func (k Keeper) DischargeATP(ctx sdk.Context, msg *MsgDischargeATP) error {
    // 1. Check citizen has allocation rights
    // 2. Verify society has sufficient ATP
    // 3. Convert society ATP to ADP
    // 4. Record work on ledger
    // 5. Update trust tensors
}

// RechargeADP - Producers add energy
func (k Keeper) RechargeADP(ctx sdk.Context, msg *MsgRechargeADP) error {
    // 1. Verify producer credentials
    // 2. Convert society ADP to ATP
    // 3. Record energy addition
    // 4. Update society metrics
}
```

### 5. Implementation Status

#### Completed ✅
- Basic proto messages for DischargeATP and RechargeADP
- Keeper methods that operate on society pool
- LCT minting for entity creation

#### In Progress 🚧
- Society treasury initialization
- Allocation rights management
- Work verification mechanisms

#### Planned 📋
- Fractal society relationships
- Cross-society token flows
- Demurrage mechanisms

## Key Differences from Traditional Tokens

| Traditional | Web4 Society Pool |
|------------|------------------|
| Individual balances | Society treasury |
| Transfer between accounts | Allocation rights |
| Permanent ownership | Temporary stewardship |
| Accumulation possible | Anti-hoarding built-in |
| Value extraction | Value circulation |

## Example: Agent Swarm Society

```json
{
  "society_lct": "lct-society-act-swarm-001",
  "citizens": [
    "lct-agent-alice",
    "lct-agent-bob",
    "lct-human-dennis"
  ],
  "treasury": {
    "ATP": "1000000uatp",
    "ADP": "500000uadp"
  },
  "laws": [
    "all_allocations_require_work_proof",
    "maximum_allocation_100_atp_per_cycle",
    "unanimous_decisions_required"
  ],
  "ledger_type": "confined"
}
```

## Migration from Individual Balances

For existing implementations assuming individual token ownership:

1. **Create Society Treasury**: Initialize society with total token supply
2. **Convert Balances to Rights**: Map existing balances to allocation rights
3. **Update Transfer Logic**: Replace transfers with allocation requests
4. **Add Work Tracking**: All ATP discharge must specify work performed
5. **Implement Producers**: Define who can recharge ADP to ATP

## Energy Investment Formula

The key insight about energy economics:

```
ATP_in - ADP_out = Confirmed_Value + Energy_Invested
```

Where `Energy_Invested` represents learning costs, not waste. This formula captures how societies invest in capability development.

## References

- [Web4 Society Specification](https://github.com/dp-web4/web4/blob/main/web4-standard/core-spec/SOCIETY_SPECIFICATION.md)
- [ATP/ADP Cycle Specification](https://github.com/dp-web4/web4/blob/main/web4-standard/core-spec/atp-adp-cycle.md)
- [Society Integration Summary](https://github.com/dp-web4/web4/blob/main/web4-standard/SOCIETY_INTEGRATION_SUMMARY.md)

---

*"Tokens belong to the society, not individuals. This isn't a limitation - it's liberation from the tyranny of accumulation."*