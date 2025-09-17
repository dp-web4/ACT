package keeper

import (
  "fmt"
  "crypto/ed25519"
  "encoding/hex"
  
  "cosmossdk.io/store/prefix"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "racecar-web/x/lctmanager/types"
)

// MintLCT creates a new LCT with Ed25519 identity
func (k Keeper) MintLCT(ctx sdk.Context, entityType string, pubKey ed25519.PublicKey) (string, error) {
  // Generate LCT ID from public key
  lctID := "lct:" + hex.EncodeToString(pubKey[:8])
  
  // Create birth certificate
  birthCert := types.BirthCertificate{
    Timestamp: ctx.BlockTime().Unix(),
    EntityType: entityType,
    GenesisBlock: ctx.BlockHeight(),
  }
  
  // Initialize T3/V3 tensors
  t3Tensor := types.T3Tensor{
    Talent: 0.5,
    Training: 0.5,
    Temperament: 0.5,
  }
  
  v3Tensor := types.V3Tensor{
    Veracity: 0.5,
    Validity: 0.5,
    Value: 0.5,
  }
  
  // Create LCT
  lct := types.LCT{
    Id: lctID,
    EntityType: entityType,
    Identity: types.LCTIdentity{
      Ed25519PubKey: pubKey,
    },
    BirthCertificate: &birthCert,
    T3Tensor: &t3Tensor,
    V3Tensor: &v3Tensor,
  }
  
  // Store LCT
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("lct/"))
  bz := k.cdc.MustMarshal(&lct)
  store.Set([]byte(lctID), bz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "lct_minted",
      sdk.NewAttribute("lct_id", lctID),
      sdk.NewAttribute("entity_type", entityType),
    ),
  )
  
  return lctID, nil
}

// BindLCT permanently binds an LCT to an entity
func (k Keeper) BindLCT(ctx sdk.Context, lctID string, entityID string, bindingProof []byte) error {
  // Get LCT
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("lct/"))
  bz := store.Get([]byte(lctID))
  if bz == nil {
    return fmt.Errorf("LCT not found: %s", lctID)
  }
  
  var lct types.LCT
  k.cdc.MustUnmarshal(bz, &lct)
  
  // Check if already bound
  if lct.BoundEntity != "" {
    return fmt.Errorf("LCT already bound to: %s", lct.BoundEntity)
  }
  
  // Verify binding proof (simplified for now)
  if len(bindingProof) < 32 {
    return fmt.Errorf("invalid binding proof")
  }
  
  // Update LCT with binding
  lct.BoundEntity = entityID
  lct.BindingTimestamp = ctx.BlockTime().Unix()
  
  // Store updated LCT
  bz = k.cdc.MustMarshal(&lct)
  store.Set([]byte(lctID), bz)
  
  // Emit event
  ctx.EventManager().EmitEvent(
    sdk.NewEvent(
      "lct_bound",
      sdk.NewAttribute("lct_id", lctID),
      sdk.NewAttribute("entity_id", entityID),
    ),
  )
  
  return nil
}

// GetLCTWithMRH retrieves an LCT with its Markov Relevancy Horizon
func (k Keeper) GetLCTWithMRH(ctx sdk.Context, lctID string) (*types.LCT, *types.MRH, error) {
  // Get LCT
  store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("lct/"))
  bz := store.Get([]byte(lctID))
  if bz == nil {
    return nil, nil, fmt.Errorf("LCT not found: %s", lctID)
  }
  
  var lct types.LCT
  k.cdc.MustUnmarshal(bz, &lct)
  
  // Get MRH
  mrhStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("mrh/"))
  mrhBz := mrhStore.Get([]byte(lctID))
  
  var mrh types.MRH
  if mrhBz != nil {
    k.cdc.MustUnmarshal(mrhBz, &mrh)
  } else {
    // Initialize empty MRH
    mrh = types.MRH{
      LctId: lctID,
      Edges: []types.MRHEdge{},
    }
  }
  
  return &lct, &mrh, nil
}