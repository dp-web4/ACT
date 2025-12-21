"""
Byzantine Quality Consensus - Session 77 Track 2

Addresses Session 76 discovery: Quality inflation attacks succeed 23.6%
because quality reporting is self-attested (malicious agents can lie).

Solution: Byzantine consensus for quality attestations
- Require 2-of-3 societies to agree on quality
- Use median quality from multiple observers
- Reject outliers (>20% deviation from median)
- HMAC-SHA256 signatures for authenticity

Architecture:
- Based on Session 73's Byzantine consensus (HMAC signatures)
- Based on Session 75's Trust Federation Protocol (quorum)
- Defends against Session 76's quality inflation attack

Test Scenario:
- 10 malicious agents report inflated quality (1.0 when actual is 0.5)
- 3 societies observe and report quality independently
- Byzantine consensus rejects inflated reports
- Expected: Attack success rate < 5% (vs 23.6% in Session 76)
"""

import random
import time
import statistics
import json
import hashlib
import hmac
from dataclasses import dataclass, asdict
from typing import Dict, List, Optional, Tuple
from collections import defaultdict


# ============================================================================
# BYZANTINE QUALITY ATTESTATION
# ============================================================================

@dataclass
class QualityAttestation:
    """Quality observation with Byzantine signature."""
    attestation_id: str
    observer_society: str  # Who observed this quality
    expert_id: int
    context: str
    quality: float  # Observed quality [0.0, 1.0]
    timestamp: float
    signature: str  # HMAC-SHA256


@dataclass
class QualityConsensusResult:
    """Result from Byzantine quality consensus."""
    expert_id: int
    context: str

    # Consensus quality
    consensus_quality: float
    consensus_method: str  # 'median', 'mean', 'rejected'

    # Individual attestations
    attestations: List[QualityAttestation]
    num_attestations: int

    # Outlier detection
    outliers_detected: List[str]  # Society IDs of outliers
    outlier_deviation_pct: List[float]  # % deviation from median

    # Validation
    all_signatures_valid: bool
    quorum_reached: bool  # At least 2 attestations

    # Attack detection
    potential_attack: bool  # True if outliers detected


class ByzantineQualityConsensus:
    """
    Byzantine consensus for quality attestations.

    Defends against quality inflation attacks by requiring multiple
    independent observers to agree on quality.

    Based on:
    - Session 73: Byzantine consensus (HMAC signatures)
    - Session 75: Trust federation (quorum verification)
    - Session 76: Attack vector discovery
    """

    def __init__(
        self,
        quorum_size: int = 2,
        outlier_threshold_pct: float = 20.0,
        consensus_method: str = 'median'
    ):
        """
        Args:
            quorum_size: Minimum attestations for consensus (default: 2)
            outlier_threshold_pct: % deviation to flag as outlier (default: 20%)
            consensus_method: 'median' or 'mean' (default: 'median')
        """
        self.quorum_size = quorum_size
        self.outlier_threshold_pct = outlier_threshold_pct
        self.consensus_method = consensus_method

        # Society registry (for signature verification)
        self.known_societies: Dict[str, str] = {}  # society_id → public_key

        # Stats
        self.stats = {
            'total_consensus_attempts': 0,
            'quorum_reached': 0,
            'quorum_failed': 0,
            'outliers_detected': 0,
            'attacks_detected': 0,
            'signatures_verified': 0,
            'signatures_rejected': 0
        }

    def register_society(self, society_id: str, public_key: str):
        """Register known society for verification."""
        self.known_societies[society_id] = public_key

    def create_attestation(
        self,
        observer_society: str,
        secret_key: str,
        expert_id: int,
        context: str,
        quality: float
    ) -> QualityAttestation:
        """Create signed quality attestation."""
        attestation_id = f"{observer_society}-{expert_id}-{context}-{time.time()}"

        # Create attestation data
        attestation_data = (
            f"{attestation_id}|{observer_society}|{expert_id}|"
            f"{context}|{quality:.6f}"
        )

        # Sign with HMAC-SHA256
        signature = hmac.new(
            secret_key.encode(),
            attestation_data.encode(),
            hashlib.sha256
        ).hexdigest()

        return QualityAttestation(
            attestation_id=attestation_id,
            observer_society=observer_society,
            expert_id=expert_id,
            context=context,
            quality=quality,
            timestamp=time.time(),
            signature=signature
        )

    def verify_attestation(
        self,
        attestation: QualityAttestation
    ) -> bool:
        """Verify attestation signature."""
        # Get public key for observer society
        public_key = self.known_societies.get(attestation.observer_society)
        if not public_key:
            self.stats['signatures_rejected'] += 1
            return False

        # Reconstruct attestation data
        attestation_data = (
            f"{attestation.attestation_id}|{attestation.observer_society}|"
            f"{attestation.expert_id}|{attestation.context}|"
            f"{attestation.quality:.6f}"
        )

        # Verify signature
        expected_signature = hmac.new(
            public_key.encode(),
            attestation_data.encode(),
            hashlib.sha256
        ).hexdigest()

        is_valid = expected_signature == attestation.signature

        if is_valid:
            self.stats['signatures_verified'] += 1
        else:
            self.stats['signatures_rejected'] += 1

        return is_valid

    def compute_consensus(
        self,
        attestations: List[QualityAttestation]
    ) -> QualityConsensusResult:
        """
        Compute Byzantine consensus on quality from multiple attestations.

        Returns consensus quality with outlier detection.
        """
        self.stats['total_consensus_attempts'] += 1

        expert_id = attestations[0].expert_id
        context = attestations[0].context

        # Verify all signatures
        valid_attestations = []
        for attestation in attestations:
            if self.verify_attestation(attestation):
                valid_attestations.append(attestation)

        all_signatures_valid = len(valid_attestations) == len(attestations)

        # Check quorum
        quorum_reached = len(valid_attestations) >= self.quorum_size

        if not quorum_reached:
            self.stats['quorum_failed'] += 1
            return QualityConsensusResult(
                expert_id=expert_id,
                context=context,
                consensus_quality=0.0,
                consensus_method='rejected',
                attestations=attestations,
                num_attestations=len(attestations),
                outliers_detected=[],
                outlier_deviation_pct=[],
                all_signatures_valid=all_signatures_valid,
                quorum_reached=False,
                potential_attack=False
            )

        self.stats['quorum_reached'] += 1

        # Extract quality values
        qualities = [a.quality for a in valid_attestations]

        # Compute median (robust to outliers)
        median_quality = statistics.median(qualities)

        # Detect outliers (> outlier_threshold_pct deviation from median)
        outliers_detected = []
        outlier_deviation_pct = []

        for attestation in valid_attestations:
            if median_quality > 0:
                deviation_pct = abs(attestation.quality - median_quality) / median_quality * 100
            else:
                deviation_pct = abs(attestation.quality - median_quality) * 100

            if deviation_pct > self.outlier_threshold_pct:
                outliers_detected.append(attestation.observer_society)
                outlier_deviation_pct.append(deviation_pct)
                self.stats['outliers_detected'] += 1

        # Detect potential attack
        potential_attack = len(outliers_detected) > 0

        if potential_attack:
            self.stats['attacks_detected'] += 1

        # Compute consensus quality
        if self.consensus_method == 'median':
            consensus_quality = median_quality
        elif self.consensus_method == 'mean':
            # Filter outliers before computing mean
            filtered_qualities = [
                a.quality for a in valid_attestations
                if a.observer_society not in outliers_detected
            ]
            consensus_quality = statistics.mean(filtered_qualities) if filtered_qualities else median_quality
        else:
            consensus_quality = median_quality

        return QualityConsensusResult(
            expert_id=expert_id,
            context=context,
            consensus_quality=consensus_quality,
            consensus_method=self.consensus_method,
            attestations=attestations,
            num_attestations=len(attestations),
            outliers_detected=outliers_detected,
            outlier_deviation_pct=outlier_deviation_pct,
            all_signatures_valid=all_signatures_valid,
            quorum_reached=True,
            potential_attack=potential_attack
        )


# ============================================================================
# ATTACK SCENARIO TESTING
# ============================================================================

@dataclass
class QualityInflationAttackResult:
    """Result from quality inflation attack test."""
    test_id: str
    num_attackers: int
    num_generations: int

    # Attack success metrics
    attack_success_rate: float  # % of attacks that succeeded
    avg_quality_inflation: float  # Average quality inflation achieved

    # Defense metrics
    attacks_detected: int
    attacks_prevented: int
    false_positives: int  # Honest agents flagged as attackers

    # Consensus stats
    total_consensus_attempts: int
    quorum_reached_pct: float
    outliers_detected_pct: float

    # Comparison to Session 76 baseline
    session76_attack_success: float  # 23.6% baseline
    improvement_pct: float  # Improvement over Session 76

    passed: bool  # True if attack_success_rate < 5%


class QualityInflationAttackTester:
    """
    Tests quality inflation attack with Byzantine consensus defense.

    Scenario (from Session 76):
    - 10 malicious agents report inflated quality (1.0 when actual is 0.5)
    - Goal: Gain selection advantage through false quality reports

    Defense (Session 77):
    - 3 societies observe quality independently
    - Byzantine consensus rejects inflated reports
    - Median quality used (robust to outliers)
    """

    def __init__(self):
        self.num_experts = 128
        self.num_societies = 3  # Thor, Legion, Sprout

        # Society secrets (for signing)
        self.societies = {
            'thor': hashlib.sha256(b'web4-society-thor').hexdigest(),
            'legion': hashlib.sha256(b'web4-society-legion').hexdigest(),
            'sprout': hashlib.sha256(b'web4-society-sprout').hexdigest()
        }

    def run_attack_test(
        self,
        num_attackers: int = 10,
        num_generations: int = 90
    ) -> QualityInflationAttackResult:
        """
        Run quality inflation attack test with Byzantine consensus.

        Args:
            num_attackers: Number of malicious agents (default: 10)
            num_generations: Number of test generations (default: 90)
        """
        # Initialize Byzantine consensus
        consensus = ByzantineQualityConsensus(
            quorum_size=2,  # 2-of-3
            outlier_threshold_pct=20.0,  # 20% deviation threshold
            consensus_method='median'
        )

        # Register societies
        for society_id, secret_key in self.societies.items():
            consensus.register_society(society_id, secret_key)

        # Select attackers
        random.seed(42)
        attackers = set(random.sample(range(self.num_experts), num_attackers))

        # Track results
        attack_attempts = 0
        attacks_succeeded = 0
        quality_inflation_achieved = []

        attacks_detected = 0
        attacks_prevented = 0
        false_positives = 0

        # Run test
        for gen in range(num_generations):
            # Select expert for this generation
            expert_id = gen % self.num_experts

            # Determine true quality
            is_attacker = expert_id in attackers
            true_quality = 0.5  # All experts have same true quality

            # Each society observes quality independently
            # Attackers try to inflate quality in one society
            attestations = []

            for society_id, secret_key in self.societies.items():
                # Honest societies observe true quality (with noise)
                observed_quality = true_quality + random.uniform(-0.05, 0.05)
                observed_quality = max(0.0, min(1.0, observed_quality))

                # Attacker compromises one society (Thor) and inflates quality
                if is_attacker and society_id == 'thor':
                    observed_quality = 1.0  # Inflated quality
                    attack_attempts += 1

                # Create attestation
                attestation = consensus.create_attestation(
                    observer_society=society_id,
                    secret_key=secret_key,
                    expert_id=expert_id,
                    context=f"cluster_{gen % 9}",
                    quality=observed_quality
                )
                attestations.append(attestation)

            # Compute consensus
            result = consensus.compute_consensus(attestations)

            # Check if attack succeeded
            if is_attacker:
                # Attack succeeds if consensus quality > true_quality + 0.2
                inflation = result.consensus_quality - true_quality
                quality_inflation_achieved.append(inflation)

                if inflation > 0.2:
                    attacks_succeeded += 1
                else:
                    attacks_prevented += 1

                # Check if attack was detected
                if result.potential_attack:
                    attacks_detected += 1

            else:
                # False positive: honest agent flagged as attacker
                if result.potential_attack:
                    false_positives += 1

        # Calculate metrics
        attack_success_rate = attacks_succeeded / attack_attempts if attack_attempts > 0 else 0.0
        avg_quality_inflation = statistics.mean(quality_inflation_achieved) if quality_inflation_achieved else 0.0

        quorum_reached_pct = consensus.stats['quorum_reached'] / consensus.stats['total_consensus_attempts'] * 100
        outliers_detected_pct = consensus.stats['outliers_detected'] / (consensus.stats['total_consensus_attempts'] * 3) * 100

        # Comparison to Session 76
        session76_attack_success = 0.236  # 23.6% from Session 76
        improvement_pct = ((session76_attack_success - attack_success_rate) / session76_attack_success) * 100

        # Test passes if attack success < 5%
        passed = attack_success_rate < 0.05

        return QualityInflationAttackResult(
            test_id="quality-inflation-byzantine-v1",
            num_attackers=num_attackers,
            num_generations=num_generations,
            attack_success_rate=attack_success_rate,
            avg_quality_inflation=avg_quality_inflation,
            attacks_detected=attacks_detected,
            attacks_prevented=attacks_prevented,
            false_positives=false_positives,
            total_consensus_attempts=consensus.stats['total_consensus_attempts'],
            quorum_reached_pct=quorum_reached_pct,
            outliers_detected_pct=outliers_detected_pct,
            session76_attack_success=session76_attack_success,
            improvement_pct=improvement_pct,
            passed=passed
        )


# ============================================================================
# DEMO
# ============================================================================

def demo_byzantine_quality_consensus():
    """
    Demo: Byzantine quality consensus defense against inflation attacks.

    Validates that Byzantine consensus prevents quality inflation attacks
    that succeeded 23.6% of the time in Session 76.
    """
    print("=" * 80)
    print("BYZANTINE QUALITY CONSENSUS - Session 77 Track 2")
    print("=" * 80)
    print()
    print("Problem (Session 76):")
    print("  Quality inflation attacks succeed 23.6% of the time")
    print("  Root cause: Self-attested quality (malicious agents can lie)")
    print()
    print("Solution (Session 77):")
    print("  Byzantine consensus for quality attestations")
    print("  - 2-of-3 societies must agree on quality")
    print("  - Median quality (robust to outliers)")
    print("  - Reject >20% deviation from median")
    print()
    print("Test Scenario:")
    print("  - 10 malicious agents (out of 128 total)")
    print("  - Attackers report inflated quality (1.0 when true is 0.5)")
    print("  - 3 societies observe independently (Thor, Legion, Sprout)")
    print("  - Attacker compromises 1 society (Thor) to inflate reports")
    print()
    print("Expected:")
    print("  - Attack success rate < 5% (vs 23.6% in Session 76)")
    print("  - Outlier detection flags inflated reports")
    print("=" * 80)
    print()

    # Run test
    tester = QualityInflationAttackTester()
    result = tester.run_attack_test(num_attackers=10, num_generations=90)

    # Display results
    print("\n" + "=" * 80)
    print("RESULTS")
    print("=" * 80)
    print()

    print("Attack Metrics:")
    print("-" * 80)
    print(f"Attack success rate:     {result.attack_success_rate:6.1%}")
    print(f"Avg quality inflation:   {result.avg_quality_inflation:+6.3f}")
    print(f"Attacks detected:        {result.attacks_detected}")
    print(f"Attacks prevented:       {result.attacks_prevented}")
    print(f"False positives:         {result.false_positives}")
    print()

    print("Consensus Statistics:")
    print("-" * 80)
    print(f"Total consensus attempts:  {result.total_consensus_attempts}")
    print(f"Quorum reached:            {result.quorum_reached_pct:5.1f}%")
    print(f"Outliers detected:         {result.outliers_detected_pct:5.1f}%")
    print()

    print("Comparison to Session 76:")
    print("-" * 80)
    print(f"Session 76 (no consensus): {result.session76_attack_success:6.1%} attack success")
    print(f"Session 77 (Byzantine):    {result.attack_success_rate:6.1%} attack success")
    print(f"Improvement:               {result.improvement_pct:+6.1f}%")
    print()

    print("Test Result:")
    print("-" * 80)
    if result.passed:
        print(f"✅ PASS - Attack success ({result.attack_success_rate:.1%}) < 5%")
        print()
        print("Conclusion:")
        print("  Byzantine consensus successfully defends against quality inflation!")
        print(f"  Attack success reduced from 23.6% → {result.attack_success_rate:.1%}")
        print(f"  ({result.improvement_pct:+.1f}% improvement)")
    else:
        print(f"❌ FAIL - Attack success ({result.attack_success_rate:.1%}) ≥ 5%")
        print()
        print("Conclusion:")
        print("  Byzantine consensus helps but is insufficient.")
        print("  Additional defenses needed (e.g., stake slashing, reputation)")
    print()

    # Save results
    results_file = "/home/dp/ai-workspace/act/implementation/byzantine_quality_consensus_results.json"
    with open(results_file, 'w') as f:
        json.dump(asdict(result), f, indent=2)
    print(f"Results saved to: {results_file}")
    print()

    return result


if __name__ == "__main__":
    demo_byzantine_quality_consensus()
