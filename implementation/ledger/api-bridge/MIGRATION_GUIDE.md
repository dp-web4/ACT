# API Bridge - Migration Guide

## 🎯 What Changed

The API Bridge project has been refactored to follow Go best practices for maximum portability and maintainability. This guide helps you understand what changed and how to work with the new structure.

## 📋 Summary of Changes

### ✅ **Fixed Multiple Main Functions Issue**
- **Before**: All binaries were in the root directory with `package main`, causing build conflicts
- **After**: Each binary is in its own directory under `cmd/` with isolated `package main`

### ✅ **New Project Structure**
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
└── config.yaml                  # Configuration file
```

### ✅ **Removed Mock Fallback Logic**
- **Before**: API bridge would fall back to mock transactions if blockchain failed
- **After**: API bridge requires real blockchain integration and fails clearly if unavailable

## 🚀 How to Work with the New Structure

### Building the Application

**Old way (no longer works):**
```bash
go build -o api-bridge.exe .
```

**New way (recommended):**
```bash
# Using Makefile
make build

# Or direct Go build
go build -o bin/api-bridge.exe ./cmd/api-bridge
```

### Running the Application

**Old way (no longer works):**
```bash
./api-bridge.exe
```

**New way (recommended):**
```bash
# Using Makefile
make run

# Or manually
./bin/api-bridge.exe
```

### Running Tests

**Old way:**
```bash
go test ./...
```

**New way:**
```bash
# Unit tests
make test

# Integration tests
make test-integration
```

## 🛠️ Available Makefile Commands

The project now includes a comprehensive Makefile:

```bash
make help          # Show all available commands
make build         # Build main API bridge
make build-all     # Build all binaries
make test          # Run unit tests
make test-integration # Run integration tests
make run           # Build and run API bridge
make run-debug     # Run with debug logging
make clean         # Clean build artifacts
```

## 📚 Updated Documentation

All documentation files have been updated to reflect the new structure:

- ✅ **README.md** - Updated build/run instructions
- ✅ **QUICK_REFERENCE.md** - Updated quick start guide
- ✅ **DEPLOYMENT_GUIDE.md** - Updated deployment instructions
- ✅ **TECHNICAL_SPEC.md** - Updated technical specifications
- ✅ **DEVELOPMENT_ROADMAP.md** - Updated development roadmap
- ✅ **README_NEW_STRUCTURE.md** - New comprehensive guide

## 🔧 Development Workflow

### 1. **Getting Started**
```bash
cd api-bridge
make build
make run
```

### 2. **Making Changes**
- Edit code in `internal/` for business logic
- Edit code in `cmd/api-bridge/` for main application
- Edit code in `cmd/debug-*` for debug tools

### 3. **Testing Changes**
```bash
make test              # Run unit tests
make test-integration  # Run integration tests
```

### 4. **Building for Production**
```bash
make build-all         # Build all binaries
# Binaries are created in bin/ directory
```

## 🐛 Troubleshooting

### Build Issues
- **Error**: `main redeclared in this block`
  - **Solution**: This is fixed! Each binary now has its own package
- **Error**: `go build -o api-bridge.exe .` fails
  - **Solution**: Use `go build -o bin/api-bridge.exe ./cmd/api-bridge`

### Runtime Issues
- **Error**: `./api-bridge: No such file or directory`
  - **Solution**: Use `./bin/api-bridge` instead
- **Error**: Mock transactions being used
  - **Solution**: This is fixed! API bridge now requires real blockchain integration

### Port Conflicts
```bash
# Stop any running instances
pkill -f api-bridge

# Use different ports
./bin/api-bridge --rest-port 8083 --grpc-port 9094
```

## 🔄 Migration Checklist

When you return to work, verify these items:

- [ ] **Build works**: `make build` succeeds
- [ ] **Run works**: `make run` starts the API bridge
- [ ] **Tests pass**: `make test` passes
- [ ] **Documentation updated**: All docs reflect new structure
- [ ] **No mock fallbacks**: API bridge fails clearly if blockchain unavailable
- [ ] **Portable**: Works on different systems

## 🎉 Benefits of the New Structure

1. **No More Build Conflicts**: Each binary is isolated
2. **Portable**: Works on any OS with Go toolchain
3. **Team-Friendly**: No more "multiple main functions" issues
4. **Standard Go Layout**: Follows Go community conventions
5. **Easy CI/CD**: Clear separation for build systems
6. **Scalable**: Easy to add new binaries or libraries

## 📞 Need Help?

If you encounter any issues:

1. **Check the Makefile**: `make help` shows all available commands
2. **Read README_NEW_STRUCTURE.md**: Comprehensive guide to the new structure
3. **Check logs**: API bridge now provides clear error messages
4. **Verify blockchain**: Ensure racecar-webd is running and accessible

The new structure is much more robust and follows Go best practices. Welcome back to a cleaner, more maintainable codebase! 🚀 