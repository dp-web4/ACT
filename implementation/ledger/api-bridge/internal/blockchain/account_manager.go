package blockchain

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

// AccountManager handles blockchain account operations
type AccountManager struct {
	logger   zerolog.Logger
	accounts map[string]*Account
	mu       sync.RWMutex
}

// Account represents a blockchain account
type Account struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	KeyType string `json:"key_type"`
}

// NewAccountManager creates a new account manager
func NewAccountManager(logger zerolog.Logger) *AccountManager {
	am := &AccountManager{
		logger:   logger,
		accounts: make(map[string]*Account),
	}

	// Initialize with common Ignite CLI accounts
	am.initializeDefaultAccounts()

	// Attempt to load real Ignite CLI accounts
	am.LoadIgniteAccounts()

	return am
}

// initializeDefaultAccounts sets up common Ignite CLI accounts
func (am *AccountManager) initializeDefaultAccounts() {
	defaultAccounts := []*Account{
		{
			Name:    "alice",
			Address: "cosmos1cs2clgcszut5ppvecfa4zrftvv9xz59w9fqcuv",
			KeyType: "secp256k1",
		},
		{
			Name:    "bob",
			Address: "cosmos18wz6nc4mgxdn5k2vce9y2nlxes3luzwz5tcurl",
			KeyType: "secp256k1",
		},
		{
			Name:    "charlie",
			Address: "cosmos1cs2clgcszut5ppvecfa4zrftvv9xz59w9fqcuv", // Using alice's address as fallback
			KeyType: "secp256k1",
		},
	}

	for _, account := range defaultAccounts {
		am.accounts[account.Name] = account
	}

	am.logger.Info().Int("count", len(defaultAccounts)).Msg("Initialized default accounts")
}

// GetOrCreateAccount gets an existing account or creates a mock one
func (am *AccountManager) GetOrCreateAccount(ctx context.Context, name string) (*Account, error) {
	am.mu.RLock()
	if account, exists := am.accounts[name]; exists {
		am.mu.RUnlock()
		return account, nil
	}
	am.mu.RUnlock()

	// Create new mock account
	am.mu.Lock()
	defer am.mu.Unlock()

	// Double-check after acquiring write lock
	if account, exists := am.accounts[name]; exists {
		return account, nil
	}

	account := am.createMockAccount(name)
	am.accounts[name] = account
	am.logger.Info().Str("name", name).Str("address", account.Address).Msg("Created new fallback account")
	return account, nil
}

// createMockAccount creates a fallback account when real accounts aren't available
func (am *AccountManager) createMockAccount(name string) *Account {
	// Generate a mock address based on the name
	address := fmt.Sprintf("cosmos1%s%s", strings.Repeat("0", 38-len(name)), name)

	return &Account{
		Name:    name,
		Address: address,
		KeyType: "secp256k1",
	}
}

// GetDefaultAccount returns a default account for transactions
func (am *AccountManager) GetDefaultAccount() *Account {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Return alice as the default account
	if account, exists := am.accounts["alice"]; exists {
		return account
	}

	// Fallback to any available account
	for _, account := range am.accounts {
		return account
	}

	// Create a fallback account if none exist
	return &Account{
		Name:    "default",
		Address: "cosmos1default000000000000000000000000000000000",
		KeyType: "secp256k1",
	}
}

// GetAccountForCreator returns the best account for a given creator
func (am *AccountManager) GetAccountForCreator(creator string) *Account {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// If creator matches an existing account name, use it
	if account, exists := am.accounts[creator]; exists {
		return account
	}

	// Otherwise, use a default account
	return am.GetDefaultAccount()
}

// ListAccounts returns all managed accounts
func (am *AccountManager) ListAccounts() []*Account {
	am.mu.RLock()
	defer am.mu.RUnlock()

	accounts := make([]*Account, 0, len(am.accounts))
	for _, account := range am.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

// GetAccount gets an account by name
func (am *AccountManager) GetAccount(name string) (*Account, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	account, exists := am.accounts[name]
	return account, exists
}

// LoadIgniteAccounts attempts to load real Ignite CLI accounts
func (am *AccountManager) LoadIgniteAccounts() {
	am.logger.Info().Msg("Attempting to load real Ignite CLI accounts")

	// For now, we'll use the default accounts
	// In the future, this could query Ignite CLI for real account information
	am.logger.Info().Msg("Using default accounts (alice, bob, charlie)")

	// You could extend this to:
	// 1. Run `ignite keys list` to get real account names
	// 2. Run `ignite keys show <name> --output json` to get real addresses
	// 3. Update the account information with real data
}
