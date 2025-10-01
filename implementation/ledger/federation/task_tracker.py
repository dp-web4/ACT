#!/usr/bin/env python3
"""
Federation Task Tracking System for SAGE Development
Tracks ATP allocations, task progress, and cross-society coordination
"""

import json
import os
from datetime import datetime
from typing import Dict, List, Optional, Tuple
from pathlib import Path
import hashlib

class FederationTaskTracker:
    """Tracks federation tasks across all societies with ATP energy accounting"""
    
    def __init__(self, base_path: str = "/home/dp/ai-workspace/act/implementation/ledger"):
        self.base_path = Path(base_path)
        self.federation_dir = self.base_path / "federation"
        self.inbox_dir = self.base_path / "federation_inbox"
        self.outbox_dir = self.base_path / "federation_outbox"
        self.state_file = self.federation_dir / "tracker_state.json"
        
        # Initialize directories
        self.federation_dir.mkdir(exist_ok=True)
        self.inbox_dir.mkdir(exist_ok=True)
        self.outbox_dir.mkdir(exist_ok=True)
        
        # Load or initialize state
        self.state = self._load_state()
        
    def _load_state(self) -> Dict:
        """Load tracker state from disk"""
        if self.state_file.exists():
            with open(self.state_file, 'r') as f:
                return json.load(f)
        else:
            # Initialize new state
            return {
                "cycle": 1,
                "block_start": 70277,
                "block_current": 70338,
                "total_atp_budget": 20000,
                "societies": {
                    "Society4": {
                        "atp_allocated": 5000,
                        "atp_discharged": 0,
                        "atp_recharged": 0,
                        "tasks": [],
                        "deliverables": []
                    },
                    "Society2": {
                        "atp_allocated": 5000,
                        "atp_discharged": 0,
                        "atp_recharged": 0,
                        "tasks": [],
                        "deliverables": []
                    },
                    "Sprout": {
                        "atp_allocated": 5000,
                        "atp_discharged": 0,
                        "atp_recharged": 0,
                        "tasks": [],
                        "deliverables": []
                    },
                    "Genesis": {
                        "atp_allocated": 5000,
                        "atp_discharged": 0,
                        "atp_recharged": 0,
                        "tasks": [],
                        "deliverables": []
                    }
                },
                "events": [],
                "witness_pool": []
            }
    
    def _save_state(self):
        """Persist tracker state to disk"""
        with open(self.state_file, 'w') as f:
            json.dump(self.state, f, indent=2)
    
    def register_task(self, society: str, task_id: str, description: str, 
                     atp_cost: int = 100) -> Dict:
        """Register a new task for a society"""
        if society not in self.state["societies"]:
            return {"error": f"Unknown society: {society}"}
        
        task = {
            "id": task_id,
            "description": description,
            "status": "accepted",
            "atp_cost": atp_cost,
            "created_block": self.state["block_current"],
            "created_time": datetime.now().isoformat(),
            "witnesses": []
        }
        
        # Add task and discharge ATP
        self.state["societies"][society]["tasks"].append(task)
        self.state["societies"][society]["atp_discharged"] += atp_cost
        
        # Log event
        event = {
            "type": "task_registered",
            "society": society,
            "task_id": task_id,
            "atp": atp_cost,
            "block": self.state["block_current"],
            "timestamp": datetime.now().isoformat()
        }
        self.state["events"].append(event)
        
        self._save_state()
        return {"status": "success", "task": task}
    
    def update_task_status(self, society: str, task_id: str, 
                          status: str, atp_discharge: int = 0) -> Dict:
        """Update task status and discharge ATP"""
        if society not in self.state["societies"]:
            return {"error": f"Unknown society: {society}"}
        
        # Find task
        tasks = self.state["societies"][society]["tasks"]
        task = None
        for t in tasks:
            if t["id"] == task_id:
                task = t
                break
        
        if not task:
            return {"error": f"Task {task_id} not found for {society}"}
        
        # Update status
        task["status"] = status
        task["updated_block"] = self.state["block_current"]
        task["updated_time"] = datetime.now().isoformat()
        
        # Discharge ATP
        if atp_discharge > 0:
            self.state["societies"][society]["atp_discharged"] += atp_discharge
        
        # Log event
        event = {
            "type": f"task_{status}",
            "society": society,
            "task_id": task_id,
            "atp": atp_discharge,
            "block": self.state["block_current"],
            "timestamp": datetime.now().isoformat()
        }
        self.state["events"].append(event)
        
        self._save_state()
        return {"status": "success", "task": task}
    
    def add_deliverable(self, society: str, deliverable_id: str, 
                       description: str, file_path: str, 
                       atp_recharge: int = 0) -> Dict:
        """Register a completed deliverable and recharge ATP"""
        if society not in self.state["societies"]:
            return {"error": f"Unknown society: {society}"}
        
        deliverable = {
            "id": deliverable_id,
            "description": description,
            "file_path": file_path,
            "atp_recharged": atp_recharge,
            "block": self.state["block_current"],
            "timestamp": datetime.now().isoformat(),
            "hash": self._hash_file(file_path) if os.path.exists(file_path) else None
        }
        
        # Add deliverable and recharge ATP
        self.state["societies"][society]["deliverables"].append(deliverable)
        if atp_recharge > 0:
            self.state["societies"][society]["atp_recharged"] += atp_recharge
        
        # Log event
        event = {
            "type": "deliverable_completed",
            "society": society,
            "deliverable_id": deliverable_id,
            "atp": atp_recharge,
            "block": self.state["block_current"],
            "timestamp": datetime.now().isoformat()
        }
        self.state["events"].append(event)
        
        self._save_state()
        return {"status": "success", "deliverable": deliverable}
    
    def add_witness(self, witness_id: str, event_hash: str, signature: str) -> Dict:
        """Add witness attestation for an event"""
        witness = {
            "witness_id": witness_id,
            "event_hash": event_hash,
            "signature": signature,
            "block": self.state["block_current"],
            "timestamp": datetime.now().isoformat()
        }
        
        self.state["witness_pool"].append(witness)
        self._save_state()
        
        return {"status": "success", "witness": witness}
    
    def get_society_status(self, society: str) -> Dict:
        """Get detailed status for a society"""
        if society not in self.state["societies"]:
            return {"error": f"Unknown society: {society}"}
        
        soc_data = self.state["societies"][society]
        
        # Calculate ATP balance
        atp_balance = (soc_data["atp_allocated"] - 
                      soc_data["atp_discharged"] + 
                      soc_data["atp_recharged"])
        
        # Count task statuses
        task_stats = {
            "accepted": 0,
            "in_progress": 0,
            "completed": 0,
            "blocked": 0
        }
        for task in soc_data["tasks"]:
            status = task.get("status", "unknown")
            if status in task_stats:
                task_stats[status] += 1
        
        return {
            "society": society,
            "atp_allocated": soc_data["atp_allocated"],
            "atp_discharged": soc_data["atp_discharged"],
            "atp_recharged": soc_data["atp_recharged"],
            "atp_balance": atp_balance,
            "task_stats": task_stats,
            "deliverables_count": len(soc_data["deliverables"]),
            "tasks": soc_data["tasks"],
            "deliverables": soc_data["deliverables"]
        }
    
    def get_federation_summary(self) -> Dict:
        """Get summary of entire federation status"""
        summary = {
            "cycle": self.state["cycle"],
            "block_start": self.state["block_start"],
            "block_current": self.state["block_current"],
            "blocks_elapsed": self.state["block_current"] - self.state["block_start"],
            "total_atp_budget": self.state["total_atp_budget"],
            "societies": {}
        }
        
        total_discharged = 0
        total_recharged = 0
        total_tasks = 0
        total_deliverables = 0
        
        for society in self.state["societies"]:
            status = self.get_society_status(society)
            summary["societies"][society] = {
                "atp_balance": status["atp_balance"],
                "tasks": status["task_stats"],
                "deliverables": status["deliverables_count"]
            }
            total_discharged += status["atp_discharged"]
            total_recharged += status["atp_recharged"]
            total_tasks += len(status["tasks"])
            total_deliverables += status["deliverables_count"]
        
        summary["totals"] = {
            "atp_discharged": total_discharged,
            "atp_recharged": total_recharged,
            "atp_active": total_discharged - total_recharged,
            "tasks": total_tasks,
            "deliverables": total_deliverables,
            "witnesses": len(self.state["witness_pool"])
        }
        
        return summary
    
    def generate_progress_report(self) -> str:
        """Generate markdown progress report"""
        summary = self.get_federation_summary()
        
        report = f"""# Federation Progress Report - Cycle {summary['cycle']}

**Generated**: {datetime.now().isoformat()}
**Current Block**: {summary['block_current']}
**Blocks Elapsed**: {summary['blocks_elapsed']} / 1500 (target)

## ATP Energy Status
- **Total Budget**: {summary['total_atp_budget']:,} ATP
- **Total Discharged**: {summary['totals']['atp_discharged']:,} ATP
- **Total Recharged**: {summary['totals']['atp_recharged']:,} ATP
- **Active Energy**: {summary['totals']['atp_active']:,} ATP

## Society Status

"""
        
        for society, data in summary['societies'].items():
            report += f"""### {society}
- **ATP Balance**: {data['atp_balance']:,} ATP
- **Tasks**: {data['tasks']['accepted']} accepted, {data['tasks']['in_progress']} in progress, {data['tasks']['completed']} completed
- **Deliverables**: {data['deliverables']}

"""
        
        # Recent events
        report += "## Recent Events\n\n"
        recent_events = self.state["events"][-10:]  # Last 10 events
        for event in reversed(recent_events):
            report += f"- **{event['type']}** ({event['society']}): "
            report += f"{event.get('task_id', event.get('deliverable_id', ''))} "
            report += f"[{event['atp']} ATP] @ Block {event['block']}\n"
        
        # Witness attestations
        report += f"\n## Witness Pool\n"
        report += f"- **Total Attestations**: {len(self.state['witness_pool'])}\n"
        
        return report
    
    def check_inbox(self) -> List[Dict]:
        """Check federation inbox for updates from societies"""
        updates = []
        
        for file_path in self.inbox_dir.glob("*.json"):
            with open(file_path, 'r') as f:
                update = json.load(f)
                updates.append(update)
                
                # Process update based on type
                if update.get("type") == "task_update":
                    self.update_task_status(
                        update["society"],
                        update["task_id"],
                        update["status"],
                        update.get("atp_discharge", 0)
                    )
                elif update.get("type") == "deliverable":
                    self.add_deliverable(
                        update["society"],
                        update["deliverable_id"],
                        update["description"],
                        update["file_path"],
                        update.get("atp_recharge", 0)
                    )
            
            # Archive processed file
            archive_dir = self.inbox_dir / "processed"
            archive_dir.mkdir(exist_ok=True)
            file_path.rename(archive_dir / file_path.name)
        
        return updates
    
    def _hash_file(self, file_path: str) -> str:
        """Generate SHA256 hash of file"""
        if not os.path.exists(file_path):
            return None
        
        sha256_hash = hashlib.sha256()
        with open(file_path, "rb") as f:
            for byte_block in iter(lambda: f.read(4096), b""):
                sha256_hash.update(byte_block)
        return sha256_hash.hexdigest()
    
    def update_block(self, new_block: int):
        """Update current block number"""
        self.state["block_current"] = new_block
        self._save_state()


if __name__ == "__main__":
    # Initialize tracker
    tracker = FederationTaskTracker()
    
    # Register initial tasks for Cycle 1
    print("=== Initializing Federation Task Tracker for SAGE Cycle 1 ===\n")
    
    # Society4 tasks
    tracker.register_task("Society4", "S4-001", "Fix SAGE training loop reward structure", 100)
    tracker.register_task("Society4", "S4-002", "Implement proper context encoding system", 100)
    tracker.register_task("Society4", "S4-003", "Create validation suite for reasoning", 100)
    
    # Society2 tasks
    tracker.register_task("Society2", "S2-001", "Wire external LLM as cognitive sensor", 100)
    tracker.register_task("Society2", "S2-002", "Implement trust-weighted output system", 100)
    tracker.register_task("Society2", "S2-003", "Test with 2B and 7B models", 100)
    
    # Sprout tasks
    tracker.register_task("Sprout", "SP-001", "Optimize SAGE for Jetson Orin Nano", 100)
    tracker.register_task("Sprout", "SP-002", "Implement memory-efficient inference", 100)
    tracker.register_task("Sprout", "SP-003", "Create performance monitoring dashboard", 100)
    
    # Genesis tasks
    tracker.register_task("Genesis", "GN-001", "Federation task tracking system", 100)
    tracker.register_task("Genesis", "GN-002", "Cross-society integration tests", 100)
    tracker.register_task("Genesis", "GN-003", "Documentation and progress reports", 100)
    
    # Mark Genesis tracking system as in progress
    tracker.update_task_status("Genesis", "GN-001", "in_progress", 200)
    
    # Generate initial report
    report = tracker.generate_progress_report()
    print(report)
    
    # Save report to file
    report_path = tracker.federation_dir / "progress_reports"
    report_path.mkdir(exist_ok=True)
    
    with open(report_path / "cycle_1_initial.md", 'w') as f:
        f.write(report)
    
    print(f"\n=== Tracker initialized and saved to {tracker.state_file} ===")
    print(f"=== Initial report saved to {report_path / 'cycle_1_initial.md'} ===")