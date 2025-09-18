#!/bin/bash
# Test the complete ATP/ADP energy cycle

echo "=== ACT Energy Cycle Test ==="
echo "Testing society pool with mint/discharge/recharge operations"
echo

# Set up aliases for easier commands
CHAIN_ID="racecarweb"
NODE="tcp://localhost:26657"
KEYRING="test"

# Test accounts (assuming alice has been set up)
ALICE="cosmos1pmnw5epy2zflns3lrmzxgcm86zsly7nr24jqrp"

echo "1. Minting ADP tokens to society pool..."
echo "   Society: society:demo"
echo "   Treasury Role: society:demo:role:treasury"
echo "   Amount: 1000000 ADP"
echo

racecar-webd tx energycycle mint-adp 1000000 society:demo society:demo:role:treasury "Genesis allocation for demo society" \
    --from alice \
    --chain-id $CHAIN_ID \
    --node $NODE \
    --keyring-backend $KEYRING \
    --yes

sleep 5

echo
echo "2. Discharging ATP for work..."
echo "   Worker: society:demo:citizen:alice"
echo "   Amount: 100 ATP"
echo "   Work: Building Web4 infrastructure"
echo

racecar-webd tx energycycle discharge-atp society:demo:citizen:alice 100 "Building Web4 infrastructure" society:demo \
    --from alice \
    --chain-id $CHAIN_ID \
    --node $NODE \
    --keyring-backend $KEYRING \
    --yes

sleep 5

echo
echo "3. Recharging ADP to ATP with energy input..."
echo "   Producer: society:demo:role:energy-producer"
echo "   Amount: 50 ADP"
echo "   Energy Source: solar"
echo

racecar-webd tx energycycle recharge-adp society:demo:role:energy-producer 50 solar "Solar farm output validation #001" \
    --from alice \
    --chain-id $CHAIN_ID \
    --node $NODE \
    --keyring-backend $KEYRING \
    --yes

sleep 5

echo
echo "4. Querying society pool balance..."
echo

# Query the society pool (when query endpoint is implemented)
# For now, we'll check the events
echo "Checking transaction events..."
racecar-webd query txs --query "message.action='/racecarweb.energycycle.v1.MsgMintADP'" --limit 10 --output json | jq '.txs[-1].logs[0].events[] | select(.type=="society_pool_mint")'

echo
echo "=== Test Complete ==="
echo "Check the events above to verify:"
echo "- ADP was minted to society pool"
echo "- ATP was discharged for work (creating ADP)"
echo "- ADP was recharged to ATP with energy"