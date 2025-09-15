package types

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"
)

// Web4LCT represents a fully Web4-compliant Linked Context Token
type Web4LCT struct {
	// Core Identity
	ID         string `json:"id"`          // lct:web4:act:<entity_type>:<uuid>
	EntityType string `json:"entity_type"` // human|agent|dictionary|society
	PublicKey  []byte `json:"public_key"`  // Ed25519 public key
	
	// MRH Context
	MRHGraphHash string `json:"mrh_graph_hash"` // Content hash of RDF graph
	ContextDepth uint32 `json:"context_depth"`  // Fractal depth
	
	// Trust Tensors (T3/V3)
	T3Competence    float64 `json:"t3_competence"`    // 0.0 to 1.0
	T3Reliability   float64 `json:"t3_reliability"`   // 0.0 to 1.0
	T3Transparency  float64 `json:"t3_transparency"`  // 0.0 to 1.0
	V3Value         float64 `json:"v3_value"`         // Value tensor
	TrustMass       float64 `json:"trust_mass"`       // Calculated from T3
	TrustRadius     float64 `json:"trust_radius"`     // Influence distance
	
	// Relationships
	ParentLCT      string   `json:"parent_lct,omitempty"`      // For agents, the human/entity that created them
	WitnessLCTs    []string `json:"witness_lcts"`               // LCT IDs that have witnessed this entity
	PairedLCTs     []string `json:"paired_lcts,omitempty"`     // For mutual relationships
	
	// Society and Governance
	SocietyID      string   `json:"society_id"`                 // Society this LCT belongs to
	BirthCertHash  string   `json:"birth_cert_hash,omitempty"` // Hash of birth certificate
	AuthorityLevel string   `json:"authority_level,omitempty"` // Level of authority in society
	
	// Lifecycle
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
	Status         string     `json:"status"` // active|suspended|revoked
	
	// Agency (for agents)
	AgencyProof    []byte   `json:"agency_proof,omitempty"`    // Proof of delegated agency
	Permissions    []string `json:"permissions,omitempty"`     // Granted permissions
	Constraints    string   `json:"constraints,omitempty"`     // JSON-encoded constraints
	
	// Metadata
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// EntityType constants
const (
	EntityTypeHuman      = "human"
	EntityTypeAgent      = "agent"
	EntityTypeDictionary = "dictionary"
	EntityTypeSociety    = "society"
)

// Status constants
const (
	StatusActive    = "active"
	StatusSuspended = "suspended"
	StatusRevoked   = "revoked"
)

// Validate checks if the LCT is valid according to Web4 spec
func (lct *Web4LCT) Validate() error {
	// Validate ID format
	if !isValidLCTID(lct.ID) {
		return fmt.Errorf("invalid LCT ID format: %s", lct.ID)
	}
	
	// Validate entity type
	if !isValidEntityType(lct.EntityType) {
		return fmt.Errorf("invalid entity type: %s", lct.EntityType)
	}
	
	// Validate public key
	if len(lct.PublicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid Ed25519 public key size: %d", len(lct.PublicKey))
	}
	
	// Validate trust scores
	if !isValidTrustScore(lct.T3Competence) || 
	   !isValidTrustScore(lct.T3Reliability) || 
	   !isValidTrustScore(lct.T3Transparency) {
		return fmt.Errorf("trust scores must be between 0.0 and 1.0")
	}
	
	// Validate agent-specific fields
	if lct.EntityType == EntityTypeAgent {
		if lct.ParentLCT == "" {
			return fmt.Errorf("agent LCT must have a parent LCT")
		}
		if len(lct.AgencyProof) == 0 {
			return fmt.Errorf("agent LCT must have agency proof")
		}
	}
	
	// Validate status
	if !isValidStatus(lct.Status) {
		return fmt.Errorf("invalid status: %s", lct.Status)
	}
	
	return nil
}

// CalculateTrustMass computes the trust mass from T3 scores
func (lct *Web4LCT) CalculateTrustMass() {
	// Trust mass is the geometric mean of the three trust dimensions
	// This gives a balanced measure that penalizes low scores in any dimension
	lct.TrustMass = cbrt(lct.T3Competence * lct.T3Reliability * lct.T3Transparency)
	
	// Trust radius is proportional to trust mass
	// Higher trust = larger influence radius
	lct.TrustRadius = lct.TrustMass * 10.0 // Max radius of 10 for perfect trust
}

// IsActive returns true if the LCT is active
func (lct *Web4LCT) IsActive() bool {
	return lct.Status == StatusActive && lct.RevokedAt == nil
}

// CanDelegate returns true if this LCT can create agent LCTs
func (lct *Web4LCT) CanDelegate() bool {
	return lct.IsActive() && 
	       (lct.EntityType == EntityTypeHuman || lct.EntityType == EntityTypeSociety)
}

// GetPublicKeyHex returns the public key as a hex string
func (lct *Web4LCT) GetPublicKeyHex() string {
	return hex.EncodeToString(lct.PublicKey)
}

// BirthCertificate represents a society-issued birth certificate for an LCT
type BirthCertificate struct {
	CertID      string            `json:"cert_id"`
	LctID       string            `json:"lct_id"`
	SocietyID   string            `json:"society_id"`
	EntityType  string            `json:"entity_type"`
	IssuedAt    time.Time         `json:"issued_at"`
	Witnesses   []string          `json:"witnesses"`    // LCT IDs of witnesses
	Metadata    map[string]string `json:"metadata"`
	Signature   []byte            `json:"signature"`    // Society's signature
}

// AgencyDelegation represents a delegation of authority from one LCT to another
type AgencyDelegation struct {
	DelegationID string    `json:"delegation_id"`
	FromLCT      string    `json:"from_lct"`      // Delegator
	ToLCT        string    `json:"to_lct"`        // Agent
	
	// Permissions and constraints
	Permissions  []string               `json:"permissions"`
	Constraints  map[string]interface{} `json:"constraints"`
	
	// Temporal validity
	ValidFrom    time.Time  `json:"valid_from"`
	ValidUntil   time.Time  `json:"valid_until"`
	
	// Revocation
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	RevocationReason string `json:"revocation_reason,omitempty"`
	
	// Cryptographic proof
	ProofOfAgency []byte    `json:"proof_of_agency"`
	Signature     []byte    `json:"signature"`
}

// LCTRelationship represents a relationship between two LCTs
type LCTRelationship struct {
	RelationshipID string    `json:"relationship_id"`
	LCT1          string    `json:"lct1"`
	LCT2          string    `json:"lct2"`
	RelationType  string    `json:"relation_type"` // parent-child, peer, witness, etc.
	
	// Bidirectional trust scores
	Trust1to2     float64   `json:"trust_1to2"`
	Trust2to1     float64   `json:"trust_2to1"`
	
	// Metadata
	EstablishedAt time.Time `json:"established_at"`
	LastUpdated   time.Time `json:"last_updated"`
	Metadata      map[string]string `json:"metadata"`
}

// Helper functions

func isValidLCTID(id string) bool {
	// Format: lct:web4:act:<entity_type>:<uuid>
	// This is a simplified check
	return len(id) > 0 && id[:4] == "lct:"
}

func isValidEntityType(entityType string) bool {
	switch entityType {
	case EntityTypeHuman, EntityTypeAgent, EntityTypeDictionary, EntityTypeSociety:
		return true
	default:
		return false
	}
}

func isValidTrustScore(score float64) bool {
	return score >= 0.0 && score <= 1.0
}

func isValidStatus(status string) bool {
	switch status {
	case StatusActive, StatusSuspended, StatusRevoked:
		return true
	default:
		return false
	}
}

// Cube root for trust mass calculation
func cbrt(x float64) float64 {
	if x < 0 {
		return -cbrt(-x)
	}
	return pow(x, 1.0/3.0)
}

// Simple power function
func pow(base, exp float64) float64 {
	// This is a placeholder - use math.Pow in production
	result := 1.0
	if exp == 1.0/3.0 {
		// Approximate cube root using Newton's method
		guess := base / 3.0
		for i := 0; i < 10; i++ {
			guess = (2*guess + base/(guess*guess)) / 3.0
		}
		return guess
	}
	return result
}