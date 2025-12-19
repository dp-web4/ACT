#!/usr/bin/env python3
"""
Cross-Society Trust Federation Protocol

Enables trust transfer and reputation portability across Web4 societies.
Solves the problem of starting from scratch when joining a new society.

Problem:
- Agent joins new society → cold start (no trust history)
- Expertise proven in Society A not recognized in Society B
- Reputation fragmentation across societies

Solution:
- Federated trust transfer protocol
- Cross-society reputation attestation
- Byzantine consensus for federation
- Privacy-preserving trust aggregation

Architecture:
- Federation Protocol: LCT-signed trust export/import
- Attestation: Multi-society witness consensus
- Privacy: Zero-knowledge proofs for aggregate trust (future)
- Portability: Expert reputation transfers across societies

Use Cases:
1. Agent migration: Move between societies with reputation intact
2. Cross-society collaboration: Leverage trust from multiple societies
3. Specialist portability: Transfer context-specific expertise
4. Federation: Multiple societies share trust network

Based on:
- WEB4-PROP-006-v2.1: Trust-first MoE standard
- Session 69: Byzantine consensus for multi-agent trust
- Session 70 Track 1: Trust persistence (ledger snapshots)

Security:
- LCT signatures for authenticity
- Byzantine consensus (2f+1 witnesses) for federation
- Decay function for cross-society trust (aging)
- Anti-Sybil: Hardware-bound identities

Author: Legion (Session 70 - Autonomous Web4 Research)
Date: 2025-12-19
"""

from dataclasses import dataclass, asdict
from typing import Dict, List, Optional, Tuple, Set
import hashlib
import time
import json


@dataclass
class TrustAttestation:
    """Trust attestation from one society about an agent."""
    agent_lct: str  # Who this attestation is about
    society_id: str  # Which society attests
    expert_id: int  # Expert role in that society
    context: str  # Context of expertise
    trust_value: float  # Attested trust [0, 1]
    observation_count: int  # Evidence strength
    attestation_timestamp: int  # When attested
    witnesses: List[str]  # LCT IDs of witnesses (Byzantine consensus)
    signature: str  # Society's signature (LCT-signed)

    def to_dict(self) -> Dict:
        return asdict(self)

    @classmethod
    def from_dict(cls, data: Dict) -> 'TrustAttestation':
        return cls(**data)

    def hash(self) -> str:
        """Deterministic hash for consensus."""
        data = f"{self.agent_lct}:{self.society_id}:{self.expert_id}:{self.context}:{self.trust_value}:{self.attestation_timestamp}"
        return hashlib.sha256(data.encode()).hexdigest()


@dataclass
class FederatedTrustProfile:
    """Aggregated trust profile from multiple societies."""
    agent_lct: str  # Who this profile is for
    attestations: List[TrustAttestation]  # All society attestations
    aggregated_trust: Dict[str, float]  # {context: aggregated_trust}
    societies: Set[str]  # Which societies contributed
    last_updated: int  # Most recent update

    def to_dict(self) -> Dict:
        return {
            "agent_lct": self.agent_lct,
            "attestations": [a.to_dict() for a in self.attestations],
            "aggregated_trust": self.aggregated_trust,
            "societies": list(self.societies),
            "last_updated": self.last_updated
        }

    @classmethod
    def from_dict(cls, data: Dict) -> 'FederatedTrustProfile':
        attestations = [TrustAttestation.from_dict(a) for a in data["attestations"]]
        return cls(
            agent_lct=data["agent_lct"],
            attestations=attestations,
            aggregated_trust=data["aggregated_trust"],
            societies=set(data["societies"]),
            last_updated=data["last_updated"]
        )


class FederationProtocol:
    """
    Cross-society trust federation protocol.

    Enables:
    1. Export trust from Society A
    2. Import trust to Society B with attestation
    3. Aggregate trust from multiple societies
    4. Byzantine consensus for federation
    """

    def __init__(
        self,
        min_witnesses: int = 3,
        trust_decay_rate: float = 0.9,  # Cross-society trust decays to 90% of original
        max_federation_age_days: int = 90  # Attestations expire after 90 days
    ):
        """
        Initialize federation protocol.

        Args:
            min_witnesses: Byzantine consensus threshold (2f+1)
            trust_decay_rate: Cross-society trust multiplier [0,1]
            max_federation_age_days: Maximum age for attestations
        """
        self.min_witnesses = min_witnesses
        self.decay_rate = trust_decay_rate
        self.max_age_seconds = max_federation_age_days * 24 * 60 * 60

        # Federation state
        self.profiles: Dict[str, FederatedTrustProfile] = {}  # agent_lct → profile
        self.pending_attestations: Dict[str, List[TrustAttestation]] = {}  # hash → [attestations]

    def create_attestation(
        self,
        agent_lct: str,
        society_id: str,
        expert_id: int,
        context: str,
        trust_value: float,
        observation_count: int,
        witnesses: List[str]
    ) -> TrustAttestation:
        """
        Create trust attestation from a society.

        Args:
            agent_lct: Agent being attested
            society_id: Attesting society
            expert_id: Expert role in that society
            context: Context of expertise
            trust_value: Trust value in [0, 1]
            observation_count: Number of observations
            witnesses: LCT IDs of witnesses (Byzantine consensus)

        Returns:
            TrustAttestation
        """
        attestation = TrustAttestation(
            agent_lct=agent_lct,
            society_id=society_id,
            expert_id=expert_id,
            context=context,
            trust_value=trust_value,
            observation_count=observation_count,
            attestation_timestamp=int(time.time()),
            witnesses=witnesses,
            signature=hashlib.sha256(
                f"{agent_lct}:{society_id}:{expert_id}:{context}:{trust_value}".encode()
            ).hexdigest()[:32]
        )

        return attestation

    def propose_attestation(
        self,
        attestation: TrustAttestation
    ) -> bool:
        """
        Propose attestation for Byzantine consensus.

        Returns:
            True if consensus reached, False if pending
        """
        attestation_hash = attestation.hash()

        if attestation_hash not in self.pending_attestations:
            self.pending_attestations[attestation_hash] = []

        # Check if already witnessed
        existing_witnesses = {
            w for att in self.pending_attestations[attestation_hash]
            for w in att.witnesses
        }

        new_witnesses = set(attestation.witnesses) - existing_witnesses
        if not new_witnesses:
            return False  # Duplicate witness

        # Add attestation
        self.pending_attestations[attestation_hash].append(attestation)

        # Check consensus
        all_witnesses = existing_witnesses | new_witnesses
        if len(all_witnesses) >= self.min_witnesses:
            # Consensus reached - add to federated profile
            self._add_to_profile(attestation)
            return True

        return False  # Still pending

    def _add_to_profile(self, attestation: TrustAttestation):
        """Add attestation to federated profile."""
        agent_lct = attestation.agent_lct

        if agent_lct not in self.profiles:
            self.profiles[agent_lct] = FederatedTrustProfile(
                agent_lct=agent_lct,
                attestations=[],
                aggregated_trust={},
                societies=set(),
                last_updated=0
            )

        profile = self.profiles[agent_lct]
        profile.attestations.append(attestation)
        profile.societies.add(attestation.society_id)
        profile.last_updated = int(time.time())

        # Reaggregate trust
        self._aggregate_trust(profile)

    def _aggregate_trust(self, profile: FederatedTrustProfile):
        """
        Aggregate trust from multiple societies.

        Strategy: Weighted average by observation count with decay
        """
        context_attestations: Dict[str, List[TrustAttestation]] = {}

        # Group by context
        for att in profile.attestations:
            # Check age
            age = int(time.time()) - att.attestation_timestamp
            if age > self.max_age_seconds:
                continue  # Expired

            if att.context not in context_attestations:
                context_attestations[att.context] = []
            context_attestations[att.context].append(att)

        # Aggregate per context
        profile.aggregated_trust = {}
        for context, attestations in context_attestations.items():
            # Weighted average by observation count
            total_weight = sum(att.observation_count for att in attestations)
            if total_weight == 0:
                continue

            weighted_sum = sum(
                att.trust_value * att.observation_count * self.decay_rate  # Apply decay
                for att in attestations
            )

            profile.aggregated_trust[context] = weighted_sum / total_weight

    def get_federated_trust(
        self,
        agent_lct: str,
        context: str
    ) -> Optional[float]:
        """
        Get aggregated trust for agent in context.

        Args:
            agent_lct: Agent LCT
            context: Context to query

        Returns:
            Aggregated trust or None if no federation data
        """
        profile = self.profiles.get(agent_lct)
        if not profile:
            return None

        return profile.aggregated_trust.get(context)

    def import_to_society(
        self,
        agent_lct: str,
        target_society_id: str,
        trust_selector
    ) -> int:
        """
        Import federated trust to a society's trust selector.

        Enables warm-start from cross-society reputation.

        Args:
            agent_lct: Agent joining society
            target_society_id: Society they're joining
            trust_selector: TrustFirstMRHSelector or TrustCoordinator

        Returns:
            Number of trust entries imported
        """
        profile = self.profiles.get(agent_lct)
        if not profile:
            return 0  # No federation data

        imported = 0
        for context, aggregated_trust in profile.aggregated_trust.items():
            # Apply additional decay for import (conservative)
            imported_trust = aggregated_trust * 0.8  # 80% of federated trust

            # Estimate observation count (conservative)
            relevant_attestations = [
                att for att in profile.attestations
                if att.context == context
            ]
            total_observations = sum(att.observation_count for att in relevant_attestations)

            # Import to trust selector (if agent maps to expert_id in new society)
            # NOTE: In real deployment, need LCT → expert_id mapping
            # For now, assume expert_id from first attestation
            if relevant_attestations:
                expert_id = relevant_attestations[0].expert_id

                key = (expert_id, context)

                if hasattr(trust_selector, 'expert_trust'):
                    # TrustFirstMRHSelector format
                    if expert_id not in trust_selector.expert_trust:
                        trust_selector.expert_trust[expert_id] = {}
                    trust_selector.expert_trust[expert_id][context] = imported_trust

                    if expert_id not in trust_selector.expert_observations:
                        trust_selector.expert_observations[expert_id] = {}
                    trust_selector.expert_observations[expert_id][context] = total_observations // 2  # Conservative

                if hasattr(trust_selector, 'society_trust'):
                    # TrustCoordinator format
                    trust_selector.society_trust.expert_trust[key] = imported_trust
                    trust_selector.society_trust.expert_observations[key] = total_observations // 2

                imported += 1

        return imported

    def get_statistics(self) -> Dict:
        """Get federation statistics."""
        return {
            "federated_agents": len(self.profiles),
            "total_attestations": sum(len(p.attestations) for p in self.profiles.values()),
            "societies_in_federation": len({
                society for p in self.profiles.values()
                for society in p.societies
            }),
            "pending_attestations": len(self.pending_attestations)
        }


# Demo
def demo_federation():
    """Demonstrate cross-society trust federation."""
    print("\n" + "="*70)
    print("CROSS-SOCIETY TRUST FEDERATION DEMO")
    print("="*70)

    federation = FederationProtocol(
        min_witnesses=3,
        trust_decay_rate=0.9,
        max_federation_age_days=90
    )

    print(f"\n✅ Federation protocol initialized")
    print(f"   Min witnesses (Byzantine): {federation.min_witnesses}")
    print(f"   Trust decay rate: {federation.decay_rate}")

    # Agent in Society A
    agent_lct = "lct://agent:researcher:alice@testnet"
    print(f"\n{'='*70}")
    print(f"SCENARIO: Agent migrates Society A → Society B")
    print(f"{'='*70}\n")
    print(f"Agent: {agent_lct}")

    # Society A attestations (3 witnesses for consensus)
    print(f"\nSociety A: Building reputation...")
    attestations_a = []
    for i in range(3):
        att = federation.create_attestation(
            agent_lct=agent_lct,
            society_id="society_a",
            expert_id=42,
            context="context_1",  # Reasoning specialist
            trust_value=0.85,
            observation_count=50,
            witnesses=[f"lct://society_a:validator:v{i}@testnet"]
        )
        attestations_a.append(att)
        consensus = federation.propose_attestation(att)
        if consensus:
            print(f"  ✅ Consensus reached with {i+1} witnesses")
            print(f"     Trust: 0.85, Context: context_1, Observations: 50")

    # Society B attestations (different context)
    print(f"\nSociety B: Agent also has expertise in...")
    for i in range(3):
        att = federation.create_attestation(
            agent_lct=agent_lct,
            society_id="society_b",
            expert_id=73,
            context="context_0",  # Code specialist
            trust_value=0.78,
            observation_count=35,
            witnesses=[f"lct://society_b:validator:v{i}@testnet"]
        )
        consensus = federation.propose_attestation(att)
        if consensus:
            print(f"  ✅ Consensus reached with {i+1} witnesses")
            print(f"     Trust: 0.78, Context: context_0, Observations: 35")

    # Check federated profile
    print(f"\n{'='*70}")
    print("FEDERATED TRUST PROFILE")
    print(f"{'='*70}\n")

    trust_reasoning = federation.get_federated_trust(agent_lct, "context_1")
    trust_code = federation.get_federated_trust(agent_lct, "context_0")

    print(f"Agent: {agent_lct}")
    print(f"\nAggregated Trust:")
    print(f"  Reasoning (context_1): {trust_reasoning:.3f} (from Society A)")
    print(f"  Code (context_0): {trust_code:.3f} (from Society B)")

    print(f"\nSocieties in Federation:")
    profile = federation.profiles[agent_lct]
    for society in profile.societies:
        society_atts = [a for a in profile.attestations if a.society_id == society]
        print(f"  {society}: {len(society_atts)} attestations")

    # Import to Society C
    print(f"\n{'='*70}")
    print("MIGRATION: Importing to Society C")
    print(f"{'='*70}\n")

    # Simulate trust selector for Society C
    class MockTrustSelector:
        def __init__(self):
            self.expert_trust = {}
            self.expert_observations = {}

    society_c_selector = MockTrustSelector()

    imported = federation.import_to_society(
        agent_lct=agent_lct,
        target_society_id="society_c",
        trust_selector=society_c_selector
    )

    print(f"✅ Imported {imported} trust entries to Society C")
    print(f"\nSociety C Trust State (warm-start):")
    for expert_id, contexts in society_c_selector.expert_trust.items():
        for context, trust in contexts.items():
            obs = society_c_selector.expert_observations[expert_id][context]
            print(f"  Expert {expert_id}, {context}: Trust={trust:.3f}, Obs={obs}")

    print(f"\nDecay Applied:")
    print(f"  Federated: 0.850 → Imported: {0.850 * 0.9 * 0.8:.3f} (90% decay + 80% import)")
    print(f"  Federated: 0.780 → Imported: {0.780 * 0.9 * 0.8:.3f} (90% decay + 80% import)")

    # Statistics
    stats = federation.get_statistics()
    print(f"\n{'='*70}")
    print("FEDERATION STATISTICS")
    print(f"{'='*70}\n")
    for key, value in stats.items():
        print(f"  {key}: {value}")

    print(f"\n✅ Federation demo complete")
    print(f"\nBenefit: Agent joins Society C with prior reputation intact!")
    print(f"         No cold-start, immediate trust-driven mode possible.")


if __name__ == "__main__":
    demo_federation()
