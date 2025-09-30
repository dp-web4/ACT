#!/bin/bash
# Start CBP Society Scheduler Daemon
# Survives sleep/restarts via nohup and background execution

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DAEMON_SCRIPT="$SCRIPT_DIR/cbp_scheduler_daemon.py"
STATE_DIR="$HOME/.cbp_scheduler"
LOG_DIR="$STATE_DIR/logs"
PID_FILE="$STATE_DIR/scheduler.pid"

# Create directories
mkdir -p "$LOG_DIR"

# Check if already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "✅ CBP Scheduler already running (PID: $PID)"
        exit 0
    else
        echo "Removing stale PID file"
        rm "$PID_FILE"
    fi
fi

echo "Starting CBP Society Scheduler..."

# Start with nohup to survive terminal close
nohup python3 "$DAEMON_SCRIPT" > "$LOG_DIR/startup.log" 2>&1 &
DAEMON_PID=$!

# Wait a moment for startup
sleep 2

# Check if it started successfully
if ps -p "$DAEMON_PID" > /dev/null 2>&1; then
    echo "✅ CBP Scheduler started successfully (PID: $DAEMON_PID)"
    echo "   Log: $LOG_DIR/scheduler_$(date +%Y%m%d).log"
    echo "   Status: python3 $DAEMON_SCRIPT status"
    echo "   Stop: python3 $DAEMON_SCRIPT stop"
else
    echo "❌ Failed to start CBP Scheduler"
    echo "   Check logs: $LOG_DIR/startup.log"
    exit 1
fi