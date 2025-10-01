#!/usr/bin/env python3
"""
CBP Governance System - Law Oracle Implementation
Web4-compliant minimal governance for CBP Society
Implements foundational laws with amendment mechanism
"""

import json
import hashlib
from datetime import datetime
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass, field, asdict
from pathlib import Path

# ============================================================================
# Law Oracle Structure (Web4 SAL Compliant)
# ============================================================================

@dataclass
class Law:
    """Individual law with versioning and metadata"""
    law_id: str
    category: str  # governance, economic, security, federation, operational
    title: str
    content: str
    rationale: str
    r6_selector: str  # R6 action grammar selector
    enforcement: str  # automatic, manual, voted
    penalties: Dict[str, Any] = field(default_factory=dict)
    version: int = 1
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    amended_at: Optional[str] = None
    author: str = "cbp:coordinator"
    status: str = "active"  # active, suspended, deprecated

@dataclass
class Amendment:
    """Proposed or enacted amendment to laws"""
    amendment_id: str
    law_id: str
    proposer: str
    changes: Dict[str, Any]
    rationale: str
    votes_for: List[str] = field(default_factory=list)
    votes_against: List[str] = field(default_factory=list)
    status: str = "proposed"  # proposed, approved, rejected
    proposed_at: str = field(default_factory=lambda: datetime.now().isoformat())
    decided_at: Optional[str] = None

@dataclass
class GovernanceEvent:
    """Immutable record of governance actions"""
    event_id: str
    event_type: str  # law_created, law_amended, vote_cast, decision_made
    actor: str
    target: str
    data: Dict[str, Any]
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    block_height: Optional[int] = None
    witnesses: List[str] = field(default_factory=list)

class CBPLawOracle:
    """
    CBP's Law Oracle - Machine-readable governance
    Implements Web4 SAL (Society-Authority-Law) specification
    """

    def __init__(self):
        self.laws: Dict[str, Law] = {}
        self.amendments: Dict[str, Amendment] = {}
        self.events: List[GovernanceEvent] = []
        self.authorities: Dict[str, List[str]] = {}  # role -> capabilities
        self.storage_path = Path("implementation/cbp-chain/cbp_governance.json")
        self._initialize_foundational_laws()
        self._initialize_authorities()

    def _initialize_foundational_laws(self):
        """Create CBP's foundational laws"""

        foundational_laws = [
            Law(
                law_id="LAW-GOV-001",
                category="governance",
                title="Queen Consensus Requirement",
                content="Decisions affecting society operations require approval from 3/5 queens",
                rationale="Ensures distributed decision-making while maintaining efficiency",
                r6_selector="r6.decision.consensus",
                enforcement="automatic",
                penalties={"violation": "decision_void", "atp_cost": 50}
            ),
            Law(
                law_id="LAW-GOV-002",
                category="governance",
                title="Security Queen Veto",
                content="Security Queen may veto any decision with security implications",
                rationale="Critical security decisions need specialized oversight",
                r6_selector="r6.security.veto",
                enforcement="automatic",
                penalties={"override": "security_breach_alert", "atp_cost": 100}
            ),
            Law(
                law_id="LAW-GOV-003",
                category="governance",
                title="Amendment Requirement",
                content="Constitutional amendments require 4/5 queens approval",
                rationale="Foundational changes need super-majority consensus",
                r6_selector="r6.governance.amendment",
                enforcement="automatic",
                penalties={"invalid_amendment": "rejection"}
            ),
            Law(
                law_id="LAW-ECON-001",
                category="economic",
                title="ATP Budget Constraint",
                content="Total society ATP allocation must equal exactly 1000",
                rationale="Web4 energy conservation principle",
                r6_selector="r6.economy.atp_total",
                enforcement="automatic",
                penalties={"violation": "automatic_rebalancing"}
            ),
            Law(
                law_id="LAW-ECON-002",
                category="economic",
                title="Cache Operations Priority",
                content="Cache Queen receives priority ATP allocation (minimum 110)",
                rationale="Cache operations are core to CBP's value proposition",
                r6_selector="r6.economy.role_allocation",
                enforcement="automatic",
                penalties={"underallocation": "performance_degradation"}
            ),
            Law(
                law_id="LAW-ECON-003",
                category="economic",
                title="Daily Recharge Rate",
                content="All roles receive 20 ATP daily recharge from ADP conversion",
                rationale="Sustainable energy regeneration for continuous operation",
                r6_selector="r6.economy.recharge",
                enforcement="automatic",
                penalties={"missed_recharge": "energy_debt"}
            ),
            Law(
                law_id="LAW-SEC-001",
                category="security",
                title="Hardware Binding Verification",
                content="All entities must maintain valid hardware binding proofs",
                rationale="Prevents Sybil attacks and ensures entity authenticity",
                r6_selector="r6.security.hardware_check",
                enforcement="automatic",
                penalties={"invalid_binding": "entity_suspension", "atp_cost": 200}
            ),
            Law(
                law_id="LAW-FED-001",
                category="federation",
                title="Federation Protocol Compliance",
                content="All federation communications must use Web4 protocol",
                rationale="Ensures interoperability with other societies",
                r6_selector="r6.federation.protocol",
                enforcement="automatic",
                penalties={"protocol_violation": "communication_blocked"}
            ),
            Law(
                law_id="LAW-OPER-001",
                category="operational",
                title="Cache Hit Rate Minimum",
                content="Cache system must maintain minimum 70% hit rate",
                rationale="Ensures cache effectiveness and value delivery",
                r6_selector="r6.operation.cache_performance",
                enforcement="manual",
                penalties={"underperformance": "optimization_required", "atp_reduction": 10}
            ),
            Law(
                law_id="LAW-OPER-002",
                category="operational",
                title="Metrics Collection Frequency",
                content="Trust tensors must be updated at least every 24 hours",
                rationale="Maintains accurate federation trust state",
                r6_selector="r6.operation.metrics_update",
                enforcement="automatic",
                penalties={"stale_metrics": "trust_decay", "accuracy_penalty": 0.01}
            )
        ]

        for law in foundational_laws:
            self.laws[law.law_id] = law
            self._record_event("law_created", "cbp:coordinator", law.law_id, {"law": asdict(law)})

    def _initialize_authorities(self):
        """Define authority capabilities for each role"""
        self.authorities = {
            "cbp:coordinator": [
                "propose_amendment",
                "call_vote",
                "execute_decision",
                "allocate_resources"
            ],
            "cbp:security_queen": [
                "veto_security",
                "audit_entities",
                "verify_hardware",
                "enforce_penalties"
            ],
            "cbp:data_queen": [
                "manage_storage",
                "optimize_cache",
                "allocate_data_atp",
                "propose_data_laws",
                "propose_amendment"
            ],
            "cbp:metrics_queen": [
                "update_tensors",
                "calculate_trust",
                "report_metrics",
                "propose_metric_laws"
            ],
            "cbp:bridge_queen": [
                "manage_federation",
                "translate_protocols",
                "propose_federation_laws"
            ],
            "cbp:cache_queen": [
                "manage_cache",
                "optimize_performance",
                "allocate_cache_atp",
                "report_hit_rates"
            ]
        }

    def propose_amendment(self, proposer: str, law_id: str, changes: Dict[str, Any], rationale: str) -> str:
        """Propose an amendment to an existing law"""
        if law_id not in self.laws:
            raise ValueError(f"Law {law_id} does not exist")

        if proposer not in self.authorities:
            raise ValueError(f"Entity {proposer} lacks authority")

        if "propose_amendment" not in self.authorities.get(proposer, []):
            raise ValueError(f"Entity {proposer} cannot propose amendments")

        amendment_id = f"AMD-{len(self.amendments):04d}"
        amendment = Amendment(
            amendment_id=amendment_id,
            law_id=law_id,
            proposer=proposer,
            changes=changes,
            rationale=rationale
        )

        self.amendments[amendment_id] = amendment
        self._record_event("amendment_proposed", proposer, amendment_id, {"amendment": asdict(amendment)})

        return amendment_id

    def vote_on_amendment(self, voter: str, amendment_id: str, vote: bool) -> Tuple[bool, str]:
        """Cast vote on an amendment"""
        if amendment_id not in self.amendments:
            raise ValueError(f"Amendment {amendment_id} does not exist")

        amendment = self.amendments[amendment_id]
        if amendment.status != "proposed":
            return False, f"Amendment already {amendment.status}"

        # Only queens can vote
        if "queen" not in voter.lower() and "coordinator" not in voter.lower():
            return False, "Only queens and coordinator can vote"

        # Remove any existing vote from this voter
        if voter in amendment.votes_for:
            amendment.votes_for.remove(voter)
        if voter in amendment.votes_against:
            amendment.votes_against.remove(voter)

        # Cast new vote
        if vote:
            amendment.votes_for.append(voter)
        else:
            amendment.votes_against.append(voter)

        self._record_event("vote_cast", voter, amendment_id, {"vote": vote})

        # Check if decision can be made
        return self._check_amendment_decision(amendment_id)

    def _check_amendment_decision(self, amendment_id: str) -> Tuple[bool, str]:
        """Check if amendment has enough votes to pass or fail"""
        amendment = self.amendments[amendment_id]
        law = self.laws[amendment.law_id]

        total_queens = 5  # CBP has 5 queens

        # Determine required threshold
        if law.category == "governance" and "constitutional" in law.title.lower():
            required = 4  # 4/5 for constitutional
        else:
            required = 3  # 3/5 for regular

        if len(amendment.votes_for) >= required:
            # Amendment passes
            amendment.status = "approved"
            amendment.decided_at = datetime.now().isoformat()
            self._apply_amendment(amendment_id)
            self._record_event("amendment_approved", "governance_system", amendment_id,
                             {"votes_for": len(amendment.votes_for), "votes_against": len(amendment.votes_against)})
            return True, "Amendment approved and applied"

        if len(amendment.votes_against) > (total_queens - required):
            # Amendment fails
            amendment.status = "rejected"
            amendment.decided_at = datetime.now().isoformat()
            self._record_event("amendment_rejected", "governance_system", amendment_id,
                             {"votes_for": len(amendment.votes_for), "votes_against": len(amendment.votes_against)})
            return True, "Amendment rejected"

        return False, f"Voting continues ({len(amendment.votes_for)} for, {len(amendment.votes_against)} against)"

    def _apply_amendment(self, amendment_id: str):
        """Apply approved amendment to law"""
        amendment = self.amendments[amendment_id]
        law = self.laws[amendment.law_id]

        # Apply changes
        for key, value in amendment.changes.items():
            if hasattr(law, key):
                setattr(law, key, value)

        law.version += 1
        law.amended_at = datetime.now().isoformat()

        self._record_event("law_amended", "governance_system", law.law_id,
                         {"amendment_id": amendment_id, "new_version": law.version})

    def enforce_law(self, law_id: str, violator: str, context: Dict[str, Any]) -> Dict[str, Any]:
        """Enforce a law violation"""
        if law_id not in self.laws:
            raise ValueError(f"Law {law_id} does not exist")

        law = self.laws[law_id]

        enforcement = {
            "law_id": law_id,
            "violator": violator,
            "timestamp": datetime.now().isoformat(),
            "penalties_applied": law.penalties,
            "context": context
        }

        self._record_event("law_enforced", "governance_system", law_id,
                         {"violator": violator, "penalties": law.penalties})

        return enforcement

    def query_laws(self, category: Optional[str] = None, r6_selector: Optional[str] = None) -> List[Law]:
        """Query laws by category or R6 selector"""
        results = []

        for law in self.laws.values():
            if law.status != "active":
                continue

            if category and law.category != category:
                continue

            if r6_selector and not law.r6_selector.startswith(r6_selector):
                continue

            results.append(law)

        return results

    def _record_event(self, event_type: str, actor: str, target: str, data: Dict[str, Any]):
        """Record governance event for immutable ledger"""
        event_id = hashlib.sha256(f"{event_type}{actor}{target}{datetime.now().isoformat()}".encode()).hexdigest()[:16]

        event = GovernanceEvent(
            event_id=event_id,
            event_type=event_type,
            actor=actor,
            target=target,
            data=data
        )

        self.events.append(event)

    def get_governance_summary(self) -> Dict[str, Any]:
        """Get summary of governance state"""
        return {
            "total_laws": len(self.laws),
            "active_laws": len([l for l in self.laws.values() if l.status == "active"]),
            "categories": list(set(l.category for l in self.laws.values())),
            "proposed_amendments": len([a for a in self.amendments.values() if a.status == "proposed"]),
            "total_amendments": len(self.amendments),
            "governance_events": len(self.events),
            "last_event": self.events[-1].timestamp if self.events else None
        }

    def save_state(self):
        """Save governance state to disk"""
        state = {
            "laws": {k: asdict(v) for k, v in self.laws.items()},
            "amendments": {k: asdict(v) for k, v in self.amendments.items()},
            "authorities": self.authorities,
            "events": [asdict(e) for e in self.events[-100:]],  # Keep last 100 events
            "metadata": {
                "version": "1.0.0",
                "compliant": "Web4 SAL specification",
                "created_at": datetime.now().isoformat()
            }
        }

        self.storage_path.parent.mkdir(parents=True, exist_ok=True)
        with open(self.storage_path, 'w') as f:
            json.dump(state, f, indent=2)

    def display_laws(self):
        """Display current laws in readable format"""
        print("\n" + "="*60)
        print("📜 CBP Law Oracle - Governance System")
        print("="*60)

        for category in ["governance", "economic", "security", "federation", "operational"]:
            laws = self.query_laws(category=category)
            if laws:
                print(f"\n{category.upper()} LAWS:")
                for law in laws:
                    print(f"  [{law.law_id}] {law.title}")
                    print(f"    {law.content}")
                    print(f"    Enforcement: {law.enforcement}")
                    if law.penalties:
                        print(f"    Penalties: {law.penalties}")

        summary = self.get_governance_summary()
        print(f"\n📊 Governance Summary:")
        print(f"  Active Laws: {summary['active_laws']}/{summary['total_laws']}")
        print(f"  Amendments: {summary['proposed_amendments']} pending, {summary['total_amendments']} total")
        print(f"  Events Recorded: {summary['governance_events']}")
        print("="*60)


def main():
    """Test CBP governance system"""
    oracle = CBPLawOracle()

    # Display initial laws
    oracle.display_laws()

    # Test amendment process
    print("\n🗳️ Testing Amendment Process:")

    # Propose an amendment
    amendment_id = oracle.propose_amendment(
        proposer="cbp:data_queen",
        law_id="LAW-OPER-001",
        changes={"content": "Cache system must maintain minimum 80% hit rate"},
        rationale="Higher standards needed as cache grows"
    )
    print(f"✅ Amendment {amendment_id} proposed")

    # Queens vote
    queens = ["cbp:data_queen", "cbp:cache_queen", "cbp:metrics_queen", "cbp:security_queen"]
    for i, queen in enumerate(queens[:3]):  # First 3 vote yes
        result, message = oracle.vote_on_amendment(queen, amendment_id, True)
        print(f"  {queen} votes YES - {message}")

    # Save state
    oracle.save_state()
    print("\n✅ Governance system initialized and saved")


if __name__ == "__main__":
    main()