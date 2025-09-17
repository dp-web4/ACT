#!/bin/bash

# Fix imports in all Go files
find . -name "*.go" -type f | while read file; do
    # Update standard cosmos-sdk module imports
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/auth|"cosmossdk.io/x/auth|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/authz|"cosmossdk.io/x/authz|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/bank|"cosmossdk.io/x/bank|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/consensus|"cosmossdk.io/x/consensus|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/distribution|"cosmossdk.io/x/distribution|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/epochs|"cosmossdk.io/x/epochs|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/gov|"cosmossdk.io/x/gov|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/group|"cosmossdk.io/x/group|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/mint|"cosmossdk.io/x/mint|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/params|"cosmossdk.io/x/params|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/slashing|"cosmossdk.io/x/slashing|g' "$file"
    sed -i 's|"github.com/cosmos/cosmos-sdk/x/staking|"cosmossdk.io/x/staking|g' "$file"
    
    # Keep genutil as it is (not moved to cosmossdk.io)
    # Keep simapp and other non-x/ modules as they are
done

echo "Import paths updated"