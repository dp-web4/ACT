package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
)

// Society represents a society with token pool
type Society struct {
	LCT  string
	Name string
}

// GetLastRechargeDay retrieves the last daily recharge timestamp
func (k Keeper) GetLastRechargeDay(ctx context.Context) (time.Time, error) {
	lastRecharge, err := k.LastRechargeDay.Get(ctx)
	if err != nil {
		if collections.IsNotFoundError(err) {
			// Return zero time if not found (first run)
			return time.Time{}, fmt.Errorf("last recharge day not initialized")
		}
		return time.Time{}, err
	}
	return lastRecharge, nil
}

// SetLastRechargeDay sets the last daily recharge timestamp
func (k Keeper) SetLastRechargeDay(ctx context.Context, day time.Time) error {
	return k.LastRechargeDay.Set(ctx, day)
}

// GetAllSocieties retrieves all societies with token pools
func (k Keeper) GetAllSocieties(ctx context.Context) ([]Society, error) {
	var societies []Society

	// Iterate through all society pools
	iter, err := k.SocietyTokenPools.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			continue
		}

		pool, err := iter.Value()
		if err != nil {
			continue
		}

		societies = append(societies, Society{
			LCT:  pool.SocietyID,
			Name: key, // Society LCT ID is the key
		})
	}

	return societies, nil
}
