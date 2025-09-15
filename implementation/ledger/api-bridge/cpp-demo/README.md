# Web4 Race Car Battery Management API Bridge - C++ Demo

A comprehensive C++ demo application for the Web4 Race Car Battery Management API Bridge, compatible with RAD Studio and other C++ development environments. This demo showcases the latest privacy-focused features and real blockchain integration.

## Features

- **Dual Interface Support**: Both REST and gRPC API interfaces
- **Privacy-Focused Features**: Anonymous component registration and verification
- **Complete API Coverage**: All blockchain operations including:
  - Account Management
  - Component Registry (Legacy)
  - Privacy-Focused Component Operations
  - LCT (Linked Context Token) Management
  - Pairing Process with Split-Key Generation
  - Pairing Queue Operations
  - Trust Tensor Operations
  - Energy Operations
- **Real-time Monitoring**: Battery status streaming via gRPC
- **Performance Testing**: REST vs gRPC comparison tools
- **RAD Studio Compatible**: Designed for Windows development environments
- **Cross-platform**: Works on Windows, Linux, and macOS
- **Real Blockchain Integration**: No mock responses, all operations use actual blockchain transactions

## Privacy Features

The demo now includes comprehensive privacy-focused features:

- **Anonymous Component Registration**: Register components using cryptographic hashes
- **Hash-based Pairing Verification**: Verify component pairings without exposing real IDs
- **Anonymous Pairing Authorization**: Create pairing authorizations using hashes
- **Revocation Events**: Track component revocations anonymously
- **Metadata Retrieval**: Access component metadata while preserving privacy

## Prerequisites

### Required Dependencies

- **C++17 compatible compiler** (MSVC, GCC, Clang)
- **CMake 3.16+**
- **Git**

### Optional Dependencies (for full functionality)

- **gRPC** - For gRPC interface support
- **Protobuf** - For protocol buffer support
- **OpenSSL** - For SSL/TLS support

## Building the Application

### Windows (RAD Studio Compatible)

1. **Install Dependencies**:
   ```powershell
   # Install vcpkg (if not already installed)
   git clone https://github.com/Microsoft/vcpkg.git
   cd vcpkg
   .\bootstrap-vcpkg.bat
   .\vcpkg integrate install
   
   # Install required packages
   .\vcpkg install nlohmann-json:x64-windows
   .\vcpkg install grpc:x64-windows
   .\vcpkg install protobuf:x64-windows
   ```

2. **Build with CMake**:
   ```powershell
   mkdir build
   cd build
   cmake .. -DCMAKE_TOOLCHAIN_FILE=C:/vcpkg/scripts/buildsystems/vcpkg.cmake
   cmake --build . --config Release
   ```

3. **Alternative: Build with Visual Studio**:
   ```powershell
   cmake .. -G "Visual Studio 16 2019" -A x64
   cmake --build . --config Release
   ```

### Linux/macOS

1. **Install Dependencies**:
   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install build-essential cmake git libssl-dev
   
   # macOS
   brew install cmake openssl
   ```

2. **Build**:
   ```bash
   mkdir build
   cd build
   cmake ..
   make -j$(nproc)
   ```

## Configuration

The application uses `config.json` for configuration. Key settings:

```json
{
  "api": {
    "rest": {
      "endpoint": "http://localhost:8080"
    },
    "grpc": {
      "endpoint": "localhost:9092"
    }
  },
  "privacy": {
    "enabled": true,
    "hash_algorithm": "sha256"
  }
}
```

## Usage

### Running the Demo

```bash
# Basic usage
./APIBridgeDemo

# With custom config
./APIBridgeDemo --config /path/to/config.json

# Test mode
./APIBridgeDemo --test

# Verbose logging
./APIBridgeDemo --verbose
```

### Menu Options

1. **Account Management**
   - List existing accounts
   - Create new accounts
   - View account details

2. **Component Registry (Legacy)**
   - Register new components (legacy method)
   - Retrieve component information
   - Verify component authenticity

3. **Privacy-Focused Features** ⭐ **NEW**
   - Register anonymous components
   - Verify pairing with hashes
   - Create anonymous pairing authorizations
   - Create revocation events
   - Get anonymous component metadata

4. **LCT Management**
   - Create Linked Context Tokens
   - Manage LCT relationships
   - Update LCT status

5. **Pairing Process**
   - Initiate device pairing
   - Complete pairing with split-key generation
   - Manage pairing sessions

6. **Pairing Queue Operations** ⭐ **NEW**
   - Queue pairing requests
   - Get queue status
   - Process offline queues
   - Cancel requests
   - List proxy queues

7. **Trust Tensor**
   - Create trust relationships
   - Update trust scores
   - Monitor trust metrics

8. **Energy Operations**
   - Create energy operations
   - Execute energy transfers
   - Monitor energy balances

9. **Real-time Monitoring** (gRPC only)
   - Stream battery status updates
   - Monitor system health

10. **Performance Testing**
    - Compare REST vs gRPC performance
    - Benchmark API operations

11. **System Information**
    - View API Bridge status
    - Check blockchain connectivity
    - Display feature availability

## API Examples

### Privacy-Focused Component Operations

```cpp
#include "RESTClient.h"

// Initialize client
RESTClient client("http://localhost:8080");

// Register anonymous component
auto anonResult = client.registerAnonymousComponent(
    "user1", 
    "real-battery-id-123", 
    "tesla-motors", 
    "lithium-ion-battery", 
    "race-car"
);
std::cout << "Component Hash: " << anonResult.componentHash << std::endl;

// Verify pairing with hashes
auto verifyResult = client.verifyComponentPairingWithHashes(
    "verifier-001",
    anonResult.componentHash,
    "motor-hash-456",
    "race-car-context"
);
std::cout << "Verification Status: " << verifyResult.status << std::endl;
```

### Pairing Queue Operations

```cpp
#include "RESTClient.h"

// Queue pairing request
auto queueResult = client.queuePairingRequest(
    "user1", 
    "battery-001", 
    "motor-001", 
    "race-car-queue"
);
std::cout << "Request ID: " << queueResult.requestId << std::endl;

// Get queue status
auto statusResult = client.getQueueStatus("default-queue");
std::cout << "Pending: " << statusResult.pendingRequests << std::endl;
```

### REST API Usage (Legacy)

```cpp
#include "RESTClient.h"

// Initialize client
RESTClient client("http://localhost:8080");

// Register a component (legacy method)
auto result = client.registerComponent("user1", "battery-module-v1", "race-car");
std::cout << "Component ID: " << result.componentId << std::endl;

// Create LCT
auto lct = client.createLCT("user1", "battery-001", "motor-001", "pairing", "proxy-1");
std::cout << "LCT ID: " << lct.lctId << std::endl;
```

### gRPC API Usage

```cpp
#include "GRPCClient.h"

// Initialize client
GRPCClient client("localhost:9092");

// Stream battery status
client.streamBatteryStatus("battery-001", 10, [](const BatteryStatusUpdate& update) {
    std::cout << "Voltage: " << update.voltage << "V" << std::endl;
    std::cout << "Current: " << update.current << "A" << std::endl;
});
```

## Integration with RAD Studio

### Project Setup

1. **Create New Project**:
   - File → New → Other → C++ → Console Application
   - Choose C++17 standard

2. **Add Source Files**:
   - Add all `.cpp` and `.h` files to your project
   - Ensure proper include paths

3. **Configure Libraries**:
   - Add library paths for httplib and nlohmann/json
   - Link required system libraries

### Build Configuration

```cpp
// In your RAD Studio project settings:
// - C++17 standard
// - Include paths: $(PROJECTDIR)/include
// - Library paths: $(PROJECTDIR)/lib
// - Link libraries: ws2_32, iphlpapi, crypt32
```

## Privacy Benefits

The privacy-focused features provide several key benefits:

- **Trade Secret Protection**: Real component IDs and manufacturer information are never exposed on-chain
- **Anonymous Verification**: Components can be verified without revealing their true identity
- **Hash-based Operations**: All sensitive operations use cryptographic hashes
- **Revocation Tracking**: Track component revocations while maintaining privacy
- **Metadata Access**: Access component metadata without exposing real identifiers

## Testing

### Unit Tests

```bash
# Run tests
cd build
ctest --verbose
```

### Integration Tests

```bash
# Start API bridge first
cd ../api-bridge
go run main.go

# Then run demo tests
./APIBridgeDemo --test
```

### Privacy Feature Tests

```bash
# Test privacy features specifically
./APIBridgeDemo --privacy-test

# Test pairing queue operations
./APIBridgeDemo --queue-test
```

## Troubleshooting

### Common Issues

1. **gRPC Not Available**:
   - Install gRPC and Protobuf
   - Check CMake configuration
   - Verify library linking

2. **Privacy Features Not Working**:
   - Ensure API Bridge is running with privacy features enabled
   - Check blockchain connection
   - Verify endpoint configuration

3. **Build Errors**:
   - Ensure C++17 standard is enabled
   - Check all dependencies are installed
   - Verify CMake configuration

4. **Connection Issues**:
   - Verify API Bridge is running on correct ports
   - Check firewall settings
   - Ensure blockchain is accessible

## API Endpoints

### Privacy-Focused Endpoints

- `POST /components/register-anonymous` - Register anonymous component
- `POST /components/verify-pairing-hashes` - Verify pairing with hashes
- `POST /components/create-pairing-authorization` - Create anonymous pairing authorization
- `POST /components/create-revocation-event` - Create revocation event
- `GET /components/anonymous/{hash}/metadata` - Get anonymous component metadata

### Pairing Queue Endpoints

- `POST /pairing/queue` - Queue pairing request
- `GET /pairing/queue/{id}/status` - Get queue status
- `GET /pairing/queue/{id}/requests` - Get queued requests
- `GET /pairing/queue/proxy/{id}` - List proxy queue
- `POST /pairing/queue/process` - Process offline queue
- `POST /pairing/queue/cancel` - Cancel request

### Legacy Endpoints

- `POST /components/register` - Register component (legacy)
- `GET /components/{id}` - Get component (legacy)
- `POST /components/verify` - Verify component (legacy)

## Performance Considerations

- **REST vs gRPC**: gRPC generally provides better performance for streaming operations
- **Privacy Operations**: Hash-based operations may have slightly higher latency
- **Queue Operations**: Offline queue processing is optimized for batch operations
- **Real Blockchain**: All operations use actual blockchain transactions

## Security Notes

- All privacy features use cryptographic hashing
- No sensitive data is stored on-chain
- Component metadata is stored off-chain
- Pairing authorizations use hash-based verification
- Revocation events maintain privacy while providing traceability

## Future Enhancements

- Enhanced privacy features with zero-knowledge proofs
- Advanced queue management with priority levels
- Real-time privacy monitoring
- Integration with additional blockchain networks
- Enhanced RAD Studio integration tools 