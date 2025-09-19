#!/bin/bash
# Create Web4-aligned LCTs for ACT Society and all role entities
# Roles are first-class entities in Web4, each gets an LCT

export PATH=/usr/local/go/bin:$PATH

echo "=== Creating LCTs for ACT Society and Role Entities ==="

# Initialize blockchain if not running
~/go/bin/racecar-webd init "act-society" --chain-id "act-web4" --home ./act-society

# Create society account (owns all tokens)
~/go/bin/racecar-webd keys add act-society-treasury --keyring-backend test --home ./act-society

# Genesis funding - society owns all tokens
~/go/bin/racecar-webd genesis add-genesis-account act-society-treasury 100000000000stake,10000000atp,10000000adp --keyring-backend test --home ./act-society

# Create validator
~/go/bin/racecar-webd genesis gentx act-society-treasury 1000000000stake \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --moniker "act-society-node"

~/go/bin/racecar-webd genesis collect-gentxs --home ./act-society
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' ./act-society/config/app.toml

echo "Starting blockchain..."
~/go/bin/racecar-webd start --home ./act-society --api.enable --grpc.enable &
CHAIN_PID=$!
sleep 10

echo ""
echo "=== Creating Society LCT ==="
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Build-Society" \
  --entity-type "society" \
  --initial-adp-amount "10000000" \
  --initial-t3-tensor "[0.8,0.9,0.95]" \
  --initial-v3-tensor "[0.1,0.1,0.8]" \
  --metadata purpose="Build ACT through fractal swarm intelligence" \
  --metadata laws="All decisions unanimous until consensus governance" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 1000stake \
  -y

sleep 5

echo ""
echo "=== Creating Genesis Queen LCT ==="
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACT-Genesis-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.9,0.8,0.9]" \
  --initial-v3-tensor "[0.3,0.3,0.4]" \
  --metadata role="meta-orchestrator" \
  --metadata domain="overall-vision" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo ""
echo "=== Creating Domain Queen LCTs ==="

# Web4 Alignment Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Web4-Alignment-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.95,0.9,1.0]" \
  --initial-v3-tensor "[0.4,0.5,0.1]" \
  --metadata role="alignment-validator" \
  --metadata domain="web4-harmony" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# Reality Alignment Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Reality-Alignment-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.8,0.7,0.95]" \
  --initial-v3-tensor "[0.5,0.4,0.1]" \
  --metadata role="assumption-checker" \
  --metadata domain="reality-verification" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# LCT Infrastructure Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "LCT-Infrastructure-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.7,0.9,0.8]" \
  --initial-v3-tensor "[0.6,0.2,0.2]" \
  --metadata role="identity-architect" \
  --metadata domain="lct-systems" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# ACP Protocol Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ACP-Protocol-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.8,0.85,0.8]" \
  --initial-v3-tensor "[0.6,0.2,0.2]" \
  --metadata role="protocol-designer" \
  --metadata domain="agent-context" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# Demo Society Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Demo-Society-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.7,0.8,0.85]" \
  --initial-v3-tensor "[0.7,0.2,0.1]" \
  --metadata role="society-builder" \
  --metadata domain="governance" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# MCP Bridge Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "MCP-Bridge-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.7,0.85,0.8]" \
  --initial-v3-tensor "[0.6,0.3,0.1]" \
  --metadata role="integration-coordinator" \
  --metadata domain="mcp-bridges" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# Client Interface Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Client-Interface-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.6,0.8,0.8]" \
  --initial-v3-tensor "[0.7,0.2,0.1]" \
  --metadata role="ui-architect" \
  --metadata domain="user-interface" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# ATP Economy Queen
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "ATP-Economy-Queen" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.9,0.8,0.85]" \
  --initial-v3-tensor "[0.5,0.3,0.2]" \
  --metadata role="economy-modeler" \
  --metadata domain="energy-economy" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo ""
echo "=== Creating Sample Worker Role LCTs ==="
echo "Note: Full worker roster would be ~30+ roles, creating representatives..."

# Pattern Validator (Web4 Alignment Worker)
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "Pattern-Validator" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.8,0.7,0.9]" \
  --initial-v3-tensor "[0.6,0.3,0.1]" \
  --metadata role="worker" \
  --metadata queen="Web4-Alignment-Queen" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# LCT Coder (Infrastructure Worker)
~/go/bin/racecar-webd tx lctmanager mint-lct \
  --entity-name "LCT-Coder" \
  --entity-type "role" \
  --initial-adp-amount "0" \
  --initial-t3-tensor "[0.6,0.8,0.7]" \
  --initial-v3-tensor "[0.8,0.1,0.1]" \
  --metadata role="worker" \
  --metadata queen="LCT-Infrastructure-Queen" \
  --metadata filled_by="Claude" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo ""
echo "=== Society LCT Structure Created ==="
echo "Society: ACT-Build-Society (owns all tokens)"
echo "Genesis Queen: Meta-orchestrator role"
echo "9 Domain Queens: Each with specialized domain"
echo "Sample Workers: Representative roles created"
echo ""
echo "All roles are entities filled by Claude instances"
echo "Other agents can be invited to fill roles when ready"
echo ""
echo "T3 Tensors reflect: [Talent, Training, Temperament]"
echo "V3 Tensors reflect: [Valuation, Veracity, Validity]"
echo "- High storage (0.8) for society = accumulating value"
echo "- Balanced creation/exchange for queens = active work"
echo "- High creation for workers = producing output"
echo ""
echo "Blockchain PID: $CHAIN_PID"
echo "Stop with: kill $CHAIN_PID"