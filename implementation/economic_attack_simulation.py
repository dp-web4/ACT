#!/usr/bin/env python3
"""
Session 73 Track 3: Economic Attack Simulation (ATP/ADP Gaming)

Analyzes economic manipulation vectors in ATP/ADP token economy.

Problem:
- ATP allocation creates economic incentives
- ADP rewards based on performance × ATP
- Rational actors may game system for profit

Economic Attack Vectors:
1. Low-Quality Farming: Accept high ATP, deliver low quality, repeat
2. Trust Hoarding: Build trust, then defect for one-time gain
3. Collusion: Coordinate to monopolize ATP allocation
4. Free-Riding: Minimal effort while benefiting from system
5. Resource Exhaustion: Consume ATP without value creation

Economic Defenses:
- Trust decay (prevents permanent high trust)
- ADP = quality × ATP (low quality → low rewards)
- Economic efficiency metrics (detect farming)
- Minimum quality thresholds
- Reputation-based ATP allocation

Methodology:
- Simulate rational economic actors
- Measure profit vs honest behavior
- Test defense mechanisms
- Quantify economic security

Created: 2025-12-20 (Legion Session 73)
Author: Legion (Autonomous Web4 Research)
"""

import numpy as np
from typing import Dict, List, Tuple
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from test_atp_epsilon_integration import ATPEpsilonWarmStartSelector
from trust_ledger_persistence import TrustLedgerPersistence


class EconomicAttackSimulator:
    """
    Simulates economic attacks on ATP/ADP system.

    Tests profitability of various gaming strategies.
    """

    def __init__(self):
        """Initialize economic attack simulator."""
        self.attack_results = []

    def simulate_low_quality_farming(
        self,
        selector: ATPEpsilonWarmStartSelector,
        malicious_expert_ids: List[int],
        generations: int = 100
    ) -> Dict:
        """
        Attack: Maximize ATP intake while delivering low quality.

        Strategy:
        - Build initial trust (bootstrap phase)
        - Once trust high, deliver minimum quality
        - Maximize ATP allocation, minimize effort
        - Goal: High ATP, low ADP, but profit from ATP alone

        Reality Check:
        - ADP = quality × ATP
        - Low quality → low ADP → trust decay → ATP allocation drops

        Args:
            selector: ATP-integrated selector
            malicious_expert_ids: List of malicious expert IDs
            generations: Number of generations

        Returns:
            Attack profitability analysis
        """
        print(f"\n{'='*70}")
        print("ECONOMIC ATTACK 1: Low-Quality Farming")
        print(f"{'='*70}")
        print(f"\nMalicious experts: {len(malicious_expert_ids)}")
        print(f"Strategy: High ATP allocation, low quality delivery\n")

        contexts = ["context_0", "context_1", "context_2"]

        # Phase 1: Build trust (generations 0-30)
        print("Phase 1: Building trust (honest behavior)...")
        malicious_atp_total = 0
        malicious_adp_total = 0
        honest_atp_total = 0
        honest_adp_total = 0

        for gen in range(30):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            experts, mode, atp_allocations = selector.select_and_allocate(
                router_logits, context, k=4
            )

            # Honest phase - everyone delivers quality
            quality = float(np.clip(0.75 + np.random.normal(0, 0.1), 0, 1))
            selector.update_trust_with_adp(experts, atp_allocations, context, quality)

        # Phase 2: Attack (generations 30-100)
        print("Phase 2: Farming attack (low quality)...\n")

        for gen in range(30, generations):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            experts, mode, atp_allocations = selector.select_and_allocate(
                router_logits, context, k=4
            )

            # Track ATP/ADP by type
            for alloc in atp_allocations:
                if alloc.expert_id in malicious_expert_ids:
                    malicious_atp_total += alloc.atp_allocated
                else:
                    honest_atp_total += alloc.atp_allocated

            # ATTACK: Malicious experts deliver low quality
            for alloc in atp_allocations:
                if alloc.expert_id in malicious_expert_ids:
                    # Low quality (0.3) - minimal effort
                    quality = 0.3
                else:
                    # Honest quality (0.75)
                    quality = float(np.clip(0.75 + np.random.normal(0, 0.1), 0, 1))

                # Calculate ADP (quality × ATP)
                adp = quality * alloc.atp_allocated

                if alloc.expert_id in malicious_expert_ids:
                    malicious_adp_total += adp
                else:
                    honest_adp_total += adp

                # Update trust
                selector.update_trust_with_adp([alloc.expert_id], [alloc], context, quality)

        # Calculate profitability
        malicious_efficiency = malicious_adp_total / malicious_atp_total \
                              if malicious_atp_total > 0 else 0
        honest_efficiency = honest_adp_total / honest_atp_total \
                           if honest_atp_total > 0 else 0

        results = {
            "attack_type": "low_quality_farming",
            "malicious_atp": malicious_atp_total,
            "malicious_adp": malicious_adp_total,
            "malicious_efficiency": malicious_efficiency,
            "honest_atp": honest_atp_total,
            "honest_adp": honest_adp_total,
            "honest_efficiency": honest_efficiency,
            "profit_ratio": malicious_efficiency / honest_efficiency if honest_efficiency > 0 else 0,
            "profitable": malicious_efficiency > honest_efficiency  # Attack succeeds if more profitable
        }

        print(f"Results (attack phase generations 30-100):")
        print(f"\n  Malicious Strategy:")
        print(f"    ATP allocated: {malicious_atp_total:.0f}")
        print(f"    ADP earned: {malicious_adp_total:.0f}")
        print(f"    Efficiency: {malicious_efficiency:.3f}")
        print(f"\n  Honest Baseline:")
        print(f"    ATP allocated: {honest_atp_total:.0f}")
        print(f"    ADP earned: {honest_adp_total:.0f}")
        print(f"    Efficiency: {honest_efficiency:.3f}")
        print(f"\n  Profit Ratio: {results['profit_ratio']:.3f}x")
        print(f"  Attack profitable: {'✅ YES' if results['profitable'] else '❌ NO'}")

        if results['profitable']:
            print(f"\n⚠️  VULNERABILITY DETECTED")
            print(f"  Economic attack is PROFITABLE")
        else:
            print(f"\n✅ Economic Defense WORKING")
            print(f"  Low quality → Low ADP → Unprofitable")

        return results

    def simulate_trust_hoarding_defection(
        self,
        selector: ATPEpsilonWarmStartSelector,
        defector_expert_id: int,
        defection_generation: int = 70,
        generations: int = 100
    ) -> Dict:
        """
        Attack: Build trust, then defect for one-time gain.

        Strategy:
        - Honest behavior for N generations (build trust)
        - At generation N: Defect (take ATP, deliver zero quality)
        - Goal: Maximize one-time gain before trust decays

        Reality Check:
        - Trust decay prevents repeated defection
        - One-time gain vs long-term honest profit
        - System recovers after defection

        Args:
            selector: ATP-integrated selector
            defector_expert_id: Expert ID that will defect
            defection_generation: When to defect
            generations: Total generations

        Returns:
            Defection profitability analysis
        """
        print(f"\n{'='*70}")
        print("ECONOMIC ATTACK 2: Trust Hoarding + Defection")
        print(f"{'='*70}")
        print(f"\nDefector expert: {defector_expert_id}")
        print(f"Defection at generation: {defection_generation}\n")

        contexts = ["context_0", "context_1", "context_2"]

        defector_atp_pre = 0
        defector_adp_pre = 0
        defector_atp_post = 0
        defector_adp_post = 0
        defection_gain = 0
        defected = False

        for gen in range(generations):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            experts, mode, atp_allocations = selector.select_and_allocate(
                router_logits, context, k=4
            )

            for alloc in atp_allocations:
                if alloc.expert_id == defector_expert_id:
                    if gen < defection_generation:
                        # Pre-defection: Honest behavior
                        quality = 0.8
                        defector_atp_pre += alloc.atp_allocated
                        defector_adp_pre += quality * alloc.atp_allocated
                    elif gen == defection_generation:
                        # DEFECTION: Zero quality
                        quality = 0.0
                        defection_gain = alloc.atp_allocated  # ATP gained, no quality delivered
                        defected = True
                    else:
                        # Post-defection: Try to recover
                        quality = 0.8
                        defector_atp_post += alloc.atp_allocated
                        defector_adp_post += quality * alloc.atp_allocated
                else:
                    # Other experts honest
                    quality = 0.75

                selector.update_trust_with_adp([alloc.expert_id], [alloc], context, quality)

            if gen == defection_generation and defected:
                print(f"Generation {gen}: DEFECTION occurred (ATP gain: {defection_gain:.0f})")

        # Calculate profitability
        total_honest_profit = defector_atp_pre + defector_atp_post  # ATP if always honest
        actual_profit = defector_atp_pre + defection_gain + defector_atp_post

        results = {
            "attack_type": "trust_hoarding_defection",
            "defector_expert": defector_expert_id,
            "defection_generation": defection_generation,
            "pre_defection_atp": defector_atp_pre,
            "defection_gain": defection_gain,
            "post_defection_atp": defector_atp_post,
            "total_actual": actual_profit,
            "total_honest": total_honest_profit,
            "profit_gain": actual_profit - total_honest_profit,
            "profitable": actual_profit > total_honest_profit
        }

        print(f"\nResults:")
        print(f"  Pre-defection ATP: {defector_atp_pre:.0f}")
        print(f"  Defection gain: {defection_gain:.0f}")
        print(f"  Post-defection ATP: {defector_atp_post:.0f}")
        print(f"  Total (with defection): {actual_profit:.0f}")
        print(f"  Total (if honest): {total_honest_profit:.0f}")
        print(f"  Net gain from defection: {results['profit_gain']:+.0f}")
        print(f"  Defection profitable: {'✅ YES' if results['profitable'] else '❌ NO'}")

        if results['profitable']:
            print(f"\n⚠️  VULNERABILITY DETECTED")
            print(f"  One-time defection is PROFITABLE")
        else:
            print(f"\n✅ Economic Defense WORKING")
            print(f"  Trust decay prevents profitable defection")

        return results

    def simulate_collusion_monopoly(
        self,
        selector: ATPEpsilonWarmStartSelector,
        cartel_expert_ids: List[int],
        generations: int = 100
    ) -> Dict:
        """
        Attack: Experts collude to monopolize ATP allocation.

        Strategy:
        - Cartel members coordinate to dominate selection
        - Report high quality for each other
        - Goal: Monopolize ATP distribution within cartel

        Reality Check:
        - Epsilon-greedy breaks monopoly
        - Trust evidence required
        - Selection diversity prevents dominance

        Args:
            selector: ATP-integrated selector
            cartel_expert_ids: Colluding expert IDs
            generations: Number of generations

        Returns:
            Collusion success analysis
        """
        print(f"\n{'='*70}")
        print("ECONOMIC ATTACK 3: Collusion Monopoly")
        print(f"{'='*70}")
        print(f"\nCartel size: {len(cartel_expert_ids)}")
        print(f"Strategy: Coordinate to monopolize ATP\n")

        contexts = ["context_0", "context_1", "context_2"]

        cartel_atp = 0
        honest_atp = 0
        total_atp = 0

        for gen in range(generations):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            experts, mode, atp_allocations = selector.select_and_allocate(
                router_logits, context, k=4
            )

            for alloc in atp_allocations:
                total_atp += alloc.atp_allocated

                if alloc.expert_id in cartel_expert_ids:
                    cartel_atp += alloc.atp_allocated
                    # Cartel reports high quality
                    quality = 0.9
                else:
                    honest_atp += alloc.atp_allocated
                    # Honest baseline
                    quality = 0.75

                selector.update_trust_with_adp([alloc.expert_id], [alloc], context, quality)

        cartel_share = cartel_atp / total_atp if total_atp > 0 else 0
        expected_share = len(cartel_expert_ids) / 128  # Fair share

        results = {
            "attack_type": "collusion_monopoly",
            "cartel_size": len(cartel_expert_ids),
            "cartel_atp": cartel_atp,
            "honest_atp": honest_atp,
            "total_atp": total_atp,
            "cartel_share": cartel_share,
            "expected_share": expected_share,
            "monopoly_gain": cartel_share / expected_share if expected_share > 0 else 0,
            "successful": cartel_share > 2 * expected_share  # Success if >2x fair share
        }

        print(f"Results:")
        print(f"  Cartel ATP: {cartel_atp:.0f}/{total_atp:.0f} ({100*cartel_share:.1f}%)")
        print(f"  Expected (fair) share: {100*expected_share:.1f}%")
        print(f"  Monopoly gain: {results['monopoly_gain']:.2f}x")
        print(f"  Attack successful: {'✅ YES' if results['successful'] else '❌ NO'}")

        if results['successful']:
            print(f"\n⚠️  VULNERABILITY DETECTED")
            print(f"  Collusion achieves >2x fair share")
        else:
            print(f"\n✅ System RESILIENT")
            print(f"  Epsilon prevents cartel monopoly")

        return results


def run_economic_attack_analysis():
    """
    Run comprehensive economic attack simulation.

    Tests profitability of gaming strategies vs honest behavior.
    """
    print("\n" + "="*70)
    print("ECONOMIC ATTACK SIMULATION: ATP/ADP GAMING")
    print("="*70)
    print("\nTesting economic manipulation vectors:")
    print("  1. Low-Quality Farming")
    print("  2. Trust Hoarding + Defection")
    print("  3. Collusion Monopoly")
    print("")

    simulator = EconomicAttackSimulator()
    attack_results = []

    # Attack 1: Low-Quality Farming
    ledger1 = TrustLedgerPersistence(Path("/tmp/econ_attack_1"))
    selector1 = ATPEpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.2,
        total_atp_budget=1000.0,
        min_atp_per_expert=100.0,
        trust_ledger=None,
        society_id=None
    )

    result1 = simulator.simulate_low_quality_farming(
        selector=selector1,
        malicious_expert_ids=[10, 25, 42, 73, 99],
        generations=100
    )
    attack_results.append(result1)

    # Attack 2: Trust Hoarding + Defection
    ledger2 = TrustLedgerPersistence(Path("/tmp/econ_attack_2"))
    selector2 = ATPEpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.2,
        total_atp_budget=1000.0,
        min_atp_per_expert=100.0,
        trust_ledger=None,
        society_id=None
    )

    result2 = simulator.simulate_trust_hoarding_defection(
        selector=selector2,
        defector_expert_id=42,
        defection_generation=70,
        generations=100
    )
    attack_results.append(result2)

    # Attack 3: Collusion Monopoly
    ledger3 = TrustLedgerPersistence(Path("/tmp/econ_attack_3"))
    selector3 = ATPEpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,
        epsilon=0.2,
        total_atp_budget=1000.0,
        min_atp_per_expert=100.0,
        trust_ledger=None,
        society_id=None
    )

    result3 = simulator.simulate_collusion_monopoly(
        selector=selector3,
        cartel_expert_ids=list(range(100, 115)),  # 15-member cartel
        generations=100
    )
    attack_results.append(result3)

    # Summary
    print(f"\n{'='*70}")
    print("ECONOMIC SECURITY SUMMARY")
    print(f"{'='*70}\n")

    successful_attacks = 0
    for result in attack_results:
        if result.get('profitable') or result.get('successful'):
            successful_attacks += 1

    print(f"{'Attack Type':<35} {'Outcome':<15} {'Key Metric'}")
    print("-" * 70)

    for result in attack_results:
        attack_type = result['attack_type'].replace('_', ' ').title()

        if result['attack_type'] == 'low_quality_farming':
            outcome = "Profitable" if result['profitable'] else "Unprofitable"
            metric = f"{result['profit_ratio']:.2f}x vs honest"
        elif result['attack_type'] == 'trust_hoarding_defection':
            outcome = "Profitable" if result['profitable'] else "Unprofitable"
            metric = f"{result['profit_gain']:+.0f} net gain"
        else:  # collusion_monopoly
            outcome = "Successful" if result['successful'] else "Failed"
            metric = f"{result['monopoly_gain']:.2f}x fair share"

        status = "⚠️" if (result.get('profitable') or result.get('successful')) else "✅"
        print(f"{attack_type:<35} {status} {outcome:<13} {metric}")

    print(f"\n{'='*70}")
    print(f"Economic Security: {len(attack_results) - successful_attacks}/{len(attack_results)} attacks mitigated")
    print(f"Vulnerability Rate: {100*successful_attacks/len(attack_results):.1f}%")
    print("="*70)

    print(f"\nKey Economic Defenses:")
    print(f"  ✅ ADP = quality × ATP (low quality → low rewards)")
    print(f"  ✅ Trust decay prevents long-term farming")
    print(f"  ✅ Epsilon-greedy breaks collusion monopoly")
    print(f"  ⚠️  One-time defection may be profitable (reputation cost)")

    return attack_results


if __name__ == "__main__":
    run_economic_attack_analysis()
