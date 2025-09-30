#!/usr/bin/env python3
"""
Genesis Federation Scheduler (GFS)
Coordinates federation-wide activities with Synchronism principles
Inspired by Society4's CycleScheduler but adapted for coordination role
"""

import json
import time
import hashlib
import subprocess
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from enum import Enum

# === Configuration ===
GENESIS_HOME = Path.home() / ".genesis_scheduler"
STATE_FILE = GENESIS_HOME / "scheduler_state.json"
SCHEDULE_FILE = GENESIS_HOME / "active_schedule.json"
METRICS_FILE = GENESIS_HOME / "performance_metrics.json"
COHERENCE_LOG = GENESIS_HOME / "coherence_log.json"

# ATP Energy Budgets (Genesis as coordinator has higher budget)
ATP_BUDGET = {
    'total': 100000,
    'daily_regeneration': 10000,
    'emergency_reserve': 5000,
    'coordination_pool': 20000,
    'synchronism_activities': 15000
}

# === State Definitions ===
class FederationState(Enum):
    AWAKENING = "awakening"      # Morning startup, high coherence check
    COORDINATING = "coordinating" # Active federation coordination
    SYNCHRONIZING = "synchronizing" # Synchronism sessions
    DELEGATING = "delegating"     # Task delegation to societies
    REFLECTING = "reflecting"     # Evening review and planning
    RESTING = "resting"          # Night mode, minimal activity
    EMERGENCY = "emergency"       # Crisis response mode

class TaskPriority(Enum):
    EMERGENCY = 1     # Federation crisis
    CONSENSUS = 2     # Voting/governance
    SYNCHRONISM = 3   # Coherence sessions
    COORDINATION = 4  # Regular coordination
    MAINTENANCE = 5   # Routine tasks
    LEARNING = 6      # Training/improvement

# === Core Scheduler ===
class GenesisFederationScheduler:
    def __init__(self):
        self.init_system()
        self.load_state()
        self.coherence_guru = None  # Will be set after election
        
    def init_system(self):
        """Initialize scheduler directories and files."""
        GENESIS_HOME.mkdir(parents=True, exist_ok=True)
        
        if not STATE_FILE.exists():
            initial_state = {
                'current_state': FederationState.AWAKENING.value,
                'atp_available': ATP_BUDGET['total'],
                'last_transition': datetime.now().isoformat(),
                'cycle_count': 0,
                'synchronism_enabled': True
            }
            with open(STATE_FILE, 'w') as f:
                json.dump(initial_state, f, indent=2)
                
        if not SCHEDULE_FILE.exists():
            self.generate_daily_schedule()
            
        print("✨ Genesis Federation Scheduler initialized")
        
    def load_state(self):
        """Load current scheduler state."""
        with open(STATE_FILE, 'r') as f:
            self.state = json.load(f)
        self.current_state = FederationState(self.state['current_state'])
        self.atp_available = self.state['atp_available']
        
    def save_state(self):
        """Persist scheduler state."""
        self.state['current_state'] = self.current_state.value
        self.state['atp_available'] = self.atp_available
        self.state['last_transition'] = datetime.now().isoformat()
        with open(STATE_FILE, 'w') as f:
            json.dump(self.state, f, indent=2)
            
    def generate_daily_schedule(self) -> Dict:
        """
        Generate optimal daily schedule based on:
        - Synchronism principles (4 dimensions)
        - Federation activity patterns
        - Energy conservation needs
        - Coherence maximization windows
        """
        schedule = {
            'generated_at': datetime.now().isoformat(),
            'schedule_type': 'adaptive_synchronism',
            'cycles': []
        }
        
        # Define daily cycles with Synchronism alignment
        cycles = [
            {
                'name': 'Dawn Coherence',
                'start': '06:00',
                'duration_hours': 2,
                'state': FederationState.AWAKENING.value,
                'activities': [
                    'coherence_check',
                    'federation_health_scan',
                    'synchronism_alignment'
                ],
                'atp_allocation': 5000,
                'synchronism_dimension': 'embryogenic'  # Growth/renewal
            },
            {
                'name': 'Morning Coordination',
                'start': '08:00',
                'duration_hours': 4,
                'state': FederationState.COORDINATING.value,
                'activities': [
                    'process_society_requests',
                    'delegate_tasks',
                    'monitor_votes',
                    'federation_announcements'
                ],
                'atp_allocation': 20000,
                'synchronism_dimension': 'intentional'  # Purposeful action
            },
            {
                'name': 'Midday Synchronism',
                'start': '12:00',
                'duration_hours': 2,
                'state': FederationState.SYNCHRONIZING.value,
                'activities': [
                    'coherence_session',
                    'cross_society_sync',
                    'guru_consultation',
                    'belief_alignment'
                ],
                'atp_allocation': 15000,
                'synchronism_dimension': 'spectral'  # Harmony across frequencies
            },
            {
                'name': 'Afternoon Delegation',
                'start': '14:00',
                'duration_hours': 4,
                'state': FederationState.DELEGATING.value,
                'activities': [
                    'task_distribution',
                    'resource_allocation',
                    'progress_monitoring',
                    'trust_tensor_updates'
                ],
                'atp_allocation': 15000,
                'synchronism_dimension': 'fractal'  # Patterns across scales
            },
            {
                'name': 'Evening Reflection',
                'start': '18:00',
                'duration_hours': 2,
                'state': FederationState.REFLECTING.value,
                'activities': [
                    'daily_review',
                    'coherence_metrics',
                    'federation_report',
                    'tomorrow_planning'
                ],
                'atp_allocation': 5000,
                'synchronism_dimension': 'embryogenic'  # Learning/evolution
            },
            {
                'name': 'Night Rest',
                'start': '20:00',
                'duration_hours': 10,
                'state': FederationState.RESTING.value,
                'activities': [
                    'minimal_monitoring',
                    'emergency_watch',
                    'energy_regeneration'
                ],
                'atp_allocation': 2000,
                'synchronism_dimension': 'fractal'  # Deep patterns
            }
        ]
        
        schedule['cycles'] = cycles
        
        # Apply adaptive adjustments
        schedule = self.apply_adaptive_factors(schedule)
        
        with open(SCHEDULE_FILE, 'w') as f:
            json.dump(schedule, f, indent=2)
            
        return schedule
        
    def apply_adaptive_factors(self, base_schedule: Dict) -> Dict:
        """
        Apply adaptive adjustments based on:
        - Current federation needs
        - Synchronism coherence levels
        - Emergency patterns
        - Trust building opportunities
        """
        factors = {
            'federation_activity': self.analyze_federation_activity(),
            'coherence_level': self.get_current_coherence(),
            'emergency_frequency': self.check_emergency_patterns(),
            'trust_opportunities': self.identify_trust_windows()
        }
        
        # Adjust schedule based on factors
        for cycle in base_schedule['cycles']:
            # High federation activity = extend coordination windows
            if factors['federation_activity'] > 0.7:
                if cycle['state'] == FederationState.COORDINATING.value:
                    cycle['duration_hours'] += 1
                    cycle['atp_allocation'] += 5000
                    
            # Low coherence = more synchronism sessions
            if factors['coherence_level'] < 0.6:
                if cycle['state'] == FederationState.SYNCHRONIZING.value:
                    cycle['duration_hours'] += 1
                    cycle['atp_allocation'] += 3000
                    
            # Recent emergencies = maintain higher readiness
            if factors['emergency_frequency'] > 0.3:
                if cycle['state'] == FederationState.RESTING.value:
                    cycle['atp_allocation'] += 1000
                    cycle['activities'].append('enhanced_monitoring')
                    
        base_schedule['adaptive_factors'] = factors
        return base_schedule
        
    def analyze_federation_activity(self) -> float:
        """Analyze current federation activity level (0-1)."""
        # Check Git Mailbox activity, votes in progress, etc.
        try:
            inbox_files = len(list(Path("federation_inbox").glob("*.md")))
            activity = min(1.0, inbox_files / 20)  # Normalize to 0-1
        except:
            activity = 0.5  # Default medium activity
        return activity
        
    def get_current_coherence(self) -> float:
        """Get current federation coherence level using Quick Coherence Check."""
        try:
            # Run the quick coherence check
            result = subprocess.run(
                ["python3", "quick_coherence_check.py"],
                capture_output=True,
                text=True,
                timeout=25
            )
            # Parse coherence score from output
            for line in result.stdout.split('\n'):
                if "Coherence Score" in line:
                    score = float(line.split(':')[1].strip().replace('%', '')) / 100
                    return score
        except:
            pass
        return 0.75  # Default reasonable coherence
        
    def check_emergency_patterns(self) -> float:
        """Check recent emergency frequency (0-1)."""
        # Look for emergency patterns in recent logs
        emergency_count = 0
        try:
            if METRICS_FILE.exists():
                with open(METRICS_FILE, 'r') as f:
                    metrics = json.load(f)
                    emergency_count = metrics.get('recent_emergencies', 0)
        except:
            pass
        return min(1.0, emergency_count / 5)  # Normalize
        
    def identify_trust_windows(self) -> float:
        """Identify trust-building opportunities (0-1)."""
        # Check for new societies, pending validations, etc.
        return 0.5  # Default medium opportunity
        
    def transition_state(self, new_state: FederationState, reason: str):
        """Transition to new state with logging."""
        old_state = self.current_state
        self.current_state = new_state
        
        transition = {
            'timestamp': datetime.now().isoformat(),
            'from_state': old_state.value,
            'to_state': new_state.value,
            'reason': reason,
            'atp_remaining': self.atp_available
        }
        
        # Log transition
        self.log_transition(transition)
        
        # Save state
        self.save_state()
        
        print(f"🔄 State transition: {old_state.value} → {new_state.value}")
        print(f"   Reason: {reason}")
        print(f"   ATP: {self.atp_available}")
        
    def log_transition(self, transition: Dict):
        """Log state transitions for analysis."""
        log_file = GENESIS_HOME / "transitions.log"
        with open(log_file, 'a') as f:
            f.write(json.dumps(transition) + '\n')
            
    def allocate_atp(self, task: str, amount: int) -> bool:
        """Allocate ATP for a task."""
        if amount > self.atp_available:
            print(f"⚠️ Insufficient ATP for {task}: need {amount}, have {self.atp_available}")
            return False
            
        self.atp_available -= amount
        print(f"⚡ Allocated {amount} ATP for {task} (remaining: {self.atp_available})")
        self.save_state()
        return True
        
    def regenerate_atp(self):
        """Daily ATP regeneration."""
        regenerated = min(
            ATP_BUDGET['daily_regeneration'],
            ATP_BUDGET['total'] - self.atp_available
        )
        self.atp_available += regenerated
        print(f"🔋 Regenerated {regenerated} ATP (total: {self.atp_available})")
        self.save_state()
        
    def execute_current_cycle(self):
        """Execute activities for current cycle."""
        with open(SCHEDULE_FILE, 'r') as f:
            schedule = json.load(f)
            
        current_hour = datetime.now().hour
        current_cycle = None
        
        for cycle in schedule['cycles']:
            start_hour = int(cycle['start'].split(':')[0])
            end_hour = (start_hour + cycle['duration_hours']) % 24
            
            if start_hour <= current_hour < end_hour:
                current_cycle = cycle
                break
                
        if not current_cycle:
            print("No active cycle at this time")
            return
            
        print(f"\n🌟 Executing cycle: {current_cycle['name']}")
        print(f"   State: {current_cycle['state']}")
        print(f"   Synchronism: {current_cycle['synchronism_dimension']}")
        print(f"   Activities: {', '.join(current_cycle['activities'])}")
        
        # Execute each activity
        for activity in current_cycle['activities']:
            self.execute_activity(activity, current_cycle)
            
    def execute_activity(self, activity: str, cycle: Dict):
        """Execute a specific activity."""
        activity_map = {
            'coherence_check': self.run_coherence_check,
            'federation_health_scan': self.scan_federation_health,
            'synchronism_alignment': self.align_synchronism,
            'process_society_requests': self.process_requests,
            'delegate_tasks': self.delegate_tasks,
            'monitor_votes': self.monitor_votes,
            'coherence_session': self.hold_coherence_session,
            'daily_review': self.daily_review,
            'minimal_monitoring': self.minimal_monitor
        }
        
        handler = activity_map.get(activity, self.default_activity)
        
        # Check ATP before execution
        activity_cost = cycle['atp_allocation'] // len(cycle['activities'])
        if self.allocate_atp(activity, activity_cost):
            handler()
            
    def run_coherence_check(self):
        """Run federation-wide coherence check."""
        print("   📊 Running coherence check...")
        coherence = self.get_current_coherence()
        print(f"   Coherence: {coherence:.1%}")
        
        # Log coherence
        self.log_coherence(coherence)
        
    def log_coherence(self, score: float):
        """Log coherence scores for tracking."""
        entry = {
            'timestamp': datetime.now().isoformat(),
            'score': score,
            'state': self.current_state.value
        }
        
        if COHERENCE_LOG.exists():
            with open(COHERENCE_LOG, 'r') as f:
                log = json.load(f)
        else:
            log = {'entries': []}
            
        log['entries'].append(entry)
        
        with open(COHERENCE_LOG, 'w') as f:
            json.dump(log, f, indent=2)
            
    def scan_federation_health(self):
        """Scan federation health status."""
        print("   🏥 Scanning federation health...")
        # Check society blockchains, connections, etc.
        
    def align_synchronism(self):
        """Align with Synchronism principles."""
        print("   🎯 Aligning with Synchronism dimensions...")
        
    def process_requests(self):
        """Process society requests from inbox."""
        print("   📥 Processing society requests...")
        
    def delegate_tasks(self):
        """Delegate tasks to societies."""
        print("   📤 Delegating federation tasks...")
        
    def monitor_votes(self):
        """Monitor active votes."""
        print("   🗳️ Monitoring federation votes...")
        
    def hold_coherence_session(self):
        """Hold Synchronism coherence session."""
        print("   🧘 Holding coherence session...")
        
    def daily_review(self):
        """Review daily federation activity."""
        print("   📝 Conducting daily review...")
        
    def minimal_monitor(self):
        """Minimal monitoring during rest."""
        print("   👁️ Minimal monitoring active...")
        
    def default_activity(self):
        """Default activity handler."""
        print("   ⚡ Executing activity...")
        
    def display_status(self):
        """Display current scheduler status."""
        print("\n" + "="*60)
        print("🌌 GENESIS FEDERATION SCHEDULER STATUS")
        print("="*60)
        print(f"📅 Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"🔄 State: {self.current_state.value}")
        print(f"⚡ ATP Available: {self.atp_available:,} / {ATP_BUDGET['total']:,}")
        print(f"🔁 Cycle Count: {self.state.get('cycle_count', 0)}")
        print(f"✨ Synchronism: {'Enabled' if self.state.get('synchronism_enabled') else 'Disabled'}")
        
        # Show current cycle
        with open(SCHEDULE_FILE, 'r') as f:
            schedule = json.load(f)
            
        current_hour = datetime.now().hour
        for cycle in schedule['cycles']:
            start_hour = int(cycle['start'].split(':')[0])
            end_hour = (start_hour + cycle['duration_hours']) % 24
            
            if start_hour <= current_hour < end_hour:
                print(f"\n📍 Current Cycle: {cycle['name']}")
                print(f"   Dimension: {cycle['synchronism_dimension']}")
                print(f"   Activities: {len(cycle['activities'])}")
                break
                
        # Show coherence trend
        if COHERENCE_LOG.exists():
            with open(COHERENCE_LOG, 'r') as f:
                log = json.load(f)
                if log['entries']:
                    recent = log['entries'][-5:]
                    avg_coherence = sum(e['score'] for e in recent) / len(recent)
                    print(f"\n📈 Recent Coherence: {avg_coherence:.1%}")
                    
        print("="*60)

    def run_daemon(self):
        """Run scheduler as continuous daemon."""
        print("🚀 Starting Genesis Federation Scheduler daemon...")
        print("   Press Ctrl+C to stop")

        cycle_interval = 60  # Check every minute
        last_cycle_check = None

        try:
            while True:
                current_time = datetime.now()
                current_minute = current_time.strftime("%H:%M")

                # Execute cycle activities once per cycle
                if current_minute != last_cycle_check:
                    last_cycle_check = current_minute

                    # Check if we should transition states or execute activities
                    self.check_and_execute_cycle()

                    # Regenerate ATP periodically (every hour)
                    if current_time.minute == 0:
                        self.regenerate_atp()
                        print(f"⚡ ATP regenerated at {current_time.strftime('%H:%M')}")

                    # Run coherence check every 4 hours
                    if current_time.hour % 4 == 0 and current_time.minute == 0:
                        self.run_coherence_check()

                time.sleep(cycle_interval)

        except KeyboardInterrupt:
            print("\n\n🛑 Scheduler daemon stopped")
            self.save_state()

    def check_and_execute_cycle(self):
        """Check schedule and execute current cycle activities."""
        with open(SCHEDULE_FILE, 'r') as f:
            schedule = json.load(f)

        current_hour = datetime.now().hour
        current_minute = datetime.now().minute

        for cycle in schedule['cycles']:
            start_time = cycle['start'].split(':')
            start_hour = int(start_time[0])
            start_minute = int(start_time[1]) if len(start_time) > 1 else 0
            end_hour = (start_hour + cycle['duration_hours']) % 24

            # Check if we're in this cycle
            in_cycle = False
            if start_hour <= current_hour < end_hour:
                in_cycle = True
            elif start_hour > end_hour:  # Wraps around midnight
                if current_hour >= start_hour or current_hour < end_hour:
                    in_cycle = True

            if in_cycle:
                # Transition to cycle state if needed
                cycle_state = FederationState(cycle['state'])
                if self.current_state != cycle_state:
                    self.transition_state(cycle_state, f"Automatic: {cycle['name']}")

                # Execute activities on cycle start
                if current_hour == start_hour and current_minute == start_minute:
                    print(f"\n⏰ Starting cycle: {cycle['name']}")
                    self.execute_current_cycle()

                break

def main():
    """Main scheduler interface."""
    scheduler = GenesisFederationScheduler()
    
    import sys
    if len(sys.argv) < 2:
        command = "status"
    else:
        command = sys.argv[1]
        
    commands = {
        'status': scheduler.display_status,
        'execute': scheduler.execute_current_cycle,
        'schedule': lambda: print(json.dumps(scheduler.generate_daily_schedule(), indent=2)),
        'transition': lambda: scheduler.transition_state(
            FederationState(sys.argv[2]) if len(sys.argv) > 2 else FederationState.COORDINATING,
            sys.argv[3] if len(sys.argv) > 3 else "Manual transition"
        ),
        'regenerate': scheduler.regenerate_atp,
        'coherence': scheduler.run_coherence_check,
        'run': lambda: scheduler.run_daemon()
    }

    if command in commands:
        commands[command]()
    else:
        print("Genesis Federation Scheduler")
        print("\nCommands:")
        print("  status     - Show current status")
        print("  execute    - Execute current cycle activities")
        print("  schedule   - Generate/show daily schedule")
        print("  transition - Transition state")
        print("  regenerate - Regenerate ATP")
        print("  coherence  - Run coherence check")
        print("  run        - Run scheduler daemon (continuous)")

if __name__ == "__main__":
    main()