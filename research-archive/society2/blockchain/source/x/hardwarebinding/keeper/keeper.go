package keeper

import (
    "fmt"

    "cosmossdk.io/log"
    "cosmossdk.io/store/storetypes"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"

    "society4chain/x/hardwarebinding/types"
)

type Keeper struct {
    cdc      codec.BinaryCodec
    storeKey storetypes.StoreKey
    logger   log.Logger

    // Genesis hardware binding
    genesisBinding *types.HardwareBinding

    // Validation enabled flag
    validationEnabled bool
}

// NewKeeper creates a new hardware binding keeper
func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey storetypes.StoreKey,
    logger log.Logger,
) Keeper {
    return Keeper{
        cdc:               cdc,
        storeKey:          storeKey,
        logger:            logger.With(log.ModuleKey, "x/"+types.ModuleName),
        validationEnabled: true, // Enable by default
    }
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
    return k.logger.With("height", ctx.BlockHeight())
}

// SetGenesisBinding sets the genesis hardware binding
func (k *Keeper) SetGenesisBinding(ctx sdk.Context, binding *types.HardwareBinding) {
    k.genesisBinding = binding

    // Store in state
    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(binding)
    store.Set([]byte("genesis_hardware"), bz)

    k.Logger(ctx).Info("Genesis hardware binding set",
        "hardware_hash", binding.HardwareHash,
        "platform", binding.Platform)
}

// GetGenesisBinding returns the genesis hardware binding
func (k Keeper) GetGenesisBinding(ctx sdk.Context) (*types.HardwareBinding, error) {
    if k.genesisBinding != nil {
        return k.genesisBinding, nil
    }

    // Try to load from store
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte("genesis_hardware"))
    if bz == nil {
        return nil, fmt.Errorf("genesis hardware binding not found")
    }

    var binding types.HardwareBinding
    k.cdc.MustUnmarshal(bz, &binding)
    k.genesisBinding = &binding

    return &binding, nil
}

// ValidateHardware validates current hardware against genesis
func (k Keeper) ValidateHardware(ctx sdk.Context) error {
    if !k.validationEnabled {
        return nil // Skip validation if disabled
    }

    genesis, err := k.GetGenesisBinding(ctx)
    if err != nil {
        return fmt.Errorf("failed to get genesis binding: %w", err)
    }

    // Verify hardware matches
    if err := genesis.VerifyHardware(); err != nil {
        return fmt.Errorf("hardware validation failed: %w", err)
    }

    return nil
}

// BeginBlocker performs hardware validation at the start of each block
func (k Keeper) BeginBlocker(ctx sdk.Context) error {
    // Only validate every 100 blocks to reduce overhead
    if ctx.BlockHeight()%100 != 0 {
        return nil
    }

    if err := k.ValidateHardware(ctx); err != nil {
        k.Logger(ctx).Error("Hardware validation failed",
            "error", err,
            "height", ctx.BlockHeight())

        // Store validation failure
        k.RecordValidationFailure(ctx, err)

        // In production, this could panic to halt the chain
        // For now, just log the error
        // panic(fmt.Sprintf("CRITICAL: Hardware validation failed at height %d: %v",
        //     ctx.BlockHeight(), err))
    } else {
        k.Logger(ctx).Info("Hardware validation successful",
            "height", ctx.BlockHeight())
    }

    return nil
}

// EndBlocker performs end-of-block processing
func (k Keeper) EndBlocker(ctx sdk.Context) error {
    // No end block processing needed for now
    return nil
}

// RecordValidationFailure records a hardware validation failure
func (k Keeper) RecordValidationFailure(ctx sdk.Context, err error) {
    store := ctx.KVStore(k.storeKey)
    key := fmt.Sprintf("failure_%d", ctx.BlockHeight())
    store.Set([]byte(key), []byte(err.Error()))
}

// SetValidationEnabled enables or disables hardware validation
func (k *Keeper) SetValidationEnabled(enabled bool) {
    k.validationEnabled = enabled
}

// InitGenesis initializes the module's genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
    // Set genesis hardware binding if provided
    if genState.HardwareBinding != nil {
        k.SetGenesisBinding(ctx, genState.HardwareBinding)
    } else {
        // Extract current hardware as genesis
        binding, err := types.ExtractCurrentHardware()
        if err != nil {
            k.Logger(ctx).Error("Failed to extract hardware for genesis", "error", err)
        } else {
            binding.BlockHeight = ctx.BlockHeight()
            binding.Timestamp = ctx.BlockTime().Unix()
            k.SetGenesisBinding(ctx, binding)
        }
    }
}

// ExportGenesis exports the module's genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
    genesis, _ := k.GetGenesisBinding(ctx)
    return &types.GenesisState{
        HardwareBinding: genesis,
    }
}