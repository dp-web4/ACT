#!/usr/bin/env python3
"""
CBP LCT (Linked Context Token) Implementation
Based on Web4 specification and Society4's reference implementation
Provides identity foundation for CBP society
"""

import json
import hashlib
import time
from datetime import datetime
from typing import Dict, List, Optional, Any
from dataclasses import dataclass, field, asdict
from pathlib import Path
import subprocess

# ============================================================================
# Web4 LCT Data Structures
# ============================================================================

@dataclass
class Web4Binding:
    """Establishes permanent link between LCT and entity"""
    entity_type: str  # device, ai, human, organization
    public_key: str  # mb64:coseKey format
    hardware_anchor: Optional[str] = None  # EAT token (RFC 9334)
    created_at: str = field(default_factory=lambda: datetime.now().isoformat())
    binding_proof: str = "pending"  # COSE signature

@dataclass
class Web4BirthCert:
    """Foundational identity and context"""
    citizen_role: str  # lct:web4:role:citizen:...
    context: str  # federation, nation, platform
    birth_timestamp: str
    parent_entity: Optional[str] = None  # lct:web4:...
    birth_witnesses: List[str] = field(default_factory=list)  # Witness LCT IDs
    founding_purpose: Optional[str] = None  # Purpose statement

@dataclass
class Web4BoundRelation:
    """Binding relationships (permanent)"""
    lct_id: str
    relation_type: str  # parent, child, sibling
    timestamp: str
    binding_context: Optional[str] = None

@dataclass
class Web4PairedRelation:
    """Active pairings (can be temporary)"""
    lct_id: str
    pairing_type: str  # birth_certificate, role, operational
    permanent: bool
    timestamp: str
    context: Optional[str] = None  # For non-birth pairings
    session_id: Optional[str] = None  # For operational pairings

@dataclass
class Web4WitnessRelation:
    """Witness relationships"""
    lct_id: str
    role: str  # time, audit, oracle
    last_attestation: str

@dataclass
class Web4MRH:
    """Markov Relevancy Horizon tracking"""
    bound: List[Web4BoundRelation] = field(default_factory=list)
    paired: List[Web4PairedRelation] = field(default_factory=list)
    witnessing: List[Web4WitnessRelation] = field(default_factory=list)
    horizon_depth: int = 3
    last_updated: str = field(default_factory=lambda: datetime.now().isoformat())

@dataclass
class Web4Policy:
    """Capabilities and constraints"""
    capabilities: List[str] = field(default_factory=list)
    constraints: Dict[str, Any] = field(default_factory=dict)

@dataclass
class Web4Attestation:
    """Witnessing events"""
    witness: str  # DID of witness
    attestation_type: str  # time, audit, oracle, existence, action, state, quality
    signature: str  # COSE signature
    timestamp: str
    evidence: Optional[str] = None

@dataclass
class Web4Lineage:
    """LCT evolution tracking"""
    reason: str  # genesis, rotation, fork, upgrade
    timestamp: str
    parent: Optional[str] = None  # Previous LCT ID

@dataclass
class Web4Revocation:
    """LCT revocation status"""
    status: str = "active"  # active, revoked
    timestamp: Optional[str] = None
    reason: Optional[str] = None  # compromise, superseded, expired

@dataclass
class Web4LCT:
    """Web4-compliant Linked Context Token"""
    lct_id: str  # lct:web4:mb32...
    subject: str  # did:web4:key:z6Mk...
    binding: Web4Binding
    mrh: Web4MRH
    policy: Web4Policy
    birth_certificate: Optional[Web4BirthCert] = None
    attestations: List[Web4Attestation] = field(default_factory=list)
    lineage: List[Web4Lineage] = field(default_factory=list)
    revocation: Optional[Web4Revocation] = None

# ============================================================================
# CBP Society Implementation
# ============================================================================

class CBPSelfLCT:
    """
    CBP's genesis 'self' LCT implementation
    Establishes CBP as a Web4 society with focus on data, metrics, and bridging
    """

    def __init__(self):
        self.hardware_hash = self._get_hardware_hash()
        self.public_key = "mb64:pending"  # Will be generated with proper keypair
        self.lct = self._create_genesis_lct()
        self.storage_path = Path("implementation/cbp-chain/cbp_self_lct.json")

    def _get_hardware_hash(self) -> str:
        """Extract WSL2 hardware hash"""
        try:
            # Get various hardware identifiers
            components = []

            # CPU info
            with open('/proc/cpuinfo', 'r') as f:
                for line in f:
                    if line.startswith('model name'):
                        components.append(line.split(':')[1].strip())
                        break

            # Memory info
            with open('/proc/meminfo', 'r') as f:
                for line in f:
                    if line.startswith('MemTotal'):
                        components.append(line.split(':')[1].strip())
                        break

            # Hostname
            result = subprocess.run(['hostname'], capture_output=True, text=True)
            if result.returncode == 0:
                components.append(result.stdout.strip())

            # Create deterministic hash
            combined = '|'.join(components)
            return hashlib.sha256(combined.encode()).hexdigest()
        except Exception as e:
            # Fallback to timestamp-based hash for development
            return hashlib.sha256(f"cbp-dev-{time.time()}".encode()).hexdigest()

    def _create_genesis_lct(self) -> Web4LCT:
        """Create CBP's genesis LCT"""
        now = datetime.now().isoformat()

        # Hardware binding
        binding = Web4Binding(
            entity_type="device",  # CBP is hardware-bound
            public_key=self.public_key,
            hardware_anchor=f"eat:mb64:hw:{self.hardware_hash}",
            created_at=now,
            binding_proof="pending"  # Will be COSE signature
        )

        # MRH relationships
        mrh = Web4MRH(
            bound=[
                Web4BoundRelation(
                    lct_id=f"lct:web4:hardware:wsl2:{self.hardware_hash[:16]}",
                    relation_type="parent",
                    timestamp=now,
                    binding_context="wsl2_hardware_sovereignty"
                )
            ],
            paired=[],  # Will add birth certificate pairing
            witnessing=[],
            horizon_depth=3,
            last_updated=now
        )

        # CBP-specific capabilities and constraints
        policy = Web4Policy(
            capabilities=[
                "pairing:initiate",
                "cache:manage",
                "metrics:collect",
                "federation:bridge",
                "blockchain:witness",
                "data:persist",
                "consensus:participate"
            ],
            constraints={
                "atp_allocation": 1000,
                "hardware_hash": self.hardware_hash,
                "network_mobility": True,
                "requires_quorum": 3,
                "cache_size_mb": 1024,
                "metrics_retention_days": 30,
                "bridge_protocols": ["web4", "mcp", "http"]
            }
        )

        # Genesis lineage
        lineage = [
            Web4Lineage(
                reason="genesis",
                timestamp=now,
                parent=None
            )
        ]

        # Active revocation status
        revocation = Web4Revocation(
            status="active",
            timestamp=now
        )

        # Create LCT
        lct = Web4LCT(
            lct_id="lct:web4:cbp:self:pending",  # Will be computed from binding proof
            subject="did:web4:cbp:coordinator",
            binding=binding,
            mrh=mrh,
            policy=policy,
            birth_certificate=None,  # Pending federation witnesses
            attestations=[],
            lineage=lineage,
            revocation=revocation
        )

        return lct

    def add_birth_certificate(self, parent_entity: str, witnesses: List[str], purpose: str) -> None:
        """Add federation birth certificate once witnessed"""
        now = datetime.now().isoformat()

        self.lct.birth_certificate = Web4BirthCert(
            citizen_role="lct:web4:federation:act:citizen:cbp",
            context="federation",
            birth_timestamp=now,
            parent_entity=parent_entity,  # lct:web4:federation:act
            birth_witnesses=witnesses,  # Genesis, Society4, Society2, Sprout
            founding_purpose=purpose or "Data persistence, metrics collection, and federation bridging"
        )

        # Add birth certificate pairing
        self.lct.mrh.paired.append(
            Web4PairedRelation(
                lct_id="lct:web4:federation:act:citizen:cbp",
                pairing_type="birth_certificate",
                permanent=True,
                timestamp=now,
                context="federation_citizenship"
            )
        )

        self.lct.mrh.last_updated = now

    def add_witness_attestation(self, witness_did: str, attestation_type: str, signature: str, evidence: Optional[str] = None) -> None:
        """Add a witness attestation to the LCT"""
        attestation = Web4Attestation(
            witness=witness_did,
            attestation_type=attestation_type,
            signature=signature,
            timestamp=datetime.now().isoformat(),
            evidence=evidence
        )
        self.lct.attestations.append(attestation)

    def compute_lct_id(self) -> str:
        """Compute the actual LCT ID from binding proof"""
        # In production, this would be SHA256 of the COSE binding proof
        # For now, use hardware hash
        lct_hash = hashlib.sha256(self.hardware_hash.encode()).hexdigest()[:32]
        self.lct.lct_id = f"lct:web4:mb32:{lct_hash}"
        return self.lct.lct_id

    def save(self) -> None:
        """Save LCT to disk"""
        self.storage_path.parent.mkdir(parents=True, exist_ok=True)

        # Convert to dict for JSON serialization
        lct_dict = self._lct_to_dict(self.lct)

        with open(self.storage_path, 'w') as f:
            json.dump({
                "lct": lct_dict,
                "hardware_hash": self.hardware_hash,
                "network_state": "federation_active",
                "created_at": datetime.now().isoformat(),
                "version": "1.0.0"
            }, f, indent=2)

    def load(self) -> bool:
        """Load LCT from disk"""
        if not self.storage_path.exists():
            return False

        with open(self.storage_path, 'r') as f:
            data = json.load(f)

        self.hardware_hash = data['hardware_hash']
        # Would reconstruct LCT from dict here
        return True

    def _lct_to_dict(self, lct: Web4LCT) -> Dict:
        """Convert LCT to dictionary for serialization"""
        result = asdict(lct)

        # Remove None values for cleaner JSON
        def remove_none(d):
            if isinstance(d, dict):
                return {k: remove_none(v) for k, v in d.items() if v is not None}
            elif isinstance(d, list):
                return [remove_none(v) for v in d]
            return d

        return remove_none(result)

    def display_status(self) -> None:
        """Display current LCT status"""
        print("\n" + "="*60)
        print("🔐 CBP Self LCT Status")
        print("="*60)
        print(f"LCT ID: {self.lct.lct_id}")
        print(f"Subject: {self.lct.subject}")
        print(f"Hardware Hash: {self.hardware_hash[:16]}...")
        print(f"Entity Type: {self.lct.binding.entity_type}")

        print("\n📋 Capabilities:")
        for cap in self.lct.policy.capabilities[:5]:
            print(f"  - {cap}")

        print("\n🔗 MRH Relations:")
        print(f"  Bound: {len(self.lct.mrh.bound)} entities")
        print(f"  Paired: {len(self.lct.mrh.paired)} entities")
        print(f"  Witnessing: {len(self.lct.mrh.witnessing)} entities")

        if self.lct.birth_certificate:
            print("\n📜 Birth Certificate:")
            print(f"  Citizen Role: {self.lct.birth_certificate.citizen_role}")
            print(f"  Context: {self.lct.birth_certificate.context}")
            print(f"  Witnesses: {len(self.lct.birth_certificate.birth_witnesses)}")
        else:
            print("\n⚠️ Birth Certificate: PENDING FEDERATION WITNESSES")

        print("\n✅ Revocation Status:", self.lct.revocation.status if self.lct.revocation else "Unknown")
        print("="*60)

# ============================================================================
# CBP Role LCTs
# ============================================================================

class CBPRoleLCT:
    """LCT for CBP roles (Queens and Workers)"""

    @staticmethod
    def create_role_lct(role_name: str, role_type: str, domain: str, atp_budget: int, parent_lct: str) -> Web4LCT:
        """Create an LCT for a CBP role"""
        now = datetime.now().isoformat()
        role_id = role_name.lower().replace(' ', '_')

        binding = Web4Binding(
            entity_type="role",
            public_key="mb64:pending",  # Derived from parent
            created_at=now,
            binding_proof="pending"
        )

        mrh = Web4MRH(
            bound=[
                Web4BoundRelation(
                    lct_id=parent_lct,
                    relation_type="child",
                    timestamp=now,
                    binding_context=f"cbp_{role_type}_role"
                )
            ],
            paired=[
                Web4PairedRelation(
                    lct_id=parent_lct,
                    pairing_type="role",
                    permanent=True,
                    timestamp=now,
                    context=f"{domain}_management"
                )
            ]
        )

        # Role-specific capabilities
        capabilities = []
        if role_type == "queen":
            capabilities = [
                f"{domain}:manage",
                "worker:coordinate",
                "atp:allocate",
                "decision:autonomous",
                "federation:interact"
            ]
        else:  # worker
            capabilities = [
                f"{domain}:execute",
                "task:process",
                "report:generate"
            ]

        policy = Web4Policy(
            capabilities=capabilities,
            constraints={
                "atp_budget": atp_budget,
                "role_type": role_type,
                "domain": domain,
                "delegation_enabled": role_type == "queen"
            }
        )

        lct = Web4LCT(
            lct_id=f"lct:web4:cbp:role:{role_id}",
            subject=f"did:web4:cbp:role:{role_id}",
            binding=binding,
            mrh=mrh,
            policy=policy,
            lineage=[
                Web4Lineage(
                    reason="role_creation",
                    timestamp=now,
                    parent=parent_lct
                )
            ],
            revocation=Web4Revocation(status="active")
        )

        return lct

# ============================================================================
# Main Interface
# ============================================================================

def main():
    """Test CBP LCT implementation"""
    import sys

    cbp = CBPSelfLCT()

    if len(sys.argv) < 2:
        command = "status"
    else:
        command = sys.argv[1]

    if command == "init":
        # Initialize CBP self LCT
        lct_id = cbp.compute_lct_id()
        cbp.save()
        print(f"✅ Created CBP Self LCT: {lct_id}")
        cbp.display_status()

    elif command == "status":
        if cbp.load():
            cbp.display_status()
        else:
            print("❌ No CBP LCT found. Run 'init' first.")

    elif command == "birth":
        # Add birth certificate (requires witnesses)
        if len(sys.argv) < 3:
            print("Usage: birth <witness1,witness2,witness3>")
            sys.exit(1)

        cbp.load()
        witnesses = sys.argv[2].split(',')
        cbp.add_birth_certificate(
            parent_entity="lct:web4:federation:act",
            witnesses=witnesses,
            purpose="Data persistence, metrics collection, and federation bridging"
        )
        cbp.save()
        print("✅ Birth certificate added")
        cbp.display_status()

    elif command == "roles":
        # Create role LCTs
        cbp.load()
        parent_lct = cbp.lct.lct_id

        roles = [
            ("Data Queen", "queen", "data", 140),
            ("Metrics Queen", "queen", "metrics", 130),
            ("Security Queen", "queen", "security", 130),
            ("Bridge Queen", "queen", "bridge", 120),
            ("Cache Queen", "queen", "cache", 110)
        ]

        for name, role_type, domain, atp in roles:
            role_lct = CBPRoleLCT.create_role_lct(name, role_type, domain, atp, parent_lct)
            print(f"✅ Created LCT for {name}: {role_lct.lct_id}")

    else:
        print("CBP LCT Management")
        print("Commands:")
        print("  init   - Initialize CBP self LCT")
        print("  status - Display LCT status")
        print("  birth  - Add birth certificate")
        print("  roles  - Create role LCTs")


if __name__ == "__main__":
    main()