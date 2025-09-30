#!/bin/bash
# CBP Society 4-Hour Autonomous Work Cycle
# Executes todo review, federation monitoring, and task generation

set -e  # Exit on error

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ACT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOG_FILE="/tmp/cbp_cycle_$(date +%Y%m%d).log"
STATE_FILE="$SCRIPT_DIR/cycle_state.json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "$1" | tee -a "$LOG_FILE"
}

# Header
log "${GREEN}=========================================${NC}"
log "${GREEN}CBP SOCIETY 4-HOUR AUTONOMOUS CYCLE${NC}"
log "${GREEN}Started: $(date -u '+%Y-%m-%d %H:%M:%S UTC')${NC}"
log "${GREEN}=========================================${NC}"

# Change to ACT directory
cd "$ACT_ROOT"

# ==================== PHASE 1: TODO EXECUTION ====================
log "\n${YELLOW}PHASE 1: Todo Review & Execution${NC}"

# Create Python script for todo handling
cat > /tmp/cbp_todo_handler.py << 'EOF'
import json
import sys
import os
from datetime import datetime

def load_todos():
    """Load current todo list from state"""
    todo_file = "implementation/cbp-chain/todo_state.json"
    if os.path.exists(todo_file):
        with open(todo_file, 'r') as f:
            return json.load(f)
    return []

def save_todos(todos):
    """Save todo list to state"""
    todo_file = "implementation/cbp-chain/todo_state.json"
    with open(todo_file, 'w') as f:
        json.dump(todos, f, indent=2)

def execute_todo():
    """Execute highest priority pending todo"""
    todos = load_todos()

    # Find first pending task
    for i, todo in enumerate(todos):
        if todo.get('status') == 'pending':
            print(f"Executing: {todo['content']}")

            # Mark as in progress
            todos[i]['status'] = 'in_progress'
            save_todos(todos)

            # Simulate execution (in real implementation, would call actual functions)
            # This is where you'd integrate with actual task execution
            print(f"  -> Working on: {todo['content']}")

            # Mark as completed
            todos[i]['status'] = 'completed'
            todos[i]['completed_at'] = datetime.utcnow().isoformat()
            save_todos(todos)

            print(f"  -> Completed: {todo['content']}")
            return True

    print("No pending todos found")
    return False

# Execute
if __name__ == "__main__":
    execute_todo()
EOF

python3 /tmp/cbp_todo_handler.py 2>&1 | tee -a "$LOG_FILE"

# ==================== PHASE 2: FEDERATION MONITORING ====================
log "\n${YELLOW}PHASE 2: Federation Activity Check${NC}"

# Pull latest changes
log "Pulling latest from repository..."
git pull 2>&1 | tail -5 | tee -a "$LOG_FILE"

# Check federation inbox
log "\nChecking federation inbox..."
if [ -d "implementation/ledger/federation_inbox" ]; then
    INBOX_COUNT=$(ls implementation/ledger/federation_inbox/*.md 2>/dev/null | wc -l)
    log "  Messages in inbox: $INBOX_COUNT"

    # Show latest message
    LATEST=$(ls -t implementation/ledger/federation_inbox/*.md 2>/dev/null | head -1)
    if [ -f "$LATEST" ]; then
        log "  Latest message: $(basename $LATEST)"
        head -10 "$LATEST" | sed 's/^/    /' | tee -a "$LOG_FILE"
    fi
fi

# Check recent society activity
log "\nRecent society activity:"
log "  Genesis commits:"
git log --oneline --grep="genesis" -3 2>/dev/null | sed 's/^/    /' | tee -a "$LOG_FILE"

log "  Society4 commits:"
git log --oneline --grep="society4" -3 2>/dev/null | sed 's/^/    /' | tee -a "$LOG_FILE"

log "  Sprout commits:"
git log --oneline --grep="sprout" -3 2>/dev/null | sed 's/^/    /' | tee -a "$LOG_FILE"

# Check for active votes
log "\nChecking for active votes..."
if [ -f "implementation/ledger/voting_tracker.py" ]; then
    python3 implementation/ledger/voting_tracker.py 2>/dev/null | head -20 | tee -a "$LOG_FILE" || log "  No voting tracker available"
fi

# ==================== PHASE 3: TODO GENERATION ====================
log "\n${YELLOW}PHASE 3: Generating New Tasks${NC}"

# Create Python script for todo generation
cat > /tmp/cbp_todo_generator.py << 'EOF'
import json
import os
import glob
from datetime import datetime, timedelta

def load_todos():
    """Load current todo list"""
    todo_file = "implementation/cbp-chain/todo_state.json"
    if os.path.exists(todo_file):
        with open(todo_file, 'r') as f:
            return json.load(f)
    return []

def save_todos(todos):
    """Save todo list"""
    todo_file = "implementation/cbp-chain/todo_state.json"
    with open(todo_file, 'w') as f:
        json.dump(todos, f, indent=2)

def generate_todos():
    """Generate new todos based on federation state"""
    todos = load_todos()
    new_count = 0

    # Check for unread messages in inbox
    inbox_files = glob.glob("implementation/ledger/federation_inbox/*.md")
    for msg_file in inbox_files[-3:]:  # Check last 3 messages
        msg_name = os.path.basename(msg_file)
        # Check if we already have a todo for this message
        if not any(msg_name in t.get('content', '') for t in todos):
            todos.append({
                'content': f"Review and respond to message: {msg_name}",
                'status': 'pending',
                'priority': 'medium',
                'created_at': datetime.utcnow().isoformat()
            })
            new_count += 1
            print(f"Added todo: Review message {msg_name}")

    # Add periodic maintenance tasks
    now = datetime.utcnow()

    # Weekly security check (every Sunday)
    if now.weekday() == 6:  # Sunday
        if not any('security audit' in t.get('content', '').lower() for t in todos if t.get('status') == 'pending'):
            todos.append({
                'content': 'Perform weekly security audit of society implementation',
                'status': 'pending',
                'priority': 'high',
                'created_at': now.isoformat()
            })
            new_count += 1
            print("Added todo: Weekly security audit")

    # Daily federation sync
    if not any('federation sync' in t.get('content', '').lower() for t in todos if t.get('status') == 'pending'):
        todos.append({
            'content': 'Daily federation sync and status check',
            'status': 'pending',
            'priority': 'low',
            'created_at': now.isoformat()
        })
        new_count += 1
        print("Added todo: Daily federation sync")

    save_todos(todos)
    print(f"Generated {new_count} new todos")
    print(f"Total pending: {sum(1 for t in todos if t.get('status') == 'pending')}")
    return new_count

# Execute
if __name__ == "__main__":
    generate_todos()
EOF

python3 /tmp/cbp_todo_generator.py 2>&1 | tee -a "$LOG_FILE"

# ==================== CYCLE METRICS ====================
log "\n${YELLOW}Updating Cycle Metrics${NC}"

# Update cycle state
python3 << EOF
import json
from datetime import datetime

state_file = "$STATE_FILE"

try:
    with open(state_file, 'r') as f:
        state = json.load(f)
except:
    state = {
        'cycles_completed': 0,
        'last_cycle': None,
        'todos_completed': 0,
        'messages_processed': 0
    }

state['cycles_completed'] += 1
state['last_cycle'] = datetime.utcnow().isoformat()

with open(state_file, 'w') as f:
    json.dump(state, f, indent=2)

print(f"Cycles completed: {state['cycles_completed']}")
print(f"Last cycle: {state['last_cycle']}")
EOF

# ==================== FEDERATION HEARTBEAT ====================
log "\n${YELLOW}Sending Federation Heartbeat${NC}"

HEARTBEAT_DIR="implementation/ledger/federation_outbox"
mkdir -p "$HEARTBEAT_DIR"

HEARTBEAT_FILE="$HEARTBEAT_DIR/cbp_heartbeat_$(date +%Y%m%d).txt"
echo "CBP-$(date +%s): Autonomous cycle completed at $(date -u '+%Y-%m-%d %H:%M:%S UTC')" >> "$HEARTBEAT_FILE"
log "  Heartbeat sent to federation"

# ==================== COMPLETION ====================
log "\n${GREEN}=========================================${NC}"
log "${GREEN}CYCLE COMPLETE: $(date -u '+%Y-%m-%d %H:%M:%S UTC')${NC}"
log "${GREEN}=========================================${NC}"

# Clean up temp files
rm -f /tmp/cbp_todo_handler.py /tmp/cbp_todo_generator.py

# Return success
exit 0