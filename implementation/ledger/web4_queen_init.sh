#!/bin/bash
# Web4 Society Initialization with Alignment Queen
# This creates a proper Web4 society with attention-based roles

export PATH=/usr/local/go/bin:$PATH

echo "=== Initializing Web4 Society with Alignment Queen ==="

# Initialize the blockchain
~/go/bin/racecar-webd init "web4-society" --chain-id "web4" --home ./web4

# Create the Web4 Alignment Queen role
# The queen is not a person but an attention partition that manages society alignment
~/go/bin/racecar-webd keys add web4-queen --keyring-backend test --home ./web4

# Create the society treasury (holds all tokens, not individuals)
~/go/bin/racecar-webd keys add society-treasury --keyring-backend test --home ./web4

# Genesis accounts - society owns tokens, roles just have execution rights
~/go/bin/racecar-webd genesis add-genesis-account web4-queen 1000000000stake --keyring-backend test --home ./web4
~/go/bin/racecar-webd genesis add-genesis-account society-treasury 10000000000stake,1000000atp,1000000adp --keyring-backend test --home ./web4

# Create validator for consensus (queen maintains attention)
~/go/bin/racecar-webd genesis gentx web4-queen 100000000stake \
  --keyring-backend test \
  --home ./web4 \
  --chain-id web4 \
  --moniker "alignment-queen" \
  --details "Web4 society alignment and attention partition manager"

# Collect genesis transactions
~/go/bin/racecar-webd genesis collect-gentxs --home ./web4

# Configure for Web4 protocol
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./web4/config/app.toml

# Set metabolic state awareness
echo "Configuring metabolic states for digital organism..."
sed -i 's/timeout_commit = "5s"/timeout_commit = "3s"/' ./web4/config/config.toml

echo "=== Web4 Society Ready ==="
echo "Queen Role: Attention partition for society alignment"
echo "Treasury: Society-owned token pools (no individual ownership)"
echo "Metabolic: 33% readiness economy enabled"
echo ""
echo "Start with: ~/go/bin/racecar-webd start --home ./web4"