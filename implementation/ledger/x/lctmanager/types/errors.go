package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/lctmanager module sentinel errors
var (
	ErrComponentNotFound    = errors.Register(ModuleName, 1201, "component not found")
	ErrLctNotFound          = errors.Register(ModuleName, 1202, "LCT not found")
	ErrInvalidLctStatus     = errors.Register(ModuleName, 1203, "invalid LCT status")
	ErrInvalidComponentPair = errors.Register(ModuleName, 1204, "invalid component pair")
	ErrInvalidSigner        = errors.Register(ModuleName, 1205, "invalid signer")
	ErrInvalidAuthority     = errors.Register(ModuleName, 1206, "invalid authority")
	ErrInvalidContext       = errors.Register(ModuleName, 1207, "invalid context")
	ErrInvalidProxy         = errors.Register(ModuleName, 1208, "invalid proxy component")
	ErrInvalidRequest       = errors.Register(ModuleName, 1100, "invalid request")
	ErrLctExists            = errors.Register(ModuleName, 1101, "LCT already exists")
)
