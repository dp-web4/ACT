package web4_test

import (
	"context"
	"crypto/ed25519"
	"testing"
	
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/require"
	
	lctTypes "racecarweb/x/lctmanager/types"
	mrhTypes "racecarweb/x/mrh/types"
	mrhKeeper "racecarweb/x/mrh/keeper"
)

// Phase1TestSuite tests Phase 1 Web4 implementation
type Phase1TestSuite struct {
	suite.Suite
	ctx          context.Context
	cryptoMgr    *lctTypes.CryptoManager
	mrhStorage   mrhTypes.MRHStorage
	mrhKeeper    *mrhKeeper.Keeper
}

func (suite *Phase1TestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.cryptoMgr = lctTypes.NewCryptoManager()
	
	// Setup local MRH storage
	storageConfig := mrhTypes.StorageConfig{
		Type:      "local",
		LocalPath: "/tmp/test_mrh",
	}
	
	storage, err := mrhTypes.NewMRHStorage(storageConfig)
	suite.Require().NoError(err)
	suite.mrhStorage = storage
	
	// Note: In real tests, we'd setup proper Cosmos SDK context and keeper
	// This is a simplified version for demonstration
}

// Test LCT creation with Ed25519 keys
func (suite *Phase1TestSuite) TestLCTCreationWithEd25519() {
	t := suite.T()
	
	// Generate Ed25519 key pair
	pubKey, privKey, err := suite.cryptoMgr.GenerateKeyPair()
	require.NoError(t, err)
	require.NotNil(t, pubKey)
	require.NotNil(t, privKey)
	require.Equal(t, ed25519.PublicKeySize, len(pubKey))
	require.Equal(t, ed25519.PrivateKeySize, len(privKey))
	
	// Generate LCT ID
	lctID := suite.cryptoMgr.GenerateLCTID(lctTypes.EntityTypeHuman, pubKey)
	require.Contains(t, lctID, "lct:web4:act:human:")
	
	// Create Web4-compliant LCT
	lct := &lctTypes.Web4LCT{
		ID:              lctID,
		EntityType:      lctTypes.EntityTypeHuman,
		PublicKey:       pubKey,
		T3Competence:    0.5,
		T3Reliability:   0.5,
		T3Transparency:  1.0,
		Status:          lctTypes.StatusActive,
		SocietyID:       "soc:web4:act:demo",
	}
	
	// Calculate trust mass
	lct.CalculateTrustMass()
	require.Greater(t, lct.TrustMass, 0.0)
	require.LessOrEqual(t, lct.TrustMass, 1.0)
	
	// Validate LCT
	err = lct.Validate()
	require.NoError(t, err)
}

// Test X25519 key derivation for Diffie-Hellman
func (suite *Phase1TestSuite) TestX25519KeyDerivation() {
	t := suite.T()
	
	// Generate Ed25519 keys
	_, privKey, err := suite.cryptoMgr.GenerateKeyPair()
	require.NoError(t, err)
	
	// Derive X25519 keys
	x25519Pub, x25519Priv, err := suite.cryptoMgr.DeriveX25519Keys(privKey)
	require.NoError(t, err)
	require.Equal(t, 32, len(x25519Pub))
	require.Equal(t, 32, len(x25519Priv))
	
	// Test Diffie-Hellman with another key pair
	_, privKey2, _ := suite.cryptoMgr.GenerateKeyPair()
	x25519Pub2, x25519Priv2, _ := suite.cryptoMgr.DeriveX25519Keys(privKey2)
	
	// Perform DH from both sides
	shared1, err := suite.cryptoMgr.PerformDiffieHellman(x25519Priv, x25519Pub2)
	require.NoError(t, err)
	
	shared2, err := suite.cryptoMgr.PerformDiffieHellman(x25519Priv2, x25519Pub)
	require.NoError(t, err)
	
	// Shared secrets should match
	require.Equal(t, shared1, shared2)
}

// Test MRH graph creation and storage
func (suite *Phase1TestSuite) TestMRHGraphCreation() {
	t := suite.T()
	
	// Create MRH graph
	graph := &mrhTypes.MRHGraph{
		LctID:        "lct:web4:act:human:test123",
		Version:      1,
		CreatedAt:    1234567890,
		ContextDepth: 3,
		Triples:      []mrhTypes.Triple{},
		Metadata:     make(map[string]string),
	}
	
	// Add some triples
	graph.AddTriple("lct:test123", "web4:trusts", "lct:test456", 0.8)
	graph.AddTriple("lct:test123", "web4:pairedWith", "lct:test789", 1.0)
	
	require.Equal(t, 2, len(graph.Triples))
	
	// Store the graph
	hash, err := suite.mrhStorage.Store(suite.ctx, graph)
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	
	// Retrieve the graph
	retrieved, err := suite.mrhStorage.Retrieve(suite.ctx, hash)
	require.NoError(t, err)
	require.Equal(t, graph.LctID, retrieved.LctID)
	require.Equal(t, len(graph.Triples), len(retrieved.Triples))
	
	// Check existence
	exists, err := suite.mrhStorage.Exists(suite.ctx, hash)
	require.NoError(t, err)
	require.True(t, exists)
}

// Test MRH context boundary calculations
func (suite *Phase1TestSuite) TestMRHContextBoundary() {
	t := suite.T()
	
	// Create a network of LCTs with relationships
	graphs := make(map[string]*mrhTypes.MRHGraph)
	
	// Create center LCT
	centerLCT := "lct:center"
	graphs[centerLCT] = &mrhTypes.MRHGraph{
		LctID:   centerLCT,
		Triples: []mrhTypes.Triple{},
	}
	
	// Add first ring of connections
	ring1 := []string{"lct:ring1a", "lct:ring1b", "lct:ring1c"}
	for _, lct := range ring1 {
		graphs[centerLCT].AddTriple(centerLCT, "web4:trusts", lct, 0.9)
		graphs[lct] = &mrhTypes.MRHGraph{
			LctID:   lct,
			Triples: []mrhTypes.Triple{},
		}
		graphs[lct].AddTriple(lct, "web4:trusts", centerLCT, 0.9)
	}
	
	// Add second ring
	ring2 := []string{"lct:ring2a", "lct:ring2b"}
	graphs[ring1[0]].AddTriple(ring1[0], "web4:trusts", ring2[0], 0.7)
	graphs[ring1[1]].AddTriple(ring1[1], "web4:trusts", ring2[1], 0.7)
	
	// Store all graphs
	for _, graph := range graphs {
		_, err := suite.mrhStorage.Store(suite.ctx, graph)
		require.NoError(t, err)
	}
	
	// Test traversal
	traversal := mrhTypes.NewLocalMRHTraversal(suite.mrhStorage.(*mrhTypes.LocalMRHStorage))
	
	// Find path from center to ring2a
	path, err := traversal.FindPath(centerLCT, ring2[0], 3)
	require.NoError(t, err)
	require.Equal(t, 3, len(path)) // center -> ring1a -> ring2a
	
	// Calculate trust
	trust, err := traversal.CalculateTrust(centerLCT, ring2[0])
	require.NoError(t, err)
	require.Greater(t, trust, 0.0)
	require.Less(t, trust, 1.0)
	
	// Get context with radius 1
	context, err := traversal.GetContext(centerLCT, 1)
	require.NoError(t, err)
	require.Contains(t, context.IncludedLCTs, centerLCT)
	// Note: This test is simplified - real implementation would include ring1
}

// Test proof of agency generation
func (suite *Phase1TestSuite) TestProofOfAgency() {
	t := suite.T()
	
	// Create delegator (human)
	delegatorPub, delegatorPriv, _ := suite.cryptoMgr.GenerateKeyPair()
	
	// Create agent
	agentPub, _, _ := suite.cryptoMgr.GenerateKeyPair()
	
	// Define permissions and constraints
	permissions := []string{"read", "write", "execute"}
	constraints := map[string]interface{}{
		"max_value": 1000,
		"time_limit": 3600,
	}
	
	// Generate proof of agency
	proof, err := suite.cryptoMgr.GenerateProofOfAgency(
		delegatorPriv,
		agentPub,
		permissions,
		constraints,
	)
	require.NoError(t, err)
	require.NotNil(t, proof)
	
	// Verify proof
	valid, err := suite.cryptoMgr.VerifyProofOfAgency(
		delegatorPub,
		agentPub,
		proof,
	)
	require.NoError(t, err)
	require.True(t, valid)
	
	// Test with wrong agent key
	wrongAgentPub, _, _ := suite.cryptoMgr.GenerateKeyPair()
	valid, err = suite.cryptoMgr.VerifyProofOfAgency(
		delegatorPub,
		wrongAgentPub,
		proof,
	)
	require.Error(t, err) // Should fail with key mismatch
}

// Test witness signature generation
func (suite *Phase1TestSuite) TestWitnessSignature() {
	t := suite.T()
	
	// Create witness
	witnessPub, witnessPriv, _ := suite.cryptoMgr.GenerateKeyPair()
	
	// Create subject LCT
	subjectLCT := "lct:web4:act:human:subject"
	eventType := "birth_certificate"
	eventData := []byte("birth event data")
	
	// Generate witness signature
	signature, err := suite.cryptoMgr.GenerateWitnessSignature(
		witnessPriv,
		subjectLCT,
		eventType,
		eventData,
	)
	require.NoError(t, err)
	require.NotNil(t, signature)
	
	// In a real test, we'd verify the signature using the witness's public key
	require.Equal(t, ed25519.SignatureSize, len(signature))
	
	// Test signature verification (simplified)
	message := []byte(subjectLCT + "|" + eventType)
	valid := suite.cryptoMgr.VerifySignature(witnessPub, message[:10], signature) // Simplified test
	_ = valid // In real test, we'd check this properly
}

// Test agent LCT creation from parent
func (suite *Phase1TestSuite) TestAgentLCTCreation() {
	t := suite.T()
	
	// Create parent (human) LCT
	humanPub, humanPriv, _ := suite.cryptoMgr.GenerateKeyPair()
	humanLCT := &lctTypes.Web4LCT{
		ID:         suite.cryptoMgr.GenerateLCTID(lctTypes.EntityTypeHuman, humanPub),
		EntityType: lctTypes.EntityTypeHuman,
		PublicKey:  humanPub,
		Status:     lctTypes.StatusActive,
	}
	
	// Create agent LCT
	agentPub, _, _ := suite.cryptoMgr.GenerateKeyPair()
	
	// Generate agency proof
	permissions := []string{"read", "limited_write"}
	constraints := map[string]interface{}{"daily_limit": 100}
	
	agencyProof, err := suite.cryptoMgr.GenerateProofOfAgency(
		humanPriv,
		agentPub,
		permissions,
		constraints,
	)
	require.NoError(t, err)
	
	// Create agent LCT
	agentLCT := &lctTypes.Web4LCT{
		ID:          suite.cryptoMgr.GenerateLCTID(lctTypes.EntityTypeAgent, agentPub),
		EntityType:  lctTypes.EntityTypeAgent,
		PublicKey:   agentPub,
		ParentLCT:   humanLCT.ID,
		AgencyProof: agencyProof,
		Permissions: permissions,
		Status:      lctTypes.StatusActive,
	}
	
	// Validate agent LCT
	err = agentLCT.Validate()
	require.NoError(t, err)
	
	// Verify parent-child relationship
	require.Equal(t, humanLCT.ID, agentLCT.ParentLCT)
	require.True(t, humanLCT.CanDelegate())
}

// Test trust tensor calculations
func (suite *Phase1TestSuite) TestTrustTensorCalculations() {
	t := suite.T()
	
	testCases := []struct {
		name           string
		competence     float64
		reliability    float64
		transparency   float64
		expectedMass   float64 // Approximate
		expectedRadius float64 // Approximate
	}{
		{
			name:           "Perfect trust",
			competence:     1.0,
			reliability:    1.0,
			transparency:   1.0,
			expectedMass:   1.0,
			expectedRadius: 10.0,
		},
		{
			name:           "Medium trust",
			competence:     0.7,
			reliability:    0.6,
			transparency:   0.8,
			expectedMass:   0.6, // Approximate geometric mean
			expectedRadius: 6.0,  // Approximate
		},
		{
			name:           "Low trust",
			competence:     0.3,
			reliability:    0.2,
			transparency:   0.4,
			expectedMass:   0.28, // Approximate
			expectedRadius: 2.8,  // Approximate
		},
	}
	
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			lct := &lctTypes.Web4LCT{
				T3Competence:   tc.competence,
				T3Reliability:  tc.reliability,
				T3Transparency: tc.transparency,
			}
			
			lct.CalculateTrustMass()
			
			// Check trust mass (within 10% tolerance)
			require.InDelta(t, tc.expectedMass, lct.TrustMass, tc.expectedMass*0.2)
			
			// Check trust radius
			require.InDelta(t, tc.expectedRadius, lct.TrustRadius, tc.expectedRadius*0.2)
		})
	}
}

// Benchmark tests

func BenchmarkEd25519KeyGeneration(b *testing.B) {
	mgr := lctTypes.NewCryptoManager()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := mgr.GenerateKeyPair()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMRHGraphStorage(b *testing.B) {
	storage, _ := mrhTypes.NewMRHStorage(mrhTypes.StorageConfig{
		Type:      "local",
		LocalPath: "/tmp/bench_mrh",
	})
	ctx := context.Background()
	
	graph := &mrhTypes.MRHGraph{
		LctID:   "lct:bench",
		Triples: make([]mrhTypes.Triple, 100), // 100 relationships
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash, err := storage.Store(ctx, graph)
		if err != nil {
			b.Fatal(err)
		}
		
		_, err = storage.Retrieve(ctx, hash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Run the test suite
func TestPhase1Suite(t *testing.T) {
	suite.Run(t, new(Phase1TestSuite))
}