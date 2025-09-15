package types

import (
	"bytes"
	"testing"
)

func TestNewSplitKey(t *testing.T) {
	lctId := "test-lct-123"
	splitKey := NewSplitKey(lctId)

	if splitKey.LctId != lctId {
		t.Errorf("Expected LCT ID %s, got %s", lctId, splitKey.LctId)
	}

	if splitKey.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", splitKey.Status)
	}

	if splitKey.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}

	if splitKey.ActivatedAt != 0 {
		t.Error("ActivatedAt should be 0 for new split key")
	}
}

func TestSetDeviceKeyHalf(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Create a test key half
	var keyHalf [32]byte
	for i := range keyHalf {
		keyHalf[i] = byte(i)
	}

	splitKey.SetDeviceKeyHalf(keyHalf)

	if len(splitKey.DeviceKeyHalf) != 32 {
		t.Errorf("Expected device key half length 32, got %d", len(splitKey.DeviceKeyHalf))
	}

	if !bytes.Equal(splitKey.DeviceKeyHalf, keyHalf[:]) {
		t.Error("Device key half not set correctly")
	}
}

func TestSetLctKeyHalf(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Create a test key half
	var keyHalf [32]byte
	for i := range keyHalf {
		keyHalf[i] = byte(i + 32)
	}

	splitKey.SetLctKeyHalf(keyHalf)

	if len(splitKey.LctKeyHalf) != 32 {
		t.Errorf("Expected LCT key half length 32, got %d", len(splitKey.LctKeyHalf))
	}

	if !bytes.Equal(splitKey.LctKeyHalf, keyHalf[:]) {
		t.Error("LCT key half not set correctly")
	}
}

func TestSetPublicKeys(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	publicKeyA := []byte{1, 2, 3, 4, 5}
	publicKeyB := []byte{6, 7, 8, 9, 10}

	splitKey.SetPublicKeyA(publicKeyA)
	splitKey.SetPublicKeyB(publicKeyB)

	if !bytes.Equal(splitKey.PublicKeyA, publicKeyA) {
		t.Error("Public key A not set correctly")
	}

	if !bytes.Equal(splitKey.PublicKeyB, publicKeyB) {
		t.Error("Public key B not set correctly")
	}
}

func TestSetCombinedKey(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Create a test combined key
	var combinedKey [32]byte
	for i := range combinedKey {
		combinedKey[i] = byte(i + 64)
	}

	splitKey.SetCombinedKey(combinedKey)

	if len(splitKey.CombinedKey) != 32 {
		t.Errorf("Expected combined key length 32, got %d", len(splitKey.CombinedKey))
	}

	if !bytes.Equal(splitKey.CombinedKey, combinedKey[:]) {
		t.Error("Combined key not set correctly")
	}
}

func TestActivate(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	splitKey.Activate()

	if splitKey.Status != "active" {
		t.Errorf("Expected status 'active', got %s", splitKey.Status)
	}

	if splitKey.ActivatedAt == 0 {
		t.Error("ActivatedAt should be set to current time")
	}
}

func TestRevoke(t *testing.T) {
	splitKey := NewSplitKey("test-lct")
	splitKey.Activate()

	splitKey.Revoke()

	if splitKey.Status != "revoked" {
		t.Errorf("Expected status 'revoked', got %s", splitKey.Status)
	}
}

func TestStatusChecks(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Test pending state
	if !splitKey.IsPending() {
		t.Error("New split key should be pending")
	}

	if splitKey.IsActive() {
		t.Error("New split key should not be active")
	}

	if splitKey.IsRevoked() {
		t.Error("New split key should not be revoked")
	}

	// Test active state
	splitKey.Activate()
	if !splitKey.IsActive() {
		t.Error("Activated split key should be active")
	}

	if splitKey.IsPending() {
		t.Error("Activated split key should not be pending")
	}

	if splitKey.IsRevoked() {
		t.Error("Activated split key should not be revoked")
	}

	// Test revoked state
	splitKey.Revoke()
	if !splitKey.IsRevoked() {
		t.Error("Revoked split key should be revoked")
	}

	if splitKey.IsActive() {
		t.Error("Revoked split key should not be active")
	}

	if splitKey.IsPending() {
		t.Error("Revoked split key should not be pending")
	}
}

func TestValidate(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Test empty LCT ID
	splitKey.LctId = ""
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for empty LCT ID")
	}

	// Test missing device key half
	splitKey.LctId = "test-lct"
	splitKey.DeviceKeyHalf = []byte{1, 2, 3} // Wrong length
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for wrong device key half length")
	}

	// Test missing LCT key half
	splitKey.DeviceKeyHalf = make([]byte, 32)
	splitKey.LctKeyHalf = []byte{1, 2, 3} // Wrong length
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for wrong LCT key half length")
	}

	// Test missing public key A
	splitKey.LctKeyHalf = make([]byte, 32)
	splitKey.PublicKeyA = nil
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for missing public key A")
	}

	// Test missing public key B
	splitKey.PublicKeyA = []byte{1, 2, 3, 4, 5}
	splitKey.PublicKeyB = nil
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for missing public key B")
	}

	// Test missing combined key when active
	splitKey.PublicKeyB = []byte{6, 7, 8, 9, 10}
	splitKey.Activate()
	splitKey.CombinedKey = []byte{1, 2, 3} // Wrong length
	if err := splitKey.Validate(); err == nil {
		t.Error("Validation should fail for wrong combined key length when active")
	}

	// Test valid split key
	splitKey.CombinedKey = make([]byte, 32)
	if err := splitKey.Validate(); err != nil {
		t.Errorf("Validation should pass for valid split key: %v", err)
	}
}

func TestGetKeyHalfAsArray(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Create test key halves
	var deviceKeyHalf [32]byte
	var lctKeyHalf [32]byte
	var combinedKey [32]byte

	for i := range deviceKeyHalf {
		deviceKeyHalf[i] = byte(i)
		lctKeyHalf[i] = byte(i + 32)
		combinedKey[i] = byte(i + 64)
	}

	splitKey.SetDeviceKeyHalf(deviceKeyHalf)
	splitKey.SetLctKeyHalf(lctKeyHalf)
	splitKey.SetCombinedKey(combinedKey)

	// Test GetDeviceKeyHalfAsArray
	retrievedDeviceKeyHalf, err := splitKey.GetDeviceKeyHalfAsArray()
	if err != nil {
		t.Errorf("GetDeviceKeyHalfAsArray failed: %v", err)
	}
	if !bytes.Equal(retrievedDeviceKeyHalf[:], deviceKeyHalf[:]) {
		t.Error("Retrieved device key half doesn't match original")
	}

	// Test GetLctKeyHalfAsArray
	retrievedLctKeyHalf, err := splitKey.GetLctKeyHalfAsArray()
	if err != nil {
		t.Errorf("GetLctKeyHalfAsArray failed: %v", err)
	}
	if !bytes.Equal(retrievedLctKeyHalf[:], lctKeyHalf[:]) {
		t.Error("Retrieved LCT key half doesn't match original")
	}

	// Test GetCombinedKeyAsArray
	retrievedCombinedKey, err := splitKey.GetCombinedKeyAsArray()
	if err != nil {
		t.Errorf("GetCombinedKeyAsArray failed: %v", err)
	}
	if !bytes.Equal(retrievedCombinedKey[:], combinedKey[:]) {
		t.Error("Retrieved combined key doesn't match original")
	}
}

func TestGetKeyHalfAsArrayErrors(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Test with wrong length device key half
	splitKey.DeviceKeyHalf = []byte{1, 2, 3}
	_, err := splitKey.GetDeviceKeyHalfAsArray()
	if err == nil {
		t.Error("GetDeviceKeyHalfAsArray should fail for wrong length")
	}

	// Test with wrong length LCT key half
	splitKey.DeviceKeyHalf = make([]byte, 32)
	splitKey.LctKeyHalf = []byte{1, 2, 3}
	_, err = splitKey.GetLctKeyHalfAsArray()
	if err == nil {
		t.Error("GetLctKeyHalfAsArray should fail for wrong length")
	}

	// Test with wrong length combined key
	splitKey.LctKeyHalf = make([]byte, 32)
	splitKey.CombinedKey = []byte{1, 2, 3}
	_, err = splitKey.GetCombinedKeyAsArray()
	if err == nil {
		t.Error("GetCombinedKeyAsArray should fail for wrong length")
	}
}

func TestZeroSensitiveData(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Set some data
	var deviceKeyHalf [32]byte
	var lctKeyHalf [32]byte
	var combinedKey [32]byte

	for i := range deviceKeyHalf {
		deviceKeyHalf[i] = byte(i)
		lctKeyHalf[i] = byte(i + 32)
		combinedKey[i] = byte(i + 64)
	}

	splitKey.SetDeviceKeyHalf(deviceKeyHalf)
	splitKey.SetLctKeyHalf(lctKeyHalf)
	splitKey.SetCombinedKey(combinedKey)
	splitKey.SetPublicKeyA([]byte{1, 2, 3, 4, 5})
	splitKey.SetPublicKeyB([]byte{6, 7, 8, 9, 10})

	// Verify data is set
	if len(splitKey.DeviceKeyHalf) == 0 {
		t.Error("Device key half should be set before zeroing")
	}

	if len(splitKey.LctKeyHalf) == 0 {
		t.Error("LCT key half should be set before zeroing")
	}

	if len(splitKey.CombinedKey) == 0 {
		t.Error("Combined key should be set before zeroing")
	}

	// Zero sensitive data
	splitKey.ZeroSensitiveData()

	// Verify sensitive data is zeroed
	if splitKey.DeviceKeyHalf != nil {
		t.Error("Device key half should be nil after zeroing")
	}

	if splitKey.LctKeyHalf != nil {
		t.Error("LCT key half should be nil after zeroing")
	}

	if splitKey.CombinedKey != nil {
		t.Error("Combined key should be nil after zeroing")
	}

	// Verify public keys are not zeroed (they're not sensitive)
	if len(splitKey.PublicKeyA) == 0 {
		t.Error("Public key A should not be zeroed")
	}

	if len(splitKey.PublicKeyB) == 0 {
		t.Error("Public key B should not be zeroed")
	}
}

func TestClone(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Set some data
	var deviceKeyHalf [32]byte
	var lctKeyHalf [32]byte
	var combinedKey [32]byte

	for i := range deviceKeyHalf {
		deviceKeyHalf[i] = byte(i)
		lctKeyHalf[i] = byte(i + 32)
		combinedKey[i] = byte(i + 64)
	}

	splitKey.SetDeviceKeyHalf(deviceKeyHalf)
	splitKey.SetLctKeyHalf(lctKeyHalf)
	splitKey.SetCombinedKey(combinedKey)
	splitKey.SetPublicKeyA([]byte{1, 2, 3, 4, 5})
	splitKey.SetPublicKeyB([]byte{6, 7, 8, 9, 10})
	splitKey.Activate()

	// Clone the split key
	clone := splitKey.Clone()

	// Verify clone has same data
	if clone.LctId != splitKey.LctId {
		t.Error("Clone LCT ID doesn't match")
	}

	if clone.Status != splitKey.Status {
		t.Error("Clone status doesn't match")
	}

	if clone.CreatedAt != splitKey.CreatedAt {
		t.Error("Clone CreatedAt doesn't match")
	}

	if clone.ActivatedAt != splitKey.ActivatedAt {
		t.Error("Clone ActivatedAt doesn't match")
	}

	if !bytes.Equal(clone.DeviceKeyHalf, splitKey.DeviceKeyHalf) {
		t.Error("Clone device key half doesn't match")
	}

	if !bytes.Equal(clone.LctKeyHalf, splitKey.LctKeyHalf) {
		t.Error("Clone LCT key half doesn't match")
	}

	if !bytes.Equal(clone.PublicKeyA, splitKey.PublicKeyA) {
		t.Error("Clone public key A doesn't match")
	}

	if !bytes.Equal(clone.PublicKeyB, splitKey.PublicKeyB) {
		t.Error("Clone public key B doesn't match")
	}

	if !bytes.Equal(clone.CombinedKey, splitKey.CombinedKey) {
		t.Error("Clone combined key doesn't match")
	}

	// Verify clone is independent (modifying original doesn't affect clone)
	splitKey.LctId = "modified-lct"
	if clone.LctId == "modified-lct" {
		t.Error("Clone should be independent of original")
	}
}

func TestToHexString(t *testing.T) {
	splitKey := NewSplitKey("test-lct")

	// Set some data
	var deviceKeyHalf [32]byte
	var lctKeyHalf [32]byte
	var combinedKey [32]byte

	for i := range deviceKeyHalf {
		deviceKeyHalf[i] = byte(i)
		lctKeyHalf[i] = byte(i + 32)
		combinedKey[i] = byte(i + 64)
	}

	splitKey.SetDeviceKeyHalf(deviceKeyHalf)
	splitKey.SetLctKeyHalf(lctKeyHalf)
	splitKey.SetCombinedKey(combinedKey)
	splitKey.SetPublicKeyA([]byte{1, 2, 3, 4, 5})
	splitKey.SetPublicKeyB([]byte{6, 7, 8, 9, 10})

	hexString := splitKey.ToHexString()

	// Verify hex string contains expected data
	if len(hexString) == 0 {
		t.Error("Hex string should not be empty")
	}

	// Verify it contains the LCT ID
	if !bytes.Contains([]byte(hexString), []byte("test-lct")) {
		t.Error("Hex string should contain LCT ID")
	}

	// Verify it contains status
	if !bytes.Contains([]byte(hexString), []byte("pending")) {
		t.Error("Hex string should contain status")
	}
}
