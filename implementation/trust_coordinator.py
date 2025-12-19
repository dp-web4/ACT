#!/usr/bin/env python3
"""
Trust-First Multi-Agent Coordinator for ACT Societies

Implements distributed trust-based expert selection for AI agent societies.
Agents share trust observations to collectively improve expert routing.

Architecture:
- Each agent maintains local trust scores (expert × context)
- Agents broadcast trust updates to society
- Society aggregates trust via Byzantine consensus (2f+1 witnesses)
- Trust-first selection uses aggregated society trust

Based on:
- Sessions 64-69 (Legion): Trust-first MoE architecture
- Sessions 69-73 (Thor): Long-term evolution and specialist emergence
- WEB4-PROP-006-v2.1: LCT-MoE Trust Standard

Integration with ACT:
- LCT identity for agent authentication
- Society ledger for trust update persistence
- ATP/ADP for resource allocation based on trust
- Byzantine consensus for Sybil resistance

Author: Legion (Session 69 - Autonomous Web4 Research)
Date: 2025-12-19
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Set, Tuple
import json
import hashlib
import time
import numpy as np


@dataclass
class LCTIdentity:
    """Linked Context Token identity for agent."""
    lct_id: str  # Format: lct://society:role:agent_id@network
    public_key: str  # For signature verification
    hardware_hash: str  # Hardware binding (TPM in production)
    role: str  # "citizen", "validator", "coordinator"


@dataclass
class TrustUpdate:
    """Trust update broadcast by agent to society."""
    agent_lct: str  # Who observed
    expert_id: int  # Which expert
    context: str  # What context
    quality: float  # Observed quality [0, 1]
    timestamp: int  # Unix timestamp
    signature: str  # Cryptographic signature (placeholder)

    def to_dict(self) -> Dict:
        return {
            "agent_lct": self.agent_lct,
            "expert_id": self.expert_id,
            "context": self.context,
            "quality": self.quality,
            "timestamp": self.timestamp,
            "signature": self.signature
        }

    def consensus_key(self) -> str:
        """Key for grouping consensus (ignores agent_lct, timestamp, quality variance)."""
        return f"expert_{self.expert_id}:context_{self.context}"

    def hash(self) -> str:
        """Deterministic hash for consensus."""
        data = f"{self.agent_lct}:{self.expert_id}:{self.context}:{self.quality}:{self.timestamp}"
        return hashlib.sha256(data.encode()).hexdigest()


@dataclass
class SocietyTrustState:
    """Aggregated trust state for entire society."""
    expert_trust: Dict[Tuple[int, str], float]  # (expert, context) → trust
    expert_observations: Dict[Tuple[int, str], int]  # (expert, context) → count
    agent_contributions: Dict[str, int]  # agent_lct → update_count
    last_update: int  # Timestamp


class ByzantineConsensus:
    """
    Byzantine fault-tolerant consensus for trust updates.

    Requirement: 2f+1 witnesses where f is max Byzantine failures.
    For society of n agents, tolerates f = (n-1)/3 failures.

    Example: 10 agents → tolerates 3 Byzantine agents
    """

    def __init__(self, min_witnesses: int = 3):
        """
        Initialize Byzantine consensus.

        Args:
            min_witnesses: Minimum witness count (2f+1 where f is failures)
        """
        self.min_witnesses = min_witnesses
        self.pending_updates: Dict[str, List[TrustUpdate]] = {}  # consensus_key → [updates]

    def propose_update(self, update: TrustUpdate) -> bool:
        """
        Propose trust update for consensus.

        Returns:
            True if consensus reached, False if pending
        """
        consensus_key = update.consensus_key()

        if consensus_key not in self.pending_updates:
            self.pending_updates[consensus_key] = []

        # Check if this agent already witnessed this (expert, context)
        witnesses = {u.agent_lct for u in self.pending_updates[consensus_key]}
        if update.agent_lct in witnesses:
            return False  # Duplicate witness

        # Add witness
        self.pending_updates[consensus_key].append(update)

        # Check consensus
        if len(self.pending_updates[consensus_key]) >= self.min_witnesses:
            return True  # Consensus reached

        return False  # Still pending

    def get_consensus_value(self, consensus_key: str) -> Optional[float]:
        """
        Get consensus value (median of witnessed values).

        Args:
            consensus_key: Consensus key for (expert, context) pair

        Returns:
            Median quality value if consensus exists, None otherwise
        """
        if consensus_key not in self.pending_updates:
            return None

        updates = self.pending_updates[consensus_key]
        if len(updates) < self.min_witnesses:
            return None

        # Use median to resist Byzantine outliers
        qualities = [u.quality for u in updates]
        return float(np.median(qualities))


class TrustCoordinator:
    """
    Multi-agent trust coordinator for ACT societies.

    Implements:
    1. Local trust tracking per agent
    2. Byzantine consensus for trust updates
    3. Society-wide trust aggregation
    4. Trust-first expert selection
    5. Specialist emergence tracking
    """

    def __init__(
        self,
        society_id: str,
        num_experts: int = 128,
        min_witnesses: int = 3,
        trust_learning_rate: float = 0.3,
        min_evidence_threshold: int = 3
    ):
        """
        Initialize trust coordinator.

        Args:
            society_id: Society identifier (e.g., "society1@testnet")
            num_experts: Total number of experts
            min_witnesses: Byzantine consensus threshold (2f+1)
            trust_learning_rate: EWMA alpha (0.3 from Session 71)
            min_evidence_threshold: Min samples before trust-driven mode
        """
        self.society_id = society_id
        self.num_experts = num_experts
        self.min_witnesses = min_witnesses
        self.trust_learning_rate = trust_learning_rate
        self.min_evidence_threshold = min_evidence_threshold

        # Society state
        self.society_trust = SocietyTrustState(
            expert_trust={},
            expert_observations={},
            agent_contributions={},
            last_update=int(time.time())
        )

        # Consensus engine
        self.consensus = ByzantineConsensus(min_witnesses=min_witnesses)

        # Agent registry
        self.agents: Dict[str, LCTIdentity] = {}

        # Statistics
        self.mode_counts = {
            "trust_driven": 0,
            "router_explore": 0,
            "quality_recovery": 0
        }
        self.total_selections = 0

    def register_agent(self, identity: LCTIdentity):
        """Register agent with society."""
        self.agents[identity.lct_id] = identity
        self.society_trust.agent_contributions[identity.lct_id] = 0

    def propose_trust_update(
        self,
        agent_lct: str,
        expert_id: int,
        context: str,
        quality: float
    ) -> bool:
        """
        Agent proposes trust update.

        Args:
            agent_lct: Agent's LCT identity
            expert_id: Expert being evaluated
            context: Context classification
            quality: Observed quality [0, 1]

        Returns:
            True if consensus reached and trust updated
        """
        # Verify agent is registered
        if agent_lct not in self.agents:
            return False

        # Create update
        update = TrustUpdate(
            agent_lct=agent_lct,
            expert_id=expert_id,
            context=context,
            quality=quality,
            timestamp=int(time.time()),
            signature=f"sig_{hashlib.sha256((agent_lct + str(time.time())).encode()).hexdigest()[:16]}"
        )

        # Propose to consensus
        consensus_key = update.consensus_key()
        consensus_reached = self.consensus.propose_update(update)

        if consensus_reached:
            # Get consensus value
            consensus_quality = self.consensus.get_consensus_value(consensus_key)
            if consensus_quality is not None:
                # Apply update to society trust
                self._apply_trust_update(expert_id, context, consensus_quality)
                # Record contribution
                self.society_trust.agent_contributions[agent_lct] += 1
                return True

        return False

    def _apply_trust_update(self, expert_id: int, context: str, quality: float):
        """
        Apply consensus-approved trust update to society state.

        Args:
            expert_id: Expert being updated
            context: Context
            quality: Consensus quality value
        """
        key = (expert_id, context)

        # Get current trust (default 0.5)
        current_trust = self.society_trust.expert_trust.get(key, 0.5)

        # EWMA update (α=0.3 from Session 71)
        new_trust = (1 - self.trust_learning_rate) * current_trust + \
                    self.trust_learning_rate * quality

        # Store
        self.society_trust.expert_trust[key] = new_trust
        self.society_trust.expert_observations[key] = \
            self.society_trust.expert_observations.get(key, 0) + 1
        self.society_trust.last_update = int(time.time())

    def select_experts(
        self,
        router_logits: np.ndarray,
        context: str,
        k: int = 4
    ) -> Dict:
        """
        Trust-first expert selection using society trust.

        Three modes (from WEB4-PROP-006-v2.1):
        1. trust_driven: Sufficient evidence exists
        2. quality_recovery: Trust declining
        3. router_explore: Bootstrap phase

        Args:
            router_logits: Router scores [num_experts]
            context: Context classification
            k: Number of experts to select

        Returns:
            Selection result with experts and metadata
        """
        self.total_selections += 1

        # Get society trust scores for all experts in this context
        trust_scores = np.array([
            self.society_trust.expert_trust.get((i, context), 0.5)
            for i in range(self.num_experts)
        ])

        # Count evidence
        evidence_counts = np.array([
            self.society_trust.expert_observations.get((i, context), 0)
            for i in range(self.num_experts)
        ])

        total_evidence = evidence_counts.sum()
        experts_with_evidence = (evidence_counts >= self.min_evidence_threshold).sum()

        # MODE SELECTION (WEB4-PROP-006-v2.1 specification)
        if experts_with_evidence >= 2 and total_evidence >= self.min_evidence_threshold * 2:
            # TRUST-DRIVEN: Sufficient evidence
            mode = "trust_driven"
            selected_indices = np.argsort(trust_scores)[-k:][::-1]

        elif trust_scores.min() < 0.3:
            # QUALITY RECOVERY: Trust declining
            mode = "quality_recovery"
            combined = 0.5 * trust_scores + 0.5 * router_logits
            selected_indices = np.argsort(combined)[-k:][::-1]

        else:
            # ROUTER EXPLORE: Bootstrap
            mode = "router_explore"
            selected_indices = np.argsort(router_logits)[-k:][::-1]

        self.mode_counts[mode] += 1

        return {
            "experts": selected_indices.tolist(),
            "mode": mode,
            "trust_scores": trust_scores[selected_indices].tolist(),
            "evidence_counts": evidence_counts[selected_indices].tolist(),
            "society_id": self.society_id,
            "total_evidence": int(total_evidence),
            "experts_with_evidence": int(experts_with_evidence)
        }

    def get_statistics(self) -> Dict:
        """Get coordinator statistics."""
        total_selections = self.total_selections
        if total_selections == 0:
            return {}

        return {
            "society_id": self.society_id,
            "registered_agents": len(self.agents),
            "total_selections": total_selections,
            "mode_distribution": {
                mode: {
                    "count": count,
                    "percentage": 100 * count / total_selections
                }
                for mode, count in self.mode_counts.items()
            },
            "trust_driven_rate": self.mode_counts["trust_driven"] / total_selections,
            "exploration_rate": (self.mode_counts["router_explore"] +
                               self.mode_counts["quality_recovery"]) / total_selections,
            "total_trust_entries": len(self.society_trust.expert_trust),
            "total_observations": sum(self.society_trust.expert_observations.values())
        }

    def analyze_specialists(self) -> Dict:
        """
        Analyze specialist vs generalist experts.

        Based on Session 69/73 findings (54.7% specialists expected).
        """
        # Build context map per expert
        context_expert_map: Dict[int, Set[str]] = {}

        for (expert_id, context), count in self.society_trust.expert_observations.items():
            if expert_id not in context_expert_map:
                context_expert_map[expert_id] = set()
            if count >= self.min_evidence_threshold:
                context_expert_map[expert_id].add(context)

        # Classify
        specialists = []
        generalists = []

        for expert_id, contexts in context_expert_map.items():
            if len(contexts) == 1:
                specialists.append({
                    "expert_id": expert_id,
                    "context": list(contexts)[0],
                    "type": "specialist"
                })
            elif len(contexts) > 1:
                generalists.append({
                    "expert_id": expert_id,
                    "contexts": list(contexts),
                    "type": "generalist"
                })

        total = len(specialists) + len(generalists)
        specialist_rate = len(specialists) / total if total > 0 else 0

        return {
            "specialists": specialists,
            "generalists": generalists,
            "specialist_count": len(specialists),
            "generalist_count": len(generalists),
            "specialist_rate": specialist_rate,
            "total_classified": total
        }

    def export_state(self) -> Dict:
        """Export complete society trust state."""
        return {
            "society_id": self.society_id,
            "expert_trust": {
                f"{expert}:{context}": trust
                for (expert, context), trust in self.society_trust.expert_trust.items()
            },
            "expert_observations": {
                f"{expert}:{context}": count
                for (expert, context), count in self.society_trust.expert_observations.items()
            },
            "agent_contributions": self.society_trust.agent_contributions,
            "last_update": self.society_trust.last_update,
            "statistics": self.get_statistics(),
            "specialists": self.analyze_specialists()
        }


# Example usage
def demo_multi_agent_coordination():
    """
    Demonstrate multi-agent trust coordination.

    Simulates 5 agents in a society sharing trust observations.
    """
    print("\n" + "="*70)
    print("MULTI-AGENT TRUST COORDINATION DEMO")
    print("="*70)

    # Create society coordinator
    coordinator = TrustCoordinator(
        society_id="demo-society@testnet",
        num_experts=128,
        min_witnesses=3,
        trust_learning_rate=0.3,
        min_evidence_threshold=3
    )

    # Register 5 agents
    agents = []
    for i in range(5):
        identity = LCTIdentity(
            lct_id=f"lct://demo-society:citizen:agent_{i}@testnet",
            public_key=f"pubkey_{i}",
            hardware_hash=f"hw_hash_{i}",
            role="citizen"
        )
        coordinator.register_agent(identity)
        agents.append(identity)

    print(f"\n✅ Registered {len(agents)} agents to society: {coordinator.society_id}")

    # Simulate multi-agent consensus over multiple rounds
    print("\n📊 Simulating multi-agent trust consensus (10 rounds)...\n")

    contexts = ["context_0", "context_1", "context_2"]
    expert_pool = [24, 42, 73, 79, 102]  # Subset of experts

    for round_num in range(10):
        # Random expert and context each round
        expert_id = np.random.choice(expert_pool)
        context = np.random.choice(contexts)
        base_quality = {"context_0": 0.7, "context_1": 0.8, "context_2": 0.6}[context]

        # Each agent observes
        for i, agent in enumerate(agents):
            # Slightly different observations (Byzantine noise)
            quality = base_quality + np.random.normal(0, 0.05)
            quality = np.clip(quality, 0, 1)

            consensus_reached = coordinator.propose_trust_update(
                agent_lct=agent.lct_id,
                expert_id=expert_id,
                context=context,
                quality=quality
            )

            if i == 0:
                print(f"Round {round_num+1}: Expert {expert_id}, {context}")

            if consensus_reached and i == coordinator.min_witnesses - 1:
                # Check society trust
                key = (expert_id, context)
                society_trust = coordinator.society_trust.expert_trust.get(key, 0.5)
                consensus_key = f"expert_{expert_id}:context_{context}"
                consensus_val = coordinator.consensus.get_consensus_value(consensus_key)
                print(f"  ✓ Consensus reached with {coordinator.min_witnesses} witnesses")
                print(f"  → Society Trust: {society_trust:.3f} (consensus value: {consensus_val:.3f})\n")
                break

    # Test trust-first selection
    print(f"\n📊 Testing trust-first expert selection...")
    router_logits = np.random.randn(128).astype(np.float32)

    for i in range(10):
        result = coordinator.select_experts(router_logits, context, k=4)

        if i == 0 or i == 9:
            print(f"\n  Generation {i+1}:")
            print(f"    Mode: {result['mode']}")
            print(f"    Selected Experts: {result['experts']}")
            print(f"    Trust Scores: {[f'{t:.3f}' for t in result['trust_scores']]}")
            print(f"    Evidence: {result['total_evidence']} total, "
                  f"{result['experts_with_evidence']} with threshold")

    # Statistics
    print(f"\n{'='*70}")
    print("SOCIETY STATISTICS")
    print(f"{'='*70}\n")

    stats = coordinator.get_statistics()
    print(f"Society: {stats['society_id']}")
    print(f"Registered Agents: {stats['registered_agents']}")
    print(f"Total Selections: {stats['total_selections']}")
    print(f"\nMode Distribution:")
    for mode, data in stats['mode_distribution'].items():
        print(f"  {mode}: {data['count']} ({data['percentage']:.1f}%)")

    print(f"\nTrust Database:")
    print(f"  Total Trust Entries: {stats['total_trust_entries']}")
    print(f"  Total Observations: {stats['total_observations']}")

    # Export state
    state = coordinator.export_state()
    print(f"\n✅ Society state exported ({len(state['expert_trust'])} trust entries)")


if __name__ == "__main__":
    demo_multi_agent_coordination()
