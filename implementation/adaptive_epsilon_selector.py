#!/usr/bin/env python3
"""
Adaptive Epsilon Decay for Trust-First MoE

Dynamically adjusts epsilon based on trust accumulation state.

Problem:
- Fixed epsilon suboptimal across session lifecycle (Thor S77)
- Early session: Need high epsilon (0.2-0.3) to break monopoly
- Late session: Need low epsilon (0.05-0.1) to exploit accumulated trust
- Current: Fixed epsilon wastes exploration once trust is strong

Solution:
- Start with high epsilon (bootstrap diversity)
- Decay epsilon as trust evidence accumulates
- Result: Efficient exploration-to-exploitation transition

Decay Strategies:
1. Linear: ε(t) = ε₀ - (ε₀ - ε_min) × (t / T)
2. Exponential: ε(t) = ε_min + (ε₀ - ε_min) × exp(-λt)
3. Evidence-based: ε(evidence) = ε₀ × (1 - evidence_ratio)
4. Hybrid: Combine time + evidence

Based on:
- Thor Session 77: Epsilon-greedy optimal at ε=0.2
- Legion Session 70: Trust persistence enables continuation
- Session 71: Epsilon + warm-start integration
- Multi-armed bandit theory (UCB, Thompson sampling)

Author: Legion (Session 71 - Autonomous Web4 Research)
Date: 2025-12-19
"""

import sys
import numpy as np
from pathlib import Path
from typing import Dict, List, Tuple, Optional
import time
import math

# Import epsilon + warm-start base
sys.path.insert(0, str(Path(__file__).parent))
from epsilon_warmstart_integration import EpsilonWarmStartSelector
from trust_ledger_persistence import TrustLedgerPersistence


class AdaptiveEpsilonSelector(EpsilonWarmStartSelector):
    """
    Trust-first selector with adaptive epsilon decay.

    Extends EpsilonWarmStartSelector with dynamic epsilon adjustment.
    """

    def __init__(
        self,
        num_experts: int = 128,
        min_evidence_threshold: int = 3,
        epsilon_start: float = 0.2,
        epsilon_min: float = 0.05,
        decay_strategy: str = "evidence",
        decay_param: float = 0.5,
        trust_ledger: Optional[TrustLedgerPersistence] = None,
        society_id: Optional[str] = None
    ):
        """
        Initialize adaptive epsilon selector.

        Args:
            num_experts: Total experts
            min_evidence_threshold: Min samples before trust-driven
            epsilon_start: Initial epsilon (high for bootstrap)
            epsilon_min: Minimum epsilon floor (always explore a bit)
            decay_strategy: "linear", "exponential", "evidence", "hybrid"
            decay_param: Strategy-specific parameter
                - linear: total_generations (T)
                - exponential: decay rate (λ)
                - evidence: target evidence ratio (threshold)
                - hybrid: weight for time vs evidence
            trust_ledger: Ledger for persistence
            society_id: Society identifier
        """
        # Initialize parent with DYNAMIC epsilon (will update)
        super().__init__(
            num_experts=num_experts,
            min_evidence_threshold=min_evidence_threshold,
            epsilon=epsilon_start,  # Current epsilon (updated each step)
            trust_ledger=trust_ledger,
            society_id=society_id
        )

        # Adaptive epsilon config
        self.epsilon_start = epsilon_start
        self.epsilon_min = epsilon_min
        self.decay_strategy = decay_strategy
        self.decay_param = decay_param

        # Tracking for decay calculation
        self.generation_count = 0  # Time-based decay

    def _calculate_epsilon(self) -> float:
        """
        Calculate current epsilon based on decay strategy.

        Returns:
            Current epsilon value
        """
        if self.decay_strategy == "linear":
            # Linear decay: ε(t) = ε₀ - (ε₀ - ε_min) × (t / T)
            T = self.decay_param  # Total generations
            progress = min(self.generation_count / T, 1.0)
            epsilon = self.epsilon_start - (self.epsilon_start - self.epsilon_min) * progress

        elif self.decay_strategy == "exponential":
            # Exponential decay: ε(t) = ε_min + (ε₀ - ε_min) × exp(-λt)
            lambda_decay = self.decay_param
            epsilon = self.epsilon_min + (self.epsilon_start - self.epsilon_min) * \
                     math.exp(-lambda_decay * self.generation_count)

        elif self.decay_strategy == "evidence":
            # Evidence-based: ε(evidence) = ε₀ × (1 - evidence_ratio)
            target_ratio = self.decay_param  # Target evidence ratio

            # Calculate current evidence ratio
            total_observations = sum(
                sum(obs.values()) for obs in self.expert_observations.values()
            )
            max_observations = self.num_experts * len(set(
                context for obs in self.expert_observations.values() for context in obs.keys()
            ))

            if max_observations == 0:
                evidence_ratio = 0.0
            else:
                evidence_ratio = min(total_observations / (max_observations * target_ratio), 1.0)

            epsilon = self.epsilon_start * (1 - evidence_ratio)
            epsilon = max(epsilon, self.epsilon_min)

        elif self.decay_strategy == "hybrid":
            # Hybrid: Combine time and evidence
            weight_time = self.decay_param  # Weight for time (0-1)
            weight_evidence = 1 - weight_time

            # Time component (exponential)
            epsilon_time = self.epsilon_min + (self.epsilon_start - self.epsilon_min) * \
                          math.exp(-0.01 * self.generation_count)

            # Evidence component
            total_observations = sum(
                sum(obs.values()) for obs in self.expert_observations.values()
            )
            num_contexts = len(set(
                context for obs in self.expert_observations.values() for context in obs.keys()
            )) or 1

            target_obs = self.num_experts * num_contexts * 0.5
            evidence_ratio = min(total_observations / target_obs, 1.0) if target_obs > 0 else 0
            epsilon_evidence = self.epsilon_start * (1 - evidence_ratio)
            epsilon_evidence = max(epsilon_evidence, self.epsilon_min)

            # Weighted combination
            epsilon = weight_time * epsilon_time + weight_evidence * epsilon_evidence

        else:
            # Unknown strategy - use fixed start value
            epsilon = self.epsilon_start

        return max(epsilon, self.epsilon_min)  # Enforce floor

    def select_experts(
        self,
        router_logits: np.ndarray,
        context: str,
        k: int = 4
    ) -> Tuple[List[int], str]:
        """
        Select experts with ADAPTIVE epsilon.

        Overrides parent to update epsilon before selection.
        """
        # Update epsilon based on current state
        self.epsilon = self._calculate_epsilon()

        # Track generation for time-based decay
        self.generation_count += 1

        # Use parent's epsilon-greedy + trust-first logic
        return super().select_experts(router_logits, context, k)

    def get_statistics(self) -> Dict:
        """Get statistics including epsilon history."""
        stats = super().get_statistics()
        stats.update({
            "epsilon_current": self.epsilon,
            "epsilon_start": self.epsilon_start,
            "epsilon_min": self.epsilon_min,
            "decay_strategy": self.decay_strategy,
            "generation_count": self.generation_count
        })
        return stats


def demo_adaptive_epsilon():
    """
    Demonstrate adaptive epsilon decay strategies.

    Compares 4 strategies:
    1. Linear decay
    2. Exponential decay
    3. Evidence-based decay
    4. Hybrid (time + evidence)
    """
    print("\n" + "="*70)
    print("ADAPTIVE EPSILON DECAY DEMO")
    print("="*70)

    strategies = [
        ("linear", 100, "Linear decay over 100 generations"),
        ("exponential", 0.03, "Exponential decay (λ=0.03)"),
        ("evidence", 0.3, "Evidence-based (target 30% coverage)"),
        ("hybrid", 0.5, "Hybrid (50% time, 50% evidence)")
    ]

    results = {}

    for strategy_name, decay_param, description in strategies:
        print(f"\n{'='*70}")
        print(f"STRATEGY: {description}")
        print(f"{'='*70}\n")

        selector = AdaptiveEpsilonSelector(
            num_experts=128,
            min_evidence_threshold=3,
            epsilon_start=0.2,
            epsilon_min=0.05,
            decay_strategy=strategy_name,
            decay_param=decay_param,
            trust_ledger=None,
            society_id=None
        )

        print(f"Configuration:")
        print(f"  Strategy: {strategy_name}")
        print(f"  ε_start: {selector.epsilon_start}")
        print(f"  ε_min: {selector.epsilon_min}")
        print(f"  Decay param: {decay_param}")

        # Simulate 100 generations
        contexts = ["context_0", "context_1", "context_2"]
        epsilon_history = []
        mode_transitions = []
        expert_usage = {}

        for gen in range(100):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            # Select (epsilon updated internally)
            experts, mode = selector.select_experts(router_logits, context, k=4)

            # Track epsilon
            epsilon_history.append(selector.epsilon)
            mode_transitions.append(mode)

            # Simulate quality
            base_quality = 0.7 + np.random.normal(0, 0.1)
            quality = float(np.clip(base_quality, 0, 1))

            # Update trust
            selector.update_trust(experts, context, quality)

            # Track usage
            for expert_id in experts:
                expert_usage[expert_id] = expert_usage.get(expert_id, 0) + 1

        stats = selector.get_statistics()

        # Calculate epsilon decay metrics
        epsilon_initial = epsilon_history[0]
        epsilon_final = epsilon_history[-1]
        epsilon_mean = np.mean(epsilon_history)
        epsilon_decay_rate = (epsilon_initial - epsilon_final) / epsilon_initial

        # Mode distribution
        mode_counts = {}
        for mode in mode_transitions:
            mode_counts[mode] = mode_counts.get(mode, 0) + 1

        print(f"\nResults:")
        print(f"  Generations: 100")
        print(f"  Unique experts: {len(expert_usage)}")
        print(f"  Trust entries: {stats['trust_entries']}")
        print(f"\n  Epsilon decay:")
        print(f"    Initial: {epsilon_initial:.3f}")
        print(f"    Final: {epsilon_final:.3f}")
        print(f"    Mean: {epsilon_mean:.3f}")
        print(f"    Decay rate: {epsilon_decay_rate*100:.1f}%")
        print(f"\n  Mode distribution:")
        for mode, count in mode_counts.items():
            print(f"    {mode}: {count} ({100*count/100:.1f}%)")

        # Store results
        results[strategy_name] = {
            "unique_experts": len(expert_usage),
            "trust_entries": stats['trust_entries'],
            "epsilon_initial": epsilon_initial,
            "epsilon_final": epsilon_final,
            "epsilon_mean": epsilon_mean,
            "decay_rate": epsilon_decay_rate,
            "mode_counts": mode_counts,
            "epsilon_history": epsilon_history
        }

    # Comparison
    print(f"\n{'='*70}")
    print("STRATEGY COMPARISON")
    print(f"{'='*70}\n")

    print(f"{'Strategy':<15} {'Experts':<10} {'Trust':<8} {'ε_final':<10} {'Decay%':<10} {'Trust%':<10}")
    print("-" * 70)

    for strategy_name in ["linear", "exponential", "evidence", "hybrid"]:
        r = results[strategy_name]
        trust_driven_pct = 100 * r['mode_counts'].get('trust_driven', 0) / 100

        print(f"{strategy_name:<15} {r['unique_experts']:<10} {r['trust_entries']:<8} "
              f"{r['epsilon_final']:<10.3f} {r['decay_rate']*100:<10.1f} {trust_driven_pct:<10.1f}")

    print(f"\n{'='*70}")
    print("KEY INSIGHTS")
    print(f"{'='*70}\n")

    # Find best strategy
    best_trust = max(results.items(), key=lambda x: x[1]['mode_counts'].get('trust_driven', 0))
    best_diversity = max(results.items(), key=lambda x: x[1]['unique_experts'])

    print(f"Best trust-driven mode: {best_trust[0]} ({100*best_trust[1]['mode_counts'].get('trust_driven', 0)/100:.1f}%)")
    print(f"Best diversity: {best_diversity[0]} ({best_diversity[1]['unique_experts']} experts)")

    print(f"\nObservations:")
    print(f"  - Evidence-based: Adapts to actual trust accumulation")
    print(f"  - Linear: Predictable decay, may over-explore late")
    print(f"  - Exponential: Fast initial decay, good for warm-start")
    print(f"  - Hybrid: Balanced, robust to varying conditions")

    print(f"\n✅ Adaptive epsilon demo complete")


if __name__ == "__main__":
    demo_adaptive_epsilon()
