package pairingqueue_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"racecar-web/x/pairingqueue/keeper"
	"racecar-web/x/pairingqueue/types"
)

// PairingQueueTestSuite provides a test suite for Pairing Queue functionality
type PairingQueueTestSuite struct {
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

// MockAddressCodec mocks the address codec for testing
type MockAddressCodec struct{}

func (m *MockAddressCodec) StringToBytes(text string) ([]byte, error) {
	return []byte(text), nil
}

func (m *MockAddressCodec) BytesToString(bz []byte) (string, error) {
	return string(bz), nil
}

// SetupTest initializes the test suite
func (suite *PairingQueueTestSuite) SetupTest() {
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

// TestQueuePairingRequest_ValidInput tests pairing request queuing with valid input
func (suite *PairingQueueTestSuite) TestQueuePairingRequest_ValidInput() {
	// Test data
	initiatorID := "comp_battery_pack_001"
	targetID := "comp_battery_module_001"
	requestType := "STANDARD"
	proxyID := "comp_gateway_001"

	// Queue pairing request
	requestID, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		initiatorID,
		targetID,
		requestType,
		proxyID,
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), requestID)

	// Verify request was stored
	request, found := suite.keeper.GetPairingRequest(suite.ctx, requestID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), requestID, request.RequestId)
	assert.Equal(suite.T(), initiatorID, request.InitiatorId)
	assert.Equal(suite.T(), targetID, request.TargetId)
	assert.Equal(suite.T(), requestType, request.RequestType)
	assert.Equal(suite.T(), proxyID, request.ProxyId)
	assert.Equal(suite.T(), "queued", request.Status)
	assert.NotZero(suite.T(), request.CreatedAt)
}

// TestQueuePairingRequest_ComponentNotFound tests queuing with non-existent components
func (suite *PairingQueueTestSuite) TestQueuePairingRequest_ComponentNotFound() {
	// Mock component registry to return not found
	suite.mockCompReg = &MockComponentRegistryKeeperNotFound{}

	// Queue pairing request with non-existent component
	_, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		"non_existent_initiator",
		"comp_battery_module_001",
		"STANDARD",
		"",
	)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "initiator component not found")
}

// MockComponentRegistryKeeperNotFound mocks component registry for not found scenarios
type MockComponentRegistryKeeperNotFound struct{}

func (m *MockComponentRegistryKeeperNotFound) GetComponentIdentity(ctx context.Context, componentID string) (types.ComponentIdentity, bool) {
	return types.ComponentIdentity{}, false
}

// TestQueueOfflineOperation tests offline operation queuing
func (suite *PairingQueueTestSuite) TestQueueOfflineOperation() {
	// Test data
	componentID := "comp_offline_module_001"
	operationType := "pairing"

	// Queue offline operation
	operationID, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		componentID,
		operationType,
	)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), operationID)

	// Verify operation was stored
	operation, found := suite.keeper.GetOfflineOperation(suite.ctx, operationID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), operationID, operation.OperationId)
	assert.Equal(suite.T(), componentID, operation.ComponentId)
	assert.Equal(suite.T(), operationType, operation.OperationType)
	assert.NotZero(suite.T(), operation.QueuedAt)
}

// TestProcessOfflineQueue tests offline queue processing
func (suite *PairingQueueTestSuite) TestProcessOfflineQueue() {
	// Queue multiple operations for a component
	componentID := "comp_process_test"

	// Queue pairing operation
	pairingOpID, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		componentID,
		"pairing",
	)
	require.NoError(suite.T(), err)

	// Queue unpairing operation
	unpairingOpID, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		componentID,
		"unpairing",
	)
	require.NoError(suite.T(), err)

	// Process offline queue
	processed, failed, err := suite.keeper.ProcessOfflineQueue(suite.ctx, componentID)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, processed, "Should process 2 operations")
	assert.Equal(suite.T(), 0, failed, "Should have no failures")

	// Verify operations were removed from queue
	_, found := suite.keeper.GetOfflineOperation(suite.ctx, pairingOpID)
	assert.False(suite.T(), found, "Pairing operation should be removed")

	_, found = suite.keeper.GetOfflineOperation(suite.ctx, unpairingOpID)
	assert.False(suite.T(), found, "Unpairing operation should be removed")
}

// TestProcessOfflineQueue_UnknownOperation tests processing unknown operation types
func (suite *PairingQueueTestSuite) TestProcessOfflineQueue_UnknownOperation() {
	// Queue unknown operation type
	componentID := "comp_unknown_op"
	operationID, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		componentID,
		"unknown_operation",
	)
	require.NoError(suite.T(), err)

	// Process offline queue
	processed, failed, err := suite.keeper.ProcessOfflineQueue(suite.ctx, componentID)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, processed, "Should not process unknown operation")
	assert.Equal(suite.T(), 1, failed, "Should mark unknown operation as failed")

	// Verify operation was removed from queue
	_, found := suite.keeper.GetOfflineOperation(suite.ctx, operationID)
	assert.False(suite.T(), found, "Failed operation should be removed")
}

// TestCancelRequest tests request cancellation
func (suite *PairingQueueTestSuite) TestCancelRequest() {
	// Queue a pairing request
	initiatorID := "comp_cancel_test_initiator"
	targetID := "comp_cancel_test_target"
	requestID, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		initiatorID,
		targetID,
		"STANDARD",
		"",
	)
	require.NoError(suite.T(), err)

	// Cancel the request
	reason := "user_cancelled"
	err = suite.keeper.CancelRequest(suite.ctx, requestID, reason)

	// Assertions
	require.NoError(suite.T(), err)

	// Verify request was cancelled
	request, found := suite.keeper.GetPairingRequest(suite.ctx, requestID)
	require.True(suite.T(), found)
	assert.Equal(suite.T(), "cancelled", request.Status)
	assert.Equal(suite.T(), reason, request.RequestData)
	assert.NotZero(suite.T(), request.ProcessedAt)
}

// TestCancelRequest_NotFound tests cancellation of non-existent request
func (suite *PairingQueueTestSuite) TestCancelRequest_NotFound() {
	err := suite.keeper.CancelRequest(suite.ctx, "non_existent_request", "reason")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "request not found")
}

// TestCancelRequest_AlreadyProcessed tests cancellation of already processed request
func (suite *PairingQueueTestSuite) TestCancelRequest_AlreadyProcessed() {
	// Queue and process a request
	requestID, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		"comp_a",
		"comp_b",
		"STANDARD",
		"",
	)
	require.NoError(suite.T(), err)

	// Mark as processed
	request, found := suite.keeper.GetPairingRequest(suite.ctx, requestID)
	require.True(suite.T(), found)
	request.Status = "completed"
	request.ProcessedAt = time.Now().Unix()
	// Note: In real implementation, this would be done through keeper methods

	// Try to cancel
	err = suite.keeper.CancelRequest(suite.ctx, requestID, "reason")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "can only cancel queued requests")
}

// TestGetQueuedRequests tests retrieving queued requests for a component
func (suite *PairingQueueTestSuite) TestGetQueuedRequests() {
	// Queue multiple requests involving the same component
	componentID := "comp_queued_test"

	// Request 1: component as initiator
	requestID1, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		componentID,
		"comp_target_1",
		"STANDARD",
		"",
	)
	require.NoError(suite.T(), err)

	// Request 2: component as target
	requestID2, err := suite.keeper.QueuePairingRequest(
		suite.ctx,
		"comp_initiator_2",
		componentID,
		"STANDARD",
		"",
	)
	require.NoError(suite.T(), err)

	// Request 3: different component (should not be included)
	_, err = suite.keeper.QueuePairingRequest(
		suite.ctx,
		"comp_other_1",
		"comp_other_2",
		"STANDARD",
		"",
	)
	require.NoError(suite.T(), err)

	// Get queued requests for the component
	requests, err := suite.keeper.GetQueuedRequests(suite.ctx, componentID)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), requests, 2, "Should return 2 requests for the component")

	// Verify both requests are included
	requestIDs := make(map[string]bool)
	for _, req := range requests {
		requestIDs[req.RequestId] = true
	}
	assert.True(suite.T(), requestIDs[requestID1])
	assert.True(suite.T(), requestIDs[requestID2])
}

// TestListProxyQueue tests listing proxy queue operations
func (suite *PairingQueueTestSuite) TestListProxyQueue() {
	// Queue operations for different proxy components
	proxyID := "comp_proxy_test"

	// Operation 1: for the proxy
	operationID1, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		"comp_offline_1",
		"pairing",
	)
	require.NoError(suite.T(), err)

	// Set proxy component for the operation
	operation1, found := suite.keeper.GetOfflineOperation(suite.ctx, operationID1)
	require.True(suite.T(), found)
	operation1.ProxyComponent = proxyID
	// Note: In real implementation, this would be done through keeper methods

	// Operation 2: for different proxy
	operationID2, err := suite.keeper.QueueOfflineOperation(
		suite.ctx,
		"comp_offline_2",
		"pairing",
	)
	require.NoError(suite.T(), err)

	operation2, found := suite.keeper.GetOfflineOperation(suite.ctx, operationID2)
	require.True(suite.T(), found)
	operation2.ProxyComponent = "different_proxy"

	// List proxy queue
	operations, err := suite.keeper.ListProxyQueue(suite.ctx, proxyID)

	// Assertions
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), operations, 1, "Should return 1 operation for the proxy")

	// Verify correct operation is included
	assert.Equal(suite.T(), operationID1, operations[0].OperationId)
}

// TestMultiTransportSupport tests support for different transport methods
func (suite *PairingQueueTestSuite) TestMultiTransportSupport() {
	testCases := []struct {
		name        string
		transport   string
		proxyID     string
		description string
	}{
		{
			name:        "SD Card Transport",
			transport:   "sd_card",
			proxyID:     "comp_sd_gateway",
			description: "Offline device with SD card data transfer",
		},
		{
			name:        "Bluetooth Transport",
			transport:   "bluetooth",
			proxyID:     "comp_bt_gateway",
			description: "Device with Bluetooth connectivity",
		},
		{
			name:        "WiFi Transport",
			transport:   "wifi",
			proxyID:     "comp_wifi_gateway",
			description: "Device with WiFi connectivity",
		},
		{
			name:        "CANBus Transport",
			transport:   "canbus",
			proxyID:     "comp_canbus_gateway",
			description: "Device with CANBus connectivity",
		},
		{
			name:        "Power-Line Communications",
			transport:   "plc",
			proxyID:     "comp_plc_gateway",
			description: "Device with PLC connectivity",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Queue pairing request with transport-specific proxy
			requestID, err := suite.keeper.QueuePairingRequest(
				suite.ctx,
				"comp_initiator",
				"comp_target",
				"STANDARD",
				tc.proxyID,
			)

			// Assertions
			require.NoError(t, err)
			assert.NotEmpty(t, requestID)

			// Verify request was stored with proxy
			request, found := suite.keeper.GetPairingRequest(suite.ctx, requestID)
			require.True(t, found)
			assert.Equal(t, tc.proxyID, request.ProxyId)
			assert.Equal(t, "queued", request.Status)
		})
	}
}

// TestAuthenticationControllerIntegration tests Authentication Controller integration
func (suite *PairingQueueTestSuite) TestAuthenticationControllerIntegration() {
	// Test scenarios where Authentication Controller mediates pairing
	testCases := []struct {
		name           string
		initiatorID    string
		targetID       string
		proxyID        string
		requestType    string
		expectedStatus string
	}{
		{
			name:           "Standard Authentication Controller Mediation",
			initiatorID:    "comp_battery_pack",
			targetID:       "comp_battery_module",
			proxyID:        "comp_auth_controller",
			requestType:    "STANDARD",
			expectedStatus: "queued",
		},
		{
			name:           "Emergency Authentication Controller Mediation",
			initiatorID:    "comp_emergency_system",
			targetID:       "comp_battery_pack",
			proxyID:        "comp_auth_controller",
			requestType:    "EMERGENCY",
			expectedStatus: "queued",
		},
		{
			name:           "Proxy Authentication Controller Mediation",
			initiatorID:    "comp_offline_device",
			targetID:       "comp_online_device",
			proxyID:        "comp_auth_controller",
			requestType:    "PROXY",
			expectedStatus: "queued",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Queue pairing request through Authentication Controller
			requestID, err := suite.keeper.QueuePairingRequest(
				suite.ctx,
				tc.initiatorID,
				tc.targetID,
				tc.requestType,
				tc.proxyID,
			)

			// Assertions
			require.NoError(t, err)
			assert.NotEmpty(t, requestID)

			// Verify Authentication Controller is involved
			request, found := suite.keeper.GetPairingRequest(suite.ctx, requestID)
			require.True(t, found)
			assert.Equal(t, tc.proxyID, request.ProxyId)
			assert.Equal(t, tc.expectedStatus, request.Status)
			assert.Equal(t, tc.requestType, request.RequestType)
		})
	}
}

// TestOfflineDeviceScenarios tests various offline device scenarios
func (suite *PairingQueueTestSuite) TestOfflineDeviceScenarios() {
	testCases := []struct {
		name           string
		operationType  string
		componentID    string
		expectedResult string
	}{
		{
			name:           "Offline Device Pairing Request",
			operationType:  "pairing",
			componentID:    "comp_offline_module_001",
			expectedResult: "queued",
		},
		{
			name:           "Offline Device Unpairing Request",
			operationType:  "unpairing",
			componentID:    "comp_offline_module_002",
			expectedResult: "queued",
		},
		{
			name:           "Offline Device Status Update",
			operationType:  "status_update",
			componentID:    "comp_offline_module_003",
			expectedResult: "queued",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Queue offline operation
			operationID, err := suite.keeper.QueueOfflineOperation(
				suite.ctx,
				tc.componentID,
				tc.operationType,
			)

			// Assertions
			require.NoError(t, err)
			assert.NotEmpty(t, operationID)

			// Verify operation was queued
			operation, found := suite.keeper.GetOfflineOperation(suite.ctx, operationID)
			require.True(t, found)
			assert.Equal(t, tc.componentID, operation.ComponentId)
			assert.Equal(t, tc.operationType, operation.OperationType)
			assert.NotZero(t, operation.QueuedAt)
		})
	}
}

// Benchmark tests for performance
func (suite *PairingQueueTestSuite) BenchmarkQueuePairingRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initiatorID := fmt.Sprintf("comp_initiator_%d", i)
		targetID := fmt.Sprintf("comp_target_%d", i)

		_, err := suite.keeper.QueuePairingRequest(
			suite.ctx,
			initiatorID,
			targetID,
			"STANDARD",
			"",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *PairingQueueTestSuite) BenchmarkProcessOfflineQueue(b *testing.B) {
	// Setup: Queue operations
	componentID := "comp_benchmark"
	for i := 0; i < 10; i++ {
		_, err := suite.keeper.QueueOfflineOperation(
			suite.ctx,
			componentID,
			"pairing",
		)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Benchmark queue processing
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := suite.keeper.ProcessOfflineQueue(suite.ctx, componentID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *PairingQueueTestSuite) BenchmarkGetQueuedRequests(b *testing.B) {
	// Setup: Queue requests
	componentID := "comp_benchmark"
	for i := 0; i < 10; i++ {
		_, err := suite.keeper.QueuePairingRequest(
			suite.ctx,
			componentID,
			fmt.Sprintf("comp_target_%d", i),
			"STANDARD",
			"",
		)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Benchmark request retrieval
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := suite.keeper.GetQueuedRequests(suite.ctx, componentID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestPairingQueueTestSuite(t *testing.T) {
	suite.Run(t, new(PairingQueueTestSuite))
}
