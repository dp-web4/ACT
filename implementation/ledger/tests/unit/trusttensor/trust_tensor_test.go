package trusttensor_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	keepertest "racecar-web/testutil/keeper"
	"racecar-web/x/trusttensor/keeper"
	"racecar-web/x/trusttensor/types"
)

type TrustTensorTestSuite struct {
	suite.Suite
	ctx    context.Context
	keeper keeper.Keeper
}

func (suite *TrustTensorTestSuite) SetupTest() {
	suite.ctx, suite.keeper = keepertest.TrustTensorKeeper(suite.T())
}

// TestCalculateRelationshipTrust tests the calculation of relationship trust scores
func (suite *TrustTensorTestSuite) TestCalculateRelationshipTrust() {
	testCases := []struct {
		name               string
		lctID              string
		operationalContext string
		expectedScoreMin   string
		expectedScoreMax   string
		expectError        bool
	}{
		{
			name:               "Valid LCT Relationship",
			lctID:              "lct_battery_pack_001",
			operationalContext: "energy_operation",
			expectedScoreMin:   "0.0",
			expectedScoreMax:   "1.0",
			expectError:        false,
		},
		{
			name:               "Emergency Context",
			lctID:              "lct_emergency_controller",
			operationalContext: "emergency_operation",
			expectedScoreMin:   "0.0",
			expectedScoreMax:   "1.0",
			expectError:        false,
		},
		{
			name:               "Empty LCT ID",
			lctID:              "",
			operationalContext: "energy_operation",
			expectedScoreMin:   "0.0",
			expectedScoreMax:   "1.0",
			expectError:        false, // Should return default score
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Calculate relationship trust
			trustScore, factors, err := suite.keeper.CalculateRelationshipTrust(
				suite.ctx,
				tc.lctID,
				tc.operationalContext,
			)

			if tc.expectError {
				require.Error(t, err)
				assert.Empty(t, trustScore)
				assert.Empty(t, factors)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, trustScore)
				assert.NotEmpty(t, factors)

				// Parse expected bounds
				minScore, err := math.LegacyNewDecFromStr(tc.expectedScoreMin)
				require.NoError(t, err)
				maxScore, err := math.LegacyNewDecFromStr(tc.expectedScoreMax)
				require.NoError(t, err)
				actualScore, err := math.LegacyNewDecFromStr(trustScore)
				require.NoError(t, err)

				// Verify trust score is within expected range
				assert.True(t, actualScore.GTE(minScore), "Trust score should be >= min")
				assert.True(t, actualScore.LTE(maxScore), "Trust score should be <= max")

				// Verify factors string contains relevant information
				assert.Contains(t, factors, "t3_score", "Factors should contain T3 score information")
			}
		})
	}
}

// TestGetRelationshipTensor tests retrieving relationship tensors
func (suite *TrustTensorTestSuite) TestGetRelationshipTensor() {
	// Setup: Create a relationship tensor
	lctID := "lct_test_relationship"
	tensor := types.RelationshipTrustTensor{
		LctId:            lctID,
		TensorType:       "relationship",
		TalentScore:      "0.8",
		TrainingScore:    "0.7",
		TemperamentScore: "0.9",
		Context:          "test_context",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		Version:          1,
		EvidenceCount:    3,
		ContextModifier:  "1.0",
	}

	// Store the tensor
	err := suite.keeper.SetRelationshipTensor(suite.ctx, lctID, tensor)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		lctID       string
		expectFound bool
	}{
		{
			name:        "Existing Tensor",
			lctID:       lctID,
			expectFound: true,
		},
		{
			name:        "Non-existent Tensor",
			lctID:       "lct_nonexistent",
			expectFound: false,
		},
		{
			name:        "Empty LCT ID",
			lctID:       "",
			expectFound: false,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Get relationship tensor
			retrievedTensor, found := suite.keeper.GetRelationshipTensor(suite.ctx, tc.lctID)

			assert.Equal(t, tc.expectFound, found)
			if tc.expectFound {
				assert.Equal(t, tc.lctID, retrievedTensor.LctId)
				assert.Equal(t, tensor.TalentScore, retrievedTensor.TalentScore)
				assert.Equal(t, tensor.TrainingScore, retrievedTensor.TrainingScore)
				assert.Equal(t, tensor.TemperamentScore, retrievedTensor.TemperamentScore)
			}
		})
	}
}

// TestSetRelationshipTensor tests storing relationship tensors
func (suite *TrustTensorTestSuite) TestSetRelationshipTensor() {
	testCases := []struct {
		name        string
		lctID       string
		tensor      types.RelationshipTrustTensor
		expectError bool
	}{
		{
			name:  "Valid Tensor",
			lctID: "lct_set_test_001",
			tensor: types.RelationshipTrustTensor{
				LctId:            "lct_set_test_001",
				TensorType:       "relationship",
				TalentScore:      "0.7",
				TrainingScore:    "0.6",
				TemperamentScore: "0.8",
				Context:          "test_context",
				CreatedAt:        time.Now().Unix(),
				UpdatedAt:        time.Now().Unix(),
				Version:          1,
				EvidenceCount:    2,
				ContextModifier:  "1.0",
			},
			expectError: false,
		},
		{
			name:  "High Trust Score",
			lctID: "lct_set_test_002",
			tensor: types.RelationshipTrustTensor{
				LctId:            "lct_set_test_002",
				TensorType:       "relationship",
				TalentScore:      "0.95",
				TrainingScore:    "0.9",
				TemperamentScore: "0.92",
				Context:          "high_performance",
				CreatedAt:        time.Now().Unix(),
				UpdatedAt:        time.Now().Unix(),
				Version:          1,
				EvidenceCount:    3,
				ContextModifier:  "1.0",
			},
			expectError: false,
		},
		{
			name:  "Empty LCT ID",
			lctID: "",
			tensor: types.RelationshipTrustTensor{
				LctId:            "",
				TensorType:       "relationship",
				TalentScore:      "0.5",
				TrainingScore:    "0.5",
				TemperamentScore: "0.5",
				Context:          "empty_context",
				CreatedAt:        time.Now().Unix(),
				UpdatedAt:        time.Now().Unix(),
				Version:          1,
				EvidenceCount:    0,
				ContextModifier:  "1.0",
			},
			expectError: false, // Should not error, but may not be useful
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Set relationship tensor
			err := suite.keeper.SetRelationshipTensor(suite.ctx, tc.lctID, tc.tensor)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify tensor was stored
				retrievedTensor, found := suite.keeper.GetRelationshipTensor(suite.ctx, tc.lctID)
				require.True(t, found)
				assert.Equal(t, tc.lctID, retrievedTensor.LctId)
				assert.Equal(t, tc.tensor.TalentScore, retrievedTensor.TalentScore)
				assert.Equal(t, tc.tensor.TrainingScore, retrievedTensor.TrainingScore)
				assert.Equal(t, tc.tensor.TemperamentScore, retrievedTensor.TemperamentScore)
			}
		})
	}
}

// TestUpdateTensorScore tests updating tensor scores
func (suite *TrustTensorTestSuite) TestUpdateTensorScore() {
	// Setup: Create a relationship tensor
	lctID := "lct_update_test"
	tensor := types.RelationshipTrustTensor{
		LctId:            lctID,
		TensorType:       "relationship",
		TalentScore:      "0.5",
		TrainingScore:    "0.6",
		TemperamentScore: "0.7",
		Context:          "update_test",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		Version:          1,
		EvidenceCount:    2,
		ContextModifier:  "1.0",
	}

	err := suite.keeper.SetRelationshipTensor(suite.ctx, lctID, tensor)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		tensorID    string
		dimension   string
		newScore    string
		evidence    string
		expectError bool
	}{
		{
			name:        "Valid Score Update",
			tensorID:    lctID,
			dimension:   "reliability",
			newScore:    "0.8",
			evidence:    "successful_operation",
			expectError: false,
		},
		{
			name:        "Performance Score Update",
			tensorID:    lctID,
			dimension:   "performance",
			newScore:    "0.9",
			evidence:    "high_efficiency",
			expectError: false,
		},
		{
			name:        "Non-existent Tensor",
			tensorID:    "lct_nonexistent",
			dimension:   "reliability",
			newScore:    "0.8",
			evidence:    "test_evidence",
			expectError: true,
		},
		{
			name:        "Invalid Score (too high)",
			tensorID:    lctID,
			dimension:   "reliability",
			newScore:    "1.5",
			evidence:    "test_evidence",
			expectError: false, // Should be clamped to 1.0
		},
		{
			name:        "Negative Score",
			tensorID:    lctID,
			dimension:   "reliability",
			newScore:    "-0.1",
			evidence:    "test_evidence",
			expectError: false, // Should be clamped to 0.0
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Parse new score
			newScore, err := math.LegacyNewDecFromStr(tc.newScore)
			require.NoError(t, err)

			// Update tensor score
			err = suite.keeper.UpdateTensorScore(
				suite.ctx,
				tc.tensorID,
				tc.dimension,
				newScore,
				tc.evidence,
			)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify tensor was updated
				updatedTensor, found := suite.keeper.GetRelationshipTensor(suite.ctx, tc.tensorID)
				if found {
					// Check that the tensor was updated (LastUpdated should be newer)
					assert.True(t, updatedTensor.UpdatedAt > tensor.UpdatedAt)
				}
			}
		})
	}
}

// TestCalculateT3CompositeScore tests T3 composite score calculation
func (suite *TrustTensorTestSuite) TestCalculateT3CompositeScore() {
	// Setup: Create a relationship tensor with some history
	lctID := "lct_t3_test"
	tensor := types.RelationshipTrustTensor{
		LctId:            lctID,
		TensorType:       "relationship",
		TalentScore:      "0.6",
		TrainingScore:    "0.7",
		TemperamentScore: "0.8",
		Context:          "test_context",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		Version:          1,
		EvidenceCount:    3,
		ContextModifier:  "1.0",
	}

	err := suite.keeper.SetRelationshipTensor(suite.ctx, lctID, tensor)
	require.NoError(suite.T(), err)

	testCases := []struct {
		name        string
		lctID       string
		expectError bool
		expectedMin string
		expectedMax string
	}{
		{
			name:        "Valid LCT",
			lctID:       lctID,
			expectError: false,
			expectedMin: "0.0",
			expectedMax: "1.0",
		},
		{
			name:        "Non-existent LCT",
			lctID:       "lct_nonexistent_t3",
			expectError: false, // Should return default score
			expectedMin: "0.0",
			expectedMax: "1.0",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Calculate T3 composite score
			t3Score, err := suite.keeper.CalculateT3CompositeScore(suite.ctx, tc.lctID)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Parse expected bounds
				minScore, err := math.LegacyNewDecFromStr(tc.expectedMin)
				require.NoError(t, err)
				maxScore, err := math.LegacyNewDecFromStr(tc.expectedMax)
				require.NoError(t, err)

				// Verify T3 score is within expected range
				assert.True(t, t3Score.GTE(minScore), "T3 score should be >= min")
				assert.True(t, t3Score.LTE(maxScore), "T3 score should be <= max")
			}
		})
	}
}

// TestGetContextModifier tests context modifier calculation
func (suite *TrustTensorTestSuite) TestGetContextModifier() {
	testCases := []struct {
		name        string
		context     string
		expectedMin string
		expectedMax string
	}{
		{
			name:        "Energy Operation Context",
			context:     "energy_operation",
			expectedMin: "0.5",
			expectedMax: "1.5",
		},
		{
			name:        "Emergency Context",
			context:     "emergency_operation",
			expectedMin: "0.8",
			expectedMax: "1.2",
		},
		{
			name:        "Standard Context",
			context:     "standard_operation",
			expectedMin: "0.9",
			expectedMax: "1.1",
		},
		{
			name:        "Empty Context",
			context:     "",
			expectedMin: "0.5",
			expectedMax: "1.5",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// Get context modifier
			modifier := suite.keeper.GetContextModifier(suite.ctx, tc.context)

			// Parse expected bounds
			minModifier, err := math.LegacyNewDecFromStr(tc.expectedMin)
			require.NoError(t, err)
			maxModifier, err := math.LegacyNewDecFromStr(tc.expectedMax)
			require.NoError(t, err)

			// Verify modifier is within expected range
			assert.True(t, modifier.GTE(minModifier), "Context modifier should be >= min")
			assert.True(t, modifier.LTE(maxModifier), "Context modifier should be <= max")
		})
	}
}

// Benchmark tests for performance
func (suite *TrustTensorTestSuite) BenchmarkCalculateRelationshipTrust(b *testing.B) {
	lctID := "lct_benchmark"
	context := "energy_operation"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := suite.keeper.CalculateRelationshipTrust(suite.ctx, lctID, context)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *TrustTensorTestSuite) BenchmarkCalculateT3CompositeScore(b *testing.B) {
	lctID := "lct_benchmark_t3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := suite.keeper.CalculateT3CompositeScore(suite.ctx, lctID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func (suite *TrustTensorTestSuite) BenchmarkUpdateTensorScore(b *testing.B) {
	// Setup: Create a tensor for benchmarking
	lctID := "lct_benchmark_update"
	tensor := types.RelationshipTrustTensor{
		LctId:            lctID,
		TensorType:       "relationship",
		TalentScore:      "0.5",
		TrainingScore:    "0.6",
		TemperamentScore: "0.7",
		Context:          "update_test",
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		Version:          1,
		EvidenceCount:    2,
		ContextModifier:  "1.0",
	}

	err := suite.keeper.SetRelationshipTensor(suite.ctx, lctID, tensor)
	require.NoError(b, err)

	newScore := math.LegacyNewDecWithPrec(8, 1) // 0.8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := suite.keeper.UpdateTensorScore(
			suite.ctx,
			lctID,
			"reliability",
			newScore,
			"benchmark_evidence",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestTrustTensorTestSuite(t *testing.T) {
	suite.Run(t, new(TrustTensorTestSuite))
}
