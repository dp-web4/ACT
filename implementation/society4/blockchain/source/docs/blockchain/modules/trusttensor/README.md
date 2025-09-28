# Trust Tensor Module

## Table of Contents

1. [Overview](#overview)
2. [Key Concepts](#key-concepts)
3. [State Management](#state-management)
4. [Messages (Transactions)](#messages-transactions)
5. [Queries](#queries)
6. [Events](#events)
7. [Parameters](#parameters)
8. [Trust Mathematics](#trust-mathematics)
9. [Integration Guide](#integration-guide)

## Overview

The Trust Tensor module implements advanced trust scoring and relationship trust management using tensor mathematics. It quantifies trust relationships between components across multiple dimensions, tracks trust evolution over time, and provides trust-based validation for energy operations and other critical system functions.

### Purpose
- Calculate and maintain multi-dimensional trust scores
- Track trust evolution through component interactions
- Provide trust-based validation for system operations
- Implement witness-based trust verification
- Support trust-aware decision making

### Dependencies
- **LCT Manager**: Uses LCT relationships as basis for trust calculations
- **Energy Cycle**: Receives energy operation outcomes to update trust

### Module Store Key
`trusttensor`

## Key Concepts

### Trust Tensor
A multi-dimensional mathematical structure that captures trust relationships:
- **Dimensions**: Different aspects of trust (reliability, authenticity, performance)
- **Tensor Values**: Quantified trust scores per dimension
- **Temporal Evolution**: How trust changes over time
- **Witness Validation**: Third-party trust confirmations

### Trust Dimensions
Core dimensions tracked in the trust tensor:
1. **Reliability**: Historical performance and consistency
2. **Authenticity**: Verification of component identity and credentials
3. **Behavioral**: Adherence to expected operational patterns
4. **Longevity**: Trust building over extended time periods
5. **Witness**: Confirmations from other trusted components

### Tensor Mathematics
Trust calculations use tensor operations:
- **Scalar Operations**: Single dimension calculations
- **Vector Operations**: Cross-dimensional analysis
- **Matrix Operations**: Relationship mapping
- **Tensor Products**: Complex trust interactions

### Witness System
Third-party validation mechanism:
- **Witnesses**: Other components that can vouch for trust
- **Witness Weight**: Reliability of the witness itself
- **Witness Consensus**: Agreement among multiple witnesses
- **Witness Decay**: Time-based reduction of witness influence

## State Management

### Stored Types

#### 1. RelationshipTrustTensor
```protobuf
message RelationshipTrustTensor {
  string tensor_id = 1;
  string lct_id = 2;
  repeated string participant_components = 3;
  string overall_trust_score = 4;  // Computed overall score (0.0 to 1.0)
  repeated TensorDimension dimensions = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp last_updated = 7;
  uint32 interaction_count = 8;
  string decay_factor = 9;
  map<string, string> tensor_metadata = 10;
}

message TensorDimension {
  string dimension_name = 1;      // reliability, authenticity, behavioral, etc.
  string dimension_score = 2;     // Score for this dimension (0.0 to 1.0)
  string dimension_weight = 3;    // Weight of this dimension in overall score
  uint32 evidence_count = 4;      // Number of data points for this dimension
  google.protobuf.Timestamp last_evidence = 5;
}
```

#### 2. TensorEntry
```protobuf
message TensorEntry {
  string entry_id = 1;
  string tensor_id = 2;
  string entry_type = 3;          // INTERACTION, WITNESS, DECAY, MANUAL
  string dimension_name = 4;
  string score_delta = 5;         // Change in score (-1.0 to 1.0)
  string evidence_data = 6;       // JSON data supporting the score change
  string source_component = 7;    // Component that generated this entry
  google.protobuf.Timestamp timestamp = 8;
  string confidence_level = 9;    // Confidence in this evidence (0.0 to 1.0)
}
```

#### 3. ValueTensor
```protobuf
message ValueTensor {
  string value_tensor_id = 1;
  string relationship_id = 2;
  repeated ValueDimension value_dimensions = 3;
  string computed_value_score = 4;
  google.protobuf.Timestamp computation_time = 5;
  string computation_algorithm = 6;  // V1, V2, V3, etc.
}

message ValueDimension {
  string dimension_name = 1;
  string raw_value = 2;
  string normalized_value = 3;
  string weight = 4;
}
```

#### 4. TensorWitness
```protobuf
message TensorWitness {
  string witness_id = 1;
  string tensor_id = 2;
  string witness_component = 3;
  string witnessed_component = 4;
  string witness_statement = 5;    // JSON with witness data
  string confidence_score = 6;     // Witness confidence (0.0 to 1.0)
  string witness_type = 7;         // POSITIVE, NEGATIVE, NEUTRAL
  google.protobuf.Timestamp witnessed_at = 8;
  string witness_signature = 9;    // Cryptographic proof
}
```

### Store Layout
```
trusttensor/
├── tensors/
│   └── {tensor_id} → RelationshipTrustTensor
├── entries/
│   └── {tensor_id}/
│       └── {entry_id} → TensorEntry
├── witnesses/
│   └── {tensor_id}/
│       └── {witness_id} → TensorWitness
├── value_tensors/
│   └── {relationship_id} → ValueTensor
├── history/
│   └── {tensor_id}/
│       └── {year}/{month} → [entry_ids]
└── params → Params
```

## Messages (Transactions)

### 1. CreateRelationshipTensor
Creates a new trust tensor for a relationship.

**Input**:
```protobuf
message MsgCreateRelationshipTensor {
  string creator = 1;
  string lct_id = 2;
  repeated TensorDimension initial_dimensions = 3;
  string initial_trust_score = 4;
  map<string, string> metadata = 5;
}
```

**Validation**:
- LCT must exist and be active
- Creator must be participant in LCT
- Initial dimensions must be valid
- Trust score must be between 0.0 and 1.0

**Effects**:
- Creates RelationshipTrustTensor
- Initializes dimension scores
- Records creation entry
- Emits `tensor_created` event

**Example**:
```bash
racecar-webd tx trusttensor create-relationship-tensor \
  --lct-id="lct_pack_module_001" \
  --initial-trust="0.5" \
  --dimensions="reliability:0.5:0.3,authenticity:0.8:0.2,behavioral:0.4:0.3,longevity:0.1:0.2" \
  --from mykey
```

### 2. UpdateTensorScore
Updates tensor scores based on new evidence.

**Input**:
```protobuf
message MsgUpdateTensorScore {
  string creator = 1;
  string tensor_id = 2;
  string dimension_name = 3;
  string score_delta = 4;          // Change to apply (-1.0 to 1.0)
  string evidence_data = 5;        // JSON evidence
  string confidence_level = 6;     // Evidence confidence (0.0 to 1.0)
}
```

**Validation**:
- Tensor must exist
- Creator must have update rights
- Dimension must be valid
- Score delta within bounds
- Evidence data properly formatted

**Effects**:
- Updates dimension score
- Recalculates overall trust score
- Creates TensorEntry record
- Updates interaction count
- Emits `tensor_score_updated` event

### 3. AddTensorWitness
Adds witness validation to a trust tensor.

**Input**:
```protobuf
message MsgAddTensorWitness {
  string creator = 1;
  string tensor_id = 2;
  string witnessed_component = 3;
  string witness_statement = 4;    // JSON witness data
  string confidence_score = 5;
  string witness_type = 6;         // POSITIVE, NEGATIVE, NEUTRAL
  string witness_signature = 7;
}
```

**Validation**:
- Tensor must exist
- Creator must be authorized witness
- Witnessed component must be tensor participant
- Witness signature must be valid
- Confidence score within bounds

**Effects**:
- Creates TensorWitness record
- Updates witness dimension in tensor
- Recalculates overall trust considering witness input
- Emits `tensor_witness_added` event

## Queries

### 1. GetRelationshipTensor
Gets the trust tensor for a relationship.

**Request**:
```protobuf
message QueryGetRelationshipTensorRequest {
  string tensor_id = 1;
  bool include_history = 2;
}
```

**Response**:
```protobuf
message QueryGetRelationshipTensorResponse {
  RelationshipTrustTensor tensor = 1;
  repeated TensorEntry history = 2;
  repeated TensorWitness witnesses = 3;
  TensorAnalytics analytics = 4;
}
```

### 2. CalculateRelationshipTrust
Calculates the current trust score for a relationship.

**Request**:
```protobuf
message QueryCalculateRelationshipTrustRequest {
  string lct_id = 1;
  string calculation_algorithm = 2;  // V1, V2, V3
  string context_data = 3;           // JSON context for calculation
}
```

**Response**:
```protobuf
message QueryCalculateRelationshipTrustResponse {
  string trust_score = 1;           // Overall calculated trust score
  repeated DimensionScore dimension_scores = 2;
  string calculation_confidence = 3;
  string calculation_metadata = 4;  // JSON with calculation details
}
```

### 3. GetTensorHistory
Gets the evolution history of a trust tensor.

**Request**:
```protobuf
message QueryGetTensorHistoryRequest {
  string tensor_id = 1;
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
  string dimension_filter = 4;      // Optional dimension to focus on
  cosmos.base.query.v1beta1.PageRequest pagination = 5;
}
```

**Response**:
```protobuf
message QueryGetTensorHistoryResponse {
  repeated TensorEntry entries = 1;
  TensorEvolution evolution_summary = 2;
  cosmos.base.query.v1beta1.PageResponse pagination = 3;
}
```

## Events

### tensor_created
Emitted when a new trust tensor is created.
```json
{
  "type": "tensor_created",
  "attributes": [
    {"key": "tensor_id", "value": "tensor_abc123"},
    {"key": "lct_id", "value": "lct_pack_module_001"},
    {"key": "initial_trust_score", "value": "0.5"},
    {"key": "dimension_count", "value": "4"}
  ]
}
```

### tensor_score_updated
Emitted when tensor scores are updated.
```json
{
  "type": "tensor_score_updated",
  "attributes": [
    {"key": "tensor_id", "value": "tensor_abc123"},
    {"key": "dimension", "value": "reliability"},
    {"key": "old_score", "value": "0.5"},
    {"key": "new_score", "value": "0.6"},
    {"key": "score_delta", "value": "0.1"},
    {"key": "evidence_confidence", "value": "0.9"}
  ]
}
```

### tensor_witness_added
Emitted when a witness validates a tensor.
```json
{
  "type": "tensor_witness_added",
  "attributes": [
    {"key": "tensor_id", "value": "tensor_abc123"},
    {"key": "witness_component", "value": "comp_controller1"},
    {"key": "witnessed_component", "value": "comp_module1"},
    {"key": "witness_type", "value": "POSITIVE"},
    {"key": "confidence_score", "value": "0.85"}
  ]
}
```

## Parameters

The module maintains the following parameters:

```protobuf
message Params {
  // Default trust score for new relationships
  string default_initial_trust_score = 1;
  
  // Trust decay factor per day without interaction
  string daily_decay_factor = 2;
  
  // Maximum trust score (usually 1.0)
  string max_trust_score = 3;
  
  // Minimum trust score (usually 0.0)
  string min_trust_score = 4;
  
  // Default dimension weights
  map<string, string> default_dimension_weights = 5;
  
  // Witness influence factor
  string witness_influence_factor = 6;
  
  // Minimum witnesses required for high trust
  uint32 min_witnesses_for_high_trust = 7;
  
  // Trust calculation algorithm version
  string default_calculation_algorithm = 8;
}
```

**Default Values**:
```json
{
  "default_initial_trust_score": "0.5",
  "daily_decay_factor": "0.999",
  "max_trust_score": "1.0",
  "min_trust_score": "0.0",
  "default_dimension_weights": {
    "reliability": "0.3",
    "authenticity": "0.2",
    "behavioral": "0.3",
    "longevity": "0.1",
    "witness": "0.1"
  },
  "witness_influence_factor": "0.1",
  "min_witnesses_for_high_trust": 3,
  "default_calculation_algorithm": "V3"
}
```

## Trust Mathematics

### Overall Trust Score Calculation (V3 Algorithm)

```go
func CalculateTrustScoreV3(tensor RelationshipTrustTensor) sdk.Dec {
    var weightedSum sdk.Dec
    var totalWeight sdk.Dec
    
    // Calculate weighted sum of dimensions
    for _, dimension := range tensor.Dimensions {
        weight := sdk.MustNewDecFromStr(dimension.DimensionWeight)
        score := sdk.MustNewDecFromStr(dimension.DimensionScore)
        
        // Apply evidence confidence factor
        confidenceFactor := calculateConfidenceFactor(dimension.EvidenceCount)
        adjustedScore := score.Mul(confidenceFactor)
        
        weightedSum = weightedSum.Add(weight.Mul(adjustedScore))
        totalWeight = totalWeight.Add(weight)
    }
    
    // Base trust score
    baseTrust := weightedSum.Quo(totalWeight)
    
    // Apply time decay
    decayFactor := calculateDecayFactor(tensor.LastUpdated)
    decayedTrust := baseTrust.Mul(decayFactor)
    
    // Apply witness boost
    witnessBoost := calculateWitnessBoost(tensor.TensorId)
    finalTrust := decayedTrust.Add(witnessBoost)
    
    // Ensure bounds [0, 1]
    if finalTrust.GT(sdk.OneDec()) {
        finalTrust = sdk.OneDec()
    }
    if finalTrust.IsNegative() {
        finalTrust = sdk.ZeroDec()
    }
    
    return finalTrust
}
```

### Decay Function

```go
func calculateDecayFactor(lastUpdated time.Time) sdk.Dec {
    daysSince := time.Since(lastUpdated).Hours() / 24
    dailyDecay := sdk.MustNewDecFromStr("0.999") // 0.1% decay per day
    
    decayFactor := dailyDecay.Power(uint64(daysSince))
    
    // Minimum decay factor to prevent complete trust loss
    minDecay := sdk.MustNewDecFromStr("0.1")
    if decayFactor.LT(minDecay) {
        decayFactor = minDecay
    }
    
    return decayFactor
}
```

### Witness Influence Calculation

```go
func calculateWitnessBoost(tensorID string) sdk.Dec {
    witnesses := getTensorWitnesses(tensorID)
    
    var positiveWitnesses, negativeWitnesses sdk.Dec
    var totalWitnessWeight sdk.Dec
    
    for _, witness := range witnesses {
        witnessReliability := getWitnessReliability(witness.WitnessComponent)
        confidence := sdk.MustNewDecFromStr(witness.ConfidenceScore)
        weight := witnessReliability.Mul(confidence)
        
        totalWitnessWeight = totalWitnessWeight.Add(weight)
        
        switch witness.WitnessType {
        case "POSITIVE":
            positiveWitnesses = positiveWitnesses.Add(weight)
        case "NEGATIVE":
            negativeWitnesses = negativeWitnesses.Add(weight)
        }
    }
    
    if totalWitnessWeight.IsZero() {
        return sdk.ZeroDec()
    }
    
    // Net witness score
    netWitness := positiveWitnesses.Sub(negativeWitnesses).Quo(totalWitnessWeight)
    
    // Apply witness influence factor
    influenceFactor := sdk.MustNewDecFromStr("0.1")
    boost := netWitness.Mul(influenceFactor)
    
    return boost
}
```

### Dimension Score Updates

```go
func updateDimensionScore(tensor *RelationshipTrustTensor, dimension string, delta sdk.Dec, confidence sdk.Dec) {
    for i, dim := range tensor.Dimensions {
        if dim.DimensionName == dimension {
            currentScore := sdk.MustNewDecFromStr(dim.DimensionScore)
            
            // Apply confidence weighting to delta
            weightedDelta := delta.Mul(confidence)
            
            // Learning rate (smaller changes for established relationships)
            learningRate := calculateLearningRate(dim.EvidenceCount)
            adjustedDelta := weightedDelta.Mul(learningRate)
            
            // Update score
            newScore := currentScore.Add(adjustedDelta)
            
            // Ensure bounds [0, 1]
            if newScore.GT(sdk.OneDec()) {
                newScore = sdk.OneDec()
            }
            if newScore.IsNegative() {
                newScore = sdk.ZeroDec()
            }
            
            tensor.Dimensions[i].DimensionScore = newScore.String()
            tensor.Dimensions[i].EvidenceCount++
            tensor.Dimensions[i].LastEvidence = time.Now()
            break
        }
    }
}
```

## Integration Guide

### For Energy Management Systems

1. **Trust-Based Energy Validation**:
```go
// Check trust before authorizing large energy transfers
trustResp, err := queryClient.CalculateRelationshipTrust(ctx,
    &types.QueryCalculateRelationshipTrustRequest{
        LctId: lctID,
        CalculationAlgorithm: "V3",
    })

if err != nil {
    return errors.New("failed to calculate trust")
}

trustScore := sdk.MustNewDecFromStr(trustResp.TrustScore)
requiredTrust := sdk.MustNewDecFromStr("0.7")

if trustScore.LT(requiredTrust) {
    return fmt.Errorf("insufficient trust: %.3f < %.3f", 
        trustScore, requiredTrust)
}
```

2. **Updating Trust from Energy Operations**:
```go
// Update trust based on energy operation outcome
var scoreDelta sdk.Dec
if energyOp.Status == "COMPLETED" {
    // Successful operation increases reliability
    efficiency := calculateEfficiency(energyOp)
    scoreDelta = sdk.NewDecWithPrec(int64(efficiency*10), 3) // 0.001 precision
} else {
    // Failed operation decreases reliability
    scoreDelta = sdk.NewDecWithPrec(-20, 3) // -0.020
}

msg := &types.MsgUpdateTensorScore{
    Creator: systemAddr,
    TensorId: tensorID,
    DimensionName: "reliability",
    ScoreDelta: scoreDelta.String(),
    EvidenceData: marshalEvidenceData(EvidenceData{
        OperationType: energyOp.OperationType,
        EnergyAmount: energyOp.EnergyAmountWh,
        TransferTime: energyOp.ExecutionTime,
        Efficiency: efficiency,
    }),
    ConfidenceLevel: "0.95",
}
```

### For Component Manufacturers

1. **Initial Trust Setup**:
```go
// Create trust tensor when components are first paired
msg := &types.MsgCreateRelationshipTensor{
    Creator: manufacturerAddr,
    LctId: lctID,
    InitialDimensions: []types.TensorDimension{
        {
            DimensionName: "authenticity",
            DimensionScore: "0.9", // High initial authenticity for factory components
            DimensionWeight: "0.3",
        },
        {
            DimensionName: "reliability",
            DimensionScore: "0.5", // Neutral until proven
            DimensionWeight: "0.4",
        },
        {
            DimensionName: "behavioral",
            DimensionScore: "0.5",
            DimensionWeight: "0.3",
        },
    },
    InitialTrustScore: "0.6",
}
```

2. **Witness Network Setup**:
```go
// Have quality assurance system witness component reliability
witnessMsg := &types.MsgAddTensorWitness{
    Creator: qaSystemAddr,
    TensorId: tensorID,
    WitnessedComponent: moduleID,
    WitnessStatement: marshalWitnessData(WitnessData{
        TestResults: qaResults,
        TestDate: time.Now(),
        TestStandard: "IEC-62619",
        QualityGrade: "A",
    }),
    ConfidenceScore: "0.95",
    WitnessType: "POSITIVE",
    WitnessSignature: signWitnessData(qaPrivateKey, witnessData),
}
```

### For Module Developers

The Trust Tensor module exposes the following keeper interface:

```go
type TrustTensorKeeper interface {
    // Tensor management
    CreateRelationshipTensor(ctx sdk.Context, tensor types.RelationshipTrustTensor) error
    GetRelationshipTensor(ctx sdk.Context, tensorID string) (types.RelationshipTrustTensor, bool)
    UpdateTensorScore(ctx sdk.Context, tensorID, dimension string, delta sdk.Dec, evidence string) error
    
    // Trust calculation
    CalculateRelationshipTrust(ctx sdk.Context, lctID string) (sdk.Dec, error)
    GetTrustScoreForOperation(ctx sdk.Context, lctID, operationType string) (sdk.Dec, error)
    
    // Witness management
    AddTensorWitness(ctx sdk.Context, witness types.TensorWitness) error
    GetTensorWitnesses(ctx sdk.Context, tensorID string) []types.TensorWitness
    ValidateWitnessReliability(ctx sdk.Context, witnessComponentID string) (sdk.Dec, error)
    
    // History and analytics
    GetTensorHistory(ctx sdk.Context, tensorID string, start, end time.Time) []types.TensorEntry
    CalculateTrustTrend(ctx sdk.Context, tensorID string, days int) (sdk.Dec, error)
}
```

Example integration in LCT Manager:

```go
func (k Keeper) CreateLCTWithTrust(ctx sdk.Context, lctReq types.MsgCreateLCTRelationship) error {
    // Create the LCT first
    lct := types.LinkedContextToken{
        LctId: generateLCTID(),
        ParticipantComponents: lctReq.ParticipantComponents,
        RelationshipType: lctReq.RelationshipType,
        Status: "ACTIVE",
        CreatedAt: ctx.BlockTime(),
    }
    
    k.SetLCT(ctx, lct)
    
    // Create associated trust tensor
    tensor := types.RelationshipTrustTensor{
        TensorId: generateTensorID(),
        LctId: lct.LctId,
        ParticipantComponents: lct.ParticipantComponents,
        OverallTrustScore: "0.5", // Default initial trust
        Dimensions: getDefaultDimensions(),
        CreatedAt: ctx.BlockTime(),
        LastUpdated: ctx.BlockTime(),
    }
    
    return k.trustTensor.CreateRelationshipTensor(ctx, tensor)
}
```

### Best Practices

1. **Trust Initialization**:
   - Start with moderate trust scores (0.4-0.6)
   - Use higher initial authenticity for factory-verified components
   - Implement gradual trust building through interactions

2. **Evidence Quality**:
   - Include detailed evidence data for trust updates
   - Use appropriate confidence levels based on evidence quality
   - Implement multiple evidence sources for critical updates

3. **Witness Management**:
   - Establish reliable witness networks
   - Regularly validate witness component reliability
   - Use cryptographic signatures for witness statements

4. **Trust Monitoring**:
   - Implement alerts for rapid trust degradation
   - Monitor trust trends over time
   - Use trust analytics for system optimization