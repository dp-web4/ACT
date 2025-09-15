#!/bin/bash

# Web4-ModBatt Complete Test Suite Runner
# This script runs all tests including blockchain modules and API bridge

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}✓ PASS${NC}: $message"
            ;;
        "FAIL")
            echo -e "${RED}✗ FAIL${NC}: $message"
            ;;
        "SKIP")
            echo -e "${YELLOW}⚠ SKIP${NC}: $message"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ INFO${NC}: $message"
            ;;
    esac
}

# Function to run a test suite
run_test_suite() {
    local test_name=$1
    local test_path=$2
    local test_args=$3
    
    print_status "INFO" "Running $test_name..."
    
    if [ ! -d "$test_path" ] && [ ! -f "$test_path" ]; then
        print_status "SKIP" "$test_name - Test path not found: $test_path"
        ((SKIPPED_TESTS++))
        return
    fi
    
    # Run the test
    if go test $test_args "$test_path" -v; then
        print_status "PASS" "$test_name completed successfully"
        ((PASSED_TESTS++))
    else
        print_status "FAIL" "$test_name failed"
        ((FAILED_TESTS++))
    fi
    
    ((TOTAL_TESTS++))
    echo ""
}

# Function to run benchmark tests
run_benchmark_suite() {
    local test_name=$1
    local test_path=$2
    
    print_status "INFO" "Running benchmarks for $test_name..."
    
    if [ ! -d "$test_path" ] && [ ! -f "$test_path" ]; then
        print_status "SKIP" "$test_name benchmarks - Test path not found: $test_path"
        return
    fi
    
    # Run benchmarks
    if go test "$test_path" -bench=. -benchmem; then
        print_status "PASS" "$test_name benchmarks completed"
    else
        print_status "FAIL" "$test_name benchmarks failed"
    fi
    
    echo ""
}

# Function to check if blockchain is running
check_blockchain() {
    local response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:1317/cosmos/base/tendermint/v1beta1/node_info 2>/dev/null || echo "000")
    if [ "$response" = "200" ]; then
        return 0
    else
        return 1
    fi
}

# Function to start blockchain if needed
start_blockchain_if_needed() {
    if ! check_blockchain; then
        print_status "INFO" "Blockchain not running. Starting blockchain node..."
        
        # Start blockchain in background
        cd ..
        ignite chain serve --reset-once > /dev/null 2>&1 &
        local blockchain_pid=$!
        cd tests
        
        # Wait for blockchain to start
        local attempts=0
        while [ $attempts -lt 30 ]; do
            if check_blockchain; then
                print_status "PASS" "Blockchain started successfully (PID: $blockchain_pid)"
                return 0
            fi
            sleep 2
            ((attempts++))
        done
        
        print_status "FAIL" "Failed to start blockchain after 60 seconds"
        kill $blockchain_pid 2>/dev/null || true
        return 1
    else
        print_status "INFO" "Blockchain already running"
        return 0
    fi
}

# Main test execution
main() {
    echo "=========================================="
    echo "Web4-ModBatt Complete Test Suite Runner"
    echo "=========================================="
    echo ""
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ] && [ ! -f "../go.mod" ]; then
        print_status "FAIL" "go.mod not found. Please run from project root or tests directory."
        exit 1
    fi
    
    # Determine project root
    if [ -f "go.mod" ]; then
        PROJECT_ROOT="."
    else
        PROJECT_ROOT=".."
    fi
    
    cd "$PROJECT_ROOT"
    
    # Check if Go is available
    if ! command -v go &> /dev/null; then
        print_status "FAIL" "Go is not installed or not in PATH"
        exit 1
    fi
    
    print_status "INFO" "Project root: $(pwd)"
    print_status "INFO" "Go version: $(go version)"
    echo ""
    
    # Start blockchain if needed for integration tests
    if [ "$1" != "--unit-only" ]; then
        if ! start_blockchain_if_needed; then
            print_status "WARN" "Blockchain not available. Some integration tests will be skipped."
        fi
    fi
    
    echo "=========================================="
    echo "Running Unit Tests"
    echo "=========================================="
    echo ""
    
    # Blockchain Module Unit Tests
    run_test_suite "Component Registry Tests" "./tests/unit/componentregistry" ""
    run_test_suite "LCT Manager Tests" "./tests/unit/lctmanager" ""
    run_test_suite "Pairing Tests" "./tests/unit/pairing" ""
    run_test_suite "Pairing Queue Tests" "./tests/unit/pairingqueue" ""
    run_test_suite "Trust Tensor Tests" "./tests/unit/trusttensor" ""
    run_test_suite "Energy Cycle Tests" "./tests/unit/energycycle" ""
    
    echo "=========================================="
    echo "Running Integration Tests"
    echo "=========================================="
    echo ""
    
    # Integration Tests
    run_test_suite "Pairing Flow Integration Tests" "./tests/integration/pairing_flow" ""
    
    echo "=========================================="
    echo "Running API Bridge Tests"
    echo "=========================================="
    echo ""
    
    # API Bridge Tests
    if [ -d "api-bridge" ]; then
        cd api-bridge
        run_test_suite "API Bridge Unit Tests" "." ""
        
        # Integration tests only if blockchain is available
        if check_blockchain; then
            run_test_suite "API Bridge Integration Tests" "." ""
        else
            print_status "SKIP" "API Bridge Integration Tests - Blockchain not available"
            ((SKIPPED_TESTS++))
        fi
        cd ..
    else
        print_status "SKIP" "API Bridge tests - api-bridge directory not found"
        ((SKIPPED_TESTS++))
    fi
    
    echo "=========================================="
    echo "Running End-to-End Tests"
    echo "=========================================="
    echo ""
    
    # E2E Tests (only if blockchain is available)
    if check_blockchain; then
        run_test_suite "Race Car E2E Tests" "./tests/e2e/race_car_scenarios" ""
    else
        print_status "SKIP" "E2E Tests - Blockchain not available"
        ((SKIPPED_TESTS++))
    fi
    
    echo "=========================================="
    echo "Running Benchmark Tests"
    echo "=========================================="
    echo ""
    
    # Benchmark Tests
    run_benchmark_suite "Component Registry Benchmarks" "./tests/unit/componentregistry"
    run_benchmark_suite "LCT Manager Benchmarks" "./tests/unit/lctmanager"
    run_benchmark_suite "Trust Tensor Benchmarks" "./tests/unit/trusttensor"
    run_benchmark_suite "Energy Cycle Benchmarks" "./tests/unit/energycycle"
    
    if [ -d "api-bridge" ]; then
        cd api-bridge
        run_benchmark_suite "API Bridge Benchmarks" "."
        cd ..
    fi
    
    echo "=========================================="
    echo "Test Summary"
    echo "=========================================="
    echo ""
    
    print_status "INFO" "Total test suites: $TOTAL_TESTS"
    print_status "PASS" "Passed: $PASSED_TESTS"
    
    if [ $FAILED_TESTS -gt 0 ]; then
        print_status "FAIL" "Failed: $FAILED_TESTS"
    else
        print_status "PASS" "Failed: $FAILED_TESTS"
    fi
    
    if [ $SKIPPED_TESTS -gt 0 ]; then
        print_status "SKIP" "Skipped: $SKIPPED_TESTS"
    fi
    
    echo ""
    
    # Calculate success rate
    if [ $TOTAL_TESTS -gt 0 ]; then
        local success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
        print_status "INFO" "Success rate: ${success_rate}%"
        
        if [ $success_rate -ge 90 ]; then
            print_status "PASS" "Excellent test coverage achieved!"
        elif [ $success_rate -ge 80 ]; then
            print_status "PASS" "Good test coverage achieved."
        elif [ $success_rate -ge 70 ]; then
            print_status "WARN" "Acceptable test coverage, but improvements needed."
        else
            print_status "FAIL" "Poor test coverage. Please improve test coverage."
        fi
    fi
    
    echo ""
    
    # Exit with appropriate code
    if [ $FAILED_TESTS -gt 0 ]; then
        print_status "FAIL" "Some tests failed. Please review and fix issues."
        exit 1
    else
        print_status "PASS" "All tests passed successfully!"
        exit 0
    fi
}

# Handle command line arguments
case "${1:-}" in
    "--unit-only")
        print_status "INFO" "Running unit tests only"
        main "$@"
        ;;
    "--help"|"-h")
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --unit-only    Run only unit tests (skip integration and E2E)"
        echo "  --help, -h     Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0              # Run all tests"
        echo "  $0 --unit-only  # Run only unit tests"
        ;;
    *)
        main "$@"
        ;;
esac 