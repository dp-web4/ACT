package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "lctmanager"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lctmanager"

	// GovModuleName duplicates the gov module's name
	GovModuleName = "gov"
)

// Collection key prefixes
var (
	ParamsKey                = collections.NewPrefix([]byte{0x00})
	LctKeyPrefix             = collections.NewPrefix([]byte{0x01})
	RelationshipKeyPrefix    = collections.NewPrefix([]byte{0x02})
	LCTMediatedPairingPrefix = collections.NewPrefix([]byte{0x03})
	SessionKeyExchangePrefix = collections.NewPrefix([]byte{0x04})
	PairingChallengePrefix   = collections.NewPrefix([]byte{0x05})
	SplitKeyPrefix           = collections.NewPrefix([]byte{0x06})
)

// KeyPrefix returns the key prefix for a specific LCT
func KeyPrefix(p string) []byte {
	return []byte(p)
}

// LctKey returns the key for storing an LCT
func LctKey(lctID string) []byte {
	return append(LctKeyPrefix.Bytes(), []byte(lctID)...)
}

// RelationshipKey returns the key for storing a component relationship
func RelationshipKey(componentID string) []byte {
	return append(RelationshipKeyPrefix.Bytes(), []byte(componentID)...)
}

// LCTMediatedPairingKey returns the key for storing an LCT-mediated pairing
func LCTMediatedPairingKey(pairingID string) []byte {
	return append(LCTMediatedPairingPrefix.Bytes(), []byte(pairingID)...)
}

// SessionKeyExchangeKey returns the key for storing a session key exchange
func SessionKeyExchangeKey(pairingID string) []byte {
	return append(SessionKeyExchangePrefix.Bytes(), []byte(pairingID)...)
}

// PairingChallengeKey returns the key for storing a pairing challenge
func PairingChallengeKey(challengeID string) []byte {
	return append(PairingChallengePrefix.Bytes(), []byte(challengeID)...)
}
