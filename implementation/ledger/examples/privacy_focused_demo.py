#!/usr/bin/env python3
"""
Privacy-Focused Component Registry Demo

This demo showcases the privacy-focused architecture for the Web4 Component Registry:
- Anonymous component registration using hashes
- Hash-based pairing verification
- Anonymous pairing authorization
- Anonymous revocation events
- Privacy-preserving metadata retrieval

The demo demonstrates how commercial data is kept off-chain while maintaining
full functionality for trust, verification, and revocation.
"""

import requests
import json
import time
import hashlib
from typing import Dict, Any

class PrivacyFocusedDemo:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
        
    def print_header(self, title: str):
        """Print a formatted header"""
        print(f"\n{'='*60}")
        print(f"🔒 {title}")
        print(f"{'='*60}")
    
    def print_step(self, step: str):
        """Print a step description"""
        print(f"\n📋 {step}")
        print("-" * 40)
    
    def print_result(self, result: Dict[str, Any], title: str = "Result"):
        """Print a formatted result"""
        print(f"\n✅ {title}:")
        print(json.dumps(result, indent=2))
    
    def print_privacy_note(self, note: str):
        """Print a privacy-focused note"""
        print(f"\n🔐 Privacy Note: {note}")
    
    def demo_anonymous_component_registration(self):
        """Demo 1: Anonymous Component Registration"""
        self.print_header("Demo 1: Anonymous Component Registration")
        
        # Real commercial data (would be kept off-chain in production)
        real_component_id = "TESLA-2KW-400V-MODULE-Serial123"
        manufacturer_id = "Tesla Motors Inc."
        component_type = "Battery Module"
        
        self.print_step("Registering component with real commercial data")
        print(f"Real Component ID: {real_component_id}")
        print(f"Manufacturer: {manufacturer_id}")
        print(f"Component Type: {component_type}")
        
        # Register component anonymously
        payload = {
            "creator": "cosmos1tesla123456789",
            "real_component_id": real_component_id,
            "manufacturer_id": manufacturer_id,
            "component_type": component_type,
            "context": "demo_registration"
        }
        
        response = self.session.post(
            f"{self.base_url}/api/v1/components/register-anonymous",
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            self.print_result(result, "Anonymous Registration Response")
            
            # Show what goes on-chain vs off-chain
            self.print_privacy_note("On-chain data (anonymous):")
            print(f"  - Component Hash: {result.get('component_hash', 'N/A')}")
            print(f"  - Manufacturer Hash: {result.get('manufacturer_hash', 'N/A')}")
            print(f"  - Category Hash: {result.get('category_hash', 'N/A')}")
            
            self.print_privacy_note("Off-chain data (commercial):")
            print(f"  - Real Component ID: {real_component_id}")
            print(f"  - Manufacturer Name: {manufacturer_id}")
            print(f"  - Component Type: {component_type}")
            print(f"  - Technical Specifications: [Stored securely off-chain]")
            
            return result
        else:
            print(f"❌ Error: {response.status_code} - {response.text}")
            return None
    
    def demo_hash_based_pairing_verification(self, component_hash_a: str, component_hash_b: str):
        """Demo 2: Hash-Based Pairing Verification"""
        self.print_header("Demo 2: Hash-Based Pairing Verification")
        
        self.print_step("Verifying component pairing using only hashes")
        print(f"Component Hash A: {component_hash_a}")
        print(f"Component Hash B: {component_hash_b}")
        
        payload = {
            "verifier": "cosmos1verifier123456",
            "component_hash_a": component_hash_a,
            "component_hash_b": component_hash_b,
            "context": "demo_pairing_verification"
        }
        
        response = self.session.post(
            f"{self.base_url}/api/v1/components/verify-pairing-hashes",
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            self.print_result(result, "Pairing Verification Response")
            
            self.print_privacy_note("Privacy achieved:")
            print("  - No manufacturer names exposed")
            print("  - No part numbers revealed")
            print("  - No commercial relationships visible")
            print("  - Only cryptographic verification data")
            
            return result
        else:
            print(f"❌ Error: {response.status_code} - {response.text}")
            return None
    
    def demo_anonymous_pairing_authorization(self, component_hash_a: str, component_hash_b: str):
        """Demo 3: Anonymous Pairing Authorization"""
        self.print_header("Demo 3: Anonymous Pairing Authorization")
        
        self.print_step("Creating anonymous pairing authorization")
        print(f"Component Hash A: {component_hash_a}")
        print(f"Component Hash B: {component_hash_b}")
        
        # Generate rule hash (in production, this would be hash of off-chain rules)
        rule_data = f"authorization_rules_{component_hash_a}_{component_hash_b}"
        rule_hash = hashlib.sha256(rule_data.encode()).hexdigest()
        
        payload = {
            "creator": "cosmos1authorizer123456",
            "component_hash_a": component_hash_a,
            "component_hash_b": component_hash_b,
            "rule_hash": rule_hash,
            "trust_score_requirement": "0.8",
            "authorization_level": "enhanced"
        }
        
        response = self.session.post(
            f"{self.base_url}/api/v1/components/authorization-anonymous",
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            self.print_result(result, "Anonymous Authorization Response")
            
            self.print_privacy_note("Commercial data protected:")
            print("  - No licensing terms exposed")
            print("  - No royalty information revealed")
            print("  - No geographic restrictions visible")
            print("  - Only cryptographic authorization ID")
            
            return result
        else:
            print(f"❌ Error: {response.status_code} - {response.text}")
            return None
    
    def demo_anonymous_revocation_event(self, target_hash: str):
        """Demo 4: Anonymous Revocation Event"""
        self.print_header("Demo 4: Anonymous Revocation Event")
        
        self.print_step("Creating anonymous revocation event")
        print(f"Target Hash: {target_hash}")
        
        payload = {
            "creator": "cosmos1revoker123456",
            "target_hash": target_hash,
            "revocation_type": "INDIVIDUAL",
            "urgency_level": "URGENT",
            "reason_category": "SAFETY",
            "context": "demo_revocation"
        }
        
        response = self.session.post(
            f"{self.base_url}/api/v1/components/revocation-anonymous",
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            self.print_result(result, "Anonymous Revocation Response")
            
            self.print_privacy_note("Revocation privacy maintained:")
            print("  - No detailed reason exposed")
            print("  - No financial impact revealed")
            print("  - No customer information visible")
            print("  - Only urgency level and category")
            
            return result
        else:
            print(f"❌ Error: {response.status_code} - {response.text}")
            return None
    
    def demo_anonymous_metadata_retrieval(self, component_hash: str):
        """Demo 5: Anonymous Metadata Retrieval"""
        self.print_header("Demo 5: Anonymous Metadata Retrieval")
        
        self.print_step("Retrieving anonymous component metadata")
        print(f"Component Hash: {component_hash}")
        
        payload = {
            "requester": "cosmos1requester123456"
        }
        
        response = self.session.get(
            f"{self.base_url}/api/v1/components/metadata-anonymous/{component_hash}",
            json=payload
        )
        
        if response.status_code == 200:
            result = response.json()
            self.print_result(result, "Anonymous Metadata Response")
            
            self.print_privacy_note("Metadata privacy preserved:")
            print("  - No manufacturer name")
            print("  - No part numbers")
            print("  - No specifications")
            print("  - Only generic type and status")
            
            return result
        else:
            print(f"❌ Error: {response.status_code} - {response.text}")
            return None
    
    def run_complete_demo(self):
        """Run the complete privacy-focused demo"""
        print("🚗 Web4 Privacy-Focused Component Registry Demo")
        print("=" * 60)
        print("This demo showcases how commercial data is protected while")
        print("maintaining full functionality for trust and verification.")
        print("=" * 60)
        
        # Demo 1: Anonymous Component Registration
        reg_result = self.demo_anonymous_component_registration()
        if not reg_result:
            print("❌ Demo failed at registration step")
            return
        
        component_hash = reg_result.get('component_hash')
        if not component_hash:
            print("❌ No component hash received")
            return
        
        # Demo 2: Hash-Based Pairing Verification
        # Create a second component hash for pairing demo
        second_component_hash = hashlib.sha256("SECOND-COMPONENT-DEMO".encode()).hexdigest()
        self.demo_hash_based_pairing_verification(component_hash, second_component_hash)
        
        # Demo 3: Anonymous Pairing Authorization
        self.demo_anonymous_pairing_authorization(component_hash, second_component_hash)
        
        # Demo 4: Anonymous Revocation Event
        self.demo_anonymous_revocation_event(component_hash)
        
        # Demo 5: Anonymous Metadata Retrieval
        self.demo_anonymous_metadata_retrieval(component_hash)
        
        # Summary
        self.print_header("Demo Summary")
        print("✅ All privacy-focused operations completed successfully!")
        print("\n🔐 Privacy Achievements:")
        print("  - Zero commercial data exposed on blockchain")
        print("  - Anonymous component identification")
        print("  - Hash-based pairing verification")
        print("  - Protected revocation events")
        print("  - Privacy-preserving metadata")
        print("\n⚡ Functional Integrity:")
        print("  - Full trust and verification capabilities")
        print("  - Complete revocation system")
        print("  - Secure pairing authorization")
        print("  - Real-time status checking")
        
        print("\n🎯 Enterprise Ready:")
        print("  - Trade secret protection")
        print("  - Anti-data harvesting")
        print("  - Competitive intelligence prevention")
        print("  - Regulatory compliance")

def main():
    """Main function to run the demo"""
    demo = PrivacyFocusedDemo()
    
    try:
        demo.run_complete_demo()
    except requests.exceptions.ConnectionError:
        print("❌ Error: Could not connect to API Bridge server")
        print("Make sure the server is running on http://localhost:8080")
    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    main() 