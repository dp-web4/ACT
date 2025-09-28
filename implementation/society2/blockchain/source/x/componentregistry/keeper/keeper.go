package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"

	"racecar-web/x/componentregistry/types"
	lctmanagertypes "racecar-web/x/lctmanager/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte

	Schema collections.Schema

	// Component storage
	Params                 collections.Item[types.Params]
	Components             collections.Map[string, types.Component]
	ComponentVerifications collections.Map[string, types.ComponentVerification]
	ComponentPairingRules  collections.Map[string, types.ComponentPairingRule]
	ManufacturerComponents collections.Map[string, types.Component] // manufacturer_id -> component (simplified)
	PairingAuthorizations  collections.Map[string, types.PairingAuthorization]

	// Pluggable verification backend
	verificationBackend types.ComponentVerificationBackend

	// Dependencies
	trusttensorKeeper types.TrusttensorKeeper
	lctmanagerKeeper  lctmanagertypes.LctmanagerKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	verificationBackend types.ComponentVerificationBackend,
	trusttensorKeeper types.TrusttensorKeeper,
	lctmanagerKeeper lctmanagertypes.LctmanagerKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:        storeService,
		cdc:                 cdc,
		addressCodec:        addressCodec,
		authority:           authority,
		verificationBackend: verificationBackend,
		trusttensorKeeper:   trusttensorKeeper,
		lctmanagerKeeper:    lctmanagerKeeper,

		Params:                 collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Components:             collections.NewMap(sb, types.ComponentPrefix, "components", collections.StringKey, codec.CollValue[types.Component](cdc)),
		ComponentVerifications: collections.NewMap(sb, types.VerificationPrefix, "verifications", collections.StringKey, codec.CollValue[types.ComponentVerification](cdc)),
		ComponentPairingRules:  collections.NewMap(sb, types.PairingRulesPrefix, "pairing_rules", collections.StringKey, codec.CollValue[types.ComponentPairingRule](cdc)),
		ManufacturerComponents: collections.NewMap(sb, types.ManufacturerComponentKey, "manufacturer_components", collections.StringKey, codec.CollValue[types.Component](cdc)),
		PairingAuthorizations:  collections.NewMap(sb, types.PairingAuthorizationKey, "pairing_authorizations", collections.StringKey, codec.CollValue[types.PairingAuthorization](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// GetParams returns the current module parameters
func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	// Since Params is empty, no validation needed
	return k.Params.Set(ctx, params)
}

// RegisterComponent registers a new component in the system
func (k Keeper) RegisterComponent(ctx context.Context, component types.Component) error {
	// Check if component already exists
	exists, err := k.Components.Has(ctx, component.ComponentId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("component %s already registered", component.ComponentId)
	}

	// Set creation timestamp
	component.CreatedAt = time.Now()
	component.LastVerifiedAt = time.Now()

	// Store component
	if err := k.Components.Set(ctx, component.ComponentId, component); err != nil {
		return err
	}

	// Update manufacturer index (store the component directly)
	return k.ManufacturerComponents.Set(ctx, component.ManufacturerId, component)
}

// GetComponent retrieves a component by ID
func (k Keeper) GetComponent(ctx context.Context, componentId string) (types.Component, error) {
	return k.Components.Get(ctx, componentId)
}

// UpdateComponent updates an existing component
func (k Keeper) UpdateComponent(ctx context.Context, component types.Component) error {
	exists, err := k.Components.Has(ctx, component.ComponentId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("component %s not found", component.ComponentId)
	}

	return k.Components.Set(ctx, component.ComponentId, component)
}

// VerifyComponent updates component verification status
func (k Keeper) VerifyComponent(ctx context.Context, verification types.ComponentVerification) error {
	// Verify component exists
	exists, err := k.Components.Has(ctx, verification.ComponentId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("component %s not found", verification.ComponentId)
	}

	// Store verification
	if err := k.ComponentVerifications.Set(ctx, verification.ComponentId, verification); err != nil {
		return err
	}

	// Update component last verified timestamp
	component, err := k.Components.Get(ctx, verification.ComponentId)
	if err != nil {
		return err
	}
	component.LastVerifiedAt = verification.VerifiedAt
	return k.Components.Set(ctx, verification.ComponentId, component)
}

// GetComponentVerification retrieves verification status for a component
func (k Keeper) GetComponentVerification(ctx context.Context, componentId string) (types.ComponentVerification, error) {
	return k.ComponentVerifications.Get(ctx, componentId)
}

// SetPairingRule sets a pairing rule for component types
func (k Keeper) SetPairingRule(ctx context.Context, rule types.ComponentPairingRule) error {
	ruleKey := fmt.Sprintf("%s-%s", rule.SourceTypeHash, rule.TargetTypeHash)
	return k.ComponentPairingRules.Set(ctx, ruleKey, rule)
}

// GetPairingRule retrieves pairing rules for component types
func (k Keeper) GetPairingRule(ctx context.Context, sourceTypeHash, targetTypeHash string) (types.ComponentPairingRule, error) {
	ruleKey := fmt.Sprintf("%s-%s", sourceTypeHash, targetTypeHash)
	return k.ComponentPairingRules.Get(ctx, ruleKey)
}

// CheckPairingAuthorization verifies if two components can pair with each other
func (k Keeper) CheckPairingAuthorization(ctx context.Context, componentAId, componentBId string) (bool, bool, error) {
	// Use the new trust-based authorization logic
	canPair, _, err := k.ValidateComponentPairing(ctx, componentAId, componentBId)
	if err != nil {
		return false, false, err
	}

	// For backward compatibility, return bidirectional result
	return canPair, canPair, nil
}

// GetManufacturerComponents retrieves all components for a manufacturer
func (k Keeper) GetManufacturerComponents(ctx context.Context, manufacturerId string) (types.Component, error) {
	return k.ManufacturerComponents.Get(ctx, manufacturerId)
}

// UpdateComponentStatus updates a component's status
func (k Keeper) UpdateComponentStatus(ctx context.Context, componentId, status string) error {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return err
	}

	component.Status = status
	return k.Components.Set(ctx, componentId, component)
}

// GetComponentRelationships retrieves all LCT relationships for a component
func (k Keeper) GetComponentRelationships(ctx context.Context, componentId string) ([]string, error) {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return nil, err
	}
	return component.RelationshipIds, nil
}

// AddComponentRelationship adds an LCT relationship to a component
func (k Keeper) AddComponentRelationship(ctx context.Context, componentId, lctId string) error {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return err
	}

	// Check if relationship already exists
	for _, existingLct := range component.RelationshipIds {
		if existingLct == lctId {
			return nil // Already exists
		}
	}

	component.RelationshipIds = append(component.RelationshipIds, lctId)
	return k.Components.Set(ctx, componentId, component)
}

// RemoveComponentRelationship removes an LCT relationship from a component
func (k Keeper) RemoveComponentRelationship(ctx context.Context, componentId, lctId string) error {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return err
	}

	// Find and remove relationship
	newRelationships := make([]string, 0, len(component.RelationshipIds))
	for _, existingLct := range component.RelationshipIds {
		if existingLct != lctId {
			newRelationships = append(newRelationships, existingLct)
		}
	}

	component.RelationshipIds = newRelationships
	return k.Components.Set(ctx, componentId, component)
}

// GetPairingRules retrieves all pairing rules for a component type
func (k Keeper) GetPairingRules(ctx context.Context, componentType string) ([]*types.ComponentPairingRule, error) {
	var rules []*types.ComponentPairingRule

	// Get pairing rules for a component type
	err := k.ComponentPairingRules.Walk(ctx, nil, func(key string, rule types.ComponentPairingRule) (bool, error) {
		if rule.SourceTypeHash == componentType {
			rules = append(rules, &rule)
		}
		return false, nil
	})

	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to walk pairing rules")
	}

	return rules, nil
}

// ListComponents retrieves all components
func (k Keeper) ListComponents(ctx context.Context) ([]*types.Component, error) {
	var components []*types.Component

	err := k.Components.Walk(ctx, nil, func(key string, component types.Component) (bool, error) {
		components = append(components, &component)
		return false, nil
	})

	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to walk components")
	}

	return components, nil
}

// VerifyComponentPairingWithBackend uses the pluggable backend to verify component pairing
func (k Keeper) VerifyComponentPairingWithBackend(ctx context.Context, componentA, componentB string) (bool, string, error) {
	if k.verificationBackend == nil {
		// Fallback to on-chain verification if no backend is configured
		aCanPairB, bCanPairA, err := k.CheckPairingAuthorization(ctx, componentA, componentB)
		if err != nil {
			return false, "", err
		}
		if aCanPairB && bCanPairA {
			return true, "pairing allowed: bidirectional authorization confirmed", nil
		}
		return false, "pairing denied: insufficient authorization", nil
	}

	return k.verificationBackend.VerifyComponentPairing(ctx, componentA, componentB)
}

// GetComponentMetadataFromBackend retrieves component metadata from the verification backend
func (k Keeper) GetComponentMetadataFromBackend(ctx context.Context, componentID string) (map[string]interface{}, error) {
	if k.verificationBackend == nil {
		// Return empty metadata if no backend is configured
		return map[string]interface{}{
			"type":   "unknown",
			"status": "no_backend_configured",
		}, nil
	}

	return k.verificationBackend.GetComponentMetadata(ctx, componentID)
}

// GetComponentIdentity returns a ComponentIdentity for the pairingqueue module interface
func (k Keeper) GetComponentIdentity(ctx context.Context, componentId string) (types.ComponentIdentity, bool) {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return types.ComponentIdentity{}, false
	}

	// Convert Component to ComponentIdentity
	// Note: For privacy-focused implementation, we use the category hash as component type
	componentIdentity := types.ComponentIdentity{
		ComponentId:      component.ComponentId,
		ComponentType:    component.CategoryHash,         // Use category hash instead of component type
		ManufacturerData: component.VerificationMetadata, // Use verification metadata instead of hardware specs
		Status:           component.Status,
		LastSeen:         component.LastVerifiedAt.Unix(),
	}

	return componentIdentity, true
}

// VerifyComponentStatus returns verification status for the pairingqueue module interface
func (k Keeper) VerifyComponentStatus(ctx context.Context, componentId string) (bool, string) {
	component, err := k.Components.Get(ctx, componentId)
	if err != nil {
		return false, "component not found"
	}

	if component.Status != "active" {
		return false, "component not active"
	}

	return true, "component verified"
}

// CheckBidirectionalPairingAuth checks bidirectional pairing authorization for the pairingqueue module interface
func (k Keeper) CheckBidirectionalPairingAuth(ctx context.Context, componentA, componentB string) (bool, bool, string) {
	aCanPairB, bCanPairA, err := k.CheckPairingAuthorization(ctx, componentA, componentB)
	if err != nil {
		return false, false, "error checking pairing authorization"
	}

	message := "pairing authorization checked"
	return aCanPairB, bCanPairA, message
}

// VerifyComponentForPairing provides the interface method expected by the pairing module
func (k Keeper) VerifyComponentForPairing(ctx context.Context, componentId string) (bool, string) {
	return k.VerifyComponentStatus(ctx, componentId)
}

// Privacy-focused methods for anonymous component operations

// GenerateComponentHash creates an anonymous hash for a real component ID
func (k Keeper) GenerateComponentHash(ctx context.Context, realComponentID string) (string, error) {
	if k.verificationBackend == nil {
		return "", fmt.Errorf("verification backend not configured")
	}
	return k.verificationBackend.GenerateComponentHash(ctx, realComponentID)
}

// ResolveComponentHash resolves a component hash to real component data
func (k Keeper) ResolveComponentHash(ctx context.Context, componentHash string) (map[string]interface{}, error) {
	if k.verificationBackend == nil {
		return nil, fmt.Errorf("verification backend not configured")
	}
	return k.verificationBackend.ResolveComponentHash(ctx, componentHash)
}

// VerifyComponentPairingWithHashes verifies pairing using component hashes
func (k Keeper) VerifyComponentPairingWithHashes(ctx context.Context, componentHashA, componentHashB string) (bool, string, error) {
	if k.verificationBackend == nil {
		// Fallback to on-chain verification if no backend is configured
		// Note: This fallback won't work with hashes, so we return an error
		return false, "", fmt.Errorf("verification backend not configured for hash-based verification")
	}
	return k.verificationBackend.VerifyComponentPairingWithHashes(ctx, componentHashA, componentHashB)
}

// GetAnonymousComponentMetadata returns only non-sensitive metadata for a component hash
func (k Keeper) GetAnonymousComponentMetadata(ctx context.Context, componentHash string) (map[string]interface{}, error) {
	if k.verificationBackend == nil {
		// Return minimal anonymous metadata if no backend is configured
		return map[string]interface{}{
			"component_hash": componentHash,
			"status":         "unknown",
			"message":        "verification backend not configured",
		}, nil
	}
	return k.verificationBackend.GetAnonymousComponentMetadata(ctx, componentHash)
}

// RegisterAnonymousComponent registers a component using anonymous hashes
func (k Keeper) RegisterAnonymousComponent(ctx context.Context, realComponentID, manufacturerID, componentType string) (types.Component, error) {
	if k.verificationBackend == nil {
		return types.Component{}, fmt.Errorf("verification backend not configured")
	}

	// Generate hashes for privacy
	componentHash, err := k.verificationBackend.GenerateComponentHash(ctx, realComponentID)
	if err != nil {
		return types.Component{}, fmt.Errorf("failed to generate component hash: %w", err)
	}

	manufacturerHash := k.generateHash(manufacturerID)
	categoryHash := k.generateHash(componentType)
	authorizationRulesHash := k.generateHash("default_rules")

	// Create anonymous component
	component := types.Component{
		ComponentId:            componentHash,
		ManufacturerHash:       manufacturerHash,
		CategoryHash:           categoryHash,
		AuthorizationRulesHash: authorizationRulesHash,
		Status:                 "active",
		RegisteredAt:           time.Now(),
		TrustAnchor:            "cryptographic_trust_anchor",
		LastVerifiedAt:         time.Now(),
		VerificationMetadata:   "anonymous_registration",
		RelationshipHashes:     []string{},
		LctHash:                "",
		EncryptedDeviceKeyHalf: []byte{},
	}

	// Store the component
	if err := k.Components.Set(ctx, componentHash, component); err != nil {
		return types.Component{}, fmt.Errorf("failed to store component: %w", err)
	}

	return component, nil
}

// generateHash creates a SHA-256 hash of the input string
func (k Keeper) generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// CreateAnonymousPairingAuthorization creates an anonymous pairing authorization
func (k Keeper) CreateAnonymousPairingAuthorization(ctx context.Context, componentHashA, componentHashB, ruleHash string) (types.AnonymousPairingAuthorization, error) {
	authID := k.generateHash(fmt.Sprintf("%s-%s-%d", componentHashA, componentHashB, time.Now().Unix()))

	auth := types.AnonymousPairingAuthorization{
		AuthId:                authID,
		ComponentHashA:        componentHashA,
		ComponentHashB:        componentHashB,
		RuleHash:              ruleHash,
		Status:                "active",
		ExpiresAt:             time.Now().AddDate(1, 0, 0), // 1 year from now
		TrustScoreRequirement: "0.7",
		AuthorizationLevel:    "basic",
	}

	// Store the authorization
	authKey := fmt.Sprintf("%s-%s", componentHashA, componentHashB)
	if err := k.PairingAuthorizations.Set(ctx, authKey, types.PairingAuthorization{
		AuthId:              authID,
		ComponentId:         componentHashA,
		AllowedPartnerTypes: componentHashB,
		AllowedPartnerIds:   componentHashB,
		ApplicationContext:  "anonymous_pairing",
		ExpiresAt:           auth.ExpiresAt.Unix(),
		Status:              "active",
		TrustScore:          auth.TrustScoreRequirement,
		AuthorizationLevel:  auth.AuthorizationLevel,
		CreatedAt:           time.Now().Unix(),
		UpdatedAt:           time.Now().Unix(),
		Version:             1,
	}); err != nil {
		return types.AnonymousPairingAuthorization{}, fmt.Errorf("failed to store pairing authorization: %w", err)
	}

	return auth, nil
}

// CreateAnonymousRevocationEvent creates an anonymous revocation event
func (k Keeper) CreateAnonymousRevocationEvent(ctx context.Context, targetHash, revocationType, urgencyLevel, reasonCategory, initiatorHash string) (types.AnonymousRevocationEvent, error) {
	revocationID := k.generateHash(fmt.Sprintf("%s-%s-%d", targetHash, initiatorHash, time.Now().Unix()))

	revocation := types.AnonymousRevocationEvent{
		RevocationId:   revocationID,
		TargetHash:     targetHash,
		RevocationType: revocationType,
		UrgencyLevel:   urgencyLevel,
		EffectiveAt:    time.Now(),
		ReasonCategory: reasonCategory,
		InitiatorHash:  initiatorHash,
	}

	// Store the revocation event
	if err := k.ComponentVerifications.Set(ctx, revocationID, types.ComponentVerification{
		ComponentId:          targetHash,
		Status:               "revoked",
		VerifiedAt:           time.Now(),
		VerificationMethod:   "revocation",
		VerificationEvidence: revocationID,
		Notes:                fmt.Sprintf("Revoked: %s - %s", revocationType, reasonCategory),
	}); err != nil {
		return types.AnonymousRevocationEvent{}, fmt.Errorf("failed to store revocation event: %w", err)
	}

	// Update component status if it's a component revocation
	if revocationType == "INDIVIDUAL" {
		if component, err := k.Components.Get(ctx, targetHash); err == nil {
			component.Status = "revoked"
			if err := k.Components.Set(ctx, targetHash, component); err != nil {
				return types.AnonymousRevocationEvent{}, fmt.Errorf("failed to update component status: %w", err)
			}
		}
	}

	return revocation, nil
}
