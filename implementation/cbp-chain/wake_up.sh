#!/bin/bash
# Claude Wake-Up Helper
# Use this when starting a new session to check if attention is needed

WAKE_SIGNAL="/tmp/claude_wake_signal.md"

# Update attention timestamp
echo "$(date +%s)" > /tmp/claude_last_attention.txt

# Check if wake signal exists
if [ -f "$WAKE_SIGNAL" ]; then
    echo "🔔 ATTENTION NEEDED - Wake signal detected!"
    echo ""
    cat "$WAKE_SIGNAL"
    echo ""
    echo "---"
    echo "To acknowledge: rm $WAKE_SIGNAL"
    exit 1
else
    echo "✅ No urgent attention needed"
    echo ""
    echo "Current federation status:"
    echo "  Inbox: $(ls /mnt/c/exe/projects/ai-agents/ACT/implementation/ledger/federation_inbox/*.md 2>/dev/null | wc -l) messages"
    echo "  Last scheduler cycle: $(cat /mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain/cycle_state.json 2>/dev/null | python3 -c "import sys,json; print(json.load(sys.stdin).get('last_cycle', 'Unknown'))" 2>/dev/null || echo "Unknown")"
    echo ""
    echo "Federation is operating normally. Continue with current tasks."
    exit 0
fi
