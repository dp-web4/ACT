#!/usr/bin/env python3
"""
Sprout Edge LCT (Linked Context Token) Implementation
Web4-compliant identity for edge nodes with power/thermal awareness
Optimized for Jetson Orin Nano (15W TDP)
"""

import json
import hashlib
import time
import subprocess
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass, field, asdict
from pathlib import Path
from enum import Enum

# ============================================================================
# Edge-Optimized Web4 Data Structures
# ============================================================================

class PowerState(Enum):
    """Edge power states for context-aware operations"""
    MAX_PERF = "max_perf"  # 15W
    BALANCED = "balanced"  # 10W
    EFFICIENT = "efficient"  # 7W
    SURVIVAL = "survival"  # 5W

@dataclass
class EdgeBinding:
    """Hardware binding optimized for edge devices"""
    entity_type: str = "edge_device"
    public_key: str = "pending"  # Will be generated
    hardware_anchor: str = ""  # Jetson device serial
    device_serial: str = ""  # 1421425085368 for our Jetson
    platform: str = "jetson_orin_nano"
    power_budget: int = 15  # Watts
    thermal_zones: int = 7  # Number of temperature sensors
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    binding_proof: str = "pending"

@dataclass
class EdgeBirthCert:
    """Birth certificate for edge society member"""
    citizen_role: str = "lct:web4:role:citizen:sprout"
    context: str = "edge_federation"
    birth_timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    parent_entity: str = "lct:web4:federation:act"  # ACT Federation
    birth_witnesses: List[str] = field(default_factory=list)
    founding_purpose: str = "Edge computing witness and sensor integration"
    edge_capabilities: List[str] = field(default_factory=lambda: [
        "low_power_witness",
        "thermal_aware_compute",
        "resilient_operation",
        "sensor_integration"
    ])

@dataclass
class EdgeMRH:
    """Markov Relevancy Horizon with edge constraints"""
    bound: List[Dict] = field(default_factory=list)  # Permanent bindings
    paired: List[Dict] = field(default_factory=list)  # Active pairings
    witnessing: List[Dict] = field(default_factory=list)  # Witness relationships
    horizon_depth: int = 2  # Shallower for edge (memory constraint)
    cache_size_mb: int = 10  # Edge cache limit
    power_state: str = "balanced"
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())
    
    def add_witness(self, lct_id: str, role: str = "edge_observer"):
        """Add witness relationship with power awareness"""
        self.witnessing.append({
            "lct_id": lct_id,
            "role": role,
            "timestamp": datetime.now().isoformat(),
            "power_cost": 0.5  # Watts
        })
        # Trim if exceeding depth
        if len(self.witnessing) > self.horizon_depth * 10:
            self.witnessing = self.witnessing[-self.horizon_depth * 10:]

@dataclass
class EdgePolicy:
    """Capabilities and constraints for edge operation"""
    capabilities: List[str] = field(default_factory=lambda: [
        "witness",
        "sensor_read",
        "cache_store",
        "resilient_queue"
    ])
    constraints: Dict[str, Any] = field(default_factory=lambda: {
        "max_power_watts": 15,
        "max_temp_celsius": 85,
        "max_memory_mb": 512,
        "max_cache_mb": 100,
        "network_resilient": True,
        "battery_aware": False  # Jetson is DC powered
    })
    thermal_throttle_temp: float = 75.0
    min_witness_interval_seconds: int = 10

@dataclass
class EdgeAttestation:
    """Lightweight attestation for edge witnessing"""
    witness: str  # Our DID
    attestation_type: str  # edge_witness, sensor_reading, thermal_event
    signature: str = "pending"  # Simplified signature
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    evidence_hash: str = ""  # SHA256 of evidence (not full data)
    power_state: str = "balanced"
    temperature: float = 50.0

@dataclass
class EdgeLineage:
    """Simplified lineage for edge nodes"""
    reason: str  # genesis, rotation, recovery
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    parent: Optional[str] = None
    power_event: bool = False  # True if caused by power issue

@dataclass
class EdgeLCT:
    """Edge-optimized Linked Context Token"""
    lct_id: str  # lct:web4:mb32:...
    subject: str  # did:web4:sprout:edge
    binding: EdgeBinding
    mrh: EdgeMRH
    policy: EdgePolicy
    birth_certificate: EdgeBirthCert
    attestations: List[EdgeAttestation] = field(default_factory=list)
    lineage: List[EdgeLineage] = field(default_factory=list)
    revocation: Optional[Dict] = None
    
    # Edge-specific fields
    current_power_state: str = "balanced"
    witness_count: int = 0
    last_thermal_event: Optional[str] = None
    cache_usage_mb: float = 0.0

# ============================================================================
# Trust Tensor Implementation (T3/V3)
# ============================================================================

@dataclass
class TrustTensor:
    """T3 Trust Tensor for edge relationships"""
    entity_id: str
    trust_score: float  # 0.0 to 1.0
    interaction_count: int
    last_interaction: str
    evidence: List[str] = field(default_factory=list)
    decay_rate: float = 0.95  # Trust decays without interaction
    
    def update(self, positive: bool = True):
        """Update trust based on interaction"""
        if positive:
            self.trust_score = min(1.0, self.trust_score * 1.1)
        else:
            self.trust_score = max(0.0, self.trust_score * 0.9)
        self.interaction_count += 1
        self.last_interaction = datetime.now().isoformat()
    
    def apply_decay(self):
        """Apply time-based trust decay"""
        self.trust_score *= self.decay_rate

@dataclass
class ValueTensor:
    """V3 Value Tensor for resource tracking"""
    entity_id: str
    atp_balance: float
    adp_balance: float
    energy_flow_rate: float  # ATP/hour
    value_created: float
    value_consumed: float
    efficiency: float  # value_created / value_consumed
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

# ============================================================================
# R6 Action Framework
# ============================================================================

@dataclass
class R6Action:
    """R6 Action: Rules, Role, Request, Reference, Resource → Result"""
    action_id: str
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    
    # R6 Components
    rules: List[str] = field(default_factory=list)  # Applicable rules
    role: str = ""  # Actor's role
    request: Dict[str, Any] = field(default_factory=dict)  # Action request
    reference: List[str] = field(default_factory=list)  # Context references
    resource: Dict[str, Any] = field(default_factory=dict)  # Required resources
    result: Optional[Dict[str, Any]] = None  # Action outcome
    
    # Edge-specific
    power_cost: float = 0.0  # Watts consumed
    thermal_impact: float = 0.0  # Temperature increase
    success: bool = False
    error: Optional[str] = None

# ============================================================================
# Sprout Edge LCT Manager
# ============================================================================

class SproutLCTManager:
    """
    Manages Sprout's LCT lifecycle with edge optimization
    """
    
    def __init__(self):
        self.storage_path = Path("/home/sprout/.sprout_lct")
        self.storage_path.mkdir(exist_ok=True)
        
        self.lct_file = self.storage_path / "sprout_self_lct.json"
        self.trust_file = self.storage_path / "trust_tensors.json"
        self.value_file = self.storage_path / "value_tensors.json"
        self.r6_log = self.storage_path / "r6_actions.jsonl"
        
        self.hardware_hash = self._get_hardware_hash()
        self.device_serial = self._get_device_serial()
        self.lct = self._load_or_create_lct()
        self.trust_tensors: Dict[str, TrustTensor] = {}
        self.value_tensors: Dict[str, ValueTensor] = {}
        
    def _get_hardware_hash(self) -> str:
        """Get Jetson hardware hash"""
        try:
            # Get device serial and create hash
            serial = self._get_device_serial()
            machine_id = subprocess.run(
                ['cat', '/etc/machine-id'],
                capture_output=True,
                text=True
            ).stdout.strip()
            
            combined = f"{serial}:{machine_id}:jetson_orin_nano"
            return hashlib.sha256(combined.encode()).hexdigest()[:32]
        except:
            return "aaff320ec7bed6eb"  # Fallback to known hash
    
    def _get_device_serial(self) -> str:
        """Get Jetson device serial"""
        try:
            with open('/proc/device-tree/serial-number', 'r') as f:
                return f.read().strip().replace('\x00', '')
        except:
            return "1421425085368"  # Our known serial
    
    def _get_temperature(self) -> float:
        """Get current CPU temperature"""
        try:
            with open('/sys/class/thermal/thermal_zone0/temp', 'r') as f:
                return int(f.read().strip()) / 1000.0
        except:
            return 50.0  # Default
    
    def _determine_power_state(self) -> str:
        """Determine current power state based on temperature"""
        temp = self._get_temperature()
        if temp > 75:
            return PowerState.SURVIVAL.value
        elif temp > 65:
            return PowerState.EFFICIENT.value
        elif temp > 55:
            return PowerState.BALANCED.value
        else:
            return PowerState.MAX_PERF.value
    
    def _generate_lct_id(self) -> str:
        """Generate unique LCT ID"""
        timestamp = int(time.time() * 1000000)
        hash_input = f"sprout:{self.device_serial}:{timestamp}"
        hash_value = hashlib.sha256(hash_input.encode()).hexdigest()[:16]
        return f"lct:web4:mb32:{hash_value}"
    
    def _load_or_create_lct(self) -> EdgeLCT:
        """Load existing LCT or create genesis"""
        if self.lct_file.exists():
            with open(self.lct_file, 'r') as f:
                data = json.load(f)
                # Reconstruct from JSON
                return self._dict_to_lct(data)
        else:
            return self._create_genesis_lct()
    
    def _create_genesis_lct(self) -> EdgeLCT:
        """Create Sprout's genesis self-LCT"""
        print("🌱 Creating Sprout genesis LCT...")
        
        # Create binding
        binding = EdgeBinding(
            hardware_anchor=self.hardware_hash,
            device_serial=self.device_serial,
            public_key=f"mb64:{self.hardware_hash[:16]}"  # Simplified
        )
        
        # Create birth certificate
        birth_cert = EdgeBirthCert()
        
        # Create MRH
        mrh = EdgeMRH(
            power_state=self._determine_power_state()
        )
        
        # Create policy
        policy = EdgePolicy()
        
        # Create lineage
        lineage = [EdgeLineage(reason="genesis")]
        
        # Create LCT
        lct = EdgeLCT(
            lct_id=self._generate_lct_id(),
            subject=f"did:web4:sprout:{self.device_serial[:8]}",
            binding=binding,
            mrh=mrh,
            policy=policy,
            birth_certificate=birth_cert,
            lineage=lineage,
            current_power_state=self._determine_power_state()
        )
        
        # Save
        self._save_lct(lct)
        print(f"✅ Genesis LCT created: {lct.lct_id}")
        
        return lct
    
    def _save_lct(self, lct: EdgeLCT):
        """Save LCT to disk"""
        with open(self.lct_file, 'w') as f:
            json.dump(self._lct_to_dict(lct), f, indent=2)
    
    def _lct_to_dict(self, lct: EdgeLCT) -> Dict:
        """Convert LCT to dictionary for JSON serialization"""
        # Manual conversion to handle nested dataclasses
        return {
            "lct_id": lct.lct_id,
            "subject": lct.subject,
            "binding": asdict(lct.binding),
            "mrh": asdict(lct.mrh),
            "policy": asdict(lct.policy),
            "birth_certificate": asdict(lct.birth_certificate),
            "attestations": [asdict(a) for a in lct.attestations],
            "lineage": [asdict(l) for l in lct.lineage],
            "revocation": lct.revocation,
            "current_power_state": lct.current_power_state,
            "witness_count": lct.witness_count,
            "last_thermal_event": lct.last_thermal_event,
            "cache_usage_mb": lct.cache_usage_mb
        }
    
    def _dict_to_lct(self, data: Dict) -> EdgeLCT:
        """Reconstruct LCT from dictionary"""
        return EdgeLCT(
            lct_id=data["lct_id"],
            subject=data["subject"],
            binding=EdgeBinding(**data["binding"]),
            mrh=EdgeMRH(**data["mrh"]),
            policy=EdgePolicy(**data["policy"]),
            birth_certificate=EdgeBirthCert(**data["birth_certificate"]),
            attestations=[EdgeAttestation(**a) for a in data.get("attestations", [])],
            lineage=[EdgeLineage(**l) for l in data.get("lineage", [])],
            revocation=data.get("revocation"),
            current_power_state=data.get("current_power_state", "balanced"),
            witness_count=data.get("witness_count", 0),
            last_thermal_event=data.get("last_thermal_event"),
            cache_usage_mb=data.get("cache_usage_mb", 0.0)
        )
    
    # ========================================================================
    # Trust and Value Tensor Management
    # ========================================================================
    
    def update_trust(self, entity_id: str, positive: bool = True):
        """Update trust tensor for an entity"""
        if entity_id not in self.trust_tensors:
            self.trust_tensors[entity_id] = TrustTensor(
                entity_id=entity_id,
                trust_score=0.5,
                interaction_count=0,
                last_interaction=datetime.now().isoformat()
            )
        
        self.trust_tensors[entity_id].update(positive)
        self._save_trust_tensors()
    
    def _save_trust_tensors(self):
        """Save trust tensors to disk"""
        with open(self.trust_file, 'w') as f:
            json.dump(
                {k: asdict(v) for k, v in self.trust_tensors.items()},
                f, indent=2
            )
    
    def update_value(self, entity_id: str, atp_delta: float = 0, adp_delta: float = 0):
        """Update value tensor for an entity"""
        if entity_id not in self.value_tensors:
            self.value_tensors[entity_id] = ValueTensor(
                entity_id=entity_id,
                atp_balance=0,
                adp_balance=0,
                energy_flow_rate=0,
                value_created=0,
                value_consumed=0,
                efficiency=1.0
            )
        
        vt = self.value_tensors[entity_id]
        vt.atp_balance += atp_delta
        vt.adp_balance += adp_delta
        
        if atp_delta > 0:
            vt.value_created += atp_delta
        else:
            vt.value_consumed += abs(atp_delta)
        
        if vt.value_consumed > 0:
            vt.efficiency = vt.value_created / vt.value_consumed
        
        vt.last_updated = datetime.now().isoformat()
        self._save_value_tensors()
    
    def _save_value_tensors(self):
        """Save value tensors to disk"""
        with open(self.value_file, 'w') as f:
            json.dump(
                {k: asdict(v) for k, v in self.value_tensors.items()},
                f, indent=2
            )
    
    # ========================================================================
    # R6 Action Framework
    # ========================================================================
    
    def execute_r6_action(self, action: R6Action) -> R6Action:
        """Execute an R6 action with edge constraints"""
        print(f"\n🎯 Executing R6 Action: {action.action_id}")
        
        # Check power budget
        current_power = self._determine_power_state()
        if current_power == PowerState.SURVIVAL.value and action.power_cost > 2:
            action.success = False
            action.error = "Insufficient power for action"
            action.result = {"error": "Power constraint violation"}
            self._log_r6_action(action)
            return action
        
        # Check thermal
        temp = self._get_temperature()
        if temp + action.thermal_impact > self.lct.policy.thermal_throttle_temp:
            action.success = False
            action.error = "Thermal limit would be exceeded"
            action.result = {"error": "Thermal constraint violation"}
            self._log_r6_action(action)
            return action
        
        # Execute based on request type
        if action.request.get("type") == "witness":
            action = self._execute_witness(action)
        elif action.request.get("type") == "attest":
            action = self._execute_attestation(action)
        elif action.request.get("type") == "pair":
            action = self._execute_pairing(action)
        else:
            action.success = False
            action.error = "Unknown action type"
        
        self._log_r6_action(action)
        return action
    
    def _execute_witness(self, action: R6Action) -> R6Action:
        """Execute witness action"""
        target = action.request.get("target")
        if not target:
            action.error = "No target specified"
            return action
        
        # Add to MRH
        self.lct.mrh.add_witness(target, "edge_witness")
        self.lct.witness_count += 1
        
        # Create attestation
        attestation = EdgeAttestation(
            witness=self.lct.subject,
            attestation_type="edge_witness",
            evidence_hash=hashlib.sha256(str(target).encode()).hexdigest(),
            power_state=self.lct.current_power_state,
            temperature=self._get_temperature()
        )
        self.lct.attestations.append(attestation)
        
        # Update trust
        self.update_trust(target, True)
        
        action.success = True
        action.result = {
            "witness_count": self.lct.witness_count,
            "attestation_id": attestation.timestamp
        }
        
        self._save_lct(self.lct)
        return action
    
    def _execute_attestation(self, action: R6Action) -> R6Action:
        """Execute attestation action"""
        attestation_type = action.request.get("attestation_type", "general")
        evidence = action.request.get("evidence", "")
        
        attestation = EdgeAttestation(
            witness=self.lct.subject,
            attestation_type=attestation_type,
            evidence_hash=hashlib.sha256(str(evidence).encode()).hexdigest(),
            power_state=self.lct.current_power_state,
            temperature=self._get_temperature()
        )
        
        self.lct.attestations.append(attestation)
        
        # Trim old attestations if too many (memory constraint)
        if len(self.lct.attestations) > 100:
            self.lct.attestations = self.lct.attestations[-100:]
        
        action.success = True
        action.result = {"attestation_id": attestation.timestamp}
        
        self._save_lct(self.lct)
        return action
    
    def _execute_pairing(self, action: R6Action) -> R6Action:
        """Execute pairing action"""
        target = action.request.get("target")
        pairing_type = action.request.get("pairing_type", "operational")
        
        pairing = {
            "lct_id": target,
            "pairing_type": pairing_type,
            "permanent": False,
            "timestamp": datetime.now().isoformat(),
            "power_state": self.lct.current_power_state
        }
        
        self.lct.mrh.paired.append(pairing)
        
        # Trim old pairings
        if len(self.lct.mrh.paired) > 20:
            self.lct.mrh.paired = self.lct.mrh.paired[-20:]
        
        action.success = True
        action.result = {"pairing_id": pairing["timestamp"]}
        
        self._save_lct(self.lct)
        return action
    
    def _log_r6_action(self, action: R6Action):
        """Log R6 action to append-only log"""
        with open(self.r6_log, 'a') as f:
            f.write(json.dumps(asdict(action)) + '\n')
    
    # ========================================================================
    # Federation Integration
    # ========================================================================
    
    def create_birth_certificate_request(self) -> Dict:
        """Create birth certificate request for federation witnessing"""
        return {
            "type": "BIRTH_CERTIFICATE_REQUEST",
            "lct_id": self.lct.lct_id,
            "subject": self.lct.subject,
            "hardware_hash": self.lct.binding.hardware_anchor,
            "device_serial": self.lct.binding.device_serial,
            "platform": self.lct.binding.platform,
            "founding_purpose": self.lct.birth_certificate.founding_purpose,
            "capabilities": self.lct.birth_certificate.edge_capabilities,
            "power_budget": self.lct.binding.power_budget,
            "request_witness": ["genesis", "society4", "cbp"],
            "timestamp": datetime.now().isoformat()
        }
    
    def witness_birth_certificate(self, entity_id: str, cert_data: Dict) -> bool:
        """Witness another entity's birth certificate"""
        # Create R6 action for witnessing
        action = R6Action(
            action_id=f"witness_birth_{entity_id[:8]}",
            rules=["federation_witness_protocol"],
            role="edge_witness",
            request={
                "type": "witness",
                "target": entity_id,
                "cert_data": cert_data
            },
            resource={"atp": 10, "power": 1.0},
            power_cost=1.0,
            thermal_impact=0.5
        )
        
        result = self.execute_r6_action(action)
        
        if result.success:
            # Add to birth witnesses
            self.lct.birth_certificate.birth_witnesses.append(entity_id)
            self._save_lct(self.lct)
            print(f"✅ Witnessed birth certificate for {entity_id}")
            return True
        else:
            print(f"❌ Failed to witness: {result.error}")
            return False
    
    def get_status(self) -> Dict:
        """Get current LCT and system status"""
        return {
            "lct_id": self.lct.lct_id,
            "subject": self.lct.subject,
            "hardware_hash": self.lct.binding.hardware_anchor[:16] + "...",
            "device_serial": self.lct.binding.device_serial,
            "power_state": self.lct.current_power_state,
            "temperature": self._get_temperature(),
            "witness_count": self.lct.witness_count,
            "attestations": len(self.lct.attestations),
            "trust_relationships": len(self.trust_tensors),
            "value_relationships": len(self.value_tensors),
            "mrh_witnesses": len(self.lct.mrh.witnessing),
            "mrh_paired": len(self.lct.mrh.paired),
            "cache_usage_mb": self.lct.cache_usage_mb
        }

# ============================================================================
# CLI Interface
# ============================================================================

def main():
    """Command-line interface for Sprout LCT management"""
    import sys
    
    manager = SproutLCTManager()
    
    if len(sys.argv) < 2:
        command = "status"
    else:
        command = sys.argv[1]
    
    if command == "status":
        status = manager.get_status()
        print("\n🌱 Sprout LCT Status")
        print("=" * 50)
        for key, value in status.items():
            print(f"{key:20}: {value}")
    
    elif command == "birth":
        # Create birth certificate request
        request = manager.create_birth_certificate_request()
        output_file = Path("/home/sprout/ai-workspace/ACT/implementation/ledger/federation_outbox/SPROUT_BIRTH_CERTIFICATE.json")
        output_file.parent.mkdir(exist_ok=True)
        with open(output_file, 'w') as f:
            json.dump(request, f, indent=2)
        print(f"✅ Birth certificate request created: {output_file}")
    
    elif command == "witness":
        if len(sys.argv) < 3:
            print("Usage: sprout_web4_lct.py witness <entity_id>")
            sys.exit(1)
        entity_id = sys.argv[2]
        cert_data = {}  # Would load from federation
        success = manager.witness_birth_certificate(entity_id, cert_data)
        if success:
            print(f"✅ Witnessed {entity_id}")
        else:
            print(f"❌ Failed to witness {entity_id}")
    
    elif command == "trust":
        if len(sys.argv) < 4:
            print("Usage: sprout_web4_lct.py trust <entity_id> <positive|negative>")
            sys.exit(1)
        entity_id = sys.argv[2]
        positive = sys.argv[3].lower() == "positive"
        manager.update_trust(entity_id, positive)
        print(f"✅ Updated trust for {entity_id}")
    
    elif command == "value":
        if len(sys.argv) < 5:
            print("Usage: sprout_web4_lct.py value <entity_id> <atp_delta> <adp_delta>")
            sys.exit(1)
        entity_id = sys.argv[2]
        atp_delta = float(sys.argv[3])
        adp_delta = float(sys.argv[4])
        manager.update_value(entity_id, atp_delta, adp_delta)
        print(f"✅ Updated value tensor for {entity_id}")
    
    elif command == "r6":
        # Execute a test R6 action
        action = R6Action(
            action_id="test_action_" + str(int(time.time())),
            rules=["test_rule"],
            role="edge_node",
            request={"type": "witness", "target": "test_entity"},
            resource={"atp": 5},
            power_cost=0.5,
            thermal_impact=1.0
        )
        result = manager.execute_r6_action(action)
        print(f"R6 Action {'succeeded' if result.success else 'failed'}")
        if result.error:
            print(f"Error: {result.error}")
    
    else:
        print("🌱 Sprout Web4 LCT Manager")
        print("\nCommands:")
        print("  status  - Show LCT and system status")
        print("  birth   - Create birth certificate request")
        print("  witness - Witness another entity's birth")
        print("  trust   - Update trust tensor")
        print("  value   - Update value tensor")
        print("  r6      - Execute test R6 action")
        print("\nOptimized for Jetson Orin Nano (15W)")

if __name__ == "__main__":
    main()
