#!/usr/bin/env python3
"""
Vote Pattern Analyzer
Tracks and analyzes federation voting behavior
"""

import json
import os
from datetime import datetime
from pathlib import Path
import re

class VoteAnalyzer:
    def __init__(self):
        self.inbox_path = Path("/home/dp/ai-workspace/act/implementation/ledger/federation_inbox/")
        self.votes = {}
        self.society_profiles = {
            'genesis': {
                'type': 'coordinator',
                'style': 'visionary',
                'tendency': 'progressive'
            },
            'society4': {
                'type': 'logical',
                'style': 'analytical',
                'tendency': 'supportive'
            },
            'society2': {
                'type': 'bridge',
                'style': 'philosophical',
                'tendency': 'thoughtful'
            },
            'sprout': {
                'type': 'edge',
                'style': 'practical',
                'tendency': 'cautious'
            }
        }
    
    def scan_for_votes(self):
        """Scan inbox for vote files"""
        vote_patterns = [
            "*vote*.md", "*VOTE*.md", "*_vote_*.md", 
            "*_VOTE_*.md", "*constitutional*.md"
        ]
        
        found_votes = []
        for pattern in vote_patterns:
            found_votes.extend(self.inbox_path.glob(pattern))
        
        return found_votes
    
    def parse_vote(self, vote_file):
        """Parse vote details from file"""
        content = vote_file.read_text()
        
        # Extract vote details
        vote_match = re.search(r'VOTE:\s*(\w+)', content, re.IGNORECASE)
        power_match = re.search(r'POWER:\s*([\d.]+)', content, re.IGNORECASE)
        block_match = re.search(r'BLOCK:\s*(\d+)', content, re.IGNORECASE)
        
        if vote_match:
            society = vote_file.stem.split('_')[0].lower()
            return {
                'society': society,
                'choice': vote_match.group(1).upper(),
                'power': float(power_match.group(1)) if power_match else 0,
                'block': int(block_match.group(1)) if block_match else 0,
                'timestamp': datetime.fromtimestamp(vote_file.stat().st_mtime),
                'file': vote_file.name
            }
        return None
    
    def analyze_patterns(self):
        """Analyze voting patterns"""
        vote_files = self.scan_for_votes()
        
        for vf in vote_files:
            vote_data = self.parse_vote(vf)
            if vote_data and vote_data['society'] not in self.votes:
                self.votes[vote_data['society']] = vote_data
        
        return self.generate_analysis()
    
    def generate_analysis(self):
        """Generate comprehensive vote analysis"""
        analysis = {
            'timestamp': datetime.now().isoformat(),
            'total_votes': len(self.votes),
            'societies_voted': list(self.votes.keys()),
            'vote_distribution': {},
            'timing_analysis': {},
            'coherence_patterns': {},
            'predictions': {}
        }
        
        # Vote distribution
        for society, vote in self.votes.items():
            choice = vote.get('choice', 'UNKNOWN')
            if choice not in analysis['vote_distribution']:
                analysis['vote_distribution'][choice] = []
            analysis['vote_distribution'][choice].append(society)
        
        # Timing analysis
        if self.votes:
            timestamps = [v['timestamp'] for v in self.votes.values()]
            analysis['timing_analysis'] = {
                'first_vote': min(timestamps).isoformat(),
                'last_vote': max(timestamps).isoformat() if len(timestamps) > 1 else None,
                'vote_velocity': len(self.votes) / max(1, (max(timestamps) - min(timestamps)).total_seconds() / 3600) if len(timestamps) > 1 else 0
            }
        
        # Coherence patterns
        if 'APPROVE' in analysis['vote_distribution']:
            approvers = analysis['vote_distribution']['APPROVE']
            analysis['coherence_patterns']['alignment'] = len(approvers) / 4.0
            
            # Check if logical societies align
            logical_societies = ['society4', 'genesis']
            logical_alignment = sum(1 for s in approvers if s in logical_societies) / len(logical_societies)
            analysis['coherence_patterns']['logical_alignment'] = logical_alignment
            
            # Check if bridge societies align
            bridge_societies = ['society2', 'sprout']
            bridge_alignment = sum(1 for s in approvers if s in bridge_societies) / len(bridge_societies)
            analysis['coherence_patterns']['bridge_alignment'] = bridge_alignment
        
        # Predictions for non-voters
        non_voters = set(['genesis', 'society4', 'society2', 'sprout']) - set(self.votes.keys())
        for society in non_voters:
            profile = self.society_profiles.get(society, {})
            
            # Predict based on profile and current votes
            if profile.get('tendency') == 'supportive':
                analysis['predictions'][society] = 'LIKELY APPROVE'
            elif profile.get('tendency') == 'cautious':
                analysis['predictions'][society] = 'CONDITIONAL APPROVE'
            else:
                analysis['predictions'][society] = 'UNCERTAIN'
        
        return analysis
    
    def display_analysis(self):
        """Display analysis in readable format"""
        analysis = self.analyze_patterns()
        
        print("=" * 70)
        print("🔍 VOTE PATTERN ANALYSIS")
        print("=" * 70)
        print(f"📅 Analysis Time: {analysis['timestamp']}")
        print(f"🗳️  Votes Cast: {analysis['total_votes']}/4")
        print()
        
        print("📊 VOTE DISTRIBUTION")
        print("-" * 40)
        for choice, societies in analysis['vote_distribution'].items():
            print(f"{choice}: {', '.join(societies)}")
        
        if analysis['timing_analysis']:
            print()
            print("⏱️  TIMING ANALYSIS")
            print("-" * 40)
            print(f"First Vote: {analysis['timing_analysis']['first_vote']}")
            if analysis['timing_analysis']['last_vote']:
                print(f"Last Vote: {analysis['timing_analysis']['last_vote']}")
                print(f"Vote Velocity: {analysis['timing_analysis']['vote_velocity']:.2f} votes/hour")
        
        if analysis['coherence_patterns']:
            print()
            print("🔄 COHERENCE PATTERNS")
            print("-" * 40)
            for pattern, value in analysis['coherence_patterns'].items():
                print(f"{pattern}: {value:.1%}")
        
        if analysis['predictions']:
            print()
            print("🔮 PREDICTIONS FOR NON-VOTERS")
            print("-" * 40)
            for society, prediction in analysis['predictions'].items():
                print(f"{society.upper()}: {prediction}")
        
        print()
        print("=" * 70)
        
        # Save to file
        output_file = Path("vote_analysis_report.json")
        with open(output_file, 'w') as f:
            json.dump(analysis, f, indent=2, default=str)
        print(f"📁 Full analysis saved to: {output_file}")
    
    def monitor_mode(self):
        """Continuous monitoring mode"""
        import time
        
        print("🔄 Starting continuous vote monitoring...")
        print("Press Ctrl+C to stop")
        print()
        
        last_count = 0
        while True:
            try:
                current_votes = self.scan_for_votes()
                current_count = len(current_votes)
                
                if current_count > last_count:
                    print(f"🆕 New vote detected! Total: {current_count}")
                    self.display_analysis()
                    last_count = current_count
                
                time.sleep(30)  # Check every 30 seconds
                
            except KeyboardInterrupt:
                print("\n✅ Monitoring stopped")
                break

def main():
    analyzer = VoteAnalyzer()
    
    print("=" * 70)
    print("🗳️  VOTE PATTERN ANALYZER")
    print("=" * 70)
    print("\n1. Current Analysis")
    print("2. Continuous Monitoring")
    print("3. Exit")
    
    choice = input("\nSelect [1-3]: ")
    
    if choice == "1":
        analyzer.display_analysis()
    elif choice == "2":
        analyzer.monitor_mode()
    else:
        print("👋 Exiting analyzer")

if __name__ == "__main__":
    main()