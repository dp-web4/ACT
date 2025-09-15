#!/bin/bash

# Test script to verify blockchain integration fixes
echo "🧪 Testing Blockchain Integration Fixes"
echo "======================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_BRIDGE_URL="http://localhost:8080"
BLOCKCHAIN_URL="http://localhost:1317"

echo ""
echo "1. Testing API Bridge Health..."
if curl -s "$API_BRIDGE_URL/health" > /dev/null; then
    echo -e "${GREEN}✅ API Bridge is running${NC}"
else
    echo -e "${RED}❌ API Bridge is not running${NC}"
    echo "   Start the API Bridge with: ./api-bridge"
    exit 1
fi

echo ""
echo "2. Testing Blockchain Status Endpoint..."
BLOCKCHAIN_STATUS=$(curl -s "$API_BRIDGE_URL/blockchain/status")
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Blockchain status endpoint responded${NC}"
    echo "   Response: $BLOCKCHAIN_STATUS"
else
    echo -e "${RED}❌ Blockchain status endpoint failed${NC}"
fi

echo ""
echo "3. Testing Direct Blockchain Connection..."
if curl -s "$BLOCKCHAIN_URL/cosmos/base/tendermint/v1beta1/node_info" > /dev/null; then
    echo -e "${GREEN}✅ Blockchain is accessible directly${NC}"
else
    echo -e "${YELLOW}⚠️  Blockchain is not accessible directly${NC}"
    echo "   Make sure to start the blockchain with: ignite chain serve"
fi

echo ""
echo "4. Testing Component Registration..."
COMPONENT_RESPONSE=$(curl -s -X POST "$API_BRIDGE_URL/api/v1/components/register" \
  -H "Content-Type: application/json" \
  -d '{
    "creator": "alice",
    "component_data": "test-battery-module-v1",
    "context": "integration-test"
  }')

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Component registration endpoint responded${NC}"
    echo "   Response: $COMPONENT_RESPONSE"
    
    # Check if it's a real transaction or mock
    if echo "$COMPONENT_RESPONSE" | grep -q "mock_tx"; then
        echo -e "${YELLOW}⚠️  Using mock transaction (blockchain integration may need debugging)${NC}"
    else
        echo -e "${GREEN}✅ Real blockchain transaction detected!${NC}"
    fi
else
    echo -e "${RED}❌ Component registration failed${NC}"
fi

echo ""
echo "🎯 Test Summary:"
echo "================"
echo "• API Bridge: $(curl -s "$API_BRIDGE_URL/health" > /dev/null && echo -e "${GREEN}✅ Running${NC}" || echo -e "${RED}❌ Not Running${NC}")"
echo "• Blockchain Direct: $(curl -s "$BLOCKCHAIN_URL/cosmos/base/tendermint/v1beta1/node_info" > /dev/null && echo -e "${GREEN}✅ Accessible${NC}" || echo -e "${YELLOW}⚠️  Not Accessible${NC}")"
echo "• Component Registration: $(curl -s -X POST "$API_BRIDGE_URL/api/v1/components/register" -H "Content-Type: application/json" -d '{"creator":"test","component_data":"test","context":"test"}' > /dev/null && echo -e "${GREEN}✅ Working${NC}" || echo -e "${YELLOW}⚠️  Failed${NC}")"

echo ""
echo "📋 Next Steps:"
echo "=============="
echo "1. If blockchain is not accessible, start it with: ignite chain serve"
echo "2. Check the blockchain status endpoint for detailed path information"
echo "3. Look at API Bridge logs for Ignite CLI and path detection messages"
echo "4. If still using mock transactions, check the logs for specific error messages"

echo ""
echo "🔧 Debugging:"
echo "============="
echo "• Check API Bridge logs: ./api-bridge 2>&1 | grep -i 'path\|ignite\|blockchain'"
echo "• Test Ignite CLI manually: ignite version"
echo "• Test racecar-webd manually: racecar-webd version"
echo "• Check blockchain status: curl http://localhost:8080/blockchain/status"
