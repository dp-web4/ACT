# Web4 Race Car Battery Management API Bridge - Technical Specification

## 📋 Document Information

- **Version**: 1.0.0
- **Date**: July 2025
- **Status**: Development/Testing
- **Author**: Web4 Development Team

## 🎯 Project Overview

The API Bridge is a middleware service that provides REST API endpoints for the racecar-webd blockchain, enabling race car battery management operations through HTTP interfaces.

### Core Objectives
1. **Abstraction Layer**: Hide blockchain complexity from client applications
2. **Standardization**: Provide consistent REST API interface
3. **Performance**: Optimize for real-time battery management operations
4. **Security**: Implement LCT-mediated trust without PKI
5. **Scalability**: Support multiple client applications simultaneously

## 🏗️ System Architecture

### High-Level Architecture
```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Applications                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ Web Dashboard│  │ Mobile App  │  │ IoT Devices │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
                                │ HTTP/REST
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Bridge Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ REST Server  │  │ Auth Layer  │  │ Rate Limiter│            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ Validator    │  │ Logger      │  │ Metrics     │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
                                │ REST/gRPC
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Blockchain Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ racecar-webd │  │ REST API    │  │ gRPC API    │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
```

### Component Architecture
```
api-bridge/
├── cmd/                           # Application entry points
│   └── root.go                   # Main application logic
├── internal/                     # Private application code
│   ├── blockchain/              # Blockchain client implementations
│   │   ├── client.go           # Main blockchain client interface
│   │   └── rest_client.go      # REST client implementation
│   ├── config/                 # Configuration management
│   │   └── config.go          # Configuration structures
│   ├── handlers/               # HTTP request handlers
│   │   └── handlers.go        # All API endpoint handlers
│   └── server/                 # HTTP server implementation
│       └── server.go          # Server setup and routing
├── config.yml                  # Configuration file
├── go.mod                     # Go module file
└── README.md                  # Documentation
```

## 🔌 API Specifications

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
**Current**: None (development mode)  
**Planned**: JWT tokens with role-based access control

### Content Type
```
Content-Type: application/json
```

### Response Format
All responses follow this structure:
```json
{
  "status": "success|error",
  "data": { ... },
  "message": "Optional message",
  "timestamp": 1751725402
}
```

### Error Handling
```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "details": { ... }
}
```

## 📡 Endpoint Specifications

### 1. Component Registry

#### Register Component
```http
POST /api/v1/components/register
```

**Request Body:**
```json
{
  "creator": "string (required)",
  "component_data": "string (required)",
  "context": "string (optional)"
}
```

**Response:**
```json
{
  "component_id": "string",
  "component_identity": "string",
  "lct_id": "string",
  "status": "string"
}
```

#### Get Component
```http
GET /api/v1/components/{id}
```

**Response:**
```json
{
  "component_id": "string",
  "component_type": "string",
  "manufacturer_id": "string",
  "status": "string"
}
```

### 2. LCT Management

#### Create LCT
```http
POST /api/v1/lct/create
```

**Request Body:**
```json
{
  "creator": "string (required)",
  "component_a": "string (required)",
  "component_b": "string (required)",
  "context": "string (required)",
  "proxy_id": "string (optional)"
}
```

**Response:**
```json
{
  "creator": "string",
  "lct_id": "string",
  "status": "string"
}
```

### 3. Trust Tensor Operations

#### Create Trust Tensor
```http
POST /api/v1/trust/tensor
```

**Request Body:**
```json
{
  "creator": "string (required)",
  "component_a": "string (required)",
  "component_b": "string (required)",
  "context": "string (required)",
  "initial_score": "number (required)"
}
```

**Response:**
```json
{
  "tensor_id": "string",
  "score": "number",
  "status": "string"
}
```

### 4. Enhanced Trust Tensor Operations

#### Calculate Relationship Trust
```http
POST /api/v1/trust-enhanced/calculate
```

**Request Body:**
```json
{
  "component_a": "string (required)",
  "component_b": "string (required)",
  "operational_context": "string (required)"
}
```

**Response:**
```json
{
  "tensor_id": "string",
  "component_a": "string",
  "component_b": "string",
  "operational_context": "string",
  "trust_score": "number",
  "status": "string",
  "calculated_at": "number",
  "txhash": "string"
}
```

#### Get Relationship Tensor
```http
GET /api/v1/trust-enhanced/relationship
```

**Query Parameters:**
- `component_a`: string (required)
- `component_b`: string (required)

**Response:**
```json
{
  "tensor_id": "string",
  "component_a": "string",
  "component_b": "string",
  "score": "number",
  "status": "string",
  "created_at": "number",
  "last_updated": "number",
  "evidence_count": "number",
  "learning_rate": "number"
}
```

#### Update Tensor Score
```http
PUT /api/v1/trust-enhanced/score
```

**Request Body:**
```json
{
  "creator": "string (required)",
  "component_a": "string (required)",
  "component_b": "string (required)",
  "score": "number (required)",
  "context": "string (optional)"
}
```

**Response:**
```json
{
  "tensor_id": "string",
  "creator": "string",
  "component_a": "string",
  "component_b": "string",
  "score": "number",
  "context": "string",
  "updated_at": "number",
  "txhash": "string"
}
```

### 5. Queue Management Operations

#### Queue Pairing Request
```http
POST /api/v1/queue/pairing-request
```

**Request Body:**
```json
{
  "component_a": "string (required)",
  "component_b": "string (required)",
  "operational_context": "string (required)",
  "proxy_id": "string (optional)"
}
```

**Response:**
```json
{
  "request_id": "string",
  "component_a": "string",
  "component_b": "string",
  "operational_context": "string",
  "proxy_id": "string",
  "status": "string",
  "created_at": "number",
  "txhash": "string"
}
```

#### Get Queue Status
```http
GET /api/v1/queue/status/{component_id}
```

**Response:**
```json
{
  "component_id": "string",
  "status": "string",
  "pending_requests": "number",
  "processed_requests": "number",
  "last_processed": "number"
}
```

#### Process Offline Queue
```http
POST /api/v1/queue/process-offline/{component_id}
```

**Response:**
```json
{
  "component_id": "string",
  "status": "string",
  "processed_requests": "number",
  "failed_requests": "number",
  "processed_at": "number",
  "txhash": "string"
}
```

#### Cancel Request
```http
DELETE /api/v1/queue/cancel/{request_id}
```

**Request Body:**
```json
{
  "reason": "string (required)"
}
```

**Response:**
```json
{
  "request_id": "string",
  "status": "string",
  "reason": "string",
  "cancelled_at": "number",
  "txhash": "string"
}
```

### 6. Authorization Management Operations

#### Create Pairing Authorization
```http
POST /api/v1/authorization/pairing
```

**Request Body:**
```json
{
  "component_a": "string (required)",
  "component_b": "string (required)",
  "operational_context": "string (required)",
  "authorization_rules": "string (required)"
}
```

**Response:**
```json
{
  "authorization_id": "string",
  "component_a": "string",
  "component_b": "string",
  "operational_context": "string",
  "authorization_rules": "string",
  "status": "string",
  "created_at": "number",
  "txhash": "string"
}
```

#### Get Component Authorizations
```http
GET /api/v1/authorization/component/{component_id}
```

**Response:**
```json
{
  "authorizations": [
    {
      "authorization_id": "string",
      "component_a": "string",
      "component_b": "string",
      "operational_context": "string",
      "authorization_rules": "string",
      "status": "string",
      "created_at": "number"
    }
  ],
  "count": "number",
  "component_id": "string"
}
```

#### Update Authorization
```http
PUT /api/v1/authorization/{authorization_id}
```

**Request Body:**
```json
{
  "operational_context": "string (optional)",
  "authorization_rules": "string (optional)",
  "status": "string (optional)"
}
```

**Response:**
```json
{
  "authorization_id": "string",
  "operational_context": "string",
  "authorization_rules": "string",
  "status": "string",
  "updated_at": "number",
  "txhash": "string"
}
```

#### Revoke Authorization
```http
DELETE /api/v1/authorization/{authorization_id}
```

**Request Body:**
```json
{
  "reason": "string (required)"
}
```

**Response:**
```json
{
  "authorization_id": "string",
  "status": "string",
  "reason": "string",
  "revoked_at": "number",
  "txhash": "string"
}
```

#### Check Pairing Authorization
```http
GET /api/v1/authorization/check
```

**Query Parameters:**
- `component_a`: string (required)
- `component_b`: string (required)
- `operational_context`: string (required)

**Response:**
```json
{
  "authorized": "boolean",
  "authorization_id": "string",
  "component_a": "string",
  "component_b": "string",
  "operational_context": "string",
  "trust_score": "number",
  "verification_status": "string",
  "checked_at": "number"
}
```
  "status": "string"
}
```

### 4. Energy Cycle Management

#### Create Energy Operation
```http
POST /api/v1/energy/operation
```

**Request Body:**
```json
{
  "creator": "string (required)",
  "component_a": "string (required)",
  "component_b": "string (required)",
  "operation_type": "string (required)",
  "amount": "number (required)",
  "context": "string (required)"
}
```

**Response:**
```json
{
  "operation_id": "string",
  "type": "string",
  "amount": "number",
  "status": "string"
}
```

## 🔧 Implementation Details

### Current Implementation Status

#### ✅ Implemented
- REST API server with Gin framework
- Mock blockchain client responses
- Basic request/response handling
- Structured logging with zerolog
- Configuration management
- Health check endpoint

#### 🔄 In Progress
- Real blockchain integration
- Error handling improvements
- Input validation

#### 📋 Planned
- gRPC client implementation
- Authentication & authorization
- Rate limiting
- Cross-platform builds

### Technology Stack

#### Backend
- **Language**: Go 1.24.0+
- **Framework**: Gin (HTTP server)
- **Logging**: zerolog
- **Configuration**: Viper + YAML
- **CLI**: Cobra

#### Blockchain Integration
- **Current**: REST client (mock)
- **Planned**: gRPC client
- **Protocol**: HTTP/2 for gRPC

#### Development Tools
- **Build**: Go modules
- **Testing**: Go testing framework
- **Documentation**: Markdown

## 🚀 Performance Requirements

### Response Times
- **Health Check**: < 100ms
- **Component Operations**: < 500ms
- **Energy Operations**: < 1000ms
- **Trust Operations**: < 500ms

### Throughput
- **Concurrent Requests**: 100+
- **Requests per Second**: 1000+
- **WebSocket Connections**: 1000+

### Scalability
- **Horizontal Scaling**: Load balancer support
- **Vertical Scaling**: Resource optimization
- **Caching**: Redis integration (planned)

## 🔒 Security Requirements

### Current Security
- **None** (development mode)

### Planned Security Features
- **Authentication**: JWT tokens
- **Authorization**: Role-based access control
- **Rate Limiting**: Per-client request limits
- **Input Validation**: Request sanitization
- **CORS**: Cross-origin resource sharing
- **HTTPS**: TLS encryption

### LCT-Mediated Trust
- **No PKI**: Split-key cryptography
- **Context Tokens**: Linked Context Tokens
- **Trust Scoring**: Dynamic trust relationships
- **Offline Support**: Offline operation capability

## 🧪 Testing Strategy

### Unit Testing
- Handler functions
- Blockchain client methods
- Configuration validation
- Utility functions

### Integration Testing
- API endpoint testing
- Blockchain integration
- Error handling scenarios
- Performance testing

### End-to-End Testing
- Complete workflow testing
- Multi-component scenarios
- Real-world use cases

## 📦 Deployment

### Current Deployment
- **Platform**: Windows
- **Method**: Direct binary execution
- **Configuration**: Local config.yml

### Planned Deployment Options

#### Linux Deployment
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/api-bridge ./cmd/api-bridge

# Systemd service
sudo systemctl enable api-bridge
sudo systemctl start api-bridge
```

#### Docker Deployment
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-bridge ./cmd/api-bridge

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-bridge .
CMD ["./api-bridge"]
```

#### ARM Deployment
```bash
# ARM64 build
GOOS=linux GOARCH=arm64 go build -o bin/api-bridge-arm64 ./cmd/api-bridge

# Raspberry Pi deployment
scp api-bridge-arm64 pi@raspberrypi:/home/pi/
```

## 🔮 Future Roadmap

### Phase 1: Production Readiness (Q3 2025)
- [ ] Real blockchain integration
- [ ] Authentication & authorization
- [ ] Input validation & security
- [ ] Error handling improvements

### Phase 2: gRPC Implementation (Q4 2025)
- [ ] gRPC client development
- [ ] Dual transport support
- [ ] gRPC streaming
- [ ] Service discovery

### Phase 3: Cross-Platform (Q1 2026)
- [ ] Linux build support
- [ ] ARM architecture support
- [ ] Docker containerization
- [ ] CI/CD pipeline

### Phase 4: Advanced Features (Q2 2026)
- [ ] Real-time WebSocket streaming
- [ ] Monitoring & observability
- [ ] High availability features
- [ ] Performance optimization

## 📊 Monitoring & Observability

### Metrics (Planned)
- Request rate and latency
- Error rates and types
- Blockchain connection status
- Resource utilization

### Logging
- Structured JSON logging
- Request/response logging
- Error logging with context
- Performance logging

### Health Checks
- Service health endpoint
- Blockchain connectivity check
- Database connectivity (future)
- External service dependencies

## 🔧 Configuration

### Environment Variables
```bash
API_BRIDGE_PORT=8080
API_BRIDGE_LOG_LEVEL=info
BLOCKCHAIN_ENDPOINT=http://localhost:1317
BLOCKCHAIN_TIMEOUT=30
```

### Configuration File (config.yml)
```yaml
server:
  port: 8080
  read_timeout: 30
  write_timeout: 30

blockchain:
  endpoint: "http://localhost:1317"
  timeout: 30

logging:
  level: "info"
```

## 📞 Support & Maintenance

### Issue Tracking
- GitHub Issues for bug reports
- Feature requests through issues
- Documentation updates

### Maintenance Schedule
- **Weekly**: Security updates
- **Monthly**: Feature releases
- **Quarterly**: Major version updates

### Support Channels
- GitHub Issues
- Documentation
- Troubleshooting guides
- Community forums (planned)

---

**Document Version**: 1.0.0  
**Last Updated**: July 2025  
**Next Review**: August 2025 