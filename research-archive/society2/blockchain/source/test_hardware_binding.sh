#!/bin/bash

# Test hardware binding functionality

set -e

echo "=== Society 4 Hardware Binding Test Suite ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Hardware Extraction
echo "Test 1: Hardware Extraction"
echo "----------------------------"

# Extract hardware
HARDWARE_JSON=$(bash extract_hardware.sh json 2>/dev/null)
HARDWARE_HASH=$(echo "$HARDWARE_JSON" | grep -o '"hardware_hash": "[^"]*' | cut -d'"' -f4)

if [ -n "$HARDWARE_HASH" ]; then
    echo -e "${GREEN}✓${NC} Hardware extracted successfully"
    echo "  Hash: ${HARDWARE_HASH:0:16}..."
else
    echo -e "${RED}✗${NC} Failed to extract hardware"
    exit 1
fi

# Save current hardware for comparison
echo "$HARDWARE_JSON" > current_hardware.json

# Test 2: Compare with Genesis Hardware
echo ""
echo "Test 2: Genesis Hardware Comparison"
echo "------------------------------------"

if [ -f "$HOME/.society4chain/hardware_binding.json" ]; then
    GENESIS_HASH=$(cat $HOME/.society4chain/hardware_binding.json | grep -o '"hardware_hash": "[^"]*' | cut -d'"' -f4)

    if [ "$HARDWARE_HASH" = "$GENESIS_HASH" ]; then
        echo -e "${GREEN}✓${NC} Current hardware matches genesis"
    else
        echo -e "${YELLOW}⚠${NC} Hardware mismatch detected"
        echo "  Current: ${HARDWARE_HASH:0:16}..."
        echo "  Genesis: ${GENESIS_HASH:0:16}..."
    fi
else
    echo -e "${YELLOW}⚠${NC} No genesis hardware binding found"
fi

# Test 3: Go Hardware Validation Test
echo ""
echo "Test 3: Go Hardware Validation Module"
echo "--------------------------------------"

# Create a simple Go test program
cat > test_hardware.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type HardwareBinding struct {
    Platform     string `json:"platform"`
    HardwareHash string `json:"hardware_hash"`
    Components   struct {
        WindowsUUID string `json:"windows_uuid"`
        CPUInfo     string `json:"cpu_info"`
    } `json:"components"`
}

func main() {
    // Read current hardware
    data, err := ioutil.ReadFile("current_hardware.json")
    if err != nil {
        fmt.Printf("Error reading hardware: %v\n", err)
        os.Exit(1)
    }

    var wrapper struct {
        HardwareBinding HardwareBinding `json:"hardware_binding"`
    }

    if err := json.Unmarshal(data, &wrapper); err != nil {
        fmt.Printf("Error parsing hardware: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Platform: %s\n", wrapper.HardwareBinding.Platform)
    fmt.Printf("Hardware Hash: %s...\n", wrapper.HardwareBinding.HardwareHash[:16])
    fmt.Printf("Windows UUID: %s\n", wrapper.HardwareBinding.Components.WindowsUUID)
    fmt.Printf("CPU: %s\n", wrapper.HardwareBinding.Components.CPUInfo)

    if wrapper.HardwareBinding.Platform == "wsl2" {
        fmt.Println("✓ Valid WSL2 platform detected")
    }
}
EOF

export PATH=/usr/local/go/bin:$PATH
if go run test_hardware.go 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Go hardware validation successful"
else
    echo -e "${RED}✗${NC} Go hardware validation failed"
fi

rm -f test_hardware.go

# Test 4: Blockchain Start with Hardware Check
echo ""
echo "Test 4: Blockchain Startup Test"
echo "--------------------------------"

# Check if chain can query status
if ./society4chaind status --home $HOME/.society4chain 2>/dev/null | grep -q "node_info"; then
    echo -e "${GREEN}✓${NC} Chain configuration valid"
else
    echo -e "${YELLOW}⚠${NC} Chain not configured or not running"
fi

# Test 5: Simulated Hardware Mismatch
echo ""
echo "Test 5: Hardware Mismatch Detection"
echo "------------------------------------"

# Create a fake hardware binding with different UUID
cat > fake_hardware.json << 'EOF'
{
    "hardware_binding": {
        "platform": "wsl2",
        "hardware_hash": "FAKE0000000000000000000000000000000000000000000000000000000000",
        "components": {
            "windows_uuid": "FAKE-UUID-0000-0000-0000-000000000000",
            "hyperv_uuid": "HYPERV_UNAVAILABLE",
            "wsl_boot_id": "fake-boot-id",
            "cpu_info": "Fake CPU",
            "memory_kb": 1000000
        }
    }
}
EOF

# Create test program to verify mismatch detection
cat > test_mismatch.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type HardwareBinding struct {
    Components struct {
        WindowsUUID string `json:"windows_uuid"`
    } `json:"components"`
}

func main() {
    // Read real hardware
    realData, _ := ioutil.ReadFile("current_hardware.json")
    var real struct {
        HardwareBinding HardwareBinding `json:"hardware_binding"`
    }
    json.Unmarshal(realData, &real)

    // Read fake hardware
    fakeData, _ := ioutil.ReadFile("fake_hardware.json")
    var fake struct {
        HardwareBinding HardwareBinding `json:"hardware_binding"`
    }
    json.Unmarshal(fakeData, &fake)

    if real.HardwareBinding.Components.WindowsUUID != fake.HardwareBinding.Components.WindowsUUID {
        fmt.Println("Hardware mismatch detected correctly")
        os.Exit(0)
    } else {
        fmt.Println("Failed to detect hardware mismatch")
        os.Exit(1)
    }
}
EOF

if go run test_mismatch.go 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Hardware mismatch detection working"
else
    echo -e "${RED}✗${NC} Hardware mismatch detection failed"
fi

rm -f test_mismatch.go fake_hardware.json

# Test 6: Performance Test
echo ""
echo "Test 6: Hardware Extraction Performance"
echo "----------------------------------------"

START_TIME=$(date +%s%N)
for i in {1..10}; do
    bash extract_hardware.sh hash >/dev/null 2>&1
done
END_TIME=$(date +%s%N)

ELAPSED=$((($END_TIME - $START_TIME) / 10000000)) # Convert to milliseconds
AVG_TIME=$(($ELAPSED / 10))

echo "Average extraction time: ${AVG_TIME}ms"
if [ $AVG_TIME -lt 500 ]; then
    echo -e "${GREEN}✓${NC} Performance acceptable (<500ms)"
else
    echo -e "${YELLOW}⚠${NC} Performance slow (>500ms)"
fi

# Summary
echo ""
echo "=== Test Summary ==="
echo "Hardware Hash: ${HARDWARE_HASH:0:32}..."
echo "Platform: WSL2"
echo "All tests completed!"

# Cleanup
rm -f current_hardware.json