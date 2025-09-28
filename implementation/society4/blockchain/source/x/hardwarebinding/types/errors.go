package types

import (
    sdkerrors "cosmossdk.io/errors"
)

// x/hardwarebinding module sentinel errors
var (
    ErrInvalidPlatform      = sdkerrors.Register(ModuleName, 1100, "invalid platform")
    ErrMissingHardwareHash  = sdkerrors.Register(ModuleName, 1101, "missing hardware hash")
    ErrHardwareMismatch     = sdkerrors.Register(ModuleName, 1102, "hardware mismatch")
    ErrHardwareExtraction   = sdkerrors.Register(ModuleName, 1103, "failed to extract hardware")
    ErrGenesisNotFound      = sdkerrors.Register(ModuleName, 1104, "genesis hardware binding not found")
)