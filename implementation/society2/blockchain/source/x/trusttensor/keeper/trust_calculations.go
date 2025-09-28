package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"

	"racecar-web/x/trusttensor/types"
)

// CalculateT3CompositeScore calculates the composite T3 trust score
func (k Keeper) CalculateT3CompositeScore(ctx context.Context, lctID string) (math.LegacyDec, error) {
	tensor, found := k.GetRelationshipTensor(ctx, lctID)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("tensor not found for LCT: %s", lctID)
	}

	// T3 Formula: (Talent × 0.3) + (Training × 0.4) + (Temperament × 0.3)
	talent, err := math.LegacyNewDecFromStr(tensor.TalentScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid talent score: %w", err)
	}

	training, err := math.LegacyNewDecFromStr(tensor.TrainingScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid training score: %w", err)
	}

	temperament, err := math.LegacyNewDecFromStr(tensor.TemperamentScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid temperament score: %w", err)
	}

	// Weighted calculation
	talentWeight := math.LegacyNewDecWithPrec(30, 2)      // 0.30
	trainingWeight := math.LegacyNewDecWithPrec(40, 2)    // 0.40
	temperamentWeight := math.LegacyNewDecWithPrec(30, 2) // 0.30

	composite := talent.Mul(talentWeight).
		Add(training.Mul(trainingWeight)).
		Add(temperament.Mul(temperamentWeight))

	// Apply context modifier if exists
	if tensor.ContextModifier != "" {
		contextModifier, err := math.LegacyNewDecFromStr(tensor.ContextModifier)
		if err == nil {
			composite = composite.Mul(contextModifier)
		}
	}

	// Ensure bounds [0, 1]
	if composite.GT(math.LegacyOneDec()) {
		composite = math.LegacyOneDec()
	}
	if composite.IsNegative() {
		composite = math.LegacyZeroDec()
	}

	return composite, nil
}

// CalculateV3CompositeScore calculates V3 validation scores
func (k Keeper) CalculateV3CompositeScore(ctx context.Context, operationID string) (math.LegacyDec, error) {
	// Get V3 tensor for operation
	tensor, found := k.GetOperationV3Tensor(ctx, operationID)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("V3 tensor not found for operation: %s", operationID)
	}

	// V3 Formula: (Valuation × 0.4) + (Veracity × 0.3) + (Validity × 0.3)
	valuation, err := math.LegacyNewDecFromStr(tensor.ValuationScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid valuation score: %w", err)
	}

	veracity, err := math.LegacyNewDecFromStr(tensor.VeracityScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid veracity score: %w", err)
	}

	validity, err := math.LegacyNewDecFromStr(tensor.ValidityScore)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("invalid validity score: %w", err)
	}

	// Weighted calculation
	valuationWeight := math.LegacyNewDecWithPrec(40, 2) // 0.40
	veracityWeight := math.LegacyNewDecWithPrec(30, 2)  // 0.30
	validityWeight := math.LegacyNewDecWithPrec(30, 2)  // 0.30

	composite := valuation.Mul(valuationWeight).
		Add(veracity.Mul(veracityWeight)).
		Add(validity.Mul(validityWeight))

	// Ensure bounds [0, 1]
	if composite.GT(math.LegacyOneDec()) {
		composite = math.LegacyOneDec()
	}
	if composite.IsNegative() {
		composite = math.LegacyZeroDec()
	}

	return composite, nil
}

// UpdateTensorScore updates a tensor dimension with evidence
func (k Keeper) UpdateTensorScore(ctx context.Context, tensorID, dimension string, newScore math.LegacyDec, evidence string) error {
	tensor, found := k.GetRelationshipTensor(ctx, tensorID)
	if !found {
		return fmt.Errorf("tensor not found: %s", tensorID)
	}

	// Calculate learning rate based on evidence count
	evidenceCount := tensor.EvidenceCount
	learningRate := k.calculateLearningRate(evidenceCount)

	// Apply weighted update
	currentScore := k.getDimensionScore(tensor, dimension)
	scoreDelta := newScore.Sub(currentScore)
	weightedDelta := scoreDelta.Mul(learningRate)

	finalScore := currentScore.Add(weightedDelta)

	// Ensure bounds [0, 1]
	if finalScore.GT(math.LegacyOneDec()) {
		finalScore = math.LegacyOneDec()
	}
	if finalScore.IsNegative() {
		finalScore = math.LegacyZeroDec()
	}

	// Update the tensor
	k.setDimensionScore(ctx, tensorID, dimension, finalScore, evidence)

	return nil
}

// calculateLearningRate decreases as evidence count increases
func (k Keeper) calculateLearningRate(evidenceCount int64) math.LegacyDec {
	// Simplified learning rate = 1 / (1 + evidenceCount/10)
	// More evidence = smaller learning rate (more stable)
	evidenceFactor := math.LegacyNewDec(evidenceCount).Quo(math.LegacyNewDec(10))
	rate := math.LegacyOneDec().Quo(math.LegacyOneDec().Add(evidenceFactor))

	// Cap learning rate between 0.01 and 0.5
	if rate.LT(math.LegacyNewDecWithPrec(1, 2)) {
		rate = math.LegacyNewDecWithPrec(1, 2) // 0.01 minimum
	}
	if rate.GT(math.LegacyNewDecWithPrec(50, 2)) {
		rate = math.LegacyNewDecWithPrec(50, 2) // 0.50 maximum
	}

	return rate
}

// getDimensionScore extracts the score for a specific dimension
func (k Keeper) getDimensionScore(tensor types.RelationshipTrustTensor, dimension string) math.LegacyDec {
	switch dimension {
	case "talent":
		score, _ := math.LegacyNewDecFromStr(tensor.TalentScore)
		return score
	case "training":
		score, _ := math.LegacyNewDecFromStr(tensor.TrainingScore)
		return score
	case "temperament":
		score, _ := math.LegacyNewDecFromStr(tensor.TemperamentScore)
		return score
	default:
		return math.LegacyZeroDec()
	}
}

// setDimensionScore updates the score for a specific dimension
func (k Keeper) setDimensionScore(ctx context.Context, tensorID, dimension string, score math.LegacyDec, evidence string) {
	tensor, found := k.GetRelationshipTensor(ctx, tensorID)
	if !found {
		return
	}

	// Update the appropriate dimension
	switch dimension {
	case "talent":
		tensor.TalentScore = score.String()
	case "training":
		tensor.TrainingScore = score.String()
	case "temperament":
		tensor.TemperamentScore = score.String()
	default:
		return
	}

	// Update metadata
	tensor.EvidenceCount++
	tensor.UpdatedAt = time.Now().Unix()
	tensor.Version++

	// Store the updated tensor
	k.SetRelationshipTensor(ctx, tensorID, tensor)
}

// GetRelationshipTensor retrieves a relationship tensor by LCT ID
func (k Keeper) GetRelationshipTensor(ctx context.Context, lctID string) (types.RelationshipTrustTensor, bool) {
	tensor, err := k.RelationshipTensors.Get(ctx, lctID)
	if err != nil {
		// Return default tensor if not found
		return types.RelationshipTrustTensor{
			TensorId:         fmt.Sprintf("tensor-%s", lctID),
			LctId:            lctID,
			TensorType:       "T3",
			TalentScore:      "0.5",
			TrainingScore:    "0.5",
			TemperamentScore: "0.5",
			Context:          "default",
			CreatedAt:        time.Now().Unix(),
			UpdatedAt:        time.Now().Unix(),
			Version:          1,
			EvidenceCount:    0,
			ContextModifier:  "1.0",
		}, false
	}
	return tensor, true
}

// SetRelationshipTensor stores a relationship tensor
func (k Keeper) SetRelationshipTensor(ctx context.Context, lctID string, tensor types.RelationshipTrustTensor) error {
	return k.RelationshipTensors.Set(ctx, lctID, tensor)
}

// GetOperationV3Tensor retrieves a V3 tensor by operation ID
func (k Keeper) GetOperationV3Tensor(ctx context.Context, operationID string) (types.ValueTensor, bool) {
	tensor, err := k.ValueTensors.Get(ctx, operationID)
	if err != nil {
		// Return default V3 tensor if not found
		return types.ValueTensor{
			TensorId:       fmt.Sprintf("v3-%s", operationID),
			OperationId:    operationID,
			ValuationScore: "0.5",
			VeracityScore:  "0.5",
			ValidityScore:  "0.5",
			CreatedAt:      time.Now().Unix(),
		}, false
	}
	return tensor, true
}

// SetOperationV3Tensor stores a V3 tensor
func (k Keeper) SetOperationV3Tensor(ctx context.Context, operationID string, tensor types.ValueTensor) error {
	return k.ValueTensors.Set(ctx, operationID, tensor)
}

// GetContextModifier retrieves a context modifier for trust calculations
func (k Keeper) GetContextModifier(ctx context.Context, context string) math.LegacyDec {
	// Context-specific modifiers for different operational contexts
	switch context {
	case "energy_operation":
		return math.LegacyNewDecWithPrec(110, 2) // 1.10 - higher trust for energy operations
	case "energy_balance":
		return math.LegacyNewDecWithPrec(105, 2) // 1.05 - slightly higher trust for balance queries
	case "critical_safety":
		return math.LegacyNewDecWithPrec(120, 2) // 1.20 - much higher trust for safety operations
	case "diagnostic":
		return math.LegacyNewDecWithPrec(95, 2) // 0.95 - slightly lower trust for diagnostics
	default:
		return math.LegacyOneDec() // 1.00 - neutral modifier
	}
}
