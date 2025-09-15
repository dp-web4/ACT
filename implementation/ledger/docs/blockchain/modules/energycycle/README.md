# Energy Cycle Module

## Table of Contents

1. [Overview](#overview)
2. [Key Concepts](#key-concepts)
3. [State Management](#state-management)
4. [Messages (Transactions)](#messages-transactions)
5. [Queries](#queries)
6. [Events](#events)
7. [Parameters](#parameters)
8. [Energy Accounting](#energy-accounting)
9. [Integration Guide](#integration-guide)

## Overview

The Energy Cycle module is the core engine for energy operations within the Web4-ModBatt system. It manages energy transfers between components, tracks energy balances, validates relationship values, and maintains comprehensive energy flow history. The module integrates closely with Trust Tensor for trust-based validation and LCT Manager for relationship verification.

### Purpose
- Execute secure energy transfers between paired components
- Track energy balances for all relationships (LCTs)
- Maintain comprehensive energy flow history
- Validate relationship values using ATP/ADP tokens
- Implement trust-based energy operation limits

### Dependencies
- **LCT Manager**: Validates relationships and access permissions
- **Trust Tensor**: Provides trust scores for operation validation
- **Component Registry**: Verifies component identities

### Module Store Key
`energycycle`

## Key Concepts

### Energy Operation
A discrete energy transaction between components:
- **Operation ID**: Unique identifier for the transaction
- **LCT Context**: The relationship enabling the transfer
- **Energy Amount**: Quantity in watt-hours (Wh)
- **Transfer Type**: DIRECT, POOLED, RESERVED, EMERGENCY
- **Validation Status**: Security and trust checks

### ATP/ADP Tokens
Attention-based tokens for relationship prioritization:
- **ATP (Attention Transfer Protocol)**: Direct attention between two components
- **ADP (Attention Distribution Protocol)**: Distributed attention across multiple relationships
- Used for operation prioritization and resource allocation

### Relationship Energy Balance
Track energy flow within each LCT relationship:
- **Total Energy In**: Energy received by relationship
- **Total Energy Out**: Energy distributed from relationship
- **Current Balance**: Net energy available
- **Reserved Energy**: Energy allocated but not yet transferred

### Energy Flow History
Immutable log of all energy operations:
- **Operation Records**: Complete transaction details
- **Balance Snapshots**: State at specific points in time
- **Performance Metrics**: Efficiency and reliability data

## State Management

### Stored Types

#### 1. EnergyOperation
```protobuf
message EnergyOperation {
  string operation_id = 1;
  string lct_id = 2;
  string source_component = 3;
  string target_component = 4;
  string energy_amount_wh = 5;  // Using string for precise decimals
  string transfer_rate_w = 6;   // Power in watts
  string operation_type = 7;    // DIRECT, POOLED, RESERVED, EMERGENCY
  string status = 8;            // PENDING, VALIDATED, EXECUTING, COMPLETED, FAILED
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp executed_at = 10;
  string trust_score_required = 11;
  string actual_trust_score = 12;
  map<string, string> metadata = 13;
}
```

#### 2. RelationshipATPToken
```protobuf
message RelationshipATPToken {
  string lct_id = 1;
  string from_component = 2;
  string to_component = 3;
  string allocated_attention = 4;  // Attention units allocated
  string consumed_attention = 5;   // Attention units used
  google.protobuf.Timestamp allocation_time = 6;
  google.protobuf.Timestamp expiry_time = 7;
  string priority_level = 8;      // LOW, NORMAL, HIGH, CRITICAL
}
```

#### 3. RelationshipADPToken
```protobuf
message RelationshipADPToken {
  string lct_id = 1;
  string controller_component = 2;
  repeated string participant_components = 3;
  string total_attention_pool = 4;
  map<string, string> attention_distribution = 5;  // component_id -> attention_amount
  google.protobuf.Timestamp created_at = 6;
  string distribution_strategy = 7; // EQUAL, WEIGHTED, DYNAMIC
}
```

#### 4. EnergyBalance
```protobuf
message EnergyBalance {
  string lct_id = 1;
  string total_energy_in_wh = 2;
  string total_energy_out_wh = 3;
  string current_balance_wh = 4;
  string reserved_energy_wh = 5;
  google.protobuf.Timestamp last_updated = 6;
  repeated EnergyFlowEntry recent_flows = 7;
}

message EnergyFlowEntry {
  string operation_id = 1;
  string amount_wh = 2;
  string direction = 3;  // IN, OUT
  google.protobuf.Timestamp timestamp = 4;
}
```

### Store Layout
```
energycycle/
├── operations/
│   └── {operation_id} → EnergyOperation
├── balances/
│   └── {lct_id} → EnergyBalance
├── atp_tokens/
│   └── {lct_id}/
│       └── {from_component}/{to_component} → RelationshipATPToken
├── adp_tokens/
│   └── {lct_id} → RelationshipADPToken
├── history/
│   └── {lct_id}/
│       └── {year}/{month}/{day} → [operation_ids]
└── params → Params
```

## Messages (Transactions)

### 1. CreateRelationshipEnergyOperation
Creates a new energy operation between components.

**Input**:
```protobuf
message MsgCreateRelationshipEnergyOperation {
  string creator = 1;
  string lct_id = 2;
  string source_component = 3;
  string target_component = 4;
  string energy_amount_wh = 5;
  string transfer_rate_w = 6;
  string operation_type = 7;
  string priority = 8;
  map<string, string> metadata = 9;
}
```

**Validation**:
- LCT must exist and be active
- Both components must be participants in LCT
- Source component must have sufficient energy balance
- Trust score must meet minimum requirements
- Transfer rate within component capabilities

**Effects**:
- Creates EnergyOperation record
- Reserves energy from source balance
- Validates trust score requirements
- Emits `energy_operation_created` event

**Example**:
```bash
racecar-webd tx energycycle create-relationship-energy-operation \
  --lct-id="lct_pack_module_001" \
  --source="comp_pack1" \
  --target="comp_module1" \
  --energy-amount="50.5" \
  --transfer-rate="100" \
  --operation-type="DIRECT" \
  --priority="NORMAL" \
  --from mykey
```

### 2. ExecuteEnergyTransfer
Executes a validated energy operation.

**Input**:
```protobuf
message MsgExecuteEnergyTransfer {
  string creator = 1;
  string operation_id = 2;
  string execution_parameters = 3;  // JSON string with execution details
}
```

**Validation**:
- Operation must exist and be in VALIDATED status
- Creator must have execution rights
- All safety checks must pass
- Trust score still valid

**Effects**:
- Updates operation status to EXECUTING then COMPLETED
- Transfers energy between component balances
- Updates relationship energy balance
- Consumes ATP tokens if applicable
- Records operation in history
- Emits `energy_transfer_executed` event

### 3. ValidateRelationshipValue
Validates and updates the value of a relationship using ATP/ADP tokens.

**Input**:
```protobuf
message MsgValidateRelationshipValue {
  string creator = 1;
  string lct_id = 2;
  string validation_type = 3;    // ATP, ADP
  string attention_amount = 4;
  string validation_data = 5;    // JSON with validation details
}
```

**Validation**:
- LCT must exist and creator must be participant
- Sufficient attention tokens available
- Validation data properly formatted

**Effects**:
- Allocates ATP or ADP tokens to relationship
- Updates relationship priority level
- Records attention consumption
- Emits `relationship_value_validated` event

## Queries

### 1. GetRelationshipEnergyBalance
Gets the current energy balance for a relationship.

**Request**:
```protobuf
message QueryGetRelationshipEnergyBalanceRequest {
  string lct_id = 1;
}
```

**Response**:
```protobuf
message QueryGetRelationshipEnergyBalanceResponse {
  EnergyBalance balance = 1;
  repeated EnergyOperation recent_operations = 2;
  string efficiency_score = 3;
}
```

### 2. GetEnergyFlowHistory
Gets the energy flow history for a relationship.

**Request**:
```protobuf
message QueryGetEnergyFlowHistoryRequest {
  string lct_id = 1;
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
  string flow_direction = 4;  // IN, OUT, ALL
  cosmos.base.query.v1beta1.PageRequest pagination = 5;
}
```

**Response**:
```protobuf
message QueryGetEnergyFlowHistoryResponse {
  repeated EnergyOperation operations = 1;
  EnergyFlowSummary summary = 2;
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

### 3. CalculateRelationshipV3
Calculates the comprehensive value of a relationship (Version 3 algorithm).

**Request**:
```protobuf
message QueryCalculateRelationshipV3Request {
  string lct_id = 1;
  string calculation_parameters = 2;  // JSON with calculation options
}
```

**Response**:
```protobuf
message QueryCalculateRelationshipV3Response {
  string relationship_value = 1;
  string energy_efficiency = 2;
  string trust_factor = 3;
  string attention_utilization = 4;
  RelationshipMetrics detailed_metrics = 5;
}
```

## Events

### energy_operation_created
Emitted when a new energy operation is created.
```json
{
  "type": "energy_operation_created",
  "attributes": [
    {"key": "operation_id", "value": "op_energy_001"},
    {"key": "lct_id", "value": "lct_pack_module_001"},
    {"key": "source_component", "value": "comp_pack1"},
    {"key": "target_component", "value": "comp_module1"},
    {"key": "energy_amount_wh", "value": "50.5"},
    {"key": "operation_type", "value": "DIRECT"}
  ]
}
```

### energy_transfer_executed
Emitted when energy transfer is completed.
```json
{
  "type": "energy_transfer_executed",
  "attributes": [
    {"key": "operation_id", "value": "op_energy_001"},
    {"key": "actual_energy_wh", "value": "50.3"},
    {"key": "transfer_efficiency", "value": "0.996"},
    {"key": "execution_time_ms", "value": "1250"},
    {"key": "trust_score_used", "value": "0.85"}
  ]
}
```

### relationship_value_validated
Emitted when relationship value is validated with attention tokens.
```json
{
  "type": "relationship_value_validated",
  "attributes": [
    {"key": "lct_id", "value": "lct_pack_module_001"},
    {"key": "validation_type", "value": "ATP"},
    {"key": "attention_allocated", "value": "100"},
    {"key": "new_priority_level", "value": "HIGH"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Minimum trust score required for energy operations
  string min_trust_score_for_transfer = 1;
  
  // Maximum energy transfer per operation (Wh)
  string max_energy_per_operation_wh = 2;
  
  // Transfer efficiency factor (0.0 to 1.0)
  string default_transfer_efficiency = 3;
  
  // Energy operation timeout (seconds)
  uint32 operation_timeout_seconds = 4;
  
  // Enable ATP/ADP token validation
  bool enable_attention_tokens = 5;
  
  // Default attention pool size
  string default_attention_pool_size = 6;
  
  // History retention period (days)
  uint32 history_retention_days = 7;
}
```

**Default Values**:
```json
{
  "min_trust_score_for_transfer": "0.5",
  "max_energy_per_operation_wh": "1000.0",
  "default_transfer_efficiency": "0.95",
  "operation_timeout_seconds": 300,
  "enable_attention_tokens": true,
  "default_attention_pool_size": "1000",
  "history_retention_days": 365
}
```

## Energy Accounting

### Balance Tracking Algorithm

```go
func UpdateEnergyBalance(lctID string, operation EnergyOperation) error {
    balance, exists := getEnergyBalance(lctID)
    if !exists {
        balance = NewEnergyBalance(lctID)
    }
    
    switch operation.Status {
    case "PENDING":
        // Reserve energy from source
        balance.ReservedEnergyWh += operation.EnergyAmountWh
        
    case "COMPLETED":
        // Complete the transfer
        balance.ReservedEnergyWh -= operation.EnergyAmountWh
        
        if operation.Source == getComponentFromLCT(lctID) {
            balance.TotalEnergyOutWh += operation.EnergyAmountWh
            balance.CurrentBalanceWh -= operation.EnergyAmountWh
        } else {
            balance.TotalEnergyInWh += operation.EnergyAmountWh
            balance.CurrentBalanceWh += operation.EnergyAmountWh
        }
        
    case "FAILED":
        // Release reserved energy
        balance.ReservedEnergyWh -= operation.EnergyAmountWh
    }
    
    balance.LastUpdated = time.Now()
    return setEnergyBalance(lctID, balance)
}
```

### Efficiency Calculation

```go
func CalculateTransferEfficiency(operation EnergyOperation) float64 {
    // Base efficiency from parameters
    baseEfficiency := params.DefaultTransferEfficiency
    
    // Trust factor adjustment
    trustFactor := operation.ActualTrustScore / operation.TrustScoreRequired
    
    // Distance/resistance factor (from metadata)
    resistanceFactor := getResistanceFactor(operation.Metadata)
    
    // Combined efficiency
    efficiency := baseEfficiency * trustFactor * resistanceFactor
    
    // Cap at 1.0
    if efficiency > 1.0 {
        efficiency = 1.0
    }
    
    return efficiency
}
```

### ATP Token Management

```go
func AllocateATPTokens(lctID, fromComp, toComp string, amount int) error {
    // Check available attention budget
    budget := getComponentAttentionBudget(fromComp)
    if budget.Available < amount {
        return errors.New("insufficient attention budget")
    }
    
    // Create ATP allocation
    atp := RelationshipATPToken{
        LctId: lctID,
        FromComponent: fromComp,
        ToComponent: toComp,
        AllocatedAttention: amount,
        AllocationTime: time.Now(),
        ExpiryTime: time.Now().Add(24 * time.Hour),
        PriorityLevel: calculatePriorityLevel(amount),
    }
    
    // Update budgets
    budget.Allocated += amount
    budget.Available -= amount
    
    return setATPToken(lctID, fromComp, toComp, atp)
}
```

## Integration Guide

### For Battery Management Systems

1. **Initiating Energy Transfers**:
```go
// Check relationship trust before transfer
trustResp, err := trustTensorClient.CalculateRelationshipTrust(ctx,
    &trusttensor.QueryCalculateRelationshipTrustRequest{
        LctId: lctID,
    })

if err != nil || trustResp.TrustScore < requiredTrustScore {
    return errors.New("insufficient trust for energy transfer")
}

// Create energy operation
msg := &types.MsgCreateRelationshipEnergyOperation{
    Creator: controllerAddr,
    LctId: lctID,
    SourceComponent: sourceModuleID,
    TargetComponent: targetModuleID,
    EnergyAmountWh: "25.5",
    TransferRateW: "50",
    OperationType: "DIRECT",
    Priority: "NORMAL",
}
```

2. **Monitoring Energy Flows**:
```go
// Get current balance
balanceResp, err := queryClient.GetRelationshipEnergyBalance(ctx,
    &types.QueryGetRelationshipEnergyBalanceRequest{
        LctId: lctID,
    })

// Check for imbalances
if balanceResp.Balance.CurrentBalanceWh < minimumBalance {
    triggerRebalancing(lctID)
}

// Monitor efficiency
if balanceResp.EfficiencyScore < 0.9 {
    investigateInefficiency(lctID)
}
```

### For Energy Optimization Systems

1. **Attention-Based Prioritization**:
```go
// Allocate higher attention to critical relationships
msg := &types.MsgValidateRelationshipValue{
    Creator: optimizerAddr,
    LctId: criticalLctID,
    ValidationType: "ATP",
    AttentionAmount: "500",  // High attention allocation
    ValidationData: marshalValidationData(ValidationData{
        Priority: "CRITICAL",
        Reason: "Battery low - emergency charging",
        Duration: "30m",
    }),
}
```

2. **Load Balancing with ADP**:
```go
// Distribute attention across multiple modules
msg := &types.MsgValidateRelationshipValue{
    Creator: balancerAddr,
    LctId: packLctID,
    ValidationType: "ADP",
    AttentionAmount: "1000",
    ValidationData: marshalADPData(ADPData{
        Strategy: "WEIGHTED",
        Weights: map[string]int{
            "module1": 300,
            "module2": 250,
            "module3": 200,
            "module4": 250,
        },
    }),
}
```

### For Module Developers

The Energy Cycle module exposes the following keeper interface:

```go
type EnergyCycleKeeper interface {
    // Operation management
    CreateEnergyOperation(ctx sdk.Context, op types.EnergyOperation) error
    GetEnergyOperation(ctx sdk.Context, opID string) (types.EnergyOperation, bool)
    ExecuteEnergyTransfer(ctx sdk.Context, opID string) error
    
    // Balance management
    GetEnergyBalance(ctx sdk.Context, lctID string) (types.EnergyBalance, bool)
    UpdateEnergyBalance(ctx sdk.Context, lctID string, op types.EnergyOperation) error
    
    // Attention token management
    AllocateATPTokens(ctx sdk.Context, lctID, from, to string, amount sdk.Int) error
    ConsumeATPTokens(ctx sdk.Context, lctID, from, to string, amount sdk.Int) error
    DistributeADPTokens(ctx sdk.Context, lctID string, distribution map[string]sdk.Int) error
    
    // History and analytics
    GetEnergyFlowHistory(ctx sdk.Context, lctID string, startTime, endTime time.Time) []types.EnergyOperation
    CalculateRelationshipValue(ctx sdk.Context, lctID string) (sdk.Dec, error)
}
```

Example integration in Trust Tensor module:

```go
func (k Keeper) UpdateTrustFromEnergyOperation(ctx sdk.Context, operation types.EnergyOperation) {
    // Get current trust tensor
    tensor, found := k.GetRelationshipTensor(ctx, operation.LctId)
    if !found {
        return
    }
    
    // Calculate trust adjustment based on operation success
    var adjustment sdk.Dec
    if operation.Status == "COMPLETED" {
        // Successful operation increases trust
        efficiency := calculateEfficiency(operation)
        adjustment = sdk.NewDecWithPrec(int64(efficiency*100), 2) // 0.01 precision
    } else if operation.Status == "FAILED" {
        // Failed operation decreases trust
        adjustment = sdk.NewDecWithPrec(-5, 2) // -0.05
    }
    
    // Update tensor score
    newScore := tensor.TrustScore.Add(adjustment)
    if newScore.GT(sdk.OneDec()) {
        newScore = sdk.OneDec()
    }
    if newScore.IsNegative() {
        newScore = sdk.ZeroDec()
    }
    
    tensor.TrustScore = newScore
    tensor.LastUpdated = ctx.BlockTime()
    
    k.SetRelationshipTensor(ctx, tensor)
}
```

### Best Practices

1. **Energy Safety**:
   - Always validate trust scores before large transfers
   - Implement maximum transfer limits
   - Monitor for abnormal energy patterns

2. **Efficiency Optimization**:
   - Use ATP tokens to prioritize critical transfers
   - Implement load balancing with ADP tokens
   - Track and optimize transfer efficiency

3. **History Management**:
   - Regularly archive old history data
   - Use history for trend analysis
   - Implement anomaly detection on energy patterns

4. **Error Handling**:
   - Handle transfer failures gracefully
   - Implement automatic retry for transient failures
   - Alert operators to systematic issues