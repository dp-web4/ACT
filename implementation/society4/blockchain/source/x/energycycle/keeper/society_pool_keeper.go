package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"racecar-web/x/energycycle/types"
)

// InitializeSociety4Pool creates the initial token pool for Society 4
// Implements LAW-ECON-001: Total ATP Budget (1000)
// Implements LAW-ECON-002: Security Queen ATP Allocation (150)
func (k Keeper) InitializeSociety4Pool(ctx context.Context, societyLCT string) error {
	// Create new pool
	pool := types.NewSocietyTokenPool(societyLCT)

	// Allocate to roles per Society4Roles definition
	for _, role := range types.Society4Roles {
		// Generate role LCT ID (in production, these would be actual LCT IDs)
		roleLCT := fmt.Sprintf("lct:web4:society4:role:%s", role.Name)

		err := pool.AllocateToRole(roleLCT, role.Name, role.InitialATP)
		if err != nil {
			return fmt.Errorf("failed to allocate ATP to %s: %w", role.Name, err)
		}
	}

	// Validate pool integrity before storing
	if err := pool.ValidatePoolIntegrity(); err != nil {
		return fmt.Errorf("pool integrity validation failed: %w", err)
	}

	// Store the pool
	if err := k.SocietyTokenPools.Set(ctx, societyLCT, *pool); err != nil {
		return fmt.Errorf("failed to store society pool: %w", err)
	}

	return nil
}

// GetSocietyPool retrieves a society's token pool
func (k Keeper) GetSocietyPool(ctx context.Context, societyLCT string) (types.SocietyTokenPool, error) {
	pool, err := k.SocietyTokenPools.Get(ctx, societyLCT)
	if err != nil {
		return types.SocietyTokenPool{}, fmt.Errorf("society pool not found: %w", err)
	}
	return pool, nil
}

// DischargeRoleATP discharges ATP for a role operation
// Implements ATP → ADP discharge cycle
func (k Keeper) DischargeRoleATP(ctx context.Context, societyLCT, roleLCT string, amount int64, operationID, reason string) error {
	// Get pool
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return err
	}

	// Discharge ATP
	tx, err := pool.DischargeATP(roleLCT, amount, operationID, reason)
	if err != nil {
		return fmt.Errorf("failed to discharge ATP: %w", err)
	}

	// Store transaction
	if err := k.AtpTransactions.Set(ctx, tx.TransactionID, *tx); err != nil {
		return fmt.Errorf("failed to store ATP transaction: %w", err)
	}

	// Update pool
	if err := k.SocietyTokenPools.Set(ctx, societyLCT, pool); err != nil {
		return fmt.Errorf("failed to update society pool: %w", err)
	}

	return nil
}

// TransferRoleATP transfers ATP between roles
func (k Keeper) TransferRoleATP(ctx context.Context, societyLCT, fromRole, toRole string, amount int64, operationID, reason string) error {
	// Get pool
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return err
	}

	// Transfer ATP
	tx, err := pool.TransferATP(fromRole, toRole, amount, operationID, reason)
	if err != nil {
		return fmt.Errorf("failed to transfer ATP: %w", err)
	}

	// Store transaction
	if err := k.AtpTransactions.Set(ctx, tx.TransactionID, *tx); err != nil {
		return fmt.Errorf("failed to store ATP transaction: %w", err)
	}

	// Update pool
	if err := k.SocietyTokenPools.Set(ctx, societyLCT, pool); err != nil {
		return fmt.Errorf("failed to update society pool: %w", err)
	}

	return nil
}

// RechargeRoleADP recharges ADP back to ATP (value creation)
func (k Keeper) RechargeRoleADP(ctx context.Context, societyLCT, roleLCT string, amount int64, operationID, reason string) error {
	// Get pool
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return err
	}

	// Recharge ADP
	tx, err := pool.RechargeADP(roleLCT, amount, operationID, reason)
	if err != nil {
		return fmt.Errorf("failed to recharge ADP: %w", err)
	}

	// Store transaction
	if err := k.AtpTransactions.Set(ctx, tx.TransactionID, *tx); err != nil {
		return fmt.Errorf("failed to store ATP transaction: %w", err)
	}

	// Update pool
	if err := k.SocietyTokenPools.Set(ctx, societyLCT, pool); err != nil {
		return fmt.Errorf("failed to update society pool: %w", err)
	}

	return nil
}

// PerformDailyRecharge performs daily ATP recharge for all roles
// Implements LAW-ECON-003: Daily ATP Recharge
// Implements PROC-ATP-RECHARGE: Daily 00:00 UTC regeneration
func (k Keeper) PerformDailyRecharge(ctx context.Context, societyLCT string) (map[string]int64, error) {
	// Get pool
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return nil, err
	}

	// Perform daily recharge
	recharged, err := pool.DailyRecharge()
	if err != nil {
		return nil, fmt.Errorf("failed to perform daily recharge: %w", err)
	}

	// Update pool
	if err := k.SocietyTokenPools.Set(ctx, societyLCT, pool); err != nil {
		return nil, fmt.Errorf("failed to update society pool: %w", err)
	}

	// Create transaction records for each recharge
	for roleLCT, amount := range recharged {
		tx := &types.AtpTransaction{
			TransactionID: fmt.Sprintf("tx-daily-recharge-%s-%d", roleLCT, pool.LastRecharge.Unix()),
			SocietyID:     societyLCT,
			FromRole:      "system",
			ToRole:        roleLCT,
			Amount:        amount,
			Type:          "daily_recharge",
			OperationID:   "PROC-ATP-RECHARGE",
			Reason:        "Daily 00:00 UTC regeneration per LAW-ECON-003",
			Timestamp:     pool.LastRecharge,
			ResultingATP:  pool.RoleBalances[roleLCT],
		}

		if err := k.AtpTransactions.Set(ctx, tx.TransactionID, *tx); err != nil {
			return nil, fmt.Errorf("failed to store recharge transaction: %w", err)
		}
	}

	return recharged, nil
}

// GetRoleBalance returns current ATP and ADP balance for a role
func (k Keeper) GetRoleBalance(ctx context.Context, societyLCT, roleLCT string) (atp int64, adp int64, error error) {
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return 0, 0, err
	}

	return pool.GetRoleBalance(roleLCT)
}

// ValidatePoolIntegrity validates society pool integrity
// Implements LAW-ECON-001: Total ATP Budget conservation
func (k Keeper) ValidatePoolIntegrity(ctx context.Context, societyLCT string) error {
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return err
	}

	return pool.ValidatePoolIntegrity()
}

// GetPoolSummary returns a summary of the pool state
func (k Keeper) GetPoolSummary(ctx context.Context, societyLCT string) (string, error) {
	pool, err := k.GetSocietyPool(ctx, societyLCT)
	if err != nil {
		return "", err
	}

	return pool.GetPoolSummary(), nil
}

// EnforceAtpStakeForTrustQuery enforces LAW-ECON-004
// Trust tensor queries require minimum 5 ATP stake (privacy)
func (k Keeper) EnforceAtpStakeForTrustQuery(ctx context.Context, societyLCT, roleLCT string, queryType string) error {
	minStake := int64(5) // LAW-ECON-004

	// Get role balance
	atp, _, err := k.GetRoleBalance(ctx, societyLCT, roleLCT)
	if err != nil {
		return err
	}

	// Check if role has sufficient ATP for stake
	if atp < minStake {
		return fmt.Errorf("insufficient ATP for trust query: %d < %d (LAW-ECON-004)", atp, minStake)
	}

	// Discharge ATP as privacy stake
	operationID := fmt.Sprintf("trust-query-%s", queryType)
	reason := fmt.Sprintf("Privacy stake for %s trust query per LAW-ECON-004", queryType)

	if err := k.DischargeRoleATP(ctx, societyLCT, roleLCT, minStake, operationID, reason); err != nil {
		return fmt.Errorf("failed to stake ATP: %w", err)
	}

	return nil
}

// GetAllTransactions retrieves all ATP transactions for a society
func (k Keeper) GetAllTransactions(ctx context.Context, societyLCT string) ([]types.AtpTransaction, error) {
	var transactions []types.AtpTransaction

	iter, err := k.AtpTransactions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		tx, err := iter.Value()
		if err != nil {
			continue
		}

		if tx.SocietyID == societyLCT {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// GetRoleTransactions retrieves all ATP transactions for a specific role
func (k Keeper) GetRoleTransactions(ctx context.Context, societyLCT, roleLCT string) ([]types.AtpTransaction, error) {
	var transactions []types.AtpTransaction

	iter, err := k.AtpTransactions.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		tx, err := iter.Value()
		if err != nil {
			continue
		}

		if tx.SocietyID == societyLCT && (tx.FromRole == roleLCT || tx.ToRole == roleLCT) {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}
