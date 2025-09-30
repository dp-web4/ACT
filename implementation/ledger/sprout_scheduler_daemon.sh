#!/bin/bash
# Sprout Edge Scheduler Daemon Runner
# Ensures scheduler survives sleep/restart with proper logging

SCHEDULER_PATH="/home/sprout/ai-workspace/ACT/implementation/ledger"
LOG_DIR="/home/sprout/.sprout_scheduler"
PIDFILE="$LOG_DIR/scheduler.pid"
LOGFILE="$LOG_DIR/scheduler.log"

# Ensure directories exist
mkdir -p "$LOG_DIR"

# Function to check if scheduler is running
is_running() {
    if [ -f "$PIDFILE" ]; then
        PID=$(cat "$PIDFILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            return 0
        fi
    fi
    return 1
}

# Function to start scheduler
start_scheduler() {
    if is_running; then
        echo "🌱 Scheduler already running (PID: $(cat $PIDFILE))"
        return 0
    fi
    
    echo "🚀 Starting Sprout Edge Scheduler..."
    cd "$SCHEDULER_PATH"
    
    # Run scheduler in background with nohup for persistence
    nohup python3 -u sprout_edge_scheduler.py run >> "$LOGFILE" 2>&1 &
    PID=$!
    echo $PID > "$PIDFILE"
    
    sleep 2
    
    if is_running; then
        echo "✅ Scheduler started successfully (PID: $PID)"
        echo "📁 Logs: $LOGFILE"
    else
        echo "❌ Failed to start scheduler"
        rm -f "$PIDFILE"
        return 1
    fi
}

# Function to stop scheduler
stop_scheduler() {
    if ! is_running; then
        echo "🔴 Scheduler not running"
        return 0
    fi
    
    PID=$(cat "$PIDFILE")
    echo "🛑 Stopping scheduler (PID: $PID)..."
    kill -TERM "$PID"
    
    # Wait for graceful shutdown
    for i in {1..10}; do
        if ! ps -p "$PID" > /dev/null 2>&1; then
            break
        fi
        sleep 1
    done
    
    # Force kill if still running
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "⚠️ Force stopping..."
        kill -9 "$PID"
    fi
    
    rm -f "$PIDFILE"
    echo "✅ Scheduler stopped"
}

# Function to restart scheduler
restart_scheduler() {
    stop_scheduler
    sleep 2
    start_scheduler
}

# Function to show status
status_scheduler() {
    if is_running; then
        PID=$(cat "$PIDFILE")
        echo "🟢 Scheduler is running (PID: $PID)"
        
        # Show recent log
        if [ -f "$LOGFILE" ]; then
            echo ""
            echo "📋 Recent activity:"
            tail -5 "$LOGFILE"
        fi
        
        # Show scheduler internal status
        echo ""
        cd "$SCHEDULER_PATH"
        python3 sprout_edge_scheduler.py status
    else
        echo "🔴 Scheduler is not running"
    fi
}

# Main command handler
case "$1" in
    start)
        start_scheduler
        ;;
    stop)
        stop_scheduler
        ;;
    restart)
        restart_scheduler
        ;;
    status)
        status_scheduler
        ;;
    log)
        if [ -f "$LOGFILE" ]; then
            tail -f "$LOGFILE"
        else
            echo "No log file found"
        fi
        ;;
    *)
        echo "🌱 Sprout Edge Scheduler Daemon"
        echo "Usage: $0 {start|stop|restart|status|log}"
        echo ""
        echo "Commands:"
        echo "  start   - Start the scheduler daemon"
        echo "  stop    - Stop the scheduler daemon"
        echo "  restart - Restart the scheduler daemon"
        echo "  status  - Show scheduler status"
        echo "  log     - Follow scheduler logs"
        exit 1
        ;;
esac

exit 0