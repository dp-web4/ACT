#!/bin/bash
# Auto-start Genesis Federation Scheduler
# Runs on WSL startup

SCHEDULER_PATH="/mnt/c/projects/ai-agents/ACT/implementation/ledger/genesis_federation_scheduler.py"
LOG_DIR="$HOME/.genesis_scheduler/logs"
PID_FILE="$HOME/.genesis_scheduler/scheduler.pid"

# Create log directory
mkdir -p "$LOG_DIR"

# Check if already running
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE")
    if ps -p "$OLD_PID" > /dev/null 2>&1; then
        echo "Scheduler already running (PID: $OLD_PID)"
        exit 0
    fi
fi

# Start scheduler in background
echo "Starting Genesis Federation Scheduler..."
nohup python3 "$SCHEDULER_PATH" run > "$LOG_DIR/scheduler_$(date +%Y%m%d).log" 2>&1 &
NEW_PID=$!

# Save PID
echo "$NEW_PID" > "$PID_FILE"

echo "Scheduler started (PID: $NEW_PID)"
echo "Log: $LOG_DIR/scheduler_$(date +%Y%m%d).log"