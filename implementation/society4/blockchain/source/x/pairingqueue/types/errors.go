package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/pairingqueue module sentinel errors
var (
	ErrComponentNotFound     = errors.Register(ModuleName, 1201, "component not found")
	ErrRequestNotFound       = errors.Register(ModuleName, 1202, "request not found")
	ErrInvalidRequestStatus  = errors.Register(ModuleName, 1203, "invalid request status")
	ErrInvalidOperationType  = errors.Register(ModuleName, 1204, "invalid operation type")
	ErrInvalidProxyComponent = errors.Register(ModuleName, 1205, "invalid proxy component")
	ErrQueueFull             = errors.Register(ModuleName, 1206, "queue is full")
	ErrInvalidSigner         = errors.Register(ModuleName, 1207, "expected gov account as only signer for proposal message")
	ErrComponentNotVerified  = errors.Register(ModuleName, 1208, "component not verified")
	ErrPairingNotAuthorized  = errors.Register(ModuleName, 1209, "pairing not authorized")
	ErrInsufficientTrust     = errors.Register(ModuleName, 1210, "insufficient trust score")
)
