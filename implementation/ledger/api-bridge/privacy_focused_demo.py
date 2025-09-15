#!/usr/bin/env python3
"""
Privacy-Focused Component Registry Demo

This script demonstrates the new privacy-focused features of the Web4 Component Registry,
including anonymous component registration, hash-based pairing verification,
anonymous pairing authorization, revocation events, and metadata retrieval.

The demo showcases how the system protects trade secrets while maintaining
verification capabilities through cryptographic hashes.
"""

import requests
import json
import time
import hashlib
import uuid
from typing import Dict, Any, Optional

class PrivacyFocusedDemo:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
        
        # Demo data
        self.demo_components = {
            "battery_module_001": {
                "real_id": "BM-2024-001",
                "manufacturer": "Tesla Motors",
                "type": "LITHIUM_ION_BATTERY",
                "specs": {
                    "capacity": "100kWh",
                    "voltage": "400V",
                    "chemistry": "NCA",
                    "cycle_life": "3000"
                }
            },
            "motor_controller_001": {
                "real_id": "MC-2024-001", 
                "manufacturer": "Bosch Automotive",
                "type": "MOTOR_CONTROLLER",
                "specs": {
                    "power": "300kW",
                    "efficiency": "95%",
                    "cooling": "liquid",
                    "protocol": "CAN_BUS"
                }
            },
            "sensor_array_001": {
                "real_id": "SA-2024-001",
                "manufacturer": "Continental AG",
                "type": "SENSOR_ARRAY",
                "specs": {
                    "sensors": ["temperature", "voltage", "current", "pressure"],
                    "accuracy": "0.1%",
                    "update_rate": "100Hz"
                }
            }
        }
        
        self.demo_users = {
            "manufacturer_1": "alice",
            "manufacturer_2": "bob", 
            "verifier": "charlie",
            "regulator": "david"
        }

    def health_check(self) -> bool:
        """Check if the API Bridge is running"""
        try:
            response = self.session.get(f"{self.base_url}/health")
            if response.status_code == 200:
                print("✅ API Bridge is running")
                return True
            else:
                print(f"❌ API Bridge health check failed: {response.status_code}")
                return False
        except Exception as e:
            print(f"❌ Cannot connect to API Bridge: {e}")
            return False

    def blockchain_status(self) -> bool:
        """Check blockchain connection status"""
        try:
            response = self.session.get(f"{self.base_url}/blockchain/status")
            if response.status_code == 200:
                status = response.json()
                print(f"✅ Blockchain status: {status.get('blockchain_status', 'unknown')}")
                return True
            else:
                print(f"❌ Blockchain status check failed: {response.status_code}")
                return False
        except Exception as e:
            print(f"❌ Cannot check blockchain status: {e}")
            return False

    def register_anonymous_component(self, component_key: str, user: str) -> Optional[str]:
        """Register a component anonymously using hashes"""
        component = self.demo_components[component_key]
        
        payload = {
            "creator": user,
            "real_component_id": component["real_id"],
            "manufacturer_id": component["manufacturer"],
            "component_type": component["type"],
            "context": "privacy_demo_registration"
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/v1/components/register-anonymous",
                json=payload
            )
            
            if response.status_code == 200:
                result = response.json()
                component_hash = result.get("component_hash")
                print(f"✅ Anonymous component registered: {component_hash}")
                print(f"   Real ID: {component['real_id']} (hidden)")
                print(f"   Manufacturer: {component['manufacturer']} (hashed)")
                print(f"   Type: {component['type']} (hashed)")
                print(f"   TX Hash: {result.get('txhash')}")
                return component_hash
            else:
                print(f"❌ Anonymous component registration failed: {response.status_code}")
                print(f"   Response: {response.text}")
                return None
                
        except Exception as e:
            print(f"❌ Error registering anonymous component: {e}")
            return None

    def verify_component_pairing_with_hashes(self, component_hash_a: str, component_hash_b: str, verifier: str) -> bool:
        """Verify if two components can pair using their hashes"""
        payload = {
            "verifier": verifier,
            "component_hash_a": component_hash_a,
            "component_hash_b": component_hash_b,
            "context": "privacy_demo_pairing_verification"
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/v1/components/verify-pairing-hashes",
                json=payload
            )
            
            if response.status_code == 200:
                result = response.json()
                can_pair = result.get("can_pair", False)
                reason = result.get("reason", "unknown")
                trust_score = result.get("trust_score", "0.0")
                
                print(f"✅ Pairing verification completed:")
                print(f"   Can pair: {can_pair}")
                print(f"   Reason: {reason}")
                print(f"   Trust score: {trust_score}")
                print(f"   TX Hash: {result.get('txhash')}")
                return can_pair
            else:
                print(f"❌ Pairing verification failed: {response.status_code}")
                print(f"   Response: {response.text}")
                return False
                
        except Exception as e:
            print(f"❌ Error verifying component pairing: {e}")
            return False

    def create_anonymous_pairing_authorization(self, creator: str, component_hash_a: str, component_hash_b: str) -> Optional[str]:
        """Create anonymous pairing authorization between two components"""
        # Generate a rule hash for the pairing rules
        rule_hash = hashlib.sha256(f"pairing_rules_{component_hash_a}_{component_hash_b}".encode()).hexdigest()
        
        payload = {
            "creator": creator,
            "component_hash_a": component_hash_a,
            "component_hash_b": component_hash_b,
            "rule_hash": rule_hash,
            "trust_score_requirement": "0.7",
            "authorization_level": "basic"
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/v1/components/authorization-anonymous",
                json=payload
            )
            
            if response.status_code == 200:
                result = response.json()
                auth_id = result.get("auth_id")
                print(f"✅ Anonymous pairing authorization created: {auth_id}")
                print(f"   Component A: {component_hash_a} (hash only)")
                print(f"   Component B: {component_hash_b} (hash only)")
                print(f"   Status: {result.get('status')}")
                print(f"   Expires: {result.get('expires_at')}")
                print(f"   TX Hash: {result.get('txhash')}")
                return auth_id
            else:
                print(f"❌ Anonymous pairing authorization failed: {response.status_code}")
                print(f"   Response: {response.text}")
                return None
                
        except Exception as e:
            print(f"❌ Error creating anonymous pairing authorization: {e}")
            return None

    def create_anonymous_revocation_event(self, creator: str, target_hash: str, revocation_type: str, urgency_level: str, reason_category: str) -> Optional[str]:
        """Create an anonymous revocation event"""
        payload = {
            "creator": creator,
            "target_hash": target_hash,
            "revocation_type": revocation_type,
            "urgency_level": urgency_level,
            "reason_category": reason_category,
            "context": "privacy_demo_revocation"
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/v1/components/revocation-anonymous",
                json=payload
            )
            
            if response.status_code == 200:
                result = response.json()
                revocation_id = result.get("revocation_id")
                print(f"✅ Anonymous revocation event created: {revocation_id}")
                print(f"   Target: {target_hash} (hash only)")
                print(f"   Type: {revocation_type}")
                print(f"   Urgency: {urgency_level}")
                print(f"   Reason: {reason_category}")
                print(f"   Status: {result.get('status')}")
                print(f"   Effective: {result.get('effective_at')}")
                print(f"   TX Hash: {result.get('txhash')}")
                return revocation_id
            else:
                print(f"❌ Anonymous revocation event failed: {response.status_code}")
                print(f"   Response: {response.text}")
                return None
                
        except Exception as e:
            print(f"❌ Error creating anonymous revocation event: {e}")
            return None

    def get_anonymous_component_metadata(self, requester: str, component_hash: str) -> Optional[Dict[str, Any]]:
        """Get anonymous metadata for a component hash"""
        try:
            response = self.session.get(
                f"{self.base_url}/api/v1/components/metadata-anonymous/{component_hash}",
                params={"requester": requester}
            )
            
            if response.status_code == 200:
                result = response.json()
                print(f"✅ Anonymous component metadata retrieved:")
                print(f"   Component hash: {result.get('component_hash')}")
                print(f"   Type: {result.get('type')} (generic only)")
                print(f"   Status: {result.get('status')}")
                print(f"   Trust anchor: {result.get('trust_anchor')}")
                print(f"   Last verified: {result.get('last_verified')}")
                print(f"   TX Hash: {result.get('txhash')}")
                return result
            else:
                print(f"❌ Anonymous component metadata retrieval failed: {response.status_code}")
                print(f"   Response: {response.text}")
                return None
                
        except Exception as e:
            print(f"❌ Error getting anonymous component metadata: {e}")
            return None

    def run_comprehensive_demo(self):
        """Run a comprehensive demo of all privacy-focused features"""
        print("🚀 Starting Privacy-Focused Component Registry Demo")
        print("=" * 60)
        
        # Step 1: Health checks
        print("\n📋 Step 1: System Health Checks")
        print("-" * 40)
        if not self.health_check():
            print("❌ Demo cannot continue - API Bridge is not available")
            return
        if not self.blockchain_status():
            print("❌ Demo cannot continue - Blockchain is not available")
            return
        
        # Step 2: Register components anonymously
        print("\n🔐 Step 2: Anonymous Component Registration")
        print("-" * 40)
        component_hashes = {}
        
        for component_key in self.demo_components.keys():
            user = self.demo_users["manufacturer_1"]
            component_hash = self.register_anonymous_component(component_key, user)
            if component_hash:
                component_hashes[component_key] = component_hash
            time.sleep(1)  # Small delay between registrations
        
        if not component_hashes:
            print("❌ No components were registered successfully")
            return
        
        # Step 3: Verify component pairing with hashes
        print("\n🔗 Step 3: Hash-Based Pairing Verification")
        print("-" * 40)
        component_keys = list(component_hashes.keys())
        if len(component_keys) >= 2:
            hash_a = component_hashes[component_keys[0]]
            hash_b = component_hashes[component_keys[1]]
            verifier = self.demo_users["verifier"]
            
            can_pair = self.verify_component_pairing_with_hashes(hash_a, hash_b, verifier)
            time.sleep(1)
        
        # Step 4: Create anonymous pairing authorization
        print("\n🔑 Step 4: Anonymous Pairing Authorization")
        print("-" * 40)
        if len(component_keys) >= 2:
            hash_a = component_hashes[component_keys[0]]
            hash_b = component_hashes[component_keys[1]]
            creator = self.demo_users["manufacturer_1"]
            
            auth_id = self.create_anonymous_pairing_authorization(creator, hash_a, hash_b)
            time.sleep(1)
        
        # Step 5: Get anonymous component metadata
        print("\n📊 Step 5: Anonymous Component Metadata Retrieval")
        print("-" * 40)
        for component_key, component_hash in component_hashes.items():
            requester = self.demo_users["verifier"]
            metadata = self.get_anonymous_component_metadata(requester, component_hash)
            time.sleep(1)
        
        # Step 6: Create anonymous revocation event
        print("\n🚫 Step 6: Anonymous Revocation Event")
        print("-" * 40)
        if component_hashes:
            target_hash = list(component_hashes.values())[0]
            creator = self.demo_users["regulator"]
            
            revocation_id = self.create_anonymous_revocation_event(
                creator=creator,
                target_hash=target_hash,
                revocation_type="INDIVIDUAL",
                urgency_level="HIGH",
                reason_category="SAFETY_CONCERN"
            )
            time.sleep(1)
        
        # Step 7: Demonstrate privacy benefits
        print("\n🛡️ Step 7: Privacy Benefits Demonstration")
        print("-" * 40)
        print("✅ Privacy Achievements:")
        print("   • Real component IDs are never exposed on-chain")
        print("   • Manufacturer identities are cryptographically hashed")
        print("   • Component specifications remain off-chain")
        print("   • Pairing verification works with hashes only")
        print("   • Revocation events protect component anonymity")
        print("   • Metadata retrieval reveals only non-sensitive data")
        print()
        print("🔒 Trade Secret Protection:")
        print("   • No commercial data on blockchain")
        print("   • Cryptographic hashes prevent reverse engineering")
        print("   • Anonymous operations maintain business confidentiality")
        print("   • Regulatory compliance without data exposure")
        
        print("\n🎉 Privacy-Focused Demo Completed Successfully!")
        print("=" * 60)

def main():
    """Main function to run the demo"""
    demo = PrivacyFocusedDemo()
    demo.run_comprehensive_demo()

if __name__ == "__main__":
    main() 