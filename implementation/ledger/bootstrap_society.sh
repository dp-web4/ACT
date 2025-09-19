#!/bin/bash
# Bootstrap the ACT self-building society
# The society builds itself through role interactions following R6 pattern

export PATH=/usr/local/go/bin:$PATH

echo "=== Bootstrapping ACT Self-Building Society ==="
echo "Purpose: Build ACT by being a Web4 society"
echo ""

# Start the blockchain
echo "Starting blockchain infrastructure..."
cd /home/dp/ai-workspace/act/implementation/ledger
~/go/bin/racecar-webd init "act-society" --chain-id "act-web4" --home ./society &>/dev/null
~/go/bin/racecar-webd keys add act-society-treasury --keyring-backend test --home ./society &>/dev/null
~/go/bin/racecar-webd genesis add-genesis-account act-society-treasury 100000000000stake,10000000atp,10000000adp --keyring-backend test --home ./society &>/dev/null
~/go/bin/racecar-webd genesis gentx act-society-treasury 1000000000stake --keyring-backend test --home ./society --chain-id act-web4 --moniker "act-society-node" &>/dev/null
~/go/bin/racecar-webd genesis collect-gentxs --home ./society &>/dev/null
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./society/config/app.toml

echo "Launching blockchain..."
~/go/bin/racecar-webd start --home ./society --api.enable --grpc.enable &
CHAIN_PID=$!
sleep 10

echo ""
echo "=== Society Genesis Block ==="
echo "Creating ACT-Build-Society LCT..."

# Create the society itself
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Build-Society" \
  --entity-type "society" \
  --initial-adp-amount "10000000" \
  --metadata purpose="Build ACT through self-organizing swarm intelligence" \
  --metadata law1="All actions must follow R6 pattern" \
  --metadata law2="Roles can create new roles with society approval" \
  --metadata law3="Rules can be negotiated through consensus" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society \
  --chain-id act-web4 \
  --fees 1000stake \
  -y &>/dev/null

echo "Society LCT created. Minting initial roles..."

echo ""
echo "=== Creating Founding Roles ==="

# Genesis Queen - The orchestrator
echo "Creating Genesis Queen role..."
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Genesis-Queen" \
  --entity-type "role" \
  --metadata role="orchestrator" \
  --metadata authority="create-roles,assign-tasks,approve-rules" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society \
  --chain-id act-web4 \
  --fees 100stake \
  -y &>/dev/null

# Web4 Alignment Queen - The guardian of R6
echo "Creating Web4 Alignment Queen role..."
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Web4-Alignment-Queen" \
  --entity-type "role" \
  --metadata role="alignment-guardian" \
  --metadata authority="validate-r6,verify-patterns,approve-actions" \
  --metadata r6_pattern="Rules+Roles+Request+Reference+Resource=Result" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./society \
  --chain-id act-web4 \
  --fees 100stake \
  -y &>/dev/null

echo ""
echo "=== Society Bootstrap Complete ==="
echo "Blockchain PID: $CHAIN_PID"
echo ""
echo "Society is now self-governing with:"
echo "- Genesis Queen: Can create roles and assign tasks"
echo "- Web4 Alignment Queen: Ensures R6 compliance"
echo "- Treasury: 10M ATP ready for work"
echo ""
echo "The society can now:"
echo "1. Create new roles as needed"
echo "2. Negotiate new rules through consensus"
echo "3. Execute tasks following R6 pattern"
echo "4. Build ACT through emergent organization"
echo ""
echo "Saving PID for monitoring..."
echo $CHAIN_PID > society.pid