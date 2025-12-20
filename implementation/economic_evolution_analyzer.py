#!/usr/bin/env python3
"""
Session 72 Track 3: Cross-Session Economic Evolution Analyzer

Analyzes trust + ATP/ADP evolution across multiple sessions with warm-start.

Demonstrates complete lifecycle:
1. Session 1 (cold, high ε): Bootstrap diversity, gather evidence
2. Session 2 (warm, med ε): Leverage trust, build specialists
3. Session 3 (warm, low ε): Exploit specialists, maximize efficiency

Economic tracking:
- ATP allocation evolution
- ADP reward trends
- Trust growth patterns
- Specialist emergence
- Economic efficiency optimization

Created: 2025-12-19 (Legion Session 72)
"""

import sys
import numpy as np
from pathlib import Path
from typing import Dict, List
import json

sys.path.insert(0, str(Path(__file__).parent))
from test_atp_epsilon_integration import ATPEpsilonWarmStartSelector
from trust_ledger_persistence import TrustLedgerPersistence


def analyze_economic_evolution():
    """
    Run 3-session economic evolution analysis.

    Session 1: Cold start, ε=0.3 (maximum exploration)
    Session 2: Warm start, ε=0.2 (moderate exploration)
    Session 3: Warm start, ε=0.1 (maximum exploitation)
    """
    print("\n" + "="*70)
    print("CROSS-SESSION ECONOMIC EVOLUTION ANALYSIS")
    print("="*70)
    print("\nLifecycle simulation:")
    print("  Session 1 (cold, ε=0.3): Bootstrap diversity")
    print("  Session 2 (warm, ε=0.2): Build specialists")
    print("  Session 3 (warm, ε=0.1): Exploit specialists")
    print("")

    # Setup ledger
    ledger_dir = Path("/tmp/economic_evolution")
    ledger = TrustLedgerPersistence(ledger_dir)

    contexts = ["context_0", "context_1", "context_2"]
    session_results = []

    # Define session configs
    sessions = [
        {"name": "Session 1", "epsilon": 0.3, "warm": False, "generations": 100},
        {"name": "Session 2", "epsilon": 0.2, "warm": True, "generations": 100},
        {"name": "Session 3", "epsilon": 0.1, "warm": True, "generations": 100},
    ]

    for session_config in sessions:
        print("="*70)
        print(f"{session_config['name']}: "
              f"{'WARM START' if session_config['warm'] else 'COLD START'} "
              f"(ε={session_config['epsilon']})")
        print("="*70)

        # Create selector
        selector = ATPEpsilonWarmStartSelector(
            num_experts=128,
            min_evidence_threshold=3,
            epsilon=session_config['epsilon'],
            total_atp_budget=1000.0,
            min_atp_per_expert=100.0,
            trust_ledger=ledger if session_config['warm'] else None,
            society_id="economic-evolution" if session_config['warm'] else None
        )

        print(f"\nConfiguration:")
        print(f"  Epsilon: {session_config['epsilon']}")
        print(f"  Warm-started: {selector.base_selector.warm_started}")
        print(f"  Generations: {session_config['generations']}")

        # Track evolution
        expert_usage = {}
        generation_data = []

        print(f"\nRunning {session_config['generations']} generations...")

        for gen in range(session_config['generations']):
            context = np.random.choice(contexts)
            router_logits = np.random.randn(128).astype(np.float32)

            # Select and allocate
            experts, mode, atp_allocations = selector.select_and_allocate(
                router_logits, context, k=4
            )

            # Simulate quality
            base_quality = 0.75 + np.random.normal(0, 0.1)
            quality = float(np.clip(base_quality, 0, 1))

            # Update trust with ADP
            selector.update_trust_with_adp(
                experts, atp_allocations, context, quality
            )

            # Track usage
            for expert_id in experts:
                expert_usage[expert_id] = expert_usage.get(expert_id, 0) + 1

            # Record generation data
            stats = selector.get_economic_statistics()
            generation_data.append({
                "generation": gen,
                "mode": mode,
                "unique_experts": len(expert_usage),
                "atp_allocated": selector.atp_allocated_total,
                "adp_earned": selector.adp_earned_total,
                "efficiency": stats['economic_efficiency'],
                "trust_driven_pct": stats['mode_distribution']['trust_driven']['percentage']
            })

            # Progress every 25 generations
            if (gen + 1) % 25 == 0:
                print(f"Gen {gen+1:3d}: {len(expert_usage)} experts, "
                      f"eff={stats['economic_efficiency']:.3f}, "
                      f"trust={stats['mode_distribution']['trust_driven']['percentage']:.1f}%")

        # Final statistics
        final_stats = selector.get_economic_statistics()

        # Analyze specialists
        specialist_count = 0
        generalist_count = 0

        for expert_id in expert_usage.keys():
            if expert_id in selector.base_selector.expert_trust:
                contexts_used = len(selector.base_selector.expert_trust[expert_id])
                if contexts_used == 1:
                    specialist_count += 1
                else:
                    generalist_count += 1

        specialist_rate = specialist_count / (specialist_count + generalist_count) \
                         if (specialist_count + generalist_count) > 0 else 0

        session_result = {
            "name": session_config['name'],
            "epsilon": session_config['epsilon'],
            "warm_started": selector.base_selector.warm_started,
            "unique_experts": len(expert_usage),
            "utilization": 100 * len(expert_usage) / 128,
            "specialist_count": specialist_count,
            "generalist_count": generalist_count,
            "specialization_rate": specialist_rate,
            "atp_allocated": final_stats['atp_allocated_total'],
            "adp_earned": final_stats['adp_earned_total'],
            "efficiency": final_stats['economic_efficiency'],
            "trust_driven_pct": final_stats['mode_distribution']['trust_driven']['percentage'],
            "router_explore_pct": final_stats['mode_distribution']['router_explore']['percentage'],
            "forced_exp_pct": final_stats['mode_distribution']['forced_exploration']['percentage'],
            "generation_data": generation_data
        }

        session_results.append(session_result)

        print(f"\n{session_config['name']} Results:")
        print(f"  Experts: {len(expert_usage)}/128 ({100*len(expert_usage)/128:.1f}%)")
        print(f"  Specialists: {specialist_count} ({100*specialist_rate:.1f}%)")
        print(f"  ATP allocated: {final_stats['atp_allocated_total']:.0f}")
        print(f"  ADP earned: {final_stats['adp_earned_total']:.0f}")
        print(f"  Efficiency: {final_stats['economic_efficiency']:.3f}")
        print(f"  Trust-driven: {final_stats['mode_distribution']['trust_driven']['percentage']:.1f}%")

        # Save snapshot for next session
        if session_config['name'] != "Session 3":  # Don't save after last
            trust_state = {}
            observation_counts = {}
            for expert_id, contexts_dict in selector.base_selector.expert_trust.items():
                for ctx, trust_val in contexts_dict.items():
                    trust_state[(expert_id, ctx)] = trust_val
                    obs_count = selector.base_selector.expert_observations.get(
                        expert_id, {}
                    ).get(ctx, 0)
                    observation_counts[(expert_id, ctx)] = obs_count

            snapshot_id = ledger.save_snapshot(
                society_id="economic-evolution",
                session_id=session_config['name'],
                trust_state=trust_state,
                observation_counts=observation_counts,
                metadata={
                    "epsilon": session_config['epsilon'],
                    "efficiency": final_stats['economic_efficiency']
                }
            )
            print(f"  ✅ Snapshot saved: {snapshot_id}")

        print("")

    # Cross-session analysis
    print("="*70)
    print("CROSS-SESSION ANALYSIS")
    print("="*70)

    print(f"\n{'Session':<15} {'Experts':<10} {'Trust%':<10} {'Eff':<10} {'Spec%':<10}")
    print("-" * 70)

    for result in session_results:
        print(f"{result['name']:<15} "
              f"{result['unique_experts']:<10} "
              f"{result['trust_driven_pct']:<10.1f} "
              f"{result['efficiency']:<10.3f} "
              f"{100*result['specialization_rate']:<10.1f}")

    print(f"\n{'='*70}")
    print("EVOLUTION TRENDS")
    print(f"{'='*70}")

    # Diversity trend
    diversity_trend = [r['unique_experts'] for r in session_results]
    print(f"\n📊 Diversity Evolution:")
    print(f"  Session 1 → 2: {diversity_trend[0]} → {diversity_trend[1]} "
          f"({100*(diversity_trend[1]-diversity_trend[0])/diversity_trend[0]:+.1f}%)")
    print(f"  Session 2 → 3: {diversity_trend[1]} → {diversity_trend[2]} "
          f"({100*(diversity_trend[2]-diversity_trend[1])/diversity_trend[1]:+.1f}%)")

    # Efficiency trend
    eff_trend = [r['efficiency'] for r in session_results]
    print(f"\n💰 Efficiency Evolution:")
    print(f"  Session 1 → 2: {eff_trend[0]:.3f} → {eff_trend[1]:.3f} "
          f"({100*(eff_trend[1]-eff_trend[0])/eff_trend[0]:+.1f}%)")
    print(f"  Session 2 → 3: {eff_trend[1]:.3f} → {eff_trend[2]:.3f} "
          f"({100*(eff_trend[2]-eff_trend[1])/eff_trend[1]:+.1f}%)")

    # Trust-driven trend
    trust_trend = [r['trust_driven_pct'] for r in session_results]
    print(f"\n🎯 Trust-Driven Evolution:")
    print(f"  Session 1 → 2: {trust_trend[0]:.1f}% → {trust_trend[1]:.1f}% "
          f"({trust_trend[1]-trust_trend[0]:+.1f}pp)")
    print(f"  Session 2 → 3: {trust_trend[1]:.1f}% → {trust_trend[2]:.1f}% "
          f"({trust_trend[2]-trust_trend[1]:+.1f}pp)")

    # Specialization trend
    spec_trend = [100*r['specialization_rate'] for r in session_results]
    print(f"\n🏆 Specialization Evolution:")
    print(f"  Session 1 → 2: {spec_trend[0]:.1f}% → {spec_trend[1]:.1f}% "
          f"({spec_trend[1]-spec_trend[0]:+.1f}pp)")
    print(f"  Session 2 → 3: {spec_trend[1]:.1f}% → {spec_trend[2]:.1f}% "
          f"({spec_trend[2]-spec_trend[1]:+.1f}pp)")

    # Key insights
    print(f"\n{'='*70}")
    print("KEY INSIGHTS")
    print(f"{'='*70}")

    print(f"\n1. Diversity Management:")
    print(f"   High ε (Session 1): {diversity_trend[0]} experts (exploration)")
    print(f"   Medium ε (Session 2): {diversity_trend[1]} experts (balance)")
    print(f"   Low ε (Session 3): {diversity_trend[2]} experts (exploitation)")

    print(f"\n2. Economic Optimization:")
    print(f"   Session 1: {eff_trend[0]:.3f} (bootstrap)")
    print(f"   Session 2: {eff_trend[1]:.3f} (building)")
    print(f"   Session 3: {eff_trend[2]:.3f} (mature)")
    if eff_trend[2] > eff_trend[0]:
        print(f"   ✅ Efficiency improved {100*(eff_trend[2]-eff_trend[0])/eff_trend[0]:.1f}% overall")

    print(f"\n3. Trust-Driven Activation:")
    print(f"   Session 1: {trust_trend[0]:.1f}% (evidence gathering)")
    print(f"   Session 2: {trust_trend[1]:.1f}% (trust building)")
    print(f"   Session 3: {trust_trend[2]:.1f}% (trust leverage)")
    if trust_trend[2] > 80:
        print(f"   ✅ Mature system (>80% trust-driven)")

    print(f"\n4. Specialist Development:")
    print(f"   Session 1: {spec_trend[0]:.1f}% specialists")
    print(f"   Session 2: {spec_trend[1]:.1f}% specialists")
    print(f"   Session 3: {spec_trend[2]:.1f}% specialists")
    if spec_trend[2] > 50:
        print(f"   ✅ Healthy specialization (>50%)")

    # Production readiness assessment
    print(f"\n{'='*70}")
    print("PRODUCTION READINESS ASSESSMENT")
    print(f"{'='*70}")

    s3 = session_results[2]  # Session 3 (mature system)

    checks = []
    checks.append(("Efficiency > 0.7", s3['efficiency'] > 0.7))
    checks.append(("Trust-driven > 80%", s3['trust_driven_pct'] > 80))
    checks.append(("Specialization 40-60%", 40 < 100*s3['specialization_rate'] < 60))
    checks.append(("Diversity > 30 experts", s3['unique_experts'] > 30))

    passed = sum(1 for _, result in checks if result)
    total = len(checks)

    print(f"\nMature System (Session 3) Checks:")
    for check_name, result in checks:
        status = "✅" if result else "❌"
        print(f"  {status} {check_name}")

    print(f"\n{'='*70}")
    if passed == total:
        print(f"✅ PRODUCTION READY ({passed}/{total} checks passed)")
        print(f"\nRecommendation:")
        print(f"  System demonstrates stable, efficient evolution")
        print(f"  Ready for deployment with 3-session lifecycle:")
        print(f"    1. Bootstrap (ε=0.3, cold)")
        print(f"    2. Build (ε=0.2, warm)")
        print(f"    3. Production (ε=0.1, warm)")
    else:
        print(f"⚠️  NEEDS TUNING ({passed}/{total} checks passed)")
        print(f"\nRecommendation:")
        print(f"  Adjust parameters or extend bootstrap phase")
    print("="*70)

    # Save results
    results_file = Path("/tmp/economic_evolution_results.json")
    with open(results_file, 'w') as f:
        json.dump(session_results, f, indent=2)
    print(f"\n✅ Detailed results saved: {results_file}")


if __name__ == "__main__":
    analyze_economic_evolution()
