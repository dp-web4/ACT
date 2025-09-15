#!/bin/bash

# Web4-ModBatt Test Coverage Generator
# This script generates comprehensive coverage reports for all modules

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
echo -e "${BLUE}  Web4-ModBatt Coverage Generator${NC}"
echo -e "${BLUE}================================${NC}"
echo "Timestamp: $TIMESTAMP"
echo ""

# Function to print section headers
print_section() {
    echo -e "${YELLOW}$1${NC}"
    echo "----------------------------------------"
}

# Function to generate coverage for a module
generate_module_coverage() {
    local module_name=$1
    local test_path=$2
    local coverage_file="$COVERAGE_DIR/${module_name}_coverage.out"
    local html_file="$COVERAGE_DIR/${module_name}_coverage.html"
    local log_file="$LOG_DIR/${module_name}_coverage_${TIMESTAMP}.log"
    
    echo "Generating coverage for $module_name..."
    echo "Coverage file: $coverage_file"
    echo "HTML file: $html_file"
    echo "Log file: $log_file"
    
    # Run tests with coverage
    go test -v -coverprofile="$coverage_file" -covermode=atomic "$test_path" 2>&1 | tee "$log_file"
    
    # Check exit status
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ $module_name tests passed${NC}"
        
        # Generate HTML coverage report
        if [ -f "$coverage_file" ]; then
            go tool cover -html="$coverage_file" -o="$html_file"
            echo -e "${GREEN}✓ HTML coverage report generated: $html_file${NC}"
            
            # Show coverage summary
            echo "Coverage Summary for $module_name:"
            go tool cover -func="$coverage_file" | tail -1
        fi
    else
        echo -e "${RED}✗ $module_name tests failed${NC}"
        return 1
    fi
    
    echo ""
}

# Function to generate integration coverage
generate_integration_coverage() {
    local test_name=$1
    local test_path=$2
    local coverage_file="$COVERAGE_DIR/integration_${test_name}_coverage.out"
    local html_file="$COVERAGE_DIR/integration_${test_name}_coverage.html"
    local log_file="$LOG_DIR/integration_${test_name}_coverage_${TIMESTAMP}.log"
    
    echo "Generating integration coverage for $test_name..."
    echo "Coverage file: $coverage_file"
    echo "HTML file: $html_file"
    echo "Log file: $log_file"
    
    # Run integration tests with coverage
    go test -v -coverprofile="$coverage_file" -covermode=atomic "$test_path" 2>&1 | tee "$log_file"
    
    # Check exit status
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        echo -e "${GREEN}✓ Integration test $test_name passed${NC}"
        
        # Generate HTML coverage report
        if [ -f "$coverage_file" ]; then
            go tool cover -html="$coverage_file" -o="$html_file"
            echo -e "${GREEN}✓ HTML coverage report generated: $html_file${NC}"
            
            # Show coverage summary
            echo "Coverage Summary for $test_name:"
            go tool cover -func="$coverage_file" | tail -1
        fi
    else
        echo -e "${RED}✗ Integration test $test_name failed${NC}"
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

# Install coverage tools
print_section "Installing Coverage Tools"
go install golang.org/x/tools/cmd/cover@latest
go install github.com/axw/gocov/gocov@latest
go install github.com/axw/gocov-xml@latest
echo -e "${GREEN}✓ Coverage tools installed${NC}"
echo ""

# Generate unit test coverage
print_section "Generating Unit Test Coverage"

# Component Registry coverage
generate_module_coverage "componentregistry" "$TEST_DIR/unit/componentregistry"

# LCT Manager coverage
generate_module_coverage "lctmanager" "$TEST_DIR/unit/lctmanager"

# Pairing coverage
generate_module_coverage "pairing" "$TEST_DIR/unit/pairing"

# Pairing Queue coverage
generate_module_coverage "pairingqueue" "$TEST_DIR/unit/pairingqueue"

# Trust Tensor coverage
generate_module_coverage "trusttensor" "$TEST_DIR/unit/trusttensor"

# Energy Cycle coverage
generate_module_coverage "energycycle" "$TEST_DIR/unit/energycycle"

echo -e "${GREEN}✓ All unit test coverage generated${NC}"
echo ""

# Generate integration test coverage
print_section "Generating Integration Test Coverage"

# Pairing flow integration coverage
generate_integration_coverage "pairing_flow" "$TEST_DIR/integration/pairing_flow"

# Energy operations integration coverage
generate_integration_coverage "energy_operations" "$TEST_DIR/integration/energy_operations"

# Trust calculations integration coverage
generate_integration_coverage "trust_calculations" "$TEST_DIR/integration/trust_calculations"

# Offline scenarios integration coverage
generate_integration_coverage "offline_scenarios" "$TEST_DIR/integration/offline_scenarios"

echo -e "${GREEN}✓ All integration test coverage generated${NC}"
echo ""

# Generate combined coverage report
print_section "Generating Combined Coverage Report"

# Create combined coverage file
combined_coverage="$COVERAGE_DIR/combined_coverage_${TIMESTAMP}.out"
echo "mode: atomic" > "$combined_coverage"

# Append all coverage files
echo "Combining coverage files..."
for coverage_file in $COVERAGE_DIR/*_coverage.out; do
    if [ -f "$coverage_file" ]; then
        echo "Adding: $coverage_file"
        tail -n +2 "$coverage_file" >> "$combined_coverage"
    fi
done

# Generate combined HTML report
combined_html="$COVERAGE_DIR/combined_coverage_${TIMESTAMP}.html"
go tool cover -html="$combined_coverage" -o="$combined_html"

echo -e "${GREEN}✓ Combined coverage report generated: $combined_html${NC}"

# Show overall coverage summary
echo ""
echo "Overall Coverage Summary:"
go tool cover -func="$combined_coverage" | tail -1

echo ""

# Generate coverage statistics
print_section "Coverage Statistics"

# Calculate coverage percentages
echo "Coverage Statistics by Module:"
echo ""

for coverage_file in $COVERAGE_DIR/*_coverage.out; do
    if [ -f "$coverage_file" ]; then
        module_name=$(basename "$coverage_file" _coverage.out)
        echo "Module: $module_name"
        go tool cover -func="$coverage_file" | tail -1
        echo ""
    fi
done

# Generate coverage trend data
print_section "Coverage Trend Analysis"

# Create coverage trend file
trend_file="$COVERAGE_DIR/coverage_trend_${TIMESTAMP}.json"
echo "{" > "$trend_file"
echo "  \"timestamp\": \"$TIMESTAMP\"," >> "$trend_file"
echo "  \"modules\": {" >> "$trend_file"

first=true
for coverage_file in $COVERAGE_DIR/*_coverage.out; do
    if [ -f "$coverage_file" ]; then
        module_name=$(basename "$coverage_file" _coverage.out)
        coverage_percent=$(go tool cover -func="$coverage_file" | tail -1 | awk '{print $3}' | sed 's/%//')
        
        if [ "$first" = true ]; then
            first=false
        else
            echo "," >> "$trend_file"
        fi
        
        echo "    \"$module_name\": $coverage_percent" >> "$trend_file"
    fi
done

echo "  }" >> "$trend_file"
echo "}" >> "$trend_file"

echo -e "${GREEN}✓ Coverage trend data saved: $trend_file${NC}"

# Generate coverage badges
print_section "Generating Coverage Badges"

# Create badges directory
mkdir -p "$COVERAGE_DIR/badges"

# Generate badges for each module
for coverage_file in $COVERAGE_DIR/*_coverage.out; do
    if [ -f "$coverage_file" ]; then
        module_name=$(basename "$coverage_file" _coverage.out)
        coverage_percent=$(go tool cover -func="$coverage_file" | tail -1 | awk '{print $3}' | sed 's/%//')
        
        # Determine badge color based on coverage
        if (( $(echo "$coverage_percent >= 90" | bc -l) )); then
            color="brightgreen"
        elif (( $(echo "$coverage_percent >= 80" | bc -l) )); then
            color="green"
        elif (( $(echo "$coverage_percent >= 70" | bc -l) )); then
            color="yellow"
        elif (( $(echo "$coverage_percent >= 60" | bc -l) )); then
            color="orange"
        else
            color="red"
        fi
        
        # Generate badge URL
        badge_url="https://img.shields.io/badge/coverage-$coverage_percent%25-$color"
        badge_file="$COVERAGE_DIR/badges/${module_name}_coverage.svg"
        
        # Download badge (if curl is available)
        if command -v curl &> /dev/null; then
            curl -s "$badge_url" -o "$badge_file"
            echo -e "${GREEN}✓ Badge generated for $module_name: $badge_file${NC}"
        else
            echo -e "${YELLOW}Warning: curl not available, skipping badge generation for $module_name${NC}"
        fi
    fi
done

# Generate overall coverage badge
overall_coverage=$(go tool cover -func="$combined_coverage" | tail -1 | awk '{print $3}' | sed 's/%//')

if (( $(echo "$overall_coverage >= 90" | bc -l) )); then
    color="brightgreen"
elif (( $(echo "$overall_coverage >= 80" | bc -l) )); then
    color="green"
elif (( $(echo "$overall_coverage >= 70" | bc -l) )); then
    color="yellow"
elif (( $(echo "$overall_coverage >= 60" | bc -l) )); then
    color="orange"
else
    color="red"
fi

badge_url="https://img.shields.io/badge/overall_coverage-$overall_coverage%25-$color"
badge_file="$COVERAGE_DIR/badges/overall_coverage.svg"

if command -v curl &> /dev/null; then
    curl -s "$badge_url" -o "$badge_file"
    echo -e "${GREEN}✓ Overall coverage badge generated: $badge_file${NC}"
fi

echo ""

# Generate coverage report summary
print_section "Coverage Report Summary"

echo "Coverage Reports Generated:"
echo "- Combined HTML Report: $combined_html"
echo "- Combined Coverage Data: $combined_coverage"
echo "- Coverage Trend Data: $trend_file"
echo "- Individual Module Reports: $COVERAGE_DIR/*_coverage.html"
echo "- Coverage Badges: $COVERAGE_DIR/badges/"
echo "- Coverage Logs: $LOG_DIR/"
echo "- Timestamp: $TIMESTAMP"

echo ""
echo "Coverage Targets:"
echo "- Unit Tests: 90% (target)"
echo "- Integration Tests: 85% (target)"
echo "- Overall Coverage: 87% (target)"
echo "- Critical Paths: 100% (required)"

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  Coverage Generation Completed${NC}"
echo -e "${GREEN}================================${NC}"

# Optional: Open combined coverage report in browser
if command -v xdg-open &> /dev/null && [ -f "$combined_html" ]; then
    echo ""
    read -p "Open combined coverage report in browser? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        xdg-open "$combined_html"
    fi
elif command -v open &> /dev/null && [ -f "$combined_html" ]; then
    echo ""
    read -p "Open combined coverage report in browser? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open "$combined_html"
    fi
fi

echo ""
echo "Coverage generation completed at: $(date)" 