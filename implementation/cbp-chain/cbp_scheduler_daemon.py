#!/usr/bin/env python3
"""
CBP Society Autonomous Scheduler Daemon
Runs persistently, survives restarts, executes 4-hour cycles
"""

import os
import sys
import time
import json
import signal
import subprocess
import logging
from datetime import datetime, timedelta
from pathlib import Path
import threading
import atexit

# Configuration
SCRIPT_DIR = Path(__file__).parent
ACT_ROOT = SCRIPT_DIR.parent.parent
STATE_DIR = Path.home() / '.cbp_scheduler'
PID_FILE = STATE_DIR / 'scheduler.pid'
LOG_DIR = STATE_DIR / 'logs'
STATE_FILE = STATE_DIR / 'state.json'
CYCLE_SCRIPT = SCRIPT_DIR / 'run_cycle.sh'
CYCLE_INTERVAL = 4 * 60 * 60  # 4 hours in seconds

# Ensure directories exist
STATE_DIR.mkdir(exist_ok=True)
LOG_DIR.mkdir(exist_ok=True)

# Setup logging
log_file = LOG_DIR / f"scheduler_{datetime.now().strftime('%Y%m%d')}.log"
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler(log_file),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class CBPScheduler:
    def __init__(self):
        self.running = True
        self.state = self.load_state()
        self.next_cycle = self.calculate_next_cycle()

        # Set up signal handlers
        signal.signal(signal.SIGTERM, self.handle_shutdown)
        signal.signal(signal.SIGINT, self.handle_shutdown)

        # Register cleanup
        atexit.register(self.cleanup)

    def load_state(self):
        """Load scheduler state from disk"""
        if STATE_FILE.exists():
            try:
                with open(STATE_FILE, 'r') as f:
                    return json.load(f)
            except:
                pass

        return {
            'cycles_completed': 0,
            'last_cycle': None,
            'scheduler_started': datetime.now().isoformat(),
            'total_runtime': 0
        }

    def save_state(self):
        """Save scheduler state to disk"""
        with open(STATE_FILE, 'w') as f:
            json.dump(self.state, f, indent=2)

    def calculate_next_cycle(self):
        """Calculate when the next cycle should run"""
        # Run at fixed times: 00:00, 04:00, 08:00, 12:00, 16:00, 20:00 UTC
        now = datetime.utcnow()

        # Find the next 4-hour mark
        hour = now.hour
        next_hour = ((hour // 4) + 1) * 4

        if next_hour >= 24:
            # Next cycle is tomorrow
            next_cycle = now.replace(hour=0, minute=0, second=0, microsecond=0)
            next_cycle += timedelta(days=1)
        else:
            next_cycle = now.replace(hour=next_hour, minute=0, second=0, microsecond=0)

        return next_cycle

    def run_cycle(self):
        """Execute a single work cycle"""
        logger.info("="*50)
        logger.info("Starting CBP Society work cycle")

        try:
            # Change to ACT directory for cycle execution
            os.chdir(ACT_ROOT)

            # Run the cycle script
            result = subprocess.run(
                ['bash', str(CYCLE_SCRIPT)],
                capture_output=True,
                text=True,
                timeout=1800  # 30 minute timeout
            )

            if result.returncode == 0:
                logger.info("Cycle completed successfully")
                self.state['cycles_completed'] += 1
            else:
                logger.error(f"Cycle failed with code {result.returncode}")
                logger.error(f"Error output: {result.stderr[:500]}")

            # Log output
            if result.stdout:
                logger.info(f"Cycle output:\n{result.stdout[:1000]}")

            # Update state
            self.state['last_cycle'] = datetime.now().isoformat()
            self.save_state()

        except subprocess.TimeoutExpired:
            logger.error("Cycle execution timed out after 30 minutes")
        except Exception as e:
            logger.error(f"Error running cycle: {e}")

        # Calculate next cycle time
        self.next_cycle = self.calculate_next_cycle()
        logger.info(f"Next cycle scheduled for {self.next_cycle.isoformat()}")
        logger.info("="*50)

    def handle_shutdown(self, signum, frame):
        """Handle shutdown signals gracefully"""
        logger.info(f"Received signal {signum}, shutting down...")
        self.running = False

    def cleanup(self):
        """Clean up on exit"""
        if PID_FILE.exists():
            PID_FILE.unlink()
        self.save_state()
        logger.info("CBP Scheduler stopped")

    def write_pid(self):
        """Write PID to file for process management"""
        with open(PID_FILE, 'w') as f:
            f.write(str(os.getpid()))

    def check_singleton(self):
        """Ensure only one instance is running"""
        if PID_FILE.exists():
            try:
                with open(PID_FILE, 'r') as f:
                    old_pid = int(f.read())

                # Check if process is still running
                os.kill(old_pid, 0)
                logger.error(f"Scheduler already running with PID {old_pid}")
                sys.exit(1)
            except (ProcessLookupError, ValueError):
                # Old process is dead, remove PID file
                PID_FILE.unlink()

    def run(self):
        """Main scheduler loop"""
        self.check_singleton()
        self.write_pid()

        logger.info("="*60)
        logger.info("CBP SOCIETY AUTONOMOUS SCHEDULER STARTED")
        logger.info(f"PID: {os.getpid()}")
        logger.info(f"State directory: {STATE_DIR}")
        logger.info(f"Next cycle: {self.next_cycle.isoformat()}")
        logger.info("="*60)

        # Run immediately if more than 4 hours since last cycle
        if self.state['last_cycle']:
            last = datetime.fromisoformat(self.state['last_cycle'])
            if datetime.now() - last > timedelta(hours=4):
                logger.info("Running immediate cycle (>4 hours since last)")
                self.run_cycle()

        # Main loop
        while self.running:
            now = datetime.utcnow()

            if now >= self.next_cycle:
                self.run_cycle()

            # Sleep for 30 seconds between checks
            time.sleep(30)

            # Update runtime
            self.state['total_runtime'] += 30

            # Log heartbeat every 30 minutes
            if self.state['total_runtime'] % 1800 == 0:
                logger.info(f"Scheduler heartbeat - Next cycle: {self.next_cycle.isoformat()}")

def status():
    """Check scheduler status"""
    if PID_FILE.exists():
        try:
            with open(PID_FILE, 'r') as f:
                pid = int(f.read())

            # Check if process is running
            os.kill(pid, 0)
            print(f"✅ CBP Scheduler is running (PID: {pid})")

            # Show state
            if STATE_FILE.exists():
                with open(STATE_FILE, 'r') as f:
                    state = json.load(f)
                print(f"   Cycles completed: {state.get('cycles_completed', 0)}")
                print(f"   Last cycle: {state.get('last_cycle', 'Never')}")

            return True

        except (ProcessLookupError, ValueError):
            print("❌ CBP Scheduler PID file exists but process is not running")
            return False
    else:
        print("❌ CBP Scheduler is not running")
        return False

def stop():
    """Stop the scheduler"""
    if PID_FILE.exists():
        try:
            with open(PID_FILE, 'r') as f:
                pid = int(f.read())

            os.kill(pid, signal.SIGTERM)
            print(f"Sent SIGTERM to scheduler (PID: {pid})")

            # Wait for it to stop
            for i in range(10):
                time.sleep(1)
                try:
                    os.kill(pid, 0)
                except ProcessLookupError:
                    print("✅ Scheduler stopped")
                    return

            # Force kill if still running
            os.kill(pid, signal.SIGKILL)
            print("⚠️  Had to force-kill scheduler")

        except Exception as e:
            print(f"Error stopping scheduler: {e}")
    else:
        print("Scheduler is not running")

def main():
    """Main entry point"""
    if len(sys.argv) > 1:
        cmd = sys.argv[1]
        if cmd == 'status':
            status()
        elif cmd == 'stop':
            stop()
        elif cmd == 'restart':
            stop()
            time.sleep(2)
            scheduler = CBPScheduler()
            scheduler.run()
        else:
            print(f"Unknown command: {cmd}")
            print("Usage: cbp_scheduler_daemon.py [start|status|stop|restart]")
    else:
        # Default: start scheduler
        scheduler = CBPScheduler()
        scheduler.run()

if __name__ == "__main__":
    main()