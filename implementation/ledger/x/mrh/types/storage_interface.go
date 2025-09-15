package types

import (
	"context"
	"encoding/json"
)

// MRHStorage defines the interface for MRH graph storage
// This allows us to swap between local and IPFS implementations
type MRHStorage interface {
	// Store saves an MRH graph and returns its content hash
	Store(ctx context.Context, graph *MRHGraph) (string, error)
	
	// Retrieve gets an MRH graph by its content hash
	Retrieve(ctx context.Context, hash string) (*MRHGraph, error)
	
	// Delete removes an MRH graph (optional for IPFS)
	Delete(ctx context.Context, hash string) error
	
	// Exists checks if a graph exists
	Exists(ctx context.Context, hash string) (bool, error)
}

// MRHGraph represents a Markov Relevancy Horizon graph
type MRHGraph struct {
	// RDF-like triple structure
	Triples []Triple `json:"triples"`
	
	// Metadata
	LctID       string            `json:"lct_id"`
	Version     int64             `json:"version"`
	CreatedAt   int64             `json:"created_at"`
	ContextDepth uint32           `json:"context_depth"`
	Metadata    map[string]string `json:"metadata"`
}

// Triple represents an RDF triple (subject-predicate-object)
type Triple struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	
	// Additional Web4 context
	Weight    float64 `json:"weight"`    // Relevancy weight
	Distance  uint32  `json:"distance"`  // Fractal distance
	Timestamp int64   `json:"timestamp"` // When this relationship was established
}

// MRHContext represents the context boundary for an entity
type MRHContext struct {
	CenterLCT    string   `json:"center_lct"`
	Radius       uint32   `json:"radius"`       // Fractal depth
	TrustDecay   float64  `json:"trust_decay"`  // How trust degrades with distance
	IncludedLCTs []string `json:"included_lcts"` // LCTs within this context
}

// WitnessRelationship represents a witness connection in the MRH
type WitnessRelationship struct {
	WitnessLCT string  `json:"witness_lct"`
	SubjectLCT string  `json:"subject_lct"`
	EventType  string  `json:"event_type"`
	Timestamp  int64   `json:"timestamp"`
	Signature  []byte  `json:"signature"`
	TrustBoost float64 `json:"trust_boost"` // How much this witnessing increases trust
}

// MRHTraversal provides methods for navigating the MRH graph
type MRHTraversal interface {
	// FindPath finds the shortest trust path between two LCTs
	FindPath(from, to string, maxDepth uint32) ([]string, error)
	
	// CalculateTrust calculates trust between two LCTs based on MRH distance
	CalculateTrust(from, to string) (float64, error)
	
	// GetContext returns all LCTs within a context boundary
	GetContext(center string, radius uint32) (*MRHContext, error)
	
	// GetWitnesses returns all witnesses for an LCT
	GetWitnesses(lctID string) ([]*WitnessRelationship, error)
}

// StorageConfig contains configuration for storage backends
type StorageConfig struct {
	Type string `json:"type"` // "local" or "ipfs"
	
	// Local storage config
	LocalPath string `json:"local_path,omitempty"`
	
	// IPFS config (for future use)
	IPFSHost string `json:"ipfs_host,omitempty"`
	IPFSPort int    `json:"ipfs_port,omitempty"`
}

// NewMRHStorage creates a storage instance based on config
func NewMRHStorage(config StorageConfig) (MRHStorage, error) {
	switch config.Type {
	case "local":
		return NewLocalMRHStorage(config.LocalPath)
	case "ipfs":
		// Future: return NewIPFSMRHStorage(config.IPFSHost, config.IPFSPort)
		return nil, ErrIPFSNotImplemented
	default:
		return nil, ErrInvalidStorageType
	}
}

// Helper functions for graph operations

// AddTriple adds a new triple to the graph
func (g *MRHGraph) AddTriple(subject, predicate, object string, weight float64) {
	triple := Triple{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Weight:    weight,
		Timestamp: getCurrentTimestamp(),
	}
	g.Triples = append(g.Triples, triple)
}

// GetSubjects returns all unique subjects in the graph
func (g *MRHGraph) GetSubjects() []string {
	seen := make(map[string]bool)
	var subjects []string
	for _, t := range g.Triples {
		if !seen[t.Subject] {
			seen[t.Subject] = true
			subjects = append(subjects, t.Subject)
		}
	}
	return subjects
}

// GetRelationships returns all triples where the given LCT is subject or object
func (g *MRHGraph) GetRelationships(lctID string) []Triple {
	var relationships []Triple
	for _, t := range g.Triples {
		if t.Subject == lctID || t.Object == lctID {
			relationships = append(relationships, t)
		}
	}
	return relationships
}

// ToJSON serializes the graph to JSON
func (g *MRHGraph) ToJSON() ([]byte, error) {
	return json.Marshal(g)
}

// FromJSON deserializes a graph from JSON
func (g *MRHGraph) FromJSON(data []byte) error {
	return json.Unmarshal(data, g)
}

// Helper function to get current timestamp (implement in keeper)
func getCurrentTimestamp() int64 {
	// This will be implemented in the keeper package
	return 0
}