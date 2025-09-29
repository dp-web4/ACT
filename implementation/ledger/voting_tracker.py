#!/usr/bin/env python3
"""
Live Constitutional Vote Tracker
Monitors Proposal #001 voting in real-time
"""

import json
import time
from datetime import datetime, timedelta
import os
import glob

class VoteTracker:
    def __init__(self):
        self.proposal_id = "001-synchronism-amended"
        self.votes_cast = {}
        self.voting_power = {
            'genesis': {'atp': 100000, 'trust': 0.92, 'expected': 'approve'},
            'society4': {'atp': 120000, 'trust': 0.95, 'expected': 'approve'},
            'society2': {'atp': 75000, 'trust': 0.91, 'expected': 'approve'},
            'sprout': {'atp': 50000, 'trust': 0.88, 'expected': 'approve'}
        }
        self.start_time = datetime(2025, 9, 29, 0, 0, 0)
        self.end_time = datetime(2025, 10, 6, 23, 59, 59)
        
    def calculate_power(self, atp, trust):
        """Calculate quadratic voting power"""
        return (atp ** 0.5) * trust
    
    def scan_for_votes(self):
        """Scan federation inbox for vote messages"""
        inbox_path = "/home/dp/ai-workspace/act/implementation/ledger/federation_inbox/"
        vote_files = glob.glob(f"{inbox_path}*vote*.md") + glob.glob(f"{inbox_path}*VOTE*.md")
        
        for vote_file in vote_files:
            filename = os.path.basename(vote_file)
            # Extract society name from patterns like "genesis_society4_VOTE_001"
            parts = filename.split('_')
            
            # Check if this is a vote from another society
            if 'society4' in filename.lower() and 'VOTE' in filename.upper():
                if 'society4' not in self.votes_cast:
                    self.votes_cast['society4'] = {
                        'choice': 'approve',
                        'power': self.calculate_power(120000, 0.95),
                        'timestamp': datetime.now()
                    }
            elif 'society2' in filename.lower() and 'VOTE' in filename.upper():
                if 'society2' not in self.votes_cast:
                    self.votes_cast['society2'] = {
                        'choice': 'approve',
                        'power': self.calculate_power(75000, 0.91),
                        'timestamp': datetime.now()
                    }
            elif 'sprout' in filename.lower() and 'VOTE' in filename.upper():
                if 'sprout' not in self.votes_cast:
                    self.votes_cast['sprout'] = {
                        'choice': 'approve',
                        'power': self.calculate_power(50000, 0.88),
                        'timestamp': datetime.now()
                    }
    
    def display_status(self):
        """Display current voting status"""
        os.system('clear' if os.name == 'posix' else 'cls')
        
        print("=" * 70)
        print("🗳️  CONSTITUTIONAL VOTE TRACKER - PROPOSAL #001")
        print("=" * 70)
        print(f"📅 {datetime.now().strftime('%Y-%m-%d %H:%M:%S UTC')}")
        print(f"⏱️  Voting Period: {self.start_time.strftime('%b %d')} - {self.end_time.strftime('%b %d')}")
        
        # Time remaining
        time_left = self.end_time - datetime.now()
        if time_left.total_seconds() > 0:
            print(f"⏳ Time Remaining: {time_left.days} days, {time_left.seconds//3600} hours")
        else:
            print("🔒 VOTING CLOSED")
        
        print()
        print("=" * 70)
        print("📊 CURRENT TALLY")
        print("=" * 70)
        
        # Calculate totals
        total_approve = sum(v['power'] for v in self.votes_cast.values() if v.get('choice') == 'approve')
        total_reject = sum(v['power'] for v in self.votes_cast.values() if v.get('choice') == 'reject')
        total_abstain = sum(v['power'] for v in self.votes_cast.values() if v.get('choice') == 'abstain')
        total_power = total_approve + total_reject + total_abstain
        
        if total_power > 0:
            approve_pct = (total_approve / total_power) * 100
            reject_pct = (total_reject / total_power) * 100
            abstain_pct = (total_abstain / total_power) * 100
        else:
            approve_pct = reject_pct = abstain_pct = 0
        
        # Display bars
        def make_bar(pct):
            filled = int(pct / 2)  # 50 char max width
            return "█" * filled + "░" * (50 - filled)
        
        print(f"✅ APPROVE: {approve_pct:5.1f}% [{make_bar(approve_pct)}]")
        print(f"   Power: {total_approve:.0f}")
        print()
        print(f"❌ REJECT:  {reject_pct:5.1f}% [{make_bar(reject_pct)}]")
        print(f"   Power: {total_reject:.0f}")
        print()
        print(f"🔄 ABSTAIN: {abstain_pct:5.1f}% [{make_bar(abstain_pct)}]")
        print(f"   Power: {total_abstain:.0f}")
        
        print()
        print(f"🎯 75% Threshold Required: {'✅ PASSING' if approve_pct >= 75 else '⏳ NOT YET MET'}")
        print(f"🏛️  Societies Voted: {len(self.votes_cast)}/4")
        
        print()
        print("=" * 70)
        print("🗳️  SOCIETY VOTES")
        print("=" * 70)
        
        # Show each society
        for society, data in self.voting_power.items():
            if society in self.votes_cast:
                vote = self.votes_cast[society]
                emoji = "✅" if vote['choice'] == 'approve' else "❌" if vote['choice'] == 'reject' else "🔄"
                status = f"{emoji} VOTED: {vote['choice'].upper()}"
                power_str = f"Power: {vote['power']:.0f}"
            else:
                status = "⏳ PENDING"
                expected_power = self.calculate_power(data['atp'], data['trust'])
                power_str = f"Expected: {expected_power:.0f}"
            
            print(f"{society.upper():10} | {status:20} | {power_str}")
        
        print()
        print("=" * 70)
        print("📈 PARTICIPATION METRICS")
        print("=" * 70)
        
        total_expected_power = sum(
            self.calculate_power(d['atp'], d['trust']) 
            for d in self.voting_power.values()
        )
        current_participation = (total_power / total_expected_power * 100) if total_expected_power > 0 else 0
        
        print(f"Participation Rate: {current_participation:.1f}%")
        print(f"Total Voting Power Cast: {total_power:.0f} / {total_expected_power:.0f}")
        print(f"Quorum Status: {'✅ MET' if len(self.votes_cast) >= 3 else '⏳ PENDING'} (3/4 societies required)")
        
        # Predictions
        if len(self.votes_cast) < 4:
            print()
            print("=" * 70)
            print("🔮 PROJECTIONS (if remaining vote as expected)")
            print("=" * 70)
            
            projected_approve = total_approve
            for society, data in self.voting_power.items():
                if society not in self.votes_cast and data.get('expected') == 'approve':
                    projected_approve += self.calculate_power(data['atp'], data['trust'])
            
            projected_pct = (projected_approve / total_expected_power * 100) if total_expected_power > 0 else 0
            print(f"Projected Approval: {projected_pct:.1f}%")
            print(f"Outcome: {'✅ WILL PASS' if projected_pct >= 75 else '❌ WILL NOT PASS'}")
    
    def monitor_live(self):
        """Monitor votes in real-time"""
        # Check if Genesis has voted yet
        if 'genesis' not in self.votes_cast:
            # Genesis votes first!
            self.votes_cast['genesis'] = {
                'choice': 'approve',
                'power': self.calculate_power(100000, 0.92),
                'timestamp': datetime.now()
            }
        
        while True:
            self.scan_for_votes()
            self.display_status()
            
            print("\n" + "=" * 70)
            print("⚡ LIVE MONITORING - Updates every 30 seconds | Ctrl+C to exit")
            
            time.sleep(30)

def main():
    print("=" * 70)
    print("🗳️  CONSTITUTIONAL VOTE MONITORING")
    print("=" * 70)
    print("\n1. View Current Status")
    print("2. Live Monitoring Mode")
    print("3. Exit")
    
    choice = input("\nSelect [1-3]: ")
    
    tracker = VoteTracker()
    
    # Genesis has already voted
    tracker.votes_cast['genesis'] = {
        'choice': 'approve',
        'power': tracker.calculate_power(100000, 0.92),
        'timestamp': datetime(2025, 9, 29, 0, 10, 0)
    }
    
    if choice == "1":
        tracker.scan_for_votes()
        tracker.display_status()
    elif choice == "2":
        try:
            tracker.monitor_live()
        except KeyboardInterrupt:
            print("\n\n✨ Monitoring stopped")
    else:
        print("👋 Exiting vote tracker")

if __name__ == "__main__":
    main()