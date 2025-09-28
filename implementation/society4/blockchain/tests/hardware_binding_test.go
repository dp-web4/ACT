package tests

import (
    "encoding/json"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    hardware "society4chain/x/hardware/types"
    lctTypes "society4chain/x/lctmanager/types"
)

func TestHardwareBinding(t *testing.T) {
    t.Run("Extract Hardware", func(t *testing.T) {
        hw, err := hardware.ExtractCurrentHardware()
        require.NoError(t, err)

        assert.Equal(t, "wsl2", hw.Platform)
        assert.NotEmpty(t, hw.HardwareHash)
        assert.NotEmpty(t, hw.Components.WindowsUUID)
        assert.NotEmpty(t, hw.Components.CPUInfo)
        assert.Greater(t, hw.Components.MemoryKB, int64(0))
    })

    t.Run("Create Self-LCT", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        assert.NotEmpty(t, selfLCT.ID)
        assert.Equal(t, "Society4-Self-LCT", selfLCT.Name)
        assert.NotNil(t, selfLCT.HardwareBinding)
        assert.NotNil(t, selfLCT.PublicKey)
        assert.NotNil(t, selfLCT.PrivateKey)
    })

    t.Run("Verify Hardware Match", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        // Should verify successfully on same machine
        err = selfLCT.VerifyHardware()
        assert.NoError(t, err)
    })

    t.Run("Sign and Verify Block", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        blockHash := []byte("test-block-hash-12345")

        // Sign block with hardware verification
        signature, err := selfLCT.SignBlock(blockHash)
        require.NoError(t, err)
        assert.NotEmpty(t, signature)

        // Verify block signature
        err = selfLCT.VerifyBlockSignature(blockHash, signature)
        assert.NoError(t, err)
    })

    t.Run("Create Role LCTs", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        // Create queen role LCT
        queenLCT, err := selfLCT.CreateRoleLCT("Treasury-Queen", "queen")
        require.NoError(t, err)

        assert.NotEmpty(t, queenLCT.ID)
        assert.Equal(t, "Treasury-Queen", queenLCT.RoleName)
        assert.Equal(t, "queen", queenLCT.RoleType)
        assert.Equal(t, selfLCT.ID, queenLCT.ParentLCT)
        assert.NotEmpty(t, queenLCT.ParentSignature)
        assert.Contains(t, selfLCT.RoleChildren, queenLCT.ID)
    })

    t.Run("Hardware Binding Persistence", func(t *testing.T) {
        // Create and save self-LCT
        selfLCT1, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        // Serialize
        data, err := json.Marshal(selfLCT1)
        require.NoError(t, err)

        // Deserialize
        var selfLCT2 lctTypes.SelfLCT
        err = json.Unmarshal(data, &selfLCT2)
        require.NoError(t, err)

        // Verify hardware still matches
        err = selfLCT2.VerifyHardware()
        assert.NoError(t, err)

        // Hashes should match
        assert.Equal(t, selfLCT1.HardwareBinding.HardwareHash,
                     selfLCT2.HardwareBinding.HardwareHash)
    })

    t.Run("Simulate Hardware Change", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        // Modify hardware binding to simulate different machine
        selfLCT.HardwareBinding.Components.WindowsUUID = "DIFFERENT-UUID"

        // Verification should fail
        err = selfLCT.VerifyHardware()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "Windows UUID mismatch")
    })

    t.Run("Block Signature Failure on Wrong Hardware", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        // Modify hardware to simulate wrong machine
        selfLCT.HardwareBinding.Components.WindowsUUID = "WRONG-UUID"

        blockHash := []byte("test-block")

        // Signing should fail due to hardware mismatch
        _, err = selfLCT.SignBlock(blockHash)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "hardware verification failed")
    })
}

func TestConsensusIntegration(t *testing.T) {
    t.Run("Validator Creation", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        validator := NewHardwareValidator(selfLCT)
        assert.NotNil(t, validator)

        // Check hardware binding
        err = validator.CheckHardwareBinding()
        assert.NoError(t, err)
    })

    t.Run("Block Validation", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        validator := NewHardwareValidator(selfLCT)

        // Create mock block
        blockHash := []byte("mock-block-hash")
        signature, err := selfLCT.SignBlock(blockHash)
        require.NoError(t, err)

        // Simulate block validation
        mockReq := createMockBeginBlockRequest(blockHash, signature)

        // Validation should pass
        err = validator.ValidateBlock(nil, mockReq)
        assert.NoError(t, err)
    })

    t.Run("Block Validation Failure", func(t *testing.T) {
        selfLCT, err := lctTypes.CreateSelfLCT()
        require.NoError(t, err)

        validator := NewHardwareValidator(selfLCT)

        // Create block with wrong signature
        blockHash := []byte("mock-block-hash")
        wrongSignature := []byte("wrong-signature")

        mockReq := createMockBeginBlockRequest(blockHash, wrongSignature)

        // Validation should fail
        err = validator.ValidateBlock(nil, mockReq)
        assert.Error(t, err)
    })
}

func BenchmarkHardwareVerification(b *testing.B) {
    selfLCT, err := lctTypes.CreateSelfLCT()
    if err != nil {
        b.Fatal(err)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = selfLCT.VerifyHardware()
    }
}

func BenchmarkBlockSigning(b *testing.B) {
    selfLCT, err := lctTypes.CreateSelfLCT()
    if err != nil {
        b.Fatal(err)
    }

    blockHash := []byte("test-block-hash")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = selfLCT.SignBlock(blockHash)
    }
}