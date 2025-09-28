# Key Exchange Protocol for Embedded Devices

This document outlines the implementation of the split-key cryptography protocol for embedded devices, particularly ARM-based microcontrollers like STM32.

## Cryptographic Requirements

### Required Libraries
- Ed25519 for digital signatures
- SHA-256 for hashing and key derivation
- ChaCha20-Poly1305 for encryption
- Curve25519 for key exchange
- Random number generator (TRNG)

### Memory Requirements
- Ed25519: ~1KB RAM for key generation and signing
- SHA-256: ~200 bytes RAM
- ChaCha20-Poly1305: ~1.5KB RAM
- Curve25519: ~1KB RAM
- Total: ~4KB RAM minimum

### Flash Requirements
- Ed25519: ~8KB
- SHA-256: ~2KB
- ChaCha20-Poly1305: ~4KB
- Curve25519: ~6KB
- Total: ~20KB minimum

## Implementation Guide

### 1. Key Generation
```c
// Example using mbedtls (common on embedded devices)
#include "mbedtls/ed25519.h"
#include "mbedtls/entropy.h"
#include "mbedtls/ctr_drbg.h"

// Generate Ed25519 key pair
mbedtls_ed25519_context ctx;
mbedtls_ed25519_init(&ctx);
mbedtls_ed25519_genkey(&ctx, MBEDTLS_ENTROPY_SOURCE_STRONG);
```

### 2. Key Share Generation
```c
// Generate 32-byte random key share
uint8_t key_share[32];
mbedtls_ctr_drbg_random(&ctr_drbg, key_share, 32);
```

### 3. Challenge-Response
```c
// Generate challenge response using SHA-256
uint8_t challenge[32];
uint8_t response[32];
mbedtls_sha256_context sha_ctx;

mbedtls_sha256_init(&sha_ctx);
mbedtls_sha256_starts(&sha_ctx, 0);
mbedtls_sha256_update(&sha_ctx, challenge, 32);
mbedtls_sha256_update(&sha_ctx, key_share, 32);
mbedtls_sha256_finish(&sha_ctx, response);
```

### 4. Key Exchange Protocol

#### Initiator Side
1. Generate Ed25519 key pair
2. Generate random key share
3. Create key exchange message:
   ```c
   struct key_exchange_msg {
       uint8_t lct_id[32];
       uint8_t component_id[32];
       uint8_t public_key[32];
       uint8_t key_share[32];
       uint8_t nonce[12];
       uint8_t signature[64];
   };
   ```
4. Sign message using Ed25519
5. Send to target component

#### Target Side
1. Verify initiator's signature
2. Generate own Ed25519 key pair
3. Generate own key share
4. Create response message
5. Sign and send response

### 5. Key Derivation
```c
// Derive combined key using SHA-256
uint8_t combined_key[32];
mbedtls_sha256_context sha_ctx;

mbedtls_sha256_init(&sha_ctx);
mbedtls_sha256_starts(&sha_ctx, 0);
mbedtls_sha256_update(&sha_ctx, key_share_a, 32);
mbedtls_sha256_update(&sha_ctx, key_share_b, 32);
mbedtls_sha256_update(&sha_ctx, shared_secret, 32);
mbedtls_sha256_finish(&sha_ctx, combined_key);
```

## Security Considerations

### 1. Random Number Generation
- Use hardware TRNG when available
- Fall back to cryptographically secure PRNG
- Example for STM32:
  ```c
  // Enable RNG peripheral
  RCC->AHB2ENR |= RCC_AHB2ENR_RNGEN;
  RNG->CR |= RNG_CR_RNGEN;
  
  // Wait for valid random number
  while(!(RNG->SR & RNG_SR_DRDY));
  uint32_t random = RNG->DR;
  ```

### 2. Key Storage
- Store private keys in secure storage if available
- Use flash memory with write protection
- Consider using STM32's RDP (Read Protection) levels

### 3. Timing Attacks
- Use constant-time operations for cryptographic functions
- Avoid branching based on secret data
- Example:
  ```c
  // Constant-time comparison
  int compare_constant_time(const uint8_t *a, const uint8_t *b, size_t len) {
      int result = 0;
      for (size_t i = 0; i < len; i++) {
          result |= a[i] ^ b[i];
      }
      return result == 0;
  }
  ```

### 4. Power Analysis
- Consider using STM32's power analysis protection features
- Implement basic power analysis countermeasures:
  - Random delays
  - Constant power consumption patterns
  - Shuffling operations

## Error Handling

### 1. Cryptographic Operations
```c
// Example error handling for cryptographic operations
int crypto_operation(void) {
    int ret;
    mbedtls_ed25519_context ctx;
    
    mbedtls_ed25519_init(&ctx);
    
    ret = mbedtls_ed25519_genkey(&ctx, MBEDTLS_ENTROPY_SOURCE_STRONG);
    if (ret != 0) {
        // Handle error
        return CRYPTO_ERROR_KEY_GEN;
    }
    
    // ... rest of operation
    
    mbedtls_ed25519_free(&ctx);
    return CRYPTO_SUCCESS;
}
```

### 2. Memory Management
- Pre-allocate buffers for cryptographic operations
- Use static allocation where possible
- Example:
  ```c
  // Pre-allocated buffers
  static uint8_t key_buffer[32];
  static uint8_t signature_buffer[64];
  static mbedtls_sha256_context sha_ctx;
  ```

## Performance Optimization

### 1. STM32-Specific Optimizations
- Use STM32's cryptographic hardware acceleration when available
- Enable instruction cache
- Use DMA for data transfers
- Example:
  ```c
  // Enable cryptographic hardware acceleration
  RCC->AHB2ENR |= RCC_AHB2ENR_CRYPEN;
  
  // Configure DMA for cryptographic operations
  // ... DMA configuration code
  ```

### 2. Memory Optimization
- Use word-aligned buffers
- Minimize stack usage
- Use static allocation for large buffers
- Example:
  ```c
  // Word-aligned buffer
  static uint32_t __attribute__((aligned(4))) crypto_buffer[8];
  ```

## Testing

### 1. Unit Tests
- Test each cryptographic operation independently
- Verify key generation and exchange
- Test error conditions
- Example:
  ```c
  void test_key_generation(void) {
      uint8_t public_key[32];
      uint8_t private_key[32];
      
      int ret = generate_key_pair(public_key, private_key);
      TEST_ASSERT_EQUAL(0, ret);
      
      // Verify key pair
      ret = verify_key_pair(public_key, private_key);
      TEST_ASSERT_EQUAL(0, ret);
  }
  ```

### 2. Integration Tests
- Test complete key exchange protocol
- Verify challenge-response mechanism
- Test key rotation
- Example:
  ```c
  void test_key_exchange(void) {
      // Simulate two components
      component_t comp_a, comp_b;
      
      // Initialize components
      init_component(&comp_a);
      init_component(&comp_b);
      
      // Perform key exchange
      int ret = perform_key_exchange(&comp_a, &comp_b);
      TEST_ASSERT_EQUAL(0, ret);
      
      // Verify shared key
      ret = verify_shared_key(&comp_a, &comp_b);
      TEST_ASSERT_EQUAL(0, ret);
  }
  ```

## Common Issues and Solutions

### 1. Memory Constraints
- Problem: Insufficient RAM for cryptographic operations
- Solution: 
  - Use streaming operations where possible
  - Implement memory-efficient algorithms
  - Consider using external secure memory

### 2. Performance Issues
- Problem: Slow cryptographic operations
- Solution:
  - Use hardware acceleration
  - Optimize critical paths
  - Implement operation queuing

### 3. Security Vulnerabilities
- Problem: Potential side-channel attacks
- Solution:
  - Implement constant-time operations
  - Use secure key storage
  - Enable security features

## References

1. STM32 Cryptographic Library
2. MbedTLS Documentation
3. ARM Cryptographic Extensions
4. STM32 Security Features 