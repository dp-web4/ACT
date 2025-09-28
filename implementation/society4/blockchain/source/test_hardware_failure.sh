#!/bin/bash

echo "=== Hardware Validation Failure Test ==="
echo ""

# Test what happens when hardware binding is modified

echo "1. Current blockchain status:"
CURRENT_HEIGHT=$(curl -s http://localhost:26657/status | grep -o '"latest_block_height":"[0-9]*"' | cut -d'"' -f4)
echo "   Block height: $CURRENT_HEIGHT"

echo ""
echo "2. Current hardware binding:"
cat $HOME/.society4chain/hardware_binding.json | grep hardware_hash | head -1

echo ""
echo "3. Creating modified hardware binding (simulating different machine)..."
cp $HOME/.society4chain/hardware_binding.json $HOME/.society4chain/hardware_binding.backup

# Modify the Windows UUID to simulate different hardware
sed -i 's/"windows_uuid": "[^"]*"/"windows_uuid": "MODIFIED-UUID-TEST-0000-000000000000"/' $HOME/.society4chain/hardware_binding.json

echo "   Modified Windows UUID in hardware binding"
cat $HOME/.society4chain/hardware_binding.json | grep windows_uuid

echo ""
echo "4. Testing if blockchain detects mismatch..."
echo "   (In production, this would halt the chain)"

# Check if chain is still producing blocks
sleep 5
NEW_HEIGHT=$(curl -s http://localhost:26657/status 2>/dev/null | grep -o '"latest_block_height":"[0-9]*"' | cut -d'"' -f4 || echo "0")

if [ "$NEW_HEIGHT" -gt "$CURRENT_HEIGHT" ]; then
    echo "   ✓ Chain still running (blocks: $CURRENT_HEIGHT -> $NEW_HEIGHT)"
    echo "   Note: Hardware validation is currently in log-only mode"
else
    echo "   ✗ Chain stopped or slowed"
fi

echo ""
echo "5. Restoring original hardware binding..."
mv $HOME/.society4chain/hardware_binding.backup $HOME/.society4chain/hardware_binding.json
echo "   ✓ Original hardware binding restored"

echo ""
echo "=== Test Complete ==="
echo "Hardware validation code is implemented but not enforced by default."
echo "To enable enforcement, modify app.go to panic on validation failure."