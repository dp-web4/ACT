package lctmanager_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"racecar-web/x/lctmanager/keeper"
	"racecar-web/x/lctmanager/types"
)

// LCTManagerTestSuite provides a test suite for LCT Manager functionality
type LCTManagerTestSuite struct {
	suite.Suite
	ctx         context.Context
	keeper      keeper.Keeper
	mockStore   *MockStoreService
	mockAuth    *MockAuthKeeper
	mockBank    *MockBankKeeper
	mockCompReg *MockComponentRegistryKeeper
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

// MockComponentRegistryKeeper mocks the component registry keeper for testing
type MockComponentRegistryKeeper struct{}

func (m *MockComponentRegistryKeeper) GetComponentIdentity(ctx context.Context, componentID string) (types.ComponentIdentity, bool) {
	// Return mock component data
	return types.ComponentIdentity{
		ComponentId:        componentID,
		ComponentType:      "battery_module",
		ManufacturerId:     "TEST_CORP",
		VerificationStatus: "verified",
		AuthorizationRules: `{"allowed_pairing_types": ["battery_pack", "battery_module"]}`,
	}, true
}

func (m *MockComponentRegistryKeeper) CheckPairingAuthorization(ctx context.Context, componentA, componentB string) (bool, bool, error) {
	// Mock bidirectional authorization
	return true, true, nil
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
func (suite *LCTManagerTestSuite) SetupTest() {
	suite.mockStore = &MockStoreService{
		store: make(map[string][]byte),
	}
	suite.mockAuth = &MockAuthKeeper{}
	suite.mockBank = &MockBankKeeper{}
	suite.mockCompReg = &MockComponentRegistryKeeper{}

	suite.ctx = context.Background()
	suite.keeper = keeper.NewKeeper(
		nil, // codec
		suite.mockStore,
		suite.mockAuth,
		suite.mockBank,
		suite.mockCompReg,
	)
}

// TestCreateLCTRelationship_ValidInput tests LCT relationship creation with valid input
func (suite *LCTManagerTestSuite) TestCreateLCTRelationship_ValidInput() {
	// Test data
	creator := "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k"
	componentA := "comp_battery_pack_001"
	componentB := "comp_battery_module_001"
	context := "race_car_battery_management"
	proxyID := "comp_gateway_001"

	// Create LCT relationship
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		creator,
		componentA,
		componentB,
		context,
		proxyID,
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), lct)
	assert.NotEmpty(suite.T(), lct.LctId)
	assert.Equal(suite.T(), componentA, lct.ComponentAId)
	assert.Equal(suite.T(), componentB, lct.ComponentBId)
	assert.Equal(suite.T(), context, lct.OperationalContext)
	assert.Equal(suite.T(), proxyID, lct.ProxyComponentId)
	assert.Equal(suite.T(), "pending", lct.PairingStatus)
	assert.Empty(suite.T(), lct.LctKeyHalf, "No key half should be stored on-chain")

	// Verify LCT was stored
	storedLct, found := suite.keeper.GetLinkedContextToken(suite.ctx, lct.LctId)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), lct.LctId, storedLct.LctId)
}

// TestCreateLCTRelationship_NoKeyStorage tests that no cryptographic keys are stored on-chain
func (suite *LCTManagerTestSuite) TestCreateLCTRelationship_NoKeyStorage() {
	// Create LCT relationship
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Empty(suite.T(), lct.LctKeyHalf, "LCT key half should not be stored on-chain")
	assert.Empty(suite.T(), lct.DeviceKeyHalf, "Device key half should not be stored on-chain")

	// Verify no key data in storage
	storedLct, found := suite.keeper.GetLinkedContextToken(suite.ctx, lct.LctId)
	require.True(suite.T(), found)
	assert.Empty(suite.T(), storedLct.LctKeyHalf, "No key half should be stored in blockchain state")
}

// TestInitiateLCTMediatedPairing tests LCT-mediated pairing initiation
func (suite *LCTManagerTestSuite) TestInitiateLCTMediatedPairing() {
	// Create two LCTs first
	lct1, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	lct2, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_c",
		"comp_d",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Initiate LCT-mediated pairing
	creator := "cosmos1creator"
	context := "pairing_context"
	proxyLctId := "lct_proxy"
	expiresAt := time.Now().Add(time.Hour).Unix()

	response, err := suite.keeper.InitiateLCTMediatedPairing(
		suite.ctx,
		creator,
		lct1.LctId,
		lct2.LctId,
		context,
		proxyLctId,
		expiresAt,
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.NotEmpty(suite.T(), response.PairingId)
	assert.Equal(suite.T(), "initiated", response.Status)
	assert.NotEmpty(suite.T(), response.ChallengeId)
	assert.NotEmpty(suite.T(), response.ChallengeData)
	assert.NotEmpty(suite.T(), response.LctRelationshipId)
}

// TestCompleteLCTMediatedPairing tests LCT-mediated pairing completion
func (suite *LCTManagerTestSuite) TestCompleteLCTMediatedPairing() {
	// Create LCTs and initiate pairing
	lct1, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	lct2, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_c",
		"comp_d",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Initiate pairing
	initResponse, err := suite.keeper.InitiateLCTMediatedPairing(
		suite.ctx,
		"cosmos1creator",
		lct1.LctId,
		lct2.LctId,
		"test_context",
		"",
		time.Now().Add(time.Hour).Unix(),
	)
	require.NoError(suite.T(), err)

	// Complete pairing
	creator := "cosmos1creator"
	initiatorResponse := "initiator_auth_response"
	targetResponse := "target_auth_response"
	sessionKeyData := []byte("session_key_data_123")

	response, err := suite.keeper.CompleteLCTMediatedPairing(
		suite.ctx,
		creator,
		initResponse.PairingId,
		initiatorResponse,
		targetResponse,
		sessionKeyData,
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), initResponse.PairingId, response.PairingId)
	assert.Equal(suite.T(), "completed", response.Status)
	assert.NotEmpty(suite.T(), response.LctRelationshipId)
	assert.NotEmpty(suite.T(), response.EncryptedSessionKeyInitiator)
	assert.NotEmpty(suite.T(), response.EncryptedSessionKeyTarget)
	assert.NotEmpty(suite.T(), response.HashedCombinedSessionKey)
	assert.NotEmpty(suite.T(), response.TrustScore)
}

// TestGenerateSplitKeyPair tests split key generation
func (suite *LCTManagerTestSuite) TestGenerateSplitKeyPair() {
	// Generate split key pair
	lctKeyHalf, deviceKeyHalf, err := suite.keeper.GenerateSplitKeyPair()

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), lctKeyHalf)
	assert.NotEmpty(suite.T(), deviceKeyHalf)
	assert.NotEqual(suite.T(), lctKeyHalf, deviceKeyHalf, "Key halves should be different")

	// Verify key format (hex string)
	assert.Len(suite.T(), lctKeyHalf, 64, "LCT key half should be 32 bytes = 64 hex chars")
	assert.Len(suite.T(), deviceKeyHalf, 64, "Device key half should be 32 bytes = 64 hex chars")
}

// TestGenerateSplitKeyPair_Uniqueness tests that generated keys are unique
func (suite *LCTManagerTestSuite) TestGenerateSplitKeyPair_Uniqueness() {
	// Generate multiple key pairs
	keyPairs := make(map[string]bool)

	for i := 0; i < 100; i++ {
		lctKeyHalf, deviceKeyHalf, err := suite.keeper.GenerateSplitKeyPair()
		require.NoError(suite.T(), err)

		// Check uniqueness
		combinedKey := lctKeyHalf + deviceKeyHalf
		assert.False(suite.T(), keyPairs[combinedKey], "Generated keys should be unique")
		keyPairs[combinedKey] = true
	}
}

// TestGetLinkedContextToken tests LCT retrieval
func (suite *LCTManagerTestSuite) TestGetLinkedContextToken() {
	// Create LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Retrieve LCT
	retrievedLct, found := suite.keeper.GetLinkedContextToken(suite.ctx, lct.LctId)

	// Assertions
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), lct.LctId, retrievedLct.LctId)
	assert.Equal(suite.T(), lct.ComponentAId, retrievedLct.ComponentAId)
	assert.Equal(suite.T(), lct.ComponentBId, retrievedLct.ComponentBId)
}

// TestGetLinkedContextToken_NotFound tests retrieval of non-existent LCT
func (suite *LCTManagerTestSuite) TestGetLinkedContextToken_NotFound() {
	lct, found := suite.keeper.GetLinkedContextToken(suite.ctx, "non_existent_lct")
	assert.False(suite.T(), found)
	assert.Empty(suite.T(), lct.LctId)
}

// TestUpdateLCTStatus tests LCT status updates
func (suite *LCTManagerTestSuite) TestUpdateLCTStatus() {
	// Create LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Update status
	newStatus := "active"
	err = suite.keeper.UpdateLCTStatus(suite.ctx, lct.LctId, newStatus)

	// Assertions
	require.NoError(suite.T(), err)

	// Verify status was updated
	updatedLct, found := suite.keeper.GetLinkedContextToken(suite.ctx, lct.LctId)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), newStatus, updatedLct.PairingStatus)
	assert.NotEqual(suite.T(), lct.UpdatedAt, updatedLct.UpdatedAt, "UpdatedAt should be changed")
}

// TestTerminateLCTRelationship tests LCT relationship termination
func (suite *LCTManagerTestSuite) TestTerminateLCTRelationship() {
	// Create LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Terminate relationship
	terminator := "cosmos1terminator"
	reason := "component_replacement"
	err = suite.keeper.TerminateLCTRelationship(suite.ctx, terminator, lct.LctId, reason)

	// Assertions
	require.NoError(suite.T(), err)

	// Verify termination
	terminatedLct, found := suite.keeper.GetLinkedContextToken(suite.ctx, lct.LctId)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), "terminated", terminatedLct.PairingStatus)
	assert.Equal(suite.T(), terminator, terminatedLct.TerminatedBy)
	assert.Equal(suite.T(), reason, terminatedLct.TerminationReason)
	assert.NotEmpty(suite.T(), terminatedLct.TerminatedAt)
}

// TestGetComponentRelationships tests component relationship tracking
func (suite *LCTManagerTestSuite) TestGetComponentRelationships() {
	// Create multiple LCTs for the same component
	componentA := "comp_a"

	lct1, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		componentA,
		"comp_b",
		"test_context_1",
		"",
	)
	require.NoError(suite.T(), err)

	lct2, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		componentA,
		"comp_c",
		"test_context_2",
		"",
	)
	require.NoError(suite.T(), err)

	// Get component relationships
	relationships, err := suite.keeper.GetComponentRelationships(suite.ctx, componentA)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), relationships, 2, "Should have 2 relationships for component A")

	// Verify both LCTs are included
	lctIds := make(map[string]bool)
	for _, rel := range relationships {
		lctIds[rel.RelatedLcts] = true
	}
	assert.True(suite.T(), lctIds[lct1.LctId])
	assert.True(suite.T(), lctIds[lct2.LctId])
}

// TestValidateLCTAccess tests LCT access validation
func (suite *LCTManagerTestSuite) TestValidateLCTAccess() {
	// Create LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	// Update status to active
	err = suite.keeper.UpdateLCTStatus(suite.ctx, lct.LctId, "active")
	require.NoError(suite.T(), err)

	// Validate access
	requester := "cosmos1requester"
	operation := "energy_transfer"
	valid, reason, err := suite.keeper.ValidateLCTAccess(suite.ctx, requester, lct.LctId, operation)

	// Assertions
	require.NoError(suite.T(), err)
	assert.True(suite.T(), valid, "Access should be valid for active LCT")
	assert.Empty(suite.T(), reason, "No reason should be provided for valid access")
}

// TestValidateLCTAccess_Terminated tests access validation for terminated LCT
func (suite *LCTManagerTestSuite) TestValidateLCTAccess_Terminated() {
	// Create and terminate LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	require.NoError(suite.T(), err)

	err = suite.keeper.TerminateLCTRelationship(suite.ctx, "cosmos1terminator", lct.LctId, "test")
	require.NoError(suite.T(), err)

	// Validate access
	requester := "cosmos1requester"
	operation := "energy_transfer"
	valid, reason, err := suite.keeper.ValidateLCTAccess(suite.ctx, requester, lct.LctId, operation)

	// Assertions
	require.NoError(suite.T(), err)
	assert.False(suite.T(), valid, "Access should be denied for terminated LCT")
	assert.NotEmpty(suite.T(), reason, "Reason should be provided for denied access")
}

// Benchmark tests for performance
func (suite *LCTManagerTestSuite) BenchmarkCreateLCTRelationship(b *testing.B) {
	for i := 0; i < b.N; i++ {
		componentA := fmt.Sprintf("comp_a_%d", i)
		componentB := fmt.Sprintf("comp_b_%d", i)

		_, err := suite.keeper.CreateLctRelationship(
			suite.ctx,
			"cosmos1creator",
			componentA,
			componentB,
			"test_context",
			"",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *LCTManagerTestSuite) BenchmarkGenerateSplitKeyPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := suite.keeper.GenerateSplitKeyPair()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *LCTManagerTestSuite) BenchmarkValidateLCTAccess(b *testing.B) {
	// Setup: Create LCT
	lct, err := suite.keeper.CreateLctRelationship(
		suite.ctx,
		"cosmos1creator",
		"comp_a",
		"comp_b",
		"test_context",
		"",
	)
	if err != nil {
		b.Fatal(err)
	}

	err = suite.keeper.UpdateLCTStatus(suite.ctx, lct.LctId, "active")
	if err != nil {
		b.Fatal(err)
	}

	// Benchmark access validation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := suite.keeper.ValidateLCTAccess(suite.ctx, "cosmos1requester", lct.LctId, "energy_transfer")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestLCTManagerTestSuite(t *testing.T) {
	suite.Run(t, new(LCTManagerTestSuite))
}
