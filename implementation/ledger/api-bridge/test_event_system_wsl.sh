#!/bin/bash

# Test script for the event system (WSL version)
# This script demonstrates how to test the event emission system in WSL

set -e

echo "=== Event System Test (WSL) ==="
echo "This script will:"
echo "1. Create a test configuration with events enabled"
echo "2. Build and start the API bridge"
echo "3. Test component registration with event emission"
echo "4. Show expected events"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if we're in WSL
check_wsl() {
    print_status "Checking WSL environment..."
    
    if grep -qi microsoft /proc/version; then
        print_status "WSL environment detected"
        WSL_VERSION=$(grep -o "WSL[0-9]*" /proc/version | head -1)
        print_status "WSL Version: $WSL_VERSION"
    else
        print_warning "Not running in WSL - this script is optimized for WSL"
    fi
}

# Check if required tools are available
check_requirements() {
    print_status "Checking requirements..."
    
    # Check for curl
    if ! command -v curl &> /dev/null; then
        print_error "curl is required but not installed"
        print_status "Installing curl..."
        sudo apt update && sudo apt install -y curl
    fi
    
    # Check for jq
    if ! command -v jq &> /dev/null; then
        print_warning "jq is not installed - installing for better JSON formatting"
        sudo apt install -y jq
    fi
    
    # Check for Go
    if ! command -v go &> /dev/null; then
        print_error "Go is required but not installed"
        print_status "Please install Go: https://golang.org/doc/install"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Go version: $GO_VERSION"
    
    print_status "Requirements check passed"
}

# Create test configuration with events enabled
create_test_config() {
    print_step "Creating test configuration with events enabled..."
    
    cat > test_config.yaml << 'EOF'
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

# Start simple webhook receiver using Python
start_webhook_receiver() {
    print_step "Starting webhook receiver on port 3000..."
    
    # Create a simple Python webhook receiver
    cat > webhook_receiver.py << 'EOF'
#!/usr/bin/env python3
import http.server
import socketserver
import json
import datetime
from urllib.parse import urlparse, parse_qs

class WebhookHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        try:
            event = json.loads(post_data.decode('utf-8'))
            
            # Log the event
            timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
            print(f"\n🔔 [{timestamp}] Event received:")
            print(f"   Type: {event.get('event_type', 'unknown')}")
            print(f"   Timestamp: {event.get('timestamp', 'unknown')}")
            
            if 'data' in event:
                data = event['data']
                print(f"   Data:")
                for key, value in data.items():
                    print(f"     {key}: {value}")
            
            # Send response
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.send_header('Access-Control-Allow-Origin', '*')
            self.end_headers()
            self.wfile.write(json.dumps({"status": "received"}).encode())
            
        except Exception as e:
            print(f"Error processing event: {e}")
            self.send_response(500)
            self.end_headers()
    
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b"<h1>Webhook Receiver</h1><p>Ready to receive events</p>")
    
    def log_message(self, format, *args):
        # Suppress default logging
        pass

if __name__ == "__main__":
    PORT = 3000
    with socketserver.TCPServer(("", PORT), WebhookHandler) as httpd:
        print(f"Webhook receiver running on port {PORT}")
        print("Press Ctrl+C to stop")
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nShutting down webhook receiver...")
EOF
    
    # Make the script executable
    chmod +x webhook_receiver.py
    
    # Start webhook receiver in background
    python3 webhook_receiver.py &
    WEBHOOK_PID=$!
    
    # Wait for webhook receiver to start
    sleep 3
    
    # Test if webhook receiver is working
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        print_status "Webhook receiver started successfully"
    else
        print_error "Failed to start webhook receiver"
        exit 1
    fi
}

# Build the API bridge
build_api_bridge() {
    print_step "Building API bridge..."
    
    # Check if we're in the right directory
    if [ ! -f "main.go" ]; then
        print_error "main.go not found. Please run this script from the api-bridge directory"
        exit 1
    fi
    
    # Build the API bridge
    if ! go build -o api-bridge-test main.go; then
        print_error "Failed to build API bridge"
        exit 1
    fi
    
    print_status "API bridge built successfully"
}

# Start API bridge with test configuration
start_api_bridge() {
    print_step "Starting API bridge with event system enabled..."
    
    # Start API bridge in background
    ./api-bridge-test -config test_config.yaml &
    API_BRIDGE_PID=$!
    
    # Wait for API bridge to start
    print_status "Waiting for API bridge to start..."
    sleep 5
    
    # Test if API bridge is running
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        print_status "API bridge started successfully"
    else
        print_error "Failed to start API bridge"
        print_status "Checking if port 8081 is available..."
        if netstat -tuln | grep -q ":8081 "; then
            print_error "Port 8081 is already in use"
        fi
        exit 1
    fi
}

# Test component registration with event emission
test_component_registration() {
    print_step "Testing component registration with event emission..."
    
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
    print_status "Waiting for event processing..."
    sleep 3
}

# Test pairing initiation with event emission
test_pairing_initiation() {
    print_step "Testing pairing initiation with event emission..."
    
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
    print_status "Waiting for event processing..."
    sleep 3
}

# Test LCT creation with event emission
test_lct_creation() {
    print_step "Testing LCT creation with event emission..."
    
    # Create LCT
    RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/lct/create \
        -H "Content-Type: application/json" \
        -d '{
            "creator": "test-user",
            "component_a": "COMP-test-user-123",
            "component_b": "COMP-test-user-456",
            "context": "event-system-test"
        }')
    
    if [ $? -eq 0 ]; then
        print_status "LCT creation successful"
        if command -v jq &> /dev/null; then
            echo "$RESPONSE" | jq .
        else
            echo "$RESPONSE"
        fi
    else
        print_error "LCT creation failed"
        return 1
    fi
    
    # Wait a moment for event processing
    print_status "Waiting for event processing..."
    sleep 3
}

# Show received events summary
show_received_events() {
    print_step "Event System Summary"
    
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
    echo "3. lct_created event:"
    echo "   - Endpoint: http://localhost:3000/webhooks/lct-created"
    echo "   - Contains: lct_id, creator, component_a, component_b, context, timestamp, tx_hash"
    echo ""
    echo "Note: Check the webhook receiver output above for actual received events"
    echo ""
    echo "In a real deployment, you would:"
    echo "- Set up proper webhook receivers"
    echo "- Store events in SQL databases"
    echo "- Implement authentication and security"
    echo "- Add monitoring and alerting"
}

# Show API bridge logs
show_api_bridge_logs() {
    print_step "API Bridge Logs"
    echo ""
    echo "Recent API bridge activity (if any):"
    echo "----------------------------------------"
    # The logs would be visible in the terminal where the API bridge is running
    echo "Check the terminal output above for API bridge logs"
    echo "Look for messages like:"
    echo "- 'Event POSTed successfully'"
    echo "- 'Event queue enabled and worker started'"
    echo "- 'Failed to POST event, will retry'"
    echo "----------------------------------------"
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    
    # Stop API bridge
    if [ ! -z "$API_BRIDGE_PID" ]; then
        print_status "Stopping API bridge..."
        kill $API_BRIDGE_PID 2>/dev/null || true
        sleep 2
    fi
    
    # Stop webhook receiver
    if [ ! -z "$WEBHOOK_PID" ]; then
        print_status "Stopping webhook receiver..."
        kill $WEBHOOK_PID 2>/dev/null || true
        sleep 1
    fi
    
    # Clean up test files
    print_status "Removing test files..."
    rm -f api-bridge-test test_config.yaml webhook_receiver.py
    
    print_status "Cleanup completed"
}

# Main execution
main() {
    # Set up cleanup on exit
    trap cleanup EXIT
    
    check_wsl
    check_requirements
    create_test_config
    start_webhook_receiver
    build_api_bridge
    start_api_bridge
    
    echo ""
    print_status "Testing event system..."
    echo ""
    
    test_component_registration
    test_pairing_initiation
    test_lct_creation
    
    echo ""
    show_received_events
    show_api_bridge_logs
    
    echo ""
    print_status "Event system test completed!"
    print_status "Check the webhook receiver output above for received events"
    echo ""
    print_warning "Press Ctrl+C to stop the test and cleanup"
    
    # Keep running to show logs
    wait
}

# Run main function
main "$@" 