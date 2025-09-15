#!/bin/bash

echo "Testing Component Registry Module..."

# Wait a moment for the chain to be ready
sleep 2

echo "1. Checking available accounts..."
./racecar-webd keys list

echo "2. Checking chain-id..."
./racecar-webd status 2>&1 | grep "chain_id"

echo "3. Registering a test component with correct arguments..."
TXHASH=$(./racecar-webd tx componentregistry register-component \
  "test-component-003" \
  "pack" \
  "test-manufacturer-001" \
  --from alice \
  --chain-id racecarweb \
  --yes | grep "txhash:" | awk '{print $2}')

echo "Transaction hash: $TXHASH"

echo "4. Waiting for transaction to be committed..."
sleep 5

echo "5. Checking transaction status..."
./racecar-webd q tx $TXHASH

echo "6. Querying component details (JSON output)..."
./racecar-webd q componentregistry get-component "test-component-003" -o json

echo "7. Querying component details (text output)..."
./racecar-webd q componentregistry get-component "test-component-003"

echo "8. Testing invalid component ID..."
./racecar-webd q componentregistry get-component "invalid-id"

echo "Test completed!" 