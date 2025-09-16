package keeper

import (
    "encoding/json"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// MRH operations for Web4 compliance
func (k Keeper) SetMRH(ctx sdk.Context, lctID string, mrh MRH) error {
    store := ctx.KVStore(k.storeKey)
    key := []byte("mrh:" + lctID)
    value, err := json.Marshal(mrh)
    if err != nil {
        return err
    }
    store.Set(key, value)
    return nil
}

func (k Keeper) GetMRH(ctx sdk.Context, lctID string) (MRH, error) {
    store := ctx.KVStore(k.storeKey)
    key := []byte("mrh:" + lctID)
    value := store.Get(key)
    
    var mrh MRH
    if value == nil {
        return mrh, nil
    }
    
    err := json.Unmarshal(value, &mrh)
    return mrh, err
}

func (k Keeper) AddWitness(ctx sdk.Context, lctID string, witnessLCT string) error {
    mrh, err := k.GetMRH(ctx, lctID)
    if err != nil {
        return err
    }
    
    mrh.Witnessing = append(mrh.Witnessing, witnessLCT)
    return k.SetMRH(ctx, lctID, mrh)
}
