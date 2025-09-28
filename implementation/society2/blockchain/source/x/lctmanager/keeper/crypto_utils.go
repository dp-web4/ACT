package keeper

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

// GenerateKeyShare generates a 32-byte random key share
func GenerateKeyShare() ([32]byte, error) {
	var keyShare [32]byte
	_, err := rand.Read(keyShare[:])
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to generate key share: %w", err)
	}
	return keyShare, nil
}

// GenerateEd25519KeyPair generates Ed25519 signing keys
func GenerateEd25519KeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}
	return publicKey, privateKey, nil
}

// GenerateCurve25519KeyPair generates Curve25519 ECDH keys
func GenerateCurve25519KeyPair() ([32]byte, [32]byte, error) {
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return [32]byte{}, [32]byte{}, fmt.Errorf("failed to generate Curve25519 private key: %w", err)
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	return publicKey, privateKey, nil
}

// DeriveSharedSecret performs Curve25519 ECDH key exchange
func DeriveSharedSecret(privateKey, publicKey [32]byte) ([32]byte, error) {
	var sharedSecret [32]byte
	curve25519.ScalarMult(&sharedSecret, &privateKey, &publicKey)
	return sharedSecret, nil
}

// CombineKeyShares derives the final encryption key from key shares and shared secret
func CombineKeyShares(shareA, shareB, sharedSecret [32]byte) [32]byte {
	// Use SHA-256 to derive the final key from all three inputs
	hash := sha256.New()
	hash.Write(shareA[:])
	hash.Write(shareB[:])
	hash.Write(sharedSecret[:])

	var combinedKey [32]byte
	copy(combinedKey[:], hash.Sum(nil))
	return combinedKey
}

// EncryptWithKey encrypts data using ChaCha20-Poly1305
func EncryptWithKey(key [32]byte, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD cipher: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	// Return nonce + ciphertext
	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result, nonce)
	copy(result[len(nonce):], ciphertext)

	return result, nil
}

// DecryptWithKey decrypts data using ChaCha20-Poly1305
func DecryptWithKey(key [32]byte, ciphertext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD cipher: %w", err)
	}

	if len(ciphertext) < aead.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce := ciphertext[:aead.NonceSize()]
	encryptedData := ciphertext[aead.NonceSize():]

	// Decrypt the data
	plaintext, err := aead.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// SignMessage signs a message using Ed25519
func SignMessage(privateKey ed25519.PrivateKey, message []byte) ([]byte, error) {
	signature := ed25519.Sign(privateKey, message)
	return signature, nil
}

// VerifyMessageSignature verifies an Ed25519 signature
func VerifyMessageSignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

// ZeroKey securely zeroes a key from memory
func ZeroKey(key *[32]byte) {
	for i := range key {
		key[i] = 0
	}
}

// ZeroEd25519Key securely zeroes an Ed25519 private key
func ZeroEd25519Key(key *ed25519.PrivateKey) {
	if key != nil {
		for i := range *key {
			(*key)[i] = 0
		}
	}
}
