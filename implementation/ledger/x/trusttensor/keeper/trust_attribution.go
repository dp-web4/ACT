package keeper

import (
  "fmt"
  "math"
  
  "cosmossdk.io/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "racecar-web/x/trusttensor/types"
  lctmanagertypes "racecar-web/x/lctmanager/types"
)

// UpdateT3 updates the talent/training/temperament tensor for an LCT
func (k Keeper) UpdateT3(ctx sdk.Context, lctID string, role string, dimension string, delta float64) error {
  // Get LCT from lctmanager
  lctStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("../lctmanager/lct/"))
  bz := lctStore.Get([]byte(lctID))
  if bz == nil {
    return fmt.Errorf("LCT not found: %s", lctID)
  }
  
  var lct lctmanagertypes.LCT
  k.cdc.MustUnmarshal(bz, &lct)
  
  // Update appropriate dimension
  switch dimension {
  case "talent":
    lct.T3Tensor.Talent = clamp(lct.T3Tensor.Talent + delta, 0, 1)
  case "training":
    lct.T3Tensor.Training = clamp(lct.T3Tensor.Training + delta, 0, 1)
  case "temperament":
    lct.T3Tensor.Temperament = clamp(lct.T3Tensor.Temperament + delta, 0, 1)
  default:
    return fmt.Errorf("invalid T3 dimension: %s", dimension)
  }
  
  // Store updated LCT
  bz = k.cdc.MustMarshal(&lct)
  lctStore.Set([]byte(lctID), bz)
  
  // Store role-specific trust
  trustStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("trust/"))
  trustKey := fmt.Sprintf("%s:%s", lctID, role)
  
  var trust types.TrustRecord
  trustBz := trustStore.Get([]byte(trustKey))
  if trustBz == nil {
    trust = types.TrustRecord{
      LctId: lctID,
      Role: role,
      T3Score: 0.5,
      V3Score: 0.5,
      LastUpdate: ctx.BlockTime().Unix(),
    }
  } else {
    k.cdc.MustUnmarshal(trustBz, &trust)
  }
  
  // Calculate new T3 score (average of dimensions)
  trust.T3Score = (lct.T3Tensor.Talent + lct.T3Tensor.Training + lct.T3Tensor.Temperament) / 3
  trust.LastUpdate = ctx.BlockTime().Unix()
  
  // Store trust record
  trustBz = k.cdc.MustMarshal(&trust)
  trustStore.Set([]byte(trustKey), trustBz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "t3_updated",
      sdk.NewAttribute("lct_id", lctID),
      sdk.NewAttribute("role", role),
      sdk.NewAttribute("dimension", dimension),
      sdk.NewAttribute("delta", fmt.Sprintf("%f", delta)),
    ),
  )
  
  return nil
}

// UpdateV3 updates the veracity/validity/value tensor from outcomes
func (k Keeper) UpdateV3(ctx sdk.Context, lctID string, context string, outcome types.Outcome, witnesses []string) error {
  // Get LCT
  lctStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("../lctmanager/lct/"))
  bz := lctStore.Get([]byte(lctID))
  if bz == nil {
    return fmt.Errorf("LCT not found: %s", lctID)
  }
  
  var lct lctmanagertypes.LCT
  k.cdc.MustUnmarshal(bz, &lct)
  
  // Calculate impact based on witnesses
  impact := 0.1 // base impact
  if len(witnesses) > 0 {
    impact = 0.1 * math.Log(float64(len(witnesses)+1))
  }
  
  // Update V3 based on outcome
  if outcome.Success {
    lct.V3Tensor.Veracity = clamp(lct.V3Tensor.Veracity + impact, 0, 1)
    lct.V3Tensor.Validity = clamp(lct.V3Tensor.Validity + impact, 0, 1)
    lct.V3Tensor.Value = clamp(lct.V3Tensor.Value + impact*outcome.ValueGenerated, 0, 1)
  } else {
    lct.V3Tensor.Veracity = clamp(lct.V3Tensor.Veracity - impact/2, 0, 1)
    lct.V3Tensor.Validity = clamp(lct.V3Tensor.Validity - impact/2, 0, 1)
  }
  
  // Store updated LCT
  bz = k.cdc.MustMarshal(&lct)
  lctStore.Set([]byte(lctID), bz)
  
  // Record outcome
  outcomeStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("outcome/"))
  outcomeKey := fmt.Sprintf("%s:%d", lctID, ctx.BlockTime().Unix())
  outcomeBz := k.cdc.MustMarshal(&outcome)
  outcomeStore.Set([]byte(outcomeKey), outcomeBz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "v3_updated",
      sdk.NewAttribute("lct_id", lctID),
      sdk.NewAttribute("context", context),
      sdk.NewAttribute("success", fmt.Sprintf("%t", outcome.Success)),
      sdk.NewAttribute("witnesses", fmt.Sprintf("%d", len(witnesses))),
    ),
  )
  
  return nil
}

// GetTrustDistance calculates trust distance between two LCTs
func (k Keeper) GetTrustDistance(ctx sdk.Context, fromLCT, toLCT, role string) (float64, error) {
  // Get trust records
  trustStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("trust/"))
  
  fromKey := fmt.Sprintf("%s:%s", fromLCT, role)
  fromBz := trustStore.Get([]byte(fromKey))
  if fromBz == nil {
    return 1.0, nil // Maximum distance if no trust record
  }
  
  toKey := fmt.Sprintf("%s:%s", toLCT, role)
  toBz := trustStore.Get([]byte(toKey))
  if toBz == nil {
    return 1.0, nil
  }
  
  var fromTrust, toTrust types.TrustRecord
  k.cdc.MustUnmarshal(fromBz, &fromTrust)
  k.cdc.MustUnmarshal(toBz, &toTrust)
  
  // Calculate Euclidean distance in trust space
  t3Diff := math.Abs(fromTrust.T3Score - toTrust.T3Score)
  v3Diff := math.Abs(fromTrust.V3Score - toTrust.V3Score)
  
  distance := math.Sqrt(t3Diff*t3Diff + v3Diff*v3Diff) / math.Sqrt(2)
  
  return distance, nil
}

// ApplyTrustGravity applies trust-based attraction/repulsion
func (k Keeper) ApplyTrustGravity(ctx sdk.Context, lctID string) error {
  // Get all trust records for this LCT
  trustStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("trust/"))
  iterator := sdk.KVStorePrefixIterator(trustStore, []byte(lctID))
  defer iterator.Close()
  
  totalGravity := 0.0
  count := 0
  
  for ; iterator.Valid(); iterator.Next() {
    var trust types.TrustRecord
    k.cdc.MustUnmarshal(iterator.Value(), &trust)
    
    // High trust creates attraction (positive gravity)
    // Low trust creates repulsion (negative gravity)
    gravity := (trust.T3Score + trust.V3Score) - 1.0 // Range: -1 to +1
    totalGravity += gravity
    count++
  }
  
  if count > 0 {
    avgGravity := totalGravity / float64(count)
    
    // Store gravity effect
    gravityStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("gravity/"))
    gravityRecord := types.GravityRecord{
      LctId: lctID,
      Gravity: avgGravity,
      Timestamp: ctx.BlockTime().Unix(),
    }
    
    bz := k.cdc.MustMarshal(&gravityRecord)
    gravityStore.Set([]byte(lctID), bz)
    
    // Emit event
    ctx.EventManager().EmitEvent(
      sdk.NewEvent(
        "trust_gravity_applied",
        sdk.NewAttribute("lct_id", lctID),
        sdk.NewAttribute("gravity", fmt.Sprintf("%f", avgGravity)),
      ),
    )
  }
  
  return nil
}

// Helper function to clamp values between min and max
func clamp(value, min, max float64) float64 {
  if value < min {
    return min
  }
  if value > max {
    return max
  }
  return value
}