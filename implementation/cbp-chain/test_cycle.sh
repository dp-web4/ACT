#!/bin/bash
# Test the CBP Society autonomous cycle immediately

echo "Testing CBP Society Autonomous Cycle..."
echo "This will execute a single cycle immediately for testing."
echo ""

# Run the cycle
bash "$(dirname "$0")/run_cycle.sh"

# Show the results
echo ""
echo "=== Cycle Results ==="

# Show cycle state
if [ -f "$(dirname "$0")/cycle_state.json" ]; then
    echo "Cycle State:"
    cat "$(dirname "$0")/cycle_state.json"
fi

echo ""

# Show todo state
if [ -f "$(dirname "$0")/todo_state.json" ]; then
    echo "Current Todos:"
    python3 -c "
import json
with open('$(dirname "$0")/todo_state.json', 'r') as f:
    todos = json.load(f)
    pending = [t for t in todos if t.get('status') == 'pending']
    completed = [t for t in todos if t.get('status') == 'completed']
    print(f'  Pending: {len(pending)}')
    print(f'  Completed: {len(completed)}')
    if pending:
        print('\nNext pending tasks:')
        for t in pending[:3]:
            print(f'  - {t["content"]}')"
fi

echo ""
echo "Test complete! Check /tmp/cbp_cycle_*.log for detailed logs."