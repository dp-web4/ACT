# API Bridge - New Go Project Structure

This document explains the new Go project structure that follows Go best practices for maximum portability and maintainability.

## Project Structure

```
api-bridge/
├── cmd/                          # Binary entrypoints
│   ├── api-bridge/              # Main API bridge service
│   │   └── main.go
│   ├── debug-grpc-client/       # Debug gRPC client tool
│   │   └── main.go
│   ├── debug-grpc-server/       # Debug gRPC server tool
│   │   └── main.go
│   ├── test-grpc-client/        # Test gRPC client tool
│   │   └── main.go
│   └── simple-grpc-test/        # Simple gRPC test tool
│       └── main.go
├── internal/                     # Private application code
│   ├── blockchain/              # Blockchain integration
│   ├── config/                  # Configuration management
│   ├── events/                  # Event system
│   ├── grpc/                    # gRPC server implementation
│   ├── handlers/                # HTTP request handlers
│   └── server/                  # Server implementation
├── proto/                       # Protocol buffer definitions
├── bin/                         # Built binaries (created by build)
├── Makefile                     # Build and development commands
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── config.yaml                  # Configuration file
```

## Benefits of This Structure

1. **Portability**: Works on any OS with Go toolchain
2. **No Package Conflicts**: Each binary is isolated
3. **Easy Building**: Clear separation of binaries and libraries
4. **Standard Go Layout**: Follows Go community conventions
5. **Scalable**: Easy to add new binaries or libraries

## Building and Running

### Using Makefile (Recommended)

```bash
# Show available commands
make help

# Build main API bridge
make build

# Build all binaries
make build-all

# Run unit tests
make test

# Run integration tests
make test-integration

# Build and run API bridge
make run

# Run with debug logging
make run-debug

# Clean build artifacts
make clean
```

### Using Go Commands Directly

```bash
# Build main API bridge
go build -o bin/api-bridge ./cmd/api-bridge

# Build debug gRPC client
go build -o bin/debug-grpc-client ./cmd/debug-grpc-client

# Build debug gRPC server
go build -o bin/debug-grpc-server ./cmd/debug-grpc-server

# Build test gRPC client
go build -o bin/test-grpc-client ./cmd/test-grpc-client

# Build simple gRPC test
go build -o bin/simple-grpc-test ./cmd/simple-grpc-test

# Run unit tests
go test -v ./api_bridge_unit_test.go

# Run integration tests
go test -v ./api_bridge_integration_test.go
```

## Running the API Bridge

```bash
# Build and run with default settings
make run

# Or run manually
./bin/api-bridge --rest-port 8082 --grpc-port 9093

# Run with debug logging
./bin/api-bridge --rest-port 8082 --grpc-port 9093 --log-level debug

# Run with custom config file
./bin/api-bridge --config custom-config.yaml --rest-port 8082 --grpc-port 9093
```

## Development Workflow

1. **Make Changes**: Edit code in `internal/` or `cmd/`
2. **Run Tests**: `make test` or `make test-integration`
3. **Build**: `make build` or `make build-all`
4. **Run**: `make run` or `./bin/api-bridge [options]`

## Adding New Binaries

To add a new binary:

1. Create a new directory under `cmd/`:
   ```bash
   mkdir cmd/my-new-tool
   ```

2. Create `main.go` in the new directory:
   ```go
   package main
   
   func main() {
       // Your binary code here
   }
   ```

3. Add build command to Makefile:
   ```makefile
   build-all: build
       @echo "Building my-new-tool..."
       go build -o bin/my-new-tool ./cmd/my-new-tool
   ```

## Configuration

The API bridge uses environment variables and configuration files:

- **Environment Variables**: Override any config setting
- **Config File**: `config.yaml` (default) or specify with `--config`
- **Command Line Flags**: Override specific settings

## Testing

- **Unit Tests**: `make test` - Fast, isolated tests
- **Integration Tests**: `make test-integration` - End-to-end tests with blockchain
- **Manual Testing**: Use the debug tools in `cmd/` for manual testing

## Troubleshooting

### Build Issues
- Ensure you're in the `api-bridge` directory
- Run `go mod tidy` to fix dependencies
- Check that all imports use the correct module path

### Runtime Issues
- Check logs for detailed error messages
- Verify blockchain node is running (for integration tests)
- Ensure ports are not already in use

### Port Conflicts
- Use different ports: `--rest-port 8083 --grpc-port 9094`
- Stop other instances: `pkill -f api-bridge`

## Migration from Old Structure

The old structure had multiple `main` functions in the same package, which caused build conflicts. The new structure:

- ✅ Separates each binary into its own package
- ✅ Follows Go best practices
- ✅ Makes the project portable across systems
- ✅ Simplifies building and deployment
- ✅ Enables easy CI/CD integration 