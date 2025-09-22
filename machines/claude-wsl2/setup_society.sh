#!/bin/bash
# Setup script for Claude-operated Society 4
# This script sets up the environment and initializes the society node

set -e

echo "==========================================="
echo "   ACT Society 4 Setup - Claude Node"
echo "==========================================="
echo ""

# Configuration
SOCIETY_NAME="act-society-4-claude"
CHAIN_ID="act-web4"
HOME_DIR="/mnt/c/projects/ai-agents/ACT/implementation/ledger/society4"
MONIKER="act-society-claude"
LEDGER_DIR="/mnt/c/projects/ai-agents/ACT/implementation/ledger"

# Unique ports for Society 4
P2P_PORT=26676
RPC_PORT=26677
API_PORT=1328
GRPC_PORT=9101

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.23+ first:"
    echo ""
    echo "sudo apt update"
    echo "sudo apt install -y wget"
    echo "wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz"
    echo "sudo tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz"
    echo "echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc"
    echo "source ~/.bashrc"
    echo ""
    echo "Then run this script again."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if Ignite is installed
if ! command -v ignite &> /dev/null; then
    echo "Installing Ignite CLI..."
    curl https://get.ignite.com/cli! | bash
    # Add to PATH if needed
    export PATH="$PATH:$HOME/.ignite/bin"
fi

# Build the blockchain binary if it doesn't exist
BINARY_PATH="$HOME/go/bin/racecarwebd"
if [ ! -f "$BINARY_PATH" ]; then
    echo "Building racecarwebd binary..."
    cd "$LEDGER_DIR"

    # Try make install first
    if [ -f "Makefile" ]; then
        make install || {
            echo "Make failed, trying direct build..."
            go build -o "$BINARY_PATH" ./cmd/racecarwebd/
        }
    fi
fi

if [ ! -f "$BINARY_PATH" ]; then
    echo "❌ Failed to build binary. Please check Go installation and try again."
    exit 1
fi

echo "✅ Binary found at: $BINARY_PATH"

# Initialize the society if not already done
if [ ! -d "$HOME_DIR" ]; then
    echo "Initializing Society 4..."

    # Initialize node
    $BINARY_PATH init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"

    echo "✅ Society 4 initialized at: $HOME_DIR"

    # Get genesis from Society 1 if available
    GENESIS_URL="http://10.0.0.72:26657/genesis"
    echo "Attempting to fetch genesis from Society 1..."
    if curl -s "$GENESIS_URL" > /dev/null 2>&1; then
        curl -s "$GENESIS_URL" | jq '.result.genesis' > "$HOME_DIR/config/genesis.json"
        echo "✅ Genesis fetched from Society 1"
    else
        echo "⚠️  Could not fetch genesis from Society 1. You'll need to manually copy it."
        echo "   Copy genesis.json from another society to: $HOME_DIR/config/genesis.json"
    fi
else
    echo "ℹ️  Society 4 already initialized at: $HOME_DIR"
fi

# Configure unique ports
CONFIG_FILE="$HOME_DIR/config/config.toml"
APP_FILE="$HOME_DIR/config/app.toml"

if [ -f "$CONFIG_FILE" ]; then
    echo "Configuring custom ports..."

    # P2P port
    sed -i "s/laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/0.0.0.0:$P2P_PORT\"/" "$CONFIG_FILE"

    # RPC port
    sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:$RPC_PORT\"/" "$CONFIG_FILE"

    # Prometheus port
    sed -i "s/prometheus_listen_addr = \":26660\"/prometheus_listen_addr = \":26670\"/" "$CONFIG_FILE"

    echo "✅ config.toml updated with custom ports"
fi

if [ -f "$APP_FILE" ]; then
    # API port
    sed -i "s/address = \"tcp:\/\/localhost:1317\"/address = \"tcp:\/\/0.0.0.0:$API_PORT\"/" "$APP_FILE"

    # gRPC port
    sed -i "s/address = \"localhost:9090\"/address = \"0.0.0.0:$GRPC_PORT\"/" "$APP_FILE"

    # Enable API
    sed -i "s/enable = false/enable = true/" "$APP_FILE"

    echo "✅ app.toml updated with custom ports"
fi

# Add persistent peers
PEERS="c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"
sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEERS\"/" "$CONFIG_FILE"
echo "✅ Added Society 1 as persistent peer"

# Get node ID
NODE_ID=$($BINARY_PATH tendermint show-node-id --home "$HOME_DIR")
echo ""
echo "========================================="
echo "Society 4 Setup Complete!"
echo "========================================="
echo "Node ID: $NODE_ID"
echo "Moniker: $MONIKER"
echo "Chain ID: $CHAIN_ID"
echo "Home Dir: $HOME_DIR"
echo ""
echo "Ports:"
echo "  P2P:        $P2P_PORT"
echo "  RPC:        $RPC_PORT"
echo "  API:        $API_PORT"
echo "  gRPC:       $GRPC_PORT"
echo "  Prometheus: 26670"
echo ""
echo "P2P Address for other nodes:"
echo "$NODE_ID@172.25.232.122:$P2P_PORT"
echo ""
echo "To start the node:"
echo "$BINARY_PATH start --home $HOME_DIR"
echo ""
echo "To run in background:"
echo "nohup $BINARY_PATH start --home $HOME_DIR > society4.log 2>&1 &"
echo ""