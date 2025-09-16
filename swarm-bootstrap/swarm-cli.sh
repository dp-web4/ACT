#!/bin/bash

# ACT Swarm CLI
# Quick commands for monitoring and controlling the swarm

SWARM_DIR="$(dirname "$0")"
MEMORY_DIR="$SWARM_DIR/swarm-memory"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
RESET='\033[0m'
BOLD='\033[1m'

# Display help
show_help() {
    echo -e "${BOLD}ACT Swarm CLI${RESET}"
    echo "Monitor and control the fractal swarm building ACT"
    echo ""
    echo "Usage: ./swarm-cli.sh [command]"
    echo ""
    echo "Commands:"
    echo "  status      - Show current swarm status"
    echo "  monitor     - Live monitoring dashboard"
    echo "  queens      - List all queen roles and status"
    echo "  workers     - Show worker activity"
    echo "  tasks       - Display current tasks"
    echo "  progress    - Show implementation progress"
    echo "  atp         - Display ATP economy status"
    echo "  witness     - View witness activity log"
    echo "  memory      - Browse swarm memory"
    echo "  logs        - Tail swarm logs"
    echo "  execute     - Execute a task on a queen"
    echo "  evolve      - Trigger swarm evolution"
    echo "  help        - Show this help message"
}

# Status command
show_status() {
    echo -e "${CYAN}${BOLD}🌟 ACT Development Swarm Status${RESET}"
    echo "================================="
    
    # Check if swarm is running
    if pgrep -f "claude-flow" > /dev/null; then
        echo -e "Status: ${GREEN}● Running${RESET}"
    else
        echo -e "Status: ${RED}● Stopped${RESET}"
    fi
    
    # Show basic stats
    echo -e "\n${BOLD}Swarm Composition:${RESET}"
    echo "  • 1 Genesis Orchestrator"
    echo "  • 6 Domain Queens"
    echo "  • 24 Worker Roles"
    echo "  • 1 Witness Network"
    echo ""
    
    # Run node monitor for snapshot
    node "$SWARM_DIR/monitor-swarm.js"
}

# Monitor command
start_monitor() {
    echo -e "${CYAN}${BOLD}Starting Live Monitor...${RESET}"
    node "$SWARM_DIR/monitor-swarm.js" live
}

# Queens command
show_queens() {
    echo -e "${MAGENTA}${BOLD}👑 Domain Queens${RESET}"
    echo "================"
    
    # Use npx claude-flow command if available
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha swarm queens --format=table 2>/dev/null || {
            # Fallback to config display
            echo "Queens configured in swarm:"
            jq -r '.queens[] | "  • \(.name) (\(.domain)) - \(.workers | length) workers"' \
                "$SWARM_DIR/swarm-config.json"
        }
    else
        echo "Claude-Flow not available. Install with: npm install -g claude-flow@alpha"
    fi
}

# Workers command
show_workers() {
    echo -e "${YELLOW}${BOLD}🔧 Worker Activity${RESET}"
    echo "=================="
    
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha swarm workers --active 2>/dev/null || {
            echo "Worker roles by queen:"
            jq -r '.queens[] | "  \(.name):", (.workers[] | "    • \(.type): \(.role)")' \
                "$SWARM_DIR/swarm-config.json"
        }
    fi
}

# Tasks command
show_tasks() {
    echo -e "${BLUE}${BOLD}🎯 Current Tasks${RESET}"
    echo "================"
    
    if [ -f "$MEMORY_DIR/implementation/tasks.json" ]; then
        jq -r '.current[] | "  [\(.status)] \(.queen): \(.task)"' \
            "$MEMORY_DIR/implementation/tasks.json" 2>/dev/null || \
            echo "No task data available yet"
    else
        echo "Tasks will appear here once swarm begins execution"
    fi
}

# Progress command
show_progress() {
    echo -e "${GREEN}${BOLD}📊 Implementation Progress${RESET}"
    echo "========================="
    
    if [ -f "$MEMORY_DIR/implementation/progress.json" ]; then
        jq '.' "$MEMORY_DIR/implementation/progress.json" 2>/dev/null || \
            echo "No progress data available"
    else
        echo "Progress tracking will begin with first task completion"
    fi
}

# ATP command
show_atp() {
    echo -e "${YELLOW}${BOLD}💰 ATP Economy${RESET}"
    echo "=============="
    
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha atp balance --all-roles 2>/dev/null || {
            echo "Treasury: 10,000 ATP"
            echo "Daily Allocation:"
            echo "  • Genesis: 1,000 ATP"
            echo "  • Queens: 600 ATP (6 × 100)"
            echo "  • Workers: ~400 ATP"
        }
    fi
}

# Witness command
show_witness() {
    echo -e "${MAGENTA}${BOLD}👁️ Witness Activity${RESET}"
    echo "==================="
    
    if [ -f "$MEMORY_DIR/witness/activity.jsonl" ]; then
        echo "Recent witness records:"
        tail -n 10 "$MEMORY_DIR/witness/activity.jsonl" | \
            jq -r '"\(.timestamp | fromdate) | \(.role_lct | split(":") | last) | \(.action)"' 2>/dev/null || \
            tail -n 10 "$MEMORY_DIR/witness/activity.jsonl"
    else
        echo "Witness activity will be recorded here"
    fi
}

# Memory browser
browse_memory() {
    echo -e "${CYAN}${BOLD}💾 Swarm Memory${RESET}"
    echo "==============="
    
    if [ -d "$MEMORY_DIR" ]; then
        echo "Memory structure:"
        tree -L 2 "$MEMORY_DIR" 2>/dev/null || ls -la "$MEMORY_DIR"
        echo ""
        echo "Browse specific area:"
        echo "  • architecture: System design decisions"
        echo "  • implementation: Code and progress"
        echo "  • decisions: Swarm consensus"
        echo "  • learnings: Patterns and optimizations"
        echo "  • witness: Action history"
    else
        echo "Creating swarm memory directories..."
        mkdir -p "$MEMORY_DIR"/{architecture,implementation,decisions,learnings,witness}
        echo "Memory directories created"
    fi
}

# Logs command
show_logs() {
    echo -e "${BOLD}📜 Swarm Logs${RESET}"
    echo "============="
    
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha swarm logs --tail=20 2>/dev/null || {
            echo "No Claude-Flow logs available"
            echo "Check system logs or run with --debug flag"
        }
    fi
}

# Execute task
execute_task() {
    local queen=${1:-"ACT-Genesis"}
    local task=${2:-"Status check"}
    
    echo -e "${GREEN}${BOLD}▶️ Executing Task${RESET}"
    echo "Queen: $queen"
    echo "Task: $task"
    echo ""
    
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha swarm execute \
            --queen="$queen" \
            --task="$task" \
            --priority="normal"
    else
        echo "Claude-Flow not available for task execution"
    fi
}

# Trigger evolution
trigger_evolution() {
    echo -e "${MAGENTA}${BOLD}🔄 Triggering Swarm Evolution${RESET}"
    echo "============================="
    
    if command -v npx &> /dev/null; then
        npx claude-flow@alpha swarm evolve \
            --analyze-performance \
            --optimize \
            --learning-rate=0.1
    else
        echo "Manual evolution trigger:"
        echo "  1. Analyze current performance"
        echo "  2. Identify bottlenecks"
        echo "  3. Apply optimizations"
        echo "  4. Update swarm configuration"
    fi
}

# Initialize memory if needed
init_memory() {
    if [ ! -d "$MEMORY_DIR" ]; then
        mkdir -p "$MEMORY_DIR"/{architecture,implementation,decisions,learnings,witness,economy}
        echo "{\"treasury\": 10000, \"spent\": 0, \"allocated\": 0}" > "$MEMORY_DIR/economy/atp-ledger.json"
        echo "{\"phases\": {}}" > "$MEMORY_DIR/implementation/progress.json"
        echo "Swarm memory initialized"
    fi
}

# Main command handler
case "$1" in
    status)
        show_status
        ;;
    monitor)
        start_monitor
        ;;
    queens)
        show_queens
        ;;
    workers)
        show_workers
        ;;
    tasks)
        show_tasks
        ;;
    progress)
        show_progress
        ;;
    atp)
        show_atp
        ;;
    witness)
        show_witness
        ;;
    memory)
        browse_memory
        ;;
    logs)
        show_logs
        ;;
    execute)
        execute_task "$2" "$3"
        ;;
    evolve)
        trigger_evolution
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        show_status
        echo ""
        echo "Run './swarm-cli.sh help' for available commands"
        ;;
esac