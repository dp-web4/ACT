package keeper

import (
	"context"
	"racecar-web/x/trusttensor/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set params if provided, otherwise use defaults
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		// Use default params if genesis params are invalid
		defaultParams := types.DefaultParams()
		return k.Params.Set(ctx, defaultParams)
	}
	return nil
}

// ExportGenesis returns the module's exported genesis state.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		// Return default params if none exist
		params = types.DefaultParams()
	}

	return &types.GenesisState{
		Params: params,
	}, nil
}
