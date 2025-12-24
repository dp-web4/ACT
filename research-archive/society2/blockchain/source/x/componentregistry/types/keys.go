package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "componentregistry"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_componentregistry"

	// GovModuleName duplicates the gov module's name
	GovModuleName = "gov"
)

// Collection key prefixes for component registry storage
var (
	ParamsKey                = collections.NewPrefix(0)
	ComponentPrefix          = collections.NewPrefix(1)
	VerificationPrefix       = collections.NewPrefix(2)
	PairingRulesPrefix       = collections.NewPrefix(3)
	ManufacturerComponentKey = collections.NewPrefix(4)
	PairingAuthorizationKey  = collections.NewPrefix(5)
)

// Component status constants
const (
	StatusActive      = "active"
	StatusInactive    = "inactive"
	StatusMaintenance = "maintenance"
	StatusRetired     = "retired"
)

// Verification status constants
const (
	VerificationStatusPending  = "pending"
	VerificationStatusVerified = "verified"
	VerificationStatusRejected = "rejected"
	VerificationStatusExpired  = "expired"
)

// Component type constants
const (
	ComponentTypeModule  = "module"
	ComponentTypePack    = "pack"
	ComponentTypeHostECU = "host_ecu"
	ComponentTypeSensor  = "sensor"
)

// ComponentKey returns the key for storing a component
func ComponentKey(componentID string) []byte {
	return append(ComponentPrefix, []byte(componentID)...)
}

// VerificationKey returns the key for storing a component verification
func VerificationKey(componentID string) []byte {
	return append(VerificationPrefix, []byte(componentID)...)
}

// PairingRulesKey returns the key for storing pairing rules for a component type
func PairingRulesKey(componentType string) []byte {
	return append(PairingRulesPrefix, []byte(componentType)...)
}

// ParsePairingRulesKey extracts the component type from a pairing rules key
func ParsePairingRulesKey(key []byte) string {
	if len(key) <= len(PairingRulesPrefix) {
		return ""
	}
	return string(key[len(PairingRulesPrefix):])
}

// KeyPrefix returns the key prefix for a module
func KeyPrefix(p string) []byte {
	return []byte(p)
}
