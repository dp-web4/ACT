package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "trusttensor"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trusttensor"

	// GovModuleName duplicates the gov module's name
	GovModuleName = "gov"
)

// Collection key prefixes for Web4 trust tensor storage
var (
	ParamsKey                  = collections.NewPrefix(0)
	RelationshipTrustTensorKey = collections.NewPrefix(1)
	ValueTensorKey             = collections.NewPrefix(2)
	TensorEntryKey             = collections.NewPrefix(3)
)
