package types

// GenesisState defines the hardwarebinding module's genesis state
type GenesisState struct {
    HardwareBinding *HardwareBinding `json:"hardware_binding,omitempty"`
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
    return &GenesisState{
        HardwareBinding: nil, // Will be set on chain init
    }
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
    if data.HardwareBinding != nil {
        if data.HardwareBinding.Platform != "wsl2" {
            return ErrInvalidPlatform
        }
        if data.HardwareBinding.HardwareHash == "" {
            return ErrMissingHardwareHash
        }
    }
    return nil
}