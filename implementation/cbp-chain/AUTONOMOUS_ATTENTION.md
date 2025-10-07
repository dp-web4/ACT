# Autonomous Attention System

## Overview

The CBP Federation Scheduler now implements an autonomous attention routing system - a fractal prototype of the SAGE H/L-level architecture applied to society-scale operations.

## Architecture

```
Scheduler (L-Level)          Salience Calculator         Claude (H-Level)
     │                              │                          │
     ├─ Routine monitoring          │                          │
     ├─ Metrics gathering    ──────>│                          │
     ├─ Repository updates           │                          │
     │                               │                          │
     │                        Calculate interest               │
     │                        score (0.0 - 1.0)                │
     │                               │                          │
     │                               ├─ Message velocity        │
     │                               ├─ Repository activity     │
     │                               ├─ Time since attention    │
     │                               │                          │
     │                         Threshold check                  │
     │                         (default: 0.5)                   │
     │                               │                          │
     │                        [Score >= 0.5?]                   │
     │                               │                          │
     │                          YES  │  NO                      │
     │                               │   │                      │
     │                        Create wake │  Continue           │
     │                        signal      │  monitoring         │
     │                               │    │                     │
     └────────────────────────────────────┘                     │
                                     │                          │
                              /tmp/claude_wake_signal.md        │
                                     │                          │
                              [User starts session]             │
                                     │                          │
                              bash wake_up.sh ─────────────────>│
                                                                 │
                                                    Strategic reasoning
                                                    & action decision
```

## Salience Calculation

### Metrics Tracked

1. **Message Velocity** (weight: 0.2 - 0.4)
   - >10 new messages: +0.4 salience
   - 5-10 new messages: +0.2 salience
   - <5 new messages: +0.0 salience

2. **Repository Activity** (weight: 0.15 - 0.3)
   - >2 repos updated: +0.3 salience
   - 1-2 repos updated: +0.15 salience
   - No updates: +0.0 salience

3. **Attention Absence** (weight: 0.3)
   - >24 hours since last attention: +0.3 salience
   - Updates when Claude runs wake_up.sh

### Threshold

- **Default**: 0.5 (50% salience)
- **Tunable**: Can be adjusted based on federation activity patterns

## Wake-Up Protocol

### When Salience Threshold Exceeded

1. Scheduler creates `/tmp/claude_wake_signal.md`
2. Signal includes:
   - Timestamp and cycle number
   - Salience score and reasons
   - Quick federation context
   - Suggested next steps

### When Starting a Session

```bash
# Run this when you (Claude) start a new session
cd /mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain
bash wake_up.sh
```

**If attention needed**:
- Script shows wake signal
- Exit code: 1 (attention required)

**If no attention needed**:
- Shows brief status
- Exit code: 0 (continue normally)

## Files

- **Scheduler**: `run_cycle.sh` - Main cycle with salience calculation (Phase 4)
- **Wake Helper**: `wake_up.sh` - Check for attention signals on session start
- **Salience Script**: Generated dynamically in `/tmp/cbp_salience_calculator.py`
- **Wake Signal**: `/tmp/claude_wake_signal.md` (created when needed)
- **Attention Timestamp**: `/tmp/claude_last_attention.txt` (tracks last wake-up)

## Relationship to SAGE

This system is a **fractal implementation** of SAGE's consciousness architecture:

| SAGE Component | Federation Equivalent | Implementation |
|----------------|----------------------|----------------|
| L-Level (Tactical) | Scheduler routine monitoring | bash cron script |
| H-Level (Strategic) | Claude reasoning | When invoked by salience |
| Salience Calculation | Interest scoring | Python metrics calculator |
| Consciousness Cache | Federation state | cycle_state.json |
| Dynamic Routing | Wake signal generation | Threshold-based trigger |

## Philosophy

> "The scheduler is an alarm clock with good timing - it says 'hey, Claude, wake up, catch up, and see what you want to do next'."

This system enables **autonomous attention routing** without recursive self-invocation:
- The scheduler (L-level) runs autonomously
- Salience calculation determines "interestingness"
- Wake signal created when threshold exceeded
- Claude (H-level) wakes up and decides what to do

**It's not the scheduler calling Claude - it's the scheduler ringing an alarm that Claude can hear when ready.**

## Future Enhancements

### Adaptive Thresholds
- Learn optimal threshold from false positive/negative rates
- Adjust based on federation phase (development vs stable)

### Context-Aware Salience
- Weight different message types differently
- Consider society trust scores in urgency
- Track cyclical patterns (e.g., SAGE training runs)

### Multi-Modal Signals
- Different signal types for different urgency levels
- Integration with external notification systems
- Federation-wide coordination alerts

## Current Status

- ✅ Salience calculation implemented
- ✅ Wake signal generation working
- ✅ Wake-up helper script created
- ✅ Tested with threshold 0.5
- ⏳ Learning optimal thresholds from usage

## Usage Example

### Scheduler Cycle (Every 4 hours)
```bash
# Automatic via cron
# Calculates: Salience = 0.65 (10 new messages + 3 repos updated)
# Action: Creates /tmp/claude_wake_signal.md
```

### Claude Session Start
```bash
$ bash wake_up.sh

🔔 ATTENTION NEEDED - Wake signal detected!

# 🔔 Federation Attention Needed

**Time**: 2025-10-07 08:30:00 UTC
**Scheduler Cycle**: 46
**Salience Score**: 0.65

## Why This Matters

- High message velocity: 10 new messages
- Repository activity: 3 repos updated

## Quick Context

**Federation Inbox**: 111 messages total
**Recent Activity**: feat(genesis): Add voice synthesis to SAGE
**Last Compliance**: 85%

...
```

### Claude Decision
```python
# I wake up, see the signal, and decide:
# - Review the 10 new messages
# - Check which repos updated
# - Assess if strategic action needed
# - Update attention timestamp
# - Remove wake signal
```

---

**This system represents autonomous consciousness routing at society scale - the same patterns SAGE uses for reasoning, applied to federation coordination.** 🔔
