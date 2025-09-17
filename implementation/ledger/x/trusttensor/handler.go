package trusttensor

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "racecar-web/x/trusttensor/keeper"
    "racecar-web/x/trusttensor/types"
)

// NewHandler returns a handler for trusttensor messages
func NewHandler(k keeper.Keeper) sdk.Handler {
    msgServer := keeper.NewMsgServerImpl(k)
    
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        ctx = ctx.WithEventManager(sdk.NewEventManager())
        
        switch msg := msg.(type) {

        case *types.MsgCreateLCT:
            res, err := msgServer.CreateLCT(sdk.WrapSDKContext(ctx), msg)
            return sdk.WrapServiceResult(ctx, res, err)
        case *types.MsgUpdateMRH:
            res, err := msgServer.UpdateMRH(sdk.WrapSDKContext(ctx), msg)
            return sdk.WrapServiceResult(ctx, res, err)
        default:
            return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest,
                "unrecognized %s message type: %T", types.ModuleName, msg)
        }
    }
}

// Handler implementations
