package keeper

import (
	"context"
	"fmt"
	"time"

	"racecar-web/x/componentregistry/types"

	"cosmossdk.io/math"
)

// Authorization levels
const (
	AuthLevelBasic    = "basic"
	AuthLevelEnhanced = "enhanced"
	AuthLevelCritical = "critical"
)

// CreatePairingAuthorization creates a new pairing authorization
func (k Keeper) CreatePairingAuthorization(ctx context.Context, componentID, lctID, authLevel, context string) (*types.PairingAuthorization, error) {
	// Verify component exists
	component, err := k.Components.Get(ctx, componentID)
	if err != nil {
		return nil, fmt.Errorf("component not found: %s", componentID)
	}

	// Verify LCT relationship exists and is active (if lctmanagerKeeper is available)
	if k.lctmanagerKeeper != nil {
		lct, found := k.lctmanagerKeeper.GetLinkedContextToken(ctx, lctID)
		if !found {
			return nil, fmt.Errorf("LCT relationship not found: %s", lctID)
		}

		if lct.PairingStatus != "active" {
			return nil, fmt.Errorf("LCT relationship not active: %s", lctID)
		}
	}

	// Get trust score for this LCT relationship (if trusttensorKeeper is available)
	var trustScore string
	if k.trusttensorKeeper != nil {
		trustScore, _, err = k.trusttensorKeeper.CalculateRelationshipTrust(ctx, lctID, "component_authorization")
		if err != nil {
			return nil, fmt.Errorf("failed to calculate trust score: %w", err)
		}
	} else {
		// Default trust score if trust tensor keeper is not available
		trustScore = "0.75"
	}

	// Validate authorization level based on trust score
	if !k.validateAuthorizationLevel(authLevel, trustScore) {
		return nil, fmt.Errorf("insufficient trust score for authorization level: %s", authLevel)
	}

	// Create authorization
	authID := fmt.Sprintf("auth-%s-%d", componentID, time.Now().Unix())
	authorization := &types.PairingAuthorization{
		AuthId:              authID,
		ComponentId:         componentID,
		AllowedPartnerTypes: k.getAllowedPartnerTypes(component.ComponentType, authLevel),
		AllowedPartnerIds:   "", // Will be populated based on context
		ParentChildRules:    k.getParentChildRules(component.ComponentType, authLevel),
		ApplicationContext:  context,
		ExpiresAt:           time.Now().Add(24 * time.Hour).Unix(), // 24 hour expiration
		Status:              "active",
		TrustScore:          trustScore,
		AuthorizationLevel:  authLevel,
		CreatedAt:           time.Now().Unix(),
		UpdatedAt:           time.Now().Unix(),
		LctId:               lctID,
		Version:             1,
	}

	// Store the authorization
	err = k.PairingAuthorizations.Set(ctx, authID, *authorization)
	if err != nil {
		return nil, fmt.Errorf("failed to store authorization: %w", err)
	}

	return authorization, nil
}

// ValidateComponentPairing validates if two components can pair based on authorization
func (k Keeper) ValidateComponentPairing(ctx context.Context, componentAID, componentBID string) (bool, string, error) {
	// Get both components
	compA, err := k.Components.Get(ctx, componentAID)
	if err != nil {
		return false, "component_a_not_found", err
	}

	compB, err := k.Components.Get(ctx, componentBID)
	if err != nil {
		return false, "component_b_not_found", err
	}

	// Check if components are verified
	verifiedA, _ := k.VerifyComponentForPairing(ctx, componentAID)
	verifiedB, _ := k.VerifyComponentForPairing(ctx, componentBID)

	if !verifiedA || !verifiedB {
		return false, "components_not_verified", fmt.Errorf("one or both components not verified")
	}

	// Get active authorizations for both components
	authA, err := k.getActiveAuthorization(ctx, componentAID)
	if err != nil {
		return false, "no_authorization_a", err
	}

	authB, err := k.getActiveAuthorization(ctx, componentBID)
	if err != nil {
		return false, "no_authorization_b", err
	}

	// Check if components can pair based on authorization rules
	canPairAtoB := k.checkPairingRules(ctx, compA, compB, authA)
	canPairBtoA := k.checkPairingRules(ctx, compB, compA, authB)

	// Both directions must be allowed for bidirectional pairing
	if !canPairAtoB || !canPairBtoA {
		return false, "pairing_rules_not_satisfied", fmt.Errorf("pairing rules not satisfied")
	}

	// Check trust scores meet minimum requirements
	minTrustScore := k.getMinimumTrustScore(authA.AuthorizationLevel, authB.AuthorizationLevel)

	trustA, err := math.LegacyNewDecFromStr(authA.TrustScore)
	if err != nil {
		return false, "invalid_trust_score_a", err
	}

	trustB, err := math.LegacyNewDecFromStr(authB.TrustScore)
	if err != nil {
		return false, "invalid_trust_score_b", err
	}

	if trustA.LT(minTrustScore) || trustB.LT(minTrustScore) {
		return false, "insufficient_trust_scores", fmt.Errorf("insufficient trust scores")
	}

	return true, "pairing_authorized", nil
}

// UpdateAuthorization updates an existing authorization
func (k Keeper) UpdateAuthorization(ctx context.Context, authID string, updates map[string]interface{}) error {
	// Get existing authorization
	auth, err := k.PairingAuthorizations.Get(ctx, authID)
	if err != nil {
		return fmt.Errorf("authorization not found: %s", authID)
	}

	// Apply updates
	if allowedTypes, ok := updates["allowed_partner_types"].(string); ok {
		auth.AllowedPartnerTypes = allowedTypes
	}

	if allowedIds, ok := updates["allowed_partner_ids"].(string); ok {
		auth.AllowedPartnerIds = allowedIds
	}

	if status, ok := updates["status"].(string); ok {
		auth.Status = status
	}

	if expiresAt, ok := updates["expires_at"].(int64); ok {
		auth.ExpiresAt = expiresAt
	}

	// Update metadata
	auth.UpdatedAt = time.Now().Unix()
	auth.Version++

	// Store updated authorization
	return k.PairingAuthorizations.Set(ctx, authID, auth)
}

// RevokeAuthorization revokes an authorization
func (k Keeper) RevokeAuthorization(ctx context.Context, authID string) error {
	auth, err := k.PairingAuthorizations.Get(ctx, authID)
	if err != nil {
		return fmt.Errorf("authorization not found: %s", authID)
	}

	auth.Status = "revoked"
	auth.UpdatedAt = time.Now().Unix()
	auth.Version++

	return k.PairingAuthorizations.Set(ctx, authID, auth)
}

// GetComponentAuthorizations gets all authorizations for a component
func (k Keeper) GetComponentAuthorizations(ctx context.Context, componentID string) ([]types.PairingAuthorization, error) {
	var authorizations []types.PairingAuthorization

	// Iterate through all authorizations (simplified implementation)
	iter, err := k.PairingAuthorizations.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		auth, err := iter.Value()
		if err != nil {
			continue
		}

		if auth.ComponentId == componentID && auth.Status == "active" {
			authorizations = append(authorizations, auth)
		}
	}

	return authorizations, nil
}

// validateAuthorizationLevel checks if trust score is sufficient for authorization level
func (k Keeper) validateAuthorizationLevel(authLevel, trustScore string) bool {
	trust, err := math.LegacyNewDecFromStr(trustScore)
	if err != nil {
		return false
	}

	switch authLevel {
	case AuthLevelBasic:
		return trust.GTE(math.LegacyNewDecWithPrec(50, 2)) // 0.50 minimum
	case AuthLevelEnhanced:
		return trust.GTE(math.LegacyNewDecWithPrec(70, 2)) // 0.70 minimum
	case AuthLevelCritical:
		return trust.GTE(math.LegacyNewDecWithPrec(85, 2)) // 0.85 minimum
	default:
		return false
	}
}

// getAllowedPartnerTypes returns allowed partner types based on component type and auth level
func (k Keeper) getAllowedPartnerTypes(componentType, authLevel string) string {
	switch componentType {
	case types.ComponentTypeModule:
		switch authLevel {
		case AuthLevelBasic:
			return "pack,host_ecu"
		case AuthLevelEnhanced:
			return "pack,host_ecu,sensor"
		case AuthLevelCritical:
			return "pack,host_ecu,sensor,module"
		}
	case types.ComponentTypePack:
		switch authLevel {
		case AuthLevelBasic:
			return "module,host_ecu"
		case AuthLevelEnhanced:
			return "module,host_ecu,sensor"
		case AuthLevelCritical:
			return "module,host_ecu,sensor,pack"
		}
	case types.ComponentTypeHostECU:
		switch authLevel {
		case AuthLevelBasic:
			return "module,pack"
		case AuthLevelEnhanced:
			return "module,pack,sensor"
		case AuthLevelCritical:
			return "module,pack,sensor,host_ecu"
		}
	default:
		return ""
	}
	return ""
}

// getParentChildRules returns parent-child relationship rules
func (k Keeper) getParentChildRules(componentType, authLevel string) string {
	switch componentType {
	case types.ComponentTypePack:
		return "can_contain_modules"
	case types.ComponentTypeHostECU:
		return "can_control_packs"
	case types.ComponentTypeModule:
		return "can_be_contained_by_packs"
	default:
		return ""
	}
}

// getActiveAuthorization gets the most recent active authorization for a component
func (k Keeper) getActiveAuthorization(ctx context.Context, componentID string) (types.PairingAuthorization, error) {
	authorizations, err := k.GetComponentAuthorizations(ctx, componentID)
	if err != nil {
		return types.PairingAuthorization{}, err
	}

	if len(authorizations) == 0 {
		return types.PairingAuthorization{}, fmt.Errorf("no active authorizations found")
	}

	// Return the most recent authorization
	latest := authorizations[0]
	for _, auth := range authorizations {
		if auth.CreatedAt > latest.CreatedAt {
			latest = auth
		}
	}

	return latest, nil
}

// checkPairingRules checks if component A can pair with component B based on authorization
func (k Keeper) checkPairingRules(ctx context.Context, compA, compB types.Component, authA types.PairingAuthorization) bool {
	// Check if component B type is in allowed partner types
	allowedTypes := authA.AllowedPartnerTypes
	if allowedTypes == "" {
		return false
	}

	// Simple string contains check (in production, use proper parsing)
	// This is a simplified implementation
	switch compB.ComponentType {
	case types.ComponentTypeModule:
		return allowedTypes == "pack,host_ecu" || allowedTypes == "pack,host_ecu,sensor" || allowedTypes == "pack,host_ecu,sensor,module"
	case types.ComponentTypePack:
		return allowedTypes == "module,host_ecu" || allowedTypes == "module,host_ecu,sensor" || allowedTypes == "module,host_ecu,sensor,pack"
	case types.ComponentTypeHostECU:
		return allowedTypes == "module,pack" || allowedTypes == "module,pack,sensor" || allowedTypes == "module,pack,sensor,host_ecu"
	default:
		return false
	}
}

// getMinimumTrustScore returns the minimum trust score required for pairing
func (k Keeper) getMinimumTrustScore(authLevelA, authLevelB string) math.LegacyDec {
	// Use the higher authorization level requirement
	var minScore math.LegacyDec

	switch authLevelA {
	case AuthLevelCritical:
		minScore = math.LegacyNewDecWithPrec(85, 2)
	case AuthLevelEnhanced:
		minScore = math.LegacyNewDecWithPrec(70, 2)
	default:
		minScore = math.LegacyNewDecWithPrec(50, 2)
	}

	switch authLevelB {
	case AuthLevelCritical:
		if math.LegacyNewDecWithPrec(85, 2).GT(minScore) {
			minScore = math.LegacyNewDecWithPrec(85, 2)
		}
	case AuthLevelEnhanced:
		if math.LegacyNewDecWithPrec(70, 2).GT(minScore) {
			minScore = math.LegacyNewDecWithPrec(70, 2)
		}
	}

	return minScore
}
