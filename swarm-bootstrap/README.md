# ACT Swarm Monitoring

## Quick Start

Monitor the fractal swarm building ACT:

```bash
# Check current status
./swarm-cli.sh status

# Live monitoring (updates every 5 seconds)
./swarm-cli.sh monitor

# Or use Node directly
node monitor-swarm.js      # Single snapshot
node monitor-swarm.js live  # Real-time monitoring
```

## Available Commands

### CLI Commands
```bash
./swarm-cli.sh status    # Current swarm status
./swarm-cli.sh monitor   # Live dashboard
./swarm-cli.sh queens    # List domain queens
./swarm-cli.sh workers   # Worker activity
./swarm-cli.sh tasks     # Current tasks
./swarm-cli.sh progress  # Implementation progress
./swarm-cli.sh atp       # ATP economy status
./swarm-cli.sh witness   # Witness activity log
./swarm-cli.sh memory    # Browse swarm memory
./swarm-cli.sh logs      # View swarm logs
./swarm-cli.sh evolve    # Trigger evolution cycle
```

### Claude-Flow Commands
```bash
# Real-time swarm status
npx claude-flow@alpha swarm status --real-time

# View swarm memory tree
npx claude-flow@alpha memory view --format=tree

# Check ATP consumption
npx claude-flow@alpha atp balance --all-roles

# View witness records
npx claude-flow@alpha witness logs --tail=20

# Execute specific task
npx claude-flow@alpha swarm execute \
  --queen="LCT-Infrastructure-Queen" \
  --task="Your task here"
```

## Monitor Dashboard

The monitoring dashboard shows:

1. **Swarm Hierarchy**
   - Genesis Orchestrator status
   - 6 Domain Queens with ATP budgets
   - 24 Worker Roles status

2. **Implementation Progress**
   - Phase-by-phase progress bars
   - Overall completion percentage
   - Estimated time to completion

3. **Current Tasks**
   - Active tasks per queen
   - Worker allocation
   - Task status (planning/in_progress/queued)

4. **ATP Economy**
   - Treasury balance
   - ATP spent and allocated
   - ADP generation efficiency
   - Recent transactions

5. **Witness Activity**
   - Last 10 witnessed actions
   - Activity statistics by role
   - Complete audit trail

## Swarm Structure

```
🌟 Genesis Orchestrator (1000 ATP/day)
    ├── 👑 LCT-Infrastructure-Queen (100 ATP/day)
    │   └── 4 Worker Roles
    ├── 👑 ACP-Protocol-Queen (100 ATP/day)
    │   └── 4 Worker Roles
    ├── 👑 Demo-Society-Queen (100 ATP/day)
    │   └── 4 Worker Roles
    ├── 👑 MCP-Bridge-Queen (100 ATP/day)
    │   └── 4 Worker Roles
    ├── 👑 Client-Interface-Queen (100 ATP/day)
    │   └── 4 Worker Roles
    └── 👑 ATP-Economy-Queen (100 ATP/day)
        └── 4 Worker Roles
```

## Memory Structure

Swarm memory is organized in:
- `architecture/` - System design decisions
- `implementation/` - Code and progress tracking
- `decisions/` - Swarm consensus records
- `learnings/` - Patterns and optimizations
- `witness/` - Complete action history
- `economy/` - ATP/ADP ledger

## Web4 Compliance

Every agent in the swarm:
- Is a Role entity with its own LCT
- Has R6 rules defining behavior
- All actions are witnessed
- Consumes ATP and generates ADP
- Follows Web4 standard protocols

## Evolution Cycles

The swarm evolves through:
- **Daily Standups** (9 AM) - Progress sync
- **Midday Reviews** (12 PM) - Integration check
- **Evening Retros** (6 PM) - Learning capture
- **Weekly Evolution** (Fridays) - Self-improvement

## Troubleshooting

If monitoring shows no data:
1. Ensure swarm is initialized: `./launch-swarm.sh`
2. Check Claude-Flow: `npx claude-flow@alpha --version`
3. Initialize memory: `./swarm-cli.sh memory`
4. Check logs: `./swarm-cli.sh logs`