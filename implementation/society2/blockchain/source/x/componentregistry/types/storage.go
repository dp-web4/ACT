package types

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ComponentVerificationBackend defines the interface for off-chain component verification
type ComponentVerificationBackend interface {
	// VerifyComponentPairing checks if two components can be paired
	VerifyComponentPairing(ctx context.Context, componentA, componentB string) (bool, string, error)

	// GetComponentMetadata retrieves metadata for a component
	GetComponentMetadata(ctx context.Context, componentID string) (map[string]interface{}, error)

	// Privacy-focused methods for anonymous component operations
	GenerateComponentHash(ctx context.Context, realComponentID string) (string, error)
	ResolveComponentHash(ctx context.Context, componentHash string) (map[string]interface{}, error)
	VerifyComponentPairingWithHashes(ctx context.Context, componentHashA, componentHashB string) (bool, string, error)
	GetAnonymousComponentMetadata(ctx context.Context, componentHash string) (map[string]interface{}, error)
}

// MockMySQLBackend implements a simple mock MySQL backend for the race car demo
type MockMySQLBackend struct {
	// Simple pairing rules: componentA -> allowed components
	pairingRules map[string][]string

	// Negative indicators that will deny pairing
	negativeIndicators []string

	// Component metadata cache
	componentMetadata map[string]map[string]interface{}

	// Privacy-focused: Hash to real component ID mapping
	hashToComponent map[string]string
	componentToHash map[string]string
}

// NewMockMySQLBackend creates a new mock MySQL backend with race car demo rules
func NewMockMySQLBackend() *MockMySQLBackend {
	backend := &MockMySQLBackend{
		pairingRules: map[string][]string{
			// Battery modules can pair with battery packs
			"MODBATT-MOD-RC001-001": {"MODBATT-PACK-RC001-A", "MODBATT-PACK-RC001-B"},
			"MODBATT-MOD-RC001-002": {"MODBATT-PACK-RC001-A", "MODBATT-PACK-RC001-B"},

			// Battery packs can pair with host ECUs
			"MODBATT-PACK-RC001-A": {"MODBATT-HOST-RC001", "MODBATT-HOST-RC002"},
			"MODBATT-PACK-RC001-B": {"MODBATT-HOST-RC001", "MODBATT-HOST-RC002"},

			// Host ECUs can pair with sensors
			"MODBATT-HOST-RC001": {"TEMP-SENSOR-RC001", "VOLTAGE-SENSOR-RC001"},
			"MODBATT-HOST-RC002": {"TEMP-SENSOR-RC001", "VOLTAGE-SENSOR-RC001"},
		},
		negativeIndicators: []string{
			"FAULTY", "RECALLED", "EXPIRED", "INCOMPATIBLE", "BLOCKED",
		},
		componentMetadata: map[string]map[string]interface{}{
			"MODBATT-MOD-RC001-001": {
				"type":          "battery_module",
				"capacity":      "25.6kWh",
				"voltage":       "400V",
				"manufacturer":  "RaceCarBatteryCo",
				"compatibility": "RC001_series",
			},
			"MODBATT-PACK-RC001-A": {
				"type":          "battery_pack",
				"capacity":      "51.2kWh",
				"voltage":       "800V",
				"manufacturer":  "RaceCarBatteryCo",
				"compatibility": "RC001_series",
			},
			"MODBATT-HOST-RC001": {
				"type":             "host_ecu",
				"firmware_version": "v2.1.0",
				"manufacturer":     "RaceCarElectronics",
				"compatibility":    "RC001_series",
			},
		},
		hashToComponent: make(map[string]string),
		componentToHash: make(map[string]string),
	}

	// Initialize hash mappings for existing components
	for componentID := range backend.componentMetadata {
		hash := backend.generateHash(componentID)
		backend.hashToComponent[hash] = componentID
		backend.componentToHash[componentID] = hash
	}

	return backend
}

// generateHash creates a SHA-256 hash of the component ID
func (m *MockMySQLBackend) generateHash(componentID string) string {
	hash := sha256.Sum256([]byte(componentID))
	return hex.EncodeToString(hash[:])
}

// GenerateComponentHash creates an anonymous hash for a real component ID
func (m *MockMySQLBackend) GenerateComponentHash(ctx context.Context, realComponentID string) (string, error) {
	if hash, exists := m.componentToHash[realComponentID]; exists {
		return hash, nil
	}

	// Generate new hash
	hash := m.generateHash(realComponentID)
	m.hashToComponent[hash] = realComponentID
	m.componentToHash[realComponentID] = hash

	return hash, nil
}

// ResolveComponentHash resolves a component hash to real component data
func (m *MockMySQLBackend) ResolveComponentHash(ctx context.Context, componentHash string) (map[string]interface{}, error) {
	realComponentID, exists := m.hashToComponent[componentHash]
	if !exists {
		return nil, fmt.Errorf("component hash not found: %s", componentHash)
	}

	// Return real component data (this would be restricted in production)
	return map[string]interface{}{
		"real_component_id": realComponentID,
		"component_hash":    componentHash,
		"metadata":          m.componentMetadata[realComponentID],
	}, nil
}

// VerifyComponentPairingWithHashes verifies pairing using component hashes
func (m *MockMySQLBackend) VerifyComponentPairingWithHashes(ctx context.Context, componentHashA, componentHashB string) (bool, string, error) {
	// Resolve hashes to real component IDs
	realComponentA, existsA := m.hashToComponent[componentHashA]
	if !existsA {
		return false, "component hash A not found", nil
	}

	realComponentB, existsB := m.hashToComponent[componentHashB]
	if !existsB {
		return false, "component hash B not found", nil
	}

	// Use existing verification logic with real component IDs
	return m.VerifyComponentPairing(ctx, realComponentA, realComponentB)
}

// GetAnonymousComponentMetadata returns only non-sensitive metadata for a component hash
func (m *MockMySQLBackend) GetAnonymousComponentMetadata(ctx context.Context, componentHash string) (map[string]interface{}, error) {
	realComponentID, exists := m.hashToComponent[componentHash]
	if !exists {
		return map[string]interface{}{
			"component_hash": componentHash,
			"status":         "not_found",
			"message":        "component hash not found in database",
		}, nil
	}

	// Get full metadata
	fullMetadata := m.componentMetadata[realComponentID]
	if fullMetadata == nil {
		fullMetadata = make(map[string]interface{})
	}

	// Return only non-sensitive metadata (privacy-focused)
	return map[string]interface{}{
		"component_hash": componentHash,
		"type":           fullMetadata["type"],         // Generic type only
		"status":         "active",                     // Generic status
		"trust_anchor":   "cryptographic_trust_anchor", // No commercial data
		// NO manufacturer, part numbers, specifications, or commercial details
	}, nil
}

// VerifyComponentPairing implements the pairing verification logic
func (m *MockMySQLBackend) VerifyComponentPairing(ctx context.Context, componentA, componentB string) (bool, string, error) {
	// Check for negative indicators in component IDs
	for _, indicator := range m.negativeIndicators {
		if strings.Contains(strings.ToUpper(componentA), indicator) ||
			strings.Contains(strings.ToUpper(componentB), indicator) {
			return false, fmt.Sprintf("pairing denied: component contains negative indicator '%s'", indicator), nil
		}
	}

	// Check if componentA is allowed to pair with componentB
	if allowedComponents, exists := m.pairingRules[componentA]; exists {
		for _, allowed := range allowedComponents {
			if allowed == componentB {
				return true, "pairing allowed: components are compatible", nil
			}
		}
	}

	// Check reverse direction (componentB -> componentA)
	if allowedComponents, exists := m.pairingRules[componentB]; exists {
		for _, allowed := range allowedComponents {
			if allowed == componentA {
				return true, "pairing allowed: components are compatible", nil
			}
		}
	}

	return false, "pairing denied: components are not compatible", nil
}

// GetComponentMetadata retrieves metadata for a component
func (m *MockMySQLBackend) GetComponentMetadata(ctx context.Context, componentID string) (map[string]interface{}, error) {
	if metadata, exists := m.componentMetadata[componentID]; exists {
		return metadata, nil
	}

	// Return default metadata for unknown components
	return map[string]interface{}{
		"type":    "unknown",
		"status":  "unverified",
		"message": "component not found in database",
	}, nil
}

// AddPairingRule allows dynamic addition of pairing rules (useful for testing)
func (m *MockMySQLBackend) AddPairingRule(componentA string, allowedComponents []string) {
	m.pairingRules[componentA] = allowedComponents
}

// AddNegativeIndicator allows dynamic addition of negative indicators
func (m *MockMySQLBackend) AddNegativeIndicator(indicator string) {
	m.negativeIndicators = append(m.negativeIndicators, strings.ToUpper(indicator))
}

// AddComponentMetadata allows dynamic addition of component metadata
func (m *MockMySQLBackend) AddComponentMetadata(componentID string, metadata map[string]interface{}) {
	m.componentMetadata[componentID] = metadata

	// Update hash mappings
	hash := m.generateHash(componentID)
	m.hashToComponent[hash] = componentID
	m.componentToHash[componentID] = hash
}
