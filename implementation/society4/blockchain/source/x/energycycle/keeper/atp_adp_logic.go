package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"

	"racecar-web/x/energycycle/types"
)

// CreateAtpToken creates a new ATP token for energy storage
func (k Keeper) CreateAtpToken(ctx context.Context, lctID, energyAmount, operationID, context string, blockHeight int64) (*types.RelationshipAtpToken, error) {
	// Validate LCT relationship exists and is active
	lct, found := k.lctmanagerKeeper.GetLinkedContextToken(ctx, lctID)
	if !found {
		return nil, fmt.Errorf("LCT relationship not found: %s", lctID)
	}

	if lct.PairingStatus != "active" {
		return nil, fmt.Errorf("LCT relationship not active: %s", lctID)
	}

	// Get trust score for this LCT relationship
	trustScore, _, err := k.trusttensorKeeper.CalculateRelationshipTrust(ctx, lctID, "energy_operation")
	if err != nil {
		return nil, fmt.Errorf("failed to calculate trust score: %w", err)
	}

	// Calculate energy efficiency rating based on trust and context
	efficiencyRating := k.calculateEnergyEfficiency(trustScore, context)

	// Create ATP token
	tokenID := fmt.Sprintf("atp-%s-%d", lctID, time.Now().Unix())
	atpToken := &types.RelationshipAtpToken{
		TokenId:             tokenID,
		LctId:               lctID,
		EnergyAmount:        energyAmount,
		CreatedAt:           time.Now().Unix(),
		OperationId:         operationID,
		Status:              types.AtpStatusActive,
		RelationshipContext: context,
		ExpirationBlock:     blockHeight + 1000, // ATP expires after 1000 blocks
		TrustScore:          trustScore,
		EfficiencyRating:    efficiencyRating,
		Version:             1,
	}

	// Store the ATP token
	err = k.RelationshipAtpTokens.Set(ctx, tokenID, *atpToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store ATP token: %w", err)
	}

	return atpToken, nil
}

// DischargeAtpToken converts an ATP token to ADP token
func (k Keeper) DischargeAtpToken(ctx context.Context, atpTokenID, operationID string, blockHeight int64) (*types.RelationshipAdpToken, error) {
	// Get the ATP token
	atpToken, err := k.RelationshipAtpTokens.Get(ctx, atpTokenID)
	if err != nil {
		return nil, fmt.Errorf("ATP token not found: %s", atpTokenID)
	}

	// Check if ATP is still active
	if atpToken.Status != types.AtpStatusActive {
		return nil, fmt.Errorf("ATP token is not active: %s", atpTokenID)
	}

	// Check if ATP has expired
	if blockHeight > atpToken.ExpirationBlock {
		// Mark ATP as expired
		atpToken.Status = types.AtpStatusExpired
		atpToken.Version++
		k.RelationshipAtpTokens.Set(ctx, atpTokenID, atpToken)
		return nil, fmt.Errorf("ATP token has expired: %s", atpTokenID)
	}

	// Calculate discharge efficiency
	dischargeEfficiency := k.calculateDischargeEfficiency(atpToken.EfficiencyRating, atpToken.TrustScore)

	// Get V3 validation score for the discharge operation
	v3Score, err := k.trusttensorKeeper.CalculateV3CompositeScore(ctx, operationID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate V3 score: %w", err)
	}

	// Create ADP token
	adpTokenID := fmt.Sprintf("adp-%s-%d", atpTokenID, time.Now().Unix())
	adpToken := &types.RelationshipAdpToken{
		TokenId:          adpTokenID,
		OriginalAtpId:    atpTokenID,
		LctId:            atpToken.LctId,
		DischargedAt:     time.Now().Unix(),
		ValueScore:       v3Score.String(),
		ConfirmationData: fmt.Sprintf("discharged_at_block_%d", blockHeight),
		EnergyEfficiency: dischargeEfficiency.String(),
		TrustValidation:  atpToken.TrustScore,
		ValidationBlock:  blockHeight,
		OperationContext: atpToken.RelationshipContext,
		Version:          1,
	}

	// Store the ADP token
	err = k.RelationshipAdpTokens.Set(ctx, adpTokenID, *adpToken)
	if err != nil {
		return nil, fmt.Errorf("failed to store ADP token: %w", err)
	}

	// Mark ATP token as discharged
	atpToken.Status = types.AtpStatusDischarged
	atpToken.Version++
	k.RelationshipAtpTokens.Set(ctx, atpTokenID, atpToken)

	return adpToken, nil
}

// CalculateEnergyBalance calculates the current energy balance for an LCT relationship
func (k Keeper) CalculateEnergyBalance(ctx context.Context, lctID string) (math.LegacyDec, error) {
	// Get all active ATP tokens for this LCT
	activeAtpTokens, err := k.getActiveAtpTokensForLct(ctx, lctID)
	if err != nil {
		return math.LegacyZeroDec(), fmt.Errorf("failed to get active ATP tokens: %w", err)
	}

	// Sum up energy amounts
	totalEnergy := math.LegacyZeroDec()
	for _, atpToken := range activeAtpTokens {
		energyAmount, err := math.LegacyNewDecFromStr(atpToken.EnergyAmount)
		if err != nil {
			continue // Skip invalid amounts
		}
		totalEnergy = totalEnergy.Add(energyAmount)
	}

	return totalEnergy, nil
}

// ValidateEnergyOperation validates an energy operation using ATP/ADP logic
func (k Keeper) ValidateEnergyOperation(ctx context.Context, operationID, sourceLct, targetLct, energyAmount, operationType string) (bool, string, error) {
	// Check if source LCT has sufficient energy balance
	if operationType == types.OperationTypeTransfer || operationType == types.OperationTypeDischarge {
		balance, err := k.CalculateEnergyBalance(ctx, sourceLct)
		if err != nil {
			return false, "balance_calculation_failed", err
		}

		requiredEnergy, err := math.LegacyNewDecFromStr(energyAmount)
		if err != nil {
			return false, "invalid_energy_amount", fmt.Errorf("invalid energy amount: %s", energyAmount)
		}

		if balance.LT(requiredEnergy) {
			return false, "insufficient_energy_balance", fmt.Errorf("insufficient energy balance: %s < %s", balance.String(), requiredEnergy.String())
		}
	}

	// Get trust scores for both LCTs
	sourceTrustStr, _, err := k.trusttensorKeeper.CalculateRelationshipTrust(ctx, sourceLct, "energy_operation")
	if err != nil {
		return false, "source_trust_calculation_failed", err
	}

	targetTrustStr, _, err := k.trusttensorKeeper.CalculateRelationshipTrust(ctx, targetLct, "energy_operation")
	if err != nil {
		return false, "target_trust_calculation_failed", err
	}

	// Convert trust scores to decimals
	sourceTrust, err := math.LegacyNewDecFromStr(sourceTrustStr)
	if err != nil {
		return false, "invalid_source_trust_score", err
	}

	targetTrust, err := math.LegacyNewDecFromStr(targetTrustStr)
	if err != nil {
		return false, "invalid_target_trust_score", err
	}

	// Calculate composite trust score
	compositeTrust := sourceTrust.Add(targetTrust).Quo(math.LegacyNewDec(2))

	// Minimum trust threshold for energy operations
	minTrustThreshold := math.LegacyNewDecWithPrec(60, 2) // 0.60
	if compositeTrust.LT(minTrustThreshold) {
		return false, "insufficient_trust_score", fmt.Errorf("insufficient trust score: %s < %s", compositeTrust.String(), minTrustThreshold.String())
	}

	return true, "operation_validated", nil
}

// calculateEnergyEfficiency calculates energy efficiency rating
func (k Keeper) calculateEnergyEfficiency(trustScore, context string) string {
	trust, err := math.LegacyNewDecFromStr(trustScore)
	if err != nil {
		return "0.5" // Default efficiency
	}

	// Base efficiency starts at 0.5 and scales with trust
	baseEfficiency := math.LegacyNewDecWithPrec(50, 2)                                              // 0.50
	trustBonus := trust.Sub(math.LegacyNewDecWithPrec(50, 2)).Mul(math.LegacyNewDecWithPrec(80, 2)) // Up to 0.40 bonus

	efficiency := baseEfficiency.Add(trustBonus)

	// Apply context-specific modifiers
	switch context {
	case "high_performance":
		efficiency = efficiency.Mul(math.LegacyNewDecWithPrec(110, 2)) // 10% boost
	case "efficiency_optimized":
		efficiency = efficiency.Mul(math.LegacyNewDecWithPrec(120, 2)) // 20% boost
	case "safety_critical":
		efficiency = efficiency.Mul(math.LegacyNewDecWithPrec(90, 2)) // 10% penalty for safety
	}

	// Ensure bounds [0.1, 1.0]
	if efficiency.GT(math.LegacyOneDec()) {
		efficiency = math.LegacyOneDec()
	}
	if efficiency.LT(math.LegacyNewDecWithPrec(10, 2)) {
		efficiency = math.LegacyNewDecWithPrec(10, 2)
	}

	return efficiency.String()
}

// calculateDischargeEfficiency calculates discharge efficiency
func (k Keeper) calculateDischargeEfficiency(atpEfficiency, trustScore string) math.LegacyDec {
	atpEff, err := math.LegacyNewDecFromStr(atpEfficiency)
	if err != nil {
		atpEff = math.LegacyNewDecWithPrec(50, 2)
	}

	trust, err := math.LegacyNewDecFromStr(trustScore)
	if err != nil {
		trust = math.LegacyNewDecWithPrec(50, 2)
	}

	// Discharge efficiency = ATP efficiency × trust score × 0.95 (discharge loss)
	dischargeEfficiency := atpEff.Mul(trust).Mul(math.LegacyNewDecWithPrec(95, 2))

	// Ensure bounds [0.1, 1.0]
	if dischargeEfficiency.GT(math.LegacyOneDec()) {
		dischargeEfficiency = math.LegacyOneDec()
	}
	if dischargeEfficiency.LT(math.LegacyNewDecWithPrec(10, 2)) {
		dischargeEfficiency = math.LegacyNewDecWithPrec(10, 2)
	}

	return dischargeEfficiency
}

// getActiveAtpTokensForLct retrieves all active ATP tokens for an LCT
func (k Keeper) getActiveAtpTokensForLct(ctx context.Context, lctID string) ([]types.RelationshipAtpToken, error) {
	var activeTokens []types.RelationshipAtpToken

	// Iterate through all ATP tokens (this is a simplified implementation)
	// In production, you'd want indexed queries
	iter, err := k.RelationshipAtpTokens.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		token, err := iter.Value()
		if err != nil {
			continue
		}

		if token.LctId == lctID && token.Status == types.AtpStatusActive {
			activeTokens = append(activeTokens, token)
		}
	}

	return activeTokens, nil
}

// GetAtpToken retrieves an ATP token by ID
func (k Keeper) GetAtpToken(ctx context.Context, tokenID string) (types.RelationshipAtpToken, error) {
	return k.RelationshipAtpTokens.Get(ctx, tokenID)
}

// GetAdpToken retrieves an ADP token by ID
func (k Keeper) GetAdpToken(ctx context.Context, tokenID string) (types.RelationshipAdpToken, error) {
	return k.RelationshipAdpTokens.Get(ctx, tokenID)
}

// GetEnergyOperation retrieves an energy operation by ID
func (k Keeper) GetEnergyOperation(ctx context.Context, operationID string) (types.EnergyOperation, error) {
	return k.EnergyOperations.Get(ctx, operationID)
}
