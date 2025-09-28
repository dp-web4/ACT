package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/componentregistry module sentinel errors
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidComponentID   = errors.Register(ModuleName, 1101, "invalid component ID")
	ErrComponentExists      = errors.Register(ModuleName, 1102, "component already exists")
	ErrComponentNotFound    = errors.Register(ModuleName, 1103, "component not found")
	ErrInvalidComponentType = errors.Register(ModuleName, 1104, "invalid component type")
	ErrInvalidAuthority     = errors.Register(ModuleName, 1105, "invalid authority")
)
