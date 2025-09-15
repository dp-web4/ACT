package keeper

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func TestGenerateKeyShare(t *testing.T) {
	// Test key share generation
	keyShare, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("GenerateKeyShare failed: %v", err)
	}

	// Verify it's 32 bytes
	if len(keyShare) != 32 {
		t.Errorf("Expected 32 bytes, got %d", len(keyShare))
	}

	// Verify it's not all zeros
	zeroKey := [32]byte{}
	if bytes.Equal(keyShare[:], zeroKey[:]) {
		t.Error("Generated key share is all zeros")
	}

	// Test multiple generations produce different keys
	keyShare2, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Second GenerateKeyShare failed: %v", err)
	}

	if bytes.Equal(keyShare[:], keyShare2[:]) {
		t.Error("Two generated key shares are identical (extremely unlikely)")
	}
}

func TestGenerateEd25519KeyPair(t *testing.T) {
	// Test Ed25519 key pair generation
	publicKey, privateKey, err := GenerateEd25519KeyPair()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPair failed: %v", err)
	}

	// Verify key sizes
	if len(publicKey) != ed25519.PublicKeySize {
		t.Errorf("Expected public key size %d, got %d", ed25519.PublicKeySize, len(publicKey))
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		t.Errorf("Expected private key size %d, got %d", ed25519.PrivateKeySize, len(privateKey))
	}

	// Test that the public key matches the private key
	derivedPublicKey := privateKey.Public().(ed25519.PublicKey)
	if !bytes.Equal(publicKey, derivedPublicKey) {
		t.Error("Generated public key doesn't match private key's public key")
	}

	// Test signing and verification
	message := []byte("test message")
	signature, err := SignMessage(privateKey, message)
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}

	if !VerifyMessageSignature(publicKey, message, signature) {
		t.Error("Signature verification failed")
	}

	// Test that wrong message fails verification
	wrongMessage := []byte("wrong message")
	if VerifyMessageSignature(publicKey, wrongMessage, signature) {
		t.Error("Signature verification should have failed for wrong message")
	}
}

func TestGenerateCurve25519KeyPair(t *testing.T) {
	// Test Curve25519 key pair generation
	publicKey, privateKey, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("GenerateCurve25519KeyPair failed: %v", err)
	}

	// Verify key sizes
	if len(publicKey) != 32 {
		t.Errorf("Expected public key size 32, got %d", len(publicKey))
	}

	if len(privateKey) != 32 {
		t.Errorf("Expected private key size 32, got %d", len(privateKey))
	}

	// Verify keys are not all zeros
	zeroKey := [32]byte{}
	if bytes.Equal(publicKey[:], zeroKey[:]) {
		t.Error("Generated public key is all zeros")
	}

	if bytes.Equal(privateKey[:], zeroKey[:]) {
		t.Error("Generated private key is all zeros")
	}

	// Test that different generations produce different keys
	publicKey2, privateKey2, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Second GenerateCurve25519KeyPair failed: %v", err)
	}

	if bytes.Equal(publicKey[:], publicKey2[:]) {
		t.Error("Two generated public keys are identical (extremely unlikely)")
	}

	if bytes.Equal(privateKey[:], privateKey2[:]) {
		t.Error("Two generated private keys are identical (extremely unlikely)")
	}
}

func TestDeriveSharedSecret(t *testing.T) {
	// Generate two key pairs
	publicKeyA, privateKeyA, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair A: %v", err)
	}

	publicKeyB, privateKeyB, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair B: %v", err)
	}

	// Derive shared secrets
	sharedSecretAB, err := DeriveSharedSecret(privateKeyA, publicKeyB)
	if err != nil {
		t.Fatalf("Failed to derive shared secret AB: %v", err)
	}

	sharedSecretBA, err := DeriveSharedSecret(privateKeyB, publicKeyA)
	if err != nil {
		t.Fatalf("Failed to derive shared secret BA: %v", err)
	}

	// Verify both parties derive the same shared secret
	if !bytes.Equal(sharedSecretAB[:], sharedSecretBA[:]) {
		t.Error("Shared secrets don't match between parties")
	}

	// Verify shared secret is not all zeros
	zeroKey := [32]byte{}
	if bytes.Equal(sharedSecretAB[:], zeroKey[:]) {
		t.Error("Derived shared secret is all zeros")
	}

	// Test that different key pairs produce different shared secrets
	publicKeyC, _, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair C: %v", err)
	}

	sharedSecretAC, err := DeriveSharedSecret(privateKeyA, publicKeyC)
	if err != nil {
		t.Fatalf("Failed to derive shared secret AC: %v", err)
	}

	if bytes.Equal(sharedSecretAB[:], sharedSecretAC[:]) {
		t.Error("Different key pairs produced identical shared secrets (extremely unlikely)")
	}
}

func TestCombineKeyShares(t *testing.T) {
	// Generate test key shares and shared secret
	shareA, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate share A: %v", err)
	}

	shareB, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate share B: %v", err)
	}

	sharedSecret, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate shared secret: %v", err)
	}

	// Combine key shares
	combinedKey := CombineKeyShares(shareA, shareB, sharedSecret)

	// Verify combined key is 32 bytes
	if len(combinedKey) != 32 {
		t.Errorf("Expected combined key size 32, got %d", len(combinedKey))
	}

	// Verify combined key is not all zeros
	zeroKey := [32]byte{}
	if bytes.Equal(combinedKey[:], zeroKey[:]) {
		t.Error("Combined key is all zeros")
	}

	// Test determinism - same inputs should produce same output
	combinedKey2 := CombineKeyShares(shareA, shareB, sharedSecret)
	if !bytes.Equal(combinedKey[:], combinedKey2[:]) {
		t.Error("CombineKeyShares is not deterministic")
	}

	// Test that different inputs produce different outputs
	shareC, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate share C: %v", err)
	}

	combinedKey3 := CombineKeyShares(shareC, shareB, sharedSecret)
	if bytes.Equal(combinedKey[:], combinedKey3[:]) {
		t.Error("Different inputs produced identical combined keys (extremely unlikely)")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Generate a key
	key, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Test data
	plaintext := []byte("Hello, World! This is a test message for encryption.")

	// Encrypt
	ciphertext, err := EncryptWithKey(key, plaintext)
	if err != nil {
		t.Fatalf("EncryptWithKey failed: %v", err)
	}

	// Verify ciphertext is different from plaintext
	if bytes.Equal(ciphertext, plaintext) {
		t.Error("Ciphertext is identical to plaintext")
	}

	// Decrypt
	decrypted, err := DecryptWithKey(key, ciphertext)
	if err != nil {
		t.Fatalf("DecryptWithKey failed: %v", err)
	}

	// Verify decrypted matches original
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Decrypted text doesn't match original plaintext")
	}

	// Test that wrong key fails
	wrongKey, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate wrong key: %v", err)
	}

	_, err = DecryptWithKey(wrongKey, ciphertext)
	if err == nil {
		t.Error("Decryption with wrong key should have failed")
	}

	// Test that tampered ciphertext fails
	tamperedCiphertext := make([]byte, len(ciphertext))
	copy(tamperedCiphertext, ciphertext)
	if len(tamperedCiphertext) > 0 {
		tamperedCiphertext[0] ^= 1 // Flip one bit
	}

	_, err = DecryptWithKey(key, tamperedCiphertext)
	if err == nil {
		t.Error("Decryption of tampered ciphertext should have failed")
	}
}

func TestZeroKey(t *testing.T) {
	// Generate a key
	key, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Verify it's not all zeros initially
	zeroKey := [32]byte{}
	if bytes.Equal(key[:], zeroKey[:]) {
		t.Error("Generated key is all zeros")
	}

	// Zero the key
	ZeroKey(&key)

	// Verify it's now all zeros
	if !bytes.Equal(key[:], zeroKey[:]) {
		t.Error("Key was not properly zeroed")
	}
}

func TestZeroEd25519Key(t *testing.T) {
	// Generate an Ed25519 key pair
	_, privateKey, err := GenerateEd25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	// Verify it's not all zeros initially
	zeroKey := make([]byte, len(privateKey))
	if bytes.Equal(privateKey, zeroKey) {
		t.Error("Generated private key is all zeros")
	}

	// Zero the key
	ZeroEd25519Key(&privateKey)

	// Verify it's now all zeros
	if !bytes.Equal(privateKey, zeroKey) {
		t.Error("Private key was not properly zeroed")
	}
}

func TestCompleteKeyExchangeFlow(t *testing.T) {
	// This test simulates a complete key exchange between two parties

	// Party A generates keys
	publicKeyA, privateKeyA, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Party A failed to generate keys: %v", err)
	}

	shareA, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Party A failed to generate share: %v", err)
	}

	// Party B generates keys
	publicKeyB, privateKeyB, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Fatalf("Party B failed to generate keys: %v", err)
	}

	shareB, err := GenerateKeyShare()
	if err != nil {
		t.Fatalf("Party B failed to generate share: %v", err)
	}

	// Both parties derive shared secret
	sharedSecretA, err := DeriveSharedSecret(privateKeyA, publicKeyB)
	if err != nil {
		t.Fatalf("Party A failed to derive shared secret: %v", err)
	}

	sharedSecretB, err := DeriveSharedSecret(privateKeyB, publicKeyA)
	if err != nil {
		t.Fatalf("Party B failed to derive shared secret: %v", err)
	}

	// Verify shared secrets match
	if !bytes.Equal(sharedSecretA[:], sharedSecretB[:]) {
		t.Error("Shared secrets don't match between parties")
	}

	// Both parties combine key shares
	combinedKeyA := CombineKeyShares(shareA, shareB, sharedSecretA)
	combinedKeyB := CombineKeyShares(shareA, shareB, sharedSecretB)

	// Verify combined keys match
	if !bytes.Equal(combinedKeyA[:], combinedKeyB[:]) {
		t.Error("Combined keys don't match between parties")
	}

	// Test encryption/decryption with the combined key
	message := []byte("Secret message from Party A to Party B")

	// Party A encrypts
	ciphertext, err := EncryptWithKey(combinedKeyA, message)
	if err != nil {
		t.Fatalf("Party A failed to encrypt: %v", err)
	}

	// Party B decrypts
	decrypted, err := DecryptWithKey(combinedKeyB, ciphertext)
	if err != nil {
		t.Fatalf("Party B failed to decrypt: %v", err)
	}

	// Verify message integrity
	if !bytes.Equal(decrypted, message) {
		t.Error("Decrypted message doesn't match original")
	}

	// Clean up sensitive data
	ZeroKey(&privateKeyA)
	ZeroKey(&privateKeyB)
	ZeroKey(&shareA)
	ZeroKey(&shareB)
	ZeroKey(&sharedSecretA)
	ZeroKey(&sharedSecretB)
	ZeroKey(&combinedKeyA)
	ZeroKey(&combinedKeyB)
}
