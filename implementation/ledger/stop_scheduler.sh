#!/bin/bash
# Stop Genesis Federation Scheduler

PID_FILE="$HOME/.genesis_scheduler/scheduler.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "Scheduler not running (no PID file)"
    exit 0
fi

PID=$(cat "$PID_FILE")

if ps -p "$PID" > /dev/null 2>&1; then
    echo "Stopping scheduler (PID: $PID)..."
    kill "$PID"
    sleep 2

    # Force kill if still running
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "Force stopping..."
        kill -9 "$PID"
    fi

    rm "$PID_FILE"
    echo "Scheduler stopped"
else
    echo "Scheduler not running (stale PID file)"
    rm "$PID_FILE"
fi