package keeper

import (
    "fmt"

    abci "github.com/cometbft/cometbft/abci/types"
    sdk "github.com/cosmos/cosmos-sdk/types"

    lctTypes "society4chain/x/lctmanager/types"
)

// HardwareValidator enforces hardware binding for block creation
type HardwareValidator struct {
    selfLCT *lctTypes.SelfLCT
}

// NewHardwareValidator creates a new hardware validator with self-LCT
func NewHardwareValidator(selfLCT *lctTypes.SelfLCT) *HardwareValidator {
    return &HardwareValidator{
        selfLCT: selfLCT,
    }
}

// ValidateBlock ensures the block is signed by the hardware-bound self-LCT
func (hv *HardwareValidator) ValidateBlock(ctx sdk.Context, req abci.RequestBeginBlock) error {
    // Extract block hash
    blockHash := req.Hash

    // Get block signature from header
    // In CometBFT, this would come from the proposer's signature
    var blockSignature []byte
    if req.Header.ProposerAddress != nil {
        // Extract signature from the block header
        // This would need integration with CometBFT's signing mechanism
        blockSignature = extractBlockSignature(req)
    }

    if len(blockSignature) == 0 {
        return fmt.Errorf("block missing required hardware signature")
    }

    // Verify the block signature with hardware check
    if err := hv.selfLCT.VerifyBlockSignature(blockHash, blockSignature); err != nil {
        return fmt.Errorf("hardware-bound signature verification failed: %w", err)
    }

    return nil
}

// SignBlock signs a block with hardware verification
func (hv *HardwareValidator) SignBlock(blockHash []byte) ([]byte, error) {
    return hv.selfLCT.SignBlock(blockHash)
}

// CheckHardwareBinding verifies current hardware matches genesis binding
func (hv *HardwareValidator) CheckHardwareBinding() error {
    if hv.selfLCT == nil {
        return fmt.Errorf("self-LCT not initialized")
    }

    return hv.selfLCT.VerifyHardware()
}

// extractBlockSignature extracts signature from block header
// This is a placeholder - actual implementation depends on CometBFT integration
func extractBlockSignature(req abci.RequestBeginBlock) []byte {
    // In production, this would extract the actual signature
    // from the block's proposer signature field
    return nil
}

// BeginBlocker performs hardware validation at the start of each block
func (k Keeper) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
    // Perform hardware validation
    if err := k.hardwareValidator.ValidateBlock(ctx, req); err != nil {
        // Log the error and panic to halt the chain
        // This ensures the chain cannot continue on wrong hardware
        ctx.Logger().Error("Hardware validation failed", "error", err)
        panic(fmt.Sprintf("CRITICAL: Hardware validation failed: %v", err))
    }

    // Log successful validation
    ctx.Logger().Info("Hardware validation successful",
        "height", ctx.BlockHeight(),
        "hardware_id", k.hardwareValidator.selfLCT.HardwareBinding.HardwareHash[:8])

    return abci.ResponseBeginBlock{}
}

// EndBlocker can perform additional hardware checks
func (k Keeper) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    // Periodic deep hardware verification (every 100 blocks)
    if ctx.BlockHeight()%100 == 0 {
        if err := k.hardwareValidator.CheckHardwareBinding(); err != nil {
            ctx.Logger().Error("Periodic hardware check failed", "error", err)
            // Store failure in state for monitoring
            k.SetHardwareCheckFailure(ctx, ctx.BlockHeight(), err.Error())
        }
    }

    return abci.ResponseEndBlock{}
}