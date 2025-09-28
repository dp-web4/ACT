package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"society4chain/x/societytodo/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RequestTodo handles citizen todo requests
func (k msgServer) RequestTodo(goCtx context.Context, msg *types.MsgRequestTodo) (*types.MsgRequestTodoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify citizen LCT exists
	if !k.lctKeeper.HasLCT(ctx, msg.CitizenLct) {
		return nil, fmt.Errorf("citizen LCT %s not found", msg.CitizenLct)
	}

	// Get society todo list
	todoList, err := k.GetSocietyTodoList(ctx, msg.SocietyLct)
	if err != nil {
		return nil, err
	}

	// Check society state allows new requests
	if !k.CanAcceptRequest(todoList.State) {
		return nil, fmt.Errorf("society in state %s cannot accept new requests", todoList.State)
	}

	// Calculate ATP cost based on priority and complexity
	atpCost := k.CalculateTodoCost(msg.Priority, msg.EstimatedComplexity)

	// Check if citizen has sufficient ATP
	citizenATP, err := k.energyKeeper.GetBalance(ctx, msg.CitizenLct)
	if err != nil {
		return nil, err
	}

	if citizenATP.LT(atpCost) {
		return nil, fmt.Errorf("insufficient ATP: need %s, have %s", atpCost, citizenATP)
	}

	// Check trust level for priority requests
	if msg.Priority == types.PRIORITY_CRITICAL {
		trustLevel, err := k.trustKeeper.GetTrustTensor(ctx, msg.CitizenLct, msg.SocietyLct)
		if err != nil {
			return nil, err
		}
		if trustLevel.T3.LT(sdk.NewDec(80)) { // Require 80+ trust for critical
			return nil, fmt.Errorf("insufficient trust level for critical priority: %s", trustLevel.T3)
		}
	}

	// Create citizen request
	requestID := k.GenerateRequestID(ctx)
	request := types.CitizenRequest{
		RequestId:           requestID,
		CitizenLct:          msg.CitizenLct,
		SocietyLct:          msg.SocietyLct,
		Title:               msg.Title,
		Description:         msg.Description,
		Priority:            msg.Priority,
		EstimatedComplexity: msg.EstimatedComplexity,
		AtpOffered:          atpCost,
		Status:              types.REQUEST_STATUS_PENDING,
		RequestedAt:         ctx.BlockTime(),
	}

	// Store request
	k.SetCitizenRequest(ctx, request)

	// Deduct ATP from citizen
	err = k.energyKeeper.DeductATP(ctx, msg.CitizenLct, atpCost)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"citizen_todo_request",
			sdk.NewAttribute("request_id", requestID),
			sdk.NewAttribute("citizen", msg.CitizenLct),
			sdk.NewAttribute("society", msg.SocietyLct),
			sdk.NewAttribute("priority", msg.Priority.String()),
			sdk.NewAttribute("atp_cost", atpCost.String()),
		),
	)

	return &types.MsgRequestTodoResponse{
		RequestId: requestID,
		AtpCost:   atpCost,
		Status:    "pending",
	}, nil
}

// DelegateATP handles ATP delegation to todo pools
func (k msgServer) DelegateATP(goCtx context.Context, msg *types.MsgDelegateATP) (*types.MsgDelegateATPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify citizen LCT
	if !k.lctKeeper.HasLCT(ctx, msg.CitizenLct) {
		return nil, fmt.Errorf("citizen LCT %s not found", msg.CitizenLct)
	}

	// Get delegation pool
	pool, err := k.GetDelegationPool(ctx, msg.SocietyLct, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Check if pool is accepting delegations
	if pool.Status != types.POOL_STATUS_ACTIVE {
		return nil, fmt.Errorf("pool %s is not active", msg.PoolId)
	}

	// Check citizen ATP balance
	balance, err := k.energyKeeper.GetBalance(ctx, msg.CitizenLct)
	if err != nil {
		return nil, err
	}

	if balance.LT(msg.Amount) {
		return nil, fmt.Errorf("insufficient ATP balance: have %s, delegating %s", balance, msg.Amount)
	}

	// Calculate delegation weight based on trust
	trustLevel, err := k.trustKeeper.GetTrustTensor(ctx, msg.CitizenLct, msg.SocietyLct)
	if err != nil {
		return nil, err
	}

	// Apply quadratic voting with trust weighting
	votingPower := k.CalculateVotingPower(msg.Amount, trustLevel.T3)

	// Create or update delegation
	delegation := types.AtpDelegation{
		CitizenLct:    msg.CitizenLct,
		PoolId:        msg.PoolId,
		Amount:        msg.Amount,
		VotingPower:   votingPower,
		DelegatedAt:   ctx.BlockTime(),
		LockedUntil:   ctx.BlockTime().Add(pool.MinLockPeriod),
	}

	// Store delegation
	k.SetDelegation(ctx, delegation)

	// Transfer ATP to pool
	err = k.energyKeeper.TransferToPool(ctx, msg.CitizenLct, msg.PoolId, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Update pool metrics
	pool.TotalDelegated = pool.TotalDelegated.Add(msg.Amount)
	pool.TotalVotingPower = pool.TotalVotingPower.Add(votingPower)
	k.SetDelegationPool(ctx, pool)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"atp_delegated",
			sdk.NewAttribute("citizen", msg.CitizenLct),
			sdk.NewAttribute("pool", msg.PoolId),
			sdk.NewAttribute("amount", msg.Amount.String()),
			sdk.NewAttribute("voting_power", votingPower.String()),
		),
	)

	return &types.MsgDelegateATPResponse{
		VotingPower: votingPower,
		LockedUntil: delegation.LockedUntil,
	}, nil
}

// ProcessTodo handles todo execution requests
func (k msgServer) ProcessTodo(goCtx context.Context, msg *types.MsgProcessTodo) (*types.MsgProcessTodoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get todo item
	todo, err := k.GetTodoItem(ctx, msg.TodoId)
	if err != nil {
		return nil, err
	}

	// Verify executor authorization
	if !k.IsAuthorizedExecutor(ctx, msg.ExecutorLct, todo) {
		return nil, fmt.Errorf("executor %s not authorized for todo %s", msg.ExecutorLct, msg.TodoId)
	}

	// Check society state allows processing
	todoList, err := k.GetSocietyTodoList(ctx, todo.SocietyLct)
	if err != nil {
		return nil, err
	}

	if !k.CanProcessTodo(todoList.State, todo.Priority) {
		return nil, fmt.Errorf("cannot process todo in current society state %s", todoList.State)
	}

	// Process based on action
	switch msg.Action {
	case types.ACTION_START:
		if todo.Status != types.TODO_STATUS_PENDING {
			return nil, fmt.Errorf("todo %s is not pending", msg.TodoId)
		}
		todo.Status = types.TODO_STATUS_IN_PROGRESS
		todo.ExecutorLct = msg.ExecutorLct
		todo.StartedAt = ctx.BlockTime()

	case types.ACTION_COMPLETE:
		if todo.Status != types.TODO_STATUS_IN_PROGRESS {
			return nil, fmt.Errorf("todo %s is not in progress", msg.TodoId)
		}
		todo.Status = types.TODO_STATUS_COMPLETED
		todo.CompletedAt = ctx.BlockTime()
		todo.CompletionProof = msg.Proof

		// Generate ADP based on quality
		adpAmount := k.CalculateADPReward(todo.AtpAllocated, msg.QualityScore)
		err = k.energyKeeper.GenerateADP(ctx, todo.ExecutorLct, adpAmount)
		if err != nil {
			return nil, err
		}

		// Update cycle metrics
		todoList.CurrentCycle.TodosCompleted++
		k.SetSocietyTodoList(ctx, todoList)

	case types.ACTION_FAIL:
		if todo.Status != types.TODO_STATUS_IN_PROGRESS {
			return nil, fmt.Errorf("todo %s is not in progress", msg.TodoId)
		}
		todo.Status = types.TODO_STATUS_FAILED
		todo.FailureReason = msg.Reason

		// Return partial ATP to pool
		refund := todo.AtpAllocated.Mul(sdk.NewDec(50)).Quo(sdk.NewDec(100)) // 50% refund
		todoList.AtpBudget.Available = todoList.AtpBudget.Available.Add(refund)
		todoList.CurrentCycle.TodosFailed++
		k.SetSocietyTodoList(ctx, todoList)

	default:
		return nil, fmt.Errorf("unknown action %s", msg.Action)
	}

	// Update todo
	k.SetTodoItem(ctx, todo)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"todo_processed",
			sdk.NewAttribute("todo_id", msg.TodoId),
			sdk.NewAttribute("action", msg.Action.String()),
			sdk.NewAttribute("executor", msg.ExecutorLct),
			sdk.NewAttribute("status", todo.Status.String()),
		),
	)

	return &types.MsgProcessTodoResponse{
		Status: todo.Status.String(),
		AdpGenerated: adpAmount,
	}, nil
}

// Helper functions

func (k Keeper) CanAcceptRequest(state types.SocietyState) bool {
	return state == types.SOCIETY_STATE_ACTIVE || 
	       state == types.SOCIETY_STATE_AWAKENING ||
	       state == types.SOCIETY_STATE_CONSERVING
}

func (k Keeper) CanProcessTodo(state types.SocietyState, priority types.Priority) bool {
	switch state {
	case types.SOCIETY_STATE_ACTIVE, types.SOCIETY_STATE_AWAKENING:
		return true
	case types.SOCIETY_STATE_CONSERVING:
		return priority >= types.PRIORITY_HIGH
	case types.SOCIETY_STATE_SLEEPING:
		return priority == types.PRIORITY_CRITICAL
	default:
		return false
	}
}

func (k Keeper) CalculateTodoCost(priority types.Priority, complexity uint32) sdk.Int {
	baseCost := sdk.NewInt(1000)
	
	// Priority multiplier
	switch priority {
	case types.PRIORITY_CRITICAL:
		baseCost = baseCost.Mul(sdk.NewInt(5))
	case types.PRIORITY_HIGH:
		baseCost = baseCost.Mul(sdk.NewInt(3))
	case types.PRIORITY_MEDIUM:
		baseCost = baseCost.Mul(sdk.NewInt(2))
	}
	
	// Complexity multiplier
	baseCost = baseCost.Mul(sdk.NewInt(int64(complexity)))
	
	return baseCost
}

func (k Keeper) CalculateVotingPower(amount sdk.Int, trustLevel sdk.Dec) sdk.Dec {
	// Quadratic voting: sqrt(amount) * trust_multiplier
	amountDec := sdk.NewDecFromInt(amount)
	sqrtAmount := amountDec.ApproxSqrt()
	trustMultiplier := trustLevel.Quo(sdk.NewDec(100)) // Normalize to 0-1
	return sqrtAmount.Mul(trustMultiplier)
}

func (k Keeper) CalculateADPReward(atpSpent sdk.Int, qualityScore sdk.Dec) sdk.Int {
	// ADP = ATP * quality_score * efficiency_bonus
	baseReward := sdk.NewDecFromInt(atpSpent).Mul(qualityScore)
	efficiencyBonus := sdk.NewDec(110).Quo(sdk.NewDec(100)) // 10% bonus
	totalReward := baseReward.Mul(efficiencyBonus)
	return totalReward.TruncateInt()
}