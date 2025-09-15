# API Bridge Makefile
# Provides convenient commands for building and running the API bridge

.PHONY: build build-all clean test run help

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the main API bridge binary"
	@echo "  build-all   - Build all binaries (api-bridge, debug clients, etc.)"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  run         - Build and run the API bridge"
	@echo "  run-debug   - Run with debug logging"

# Build the main API bridge binary
build:
	@echo "Building API bridge..."
	go build -o bin/api-bridge ./cmd/api-bridge

# Build all binaries
build-all: build
	@echo "Creating organized directory structure..."
	mkdir -p bin/debug bin/test bin/windows
	@echo "Building debug gRPC client..."
	go build -o bin/debug/debug_grpc_client ./cmd/debug-grpc-client
	@echo "Building debug gRPC server..."
	go build -o bin/debug/debug_grpc_server ./cmd/debug-grpc-server
	@echo "Building test gRPC client..."
	go build -o bin/test/test_grpc_client ./cmd/test-grpc-client
	@echo "Building API bridge test binary..."
	go build -o bin/test/api-bridge-test ./cmd/api-bridge
	@echo "Building simple gRPC test..."
	go build -o bin/test/simple-grpc-test ./cmd/simple-grpc-test
	@echo "Building Windows binaries..."
	GOOS=windows GOARCH=amd64 go build -o bin/windows/api-bridge.exe ./cmd/api-bridge
	GOOS=windows GOARCH=amd64 go build -o bin/windows/debug_grpc_server.exe ./cmd/debug-grpc-server
	@echo "All binaries built successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Build artifacts cleaned!"

# Run unit tests
test:
	@echo "Running unit tests..."
	go test -v ./api_bridge_unit_test.go

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./api_bridge_integration_test.go

# Build and run the API bridge
run: build
	@echo "Starting API bridge..."
	./bin/api-bridge --rest-port 8082 --grpc-port 9093

# Run with debug logging
run-debug: build
	@echo "Starting API bridge with debug logging..."
	./bin/api-bridge --rest-port 8082 --grpc-port 9093 --log-level debug

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Show project structure
tree:
	@echo "Project structure:"
	@tree -I 'bin|.git|node_modules' -a 