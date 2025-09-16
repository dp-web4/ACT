package keeper

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/dp-web4/act/x/energycycle/types"
)

// Keeper maintains state for energycycle
type Keeper struct {
    storeKey sdk.StoreKey
    cdc      codec.BinaryCodec
    
}

// NewKeeper creates a new keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec) *Keeper {
    return &Keeper{
        storeKey: storeKey,
        cdc:      cdc,
    }
}


// SetLCT stores an LCT in the store
func (k Keeper) SetLCT(ctx sdk.Context, lct types.LCT) {
    store := ctx.KVStore(k.storeKey)
    b := k.cdc.MustMarshal(&lct)
    store.Set(types.LCTKey(lct.Id), b)
}

// GetLCT retrieves an LCT from the store
func (k Keeper) GetLCT(ctx sdk.Context, id string) (val types.LCT, found bool) {
    store := ctx.KVStore(k.storeKey)
    b := store.Get(types.LCTKey(id))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}

// RemoveLCT removes an LCT from the store
func (k Keeper) RemoveLCT(ctx sdk.Context, id string) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(types.LCTKey(id))
}

// GetAllLCT returns all LCTs
func (k Keeper) GetAllLCT(ctx sdk.Context) (list []types.LCT) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, []byte{})
    defer iterator.Close()
    
    for ; iterator.Valid(); iterator.Next() {
        var val types.LCT
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }
    return
}


// Queries implementation would go here


// ValidateLCT validates an LCT
func (k Keeper) ValidateLCT(ctx sdk.Context, lct types.LCT) error {
    if lct.Id == "" {
        return fmt.Errorf("LCT ID cannot be empty")
    }
    if lct.EntityType == "" {
        return fmt.Errorf("entity type cannot be empty")
    }
    return nil
}

