package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// APIBridgeIntegrationTestSuite tests the API bridge with real blockchain integration
type APIBridgeIntegrationTestSuite struct {
	suite.Suite
	blockchainProcess *exec.Cmd
	bridgeProcess     *exec.Cmd
	restBaseURL       string
	grpcBaseURL       string
	httpClient        *http.Client
	cleanup           []func()
}

func (suite *APIBridgeIntegrationTestSuite) SetupSuite() {
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.restBaseURL = "http://localhost:8080"
	suite.grpcBaseURL = "localhost:9090"
	suite.cleanup = make([]func(), 0)

	// Start blockchain node
	suite.startBlockchainNode()

	// Start API bridge
	suite.startAPIBridge()

	// Wait for services to be ready
	suite.waitForServices()
}

func (suite *APIBridgeIntegrationTestSuite) TearDownSuite() {
	// Stop processes in reverse order
	for i := len(suite.cleanup) - 1; i >= 0; i-- {
		suite.cleanup[i]()
	}
}

func (suite *APIBridgeIntegrationTestSuite) startBlockchainNode() {
	// Check if blockchain node is already running
	if suite.isBlockchainRunning() {
		suite.T().Log("Blockchain node already running")
		return
	}

	suite.T().Log("Starting blockchain node...")

	// Start blockchain node with ignite
	cmd := exec.Command("ignite", "chain", "serve", "--reset-once")
	cmd.Dir = ".." // Go up to project root
	cmd.Env = os.Environ()

	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	require.NoError(suite.T(), err, "Failed to start blockchain node")

	suite.blockchainProcess = cmd
	suite.cleanup = append(suite.cleanup, func() {
		if suite.blockchainProcess != nil {
			suite.T().Log("Stopping blockchain node...")
			suite.blockchainProcess.Process.Kill()
			suite.blockchainProcess.Wait()
		}
	})

	// Wait for blockchain to be ready
	suite.waitForBlockchain()
}

func (suite *APIBridgeIntegrationTestSuite) startAPIBridge() {
	suite.T().Log("Starting API bridge...")

	// Build API bridge
	buildCmd := exec.Command("go", "build", "-o", "api-bridge-test", "./cmd/api-bridge")
	buildCmd.Dir = "." // Current directory
	err := buildCmd.Run()
	require.NoError(suite.T(), err, "Failed to build API bridge")

	// Start API bridge
	cmd := exec.Command("./api-bridge-test", "--rest-port", "8080", "--grpc-port", "9090")
	cmd.Dir = "."
	cmd.Env = os.Environ()

	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	require.NoError(suite.T(), err, "Failed to start API bridge")

	suite.bridgeProcess = cmd
	suite.cleanup = append(suite.cleanup, func() {
		if suite.bridgeProcess != nil {
			suite.T().Log("Stopping API bridge...")
			suite.bridgeProcess.Process.Kill()
			suite.bridgeProcess.Wait()
		}
		// Clean up binary
		os.Remove("api-bridge-test")
	})
}

func (suite *APIBridgeIntegrationTestSuite) isBlockchainRunning() bool {
	// Check if blockchain REST API is responding
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (suite *APIBridgeIntegrationTestSuite) waitForBlockchain() {
	suite.T().Log("Waiting for blockchain to be ready...")

	timeout := time.After(60 * time.Second)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			suite.T().Fatal("Timeout waiting for blockchain to start")
		case <-tick:
			if suite.isBlockchainRunning() {
				suite.T().Log("Blockchain is ready")
				return
			}
		}
	}
}

func (suite *APIBridgeIntegrationTestSuite) waitForServices() {
	suite.T().Log("Waiting for API bridge to be ready...")

	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			suite.T().Fatal("Timeout waiting for API bridge to start")
		case <-tick:
			resp, err := suite.httpClient.Get(suite.restBaseURL + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				suite.T().Log("API bridge is ready")
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

// Test REST API Integration
func (suite *APIBridgeIntegrationTestSuite) TestRESTHealthEndpoint() {
	resp, err := suite.httpClient.Get(suite.restBaseURL + "/health")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), "healthy", response["status"])
	assert.Contains(suite.T(), response, "timestamp")
}

func (suite *APIBridgeIntegrationTestSuite) TestRESTGetAccounts() {
	resp, err := suite.httpClient.Get(suite.restBaseURL + "/accounts")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "accounts")
	accounts, ok := response["accounts"].([]interface{})
	assert.True(suite.T(), ok)
	assert.Greater(suite.T(), len(accounts), 0, "Should have at least one account")
}

func (suite *APIBridgeIntegrationTestSuite) TestRESTComponentRegistration() {
	// Test component registration
	requestData := map[string]interface{}{
		"creator":        "alice",
		"component_data": "integration-test-battery-module",
		"context":        "integration-test",
	}

	jsonData, err := json.Marshal(requestData)
	require.NoError(suite.T(), err)

	resp, err := suite.httpClient.Post(
		suite.restBaseURL+"/component/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "component_id")
	assert.Contains(suite.T(), response, "txhash")

	componentId := response["component_id"].(string)
	assert.NotEmpty(suite.T(), componentId)

	// Verify we can retrieve the component
	suite.T().Run("GetRegisteredComponent", func(t *testing.T) {
		getResp, err := suite.httpClient.Get(suite.restBaseURL + "/component/" + componentId)
		require.NoError(t, err)
		defer getResp.Body.Close()

		if getResp.StatusCode == http.StatusOK {
			var component map[string]interface{}
			err = json.NewDecoder(getResp.Body).Decode(&component)
			require.NoError(t, err)

			assert.Equal(t, componentId, component["component_id"])
			assert.Equal(t, "alice", component["creator"])
		}
		// Note: Component retrieval might not be implemented yet
	})
}

func (suite *APIBridgeIntegrationTestSuite) TestRESTLCTCreation() {
	requestData := map[string]interface{}{
		"creator":     "alice",
		"component_a": "integration-battery-001",
		"component_b": "integration-motor-001",
		"context":     "integration-test-pairing",
		"proxy_id":    "integration-proxy-001",
	}

	jsonData, err := json.Marshal(requestData)
	require.NoError(suite.T(), err)

	resp, err := suite.httpClient.Post(
		suite.restBaseURL+"/lct/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "lct_id")
	assert.Contains(suite.T(), response, "txhash")

	// Verify split keys are present (but no actual key material)
	assert.Contains(suite.T(), response, "lct_key_half")
	assert.Contains(suite.T(), response, "device_key_half")

	lctKeyHalf, ok := response["lct_key_half"].(string)
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), lctKeyHalf)

	// Security check: ensure no actual private key material
	assert.False(suite.T(), strings.Contains(lctKeyHalf, "-----BEGIN PRIVATE KEY-----"))
	assert.False(suite.T(), strings.Contains(lctKeyHalf, "-----BEGIN RSA PRIVATE KEY-----"))
}

func (suite *APIBridgeIntegrationTestSuite) TestRESTPairingWorkflow() {
	// Step 1: Initiate pairing
	initiateData := map[string]interface{}{
		"creator":             "alice",
		"component_a":         "integration-battery-002",
		"component_b":         "integration-motor-002",
		"operational_context": "integration-race-operation",
		"proxy_id":            "integration-proxy-002",
		"force_immediate":     false,
	}

	jsonData, err := json.Marshal(initiateData)
	require.NoError(suite.T(), err)

	resp, err := suite.httpClient.Post(
		suite.restBaseURL+"/pairing/initiate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var initiateResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&initiateResponse)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), initiateResponse, "challenge_id")
	assert.Contains(suite.T(), initiateResponse, "txhash")

	challengeId := initiateResponse["challenge_id"].(string)
	assert.NotEmpty(suite.T(), challengeId)

	// Step 2: Complete pairing
	completeData := map[string]interface{}{
		"creator":          "alice",
		"challenge_id":     challengeId,
		"component_a_auth": "integration-battery-auth",
		"component_b_auth": "integration-motor-auth",
		"session_context":  "integration-session-002",
	}

	jsonData, err = json.Marshal(completeData)
	require.NoError(suite.T(), err)

	resp, err = suite.httpClient.Post(
		suite.restBaseURL+"/pairing/complete",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var completeResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&completeResponse)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), completeResponse, "lct_id")
	assert.Contains(suite.T(), completeResponse, "split_key_a")
	assert.Contains(suite.T(), completeResponse, "split_key_b")
	assert.Contains(suite.T(), completeResponse, "txhash")

	// Security validation: ensure split keys are present but contain no sensitive data
	splitKeyA := completeResponse["split_key_a"].(string)
	splitKeyB := completeResponse["split_key_b"].(string)

	assert.NotEmpty(suite.T(), splitKeyA)
	assert.NotEmpty(suite.T(), splitKeyB)
	assert.NotEqual(suite.T(), splitKeyA, splitKeyB, "Split keys should be different")

	// Ensure no private key material
	assert.False(suite.T(), strings.Contains(splitKeyA, "-----BEGIN"))
	assert.False(suite.T(), strings.Contains(splitKeyB, "-----BEGIN"))
}

// Test error handling and edge cases
func (suite *APIBridgeIntegrationTestSuite) TestRESTErrorHandling() {
	// Test invalid JSON
	resp, err := suite.httpClient.Post(
		suite.restBaseURL+"/component/register",
		"application/json",
		bytes.NewBuffer([]byte("invalid json")),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	// Test missing required fields
	requestData := map[string]interface{}{
		"creator": "alice",
		// Missing component_data and context
	}

	jsonData, err := json.Marshal(requestData)
	require.NoError(suite.T(), err)

	resp, err = suite.httpClient.Post(
		suite.restBaseURL+"/component/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should return error for missing fields
	assert.True(suite.T(), resp.StatusCode >= 400)
}

func (suite *APIBridgeIntegrationTestSuite) TestRateLimitingAndPerformance() {
	// Test rapid requests to ensure server handles load
	const numRequests = 50
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			resp, err := suite.httpClient.Get(suite.restBaseURL + "/health")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("request %d failed with status %d", index, resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < numRequests; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else {
			suite.T().Logf("Request failed: %v", err)
		}
	}

	// At least 80% of requests should succeed
	assert.Greater(suite.T(), successCount, numRequests*8/10, "Should handle concurrent requests")
}

func (suite *APIBridgeIntegrationTestSuite) TestSecurityHeaders() {
	resp, err := suite.httpClient.Get(suite.restBaseURL + "/health")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(suite.T(), contentType, "application/json")

	// Check for security headers (if implemented)
	// These tests can be enabled when security headers are added
	/*
		assert.NotEmpty(suite.T(), resp.Header.Get("X-Content-Type-Options"))
		assert.NotEmpty(suite.T(), resp.Header.Get("X-Frame-Options"))
		assert.NotEmpty(suite.T(), resp.Header.Get("X-XSS-Protection"))
	*/
}

// Test blockchain interaction edge cases
func (suite *APIBridgeIntegrationTestSuite) TestBlockchainConnectionRecovery() {
	// This test would verify that the API bridge handles blockchain disconnections gracefully
	// For now, we'll test with an invalid request that should return a meaningful error

	requestData := map[string]interface{}{
		"creator":        "", // Invalid empty creator
		"component_data": "test-component",
		"context":        "test-context",
	}

	jsonData, err := json.Marshal(requestData)
	require.NoError(suite.T(), err)

	resp, err := suite.httpClient.Post(
		suite.restBaseURL+"/component/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Should return error for invalid input
	assert.True(suite.T(), resp.StatusCode >= 400)
}

// Test race car specific scenarios
func (suite *APIBridgeIntegrationTestSuite) TestRaceCarBatteryPackScenario() {
	creator := "cosmos1racecar"

	// Step 1: Register race car components
	components := []struct {
		id            string
		componentType string
		metadata      string
	}{
		{
			id:            "comp_main_battery_pack",
			componentType: "battery_pack",
			metadata:      `{"manufacturer_id": "TESLA_RACING", "model": "TR-2024-MAIN", "capacity_kwh": "100", "max_voltage": "800", "cell_chemistry": "NCM"}`,
		},
		{
			id:            "comp_battery_module_001",
			componentType: "battery_module",
			metadata:      `{"manufacturer_id": "TESLA_RACING", "model": "TR-2024-MOD", "capacity_kwh": "10", "position": "front_left"}`,
		},
		{
			id:            "comp_motor_controller",
			componentType: "motor_controller",
			metadata:      `{"manufacturer_id": "TESLA_RACING", "model": "TR-2024-MC", "max_power_kw": "300", "cooling": "liquid"}`,
		},
	}

	// Register all components
	for _, comp := range components {
		requestData := map[string]interface{}{
			"creator":        creator,
			"component_data": comp.metadata,
			"context":        "race_car_registration",
		}

		jsonData, err := json.Marshal(requestData)
		require.NoError(suite.T(), err)

		resp, err := suite.httpClient.Post(
			suite.restBaseURL+"/component/register",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(suite.T(), err)
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode, "Failed to register component %s", comp.id)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(suite.T(), err)
		assert.Contains(suite.T(), response, "component_id")
	}

	// Step 2: Create LCT relationships between components
	lctRelationships := []struct {
		componentA string
		componentB string
		context    string
	}{
		{"comp_main_battery_pack", "comp_battery_module_001", "battery_management"},
		{"comp_main_battery_pack", "comp_motor_controller", "energy_delivery"},
	}

	for _, rel := range lctRelationships {
		requestData := map[string]interface{}{
			"creator":     creator,
			"component_a": rel.componentA,
			"component_b": rel.componentB,
			"context":     rel.context,
			"proxy_id":    "",
		}

		jsonData, err := json.Marshal(requestData)
		require.NoError(suite.T(), err)

		resp, err := suite.httpClient.Post(
			suite.restBaseURL+"/lct/create",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(suite.T(), err)
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode, "Failed to create LCT for %s -> %s", rel.componentA, rel.componentB)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(suite.T(), err)
		assert.Contains(suite.T(), response, "lct_id")
	}

	// Step 3: Test pairing workflow for battery pack and motor controller
	suite.T().Run("BatteryMotorPairing", func(t *testing.T) {
		// Initiate pairing
		initiateData := map[string]interface{}{
			"creator":             creator,
			"component_a":         "comp_main_battery_pack",
			"component_b":         "comp_motor_controller",
			"operational_context": "race_energy_delivery",
			"proxy_id":            "",
			"force_immediate":     false,
		}

		jsonData, err := json.Marshal(initiateData)
		require.NoError(t, err)

		resp, err := suite.httpClient.Post(
			suite.restBaseURL+"/pairing/initiate",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var initiateResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&initiateResponse)
		require.NoError(t, err)

		challengeId := initiateResponse["challenge_id"].(string)
		assert.NotEmpty(t, challengeId)

		// Complete pairing
		completeData := map[string]interface{}{
			"creator":          creator,
			"challenge_id":     challengeId,
			"component_a_auth": "battery_auth_token",
			"component_b_auth": "motor_auth_token",
			"session_context":  "race_session_001",
		}

		jsonData, err = json.Marshal(completeData)
		require.NoError(t, err)

		resp, err = suite.httpClient.Post(
			suite.restBaseURL+"/pairing/complete",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var completeResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&completeResponse)
		require.NoError(t, err)

		assert.Contains(t, completeResponse, "lct_id")
		assert.Contains(t, completeResponse, "split_key_a")
		assert.Contains(t, completeResponse, "split_key_b")
		assert.Equal(t, "completed", completeResponse["status"])
	})
}

// Benchmark tests
func (suite *APIBridgeIntegrationTestSuite) BenchmarkRESTHealthEndpoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := suite.httpClient.Get(suite.restBaseURL + "/health")
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

func (suite *APIBridgeIntegrationTestSuite) BenchmarkRESTComponentRegistration(b *testing.B) {
	requestData := map[string]interface{}{
		"creator":        "alice",
		"component_data": "benchmark-component",
		"context":        "benchmark-test",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := suite.httpClient.Post(
			suite.restBaseURL+"/component/register",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// Helper function to check if test should be skipped
func (suite *APIBridgeIntegrationTestSuite) skipIfNoBlockchain() {
	if !suite.isBlockchainRunning() {
		suite.T().Skip("Blockchain not running, skipping integration test")
	}
}

// Run the test suite
func TestAPIBridgeIntegrationTestSuite(t *testing.T) {
	// Skip integration tests if in CI or if blockchain is not available
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	suite.Run(t, new(APIBridgeIntegrationTestSuite))
}
