# Web4 Race Car Battery Management API Bridge - C++ Integration Guide

This guide provides comprehensive instructions for integrating the Web4 Race Car Battery Management API Bridge into C++ applications, with a focus on the latest privacy-focused features and real blockchain integration.

## Overview

The API Bridge provides a unified interface for interacting with the Web4 blockchain, supporting both REST and gRPC protocols. The latest version includes comprehensive privacy features that protect trade secrets while maintaining full verification capabilities.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   C++ Client    │    │   API Bridge    │    │   Blockchain    │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ REST Client │◄┼────┼►│ REST Server │ │    │ │ Ignite CLI  │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│                 │    │                 │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │ gRPC Client │◄┼────┼►│ gRPC Server │ │    │ │ Real Chain  │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Privacy Features

The API Bridge now includes comprehensive privacy features:

### Anonymous Component Registration
- Components are registered using cryptographic hashes
- Real component IDs and manufacturer information are never exposed on-chain
- Full verification capabilities are maintained

### Hash-based Operations
- All sensitive operations use cryptographic hashes
- Pairing verification uses hash-based matching
- Authorization and revocation events maintain privacy

### Metadata Management
- Component metadata is stored off-chain
- Access is provided through hash-based lookups
- No sensitive data is exposed in blockchain transactions

## Quick Start

### 1. Include Headers

```cpp
#include "RESTClient.h"
#include "GRPCClient.h"
#include "DemoUI.h"
```

### 2. Initialize Clients

```cpp
// REST Client
RESTClient restClient("http://localhost:8080");

// gRPC Client (optional)
GRPCClient grpcClient("localhost:9092");
```

### 3. Basic Operations

```cpp
// Privacy-focused component registration
auto anonResult = restClient.registerAnonymousComponent(
    "user1", 
    "real-battery-id-123", 
    "tesla-motors", 
    "lithium-ion-battery", 
    "race-car"
);

// Verify pairing with hashes
auto verifyResult = restClient.verifyComponentPairingWithHashes(
    "verifier-001",
    anonResult.componentHash,
    "motor-hash-456",
    "race-car-context"
);
```

## API Reference

### Privacy-Focused Operations

#### Register Anonymous Component

```cpp
AnonymousComponentResult registerAnonymousComponent(
    const std::string& creator,
    const std::string& realComponentId,
    const std::string& manufacturerId,
    const std::string& componentType,
    const std::string& context
);
```

**Parameters:**
- `creator`: User identifier
- `realComponentId`: Actual component ID (not stored on-chain)
- `manufacturerId`: Manufacturer identifier (not stored on-chain)
- `componentType`: Component type/category (not stored on-chain)
- `context`: Operational context

**Returns:**
- `componentHash`: Cryptographic hash of the component
- `manufacturerHash`: Hash of manufacturer information
- `categoryHash`: Hash of component category
- `txHash`: Blockchain transaction hash

#### Verify Pairing with Hashes

```cpp
PairingVerificationResult verifyComponentPairingWithHashes(
    const std::string& verifier,
    const std::string& sourceHash,
    const std::string& targetHash,
    const std::string& context
);
```

**Parameters:**
- `verifier`: Verifier identifier
- `sourceHash`: Hash of source component
- `targetHash`: Hash of target component
- `context`: Verification context

**Returns:**
- `sourceHash`: Confirmed source component hash
- `targetHash`: Confirmed target component hash
- `status`: Verification status
- `txHash`: Blockchain transaction hash

#### Create Anonymous Pairing Authorization

```cpp
PairingAuthorizationResult createAnonymousPairingAuthorization(
    const std::string& creator,
    const std::string& sourceHash,
    const std::string& targetHash,
    const std::string& context
);
```

**Parameters:**
- `creator`: Authorization creator
- `sourceHash`: Hash of source component
- `targetHash`: Hash of target component
- `context`: Authorization context

**Returns:**
- `authorizationId`: Unique authorization identifier
- `sourceHash`: Source component hash
- `targetHash`: Target component hash
- `status`: Authorization status
- `txHash`: Blockchain transaction hash

#### Create Revocation Event

```cpp
RevocationEventResult createAnonymousRevocationEvent(
    const std::string& creator,
    const std::string& componentHash,
    const std::string& reason,
    const std::string& context
);
```

**Parameters:**
- `creator`: Revocation creator
- `componentHash`: Hash of component to revoke
- `reason`: Revocation reason
- `context`: Revocation context

**Returns:**
- `revocationId`: Unique revocation identifier
- `componentHash`: Component hash
- `reason`: Revocation reason
- `status`: Revocation status
- `txHash`: Blockchain transaction hash

#### Get Anonymous Component Metadata

```cpp
ComponentMetadataResult getAnonymousComponentMetadata(
    const std::string& componentHash
);
```

**Parameters:**
- `componentHash`: Hash of the component

**Returns:**
- `componentHash`: Component hash
- `metadata`: Component metadata (off-chain)
- `status`: Retrieval status
- `txHash`: Blockchain transaction hash

### Pairing Queue Operations

#### Queue Pairing Request

```cpp
PairingRequestResult queuePairingRequest(
    const std::string& creator,
    const std::string& componentA,
    const std::string& componentB,
    const std::string& context
);
```

**Parameters:**
- `creator`: Request creator
- `componentA`: First component identifier
- `componentB`: Second component identifier
- `context`: Request context

**Returns:**
- `requestId`: Unique request identifier
- `componentA`: First component
- `componentB`: Second component
- `status`: Request status
- `createdAt`: Creation timestamp
- `txHash`: Blockchain transaction hash

#### Get Queue Status

```cpp
QueueStatusResult getQueueStatus(const std::string& queueId);
```

**Parameters:**
- `queueId`: Queue identifier

**Returns:**
- `queueId`: Queue identifier
- `pendingRequests`: Number of pending requests
- `processedRequests`: Number of processed requests
- `status`: Queue status
- `txHash`: Blockchain transaction hash

#### Process Offline Queue

```cpp
std::string processOfflineQueue(
    const std::string& processor,
    const std::string& queueId,
    const std::string& context
);
```

**Parameters:**
- `processor`: Processor identifier
- `queueId`: Queue to process
- `context`: Processing context

**Returns:**
- Processing result string

#### Cancel Request

```cpp
std::string cancelRequest(
    const std::string& creator,
    const std::string& requestId,
    const std::string& reason
);
```

**Parameters:**
- `creator`: Request creator
- `requestId`: Request to cancel
- `reason`: Cancellation reason

**Returns:**
- Cancellation result string

### Legacy Operations

#### Register Component (Legacy)

```cpp
ComponentRegistrationResult registerComponent(
    const std::string& creator,
    const std::string& componentData,
    const std::string& context
);
```

**Note:** This method stores component data directly on-chain and should be avoided for privacy-sensitive applications.

## Error Handling

### Exception Types

```cpp
try {
    auto result = restClient.registerAnonymousComponent(...);
} catch (const std::exception& e) {
    std::cerr << "Error: " << e.what() << std::endl;
}
```

### Common Error Scenarios

1. **Network Errors**: Connection timeouts, DNS resolution failures
2. **Authentication Errors**: Invalid credentials, expired tokens
3. **Validation Errors**: Invalid parameters, missing required fields
4. **Blockchain Errors**: Transaction failures, insufficient funds
5. **Privacy Errors**: Hash validation failures, metadata access denied

## Performance Optimization

### Connection Pooling

```cpp
// REST client automatically manages connection pooling
RESTClient client("http://localhost:8080");
client.set_connection_timeout(10);  // 10 seconds
client.set_read_timeout(30);        // 30 seconds
```

### Batch Operations

```cpp
// Process multiple requests efficiently
std::vector<PairingRequestResult> results;
for (const auto& request : requests) {
    try {
        auto result = restClient.queuePairingRequest(
            request.creator,
            request.componentA,
            request.componentB,
            request.context
        );
        results.push_back(result);
    } catch (const std::exception& e) {
        // Handle individual request failures
    }
}
```

### Caching

```cpp
// Cache frequently accessed data
std::unordered_map<std::string, ComponentMetadataResult> metadataCache;

auto getCachedMetadata = [&](const std::string& componentHash) {
    auto it = metadataCache.find(componentHash);
    if (it != metadataCache.end()) {
        return it->second;
    }
    
    auto result = restClient.getAnonymousComponentMetadata(componentHash);
    metadataCache[componentHash] = result;
    return result;
};
```

## Security Considerations

### Hash Generation

```cpp
// Use secure hash generation for component IDs
#include <openssl/sha.h>

std::string generateComponentHash(const std::string& realComponentId) {
    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256_CTX sha256;
    SHA256_Init(&sha256);
    SHA256_Update(&sha256, realComponentId.c_str(), realComponentId.length());
    SHA256_Final(hash, &sha256);
    
    std::stringstream ss;
    for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << (int)hash[i];
    }
    return ss.str();
}
```

### Input Validation

```cpp
// Validate all inputs before sending to API
bool validateComponentData(const std::string& componentData) {
    if (componentData.empty() || componentData.length() > 1000) {
        return false;
    }
    
    // Check for valid characters
    for (char c : componentData) {
        if (!std::isalnum(c) && c != '-' && c != '_' && c != '.') {
            return false;
        }
    }
    
    return true;
}
```

## Testing

### Unit Tests

```cpp
#include <gtest/gtest.h>

class PrivacyFeaturesTest : public ::testing::Test {
protected:
    RESTClient client;
    
    PrivacyFeaturesTest() : client("http://localhost:8080") {}
};

TEST_F(PrivacyFeaturesTest, AnonymousRegistration) {
    auto result = client.registerAnonymousComponent(
        "test-user",
        "test-component-123",
        "test-manufacturer",
        "test-type",
        "test-context"
    );
    
    EXPECT_FALSE(result.componentHash.empty());
    EXPECT_FALSE(result.manufacturerHash.empty());
    EXPECT_FALSE(result.categoryHash.empty());
    EXPECT_FALSE(result.txHash.empty());
}
```

### Integration Tests

```cpp
TEST_F(PrivacyFeaturesTest, FullPrivacyFlow) {
    // 1. Register anonymous component
    auto anonResult = client.registerAnonymousComponent(...);
    
    // 2. Create pairing authorization
    auto authResult = client.createAnonymousPairingAuthorization(...);
    
    // 3. Verify pairing
    auto verifyResult = client.verifyComponentPairingWithHashes(...);
    
    // 4. Get metadata
    auto metadataResult = client.getAnonymousComponentMetadata(...);
    
    // 5. Create revocation event
    auto revokeResult = client.createAnonymousRevocationEvent(...);
    
    // Verify all operations succeeded
    EXPECT_EQ(anonResult.status, "success");
    EXPECT_EQ(authResult.status, "success");
    EXPECT_EQ(verifyResult.status, "success");
    EXPECT_EQ(metadataResult.status, "success");
    EXPECT_EQ(revokeResult.status, "success");
}
```

## Deployment

### Production Configuration

```json
{
  "api": {
    "rest": {
      "endpoint": "https://api-bridge.production.com",
      "timeout": 30,
      "retries": 3
    },
    "grpc": {
      "endpoint": "api-bridge.production.com:9092",
      "timeout": 30
    }
  },
  "privacy": {
    "enabled": true,
    "hash_algorithm": "sha256",
    "metadata_storage": "secure"
  },
  "security": {
    "ssl_verify": true,
    "certificate_path": "/path/to/cert.pem",
    "key_path": "/path/to/key.pem"
  }
}
```

### Monitoring

```cpp
// Implement health checks
auto healthStatus = restClient.getHealthStatus();
auto blockchainStatus = restClient.getBlockchainStatus();

if (healthStatus.find("healthy") == std::string::npos) {
    // Handle unhealthy API Bridge
}

if (blockchainStatus.find("connected") == std::string::npos) {
    // Handle blockchain connection issues
}
```

## Troubleshooting

### Common Issues

1. **Privacy Features Not Working**
   - Ensure API Bridge is running with privacy features enabled
   - Check blockchain connection status
   - Verify hash generation is working correctly

2. **Queue Operations Failing**
   - Check queue configuration
   - Verify processor permissions
   - Ensure queue exists and is accessible

3. **Performance Issues**
   - Monitor network latency
   - Check blockchain transaction times
   - Implement connection pooling
   - Use batch operations where possible

4. **Security Issues**
   - Verify SSL/TLS configuration
   - Check certificate validity
   - Ensure proper input validation
   - Monitor for suspicious activity

### Debug Mode

```cpp
// Enable debug logging
client.set_debug_mode(true);

// Check detailed error information
try {
    auto result = client.registerAnonymousComponent(...);
} catch (const std::exception& e) {
    std::cerr << "Detailed error: " << e.what() << std::endl;
    // Check logs for additional information
}
```

## Best Practices

1. **Always use privacy features** for sensitive component data
2. **Validate all inputs** before sending to API
3. **Implement proper error handling** for all operations
4. **Use connection pooling** for better performance
5. **Cache frequently accessed data** when appropriate
6. **Monitor API Bridge health** regularly
7. **Implement retry logic** for transient failures
8. **Use secure hash generation** for component IDs
9. **Keep dependencies updated** for security patches
10. **Test thoroughly** before production deployment

## Support

For additional support and questions:

- **Documentation**: See the main README.md
- **Examples**: Check the demo application
- **Issues**: Report problems via GitHub Issues
- **Community**: Join the developer community forum 