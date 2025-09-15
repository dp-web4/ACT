package keeper

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	energycyclekeeper "racecar-web/x/energycycle/keeper"
	energycyclemodule "racecar-web/x/energycycle/module"
	energycycletypes "racecar-web/x/energycycle/types"
	lctmanagertypes "racecar-web/x/lctmanager/types"
	trusttensorkeeper "racecar-web/x/trusttensor/keeper"
	trusttensormodule "racecar-web/x/trusttensor/module"
	trusttensortypes "racecar-web/x/trusttensor/types"
)

// TrustTensorKeeper returns a test fixture specifically for trust tensor tests
func TrustTensorKeeper(t *testing.T) (context.Context, trusttensorkeeper.Keeper) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(trusttensormodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(trusttensortypes.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(trusttensortypes.GovModuleName)

	// Create mock keepers for dependencies
	mockLctManagerKeeper := &MockLctManagerKeeper{}

	k := trusttensorkeeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, // bank keeper
		mockLctManagerKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, trusttensortypes.DefaultParams()); err != nil {
		t.Fatalf("failed to set trust tensor params: %v", err)
	}

	return ctx, k
}

// EnergyCycleKeeper returns a test fixture specifically for energy cycle tests
func EnergyCycleKeeper(t *testing.T) (context.Context, energycyclekeeper.Keeper) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(energycyclemodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(energycycletypes.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(energycycletypes.GovModuleName)

	// Create mock keepers for dependencies
	mockLctManagerKeeper := &MockLctManagerKeeper{}
	mockTrustTensorKeeper := &MockTrustTensorKeeper{}

	k := energycyclekeeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, // bank keeper
		mockLctManagerKeeper,
		mockTrustTensorKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, energycycletypes.DefaultParams()); err != nil {
		t.Fatalf("failed to set energy cycle params: %v", err)
	}

	return ctx, k
}

// MockLctManagerKeeper provides a mock implementation for LCT manager keeper
type MockLctManagerKeeper struct{}

func (m *MockLctManagerKeeper) GetLinkedContextToken(ctx context.Context, lctId string) (lctmanagertypes.LinkedContextToken, bool) {
	// Return a mock LCT for testing
	return lctmanagertypes.LinkedContextToken{
		LctId:         lctId,
		PairingStatus: "active",
	}, true
}

func (m *MockLctManagerKeeper) CreateLCTRelationship(ctx context.Context, componentA, componentB, operationalContext, proxyId string) (string, string, error) {
	// Return mock LCT ID and key reference for testing
	lctId := fmt.Sprintf("lct_%s_%s", componentA, componentB)
	keyReference := fmt.Sprintf("key_ref_%s_%s", componentA, componentB)
	return lctId, keyReference, nil
}

func (m *MockLctManagerKeeper) GetComponentRelationships(ctx context.Context, componentId string) ([]lctmanagertypes.LinkedContextToken, error) {
	// Return empty list for testing
	return []lctmanagertypes.LinkedContextToken{}, nil
}

func (m *MockLctManagerKeeper) TerminateLCTRelationship(ctx context.Context, lctId, reason string, notifyOffline bool) error {
	// Mock successful termination
	return nil
}

// MockTrustTensorKeeper provides a mock implementation for trust tensor keeper
type MockTrustTensorKeeper struct{}

func (m *MockTrustTensorKeeper) CalculateRelationshipTrust(ctx context.Context, lctId, operationalContext string) (string, string, error) {
	// Return a default trust score for testing
	return "0.5", "default_trust", nil
}

func (m *MockTrustTensorKeeper) CalculateV3CompositeScore(ctx context.Context, operationID string) (math.LegacyDec, error) {
	// Return a default V3 score for testing
	return math.LegacyNewDecWithPrec(7, 1), nil // 0.7
}
