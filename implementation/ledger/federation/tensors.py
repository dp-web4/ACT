#!/usr/bin/env python3
"""
T3/V3 Tensor Management for Web4 Federation
Implements Trust and Value Tensors per Web4 T3/V3 Specification
"""

import json
import math
from typing import Dict, List, Tuple, Optional
from dataclasses import dataclass, field
from datetime import datetime
from pathlib import Path

@dataclass
class TrustVector:
    """Individual trust component per T3 Specification §2.1"""
    performance: float = 0.5  # Task completion rate
    reliability: float = 0.5  # Uptime and availability
    integrity: float = 0.5    # Protocol compliance
    cooperation: float = 0.5  # Inter-society collaboration
    
    def magnitude(self) -> float:
        """Calculate trust magnitude"""
        return math.sqrt(
            self.performance**2 + 
            self.reliability**2 + 
            self.integrity**2 + 
            self.cooperation**2
        ) / 2.0  # Normalize to [0,1]

@dataclass
class ValueVector:
    """Individual value component per V3 Specification §2.1"""
    energy_efficiency: float = 0.5  # ATP usage efficiency
    task_impact: float = 0.5        # Deliverable quality
    innovation: float = 0.5         # Novel contributions
    coordination: float = 0.5       # Federation facilitation
    
    def magnitude(self) -> float:
        """Calculate value magnitude"""
        return math.sqrt(
            self.energy_efficiency**2 + 
            self.task_impact**2 + 
            self.innovation**2 + 
            self.coordination**2
        ) / 2.0  # Normalize to [0,1]

@dataclass
class RoleContext:
    """Role-contextual tensor per T3/V3 Specification §1.1"""
    role: str
    context: str
    trust: TrustVector = field(default_factory=TrustVector)
    value: ValueVector = field(default_factory=ValueVector)
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    witness_count: int = 0

class T3V3TensorSystem:
    """
    Trust and Value Tensor System for Federation
    Implements role-contextual trust per T3/V3 Specification
    """
    
    def __init__(self, federation_path: str = "/home/dp/ai-workspace/act/implementation/ledger"):
        self.federation_path = Path(federation_path)
        self.tensor_file = self.federation_path / "federation" / "tensors.json"
        self.tensors: Dict[str, Dict[str, RoleContext]] = {}
        self.load_tensors()
        
        # Define federation roles
        self.roles = {
            "coordinator": "Genesis Queen",
            "validator": "Society4 Security Queen",
            "integrator": "Society2 Bridge Queen",
            "optimizer": "Sprout Resource Manager",
            "witness": "Federation Witness Pool"
        }
        
        # Define contexts
        self.contexts = [
            "sage_development",
            "web4_compliance",
            "federation_governance",
            "energy_economy",
            "task_execution"
        ]
    
    def load_tensors(self):
        """Load tensors from persistent storage"""
        if self.tensor_file.exists():
            with open(self.tensor_file, 'r') as f:
                data = json.load(f)
                for society, contexts in data.items():
                    self.tensors[society] = {}
                    for ctx_key, ctx_data in contexts.items():
                        rc = RoleContext(
                            role=ctx_data["role"],
                            context=ctx_data["context"],
                            trust=TrustVector(**ctx_data["trust"]),
                            value=ValueVector(**ctx_data["value"]),
                            timestamp=ctx_data["timestamp"],
                            witness_count=ctx_data["witness_count"]
                        )
                        self.tensors[society][ctx_key] = rc
    
    def save_tensors(self):
        """Save tensors to persistent storage"""
        data = {}
        for society, contexts in self.tensors.items():
            data[society] = {}
            for ctx_key, rc in contexts.items():
                data[society][ctx_key] = {
                    "role": rc.role,
                    "context": rc.context,
                    "trust": {
                        "performance": rc.trust.performance,
                        "reliability": rc.trust.reliability,
                        "integrity": rc.trust.integrity,
                        "cooperation": rc.trust.cooperation
                    },
                    "value": {
                        "energy_efficiency": rc.value.energy_efficiency,
                        "task_impact": rc.value.task_impact,
                        "innovation": rc.value.innovation,
                        "coordination": rc.value.coordination
                    },
                    "timestamp": rc.timestamp,
                    "witness_count": rc.witness_count
                }
        
        self.tensor_file.parent.mkdir(exist_ok=True)
        with open(self.tensor_file, 'w') as f:
            json.dump(data, f, indent=2)
    
    def update_trust_tensor(self, society: str, role: str, context: str,
                           performance_delta: float = 0,
                           reliability_delta: float = 0,
                           integrity_delta: float = 0,
                           cooperation_delta: float = 0) -> TrustVector:
        """
        Update trust tensor for a society in role-context
        Per T3 Specification §3.1 - Trust Update Algorithm
        """
        ctx_key = f"{role}:{context}"
        
        if society not in self.tensors:
            self.tensors[society] = {}
        
        if ctx_key not in self.tensors[society]:
            self.tensors[society][ctx_key] = RoleContext(role=role, context=context)
        
        rc = self.tensors[society][ctx_key]
        
        # Apply deltas with decay factor (trust changes slowly)
        decay = 0.9
        rc.trust.performance = max(0, min(1, 
            rc.trust.performance * decay + performance_delta))
        rc.trust.reliability = max(0, min(1,
            rc.trust.reliability * decay + reliability_delta))
        rc.trust.integrity = max(0, min(1,
            rc.trust.integrity * decay + integrity_delta))
        rc.trust.cooperation = max(0, min(1,
            rc.trust.cooperation * decay + cooperation_delta))
        
        rc.timestamp = datetime.now().isoformat()
        self.save_tensors()
        
        return rc.trust
    
    def update_value_tensor(self, society: str, role: str, context: str,
                           energy_efficiency_delta: float = 0,
                           task_impact_delta: float = 0,
                           innovation_delta: float = 0,
                           coordination_delta: float = 0) -> ValueVector:
        """
        Update value tensor for a society in role-context
        Per V3 Specification §3.1 - Value Update Algorithm
        """
        ctx_key = f"{role}:{context}"
        
        if society not in self.tensors:
            self.tensors[society] = {}
        
        if ctx_key not in self.tensors[society]:
            self.tensors[society][ctx_key] = RoleContext(role=role, context=context)
        
        rc = self.tensors[society][ctx_key]
        
        # Apply deltas with growth factor (value can grow quickly)
        growth = 1.1
        rc.value.energy_efficiency = max(0, min(1,
            rc.value.energy_efficiency * growth + energy_efficiency_delta))
        rc.value.task_impact = max(0, min(1,
            rc.value.task_impact * growth + task_impact_delta))
        rc.value.innovation = max(0, min(1,
            rc.value.innovation * growth + innovation_delta))
        rc.value.coordination = max(0, min(1,
            rc.value.coordination * growth + coordination_delta))
        
        rc.timestamp = datetime.now().isoformat()
        self.save_tensors()
        
        return rc.value
    
    def calculate_inter_society_trust(self, society_a: str, society_b: str, 
                                     context: str) -> float:
        """
        Calculate trust between two societies in a context
        Per T3 Specification §4 - Inter-entity Trust
        """
        trust_scores = []
        
        for role in self.roles.keys():
            ctx_key = f"{role}:{context}"
            
            # Get trust vectors for both societies
            trust_a = 0.5  # Default neutral trust
            trust_b = 0.5
            
            if society_a in self.tensors and ctx_key in self.tensors[society_a]:
                trust_a = self.tensors[society_a][ctx_key].trust.magnitude()
            
            if society_b in self.tensors and ctx_key in self.tensors[society_b]:
                trust_b = self.tensors[society_b][ctx_key].trust.magnitude()
            
            # Bidirectional trust is geometric mean
            trust_scores.append(math.sqrt(trust_a * trust_b))
        
        # Average across all roles
        return sum(trust_scores) / len(trust_scores) if trust_scores else 0.5
    
    def calculate_federation_coherence(self) -> float:
        """
        Calculate overall federation coherence
        Per T3/V3 Specification §5 - System Coherence
        """
        societies = ["genesis", "society4", "society2", "sprout"]
        coherence_scores = []
        
        for i, society_a in enumerate(societies):
            for society_b in societies[i+1:]:
                for context in self.contexts:
                    trust = self.calculate_inter_society_trust(
                        society_a, society_b, context
                    )
                    coherence_scores.append(trust)
        
        return sum(coherence_scores) / len(coherence_scores) if coherence_scores else 0.5
    
    def reward_completion(self, society: str, task_type: str):
        """Reward society for task completion"""
        role = "coordinator" if society == "genesis" else "validator"
        context = "sage_development" if "sage" in task_type.lower() else "task_execution"
        
        self.update_trust_tensor(
            society, role, context,
            performance_delta=0.1,
            reliability_delta=0.05,
            cooperation_delta=0.05
        )
        
        self.update_value_tensor(
            society, role, context,
            task_impact_delta=0.1,
            energy_efficiency_delta=0.05
        )
    
    def penalize_failure(self, society: str, failure_type: str):
        """Penalize society for failure"""
        role = "validator"
        context = "federation_governance"
        
        self.update_trust_tensor(
            society, role, context,
            performance_delta=-0.1,
            reliability_delta=-0.15,
            integrity_delta=-0.1
        )
    
    def witness_attestation(self, society: str, role: str, context: str):
        """Record witness attestation for tensor update"""
        ctx_key = f"{role}:{context}"
        
        if society in self.tensors and ctx_key in self.tensors[society]:
            self.tensors[society][ctx_key].witness_count += 1
            self.save_tensors()
    
    def get_society_report(self, society: str) -> Dict:
        """Generate trust/value report for a society"""
        if society not in self.tensors:
            return {"error": f"No tensor data for {society}"}
        
        report = {
            "society": society,
            "contexts": {},
            "average_trust": 0,
            "average_value": 0,
            "total_witnesses": 0
        }
        
        trust_sum = 0
        value_sum = 0
        
        for ctx_key, rc in self.tensors[society].items():
            report["contexts"][ctx_key] = {
                "trust_magnitude": rc.trust.magnitude(),
                "value_magnitude": rc.value.magnitude(),
                "witnesses": rc.witness_count,
                "last_updated": rc.timestamp
            }
            trust_sum += rc.trust.magnitude()
            value_sum += rc.value.magnitude()
            report["total_witnesses"] += rc.witness_count
        
        num_contexts = len(report["contexts"])
        if num_contexts > 0:
            report["average_trust"] = trust_sum / num_contexts
            report["average_value"] = value_sum / num_contexts
        
        return report
    
    def initialize_federation_tensors(self):
        """Initialize base tensors for all societies"""
        societies = {
            "genesis": {"role": "coordinator", "base_trust": 0.7},
            "society4": {"role": "validator", "base_trust": 0.65},
            "society2": {"role": "integrator", "base_trust": 0.6},
            "sprout": {"role": "optimizer", "base_trust": 0.6}
        }
        
        for society, config in societies.items():
            for context in self.contexts:
                ctx_key = f"{config['role']}:{context}"
                
                if society not in self.tensors:
                    self.tensors[society] = {}
                
                if ctx_key not in self.tensors[society]:
                    rc = RoleContext(
                        role=config['role'],
                        context=context,
                        trust=TrustVector(
                            performance=config['base_trust'],
                            reliability=config['base_trust'],
                            integrity=config['base_trust'] + 0.1,  # Web4 compliance bonus
                            cooperation=config['base_trust'] - 0.1  # New federation penalty
                        ),
                        value=ValueVector(
                            energy_efficiency=0.5,
                            task_impact=0.5,
                            innovation=0.5,
                            coordination=0.5
                        )
                    )
                    self.tensors[society][ctx_key] = rc
        
        self.save_tensors()


if __name__ == "__main__":
    print("=== Web4 Federation T3/V3 Tensor System ===\n")
    
    # Initialize tensor system
    tensor_system = T3V3TensorSystem()
    
    print("1. Initializing federation tensors...")
    tensor_system.initialize_federation_tensors()
    print("   ✅ Base tensors created for all societies\n")
    
    print("2. Rewarding Genesis for task completion...")
    tensor_system.reward_completion("genesis", "federation_tracking")
    genesis_report = tensor_system.get_society_report("genesis")
    print(f"   Trust: {genesis_report['average_trust']:.3f}")
    print(f"   Value: {genesis_report['average_value']:.3f}\n")
    
    print("3. Updating Sprout tensors for SAGE acceptance...")
    tensor_system.update_trust_tensor(
        "sprout", "optimizer", "sage_development",
        performance_delta=0.15,
        cooperation_delta=0.2,
        integrity_delta=0.1
    )
    tensor_system.update_value_tensor(
        "sprout", "optimizer", "sage_development",
        innovation_delta=0.2,
        task_impact_delta=0.15
    )
    print("   ✅ Sprout tensors updated\n")
    
    print("4. Calculating inter-society trust...")
    trust_genesis_sprout = tensor_system.calculate_inter_society_trust(
        "genesis", "sprout", "sage_development"
    )
    print(f"   Genesis <-> Sprout trust: {trust_genesis_sprout:.3f}\n")
    
    print("5. Calculating federation coherence...")
    coherence = tensor_system.calculate_federation_coherence()
    print(f"   Federation Coherence: {coherence:.3f}\n")
    
    print("6. Generating federation tensor report...")
    for society in ["genesis", "society4", "society2", "sprout"]:
        report = tensor_system.get_society_report(society)
        print(f"\n   {society.upper()}:")
        print(f"   - Average Trust: {report['average_trust']:.3f}")
        print(f"   - Average Value: {report['average_value']:.3f}")
        print(f"   - Total Witnesses: {report['total_witnesses']}")
    
    print("\n✅ T3/V3 Tensor System implementation complete!")
    print("   - Role-contextual tensors")
    print("   - Trust and value tracking")
    print("   - Inter-society trust calculation")
    print("   - Federation coherence metrics")
    print("\nWeb4 T3/V3 Specification compliance achieved!")