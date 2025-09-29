#!/usr/bin/env python3
"""
Quick Coherence Check Tool
20-Second Decision Framework for Synchronism
Federation Constitutional Tool v1.0
"""

import time
import json
from datetime import datetime

class QuickCoherenceCheck:
    def __init__(self):
        self.checks = {
            'fractal': (0, "Fractal Alignment (pattern similarity)"),
            'intent': (0, "Intent Clarity (purpose definition)"),
            'spectrum': (0, "Spectral Position (state awareness)"),
            'boundary': (0, "Relevance Boundary (scope clarity)")
        }
        
    def run_check(self, decision_name=""):
        """Run a 20-second coherence check on a decision"""
        print("\n" + "="*50)
        print("🌌 QUICK COHERENCE CHECK - 20 Second Framework")
        print("="*50)
        
        if not decision_name:
            decision_name = input("Decision/Action to evaluate: ")
        
        print(f"\nEvaluating: {decision_name}")
        print("Rate each dimension 1-10 (5 seconds each):\n")
        
        start_time = time.time()
        
        # Collect ratings
        for key, (_, description) in self.checks.items():
            timer_start = time.time()
            
            while True:
                try:
                    rating = int(input(f"📊 {description} [1-10]: "))
                    if 1 <= rating <= 10:
                        self.checks[key] = (rating, description)
                        break
                    else:
                        print("Please enter 1-10")
                except ValueError:
                    print("Please enter a number 1-10")
            
            elapsed = time.time() - timer_start
            if elapsed < 5:
                print(f"   ⏱️  {5-elapsed:.1f}s remaining...")
                time.sleep(max(0, 5-elapsed))
        
        # Calculate coherence
        scores = [v[0] for v in self.checks.values()]
        
        # Multiplicative coherence (normalized)
        multiplicative = 1.0
        for score in scores:
            multiplicative *= (score / 10.0)
        
        # Average coherence
        average = sum(scores) / len(scores) / 10.0
        
        # Combined coherence (weighted)
        coherence = (multiplicative * 0.6) + (average * 0.4)
        
        total_time = time.time() - start_time
        
        # Display results
        print("\n" + "="*50)
        print("📈 COHERENCE ANALYSIS COMPLETE")
        print("="*50)
        
        print(f"\n🎯 Decision: {decision_name}")
        print(f"⏱️  Time taken: {total_time:.1f} seconds")
        
        print("\n📊 Dimensional Scores:")
        for key, (score, desc) in self.checks.items():
            bar = "█" * score + "░" * (10-score)
            print(f"  {desc[:20]:20} [{bar}] {score}/10")
        
        print(f"\n🔮 Coherence Metrics:")
        print(f"  Multiplicative: {multiplicative:.3f}")
        print(f"  Average: {average:.3f}")
        print(f"  Combined: {coherence:.3f}")
        
        # Interpretation
        print(f"\n💭 Interpretation:")
        if coherence >= 0.7:
            print("  ✅ HIGH COHERENCE - Strong alignment with Synchronism")
            print("  Recommendation: Proceed with confidence")
        elif coherence >= 0.5:
            print("  🔄 MODERATE COHERENCE - Acceptable alignment")
            print("  Recommendation: Proceed with attention to weak areas")
        else:
            print("  ⚠️  LOW COHERENCE - Weak alignment detected")
            print("  Recommendation: Reconsider or refine approach")
        
        # Save result
        self.save_result(decision_name, coherence, scores, total_time)
        
        return coherence
    
    def save_result(self, decision, coherence, scores, time_taken):
        """Save coherence check to federation log"""
        result = {
            'timestamp': datetime.now().isoformat(),
            'decision': decision,
            'coherence': coherence,
            'scores': {
                'fractal': scores[0],
                'intent': scores[1],
                'spectrum': scores[2],
                'boundary': scores[3]
            },
            'time_seconds': round(time_taken, 1),
            'society': 'genesis'
        }
        
        # Append to log file
        try:
            with open('coherence_checks.jsonl', 'a') as f:
                f.write(json.dumps(result) + '\n')
            print(f"\n💾 Result saved to coherence_checks.jsonl")
        except:
            print(f"\n⚠️  Could not save result to file")
        
        return result

def main():
    print("="*50)
    print("🌌 SYNCHRONISM QUICK COHERENCE CHECK")
    print("Federation Constitutional Tool v1.0")
    print("="*50)
    
    checker = QuickCoherenceCheck()
    
    while True:
        print("\nOptions:")
        print("1. Run coherence check")
        print("2. View recent checks")
        print("3. Exit")
        
        choice = input("\nSelect [1-3]: ")
        
        if choice == "1":
            checker.run_check()
        elif choice == "2":
            try:
                print("\n📋 Recent Coherence Checks:")
                with open('coherence_checks.jsonl', 'r') as f:
                    lines = f.readlines()[-5:]  # Last 5 checks
                    for line in lines:
                        data = json.loads(line)
                        print(f"  {data['timestamp']}: {data['decision']} = {data['coherence']:.3f}")
            except:
                print("  No previous checks found")
        elif choice == "3":
            print("\n✨ May your decisions be coherent!\n")
            break
        else:
            print("Invalid option")

if __name__ == "__main__":
    main()