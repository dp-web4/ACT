package energycycle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	keepertest "racecar-web/testutil/keeper"
	"racecar-web/x/energycycle/keeper"
)

type EnergyCycleTestSuite struct {
	suite.Suite
	ctx    context.Context
	keeper keeper.Keeper
}

func (suite *EnergyCycleTestSuite) SetupTest() {
	suite.ctx, suite.keeper = keepertest.EnergyCycleKeeper(suite.T())
}

// TestCreateAtpToken tests ATP token creation for energy storage
func (suite *EnergyCycleTestSuite) TestCreateAtpToken() {
	testCases := []struct {
		name           string
		lctID          string
		energyAmount   string
		operationID    string
		context        string
		blockHeight    int64
		expectError    bool
		expectedStatus string
	}{
		{
			name:           "Valid ATP Token Creation",
			lctID:          "lct_battery_pack_001",
			energyAmount:   "100.5",
			operationID:    "op_charge_001",
			context:        "charging_operation",
			blockHeight:    1000,
			expectError:    false,
			expectedStatus: "active",
		},
		{
			name:           "Large Energy Amount",
			lctID:          "lct_battery_pack_002",
			energyAmount:   "1000.0",
			operationID:    "op_charge_002",
			context:        "fast_charging",
			blockHeight:    1001,
			expectError:    false,
			expectedStatus: "active",
		},
		{
			name:           "Zero Energy Amount",
			lctID:          "lct_battery_pack_003",
			energyAmount:   "0.0",
			operationID:    "op_charge_003",
			context:        "test_operation",
			blockHeight:    1002,
			expectError:    true,
			expectedStatus: "",
		},
		{
			name:           "Negative Energy Amount",
			lctID:          "lct_battery_pack_004",
			energyAmount:   "-50.0",
			operationID:    "op_charge_004",
			context:        "invalid_operation",
			blockHeight:    1003,
			expectError:    true,
			expectedStatus: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			atpToken, err := suite.keeper.CreateAtpToken(
				suite.ctx,
				tc.lctID,
				tc.energyAmount,
				tc.operationID,
				tc.context,
				tc.blockHeight,
			)

			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, atpToken)
			} else {
				require.NoError(t, err)
				require.NotNil(t, atpToken)

				// Verify ATP token properties
				assert.NotEmpty(t, atpToken.TokenId)
				assert.Equal(t, tc.lctID, atpToken.LctId)
				assert.Equal(t, tc.energyAmount, atpToken.EnergyAmount)
				assert.Equal(t, tc.operationID, atpToken.OperationId)
				assert.Equal(t, tc.context, atpToken.Context)
				assert.Equal(t, tc.blockHeight, atpToken.BlockHeight)
				assert.Equal(t, tc.expectedStatus, atpToken.Status)
				assert.NotZero(t, atpToken.CreatedAt)

				// Verify token was stored
				storedToken, err := suite.keeper.GetAtpToken(suite.ctx, atpToken.TokenId)
				require.NoError(t, err)
				assert.Equal(t, atpToken.TokenId, storedToken.TokenId)
			}
		})
	}
}

// TestDischargeAtpToken tests ADP token creation through ATP discharge
func (suite *EnergyCycleTestSuite) TestDischargeAtpToken() {
	// Setup: Create an ATP token first
	lctID := "lct_discharge_test"
	energyAmount := "100.0"
	operationID := "op_charge_test"
	context := "test_charging"
	blockHeight := int64(1000)

	atpToken, err := suite.keeper.CreateAtpToken(
		suite.ctx,
		lctID,
		energyAmount,
		operationID,
		context,
		blockHeight,
	)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name           string
		atpTokenID     string
		operationID    string
		blockHeight    int64
		expectError    bool
		expectedStatus string
	}{
		{
			name:           "Valid ATP Discharge",
			atpTokenID:     atpToken.TokenId,
			operationID:    "op_discharge_test",
			blockHeight:    1001,
			expectError:    false,
			expectedStatus: "active",
		},
		{
			name:           "Invalid ATP Token ID",
			atpTokenID:     "invalid_token_id",
			operationID:    "op_discharge_invalid",
			blockHeight:    1002,
			expectError:    true,
			expectedStatus: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			adpToken, err := suite.keeper.DischargeAtpToken(
				suite.ctx,
				tc.atpTokenID,
				tc.operationID,
				tc.blockHeight,
			)

			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, adpToken)
			} else {
				require.NoError(t, err)
				require.NotNil(t, adpToken)

				// Verify ADP token properties
				assert.NotEmpty(t, adpToken.TokenId)
				assert.Equal(t, lctID, adpToken.LctId)
				assert.Equal(t, tc.operationID, adpToken.OperationId)
				assert.Equal(t, tc.blockHeight, adpToken.BlockHeight)
				assert.Equal(t, tc.expectedStatus, adpToken.Status)
				assert.NotZero(t, adpToken.CreatedAt)

				// Verify token was stored
				storedToken, err := suite.keeper.GetAdpToken(suite.ctx, adpToken.TokenId)
				require.NoError(t, err)
				assert.Equal(t, adpToken.TokenId, storedToken.TokenId)
			}
		})
	}
}

// TestValidateEnergyOperation tests energy operation validation logic
func (suite *EnergyCycleTestSuite) TestValidateEnergyOperation() {
	testCases := []struct {
		name            string
		operationID     string
		sourceLct       string
		targetLct       string
		energyAmount    string
		operationType   string
		expectedValid   bool
		expectedMessage string
	}{
		{
			name:            "Valid Transfer Operation",
			operationID:     "op_validate_001",
			sourceLct:       "lct_source_001",
			targetLct:       "lct_target_001",
			energyAmount:    "50.0",
			operationType:   "transfer",
			expectedValid:   true,
			expectedMessage: "validation successful",
		},
		{
			name:            "Valid Charge Operation",
			operationID:     "op_validate_002",
			sourceLct:       "lct_charger_001",
			targetLct:       "lct_battery_001",
			energyAmount:    "100.0",
			operationType:   "charge",
			expectedValid:   true,
			expectedMessage: "validation successful",
		},
		{
			name:            "Invalid Energy Amount",
			operationID:     "op_validate_003",
			sourceLct:       "lct_source_002",
			targetLct:       "lct_target_002",
			energyAmount:    "0.0",
			operationType:   "transfer",
			expectedValid:   false,
			expectedMessage: "energy amount must be greater than zero",
		},
		{
			name:            "Same Source and Target",
			operationID:     "op_validate_004",
			sourceLct:       "lct_same_001",
			targetLct:       "lct_same_001",
			energyAmount:    "50.0",
			operationType:   "transfer",
			expectedValid:   false,
			expectedMessage: "source and target cannot be the same",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			isValid, message, err := suite.keeper.ValidateEnergyOperation(
				suite.ctx,
				tc.operationID,
				tc.sourceLct,
				tc.targetLct,
				tc.energyAmount,
				tc.operationType,
			)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedValid, isValid)
			assert.Contains(t, message, tc.expectedMessage)
		})
	}
}

// TestCalculateEnergyBalance tests energy balance calculation
func (suite *EnergyCycleTestSuite) TestCalculateEnergyBalance() {
	// Setup: Create some ATP tokens for an LCT
	lctID := "lct_balance_test"
	energyAmounts := []string{"100.0", "50.0", "75.0"}
	operationIDs := []string{"op_charge_001", "op_charge_002", "op_charge_003"}

	for i, amount := range energyAmounts {
		_, err := suite.keeper.CreateAtpToken(
			suite.ctx,
			lctID,
			amount,
			operationIDs[i],
			"balance_test",
			int64(1000+i),
		)
		require.NoError(suite.T(), err)
	}

	testCases := []struct {
		name        string
		lctID       string
		expectError bool
		expectedMin string
		expectedMax string
	}{
		{
			name:        "Valid LCT Balance",
			lctID:       lctID,
			expectError: false,
			expectedMin: "200.0", // Sum of energy amounts
			expectedMax: "250.0",
		},
		{
			name:        "Non-existent LCT",
			lctID:       "lct_nonexistent",
			expectError: false, // Should return zero balance
			expectedMin: "0.0",
			expectedMax: "0.0",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			balance, err := suite.keeper.CalculateEnergyBalance(suite.ctx, tc.lctID)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Parse expected bounds
				minBalance, err := math.LegacyNewDecFromStr(tc.expectedMin)
				require.NoError(t, err)
				maxBalance, err := math.LegacyNewDecFromStr(tc.expectedMax)
				require.NoError(t, err)

				// Verify balance is within expected range
				assert.True(t, balance.GTE(minBalance), "Balance should be >= min")
				assert.True(t, balance.LTE(maxBalance), "Balance should be <= max")
			}
		})
	}
}

// TestGetAtpToken tests ATP token retrieval
func (suite *EnergyCycleTestSuite) TestGetAtpToken() {
	// Setup: Create an ATP token
	lctID := "lct_get_test"
	energyAmount := "100.0"
	operationID := "op_get_test"
	context := "get_test"
	blockHeight := int64(1000)

	atpToken, err := suite.keeper.CreateAtpToken(
		suite.ctx,
		lctID,
		energyAmount,
		operationID,
		context,
		blockHeight,
	)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		tokenID     string
		expectError bool
	}{
		{
			name:        "Valid Token ID",
			tokenID:     atpToken.TokenId,
			expectError: false,
		},
		{
			name:        "Invalid Token ID",
			tokenID:     "invalid_token_id",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			retrievedToken, err := suite.keeper.GetAtpToken(suite.ctx, tc.tokenID)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, atpToken.TokenId, retrievedToken.TokenId)
				assert.Equal(t, lctID, retrievedToken.LctId)
				assert.Equal(t, energyAmount, retrievedToken.EnergyAmount)
				assert.Equal(t, operationID, retrievedToken.OperationId)
			}
		})
	}
}

// TestGetAdpToken tests ADP token retrieval
func (suite *EnergyCycleTestSuite) TestGetAdpToken() {
	// Setup: Create an ATP token and discharge it to create ADP token
	lctID := "lct_adp_test"
	energyAmount := "100.0"
	chargeOpID := "op_charge_adp"
	dischargeOpID := "op_discharge_adp"

	atpToken, err := suite.keeper.CreateAtpToken(
		suite.ctx,
		lctID,
		energyAmount,
		chargeOpID,
		"adp_test",
		int64(1000),
	)
	require.NoError(suite.T(), err)

	adpToken, err := suite.keeper.DischargeAtpToken(
		suite.ctx,
		atpToken.TokenId,
		dischargeOpID,
		int64(1001),
	)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		tokenID     string
		expectError bool
	}{
		{
			name:        "Valid ADP Token ID",
			tokenID:     adpToken.TokenId,
			expectError: false,
		},
		{
			name:        "Invalid Token ID",
			tokenID:     "invalid_adp_token_id",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			retrievedToken, err := suite.keeper.GetAdpToken(suite.ctx, tc.tokenID)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, adpToken.TokenId, retrievedToken.TokenId)
				assert.Equal(t, lctID, retrievedToken.LctId)
				assert.Equal(t, dischargeOpID, retrievedToken.OperationId)
			}
		})
	}
}

// TestGetEnergyOperation tests energy operation retrieval
func (suite *EnergyCycleTestSuite) TestGetEnergyOperation() {
	// This test would require creating energy operations through the message server
	// For now, we'll test the method signature and error handling
	testCases := []struct {
		name        string
		operationID string
		expectError bool
	}{
		{
			name:        "Non-existent Operation",
			operationID: "op_nonexistent",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			_, err := suite.keeper.GetEnergyOperation(suite.ctx, tc.operationID)

			if tc.expectError {
				require.Error(t, err)
			}
		})
	}
}

// TestGetActiveAtpTokensForLct tests retrieving active ATP tokens for an LCT
func (suite *EnergyCycleTestSuite) TestGetActiveAtpTokensForLct() {
	// Setup: Create multiple ATP tokens for an LCT
	lctID := "lct_active_test"
	energyAmounts := []string{"50.0", "75.0", "100.0"}
	operationIDs := []string{"op_active_001", "op_active_002", "op_active_003"}

	for i, amount := range energyAmounts {
		_, err := suite.keeper.CreateAtpToken(
			suite.ctx,
			lctID,
			amount,
			operationIDs[i],
			"active_test",
			int64(1000+i),
		)
		require.NoError(suite.T(), err)
	}

	testCases := []struct {
		name        string
		lctID       string
		expectError bool
		expectedMin int
		expectedMax int
	}{
		{
			name:        "Valid LCT with Active Tokens",
			lctID:       lctID,
			expectError: false,
			expectedMin: 3,
			expectedMax: 3,
		},
		{
			name:        "Non-existent LCT",
			lctID:       "lct_nonexistent",
			expectError: false,
			expectedMin: 0,
			expectedMax: 0,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tokens, err := suite.keeper.getActiveAtpTokensForLct(suite.ctx, tc.lctID)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(tokens), tc.expectedMin)
				assert.LessOrEqual(t, len(tokens), tc.expectedMax)

				// Verify all tokens belong to the correct LCT
				for _, token := range tokens {
					assert.Equal(t, tc.lctID, token.LctId)
					assert.Equal(t, "active", token.Status)
				}
			}
		})
	}
}

// Benchmark tests for performance
func (suite *EnergyCycleTestSuite) BenchmarkCreateAtpToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lctID := fmt.Sprintf("lct_bench_atp_%d", i)
		operationID := fmt.Sprintf("op_bench_atp_%d", i)

		_, err := suite.keeper.CreateAtpToken(
			suite.ctx,
			lctID,
			"100.0",
			operationID,
			"benchmark",
			int64(i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *EnergyCycleTestSuite) BenchmarkDischargeAtpToken(b *testing.B) {
	// Setup: Create ATP token for benchmarking
	lctID := "lct_bench_discharge"
	atpToken, err := suite.keeper.CreateAtpToken(
		suite.ctx,
		lctID,
		"100.0",
		"op_bench_setup",
		"benchmark_setup",
		int64(1000),
	)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		operationID := fmt.Sprintf("op_bench_discharge_%d", i)
		_, err := suite.keeper.DischargeAtpToken(
			suite.ctx,
			atpToken.TokenId,
			operationID,
			int64(1001+i),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *EnergyCycleTestSuite) BenchmarkCalculateEnergyBalance(b *testing.B) {
	// Setup: Create some ATP tokens
	lctID := "lct_bench_balance"
	for i := 0; i < 10; i++ {
		_, err := suite.keeper.CreateAtpToken(
			suite.ctx,
			lctID,
			"50.0",
			fmt.Sprintf("op_bench_balance_%d", i),
			"benchmark_balance",
			int64(1000+i),
		)
		require.NoError(b, err)
	}

	// Benchmark balance calculation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := suite.keeper.CalculateEnergyBalance(suite.ctx, lctID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestEnergyCycleTestSuite(t *testing.T) {
	suite.Run(t, new(EnergyCycleTestSuite))
}
