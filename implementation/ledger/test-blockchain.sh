#!/bin/bash
# ACT Blockchain Test Script

echo "Starting ACT blockchain test..."

# Build the chain
cd /mnt/c/exe/projects/ai-agents/ACT/implementation/ledger
echo "Building blockchain..."
go build -o actd ./cmd/actd

# Initialize chain
echo "Initializing chain..."
./actd init test-validator --chain-id act-testnet-1

# Copy our genesis
cp config/genesis.json ~/.act/config/genesis.json

# Start chain in background
echo "Starting chain..."
./actd start &
CHAIN_PID=$!

# Wait for chain to start
sleep 5

# Test transactions
echo "Testing LCT mint transaction..."
./actd tx lctmanager mint-lct AGENT --from validator --chain-id act-testnet-1 --yes

echo "Testing ATP discharge..."
./actd tx energycycle discharge-atp 100 --from validator --chain-id act-testnet-1 --yes

# Query state
echo "Querying LCTs..."
./actd query lctmanager list-lct

echo "Querying energy pools..."
./actd query energycycle list-pool

# Stop chain
kill $CHAIN_PID

echo "Test complete!"
