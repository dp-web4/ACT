#!/usr/bin/env python3
"""
CBP Trust Tensor Implementation (T3/V3)
Web4-compliant multi-dimensional trust and value tracking
Based on Web4 specification and Society4 patterns
"""

import json
import math
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Tuple, Any
from dataclasses import dataclass, field, asdict
from pathlib import Path
import numpy as np

# ============================================================================
# T3: Talent, Training, Temperament (Trust Tensors)
# ============================================================================

@dataclass
class TalentDimension:
    """Natural capabilities and aptitudes"""
    technical_skill: float = 0.5  # 0-1: Technical competence
    innovation: float = 0.5       # 0-1: Creative problem solving
    efficiency: float = 0.5       # 0-1: Resource optimization
    accuracy: float = 0.5         # 0-1: Error rate inverse
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class TrainingDimension:
    """Learned behaviors and improvements"""
    tasks_completed: int = 0
    success_rate: float = 0.0
    learning_rate: float = 0.5    # 0-1: Speed of improvement
    adaptation_score: float = 0.5  # 0-1: Flexibility to new tasks
    specializations: List[str] = field(default_factory=list)
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class TemperamentDimension:
    """Behavioral patterns and reliability"""
    consistency: float = 0.5      # 0-1: Predictability of behavior
    cooperation: float = 0.5      # 0-1: Federation collaboration
    responsiveness: float = 0.5   # 0-1: Reaction time to requests
    stability: float = 0.5        # 0-1: Emotional/operational stability
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class T3Tensor:
    """Trust Tensor combining three dimensions"""
    entity_id: str
    talent: TalentDimension = field(default_factory=TalentDimension)
    training: TrainingDimension = field(default_factory=TrainingDimension)
    temperament: TemperamentDimension = field(default_factory=TemperamentDimension)
    composite_trust: float = 0.5
    tensor_version: int = 1
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now().isoformat())

    def calculate_composite(self) -> float:
        """Calculate composite trust score from dimensions"""
        # Weighted average with CBP-specific priorities
        talent_score = (
            self.talent.technical_skill * 0.3 +
            self.talent.innovation * 0.2 +
            self.talent.efficiency * 0.3 +
            self.talent.accuracy * 0.2
        )

        training_score = (
            min(self.training.success_rate, 1.0) * 0.4 +
            self.training.learning_rate * 0.3 +
            self.training.adaptation_score * 0.3
        )

        temperament_score = (
            self.temperament.consistency * 0.3 +
            self.temperament.cooperation * 0.3 +
            self.temperament.responsiveness * 0.2 +
            self.temperament.stability * 0.2
        )

        # CBP weights: Training > Talent > Temperament (data-focused)
        self.composite_trust = (
            talent_score * 0.3 +
            training_score * 0.4 +
            temperament_score * 0.3
        )

        return self.composite_trust

# ============================================================================
# V3: Value, Veracity, Velocity (Value Tensors)
# ============================================================================

@dataclass
class ValueDimension:
    """Economic and utility value created"""
    atp_generated: float = 0.0    # Total ATP value created
    cache_hits: int = 0           # Cache value provided
    metrics_collected: int = 0     # Metrics value generated
    bridges_established: int = 0   # Federation connections
    utility_score: float = 0.5    # 0-1: Overall utility to federation
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class VeracityDimension:
    """Truth and accuracy of information"""
    data_accuracy: float = 0.5    # 0-1: Correctness of stored data
    metric_precision: float = 0.5 # 0-1: Precision of measurements
    false_positive_rate: float = 0.0  # Error rate
    attestation_validity: float = 0.5  # 0-1: Witness trustworthiness
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class VelocityDimension:
    """Speed and momentum of value creation"""
    transaction_rate: float = 0.0  # Transactions per hour
    response_time_ms: float = 1000.0  # Average response time
    throughput: float = 0.0       # Operations per second
    growth_rate: float = 0.0      # Month-over-month improvement
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class V3Tensor:
    """Value Tensor combining three dimensions"""
    entity_id: str
    value: ValueDimension = field(default_factory=ValueDimension)
    veracity: VeracityDimension = field(default_factory=VeracityDimension)
    velocity: VelocityDimension = field(default_factory=VelocityDimension)
    composite_value: float = 0.5
    tensor_version: int = 1
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now().isoformat())

    def calculate_composite(self) -> float:
        """Calculate composite value score from dimensions"""
        # Normalize value dimension
        value_score = min(
            (self.value.atp_generated / 1000.0) * 0.3 +
            (self.value.cache_hits / 1000.0) * 0.2 +
            (self.value.metrics_collected / 100.0) * 0.2 +
            (self.value.bridges_established / 10.0) * 0.1 +
            self.value.utility_score * 0.2,
            1.0
        )

        veracity_score = (
            self.veracity.data_accuracy * 0.4 +
            self.veracity.metric_precision * 0.3 +
            (1.0 - self.veracity.false_positive_rate) * 0.2 +
            self.veracity.attestation_validity * 0.1
        )

        # Normalize velocity (lower response time is better)
        velocity_score = (
            min(self.velocity.transaction_rate / 100.0, 1.0) * 0.3 +
            max(0, 1.0 - (self.velocity.response_time_ms / 5000.0)) * 0.3 +
            min(self.velocity.throughput / 1000.0, 1.0) * 0.2 +
            min(max(self.velocity.growth_rate, -1.0), 1.0) * 0.2
        )

        # CBP weights: Veracity > Value > Velocity (accuracy critical for cache)
        self.composite_value = (
            value_score * 0.3 +
            veracity_score * 0.4 +
            velocity_score * 0.3
        )

        return self.composite_value

# ============================================================================
# Fractal Attribution System
# ============================================================================

@dataclass
class FractalAttribution:
    """Track value flow through fractal relationships"""
    source_entity: str
    target_entity: str
    attribution_type: str  # direct, indirect, cascading
    value_transferred: float
    trust_impact: float
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    depth: int = 0  # Fractal depth (0=direct, 1=first indirect, etc.)

class CBPTrustTensorSystem:
    """
    CBP's implementation of Web4 Trust and Value Tensors
    Focuses on cache, metrics, and bridging operations
    """

    def __init__(self):
        self.t3_tensors: Dict[str, T3Tensor] = {}
        self.v3_tensors: Dict[str, V3Tensor] = {}
        self.attributions: List[FractalAttribution] = []
        self.storage_path = Path("implementation/cbp-chain/trust_tensors.json")

    def create_entity_tensors(self, entity_id: str) -> Tuple[T3Tensor, V3Tensor]:
        """Create initial T3 and V3 tensors for an entity"""
        t3 = T3Tensor(entity_id=entity_id)
        v3 = V3Tensor(entity_id=entity_id)

        self.t3_tensors[entity_id] = t3
        self.v3_tensors[entity_id] = v3

        return t3, v3

    def update_talent(self, entity_id: str, skill: float = None, innovation: float = None,
                     efficiency: float = None, accuracy: float = None) -> T3Tensor:
        """Update talent dimension of trust tensor"""
        if entity_id not in self.t3_tensors:
            self.create_entity_tensors(entity_id)

        t3 = self.t3_tensors[entity_id]
        if skill is not None:
            t3.talent.technical_skill = max(0, min(1, skill))
        if innovation is not None:
            t3.talent.innovation = max(0, min(1, innovation))
        if efficiency is not None:
            t3.talent.efficiency = max(0, min(1, efficiency))
        if accuracy is not None:
            t3.talent.accuracy = max(0, min(1, accuracy))

        t3.talent.last_updated = datetime.now().isoformat()
        t3.calculate_composite()
        t3.updated_at = datetime.now().isoformat()
        t3.tensor_version += 1

        return t3

    def update_training(self, entity_id: str, task_completed: bool = None,
                       learning_delta: float = 0.0, new_specialization: str = None) -> T3Tensor:
        """Update training dimension based on task outcomes"""
        if entity_id not in self.t3_tensors:
            self.create_entity_tensors(entity_id)

        t3 = self.t3_tensors[entity_id]

        if task_completed is not None:
            t3.training.tasks_completed += 1
            if task_completed:
                # Update success rate with exponential moving average
                alpha = 0.1  # Learning rate
                t3.training.success_rate = (1 - alpha) * t3.training.success_rate + alpha
            else:
                t3.training.success_rate = (1 - alpha) * t3.training.success_rate

        if learning_delta != 0:
            t3.training.learning_rate = max(0, min(1, t3.training.learning_rate + learning_delta))

        if new_specialization and new_specialization not in t3.training.specializations:
            t3.training.specializations.append(new_specialization)
            t3.training.adaptation_score = min(1, t3.training.adaptation_score + 0.1)

        t3.training.last_updated = datetime.now().isoformat()
        t3.calculate_composite()
        t3.updated_at = datetime.now().isoformat()
        t3.tensor_version += 1

        return t3

    def update_value_from_cache_operation(self, entity_id: str, cache_hit: bool,
                                         response_time_ms: float) -> V3Tensor:
        """Update value tensor from cache operation"""
        if entity_id not in self.v3_tensors:
            self.create_entity_tensors(entity_id)

        v3 = self.v3_tensors[entity_id]

        if cache_hit:
            v3.value.cache_hits += 1
            v3.value.atp_generated += 0.1  # Small ATP value for cache hit

        # Update velocity metrics
        v3.velocity.transaction_rate = v3.value.cache_hits / max(1,
            (datetime.now() - datetime.fromisoformat(v3.created_at)).total_seconds() / 3600)

        # Exponential moving average for response time
        alpha = 0.1
        v3.velocity.response_time_ms = (1 - alpha) * v3.velocity.response_time_ms + alpha * response_time_ms

        v3.value.last_updated = datetime.now().isoformat()
        v3.velocity.last_updated = datetime.now().isoformat()
        v3.calculate_composite()
        v3.updated_at = datetime.now().isoformat()
        v3.tensor_version += 1

        return v3

    def add_fractal_attribution(self, source: str, target: str, value: float,
                               attribution_type: str = "direct", depth: int = 0) -> None:
        """Add fractal value attribution between entities"""
        # Update source entity's value contribution
        if source in self.v3_tensors:
            self.v3_tensors[source].value.atp_generated += value * 0.1

        # Update target entity's trust based on value received
        if target in self.t3_tensors:
            trust_delta = value * 0.001  # Small trust increase from value
            self.t3_tensors[target].temperament.cooperation = min(1,
                self.t3_tensors[target].temperament.cooperation + trust_delta)

        attribution = FractalAttribution(
            source_entity=source,
            target_entity=target,
            attribution_type=attribution_type,
            value_transferred=value,
            trust_impact=value * 0.001,
            depth=depth
        )

        self.attributions.append(attribution)

        # Cascade attribution fractally (reduced value at each level)
        if depth < 3 and value > 0.1:  # Max depth 3, min value 0.1
            # Find entities connected to target
            connected = self._find_connected_entities(target)
            for connected_entity in connected[:3]:  # Max 3 cascades
                self.add_fractal_attribution(
                    source=target,
                    target=connected_entity,
                    value=value * 0.5,  # Halve value at each level
                    attribution_type="cascading",
                    depth=depth + 1
                )

    def _find_connected_entities(self, entity_id: str) -> List[str]:
        """Find entities connected through attributions"""
        connected = set()
        for attr in self.attributions[-100:]:  # Look at recent attributions
            if attr.source_entity == entity_id:
                connected.add(attr.target_entity)
            elif attr.target_entity == entity_id:
                connected.add(attr.source_entity)
        return list(connected)

    def get_federation_trust_summary(self) -> Dict[str, Any]:
        """Get summary of trust across federation"""
        if not self.t3_tensors:
            return {"status": "no_entities"}

        trust_scores = [t3.composite_trust for t3 in self.t3_tensors.values()]
        value_scores = [v3.composite_value for v3 in self.v3_tensors.values()]

        return {
            "entity_count": len(self.t3_tensors),
            "average_trust": sum(trust_scores) / len(trust_scores),
            "average_value": sum(value_scores) / len(value_scores) if value_scores else 0,
            "highest_trust": max(trust_scores) if trust_scores else 0,
            "lowest_trust": min(trust_scores) if trust_scores else 0,
            "total_attributions": len(self.attributions),
            "timestamp": datetime.now().isoformat()
        }

    def save_state(self) -> None:
        """Save tensor state to disk"""
        state = {
            "t3_tensors": {k: asdict(v) for k, v in self.t3_tensors.items()},
            "v3_tensors": {k: asdict(v) for k, v in self.v3_tensors.items()},
            "attributions": [asdict(a) for a in self.attributions[-1000:]],  # Keep last 1000
            "summary": self.get_federation_trust_summary()
        }

        self.storage_path.parent.mkdir(parents=True, exist_ok=True)
        with open(self.storage_path, 'w') as f:
            json.dump(state, f, indent=2)

    def display_entity_trust(self, entity_id: str) -> None:
        """Display trust and value tensors for an entity"""
        if entity_id not in self.t3_tensors:
            print(f"Entity {entity_id} not found")
            return

        t3 = self.t3_tensors[entity_id]
        v3 = self.v3_tensors.get(entity_id)

        print(f"\n{'='*60}")
        print(f"🔮 Trust & Value Tensors: {entity_id}")
        print(f"{'='*60}")

        print(f"\n📊 T3 Trust Tensor (Composite: {t3.composite_trust:.3f})")
        print(f"  Talent:")
        print(f"    Technical: {t3.talent.technical_skill:.2f}  Innovation: {t3.talent.innovation:.2f}")
        print(f"    Efficiency: {t3.talent.efficiency:.2f}  Accuracy: {t3.talent.accuracy:.2f}")
        print(f"  Training:")
        print(f"    Tasks: {t3.training.tasks_completed}  Success Rate: {t3.training.success_rate:.2%}")
        print(f"    Learning: {t3.training.learning_rate:.2f}  Adaptation: {t3.training.adaptation_score:.2f}")
        print(f"  Temperament:")
        print(f"    Consistency: {t3.temperament.consistency:.2f}  Cooperation: {t3.temperament.cooperation:.2f}")
        print(f"    Responsive: {t3.temperament.responsiveness:.2f}  Stability: {t3.temperament.stability:.2f}")

        if v3:
            print(f"\n💎 V3 Value Tensor (Composite: {v3.composite_value:.3f})")
            print(f"  Value:")
            print(f"    ATP Generated: {v3.value.atp_generated:.1f}  Cache Hits: {v3.value.cache_hits}")
            print(f"    Metrics: {v3.value.metrics_collected}  Bridges: {v3.value.bridges_established}")
            print(f"  Veracity:")
            print(f"    Data Accuracy: {v3.veracity.data_accuracy:.2f}  Precision: {v3.veracity.metric_precision:.2f}")
            print(f"    False Positives: {v3.veracity.false_positive_rate:.2%}")
            print(f"  Velocity:")
            print(f"    TPS: {v3.velocity.transaction_rate:.1f}  Response: {v3.velocity.response_time_ms:.0f}ms")
            print(f"    Throughput: {v3.velocity.throughput:.1f}  Growth: {v3.velocity.growth_rate:+.2%}")

        print(f"{'='*60}")


def main():
    """Test trust tensor system"""
    import sys

    system = CBPTrustTensorSystem()

    if len(sys.argv) < 2:
        command = "help"
    else:
        command = sys.argv[1]

    if command == "init":
        # Initialize tensors for CBP roles
        roles = ["data_queen", "metrics_queen", "cache_queen", "coordinator"]

        for role in roles:
            t3, v3 = system.create_entity_tensors(f"cbp:{role}")
            print(f"✅ Created tensors for {role}")

        # Set initial values based on role specializations
        system.update_talent("cbp:data_queen", skill=0.8, accuracy=0.9, efficiency=0.7)
        system.update_talent("cbp:metrics_queen", skill=0.7, innovation=0.8, accuracy=0.85)
        system.update_talent("cbp:cache_queen", efficiency=0.9, accuracy=0.8, skill=0.75)

        system.save_state()
        print("\n📊 Federation Trust Summary:")
        print(json.dumps(system.get_federation_trust_summary(), indent=2))

    elif command == "simulate":
        # Simulate some operations
        system.create_entity_tensors("cbp:cache_queen")

        # Simulate cache operations
        for i in range(10):
            hit = i % 3 != 0  # 70% hit rate
            response = 50 if hit else 200
            system.update_value_from_cache_operation("cbp:cache_queen", hit, response)

        # Simulate training improvement
        system.update_training("cbp:cache_queen", task_completed=True, learning_delta=0.05)

        # Add fractal attribution
        system.add_fractal_attribution("cbp:cache_queen", "cbp:data_queen", 10.0)

        system.display_entity_trust("cbp:cache_queen")
        system.save_state()

    elif command == "show":
        if len(sys.argv) < 3:
            print("Usage: show <entity_id>")
        else:
            # Load state first
            if system.storage_path.exists():
                # Would implement load_state here
                pass
            system.display_entity_trust(sys.argv[2])

    else:
        print("CBP Trust Tensor System")
        print("Commands:")
        print("  init     - Initialize trust tensors for CBP roles")
        print("  simulate - Run simulation of tensor updates")
        print("  show     - Display entity tensors")


if __name__ == "__main__":
    main()