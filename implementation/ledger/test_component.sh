#!/bin/bash

# Register a new component
echo "Registering component..."
./racecar-webd tx componentregistry register-component test-module-005 module '{"manufacturer_id": "tesla"}' --from alice --chain-id racecarweb --yes

echo ""
echo "Waiting for transaction to be committed..."
sleep 5

echo ""
echo "Testing improved get-component query (should return full component details)..."
./racecar-webd q componentregistry get-component test-module-005 --output json

echo ""
echo "Testing with invalid component ID (should fail)..."
./racecar-webd q componentregistry get-component invalid@component --output json