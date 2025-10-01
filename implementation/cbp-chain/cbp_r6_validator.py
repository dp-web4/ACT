#!/usr/bin/env python3
"""
CBP R6 Action Framework Validator
Implements basic R6 action validation per Web4 specification
Based on Rules-Role-Resource-Reality-Relevance-Relationship framework
"""

import json
from datetime import datetime
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass, field
from pathlib import Path

@dataclass
class R6Action:
    """R6 action structure following Web4 specification"""
    action_id: str
    actor: str  # Entity performing action
    rules: List[str]  # Law IDs that govern this action
    role: str  # Role context for the action
    resource: str  # Resource being acted upon
    reality: str  # Reality context (federation, local, cross-chain)
    relevance: float  # 0-1 relevance score
    relationships: List[str]  # Affected entities/relationships
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    status: str = "pending"  # pending, validated, rejected

class CBP_R6_Validator:
    """
    Basic R6 validation framework for CBP Society
    Validates actions against governance rules and role capabilities
    """

    def __init__(self):
        # Load governance for rule validation
        self.governance_path = Path("implementation/cbp-chain/cbp_governance.json")
        self.governance_data = self._load_governance()

    def _load_governance(self) -> Dict:
        """Load governance rules for validation"""
        if self.governance_path.exists():
            with open(self.governance_path, 'r') as f:
                return json.load(f)
        return {"laws": {}, "authorities": {}}

    def validate_r6_action(self, action: R6Action) -> Tuple[bool, str, List[str]]:
        """
        Validate an R6 action against all 6 dimensions
        Returns: (is_valid, reason, violations)
        """
        violations = []

        # 1. Rules validation
        rules_valid, rules_msg = self._validate_rules(action)
        if not rules_valid:
            violations.append(f"Rules: {rules_msg}")

        # 2. Role validation
        role_valid, role_msg = self._validate_role(action)
        if not role_valid:
            violations.append(f"Role: {role_msg}")

        # 3. Resource validation
        resource_valid, resource_msg = self._validate_resource(action)
        if not resource_valid:
            violations.append(f"Resource: {resource_msg}")

        # 4. Reality validation
        reality_valid, reality_msg = self._validate_reality(action)
        if not reality_valid:
            violations.append(f"Reality: {reality_msg}")

        # 5. Relevance validation
        relevance_valid, relevance_msg = self._validate_relevance(action)
        if not relevance_valid:
            violations.append(f"Relevance: {relevance_msg}")

        # 6. Relationships validation
        relationships_valid, relationships_msg = self._validate_relationships(action)
        if not relationships_valid:
            violations.append(f"Relationships: {relationships_msg}")

        is_valid = len(violations) == 0
        summary = "Action validated" if is_valid else f"{len(violations)} violations found"

        return is_valid, summary, violations

    def _validate_rules(self, action: R6Action) -> Tuple[bool, str]:
        """Validate action against governance rules"""
        if not action.rules:
            return False, "No governing rules specified"

        laws = self.governance_data.get("laws", {})

        for rule_id in action.rules:
            if rule_id not in laws:
                return False, f"Unknown rule: {rule_id}"

            law = laws[rule_id]
            if law.get("status") != "active":
                return False, f"Inactive rule: {rule_id}"

            # Basic R6 selector validation
            if "r6." not in law.get("r6_selector", ""):
                return False, f"Rule {rule_id} missing R6 selector"

        return True, f"All {len(action.rules)} rules validated"

    def _validate_role(self, action: R6Action) -> Tuple[bool, str]:
        """Validate actor has required role and capabilities"""
        authorities = self.governance_data.get("authorities", {})

        if action.actor not in authorities:
            return False, f"Actor {action.actor} has no defined authorities"

        # Extract action type from rules
        action_type = self._extract_action_type(action)
        actor_capabilities = authorities[action.actor]

        # Check if actor has capability for this action type
        required_capability = self._map_action_to_capability(action_type)
        if required_capability and required_capability not in actor_capabilities:
            return False, f"Actor lacks required capability: {required_capability}"

        return True, f"Role {action.role} validated for {action.actor}"

    def _validate_resource(self, action: R6Action) -> Tuple[bool, str]:
        """Validate resource access and constraints"""
        # Basic resource validation
        if not action.resource:
            return False, "No resource specified"

        # Check resource format
        if not any(prefix in action.resource for prefix in ["atp:", "lct:", "cache:", "metric:", "law:"]):
            return False, f"Invalid resource format: {action.resource}"

        return True, f"Resource {action.resource} validated"

    def _validate_reality(self, action: R6Action) -> Tuple[bool, str]:
        """Validate reality context"""
        valid_realities = ["federation", "local", "cross-chain", "simulation"]

        if action.reality not in valid_realities:
            return False, f"Invalid reality context: {action.reality}"

        return True, f"Reality context {action.reality} validated"

    def _validate_relevance(self, action: R6Action) -> Tuple[bool, str]:
        """Validate relevance score"""
        if not 0 <= action.relevance <= 1:
            return False, f"Relevance must be 0-1, got {action.relevance}"

        # Low relevance actions should be questioned
        if action.relevance < 0.3:
            return False, f"Relevance too low: {action.relevance}"

        return True, f"Relevance {action.relevance:.2f} validated"

    def _validate_relationships(self, action: R6Action) -> Tuple[bool, str]:
        """Validate relationship impacts"""
        if not action.relationships:
            return True, "No relationships affected"

        # Basic relationship format validation
        for rel in action.relationships:
            if not any(prefix in rel for prefix in ["lct:", "role:", "society:", "federation:"]):
                return False, f"Invalid relationship format: {rel}"

        return True, f"{len(action.relationships)} relationships validated"

    def _extract_action_type(self, action: R6Action) -> str:
        """Extract action type from rules"""
        for rule_id in action.rules:
            if "GOV" in rule_id:
                return "governance"
            elif "ECON" in rule_id:
                return "economic"
            elif "SEC" in rule_id:
                return "security"
            elif "FED" in rule_id:
                return "federation"
            elif "OPER" in rule_id:
                return "operational"
        return "general"

    def _map_action_to_capability(self, action_type: str) -> Optional[str]:
        """Map action type to required capability"""
        mapping = {
            "governance": "propose_amendment",
            "economic": "allocate_atp",
            "security": "verify_hardware",
            "federation": "manage_federation",
            "operational": "manage_cache"
        }
        return mapping.get(action_type)

    def create_r6_action(self, actor: str, action_type: str, target: str,
                        reality: str = "federation") -> R6Action:
        """Helper to create R6 action with proper structure"""
        action_id = f"r6-{action_type}-{int(datetime.now().timestamp())}"

        # Map action type to relevant rules
        rule_mapping = {
            "atp_discharge": ["LAW-ECON-001", "LAW-ECON-002"],
            "cache_operation": ["LAW-OPER-001"],
            "governance_vote": ["LAW-GOV-001", "LAW-GOV-003"],
            "federation_bridge": ["LAW-FED-001"],
            "security_check": ["LAW-SEC-001"]
        }

        rules = rule_mapping.get(action_type, ["LAW-GOV-001"])
        role = actor.split(":")[-1] if ":" in actor else "unknown"

        # Calculate basic relevance
        relevance = 0.8 if action_type in rule_mapping else 0.5

        action = R6Action(
            action_id=action_id,
            actor=actor,
            rules=rules,
            role=role,
            resource=target,
            reality=reality,
            relevance=relevance,
            relationships=[target] if target else []
        )

        return action

def main():
    """Test R6 validation framework"""
    validator = CBP_R6_Validator()

    # Test valid action
    action = validator.create_r6_action(
        actor="cbp:cache_queen",
        action_type="cache_operation",
        target="cache:federation_metrics",
        reality="federation"
    )

    print("🔍 Testing R6 Action Validation")
    print("=" * 50)
    print(f"Action ID: {action.action_id}")
    print(f"Actor: {action.actor}")
    print(f"Rules: {action.rules}")
    print(f"Resource: {action.resource}")

    valid, message, violations = validator.validate_r6_action(action)

    if valid:
        print(f"\n✅ {message}")
    else:
        print(f"\n❌ {message}")
        for violation in violations:
            print(f"  - {violation}")

    # Test invalid action
    invalid_action = R6Action(
        action_id="invalid-test",
        actor="unknown_actor",
        rules=["NONEXISTENT_LAW"],
        role="invalid",
        resource="bad_format",
        reality="metaverse",
        relevance=1.5,  # Invalid range
        relationships=["bad_format"]
    )

    print(f"\n🔍 Testing Invalid Action")
    print("=" * 30)
    valid, message, violations = validator.validate_r6_action(invalid_action)

    print(f"❌ {message}")
    for violation in violations:
        print(f"  - {violation}")

if __name__ == "__main__":
    main()