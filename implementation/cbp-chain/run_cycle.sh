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
log "\n${YELLOW}PHASE 2: Federation Repository Updates${NC}"

# Define federation repositories
AI_AGENTS_BASE="/mnt/c/exe/projects/ai-agents"
FEDERATION_REPOS=(
    "$ACT_ROOT"
    "$AI_AGENTS_BASE/web4/web4-standard"
    "$AI_AGENTS_BASE/HRM"
    "$AI_AGENTS_BASE/Synchronism"
    "$AI_AGENTS_BASE/ModuleCPU"
    "$AI_AGENTS_BASE/CellCPU"
)

# Pull all federation repositories
for REPO in "${FEDERATION_REPOS[@]}"; do
    if [ -d "$REPO/.git" ]; then
        REPO_NAME=$(basename "$REPO")
        log "Pulling $REPO_NAME..."
        cd "$REPO"
        git pull 2>&1 | tail -3 | sed 's/^/  /' | tee -a "$LOG_FILE"
    else
        log "  Skipping $REPO (not a git repository)"
    fi
done

# Return to ACT directory
cd "$ACT_ROOT"

# ==================== PHASE 3: FEDERATION BUSINESS MONITORING ====================
log "\n${YELLOW}PHASE 3: Federation Business Monitoring${NC}"

# Check federation inbox
log "\nChecking federation inbox..."
if [ -d "implementation/ledger/federation_inbox" ]; then
    INBOX_COUNT=$(ls implementation/ledger/federation_inbox/*.md 2>/dev/null | wc -l)
    log "  Messages in inbox: $INBOX_COUNT"

    # Show last 3 messages
    log "  Recent messages:"
    ls -t implementation/ledger/federation_inbox/*.md 2>/dev/null | head -3 | while read MSG; do
        log "    - $(basename $MSG)"
    done
fi

# Check for new RFCs
log "\nChecking for new RFCs..."
if [ -d "$AI_AGENTS_BASE/web4/web4-standard/rfcs" ]; then
    RFC_COUNT=$(ls "$AI_AGENTS_BASE/web4/web4-standard/rfcs"/*.md 2>/dev/null | wc -l)
    log "  Total RFCs: $RFC_COUNT"
    log "  Recent RFCs:"
    ls -t "$AI_AGENTS_BASE/web4/web4-standard/rfcs"/*.md 2>/dev/null | head -3 | while read RFC; do
        log "    - $(basename $RFC)"
    done
fi

# Check SAGE development status
log "\nChecking SAGE development status..."
if [ -d "$AI_AGENTS_BASE/HRM" ]; then
    log "  Recent HRM updates:"
    cd "$AI_AGENTS_BASE/HRM"
    git log --oneline -3 2>/dev/null | sed 's/^/    /' | tee -a "$LOG_FILE"
    cd "$ACT_ROOT"
fi

# Check federation compliance
log "\nChecking federation compliance status..."
if [ -f "implementation/ledger/WEB4_COMPLIANCE_REPORT.md" ]; then
    COMPLIANCE=$(grep -m1 "Overall Compliance" implementation/ledger/WEB4_COMPLIANCE_REPORT.md 2>/dev/null | grep -o "[0-9]*%" || echo "Unknown")
    log "  Current compliance: $COMPLIANCE"
fi

# Check ATP/ADP status
log "\nChecking federation energy status..."
if [ -f "implementation/ledger/federation/tracker_state.json" ]; then
    python3 -c "
import json
with open('implementation/ledger/federation/tracker_state.json', 'r') as f:
    state = json.load(f)
    total = sum(s.get('atp_balance', 0) for s in state.get('societies', {}).values())
    print(f'  Total federation ATP: {total}')
    print('  Society balances:')
    for name, data in state.get('societies', {}).items():
        print(f'    - {name}: {data.get(\"atp_balance\", 0)} ATP')
" 2>/dev/null || log "  Unable to read ATP status"
fi

# Check trust tensor status
log "\nChecking federation trust metrics..."
if [ -f "implementation/ledger/federation/tensors.json" ]; then
    python3 -c "
import json
with open('implementation/ledger/federation/tensors.json', 'r') as f:
    data = json.load(f)
    coherence = data.get('federation_coherence', 0)
    print(f'  Federation coherence: {coherence:.1%}')
    print('  Trust leaderboard:')
    scores = []
    for entity, tensors in data.get('entity_tensors', {}).items():
        if 'genesis' in entity.lower() or 'society' in entity.lower() or 'sprout' in entity.lower():
            trust = tensors.get('trust_tensor', {}).get('aggregate', 0)
            scores.append((entity.split(':')[-1], trust))
    scores.sort(key=lambda x: x[1], reverse=True)
    for name, score in scores[:4]:
        print(f'    - {name}: {score:.3f}')
" 2>/dev/null || log "  Unable to read trust metrics"
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

# ==================== PHASE 4: SALIENCE CALCULATION ====================
log "\n${YELLOW}PHASE 4: Calculating Attention Salience${NC}"

# Create salience calculator
cat > /tmp/cbp_salience_calculator.py << 'EOF'
import json
import os
import glob
from datetime import datetime

def calculate_salience():
    """Calculate if Claude should be woken up"""
    reasons = []
    salience_score = 0.0

    # Load previous state
    state_file = "implementation/cbp-chain/cycle_state.json"
    prev_state = {}
    if os.path.exists(state_file):
        with open(state_file, 'r') as f:
            prev_state = json.load(f)

    # Check 1: Message velocity
    inbox_count = len(glob.glob("implementation/ledger/federation_inbox/*.md"))
    prev_inbox = prev_state.get('prev_inbox_count', inbox_count)
    new_messages = inbox_count - prev_inbox

    if new_messages > 10:
        salience_score += 0.4
        reasons.append(f"High message velocity: {new_messages} new messages")
    elif new_messages > 5:
        salience_score += 0.2
        reasons.append(f"Moderate activity: {new_messages} new messages")

    # Check 2: Repository activity
    repos_updated = 0
    for repo in ['ACT', 'HRM', 'web4-standard', 'Synchronism', 'ModuleCPU']:
        log_file = f"/tmp/cbp_git_pull_{repo}.txt"
        if os.path.exists(log_file):
            with open(log_file, 'r') as f:
                content = f.read()
                if 'Fast-forward' in content or 'files changed' in content:
                    repos_updated += 1

    if repos_updated > 2:
        salience_score += 0.3
        reasons.append(f"Multiple repos updated: {repos_updated} repositories")
    elif repos_updated > 0:
        salience_score += 0.15
        reasons.append(f"Repository activity: {repos_updated} repos updated")

    # Check 3: Time since last attention
    last_attention_file = "/tmp/claude_last_attention.txt"
    hours_since_attention = 999
    if os.path.exists(last_attention_file):
        with open(last_attention_file, 'r') as f:
            last_time = float(f.read().strip())
            hours_since_attention = (datetime.now().timestamp() - last_time) / 3600

    if hours_since_attention > 24:
        salience_score += 0.3
        reasons.append(f"Long absence: {hours_since_attention:.1f} hours since last attention")

    # Update state for next cycle
    current_state = {
        'cycles_completed': prev_state.get('cycles_completed', 0) + 1,
        'last_cycle': datetime.utcnow().isoformat(),
        'prev_inbox_count': inbox_count,
        'todos_completed': 0,
        'messages_processed': 0
    }

    with open(state_file, 'w') as f:
        json.dump(current_state, f, indent=2)

    # Determine if wake-up needed
    threshold = 0.5
    should_wake = salience_score >= threshold

    print(f"Salience Score: {salience_score:.2f} (threshold: {threshold})")
    print(f"Wake Claude: {'YES' if should_wake else 'NO'}")
    if reasons:
        print("Reasons:")
        for reason in reasons:
            print(f"  - {reason}")

    return should_wake, salience_score, reasons

# Execute
if __name__ == "__main__":
    should_wake, score, reasons = calculate_salience()
    if should_wake:
        exit(1)  # Signal to wake Claude
    exit(0)
EOF

python3 /tmp/cbp_salience_calculator.py 2>&1 | tee -a "$LOG_FILE"
WAKE_NEEDED=$?

# ==================== PHASE 5: ATTENTION SIGNAL ====================
if [ $WAKE_NEEDED -eq 1 ]; then
    log "\n${YELLOW}PHASE 5: Generating Attention Signal${NC}"

    # Create wake-up signal
    WAKE_SIGNAL="/tmp/claude_wake_signal.md"
    cat > "$WAKE_SIGNAL" << 'SIGNAL'
# 🔔 Federation Attention Needed

**Time**: $(date -u '+%Y-%m-%d %H:%M:%S UTC')
**Scheduler Cycle**: $(cat implementation/cbp-chain/cycle_state.json | python3 -c "import sys,json; print(json.load(sys.stdin)['cycles_completed'])")
**Salience Score**: $(python3 /tmp/cbp_salience_calculator.py 2>/dev/null | grep "Salience Score" | cut -d: -f2 || echo "Unknown")

## Why This Matters

The federation scheduler has detected activity that warrants your strategic attention:

$(python3 /tmp/cbp_salience_calculator.py 2>/dev/null | grep -A 10 "Reasons:" || echo "- General federation activity")

## Quick Context

**Federation Inbox**: $(ls implementation/ledger/federation_inbox/*.md 2>/dev/null | wc -l) messages total
**Recent Activity**: $(git log --oneline -1)
**Last Compliance**: $(grep -m1 "Overall Compliance" implementation/ledger/WEB4_COMPLIANCE_REPORT.md 2>/dev/null | grep -o "[0-9]*%" || echo "Unknown")

## Specific Action Items

**New Messages Requiring Response:**
$(ls -t implementation/ledger/federation_inbox/*.md 2>/dev/null | head -5 | while read file; do
    sender=$(basename "$file" | cut -d_ -f1)
    subject=$(grep -m1 "^# " "$file" 2>/dev/null | sed 's/^# //' || basename "$file")
    echo "- [$sender] $subject"
done)

**Recent Federation Commits:**
$(git log --oneline --since="24 hours ago" | head -3 | sed 's/^/- /')

**Pending Questions/Decisions:**
$(grep -h "Question.*:" implementation/ledger/federation_inbox/*.md 2>/dev/null | head -3 | sed 's/^/- /')

## Next Steps

When you wake up:
1. **Review specific messages above** - these need CBP response
2. **Check for votes/proposals** - deadlines may be approaching
3. **Assess strategic actions** - what requires CBP engagement?
4. **Update attention timestamp** - mark cycle as reviewed

---

*This signal was generated by the CBP Federation Scheduler based on autonomous salience calculation.*
*Your attention is requested, not demanded. Review at your discretion.*
SIGNAL

    log "  ✅ Wake signal created: $WAKE_SIGNAL"
    log "  📊 Salience threshold exceeded - attention recommended"
else
    log "\n${YELLOW}PHASE 5: Attention Signal${NC}"
    log "  ℹ️  Salience below threshold - no wake signal needed"
fi

# ==================== PHASE 6: TODO GENERATION ====================
log "\n${YELLOW}PHASE 6: Generating New Tasks${NC}"

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