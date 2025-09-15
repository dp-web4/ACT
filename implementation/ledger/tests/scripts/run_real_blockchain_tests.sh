#!/bin/bash

# Web4-ModBatt Real Blockchain Test Runner
# This script runs real blockchain integration tests for the Web4-ModBatt system

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
BLOCKCHAIN_RPC_URL="http://localhost:1317"
BLOCKCHAIN_TENDERMINT_URL="http://localhost:26657"

# Create directories
mkdir -p $COVERAGE_DIR
mkdir -p $LOG_DIR

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Web4-ModBatt Real Blockchain Test Runner${NC}"
echo -e "${BLUE}================================${NC}"
echo "Timestamp: $TIMESTAMP"
echo ""

# Function to print section headers
print_section() {
    echo -e "${YELLOW}$1${NC}"
    echo "----------------------------------------"
}

# Function to check blockchain health
check_blockchain_health() {
    print_section "Checking Blockchain Health"
    
    # Check if blockchain node is running
    if ! pgrep -f "racecar-webd" > /dev/null; then
        echo -e "${RED}❌ Blockchain node not running${NC}"
        echo "To start the blockchain node, run:"
        echo "  ignite chain serve"
        echo ""
        echo "Or if using Docker:"
        echo "  docker-compose up -d"
        echo ""
        return 1
    fi
    
    echo -e "${GREEN}✅ Blockchain node is running${NC}"
    
    # Check RPC endpoint
    if curl -s --max-time 5 "$BLOCKCHAIN_RPC_URL/cosmos/base/tendermint/v1beta1/node_info" > /dev/null; then
        echo -e "${GREEN}✅ RPC endpoint is accessible${NC}"
    else
        echo -e "${RED}❌ RPC endpoint not accessible at $BLOCKCHAIN_RPC_URL${NC}"
        return 1
    fi
    
    # Check Tendermint endpoint
    if curl -s --max-time 5 "$BLOCKCHAIN_TENDERMINT_URL/status" > /dev/null; then
        echo -e "${GREEN}✅ Tendermint endpoint is accessible${NC}"
    else
        echo -e "${RED}❌ Tendermint endpoint not accessible at $BLOCKCHAIN_TENDERMINT_URL${NC}"
        return 1
    fi
    
    # Get blockchain info
    echo ""
    echo -e "${BLUE}Blockchain Information:${NC}"
    NODE_INFO=$(curl -s "$BLOCKCHAIN_RPC_URL/cosmos/base/tendermint/v1beta1/node_info" | jq -r '.node_info.moniker // "Unknown"')
    echo "Node: $NODE_INFO"
    
    CHAIN_STATUS=$(curl -s "$BLOCKCHAIN_TENDERMINT_URL/status" | jq -r '.result.sync_info.latest_block_height // "Unknown"')
    echo "Latest Block: $CHAIN_STATUS"
    
    echo ""
    return 0
}

# Function to run real blockchain unit tests
run_real_blockchain_unit_tests() {
    local test_path=$1
    local module_name=$2
    local coverage_file="$COVERAGE_DIR/${module_name}_real_blockchain_coverage.out"
    local log_file="$LOG_DIR/${module_name}_real_blockchain_test_${TIMESTAMP}.log"
    
    echo "Running real blockchain tests for $module_name..."
    echo "Coverage file: $coverage_file"
    echo "Log file: $log_file"
    
    # Run tests with coverage - use specific test file if provided
    if [[ "$module_name" == "trusttensor" ]]; then
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path -run "RealBlockchain" 2>&1 | tee "$log_file"
    elif [[ "$module_name" == "energycycle_simple" ]]; then
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path/energy_cycle_simple_real_blockchain_test.go -run "SimpleRealBlockchain" 2>&1 | tee "$log_file"
    elif [[ "$module_name" == "lctmanager_simple" ]]; then
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path/lct_manager_simple_real_blockchain_test.go -run "SimpleRealBlockchain" 2>&1 | tee "$log_file"
    elif [[ "$module_name" == "pairingqueue_simple" ]]; then
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path/pairing_queue_simple_real_blockchain_test.go -run "SimpleRealBlockchain" 2>&1 | tee "$log_file"
    elif [[ "$module_name" == "componentregistry_simple" ]]; then
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path/component_registry_simple_real_blockchain_test.go -run "SimpleRealBlockchain" 2>&1 | tee "$log_file"
    else
        go test -v -coverprofile="$coverage_file" -covermode=atomic ./$test_path -run "RealBlockchain" 2>&1 | tee "$log_file"
    fi
    
    # Check exit status
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ $module_name real blockchain tests passed${NC}"
        
        # Generate coverage report
        if [ -f "$coverage_file" ]; then
            go tool cover -func="$coverage_file" | tee "$LOG_DIR/${module_name}_real_blockchain_coverage_${TIMESTAMP}.log"
            echo "Coverage report saved to: $LOG_DIR/${module_name}_real_blockchain_coverage_${TIMESTAMP}.log"
        fi
    else
        echo -e "${RED}✗ $module_name real blockchain tests failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to run real blockchain integration tests
run_real_blockchain_integration_tests() {
    local test_path=$1
    local test_name=$2
    local log_file="$LOG_DIR/integration_real_blockchain_${test_name}_${TIMESTAMP}.log"
    
    echo "Running real blockchain integration test: $test_name..."
    echo "Log file: $log_file"
    
    go test -v ./"$test_path" -run "RealBlockchain" 2>&1 | tee "$log_file"
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ Real blockchain integration test $test_name passed${NC}"
    else
        echo -e "${RED}✗ Real blockchain integration test $test_name failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to run real blockchain E2E tests
run_real_blockchain_e2e_tests() {
    local test_path=$1
    local test_name=$2
    local log_file="$LOG_DIR/e2e_real_blockchain_${test_name}_${TIMESTAMP}.log"
    
    echo "Running real blockchain E2E test: $test_name..."
    echo "Log file: $log_file"
    
    go test -v ./"$test_path" -run "RealBlockchain" 2>&1 | tee "$log_file"
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ Real blockchain E2E test $test_name passed${NC}"
    else
        echo -e "${RED}✗ Real blockchain E2E test $test_name failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to run performance benchmarks
run_real_blockchain_benchmarks() {
    local test_path=$1
    local module_name=$2
    local log_file="$LOG_DIR/benchmark_real_blockchain_${module_name}_${TIMESTAMP}.log"
    
    echo "Running real blockchain benchmarks for $module_name..."
    echo "Log file: $log_file"
    
    go test -v -bench=Benchmark -benchmem ./"$test_path" -run "RealBlockchain" 2>&1 | tee "$log_file"
    
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ $module_name real blockchain benchmarks completed${NC}"
    else
        echo -e "${RED}✗ $module_name real blockchain benchmarks failed${NC}"
        return 1
    fi
    
    echo ""
}

# Check prerequisites
print_section "Checking Prerequisites"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

# Check if jq is installed (for JSON parsing)
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}Warning: jq is not installed. Installing...${NC}"
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y jq
    elif command -v yum &> /dev/null; then
        sudo yum install -y jq
    elif command -v brew &> /dev/null; then
        brew install jq
    else
        echo -e "${RED}Error: Cannot install jq automatically. Please install it manually.${NC}"
        exit 1
    fi
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

# Check blockchain health
if ! check_blockchain_health; then
    echo -e "${RED}❌ Blockchain health check failed. Cannot run real blockchain tests.${NC}"
    echo ""
    echo -e "${YELLOW}To fix this:${NC}"
    echo "1. Start the blockchain node:"
    echo "   ignite chain serve"
    echo ""
    echo "2. Wait for the node to fully sync"
    echo ""
    echo "3. Run this script again"
    echo ""
    exit 1
fi

# Install test dependencies
print_section "Installing Test Dependencies"
go mod tidy
echo -e "${GREEN}✓ Dependencies updated${NC}"
echo ""

# Run Phase 1: Critical Real Blockchain Unit Tests
print_section "Phase 1: Critical Real Blockchain Unit Tests"

# Trust Tensor real blockchain tests
run_real_blockchain_unit_tests "$TEST_DIR/unit/trusttensor" "trusttensor"

# Energy Cycle real blockchain tests (using simplified version)
run_real_blockchain_unit_tests "$TEST_DIR/unit/energycycle" "energycycle_simple"

# LCT Manager real blockchain tests (using simplified version)
run_real_blockchain_unit_tests "$TEST_DIR/unit/lctmanager" "lctmanager_simple"

# Pairing Queue real blockchain tests (using simplified version)
run_real_blockchain_unit_tests "$TEST_DIR/unit/pairingqueue" "pairingqueue_simple"

# Component Registry real blockchain tests (using simplified version)
run_real_blockchain_unit_tests "$TEST_DIR/unit/componentregistry" "componentregistry_simple"

echo -e "${GREEN}✓ Phase 1 completed${NC}"
echo ""

# Run Phase 2: Real Blockchain Integration Tests
print_section "Phase 2: Real Blockchain Integration Tests"

# Energy operations integration
if [ -d "$TEST_DIR/integration/energy_operations" ]; then
    run_real_blockchain_integration_tests "$TEST_DIR/integration/energy_operations" "energy_operations"
fi

# Pairing flow integration
if [ -d "$TEST_DIR/integration/pairing_flow" ]; then
    run_real_blockchain_integration_tests "$TEST_DIR/integration/pairing_flow" "pairing_flow"
fi

# Trust calculations integration
if [ -d "$TEST_DIR/integration/trust_calculations" ]; then
    run_real_blockchain_integration_tests "$TEST_DIR/integration/trust_calculations" "trust_calculations"
fi

echo -e "${GREEN}✓ Phase 2 completed${NC}"
echo ""

# Run Phase 3: Real Blockchain E2E Tests
print_section "Phase 3: Real Blockchain E2E Tests"

# Race car scenarios
if [ -d "$TEST_DIR/e2e/race_car_scenarios" ]; then
    run_real_blockchain_e2e_tests "$TEST_DIR/e2e/race_car_scenarios" "race_car_scenarios"
fi

# Battery pack tests
if [ -d "$TEST_DIR/e2e/battery_pack_tests" ]; then
    run_real_blockchain_e2e_tests "$TEST_DIR/e2e/battery_pack_tests" "battery_pack_tests"
fi

echo -e "${GREEN}✓ Phase 3 completed${NC}"
echo ""

# Run Phase 4: Real Blockchain Performance Tests
print_section "Phase 4: Real Blockchain Performance Tests"

# Trust Tensor benchmarks
run_real_blockchain_benchmarks "$TEST_DIR/unit/trusttensor" "trusttensor"

# Energy Cycle benchmarks
run_real_blockchain_benchmarks "$TEST_DIR/unit/energycycle" "energycycle"

# LCT Manager benchmarks
run_real_blockchain_benchmarks "$TEST_DIR/unit/lctmanager" "lctmanager"

echo -e "${GREEN}✓ Phase 4 completed${NC}"
echo ""

# Generate overall coverage report
print_section "Generating Overall Coverage Report"

# Combine coverage files
COMBINED_COVERAGE="$COVERAGE_DIR/combined_real_blockchain_coverage.out"
echo "mode: atomic" > "$COMBINED_COVERAGE"

for coverage_file in "$COVERAGE_DIR"/*_real_blockchain_coverage.out; do
    if [ -f "$coverage_file" ]; then
        tail -n +2 "$coverage_file" >> "$COMBINED_COVERAGE"
    fi
done

# Generate HTML coverage report
if [ -f "$COMBINED_COVERAGE" ]; then
    go tool cover -html="$COMBINED_COVERAGE" -o "$COVERAGE_DIR/real_blockchain_coverage.html"
    echo -e "${GREEN}✓ HTML coverage report generated: $COVERAGE_DIR/real_blockchain_coverage.html${NC}"
fi

# Generate summary
print_section "Test Summary"

echo -e "${GREEN}✅ Real blockchain tests completed successfully!${NC}"
echo ""
echo -e "${BLUE}Test Results:${NC}"
echo "- Unit Tests: Real blockchain integration validated"
echo "- Integration Tests: Cross-module operations validated"
echo "- E2E Tests: End-to-end workflows validated"
echo "- Performance Tests: Real blockchain performance measured"
echo ""
echo -e "${BLUE}Coverage Reports:${NC}"
echo "- Individual module coverage: $COVERAGE_DIR/"
echo "- Combined coverage: $COMBINED_COVERAGE"
echo "- HTML report: $COVERAGE_DIR/real_blockchain_coverage.html"
echo ""
echo -e "${BLUE}Log Files:${NC}"
echo "- Test logs: $LOG_DIR/"
echo "- Timestamp: $TIMESTAMP"
echo ""

# Final blockchain health check
print_section "Final Blockchain Health Check"
if check_blockchain_health; then
    echo -e "${GREEN}✅ Blockchain remains healthy after all tests${NC}"
else
    echo -e "${YELLOW}⚠️  Blockchain health degraded during testing${NC}"
fi

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Real Blockchain Testing Complete${NC}"
echo -e "${BLUE}================================${NC}" 