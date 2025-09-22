#!/bin/bash
# Build script for Sprout (Jetson Orin Nano)
# Machine-specific build configuration

# Load machine config
MACHINE_CONFIG="$(dirname "$0")/machine-config.json"
ACT_ROOT="/home/sprout/ai-workspace/ACT"
LEDGER_DIR="${ACT_ROOT}/implementation/ledger"
GO_PATH="/usr/local/go/bin"
GO_BIN_DIR="/home/sprout/go/bin"

# Ensure Go 1.24 is in PATH
export PATH="${GO_PATH}:${PATH}"

echo "=== Building ACT for Sprout (Jetson Orin Nano) ==="
echo "Machine: ARM64 architecture"
echo "Go version: $(go version)"

# Navigate to ledger directory
cd "${LEDGER_DIR}" || exit 1

# Build the blockchain
echo "Building racecar-webd..."
if [ -f "Makefile" ]; then
    make install
else
    # Direct build if Makefile fails
    go build -o "${GO_BIN_DIR}/racecar-webd" ./cmd/racecar-webd || {
        echo "Build failed. Attempting Ignite build..."
        ignite chain build --skip-proto
    }
fi

# Check if binary exists
if [ -f "${GO_BIN_DIR}/racecar-webd" ]; then
    echo "✅ Build successful! Binary at: ${GO_BIN_DIR}/racecar-webd"
else
    echo "❌ Build failed. Binary not found."
    exit 1
fi

echo "Build complete for Sprout machine"