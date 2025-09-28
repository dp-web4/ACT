package keeper

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "society4chain/x/hardwarebinding/types"
)

type msgServer struct {
    Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
    return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateHardware handles hardware validation toggle (admin only)
func (k msgServer) UpdateHardware(ctx sdk.Context, msg *types.MsgUpdateHardware) (*types.MsgUpdateHardwareResponse, error) {
    // In production, this should check for admin/governance permissions
    // For now, we'll allow it for testing

    k.SetValidationEnabled(msg.Enable)

    k.Logger(ctx).Info("Hardware validation updated",
        "enabled", msg.Enable,
        "updater", msg.Creator)

    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "hardware_validation_updated",
            sdk.NewAttribute("enabled", fmt.Sprintf("%v", msg.Enable)),
            sdk.NewAttribute("updater", msg.Creator),
        ),
    )

    return &types.MsgUpdateHardwareResponse{}, nil
}

// Hardware handles hardware query
func (k Keeper) Hardware(ctx sdk.Context, req *types.QueryHardwareRequest) (*types.QueryHardwareResponse, error) {
    genesis, err := k.GetGenesisBinding(ctx)
    if err != nil {
        return nil, err
    }

    // Check if current hardware matches
    currentMatch := genesis.VerifyHardware() == nil

    return &types.QueryHardwareResponse{
        HardwareBinding:   genesis,
        CurrentMatch:      currentMatch,
        ValidationEnabled: k.validationEnabled,
    }, nil
}