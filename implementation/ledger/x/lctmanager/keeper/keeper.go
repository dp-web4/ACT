package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	componentregistrytypes "racecar-web/x/componentregistry/types"
	"racecar-web/x/lctmanager/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.BinaryCodec
	addressCodec address.Codec
	logger       log.Logger

	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]
	// Collections for Web4 LCT relationships
	LinkedContextToken    collections.Map[string, types.LinkedContextToken]
	ComponentRelationship collections.Map[string, types.ComponentRelationship]
	LCTMediatedPairings   collections.Map[string, types.LCTMediatedPairing]
	SessionKeyExchanges   collections.Map[string, types.SessionKeyExchange]
	PairingChallenges     collections.Map[string, types.PairingChallenge]
	SplitKeys             collections.Map[string, types.SplitKey]

	bankKeeper              types.BankKeeper
	componentregistryKeeper componentregistrytypes.ComponentregistryKeeper
	pairingqueueKeeper      types.PairingqueueKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	componentregistryKeeper componentregistrytypes.ComponentregistryKeeper,
	pairingqueueKeeper types.PairingqueueKeeper,
	logger log.Logger,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:            storeService,
		cdc:                     cdc,
		addressCodec:            addressCodec,
		authority:               authority,
		bankKeeper:              bankKeeper,
		componentregistryKeeper: componentregistryKeeper,
		pairingqueueKeeper:      pairingqueueKeeper,
		logger:                  logger,

		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		LinkedContextToken:    collections.NewMap(sb, types.LctKeyPrefix, "linked_context_token", collections.StringKey, codec.CollValue[types.LinkedContextToken](cdc)),
		ComponentRelationship: collections.NewMap(sb, types.RelationshipKeyPrefix, "component_relationship", collections.StringKey, codec.CollValue[types.ComponentRelationship](cdc)),
		LCTMediatedPairings:   collections.NewMap(sb, types.LCTMediatedPairingPrefix, "lct_mediated_pairings", collections.StringKey, codec.CollValue[types.LCTMediatedPairing](cdc)),
		SessionKeyExchanges:   collections.NewMap(sb, types.SessionKeyExchangePrefix, "session_key_exchanges", collections.StringKey, codec.CollValue[types.SessionKeyExchange](cdc)),
		PairingChallenges:     collections.NewMap(sb, types.PairingChallengePrefix, "pairing_challenges", collections.StringKey, codec.CollValue[types.PairingChallenge](cdc)),
		SplitKeys:             collections.NewMap(sb, types.SplitKeyPrefix, "split_keys", collections.StringKey, codec.CollValue[types.SplitKey](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// Params get all parameters as types.Params
func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

// SetParams set the params
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.Params.Set(ctx, params)
}

// Web4 Interface Methods - Required by other modules

// GetLinkedContextToken retrieves LCT relationship information
func (k Keeper) GetLinkedContextToken(ctx context.Context, lctId string) (types.LinkedContextToken, bool) {
	lct, err := k.LinkedContextToken.Get(ctx, lctId)
	if err != nil {
		return types.LinkedContextToken{}, false
	}
	return lct, true
}

// SetLinkedContextToken stores LCT relationship information
func (k Keeper) SetLinkedContextToken(ctx context.Context, lct types.LinkedContextToken) error {
	return k.LinkedContextToken.Set(ctx, lct.LctId, lct)
}

// CreateLCTRelationship creates a new LCT representing the relationship between two components
func (k Keeper) CreateLCTRelationship(ctx context.Context, componentA, componentB, operationalContext, proxyId string) (string, string, error) {
	// Generate unique LCT ID for this relationship
	lctId := k.generateLCTId(componentA, componentB)

	// Generate split-key pair for this relationship
	lctKeyHalf, deviceKeyHalf, err := k.generateSplitKeyPair()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate split keys: %w", err)
	}

	// Convert key halves to hex strings for storage
	lctKeyHalfHex := fmt.Sprintf("%x", lctKeyHalf[:])
	deviceKeyHalfHex := fmt.Sprintf("%x", deviceKeyHalf[:])

	// Create the LCT relationship
	lct := types.LinkedContextToken{
		LctId:              lctId,
		ComponentAId:       componentA,
		ComponentBId:       componentB,
		LctKeyHalf:         lctKeyHalfHex,
		PairingStatus:      "active",
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		LastContactAt:      time.Now().Unix(),
		TrustAnchor:        k.generateTrustAnchor(componentA, componentB),
		OperationalContext: operationalContext,
		ProxyComponentId:   proxyId,
		AuthorizationRules: "", // Will be populated from component registry
	}

	// Store the LCT relationship
	err = k.SetLinkedContextToken(ctx, lct)
	if err != nil {
		return "", "", fmt.Errorf("failed to store LCT relationship: %w", err)
	}

	// Update component relationship tracking for both components
	err = k.updateComponentRelationships(ctx, componentA, lctId, "paired")
	if err != nil {
		return "", "", fmt.Errorf("failed to update component A relationships: %w", err)
	}

	err = k.updateComponentRelationships(ctx, componentB, lctId, "paired")
	if err != nil {
		return "", "", fmt.Errorf("failed to update component B relationships: %w", err)
	}

	return lctId, deviceKeyHalfHex, nil
}

// GetComponentRelationships returns all LCT relationships for a component (many-to-many support)
func (k Keeper) GetComponentRelationships(ctx context.Context, componentId string) ([]types.LinkedContextToken, error) {
	var relationships []types.LinkedContextToken

	err := k.LinkedContextToken.Walk(ctx, nil, func(key string, lct types.LinkedContextToken) (bool, error) {
		if lct.ComponentAId == componentId || lct.ComponentBId == componentId {
			relationships = append(relationships, lct)
		}
		return false, nil
	})

	return relationships, err
}

// TerminateLCTRelationship ends an LCT relationship between components
func (k Keeper) TerminateLCTRelationship(ctx context.Context, lctId, reason string, notifyOffline bool) error {
	lct, err := k.LinkedContextToken.Get(ctx, lctId)
	if err != nil {
		return fmt.Errorf("LCT relationship not found: %s", lctId)
	}

	// Update LCT status
	lct.PairingStatus = types.StatusTerminated
	lct.UpdatedAt = time.Now().Unix()

	err = k.SetLinkedContextToken(ctx, lct)
	if err != nil {
		return fmt.Errorf("failed to update LCT status: %w", err)
	}

	// Queue offline operations if requested and pairingqueueKeeper is available
	if notifyOffline && k.pairingqueueKeeper != nil {
		_, _ = k.pairingqueueKeeper.QueueOfflineOperation(ctx, lct.ComponentAId, "termination_notification")
		_, _ = k.pairingqueueKeeper.QueueOfflineOperation(ctx, lct.ComponentBId, "termination_notification")
	}

	return nil
}

// Helper Methods

// generateLCTId creates a unique identifier for the LCT relationship
func (k Keeper) generateLCTId(componentA, componentB string) string {
	return fmt.Sprintf("lct-%s-%s-%d", componentA, componentB, time.Now().UnixNano())
}

// generateSplitKeyPair creates the split key pair for Web4 cryptography using our crypto primitives
func (k Keeper) generateSplitKeyPair() ([32]byte, [32]byte, error) {
	// Generate two 32-byte key shares using our crypto primitives
	keyShareA, err := GenerateKeyShare()
	if err != nil {
		return [32]byte{}, [32]byte{}, fmt.Errorf("failed to generate key share A: %w", err)
	}

	keyShareB, err := GenerateKeyShare()
	if err != nil {
		return [32]byte{}, [32]byte{}, fmt.Errorf("failed to generate key share B: %w", err)
	}

	return keyShareA, keyShareB, nil
}

// combineKeyShares combines two key shares to create the encryption key
func (k Keeper) combineKeyShares(keyShareA, keyShareB [32]byte) [32]byte {
	// For split-key cryptography, we XOR the two shares
	var combinedKey [32]byte
	for i := 0; i < 32; i++ {
		combinedKey[i] = keyShareA[i] ^ keyShareB[i]
	}
	return combinedKey
}

// encryptWithSplitKey encrypts data using a combined key from two shares
func (k Keeper) encryptWithSplitKey(keyShareA, keyShareB [32]byte, plaintext []byte) ([]byte, error) {
	combinedKey := k.combineKeyShares(keyShareA, keyShareB)
	return EncryptWithKey(combinedKey, plaintext)
}

// decryptWithSplitKey decrypts data using a combined key from two shares
func (k Keeper) decryptWithSplitKey(keyShareA, keyShareB [32]byte, ciphertext []byte) ([]byte, error) {
	combinedKey := k.combineKeyShares(keyShareA, keyShareB)
	return DecryptWithKey(combinedKey, ciphertext)
}

// GetSplitKey retrieves split key metadata by ID
func (k Keeper) GetSplitKey(ctx context.Context, splitKeyId string) (types.SplitKey, bool) {
	splitKey, err := k.SplitKeys.Get(ctx, splitKeyId)
	if err != nil {
		return types.SplitKey{}, false
	}
	return splitKey, true
}

// SetSplitKey stores split key metadata
func (k Keeper) SetSplitKey(ctx context.Context, splitKey types.SplitKey) error {
	return k.SplitKeys.Set(ctx, splitKey.LctId, splitKey)
}

// UpdateSplitKeyStatus updates the status of a split key
func (k Keeper) UpdateSplitKeyStatus(ctx context.Context, splitKeyId, newStatus string) error {
	splitKey, err := k.SplitKeys.Get(ctx, splitKeyId)
	if err != nil {
		return fmt.Errorf("split key not found: %s", splitKeyId)
	}

	splitKey.Status = newStatus
	if newStatus == "active" {
		splitKey.ActivatedAt = time.Now().Unix()
	}

	return k.SplitKeys.Set(ctx, splitKeyId, splitKey)
}

// Cryptographic Operations for LCT Manager

// EncryptMessageForLCT encrypts a message for secure communication between LCT components
func (k Keeper) EncryptMessageForLCT(ctx context.Context, lctId string, message []byte) ([]byte, error) {
	// Get the LCT relationship
	lct, found := k.GetLinkedContextToken(ctx, lctId)
	if !found {
		return nil, fmt.Errorf("LCT not found: %s", lctId)
	}

	// Check if LCT is active
	if lct.PairingStatus != "active" {
		return nil, fmt.Errorf("LCT is not active: %s", lctId)
	}

	// Generate ephemeral key shares for this message
	keyShareA, keyShareB, err := k.generateSplitKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message keys: %w", err)
	}

	// Encrypt the message using split-key cryptography
	encryptedMessage, err := k.encryptWithSplitKey(keyShareA, keyShareB, message)
	if err != nil {
		// Clean up keys on error
		ZeroKey(&keyShareA)
		ZeroKey(&keyShareB)
		return nil, fmt.Errorf("failed to encrypt message: %w", err)
	}

	// Securely zero the key shares
	ZeroKey(&keyShareA)
	ZeroKey(&keyShareB)

	return encryptedMessage, nil
}

// DecryptMessageForLCT decrypts a message for secure communication between LCT components
func (k Keeper) DecryptMessageForLCT(ctx context.Context, lctId string, encryptedMessage []byte) ([]byte, error) {
	// Get the LCT relationship
	lct, found := k.GetLinkedContextToken(ctx, lctId)
	if !found {
		return nil, fmt.Errorf("LCT not found: %s", lctId)
	}

	// Check if LCT is active
	if lct.PairingStatus != "active" {
		return nil, fmt.Errorf("LCT is not active: %s", lctId)
	}

	// In a real implementation, the key shares would be provided by the components
	// For now, we'll simulate the decryption process
	// This would typically involve the components providing their key shares

	// For demonstration, we'll return an error indicating that key shares are needed
	return nil, fmt.Errorf("key shares required for decryption - this should be handled by components")
}

// ValidateLCTCryptographicIntegrity validates the cryptographic integrity of an LCT relationship
func (k Keeper) ValidateLCTCryptographicIntegrity(ctx context.Context, lctId string) (bool, error) {
	// Get the LCT relationship
	lct, found := k.GetLinkedContextToken(ctx, lctId)
	if !found {
		return false, fmt.Errorf("LCT not found: %s", lctId)
	}

	// Check if LCT has required fields
	if lct.LctId == "" || lct.ComponentAId == "" || lct.ComponentBId == "" {
		return false, fmt.Errorf("LCT missing required fields: %s", lctId)
	}

	// Check if LCT is in a valid state for cryptographic operations
	if lct.PairingStatus != "active" && lct.PairingStatus != "pending" {
		return false, fmt.Errorf("LCT not in valid state for crypto operations: %s", lctId)
	}

	// Check if LCT has a trust anchor
	if lct.TrustAnchor == "" {
		return false, fmt.Errorf("LCT missing trust anchor: %s", lctId)
	}

	return true, nil
}

// GenerateLCTChallenge creates a cryptographic challenge for LCT authentication
func (k Keeper) GenerateLCTChallenge(ctx context.Context, lctId string) ([]byte, error) {
	// Validate LCT integrity
	valid, err := k.ValidateLCTCryptographicIntegrity(ctx, lctId)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("LCT failed integrity validation: %s", lctId)
	}

	// Generate a cryptographically secure challenge
	challengeKey, err := GenerateKeyShare()
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge key: %w", err)
	}

	// Create challenge data with LCT context
	challengeData := append(challengeKey[:], []byte(fmt.Sprintf("-lct-%s-%d", lctId, time.Now().Unix()))...)

	// Securely zero the challenge key
	ZeroKey(&challengeKey)

	return challengeData, nil
}

// VerifyLCTChallengeResponse verifies a response to an LCT challenge
func (k Keeper) VerifyLCTChallengeResponse(ctx context.Context, lctId string, challenge []byte, response []byte) (bool, error) {
	// Validate LCT integrity
	valid, err := k.ValidateLCTCryptographicIntegrity(ctx, lctId)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, fmt.Errorf("LCT failed integrity validation: %s", lctId)
	}

	// Hash the challenge to get expected response
	expectedResponse := sha256.Sum256(challenge)

	// Compare the provided response with the expected response
	if len(response) != len(expectedResponse) {
		return false, nil
	}

	for i := 0; i < len(response); i++ {
		if response[i] != expectedResponse[i] {
			return false, nil
		}
	}

	return true, nil
}

// generateTrustAnchor creates the genesis trust binding for the relationship
func (k Keeper) generateTrustAnchor(componentA, componentB string) string {
	return fmt.Sprintf("trust-anchor-%s-%s", componentA, componentB)
}

// updateComponentRelationships maintains many-to-many relationship tracking
func (k Keeper) updateComponentRelationships(ctx context.Context, componentId, lctId, status string) error {
	relationshipId := fmt.Sprintf("rel-%s", componentId)

	relationship, err := k.ComponentRelationship.Get(ctx, relationshipId)
	if err != nil {
		// Create new relationship tracking
		relationship = types.ComponentRelationship{
			RelationshipId:   relationshipId,
			ComponentId:      componentId,
			RelatedLcts:      lctId,
			RelationshipType: "pairing",
			Status:           status,
		}
	} else {
		// Update existing relationship tracking
		if relationship.RelatedLcts == "" {
			relationship.RelatedLcts = lctId
		} else {
			relationship.RelatedLcts = relationship.RelatedLcts + "," + lctId
		}
		relationship.Status = status
	}

	return k.SetComponentRelationship(ctx, relationship)
}

// GetComponentRelationship retrieves component relationship tracking
func (k Keeper) GetComponentRelationship(ctx context.Context, relationshipId string) (types.ComponentRelationship, bool) {
	relationship, err := k.ComponentRelationship.Get(ctx, relationshipId)
	if err != nil {
		return types.ComponentRelationship{}, false
	}
	return relationship, true
}

// SetComponentRelationship stores component relationship tracking
func (k Keeper) SetComponentRelationship(ctx context.Context, relationship types.ComponentRelationship) error {
	return k.ComponentRelationship.Set(ctx, relationship.RelationshipId, relationship)
}

// CreateLctRelationship creates a new LCT relationship between two components
func (k Keeper) CreateLctRelationship(ctx context.Context, creator sdk.AccAddress, componentA, componentB, context, proxyID string) (*types.LinkedContextToken, error) {
	// Generate LCT ID
	lctID := fmt.Sprintf("lct_%s_%s_%d", componentA, componentB, time.Now().UnixNano())

	// Create LCT
	lct := types.LinkedContextToken{
		LctId:              lctID,
		ComponentAId:       componentA,
		ComponentBId:       componentB,
		LctKeyHalf:         "", // To be set during key exchange
		PairingStatus:      types.StatusPending,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
		LastContactAt:      time.Now().Unix(),
		TrustAnchor:        creator.String(),
		OperationalContext: context,
		ProxyComponentId:   proxyID,
		AuthorizationRules: "{}", // Default empty rules
	}

	// Store LCT
	if err := k.LinkedContextToken.Set(ctx, lctID, lct); err != nil {
		return nil, err
	}

	// Update component relationships
	if err := k.updateComponentRelationships(ctx, componentA, lctID, types.StatusActive); err != nil {
		return nil, err
	}
	if err := k.updateComponentRelationships(ctx, componentB, lctID, types.StatusActive); err != nil {
		return nil, err
	}

	return &lct, nil
}

// GetLct retrieves an LCT by ID
func (k Keeper) GetLct(ctx context.Context, lctID string) (types.LinkedContextToken, bool) {
	lct, err := k.LinkedContextToken.Get(ctx, lctID)
	if err != nil {
		return types.LinkedContextToken{}, false
	}
	return lct, true
}

// UpdateLctStatus updates the status of an LCT
func (k Keeper) UpdateLctStatus(ctx context.Context, lctID, newStatus, reason string) error {
	lct, err := k.LinkedContextToken.Get(ctx, lctID)
	if err != nil {
		return types.ErrLctNotFound
	}

	// Validate the new status
	if !types.IsValidLCTStatus(newStatus) {
		return types.ErrInvalidLctStatus
	}

	lct.PairingStatus = newStatus
	lct.UpdatedAt = time.Now().Unix()

	if err := k.LinkedContextToken.Set(ctx, lctID, lct); err != nil {
		return err
	}

	// If terminating, update relationships
	if newStatus == types.StatusTerminated {
		if err := k.removeLctFromRelationships(ctx, lctID); err != nil {
			return err
		}
	}

	return nil
}

// ValidateLctAccess checks if a component has access to an LCT
func (k Keeper) ValidateLctAccess(ctx context.Context, lctID, requestorID string) (bool, string, error) {
	lct, err := k.LinkedContextToken.Get(ctx, lctID)
	if err != nil {
		return false, "", types.ErrLctNotFound
	}

	// Check if requestor is one of the components
	if requestorID != lct.ComponentAId && requestorID != lct.ComponentBId {
		return false, "", nil
	}

	// Check LCT status
	if lct.PairingStatus != types.StatusActive {
		return false, "", nil
	}

	// Parse authorization rules
	var rules map[string]interface{}
	if err := json.Unmarshal([]byte(lct.AuthorizationRules), &rules); err != nil {
		return false, "", err
	}

	// Get component's role in the LCT
	role := "component_a"
	if requestorID == lct.ComponentBId {
		role = "component_b"
	}

	// Check if component has explicit access rules
	if componentRules, ok := rules[role].(map[string]interface{}); ok {
		// Check if access is explicitly denied
		if denied, ok := componentRules["denied"].(bool); ok && denied {
			return false, "", nil
		}

		// Check access level
		if accessLevel, ok := componentRules["access_level"].(string); ok {
			return true, accessLevel, nil
		}
	}

	// Check default access rules
	if defaultRules, ok := rules["default"].(map[string]interface{}); ok {
		if accessLevel, ok := defaultRules["access_level"].(string); ok {
			return true, accessLevel, nil
		}
	}

	// If no specific rules are found, check operational context
	if lct.OperationalContext != "" {
		// Components in the same operational context get standard access
		return true, "standard", nil
	}

	// Default to restricted access if no rules are specified
	return true, "restricted", nil
}

// Helper function to remove LCT from relationships
func (k Keeper) removeLctFromRelationships(ctx context.Context, lctID string) error {
	lct, err := k.LinkedContextToken.Get(ctx, lctID)
	if err != nil {
		return types.ErrLctNotFound
	}

	// Remove from both components' relationships
	for _, componentID := range []string{lct.ComponentAId, lct.ComponentBId} {
		rel, err := k.ComponentRelationship.Get(ctx, componentID)
		if err != nil {
			continue
		}

		var lctIDs []string
		if err := json.Unmarshal([]byte(rel.RelatedLcts), &lctIDs); err != nil {
			return err
		}

		// Remove LCT from list
		newLctIDs := make([]string, 0, len(lctIDs))
		for _, id := range lctIDs {
			if id != lctID {
				newLctIDs = append(newLctIDs, id)
			}
		}

		// Update relationships
		lctIDsJSON, err := json.Marshal(newLctIDs)
		if err != nil {
			return err
		}
		rel.RelatedLcts = string(lctIDsJSON)

		if err := k.ComponentRelationship.Set(ctx, componentID, rel); err != nil {
			return err
		}
	}

	return nil
}
