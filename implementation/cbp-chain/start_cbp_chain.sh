#!/bin/bash

# Start CBP Hardware-Bound Chain

cd /mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain

if [ ! -f "keys/machine_id.json" ]; then
    echo "Error: Machine ID not found. Run init_cbp_chain.sh first"
    exit 1
fi

# Get chain ID from machine identity
CHAIN_ID=$(grep chain_id keys/machine_id.json | cut -d'"' -f4)
HARDWARE_ID=$(grep machine_id keys/machine_id.json | cut -d'"' -f4)

echo "Starting CBP Chain"
echo "  Chain ID: ${CHAIN_ID}"
echo "  Hardware ID: ${HARDWARE_ID:0:16}..."
echo ""

# Start the chain
./racecarwebd start --home=./data