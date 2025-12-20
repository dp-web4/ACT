#!/usr/bin/env python3
"""
Epsilon-Greedy + Warm-Start Integration

Combines Thor's Session 77 epsilon-greedy discovery with Legion's Session 70 persistence.

Problem:
- Real model router monopoly prevents trust accumulation (Thor S76)
- Epsilon-greedy breaks monopoly but starts from scratch each session (Thor S77)
- Warm-start loads prior trust but doesn't break monopoly (Legion S70)

Solution:
- Epsilon-greedy for initial evidence gathering (breaks monopoly)
- Persistence to carry trust forward across sessions
- Result: Best of both - monopoly breaking + continuous evolution

Architecture:
1. Session N: High epsilon (0.2-0.3) → gather diverse evidence
2. Save snapshot to ledger → persist accumulated trust
3. Session N+1: Load snapshot (warm-start) + lower epsilon (0.1)
4. Repeat → trust evolves continuously without monopoly

Based on:
- Thor Sessions 76-77: Real model monopoly + epsilon solution
- Legion Session 70: Trust persistence infrastructure
- WEB4-PROP-006-v2.1: Trust-first MoE standard

Author: Legion (Session 71 - Autonomous Web4 Research)
Date: 2025-12-19
"""

import sys
import numpy as np
from pathlib import Path
from typing import Dict, List, Tuple, Optional
import time

# Import Session 70 persistence
sys.path.insert(0, str(Path(__file__).parent))
from trust_ledger_persistence import TrustLedgerPersistence, warm_start_trust_selector


class EpsilonWarmStartSelector:
    """
    Trust-first selector with epsilon-greedy + warm-start.

    Combines:
    - Epsilon-greedy forced exploration (Thor S77)
    - Trust persistence/warm-start (Legion S70)
    - Trust-first conditional logic (WEB4-v2.1)
    """

    def __init__(
        self,
        num_experts: int = 128,
        min_evidence_threshold: int = 3,
        epsilon: float = 0.2,  # Thor S77 optimal
        trust_ledger: Optional[TrustLedgerPersistence] = None,
        society_id: Optional[str] = None
    ):
        """
        Initialize epsilon + warm-start selector.

        Args:
            num_experts: Total experts
            min_evidence_threshold: Min samples before trust-driven
            epsilon: Forced exploration probability (0.2 optimal)
            trust_ledger: Ledger for persistence (None = no persistence)
            society_id: Society identifier for warm-start
        """
        self.num_experts = num_experts
        self.min_evidence_threshold = min_evidence_threshold
        self.epsilon = epsilon

        # Trust state
        self.expert_trust: Dict[int, Dict[str, float]] = {}
        self.expert_observations: Dict[int, Dict[str, int]] = {}

        # Mode tracking
        self.mode_counts = {
            "trust_driven": 0,
            "router_explore": 0,
            "forced_exploration": 0
        }
        self.total_selections = 0

        # Warm-start if ledger provided
        if trust_ledger and society_id:
            self.warm_started = warm_start_trust_selector(
                self, trust_ledger, society_id
            )
        else:
            self.warm_started = False

    def select_experts(
        self,
        router_logits: np.ndarray,
        context: str,
        k: int = 4
    ) -> Tuple[List[int], str]:
        """
        Select experts with epsilon-greedy + trust-first logic.

        Selection priority:
        1. With probability epsilon → forced_exploration (random)
        2. If trust evidence exists → trust_driven
        3. Else → router_explore

        Args:
            router_logits: Router scores [num_experts]
            context: Context classification
            k: Number of experts to select

        Returns:
            (selected_expert_ids, selection_mode)
        """
        self.total_selections += 1

        # 1. Epsilon-greedy forced exploration (Thor S77)
        if self.epsilon > 0 and np.random.random() < self.epsilon:
            selected = np.random.choice(self.num_experts, size=k, replace=False).tolist()
            self.mode_counts["forced_exploration"] += 1
            return selected, "forced_exploration"

        # Get trust scores
        trust_scores = np.array([
            self.expert_trust.get(i, {}).get(context, 0.5)
            for i in range(self.num_experts)
        ])

        # Count evidence
        evidence_counts = np.array([
            self.expert_observations.get(i, {}).get(context, 0)
            for i in range(self.num_experts)
        ])

        total_evidence = evidence_counts.sum()
        experts_with_evidence = (evidence_counts >= self.min_evidence_threshold).sum()

        # 2. Trust-driven (if evidence exists)
        if experts_with_evidence >= 2 and total_evidence >= self.min_evidence_threshold * 2:
            mode = "trust_driven"
            selected_indices = np.argsort(trust_scores)[-k:][::-1]
        else:
            # 3. Router explore (bootstrap/fallback)
            mode = "router_explore"
            selected_indices = np.argsort(router_logits)[-k:][::-1]

        self.mode_counts[mode] += 1
        return selected_indices.tolist(), mode

    def update_trust(self, expert_ids: List[int], context: str, quality: float):
        """Update trust based on observed quality (EWMA α=0.3)."""
        alpha = 0.3

        for expert_id in expert_ids:
            if expert_id not in self.expert_trust:
                self.expert_trust[expert_id] = {}
            if expert_id not in self.expert_observations:
                self.expert_observations[expert_id] = {}

            current_trust = self.expert_trust[expert_id].get(context, 0.5)
            new_trust = (1 - alpha) * current_trust + alpha * quality

            self.expert_trust[expert_id][context] = new_trust
            self.expert_observations[expert_id][context] = \
                self.expert_observations[expert_id].get(context, 0) + 1

    def save_snapshot(
        self,
        ledger: TrustLedgerPersistence,
        society_id: str,
        session_id: str,
        metadata: Optional[Dict] = None
    ) -> str:
        """Save current trust state to ledger."""
        # Convert to (expert, context) → trust format
        trust_state = {}
        observation_counts = {}

        for expert_id, contexts in self.expert_trust.items():
            for context, trust_value in contexts.items():
                key = (expert_id, context)
                trust_state[key] = trust_value
                observation_counts[key] = self.expert_observations[expert_id][context]

        return ledger.save_snapshot(
            society_id=society_id,
            session_id=session_id,
            trust_state=trust_state,
            observation_counts=observation_counts,
            metadata=metadata or self.get_statistics()
        )

    def get_statistics(self) -> Dict:
        """Get selection statistics."""
        total = self.total_selections
        if total == 0:
            return {}

        return {
            "total_selections": total,
            "mode_distribution": {
                mode: {
                    "count": count,
                    "percentage": 100 * count / total
                }
                for mode, count in self.mode_counts.items()
            },
            "warm_started": self.warm_started,
            "epsilon": self.epsilon,
            "trust_entries": sum(len(contexts) for contexts in self.expert_trust.values()),
            "total_observations": sum(
                sum(obs.values()) for obs in self.expert_observations.values()
            )
        }


def demo_epsilon_warmstart():
    """
    Demonstrate epsilon + warm-start integration.

    Simulates 2 sessions:
    - Session A: Cold start with ε=0.2 (gather evidence)
    - Session B: Warm start with ε=0.1 (leverage prior trust)
    """
    print("\n" + "="*70)
    print("EPSILON-GREEDY + WARM-START INTEGRATION DEMO")
    print("="*70)

    # Setup ledger
    ledger_dir = Path("/tmp/epsilon_warmstart_demo")
    ledger = TrustLedgerPersistence(ledger_dir)

    print(f"\n✅ Trust ledger initialized: {ledger_dir}")

    # Session A: Cold start with ε=0.2
    print(f"\n{'='*70}")
    print("SESSION A: COLD START (ε=0.2 - High Exploration)")
    print(f"{'='*70}\n")

    selector_a = EpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.2,  # High epsilon to break monopoly
        trust_ledger=None,  # Cold start
        society_id=None
    )

    print(f"Configuration:")
    print(f"  Epsilon: {selector_a.epsilon}")
    print(f"  Warm-started: {selector_a.warm_started}")

    # Simulate 50 generations
    contexts = ["context_0", "context_1", "context_2"]
    expert_usage_a = {}

    for gen in range(50):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        # Select
        experts, mode = selector_a.select_experts(router_logits, context, k=4)

        # Simulate quality
        base_quality = 0.7 + np.random.normal(0, 0.1)
        quality = float(np.clip(base_quality, 0, 1))

        # Update trust
        selector_a.update_trust(experts, context, quality)

        # Track usage
        for expert_id in experts:
            expert_usage_a[expert_id] = expert_usage_a.get(expert_id, 0) + 1

    stats_a = selector_a.get_statistics()
    print(f"\nSession A Results:")
    print(f"  Generations: 50")
    print(f"  Unique experts: {len(expert_usage_a)}")
    print(f"  Trust entries: {stats_a['trust_entries']}")
    print(f"  Mode distribution:")
    for mode, data in stats_a['mode_distribution'].items():
        print(f"    {mode}: {data['count']} ({data['percentage']:.1f}%)")

    # Save snapshot
    snapshot_id = selector_a.save_snapshot(
        ledger=ledger,
        society_id="demo-society",
        session_id="session_a",
        metadata={"generations": 50, "epsilon": 0.2}
    )

    print(f"\n  ✅ Snapshot saved: {snapshot_id}")

    # Session B: Warm start with ε=0.1
    print(f"\n{'='*70}")
    print("SESSION B: WARM START (ε=0.1 - Lower Exploration)")
    print(f"{'='*70}\n")

    selector_b = EpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.1,  # Lower epsilon (trust already accumulated)
        trust_ledger=ledger,
        society_id="demo-society"
    )

    print(f"Configuration:")
    print(f"  Epsilon: {selector_b.epsilon}")
    print(f"  Warm-started: {selector_b.warm_started}")

    # Count pre-loaded trust
    preloaded_entries = sum(len(contexts) for contexts in selector_b.expert_trust.values())
    print(f"  Pre-loaded trust entries: {preloaded_entries}")

    # Simulate 50 more generations
    expert_usage_b = {}

    for gen in range(50):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        experts, mode = selector_b.select_experts(router_logits, context, k=4)

        base_quality = 0.7 + np.random.normal(0, 0.1)
        quality = float(np.clip(base_quality, 0, 1))

        selector_b.update_trust(experts, context, quality)

        for expert_id in experts:
            expert_usage_b[expert_id] = expert_usage_b.get(expert_id, 0) + 1

    stats_b = selector_b.get_statistics()
    print(f"\nSession B Results:")
    print(f"  Generations: 50")
    print(f"  Unique experts: {len(expert_usage_b)}")
    print(f"  Trust entries: {stats_b['trust_entries']}")
    print(f"  Mode distribution:")
    for mode, data in stats_b['mode_distribution'].items():
        print(f"    {mode}: {data['count']} ({data['percentage']:.1f}%)")

    # Comparison
    print(f"\n{'='*70}")
    print("COMPARISON: SESSION A vs SESSION B")
    print(f"{'='*70}\n")

    print(f"Session A (Cold, ε=0.2):")
    print(f"  Unique experts: {len(expert_usage_a)}")
    print(f"  Forced exploration: {stats_a['mode_distribution']['forced_exploration']['percentage']:.1f}%")

    print(f"\nSession B (Warm, ε=0.1):")
    print(f"  Unique experts: {len(expert_usage_b)}")
    print(f"  Forced exploration: {stats_b['mode_distribution']['forced_exploration']['percentage']:.1f}%")
    print(f"  Trust-driven: {stats_b['mode_distribution']['trust_driven']['percentage']:.1f}%")

    combined_experts = len(set(expert_usage_a.keys()) | set(expert_usage_b.keys()))
    print(f"\nCombined across both sessions:")
    print(f"  Total unique experts: {combined_experts}")

    print(f"\n✅ Integration demo complete")
    print(f"\nKey Insight:")
    print(f"  Session A: High epsilon breaks monopoly, gathers diverse evidence")
    print(f"  Session B: Warm-start + lower epsilon leverages prior trust")
    print(f"  Result: Continuous evolution without monopoly!")


if __name__ == "__main__":
    demo_epsilon_warmstart()
