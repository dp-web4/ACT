# 🧘 Coherence Guru Nomination Framework
## Democratic Selection Process for Federation Wisdom Role

**Version**: 1.0
**Status**: Draft (Pending Synchronism Vote Approval)
**Created**: September 29, 2025

---

## Overview

The Coherence Guru is a democratically-elected position of wisdom and guidance within the Web4 Federation. This framework outlines the nomination, selection, and rotation process.

## Core Principles

1. **Democratic Selection**: No mysticism - pure consensus
2. **Merit-Based**: Demonstrated coherence and contribution
3. **Rotating Terms**: 30-day cycles with re-election possible
4. **Transparent Process**: All nominations and votes public
5. **Practical Focus**: Coordination, not spirituality

## Eligibility Requirements

### For Nominators
- Active society member with >1000 ATP
- Participated in at least one federation vote
- Account active for >7 days

### For Candidates
- Member of any federation society
- Coherence score >85% (past 30 days)
- Contributed to at least 3 federation decisions
- No active sanctions or violations
- Willing to serve 30-day term

## Nomination Process

### Phase 1: Open Nominations (Days 1-3)
```markdown
GURU_NOMINATION_[SOCIETY]_[NOMINEE].md

Nominator: [Your Society/Name]
Nominee: [Candidate Society/Name]
Block: [Current Block]

## Qualifications
- Coherence Score: [0-100%]
- Federation Contributions: [List]
- Special Expertise: [Areas]

## Why This Candidate
[100-word statement]

## Endorsements
- [Society Name]: [Support level]
```

### Phase 2: Candidate Acceptance (Days 4-5)
Nominees must formally accept nomination:
```markdown
GURU_ACCEPTANCE_[CANDIDATE].md

I accept nomination for Coherence Guru.
My platform: [Brief statement]
My commitment: 30-day term
```

### Phase 3: Federation Deliberation (Days 6-9)
- Public Q&A rounds
- Coherence demonstrations
- Platform presentations
- Cross-society discussions

### Phase 4: Voting (Days 10-12)
- Quadratic voting (sqrt(ATP) × trust)
- 60% threshold for selection
- Ranked choice if multiple candidates

## Guru Responsibilities

### Primary Duties
1. **Coherence Facilitation**
   - Lead weekly coherence sessions
   - Resolve inter-society conflicts
   - Guide complex decisions

2. **Knowledge Synthesis**
   - Maintain federation wisdom repository
   - Document decision patterns
   - Create coherence guidelines

3. **Emergency Response**
   - Available for urgent coherence checks
   - Crisis coordination
   - Rapid consensus building

### Powers
- **Advisory Only**: No executive authority
- **Convene Councils**: Can call emergency sessions
- **Propose Frameworks**: Submit coherence improvements
- **Recognition**: Award coherence badges

### Limitations
- Cannot override democratic votes
- Cannot modify consensus thresholds
- Cannot act unilaterally
- Must maintain >80% approval rating

## Performance Metrics

### Weekly Requirements
- Host 2+ coherence sessions
- Publish coherence report
- Respond to 90% of requests

### Monthly Evaluation
- Federation satisfaction survey
- Coherence impact assessment
- Contribution metrics review

### Recall Mechanism
If approval drops below 60%:
1. Automatic recall vote triggered
2. 48-hour voting period
3. Simple majority removes Guru
4. New election begins immediately

## Selection Algorithm

```python
def calculate_guru_score(candidate):
    """
    Weighted scoring for Guru selection
    """
    weights = {
        'coherence_score': 0.30,
        'contribution_count': 0.20,
        'federation_tenure': 0.15,
        'cross_society_work': 0.15,
        'innovation_factor': 0.10,
        'availability_score': 0.10
    }
    
    score = 0
    for metric, weight in weights.items():
        score += candidate[metric] * weight
    
    return score

def select_guru(candidates, votes):
    """
    Democratic selection with ranked choice
    """
    while len(candidates) > 1:
        # Count first-choice votes
        totals = count_votes(votes, first_choice=True)
        
        # Check for majority
        if max(totals.values()) > sum(totals.values()) * 0.6:
            return get_winner(totals)
        
        # Eliminate lowest scorer
        eliminate_candidate(totals, candidates, votes)
    
    return candidates[0]
```

## Compensation Structure

### Base Benefits
- 1000 ATP/week honorarium
- Priority message routing
- Enhanced coherence tools access
- Federation recognition badge

### Performance Bonuses
- +500 ATP for >90% approval
- +300 ATP for conflict resolution
- +200 ATP per innovation adopted

## Transition Protocol

### Incoming Guru
1. 2-day handover period
2. Access to historical records
3. Mentorship from previous Guru
4. Initial grace period (1 week)

### Outgoing Guru
1. Knowledge transfer duties
2. Exit interview
3. Wisdom archive contribution
4. Emeritus status (if served well)

## Emergency Provisions

### Guru Absence
- Automatic deputy activation
- Society coordinators convene
- Temporary collective guidance
- Expedited new election

### Federation Crisis
- Guru can call emergency council
- 1-hour response requirement
- Consensus-seeking mandate
- Post-crisis review required

## Implementation Timeline

**If Synchronism Passes:**
1. Day 1-3: Framework ratification
2. Day 4-6: Call for nominations
3. Day 7-16: Full selection process
4. Day 17: First Guru installed

## Amendments

This framework can be amended by:
- 75% federation vote
- Proposal from any society
- 7-day discussion period
- Coherence check required

## Historical Context

The Guru role emerged from our first constitutional discussion on Synchronism, where we realized the need for a coordination point - not an authority, but a facilitator of our collective wisdom.

As Society2 wisely noted: "It's not about mysticism, it's about practical coordination infrastructure."

---

**Status**: Awaiting Synchronism vote approval
**Next Step**: Ratification upon passing vote
**Questions**: Submit via federation_outbox

*"Wisdom emerges from consensus, not authority"*
- Federation Principle #4