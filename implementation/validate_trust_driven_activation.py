#!/usr/bin/env python3
"""
Session 72 Track 1: Validate Trust-Driven Activation

Verifies that trust_driven mode activates correctly when using
UNWEIGHTED quality (fixing Thor Sessions 77-78 bug).

Bug discovered by Thor Session 79:
- Sessions 77-78 stored weighted_quality = quality × weight ≈ 0.19
- Threshold check: 0.19 < low_trust_threshold (0.3) → FAILED
- Fix: Store unweighted quality ≈ 0.75 > 0.3 → PASSES

This test validates that Legion's implementation (which uses unweighted
quality) correctly activates trust_driven mode.

Created: 2025-12-19 (Legion Session 72)
"""

import sys
import numpy as np
from pathlib import Path
from typing import Dict, List

sys.path.insert(0, str(Path(__file__).parent))
from epsilon_warmstart_integration import EpsilonWarmStartSelector


def test_trust_driven_activation():
    """
    Test that trust_driven mode activates after sufficient evidence.

    Expected behavior:
    1. Bootstrap (gen 0-30): router_explore dominant
    2. Evidence accumulation (gen 10-40): forced_exploration gathers diversity
    3. Trust_driven activation (gen 30-50): trust_driven should activate
    """
    print("\n" + "="*70)
    print("TRUST-DRIVEN ACTIVATION VALIDATION")
    print("="*70)
    print("\nValidating fix for Thor Sessions 77-78 bug:")
    print("  Bug: Stored weighted_quality ≈ 0.19 < threshold (0.3)")
    print("  Fix: Store unweighted quality ≈ 0.75 > threshold (0.3)")
    print("  Expected: trust_driven activates after ~30-40 generations")
    print("")

    # Configuration
    selector = EpsilonWarmStartSelector(
        num_experts=128,
        min_evidence_threshold=3,  # Same as Thor S78
        epsilon=0.2,  # Same as Thor S77-78
        trust_ledger=None,
        society_id=None
    )

    contexts = ["context_0", "context_1", "context_2"]

    # Track first trust_driven activation
    first_trust_driven = None
    generation_modes = []
    expert_usage = {}

    # Track trust values to verify they're > threshold
    trust_values_log = []

    print("Running 100 generations...")
    print("=" * 70)

    for gen in range(100):
        context = np.random.choice(contexts)
        router_logits = np.random.randn(128).astype(np.float32)

        # Select
        experts, mode = selector.select_experts(router_logits, context, k=4)
        generation_modes.append(mode)

        # Track first trust_driven
        if mode == "trust_driven" and first_trust_driven is None:
            first_trust_driven = gen
            print(f"\n✅ FIRST trust_driven activation at generation {gen}")

        # Simulate quality (same as Sessions 77-78)
        base_quality = 0.75 + np.random.normal(0, 0.1)
        quality = float(np.clip(base_quality, 0, 1))

        # CRITICAL: Update with UNWEIGHTED quality (Legion implementation)
        selector.update_trust(experts, context, quality)

        # Log trust values every 10 generations
        if gen % 10 == 0:
            trust_snapshot = {}
            for expert_id in range(128):
                if expert_id in selector.expert_trust:
                    for ctx, trust_val in selector.expert_trust[expert_id].items():
                        if ctx == context:
                            trust_snapshot[expert_id] = trust_val

            if trust_snapshot:
                avg_trust = np.mean(list(trust_snapshot.values()))
                max_trust = np.max(list(trust_snapshot.values()))
                trust_values_log.append({
                    "generation": gen,
                    "avg_trust": avg_trust,
                    "max_trust": max_trust,
                    "num_experts": len(trust_snapshot)
                })

        # Track usage
        for expert_id in experts:
            expert_usage[expert_id] = expert_usage.get(expert_id, 0) + 1

        # Print progress every 20 generations
        if (gen + 1) % 20 == 0:
            mode_dist = {}
            for m in generation_modes:
                mode_dist[m] = mode_dist.get(m, 0) + 1

            trust_pct = 100 * mode_dist.get("trust_driven", 0) / (gen + 1)
            print(f"Gen {gen+1:3d}: {len(expert_usage)} experts, "
                  f"{trust_pct:.1f}% trust_driven")

    # Final statistics
    print("\n" + "="*70)
    print("RESULTS")
    print("="*70)

    stats = selector.get_statistics()

    # Mode distribution
    mode_counts = {mode: 0 for mode in ["trust_driven", "router_explore", "forced_exploration"]}
    for mode in generation_modes:
        mode_counts[mode] += 1

    print(f"\n📊 Expert Diversity:")
    print(f"  Unique experts: {len(expert_usage)}/128 ({100*len(expert_usage)/128:.1f}%)")
    print(f"  Total selections: {sum(expert_usage.values())}")

    print(f"\n🔄 Mode Distribution:")
    for mode, count in mode_counts.items():
        pct = 100 * count / 100
        print(f"  {mode}: {count}/100 ({pct:.1f}%)")

    print(f"\n🎯 Trust-Driven Activation:")
    if first_trust_driven is not None:
        print(f"  ✅ First activation: Generation {first_trust_driven}")
        print(f"  ✅ Total trust_driven: {mode_counts['trust_driven']}/100 "
              f"({100*mode_counts['trust_driven']/100:.1f}%)")
    else:
        print(f"  ❌ FAILED: trust_driven never activated")
        print(f"  This indicates the bug is still present!")

    print(f"\n📈 Trust Value Evolution:")
    for log_entry in trust_values_log:
        print(f"  Gen {log_entry['generation']:3d}: "
              f"avg={log_entry['avg_trust']:.3f}, "
              f"max={log_entry['max_trust']:.3f}, "
              f"experts={log_entry['num_experts']}")

    # Validation checks
    print(f"\n" + "="*70)
    print("VALIDATION")
    print("="*70)

    checks_passed = 0
    total_checks = 4

    # Check 1: Trust values > threshold
    if trust_values_log:
        final_avg = trust_values_log[-1]['avg_trust']
        final_max = trust_values_log[-1]['max_trust']

        if final_avg > 0.3:
            print(f"✅ Check 1: Average trust ({final_avg:.3f}) > threshold (0.3)")
            checks_passed += 1
        else:
            print(f"❌ Check 1: Average trust ({final_avg:.3f}) <= threshold (0.3)")
    else:
        print(f"❌ Check 1: No trust values logged")

    # Check 2: Trust_driven activated
    if first_trust_driven is not None:
        print(f"✅ Check 2: trust_driven activated at generation {first_trust_driven}")
        checks_passed += 1
    else:
        print(f"❌ Check 2: trust_driven never activated")

    # Check 3: Trust_driven percentage > 20%
    trust_pct = 100 * mode_counts['trust_driven'] / 100
    if trust_pct > 20:
        print(f"✅ Check 3: trust_driven percentage ({trust_pct:.1f}%) > 20%")
        checks_passed += 1
    else:
        print(f"❌ Check 3: trust_driven percentage ({trust_pct:.1f}%) <= 20%")

    # Check 4: Expert diversity > 50
    if len(expert_usage) > 50:
        print(f"✅ Check 4: Expert diversity ({len(expert_usage)}) > 50")
        checks_passed += 1
    else:
        print(f"❌ Check 4: Expert diversity ({len(expert_usage)}) <= 50")

    # Final verdict
    print(f"\n" + "="*70)
    if checks_passed == total_checks:
        print(f"✅ ALL CHECKS PASSED ({checks_passed}/{total_checks})")
        print(f"\nConclusion:")
        print(f"  Legion implementation is CORRECT - uses unweighted quality")
        print(f"  Trust_driven mode activates as expected")
        print(f"  Thor's bug (weighted_quality) is NOT present in Legion")
    else:
        print(f"⚠️  SOME CHECKS FAILED ({checks_passed}/{total_checks})")
        print(f"\nConclusion:")
        print(f"  Investigation needed - unexpected behavior")
    print("="*70)


def compare_weighted_vs_unweighted():
    """
    Demonstrate the bug: weighted vs unweighted quality storage.
    """
    print("\n" + "="*70)
    print("BUG DEMONSTRATION: Weighted vs Unweighted Quality")
    print("="*70)

    print("\nScenario: k=4 experts selected, quality=0.75")
    print("")

    quality = 0.75
    k = 4
    weight = 1.0 / k  # Uniform weighting

    weighted_quality = quality * weight
    threshold = 0.3

    print(f"Unweighted approach (CORRECT - Legion):")
    print(f"  quality = {quality}")
    print(f"  threshold = {threshold}")
    print(f"  Check: {quality} > {threshold} → {'✅ PASSES' if quality > threshold else '❌ FAILS'}")
    print("")

    print(f"Weighted approach (BUG - Thor S77-78):")
    print(f"  quality = {quality}")
    print(f"  weight = 1/{k} = {weight}")
    print(f"  weighted_quality = {quality} × {weight} = {weighted_quality}")
    print(f"  threshold = {threshold}")
    print(f"  Check: {weighted_quality} > {threshold} → {'✅ PASSES' if weighted_quality > threshold else '❌ FAILS'}")
    print("")

    print(f"Result:")
    print(f"  Unweighted: trust_driven CAN activate ✅")
    print(f"  Weighted: trust_driven CANNOT activate ❌")
    print(f"\nThis explains why Thor S77-78 had 0% trust_driven!")
    print("="*70)


if __name__ == "__main__":
    # First demonstrate the bug
    compare_weighted_vs_unweighted()

    print("\n\n")

    # Then validate Legion implementation is correct
    test_trust_driven_activation()
