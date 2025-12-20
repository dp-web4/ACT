#!/usr/bin/env python3
"""
Session 74 Track 2: Trust Monitoring Dashboard

Provides real-time monitoring and visualization of trust-first MoE architecture.

Problem:
- Production deployment needs operational visibility
- Need to track trust accumulation, expert utilization, ATP allocation
- Need to detect anomalies and performance degradation
- Need to validate Session 80 trust fix effectiveness

Features:
1. Real-time metrics collection
2. Trust evolution tracking
3. Expert utilization analysis
4. Mode distribution monitoring
5. ATP/ADP economic metrics
6. Anomaly detection
7. Performance benchmarking

Architecture:
- MetricsCollector: Collects metrics from TrustFirstMRHSelector
- DashboardMonitor: Aggregates and analyzes metrics
- AlertSystem: Detects anomalies and triggers alerts
- ReportGenerator: Creates summary reports

Based on:
- Session 80: Trust fix validation (73.3% trust_driven)
- Session 73: Security analysis and attack detection
- Session 71: Epsilon + warm-start integration
- WEB4-PROP-006-v2.2: Trust-first standard

Created: 2025-12-20 (Legion Session 74)
Author: Legion (Autonomous Web4 Research)
"""

import time
import json
from dataclasses import dataclass, field, asdict
from typing import Dict, List, Optional, Tuple, Any
from collections import defaultdict, deque
import statistics


@dataclass
class TrustMetrics:
    """Snapshot of trust state at a point in time."""
    timestamp: float
    generation: int

    # Trust state
    total_trust_entries: int
    contexts_with_evidence: int
    avg_trust_value: float
    max_trust_value: float
    min_trust_value: float

    # Expert utilization
    unique_experts_used: int
    total_experts: int
    utilization_rate: float
    specialists_count: int
    generalists_count: int

    # Mode distribution
    trust_driven_count: int
    router_explore_count: int
    forced_exploration_count: int
    trust_driven_rate: float

    # ATP/ADP economics (if available)
    total_atp_allocated: float = 0.0
    total_adp_earned: float = 0.0
    avg_quality: float = 0.0

    # Performance
    avg_selection_time_ms: float = 0.0


@dataclass
class Alert:
    """Alert for anomaly detection."""
    timestamp: float
    severity: str  # "info", "warning", "critical"
    category: str  # "trust", "utilization", "performance", "economics"
    message: str
    details: Dict[str, Any] = field(default_factory=dict)


class MetricsCollector:
    """
    Collects metrics from TrustFirstMRHSelector during execution.

    Integrates with selector to capture real-time metrics.
    """

    def __init__(self, selector, window_size: int = 100):
        """
        Initialize metrics collector.

        Args:
            selector: TrustFirstMRHSelector instance
            window_size: Sliding window for metrics
        """
        self.selector = selector
        self.window_size = window_size

        # Metrics history
        self.metrics_history: List[TrustMetrics] = []

        # Real-time tracking
        self.selection_times: deque = deque(maxlen=window_size)
        self.recent_qualities: deque = deque(maxlen=window_size)
        self.mode_counts = defaultdict(int)

        # ATP/ADP tracking
        self.atp_history: deque = deque(maxlen=window_size)
        self.adp_history: deque = deque(maxlen=window_size)

    def collect_snapshot(self, generation: int) -> TrustMetrics:
        """
        Collect current metrics snapshot.

        Args:
            generation: Current generation number

        Returns:
            TrustMetrics snapshot
        """
        # Count trust entries
        total_entries = sum(
            len(experts)
            for context_experts in self.selector.trust_scores.values()
            for experts in context_experts.values()
        )

        # Count contexts with evidence
        contexts_with_evidence = sum(
            1 for context in range(self.selector.num_contexts)
            if self.selector._has_sufficient_evidence(context)
        )

        # Calculate trust statistics
        all_trust_values = []
        for context_experts in self.selector.trust_scores.values():
            for expert_scores in context_experts.values():
                all_trust_values.extend(expert_scores)

        avg_trust = statistics.mean(all_trust_values) if all_trust_values else 0.0
        max_trust = max(all_trust_values) if all_trust_values else 0.0
        min_trust = min(all_trust_values) if all_trust_values else 0.0

        # Expert utilization
        all_experts_used = set()
        for context_experts in self.selector.trust_scores.values():
            all_experts_used.update(context_experts.keys())

        unique_experts = len(all_experts_used)
        utilization_rate = unique_experts / self.selector.num_experts

        # Count specialists vs generalists
        specialist_threshold = self.selector.num_contexts * 0.7
        specialists = 0
        generalists = 0

        for expert_id in all_experts_used:
            contexts_count = sum(
                1 for context_experts in self.selector.trust_scores.values()
                if expert_id in context_experts
            )
            if contexts_count == 1:
                specialists += 1
            else:
                generalists += 1

        # Mode distribution
        stats = self.selector.get_stats()
        trust_driven = stats.get("trust_driven", 0)
        router_explore = stats.get("router_explore", 0)
        forced_exploration = stats.get("forced_exploration", 0)
        total_selections = stats.get("total_selections", 1)

        trust_driven_rate = trust_driven / total_selections if total_selections > 0 else 0.0

        # ATP/ADP economics
        total_atp = sum(self.atp_history) if self.atp_history else 0.0
        total_adp = sum(self.adp_history) if self.adp_history else 0.0
        avg_quality = statistics.mean(self.recent_qualities) if self.recent_qualities else 0.0

        # Performance
        avg_selection_time = statistics.mean(self.selection_times) if self.selection_times else 0.0

        return TrustMetrics(
            timestamp=time.time(),
            generation=generation,
            total_trust_entries=total_entries,
            contexts_with_evidence=contexts_with_evidence,
            avg_trust_value=avg_trust,
            max_trust_value=max_trust,
            min_trust_value=min_trust,
            unique_experts_used=unique_experts,
            total_experts=self.selector.num_experts,
            utilization_rate=utilization_rate,
            specialists_count=specialists,
            generalists_count=generalists,
            trust_driven_count=trust_driven,
            router_explore_count=router_explore,
            forced_exploration_count=forced_exploration,
            trust_driven_rate=trust_driven_rate,
            total_atp_allocated=total_atp,
            total_adp_earned=total_adp,
            avg_quality=avg_quality,
            avg_selection_time_ms=avg_selection_time
        )

    def record_selection(self, mode: str, selection_time_ms: float):
        """Record expert selection event."""
        self.mode_counts[mode] += 1
        self.selection_times.append(selection_time_ms)

    def record_quality(self, quality: float):
        """Record quality feedback."""
        self.recent_qualities.append(quality)

    def record_atp_adp(self, atp: float, adp: float):
        """Record ATP allocation and ADP earned."""
        self.atp_history.append(atp)
        self.adp_history.append(adp)


class AlertSystem:
    """
    Monitors metrics and generates alerts for anomalies.
    """

    def __init__(
        self,
        min_utilization: float = 0.1,
        min_trust_driven_rate: float = 0.3,
        max_selection_time_ms: float = 100.0,
        min_avg_quality: float = 0.5
    ):
        """
        Initialize alert system.

        Args:
            min_utilization: Minimum acceptable expert utilization rate
            min_trust_driven_rate: Minimum trust_driven mode rate
            max_selection_time_ms: Maximum acceptable selection time
            min_avg_quality: Minimum acceptable average quality
        """
        self.min_utilization = min_utilization
        self.min_trust_driven_rate = min_trust_driven_rate
        self.max_selection_time_ms = max_selection_time_ms
        self.min_avg_quality = min_avg_quality

        self.alerts: List[Alert] = []

    def check_metrics(self, metrics: TrustMetrics) -> List[Alert]:
        """
        Check metrics for anomalies and generate alerts.

        Args:
            metrics: Current metrics snapshot

        Returns:
            List of new alerts
        """
        new_alerts = []

        # Check expert utilization
        if metrics.utilization_rate < self.min_utilization:
            new_alerts.append(Alert(
                timestamp=metrics.timestamp,
                severity="warning",
                category="utilization",
                message=f"Low expert utilization: {metrics.utilization_rate:.1%}",
                details={
                    "utilization_rate": metrics.utilization_rate,
                    "unique_experts": metrics.unique_experts_used,
                    "total_experts": metrics.total_experts
                }
            ))

        # Check trust_driven activation
        if metrics.generation > 20 and metrics.trust_driven_rate < self.min_trust_driven_rate:
            new_alerts.append(Alert(
                timestamp=metrics.timestamp,
                severity="warning",
                category="trust",
                message=f"Low trust_driven rate: {metrics.trust_driven_rate:.1%} at gen {metrics.generation}",
                details={
                    "trust_driven_rate": metrics.trust_driven_rate,
                    "trust_driven_count": metrics.trust_driven_count,
                    "contexts_with_evidence": metrics.contexts_with_evidence
                }
            ))

        # Check performance
        if metrics.avg_selection_time_ms > self.max_selection_time_ms:
            new_alerts.append(Alert(
                timestamp=metrics.timestamp,
                severity="warning",
                category="performance",
                message=f"Slow selection time: {metrics.avg_selection_time_ms:.2f}ms",
                details={
                    "avg_selection_time_ms": metrics.avg_selection_time_ms
                }
            ))

        # Check quality
        if metrics.avg_quality > 0 and metrics.avg_quality < self.min_avg_quality:
            new_alerts.append(Alert(
                timestamp=metrics.timestamp,
                severity="critical",
                category="economics",
                message=f"Low average quality: {metrics.avg_quality:.2f}",
                details={
                    "avg_quality": metrics.avg_quality,
                    "total_adp_earned": metrics.total_adp_earned
                }
            ))

        # Session 80 validation: Check if trust fix is working
        if metrics.generation == 10:
            if metrics.trust_driven_rate == 0:
                new_alerts.append(Alert(
                    timestamp=metrics.timestamp,
                    severity="critical",
                    category="trust",
                    message="Trust fix validation FAILED: 0% trust_driven at gen 10",
                    details={
                        "expected": "> 0%",
                        "actual": "0%",
                        "contexts_with_evidence": metrics.contexts_with_evidence
                    }
                ))
            else:
                new_alerts.append(Alert(
                    timestamp=metrics.timestamp,
                    severity="info",
                    category="trust",
                    message=f"Trust fix validation PASSED: {metrics.trust_driven_rate:.1%} trust_driven at gen 10",
                    details={
                        "trust_driven_rate": metrics.trust_driven_rate
                    }
                ))

        self.alerts.extend(new_alerts)
        return new_alerts


class DashboardMonitor:
    """
    Main dashboard for monitoring trust-first MoE architecture.
    """

    def __init__(
        self,
        selector,
        snapshot_interval: int = 10,
        alert_config: Optional[Dict] = None
    ):
        """
        Initialize dashboard monitor.

        Args:
            selector: TrustFirstMRHSelector instance
            snapshot_interval: Generations between snapshots
            alert_config: Configuration for alert system
        """
        self.collector = MetricsCollector(selector)
        self.alert_system = AlertSystem(**(alert_config or {}))
        self.snapshot_interval = snapshot_interval

        self.start_time = time.time()

    def collect_and_alert(self, generation: int) -> Tuple[TrustMetrics, List[Alert]]:
        """
        Collect metrics and check for alerts.

        Args:
            generation: Current generation number

        Returns:
            (metrics, alerts) tuple
        """
        metrics = self.collector.collect_snapshot(generation)
        alerts = self.alert_system.check_metrics(metrics)

        self.collector.metrics_history.append(metrics)

        return metrics, alerts

    def should_snapshot(self, generation: int) -> bool:
        """Check if snapshot should be taken at this generation."""
        return generation % self.snapshot_interval == 0

    def generate_report(self) -> Dict[str, Any]:
        """
        Generate comprehensive monitoring report.

        Returns:
            Report dictionary
        """
        if not self.collector.metrics_history:
            return {"error": "No metrics collected"}

        latest = self.collector.metrics_history[-1]
        first = self.collector.metrics_history[0]

        # Calculate trends
        utilization_trend = latest.utilization_rate - first.utilization_rate
        trust_driven_trend = latest.trust_driven_rate - first.trust_driven_rate

        # Aggregate alerts by severity
        alerts_by_severity = defaultdict(int)
        for alert in self.alert_system.alerts:
            alerts_by_severity[alert.severity] += 1

        return {
            "session_duration_seconds": time.time() - self.start_time,
            "total_generations": latest.generation,
            "snapshots_collected": len(self.collector.metrics_history),

            "current_state": {
                "utilization_rate": latest.utilization_rate,
                "unique_experts": latest.unique_experts_used,
                "specialists": latest.specialists_count,
                "generalists": latest.generalists_count,
                "trust_driven_rate": latest.trust_driven_rate,
                "contexts_with_evidence": latest.contexts_with_evidence,
                "avg_quality": latest.avg_quality
            },

            "trends": {
                "utilization_change": utilization_trend,
                "trust_driven_change": trust_driven_trend
            },

            "performance": {
                "avg_selection_time_ms": latest.avg_selection_time_ms,
                "total_trust_entries": latest.total_trust_entries
            },

            "economics": {
                "total_atp_allocated": latest.total_atp_allocated,
                "total_adp_earned": latest.total_adp_earned,
                "efficiency": latest.total_adp_earned / latest.total_atp_allocated if latest.total_atp_allocated > 0 else 0.0
            },

            "alerts": {
                "total_alerts": len(self.alert_system.alerts),
                "by_severity": dict(alerts_by_severity),
                "recent_critical": [
                    asdict(alert) for alert in self.alert_system.alerts
                    if alert.severity == "critical"
                ][-5:]  # Last 5 critical alerts
            },

            "session_80_validation": {
                "expected_trust_driven_rate": "> 70%",
                "actual_trust_driven_rate": f"{latest.trust_driven_rate:.1%}",
                "status": "PASS" if latest.trust_driven_rate > 0.7 else "FAIL"
            }
        }

    def print_dashboard(self, generation: int):
        """Print real-time dashboard to console."""
        metrics = self.collector.metrics_history[-1] if self.collector.metrics_history else None

        if not metrics:
            print("No metrics available")
            return

        print("\n" + "="*70)
        print(f"TRUST-FIRST MoE MONITORING DASHBOARD - Generation {generation}")
        print("="*70)

        print(f"\n📊 EXPERT UTILIZATION")
        print(f"  Unique Experts: {metrics.unique_experts_used}/{metrics.total_experts} ({metrics.utilization_rate:.1%})")
        print(f"  Specialists: {metrics.specialists_count}")
        print(f"  Generalists: {metrics.generalists_count}")

        print(f"\n🔒 TRUST STATE")
        print(f"  Trust Entries: {metrics.total_trust_entries}")
        print(f"  Contexts w/ Evidence: {metrics.contexts_with_evidence}")
        print(f"  Avg Trust Value: {metrics.avg_trust_value:.3f}")
        print(f"  Trust Range: [{metrics.min_trust_value:.3f}, {metrics.max_trust_value:.3f}]")

        print(f"\n🎯 MODE DISTRIBUTION")
        print(f"  Trust-Driven: {metrics.trust_driven_count} ({metrics.trust_driven_rate:.1%})")
        print(f"  Router Explore: {metrics.router_explore_count}")
        print(f"  Forced Exploration: {metrics.forced_exploration_count}")

        if metrics.avg_quality > 0:
            print(f"\n💰 ECONOMICS")
            print(f"  Total ATP Allocated: {metrics.total_atp_allocated:.2f}")
            print(f"  Total ADP Earned: {metrics.total_adp_earned:.2f}")
            print(f"  Avg Quality: {metrics.avg_quality:.3f}")
            efficiency = metrics.total_adp_earned / metrics.total_atp_allocated if metrics.total_atp_allocated > 0 else 0.0
            print(f"  Efficiency (ADP/ATP): {efficiency:.3f}")

        print(f"\n⚡ PERFORMANCE")
        print(f"  Avg Selection Time: {metrics.avg_selection_time_ms:.2f}ms")

        # Recent alerts
        recent_alerts = self.alert_system.alerts[-3:] if self.alert_system.alerts else []
        if recent_alerts:
            print(f"\n⚠️  RECENT ALERTS")
            for alert in recent_alerts:
                icon = "🔴" if alert.severity == "critical" else "🟡" if alert.severity == "warning" else "ℹ️ "
                print(f"  {icon} [{alert.category}] {alert.message}")

        print("="*70)


def demo_monitoring_dashboard():
    """
    Demonstrate monitoring dashboard with simulated selector.

    Simulates Session 80 scenario to validate monitoring.
    """
    print("\n" + "="*70)
    print("TRUST MONITORING DASHBOARD DEMONSTRATION")
    print("="*70)

    # Mock selector for demo
    class MockSelector:
        def __init__(self):
            self.num_experts = 128
            self.num_contexts = 3
            self.trust_scores = defaultdict(lambda: defaultdict(list))
            self.stats_data = {
                "total_selections": 0,
                "trust_driven": 0,
                "router_explore": 0,
                "forced_exploration": 0
            }

        def _has_sufficient_evidence(self, context: int) -> bool:
            # Simulate evidence accumulation
            return len(self.trust_scores.get(context, {})) >= 2

        def get_stats(self):
            return self.stats_data

        def simulate_generation(self, gen: int):
            """Simulate one generation of expert selection."""
            # Simulate epsilon-greedy + trust accumulation
            import random

            if gen < 5:
                mode = "router_explore"
                experts = random.sample(range(self.num_experts), 4)
            elif random.random() < 0.2:
                mode = "forced_exploration"
                experts = random.sample(range(self.num_experts), 4)
            elif gen >= 8:
                mode = "trust_driven"
                # Prefer experts with high trust
                experts = random.sample(range(min(60, self.num_experts)), 4)
            else:
                mode = "router_explore"
                experts = random.sample(range(self.num_experts), 4)

            # Update stats
            self.stats_data["total_selections"] += 1
            self.stats_data[mode] += 1

            # Update trust for selected experts
            context = gen % self.num_contexts
            quality = random.uniform(0.6, 0.9)

            for expert_id in experts:
                self.trust_scores[context][expert_id].append(quality)

            return mode, experts, quality

    selector = MockSelector()
    monitor = DashboardMonitor(
        selector,
        snapshot_interval=10,
        alert_config={
            "min_utilization": 0.1,
            "min_trust_driven_rate": 0.5,
            "max_selection_time_ms": 50.0,
            "min_avg_quality": 0.5
        }
    )

    print("\n" + "="*70)
    print("SIMULATION: Session 80 Scenario (90 generations)")
    print("="*70)

    # Simulate 90 generations
    for gen in range(1, 91):
        start_time = time.time()
        mode, experts, quality = selector.simulate_generation(gen)
        selection_time = (time.time() - start_time) * 1000  # ms

        # Record metrics
        monitor.collector.record_selection(mode, selection_time)
        monitor.collector.record_quality(quality)
        monitor.collector.record_atp_adp(atp=1.0, adp=quality)

        # Take snapshot every 10 generations
        if monitor.should_snapshot(gen):
            metrics, alerts = monitor.collect_and_alert(gen)

            # Print dashboard at key generations
            if gen in [10, 30, 50, 70, 90]:
                monitor.print_dashboard(gen)

                # Show alerts
                if alerts:
                    print(f"\n🔔 NEW ALERTS:")
                    for alert in alerts:
                        print(f"  [{alert.severity.upper()}] {alert.message}")

    # Final report
    print("\n" + "="*70)
    print("FINAL MONITORING REPORT")
    print("="*70)

    report = monitor.generate_report()
    print(json.dumps(report, indent=2))

    print("\n" + "="*70)
    print("KEY FEATURES VALIDATED")
    print("="*70)

    print("\n✅ Real-Time Metrics:")
    print("   - Expert utilization tracking")
    print("   - Trust accumulation monitoring")
    print("   - Mode distribution analysis")
    print("   - ATP/ADP economic metrics")

    print("\n✅ Anomaly Detection:")
    print("   - Low utilization alerts")
    print("   - Trust_driven activation monitoring")
    print("   - Performance degradation detection")
    print("   - Quality threshold enforcement")

    print("\n✅ Session 80 Validation:")
    print("   - Trust fix effectiveness tracking")
    print("   - Expected vs actual trust_driven rate")
    print("   - Automatic validation at gen 10")

    print("\n✅ Reporting:")
    print("   - Comprehensive session reports")
    print("   - Trend analysis")
    print("   - Alert aggregation by severity")
    print("   - Economics efficiency tracking")

    print("="*70)


if __name__ == "__main__":
    demo_monitoring_dashboard()
