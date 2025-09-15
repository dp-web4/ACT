# API Bridge Binaries

This directory contains all compiled binaries for the API Bridge service, organized by purpose and platform.

## 📁 Directory Structure

```
bin/
├── api-bridge              # Main API Bridge binary (Linux)
├── debug/                  # Debug and development binaries
│   ├── debug_grpc_client   # gRPC client for debugging
│   └── debug_grpc_server   # gRPC server for debugging
├── test/                   # Test binaries
│   ├── api-bridge-test     # API Bridge test binary
│   └── test_grpc_client    # gRPC client for testing
└── windows/                # Windows platform binaries
    ├── api-bridge.exe      # Main API Bridge binary (Windows)
    └── debug_grpc_server.exe # Debug gRPC server (Windows)
```

## 🚀 Usage

### Main API Bridge
```bash
# Linux/macOS
./bin/api-bridge --grpc-port 9092 --rest-port 8080

# Windows
./bin/windows/api-bridge.exe --grpc-port 9092 --rest-port 8080
```

### Test Binaries
```bash
# Run API Bridge tests
./bin/test/api-bridge-test

# Test gRPC client
./bin/test/test_grpc_client
```

### Debug Binaries
```bash
# Debug gRPC client
./bin/debug/debug_grpc_client

# Debug gRPC server
./bin/debug/debug_grpc_server
```

## 🔧 Building

To rebuild all binaries:

```bash
# Build main binary
go build -o bin/api-bridge cmd/api-bridge/main.go

# Build test binary
go build -o bin/test/api-bridge-test cmd/api-bridge/main.go

# Build debug binaries
go build -o bin/debug/debug_grpc_client cmd/debug-grpc-client/main.go
go build -o bin/debug/debug_grpc_server cmd/debug-grpc-server/main.go

# Build test client
go build -o bin/test/test_grpc_client cmd/test-grpc-client/main.go

# Build Windows binaries (on Windows or with cross-compilation)
GOOS=windows GOARCH=amd64 go build -o bin/windows/api-bridge.exe cmd/api-bridge/main.go
GOOS=windows GOARCH=amd64 go build -o bin/windows/debug_grpc_server.exe cmd/debug-grpc-server/main.go
```

## 📋 Binary Descriptions

### Main Binaries
- **api-bridge**: Main API Bridge service with REST and gRPC endpoints
- **api-bridge.exe**: Windows version of the main API Bridge service

### Test Binaries
- **api-bridge-test**: Test binary for API Bridge functionality
- **test_grpc_client**: gRPC client for testing API Bridge endpoints

### Debug Binaries
- **debug_grpc_client**: Debug gRPC client for development and troubleshooting
- **debug_grpc_server**: Debug gRPC server for development and troubleshooting
- **debug_grpc_server.exe**: Windows version of debug gRPC server

## 🧹 Cleanup

This organization was created to clean up duplicate binaries that were previously scattered in the root directory. All binaries are now properly organized by purpose and platform.

## 🔍 Verification

To verify all binaries are working:

```bash
# Test main binary
./bin/api-bridge --help

# Test Windows binary (if on Windows)
./bin/windows/api-bridge.exe --help

# Test debug binaries
./bin/debug/debug_grpc_client --help
./bin/debug/debug_grpc_server --help
``` 