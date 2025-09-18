#!/bin/bash

# Genesis Society Creation Script
# Creates the foundational society with role-based LCTs

set -e

echo "=== Creating Genesis Society for ACT ==="

# Configuration
CHAIN_ID="racecar-web"
NODE="http://localhost:26657"
KEYRING="test"
GAS_PRICES="0.025stake"

# Check if chain is running
if ! curl -s "$NODE/status" > /dev/null 2>&1; then
    echo "Error: Chain not running at $NODE"
    echo "Start with: racecar-webd start --pruning nothing"
    exit 1
fi

echo "Chain is running. Creating society structure..."

# Step 1: Create Society LCT with 'society' genesis role
echo "1. Creating Society LCT..."
SOCIETY_TX=$(racecar-webd tx lctmanager mint-lct \
    --entity-name "ACT-Demo-Society" \
    --entity-type "society" \
    --metadata "genesis_date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    --metadata "founding_principle=all_decisions_unanimous" \
    --metadata "metabolic_state=active" \
    --initial-adp "0" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

echo "Society LCT transaction submitted"
sleep 5

# Get the Society LCT ID from events (we'll parse this from the response)
# For now, we'll use a predictable ID format
SOCIETY_LCT="lct-society-act-demo-001"
echo "Society LCT: $SOCIETY_LCT"

# Step 2: Create Treasury Role LCT
echo "2. Creating Treasury Role LCT..."
TREASURY_TX=$(racecar-webd tx lctmanager mint-lct \
    --entity-name "Treasury-Role" \
    --entity-type "role" \
    --metadata "role_type=treasury" \
    --metadata "parent_society=$SOCIETY_LCT" \
    --metadata "permissions=mint_adp,charge_atp,manage_pool" \
    --initial-adp "0" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

TREASURY_LCT="lct-role-treasury-001"
echo "Treasury Role LCT: $TREASURY_LCT"
sleep 5

# Step 3: Create Witness Role LCT
echo "3. Creating Witness Role LCT..."
WITNESS_TX=$(racecar-webd tx lctmanager mint-lct \
    --entity-name "Witness-Role" \
    --entity-type "role" \
    --metadata "role_type=witness" \
    --metadata "parent_society=$SOCIETY_LCT" \
    --metadata "permissions=witness_work,validate_claims,update_trust" \
    --initial-adp "0" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

WITNESS_LCT="lct-role-witness-001"
echo "Witness Role LCT: $WITNESS_LCT"
sleep 5

# Step 4: Create Citizen Role LCT
echo "4. Creating Citizen Role LCT..."
CITIZEN_TX=$(racecar-webd tx lctmanager mint-lct \
    --entity-name "Citizen-Role" \
    --entity-type "role" \
    --metadata "role_type=citizen" \
    --metadata "parent_society=$SOCIETY_LCT" \
    --metadata "permissions=request_allocation,perform_work,receive_witness" \
    --initial-adp "0" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

CITIZEN_LCT="lct-role-citizen-001"
echo "Citizen Role LCT: $CITIZEN_LCT"
sleep 5

# Step 5: Create Law Oracle Role LCT
echo "5. Creating Law Oracle Role LCT..."
LAW_TX=$(racecar-webd tx lctmanager mint-lct \
    --entity-name "Law-Oracle-Role" \
    --entity-type "role" \
    --metadata "role_type=law_oracle" \
    --metadata "parent_society=$SOCIETY_LCT" \
    --metadata "permissions=publish_laws,validate_compliance,enforce_rules" \
    --initial-adp "0" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

LAW_LCT="lct-role-law-oracle-001"
echo "Law Oracle Role LCT: $LAW_LCT"
sleep 5

# Step 6: Now bind all roles to the Society LCT
echo "6. Binding roles to Society..."
# Note: We need to implement the binding mechanism in the chain
# For now, we'll record these relationships in metadata

# Step 7: Treasury mints initial ADP pool
echo "7. Treasury Role minting initial ADP pool (1,000,000 ADP)..."
MINT_TX=$(racecar-webd tx energycycle mint-a-d-p \
    "1000000" \
    "$SOCIETY_LCT" \
    "$TREASURY_LCT" \
    "Genesis allocation for society treasury" \
    --from validator \
    --chain-id $CHAIN_ID \
    --gas-prices $GAS_PRICES \
    --keyring-backend $KEYRING \
    --yes \
    --output json)

echo "Initial ADP pool minted"
sleep 5

# Step 8: Query the society state
echo ""
echo "=== Society Genesis Complete ==="
echo "Society LCT: $SOCIETY_LCT"
echo "Treasury Role: $TREASURY_LCT"
echo "Witness Role: $WITNESS_LCT"
echo "Citizen Role: $CITIZEN_LCT"
echo "Law Oracle Role: $LAW_LCT"
echo "Initial ADP Pool: 1,000,000 ADP"
echo ""
echo "Next steps:"
echo "1. Create citizen LCTs for human and agents"
echo "2. Bind citizens to citizen role"
echo "3. Create tasks that allocate ATP"
echo "4. Discharge ATP through work completion"

# Save the society configuration
cat > society_config.json << EOF
{
  "society_lct": "$SOCIETY_LCT",
  "roles": {
    "treasury": "$TREASURY_LCT",
    "witness": "$WITNESS_LCT",
    "citizen": "$CITIZEN_LCT",
    "law_oracle": "$LAW_LCT"
  },
  "initial_pool": {
    "adp": 1000000,
    "atp": 0
  },
  "laws": [
    "all_decisions_unanimous"
  ],
  "metabolic_state": "active",
  "genesis_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
EOF

echo "Society configuration saved to society_config.json"