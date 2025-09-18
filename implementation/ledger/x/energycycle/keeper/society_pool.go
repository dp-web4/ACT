package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	runtime "github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"racecar-web/x/energycycle/types"
)

// GetSocietyPool retrieves a society pool from the store
func (k Keeper) GetSocietyPool(ctx context.Context, societyLct string) (types.SocietyPool, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.SocietyPoolKey)

	key := []byte(societyLct)
	bz := store.Get(key)

	if bz == nil {
		// Return empty pool if not exists
		return types.SocietyPool{
			SocietyLct:  societyLct,
			AtpBalance:  sdk.NewCoin("atp", math.ZeroInt()),
			AdpBalance:  sdk.NewCoin("adp", math.ZeroInt()),
			LastUpdate:  0,
			TotalMinted: "0",
			TotalDischarged: "0",
			TotalRecharged: "0",
			MetabolicState: "active",
			Metadata: make(map[string]string),
		}, nil
	}

	var pool types.SocietyPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, nil
}

// SetSocietyPool stores a society pool
func (k Keeper) SetSocietyPool(ctx context.Context, pool types.SocietyPool) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.SocietyPoolKey)

	key := []byte(pool.SocietyLct)
	pool.LastUpdate = sdkCtx.BlockTime().Unix()

	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
	return nil
}

// UpdateSocietyBalance updates ATP/ADP balances for a society
func (k Keeper) UpdateSocietyBalance(ctx context.Context, societyLct string, atpDelta, adpDelta math.Int) error {
	pool, err := k.GetSocietyPool(ctx, societyLct)
	if err != nil {
		return fmt.Errorf("failed to get society pool: %w", err)
	}

	// Initialize if first time
	if pool.LastUpdate == 0 {
		pool.SocietyLct = societyLct
		pool.AtpBalance = sdk.NewCoin("atp", math.ZeroInt())
		pool.AdpBalance = sdk.NewCoin("adp", math.ZeroInt())
		pool.MetabolicState = "active"
		pool.Metadata = make(map[string]string)
	}

	// Update ATP balance
	if !atpDelta.IsZero() {
		newAtpAmount := pool.AtpBalance.Amount.Add(atpDelta)
		if newAtpAmount.IsNegative() {
			return fmt.Errorf("insufficient ATP balance: have %s, need %s", pool.AtpBalance.Amount.String(), atpDelta.Abs().String())
		}
		pool.AtpBalance = sdk.NewCoin("atp", newAtpAmount)
	}

	// Update ADP balance
	if !adpDelta.IsZero() {
		newAdpAmount := pool.AdpBalance.Amount.Add(adpDelta)
		if newAdpAmount.IsNegative() {
			return fmt.Errorf("insufficient ADP balance: have %s, need %s", pool.AdpBalance.Amount.String(), adpDelta.Abs().String())
		}
		pool.AdpBalance = sdk.NewCoin("adp", newAdpAmount)
	}

	return k.SetSocietyPool(ctx, pool)
}

// MintADPToPool mints new ADP tokens to a society pool
func (k Keeper) MintADPToPool(ctx context.Context, societyLct string, amount math.Int, treasuryRole string) error {
	pool, err := k.GetSocietyPool(ctx, societyLct)
	if err != nil {
		return fmt.Errorf("failed to get society pool: %w", err)
	}

	// Initialize if first time
	if pool.LastUpdate == 0 {
		pool.SocietyLct = societyLct
		pool.AtpBalance = sdk.NewCoin("atp", math.ZeroInt())
		pool.AdpBalance = sdk.NewCoin("adp", math.ZeroInt())
		pool.TotalMinted = "0"
		pool.MetabolicState = "active"
		pool.TreasuryRole = treasuryRole
		pool.Metadata = make(map[string]string)
	}

	// TODO: Verify treasuryRole has minting authority
	// For now, we trust the caller has validated this

	// Add to ADP balance
	pool.AdpBalance = sdk.NewCoin("adp", pool.AdpBalance.Amount.Add(amount))

	// Update total minted
	totalMinted := math.LegacyMustNewDecFromStr(pool.TotalMinted)
	newTotal := totalMinted.Add(math.LegacyNewDecFromInt(amount))
	pool.TotalMinted = newTotal.String()

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("society_pool_mint",
			sdk.NewAttribute("society_lct", societyLct),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("treasury_role", treasuryRole),
			sdk.NewAttribute("new_balance", pool.AdpBalance.String()),
		),
	)

	return k.SetSocietyPool(ctx, pool)
}

// DischargeATPFromPool converts ATP to ADP for work performed
func (k Keeper) DischargeATPFromPool(ctx context.Context, societyLct string, amount math.Int, workerLct string, workDescription string) error {
	// Convert ATP to ADP (energy state change)
	err := k.UpdateSocietyBalance(ctx, societyLct, amount.Neg(), amount)
	if err != nil {
		return fmt.Errorf("failed to discharge ATP: %w", err)
	}

	// Update total discharged
	pool, _ := k.GetSocietyPool(ctx, societyLct)
	totalDischarged := math.LegacyMustNewDecFromStr(pool.TotalDischarged)
	pool.TotalDischarged = totalDischarged.Add(math.LegacyNewDecFromInt(amount)).String()
	k.SetSocietyPool(ctx, pool)

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("society_pool_discharge",
			sdk.NewAttribute("society_lct", societyLct),
			sdk.NewAttribute("worker_lct", workerLct),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("work", workDescription),
		),
	)

	return nil
}

// RechargeADPToATP converts ADP back to ATP through energy input
func (k Keeper) RechargeADPToATP(ctx context.Context, societyLct string, amount math.Int, producerLct string, energySource string) error {
	// Convert ADP to ATP (energy recharge)
	err := k.UpdateSocietyBalance(ctx, societyLct, amount, amount.Neg())
	if err != nil {
		return fmt.Errorf("failed to recharge ADP: %w", err)
	}

	// Update total recharged
	pool, _ := k.GetSocietyPool(ctx, societyLct)
	totalRecharged := math.LegacyMustNewDecFromStr(pool.TotalRecharged)
	pool.TotalRecharged = totalRecharged.Add(math.LegacyNewDecFromInt(amount)).String()
	k.SetSocietyPool(ctx, pool)

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent("society_pool_recharge",
			sdk.NewAttribute("society_lct", societyLct),
			sdk.NewAttribute("producer_lct", producerLct),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("energy_source", energySource),
		),
	)

	return nil
}

// GetAllSocietyPools returns all society pools (for genesis export)
func (k Keeper) GetAllSocietyPools(ctx context.Context) ([]types.SocietyPool, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.SocietyPoolKey)

	var pools []types.SocietyPool

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.SocietyPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools, nil
}