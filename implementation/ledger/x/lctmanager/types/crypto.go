package types

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	
	"golang.org/x/crypto/curve25519"
)

// CryptoManager handles cryptographic operations for LCTs
type CryptoManager struct {
	// We don't store private keys - they're managed externally
}

// NewCryptoManager creates a new crypto manager
func NewCryptoManager() *CryptoManager {
	return &CryptoManager{}
}

// GenerateKeyPair generates a new Ed25519 key pair for an LCT
// NOTE: Private key should be stored securely off-chain
func (cm *CryptoManager) GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}
	return pub, priv, nil
}

// DeriveX25519Keys derives X25519 keys from Ed25519 for Diffie-Hellman
func (cm *CryptoManager) DeriveX25519Keys(edPriv ed25519.PrivateKey) ([]byte, []byte, error) {
	if len(edPriv) != ed25519.PrivateKeySize {
		return nil, nil, fmt.Errorf("invalid Ed25519 private key size")
	}
	
	// Extract the 32-byte seed from the Ed25519 private key
	seed := edPriv[:32]
	
	// Hash the seed to get X25519 private key
	hash := sha256.Sum256(seed)
	x25519Priv := hash[:]
	
	// Clamp the private key as per X25519 spec
	x25519Priv[0] &= 248
	x25519Priv[31] &= 127
	x25519Priv[31] |= 64
	
	// Derive X25519 public key
	x25519Pub, err := curve25519.X25519(x25519Priv, curve25519.Basepoint)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive X25519 public key: %w", err)
	}
	
	return x25519Pub, x25519Priv, nil
}

// PerformDiffieHellman performs Diffie-Hellman key exchange
func (cm *CryptoManager) PerformDiffieHellman(privateKey, peerPublicKey []byte) ([]byte, error) {
	if len(privateKey) != 32 || len(peerPublicKey) != 32 {
		return nil, fmt.Errorf("invalid key sizes for Diffie-Hellman")
	}
	
	sharedSecret, err := curve25519.X25519(privateKey, peerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Diffie-Hellman failed: %w", err)
	}
	
	return sharedSecret, nil
}

// SignMessage signs a message with an Ed25519 private key
func (cm *CryptoManager) SignMessage(privateKey ed25519.PrivateKey, message []byte) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid Ed25519 private key size")
	}
	
	signature := ed25519.Sign(privateKey, message)
	return signature, nil
}

// VerifySignature verifies an Ed25519 signature
func (cm *CryptoManager) VerifySignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
	if len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	
	return ed25519.Verify(publicKey, message, signature)
}

// GenerateLCTID generates a unique LCT identifier
func (cm *CryptoManager) GenerateLCTID(entityType string, publicKey ed25519.PublicKey) string {
	// Create a deterministic ID from the public key
	hash := sha256.Sum256(publicKey)
	shortHash := hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter ID
	
	// Format: lct:web4:act:<entity_type>:<hash>
	return fmt.Sprintf("lct:web4:act:%s:%s", entityType, shortHash)
}

// GenerateProofOfAgency creates a cryptographic proof of delegated agency
func (cm *CryptoManager) GenerateProofOfAgency(
	delegatorPrivKey ed25519.PrivateKey,
	agentPubKey ed25519.PublicKey,
	permissions []string,
	constraints map[string]interface{},
) ([]byte, error) {
	// Create a structured message for the proof
	message := ProofOfAgencyMessage{
		AgentPublicKey: hex.EncodeToString(agentPubKey),
		Permissions:    permissions,
		Constraints:    constraints,
		Timestamp:      getCurrentTimestamp(),
		Nonce:          generateNonce(),
	}
	
	// Serialize the message
	messageBytes, err := serializeProofMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proof message: %w", err)
	}
	
	// Sign with delegator's private key
	signature := ed25519.Sign(delegatorPrivKey, messageBytes)
	
	// Combine message and signature
	proof := ProofOfAgency{
		Message:   messageBytes,
		Signature: signature,
	}
	
	return serializeProof(proof)
}

// VerifyProofOfAgency verifies a proof of delegated agency
func (cm *CryptoManager) VerifyProofOfAgency(
	delegatorPubKey ed25519.PublicKey,
	agentPubKey ed25519.PublicKey,
	proofBytes []byte,
) (bool, error) {
	// Deserialize the proof
	proof, err := deserializeProof(proofBytes)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize proof: %w", err)
	}
	
	// Verify the signature
	if !ed25519.Verify(delegatorPubKey, proof.Message, proof.Signature) {
		return false, nil
	}
	
	// Deserialize and verify the message content
	message, err := deserializeProofMessage(proof.Message)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize proof message: %w", err)
	}
	
	// Verify the agent public key matches
	expectedKey := hex.EncodeToString(agentPubKey)
	if message.AgentPublicKey != expectedKey {
		return false, fmt.Errorf("agent public key mismatch")
	}
	
	return true, nil
}

// GenerateWitnessSignature creates a witness signature for an event
func (cm *CryptoManager) GenerateWitnessSignature(
	witnessPrivKey ed25519.PrivateKey,
	subjectLCT string,
	eventType string,
	eventData []byte,
) ([]byte, error) {
	// Create witness message
	message := WitnessMessage{
		SubjectLCT: subjectLCT,
		EventType:  eventType,
		EventData:  hex.EncodeToString(eventData),
		Timestamp:  getCurrentTimestamp(),
	}
	
	// Serialize and sign
	messageBytes, err := serializeWitnessMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize witness message: %w", err)
	}
	
	signature := ed25519.Sign(witnessPrivKey, messageBytes)
	return signature, nil
}

// Helper structures

type ProofOfAgencyMessage struct {
	AgentPublicKey string                 `json:"agent_public_key"`
	Permissions    []string               `json:"permissions"`
	Constraints    map[string]interface{} `json:"constraints"`
	Timestamp      int64                  `json:"timestamp"`
	Nonce          string                 `json:"nonce"`
}

type ProofOfAgency struct {
	Message   []byte `json:"message"`
	Signature []byte `json:"signature"`
}

type WitnessMessage struct {
	SubjectLCT string `json:"subject_lct"`
	EventType  string `json:"event_type"`
	EventData  string `json:"event_data"`
	Timestamp  int64  `json:"timestamp"`
}

// Serialization helpers (implement with proper encoding)

func serializeProofMessage(msg ProofOfAgencyMessage) ([]byte, error) {
	// In production, use proper CBOR or protobuf encoding
	// For now, using simple concatenation
	data := fmt.Sprintf("%s|%v|%v|%d|%s", 
		msg.AgentPublicKey, 
		msg.Permissions, 
		msg.Constraints,
		msg.Timestamp,
		msg.Nonce,
	)
	return []byte(data), nil
}

func deserializeProofMessage(data []byte) (*ProofOfAgencyMessage, error) {
	// Placeholder - implement proper deserialization
	return &ProofOfAgencyMessage{}, nil
}

func serializeProof(proof ProofOfAgency) ([]byte, error) {
	// Placeholder - implement proper serialization
	return append(proof.Message, proof.Signature...), nil
}

func deserializeProof(data []byte) (*ProofOfAgency, error) {
	// Placeholder - implement proper deserialization
	if len(data) < ed25519.SignatureSize {
		return nil, fmt.Errorf("invalid proof data")
	}
	
	sigStart := len(data) - ed25519.SignatureSize
	return &ProofOfAgency{
		Message:   data[:sigStart],
		Signature: data[sigStart:],
	}, nil
}

func serializeWitnessMessage(msg WitnessMessage) ([]byte, error) {
	// Placeholder - implement proper serialization
	data := fmt.Sprintf("%s|%s|%s|%d",
		msg.SubjectLCT,
		msg.EventType,
		msg.EventData,
		msg.Timestamp,
	)
	return []byte(data), nil
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func getCurrentTimestamp() int64 {
	// This will be properly implemented in the keeper
	return 0
}