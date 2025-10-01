package types

import (
	"fmt"
	"time"
)

// SocietyTokenPool represents Society 4's ATP/ADP treasury
// Per LAW-ECON-001: Total ATP budget is 1000
type SocietyTokenPool struct {
	SocietyID       string            `json:"society_id"`        // LCT ID of society
	TotalATP        int64             `json:"total_atp"`         // Fixed at 1000 per LAW-ECON-001
	AllocatedATP    int64             `json:"allocated_atp"`     // Currently allocated to roles
	AvailableATP    int64             `json:"available_atp"`     // Unallocated ATP
	TotalADP        int64             `json:"total_adp"`         // Discharged ATP
	RoleAllocations map[string]int64  `json:"role_allocations"`  // ATP allocated per role
	RoleBalances    map[string]int64  `json:"role_balances"`     // Current ATP balance per role
	ADPBalances     map[string]int64  `json:"adp_balances"`      // ADP balance per role
	LastRecharge    time.Time         `json:"last_recharge"`     // Last daily recharge timestamp
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Version         int64             `json:"version"`
}

// RoleAllocation represents ATP allocation for a role
// Per LAW-ECON-002: Security Queen gets 150 ATP (highest)
type RoleAllocation struct {
	RoleID          string    `json:"role_id"`           // Role LCT ID
	RoleName        string    `json:"role_name"`         // Human-readable name
	InitialATP      int64     `json:"initial_atp"`       // Initial allocation
	CurrentATP      int64     `json:"current_atp"`       // Current ATP balance
	CurrentADP      int64     `json:"current_adp"`       // Current ADP balance
	DailyRecharge   int64     `json:"daily_recharge"`    // Daily recharge amount (20 per LAW-ECON-003)
	LastRecharge    time.Time `json:"last_recharge"`
	TotalSpent      int64     `json:"total_spent"`       // Total ATP spent (discharged)
	TotalRecharged  int64     `json:"total_recharged"`   // Total ATP recharged
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// AtpTransaction represents an ATP transfer or discharge
type AtpTransaction struct {
	TransactionID   string    `json:"transaction_id"`
	SocietyID       string    `json:"society_id"`
	FromRole        string    `json:"from_role"`         // Source role LCT ID
	ToRole          string    `json:"to_role,omitempty"` // Target role LCT ID (if transfer)
	Amount          int64     `json:"amount"`
	Type            string    `json:"type"`              // "discharge", "transfer", "recharge"
	OperationID     string    `json:"operation_id"`      // Associated operation
	Reason          string    `json:"reason"`
	BlockHeight     int64     `json:"block_height"`
	Timestamp       time.Time `json:"timestamp"`
	ResultingATP    int64     `json:"resulting_atp"`     // Balance after transaction
	ResultingADP    int64     `json:"resulting_adp"`     // ADP balance after transaction
}

// Society4Roles defines the 8 queens + king per governance structure
var Society4Roles = []struct {
	Name          string
	InitialATP    int64
	DailyRecharge int64
}{
	{"King Claudius", 100, 20},           // Monarch
	{"Security Queen", 150, 20},          // LAW-ECON-002: Highest allocation
	{"Law Oracle Queen", 120, 20},        // Critical governance role
	{"Treasury Queen", 130, 20},          // Economic management
	{"Hardware Binding Queen", 110, 20},  // Identity management
	{"Federation Bridge Queen", 100, 20}, // External coordination
	{"Reality Cache Queen", 90, 20},      // Knowledge management
	{"Consensus Queen", 100, 20},         // Decision coordination
	{"Temporal Auth Queen", 100, 20},     // Time-based security
}

// NewSocietyTokenPool creates a new token pool for Society 4
func NewSocietyTokenPool(societyID string) *SocietyTokenPool {
	pool := &SocietyTokenPool{
		SocietyID:       societyID,
		TotalATP:        1000, // LAW-ECON-001: Fixed total
		AllocatedATP:    0,
		AvailableATP:    1000,
		TotalADP:        0,
		RoleAllocations: make(map[string]int64),
		RoleBalances:    make(map[string]int64),
		ADPBalances:     make(map[string]int64),
		LastRecharge:    time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Version:         1,
	}

	return pool
}

// AllocateToRole allocates initial ATP to a role
func (p *SocietyTokenPool) AllocateToRole(roleLCT string, roleName string, amount int64) error {
	// Check if we have enough available ATP
	if p.AvailableATP < amount {
		return fmt.Errorf("insufficient available ATP: %d < %d", p.AvailableATP, amount)
	}

	// Check if total would exceed LAW-ECON-001 limit
	if p.AllocatedATP+amount > p.TotalATP {
		return fmt.Errorf("allocation would exceed total ATP budget: %d + %d > %d", p.AllocatedATP, amount, p.TotalATP)
	}

	// Allocate ATP
	p.RoleAllocations[roleLCT] = amount
	p.RoleBalances[roleLCT] = amount
	p.ADPBalances[roleLCT] = 0
	p.AllocatedATP += amount
	p.AvailableATP -= amount
	p.UpdatedAt = time.Now()
	p.Version++

	return nil
}

// DischargeATP discharges ATP to ADP for a role
func (p *SocietyTokenPool) DischargeATP(roleLCT string, amount int64, operationID string, reason string) (*AtpTransaction, error) {
	// Check if role exists
	balance, exists := p.RoleBalances[roleLCT]
	if !exists {
		return nil, fmt.Errorf("role not found: %s", roleLCT)
	}

	// Check if role has sufficient ATP
	if balance < amount {
		return nil, fmt.Errorf("insufficient ATP balance for role %s: %d < %d", roleLCT, balance, amount)
	}

	// Discharge ATP to ADP
	p.RoleBalances[roleLCT] -= amount
	p.ADPBalances[roleLCT] += amount
	p.TotalADP += amount
	p.UpdatedAt = time.Now()
	p.Version++

	// Create transaction record
	tx := &AtpTransaction{
		TransactionID: fmt.Sprintf("tx-%s-%d", roleLCT, time.Now().Unix()),
		SocietyID:     p.SocietyID,
		FromRole:      roleLCT,
		Amount:        amount,
		Type:          "discharge",
		OperationID:   operationID,
		Reason:        reason,
		Timestamp:     time.Now(),
		ResultingATP:  p.RoleBalances[roleLCT],
		ResultingADP:  p.ADPBalances[roleLCT],
	}

	return tx, nil
}

// TransferATP transfers ATP between roles
func (p *SocietyTokenPool) TransferATP(fromRole string, toRole string, amount int64, operationID string, reason string) (*AtpTransaction, error) {
	// Check if both roles exist
	fromBalance, fromExists := p.RoleBalances[fromRole]
	if !fromExists {
		return nil, fmt.Errorf("source role not found: %s", fromRole)
	}

	_, toExists := p.RoleBalances[toRole]
	if !toExists {
		return nil, fmt.Errorf("target role not found: %s", toRole)
	}

	// Check if source has sufficient ATP
	if fromBalance < amount {
		return nil, fmt.Errorf("insufficient ATP balance for role %s: %d < %d", fromRole, fromBalance, amount)
	}

	// Transfer ATP
	p.RoleBalances[fromRole] -= amount
	p.RoleBalances[toRole] += amount
	p.UpdatedAt = time.Now()
	p.Version++

	// Create transaction record
	tx := &AtpTransaction{
		TransactionID: fmt.Sprintf("tx-%s-%s-%d", fromRole, toRole, time.Now().Unix()),
		SocietyID:     p.SocietyID,
		FromRole:      fromRole,
		ToRole:        toRole,
		Amount:        amount,
		Type:          "transfer",
		OperationID:   operationID,
		Reason:        reason,
		Timestamp:     time.Now(),
		ResultingATP:  p.RoleBalances[fromRole],
	}

	return tx, nil
}

// DailyRecharge recharges all roles per LAW-ECON-003
// All queens receive 20 ATP daily, capped at initial allocation
func (p *SocietyTokenPool) DailyRecharge() (map[string]int64, error) {
	now := time.Now()

	// Check if 24 hours have passed since last recharge
	if now.Sub(p.LastRecharge) < 24*time.Hour {
		return nil, fmt.Errorf("recharge not yet due: last recharge was %v", p.LastRecharge)
	}

	rechargeAmount := int64(20) // LAW-ECON-003
	recharged := make(map[string]int64)

	// Recharge each role, capped at initial allocation
	for roleLCT, initialAllocation := range p.RoleAllocations {
		currentBalance := p.RoleBalances[roleLCT]

		// Only recharge if below initial allocation
		if currentBalance < initialAllocation {
			// Calculate recharge amount (capped)
			maxRecharge := initialAllocation - currentBalance
			actualRecharge := rechargeAmount
			if actualRecharge > maxRecharge {
				actualRecharge = maxRecharge
			}

			// Apply recharge
			p.RoleBalances[roleLCT] += actualRecharge
			recharged[roleLCT] = actualRecharge
		}
	}

	p.LastRecharge = now
	p.UpdatedAt = now
	p.Version++

	return recharged, nil
}

// RechargeADP recharges ADP back to ATP (value creation)
func (p *SocietyTokenPool) RechargeADP(roleLCT string, amount int64, operationID string, reason string) (*AtpTransaction, error) {
	// Check if role exists
	adpBalance, exists := p.ADPBalances[roleLCT]
	if !exists {
		return nil, fmt.Errorf("role not found: %s", roleLCT)
	}

	// Check if role has sufficient ADP
	if adpBalance < amount {
		return nil, fmt.Errorf("insufficient ADP balance for role %s: %d < %d", roleLCT, adpBalance, amount)
	}

	// Check initial allocation cap
	initialAllocation := p.RoleAllocations[roleLCT]
	currentATP := p.RoleBalances[roleLCT]
	if currentATP+amount > initialAllocation {
		return nil, fmt.Errorf("recharge would exceed initial allocation: %d + %d > %d", currentATP, amount, initialAllocation)
	}

	// Recharge ADP to ATP
	p.ADPBalances[roleLCT] -= amount
	p.RoleBalances[roleLCT] += amount
	p.TotalADP -= amount
	p.UpdatedAt = time.Now()
	p.Version++

	// Create transaction record
	tx := &AtpTransaction{
		TransactionID: fmt.Sprintf("tx-recharge-%s-%d", roleLCT, time.Now().Unix()),
		SocietyID:     p.SocietyID,
		FromRole:      roleLCT,
		Amount:        amount,
		Type:          "recharge",
		OperationID:   operationID,
		Reason:        reason,
		Timestamp:     time.Now(),
		ResultingATP:  p.RoleBalances[roleLCT],
		ResultingADP:  p.ADPBalances[roleLCT],
	}

	return tx, nil
}

// GetRoleBalance returns current ATP and ADP balance for a role
func (p *SocietyTokenPool) GetRoleBalance(roleLCT string) (atp int64, adp int64, err error) {
	atpBalance, exists := p.RoleBalances[roleLCT]
	if !exists {
		return 0, 0, fmt.Errorf("role not found: %s", roleLCT)
	}

	adpBalance := p.ADPBalances[roleLCT]
	return atpBalance, adpBalance, nil
}

// ValidatePoolIntegrity checks pool integrity per LAW-ECON-001
func (p *SocietyTokenPool) ValidatePoolIntegrity() error {
	// Check total ATP budget
	if p.TotalATP != 1000 {
		return fmt.Errorf("total ATP must be 1000 per LAW-ECON-001, got %d", p.TotalATP)
	}

	// Sum all allocated ATP
	totalAllocated := int64(0)
	for _, amount := range p.RoleAllocations {
		totalAllocated += amount
	}

	if totalAllocated != p.AllocatedATP {
		return fmt.Errorf("allocated ATP mismatch: computed %d != stored %d", totalAllocated, p.AllocatedATP)
	}

	// Check conservation: AllocatedATP + AvailableATP = TotalATP
	if p.AllocatedATP+p.AvailableATP != p.TotalATP {
		return fmt.Errorf("ATP conservation violated: %d + %d != %d", p.AllocatedATP, p.AvailableATP, p.TotalATP)
	}

	// Sum all role balances
	totalRoleATP := int64(0)
	for _, balance := range p.RoleBalances {
		totalRoleATP += balance
	}

	// Sum all ADP balances
	totalRoleADP := int64(0)
	for _, balance := range p.ADPBalances {
		totalRoleADP += balance
	}

	// Energy conservation: RoleATP + RoleADP = AllocatedATP
	totalEnergy := totalRoleATP + totalRoleADP
	if totalEnergy != p.AllocatedATP {
		return fmt.Errorf("energy conservation violated: ATP(%d) + ADP(%d) = %d != allocated(%d)",
			totalRoleATP, totalRoleADP, totalEnergy, p.AllocatedATP)
	}

	return nil
}

// GetPoolSummary returns a summary of the pool state
func (p *SocietyTokenPool) GetPoolSummary() string {
	return fmt.Sprintf(
		"Society %s Pool: Total=%d, Allocated=%d, Available=%d, TotalADP=%d, Roles=%d, Version=%d",
		p.SocietyID, p.TotalATP, p.AllocatedATP, p.AvailableATP, p.TotalADP, len(p.RoleBalances), p.Version,
	)
}
