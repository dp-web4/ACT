#!/bin/bash

# Society2 - Alternative Web4 Society Implementation
# This creates a second, independent society with different governance and characteristics

set -e

SOCIETY_NAME="society2"
CHAIN_ID="web4-society2-001"
MONIKER="society2-validator"
HOME_DIR="./data"

echo "======================================"
echo "   Initializing Society2 Blockchain  "
echo "======================================"
echo ""
echo "Society Name: ${SOCIETY_NAME}"
echo "Chain ID: ${CHAIN_ID}"
echo ""

cd /mnt/c/exe/projects/ai-agents/ACT/implementation/society2

# Extract unique hardware fingerprint for society2
echo "=== Generating Society2 Hardware Identity ==="
HOSTNAME=$(hostname)
TIMESTAMP=$(date +%s)
HW_UUID=$(uuidgen || echo "society2-${TIMESTAMP}")
CPU_INFO=$(grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | xargs)
CPU_CORES=$(nproc)
MEM_TOTAL=$(grep MemTotal /proc/meminfo | awk '{print $2}')

# Create unique society2 hardware ID (different from society4)
HARDWARE_ID=$(echo -n "society2-${HOSTNAME}-${HW_UUID}-${CPU_CORES}" | sha256sum | cut -d' ' -f1)
SEED=$(echo -n "${HARDWARE_ID}-society2-${USER}" | sha256sum | cut -d' ' -f1)

echo "Hardware ID: ${HARDWARE_ID}"
echo "Seed: ${SEED:0:16}..."

# Save society2 identity
mkdir -p keys
cat > keys/society2_identity.json << EOF
{
  "society_name": "${SOCIETY_NAME}",
  "chain_id": "${CHAIN_ID}",
  "hardware_id": "${HARDWARE_ID}",
  "created": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "characteristics": {
    "governance": "democratic",
    "consensus": "proof-of-contribution",
    "energy_model": "renewable",
    "federation": "enabled",
    "interchain": true
  },
  "network": {
    "p2p_port": 26656,
    "rpc_port": 26657,
    "api_port": 1317,
    "grpc_port": 9090
  },
  "seed": "${SEED}"
}
EOF

# Create society2 specific configuration
echo "=== Configuring Society2 Governance ==="
mkdir -p laws
cat > laws/society2_constitution.md << 'EOF'
# Society2 Constitutional Laws

## Core Principles
1. **Democratic Governance**: All citizens have equal voting rights
2. **Proof of Contribution**: Consensus based on value creation
3. **Open Federation**: Welcomes inter-society collaboration
4. **Energy Sustainability**: Prioritizes renewable ATP sources

## Governance Model
- **Voting**: 1 citizen = 1 vote
- **Proposals**: Any citizen can propose changes
- **Quorum**: 51% participation required
- **Passage**: 60% approval needed

## Inter-Society Relations
- Automatic federation with compatible societies
- Cross-chain atomic swaps enabled
- Shared trust tensor calculations
- Energy trading permitted

## Differences from Society4
- Society4: Hierarchical roles → Society2: Flat democracy
- Society4: Proof of stake → Society2: Proof of contribution
- Society4: Isolated → Society2: Federation-first
EOF

# Create custom genesis configuration
echo "=== Building Society2 Genesis Configuration ==="
cat > config/genesis_template.json << EOF
{
  "genesis_time": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "chain_id": "${CHAIN_ID}",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1"
    }
  },
  "app_state": {
    "society": {
      "name": "${SOCIETY_NAME}",
      "type": "democratic",
      "federation_enabled": true,
      "energy_model": "renewable"
    },
    "lctmanager": {
      "society_lct": {
        "id": "society2-root-lct",
        "hardware_bound": true,
        "hardware_id": "${HARDWARE_ID}"
      }
    },
    "energycycle": {
      "society_pool": {
        "atp_balance": "1000000",
        "adp_balance": "0",
        "renewable_sources": true
      }
    },
    "trusttensor": {
      "initial_trust": {
        "self_trust": 1.0,
        "federation_trust": 0.5
      }
    }
  }
}
EOF

# Build the blockchain binary with society2 customizations
echo "=== Building Society2 Blockchain ==="
cd blockchain/source

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21+"
    exit 1
fi

# Build with society2 tag
echo "Building blockchain binary..."
go build -tags society2 -o society2d ./app || {
    echo "Build failed, trying without tags..."
    go build -o society2d ./app
}

if [ ! -f society2d ]; then
    echo "Error: Failed to build society2d binary"
    exit 1
fi

echo "Binary built successfully: society2d"

# Initialize the chain
echo "=== Initializing Society2 Chain ==="
./society2d init ${MONIKER} --chain-id=${CHAIN_ID} --home=${HOME_DIR}

# Create validator key from seed
echo "=== Creating Society2 Validator ==="
echo "${SEED}" | ./society2d keys add ${MONIKER} --home=${HOME_DIR} --recover 2>/dev/null || {
    # If recovery fails, create new key
    ./society2d keys add ${MONIKER} --home=${HOME_DIR}
}

# Add genesis account
./society2d genesis add-genesis-account ${MONIKER} 1000000000stake --home=${HOME_DIR}

# Create genesis transaction
./society2d genesis gentx ${MONIKER} 100000000stake \
    --chain-id=${CHAIN_ID} \
    --moniker=${MONIKER} \
    --home=${HOME_DIR}

# Collect genesis transactions
./society2d genesis collect-gentxs --home=${HOME_DIR}

# Configure for inter-chain communication
echo "=== Configuring Inter-Chain Communication ==="
CONFIG_FILE="${HOME_DIR}/config/config.toml"
APP_CONFIG="${HOME_DIR}/config/app.toml"

# Update ports to avoid conflicts with society4
sed -i 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/127.0.0.1:26558"/' ${CONFIG_FILE}
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26557"/' ${CONFIG_FILE}
sed -i 's/laddr = "tcp:\/\/0.0.0.0:26656"/laddr = "tcp:\/\/0.0.0.0:26556"/' ${CONFIG_FILE}
sed -i 's/pprof_laddr = "localhost:6060"/pprof_laddr = "localhost:6061"/' ${CONFIG_FILE}

# Enable CORS for inter-chain
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/' ${CONFIG_FILE}

# Update app ports
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1217"/' ${APP_CONFIG}
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9091"/' ${APP_CONFIG}
sed -i 's/address = "localhost:9091"/address = "localhost:9092"/' ${APP_CONFIG}

# Enable API and swagger
sed -i 's/enable = false/enable = true/' ${APP_CONFIG}
sed -i 's/swagger = false/swagger = true/' ${APP_CONFIG}

echo "=== Society2 Initialization Complete ==="
echo ""
echo "Society2 Configuration:"
echo "  Name: ${SOCIETY_NAME}"
echo "  Chain ID: ${CHAIN_ID}"
echo "  Governance: Democratic"
echo "  Federation: Enabled"
echo "  Hardware ID: ${HARDWARE_ID:0:16}..."
echo ""
echo "Network Ports (different from society4):"
echo "  P2P: 26556"
echo "  RPC: 26557"
echo "  API: 1217"
echo "  gRPC: 9091"
echo ""
echo "To start Society2:"
echo "  cd blockchain/source"
echo "  ./society2d start --home=${HOME_DIR}"
echo ""
echo "To connect Society2 with Society4:"
echo "  Use the federation bridge configuration"
echo "  See: federation_config.json"