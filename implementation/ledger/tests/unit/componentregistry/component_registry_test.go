package componentregistry_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"racecar-web/x/componentregistry/keeper"
	"racecar-web/x/componentregistry/types"
)

// ComponentRegistryTestSuite provides a test suite for Component Registry functionality
type ComponentRegistryTestSuite struct {
	suite.Suite
	ctx       context.Context
	keeper    keeper.Keeper
	mockStore *MockStoreService
	mockAuth  *MockAuthKeeper
	mockBank  *MockBankKeeper
}

// MockStoreService mocks the store service for testing
type MockStoreService struct {
	store map[string][]byte
}

func (m *MockStoreService) OpenKVStore(name string) (types.KVStore, error) {
	return &MockKVStore{store: m.store}, nil
}

// MockKVStore mocks the key-value store for testing
type MockKVStore struct {
	store map[string][]byte
}

func (m *MockKVStore) Get(key []byte) ([]byte, error) {
	if val, exists := m.store[string(key)]; exists {
		return val, nil
	}
	return nil, nil
}

func (m *MockKVStore) Set(key, value []byte) error {
	m.store[string(key)] = value
	return nil
}

func (m *MockKVStore) Delete(key []byte) error {
	delete(m.store, string(key))
	return nil
}

func (m *MockKVStore) Iterator(start, end []byte) (types.Iterator, error) {
	return &MockIterator{}, nil
}

func (m *MockKVStore) ReverseIterator(start, end []byte) (types.Iterator, error) {
	return &MockIterator{}, nil
}

// MockIterator mocks the iterator for testing
type MockIterator struct {
	keys   []string
	values [][]byte
	index  int
}

func (m *MockIterator) Domain() ([]byte, []byte) {
	return nil, nil
}

func (m *MockIterator) Valid() bool {
	return m.index < len(m.keys)
}

func (m *MockIterator) Next() {
	m.index++
}

func (m *MockIterator) Key() []byte {
	if m.Valid() {
		return []byte(m.keys[m.index])
	}
	return nil
}

func (m *MockIterator) Value() []byte {
	if m.Valid() {
		return m.values[m.index]
	}
	return nil
}

func (m *MockIterator) Close() error {
	return nil
}

// MockAuthKeeper mocks the auth keeper for testing
type MockAuthKeeper struct{}

func (m *MockAuthKeeper) AddressCodec() types.AddressCodec {
	return &MockAddressCodec{}
}

// MockBankKeeper mocks the bank keeper for testing
type MockBankKeeper struct{}

func (m *MockBankKeeper) SpendableCoins(ctx context.Context, addr types.AccAddress) types.Coins {
	return types.NewCoins()
}

func (m *MockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr types.AccAddress, amt types.Coins) error {
	return nil
}

// MockAddressCodec mocks the address codec for testing
type MockAddressCodec struct{}

func (m *MockAddressCodec) StringToBytes(text string) ([]byte, error) {
	return []byte(text), nil
}

func (m *MockAddressCodec) BytesToString(bz []byte) (string, error) {
	return string(bz), nil
}

// SetupTest initializes the test suite
func (suite *ComponentRegistryTestSuite) SetupTest() {
	suite.mockStore = &MockStoreService{
		store: make(map[string][]byte),
	}
	suite.mockAuth = &MockAuthKeeper{}
	suite.mockBank = &MockBankKeeper{}

	suite.ctx = context.Background()
	suite.keeper = keeper.NewKeeper(
		nil, // codec
		suite.mockStore,
		suite.mockAuth,
		suite.mockBank,
	)
}

// TestRegisterComponent_ValidInput tests component registration with valid input
func (suite *ComponentRegistryTestSuite) TestRegisterComponent_ValidInput() {
	// Test data
	creator := "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k"
	componentID := "comp_battery_module_001"
	componentType := "battery_module"
	manufacturerData := `{
		"manufacturer_id": "MODBATT_CORP",
		"model": "BM-2000",
		"serial_number": "SN123456789",
		"capacity_wh": 2000,
		"voltage_nominal": 48,
		"chemistry": "lithium_ion"
	}`

	// Register component
	err := suite.keeper.RegisterComponent(
		suite.ctx,
		creator,
		componentID,
		componentType,
		manufacturerData,
	)

	// Assertions
	require.NoError(suite.T(), err)

	// Verify component was stored
	component, found := suite.keeper.GetComponentIdentity(suite.ctx, componentID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), componentID, component.ComponentId)
	assert.Equal(suite.T(), componentType, component.ComponentType)
	assert.Equal(suite.T(), "MODBATT_CORP", component.ManufacturerId)
	assert.Equal(suite.T(), "pending", component.VerificationStatus)
}

// TestRegisterComponent_ExtractManufacturerID tests manufacturer ID extraction from JSON
func (suite *ComponentRegistryTestSuite) TestRegisterComponent_ExtractManufacturerID() {
	// Test data with different manufacturer ID formats
	testCases := []struct {
		name             string
		manufacturerData string
		expectedID       string
	}{
		{
			name: "Standard manufacturer ID",
			manufacturerData: `{
				"manufacturer_id": "MODBATT_CORP",
				"model": "BM-2000"
			}`,
			expectedID: "MODBATT_CORP",
		},
		{
			name: "Lowercase manufacturer ID",
			manufacturerData: `{
				"manufacturer_id": "modbatt_corp",
				"model": "BM-2000"
			}`,
			expectedID: "modbatt_corp",
		},
		{
			name: "Manufacturer ID with spaces",
			manufacturerData: `{
				"manufacturer_id": "MODBATT CORP",
				"model": "BM-2000"
			}`,
			expectedID: "MODBATT CORP",
		},
		{
			name: "Missing manufacturer ID",
			manufacturerData: `{
				"model": "BM-2000",
				"serial_number": "SN123456789"
			}`,
			expectedID: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			componentID := "comp_test_" + tc.name

			err := suite.keeper.RegisterComponent(
				suite.ctx,
				"cosmos1creator",
				componentID,
				"battery_module",
				tc.manufacturerData,
			)

			if tc.expectedID == "" {
				// Should fail if no manufacturer ID
				assert.Error(t, err)
			} else {
				// Should succeed and extract manufacturer ID
				require.NoError(t, err)
				component, found := suite.keeper.GetComponentIdentity(suite.ctx, componentID)
				require.True(t, found)
				assert.Equal(t, tc.expectedID, component.ManufacturerId)
			}
		})
	}
}

// TestRegisterComponent_InvalidInput tests component registration with invalid input
func (suite *ComponentRegistryTestSuite) TestRegisterComponent_InvalidInput() {
	testCases := []struct {
		name             string
		creator          string
		componentID      string
		componentType    string
		manufacturerData string
		expectError      bool
	}{
		{
			name:             "Empty creator",
			creator:          "",
			componentID:      "comp_test",
			componentType:    "battery_module",
			manufacturerData: `{"manufacturer_id": "TEST"}`,
			expectError:      true,
		},
		{
			name:             "Empty component ID",
			creator:          "cosmos1creator",
			componentID:      "",
			componentType:    "battery_module",
			manufacturerData: `{"manufacturer_id": "TEST"}`,
			expectError:      true,
		},
		{
			name:             "Empty component type",
			creator:          "cosmos1creator",
			componentID:      "comp_test",
			componentType:    "",
			manufacturerData: `{"manufacturer_id": "TEST"}`,
			expectError:      true,
		},
		{
			name:             "Invalid JSON",
			creator:          "cosmos1creator",
			componentID:      "comp_test",
			componentType:    "battery_module",
			manufacturerData: `{invalid json}`,
			expectError:      true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.keeper.RegisterComponent(
				suite.ctx,
				tc.creator,
				tc.componentID,
				tc.componentType,
				tc.manufacturerData,
			)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestVerifyComponent tests component verification
func (suite *ComponentRegistryTestSuite) TestVerifyComponent() {
	// Register a component first
	componentID := "comp_verify_test"
	err := suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentID,
		"battery_module",
		`{"manufacturer_id": "TEST_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Verify the component
	verifier := "cosmos1verifier"
	verificationData := `{
		"verification_method": "digital_signature",
		"signature": "abc123",
		"timestamp": "2024-01-01T00:00:00Z"
	}`

	err = suite.keeper.VerifyComponent(
		suite.ctx,
		verifier,
		componentID,
		"verified",
		verificationData,
	)

	// Assertions
	require.NoError(suite.T(), err)

	// Check component status was updated
	component, found := suite.keeper.GetComponentIdentity(suite.ctx, componentID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), "verified", component.VerificationStatus)
	assert.Equal(suite.T(), verifier, component.VerifiedBy)
	assert.NotEmpty(suite.T(), component.VerifiedAt)
}

// TestUpdateAuthorization tests authorization rule updates
func (suite *ComponentRegistryTestSuite) TestUpdateAuthorization() {
	// Register a component first
	componentID := "comp_auth_test"
	err := suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentID,
		"battery_module",
		`{"manufacturer_id": "TEST_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Update authorization rules
	authorizer := "cosmos1authorizer"
	authorizationRules := `{
		"allowed_pairing_types": ["battery_pack", "battery_module"],
		"max_energy_transfer_wh": 1000,
		"requires_trust_score": 0.8,
		"allowed_operations": ["read", "write", "energy_transfer"]
	}`

	err = suite.keeper.UpdateAuthorization(
		suite.ctx,
		authorizer,
		componentID,
		authorizationRules,
	)

	// Assertions
	require.NoError(suite.T(), err)

	// Check authorization was updated
	component, found := suite.keeper.GetComponentIdentity(suite.ctx, componentID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), authorizationRules, component.AuthorizationRules)
	assert.Equal(suite.T(), authorizer, component.AuthorizedBy)
	assert.NotEmpty(suite.T(), component.AuthorizedAt)
}

// TestCheckPairingAuthorization tests bidirectional pairing authorization
func (suite *ComponentRegistryTestSuite) TestCheckPairingAuthorization() {
	// Register two components
	componentA := "comp_a"
	componentB := "comp_b"

	// Register component A
	err := suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentA,
		"battery_pack",
		`{"manufacturer_id": "PACK_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Register component B
	err = suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentB,
		"battery_module",
		`{"manufacturer_id": "MODULE_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Set authorization rules for component A
	err = suite.keeper.UpdateAuthorization(
		suite.ctx,
		"cosmos1authorizer",
		componentA,
		`{"allowed_pairing_types": ["battery_module"]}`,
	)
	require.NoError(suite.T(), err)

	// Set authorization rules for component B
	err = suite.keeper.UpdateAuthorization(
		suite.ctx,
		"cosmos1authorizer",
		componentB,
		`{"allowed_pairing_types": ["battery_pack"]}`,
	)
	require.NoError(suite.T(), err)

	// Check bidirectional authorization
	aCanPairB, bCanPairA, err := suite.keeper.CheckPairingAuthorization(suite.ctx, componentA, componentB)

	// Assertions
	require.NoError(suite.T(), err)
	assert.True(suite.T(), aCanPairB, "Component A should be able to pair with Component B")
	assert.True(suite.T(), bCanPairA, "Component B should be able to pair with Component A")
}

// TestCheckPairingAuthorization_Unauthorized tests unauthorized pairing scenarios
func (suite *ComponentRegistryTestSuite) TestCheckPairingAuthorization_Unauthorized() {
	// Register two components with incompatible authorization rules
	componentA := "comp_unauthorized_a"
	componentB := "comp_unauthorized_b"

	// Register component A
	err := suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentA,
		"battery_pack",
		`{"manufacturer_id": "PACK_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Register component B
	err = suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentB,
		"battery_module",
		`{"manufacturer_id": "MODULE_CORP"}`,
	)
	require.NoError(suite.T(), err)

	// Set incompatible authorization rules
	err = suite.keeper.UpdateAuthorization(
		suite.ctx,
		"cosmos1authorizer",
		componentA,
		`{"allowed_pairing_types": ["controller"]}`, // Only allows controllers
	)
	require.NoError(suite.T(), err)

	err = suite.keeper.UpdateAuthorization(
		suite.ctx,
		"cosmos1authorizer",
		componentB,
		`{"allowed_pairing_types": ["battery_pack"]}`, // Allows battery packs
	)
	require.NoError(suite.T(), err)

	// Check bidirectional authorization
	aCanPairB, bCanPairA, err := suite.keeper.CheckPairingAuthorization(suite.ctx, componentA, componentB)

	// Assertions
	require.NoError(suite.T(), err)
	assert.False(suite.T(), aCanPairB, "Component A should not be able to pair with Component B")
	assert.True(suite.T(), bCanPairA, "Component B should be able to pair with Component A")
}

// TestGetComponentIdentity_NotFound tests getting non-existent component
func (suite *ComponentRegistryTestSuite) TestGetComponentIdentity_NotFound() {
	component, found := suite.keeper.GetComponentIdentity(suite.ctx, "non_existent_component")
	assert.False(suite.T(), found)
	assert.Empty(suite.T(), component.ComponentId)
}

// TestListAuthorizedPartners tests listing authorized partners
func (suite *ComponentRegistryTestSuite) TestListAuthorizedPartners() {
	// Register multiple components
	components := []string{"comp_1", "comp_2", "comp_3", "comp_4"}

	for i, compID := range components {
		err := suite.keeper.RegisterComponent(
			suite.ctx,
			"cosmos1creator",
			compID,
			"battery_module",
			`{"manufacturer_id": "TEST_CORP"}`,
		)
		require.NoError(suite.T(), err)

		// Set authorization rules
		authRules := `{"allowed_pairing_types": ["battery_module", "battery_pack"]}`
		err = suite.keeper.UpdateAuthorization(
			suite.ctx,
			"cosmos1authorizer",
			compID,
			authRules,
		)
		require.NoError(suite.T(), err)
	}

	// List authorized partners for component 1
	partners, err := suite.keeper.ListAuthorizedPartners(suite.ctx, components[0])

	// Assertions
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), partners, len(components)-1) // Should include all other components
}

// Benchmark tests for performance
func (suite *ComponentRegistryTestSuite) BenchmarkRegisterComponent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		componentID := fmt.Sprintf("comp_bench_%d", i)
		err := suite.keeper.RegisterComponent(
			suite.ctx,
			"cosmos1creator",
			componentID,
			"battery_module",
			`{"manufacturer_id": "BENCH_CORP"}`,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *ComponentRegistryTestSuite) BenchmarkCheckPairingAuthorization(b *testing.B) {
	// Setup: Register two components
	componentA := "comp_bench_a"
	componentB := "comp_bench_b"

	err := suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentA,
		"battery_pack",
		`{"manufacturer_id": "PACK_CORP"}`,
	)
	if err != nil {
		b.Fatal(err)
	}

	err = suite.keeper.RegisterComponent(
		suite.ctx,
		"cosmos1creator",
		componentB,
		"battery_module",
		`{"manufacturer_id": "MODULE_CORP"}`,
	)
	if err != nil {
		b.Fatal(err)
	}

	// Benchmark the authorization check
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := suite.keeper.CheckPairingAuthorization(suite.ctx, componentA, componentB)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestComponentRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentRegistryTestSuite))
}
