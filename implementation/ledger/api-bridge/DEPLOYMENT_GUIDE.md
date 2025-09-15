# Web4 Race Car Battery Management API Bridge - Deployment Guide

## 📋 Overview

This guide covers deployment of the API Bridge across different platforms and environments, from development to production. **The API Bridge now follows Go best practices for maximum portability and maintainability. All build, run, and test operations should use the provided Makefile.**

---

## 🚨 Important: New Project Structure & Workflow

- All binaries are now under `cmd/` (e.g., `cmd/api-bridge/`, `cmd/debug-grpc-client/`)
- The `bin/` directory contains all built binaries
- Use the Makefile for all build, run, and test operations
- **Mock fallback logic has been removed**: the API bridge requires a real blockchain and fails clearly if unavailable
- All configuration should use `config.yaml` (not `config.yml`)

### Example Project Structure
```
api-bridge/
├── cmd/                          # Binary entrypoints
│   ├── api-bridge/              # Main API bridge service
│   ├── debug-grpc-client/       # Debug gRPC client tool
│   ├── debug-grpc-server/       # Debug gRPC server tool
│   ├── test-grpc-client/        # Test gRPC client tool
│   └── simple-grpc-test/        # Simple gRPC test tool
├── internal/                     # Private application code
├── proto/                       # Protocol buffer definitions
├── bin/                         # Built binaries (created by build)
├── Makefile                     # Build and development commands
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── config.yaml                  # Configuration file
```

---

## 🛠️ Makefile Commands

The Makefile provides all necessary commands for development and deployment:

```bash
make help             # Show all available commands
make build            # Build main API bridge
make build-all        # Build all binaries
make test             # Run unit tests
make test-integration # Run integration tests
make run              # Build and run API bridge
make run-debug        # Run with debug logging
make clean            # Clean build artifacts
```

---

## 🚀 Deployment Options

### 1. Development Deployment (Windows)
- **Platform**: Windows 10/11
- **Method**: Direct binary execution
- **Use Case**: Development and testing

#### Prerequisites
- Go 1.24.0 or higher
- racecar-webd blockchain running

#### Steps
```bash
cd api-bridge
make build
make run
```

- To run manually:
  ```bash
  ./bin/api-bridge.exe
  ```
- To specify custom ports:
  ```bash
  ./bin/api-bridge.exe --rest-port 8080 --grpc-port 9092
  ```

#### Configuration
Create `config.yaml` in the api-bridge directory:
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

#### Health Check
```bash
curl http://localhost:8080/health
```
Expected response:
```json
{
  "status": "healthy",
  "timestamp": 1751725402
}
```

#### Testing
- Run unit tests:
  ```bash
  make test
  ```
- Run integration tests:
  ```bash
  make test-integration
  ```
- (Legacy) Run integration script:
  ```bash
  ./test_blockchain_fix.sh
  ```

---

### 2. Production Deployment (Linux)
- **Platform**: Ubuntu 20.04+ or CentOS 8+
- **Method**: Systemd service
- **Use Case**: Production servers

#### Prerequisites
- Go 1.24.0+
- racecar-webd blockchain
- Systemd

#### Steps
```bash
make build
# Or cross-compile:
GOOS=linux GOARCH=amd64 go build -o bin/api-bridge ./cmd/api-bridge
```

- Create a system user and install the binary:
  ```bash
  sudo useradd -r -s /bin/false api-bridge
  sudo mkdir -p /opt/api-bridge
  sudo chown api-bridge:api-bridge /opt/api-bridge
  sudo cp bin/api-bridge /opt/api-bridge/
  sudo cp config.yaml /opt/api-bridge/
  sudo chown api-bridge:api-bridge /opt/api-bridge/*
  ```

- Create `/etc/systemd/system/api-bridge.service`:
  ```ini
  [Unit]
  Description=Web4 Race Car API Bridge
  After=network.target
  Wants=network.target

  [Service]
  Type=simple
  User=api-bridge
  Group=api-bridge
  WorkingDirectory=/opt/api-bridge
  ExecStart=/opt/api-bridge/api-bridge
  Restart=always
  RestartSec=5
  StandardOutput=journal
  StandardError=journal
  SyslogIdentifier=api-bridge
  Environment=API_BRIDGE_PORT=8080
  Environment=API_BRIDGE_LOG_LEVEL=info
  NoNewPrivileges=true
  PrivateTmp=true
  ProtectSystem=strict
  ProtectHome=true
  ReadWritePaths=/opt/api-bridge

  [Install]
  WantedBy=multi-user.target
  ```

- Enable and start the service:
  ```bash
  sudo systemctl daemon-reload
  sudo systemctl enable api-bridge
  sudo systemctl start api-bridge
  ```

- Verify deployment:
  ```bash
  sudo systemctl status api-bridge
  curl http://localhost:8080/health
  sudo journalctl -u api-bridge -f
  ```

---

### 3. Docker Deployment
- **Platform**: Docker
- **Method**: Container orchestration
- **Use Case**: Cloud deployment

#### Steps
- Create a `Dockerfile`:
  ```dockerfile
  FROM golang:1.24-alpine AS builder
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-bridge ./cmd/api-bridge

  FROM alpine:latest
  RUN apk --no-cache add ca-certificates tzdata
  WORKDIR /root/
  COPY --from=builder /app/api-bridge .
  COPY --from=builder /app/config.yaml .
  RUN addgroup -g 1001 -S api-bridge && \
      adduser -u 1001 -S api-bridge -G api-bridge
  RUN chown -R api-bridge:api-bridge /root/
  USER api-bridge
  EXPOSE 8080
  CMD ["./api-bridge"]
  ```

- Build and run:
  ```bash
  docker build -t web4-api-bridge:latest .
  docker run -d -p 8080:8080 --name api-bridge web4-api-bridge:latest
  ```

- (Optional) Use Docker Compose for multi-service deployments.

---

### 4. ARM64 Edge Deployment
- **Platform**: ARM64 (Raspberry Pi, etc.)
- **Method**: Embedded deployment
- **Use Case**: Edge devices

#### Steps
- Cross-compile:
  ```bash
  GOOS=linux GOARCH=arm64 go build -o bin/api-bridge-arm64 ./cmd/api-bridge
  ```
- Deploy to device and follow Linux systemd steps above.

---

## 🔒 Security & Environment
- Only real blockchain integration is supported—**no mock fallback**
- Use environment variables or config file for all paths and secrets
- For Ignite CLI and racecar-webd, ensure their directories are in your PATH or set environment variables:
  ```bash
  export IGNITE_CLI_PATH=/custom/path/ignite
  export RACECAR_WEBD_PATH=/custom/path/racecar-webd
  ```

---

## 🔧 Ignite CLI Path Detection and Customization

The API bridge will automatically detect the Ignite CLI in the following order:
- If the `IGNITE_CLI_PATH` environment variable is set, it will use that path.
- Otherwise, it will search for `ignite` in your system PATH and common locations, including `/snap/bin/ignite` (for Snap installations), `/usr/local/bin/ignite`, `/usr/bin/ignite`, `$GOPATH/bin/ignite`, and `$HOME/go/bin/ignite`.

If you have Ignite installed elsewhere, set the environment variable before running the API bridge:

```bash
export IGNITE_CLI_PATH=/custom/path/ignite
```

Or add the directory to your PATH:

```bash
export PATH=$PATH:/path/to/ignite
```

**Note:** The API bridge no longer uses any hard-coded user or directory paths. All environment variables (HOME, GOPATH, PATH) are detected dynamically. If you are using Snap, ensure `/snap/bin` is in your PATH or set `IGNITE_CLI_PATH` accordingly.

---

## 📋 Deployment Checklist

### Pre-Deployment
- [ ] Code review completed
- [ ] Tests passing
- [ ] Security scan completed
- [ ] Configuration validated
- [ ] Backup strategy in place

### Deployment
- [ ] Service stopped
- [ ] Backup created
- [ ] New version deployed
- [ ] Configuration updated
- [ ] Service started
- [ ] Health check passed
- [ ] Logs reviewed

### Post-Deployment
- [ ] Monitoring alerts configured
- [ ] Performance metrics collected
- [ ] User acceptance testing
- [ ] Documentation updated

---

## 📞 Support

For deployment issues:
1. Check service logs: `sudo journalctl -u api-bridge -f`
2. Verify configuration: `./api-bridge --validate-config`
3. Test connectivity: `curl http://localhost:8080/health`
4. Review troubleshooting section above

---

**Document Version**: 2.0.0  
**Last Updated**: July 2025  
**Next Review**: August 2025 