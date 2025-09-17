package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"racecar-web/x/mrh/types"
)

// AddWitnessRelationship records a witness relationship between LCTs
func (k Keeper) AddWitnessRelationship(
	ctx context.Context,
	witnessLCT string,
	subjectLCT string,
	eventType string,
	signature []byte,
) error {
	// Create witness relationship
	witness := &types.WitnessRelationship{
		WitnessLCT: witnessLCT,
		SubjectLCT: subjectLCT,
		EventType:  eventType,
		Timestamp:  time.Now().Unix(),
		Signature:  signature,
		TrustBoost: k.calculateTrustBoost(eventType),
	}
	
	// Store witness relationship
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	key := types.GetWitnessRelationshipKey(witnessLCT, subjectLCT)
	
	// Serialize witness relationship
	data, err := json.Marshal(witness)
	if err != nil {
		return fmt.Errorf("failed to serialize witness relationship: %w", err)
	}
	
	if err := store.Set(key, data); err != nil {
		return fmt.Errorf("failed to store witness relationship: %w", err)
	}
	
	// Update MRH graphs for both LCTs
	// Add witness relationship to subject's graph
	if err := k.AddRelationship(ctx, subjectLCT, "web4:witnessedBy", witnessLCT, witness.TrustBoost); err != nil {
		return fmt.Errorf("failed to add witness relationship to subject graph: %w", err)
	}
	
	// Add witness relationship to witness's graph
	if err := k.AddRelationship(ctx, witnessLCT, "web4:witnessed", subjectLCT, witness.TrustBoost); err != nil {
		return fmt.Errorf("failed to add witness relationship to witness graph: %w", err)
	}
	
	// Emit event
	k.emitWitnessEvent(sdkCtx, witness)
	
	return nil
}

// GetWitnesses returns all witnesses for a subject LCT
func (k Keeper) GetWitnesses(ctx context.Context, subjectLCT string) ([]*types.WitnessRelationship, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	
	var witnesses []*types.WitnessRelationship
	
	// Iterate through all witness relationships
	// This is inefficient - in production, use an index
	iterator := store.Iterator(types.WitnessRelationshipPrefix, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var witness types.WitnessRelationship
		if err := json.Unmarshal(iterator.Value(), &witness); err != nil {
			continue
		}
		
		if witness.SubjectLCT == subjectLCT {
			witnesses = append(witnesses, &witness)
		}
	}
	
	return witnesses, nil
}

// GetWitnessedEntities returns all entities witnessed by a specific LCT
func (k Keeper) GetWitnessedEntities(ctx context.Context, witnessLCT string) ([]*types.WitnessRelationship, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	
	var witnessed []*types.WitnessRelationship
	
	// Iterate through witness relationships
	iterator := store.Iterator(types.WitnessRelationshipPrefix, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var witness types.WitnessRelationship
		if err := json.Unmarshal(iterator.Value(), &witness); err != nil {
			continue
		}
		
		if witness.WitnessLCT == witnessLCT {
			witnessed = append(witnessed, &witness)
		}
	}
	
	return witnessed, nil
}

// VerifyWitnessThreshold checks if an LCT has enough witnesses
func (k Keeper) VerifyWitnessThreshold(
	ctx context.Context,
	subjectLCT string,
	requiredWitnesses int,
	eventType string,
) (bool, error) {
	witnesses, err := k.GetWitnesses(ctx, subjectLCT)
	if err != nil {
		return false, fmt.Errorf("failed to get witnesses: %w", err)
	}
	
	// Count witnesses for specific event type
	count := 0
	for _, w := range witnesses {
		if eventType == "" || w.EventType == eventType {
			count++
		}
	}
	
	return count >= requiredWitnesses, nil
}

// CalculateWitnessedTrust calculates trust score with witness boost
func (k Keeper) CalculateWitnessedTrust(
	ctx context.Context,
	fromLCT string,
	toLCT string,
) (float64, error) {
	// Get base trust from path
	_, baseTrust, err := k.CalculateTrustPath(ctx, fromLCT, toLCT, 6)
	if err != nil {
		return 0, err
	}
	
	// Get witnesses for the target LCT
	witnesses, err := k.GetWitnesses(ctx, toLCT)
	if err != nil {
		return baseTrust, nil // Return base trust if can't get witnesses
	}
	
	// Calculate witness boost
	totalBoost := 0.0
	for _, w := range witnesses {
		// Check if witness is trusted by the source
		_, witnessТrust, err := k.CalculateTrustPath(ctx, fromLCT, w.WitnessLCT, 3)
		if err == nil && witnessТrust > 0.5 {
			// Witness is trusted, apply their boost
			totalBoost += w.TrustBoost * witnessТrust
		}
	}
	
	// Apply witness boost with diminishing returns
	boostedTrust := baseTrust + (1-baseTrust)*totalBoost*0.1 // Max 10% boost per witness
	
	// Cap at 1.0
	if boostedTrust > 1.0 {
		boostedTrust = 1.0
	}
	
	return boostedTrust, nil
}

// CreateBirthCertificate creates a birth certificate with witnesses
func (k Keeper) CreateBirthCertificate(
	ctx context.Context,
	lctID string,
	societyID string,
	entityType string,
	witnesses []string,
	metadata map[string]string,
) (*BirthCertificate, error) {
	// Verify minimum witnesses (typically 3)
	if len(witnesses) < 3 {
		return nil, fmt.Errorf("insufficient witnesses: need at least 3, got %d", len(witnesses))
	}
	
	// Create birth certificate
	cert := &BirthCertificate{
		CertID:     generateCertID(),
		LctID:      lctID,
		SocietyID:  societyID,
		EntityType: entityType,
		IssuedAt:   time.Now(),
		Witnesses:  witnesses,
		Metadata:   metadata,
	}
	
	// Store birth certificate
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	
	certKey := getBirthCertKey(cert.CertID)
	certData, err := json.Marshal(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize birth certificate: %w", err)
	}
	
	if err := store.Set(certKey, certData); err != nil {
		return nil, fmt.Errorf("failed to store birth certificate: %w", err)
	}
	
	// Map LCT to birth certificate
	lctCertKey := getLCTBirthCertKey(lctID)
	if err := store.Set(lctCertKey, []byte(cert.CertID)); err != nil {
		return nil, fmt.Errorf("failed to store LCT-cert mapping: %w", err)
	}
	
	// Record witness relationships
	for _, witnessLCT := range witnesses {
		if err := k.AddWitnessRelationship(ctx, witnessLCT, lctID, "birth_certificate", nil); err != nil {
			// Log error but don't fail
			k.Logger().Error("failed to record witness relationship", "witness", witnessLCT, "error", err)
		}
	}
	
	// Emit event
	k.emitBirthCertEvent(sdkCtx, cert)
	
	return cert, nil
}

// GetBirthCertificate retrieves a birth certificate for an LCT
func (k Keeper) GetBirthCertificate(ctx context.Context, lctID string) (*BirthCertificate, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	
	// Get certificate ID from mapping
	lctCertKey := getLCTBirthCertKey(lctID)
	certIDBytes, err := store.Get(lctCertKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get birth certificate mapping: %w", err)
	}
	if certIDBytes == nil {
		return nil, fmt.Errorf("no birth certificate found for LCT %s", lctID)
	}
	
	// Get certificate
	certKey := getBirthCertKey(string(certIDBytes))
	certData, err := store.Get(certKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get birth certificate: %w", err)
	}
	if certData == nil {
		return nil, fmt.Errorf("birth certificate not found: %s", string(certIDBytes))
	}
	
	var cert BirthCertificate
	if err := json.Unmarshal(certData, &cert); err != nil {
		return nil, fmt.Errorf("failed to deserialize birth certificate: %w", err)
	}
	
	return &cert, nil
}

// Helper functions

func (k Keeper) calculateTrustBoost(eventType string) float64 {
	// Different event types provide different trust boosts
	switch eventType {
	case "birth_certificate":
		return 0.5 // High boost for birth certificate witnessing
	case "pairing":
		return 0.3 // Medium boost for pairing
	case "transaction":
		return 0.1 // Low boost for regular transactions
	default:
		return 0.05 // Minimal boost for other events
	}
}

func (k Keeper) emitWitnessEvent(ctx sdk.Context, witness *types.WitnessRelationship) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"witness_relationship_added",
			sdk.NewAttribute("witness_lct", witness.WitnessLCT),
			sdk.NewAttribute("subject_lct", witness.SubjectLCT),
			sdk.NewAttribute("event_type", witness.EventType),
			sdk.NewAttribute("timestamp", fmt.Sprintf("%d", witness.Timestamp)),
		),
	)
}

func (k Keeper) emitBirthCertEvent(ctx sdk.Context, cert *BirthCertificate) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"birth_certificate_issued",
			sdk.NewAttribute("cert_id", cert.CertID),
			sdk.NewAttribute("lct_id", cert.LctID),
			sdk.NewAttribute("society_id", cert.SocietyID),
			sdk.NewAttribute("entity_type", cert.EntityType),
			sdk.NewAttribute("witness_count", fmt.Sprintf("%d", len(cert.Witnesses))),
		),
	)
}

// BirthCertificate type (should be in types package)
type BirthCertificate struct {
	CertID     string            `json:"cert_id"`
	LctID      string            `json:"lct_id"`
	SocietyID  string            `json:"society_id"`
	EntityType string            `json:"entity_type"`
	IssuedAt   time.Time         `json:"issued_at"`
	Witnesses  []string          `json:"witnesses"`
	Metadata   map[string]string `json:"metadata"`
	Signature  []byte            `json:"signature,omitempty"`
}

// Storage key helpers
func getBirthCertKey(certID string) []byte {
	return append([]byte("birth_cert:"), []byte(certID)...)
}

func getLCTBirthCertKey(lctID string) []byte {
	return append([]byte("lct_birth_cert:"), []byte(lctID)...)
}

func generateCertID() string {
	return fmt.Sprintf("cert:%d", time.Now().UnixNano())
}