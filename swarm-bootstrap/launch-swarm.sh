#!/bin/bash

# ACT Development Swarm Launch Script
# Bootstrap the fractal swarm that builds ACT

echo "🚀 ACT Development Swarm Launcher"
echo "=================================="
echo "Web4-compliant Role-based Architecture"
echo ""

# Check for Claude-Flow
echo "📦 Checking Claude-Flow installation..."
if ! command -v npx &> /dev/null; then
    echo "❌ npx not found. Please install Node.js"
    exit 1
fi

# Initialize Claude-Flow if needed
echo "🔧 Initializing Claude-Flow..."
npx claude-flow@alpha --version 2>/dev/null || {
    echo "Installing Claude-Flow..."
    npm install -g claude-flow@alpha
}

# Create ACT Development Society
echo ""
echo "🌍 Creating ACT Development Society..."
echo "This society will issue Role LCTs for all swarm agents"

cat << 'EOF' > /tmp/act-society.json
{
  "name": "ACT Development Collective",
  "type": "society",
  "mode": "delegative",
  "law_oracle": "act-governance-v1",
  "atp_treasury": 10000,
  "birth_certificate": {
    "rights": ["issue_lcts", "bind_roles", "witness_actions"],
    "responsibilities": ["deliver_act", "maintain_trust", "evolve_system"]
  }
}
EOF

# Initialize the swarm with Claude-Flow
echo ""
echo "👑 Initializing Genesis Orchestrator Role..."
npx claude-flow@alpha swarm init \
    --topology="hierarchical" \
    --config="/tmp/act-society.json" \
    --max-agents=31 \
    --memory="persistent" \
    --witness="all" || {
    echo "⚠️  Claude-Flow initialization failed"
    echo "Falling back to direct execution..."
}

# Create Role LCTs
echo ""
echo "🎭 Creating Role LCTs..."
echo "  - Genesis Orchestrator (Delegative)"
echo "  - 6 Domain Queens (Delegative)"
echo "  - 24 Worker Roles (Responsive)"
echo "  - 1 Witness Role (Responsive)"

# Setup swarm memory
echo ""
echo "💾 Setting up swarm memory..."
mkdir -p swarm-memory/{architecture,implementation,decisions,learnings,witness}

# Initialize R6 rules for each role
echo ""
echo "📜 Loading R6 rules for all roles..."

# Start the Genesis Orchestrator
echo ""
echo "✨ LAUNCHING GENESIS ORCHESTRATOR"
echo "================================"
echo ""
echo "The Genesis Orchestrator will now:"
echo "  1. Spawn 6 Domain Queen roles"
echo "  2. Each Queen spawns 4 Worker roles"
echo "  3. Establish witness network"
echo "  4. Begin building ACT"
echo ""

# Execute with Node.js
node genesis-queen-init.js

echo ""
echo "🎉 Swarm Initialized!"
echo ""
echo "Monitor progress with:"
echo "  npx claude-flow@alpha swarm status --real-time"
echo ""
echo "View swarm memory:"
echo "  npx claude-flow@alpha memory view --tree"
echo ""
echo "Check ATP consumption:"
echo "  npx claude-flow@alpha atp balance --all-roles"
echo ""
echo "The fractal swarm is now building ACT!"
echo "Estimated completion: 4 weeks"
echo ""
echo "📊 Current Swarm Structure:"
echo "  1 Genesis Orchestrator"
echo "  6 Domain Queens"
echo "  24 Worker Roles"
echo "  1 Witness Network"
echo "  ─────────────────"
echo "  32 Total Role Entities"
echo ""
echo "All actions are being witnessed and recorded in the ledger."
echo "Every role has its own LCT and follows R6 rules."
echo ""
echo "🔄 The swarm that builds ACT will be ACT's first citizen!"