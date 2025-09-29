#!/usr/bin/env python3
"""
Pending Consensus System for Society 4
Manages decisions made while network-isolated from federation
"""

import json
import os
import subprocess
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional

# Configuration
PENDING_DIR = Path.home() / ".society4chain" / "pending_consensus"
PENDING_FILE = PENDING_DIR / "pending.json"
PROCESSED_FILE = PENDING_DIR / "processed.json"
NETWORK_STATE_FILE = PENDING_DIR / "network_state"

# Colors for terminal output
GREEN = '\033[0;32m'
YELLOW = '\033[1;33m'
RED = '\033[0;31m'
NC = '\033[0m'  # No Color

class PendingConsensus:
    def __init__(self):
        self.init_system()

    def init_system(self):
        """Initialize the pending consensus directory and files."""
        PENDING_DIR.mkdir(parents=True, exist_ok=True)

        if not PENDING_FILE.exists():
            with open(PENDING_FILE, 'w') as f:
                json.dump({"pending_decisions": [], "metadata": {}}, f, indent=2)

        if not PROCESSED_FILE.exists():
            with open(PROCESSED_FILE, 'w') as f:
                json.dump({"processed": []}, f, indent=2)

        print(f"{GREEN}✓{NC} Pending consensus system initialized")

    def detect_network(self) -> str:
        """Detect current network state."""
        try:
            # Get current IP
            result = subprocess.run(
                ["ip", "addr", "show", "eth0"],
                capture_output=True,
                text=True
            )

            ip_line = [line for line in result.stdout.split('\n') if 'inet ' in line]
            if ip_line:
                current_ip = ip_line[0].split()[1].split('/')[0]
            else:
                current_ip = "unknown"

            # Determine network type
            if current_ip.startswith("10.0.0."):
                network_type = "home_federation"
                print(f"{GREEN}✓{NC} Connected to HOME network (Federation accessible)")
            elif current_ip.startswith("172.25."):
                network_type = "work_isolated"
                print(f"{YELLOW}⚠{NC} Connected to WORK network (Federation isolated)")
            else:
                network_type = "unknown"
                print(f"{RED}✗{NC} Unknown network: {current_ip}")

            # Save network state
            with open(NETWORK_STATE_FILE, 'w') as f:
                f.write(f"{network_type}|{current_ip}|{datetime.now().isoformat()}")

            return network_type

        except Exception as e:
            print(f"{RED}Error detecting network: {e}{NC}")
            return "unknown"

    def get_hardware_hash(self) -> str:
        """Extract hardware hash."""
        try:
            script_path = "/mnt/c/projects/ai-agents/ACT/implementation/society4/blockchain/source/extract_hardware.sh"
            result = subprocess.run(
                ["bash", script_path],
                capture_output=True,
                text=True
            )

            for line in result.stdout.split('\n'):
                if "Hardware Hash" in line:
                    return line.split()[-1]
        except:
            pass
        return "unavailable"

    def add_pending_decision(self, decision_type: str, data: str, reason: str):
        """Add a pending decision."""
        network_state = self.detect_network()
        timestamp = datetime.now().isoformat()
        hardware_hash = self.get_hardware_hash()

        # Create decision ID from timestamp
        decision_id = timestamp.replace(":", "").replace("-", "").split(".")[0]

        decision = {
            "id": decision_id,
            "type": decision_type,
            "data": data,
            "reason": reason,
            "network_state": network_state,
            "timestamp": timestamp,
            "hardware_attestation": hardware_hash,
            "status": "pending"
        }

        # Load and update pending file
        with open(PENDING_FILE, 'r') as f:
            pending_data = json.load(f)

        pending_data["pending_decisions"].append(decision)

        with open(PENDING_FILE, 'w') as f:
            json.dump(pending_data, f, indent=2)

        print(f"{GREEN}✓{NC} Pending decision added: {decision_type}")
        print(json.dumps(decision, indent=2))

    def list_pending(self):
        """List all pending decisions."""
        with open(PENDING_FILE, 'r') as f:
            pending_data = json.load(f)

        pending = [d for d in pending_data["pending_decisions"] if d["status"] == "pending"]

        if not pending:
            print(f"{YELLOW}No pending decisions{NC}")
            return

        print(f"{GREEN}Pending Decisions ({len(pending)}):{NC}")
        for decision in pending:
            date = decision["timestamp"].split("T")[0]
            print(f"  {date} | {decision['type']} | {decision['reason']} | Status: {decision['status']}")

    def process_pending(self):
        """Process pending decisions when reconnected to federation."""
        network_state = self.detect_network()

        if network_state != "home_federation":
            print(f"{RED}✗{NC} Cannot process pending decisions - not connected to federation")
            print(f"Current network: {network_state}")
            return False

        with open(PENDING_FILE, 'r') as f:
            pending_data = json.load(f)

        pending = [d for d in pending_data["pending_decisions"] if d["status"] == "pending"]

        if not pending:
            print(f"{YELLOW}No pending decisions to process{NC}")
            return True

        print(f"{GREEN}Processing {len(pending)} pending decisions...{NC}")

        # Process each decision
        for decision in pending:
            print(f"Processing decision {decision['id']} (type: {decision['type']})...")

            # Here you would submit to the actual federation
            # For now, mark as processed
            decision["status"] = "processed"
            decision["processed_at"] = datetime.now().isoformat()

            # Update processed log
            with open(PROCESSED_FILE, 'r') as f:
                processed_data = json.load(f)

            processed_data["processed"].append(decision)

            with open(PROCESSED_FILE, 'w') as f:
                json.dump(processed_data, f, indent=2)

            print(f"{GREEN}✓{NC} Decision {decision['id']} processed")

        # Save updated pending file
        with open(PENDING_FILE, 'w') as f:
            json.dump(pending_data, f, indent=2)

        return True

    def export_for_git(self):
        """Export pending decisions for git commit."""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        export_path = Path(f"/mnt/c/projects/ai-agents/ACT/implementation/society4/public")
        export_path.mkdir(parents=True, exist_ok=True)
        export_file = export_path / f"pending_consensus_{timestamp}.json"

        with open(PENDING_FILE, 'r') as f:
            pending_data = json.load(f)

        # Add export metadata
        pending_data["export_metadata"] = {
            "timestamp": datetime.now().isoformat(),
            "network_state": self.detect_network(),
            "hardware_attestation": self.get_hardware_hash(),
            "pending_count": len([d for d in pending_data["pending_decisions"] if d["status"] == "pending"]),
            "processed_count": len([d for d in pending_data["pending_decisions"] if d["status"] == "processed"])
        }

        with open(export_file, 'w') as f:
            json.dump(pending_data, f, indent=2)

        print(f"{GREEN}✓{NC} Exported pending consensus to:")
        print(f"  {export_file}")
        print("\nYou can now commit this to git for federation visibility")

    def show_status(self):
        """Show current system status."""
        print(f"\n{GREEN}=== Society 4 Pending Consensus Status ==={NC}")

        # Network state
        print(f"\n{YELLOW}Network State:{NC}")
        network = self.detect_network()

        # Hardware binding
        print(f"\n{YELLOW}Hardware Binding:{NC}")
        hw_hash = self.get_hardware_hash()
        print(f"  Hash: {hw_hash}")

        # Pending decisions
        print(f"\n{YELLOW}Pending Decisions:{NC}")
        with open(PENDING_FILE, 'r') as f:
            pending_data = json.load(f)

        pending = len([d for d in pending_data["pending_decisions"] if d["status"] == "pending"])
        processed = len([d for d in pending_data["pending_decisions"] if d["status"] == "processed"])
        print(f"  Pending: {pending}")
        print(f"  Processed: {processed}")

        # Last network transition
        if NETWORK_STATE_FILE.exists():
            print(f"\n{YELLOW}Last Network State:{NC}")
            with open(NETWORK_STATE_FILE, 'r') as f:
                print(f"  {f.read()}")

def main():
    """Main command handler."""
    pc = PendingConsensus()

    if len(sys.argv) < 2:
        command = "help"
    else:
        command = sys.argv[1]

    if command == "init":
        pc.init_system()
    elif command == "add":
        if len(sys.argv) < 5:
            print("Usage: pending_consensus.py add <type> <data> <reason>")
            sys.exit(1)
        pc.add_pending_decision(sys.argv[2], sys.argv[3], sys.argv[4])
    elif command == "list":
        pc.list_pending()
    elif command == "process":
        pc.process_pending()
    elif command == "export":
        pc.export_for_git()
    elif command == "status":
        pc.show_status()
    elif command == "network":
        pc.detect_network()
    else:
        print("Society 4 Pending Consensus System")
        print("")
        print("Usage: python3 pending_consensus.py [command]")
        print("")
        print("Commands:")
        print("  init     - Initialize pending consensus system")
        print("  add      - Add a pending decision: add <type> <data> <reason>")
        print("  list     - List all pending decisions")
        print("  process  - Process pending decisions (requires federation network)")
        print("  export   - Export pending decisions for git commit")
        print("  status   - Show current system status")
        print("  network  - Detect current network state")
        print("")
        print("Examples:")
        print('  python3 pending_consensus.py add vote \'{"proposal":"123","vote":"yes"}\' "Support Web4"')
        print('  python3 pending_consensus.py add transaction \'{"to":"society1","amount":10}\' "ATP delegation"')

if __name__ == "__main__":
    main()