#!/bin/bash

# Test script for the event system
# This script demonstrates how to test the event emission system

set -e

echo "=== Event System Test ==="
echo "This script will:"
echo "1. Start a mock webhook receiver"
echo "2. Configure the API bridge to emit events"
echo "3. Test component registration with event emission"
echo "4. Show received events"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are available
check_requirements() {
    print_status "Checking requirements..."
    
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        print_warning "jq is not installed - JSON output will not be formatted"
    fi
    
    print_status "Requirements check passed"
}

# Start mock webhook receiver
start_webhook_receiver() {
    print_status "Starting mock webhook receiver on port 3000..."
    
    # Create a simple webhook receiver using netcat
    # This will listen for POST requests and log them
    (
        echo "HTTP/1.1 200 OK"
        echo "Content-Type: application/json"
        echo "Access-Control-Allow-Origin: *"
        echo ""
        echo '{"status": "received"}'
    ) | nc -l 3000 &
    
    WEBHOOK_PID=$!
    sleep 2
    
    # Test if webhook receiver is working
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        print_status "Webhook receiver started successfully"
    else
        print_error "Failed to start webhook receiver"
        exit 1
    fi
}

# Create test configuration with events enabled
create_test_config() {
    print_status "Creating test configuration with events enabled..."
    
    cat > test_config.yaml << EOF
blockchain:
  rest_endpoint: "http://0.0.0.0:1317"
  grpc_endpoint: "localhost:9090"
  chain_id: "racecarweb"
  timeout: 30

server:
  port: 8081  # Use different port to avoid conflicts
  host: "0.0.0.0"
  read_timeout: 30
  write_timeout: 30

logging:
  level: "info"
  format: "json"

events:
  enabled: true
  max_retries: 3
  retry_delay: 2
  queue_size: 1000
  endpoints:
    component_registered:
      - "http://localhost:3000/webhooks/component-registered"
    component_verified:
      - "http://localhost:3000/webhooks/component-verified"
    pairing_initiated:
      - "http://localhost:3000/webhooks/pairing-initiated"
    pairing_completed:
      - "http://localhost:3000/webhooks/pairing-completed"
    lct_created:
      - "http://localhost:3000/webhooks/lct-created"
    trust_tensor_created:
      - "http://localhost:3000/webhooks/trust-tensor-created"
    energy_transfer:
      - "http://localhost:3000/webhooks/energy-transfer"
EOF
    
    print_status "Test configuration created: test_config.yaml"
}

# Start API bridge with test configuration
start_api_bridge() {
    print_status "Starting API bridge with event system enabled..."
    
    # Build the API bridge
    if ! go build -o api-bridge-test main.go; then
        print_error "Failed to build API bridge"
        exit 1
    fi
    
    # Start API bridge in background
    ./api-bridge-test -config test_config.yaml &
    API_BRIDGE_PID=$!
    
    # Wait for API bridge to start
    sleep 5
    
    # Test if API bridge is running
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        print_status "API bridge started successfully"
    else
        print_error "Failed to start API bridge"
        exit 1
    fi
}

# Test component registration with event emission
test_component_registration() {
    print_status "Testing component registration with event emission..."
    
    # Register a test component
    RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/components/register \
        -H "Content-Type: application/json" \
        -d '{
            "creator": "test-user",
            "component_data": "test-battery-module",
            "context": "event-system-test"
        }')
    
    if [ $? -eq 0 ]; then
        print_status "Component registration successful"
        if command -v jq &> /dev/null; then
            echo "$RESPONSE" | jq .
        else
            echo "$RESPONSE"
        fi
    else
        print_error "Component registration failed"
        return 1
    fi
    
    # Wait a moment for event processing
    sleep 3
}

# Test pairing initiation with event emission
test_pairing_initiation() {
    print_status "Testing pairing initiation with event emission..."
    
    # Initiate pairing
    RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/pairing/initiate \
        -H "Content-Type: application/json" \
        -d '{
            "creator": "test-user",
            "component_a": "COMP-test-user-123",
            "component_b": "COMP-test-user-456",
            "operational_context": "event-system-test",
            "force_immediate": true
        }')
    
    if [ $? -eq 0 ]; then
        print_status "Pairing initiation successful"
        if command -v jq &> /dev/null; then
            echo "$RESPONSE" | jq .
        else
            echo "$RESPONSE"
        fi
    else
        print_error "Pairing initiation failed"
        return 1
    fi
    
    # Wait a moment for event processing
    sleep 3
}

# Show received events (mock implementation)
show_received_events() {
    print_status "Checking for received events..."
    
    # In a real implementation, you would check the webhook receiver logs
    # For this demo, we'll just show what events should have been sent
    
    echo ""
    echo "=== Expected Events ==="
    echo "The following events should have been sent to webhook endpoints:"
    echo ""
    echo "1. component_registered event:"
    echo "   - Endpoint: http://localhost:3000/webhooks/component-registered"
    echo "   - Contains: component_id, creator, component_data, context, timestamp, tx_hash"
    echo ""
    echo "2. pairing_initiated event:"
    echo "   - Endpoint: http://localhost:3000/webhooks/pairing-initiated"
    echo "   - Contains: challenge_id, creator, component_a, component_b, operational_context, timestamp, tx_hash"
    echo ""
    echo "Note: In a real deployment, you would:"
    echo "- Set up proper webhook receivers"
    echo "- Store events in SQL databases"
    echo "- Implement authentication and security"
    echo "- Add monitoring and alerting"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    
    if [ ! -z "$API_BRIDGE_PID" ]; then
        kill $API_BRIDGE_PID 2>/dev/null || true
    fi
    
    if [ ! -z "$WEBHOOK_PID" ]; then
        kill $WEBHOOK_PID 2>/dev/null || true
    fi
    
    # Clean up test files
    rm -f api-bridge-test test_config.yaml
    
    print_status "Cleanup completed"
}

# Main execution
main() {
    # Set up cleanup on exit
    trap cleanup EXIT
    
    check_requirements
    start_webhook_receiver
    create_test_config
    start_api_bridge
    
    echo ""
    print_status "Testing event system..."
    echo ""
    
    test_component_registration
    test_pairing_initiation
    
    echo ""
    show_received_events
    
    echo ""
    print_status "Event system test completed!"
    print_status "Check the API bridge logs for event emission details"
    echo ""
    print_warning "Press Ctrl+C to stop the test and cleanup"
    
    # Keep running to show logs
    wait
}

# Run main function
main "$@" 