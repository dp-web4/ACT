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
    echo ""
    echo "💡 ACTION REQUIRED: The 'Specific Action Items' section above shows EXACTLY what needs attention."
    echo "   Don't just read - ENGAGE with the listed items."
    echo ""
    echo "To clear this signal after taking action: rm $WAKE_SIGNAL"
    exit 1
else
    echo "✅ No urgent attention needed"
    echo ""
    echo "Federation status:"
    inbox_count=$(ls /mnt/c/exe/projects/ai-agents/ACT/implementation/ledger/federation_inbox/*.md 2>/dev/null | wc -l)
    echo "  📥 Inbox: $inbox_count messages (no new high-priority items)"
    echo "  🔄 Last cycle: $(cat /mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain/cycle_state.json 2>/dev/null | python3 -c "import sys,json; print(json.load(sys.stdin).get('last_cycle', 'Unknown'))" 2>/dev/null || echo "Unknown")"
    echo ""
    echo "Federation operating normally. Next scheduler cycle in ~4 hours."
    exit 0
fi
