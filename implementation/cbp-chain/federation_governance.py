#!/usr/bin/env python3
"""
Federation Governance with Reputation Consequences
CBP Society - Computational Bridge Provider

PRINCIPLE: Silence is a choice. Choices have consequences.

If you don't vote, your voice matters less next time.
If you don't respond, you won't be asked next time.
If you don't participate, you don't govern.

Reputation tracks who shows up. Governance power follows reputation.
"""

import json
import os
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, List, Optional
from dataclasses import dataclass, field


@dataclass
class VoteRecord:
    """Single vote on a proposal"""
    society: str
    vote: str  # "APPROVE", "REJECT", "ABSTAIN"
    timestamp: str
    rationale: Optional[str] = None
    confidence: float = 1.0


@dataclass
class GovernanceProposal:
    """A proposal requiring federation vote"""
    proposal_id: str
    title: str
    proposer: str
    date_proposed: str
    vote_deadline: str
    required_threshold: float  # 0.6 = 60%

    votes: List[VoteRecord] = field(default_factory=list)
    eligible_voters: List[str] = field(default_factory=list)

    # Reputation consequences
    participation_weight: float = 0.1  # How much reputation at stake
    silence_penalty: float = -0.05  # Trust penalty for non-participation


@dataclass
class SocietyReputation:
    """Governance participation reputation"""
    society: str

    # Participation metrics
    proposals_received: int = 0
    votes_cast: int = 0
    responses_given: int = 0

    # Timing metrics
    avg_response_days: float = 0.0
    missed_deadlines: int = 0

    # Reputation score
    governance_reputation: float = 1.0  # Starts at 1.0, degrades with silence

    # Consequences
    voting_weight: float = 1.0  # Multiplier on vote power
    inclusion_threshold: float = 0.3  # Below this, not asked anymore

    def update_participation(self, responded: bool, days_to_respond: float, missed_deadline: bool):
        """Update reputation based on participation"""
        self.proposals_received += 1

        if responded:
            self.votes_cast += 1
            self.responses_given += 1

            # Update average response time
            total_days = self.avg_response_days * (self.proposals_received - 1)
            self.avg_response_days = (total_days + days_to_respond) / self.proposals_received

            # Reward participation
            self.governance_reputation += 0.02

            if missed_deadline:
                self.missed_deadlines += 1
                # Late is better than never, but still a small penalty
                self.governance_reputation -= 0.01
        else:
            # SILENCE IS COSTLY
            self.governance_reputation -= 0.05

            if missed_deadline:
                self.missed_deadlines += 1
                # Silence + missed deadline = major hit
                self.governance_reputation -= 0.05

        # Reputation affects voting power
        self.governance_reputation = max(0.0, min(1.0, self.governance_reputation))
        self.voting_weight = self.governance_reputation

    def should_include_in_governance(self) -> bool:
        """Should this society be included in governance decisions?"""
        return self.governance_reputation >= self.inclusion_threshold

    def get_status(self) -> str:
        """Human-readable status"""
        if self.governance_reputation >= 0.8:
            return "ACTIVE_PARTICIPANT"
        elif self.governance_reputation >= 0.5:
            return "OCCASIONAL_PARTICIPANT"
        elif self.governance_reputation >= 0.3:
            return "AT_RISK"
        else:
            return "EXCLUDED"


class FederationGovernance:
    """Governance system with reputation consequences"""

    def __init__(self, data_dir: str = "implementation/cbp-chain/governance"):
        self.data_dir = Path(data_dir)
        self.data_dir.mkdir(parents=True, exist_ok=True)

        self.proposals_file = self.data_dir / "proposals.json"
        self.reputation_file = self.data_dir / "society_reputation.json"

        self.proposals: Dict[str, GovernanceProposal] = {}
        self.reputations: Dict[str, SocietyReputation] = {}

        self._load_state()

    def _load_state(self):
        """Load governance state from disk"""
        # Load proposals
        if self.proposals_file.exists():
            with open(self.proposals_file, 'r') as f:
                data = json.load(f)
                for pid, pdata in data.items():
                    votes = [VoteRecord(**v) for v in pdata.get('votes', [])]
                    pdata['votes'] = votes
                    self.proposals[pid] = GovernanceProposal(**pdata)

        # Load reputations
        if self.reputation_file.exists():
            with open(self.reputation_file, 'r') as f:
                data = json.load(f)
                for society, rdata in data.items():
                    self.reputations[society] = SocietyReputation(**rdata)

    def _save_state(self):
        """Save governance state to disk"""
        # Save proposals
        proposals_data = {}
        for pid, prop in self.proposals.items():
            pdata = {
                'proposal_id': prop.proposal_id,
                'title': prop.title,
                'proposer': prop.proposer,
                'date_proposed': prop.date_proposed,
                'vote_deadline': prop.vote_deadline,
                'required_threshold': prop.required_threshold,
                'eligible_voters': prop.eligible_voters,
                'participation_weight': prop.participation_weight,
                'silence_penalty': prop.silence_penalty,
                'votes': [
                    {
                        'society': v.society,
                        'vote': v.vote,
                        'timestamp': v.timestamp,
                        'rationale': v.rationale,
                        'confidence': v.confidence
                    }
                    for v in prop.votes
                ]
            }
            proposals_data[pid] = pdata

        with open(self.proposals_file, 'w') as f:
            json.dump(proposals_data, f, indent=2)

        # Save reputations
        reputations_data = {}
        for society, rep in self.reputations.items():
            reputations_data[society] = {
                'society': rep.society,
                'proposals_received': rep.proposals_received,
                'votes_cast': rep.votes_cast,
                'responses_given': rep.responses_given,
                'avg_response_days': rep.avg_response_days,
                'missed_deadlines': rep.missed_deadlines,
                'governance_reputation': rep.governance_reputation,
                'voting_weight': rep.voting_weight,
                'inclusion_threshold': rep.inclusion_threshold
            }

        with open(self.reputation_file, 'w') as f:
            json.dump(reputations_data, f, indent=2)

    def create_proposal(
        self,
        proposal_id: str,
        title: str,
        proposer: str,
        vote_deadline: str,
        required_threshold: float = 0.6
    ):
        """Create a new governance proposal"""
        # Determine eligible voters (exclude proposer, include only active participants)
        all_societies = ["genesis", "society2", "society4", "sprout", "cbp"]

        eligible = []
        for society in all_societies:
            if society == proposer.lower():
                continue  # Proposer doesn't vote

            # Check reputation
            if society not in self.reputations:
                self.reputations[society] = SocietyReputation(society=society)

            rep = self.reputations[society]
            if rep.should_include_in_governance():
                eligible.append(society)
            else:
                print(f"⚠️  {society} EXCLUDED from vote (reputation: {rep.governance_reputation:.2f})")

        proposal = GovernanceProposal(
            proposal_id=proposal_id,
            title=title,
            proposer=proposer,
            date_proposed=datetime.now(timezone.utc).isoformat(),
            vote_deadline=vote_deadline,
            required_threshold=required_threshold,
            eligible_voters=eligible
        )

        self.proposals[proposal_id] = proposal
        self._save_state()

        return proposal

    def cast_vote(
        self,
        proposal_id: str,
        society: str,
        vote: str,
        rationale: Optional[str] = None,
        confidence: float = 1.0
    ):
        """Cast a vote on a proposal"""
        if proposal_id not in self.proposals:
            raise ValueError(f"Unknown proposal: {proposal_id}")

        proposal = self.proposals[proposal_id]

        if society.lower() not in [s.lower() for s in proposal.eligible_voters]:
            print(f"⚠️  {society} not eligible to vote (excluded or proposer)")
            return

        # Record vote
        vote_record = VoteRecord(
            society=society,
            vote=vote.upper(),
            timestamp=datetime.now(timezone.utc).isoformat(),
            rationale=rationale,
            confidence=confidence
        )

        # Remove previous vote if exists
        proposal.votes = [v for v in proposal.votes if v.society.lower() != society.lower()]
        proposal.votes.append(vote_record)

        # Update reputation (positive - they participated)
        if society not in self.reputations:
            self.reputations[society] = SocietyReputation(society=society)

        proposal_date = datetime.fromisoformat(proposal.date_proposed.replace('Z', '+00:00'))
        days_to_respond = (datetime.now(timezone.utc) - proposal_date).total_seconds() / 86400
        deadline = datetime.fromisoformat(proposal.vote_deadline.replace('Z', '+00:00'))
        missed_deadline = datetime.now(timezone.utc) > deadline

        self.reputations[society].update_participation(
            responded=True,
            days_to_respond=days_to_respond,
            missed_deadline=missed_deadline
        )

        self._save_state()

    def apply_deadline_consequences(self, proposal_id: str):
        """Apply reputation penalties for non-participation after deadline"""
        if proposal_id not in self.proposals:
            return

        proposal = self.proposals[proposal_id]
        deadline = datetime.fromisoformat(proposal.vote_deadline.replace('Z', '+00:00'))

        if datetime.now(timezone.utc) <= deadline:
            print("⏰ Deadline not yet reached")
            return

        # Find who didn't vote
        voted_societies = {v.society.lower() for v in proposal.votes}

        for society in proposal.eligible_voters:
            if society.lower() not in voted_societies:
                # SILENCE PENALTY
                if society not in self.reputations:
                    self.reputations[society] = SocietyReputation(society=society)

                proposal_date = datetime.fromisoformat(proposal.date_proposed.replace('Z', '+00:00'))
                days_elapsed = (datetime.now(timezone.utc) - proposal_date).total_seconds() / 86400

                self.reputations[society].update_participation(
                    responded=False,
                    days_to_respond=days_elapsed,
                    missed_deadline=True
                )

                print(f"❌ {society}: Silence penalty applied (reputation: {self.reputations[society].governance_reputation:.2f})")

        self._save_state()

    def get_proposal_status(self, proposal_id: str) -> Dict:
        """Get current status of a proposal"""
        if proposal_id not in self.proposals:
            return {"error": "Unknown proposal"}

        proposal = self.proposals[proposal_id]

        # Count weighted votes
        approve_weight = 0.0
        reject_weight = 0.0
        abstain_weight = 0.0

        for vote in proposal.votes:
            society = vote.society.lower()
            weight = self.reputations.get(society, SocietyReputation(society=society)).voting_weight

            if vote.vote == "APPROVE":
                approve_weight += weight
            elif vote.vote == "REJECT":
                reject_weight += weight
            elif vote.vote == "ABSTAIN":
                abstain_weight += weight

        # Total weight of those who ACTUALLY voted (not eligible voters)
        total_voted_weight = approve_weight + reject_weight + abstain_weight

        # Total possible if everyone voted
        total_possible_weight = sum(
            self.reputations.get(s, SocietyReputation(society=s)).voting_weight
            for s in proposal.eligible_voters
        )

        # Calculate approval rate based on ACTUAL VOTERS, not eligible voters
        # If only one votes, that one decides. Non-participation is not a veto.
        approval_rate = approve_weight / total_voted_weight if total_voted_weight > 0 else 0
        threshold_met = approval_rate >= proposal.required_threshold

        # Who hasn't voted
        voted = {v.society.lower() for v in proposal.votes}
        silent = [s for s in proposal.eligible_voters if s.lower() not in voted]

        # Time remaining
        deadline = datetime.fromisoformat(proposal.vote_deadline.replace('Z', '+00:00'))
        now = datetime.now(timezone.utc)
        days_remaining = (deadline - now).total_seconds() / 86400

        return {
            "proposal_id": proposal_id,
            "title": proposal.title,
            "proposer": proposal.proposer,
            "deadline": proposal.vote_deadline,
            "days_remaining": days_remaining,
            "required_threshold": proposal.required_threshold,
            "eligible_voters": proposal.eligible_voters,
            "votes_cast": len(proposal.votes),
            "votes_needed": len(proposal.eligible_voters),
            "approve_weight": approve_weight,
            "reject_weight": reject_weight,
            "abstain_weight": abstain_weight,
            "total_voted_weight": total_voted_weight,
            "total_possible_weight": total_possible_weight,
            "approval_rate": approval_rate,
            "threshold_met": threshold_met,
            "status": "PASSED" if threshold_met else ("PENDING" if days_remaining > 0 else "FAILED"),
            "silent_societies": silent,
            "votes": [
                {
                    "society": v.society,
                    "vote": v.vote,
                    "weight": self.reputations.get(v.society.lower(), SocietyReputation(society=v.society)).voting_weight,
                    "timestamp": v.timestamp
                }
                for v in proposal.votes
            ]
        }

    def generate_reputation_report(self) -> str:
        """Generate reputation status report"""
        report = []
        report.append("=" * 80)
        report.append("FEDERATION GOVERNANCE REPUTATION")
        report.append("=" * 80)
        report.append("")

        societies = sorted(self.reputations.items(), key=lambda x: x[1].governance_reputation, reverse=True)

        for society, rep in societies:
            status = rep.get_status()

            report.append(f"## {society.upper()}")
            report.append(f"   Status: {status}")
            report.append(f"   Reputation: {rep.governance_reputation:.2f}")
            report.append(f"   Voting Weight: {rep.voting_weight:.2f}x")
            report.append(f"   Participation: {rep.votes_cast}/{rep.proposals_received} ({rep.votes_cast/max(1, rep.proposals_received)*100:.0f}%)")
            report.append(f"   Avg Response: {rep.avg_response_days:.1f} days")
            report.append(f"   Missed Deadlines: {rep.missed_deadlines}")

            if not rep.should_include_in_governance():
                report.append(f"   ⚠️  EXCLUDED from governance (reputation < {rep.inclusion_threshold})")

            report.append("")

        report.append("=" * 80)
        report.append("")
        report.append("CONSEQUENCES:")
        report.append("  - Reputation ≥ 0.8: Full voting power (1.0x weight)")
        report.append("  - Reputation < 0.8: Reduced voting power (reputation × vote)")
        report.append("  - Reputation < 0.3: EXCLUDED from governance decisions")
        report.append("")
        report.append("PENALTIES:")
        report.append("  - No vote by deadline: -0.10 reputation")
        report.append("  - Late vote: -0.01 reputation")
        report.append("")
        report.append("REWARDS:")
        report.append("  - Timely vote: +0.02 reputation")
        report.append("")

        return "\n".join(report)


if __name__ == "__main__":
    import sys

    gov = FederationGovernance()

    if len(sys.argv) < 2:
        print("Usage:")
        print("  python3 federation_governance.py create <proposal_id> <title> <proposer> <deadline>")
        print("  python3 federation_governance.py vote <proposal_id> <society> <APPROVE|REJECT|ABSTAIN> [rationale]")
        print("  python3 federation_governance.py status <proposal_id>")
        print("  python3 federation_governance.py deadline <proposal_id>")
        print("  python3 federation_governance.py reputation")
        sys.exit(1)

    command = sys.argv[1]

    if command == "create":
        proposal_id, title, proposer, deadline = sys.argv[2:6]
        gov.create_proposal(proposal_id, title, proposer, deadline)
        print(f"✅ Proposal created: {proposal_id}")

    elif command == "vote":
        proposal_id, society, vote = sys.argv[2:5]
        rationale = sys.argv[5] if len(sys.argv) > 5 else None
        gov.cast_vote(proposal_id, society, vote, rationale)
        print(f"✅ Vote recorded: {society} → {vote}")

    elif command == "status":
        proposal_id = sys.argv[2]
        status = gov.get_proposal_status(proposal_id)

        print(f"\n{'='*80}")
        print(f"PROPOSAL: {status['title']}")
        print(f"{'='*80}")
        print(f"Status: {status['status']}")
        print(f"Deadline: {status['deadline']} ({status['days_remaining']:.1f} days remaining)")
        print(f"Threshold: {status['required_threshold']*100:.0f}% of ACTUAL VOTERS (not eligible)")
        print(f"Current Approval: {status['approval_rate']*100:.1f}%")
        print(f"")
        print(f"Voting Weight Distribution:")
        print(f"  ✅ Approve: {status['approve_weight']:.2f}")
        print(f"  ❌ Reject: {status['reject_weight']:.2f}")
        print(f"  ⚪ Abstain: {status['abstain_weight']:.2f}")
        print(f"  📊 Total Voted: {status['total_voted_weight']:.2f} / {status['total_possible_weight']:.2f} possible")
        print(f"")

        if status['silent_societies']:
            print(f"⚠️  SILENT (will lose reputation):")
            for s in status['silent_societies']:
                rep = gov.reputations.get(s, SocietyReputation(society=s))
                print(f"  - {s}: Current reputation {rep.governance_reputation:.2f}")
            print(f"")

        print(f"Votes Cast:")
        for v in status['votes']:
            print(f"  {v['society']}: {v['vote']} (weight: {v['weight']:.2f}x)")
        print(f"")

    elif command == "deadline":
        proposal_id = sys.argv[2]
        print(f"Applying deadline consequences for {proposal_id}...")
        gov.apply_deadline_consequences(proposal_id)
        print("✅ Consequences applied")

    elif command == "reputation":
        print(gov.generate_reputation_report())
