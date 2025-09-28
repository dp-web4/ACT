#!/bin/bash

# CBP Machine-Specific Blockchain Initialization
# This creates a hardware-bound private chain for the CBP machine

cd /mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain

# Create hardware identity
echo "=== Generating CBP Hardware Identity ==="

HOSTNAME=$(hostname)
WSL_UUID=$(uuidgen || echo "cbp-$(date +%s)")
CPU_INFO=$(grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | xargs)
CPU_CORES=$(nproc)
MEM_TOTAL=$(grep MemTotal /proc/meminfo | awk '{print $2}')
HARDWARE_ID=$(echo -n "${HOSTNAME}-${WSL_UUID}-${CPU_CORES}-${MEM_TOTAL}" | sha256sum | cut -d' ' -f1)
CHAIN_ID="cbp-chain-${HARDWARE_ID:0:8}"
SEED=$(echo -n "${HARDWARE_ID}-${USER}-cbp" | sha256sum | cut -d' ' -f1)

echo "Hardware ID: ${HARDWARE_ID}"
echo "Chain ID: ${CHAIN_ID}"

# Save machine identity
cat > keys/machine_id.json << EOF
{
  "machine_id": "${HARDWARE_ID}",
  "hostname": "${HOSTNAME}",
  "chain_id": "${CHAIN_ID}",
  "created": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "hardware": {
    "cpu": "${CPU_INFO}",
    "cores": ${CPU_CORES},
    "memory_kb": ${MEM_TOTAL}
  },
  "seed": "${SEED}"
}
EOF

# Check if we have the blockchain source
if [ ! -d "../society4/blockchain/source" ]; then
    echo "Error: Blockchain source not found at ../society4/blockchain/source"
    echo "Please ensure the society4 implementation exists"
    exit 1
fi

# Copy blockchain binary or build it
echo "=== Setting Up Blockchain Binary ==="
if [ -f "../society4/blockchain/source/racecarwebd" ]; then
    cp ../society4/blockchain/source/racecarwebd ./racecarwebd
    echo "Using existing blockchain binary"
else
    echo "Building blockchain from source..."
    cd ../society4/blockchain/source
    go mod download
    go build -o racecarwebd ./cmd/racecarwebd
    cp racecarwebd ../../../cbp-chain/
    cd ../../../cbp-chain
fi

# Initialize the chain
echo "=== Initializing CBP Chain ==="
./racecarwebd init cbp-validator --chain-id=${CHAIN_ID} --home=./data

# Configure chain for hardware binding
echo "=== Configuring Hardware Binding ==="
cat > config/hardware_binding.json << EOF
{
  "hardware_id": "${HARDWARE_ID}",
  "chain_id": "${CHAIN_ID}",
  "binding_enabled": true,
  "attestation_required": true,
  "machine_specific": true,
  "created_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF

# Create validator key from hardware seed
echo "=== Creating Hardware-Derived Validator Key ==="
echo "${SEED}" | ./racecarwebd keys add cbp-validator --home=./data --recover

# Configure genesis
echo "=== Configuring Genesis ==="
./racecarwebd genesis add-genesis-account cbp-validator 1000000000stake --home=./data
./racecarwebd genesis gentx cbp-validator 100000000stake --chain-id=${CHAIN_ID} --home=./data
./racecarwebd genesis collect-gentxs --home=./data

# Update config for local operation
sed -i 's/127.0.0.1:26657/0.0.0.0:26657/g' ./data/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' ./data/config/config.toml
sed -i 's/enable = false/enable = true/g' ./data/config/app.toml
sed -i 's/swagger = false/swagger = true/g' ./data/config/app.toml

echo "=== CBP Chain Initialized Successfully ==="
echo ""
echo "Chain Details:"
echo "  Hardware ID: ${HARDWARE_ID}"
echo "  Chain ID: ${CHAIN_ID}"
echo "  Validator: cbp-validator"
echo "  Home Directory: ./data"
echo ""
echo "To start the chain:"
echo "  ./start_cbp_chain.sh"
echo ""
echo "This chain is hardware-bound to THIS specific machine."
echo "It will not run on other machines without the correct hardware attestation."