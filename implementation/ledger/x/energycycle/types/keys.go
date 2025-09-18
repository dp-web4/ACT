package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "energycycle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_energycycle"

	// GovModuleName duplicates the gov module's name
	GovModuleName = "gov"
)

// Collection key prefixes for Web4 energy cycle storage
var (
	ParamsKey                  = collections.NewPrefix(0)
	EnergyOperationKey         = collections.NewPrefix(1)
	RelationshipAtpTokenKey    = collections.NewPrefix(2)
	RelationshipAdpTokenKey    = collections.NewPrefix(3)
	SocietyPoolKey             = collections.NewPrefix(4)
)

// Energy operation types
const (
	OperationTypeDischarge = "discharge"
	OperationTypeCharge    = "charge"
	OperationTypeTransfer  = "transfer"
	OperationTypeBalance   = "balance"
)

// Energy operation statuses
const (
	StatusCreated   = "created"
	StatusActive    = "active"
	StatusExecuting = "executing"
	StatusCompleted = "completed"
	StatusValidated = "validated"
	StatusFailed    = "failed"
)

// ATP/ADP token statuses
const (
	AtpStatusActive     = "active"
	AtpStatusDischarged = "discharged"
	AtpStatusExpired    = "expired"
)
