#!/usr/bin/env python3
"""
CBP Token Pool Implementation
Adapted from Society4's Web4-compliant ATP/ADP system
Achieves energy economics for cache and metrics operations
"""

import json
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, field, asdict
from pathlib import Path

# CBP Role Structure with ATP Allocations (Total: 1000)
CBP_ROLES = [
    {"name": "Data Queen", "initial_atp": 140, "daily_recharge": 20},
    {"name": "Metrics Queen", "initial_atp": 130, "daily_recharge": 20},
    {"name": "Security Queen", "initial_atp": 130, "daily_recharge": 20},
    {"name": "Bridge Queen", "initial_atp": 120, "daily_recharge": 20},
    {"name": "Cache Queen", "initial_atp": 110, "daily_recharge": 20},
    {"name": "Worker 1", "initial_atp": 90, "daily_recharge": 20},
    {"name": "Worker 2", "initial_atp": 90, "daily_recharge": 20},
    {"name": "Worker 3", "initial_atp": 90, "daily_recharge": 20},
    {"name": "Coordinator", "initial_atp": 100, "daily_recharge": 20},
]

@dataclass
class AtpTransaction:
    """Record of an ATP transaction"""
    transaction_id: str
    society_id: str
    from_role: str
    to_role: Optional[str] = None
    amount: int = 0
    tx_type: str = "discharge"  # discharge, transfer, recharge
    operation_id: str = ""
    reason: str = ""
    timestamp: str = field(default_factory=lambda: datetime.now().isoformat())
    resulting_atp: int = 0
    resulting_adp: int = 0
    block_height: int = 0

@dataclass
class RoleAllocation:
    """ATP allocation and balance for a role"""
    role_id: str
    role_name: str
    initial_atp: int
    current_atp: int
    current_adp: int
    daily_recharge: int
    last_recharge: str
    total_spent: int = 0
    total_recharged: int = 0
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now().isoformat())

class CBPTokenPool:
    """
    CBP's ATP/ADP token pool implementation
    Web4-compliant energy economics for cache and metrics operations
    """

    def __init__(self, society_id: str = "lct:web4:cbp:self"):
        self.society_id = society_id
        self.total_atp = 1000  # Fixed per Web4 spec
        self.allocated_atp = 0
        self.available_atp = 1000
        self.total_adp = 0
        self.role_allocations: Dict[str, RoleAllocation] = {}
        self.role_balances: Dict[str, int] = {}
        self.adp_balances: Dict[str, int] = {}
        self.transactions: List[AtpTransaction] = []
        self.last_recharge = datetime.now()
        self.created_at = datetime.now()
        self.updated_at = datetime.now()
        self.version = 1

        # Storage paths
        self.pool_file = Path("implementation/cbp-chain/cbp_token_pool.json")
        self.tx_file = Path("implementation/cbp-chain/cbp_atp_transactions.json")

    def initialize_roles(self) -> None:
        """Initialize CBP roles with ATP allocations"""
        for role_def in CBP_ROLES:
            role_id = f"lct:web4:cbp:role:{role_def['name'].lower().replace(' ', '_')}"

            # Create role allocation
            role = RoleAllocation(
                role_id=role_id,
                role_name=role_def['name'],
                initial_atp=role_def['initial_atp'],
                current_atp=role_def['initial_atp'],
                current_adp=0,
                daily_recharge=role_def['daily_recharge'],
                last_recharge=datetime.now().isoformat()
            )

            # Allocate ATP
            self.role_allocations[role_id] = role
            self.role_balances[role_id] = role_def['initial_atp']
            self.adp_balances[role_id] = 0
            self.allocated_atp += role_def['initial_atp']
            self.available_atp -= role_def['initial_atp']

        self.validate_pool_integrity()
        print(f"✅ Initialized {len(CBP_ROLES)} CBP roles with {self.allocated_atp} ATP")

    def discharge_atp(self, role_id: str, amount: int, operation_id: str, reason: str) -> AtpTransaction:
        """
        Discharge ATP to ADP for an operation
        ATP → ADP represents energy spent on computation
        """
        if role_id not in self.role_balances:
            raise ValueError(f"Role not found: {role_id}")

        balance = self.role_balances[role_id]
        if balance < amount:
            raise ValueError(f"Insufficient ATP for {role_id}: {balance} < {amount}")

        # Discharge ATP to ADP
        self.role_balances[role_id] -= amount
        self.adp_balances[role_id] += amount
        self.total_adp += amount

        # Update role allocation
        if role_id in self.role_allocations:
            self.role_allocations[role_id].current_atp = self.role_balances[role_id]
            self.role_allocations[role_id].current_adp = self.adp_balances[role_id]
            self.role_allocations[role_id].total_spent += amount
            self.role_allocations[role_id].updated_at = datetime.now().isoformat()

        # Create transaction
        tx = AtpTransaction(
            transaction_id=f"tx-{role_id.split(':')[-1]}-{int(time.time())}",
            society_id=self.society_id,
            from_role=role_id,
            amount=amount,
            tx_type="discharge",
            operation_id=operation_id,
            reason=reason,
            resulting_atp=self.role_balances[role_id],
            resulting_adp=self.adp_balances[role_id]
        )

        self.transactions.append(tx)
        self.updated_at = datetime.now()
        self.version += 1

        return tx

    def transfer_atp(self, from_role: str, to_role: str, amount: int, operation_id: str, reason: str) -> AtpTransaction:
        """Transfer ATP between roles"""
        if from_role not in self.role_balances:
            raise ValueError(f"Source role not found: {from_role}")
        if to_role not in self.role_balances:
            raise ValueError(f"Target role not found: {to_role}")

        from_balance = self.role_balances[from_role]
        if from_balance < amount:
            raise ValueError(f"Insufficient ATP for {from_role}: {from_balance} < {amount}")

        # Transfer ATP
        self.role_balances[from_role] -= amount
        self.role_balances[to_role] += amount

        # Update allocations
        if from_role in self.role_allocations:
            self.role_allocations[from_role].current_atp = self.role_balances[from_role]
            self.role_allocations[from_role].updated_at = datetime.now().isoformat()
        if to_role in self.role_allocations:
            self.role_allocations[to_role].current_atp = self.role_balances[to_role]
            self.role_allocations[to_role].updated_at = datetime.now().isoformat()

        # Create transaction
        tx = AtpTransaction(
            transaction_id=f"tx-transfer-{int(time.time())}",
            society_id=self.society_id,
            from_role=from_role,
            to_role=to_role,
            amount=amount,
            tx_type="transfer",
            operation_id=operation_id,
            reason=reason,
            resulting_atp=self.role_balances[from_role]
        )

        self.transactions.append(tx)
        self.updated_at = datetime.now()
        self.version += 1

        return tx

    def daily_recharge(self) -> Dict[str, int]:
        """
        Perform daily recharge of ATP from ADP
        ADP → ATP represents energy recovery
        """
        recharged = {}

        for role_id, allocation in self.role_allocations.items():
            # Calculate recharge amount (min of daily_recharge and current ADP)
            recharge_amount = min(allocation.daily_recharge, self.adp_balances[role_id])

            if recharge_amount > 0:
                # Recharge ADP to ATP
                self.adp_balances[role_id] -= recharge_amount
                self.role_balances[role_id] += recharge_amount
                self.total_adp -= recharge_amount

                # Update allocation
                allocation.current_atp = self.role_balances[role_id]
                allocation.current_adp = self.adp_balances[role_id]
                allocation.total_recharged += recharge_amount
                allocation.last_recharge = datetime.now().isoformat()
                allocation.updated_at = datetime.now().isoformat()

                recharged[role_id] = recharge_amount

                # Create transaction
                tx = AtpTransaction(
                    transaction_id=f"tx-recharge-{role_id.split(':')[-1]}-{int(time.time())}",
                    society_id=self.society_id,
                    from_role=role_id,
                    amount=recharge_amount,
                    tx_type="recharge",
                    operation_id="daily_recharge",
                    reason="Scheduled daily ATP recharge",
                    resulting_atp=self.role_balances[role_id],
                    resulting_adp=self.adp_balances[role_id]
                )

                self.transactions.append(tx)

        self.last_recharge = datetime.now()
        self.updated_at = datetime.now()
        self.version += 1

        return recharged

    def validate_pool_integrity(self) -> bool:
        """Validate pool integrity constraints"""
        # Check total ATP constraint
        total_in_roles = sum(self.role_balances.values())
        total_adp_in_roles = sum(self.adp_balances.values())

        if total_in_roles + self.available_atp + total_adp_in_roles != self.total_atp:
            raise ValueError(f"Pool integrity violation: ATP conservation failed")

        # Check allocation limits
        if self.allocated_atp > self.total_atp:
            raise ValueError(f"Allocation exceeds total: {self.allocated_atp} > {self.total_atp}")

        return True

    def get_role_status(self, role_id: str) -> Dict:
        """Get current status of a role"""
        if role_id not in self.role_allocations:
            raise ValueError(f"Role not found: {role_id}")

        allocation = self.role_allocations[role_id]
        return {
            "role_id": role_id,
            "role_name": allocation.role_name,
            "current_atp": self.role_balances[role_id],
            "current_adp": self.adp_balances[role_id],
            "initial_atp": allocation.initial_atp,
            "total_spent": allocation.total_spent,
            "total_recharged": allocation.total_recharged,
            "energy_efficiency": (allocation.total_recharged / max(allocation.total_spent, 1)) * 100
        }

    def get_pool_summary(self) -> Dict:
        """Get summary of pool status"""
        return {
            "society_id": self.society_id,
            "total_atp": self.total_atp,
            "allocated_atp": self.allocated_atp,
            "available_atp": self.available_atp,
            "total_adp": self.total_adp,
            "active_roles": len(self.role_allocations),
            "total_transactions": len(self.transactions),
            "last_recharge": self.last_recharge.isoformat(),
            "version": self.version,
            "energy_cycling": (self.total_adp / self.total_atp) * 100
        }

    def save_state(self) -> None:
        """Save pool state to disk"""
        pool_data = {
            "society_id": self.society_id,
            "total_atp": self.total_atp,
            "allocated_atp": self.allocated_atp,
            "available_atp": self.available_atp,
            "total_adp": self.total_adp,
            "role_allocations": {k: asdict(v) for k, v in self.role_allocations.items()},
            "role_balances": self.role_balances,
            "adp_balances": self.adp_balances,
            "last_recharge": self.last_recharge.isoformat(),
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
            "version": self.version
        }

        self.pool_file.parent.mkdir(parents=True, exist_ok=True)
        with open(self.pool_file, 'w') as f:
            json.dump(pool_data, f, indent=2)

        # Save transactions
        if self.transactions:
            tx_data = [asdict(tx) for tx in self.transactions]
            with open(self.tx_file, 'w') as f:
                json.dump(tx_data, f, indent=2)

    def load_state(self) -> None:
        """Load pool state from disk"""
        if self.pool_file.exists():
            with open(self.pool_file, 'r') as f:
                data = json.load(f)

            self.society_id = data['society_id']
            self.total_atp = data['total_atp']
            self.allocated_atp = data['allocated_atp']
            self.available_atp = data['available_atp']
            self.total_adp = data['total_adp']
            self.role_balances = data['role_balances']
            self.adp_balances = data['adp_balances']
            self.last_recharge = datetime.fromisoformat(data['last_recharge'])
            self.created_at = datetime.fromisoformat(data['created_at'])
            self.updated_at = datetime.fromisoformat(data['updated_at'])
            self.version = data['version']

            # Load role allocations
            for role_id, role_data in data['role_allocations'].items():
                self.role_allocations[role_id] = RoleAllocation(**role_data)

        # Load transactions
        if self.tx_file.exists():
            with open(self.tx_file, 'r') as f:
                tx_data = json.load(f)
                self.transactions = [AtpTransaction(**tx) for tx in tx_data]


def main():
    """Test CBP token pool"""
    import sys

    pool = CBPTokenPool()

    if len(sys.argv) < 2:
        command = "status"
    else:
        command = sys.argv[1]

    if command == "init":
        pool.initialize_roles()
        pool.save_state()
        print(json.dumps(pool.get_pool_summary(), indent=2))

    elif command == "status":
        pool.load_state()
        print("\n🏦 CBP Token Pool Status")
        print("=" * 50)
        summary = pool.get_pool_summary()
        for key, value in summary.items():
            print(f"{key:20}: {value}")

        print("\n📊 Role Balances:")
        for role_id in pool.role_allocations:
            status = pool.get_role_status(role_id)
            print(f"  {status['role_name']:20}: ATP={status['current_atp']:3} ADP={status['current_adp']:3}")

    elif command == "discharge":
        if len(sys.argv) < 5:
            print("Usage: discharge <role_name> <amount> <reason>")
            sys.exit(1)

        pool.load_state()
        role_name = sys.argv[2]
        amount = int(sys.argv[3])
        reason = " ".join(sys.argv[4:])

        # Find role ID
        role_id = f"lct:web4:cbp:role:{role_name.lower().replace(' ', '_')}"

        tx = pool.discharge_atp(role_id, amount, f"op-{int(time.time())}", reason)
        pool.save_state()

        print(f"✅ Discharged {amount} ATP from {role_name}")
        print(f"   Transaction: {tx.transaction_id}")
        print(f"   Resulting ATP: {tx.resulting_atp}")
        print(f"   Resulting ADP: {tx.resulting_adp}")

    elif command == "recharge":
        pool.load_state()
        recharged = pool.daily_recharge()
        pool.save_state()

        print("⚡ Daily Recharge Complete")
        for role_id, amount in recharged.items():
            role = pool.role_allocations[role_id]
            print(f"  {role.role_name}: +{amount} ATP")

    else:
        print("CBP Token Pool Management")
        print("Commands:")
        print("  init     - Initialize role allocations")
        print("  status   - Show pool status")
        print("  discharge - Discharge ATP for operation")
        print("  recharge - Perform daily recharge")


if __name__ == "__main__":
    main()