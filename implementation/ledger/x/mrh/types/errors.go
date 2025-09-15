package types

import (
	"cosmossdk.io/errors"
)

// MRH module errors
var (
	ErrInvalidStorageType = errors.Register(ModuleName, 1100, "invalid storage type")
	ErrIPFSNotImplemented = errors.Register(ModuleName, 1101, "IPFS storage not yet implemented")
	ErrGraphNotFound      = errors.Register(ModuleName, 1102, "MRH graph not found")
	ErrInvalidGraph       = errors.Register(ModuleName, 1103, "invalid MRH graph")
	ErrInvalidTriple      = errors.Register(ModuleName, 1104, "invalid RDF triple")
	ErrPathNotFound       = errors.Register(ModuleName, 1105, "no path found in MRH")
	ErrContextTooLarge    = errors.Register(ModuleName, 1106, "context radius too large")
	ErrInvalidLCT         = errors.Register(ModuleName, 1107, "invalid LCT identifier")
)