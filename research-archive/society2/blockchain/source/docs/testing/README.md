# Web4-ModBatt Testing Framework

## Overview

This document provides comprehensive testing guidelines for the Web4-ModBatt blockchain system. The testing framework covers all six modules with unit tests, integration tests, and end-to-end scenarios.

## Testing Philosophy

### Security-First Testing
- **No Key Storage**: Verify that no cryptographic key halves are stored on-chain
- **LCT-Mediated Trust**: Test all pairing operations use LCT relationships
- **Offline Resilience**: Validate queue operations for offline devices
- **Audit Trail**: Ensure all operations emit proper events for off-chain audit

### Battery Management Focus
- **Component Lifecycle**: Test component registration, verification, and authorization
- **Energy Operations**: Validate ATP/ADP token-based energy transfers
- **Trust Evolution**: Test trust tensor calculations and updates
- **Multi-Transport**: Verify support for SD card, Bluetooth, WiFi, CANBus operations

## Test Structure

```
tests/
├── unit/                    # Unit tests for individual functions
│   ├── componentregistry/   # Component registry unit tests
│   ├── lctmanager/         # LCT manager unit tests
│   ├── pairing/            # Pairing module unit tests
│   ├── pairingqueue/       # Pairing queue unit tests
│   ├── trusttensor/        # Trust tensor unit tests
│   └── energycycle/        # Energy cycle unit tests
├── integration/            # Integration tests between modules
│   ├── pairing_flow/       # Complete pairing workflows
│   ├── energy_operations/  # Energy transfer scenarios
│   ├── trust_calculations/ # Trust tensor integration
│   └── offline_scenarios/  # Offline device handling
├── e2e/                    # End-to-end system tests
│   ├── race_car_scenarios/ # Real-world race car scenarios
│   ├── battery_pack_tests/ # Complete battery pack operations
│   └── stress_tests/       # Performance and load testing
├── fixtures/               # Test data and mock objects
│   ├── components/         # Mock battery components
│   ├── lcts/              # Sample LCT relationships
│   └── scenarios/         # Complete test scenarios
└── scripts/               # Test automation scripts
    ├── run_all_tests.sh   # Run complete test suite
    ├── test_coverage.sh   # Generate coverage reports
    └── benchmark_tests.sh # Performance benchmarking
```

## Module Testing Strategy

### 1. Component Registry Module
**Focus Areas:**
- Component registration with manufacturer ID extraction
- Component verification and status management
- Authorization rules and pairing permissions
- Bidirectional pairing authorization checks

**Key Test Scenarios:**
- Register new battery module with JSON metadata
- Verify component authenticity
- Update authorization rules
- Check pairing permissions between components

### 2. LCT Manager Module
**Focus Areas:**
- LCT relationship creation and management
- Split-key generation (off-chain storage)
- Key exchange protocols
- Relationship lifecycle management

**Key Test Scenarios:**
- Create LCT relationship between two components
- Verify no key halves stored on-chain
- Test LCT-mediated pairing initiation
- Validate relationship termination

### 3. Pairing Module
**Focus Areas:**
- Bidirectional authentication
- Challenge-response protocols
- Session key exchange
- Pairing completion and validation

**Key Test Scenarios:**
- Initiate bidirectional pairing
- Complete pairing with authentication
- Revoke pairing relationships
- Test pairing status queries

### 4. Pairing Queue Module
**Focus Areas:**
- Offline device queue management
- Authentication Controller integration
- Multi-transport support
- Queue processing and status tracking

**Key Test Scenarios:**
- Queue pairing request for offline device
- Process offline queue when device comes online
- Test proxy component operations
- Validate queue cancellation

### 5. Trust Tensor Module
**Focus Areas:**
- Multi-dimensional trust calculations
- Trust evolution over time
- Witness-based validation
- T3/V3 tensor algorithms

**Key Test Scenarios:**
- Create relationship trust tensor
- Update tensor scores based on operations
- Add tensor witnesses
- Calculate relationship trust scores

### 6. Energy Cycle Module
**Focus Areas:**
- ATP/ADP token management
- Energy operation creation and execution
- Trust-based validation
- Energy flow history tracking

**Key Test Scenarios:**
- Create energy operation between components
- Execute energy transfer with trust validation
- Validate relationship value with attention tokens
- Track energy flow history

## Running Tests

### Prerequisites
```bash
# Install Go testing tools
go install github.com/stretchr/testify/assert@latest
go install github.com/stretchr/testify/suite@latest
go install github.com/stretchr/testify/mock@latest

# Install coverage tools
go install golang.org/x/tools/cmd/cover@latest
go install github.com/axw/gocov/gocov@latest
```

### Unit Tests
```bash
# Run all unit tests
go test ./tests/unit/... -v

# Run specific module tests
go test ./tests/unit/componentregistry/... -v
go test ./tests/unit/lctmanager/... -v
go test ./tests/unit/pairing/... -v
go test ./tests/unit/pairingqueue/... -v
go test ./tests/unit/trusttensor/... -v
go test ./tests/unit/energycycle/... -v
```

### Integration Tests
```bash
# Run integration tests
go test ./tests/integration/... -v

# Run specific integration scenarios
go test ./tests/integration/pairing_flow/... -v
go test ./tests/integration/energy_operations/... -v
go test ./tests/integration/trust_calculations/... -v
go test ./tests/integration/offline_scenarios/... -v
```

### End-to-End Tests
```bash
# Run E2E tests (requires blockchain node)
go test ./tests/e2e/... -v

# Run specific E2E scenarios
go test ./tests/e2e/race_car_scenarios/... -v
go test ./tests/e2e/battery_pack_tests/... -v
go test ./tests/e2e/stress_tests/... -v
```

### Test Scripts
```bash
# Run complete test suite
./tests/scripts/run_all_tests.sh

# Generate coverage report
./tests/scripts/test_coverage.sh

# Run performance benchmarks
./tests/scripts/benchmark_tests.sh
```

## Test Coverage Requirements

### Minimum Coverage Targets
- **Unit Tests**: 90% code coverage
- **Integration Tests**: 85% module interaction coverage
- **E2E Tests**: 80% business scenario coverage

### Critical Path Coverage
- **Security Functions**: 100% coverage required
- **Key Exchange**: 100% coverage required
- **Trust Calculations**: 95% coverage required
- **Energy Operations**: 90% coverage required

## Test Data Management

### Fixtures
Test fixtures are stored in `tests/fixtures/` and include:
- **Mock Components**: Sample battery components with realistic metadata
- **LCT Relationships**: Pre-configured LCT relationships for testing
- **Test Scenarios**: Complete end-to-end scenarios for validation

### Test Data Generation
```bash
# Generate test fixtures
go run tests/scripts/generate_fixtures.go

# Validate test data
go run tests/scripts/validate_fixtures.go
```

## Continuous Integration

### GitHub Actions
Tests are automatically run on:
- **Pull Requests**: All tests must pass
- **Main Branch**: Full test suite + coverage reporting
- **Release Tags**: Performance benchmarks + security scans

### Test Reports
- **Coverage Reports**: Generated for each module
- **Performance Metrics**: Benchmark results tracked over time
- **Security Scans**: Automated security testing

## Debugging Tests

### Common Issues
1. **Test Database**: Ensure test database is properly initialized
2. **Mock Objects**: Verify mock objects are properly configured
3. **Async Operations**: Handle asynchronous blockchain operations
4. **Resource Cleanup**: Ensure proper cleanup after tests

### Debug Commands
```bash
# Run tests with verbose output
go test -v -run TestName

# Run tests with race detection
go test -race ./tests/...

# Run tests with memory profiling
go test -memprofile=mem.prof ./tests/...

# Run tests with CPU profiling
go test -cpuprofile=cpu.prof ./tests/...
```

## Contributing to Tests

### Adding New Tests
1. **Unit Tests**: Add to appropriate module directory
2. **Integration Tests**: Create new scenario in integration directory
3. **E2E Tests**: Add to e2e directory with realistic scenarios
4. **Fixtures**: Update test data as needed

### Test Naming Conventions
- **Unit Tests**: `TestFunctionName_Scenario`
- **Integration Tests**: `TestModule_Integration_Scenario`
- **E2E Tests**: `TestEndToEnd_Scenario`

### Documentation
- **Test Purpose**: Document what each test validates
- **Test Data**: Explain test fixtures and mock objects
- **Expected Results**: Document expected outcomes
- **Dependencies**: List any external dependencies

## Performance Testing

### Benchmarks
- **Transaction Throughput**: Measure TPS for different operations
- **Memory Usage**: Track memory consumption during operations
- **Response Time**: Measure latency for key operations
- **Scalability**: Test with increasing load

### Load Testing
- **Concurrent Users**: Test with multiple simultaneous users
- **Large Datasets**: Test with realistic data volumes
- **Stress Testing**: Test system limits and failure modes

## Security Testing

### Security Scenarios
- **Key Management**: Verify no keys stored on-chain
- **Authentication**: Test authentication bypass attempts
- **Authorization**: Test unauthorized access attempts
- **Data Integrity**: Test data tampering scenarios

### Penetration Testing
- **API Security**: Test API endpoint security
- **Blockchain Security**: Test blockchain-specific vulnerabilities
- **Network Security**: Test network-level security

## Maintenance

### Test Maintenance
- **Regular Updates**: Update tests when code changes
- **Dependency Updates**: Keep test dependencies current
- **Performance Monitoring**: Track test performance over time
- **Coverage Monitoring**: Monitor coverage trends

### Test Environment
- **Docker Support**: Containerized test environment
- **CI/CD Integration**: Automated test execution
- **Reporting**: Automated test reporting and notifications 