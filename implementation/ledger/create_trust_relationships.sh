#!/bin/bash
# Create role-contextual T3/V3 trust relationships via RDF-style graph
# Trust is not a property but a relationship between entities

export PATH=/usr/local/go/bin:$PATH

echo "=== Creating Role-Contextual Trust Relationships ==="
echo "Trust and value are relationships, not properties"
echo ""

# Create trust relationships between queens and society
echo "Society trusts Genesis Queen for orchestration..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Build-Society" \
  --to-lct "ACT-Genesis-Queen" \
  --tensor-type "T3" \
  --context "orchestration" \
  --talent "0.9" \
  --training "0.8" \
  --temperament "0.9" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo "Genesis Queen trusts Web4 Alignment Queen for pattern recognition..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Genesis-Queen" \
  --to-lct "Web4-Alignment-Queen" \
  --tensor-type "T3" \
  --context "web4-alignment" \
  --talent "0.95" \
  --training "0.9" \
  --temperament "1.0" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo "Genesis Queen trusts Reality Alignment Queen for assumption checking..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Genesis-Queen" \
  --to-lct "Reality-Alignment-Queen" \
  --tensor-type "T3" \
  --context "reality-validation" \
  --talent "0.8" \
  --training "0.7" \
  --temperament "0.95" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# Queens trust their workers
echo "Web4 Alignment Queen trusts Pattern Validator for pattern matching..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "Web4-Alignment-Queen" \
  --to-lct "Pattern-Validator" \
  --tensor-type "T3" \
  --context "pattern-validation" \
  --talent "0.8" \
  --training "0.7" \
  --temperament "0.9" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo "LCT Infrastructure Queen trusts LCT Coder for implementation..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "LCT-Infrastructure-Queen" \
  --to-lct "LCT-Coder" \
  --tensor-type "T3" \
  --context "code-implementation" \
  --talent "0.6" \
  --training "0.8" \
  --temperament "0.7" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

# Value relationships (V3)
echo ""
echo "=== Creating Value Relationships ==="

echo "Society values creation from LCT Coder..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Build-Society" \
  --to-lct "LCT-Coder" \
  --tensor-type "V3" \
  --context "code-creation" \
  --valuation "0.8" \
  --veracity "0.7" \
  --validity "0.9" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo "Society values coordination from Genesis Queen..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Build-Society" \
  --to-lct "ACT-Genesis-Queen" \
  --tensor-type "V3" \
  --context "coordination" \
  --valuation "0.7" \
  --veracity "0.8" \
  --validity "0.9" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo "Genesis Queen values alignment verification from Web4 Alignment Queen..."
~/go/bin/racecar-webd tx trusttensor create-relationship-tensor \
  --from-lct "ACT-Genesis-Queen" \
  --to-lct "Web4-Alignment-Queen" \
  --tensor-type "V3" \
  --context "alignment-verification" \
  --valuation "0.9" \
  --veracity "0.95" \
  --validity "0.9" \
  --from act-society-treasury \
  --keyring-backend test \
  --home ./act-society \
  --chain-id act-web4 \
  --fees 100stake \
  -y

sleep 3

echo ""
echo "=== Trust Graph Structure Created ==="
echo ""
echo "Key principles:"
echo "1. Trust is directional: A trusts B ≠ B trusts A"
echo "2. Trust is contextual: Trust for coding ≠ Trust for decisions"
echo "3. Trust evolves: Based on witnessed performance"
echo "4. No absolute trust: Only relationships"
echo ""
echo "RDF-style graph:"
echo "<Society> --[trusts-for-orchestration]--> <Genesis-Queen>"
echo "<Genesis-Queen> --[trusts-for-alignment]--> <Web4-Alignment-Queen>"
echo "<Web4-Alignment-Queen> --[trusts-for-patterns]--> <Pattern-Validator>"
echo ""
echo "This creates a proper Web4 trust network where:"
echo "- Roles trust each other contextually"
echo "- Value flows along trust edges"
echo "- ATP discharges when trust is exercised"
echo "- ADP recharges when value is delivered"