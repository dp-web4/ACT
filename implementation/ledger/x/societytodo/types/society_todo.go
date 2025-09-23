package types

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SocietyTodoList represents a society's todo management system
type SocietyTodoList struct {
	SocietyLCT      string       `json:"society_lct"`
	State           SocietyState `json:"state"`
	ATPBudget       sdk.Int      `json:"atp_budget"`
	ATPAvailable    sdk.Int      `json:"atp_available"`
	CycleStartTime  time.Time    `json:"cycle_start_time"`
	CycleDuration   time.Duration `json:"cycle_duration"`
	CycleNumber     uint64       `json:"cycle_number"`
	EnergyEfficiency float64     `json:"energy_efficiency"`
	AcceptingRequests bool       `json:"accepting_requests"`
	EmergencyReserve sdk.Int     `json:"emergency_reserve"`
	ActiveTodos     []string     `json:"active_todos"`
	CompletedCount  uint64       `json:"completed_count"`
	FailedCount     uint64       `json:"failed_count"`
}

// TodoItem represents a single task in the society
type TodoItem struct {
	ID          string       `json:"id"`
	SocietyLCT  string       `json:"society_lct"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TodoPriority `json:"priority"`
	Status      TodoStatus   `json:"status"`
	Complexity  uint32       `json:"complexity"`
	ATPCost     sdk.Int      `json:"atp_cost"`
	ADPReward   sdk.Int      `json:"adp_reward"`
	RequestedBy string       `json:"requested_by"`
	ExecutedBy  string       `json:"executed_by"`
	CreatedAt   time.Time    `json:"created_at"`
	StartedAt   *time.Time   `json:"started_at,omitempty"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	ProofHash   string       `json:"proof_hash,omitempty"`
	QualityScore float64     `json:"quality_score"`
	SharedWith  []string     `json:"shared_with,omitempty"`
}

// CitizenRequest represents a request from a society citizen
type CitizenRequest struct {
	ID          string       `json:"id"`
	CitizenLCT  string       `json:"citizen_lct"`
	SocietyLCT  string       `json:"society_lct"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TodoPriority `json:"priority"`
	Complexity  uint32       `json:"complexity"`
	EstimatedATP sdk.Int     `json:"estimated_atp"`
	TrustLevel  sdk.Dec      `json:"trust_level"`
	RequestTime time.Time    `json:"request_time"`
	Status      string       `json:"status"`
	TodoID      string       `json:"todo_id,omitempty"`
	RejectionReason string   `json:"rejection_reason,omitempty"`
}

// ATPDelegation represents a citizen's ATP delegation to a pool
type ATPDelegation struct {
	ID           string         `json:"id"`
	CitizenLCT   string         `json:"citizen_lct"`
	SocietyLCT   string         `json:"society_lct"`
	Pool         DelegationPool `json:"pool"`
	Amount       sdk.Int        `json:"amount"`
	VotingPower  sdk.Dec        `json:"voting_power"`
	LockPeriod   time.Duration  `json:"lock_period"`
	DelegatedAt  time.Time      `json:"delegated_at"`
	UnlocksAt    time.Time      `json:"unlocks_at"`
	IsActive     bool           `json:"is_active"`
}

// WakeSleepCycle represents a society's energy cycle history
type WakeSleepCycle struct {
	SocietyLCT   string       `json:"society_lct"`
	CycleNumber  uint64       `json:"cycle_number"`
	StartState   SocietyState `json:"start_state"`
	EndState     SocietyState `json:"end_state"`
	StartATP     sdk.Int      `json:"start_atp"`
	EndATP       sdk.Int      `json:"end_atp"`
	TodosCreated uint32       `json:"todos_created"`
	TodosCompleted uint32     `json:"todos_completed"`
	TodosFailed  uint32       `json:"todos_failed"`
	TotalATPSpent sdk.Int     `json:"total_atp_spent"`
	TotalADPEarned sdk.Int    `json:"total_adp_earned"`
	Efficiency   float64      `json:"efficiency"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      time.Time    `json:"end_time"`
}

// FederationShare represents cross-society todo sharing
type FederationShare struct {
	TodoID         string    `json:"todo_id"`
	OriginSociety  string    `json:"origin_society"`
	SharedSocieties []string `json:"shared_societies"`
	ShareModel     string    `json:"share_model"`
	CostSplit      map[string]sdk.Dec `json:"cost_split"`
	RewardSplit    map[string]sdk.Dec `json:"reward_split"`
	TrustThreshold sdk.Dec   `json:"trust_threshold"`
	SharedAt       time.Time `json:"shared_at"`
	AcceptedBy     []string  `json:"accepted_by"`
}

// DelegationPoolStats tracks statistics for each delegation pool
type DelegationPoolStats struct {
	Pool            DelegationPool `json:"pool"`
	TotalDelegated  sdk.Int        `json:"total_delegated"`
	ActiveDelegators uint32        `json:"active_delegators"`
	TotalVotingPower sdk.Dec       `json:"total_voting_power"`
	TodosSupported  uint32         `json:"todos_supported"`
	ATPSpent        sdk.Int        `json:"atp_spent"`
	ADPEarned       sdk.Int        `json:"adp_earned"`
	LastUpdated     time.Time      `json:"last_updated"`
}