package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"racecarweb/x/societytodo/types"
	energycyclekeeper "racecarweb/x/energycycle/keeper"
	lctmanagerkeeper "racecarweb/x/lctmanager/keeper"
	trustkeeper "racecarweb/x/trusttensor/keeper"
	mrhkeeper "racecarweb/x/mrh/keeper"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey

	// Integration with other modules
	lctKeeper    lctmanagerkeeper.Keeper
	energyKeeper energycyclekeeper.Keeper
	trustKeeper  trustkeeper.Keeper
	mrhKeeper    mrhkeeper.Keeper
}

// NewKeeper creates a new societytodo Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	lctKeeper lctmanagerkeeper.Keeper,
	energyKeeper energycyclekeeper.Keeper,
	trustKeeper trustkeeper.Keeper,
	mrhKeeper mrhkeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		lctKeeper:    lctKeeper,
		energyKeeper: energKeeper,
		trustKeeper:  trustKeeper,
		mrhKeeper:    mrhKeeper,
	}
}

// CreateSocietyTodoList initializes a new todo list for a society
func (k Keeper) CreateSocietyTodoList(ctx sdk.Context, societyLCT string) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSocietyTodoListKey(societyLCT)

	// Check if society todo list already exists
	if store.Has(key) {
		return fmt.Errorf("society todo list already exists for %s", societyLCT)
	}

	// Verify society LCT exists
	if !k.lctKeeper.HasLCT(ctx, societyLCT) {
		return fmt.Errorf("society LCT %s not found", societyLCT)
	}

	// Initialize todo list with default state
	todoList := types.SocietyTodoList{
		SocietyLct: societyLCT,
		ListId:     fmt.Sprintf("list_%s_%d", societyLCT, ctx.BlockHeight()),
		State:      types.SOCIETY_STATE_AWAKENING,
		CurrentCycle: types.CycleInfo{
			CycleNumber:        1,
			CycleStart:        ctx.BlockTime(),
			CycleDurationSeconds: 86400, // 24 hours default
			EnergyEfficiency:   sdk.NewDec(100),
		},
		AtpBudget: types.AtpBudget{
			TotalAllocated: sdk.NewInt(1000000), // 1M ATP initial
			Available:      sdk.NewInt(1000000),
			Reserved:       sdk.NewInt(100000),  // 10% reserved
			MaxPerTodo:     sdk.NewInt(10000),   // 10K ATP max per todo
			MaxPerCycle:    sdk.NewInt(500000),  // 500K ATP max per cycle
		},
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}

	// Store the todo list
	bz := k.cdc.MustMarshal(&todoList)
	store.Set(key, bz)

	return nil
}

// ProcessWakeSleepCycle manages society state transitions
func (k Keeper) ProcessWakeSleepCycle(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SocietyTodoListPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var todoList types.SocietyTodoList
		k.cdc.MustUnmarshal(iterator.Value(), &todoList)

		// Check ATP levels and adjust state
		atpLevel := todoList.AtpBudget.Available.Quo(todoList.AtpBudget.TotalAllocated)
		
		switch {
		case atpLevel.LT(sdk.NewDec(10)): // Less than 10% ATP
			if todoList.State != types.SOCIETY_STATE_HIBERNATING {
				k.TransitionToState(ctx, &todoList, types.SOCIETY_STATE_HIBERNATING)
			}
		case atpLevel.LT(sdk.NewDec(25)): // Less than 25% ATP
			if todoList.State != types.SOCIETY_STATE_SLEEPING {
				k.TransitionToState(ctx, &todoList, types.SOCIETY_STATE_SLEEPING)
			}
		case atpLevel.LT(sdk.NewDec(50)): // Less than 50% ATP
			if todoList.State == types.SOCIETY_STATE_ACTIVE {
				k.TransitionToState(ctx, &todoList, types.SOCIETY_STATE_CONSERVING)
			}
		case atpLevel.GT(sdk.NewDec(75)): // More than 75% ATP
			if todoList.State != types.SOCIETY_STATE_ACTIVE {
				k.TransitionToState(ctx, &todoList, types.SOCIETY_STATE_AWAKENING)
			}
		}

		// Update cycle info
		k.UpdateCycleMetrics(ctx, &todoList)

		// Store updated todo list
		bz := k.cdc.MustMarshal(&todoList)
		store.Set(iterator.Key(), bz)
	}
}

// TransitionToState changes society operational state
func (k Keeper) TransitionToState(ctx sdk.Context, todoList *types.SocietyTodoList, newState types.SocietyState) {
	oldState := todoList.State
	todoList.State = newState
	todoList.UpdatedAt = ctx.BlockTime()

	// Emit state transition event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"society_state_transition",
			sdk.NewAttribute("society_lct", todoList.SocietyLct),
			sdk.NewAttribute("old_state", oldState.String()),
			sdk.NewAttribute("new_state", newState.String()),
			sdk.NewAttribute("atp_available", todoList.AtpBudget.Available.String()),
		),
	)

	// Adjust energy consumption based on state
	switch newState {
	case types.SOCIETY_STATE_HIBERNATING:
		// Minimal energy consumption
		todoList.CurrentCycle.EnergyEfficiency = sdk.NewDec(10)
	case types.SOCIETY_STATE_SLEEPING:
		todoList.CurrentCycle.EnergyEfficiency = sdk.NewDec(25)
	case types.SOCIETY_STATE_CONSERVING:
		todoList.CurrentCycle.EnergyEfficiency = sdk.NewDec(50)
	case types.SOCIETY_STATE_ACTIVE:
		todoList.CurrentCycle.EnergyEfficiency = sdk.NewDec(100)
	case types.SOCIETY_STATE_AWAKENING:
		todoList.CurrentCycle.EnergyEfficiency = sdk.NewDec(75)
	}
}

// UpdateCycleMetrics updates performance metrics for the current cycle
func (k Keeper) UpdateCycleMetrics(ctx sdk.Context, todoList *types.SocietyTodoList) {
	cycle := &todoList.CurrentCycle
	
	// Calculate cycle duration
	elapsed := ctx.BlockTime().Sub(cycle.CycleStart)
	if elapsed.Seconds() >= float64(cycle.CycleDurationSeconds) {
		// Start new cycle
		k.StartNewCycle(ctx, todoList)
	}

	// Update efficiency based on completion rate
	if cycle.TodosCompleted+cycle.TodosFailed > 0 {
		completionRate := sdk.NewDec(int64(cycle.TodosCompleted)).Quo(
			sdk.NewDec(int64(cycle.TodosCompleted + cycle.TodosFailed)),
		)
		cycle.EnergyEfficiency = cycle.EnergyEfficiency.Mul(completionRate)
	}
}

// StartNewCycle initializes a new operational cycle
func (k Keeper) StartNewCycle(ctx sdk.Context, todoList *types.SocietyTodoList) {
	// Store previous cycle in history
	k.StoreCycleHistory(ctx, todoList.SocietyLct, todoList.CurrentCycle)

	// Initialize new cycle
	todoList.CurrentCycle = types.CycleInfo{
		CycleNumber:          todoList.CurrentCycle.CycleNumber + 1,
		CycleStart:          ctx.BlockTime(),
		CycleDurationSeconds: todoList.CurrentCycle.CycleDurationSeconds,
		EnergyEfficiency:     todoList.CurrentCycle.EnergyEfficiency,
		TodosCompleted:       0,
		TodosFailed:          0,
	}

	// Predict next wake/sleep times based on state
	switch todoList.State {
	case types.SOCIETY_STATE_ACTIVE:
		todoList.CurrentCycle.NextSleepTime = ctx.BlockTime().Add(time.Hour * 16) // 16 hours active
	case types.SOCIETY_STATE_SLEEPING, types.SOCIETY_STATE_HIBERNATING:
		todoList.CurrentCycle.NextWakeTime = ctx.BlockTime().Add(time.Hour * 8) // 8 hours sleep
	}
}

// StoreCycleHistory saves completed cycle data
func (k Keeper) StoreCycleHistory(ctx sdk.Context, societyLCT string, cycle types.CycleInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCycleHistoryKey(societyLCT, cycle.CycleNumber)
	bz := k.cdc.MustMarshal(&cycle)
	store.Set(key, bz)
}