# Web4 Race Car Battery Management API Bridge - Quick Reference

## 🚀 Quick Start

### Build & Run (Windows)
```bash
cd api-bridge

# Option 1: Using Makefile (Recommended)
make build
./bin/api-bridge.exe

# Option 2: Direct Go build
go build -o bin/api-bridge.exe ./cmd/api-bridge
./bin/api-bridge.exe
```

### Test Health
```bash
curl http://localhost:8080/health
```

### Running with Custom Ports
You can set the REST and gRPC ports at launch:

```bash
# Using Makefile
make run

# Or manually
./bin/api-bridge --rest-port 8080 --grpc-port 9092
```

- `--rest-port` (or `-p`): REST API port (default: 8080)
- `--grpc-port` (or `-g`): gRPC API port (default: 9090)

### Testing Blockchain Integration
To validate your deployment, use the provided test script:

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Or use the legacy test script
./test_blockchain_fix.sh
```

This script checks API bridge health, blockchain connectivity, and component registration.

## 📡 API Endpoints Quick Reference

### Component Registry
```bash
# Register component
curl -X POST http://localhost:8080/api/v1/components/register \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_data": "battery_pack_v1", "context": "test"}'

# Get component
curl http://localhost:8080/api/v1/components/comp_1751725111

# Verify component
curl -X POST http://localhost:8080/api/v1/components/comp_1751725111/verify \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice"}'
```

### LCT Management
```bash
# Create LCT
curl -X POST http://localhost:8080/api/v1/lct/create \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_a": "comp_1751725111", "component_b": "comp_0987654321", "context": "battery_management"}'

# Get LCT
curl http://localhost:8080/api/v1/lct/lct_comp_1751725111
```

### Trust Tensor
```bash
# Create trust tensor
curl -X POST http://localhost:8080/api/v1/trust/tensor \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_a": "comp_1751725111", "component_b": "comp_0987654321", "context": "battery_management", "initial_score": 0.85}'

# Get trust tensor
curl http://localhost:8080/api/v1/trust/tensor/tensor_comp_1751725111_comp_0987654321
```

### Energy Operations
```bash
# Create energy operation
curl -X POST http://localhost:8080/api/v1/energy/operation \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_a": "comp_1751725111", "component_b": "comp_0987654321", "operation_type": "charge", "amount": 50.0, "context": "battery_management"}'

# Execute energy transfer
curl -X POST http://localhost:8080/api/v1/energy/transfer \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "operation_id": "op_1751725363"}'

# Get energy balance
curl http://localhost:8080/api/v1/energy/balance/comp_1751725111
```

### Pairing
```bash
# Initiate pairing
curl -X POST http://localhost:8080/api/v1/pairing/initiate \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_a": "comp_1751725111", "component_b": "comp_0987654321", "operational_context": "battery_management"}'

# Complete pairing
curl -X POST http://localhost:8080/api/v1/pairing/complete \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "challenge_id": "challenge_1751725402", "response": "response_data"}'

# Get pairing status
curl http://localhost:8080/api/v1/pairing/status/challenge_1751725402
```

### Queue Management
```bash
# Queue pairing request
curl -X POST http://localhost:8080/api/v1/queue/pairing-request \
  -H "Content-Type: application/json" \
  -d '{"component_a": "battery_001", "component_b": "battery_002", "operational_context": "energy_transfer", "proxy_id": "proxy_001"}'

# Get queue status
curl http://localhost:8080/api/v1/queue/status/battery_001

# Process offline queue
curl -X POST http://localhost:8080/api/v1/queue/process-offline/battery_001

# Cancel request
curl -X DELETE http://localhost:8080/api/v1/queue/cancel/queue_battery_001_battery_002_1704067200 \
  -H "Content-Type: application/json" \
  -d '{"reason": "component_offline"}'

# Get queued requests
curl http://localhost:8080/api/v1/queue/requests/battery_001

# List proxy queue
curl http://localhost:8080/api/v1/queue/proxy/proxy_001
```

### Authorization Management
```bash
# Create pairing authorization
curl -X POST http://localhost:8080/api/v1/authorization/pairing \
  -H "Content-Type: application/json" \
  -d '{"component_a": "battery_001", "component_b": "battery_002", "operational_context": "energy_transfer", "authorization_rules": "trust_score > 0.7"}'

# Get component authorizations
curl http://localhost:8080/api/v1/authorization/component/battery_001

# Update authorization
curl -X PUT http://localhost:8080/api/v1/authorization/auth_001 \
  -H "Content-Type: application/json" \
  -d '{"operational_context": "high_priority_energy_transfer", "authorization_rules": "trust_score > 0.8"}'

# Revoke authorization
curl -X DELETE http://localhost:8080/api/v1/authorization/auth_001 \
  -H "Content-Type: application/json" \
  -d '{"reason": "security_concern"}'

# Check pairing authorization
curl "http://localhost:8080/api/v1/authorization/check?component_a=battery_001&component_b=battery_002&operational_context=energy_transfer"
```

### Enhanced Trust Tensor
```bash
# Calculate relationship trust
curl -X POST http://localhost:8080/api/v1/trust-enhanced/calculate \
  -H "Content-Type: application/json" \
  -d '{"component_a": "battery_001", "component_b": "battery_002", "operational_context": "energy_transfer"}'

# Get relationship tensor
curl "http://localhost:8080/api/v1/trust-enhanced/relationship?component_a=battery_001&component_b=battery_002"

# Update tensor score
curl -X PUT http://localhost:8080/api/v1/trust-enhanced/score \
  -H "Content-Type: application/json" \
  -d '{"creator": "alice", "component_a": "battery_001", "component_b": "battery_002", "score": 0.85, "context": "successful_energy_transfer"}'
```

## ⚙️ Configuration

### config.yml
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

### Environment Variables
```bash
export API_BRIDGE_PORT=8080
export API_BRIDGE_LOG_LEVEL=info
export BLOCKCHAIN_ENDPOINT=http://localhost:1317
export BLOCKCHAIN_TIMEOUT=30
```

### Event System Configuration
```yaml
events:
  enabled: true
  max_retries: 3
  retry_delay: 5
  queue_size: 1000
  endpoints:
    component_registered:
      - "http://localhost:3000/webhooks/component-registered"
    energy_transfer:
      - "http://sql-audit-service:8080/events"
```

**Test Event System:**
```bash
# Windows
test_event_system.bat

# Linux/macOS
./test_event_system.sh
```

## 🔧 Development Commands

### Build Commands
```bash
# Build for current platform
go build -o api-bridge .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o api-bridge.exe .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o api-bridge .

# Build for ARM64
GOOS=linux GOARCH=arm64 go build -o api-bridge-arm64 .
```

### Test Commands
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/handlers

# Run benchmarks
go test -bench=. ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for security issues
gosec ./...

# Update dependencies
go mod tidy
go mod verify
```

## 🐛 Troubleshooting

### Common Issues

#### Service Won't Start
```bash
# Check if port is in use
netstat -an | grep :8080

# Check logs
./api-bridge.exe --verbose

# Check configuration
./api-bridge.exe --validate-config
```

#### Blockchain Connection Issues
```bash
# Test blockchain connectivity
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info

# Check blockchain status
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/status
```

#### Build Errors
```bash
# Clean and rebuild
go clean
go mod tidy
go build -o api-bridge .

# Check Go version
go version

# Update Go if needed
# Download from https://golang.org/dl/
```

### Log Levels
```bash
# Debug logging
./api-bridge.exe --log-level=debug

# Info logging (default)
./api-bridge.exe --log-level=info

# Warning only
./api-bridge.exe --log-level=warn

# Error only
./api-bridge.exe --log-level=error
```

## 📊 Monitoring

### Health Check
```bash
# Basic health check
curl http://localhost:8080/health

# Detailed health check
curl http://localhost:8080/health/detailed
```

### Metrics (Planned)
```bash
# Prometheus metrics
curl http://localhost:8080/metrics

# Application metrics
curl http://localhost:8080/api/v1/metrics
```

### Logs
```bash
# View logs in real-time
tail -f /var/log/api-bridge/app.log

# Search logs
grep "ERROR" /var/log/api-bridge/app.log

# Log rotation
logrotate /etc/logrotate.d/api-bridge
```

## 🔒 Security

### Authentication (Planned)
```bash
# Get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'

# Use token
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/components/comp_1751725111
```

### Rate Limiting (Planned)
```bash
# Check rate limit headers
curl -I http://localhost:8080/api/v1/components/register

# Response headers:
# X-RateLimit-Limit: 1000
# X-RateLimit-Remaining: 999
# X-RateLimit-Reset: 1751725402
```

## 🚀 Deployment

### Windows Service
```bash
# Install as Windows service
sc create "Web4APIBridge" binPath="C:\path\to\api-bridge.exe"
sc start "Web4APIBridge"
sc stop "Web4APIBridge"
```

### Linux Systemd
```bash
# Enable service
sudo systemctl enable api-bridge
sudo systemctl start api-bridge
sudo systemctl status api-bridge

# View logs
sudo journalctl -u api-bridge -f
```

### Docker
```bash
# Build image
docker build -t web4-api-bridge:latest .

# Run container
docker run -d -p 8080:8080 --name api-bridge web4-api-bridge:latest

# View logs
docker logs -f api-bridge
```

## 📋 Response Codes

### HTTP Status Codes
- `200 OK`: Success
- `201 Created`: Resource created
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Error Response Format
```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "details": {
    "field": "validation error"
  },
  "timestamp": 1751725402
}
```

## 🔗 Related Documentation

- [README.md](./README.md) - Complete documentation
- [TECHNICAL_SPEC.md](./TECHNICAL_SPEC.md) - Technical specifications
- [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) - Deployment instructions
- [DEVELOPMENT_ROADMAP.md](./DEVELOPMENT_ROADMAP.md) - Development roadmap

## 📞 Support

### Getting Help
1. Check this quick reference
2. Review troubleshooting section
3. Check logs for errors
4. Search GitHub issues
5. Contact development team

### Useful Links
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Cosmos SDK](https://docs.cosmos.network/)
- [gRPC Documentation](https://grpc.io/docs/)

---

**Last Updated**: July 2025  
**Version**: 1.0.0 