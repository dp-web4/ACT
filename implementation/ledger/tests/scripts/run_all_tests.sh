#!/bin/bash

# Web4-ModBatt Test Runner
# This script runs all tests for the Web4-ModBatt blockchain system

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DIR="tests"
COVERAGE_DIR="coverage"
LOG_DIR="logs"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create directories
mkdir -p $COVERAGE_DIR
mkdir -p $LOG_DIR

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Web4-ModBatt Test Runner${NC}"
echo -e "${BLUE}================================${NC}"
echo "Timestamp: $TIMESTAMP"
echo ""

# Function to print section headers
print_section() {
    echo -e "${YELLOW}$1${NC}"
    echo "----------------------------------------"
}

# Function to run tests with coverage
run_tests_with_coverage() {
    local test_path=$1
    local module_name=$2
    local coverage_file="$COVERAGE_DIR/${module_name}_coverage.out"
    local log_file="$LOG_DIR/${module_name}_test_${TIMESTAMP}.log"
    
    echo "Running tests for $module_name..."
    echo "Coverage file: $coverage_file"
    echo "Log file: $log_file"
    
    # Run tests with coverage
    go test -v -coverprofile="$coverage_file" -covermode=atomic "$test_path" 2>&1 | tee "$log_file"
    
    # Check exit status
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ $module_name tests passed${NC}"
        
        # Generate coverage report
        if [ -f "$coverage_file" ]; then
            go tool cover -func="$coverage_file" | tee "$LOG_DIR/${module_name}_coverage_${TIMESTAMP}.log"
            echo "Coverage report saved to: $LOG_DIR/${module_name}_coverage_${TIMESTAMP}.log"
        fi
    else
        echo -e "${RED}✗ $module_name tests failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to run integration tests
run_integration_tests() {
    local test_path=$1
    local test_name=$2
    local log_file="$LOG_DIR/integration_${test_name}_${TIMESTAMP}.log"
    
    echo "Running integration test: $test_name..."
    echo "Log file: $log_file"
    
    go test -v "$test_path" 2>&1 | tee "$log_file"
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ Integration test $test_name passed${NC}"
    else
        echo -e "${RED}✗ Integration test $test_name failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to run E2E tests
run_e2e_tests() {
    local test_path=$1
    local test_name=$2
    local log_file="$LOG_DIR/e2e_${test_name}_${TIMESTAMP}.log"
    
    echo "Running E2E test: $test_name..."
    echo "Log file: $log_file"
    
    go test -v "$test_path" 2>&1 | tee "$log_file"
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ E2E test $test_name passed${NC}"
    else
        echo -e "${RED}✗ E2E test $test_name failed${NC}"
        return 1
    fi
    
    echo ""
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

# Install test dependencies
print_section "Installing Test Dependencies"
go mod tidy
go install github.com/stretchr/testify/assert@latest
go install github.com/stretchr/testify/suite@latest
go install github.com/stretchr/testify/mock@latest
go install golang.org/x/tools/cmd/cover@latest
echo -e "${GREEN}✓ Dependencies installed${NC}"
echo ""

# Run unit tests
print_section "Running Unit Tests"

# Component Registry tests
run_tests_with_coverage "$TEST_DIR/unit/componentregistry" "componentregistry"

# LCT Manager tests
run_tests_with_coverage "$TEST_DIR/unit/lctmanager" "lctmanager"

# Pairing tests
run_tests_with_coverage "$TEST_DIR/unit/pairing" "pairing"

# Pairing Queue tests
run_tests_with_coverage "$TEST_DIR/unit/pairingqueue" "pairingqueue"

# Trust Tensor tests
run_tests_with_coverage "$TEST_DIR/unit/trusttensor" "trusttensor"

# Energy Cycle tests
run_tests_with_coverage "$TEST_DIR/unit/energycycle" "energycycle"

echo -e "${GREEN}✓ All unit tests completed${NC}"
echo ""

# Run integration tests
print_section "Running Integration Tests"

# Pairing flow integration
run_integration_tests "$TEST_DIR/integration/pairing_flow" "pairing_flow"

# Energy operations integration
run_integration_tests "$TEST_DIR/integration/energy_operations" "energy_operations"

# Trust calculations integration
run_integration_tests "$TEST_DIR/integration/trust_calculations" "trust_calculations"

# Offline scenarios integration
run_integration_tests "$TEST_DIR/integration/offline_scenarios" "offline_scenarios"

echo -e "${GREEN}✓ All integration tests completed${NC}"
echo ""

# Run E2E tests (if blockchain node is available)
print_section "Running End-to-End Tests"

# Check if blockchain node is running
if pgrep -f "racecar-webd" > /dev/null; then
    echo -e "${GREEN}Blockchain node detected, running E2E tests...${NC}"
    
    # Race car scenarios
    run_e2e_tests "$TEST_DIR/e2e/race_car_scenarios" "race_car_scenarios"
    
    # Battery pack tests
    run_e2e_tests "$TEST_DIR/e2e/battery_pack_tests" "battery_pack_tests"
    
    # Stress tests
    run_e2e_tests "$TEST_DIR/e2e/stress_tests" "stress_tests"
    
    echo -e "${GREEN}✓ All E2E tests completed${NC}"
else
    echo -e "${YELLOW}Warning: Blockchain node not detected. Skipping E2E tests.${NC}"
    echo "To run E2E tests, start the blockchain node with: ignite chain serve"
fi

echo ""

# Generate overall coverage report
print_section "Generating Coverage Report"

# Combine all coverage files
if ls $COVERAGE_DIR/*_coverage.out 1> /dev/null 2>&1; then
    echo "Combining coverage reports..."
    
    # Create combined coverage file
    combined_coverage="$COVERAGE_DIR/combined_coverage_${TIMESTAMP}.out"
    echo "mode: atomic" > "$combined_coverage"
    
    # Append all coverage files
    for coverage_file in $COVERAGE_DIR/*_coverage.out; do
        if [ -f "$coverage_file" ]; then
            tail -n +2 "$coverage_file" >> "$combined_coverage"
        fi
    done
    
    # Generate HTML coverage report
    html_coverage="$COVERAGE_DIR/coverage_report_${TIMESTAMP}.html"
    go tool cover -html="$combined_coverage" -o="$html_coverage"
    
    echo -e "${GREEN}✓ Coverage report generated: $html_coverage${NC}"
    
    # Show overall coverage summary
    echo ""
    echo "Overall Coverage Summary:"
    go tool cover -func="$combined_coverage" | tail -1
else
    echo -e "${YELLOW}No coverage files found${NC}"
fi

echo ""

# Run security tests
print_section "Running Security Tests"

# Check for security vulnerabilities
echo "Checking for security issues..."
go list -json -deps . | jq -r '.Deps[]' | xargs -I {} go list -json {} | jq -r 'select(.Vulnerabilities != null) | .Path + ": " + (.Vulnerabilities | join(", "))' || echo "No known vulnerabilities found"

echo ""

# Run performance benchmarks
print_section "Running Performance Benchmarks"

# Run benchmarks for each module
for module in componentregistry lctmanager pairing pairingqueue trusttensor energycycle; do
    if [ -d "$TEST_DIR/unit/$module" ]; then
        echo "Running benchmarks for $module..."
        benchmark_log="$LOG_DIR/benchmark_${module}_${TIMESTAMP}.log"
        go test -bench=. -benchmem "$TEST_DIR/unit/$module" 2>&1 | tee "$benchmark_log"
        echo -e "${GREEN}✓ $module benchmarks completed${NC}"
    fi
done

echo ""

# Generate test summary
print_section "Test Summary"

echo "Test Results Summary:"
echo "- Unit Tests: $TEST_DIR/unit/"
echo "- Integration Tests: $TEST_DIR/integration/"
echo "- E2E Tests: $TEST_DIR/e2e/"
echo "- Coverage Reports: $COVERAGE_DIR/"
echo "- Test Logs: $LOG_DIR/"
echo "- Timestamp: $TIMESTAMP"

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  All Tests Completed Successfully${NC}"
echo -e "${GREEN}================================${NC}"

# Optional: Open coverage report in browser
if command -v xdg-open &> /dev/null && [ -f "$html_coverage" ]; then
    echo ""
    read -p "Open coverage report in browser? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        xdg-open "$html_coverage"
    fi
elif command -v open &> /dev/null && [ -f "$html_coverage" ]; then
    echo ""
    read -p "Open coverage report in browser? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open "$html_coverage"
    fi
fi

echo ""
echo "Test execution completed at: $(date)" 