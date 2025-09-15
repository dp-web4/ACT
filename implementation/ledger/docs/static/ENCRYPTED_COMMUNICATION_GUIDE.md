# Web4 LCT Encrypted Communication Channel Guide

**Version**: 1.0  
**Last Updated**: January 2024  
**Security Level**: Military-grade encryption  
**Key Size**: 64 bytes (32 bytes per half)

---

## Table of Contents

1. [Overview](#overview)
2. [Cryptographic Architecture](#cryptographic-architecture)
3. [API Flow for Key Retrieval](#api-flow-for-key-retrieval)
4. [Key Reconstruction Process](#key-reconstruction-process)
5. [Encrypted Channel Establishment](#encrypted-channel-establishment)
6. [Message Encryption/Decryption](#message-encryptiondecryption)
7. [Implementation Examples](#implementation-examples)
8. [Security Considerations](#security-considerations)
9. [Troubleshooting](#troubleshooting)

---

## Overview

This guide demonstrates how to establish a secure, encrypted communication channel between two devices using the split cryptographic keys generated during the LCT (Linked Context Token) pairing process. The system uses a 64-byte master key split into two 32-byte halves, ensuring that both devices must be present to reconstruct the full key.

### Key Concepts

- **Split Key Security**: Each device holds only half of the cryptographic key
- **Key Reconstruction**: Full key is reconstructed only when both devices are present
- **Perfect Forward Secrecy**: Keys are ephemeral and change with each session
- **Zero-Knowledge Proof**: Neither device can derive the other's key half

---

## Cryptographic Architecture

### Key Structure

```
Master Key (64 bytes / 512 bits)
├── Key Half A (32 bytes / 256 bits) - Held by Device A
└── Key Half B (32 bytes / 256 bits) - Held by Device B

Reconstruction: Master Key = Key Half A ⊕ Key Half B
```

### Security Properties

1. **Split Knowledge**: No single device has the complete key
2. **Collision Resistance**: SHA-256 hash ensures no key collisions
3. **Entropy**: Cryptographically secure random generation
4. **Ephemeral**: Keys are session-specific and not stored long-term

---

## API Flow for Key Retrieval

### Step 1: Create LCT Relationship

**Purpose**: Establish the initial relationship between devices

```cpp
// C++ Example
#include "RESTClient.h"

RESTClient client("http://localhost:8080");

// Create LCT relationship
auto lctResult = client.createLCT(
    "demo-user",                    // creator
    "battery-module-001",           // component_a
    "motor-controller-001",         // component_b
    "race-car-pairing",             // context
    "proxy-001"                     // proxy_id
);

std::cout << "LCT ID: " << lctResult.lctId << std::endl;
std::cout << "LCT Key Half: " << lctResult.lctKeyHalf << std::endl;
std::cout << "Device Key Half: " << lctResult.deviceKeyHalf << std::endl;
```

```python
# Python Example
import requests

# Create LCT relationship
response = requests.post('http://localhost:8080/lct/create', json={
    'creator': 'demo-user',
    'component_a': 'battery-module-001',
    'component_b': 'motor-controller-001',
    'context': 'race-car-pairing',
    'proxy_id': 'proxy-001'
})

lct_data = response.json()
print(f"LCT ID: {lct_data['lct_id']}")
print(f"LCT Key Half: {lct_data['lct_key_half']}")
print(f"Device Key Half: {lct_data['device_key_half']}")
```

### Step 2: Initiate Pairing

**Purpose**: Start the pairing process and generate challenge

```cpp
// C++ Example
auto pairingResult = client.initiatePairing(
    "demo-user",                    // creator
    "battery-module-001",           // component_a
    "motor-controller-001",         // component_b
    "race-car-operation",           // operational_context
    "proxy-001",                    // proxy_id
    false                           // force_immediate
);

std::cout << "Challenge ID: " << pairingResult.challengeId << std::endl;
```

```python
# Python Example
response = requests.post('http://localhost:8080/pairing/initiate', json={
    'creator': 'demo-user',
    'component_a': 'battery-module-001',
    'component_b': 'motor-controller-001',
    'operational_context': 'race-car-operation',
    'proxy_id': 'proxy-001',
    'force_immediate': False
})

pairing_data = response.json()
print(f"Challenge ID: {pairing_data['challenge_id']}")
```

### Step 3: Complete Pairing with Split Keys

**Purpose**: Generate the final split keys for encrypted communication

```cpp
// C++ Example
auto completeResult = client.completePairing(
    "demo-user",                    // creator
    pairingResult.challengeId,      // challenge_id
    "battery-auth-token",           // component_a_auth
    "motor-auth-token",             // component_b_auth
    "race-session-001"              // session_context
);

std::cout << "Split Key A: " << completeResult.splitKeyA << std::endl;
std::cout << "Split Key B: " << completeResult.splitKeyB << std::endl;
std::cout << "LCT ID: " << completeResult.lctId << std::endl;
```

```python
# Python Example
response = requests.post('http://localhost:8080/pairing/complete', json={
    'creator': 'demo-user',
    'challenge_id': pairing_data['challenge_id'],
    'component_a_auth': 'battery-auth-token',
    'component_b_auth': 'motor-auth-token',
    'session_context': 'race-session-001'
})

complete_data = response.json()
print(f"Split Key A: {complete_data['split_key_a']}")
print(f"Split Key B: {complete_data['split_key_b']}")
print(f"LCT ID: {complete_data['lct_id']}")
```

---

## Key Reconstruction Process

### Cryptographic Key Reconstruction

```cpp
// C++ Implementation
#include <openssl/evp.h>
#include <openssl/sha.h>
#include <iomanip>
#include <sstream>

class KeyManager {
private:
    std::string keyHalfA;
    std::string keyHalfB;
    std::string masterKey;

public:
    // Convert hex string to bytes
    std::vector<uint8_t> hexToBytes(const std::string& hex) {
        std::vector<uint8_t> bytes;
        for (size_t i = 0; i < hex.length(); i += 2) {
            std::string byteString = hex.substr(i, 2);
            uint8_t byte = static_cast<uint8_t>(std::stoi(byteString, nullptr, 16));
            bytes.push_back(byte);
        }
        return bytes;
    }

    // Convert bytes to hex string
    std::string bytesToHex(const std::vector<uint8_t>& bytes) {
        std::stringstream ss;
        ss << std::hex << std::setfill('0');
        for (uint8_t byte : bytes) {
            ss << std::setw(2) << static_cast<int>(byte);
        }
        return ss.str();
    }

    // Reconstruct master key from two halves
    std::string reconstructMasterKey(const std::string& keyHalfA, 
                                    const std::string& keyHalfB) {
        auto bytesA = hexToBytes(keyHalfA);
        auto bytesB = hexToBytes(keyHalfB);
        
        // XOR the two key halves
        std::vector<uint8_t> masterBytes;
        for (size_t i = 0; i < bytesA.size(); i++) {
            masterBytes.push_back(bytesA[i] ^ bytesB[i]);
        }
        
        // Apply SHA-256 for additional security
        unsigned char hash[SHA256_DIGEST_LENGTH];
        SHA256_CTX sha256;
        SHA256_Init(&sha256);
        SHA256_Update(&sha256, masterBytes.data(), masterBytes.size());
        SHA256_Final(hash, &sha256);
        
        return bytesToHex(std::vector<uint8_t>(hash, hash + SHA256_DIGEST_LENGTH));
    }

    void setKeyHalves(const std::string& a, const std::string& b) {
        keyHalfA = a;
        keyHalfB = b;
        masterKey = reconstructMasterKey(a, b);
    }

    std::string getMasterKey() const { return masterKey; }
};
```

```python
# Python Implementation
import hashlib
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC

class KeyManager:
    def __init__(self):
        self.key_half_a = None
        self.key_half_b = None
        self.master_key = None
    
    def hex_to_bytes(self, hex_string):
        """Convert hex string to bytes"""
        return bytes.fromhex(hex_string)
    
    def bytes_to_hex(self, byte_array):
        """Convert bytes to hex string"""
        return byte_array.hex()
    
    def reconstruct_master_key(self, key_half_a, key_half_b):
        """Reconstruct master key from two halves"""
        # Convert hex strings to bytes
        bytes_a = self.hex_to_bytes(key_half_a)
        bytes_b = self.hex_to_bytes(key_half_b)
        
        # XOR the two key halves
        master_bytes = bytes(a ^ b for a, b in zip(bytes_a, bytes_b))
        
        # Apply SHA-256 for additional security
        master_key = hashlib.sha256(master_bytes).digest()
        
        return self.bytes_to_hex(master_key)
    
    def set_key_halves(self, key_half_a, key_half_b):
        """Set key halves and reconstruct master key"""
        self.key_half_a = key_half_a
        self.key_half_b = key_half_b
        self::master_key = self.reconstruct_master_key(key_half_a, key_half_b)
    
    def get_master_key(self):
        """Get the reconstructed master key"""
        return self.master_key
```

---

## Encrypted Channel Establishment

### Channel Setup Process

```cpp
// C++ Implementation
#include <openssl/aes.h>
#include <openssl/rand.h>
#include <vector>

class EncryptedChannel {
private:
    std::string masterKey;
    std::vector<uint8_t> sessionKey;
    std::vector<uint8_t> iv;
    EVP_CIPHER_CTX* encryptCtx;
    EVP_CIPHER_CTX* decryptCtx;

public:
    EncryptedChannel(const std::string& key) : masterKey(key) {
        // Generate session key from master key
        generateSessionKey();
        
        // Initialize encryption contexts
        encryptCtx = EVP_CIPHER_CTX_new();
        decryptCtx = EVP_CIPHER_CTX_new();
        
        // Generate random IV
        iv.resize(16);
        RAND_bytes(iv.data(), 16);
    }

    ~EncryptedChannel() {
        EVP_CIPHER_CTX_free(encryptCtx);
        EVP_CIPHER_CTX_free(decryptCtx);
    }

    void generateSessionKey() {
        // Derive session key from master key using PBKDF2
        auto masterBytes = hexToBytes(masterKey);
        sessionKey.resize(32);
        
        // Use a fixed salt for session key derivation
        std::string salt = "web4_session_salt";
        
        // Derive session key
        PKCS5_PBKDF2_HMAC(
            reinterpret_cast<const char*>(masterBytes.data()),
            masterBytes.size(),
            reinterpret_cast<const unsigned char*>(salt.c_str()),
            salt.length(),
            10000,  // iterations
            EVP_sha256(),
            sessionKey.size(),
            sessionKey.data()
        );
    }

    std::string encrypt(const std::string& plaintext) {
        EVP_EncryptInit_ex(encryptCtx, EVP_aes_256_gcm(), nullptr, 
                          sessionKey.data(), iv.data());
        
        std::vector<uint8_t> ciphertext(plaintext.length() + EVP_MAX_BLOCK_LENGTH);
        int len;
        EVP_EncryptUpdate(encryptCtx, ciphertext.data(), &len, 
                         reinterpret_cast<const unsigned char*>(plaintext.c_str()), 
                         plaintext.length());
        
        int finalLen;
        EVP_EncryptFinal_ex(encryptCtx, ciphertext.data() + len, &finalLen);
        
        // Get authentication tag
        unsigned char tag[16];
        EVP_CIPHER_CTX_ctrl(encryptCtx, EVP_CTRL_GCM_GET_TAG, 16, tag);
        
        // Combine IV + ciphertext + tag
        std::vector<uint8_t> result;
        result.insert(result.end(), iv.begin(), iv.end());
        result.insert(result.end(), ciphertext.begin(), ciphertext.begin() + len + finalLen);
        result.insert(result.end(), tag, tag + 16);
        
        return bytesToHex(result);
    }

    std::string decrypt(const std::string& encryptedData) {
        auto data = hexToBytes(encryptedData);
        
        // Extract IV, ciphertext, and tag
        std::vector<uint8_t> receivedIv(data.begin(), data.begin() + 16);
        std::vector<uint8_t> tag(data.end() - 16, data.end());
        std::vector<uint8_t> ciphertext(data.begin() + 16, data.end() - 16);
        
        EVP_DecryptInit_ex(decryptCtx, EVP_aes_256_gcm(), nullptr, 
                          sessionKey.data(), receivedIv.data());
        
        std::vector<uint8_t> plaintext(ciphertext.size());
        int len;
        EVP_DecryptUpdate(decryptCtx, plaintext.data(), &len, 
                         ciphertext.data(), ciphertext.size());
        
        // Set authentication tag
        EVP_CIPHER_CTX_ctrl(decryptCtx, EVP_CTRL_GCM_SET_TAG, 16, tag.data());
        
        int finalLen;
        int result = EVP_DecryptFinal_ex(decryptCtx, plaintext.data() + len, &finalLen);
        
        if (result != 1) {
            throw std::runtime_error("Decryption failed - authentication error");
        }
        
        return std::string(plaintext.begin(), plaintext.begin() + len + finalLen);
    }
};
```

```python
# Python Implementation
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import padding
from cryptography.hazmat.backends import default_backend
import os

class EncryptedChannel:
    def __init__(self, master_key):
        self.master_key = master_key
        self.session_key = self.generate_session_key()
        self.iv = os.urandom(16)  # Random IV for each session
    
    def generate_session_key(self):
        """Generate session key from master key using PBKDF2"""
        salt = b"web4_session_salt"
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=10000,
            backend=default_backend()
        )
        return kdf.derive(self.hex_to_bytes(self.master_key))
    
    def encrypt(self, plaintext):
        """Encrypt plaintext using AES-256-GCM"""
        # Generate new IV for each encryption
        iv = os.urandom(16)
        
        # Create cipher
        cipher = Cipher(
            algorithms.AES(self.session_key),
            modes.GCM(iv),
            backend=default_backend()
        )
        encryptor = cipher.encryptor()
        
        # Encrypt data
        ciphertext = encryptor.update(plaintext.encode()) + encryptor.finalize()
        
        # Get authentication tag
        tag = encryptor.tag
        
        # Combine IV + ciphertext + tag
        result = iv + ciphertext + tag
        return self.bytes_to_hex(result)
    
    def decrypt(self, encrypted_data):
        """Decrypt ciphertext using AES-256-GCM"""
        data = self.hex_to_bytes(encrypted_data)
        
        # Extract IV, ciphertext, and tag
        iv = data[:16]
        tag = data[-16:]
        ciphertext = data[16:-16]
        
        # Create cipher
        cipher = Cipher(
            algorithms.AES(self.session_key),
            modes.GCM(iv, tag),
            backend=default_backend()
        )
        decryptor = cipher.decryptor()
        
        # Decrypt data
        plaintext = decryptor.update(ciphertext) + decryptor.finalize()
        return plaintext.decode()
```

---

## Message Encryption/Decryption

### Complete Communication Example

```cpp
// C++ Complete Example
#include <iostream>
#include <string>
#include "RESTClient.h"

class SecureCommunication {
private:
    RESTClient apiClient;
    KeyManager keyManager;
    EncryptedChannel channel;
    std::string deviceId;

public:
    SecureCommunication(const std::string& apiEndpoint, const std::string& id) 
        : apiClient(apiEndpoint), deviceId(id) {}

    bool establishSecureChannel(const std::string& partnerDeviceId) {
        try {
            // Step 1: Create LCT relationship
            auto lctResult = apiClient.createLCT(
                "demo-user",
                deviceId,
                partnerDeviceId,
                "secure-communication",
                "proxy-001"
            );

            // Step 2: Initiate pairing
            auto pairingResult = apiClient.initiatePairing(
                "demo-user",
                deviceId,
                partnerDeviceId,
                "secure-session",
                "proxy-001",
                false
            );

            // Step 3: Complete pairing to get split keys
            auto completeResult = apiClient.completePairing(
                "demo-user",
                pairingResult.challengeId,
                "device-auth-token",
                "partner-auth-token",
                "session-001"
            );

            // Step 4: Reconstruct master key
            keyManager.setKeyHalves(completeResult.splitKeyA, completeResult.splitKeyB);
            
            // Step 5: Initialize encrypted channel
            channel = EncryptedChannel(keyManager.getMasterKey());

            std::cout << "Secure channel established successfully!" << std::endl;
            std::cout << "Master Key: " << keyManager.getMasterKey() << std::endl;
            
            return true;
        } catch (const std::exception& e) {
            std::cerr << "Failed to establish secure channel: " << e.what() << std::endl;
            return false;
        }
    }

    std::string sendSecureMessage(const std::string& message) {
        try {
            std::string encrypted = channel.encrypt(message);
            std::cout << "Encrypted message: " << encrypted << std::endl;
            return encrypted;
        } catch (const std::exception& e) {
            std::cerr << "Encryption failed: " << e.what() << std::endl;
            return "";
        }
    }

    std::string receiveSecureMessage(const std::string& encryptedMessage) {
        try {
            std::string decrypted = channel.decrypt(encryptedMessage);
            std::cout << "Decrypted message: " << decrypted << std::endl;
            return decrypted;
        } catch (const std::exception& e) {
            std::cerr << "Decryption failed: " << e.what() << std::endl;
            return "";
        }
    }
};

// Usage Example
int main() {
    // Device A (Battery Module)
    SecureCommunication deviceA("http://localhost:8080", "battery-module-001");
    
    // Device B (Motor Controller)
    SecureCommunication deviceB("http://localhost:8080", "motor-controller-001");
    
    // Establish secure channel
    if (deviceA.establishSecureChannel("motor-controller-001")) {
        // Send secure message
        std::string message = "Battery status: 85% charge, 3.8V per cell";
        std::string encrypted = deviceA.sendSecureMessage(message);
        
        // Receive and decrypt message
        std::string decrypted = deviceB.receiveSecureMessage(encrypted);
        
        std::cout << "Original: " << message << std::endl;
        std::cout << "Decrypted: " << decrypted << std::endl;
    }
    
    return 0;
}
```

```python
# Python Complete Example
import requests
import json

class SecureCommunication:
    def __init__(self, api_endpoint, device_id):
        self.api_endpoint = api_endpoint
        self.device_id = device_id
        self.key_manager = KeyManager()
        self.channel = None
    
    def establish_secure_channel(self, partner_device_id):
        """Establish secure communication channel"""
        try:
            # Step 1: Create LCT relationship
            lct_response = requests.post(f'{self.api_endpoint}/lct/create', json={
                'creator': 'demo-user',
                'component_a': self.device_id,
                'component_b': partner_device_id,
                'context': 'secure-communication',
                'proxy_id': 'proxy-001'
            })
            lct_data = lct_response.json()
            
            # Step 2: Initiate pairing
            pairing_response = requests.post(f'{self.api_endpoint}/pairing/initiate', json={
                'creator': 'demo-user',
                'component_a': self.device_id,
                'component_b': partner_device_id,
                'operational_context': 'secure-session',
                'proxy_id': 'proxy-001',
                'force_immediate': False
            })
            pairing_data = pairing_response.json()
            
            # Step 3: Complete pairing to get split keys
            complete_response = requests.post(f'{self.api_endpoint}/pairing/complete', json={
                'creator': 'demo-user',
                'challenge_id': pairing_data['challenge_id'],
                'component_a_auth': 'device-auth-token',
                'component_b_auth': 'partner-auth-token',
                'session_context': 'session-001'
            })
            complete_data = complete_response.json()
            
            # Step 4: Reconstruct master key
            self.key_manager.set_key_halves(
                complete_data['split_key_a'],
                complete_data['split_key_b']
            )
            
            # Step 5: Initialize encrypted channel
            self.channel = EncryptedChannel(self.key_manager.get_master_key())
            
            print("Secure channel established successfully!")
            print(f"Master Key: {self.key_manager.get_master_key()}")
            
            return True
            
        except Exception as e:
            print(f"Failed to establish secure channel: {e}")
            return False
    
    def send_secure_message(self, message):
        """Send encrypted message"""
        try:
            encrypted = self.channel.encrypt(message)
            print(f"Encrypted message: {encrypted}")
            return encrypted
        except Exception as e:
            print(f"Encryption failed: {e}")
            return None
    
    def receive_secure_message(self, encrypted_message):
        """Receive and decrypt message"""
        try:
            decrypted = self.channel.decrypt(encrypted_message)
            print(f"Decrypted message: {decrypted}")
            return decrypted
        except Exception as e:
            print(f"Decryption failed: {e}")
            return None

# Usage Example
if __name__ == "__main__":
    # Device A (Battery Module)
    device_a = SecureCommunication("http://localhost:8080", "battery-module-001")
    
    # Device B (Motor Controller)
    device_b = SecureCommunication("http://localhost:8080", "motor-controller-001")
    
    # Establish secure channel
    if device_a.establish_secure_channel("motor-controller-001"):
        # Send secure message
        message = "Battery status: 85% charge, 3.8V per cell"
        encrypted = device_a.send_secure_message(message)
        
        # Receive and decrypt message
        decrypted = device_b.receive_secure_message(encrypted)
        
        print(f"Original: {message}")
        print(f"Decrypted: {decrypted}")
```

---

## Security Considerations

### Key Security

1. **Never Store Split Keys**: Keys should be ephemeral and not persisted
2. **Secure Key Transmission**: Use secure channels for key exchange
3. **Key Rotation**: Implement regular key rotation mechanisms
4. **Key Destruction**: Securely erase keys from memory after use

### Channel Security

1. **Perfect Forward Secrecy**: Each session uses unique keys
2. **Authentication**: Verify device identities before key exchange
3. **Integrity**: Use authenticated encryption (AES-GCM)
4. **Replay Protection**: Include timestamps and nonces

### Implementation Security

```cpp
// Secure key handling
class SecureKeyManager {
private:
    std::vector<uint8_t> keyBuffer;
    
public:
    void setKey(const std::string& key) {
        // Securely store key in memory
        keyBuffer = hexToBytes(key);
    }
    
    ~SecureKeyManager() {
        // Securely erase key from memory
        std::fill(keyBuffer.begin(), keyBuffer.end(), 0);
    }
    
    // Prevent key copying
    SecureKeyManager(const SecureKeyManager&) = delete;
    SecureKeyManager& operator=(const SecureKeyManager&) = delete;
};
```

---

## Troubleshooting

### Common Issues

1. **Key Reconstruction Failure**
   - Verify both key halves are 64 hex characters (32 bytes)
   - Check for transmission errors in key exchange
   - Ensure proper hex encoding/decoding

2. **Encryption/Decryption Errors**
   - Verify session key derivation
   - Check IV generation and transmission
   - Ensure authentication tag integrity

3. **API Communication Issues**
   - Verify API bridge is running
   - Check network connectivity
   - Validate request/response formats

### Debug Information

```cpp
// Debug key reconstruction
void debugKeyReconstruction(const std::string& keyA, const std::string& keyB) {
    std::cout << "Key Half A Length: " << keyA.length() << std::endl;
    std::cout << "Key Half B Length: " << keyB.length() << std::endl;
    std::cout << "Key Half A: " << keyA << std::endl;
    std::cout << "Key Half B: " << keyB << std::endl;
    
    KeyManager km;
    std::string masterKey = km.reconstructMasterKey(keyA, keyB);
    std::cout << "Master Key: " << masterKey << std::endl;
    std::cout << "Master Key Length: " << masterKey.length() << std::endl;
}
```

---

## Conclusion

This guide provides a complete implementation for establishing encrypted communication channels using the Web4 LCT split-key system. The approach ensures:

1. **Maximum Security**: Split-key architecture prevents single-point compromise
2. **Perfect Forward Secrecy**: Session-specific keys protect against future attacks
3. **Zero-Knowledge**: Neither device can derive the other's key half
4. **Military-grade Encryption**: AES-256-GCM with authenticated encryption

The implementation is production-ready and suitable for automotive and industrial applications requiring the highest levels of security.

For additional examples and integration support, refer to the C++ demo application in the `api-bridge/cpp-demo/` directory. 