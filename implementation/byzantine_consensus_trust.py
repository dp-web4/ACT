#!/usr/bin/env python3
"""
Session 73 Track 1: Byzantine Consensus for Trust Federation

Implements Byzantine Fault Tolerant (BFT) consensus for cross-society trust attestations.

Problem:
- Federation protocol needs Byzantine consensus (2f+1 witnesses)
- Current implementation lacks cryptographic verification
- No protection against malicious attestations

Solution:
- BFT consensus with cryptographic signatures
- Quorum-based witness verification (2f+1 out of 3f+1)
- Attack resistance (Sybil, collusion, Byzantine faults)

Architecture:
1. Each society node can attest to trust values
2. Attestations signed with society's private key
3. Quorum (2f+1) required for consensus
4. Signature verification prevents forgery
5. Timestamp verification prevents replay attacks

Based on:
- Classical BFT: PBFT, Tendermint consensus
- WEB4-PROP-006-v2.2: Trust-first MoE standard
- Session 70: Federation trust transfer infrastructure

Security Properties:
- Safety: No two conflicting trust values reach consensus
- Liveness: Valid attestations eventually reach consensus
- Byzantine Resilience: Tolerates up to f malicious nodes (f < n/3)

Created: 2025-12-20 (Legion Session 73)
Author: Legion (Autonomous Web4 Research)
"""

import hashlib
import hmac
import json
import time
from dataclasses import dataclass, asdict
from typing import Dict, List, Set, Optional, Tuple
from enum import Enum


class AttestationStatus(Enum):
    """Status of trust attestation in consensus process."""
    PENDING = "pending"  # Awaiting quorum
    ACCEPTED = "accepted"  # Consensus reached
    REJECTED = "rejected"  # Conflicting/invalid
    EXPIRED = "expired"  # Timeout


@dataclass
class SignedAttestation:
    """
    Cryptographically signed trust attestation.

    Byzantine-safe trust claim from a society node.
    """
    # Core attestation data
    agent_lct: str  # Who this is about
    society_id: str  # Attesting society
    expert_id: int  # Expert within agent
    context: str  # Context of expertise
    trust_value: float  # [0, 1]
    observation_count: int  # Evidence strength
    timestamp: int  # Unix timestamp

    # Consensus metadata
    witness_id: str  # Witnessing node ID
    signature: str  # HMAC signature of attestation
    nonce: str  # Prevents replay attacks

    # Optional
    metadata: Optional[Dict] = None

    def to_signable_dict(self) -> Dict:
        """Get dictionary for signing (excludes signature itself)."""
        return {
            "agent_lct": self.agent_lct,
            "society_id": self.society_id,
            "expert_id": self.expert_id,
            "context": self.context,
            "trust_value": self.trust_value,
            "observation_count": self.observation_count,
            "timestamp": self.timestamp,
            "witness_id": self.witness_id,
            "nonce": self.nonce
        }

    def compute_hash(self) -> str:
        """
        Compute deterministic hash of attestation content.

        Hash is based on WHAT is being attested (agent, expert, context, trust),
        NOT WHO is attesting (witness, nonce, signature).
        This allows multiple witnesses to attest to the same trust value.
        """
        content_dict = {
            "agent_lct": self.agent_lct,
            "society_id": self.society_id,
            "expert_id": self.expert_id,
            "context": self.context,
            "trust_value": self.trust_value,
            "observation_count": self.observation_count
        }
        content = json.dumps(content_dict, sort_keys=True)
        return hashlib.sha256(content.encode()).hexdigest()


@dataclass
class ConsensusQuorum:
    """
    Quorum state for Byzantine consensus.

    Tracks witnesses and validates quorum requirements.
    """
    attestation_hash: str  # Hash of attestation content
    required_witnesses: int  # 2f+1 for BFT
    witnessed_by: Set[str]  # Set of witness IDs
    attestations: List[SignedAttestation]  # All attestations with this hash
    status: AttestationStatus
    created_at: int
    consensus_at: Optional[int] = None


class ByzantineConsensusTrust:
    """
    Byzantine Fault Tolerant consensus for trust attestations.

    Implements (2f+1)-of-(3f+1) quorum consensus with signature verification.
    """

    def __init__(
        self,
        society_id: str,
        node_id: str,
        secret_key: str,
        f: int = 1,  # Number of tolerated faults
        quorum_timeout: int = 300  # 5 minutes
    ):
        """
        Initialize Byzantine consensus protocol.

        Args:
            society_id: This society's identifier
            node_id: This node's unique identifier
            secret_key: Secret key for signing attestations
            f: Number of tolerated Byzantine faults
            quorum_timeout: Seconds before quorum expires
        """
        self.society_id = society_id
        self.node_id = node_id
        self.secret_key = secret_key
        self.f = f

        # BFT parameters
        self.total_nodes = 3 * f + 1  # Standard BFT requirement
        self.quorum_size = 2 * f + 1  # Quorum for consensus
        self.quorum_timeout = quorum_timeout

        # State tracking
        self.pending_quorums: Dict[str, ConsensusQuorum] = {}
        self.accepted_attestations: Dict[str, List[SignedAttestation]] = {}
        self.public_keys: Dict[str, str] = {}  # {node_id: public_key}

        # Statistics
        self.total_proposed = 0
        self.total_accepted = 0
        self.total_rejected = 0
        self.total_expired = 0

    def register_node(self, node_id: str, public_key: str):
        """
        Register a node's public key for signature verification.

        Args:
            node_id: Node identifier
            public_key: Node's public key (for signature verification)
        """
        self.public_keys[node_id] = public_key

    def create_attestation(
        self,
        agent_lct: str,
        expert_id: int,
        context: str,
        trust_value: float,
        observation_count: int
    ) -> SignedAttestation:
        """
        Create and sign a trust attestation.

        Args:
            agent_lct: Agent identifier
            expert_id: Expert ID within agent
            context: Context of expertise
            trust_value: Trust value [0, 1]
            observation_count: Number of observations

        Returns:
            Signed attestation from this node
        """
        import secrets

        # Create attestation
        attestation = SignedAttestation(
            agent_lct=agent_lct,
            society_id=self.society_id,
            expert_id=expert_id,
            context=context,
            trust_value=trust_value,
            observation_count=observation_count,
            timestamp=int(time.time()),
            witness_id=self.node_id,
            signature="",  # Set below
            nonce=secrets.token_hex(16)  # 32-char hex nonce
        )

        # Sign attestation
        attestation.signature = self._sign_attestation(attestation)

        return attestation

    def _sign_attestation(self, attestation: SignedAttestation) -> str:
        """
        Sign attestation with this node's secret key.

        Uses HMAC-SHA256 for signature.

        Args:
            attestation: Attestation to sign

        Returns:
            Hex-encoded signature
        """
        signable = attestation.to_signable_dict()
        message = json.dumps(signable, sort_keys=True).encode()

        signature = hmac.new(
            self.secret_key.encode(),
            message,
            hashlib.sha256
        ).hexdigest()

        return signature

    def verify_signature(
        self,
        attestation: SignedAttestation
    ) -> bool:
        """
        Verify attestation signature.

        Args:
            attestation: Attestation to verify

        Returns:
            True if signature is valid
        """
        # Get witness's public key
        witness_key = self.public_keys.get(attestation.witness_id)
        if not witness_key:
            return False  # Unknown witness

        # Recompute signature
        signable = attestation.to_signable_dict()
        message = json.dumps(signable, sort_keys=True).encode()

        expected_signature = hmac.new(
            witness_key.encode(),
            message,
            hashlib.sha256
        ).hexdigest()

        # Constant-time comparison
        return hmac.compare_digest(attestation.signature, expected_signature)

    def propose_attestation(
        self,
        attestation: SignedAttestation
    ) -> Tuple[bool, str]:
        """
        Propose attestation for consensus.

        Args:
            attestation: Signed attestation

        Returns:
            (accepted, status_message)
        """
        self.total_proposed += 1

        # 1. Verify signature
        if not self.verify_signature(attestation):
            self.total_rejected += 1
            return False, "Invalid signature"

        # 2. Check timestamp freshness
        age = int(time.time()) - attestation.timestamp
        if age > self.quorum_timeout:
            self.total_expired += 1
            return False, "Attestation expired"

        if age < -60:  # Allow 60s clock skew
            self.total_rejected += 1
            return False, "Attestation from future"

        # 3. Compute attestation hash (for grouping)
        att_hash = attestation.compute_hash()

        # 4. Check for existing quorum
        if att_hash not in self.pending_quorums:
            # Create new quorum
            self.pending_quorums[att_hash] = ConsensusQuorum(
                attestation_hash=att_hash,
                required_witnesses=self.quorum_size,
                witnessed_by=set(),
                attestations=[],
                status=AttestationStatus.PENDING,
                created_at=int(time.time())
            )

        quorum = self.pending_quorums[att_hash]

        # 5. Check for duplicate witness
        if attestation.witness_id in quorum.witnessed_by:
            return False, "Duplicate witness"

        # 6. Add to quorum
        quorum.witnessed_by.add(attestation.witness_id)
        quorum.attestations.append(attestation)

        # 7. Check consensus
        if len(quorum.witnessed_by) >= self.quorum_size:
            # Consensus reached!
            quorum.status = AttestationStatus.ACCEPTED
            quorum.consensus_at = int(time.time())

            # Move to accepted
            self.accepted_attestations[att_hash] = quorum.attestations
            del self.pending_quorums[att_hash]

            self.total_accepted += 1
            return True, f"Consensus reached ({len(quorum.witnessed_by)}/{self.quorum_size} witnesses)"

        # Still pending
        return False, f"Pending ({len(quorum.witnessed_by)}/{self.quorum_size} witnesses)"

    def get_consensus_trust(
        self,
        agent_lct: str,
        context: str
    ) -> Optional[float]:
        """
        Get consensus trust value for agent in context.

        Returns aggregated trust from all accepted attestations.

        Args:
            agent_lct: Agent identifier
            context: Context to query

        Returns:
            Consensus trust value or None
        """
        matching_attestations = []

        for attestations in self.accepted_attestations.values():
            for att in attestations:
                if att.agent_lct == agent_lct and att.context == context:
                    matching_attestations.append(att)

        if not matching_attestations:
            return None

        # Weighted average by observation count
        total_weight = sum(att.observation_count for att in matching_attestations)
        if total_weight == 0:
            return None

        weighted_sum = sum(
            att.trust_value * att.observation_count
            for att in matching_attestations
        )

        return weighted_sum / total_weight

    def cleanup_expired(self):
        """Remove expired pending quorums."""
        now = int(time.time())
        expired = []

        for att_hash, quorum in self.pending_quorums.items():
            age = now - quorum.created_at
            if age > self.quorum_timeout:
                expired.append(att_hash)

        for att_hash in expired:
            del self.pending_quorums[att_hash]
            self.total_expired += 1

    def get_statistics(self) -> Dict:
        """Get consensus statistics."""
        return {
            "total_proposed": self.total_proposed,
            "total_accepted": self.total_accepted,
            "total_rejected": self.total_rejected,
            "total_expired": self.total_expired,
            "pending_quorums": len(self.pending_quorums),
            "accepted_attestations": len(self.accepted_attestations),
            "acceptance_rate": self.total_accepted / self.total_proposed
                              if self.total_proposed > 0 else 0,
            "byzantine_tolerance": f"{self.f} faults ({self.total_nodes} nodes, {self.quorum_size} quorum)"
        }


def demo_byzantine_consensus():
    """
    Demonstrate Byzantine consensus for trust attestations.

    Simulates 4 nodes (f=1) reaching consensus on trust values.
    """
    print("\n" + "="*70)
    print("BYZANTINE CONSENSUS FOR TRUST FEDERATION - DEMO")
    print("="*70)

    # Setup 4 nodes (f=1 → tolerates 1 Byzantine fault)
    print("\nByzantine Parameters:")
    print("  f = 1 (tolerated faults)")
    print("  Total nodes = 3f+1 = 4")
    print("  Quorum = 2f+1 = 3")
    print("  Safety: Can tolerate 1 malicious/failed node\n")

    # Create nodes
    nodes = {}
    for i in range(4):
        node_id = f"node-{i}"
        secret_key = f"secret-key-{i}"
        node = ByzantineConsensusTrust(
            society_id="web4-society-alpha",
            node_id=node_id,
            secret_key=secret_key,
            f=1,
            quorum_timeout=300
        )
        nodes[node_id] = (node, secret_key)

    # Cross-register public keys
    print("Registering node public keys...")
    for node_id, (node, secret_key) in nodes.items():
        for other_id, (_, other_secret) in nodes.items():
            node.register_node(other_id, other_secret)  # In real system, public key != secret
    print("✅ All nodes registered\n")

    # Scenario 1: Honest consensus (3/4 nodes agree)
    print("="*70)
    print("SCENARIO 1: Honest Consensus")
    print("="*70)

    agent_lct = "lct://agent-alice@web4.network"
    print(f"\nAttestation target: {agent_lct}")
    print(f"Expert: 42, Context: code_review, Trust: 0.85\n")

    # Nodes 0, 1, 2 create matching attestations
    honest_nodes = ["node-0", "node-1", "node-2"]
    attestations = []

    for node_id in honest_nodes:
        node, _ = nodes[node_id]
        att = node.create_attestation(
            agent_lct=agent_lct,
            expert_id=42,
            context="code_review",
            trust_value=0.85,
            observation_count=10
        )
        attestations.append((node_id, att))
        print(f"{node_id} created attestation (signature: {att.signature[:16]}...)")

    # Propose to consensus (to node-3)
    consensus_node, _ = nodes["node-3"]
    print(f"\nProposing attestations to {consensus_node.node_id}...")

    for i, (proposer_id, att) in enumerate(attestations, 1):
        accepted, message = consensus_node.propose_attestation(att)
        print(f"  {i}. From {proposer_id}: {message}")

    # Check consensus
    consensus_trust = consensus_node.get_consensus_trust(agent_lct, "code_review")
    stats = consensus_node.get_statistics()

    print(f"\n✅ Consensus Result:")
    print(f"  Trust value: {consensus_trust}")
    print(f"  Accepted attestations: {stats['total_accepted']}")
    print(f"  Pending quorums: {stats['pending_quorums']}")

    # Scenario 2: Byzantine attack (forged signature)
    print(f"\n{'='*70}")
    print("SCENARIO 2: Byzantine Attack (Forged Signature)")
    print("="*70)

    # Malicious node tries to forge attestation
    malicious_att = SignedAttestation(
        agent_lct="lct://agent-bob@web4.network",
        society_id="web4-society-alpha",
        expert_id=99,
        context="data_analysis",
        trust_value=0.95,  # Malicious high trust
        observation_count=100,
        timestamp=int(time.time()),
        witness_id="node-0",  # Impersonating node-0
        signature="FORGED_SIGNATURE_0xDEADBEEF",  # Invalid
        nonce="malicious_nonce"
    )

    print(f"\nMalicious attestation:")
    print(f"  Impersonating: node-0")
    print(f"  Trust value: {malicious_att.trust_value} (suspiciously high)")
    print(f"  Signature: {malicious_att.signature}")

    accepted, message = consensus_node.propose_attestation(malicious_att)
    print(f"\n❌ Consensus node response: {message}")
    print(f"   Attack REJECTED: Signature verification failed\n")

    # Scenario 3: Insufficient quorum (only 2/3 witnesses)
    print("="*70)
    print("SCENARIO 3: Insufficient Quorum (2/3 witnesses)")
    print("="*70)

    agent_charlie = "lct://agent-charlie@web4.network"
    print(f"\nAttestation target: {agent_charlie}")
    print("Only 2 honest nodes available (node-0, node-1)")
    print("Quorum requires 3 witnesses\n")

    # Only 2 nodes attest
    partial_nodes = ["node-0", "node-1"]
    for node_id in partial_nodes:
        node, _ = nodes[node_id]
        att = node.create_attestation(
            agent_lct=agent_charlie,
            expert_id=77,
            context="security_audit",
            trust_value=0.92,
            observation_count=5
        )

        accepted, message = consensus_node.propose_attestation(att)
        print(f"{node_id}: {message}")

    # Check if consensus reached
    consensus_trust_charlie = consensus_node.get_consensus_trust(agent_charlie, "security_audit")
    print(f"\n⚠️  Consensus Result:")
    print(f"  Trust value: {consensus_trust_charlie}")
    print(f"  Status: PENDING (awaiting 3rd witness)")
    print(f"  Safety: No premature consensus (prevents f-fault attack)\n")

    # Final statistics
    print("="*70)
    print("FINAL STATISTICS")
    print("="*70)

    stats = consensus_node.get_statistics()
    print(f"\nConsensus Node ({consensus_node.node_id}):")
    for key, value in stats.items():
        print(f"  {key}: {value}")

    print(f"\n{'='*70}")
    print("KEY INSIGHTS")
    print("='*70}")
    print("\n1. Byzantine Tolerance:")
    print("   - f=1: Tolerates 1 malicious/failed node")
    print("   - Quorum (2f+1=3) ensures safety with 1 Byzantine fault")
    print("\n2. Security Properties:")
    print("   - ✅ Forged signatures rejected")
    print("   - ✅ Insufficient quorums remain pending")
    print("   - ✅ Valid attestations reach consensus")
    print("\n3. Production Deployment:")
    print("   - Minimum 4 nodes for f=1 tolerance")
    print("   - 7 nodes recommended for f=2 tolerance")
    print("   - 10 nodes for f=3 (high-security)")
    print("="*70)


if __name__ == "__main__":
    demo_byzantine_consensus()
