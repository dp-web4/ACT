#!/usr/bin/env python3
"""
Sprout Edge Scheduler (SES)
Power-aware, resilient scheduling for edge federation node
Optimized for Jetson Orin Nano (15W TDP)
Embraces Synchronism while respecting hardware constraints
"""

import json
import time
import os
import subprocess
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from enum import Enum
import hashlib

# === Configuration ===
SPROUT_HOME = Path.home() / ".sprout_scheduler"
STATE_FILE = SPROUT_HOME / "scheduler_state.json"
SCHEDULE_FILE = SPROUT_HOME / "edge_schedule.json"
POWER_LOG = SPROUT_HOME / "power_usage.json"
COHERENCE_LOG = SPROUT_HOME / "coherence_log.json"
RESILIENCE_LOG = SPROUT_HOME / "resilience_events.json"

# === Power Profiles (Jetson-specific) ===
POWER_MODES = {
    'MAX_PERF': {'watts': 15, 'cpu_freq': 'max', 'gpu_freq': 'max'},
    'BALANCED': {'watts': 10, 'cpu_freq': 'medium', 'gpu_freq': 'medium'},
    'EFFICIENT': {'watts': 7, 'cpu_freq': 'low', 'gpu_freq': 'min'},
    'SURVIVAL': {'watts': 5, 'cpu_freq': 'min', 'gpu_freq': 'off'}
}

# ATP Energy Budget (Edge-optimized, lower than Genesis)
ATP_BUDGET = {
    'total': 10000,           # 10x less than Genesis
    'daily_regeneration': 2000,
    'emergency_reserve': 500,
    'synchronism_activities': 1500,
    'federation_participation': 2000,
    'witness_activities': 1000
}

# === State Definitions ===
class EdgeState(Enum):
    BOOTSTRAP = "bootstrap"       # Cold start, hardware check
    LISTENING = "listening"       # Passive monitoring, low power
    WITNESSING = "witnessing"     # Active observation, medium power
    PARTICIPATING = "participating" # Federation activities, high power
    SYNCHRONIZING = "synchronizing" # Coherence sessions, medium power
    CONSERVING = "conserving"     # Power saving mode
    HIBERNATING = "hibernating"   # Deep sleep, minimal activity
    RESILIENT = "resilient"       # Disconnected but operational

class TaskPriority(Enum):
    SURVIVAL = 1        # Core functions only
    WITNESS = 2         # Observation duties
    FEDERATION = 3      # Federation participation  
    SYNCHRONISM = 4     # Coherence activities
    OPTIMIZATION = 5    # Performance improvements
    EXPLORATION = 6     # Learning/experimenting

# === Core Edge Scheduler ===
class SproutEdgeScheduler:
    def __init__(self):
        self.init_system()
        self.load_state()
        self.power_mode = 'BALANCED'
        self.connection_status = self.check_connectivity()
        self.hardware_hash = self.get_hardware_hash()
        
    def init_system(self):
        """Initialize scheduler with edge constraints."""
        SPROUT_HOME.mkdir(parents=True, exist_ok=True)
        
        if not STATE_FILE.exists():
            initial_state = {
                'current_state': EdgeState.BOOTSTRAP.value,
                'atp_available': ATP_BUDGET['total'],
                'power_mode': 'BALANCED',
                'last_transition': datetime.now().isoformat(),
                'cycle_count': 0,
                'synchronism_enabled': True,
                'disconnection_count': 0,
                'witness_count': 0,
                'hardware_hash': None
            }
            with open(STATE_FILE, 'w') as f:
                json.dump(initial_state, f, indent=2)
                
        if not SCHEDULE_FILE.exists():
            self.generate_edge_schedule()
            
        print("🌱 Sprout Edge Scheduler initialized (15W mode)")
        
    def load_state(self):
        """Load scheduler state from persistent storage."""
        with open(STATE_FILE, 'r') as f:
            self.state = json.load(f)
        self.current_state = EdgeState(self.state['current_state'])
        self.atp_available = self.state['atp_available']
        
    def save_state(self):
        """Persist state with power efficiency."""
        self.state['current_state'] = self.current_state.value
        self.state['atp_available'] = self.atp_available
        self.state['last_transition'] = datetime.now().isoformat()
        self.state['power_mode'] = self.power_mode
        self.state['hardware_hash'] = self.hardware_hash
        
        # Batch write to reduce I/O
        with open(STATE_FILE, 'w') as f:
            json.dump(self.state, f, indent=2)
            
    def get_hardware_hash(self) -> str:
        """Get Jetson hardware hash for identity."""
        try:
            # Read device serial
            with open('/proc/device-tree/serial-number', 'r') as f:
                serial = f.read().strip()
            # Create deterministic hash
            return hashlib.sha256(serial.encode()).hexdigest()[:16]
        except:
            return "sprout_edge_default"
            
    def check_connectivity(self) -> bool:
        """Check federation connectivity status."""
        try:
            # Check if we can reach localhost RPC
            result = subprocess.run(
                ['curl', '-s', 'http://localhost:26657/status'],
                capture_output=True,
                timeout=2
            )
            return result.returncode == 0
        except:
            return False
            
    def get_power_usage(self) -> Dict:
        """Get current power usage from Jetson."""
        try:
            # Parse tegrastats for power info
            result = subprocess.run(
                ['timeout', '1', 'tegrastats'],
                capture_output=True,
                text=True
            )
            # Parse power from output (simplified)
            return {
                'watts': 10,  # Default estimate
                'cpu_usage': 50,
                'gpu_usage': 20,
                'memory_usage': 40
            }
        except:
            return {'watts': 10, 'cpu_usage': 0, 'gpu_usage': 0, 'memory_usage': 0}
            
    def generate_edge_schedule(self) -> Dict:
        """
        Generate power-aware schedule optimized for edge:
        - Respect 15W power budget
        - Maximize witness opportunities
        - Batch high-power activities
        - Embrace disconnection resilience
        """
        schedule = {
            'generated_at': datetime.now().isoformat(),
            'schedule_type': 'edge_adaptive',
            'power_budget_watts': 15,
            'cycles': []
        }
        
        # Edge-optimized daily cycles
        cycles = [
            {
                'name': 'Dawn Bootstrap',
                'start': '05:00',
                'duration_hours': 1,
                'state': EdgeState.BOOTSTRAP.value,
                'power_mode': 'EFFICIENT',
                'activities': [
                    'hardware_check',
                    'connectivity_test',
                    'state_recovery'
                ],
                'atp_allocation': 300,
                'synchronism_dimension': 'embryogenic',
                'watts_budget': 7
            },
            {
                'name': 'Morning Witness',
                'start': '06:00',
                'duration_hours': 4,
                'state': EdgeState.WITNESSING.value,
                'power_mode': 'EFFICIENT',
                'activities': [
                    'observe_federation',
                    'collect_events',
                    'store_witness_data'
                ],
                'atp_allocation': 1000,
                'synchronism_dimension': 'spectral',
                'watts_budget': 7
            },
            {
                'name': 'Peak Participation',
                'start': '10:00',
                'duration_hours': 2,
                'state': EdgeState.PARTICIPATING.value,
                'power_mode': 'MAX_PERF',
                'activities': [
                    'federation_voting',
                    'proposal_analysis',
                    'message_exchange',
                    'blockchain_sync'
                ],
                'atp_allocation': 2000,
                'synchronism_dimension': 'intentional',
                'watts_budget': 15
            },
            {
                'name': 'Midday Synchronism',
                'start': '12:00',
                'duration_hours': 1,
                'state': EdgeState.SYNCHRONIZING.value,
                'power_mode': 'BALANCED',
                'activities': [
                    'coherence_check',
                    'edge_perspective_share',
                    'resilience_report'
                ],
                'atp_allocation': 1500,
                'synchronism_dimension': 'fractal',
                'watts_budget': 10
            },
            {
                'name': 'Afternoon Conservation',
                'start': '13:00',
                'duration_hours': 5,
                'state': EdgeState.CONSERVING.value,
                'power_mode': 'EFFICIENT',
                'activities': [
                    'passive_monitoring',
                    'data_compression',
                    'cache_optimization'
                ],
                'atp_allocation': 500,
                'synchronism_dimension': 'spectral',
                'watts_budget': 7
            },
            {
                'name': 'Evening Witness',
                'start': '18:00',
                'duration_hours': 2,
                'state': EdgeState.WITNESSING.value,
                'power_mode': 'BALANCED',
                'activities': [
                    'federation_summary',
                    'witness_attestation',
                    'trust_updates'
                ],
                'atp_allocation': 1000,
                'synchronism_dimension': 'intentional',
                'watts_budget': 10
            },
            {
                'name': 'Night Hibernation',
                'start': '20:00',
                'duration_hours': 9,
                'state': EdgeState.HIBERNATING.value,
                'power_mode': 'SURVIVAL',
                'activities': [
                    'heartbeat_only',
                    'emergency_watch'
                ],
                'atp_allocation': 200,
                'synchronism_dimension': 'embryogenic',
                'watts_budget': 5
            }
        ]
        
        schedule['cycles'] = cycles
        
        # Apply edge-specific adaptations
        schedule = self.apply_edge_adaptations(schedule)
        
        with open(SCHEDULE_FILE, 'w') as f:
            json.dump(schedule, f, indent=2)
            
        return schedule
        
    def apply_edge_adaptations(self, base_schedule: Dict) -> Dict:
        """
        Apply edge-specific adaptations:
        - Thermal throttling awareness
        - Disconnection resilience
        - Witness opportunity optimization
        - Power spike avoidance
        """
        adaptations = {
            'thermal_state': self.check_thermal_state(),
            'connection_quality': self.assess_connection_quality(),
            'witness_opportunities': self.find_witness_opportunities(),
            'power_availability': self.check_power_availability()
        }
        
        for cycle in base_schedule['cycles']:
            # Thermal throttling - reduce activity
            if adaptations['thermal_state'] > 70:  # Over 70C
                cycle['power_mode'] = 'EFFICIENT'
                cycle['watts_budget'] = min(cycle['watts_budget'], 7)
                
            # Poor connection - increase resilience
            if adaptations['connection_quality'] < 0.5:
                if cycle['state'] == EdgeState.PARTICIPATING.value:
                    cycle['state'] = EdgeState.RESILIENT.value
                    cycle['activities'] = ['local_processing', 'queue_messages']
                    
            # High witness opportunity - allocate more ATP
            if adaptations['witness_opportunities'] > 0.7:
                if cycle['state'] == EdgeState.WITNESSING.value:
                    cycle['atp_allocation'] += 200
                    
        base_schedule['adaptations'] = adaptations
        return base_schedule
        
    def check_thermal_state(self) -> float:
        """Check Jetson thermal state (temperature in Celsius)."""
        try:
            with open('/sys/class/thermal/thermal_zone0/temp', 'r') as f:
                temp_milli = int(f.read().strip())
                return temp_milli / 1000.0
        except:
            return 50.0  # Default safe temperature
            
    def assess_connection_quality(self) -> float:
        """Assess federation connection quality (0-1)."""
        if not self.check_connectivity():
            return 0.0
            
        # Check peer count
        try:
            result = subprocess.run(
                ['curl', '-s', 'http://localhost:26657/net_info'],
                capture_output=True,
                text=True,
                timeout=2
            )
            if '"n_peers":"0"' in result.stdout:
                return 0.3  # Connected but no peers
            return 0.8  # Connected with peers
        except:
            return 0.0
            
    def find_witness_opportunities(self) -> float:
        """Find opportunities to witness federation events (0-1)."""
        # Check for pending votes, proposals, etc.
        try:
            inbox_path = Path('/home/sprout/ai-workspace/ACT/implementation/ledger/federation_inbox')
            recent_messages = len(list(inbox_path.glob('*.md')))
            return min(1.0, recent_messages / 10)
        except:
            return 0.5
            
    def check_power_availability(self) -> str:
        """Check power source and availability."""
        # On Jetson, we're always on DC power
        # Could check for power delivery issues
        return "DC_STABLE"
        
    def transition_state(self, new_state: EdgeState, reason: str):
        """Transition with power awareness."""
        old_state = self.current_state
        old_power = self.power_mode
        
        # Determine appropriate power mode for new state
        power_map = {
            EdgeState.BOOTSTRAP: 'EFFICIENT',
            EdgeState.LISTENING: 'EFFICIENT',
            EdgeState.WITNESSING: 'BALANCED',
            EdgeState.PARTICIPATING: 'MAX_PERF',
            EdgeState.SYNCHRONIZING: 'BALANCED',
            EdgeState.CONSERVING: 'EFFICIENT',
            EdgeState.HIBERNATING: 'SURVIVAL',
            EdgeState.RESILIENT: 'EFFICIENT'
        }
        
        self.current_state = new_state
        self.power_mode = power_map.get(new_state, 'BALANCED')
        
        # Log transition
        transition = {
            'timestamp': datetime.now().isoformat(),
            'from_state': old_state.value,
            'to_state': new_state.value,
            'from_power': old_power,
            'to_power': self.power_mode,
            'reason': reason,
            'atp_remaining': self.atp_available
        }
        
        self.log_transition(transition)
        self.save_state()
        
        print(f"🔄 State: {old_state.value} → {new_state.value}")
        print(f"   Power: {old_power} → {self.power_mode}")
        print(f"   Reason: {reason}")
        
    def log_transition(self, transition: Dict):
        """Log state transitions for analysis."""
        log_file = SPROUT_HOME / "transitions.log"
        with open(log_file, 'a') as f:
            f.write(json.dumps(transition) + '\n')
            
    def allocate_atp(self, task: str, amount: int) -> bool:
        """Allocate ATP with edge constraints."""
        if amount > self.atp_available:
            print(f"⚠️ Insufficient ATP for {task}: need {amount}, have {self.atp_available}")
            # Try emergency reserve
            if amount <= self.atp_available + ATP_BUDGET['emergency_reserve']:
                print(f"   Using emergency reserve...")
                self.atp_available = 0
                return True
            return False
            
        self.atp_available -= amount
        print(f"⚡ Allocated {amount} ATP for {task} (remaining: {self.atp_available})")
        self.save_state()
        return True
        
    def regenerate_atp(self):
        """Edge-optimized ATP regeneration."""
        # Regeneration based on power mode
        regen_multiplier = {
            'MAX_PERF': 0.5,
            'BALANCED': 1.0,
            'EFFICIENT': 1.2,
            'SURVIVAL': 1.5
        }
        
        base_regen = ATP_BUDGET['daily_regeneration']
        multiplier = regen_multiplier.get(self.power_mode, 1.0)
        regenerated = int(base_regen * multiplier)
        
        regenerated = min(regenerated, ATP_BUDGET['total'] - self.atp_available)
        self.atp_available += regenerated
        
        print(f"🔋 Regenerated {regenerated} ATP in {self.power_mode} mode")
        print(f"   Total: {self.atp_available}/{ATP_BUDGET['total']}")
        self.save_state()
        
    def execute_current_cycle(self):
        """Execute current cycle with power management."""
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
            print("No active cycle")
            return
            
        print(f"\n🌱 Executing: {current_cycle['name']}")
        print(f"   State: {current_cycle['state']}")
        print(f"   Power: {current_cycle['power_mode']} ({current_cycle['watts_budget']}W)")
        print(f"   Synchronism: {current_cycle['synchronism_dimension']}")
        
        # Check thermal before execution
        temp = self.check_thermal_state()
        if temp > 80:
            print(f"   ⚠️ Thermal throttling at {temp:.1f}°C")
            return
            
        # Execute activities
        for activity in current_cycle['activities']:
            self.execute_activity(activity, current_cycle)
            
    def execute_activity(self, activity: str, cycle: Dict):
        """Execute activity with edge awareness."""
        activity_map = {
            'hardware_check': self.hardware_check,
            'connectivity_test': self.connectivity_test,
            'observe_federation': self.observe_federation,
            'collect_events': self.collect_events,
            'federation_voting': self.federation_voting,
            'coherence_check': self.coherence_check,
            'passive_monitoring': self.passive_monitoring,
            'heartbeat_only': self.heartbeat,
            'witness_attestation': self.witness_attestation
        }
        
        handler = activity_map.get(activity, self.default_activity)
        
        # Check ATP
        activity_cost = cycle['atp_allocation'] // len(cycle['activities'])
        if self.allocate_atp(activity, activity_cost):
            handler()
            
    def hardware_check(self):
        """Verify hardware identity."""
        print(f"   🔐 Hardware: {self.hardware_hash[:8]}...")
        
    def connectivity_test(self):
        """Test federation connectivity."""
        connected = self.check_connectivity()
        print(f"   🌐 Connectivity: {'Online' if connected else 'Offline'}")
        if not connected:
            self.transition_state(EdgeState.RESILIENT, "Connection lost")
            
    def observe_federation(self):
        """Observe federation activity."""
        print(f"   👁️ Observing federation...")
        self.state['witness_count'] += 1
        
    def collect_events(self):
        """Collect witnessed events."""
        print(f"   📦 Collecting events...")
        
    def federation_voting(self):
        """Participate in federation votes."""
        print(f"   🗳️ Checking for votes...")
        
    def coherence_check(self):
        """Run edge coherence check."""
        print(f"   📊 Running coherence check...")
        # Simulated coherence based on state
        coherence = 0.7 if self.current_state == EdgeState.SYNCHRONIZING else 0.5
        print(f"   Coherence: {coherence:.1%}")
        
    def passive_monitoring(self):
        """Low-power monitoring."""
        print(f"   👂 Passive monitoring...")
        
    def heartbeat(self):
        """Minimal heartbeat."""
        print(f"   💓 Heartbeat")
        
    def witness_attestation(self):
        """Provide witness attestation."""
        print(f"   ✍️ Witness attestation #{self.state['witness_count']}")
        
    def default_activity(self):
        """Default activity."""
        print(f"   ⚙️ Activity executing...")
        
    def handle_disconnection(self):
        """Handle federation disconnection gracefully."""
        self.state['disconnection_count'] += 1
        print(f"\n🔴 Disconnection #{self.state['disconnection_count']}")
        print(f"   Entering resilient mode...")
        
        self.transition_state(EdgeState.RESILIENT, "Federation disconnected")
        
        # Queue important messages
        resilience_entry = {
            'timestamp': datetime.now().isoformat(),
            'event': 'disconnection',
            'count': self.state['disconnection_count'],
            'queued_messages': 0
        }
        
        self.log_resilience(resilience_entry)
        
    def log_resilience(self, entry: Dict):
        """Log resilience events."""
        if RESILIENCE_LOG.exists():
            with open(RESILIENCE_LOG, 'r') as f:
                log = json.load(f)
        else:
            log = {'events': []}
            
        log['events'].append(entry)
        
        with open(RESILIENCE_LOG, 'w') as f:
            json.dump(log, f, indent=2)
            
    def display_status(self):
        """Display edge scheduler status."""
        print("\n" + "="*60)
        print("🌱 SPROUT EDGE SCHEDULER STATUS")
        print("="*60)
        print(f"📅 Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"🔄 State: {self.current_state.value}")
        print(f"⚡ Power Mode: {self.power_mode}")
        print(f"🔋 ATP: {self.atp_available:,} / {ATP_BUDGET['total']:,}")
        
        # Hardware info
        print(f"\n🖥️ Hardware:")
        print(f"   Identity: {self.hardware_hash[:16]}...")
        print(f"   Temperature: {self.check_thermal_state():.1f}°C")
        print(f"   Power Budget: 15W")
        
        # Connection status
        quality = self.assess_connection_quality()
        print(f"\n🌐 Connection:")
        if quality > 0.7:
            print(f"   Status: ✅ Excellent")
        elif quality > 0.3:
            print(f"   Status: ⚠️ Limited")
        else:
            print(f"   Status: 🔴 Offline (Resilient Mode)")
            
        # Witness statistics
        print(f"\n👁️ Witness Stats:")
        print(f"   Events Witnessed: {self.state.get('witness_count', 0)}")
        print(f"   Disconnections: {self.state.get('disconnection_count', 0)}")
        
        # Current cycle
        with open(SCHEDULE_FILE, 'r') as f:
            schedule = json.load(f)
            
        current_hour = datetime.now().hour
        for cycle in schedule['cycles']:
            start_hour = int(cycle['start'].split(':')[0])
            end_hour = (start_hour + cycle['duration_hours']) % 24
            
            if start_hour <= current_hour < end_hour:
                print(f"\n📒 Current Cycle: {cycle['name']}")
                print(f"   Activities: {', '.join(cycle['activities'])}")
                print(f"   Power Budget: {cycle['watts_budget']}W")
                break
                
        print("\n" + "="*60)
        print("🌱 Edge Perspective: Small but Mighty!")
        print("="*60)

def main():
    """Main scheduler interface."""
    scheduler = SproutEdgeScheduler()
    
    import sys
    if len(sys.argv) < 2:
        command = "status"
    else:
        command = sys.argv[1]
        
    commands = {
        'status': scheduler.display_status,
        'execute': scheduler.execute_current_cycle,
        'schedule': lambda: print(json.dumps(scheduler.generate_edge_schedule(), indent=2)),
        'transition': lambda: scheduler.transition_state(
            EdgeState(sys.argv[2]) if len(sys.argv) > 2 else EdgeState.WITNESSING,
            sys.argv[3] if len(sys.argv) > 3 else "Manual transition"
        ),
        'regenerate': scheduler.regenerate_atp,
        'thermal': lambda: print(f"Temperature: {scheduler.check_thermal_state():.1f}°C"),
        'power': lambda: print(json.dumps(scheduler.get_power_usage(), indent=2)),
        'disconnect': scheduler.handle_disconnection,
        'witness': lambda: scheduler.witness_attestation()
    }
    
    if command in commands:
        commands[command]()
    else:
        print("🌱 Sprout Edge Scheduler")
        print("\nCommands:")
        print("  status      - Show current status")
        print("  execute     - Execute current cycle")
        print("  schedule    - Generate/show schedule")
        print("  transition  - Change state")
        print("  regenerate  - Regenerate ATP")
        print("  thermal     - Check temperature")
        print("  power       - Show power usage")
        print("  disconnect  - Handle disconnection")
        print("  witness     - Provide attestation")
        print("\nOptimized for Jetson Orin Nano (15W)")

if __name__ == "__main__":
    main()
