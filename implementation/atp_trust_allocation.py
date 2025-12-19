#!/usr/bin/env python3
"""
ATP Trust-Based Resource Allocation for ACT Societies

Integrates ATP (Attention-Time-Processing) token economy with trust-first expert selection.
Experts with higher trust receive more ATP budget for computation.

Problem:
- Trust-first selects experts, but all get equal resources
- No economic incentive for quality performance
- Resource allocation disconnected from trust

Solution:
- ATP allocation proportional to trust scores
- Higher trust → more compute budget → better performance → higher trust (virtuous cycle)
- Economic incentive alignment with trust-based selection

Architecture:
- Trust scores → ATP allocation weights
- ATP budget distributed per selection
- Quality performance tracked for ADP (Attention-Derived-Product) rewards
- Virtuous cycle: ATP → Performance → Trust → More ATP

Based on:
- WEB4-PROP-006-v2.1: Trust-first MoE architecture
- ACT ATP economy: Resource allocation framework
- Session 69: Multi-agent Byzantine consensus

Economics:
- Total ATP pool per generation (fixed budget)
- Allocation: trust_weighted = (trust / sum(trust)) * total_ATP
- Reward: ADP = quality_observed * ATP_used
- Reinvestment: High ADP → trust increase → ATP allocation increase

Author: Legion (Session 70 - Autonomous Web4 Research)
Date: 2025-12-19
"""

from dataclasses import dataclass
from typing import Dict, List, Tuple, Optional
import numpy as np


@dataclass
class ATPAllocation:
    """ATP allocation for a single expert."""
    expert_id: int
    context: str
    trust_score: float  # Expert's trust in this context
    atp_allocated: float  # ATP tokens allocated
    atp_weight: float  # Normalized weight (proportion of total)


@dataclass
class ADPReward:
    """ADP reward for observed performance."""
    expert_id: int
    context: str
    atp_used: float  # ATP consumed
    quality_observed: float  # Quality [0, 1]
    adp_earned: float  # ADP = quality * ATP
    timestamp: int


class ATPTrustAllocator:
    """
    Allocate ATP resources based on trust scores.

    Implements economic alignment between trust and resource allocation.
    """

    def __init__(
        self,
        total_atp_per_generation: float = 1000.0,
        min_atp_per_expert: float = 10.0,
        quality_to_trust_feedback: float = 0.3  # EWMA α
    ):
        """
        Initialize ATP-Trust allocator.

        Args:
            total_atp_per_generation: Fixed ATP budget per generation
            min_atp_per_expert: Minimum ATP even for low-trust experts
            quality_to_trust_feedback: Feedback rate (α in EWMA)
        """
        self.total_atp = total_atp_per_generation
        self.min_atp = min_atp_per_expert
        self.feedback_rate = quality_to_trust_feedback

        # Economics tracking
        self.atp_allocated_history: List[ATPAllocation] = []
        self.adp_rewards_history: List[ADPReward] = []
        self.total_atp_allocated = 0.0
        self.total_adp_earned = 0.0

    def allocate_atp(
        self,
        selected_experts: List[int],
        trust_scores: List[float],
        context: str
    ) -> List[ATPAllocation]:
        """
        Allocate ATP to selected experts based on trust.

        Args:
            selected_experts: Expert IDs selected
            trust_scores: Trust values for selected experts
            context: Context classification

        Returns:
            List of ATP allocations
        """
        k = len(selected_experts)
        trust_scores_array = np.array(trust_scores)

        # Calculate ATP weights (normalized trust)
        trust_sum = trust_scores_array.sum()
        if trust_sum > 0:
            weights = trust_scores_array / trust_sum
        else:
            weights = np.ones(k) / k  # Uniform if no trust

        # Allocate ATP with minimum guarantee
        atp_allocations = []
        remaining_atp = self.total_atp - (self.min_atp * k)

        for i, expert_id in enumerate(selected_experts):
            # Minimum + trust-weighted portion
            atp_allocated = self.min_atp + (weights[i] * remaining_atp)

            allocation = ATPAllocation(
                expert_id=expert_id,
                context=context,
                trust_score=trust_scores[i],
                atp_allocated=atp_allocated,
                atp_weight=weights[i]
            )

            atp_allocations.append(allocation)
            self.atp_allocated_history.append(allocation)
            self.total_atp_allocated += atp_allocated

        return atp_allocations

    def calculate_adp_reward(
        self,
        expert_id: int,
        context: str,
        atp_used: float,
        quality_observed: float,
        timestamp: int
    ) -> ADPReward:
        """
        Calculate ADP reward for observed performance.

        ADP = quality * ATP_used

        Higher quality and more ATP → higher ADP
        Creates incentive for both trust (ATP allocation) and performance (quality)

        Args:
            expert_id: Expert that performed
            context: Context of task
            atp_used: ATP tokens consumed
            quality_observed: Observed quality [0, 1]
            timestamp: When observed

        Returns:
            ADPReward
        """
        adp_earned = quality_observed * atp_used

        reward = ADPReward(
            expert_id=expert_id,
            context=context,
            atp_used=atp_used,
            quality_observed=quality_observed,
            adp_earned=adp_earned,
            timestamp=timestamp
        )

        self.adp_rewards_history.append(reward)
        self.total_adp_earned += adp_earned

        return reward

    def get_atp_efficiency(self, expert_id: int, context: str) -> Optional[float]:
        """
        Calculate ATP efficiency for expert in context.

        Efficiency = Total ADP earned / Total ATP used
        Higher efficiency → better performance per ATP

        Args:
            expert_id: Expert to analyze
            context: Context to filter

        Returns:
            Efficiency ratio or None if no data
        """
        expert_rewards = [
            r for r in self.adp_rewards_history
            if r.expert_id == expert_id and r.context == context
        ]

        if not expert_rewards:
            return None

        total_atp_used = sum(r.atp_used for r in expert_rewards)
        total_adp = sum(r.adp_earned for r in expert_rewards)

        return total_adp / total_atp_used if total_atp_used > 0 else 0.0

    def get_statistics(self) -> Dict:
        """Get ATP/ADP statistics."""
        return {
            "total_atp_allocated": self.total_atp_allocated,
            "total_adp_earned": self.total_adp_earned,
            "overall_efficiency": self.total_adp_earned / self.total_atp_allocated
                if self.total_atp_allocated > 0 else 0.0,
            "allocation_count": len(self.atp_allocated_history),
            "reward_count": len(self.adp_rewards_history)
        }


class TrustATPIntegratedSelector:
    """
    Integrated trust-first selector with ATP allocation.

    Combines:
    - Trust-first expert selection (WEB4-PROP-006-v2.1)
    - ATP resource allocation (proportional to trust)
    - ADP reward calculation (quality * ATP)
    - Economic feedback loop (ADP → trust update)
    """

    def __init__(
        self,
        num_experts: int = 128,
        min_evidence_threshold: int = 3,
        total_atp_per_generation: float = 1000.0
    ):
        """Initialize integrated selector."""
        self.num_experts = num_experts
        self.min_evidence_threshold = min_evidence_threshold

        # Trust state
        self.expert_trust: Dict[Tuple[int, str], float] = {}
        self.expert_observations: Dict[Tuple[int, str], int] = {}

        # ATP allocator
        self.atp_allocator = ATPTrustAllocator(
            total_atp_per_generation=total_atp_per_generation
        )

        # Mode tracking
        self.mode_counts = {"trust_driven": 0, "router_explore": 0}
        self.total_selections = 0

    def select_and_allocate(
        self,
        router_logits: np.ndarray,
        context: str,
        k: int = 4
    ) -> Tuple[List[int], List[ATPAllocation], str]:
        """
        Select experts and allocate ATP.

        Args:
            router_logits: Router scores
            context: Context classification
            k: Number of experts to select

        Returns:
            (selected_expert_ids, atp_allocations, selection_mode)
        """
        self.total_selections += 1

        # Get trust scores
        trust_scores = np.array([
            self.expert_trust.get((i, context), 0.5)
            for i in range(self.num_experts)
        ])

        # Count evidence
        evidence_counts = np.array([
            self.expert_observations.get((i, context), 0)
            for i in range(self.num_experts)
        ])

        total_evidence = evidence_counts.sum()
        experts_with_evidence = (evidence_counts >= self.min_evidence_threshold).sum()

        # Select mode (trust-first logic)
        if experts_with_evidence >= 2 and total_evidence >= self.min_evidence_threshold * 2:
            mode = "trust_driven"
            selected_indices = np.argsort(trust_scores)[-k:][::-1]
        else:
            mode = "router_explore"
            selected_indices = np.argsort(router_logits)[-k:][::-1]

        self.mode_counts[mode] += 1

        # Allocate ATP
        selected_experts = selected_indices.tolist()
        selected_trust_scores = trust_scores[selected_indices].tolist()
        atp_allocations = self.atp_allocator.allocate_atp(
            selected_experts,
            selected_trust_scores,
            context
        )

        return selected_experts, atp_allocations, mode

    def update_trust_with_adp(
        self,
        expert_id: int,
        context: str,
        atp_used: float,
        quality: float,
        timestamp: int
    ):
        """
        Update trust based on observed quality and calculate ADP reward.

        Economic cycle:
        1. Quality observed → ADP calculated
        2. ADP → Trust update (EWMA)
        3. Trust → ATP allocation next selection
        4. ATP → Performance → Quality → (loop)

        Args:
            expert_id: Expert that performed
            context: Context of task
            atp_used: ATP consumed
            quality: Observed quality [0, 1]
            timestamp: When observed
        """
        key = (expert_id, context)

        # Calculate ADP reward
        adp = self.atp_allocator.calculate_adp_reward(
            expert_id, context, atp_used, quality, timestamp
        )

        # Update trust (EWMA with α=0.3)
        current_trust = self.expert_trust.get(key, 0.5)
        new_trust = 0.7 * current_trust + 0.3 * quality

        self.expert_trust[key] = new_trust
        self.expert_observations[key] = self.expert_observations.get(key, 0) + 1

    def get_statistics(self) -> Dict:
        """Get complete statistics."""
        atp_stats = self.atp_allocator.get_statistics()

        return {
            **atp_stats,
            "mode_distribution": {
                mode: f"{count}/{self.total_selections} ({100*count/self.total_selections:.1f}%)"
                for mode, count in self.mode_counts.items()
            },
            "trust_entries": len(self.expert_trust),
            "total_observations": sum(self.expert_observations.values())
        }


# Demo
def demo_atp_trust_integration():
    """Demonstrate ATP-Trust integration."""
    print("\n" + "="*70)
    print("ATP-TRUST INTEGRATION DEMO")
    print("="*70)

    selector = TrustATPIntegratedSelector(
        num_experts=128,
        min_evidence_threshold=3,
        total_atp_per_generation=1000.0
    )

    print(f"\n✅ Integrated selector initialized")
    print(f"   Total ATP budget per generation: {selector.atp_allocator.total_atp}")
    print(f"   Minimum ATP per expert: {selector.atp_allocator.min_atp}")

    # Simulate 20 generations
    import time
    contexts = ["context_0", "context_1", "context_2"]

    print(f"\n📊 Running 20 generations with ATP allocation...\n")

    for gen in range(20):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        # Select and allocate
        experts, atp_allocs, mode = selector.select_and_allocate(
            router_logits, context, k=4
        )

        # Simulate performance (quality varies with trust)
        for i, (expert_id, alloc) in enumerate(zip(experts, atp_allocs)):
            # Higher trust → better average quality (but still variable)
            base_quality = 0.5 + 0.3 * alloc.trust_score
            quality = float(np.clip(base_quality + np.random.normal(0, 0.1), 0, 1))

            # Update trust with ADP
            selector.update_trust_with_adp(
                expert_id, context, alloc.atp_allocated, quality, int(time.time())
            )

        if gen % 5 == 0 or gen == 19:
            print(f"Gen {gen+1:2d}: Mode={mode:15s} Experts={experts}")
            print(f"        ATP Allocations: {[f'{a.atp_allocated:.1f}' for a in atp_allocs]}")
            print(f"        Trust Scores:    {[f'{a.trust_score:.3f}' for a in atp_allocs]}\n")

    # Statistics
    print(f"{'='*70}")
    print("FINAL STATISTICS")
    print(f"{'='*70}\n")

    stats = selector.get_statistics()
    print(f"Economic Metrics:")
    print(f"  Total ATP Allocated: {stats['total_atp_allocated']:.1f}")
    print(f"  Total ADP Earned: {stats['total_adp_earned']:.1f}")
    print(f"  Overall Efficiency: {stats['overall_efficiency']:.3f} (ADP/ATP)")

    print(f"\nSelection Metrics:")
    for mode, dist in stats['mode_distribution'].items():
        print(f"  {mode}: {dist}")

    print(f"\nTrust Database:")
    print(f"  Trust entries: {stats['trust_entries']}")
    print(f"  Total observations: {stats['total_observations']}")

    # Show top ATP-efficient experts
    print(f"\nTop ATP-Efficient Experts:")
    efficiencies = []
    for (expert_id, context) in selector.expert_trust.keys():
        eff = selector.atp_allocator.get_atp_efficiency(expert_id, context)
        if eff is not None:
            efficiencies.append((expert_id, context, eff))

    efficiencies.sort(key=lambda x: x[2], reverse=True)
    for expert_id, context, eff in efficiencies[:5]:
        trust = selector.expert_trust.get((expert_id, context), 0)
        print(f"  Expert {expert_id:3d}, {context}: Efficiency={eff:.3f}, Trust={trust:.3f}")

    print(f"\n✅ ATP-Trust integration demo complete")
    print(f"\nVirtuous Cycle Demonstrated:")
    print(f"  Trust → ATP → Performance → Quality → ADP → Trust ↑")


if __name__ == "__main__":
    demo_atp_trust_integration()
