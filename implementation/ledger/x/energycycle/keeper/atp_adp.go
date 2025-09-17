package keeper

import (
  "fmt"
  
  "cosmossdk.io/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "racecar-web/x/energycycle/types"
)

// MintATPADP creates new ATP/ADP token pairs in society pool
func (k Keeper) MintATPADP(ctx sdk.Context, societyPool string, amount uint64) error {
  // Get or create pool
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("pool/"))
  
  var pool types.EnergyPool
  bz := store.Get([]byte(societyPool))
  if bz == nil {
    pool = types.EnergyPool{
      Id: societyPool,
      AtpBalance: 0,
      AdpBalance: 0,
      VelocityRequirement: 0.1,
      DemurrageRate: 0.001,
    }
  } else {
    k.cdc.MustUnmarshal(bz, &pool)
  }
  
  // Mint ATP tokens (start in charged state)
  pool.AtpBalance += amount
  
  // Store updated pool
  bz = k.cdc.MustMarshal(&pool)
  store.Set([]byte(societyPool), bz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "atp_minted",
      sdk.NewAttribute("pool", societyPool),
      sdk.NewAttribute("amount", fmt.Sprintf("%d", amount)),
    ),
  )
  
  return nil
}

// DischargeATP converts ATP to ADP through R6 action
func (k Keeper) DischargeATP(ctx sdk.Context, fromLCT string, amount uint64, r6Action types.R6Action) (*types.ADPToken, error) {
  // Validate R6 action
  if err := k.validateR6Action(ctx, r6Action); err != nil {
    return nil, fmt.Errorf("invalid R6 action: %w", err)
  }
  
  // Get society pool (simplified - using default)
  poolStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("pool/"))
  var pool types.EnergyPool
  bz := poolStore.Get([]byte("default"))
  if bz == nil {
    return nil, fmt.Errorf("no energy pool found")
  }
  k.cdc.MustUnmarshal(bz, &pool)
  
  // Check ATP balance
  if pool.AtpBalance < amount {
    return nil, fmt.Errorf("insufficient ATP: have %d, need %d", pool.AtpBalance, amount)
  }
  
  // Discharge ATP to ADP
  pool.AtpBalance -= amount
  pool.AdpBalance += amount
  
  // Create ADP token
  adpToken := &types.ADPToken{
    Id: fmt.Sprintf("adp:%s:%d", fromLCT, ctx.BlockTime().Unix()),
    Amount: amount,
    DischargedBy: fromLCT,
    DischargeTime: ctx.BlockTime().Unix(),
    R6Action: &r6Action,
  }
  
  // Store updated pool
  bz = k.cdc.MustMarshal(&pool)
  poolStore.Set([]byte("default"), bz)
  
  // Store ADP token
  adpStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("adp/"))
  adpBz := k.cdc.MustMarshal(adpToken)
  adpStore.Set([]byte(adpToken.Id), adpBz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "atp_discharged",
      sdk.NewAttribute("from_lct", fromLCT),
      sdk.NewAttribute("amount", fmt.Sprintf("%d", amount)),
      sdk.NewAttribute("adp_id", adpToken.Id),
    ),
  )
  
  return adpToken, nil
}

// RechargeADP converts ADP back to ATP through productive work
func (k Keeper) RechargeADP(ctx sdk.Context, toLCT string, adpTokenID string, workProof []byte) (*types.ATPToken, error) {
  // Get ADP token
  adpStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("adp/"))
  adpBz := adpStore.Get([]byte(adpTokenID))
  if adpBz == nil {
    return nil, fmt.Errorf("ADP token not found: %s", adpTokenID)
  }
  
  var adpToken types.ADPToken
  k.cdc.MustUnmarshal(adpBz, &adpToken)
  
  // Validate work proof (simplified)
  if len(workProof) < 32 {
    return nil, fmt.Errorf("invalid work proof")
  }
  
  // Get pool
  poolStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("pool/"))
  var pool types.EnergyPool
  bz := poolStore.Get([]byte("default"))
  k.cdc.MustUnmarshal(bz, &pool)
  
  // Recharge ADP to ATP
  pool.AdpBalance -= adpToken.Amount
  pool.AtpBalance += adpToken.Amount
  
  // Create ATP token
  atpToken := &types.ATPToken{
    Id: fmt.Sprintf("atp:%s:%d", toLCT, ctx.BlockTime().Unix()),
    Amount: adpToken.Amount,
    RechargedBy: toLCT,
    RechargeTime: ctx.BlockTime().Unix(),
    WorkProof: workProof,
  }
  
  // Store updated pool
  bz = k.cdc.MustMarshal(&pool)
  poolStore.Set([]byte("default"), bz)
  
  // Delete ADP token (consumed)
  adpStore.Delete([]byte(adpTokenID))
  
  // Store ATP token record
  atpStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("atp/"))
  atpBz := k.cdc.MustMarshal(atpToken)
  atpStore.Set([]byte(atpToken.Id), atpBz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "adp_recharged",
      sdk.NewAttribute("to_lct", toLCT),
      sdk.NewAttribute("amount", fmt.Sprintf("%d", adpToken.Amount)),
      sdk.NewAttribute("atp_id", atpToken.Id),
    ),
  )
  
  return atpToken, nil
}

// validateR6Action validates an R6 framework action
func (k Keeper) validateR6Action(ctx sdk.Context, action types.R6Action) error {
  // Check all R6 fields are present
  if action.Rules == "" || action.Roles == "" || action.Request == "" ||
     action.Reference == "" || action.Resource == "" || action.Result == "" {
    return fmt.Errorf("incomplete R6 action")
  }
  return nil
}