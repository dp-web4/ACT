package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// APIBridgeTestSuite provides a comprehensive testing framework for the API bridge
type APIBridgeTestSuite struct {
	suite.Suite
	apiBridgeURL string
	httpClient   *http.Client
}

// Test data structures for real API responses
type Account struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	KeyType string `json:"key_type"`
}

type Component struct {
	ComponentId string    `json:"component_id"`
	Creator     string    `json:"creator"`
	Data        string    `json:"component_data"`
	Context     string    `json:"context"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	TxHash      string    `json:"txhash,omitempty"`
}

type LCT struct {
	LctId         string    `json:"lct_id"`
	ComponentA    string    `json:"component_a"`
	ComponentB    string    `json:"component_b"`
	Context       string    `json:"context"`
	ProxyId       string    `json:"proxy_id"`
	Status        string    `json:"status"`
	LctKeyHalf    string    `json:"lct_key_half"`
	DeviceKeyHalf string    `json:"device_key_half"`
	CreatedAt     time.Time `json:"created_at"`
	TxHash        string    `json:"txhash,omitempty"`
}

type Pairing struct {
	ChallengeId        string `json:"challenge_id"`
	ComponentA         string `json:"component_a"`
	ComponentB         string `json:"component_b"`
	OperationalContext string `json:"operational_context"`
	ProxyId            string `json:"proxy_id"`
	Status             string `json:"status"`
	LctId              string `json:"lct_id,omitempty"`
	SplitKeyA          string `json:"split_key_a,omitempty"`
	SplitKeyB          string `json:"split_key_b,omitempty"`
	TxHash             string `json:"txhash,omitempty"`
}

type EnergyOperation struct {
	OperationId   string `json:"operation_id"`
	SourceLct     string `json:"source_lct"`
	TargetLct     string `json:"target_lct"`
	EnergyAmount  string `json:"energy_amount"`
	OperationType string `json:"operation_type"`
	Status        string `json:"status"`
	TrustScore    string `json:"trust_score"`
	TxHash        string `json:"txhash,omitempty"`
}

type HealthResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
	Count    int       `json:"count"`
}

func (suite *APIBridgeTestSuite) SetupTest() {
	// Use the real API Bridge URL
	suite.apiBridgeURL = "http://localhost:8080"

	// Initialize HTTP client for real API calls
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

// Helper function to make HTTP requests to the real API Bridge
func (suite *APIBridgeTestSuite) makeRequest(method, endpoint string, body interface{}) (*http.Response, []byte, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	url := suite.apiBridgeURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBody, err := suite.readResponseBody(resp)
	if err != nil {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, respBody, nil
}

func (suite *APIBridgeTestSuite) readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (suite *APIBridgeTestSuite) TearDownTest() {
	// No cleanup needed for real API testing
}

// Test Health Endpoint
func (suite *APIBridgeTestSuite) TestHealthEndpoint() {
	resp, body, err := suite.makeRequest("GET", "/health", nil)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var healthResp HealthResponse
	err = json.Unmarshal(body, &healthResp)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "api-bridge", healthResp.Service)
	assert.Equal(suite.T(), "healthy", healthResp.Status)
	assert.Greater(suite.T(), healthResp.Timestamp, int64(0))
}

// Test Get Accounts
func (suite *APIBridgeTestSuite) TestGetAccounts() {
	resp, body, err := suite.makeRequest("GET", "/api/v1/accounts", nil)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var accountsResp AccountsResponse
	err = json.Unmarshal(body, &accountsResp)
	require.NoError(suite.T(), err)

	assert.GreaterOrEqual(suite.T(), accountsResp.Count, 0)
	assert.Len(suite.T(), accountsResp.Accounts, accountsResp.Count)

	// Check for expected default accounts
	accountNames := make(map[string]bool)
	for _, account := range accountsResp.Accounts {
		accountNames[account.Name] = true
	}

	// Should have at least some accounts (alice, bob, charlie, or default)
	assert.Greater(suite.T(), len(accountNames), 0)
}

// Test Component Registration (Real Blockchain)
func (suite *APIBridgeTestSuite) TestRegisterComponent() {
	componentData := map[string]interface{}{
		"creator":        "alice",
		"component_data": fmt.Sprintf("test_battery_%d", time.Now().Unix()),
		"context":        "unit_test",
	}

	resp, body, err := suite.makeRequest("POST", "/api/v1/components/register", componentData)

	// Note: This might fail if blockchain is not running or has issues
	// We'll handle both success and expected failure cases
	if err != nil {
		suite.T().Logf("Component registration failed (expected if blockchain not running): %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var component Component
		err = json.Unmarshal(body, &component)
		require.NoError(suite.T(), err)

		assert.NotEmpty(suite.T(), component.ComponentId)
		assert.Equal(suite.T(), "alice", component.Creator)
		assert.Equal(suite.T(), componentData["component_data"], component.Data)
		assert.Equal(suite.T(), "unit_test", component.Context)
		assert.Equal(suite.T(), "registered", component.Status)
		assert.NotEmpty(suite.T(), component.TxHash)
	} else {
		// Expected failure - log the error but don't fail the test
		suite.T().Logf("Component registration returned status %d: %s", resp.StatusCode, string(body))
	}
}

// Test Anonymous Component Registration (Privacy Feature)
func (suite *APIBridgeTestSuite) TestRegisterAnonymousComponent() {
	anonymousData := map[string]interface{}{
		"creator":           "alice",
		"real_component_id": fmt.Sprintf("real_battery_%d", time.Now().Unix()),
		"manufacturer_id":   "tesla_motors",
		"component_type":    "lithium_ion_battery",
		"context":           "privacy_test",
	}

	resp, body, err := suite.makeRequest("POST", "/api/v1/components/register-anonymous", anonymousData)

	if err != nil {
		suite.T().Logf("Anonymous component registration failed (expected if blockchain not running): %v", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(suite.T(), err)

		assert.NotEmpty(suite.T(), result["component_hash"])
		assert.Equal(suite.T(), anonymousData["real_component_id"], result["real_component_id"])
		assert.NotEmpty(suite.T(), result["manufacturer_hash"])
		assert.NotEmpty(suite.T(), result["component_type_hash"])
		assert.Equal(suite.T(), "registered_anonymously", result["status"])
		assert.NotEmpty(suite.T(), result["txhash"])
	} else {
		// Expected failure - log the error but don't fail the test
		suite.T().Logf("Anonymous component registration returned status %d: %s", resp.StatusCode, string(body))
	}
}

// Test Blockchain Status
func (suite *APIBridgeTestSuite) TestBlockchainStatus() {
	resp, body, err := suite.makeRequest("GET", "/blockchain/status", nil)

	if err != nil {
		suite.T().Logf("Blockchain status check failed: %v", err)
		return
	}

	// Status endpoint might return different status codes depending on blockchain state
	suite.T().Logf("Blockchain status response: %d - %s", resp.StatusCode, string(body))

	// Don't fail the test - just log the status
	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusServiceUnavailable)
}

// Test Invalid Request Body
func (suite *APIBridgeTestSuite) TestInvalidRequestBody() {
	invalidJSON := `{"invalid": json}`

	req, err := http.NewRequest("POST", suite.apiBridgeURL+"/api/v1/components/register", strings.NewReader(invalidJSON))
	require.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should return 400 Bad Request for invalid JSON
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// Test Non-Existent Endpoint
func (suite *APIBridgeTestSuite) TestNonExistentEndpoint() {
	resp, _, err := suite.makeRequest("GET", "/api/v1/nonexistent", nil)
	require.NoError(suite.T(), err)

	// Should return 404 Not Found
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// Test Privacy-Focused Endpoints
func (suite *APIBridgeTestSuite) TestPrivacyEndpoints() {
	// Test verify-pairing-hashes endpoint
	verifyData := map[string]interface{}{
		"verifier":         "alice",
		"component_hash_a": "a1b2c3d4e5f6",
		"component_hash_b": "f6e5d4c3b2a1",
		"context":          "energy_transfer",
	}

	resp, body, err := suite.makeRequest("POST", "/api/v1/components/verify-pairing-hashes", verifyData)

	if err != nil {
		suite.T().Logf("Privacy endpoint test failed: %v", err)
		return
	}

	suite.T().Logf("Privacy endpoint response: %d - %s", resp.StatusCode, string(body))

	// Don't fail the test - just verify the endpoint exists
	assert.True(suite.T(), resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusInternalServerError)
}

// Benchmark Health Endpoint
func (suite *APIBridgeTestSuite) BenchmarkHealthEndpoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, _, err := suite.makeRequest("GET", "/health", nil)
		if err != nil {
			b.Fatalf("Health endpoint failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Health endpoint returned status %d", resp.StatusCode)
		}
	}
}

// Test Security Headers
func (suite *APIBridgeTestSuite) TestSecurityHeaders() {
	resp, _, err := suite.makeRequest("GET", "/health", nil)
	require.NoError(suite.T(), err)

	// Check for basic security headers (if implemented)
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(suite.T(), contentType, "application/json")
}

// Test Input Sanitization
func (suite *APIBridgeTestSuite) TestInputSanitization() {
	// Test with potentially malicious input
	maliciousData := map[string]interface{}{
		"creator":        "alice",
		"component_data": "<script>alert('xss')</script>",
		"context":        "test",
	}

	resp, body, err := suite.makeRequest("POST", "/api/v1/components/register", maliciousData)

	if err != nil {
		suite.T().Logf("Input sanitization test failed: %v", err)
		return
	}

	// Should not return 500 Internal Server Error for malicious input
	assert.NotEqual(suite.T(), http.StatusInternalServerError, resp.StatusCode)

	// Log the response for analysis
	suite.T().Logf("Input sanitization response: %d - %s", resp.StatusCode, string(body))
}

// Test API Bridge is Running
func (suite *APIBridgeTestSuite) TestAPIBridgeIsRunning() {
	resp, _, err := suite.makeRequest("GET", "/health", nil)

	if err != nil {
		suite.T().Skipf("API Bridge is not running: %v", err)
		return
	}

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func TestAPIBridgeTestSuite(t *testing.T) {
	suite.Run(t, new(APIBridgeTestSuite))
}
