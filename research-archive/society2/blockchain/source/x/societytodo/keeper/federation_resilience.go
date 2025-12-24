package keeper

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"racecar-web/x/societytodo/types"
)

// PeerHealth tracks the health status of federation peers
type PeerHealth struct {
	PeerID       string
	LastSeen     time.Time
	IsHealthy    bool
	FailureCount int
	TrustLevel   sdk.Dec
}

// MonitorFederationHealth checks peer connectivity and handles dropouts
func (k Keeper) MonitorFederationHealth(ctx sdk.Context) {
	peers := k.GetKnownPeers(ctx)
	currentTime := ctx.BlockTime()
	
	for _, peer := range peers {
		health := k.GetPeerHealth(ctx, peer.ID)
		
		// Check if peer has timed out
		timeSinceLastSeen := currentTime.Sub(health.LastSeen).Seconds()
		timeoutThreshold := float64(120) // 2 minutes
		
		if timeSinceLastSeen > timeoutThreshold {
			if health.IsHealthy {
				// Peer just went unhealthy
				k.HandlePeerDropout(ctx, peer)
				health.IsHealthy = false
				health.FailureCount++
			}
		} else {
			if !health.IsHealthy {
				// Peer recovered
				k.HandlePeerRecovery(ctx, peer)
				health.IsHealthy = true
			}
		}
		
		k.SetPeerHealth(ctx, health)
	}
}

// HandlePeerDropout manages TODO migration when a peer goes offline
func (k Keeper) HandlePeerDropout(ctx sdk.Context, peer types.FederationPeer) {
	k.Logger(ctx).Info("Peer dropout detected", "peer_id", peer.ID)
	
	// Find todos assigned to this peer
	todos := k.GetTodosByExecutor(ctx, peer.ID)
	
	for _, todo := range todos {
		if todo.Status == types.TodoStatusInProgress {
			// Check if grace period has expired
			gracePeriod := time.Duration(120) * time.Second
			timeSinceStart := ctx.BlockTime().Sub(*todo.StartedAt)
			
			if timeSinceStart > gracePeriod {
				// Migrate todo to another peer
				k.MigrateTodo(ctx, todo)
			} else {
				// Mark for monitoring
				k.AddToWatchList(ctx, todo.ID)
			}
		}
	}
	
	// Update federation state
	k.UpdateFederationQuorum(ctx)
	
	// Emit dropout event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"federation_peer_dropout",
			sdk.NewAttribute("peer_id", peer.ID),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
			sdk.NewAttribute("todos_affected", fmt.Sprintf("%d", len(todos))),
		),
	)
}

// HandlePeerRecovery reconciles state when a peer comes back online
func (k Keeper) HandlePeerRecovery(ctx sdk.Context, peer types.FederationPeer) {
	k.Logger(ctx).Info("Peer recovery detected", "peer_id", peer.ID)
	
	// Request state hash from recovered peer
	stateHash := k.RequestStateHash(ctx, peer.ID)
	localHash := k.ComputeLocalStateHash(ctx)
	
	if stateHash != localHash {
		// State divergence detected
		k.InitiateStateReconciliation(ctx, peer.ID)
	}
	
	// Update federation state
	k.UpdateFederationQuorum(ctx)
	
	// Emit recovery event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"federation_peer_recovery",
			sdk.NewAttribute("peer_id", peer.ID),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
			sdk.NewAttribute("state_match", fmt.Sprintf("%v", stateHash == localHash)),
		),
	)
}

// MigrateTodo transfers a todo to another available society
func (k Keeper) MigrateTodo(ctx sdk.Context, todo types.TodoItem) error {
	// Find eligible societies
	candidates := k.FindEligibleExecutors(ctx, todo)
	
	if len(candidates) == 0 {
		// No candidates available, enter degraded mode
		k.EnterDegradedMode(ctx, todo)
		return fmt.Errorf("no eligible executors for todo %s", todo.ID)
	}
	
	// Select best candidate based on trust and capacity
	selected := k.SelectBestCandidate(ctx, candidates, todo)
	
	// Create checkpoint of current progress
	checkpoint := types.TodoCheckpoint{
		TodoID:      todo.ID,
		Timestamp:   ctx.BlockTime(),
		Progress:    k.EstimateProgress(ctx, todo),
		FromSociety: todo.ExecutedBy,
		ToSociety:   selected.ID,
		Reason:      "peer_dropout",
	}
	k.SaveCheckpoint(ctx, checkpoint)
	
	// Update todo assignment
	todo.ExecutedBy = selected.ID
	todo.MigrationCount++
	k.SetTodoItem(ctx, todo)
	
	// Transfer ATP allocation
	k.TransferATPAllocation(ctx, todo, selected.ID)
	
	// Emit migration event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"todo_migration",
			sdk.NewAttribute("todo_id", todo.ID),
			sdk.NewAttribute("from_society", checkpoint.FromSociety),
			sdk.NewAttribute("to_society", checkpoint.ToSociety),
			sdk.NewAttribute("reason", checkpoint.Reason),
		),
	)
	
	return nil
}

// UpdateFederationQuorum adjusts consensus requirements based on available peers
func (k Keeper) UpdateFederationQuorum(ctx sdk.Context) {
	healthyPeers := k.CountHealthyPeers(ctx)
	totalPeers := k.CountTotalPeers(ctx)
	
	quorum := float64(healthyPeers) / float64(totalPeers)
	
	var mode types.ConsensusMode
	switch {
	case quorum >= 0.75:
		mode = types.ConsensusModeNormal
	case quorum >= 0.5:
		mode = types.ConsensusModeReduced
	case quorum >= 0.25:
		mode = types.ConsensusModeDegraded
	default:
		mode = types.ConsensusModeEmergency
	}
	
	k.SetConsensusMode(ctx, mode)
	
	// Adjust TODO acceptance based on mode
	switch mode {
	case types.ConsensusModeDegraded:
		k.SetAcceptedPriorities(ctx, []types.TodoPriority{types.PriorityCritical})
	case types.ConsensusModeEmergency:
		k.SetAcceptingRequests(ctx, false)
	}
}

// InitiateStateReconciliation synchronizes state with a recovered peer
func (k Keeper) InitiateStateReconciliation(ctx sdk.Context, peerID string) {
	// Get divergence point
	divergenceBlock := k.FindDivergencePoint(ctx, peerID)
	
	// Get todos created after divergence
	divergentTodos := k.GetTodosAfterBlock(ctx, divergenceBlock)
	
	// Exchange todo lists with peer
	peerTodos := k.RequestTodoList(ctx, peerID, divergenceBlock)
	
	// Reconcile based on consensus
	reconciledTodos := k.ReconcileTodos(ctx, divergentTodos, peerTodos)
	
	// Apply reconciled state
	for _, todo := range reconciledTodos {
		k.SetTodoItem(ctx, todo)
	}
	
	k.Logger(ctx).Info("State reconciliation completed",
		"peer_id", peerID,
		"divergence_block", divergenceBlock,
		"todos_reconciled", len(reconciledTodos),
	)
}

// EnterDegradedMode switches to limited operations during low connectivity
func (k Keeper) EnterDegradedMode(ctx sdk.Context, todo types.TodoItem) {
	todoList, _ := k.GetSocietyTodoList(ctx, todo.SocietyLCT)
	
	// Increase ATP costs for new requests
	k.SetDegradedMultiplier(ctx, sdk.NewDec(2))
	
	// Queue non-critical todos
	if todo.Priority != types.PriorityCritical {
		k.AddToQueue(ctx, todo)
		todo.Status = types.TodoStatusQueued
		k.SetTodoItem(ctx, todo)
	}
	
	// Log degraded operation
	k.LogDegradedOperation(ctx, types.DegradedOperation{
		Timestamp:   ctx.BlockTime(),
		TodoID:      todo.ID,
		Action:      "queued_for_quorum",
		QuorumLevel: k.GetCurrentQuorum(ctx),
	})
	
	// Update society state
	if todoList.State != types.StateConserving {
		k.TransitionToState(ctx, &todoList, types.StateConserving)
	}
}

// SaveCheckpoint creates a recovery point for a todo
func (k Keeper) SaveCheckpoint(ctx sdk.Context, checkpoint types.TodoCheckpoint) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCheckpointKey(checkpoint.TodoID, checkpoint.Timestamp)
	bz := k.cdc.MustMarshal(&checkpoint)
	store.Set(key, bz)
}

// FindEligibleExecutors returns societies that can handle a todo
func (k Keeper) FindEligibleExecutors(ctx sdk.Context, todo types.TodoItem) []types.FederationPeer {
	var eligible []types.FederationPeer
	peers := k.GetHealthyPeers(ctx)
	
	for _, peer := range peers {
		// Check trust level
		if peer.TrustLevel.LT(todo.MinTrustLevel) {
			continue
		}
		
		// Check capacity
		capacity := k.GetPeerCapacity(ctx, peer.ID)
		if capacity.ActiveTodos >= capacity.MaxConcurrent {
			continue
		}
		
		// Check ATP availability
		atpAvailable := k.GetPeerATPBalance(ctx, peer.ID)
		if atpAvailable.LT(todo.ATPCost) {
			continue
		}
		
		eligible = append(eligible, peer)
	}
	
	return eligible
}

// SelectBestCandidate chooses optimal executor based on multiple factors
func (k Keeper) SelectBestCandidate(ctx sdk.Context, candidates []types.FederationPeer, todo types.TodoItem) types.FederationPeer {
	var bestPeer types.FederationPeer
	bestScore := sdk.ZeroDec()
	
	for _, peer := range candidates {
		// Calculate score based on:
		// - Trust level (40%)
		// - Available capacity (30%)
		// - Historical success rate (20%)
		// - Network latency (10%)
		
		trustScore := peer.TrustLevel.Mul(sdk.NewDec(40))
		
		capacity := k.GetPeerCapacity(ctx, peer.ID)
		capacityRatio := sdk.NewDec(capacity.MaxConcurrent - capacity.ActiveTodos).
			Quo(sdk.NewDec(capacity.MaxConcurrent))
		capacityScore := capacityRatio.Mul(sdk.NewDec(30))
		
		successRate := k.GetPeerSuccessRate(ctx, peer.ID)
		successScore := successRate.Mul(sdk.NewDec(20))
		
		latency := k.GetPeerLatency(ctx, peer.ID)
		latencyScore := sdk.NewDec(100).Sub(sdk.NewDec(latency)).
			Quo(sdk.NewDec(100)).Mul(sdk.NewDec(10))
		
		totalScore := trustScore.Add(capacityScore).Add(successScore).Add(latencyScore)
		
		if totalScore.GT(bestScore) {
			bestScore = totalScore
			bestPeer = peer
		}
	}
	
	return bestPeer
}