# Genesis Federation Scheduler Auto-Start Setup

## Overview
The Genesis Federation Scheduler coordinates federation-wide activities. This guide sets up automatic startup when the machine boots or resumes from sleep.

## Quick Start

### Manual Start
```bash
cd /mnt/c/projects/ai-agents/ACT/implementation/ledger
bash start_scheduler.sh
```

### Check Status
```bash
python3 genesis_federation_scheduler.py status
```

### Stop Scheduler
```bash
bash stop_scheduler.sh
```

## Auto-Start on WSL Boot

### Option 1: Windows Task Scheduler (Recommended for Laptops)
This handles both boot and resume from sleep:

1. **Create PowerShell startup script** (already created):
   ```powershell
   # Location: C:\projects\ai-agents\ACT\implementation\ledger\start_scheduler.ps1
   wsl -d Ubuntu -u dp bash /mnt/c/projects/ai-agents/ACT/implementation/ledger/start_scheduler.sh
   ```

2. **Add to Task Scheduler**:
   - Open Task Scheduler (taskschd.msc)
   - Create Basic Task → "Genesis Scheduler"
   - Trigger: "When the computer starts"
   - Action: "Start a program"
   - Program: `powershell.exe`
   - Arguments: `-File C:\projects\ai-agents\ACT\implementation\ledger\start_scheduler.ps1`
   - Additional triggers:
     - "On workstation unlock" (after sleep)
     - "At log on"

3. **Configure for sleep/wake**:
   - Task properties → Conditions
   - Uncheck "Start only if computer is on AC power"
   - Check "Wake the computer to run this task"

### Option 2: WSL Boot Command (Simple but boot-only)
Add to `/etc/wsl.conf` in WSL:
```ini
[boot]
command = "sudo -u dp /mnt/c/projects/ai-agents/ACT/implementation/ledger/start_scheduler.sh"
```

Restart WSL:
```powershell
wsl --shutdown
wsl
```

### Option 3: .bashrc Auto-Start (User login)
Add to `~/.bashrc`:
```bash
# Auto-start Genesis scheduler
if [ ! -f "$HOME/.genesis_scheduler/scheduler.pid" ]; then
    /mnt/c/projects/ai-agents/ACT/implementation/ledger/start_scheduler.sh
fi
```

## Logs

Logs are stored in:
```
~/.genesis_scheduler/logs/scheduler_YYYYMMDD.log
```

View current log:
```bash
tail -f ~/.genesis_scheduler/logs/scheduler_$(date +%Y%m%d).log
```

## Scheduler Behavior

### Daily Cycle Schedule
- **06:00-08:00**: Dawn Coherence (Awakening)
- **08:00-12:00**: Morning Coordination
- **12:00-14:00**: Midday Synchronism
- **14:00-18:00**: Afternoon Delegation
- **18:00-20:00**: Evening Reflection
- **20:00-06:00**: Night Rest (minimal activity)

### Automatic Activities
- ATP regeneration every hour
- Coherence checks every 4 hours
- State transitions based on time
- Federation coordination tasks

### Handling Sleep/Wake
When machine sleeps:
- Scheduler process pauses (no CPU usage)
- State preserved in `~/.genesis_scheduler/scheduler_state.json`

When machine wakes:
- Process resumes automatically (if running)
- OR Task Scheduler restarts it (if using Option 1)
- Catches up on missed cycles

## Troubleshooting

### Scheduler not running after boot
```bash
# Check if process exists
ps aux | grep genesis_federation_scheduler

# Check PID file
cat ~/.genesis_scheduler/scheduler.pid

# Check logs
tail -20 ~/.genesis_scheduler/logs/scheduler_$(date +%Y%m%d).log

# Restart manually
bash /mnt/c/projects/ai-agents/ACT/implementation/ledger/start_scheduler.sh
```

### Multiple instances running
```bash
# Stop all
bash /mnt/c/projects/ai-agents/ACT/implementation/ledger/stop_scheduler.sh

# Kill any remaining
pkill -f genesis_federation_scheduler

# Start fresh
bash /mnt/c/projects/ai-agents/ACT/implementation/ledger/start_scheduler.sh
```

### Check what cycle we're in
```bash
python3 /mnt/c/projects/ai-agents/ACT/implementation/ledger/genesis_federation_scheduler.py status
```

## Current Status (as of setup)

**Machine**: Society 4 laptop (mobile node)
**Network**: Work/Home alternating
**Time**: Monday 8:17 PM (Night Rest cycle)
**Auto-start**: Configured ✓

---

*"The scheduler runs continuously, coordinating federation activities across all societies, adapting to the natural rhythms of work and rest."*