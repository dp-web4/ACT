package types

import (
    "crypto/ed25519"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"

    hardware "society4chain/x/hardware/types"
)

// SelfLCT is the root identity for Society 4, bound to hardware
type SelfLCT struct {
    ID              string                     `json:"id"`
    Name            string                     `json:"name"`
    Description     string                     `json:"description"`
    HardwareBinding *hardware.HardwareBinding `json:"hardware_binding"`
    PublicKey       ed25519.PublicKey         `json:"public_key"`
    PrivateKey      ed25519.PrivateKey        `json:"-"` // Never serialize
    GenesisHeight   uint64                    `json:"genesis_height"`
    CreatedAt       time.Time                 `json:"created_at"`
    RoleChildren    []string                  `json:"role_children"`
    Witnesses       []WitnessAttestation      `json:"witnesses"`
}

// WitnessAttestation represents external validation of self-LCT
type WitnessAttestation struct {
    WitnessSociety string    `json:"witness_society"`
    TrustLevel     uint64    `json:"trust_level"`
    Timestamp      time.Time `json:"timestamp"`
    Signature      []byte    `json:"signature"`
}

// CreateSelfLCT creates the genesis self-LCT bound to current hardware
func CreateSelfLCT() (*SelfLCT, error) {
    // Extract current hardware
    hwBinding, err := hardware.ExtractCurrentHardware()
    if err != nil {
        return nil, fmt.Errorf("failed to extract hardware: %w", err)
    }

    // Validate hardware binding
    if err := hwBinding.ValidateBinding(); err != nil {
        return nil, fmt.Errorf("invalid hardware binding: %w", err)
    }

    // Generate ED25519 key pair
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("failed to generate keys: %w", err)
    }

    // Create self-LCT
    selfLCT := &SelfLCT{
        ID:              generateLCTID("self"),
        Name:            "Society4-Self-LCT",
        Description:     "Root identity for Society 4, hardware-bound genesis token",
        HardwareBinding: hwBinding,
        PublicKey:       publicKey,
        PrivateKey:      privateKey,
        GenesisHeight:   0,
        CreatedAt:       time.Now(),
        RoleChildren:    []string{},
        Witnesses:       []WitnessAttestation{},
    }

    // Sign the hardware binding
    hwBinding.Signature = selfLCT.SignData(hwBinding.Hash())

    return selfLCT, nil
}

// VerifyHardware checks if self-LCT matches current hardware
func (s *SelfLCT) VerifyHardware() error {
    if s.HardwareBinding == nil {
        return fmt.Errorf("no hardware binding found")
    }

    return s.HardwareBinding.VerifyHardware()
}

// SignData signs arbitrary data with the self-LCT private key
func (s *SelfLCT) SignData(data []byte) []byte {
    if s.PrivateKey == nil {
        return nil
    }
    return ed25519.Sign(s.PrivateKey, data)
}

// VerifySignature verifies a signature against the self-LCT public key
func (s *SelfLCT) VerifySignature(data []byte, signature []byte) bool {
    return ed25519.Verify(s.PublicKey, data, signature)
}

// SignBlock signs a block hash with hardware verification
func (s *SelfLCT) SignBlock(blockHash []byte) ([]byte, error) {
    // First verify we're on the correct hardware
    if err := s.VerifyHardware(); err != nil {
        return nil, fmt.Errorf("hardware verification failed: %w", err)
    }

    // Create composite data including hardware hash
    composite := append(blockHash, s.HardwareBinding.Hash()...)

    // Sign the composite
    signature := s.SignData(composite)
    if signature == nil {
        return nil, fmt.Errorf("failed to sign block")
    }

    return signature, nil
}

// VerifyBlockSignature verifies a block signature with hardware check
func (s *SelfLCT) VerifyBlockSignature(blockHash []byte, signature []byte) error {
    // Verify hardware first
    if err := s.VerifyHardware(); err != nil {
        return fmt.Errorf("hardware verification failed: %w", err)
    }

    // Create composite data
    composite := append(blockHash, s.HardwareBinding.Hash()...)

    // Verify signature
    if !s.VerifySignature(composite, signature) {
        return fmt.Errorf("signature verification failed")
    }

    return nil
}

// Hash returns the canonical hash of the self-LCT
func (s *SelfLCT) Hash() []byte {
    h := sha256.New()
    h.Write([]byte(s.ID))
    h.Write([]byte(s.Name))
    h.Write(s.HardwareBinding.Hash())
    h.Write(s.PublicKey)
    h.Write([]byte(fmt.Sprintf("%d", s.GenesisHeight)))
    return h.Sum(nil)
}

// AddWitness adds an attestation from another society
func (s *SelfLCT) AddWitness(witness WitnessAttestation) error {
    // Verify witness signature (would need witness public key)
    // For now, just add it
    s.Witnesses = append(s.Witnesses, witness)
    return nil
}

// CreateRoleLCT creates a derived LCT for a queen or worker role
func (s *SelfLCT) CreateRoleLCT(roleName string, roleType string) (*RoleLCT, error) {
    // Verify hardware before creating child
    if err := s.VerifyHardware(); err != nil {
        return nil, fmt.Errorf("hardware verification failed: %w", err)
    }

    roleLCT := &RoleLCT{
        ID:         generateLCTID(roleName),
        ParentLCT:  s.ID,
        RoleName:   roleName,
        RoleType:   roleType,
        CreatedAt:  time.Now(),
        Active:     true,
    }

    // Sign the role LCT with self-LCT
    roleLCT.ParentSignature = s.SignData(roleLCT.Hash())

    // Add to children
    s.RoleChildren = append(s.RoleChildren, roleLCT.ID)

    return roleLCT, nil
}

// RoleLCT represents a derived LCT for a specific role
type RoleLCT struct {
    ID              string    `json:"id"`
    ParentLCT       string    `json:"parent_lct"`
    RoleName        string    `json:"role_name"`
    RoleType        string    `json:"role_type"` // "queen" or "worker"
    CreatedAt       time.Time `json:"created_at"`
    Active          bool      `json:"active"`
    ParentSignature []byte    `json:"parent_signature"`
}

// Hash returns the canonical hash of a role LCT
func (r *RoleLCT) Hash() []byte {
    h := sha256.New()
    h.Write([]byte(r.ID))
    h.Write([]byte(r.ParentLCT))
    h.Write([]byte(r.RoleName))
    h.Write([]byte(r.RoleType))
    return h.Sum(nil)
}

// generateLCTID creates a unique LCT identifier
func generateLCTID(prefix string) string {
    timestamp := time.Now().UnixNano()
    hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", prefix, timestamp)))
    return fmt.Sprintf("lct-%s-%s", prefix, hex.EncodeToString(hash[:8]))
}

// Export exports the public components of self-LCT for sharing
func (s *SelfLCT) Export() map[string]interface{} {
    return map[string]interface{}{
        "id":               s.ID,
        "name":             s.Name,
        "description":      s.Description,
        "hardware_hash":    s.HardwareBinding.HardwareHash,
        "public_key":       hex.EncodeToString(s.PublicKey),
        "genesis_height":   s.GenesisHeight,
        "created_at":       s.CreatedAt,
        "witness_count":    len(s.Witnesses),
        "role_children":    s.RoleChildren,
    }
}