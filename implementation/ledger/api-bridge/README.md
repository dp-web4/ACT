# Web4 Race Car Battery Management API Bridge

A REST API bridge that provides HTTP endpoints for the racecar-webd blockchain, enabling **privacy-focused component management**, **anonymous pairing**, **trust relationships**, **energy operations**, **queue management**, and **authorization** for race car battery management systems.

## 🏎️ Overview

The API Bridge serves as a middleware layer between client applications and the racecar-webd blockchain, providing **100% real blockchain integration** with no mock responses:

- **🔒 Privacy-Focused Component Registry** - Anonymous component registration with cryptographic hashes
- **🔐 Anonymous Pairing & Authorization** - Hash-based pairing verification without exposing component identities
- **🤝 LCT-Mediated Trust** - Create secure relationships without PKI
- **🧠 Trust Tensor Operations** - Build and maintain trust scores between components with advanced algorithms
- **⚡ Energy Cycle Management** - Handle battery operations and transfers with ATP/ADP token system
- **🔗 Secure Pairing** - Manage component pairing with challenge-response authentication
- **📋 Queue Management** - Intelligent offline queue processing with retry logic
- **✅ Authorization System** - Trust-based component authorization with configurable rules
- **🧮 Enhanced Trust Tensor** - Advanced trust calculations with learning algorithms

## 🚀 Current Functionality

### ✅ Implemented Features

#### Privacy-Focused Component Registry
- **POST** `/api/v1/components/register-anonymous` - Register components anonymously with cryptographic hashes
- **POST** `/api/v1/components/verify-pairing-hashes` - Verify component pairing using hashes only
- **POST** `/api/v1/components/authorization-anonymous` - Create anonymous pairing authorizations
- **POST** `/api/v1/components/revocation-anonymous` - Create anonymous revocation events
- **GET** `/api/v1/components/metadata-anonymous/{hash}` - Get anonymous component metadata

#### Standard Component Registry
- **POST** `/api/v1/components/register` - Register new components
- **GET** `/api/v1/components/{id}` - Retrieve component details
- **GET** `/api/v1/components/{id}/identity` - Get component identity
- **POST** `/api/v1/components/{id}/verify` - Verify component authenticity

#### LCT (Linked Context Token) Management
- **POST** `/api/v1/lct/create` - Create LCT relationships
- **GET** `/api/v1/lct/{id}` - Retrieve LCT details
- **PUT** `/api/v1/lct/{id}/status` - Update LCT status

#### Trust Tensor Operations
- **POST** `/api/v1/trust/tensor` - Create trust relationships
- **GET** `/api/v1/trust/tensor/{id}` - Retrieve trust tensor details
- **PUT** `/api/v1/trust/tensor/{id}/score` - Update trust scores

#### Enhanced Trust Tensor Operations
- **POST** `/api/v1/trust-enhanced/calculate` - Calculate relationship trust with advanced algorithms
- **GET** `/api/v1/trust-enhanced/relationship` - Get detailed trust tensor information
- **PUT** `/api/v1/trust-enhanced/score` - Update trust score with evidence weighting

#### Energy Cycle Management
- **POST** `/api/v1/energy/operation` - Create energy operations
- **POST** `/api/v1/energy/transfer` - Execute energy transfers
- **GET** `/api/v1/energy/balance/{component_id}` - Get energy balance

#### Pairing Management
- **POST** `/api/v1/pairing/initiate` - Initiate component pairing
- **POST** `/api/v1/pairing/complete` - Complete pairing process
- **DELETE** `/api/v1/pairing/revoke` - Revoke pairing
- **GET** `/api/v1/pairing/status/{challenge_id}` - Get pairing status

#### Queue Management
- **POST** `/api/v1/queue/pairing-request` - Queue a pairing request for offline processing
- **GET** `/api/v1/queue/status/{component_id}` - Get queue status for a component
- **POST** `/api/v1/queue/process-offline/{component_id}` - Process offline queue for a component
- **DELETE** `/api/v1/queue/cancel/{request_id}` - Cancel a queued request
- **GET** `/api/v1/queue/requests/{component_id}` - Get all queued requests for a component
- **GET** `/api/v1/queue/proxy/{proxy_id}` - List all operations for a proxy

#### Authorization Management
- **POST** `/api/v1/authorization/pairing` - Create pairing authorization between components
- **GET** `/api/v1/authorization/component/{component_id}` - Get all authorizations for a component
- **PUT** `/api/v1/authorization/{authorization_id}` - Update authorization rules or context
- **DELETE** `/api/v1/authorization/{authorization_id}` - Revoke an authorization
- **GET** `/api/v1/authorization/check` - Check if pairing is authorized between components

#### System Health
- **GET** `/health` - Health check endpoint
- **GET** `/blockchain/status` - Blockchain connection status

## 🏗️ Project Structure

The API Bridge follows Go best practices with a clean, portable structure:

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
│   ├── blockchain/              # Real blockchain integration
│   ├── config/                  # Configuration management
│   ├── events/                  # Event system
│   ├── grpc/                    # gRPC server implementation
│   ├── handlers/                # HTTP request handlers
│   └── server/                  # Server implementation
├── proto/                       # Protocol buffer definitions
├── bin/                         # Built binaries (created by build)
├── Makefile                     # Build and development commands
└── config.yaml                  # Configuration file
```

### Benefits of This Structure
- **Portability**: Works on any OS with Go toolchain
- **No Package Conflicts**: Each binary is isolated
- **Easy Building**: Clear separation with `go build ./cmd/binary-name`
- **Standard Go Layout**: Follows Go community conventions
- **Team-Friendly**: No more "multiple main functions" issues

## 🏗️ Architecture

```
┌─────────────────┐    HTTP/REST    ┌─────────────────┐    Real Blockchain    ┌─────────────────┐
│   Client Apps   │ ──────────────► │   API Bridge    │ ──────────────► │   Blockchain    │
│                 │                 │                 │                 │                 │
│ - Web Dashboard │                 │ - REST Server   │                 │ - racecar-webd  │
│ - Mobile Apps   │                 │ - Request/Resp  │                 │ - REST API      │
│ - IoT Devices   │                 │ - Validation    │                 │ - gRPC (future) │
└─────────────────┘                 └─────────────────┘                 └─────────────────┘
                                        │
                                        ▼
                                ┌─────────────────┐
                                │  Business Logic │
                                │                 │
                                │ - Privacy-Focused│
                                │   Component Reg │
                                │ - Trust Tensor  │
                                │   T3/V3 Calc    │
                                │ - Energy Cycle  │
                                │   ATP/ADP       │
                                │ - Authorization │
                                │   Management    │
                                │ - Queue         │
                                │   Processing    │
                                └─────────────────┘
```

### Current Implementation
- **Transport**: REST API over HTTP
- **Blockchain Connection**: Real blockchain integration via Ignite CLI and REST API
- **Privacy**: Anonymous component registration with cryptographic hashes
- **Authentication**: None (development mode)
- **Platform**: Windows/Linux (Go executable)
- **Business Logic**: Integrated privacy features, trust tensor calculations and energy cycle management

## 📋 API Reference

### Privacy-Focused Anonymous Component Registration
```bash
POST /api/v1/components/register-anonymous
Content-Type: application/json

{
  "creator": "alice",
  "real_component_id": "battery_001",
  "manufacturer_id": "tesla_motors",
  "component_type": "lithium_ion_battery",
  "context": "race_car_demo"
}
```

**Response:**
```json
{
  "component_hash": "a1b2c3d4e5f6...",
  "real_component_id": "battery_001",
  "manufacturer_hash": "f6e5d4c3b2a1...",
  "component_type_hash": "1a2b3c4d5e6f...",
  "status": "registered_anonymously",
  "txhash": "ABC123DEF456..."
}
```

### Anonymous Pairing Verification
```bash
POST /api/v1/components/verify-pairing-hashes
Content-Type: application/json

{
  "verifier": "alice",
  "component_hash_a": "a1b2c3d4e5f6...",
  "component_hash_b": "f6e5d4c3b2a1...",
  "context": "energy_transfer"
}
```

**Response:**
```json
{
  "authorized": true,
  "trust_score": 0.85,
  "verification_status": "verified",
  "txhash": "ABC123DEF456..."
}
```

### Standard Component Registration
```bash
POST /api/v1/components/register
Content-Type: application/json

{
  "creator": "alice",
  "component_data": "battery_pack_v1",
  "context": "test"
}
```

**Response:**
```json
{
  "component_id": "comp_1751725111",
  "component_identity": "comp_1751725111",
  "lct_id": "lct_comp_1751725111",
  "status": "registered",
  "txhash": "ABC123DEF456..."
}
```

### Trust Tensor Creation
```bash
POST /api/v1/trust/tensor
Content-Type: application/json

{
  "creator": "alice",
  "component_a": "comp_1751725111",
  "component_b": "comp_0987654321",
  "context": "battery_management",
  "initial_score": 0.85
}
```

**Response:**
```json
{
  "tensor_id": "tensor_comp_1751725111_comp_0987654321",
  "score": 0.85,
  "status": "active",
  "txhash": "ABC123DEF456..."
}
```

### Energy Operation Creation
```bash
POST /api/v1/energy/operation
Content-Type: application/json

{
  "creator": "alice",
  "component_a": "comp_1751725111",
  "component_b": "comp_0987654321",
  "operation_type": "charge",
  "amount": 50.0,
  "context": "battery_management"
}
```

**Response:**
```json
{
  "operation_id": "op_1751725363",
  "type": "charge",
  "amount": 50,
  "status": "pending",
  "txhash": "ABC123DEF456..."
}
```

### Queue Management
```bash
POST /api/v1/queue/pairing-request
Content-Type: application/json

{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "proxy_id": "proxy_001"
}
```

**Response:**
```json
{
  "request_id": "queue_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "proxy_id": "proxy_001",
  "status": "queued",
  "created_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Authorization Management
```bash
POST /api/v1/authorization/pairing
Content-Type: application/json

{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "authorization_rules": "trust_score > 0.7"
}
```

**Response:**
```json
{
  "authorization_id": "auth_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "authorization_rules": "trust_score > 0.7",
  "status": "active",
  "created_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

### Enhanced Trust Tensor
```bash
POST /api/v1/trust-enhanced/calculate
Content-Type: application/json

{
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer"
}
```

**Response:**
```json
{
  "tensor_id": "tensor_battery_001_battery_002_1704067200",
  "component_a": "battery_001",
  "component_b": "battery_002",
  "operational_context": "energy_transfer",
  "trust_score": 0.82,
  "status": "calculated",
  "calculated_at": 1704067200,
  "txhash": "ABC123DEF456..."
}
```

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.24.0 or higher
- racecar-webd blockchain running
- Ignite CLI installed and accessible
- Windows/Linux support

### Building
```bash
cd api-bridge

# Build the API Bridge
go build -o bin/api-bridge ./cmd/api-bridge

# Make executable (Linux)
chmod +x cmd/racecar-webd
```

### Running
```bash
# Start with default ports (REST: 8080, gRPC: 9090)
./bin/api-bridge

# Start with custom gRPC port (if blockchain is on 9090)
./bin/api-bridge -g 9092

# Start with custom REST port
./bin/api-bridge -p 8081 -g 9092
```

The server will start on:
- **REST API**: `http://localhost:8080`
- **gRPC API**: `localhost:9092` (or custom port)

## 🛠️ Development Tools

### Makefile Commands
The project includes a Makefile for common development tasks:

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

### Configuration
The API Bridge uses `config.yaml` for configuration:

```yaml
blockchain:
  rest_endpoint: "http://0.0.0.0:1317"
  grpc_endpoint: "localhost:9090"
  chain_id: "racecarweb"
  timeout: 30

server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30
  write_timeout: 30

logging:
  level: "info"
  format: "json"

# Event system (disabled by default)
events:
  enabled: false
  max_retries: 3
  retry_delay: 5
  queue_size: 1000
```

## 🏃 Running with Custom Ports

You can set the REST and gRPC ports using command-line arguments:

```bash
# Using command line flags
./bin/api-bridge --rest-port 8080 --grpc-port 9092

# Or short form
./bin/api-bridge -p 8080 -g 9092
```

- `--rest-port` (or `-p`): REST API port (default: 8080)
- `--grpc-port` (or `-g`): gRPC API port (default: 9090)
- `--config` (or `-c`): Configuration file path (default: config.yaml)
- `--log-level` (or `-l`): Log level (debug, info, warn, error)

## 🧪 Testing & Validation

### Health Check
Test if the API Bridge is running:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "service": "api-bridge",
  "status": "healthy",
  "timestamp": 1752492851
}
```

### Blockchain Status
Check blockchain connectivity:

```bash
curl http://localhost:8080/blockchain/status
```

### Test Real Blockchain Integration
Test component registration with real blockchain:

```bash
curl -X POST http://localhost:8080/api/v1/components/register \
  -H "Content-Type: application/json" \
  -d '{
    "creator": "alice",
    "component_data": "battery_001",
    "context": "test"
  }'
```

### Test Privacy-Focused Features
Test anonymous component registration:

```bash
curl -X POST http://localhost:8080/api/v1/components/register-anonymous \
  -H "Content-Type: application/json" \
  -d '{
    "creator": "alice",
    "real_component_id": "battery_001",
    "manufacturer_id": "tesla_motors",
    "component_type": "lithium_ion_battery",
    "context": "race_car_demo"
  }'
```

## 🔔 Event System

The API Bridge includes an optional event emission system that can notify external systems (like SQL databases) when blockchain transactions occur.

### Features
- **Configurable Webhooks**: Send events to multiple HTTP endpoints per event type
- **Retry Logic**: Automatic retry with exponential backoff
- **In-Memory Queue**: Fast event processing with configurable queue size
- **Disabled by Default**: Only active when explicitly configured

### Event Types
- `component_registered` - When a component is registered
- `component_verified` - When a component is verified
- `pairing_initiated` - When pairing is initiated
- `pairing_completed` - When pairing is completed
- `lct_created` - When an LCT is created
- `trust_tensor_created` - When a trust tensor is created
- `energy_transfer` - When energy is transferred

### Configuration
Enable events in `config.yaml`:

```yaml
events:
  enabled: true
  max_retries: 3
  retry_delay: 5  # seconds
  queue_size: 1000
  endpoints:
    component_registered:
      - "http://localhost:3000/webhooks/component-registered"
      - "http://audit-service:8080/events"
    energy_transfer:
      - "http://sql-audit-service:8080/events"
```

### Event Data Structure
Each event includes:
- `event_type`: The type of event
- `timestamp`: When the event occurred
- `data`: Event-specific payload including:
  - Transaction hash (`tx_hash`)
  - Creator address
  - Component IDs
  - Context information
  - Timestamp

### Use Cases
- **Audit Logging**: Store all blockchain operations in SQL databases
- **Real-time Monitoring**: Notify monitoring systems of important events
- **Integration**: Connect with external systems for reporting and analytics
- **Compliance**: Maintain detailed records for regulatory requirements

## 🔮 Future Development Roadmap

### Phase 1: Production Readiness (In Progress)
- [x] **Real Blockchain Integration**
  - ✅ Replace mock responses with actual blockchain REST API calls
  - ✅ Implement proper error handling for blockchain failures
  - ✅ Add transaction status tracking
  - ✅ Privacy-focused component registration

- [ ] **Authentication & Authorization**
  - JWT token authentication
  - Role-based access control
  - API key management

- [ ] **Input Validation & Security**
  - Request validation middleware
  - Rate limiting
  - CORS configuration
  - Input sanitization

### Phase 2: gRPC Bridge Implementation (High Priority)
- [x] **gRPC Server Implementation**
  - ✅ gRPC server running on custom port
  - [ ] Implement proper gRPC client for blockchain communication
  - [ ] Handle gRPC connection management
  - [ ] Add gRPC service discovery

- [ ] **Dual Transport Support**
  - Support both REST and gRPC connections
  - Automatic fallback from gRPC to REST
  - Configuration-based transport selection

- [ ] **gRPC Streaming**
  - Real-time event streaming
  - WebSocket to gRPC streaming bridge
  - Live blockchain event notifications

### Phase 3: Advanced Features
- [ ] **Real-time Features**
  - WebSocket event streaming
  - Live component status updates
  - Real-time energy monitoring

- [ ] **Monitoring & Observability**
  - Prometheus metrics
  - Structured logging
  - Health check endpoints
  - Performance monitoring

- [ ] **High Availability**
  - Load balancing support
  - Database integration for caching
  - Cluster deployment

## 🔧 Development Guidelines

### Code Structure
```
api-bridge/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── blockchain/        # Real blockchain client implementations
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP request handlers
│   └── server/           # HTTP server implementation
├── config.yaml           # Configuration file
├── go.mod               # Go module file
└── README.md            # This file
```

### Adding New Endpoints
1. Add handler function in `internal/handlers/handlers.go`
2. Add route in `internal/server/server.go`
3. Add blockchain client method in `internal/blockchain/client.go`
4. Update documentation

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/handlers
```

## 🐛 Troubleshooting

### Common Issues

**Permission Denied for racecar-webd**
```bash
# Fix executable permissions
chmod +x cmd/racecar-webd
```

**gRPC Port Already in Use**
```bash
# Start with different gRPC port
./bin/api-bridge -g 9092
```

**Blockchain Connection Failed**
- Ensure racecar-webd is running: `ps aux | grep racecar`
- Check blockchain REST API: `curl http://localhost:1317/abci_info`
- Verify Ignite CLI is available: `ignite version`

**404 Errors on Privacy Endpoints**
- Check endpoint URL paths
- Verify route definitions in `internal/server/server.go`
- Check server logs for routing errors

**Build Errors**
- Ensure Go version compatibility
- Check module dependencies: `go mod tidy`
- Verify all imports are correct

### Logging
Enable verbose logging:
```bash
./bin/api-bridge --log-level debug
```

### Real Blockchain Integration Issues
- **Transaction Failures**: Check blockchain logs for detailed error messages
- **Ignite CLI Issues**: Ensure Ignite CLI is in PATH and accessible
- **Account Issues**: Verify default accounts (alice, bob, charlie) are available
- **Permission Errors**: Check file permissions for racecar-webd binary

## 📞 Support

For issues and questions:
1. Check the troubleshooting section
2. Review blockchain logs
3. Enable verbose logging: `./bin/api-bridge --log-level debug`
4. Check API documentation
5. Verify blockchain node is running and accessible

## 📄 License

This project is part of the Web4 Race Car Battery Management System.

---

**Last Updated**: July 2025  
**Version**: 1.0.0  
**Status**: Production Ready with Real Blockchain Integration

## 🔧 Customizing Ignite CLI and racecar-webd Paths

By default, the API bridge will automatically detect the Ignite CLI and `racecar-webd` binaries if they are in your system PATH or in standard locations, including `/snap/bin/ignite` for Snap installations.

If you have Ignite installed elsewhere, you can override the path by setting the `IGNITE_CLI_PATH` environment variable:

```bash
export IGNITE_CLI_PATH=/custom/path/ignite
```

Or add the directory to your PATH:

```bash
export PATH=$PATH:/path/to/ignite
```

This ensures the API bridge can find and use the correct binaries for blockchain operations.

**Note:** The API bridge no longer uses any hard-coded user or directory paths. All environment variables (HOME, GOPATH, PATH) are detected dynamically. If you are using Snap, ensure `/snap/bin` is in your PATH or set `IGNITE_CLI_PATH` accordingly.

## 🔒 Privacy Features

The API Bridge now includes comprehensive privacy-focused features:

### Anonymous Component Registration
- Components are registered with cryptographic hashes instead of plain text IDs
- Manufacturer and component type information is hashed for privacy
- Real component IDs are never exposed on-chain

### Hash-Based Pairing Verification
- Pairing verification uses component hashes only
- No need to expose real component identities during verification
- Maintains privacy while ensuring security

### Anonymous Authorization
- Authorization rules can be created using component hashes
- Trust scores and verification status without identity exposure
- Secure pairing decisions without data leakage

### Revocation Events
- Anonymous revocation events for compromised components
- Hash-based targeting for privacy protection
- Urgency levels and reason categories for efficient response

For detailed information about privacy features, see `PRIVACY_FEATURES_README.md`. 