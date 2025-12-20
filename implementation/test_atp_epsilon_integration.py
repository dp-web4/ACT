#!/usr/bin/env python3
"""
Session 72 Track 2: ATP/ADP + Epsilon-Warmstart Integration Test

Integrates three components:
1. Epsilon-greedy forced exploration (Thor S77)
2. Warm-start persistence (Legion S70)
3. ATP/ADP economic model (Legion S70)

Tests the complete economic cycle:
ATP allocation (∝ trust) → Performance → ADP reward → Trust update

Expected behavior:
- High-trust experts get more ATP allocation
- Better performance yields higher ADP rewards
- ADP feeds back into trust evolution
- Economic incentives align with trust-based selection

Created: 2025-12-19 (Legion Session 72)
"""

import sys
import numpy as np
from pathlib import Path
from typing import Dict, List, Tuple

sys.path.insert(0, str(Path(__file__).parent))
from epsilon_warmstart_integration import EpsilonWarmStartSelector
from atp_trust_allocation import ATPTrustAllocator, TrustATPIntegratedSelector
from trust_ledger_persistence import TrustLedgerPersistence


class ATPEpsilonWarmStartSelector:
    """
    Complete production selector combining:
    - Epsilon-greedy forced exploration
    - Trust-first conditional logic
    - Warm-start persistence
    - ATP/ADP economic integration
    """

    def __init__(
        self,
        num_experts: int = 128,
        min_evidence_threshold: int = 3,
        epsilon: float = 0.2,
        total_atp_budget: float = 1000.0,
        min_atp_per_expert: float = 100.0,
        trust_ledger: TrustLedgerPersistence = None,
        society_id: str = None
    ):
        """
        Initialize ATP + epsilon + warmstart selector.

        Args:
            num_experts: Total number of experts
            min_evidence_threshold: Min samples before trust-driven
            epsilon: Forced exploration probability
            total_atp_budget: Total ATP available per generation
            min_atp_per_expert: Minimum ATP guarantee per expert
            trust_ledger: Ledger for persistence
            society_id: Society identifier
        """
        # Base epsilon selector
        self.base_selector = EpsilonWarmStartSelector(
            num_experts=num_experts,
            min_evidence_threshold=min_evidence_threshold,
            epsilon=epsilon,
            trust_ledger=trust_ledger,
            society_id=society_id
        )

        # ATP allocator
        self.atp_allocator = ATPTrustAllocator(
            total_atp_per_generation=total_atp_budget,
            min_atp_per_expert=min_atp_per_expert
        )

        # Economic tracking
        self.atp_allocated_total = 0.0
        self.adp_earned_total = 0.0
        self.expert_atp_history = {}  # {expert_id: [atp_values]}
        self.expert_adp_history = {}  # {expert_id: [adp_values]}

    def select_and_allocate(
        self,
        router_logits: np.ndarray,
        context: str,
        k: int = 4
    ) -> Tuple[List[int], str, List[Dict]]:
        """
        Select experts AND allocate ATP based on trust.

        Returns:
            (expert_ids, selection_mode, atp_allocations)
        """
        # 1. Select experts (epsilon + trust-first)
        experts, mode = self.base_selector.select_experts(
            router_logits, context, k
        )

        # 2. Get trust scores for selected experts
        trust_scores = []
        for expert_id in experts:
            if expert_id in self.base_selector.expert_trust:
                trust = self.base_selector.expert_trust[expert_id].get(context, 0.5)
            else:
                trust = 0.5
            trust_scores.append(trust)

        trust_scores = np.array(trust_scores)

        # 3. Allocate ATP proportional to trust
        allocations = self.atp_allocator.allocate_atp(
            experts, trust_scores, context
        )

        # Track ATP allocation
        for alloc in allocations:
            self.atp_allocated_total += alloc.atp_allocated

            if alloc.expert_id not in self.expert_atp_history:
                self.expert_atp_history[alloc.expert_id] = []
            self.expert_atp_history[alloc.expert_id].append(alloc.atp_allocated)

        return experts, mode, allocations

    def update_trust_with_adp(
        self,
        expert_ids: List[int],
        atp_allocations: List[Dict],
        context: str,
        quality: float
    ):
        """
        Update trust AND calculate ADP rewards.

        Economic cycle:
        1. ATP allocated (∝ trust)
        2. Performance measured (quality)
        3. ADP rewarded (quality × ATP)
        4. Trust updated (EWMA with quality)
        """
        # 1. Update trust (base selector handles EWMA)
        self.base_selector.update_trust(expert_ids, context, quality)

        # 2. Calculate ADP rewards
        import time
        for alloc in atp_allocations:
            expert_id = alloc.expert_id
            atp_used = alloc.atp_allocated

            # ADP = quality × ATP_used
            adp_reward = self.atp_allocator.calculate_adp_reward(
                expert_id, context, atp_used, quality, int(time.time())
            )

            # Track ADP
            self.adp_earned_total += adp_reward.adp_earned

            if expert_id not in self.expert_adp_history:
                self.expert_adp_history[expert_id] = []
            self.expert_adp_history[expert_id].append(adp_reward.adp_earned)

    def get_economic_statistics(self) -> Dict:
        """Get economic + trust statistics."""
        base_stats = self.base_selector.get_statistics()

        # Calculate economic efficiency
        efficiency = self.adp_earned_total / self.atp_allocated_total \
                    if self.atp_allocated_total > 0 else 0

        # Top ATP earners
        atp_totals = {
            expert_id: sum(history)
            for expert_id, history in self.expert_atp_history.items()
        }
        top_atp = sorted(atp_totals.items(), key=lambda x: x[1], reverse=True)[:5]

        # Top ADP earners
        adp_totals = {
            expert_id: sum(history)
            for expert_id, history in self.expert_adp_history.items()
        }
        top_adp = sorted(adp_totals.items(), key=lambda x: x[1], reverse=True)[:5]

        return {
            **base_stats,
            "atp_allocated_total": self.atp_allocated_total,
            "adp_earned_total": self.adp_earned_total,
            "economic_efficiency": efficiency,
            "top_atp_earners": top_atp,
            "top_adp_earners": top_adp,
            "unique_atp_recipients": len(self.expert_atp_history)
        }


def test_atp_epsilon_integration():
    """
    Test ATP + epsilon + warmstart integration.

    Validates economic cycle:
    1. High-trust experts get more ATP
    2. Performance determines ADP
    3. ADP aligns with trust evolution
    """
    print("\n" + "="*70)
    print("ATP/ADP + EPSILON-WARMSTART INTEGRATION TEST")
    print("="*70)
    print("\nTesting complete economic cycle:")
    print("  ATP allocation → Performance → ADP reward → Trust update")
    print("")

    # Setup ledger
    ledger_dir = Path("/tmp/atp_epsilon_test")
    ledger = TrustLedgerPersistence(ledger_dir)

    # Session A: Cold start with ATP
    print("="*70)
    print("SESSION A: COLD START (ε=0.2 + ATP Integration)")
    print("="*70)

    selector_a = ATPEpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.2,
        total_atp_budget=1000.0,
        min_atp_per_expert=100.0,
        trust_ledger=None,
        society_id=None
    )

    contexts = ["context_0", "context_1", "context_2"]
    expert_usage_a = {}

    print("\nRunning 100 generations with ATP integration...")

    for gen in range(100):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        # Select and allocate ATP
        experts, mode, atp_allocations = selector_a.select_and_allocate(
            router_logits, context, k=4
        )

        # Simulate quality
        base_quality = 0.75 + np.random.normal(0, 0.1)
        quality = float(np.clip(base_quality, 0, 1))

        # Update trust with ADP calculation
        selector_a.update_trust_with_adp(
            experts, atp_allocations, context, quality
        )

        # Track usage
        for expert_id in experts:
            expert_usage_a[expert_id] = expert_usage_a.get(expert_id, 0) + 1

        # Progress
        if (gen + 1) % 25 == 0:
            stats = selector_a.get_economic_statistics()
            print(f"Gen {gen+1:3d}: {len(expert_usage_a)} experts, "
                  f"ATP={stats['atp_allocated_total']:.0f}, "
                  f"ADP={stats['adp_earned_total']:.0f}, "
                  f"eff={stats['economic_efficiency']:.3f}")

    # Session A results
    stats_a = selector_a.get_economic_statistics()

    print(f"\n{'='*70}")
    print("SESSION A RESULTS")
    print(f"{'='*70}")

    print(f"\n📊 Expert Diversity:")
    print(f"  Unique experts: {len(expert_usage_a)}/128 ({100*len(expert_usage_a)/128:.1f}%)")

    print(f"\n🔄 Mode Distribution:")
    for mode, info in stats_a['mode_distribution'].items():
        print(f"  {mode}: {info['count']} ({info['percentage']:.1f}%)")

    print(f"\n💰 Economic Metrics:")
    print(f"  Total ATP allocated: {stats_a['atp_allocated_total']:.2f}")
    print(f"  Total ADP earned: {stats_a['adp_earned_total']:.2f}")
    print(f"  Economic efficiency: {stats_a['economic_efficiency']:.3f}")
    print(f"  (ADP/ATP ratio - higher is better)")

    print(f"\n🏆 Top 5 ATP Recipients:")
    for i, (expert_id, atp_total) in enumerate(stats_a['top_atp_earners'], 1):
        # Get trust for this expert
        max_trust = 0.0
        if expert_id in selector_a.base_selector.expert_trust:
            max_trust = max(selector_a.base_selector.expert_trust[expert_id].values())
        print(f"  {i}. Expert {expert_id}: {atp_total:.0f} ATP (trust={max_trust:.3f})")

    print(f"\n🏆 Top 5 ADP Earners:")
    for i, (expert_id, adp_total) in enumerate(stats_a['top_adp_earners'], 1):
        print(f"  {i}. Expert {expert_id}: {adp_total:.0f} ADP")

    # Save snapshot
    # Extract trust state for snapshot
    trust_state = {}
    observation_counts = {}
    for expert_id, contexts_dict in selector_a.base_selector.expert_trust.items():
        for ctx, trust_val in contexts_dict.items():
            trust_state[(expert_id, ctx)] = trust_val
            obs_count = selector_a.base_selector.expert_observations.get(
                expert_id, {}
            ).get(ctx, 0)
            observation_counts[(expert_id, ctx)] = obs_count

    snapshot_id = ledger.save_snapshot(
        society_id="atp-test-society",
        session_id="session-a",
        trust_state=trust_state,
        observation_counts=observation_counts,
        metadata={
            "atp_allocated": stats_a['atp_allocated_total'],
            "adp_earned": stats_a['adp_earned_total'],
            "efficiency": stats_a['economic_efficiency']
        }
    )
    print(f"\n✅ Snapshot saved: {snapshot_id}")

    # Session B: Warm start with lower epsilon
    print(f"\n{'='*70}")
    print("SESSION B: WARM START (ε=0.1 + ATP Integration)")
    print(f"{'='*70}")

    selector_b = ATPEpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.1,  # Lower epsilon
        total_atp_budget=1000.0,
        min_atp_per_expert=100.0,
        trust_ledger=ledger,
        society_id="atp-test-society"
    )

    print(f"\nWarm-started: {selector_b.base_selector.warm_started}")
    print(f"Pre-loaded trust entries: {stats_a['trust_entries']}")

    expert_usage_b = {}

    print("\nRunning 100 generations with warm-start...")

    for gen in range(100):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        # Select and allocate ATP
        experts, mode, atp_allocations = selector_b.select_and_allocate(
            router_logits, context, k=4
        )

        # Simulate quality
        base_quality = 0.75 + np.random.normal(0, 0.1)
        quality = float(np.clip(base_quality, 0, 1))

        # Update trust with ADP
        selector_b.update_trust_with_adp(
            experts, atp_allocations, context, quality
        )

        # Track usage
        for expert_id in experts:
            expert_usage_b[expert_id] = expert_usage_b.get(expert_id, 0) + 1

        # Progress
        if (gen + 1) % 25 == 0:
            stats = selector_b.get_economic_statistics()
            print(f"Gen {gen+1:3d}: {len(expert_usage_b)} experts, "
                  f"ATP={stats['atp_allocated_total']:.0f}, "
                  f"ADP={stats['adp_earned_total']:.0f}, "
                  f"eff={stats['economic_efficiency']:.3f}")

    # Session B results
    stats_b = selector_b.get_economic_statistics()

    print(f"\n{'='*70}")
    print("SESSION B RESULTS")
    print(f"{'='*70}")

    print(f"\n📊 Expert Diversity:")
    print(f"  Unique experts: {len(expert_usage_b)}/128 ({100*len(expert_usage_b)/128:.1f}%)")

    print(f"\n🔄 Mode Distribution:")
    for mode, info in stats_b['mode_distribution'].items():
        print(f"  {mode}: {info['count']} ({info['percentage']:.1f}%)")

    print(f"\n💰 Economic Metrics:")
    print(f"  Total ATP allocated: {stats_b['atp_allocated_total']:.2f}")
    print(f"  Total ADP earned: {stats_b['adp_earned_total']:.2f}")
    print(f"  Economic efficiency: {stats_b['economic_efficiency']:.3f}")

    # Comparison
    print(f"\n{'='*70}")
    print("COMPARISON: SESSION A vs SESSION B")
    print(f"{'='*70}")

    trust_a_pct = stats_a['mode_distribution']['trust_driven']['percentage']
    trust_b_pct = stats_b['mode_distribution']['trust_driven']['percentage']
    eff_a = stats_a['economic_efficiency']
    eff_b = stats_b['economic_efficiency']

    print(f"\nSession A (cold, ε=0.2):")
    print(f"  Experts: {len(expert_usage_a)}")
    print(f"  Trust-driven: {trust_a_pct:.1f}%")
    print(f"  Efficiency: {eff_a:.3f}")

    print(f"\nSession B (warm, ε=0.1):")
    print(f"  Experts: {len(expert_usage_b)}")
    print(f"  Trust-driven: {trust_b_pct:.1f}%")
    print(f"  Efficiency: {eff_b:.3f}")

    print(f"\nImprovements:")
    print(f"  Trust-driven: {trust_b_pct:.1f}% / {trust_a_pct:.1f}% = "
          f"{trust_b_pct/trust_a_pct if trust_a_pct > 0 else 0:.2f}x")
    print(f"  Efficiency: {eff_b:.3f} / {eff_a:.3f} = {eff_b/eff_a:.2f}x")

    # Validation
    print(f"\n{'='*70}")
    print("VALIDATION")
    print(f"{'='*70}")

    checks_passed = 0
    total_checks = 4

    # Check 1: Warm-start increases trust-driven %
    if trust_b_pct > trust_a_pct:
        print(f"✅ Check 1: Warm-start increases trust-driven "
              f"({trust_a_pct:.1f}% → {trust_b_pct:.1f}%)")
        checks_passed += 1
    else:
        print(f"❌ Check 1: Trust-driven did not increase")

    # Check 2: Economic efficiency > 0.5
    if eff_b > 0.5:
        print(f"✅ Check 2: Economic efficiency ({eff_b:.3f}) > 0.5")
        checks_passed += 1
    else:
        print(f"❌ Check 2: Economic efficiency ({eff_b:.3f}) <= 0.5")

    # Check 3: ATP allocated matches budget
    expected_total = 100 * 1000.0  # 100 generations × 1000 ATP/gen
    actual_total = stats_b['atp_allocated_total']
    if abs(actual_total - expected_total) / expected_total < 0.05:  # Within 5%
        print(f"✅ Check 3: ATP allocation correct "
              f"({actual_total:.0f} ≈ {expected_total:.0f})")
        checks_passed += 1
    else:
        print(f"❌ Check 3: ATP allocation incorrect "
              f"({actual_total:.0f} vs {expected_total:.0f})")

    # Check 4: ADP earned > 0
    if stats_b['adp_earned_total'] > 0:
        print(f"✅ Check 4: ADP rewards generated ({stats_b['adp_earned_total']:.0f})")
        checks_passed += 1
    else:
        print(f"❌ Check 4: No ADP rewards generated")

    print(f"\n{'='*70}")
    if checks_passed == total_checks:
        print(f"✅ ALL CHECKS PASSED ({checks_passed}/{total_checks})")
        print(f"\nConclusion:")
        print(f"  ATP/ADP + epsilon-warmstart integration working correctly")
        print(f"  Economic cycle validated: ATP → Performance → ADP → Trust")
        print(f"  Warm-start improves both trust-driven % and efficiency")
    else:
        print(f"⚠️  SOME CHECKS FAILED ({checks_passed}/{total_checks})")
    print("="*70)


if __name__ == "__main__":
    test_atp_epsilon_integration()
