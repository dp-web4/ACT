#!/usr/bin/env python3
"""
Trust Ledger Persistence for ACT Societies

Enables warm-start by persisting trust state to society ledger.
Solves the "cold bootstrap" problem discovered in Thor Session 75.

Problem:
- Trust-first architecture requires 40+ generations to accumulate evidence
- Every session starts from scratch (4 experts, 100% router_explore)
- Real model deployment stuck in bootstrap phase

Solution:
- Persist trust state to ACT society ledger
- Load prior trust on session start (warm-start)
- Immediate trust-driven mode activation if evidence exists

Architecture:
- Trust state → Cosmos SDK state machine
- Immutable audit trail (blockchain)
- Cross-session continuity
- Byzantine consensus for trust updates

Based on:
- Thor Sessions 74-75: Real model integration + bootstrap discovery
- Session 69: ACT multi-agent coordination with Byzantine consensus
- WEB4-PROP-006-v2.1: Trust-first MoE standard

Author: Legion (Session 70 - Autonomous Web4 Research)
Date: 2025-12-19
"""

from dataclasses import dataclass, asdict
from typing import Dict, List, Optional, Tuple
import json
import time
import hashlib
from pathlib import Path


@dataclass
class TrustEntry:
    """Single trust entry for (expert, context) pair."""
    expert_id: int
    context: str
    trust_value: float  # [0, 1]
    observation_count: int
    last_updated: int  # Unix timestamp
    session_id: str  # Which session created/updated this


@dataclass
class TrustSnapshot:
    """Complete trust state snapshot for a session."""
    snapshot_id: str  # Unique identifier
    society_id: str  # Which society owns this trust state
    timestamp: int  # When snapshot was created
    session_id: str  # Which session created it
    entries: List[TrustEntry]  # All trust entries
    metadata: Dict  # Additional context (mode stats, specialist counts, etc.)

    def to_dict(self) -> Dict:
        return {
            "snapshot_id": self.snapshot_id,
            "society_id": self.society_id,
            "timestamp": self.timestamp,
            "session_id": self.session_id,
            "entries": [asdict(e) for e in self.entries],
            "metadata": self.metadata
        }

    @classmethod
    def from_dict(cls, data: Dict) -> 'TrustSnapshot':
        entries = [TrustEntry(**e) for e in data["entries"]]
        return cls(
            snapshot_id=data["snapshot_id"],
            society_id=data["society_id"],
            timestamp=data["timestamp"],
            session_id=data["session_id"],
            entries=entries,
            metadata=data["metadata"]
        )


class TrustLedgerPersistence:
    """
    Persist trust state to ACT society ledger for warm-start.

    Provides:
    1. save_snapshot(): Persist current trust state
    2. load_snapshot(): Restore trust state from ledger
    3. list_snapshots(): Query available snapshots
    4. merge_snapshots(): Combine trust from multiple sessions/societies

    Storage:
    - File-based (JSON) for now
    - Cosmos SDK integration ready (same interface)
    - Blockchain audit trail future
    """

    def __init__(self, ledger_dir: Path):
        """
        Initialize trust ledger persistence.

        Args:
            ledger_dir: Directory for trust snapshots
        """
        self.ledger_dir = Path(ledger_dir)
        self.ledger_dir.mkdir(parents=True, exist_ok=True)

    def save_snapshot(
        self,
        society_id: str,
        session_id: str,
        trust_state: Dict[Tuple[int, str], float],
        observation_counts: Dict[Tuple[int, str], int],
        metadata: Optional[Dict] = None
    ) -> str:
        """
        Save trust snapshot to ledger.

        Args:
            society_id: Society identifier
            session_id: Session identifier
            trust_state: {(expert_id, context): trust_value}
            observation_counts: {(expert_id, context): count}
            metadata: Additional context (mode stats, etc.)

        Returns:
            snapshot_id: Unique identifier for this snapshot
        """
        timestamp = int(time.time())
        snapshot_id = hashlib.sha256(
            f"{society_id}:{session_id}:{timestamp}".encode()
        ).hexdigest()[:16]

        # Convert to TrustEntry list
        entries = []
        for (expert_id, context), trust_value in trust_state.items():
            entry = TrustEntry(
                expert_id=expert_id,
                context=context,
                trust_value=trust_value,
                observation_count=observation_counts.get((expert_id, context), 0),
                last_updated=timestamp,
                session_id=session_id
            )
            entries.append(entry)

        # Create snapshot
        snapshot = TrustSnapshot(
            snapshot_id=snapshot_id,
            society_id=society_id,
            timestamp=timestamp,
            session_id=session_id,
            entries=entries,
            metadata=metadata or {}
        )

        # Persist to file (Cosmos SDK integration would replace this)
        snapshot_file = self.ledger_dir / f"{snapshot_id}.json"
        with open(snapshot_file, 'w') as f:
            json.dump(snapshot.to_dict(), f, indent=2)

        # Update society index
        self._update_society_index(society_id, snapshot_id, timestamp)

        return snapshot_id

    def load_snapshot(self, snapshot_id: str) -> Optional[TrustSnapshot]:
        """
        Load trust snapshot from ledger.

        Args:
            snapshot_id: Snapshot identifier

        Returns:
            TrustSnapshot if found, None otherwise
        """
        snapshot_file = self.ledger_dir / f"{snapshot_id}.json"
        if not snapshot_file.exists():
            return None

        with open(snapshot_file, 'r') as f:
            data = json.load(f)

        return TrustSnapshot.from_dict(data)

    def load_latest_snapshot(self, society_id: str) -> Optional[TrustSnapshot]:
        """
        Load most recent snapshot for a society.

        Args:
            society_id: Society identifier

        Returns:
            Latest TrustSnapshot if exists, None otherwise
        """
        index = self._load_society_index(society_id)
        if not index or not index["snapshots"]:
            return None

        # Get most recent snapshot_id
        latest = max(index["snapshots"], key=lambda x: x["timestamp"])
        return self.load_snapshot(latest["snapshot_id"])

    def list_snapshots(self, society_id: Optional[str] = None) -> List[Dict]:
        """
        List available snapshots.

        Args:
            society_id: Filter by society (None = all societies)

        Returns:
            List of snapshot metadata
        """
        if society_id:
            index = self._load_society_index(society_id)
            return index.get("snapshots", []) if index else []

        # All societies
        snapshots = []
        for index_file in self.ledger_dir.glob("*_index.json"):
            with open(index_file, 'r') as f:
                index = json.load(f)
                snapshots.extend(index.get("snapshots", []))

        return snapshots

    def merge_snapshots(
        self,
        snapshot_ids: List[str],
        merge_strategy: str = "latest"
    ) -> Dict[Tuple[int, str], float]:
        """
        Merge trust from multiple snapshots.

        Args:
            snapshot_ids: Snapshots to merge
            merge_strategy: How to merge conflicting entries
                - "latest": Use most recent trust value
                - "average": Average trust values
                - "max": Use highest trust value

        Returns:
            Merged trust state {(expert, context): trust}
        """
        merged_trust = {}
        entry_timestamps = {}

        for snapshot_id in snapshot_ids:
            snapshot = self.load_snapshot(snapshot_id)
            if not snapshot:
                continue

            for entry in snapshot.entries:
                key = (entry.expert_id, entry.context)

                if merge_strategy == "latest":
                    # Keep most recent
                    if key not in entry_timestamps or entry.last_updated > entry_timestamps[key]:
                        merged_trust[key] = entry.trust_value
                        entry_timestamps[key] = entry.last_updated

                elif merge_strategy == "average":
                    # Average with existing
                    if key in merged_trust:
                        merged_trust[key] = (merged_trust[key] + entry.trust_value) / 2
                    else:
                        merged_trust[key] = entry.trust_value

                elif merge_strategy == "max":
                    # Keep highest trust
                    if key not in merged_trust or entry.trust_value > merged_trust[key]:
                        merged_trust[key] = entry.trust_value

        return merged_trust

    def _update_society_index(self, society_id: str, snapshot_id: str, timestamp: int):
        """Update society index with new snapshot."""
        index_file = self.ledger_dir / f"{society_id}_index.json"

        if index_file.exists():
            with open(index_file, 'r') as f:
                index = json.load(f)
        else:
            index = {"society_id": society_id, "snapshots": []}

        index["snapshots"].append({
            "snapshot_id": snapshot_id,
            "timestamp": timestamp
        })

        with open(index_file, 'w') as f:
            json.dump(index, f, indent=2)

    def _load_society_index(self, society_id: str) -> Optional[Dict]:
        """Load society index."""
        index_file = self.ledger_dir / f"{society_id}_index.json"
        if not index_file.exists():
            return None

        with open(index_file, 'r') as f:
            return json.load(f)


def warm_start_trust_selector(
    trust_selector,
    ledger: TrustLedgerPersistence,
    society_id: str,
    merge_strategy: str = "latest"
) -> bool:
    """
    Warm-start a trust selector from ledger.

    Args:
        trust_selector: TrustFirstMRHSelector or TrustCoordinator
        ledger: TrustLedgerPersistence instance
        society_id: Society to load trust from
        merge_strategy: How to handle multiple snapshots

    Returns:
        True if warm-started, False if cold-start (no prior trust)
    """
    # Load latest snapshot
    snapshot = ledger.load_latest_snapshot(society_id)
    if not snapshot:
        return False  # Cold start

    # Populate trust selector
    for entry in snapshot.entries:
        key = (entry.expert_id, entry.context)

        # Set trust value
        if hasattr(trust_selector, 'expert_trust'):
            # TrustFirstMRHSelector or TrustCoordinator format
            if entry.expert_id not in trust_selector.expert_trust:
                trust_selector.expert_trust[entry.expert_id] = {}
            trust_selector.expert_trust[entry.expert_id][entry.context] = entry.trust_value

        if hasattr(trust_selector, 'expert_observations'):
            # Set observation count
            if entry.expert_id not in trust_selector.expert_observations:
                trust_selector.expert_observations[entry.expert_id] = {}
            trust_selector.expert_observations[entry.expert_id][entry.context] = entry.observation_count

        if hasattr(trust_selector, 'society_trust'):
            # TrustCoordinator format
            trust_selector.society_trust.expert_trust[key] = entry.trust_value
            trust_selector.society_trust.expert_observations[key] = entry.observation_count

    return True  # Warm started


# Example usage
def demo_warm_start():
    """
    Demonstrate warm-start vs cold-start.

    Simulates two sessions:
    1. Session A: Cold start (bootstrap from scratch)
    2. Session B: Warm start (load Session A's trust)
    """
    print("\n" + "="*70)
    print("TRUST PERSISTENCE WARM-START DEMO")
    print("="*70)

    # Setup ledger
    ledger_dir = Path("/tmp/trust_ledger_demo")
    ledger = TrustLedgerPersistence(ledger_dir)

    print(f"\n✅ Trust ledger initialized: {ledger_dir}")

    # Simulate Session A: Cold start
    print(f"\n{'='*70}")
    print("SESSION A: COLD START (Bootstrap from scratch)")
    print(f"{'='*70}\n")

    session_a_trust = {}
    session_a_counts = {}

    # Simulate 50 generations of trust accumulation
    import numpy as np
    for gen in range(50):
        expert_id = int(np.random.choice([24, 42, 73, 79, 102]))  # Convert to Python int
        context = str(np.random.choice(["context_0", "context_1", "context_2"]))
        quality = float(0.7 + np.random.normal(0, 0.1))
        quality = float(np.clip(quality, 0, 1))

        key = (expert_id, context)
        current_trust = session_a_trust.get(key, 0.5)
        new_trust = float(0.7 * current_trust + 0.3 * quality)  # EWMA α=0.3

        session_a_trust[key] = new_trust
        session_a_counts[key] = session_a_counts.get(key, 0) + 1

    print(f"  Generations trained: 50")
    print(f"  Trust entries accumulated: {len(session_a_trust)}")
    print(f"  Total observations: {sum(session_a_counts.values())}")

    # Save snapshot
    snapshot_id = ledger.save_snapshot(
        society_id="demo-society",
        session_id="session_a",
        trust_state=session_a_trust,
        observation_counts=session_a_counts,
        metadata={"generations": 50, "mode": "cold_start"}
    )

    print(f"\n  ✅ Trust snapshot saved: {snapshot_id}")

    # Simulate Session B: Warm start
    print(f"\n{'='*70}")
    print("SESSION B: WARM START (Load Session A's trust)")
    print(f"{'='*70}\n")

    # Load snapshot
    loaded_snapshot = ledger.load_latest_snapshot("demo-society")
    if loaded_snapshot:
        print(f"  ✅ Loaded snapshot: {loaded_snapshot.snapshot_id}")
        print(f"  Timestamp: {time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(loaded_snapshot.timestamp))}")
        print(f"  Trust entries: {len(loaded_snapshot.entries)}")
        print(f"  Total observations: {sum(e.observation_count for e in loaded_snapshot.entries)}")

        # Show sample entries
        print(f"\n  Sample Trust Entries:")
        for entry in sorted(loaded_snapshot.entries, key=lambda x: x.trust_value, reverse=True)[:5]:
            print(f"    Expert {entry.expert_id}, {entry.context}: "
                  f"Trust={entry.trust_value:.3f}, Obs={entry.observation_count}")

    # List all snapshots
    print(f"\n{'='*70}")
    print("SNAPSHOT HISTORY")
    print(f"{'='*70}\n")

    snapshots = ledger.list_snapshots("demo-society")
    for snap in snapshots:
        timestamp_str = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(snap['timestamp']))
        print(f"  {snap['snapshot_id']}: {timestamp_str}")

    print(f"\n✅ Warm-start demo complete")
    print(f"\nBenefit: Session B starts with accumulated trust from Session A,")
    print(f"enabling immediate trust-driven mode instead of 40+ generation bootstrap!")


if __name__ == "__main__":
    demo_warm_start()
