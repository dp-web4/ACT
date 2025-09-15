package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "pairingqueue"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_pairingqueue"

	// GovModuleName duplicates the gov module's name
	GovModuleName = "gov"
)

// Collection key prefixes
var (
	ParamsKey               = collections.NewPrefix(0)
	PairingRequestsPrefix   = collections.NewPrefix(1)
	OfflineOperationsPrefix = collections.NewPrefix(2)
)
