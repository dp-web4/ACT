#!/usr/bin/env python3
"""
CBP Trust Tensor Implementation V2 - Role-Contextual
CRITICAL FIX: Trust exists only within role contexts per Web4 spec
Based on Web4 specification sections 2.3, 5.1, 8.3 of t3-v3-tensors.md
"""

import json
from datetime import datetime
from typing import Dict, List, Optional, Tuple, Any
from dataclasses import dataclass, field, asdict
from pathlib import Path

# ============================================================================
# Role-Contextual Trust Implementation (Web4 Compliant)
# ============================================================================

@dataclass
class RoleContext:
    """Defines a specific role context for trust calculation"""
    role_id: str  # e.g., "web4:CacheManager"
    role_type: str  # queen, worker, coordinator
    domain: str  # cache, metrics, bridge, security, data
    authority_level: int  # 1-10
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class RoleContextualT3:
    """Trust tensor that exists only within a role context"""
    entity_id: str
    role_id: str
    talent: Dict[str, float] = field(default_factory=lambda: {
        "technical_skill": 0.5,
        "innovation": 0.5,
        "efficiency": 0.5,
        "accuracy": 0.5
    })
    training: Dict[str, Any] = field(default_factory=lambda: {
        "tasks_completed": 0,
        "success_rate": 0.0,
        "learning_rate": 0.5,
        "adaptation_score": 0.5,
        "specializations": []
    })
    temperament: Dict[str, float] = field(default_factory=lambda: {
        "consistency": 0.5,
        "cooperation": 0.5,
        "responsiveness": 0.5,
        "stability": 0.5
    })
    composite_trust: float = 0.5
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

    def calculate_composite_for_role(self, role_context: RoleContext) -> float:
        """Calculate trust score specific to role context"""
        # Role-specific weights based on domain
        if role_context.domain == "cache":
            weights = {"accuracy": 0.4, "efficiency": 0.3, "consistency": 0.3}
        elif role_context.domain == "metrics":
            weights = {"accuracy": 0.5, "technical_skill": 0.3, "consistency": 0.2}
        elif role_context.domain == "bridge":
            weights = {"cooperation": 0.4, "responsiveness": 0.3, "stability": 0.3}
        elif role_context.domain == "security":
            weights = {"accuracy": 0.3, "consistency": 0.4, "stability": 0.3}
        else:  # data
            weights = {"technical_skill": 0.3, "efficiency": 0.4, "accuracy": 0.3}

        # Calculate weighted score
        score = 0.0
        total_weight = 0.0

        for key, weight in weights.items():
            if key in self.talent:
                score += self.talent[key] * weight
                total_weight += weight
            elif key in self.temperament:
                score += self.temperament[key] * weight
                total_weight += weight

        # Add training contribution
        training_score = (self.training["success_rate"] * 0.5 +
                         self.training["learning_rate"] * 0.5)

        self.composite_trust = (score / max(total_weight, 0.1)) * 0.7 + training_score * 0.3
        return self.composite_trust

@dataclass
class RoleContextualV3:
    """Value tensor that exists only within a role context"""
    entity_id: str
    role_id: str
    value: Dict[str, Any] = field(default_factory=lambda: {
        "atp_generated": 0.0,
        "role_specific_value": 0.0,
        "operations_count": 0,
        "utility_score": 0.5
    })
    veracity: Dict[str, float] = field(default_factory=lambda: {
        "accuracy": 0.5,
        "precision": 0.5,
        "false_positive_rate": 0.0,
        "attestation_validity": 0.5
    })
    velocity: Dict[str, float] = field(default_factory=lambda: {
        "transaction_rate": 0.0,
        "response_time_ms": 1000.0,
        "throughput": 0.0,
        "growth_rate": 0.0
    })
    composite_value: float = 0.5
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

    def calculate_composite_for_role(self, role_context: RoleContext) -> float:
        """Calculate value score specific to role context"""
        # Role-specific value calculations
        if role_context.domain == "cache":
            # Cache values hit rate and response time
            specific_value = (self.value["operations_count"] / 1000.0) * 0.5
            speed_value = max(0, 1.0 - (self.velocity["response_time_ms"] / 1000.0)) * 0.5
        elif role_context.domain == "metrics":
            # Metrics values accuracy and throughput
            specific_value = self.veracity["accuracy"] * 0.6
            speed_value = min(1.0, self.velocity["throughput"] / 100.0) * 0.4
        else:
            specific_value = self.value["utility_score"]
            speed_value = min(1.0, self.velocity["transaction_rate"] / 100.0)

        self.composite_value = specific_value * 0.6 + speed_value * 0.4
        return self.composite_value

@dataclass
class EntityRolePair:
    """Binds an entity to a role with trust and value tensors"""
    entity_id: str
    role: RoleContext
    t3: RoleContextualT3
    v3: RoleContextualV3
    pairing_type: str = "operational"  # operational, permanent, temporary
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())

    def get_composite_trust(self) -> float:
        """Get role-specific trust score"""
        return self.t3.calculate_composite_for_role(self.role)

    def get_composite_value(self) -> float:
        """Get role-specific value score"""
        return self.v3.calculate_composite_for_role(self.role)

class CBPRoleContextualTensorSystem:
    """
    Web4-compliant trust system where trust only exists within role contexts
    Fixes critical violation identified in compliance review
    """

    def __init__(self):
        self.roles: Dict[str, RoleContext] = {}
        self.entity_role_pairs: Dict[str, EntityRolePair] = {}
        self.storage_path = Path("implementation/cbp-chain/role_contextual_tensors.json")
        self._initialize_cbp_roles()

    def _initialize_cbp_roles(self):
        """Initialize CBP-specific roles"""
        cbp_roles = [
            RoleContext("web4:cbp:DataQueen", "queen", "data", 8),
            RoleContext("web4:cbp:MetricsQueen", "queen", "metrics", 8),
            RoleContext("web4:cbp:SecurityQueen", "queen", "security", 9),
            RoleContext("web4:cbp:BridgeQueen", "queen", "bridge", 7),
            RoleContext("web4:cbp:CacheQueen", "queen", "cache", 7),
            RoleContext("web4:cbp:Worker", "worker", "general", 3),
            RoleContext("web4:cbp:Coordinator", "coordinator", "general", 6),
        ]

        for role in cbp_roles:
            self.roles[role.role_id] = role

    def create_entity_role_pairing(self, entity_id: str, role_id: str,
                                   pairing_type: str = "operational") -> EntityRolePair:
        """Create a new entity-role pairing with contextual tensors"""
        if role_id not in self.roles:
            raise ValueError(f"Unknown role: {role_id}")

        role = self.roles[role_id]

        # Create role-specific tensors
        t3 = RoleContextualT3(entity_id=entity_id, role_id=role_id)
        v3 = RoleContextualV3(entity_id=entity_id, role_id=role_id)

        # Initialize role-appropriate starting values
        if role.domain == "cache":
            t3.talent["efficiency"] = 0.7
            t3.talent["accuracy"] = 0.8
        elif role.domain == "metrics":
            t3.talent["accuracy"] = 0.85
            t3.talent["technical_skill"] = 0.75
        elif role.domain == "security":
            t3.temperament["consistency"] = 0.9
            t3.temperament["stability"] = 0.85

        pair = EntityRolePair(
            entity_id=entity_id,
            role=role,
            t3=t3,
            v3=v3,
            pairing_type=pairing_type
        )

        # Store with composite key
        pair_key = f"{entity_id}:{role_id}"
        self.entity_role_pairs[pair_key] = pair

        return pair

    def update_trust_in_role(self, entity_id: str, role_id: str,
                             dimension: str, updates: Dict[str, Any]) -> float:
        """Update trust tensor within specific role context"""
        pair_key = f"{entity_id}:{role_id}"

        if pair_key not in self.entity_role_pairs:
            raise ValueError(f"No pairing exists for {entity_id} in role {role_id}")

        pair = self.entity_role_pairs[pair_key]

        # Update appropriate dimension
        if dimension == "talent":
            for key, value in updates.items():
                if key in pair.t3.talent:
                    pair.t3.talent[key] = max(0, min(1, value))
        elif dimension == "training":
            for key, value in updates.items():
                if key in pair.t3.training:
                    if key == "tasks_completed":
                        pair.t3.training[key] += value
                    else:
                        pair.t3.training[key] = value
        elif dimension == "temperament":
            for key, value in updates.items():
                if key in pair.t3.temperament:
                    pair.t3.temperament[key] = max(0, min(1, value))

        pair.t3.last_updated = datetime.now().isoformat()
        return pair.get_composite_trust()

    def update_value_in_role(self, entity_id: str, role_id: str,
                             dimension: str, updates: Dict[str, Any]) -> float:
        """Update value tensor within specific role context"""
        pair_key = f"{entity_id}:{role_id}"

        if pair_key not in self.entity_role_pairs:
            raise ValueError(f"No pairing exists for {entity_id} in role {role_id}")

        pair = self.entity_role_pairs[pair_key]

        # Update appropriate dimension
        if dimension == "value":
            for key, value in updates.items():
                if key in pair.v3.value:
                    if key in ["atp_generated", "operations_count"]:
                        pair.v3.value[key] += value
                    else:
                        pair.v3.value[key] = value
        elif dimension == "veracity":
            for key, value in updates.items():
                if key in pair.v3.veracity:
                    pair.v3.veracity[key] = max(0, min(1, value))
        elif dimension == "velocity":
            for key, value in updates.items():
                if key in pair.v3.velocity:
                    pair.v3.velocity[key] = value

        pair.v3.last_updated = datetime.now().isoformat()
        return pair.get_composite_value()

    def get_entity_trust_across_roles(self, entity_id: str) -> Dict[str, float]:
        """Get trust scores for entity across all their roles"""
        trust_scores = {}

        for pair_key, pair in self.entity_role_pairs.items():
            if pair.entity_id == entity_id:
                trust_scores[pair.role.role_id] = pair.get_composite_trust()

        return trust_scores

    def get_role_trust_across_entities(self, role_id: str) -> Dict[str, float]:
        """Get trust scores for all entities in a specific role"""
        trust_scores = {}

        for pair_key, pair in self.entity_role_pairs.items():
            if pair.role.role_id == role_id:
                trust_scores[pair.entity_id] = pair.get_composite_trust()

        return trust_scores

    def apply_tensor_decay(self, days_elapsed: int = 30):
        """Apply time-based decay to tensors per Web4 spec"""
        decay_rate = 0.001 * days_elapsed  # -0.001/month for training
        recovery_rate = 0.01 * days_elapsed  # +0.01/month for temperament

        for pair in self.entity_role_pairs.values():
            # Decay training scores
            pair.t3.training["learning_rate"] = max(0,
                pair.t3.training["learning_rate"] - decay_rate)
            pair.t3.training["adaptation_score"] = max(0,
                pair.t3.training["adaptation_score"] - decay_rate)

            # Recover temperament scores
            for key in pair.t3.temperament:
                current = pair.t3.temperament[key]
                pair.t3.temperament[key] = min(1.0, current + recovery_rate * (0.5 - current))

    def save_state(self):
        """Save all tensor states to disk"""
        state = {
            "roles": {k: asdict(v) for k, v in self.roles.items()},
            "entity_role_pairs": {},
            "metadata": {
                "version": "2.0",
                "compliant": "Web4 sections 2.3, 5.1, 8.3",
                "created_at": datetime.now().isoformat()
            }
        }

        for key, pair in self.entity_role_pairs.items():
            state["entity_role_pairs"][key] = {
                "entity_id": pair.entity_id,
                "role_id": pair.role.role_id,
                "pairing_type": pair.pairing_type,
                "t3": asdict(pair.t3),
                "v3": asdict(pair.v3),
                "created_at": pair.created_at,
                "composite_trust": pair.get_composite_trust(),
                "composite_value": pair.get_composite_value()
            }

        self.storage_path.parent.mkdir(parents=True, exist_ok=True)
        with open(self.storage_path, 'w') as f:
            json.dump(state, f, indent=2)

    def display_compliance_status(self):
        """Display Web4 compliance status for trust tensors"""
        print("\n" + "="*60)
        print("🔐 CBP Role-Contextual Trust System (Web4 Compliant)")
        print("="*60)
        print(f"✅ Role Contexts Defined: {len(self.roles)}")
        print(f"✅ Entity-Role Pairings: {len(self.entity_role_pairs)}")
        print(f"✅ Trust exists ONLY within role contexts")
        print(f"✅ Compliant with Web4 spec sections 2.3, 5.1, 8.3")

        print("\n📊 Active Entity-Role Pairings:")
        for pair_key, pair in list(self.entity_role_pairs.items())[:5]:
            trust = pair.get_composite_trust()
            value = pair.get_composite_value()
            print(f"  {pair_key}: Trust={trust:.3f}, Value={value:.3f}")

        print("\n🎯 Compliance Features:")
        print("  ✅ Role-contextual trust calculation")
        print("  ✅ Domain-specific weight adjustments")
        print("  ✅ Tensor decay mechanisms implemented")
        print("  ✅ Separate trust scores per role")
        print("="*60)


def main():
    """Test role-contextual trust system"""
    system = CBPRoleContextualTensorSystem()

    # Create entity-role pairings
    entities = [
        ("alice", "web4:cbp:CacheQueen"),
        ("alice", "web4:cbp:MetricsQueen"),  # Same entity, different role!
        ("bob", "web4:cbp:DataQueen"),
        ("charlie", "web4:cbp:SecurityQueen"),
    ]

    for entity, role in entities:
        pair = system.create_entity_role_pairing(entity, role, "permanent")
        print(f"✅ Created pairing: {entity} as {role}")

    # Update trust in specific role context
    system.update_trust_in_role("alice", "web4:cbp:CacheQueen",
                                "talent", {"efficiency": 0.9, "accuracy": 0.95})

    # Alice has different trust in different roles!
    alice_trust = system.get_entity_trust_across_roles("alice")
    print(f"\n📊 Alice's trust across roles:")
    for role, trust in alice_trust.items():
        print(f"  {role}: {trust:.3f}")

    # Display compliance status
    system.display_compliance_status()
    system.save_state()


if __name__ == "__main__":
    main()