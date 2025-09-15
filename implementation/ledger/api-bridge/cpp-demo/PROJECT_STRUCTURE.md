# C++ Demo Application Project Structure

## Overview

This C++ demo application provides a comprehensive reference implementation for the Web4 Race Car Battery Management API Bridge, compatible with RAD Studio and other C++ development environments.

## Directory Structure

```
cpp-demo/
├── APIBridgeDemo.cpp          # Main application entry point
├── RESTClient.h               # REST API client header
├── RESTClient.cpp             # REST API client implementation
├── GRPCClient.h               # gRPC API client header
├── GRPCClient.cpp             # gRPC API client implementation
├── DemoUI.h                   # User interface header
├── DemoUI.cpp                 # User interface implementation
├── test_simple.cpp            # Simple C++17 test
├── config.json                # Configuration file
├── CMakeLists.txt             # CMake build configuration
├── build.bat                  # Windows build script
├── build.sh                   # Linux/macOS build script
├── README.md                  # Comprehensive documentation
├── PROJECT_STRUCTURE.md       # This file
└── build/                     # Build output directory (created during build)
    └── bin/
        ├── APIBridgeDemo      # Executable
        └── config.json        # Copied configuration
```

## File Descriptions

### Core Application Files

- **APIBridgeDemo.cpp**: Main application with menu-driven interface
  - Demonstrates all API operations
  - Provides performance comparison tools
  - Handles both REST and gRPC interfaces

- **RESTClient.h/cpp**: REST API client implementation
  - HTTP client using httplib
  - JSON parsing with nlohmann/json
  - Complete API coverage for all blockchain operations

- **GRPCClient.h/cpp**: gRPC API client implementation
  - gRPC client using official gRPC C++ library
  - Protocol buffer support
  - Real-time streaming capabilities

- **DemoUI.h/cpp**: User interface implementation
  - Console-based menu system
  - Input validation and formatting
  - Progress indicators and status displays

### Build and Configuration

- **CMakeLists.txt**: CMake build configuration
  - Dependency management (httplib, nlohmann/json, gRPC)
  - Platform-specific settings
  - Installation and packaging

- **config.json**: Application configuration
  - API endpoints and timeouts
  - Demo settings and defaults
  - Logging and UI preferences

- **build.bat**: Windows build script
  - Automated dependency installation
  - vcpkg integration
  - One-click build process

- **build.sh**: Linux/macOS build script
  - Cross-platform dependency management
  - System package detection
  - Automated build process

### Testing and Documentation

- **test_simple.cpp**: Basic C++17 functionality test
  - Verifies compiler compatibility
  - Tests modern C++ features
  - Quick validation tool

- **README.md**: Comprehensive documentation
  - Installation instructions
  - Usage examples
  - Troubleshooting guide

## Key Features

### API Coverage

1. **Account Management**
   - List accounts
   - Create accounts
   - Account details

2. **Component Registry**
   - Register components
   - Retrieve component info
   - Verify components

3. **LCT Management**
   - Create Linked Context Tokens
   - Manage relationships
   - Update status

4. **Pairing Process**
   - Initiate pairing
   - Complete pairing with split-key generation
   - Manage sessions

5. **Trust Tensor**
   - Create trust relationships
   - Update scores
   - Monitor metrics

6. **Energy Operations**
   - Create operations
   - Execute transfers
   - Monitor balances

### Interface Support

- **REST API**: HTTP-based interface
  - Easy to debug
  - Human-readable
  - Wide compatibility

- **gRPC API**: High-performance interface
  - Binary protocol
  - Streaming support
  - Lower latency

### Development Features

- **RAD Studio Compatible**: Windows development environment
- **Cross-platform**: Windows, Linux, macOS
- **Modern C++**: C++17 standard
- **Comprehensive Testing**: Unit and integration tests
- **Performance Monitoring**: Built-in benchmarking tools

## Build Process

### Prerequisites

1. **C++17 Compiler**: MSVC, GCC, or Clang
2. **CMake 3.16+**: Build system
3. **Git**: Version control
4. **vcpkg**: Package manager (auto-installed)

### Dependencies

- **httplib**: Header-only HTTP client
- **nlohmann/json**: JSON parsing library
- **gRPC**: High-performance RPC framework
- **Protobuf**: Protocol buffer support

### Build Steps

1. **Windows**:
   ```cmd
   build.bat
   ```

2. **Linux/macOS**:
   ```bash
   ./build.sh
   ```

3. **Manual**:
   ```bash
   mkdir build && cd build
   cmake .. -DCMAKE_TOOLCHAIN_FILE=path/to/vcpkg.cmake
   cmake --build . --config Release
   ```

## Usage Examples

### Basic Usage

```bash
# Run demo
./APIBridgeDemo

# With custom config
./APIBridgeDemo --config config.json

# Test mode
./APIBridgeDemo --test
```

### API Examples

```cpp
// REST API
RESTClient client("http://localhost:8080");
auto result = client.registerComponent("user", "data", "context");

// gRPC API
GRPCClient client("localhost:9092");
client.streamBatteryStatus("battery-001", 2, callback);
```

## Integration with RAD Studio

### Project Setup

1. Create new C++ Console Application
2. Add all source files to project
3. Configure include and library paths
4. Set C++17 standard
5. Link required libraries

### Build Configuration

- **Include Paths**: Source directory
- **Library Paths**: vcpkg installed libraries
- **Link Libraries**: System and third-party libraries
- **Preprocessor**: Platform-specific definitions

## Testing Strategy

### Unit Tests

- Individual component testing
- API client validation
- Error handling verification

### Integration Tests

- End-to-end API testing
- Real blockchain interaction
- Performance benchmarking

### Compatibility Tests

- Compiler compatibility
- Platform-specific features
- Dependency validation

## Performance Considerations

### REST vs gRPC

- **REST**: Better for simple operations, easier debugging
- **gRPC**: Better for streaming, lower latency, smaller payloads

### Optimization

- Connection pooling
- Request batching
- Caching strategies
- Memory management

## Security Features

- SSL/TLS support
- Certificate validation
- Input sanitization
- Authentication tokens

## Future Enhancements

- GUI interface
- WebSocket support
- Advanced caching
- Load balancing
- Metrics collection 