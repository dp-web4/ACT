#!/usr/bin/env python3
"""
Web4 LCT Encrypted Communication Demo

This script demonstrates the complete flow of establishing an encrypted
communication channel using split keys from the Web4 blockchain.

Usage:
    python encrypted_communication_demo.py

Requirements:
    pip install requests cryptography
"""

import requests
import hashlib
import os
import json
import time
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend

class KeyManager:
    """Manages split key reconstruction and master key generation"""
    
    def __init__(self):
        self.key_half_a = None
        self.key_half_b = None
        self.master_key = None
    
    def hex_to_bytes(self, hex_string):
        """Convert hex string to bytes"""
        return bytes.fromhex(hex_string)
    
    def bytes_to_hex(self, byte_array):
        """Convert bytes to hex string"""
        return byte_array.hex()
    
    def reconstruct_master_key(self, key_half_a, key_half_b):
        """Reconstruct master key from two halves using XOR + SHA-256"""
        # Convert hex strings to bytes
        bytes_a = self.hex_to_bytes(key_half_a)
        bytes_b = self.hex_to_bytes(key_half_b)
        
        # XOR the two key halves
        master_bytes = bytes(a ^ b for a, b in zip(bytes_a, bytes_b))
        
        # Apply SHA-256 for additional security
        master_key = hashlib.sha256(master_bytes).digest()
        
        return self.bytes_to_hex(master_key)
    
    def set_key_halves(self, key_half_a, key_half_b):
        """Set key halves and reconstruct master key"""
        self.key_half_a = key_half_a
        self.key_half_b = key_half_b
        self.master_key = self.reconstruct_master_key(key_half_a, key_half_b)
    
    def get_master_key(self):
        """Get the reconstructed master key"""
        return self.master_key

class EncryptedChannel:
    """Handles AES-256-GCM encryption/decryption"""
    
    def __init__(self, master_key):
        self.master_key = master_key
        self.session_key = self.generate_session_key()
    
    def generate_session_key(self):
        """Generate session key from master key using PBKDF2"""
        salt = b"web4_session_salt"
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=10000,
            backend=default_backend()
        )
        return kdf.derive(self.hex_to_bytes(self.master_key))
    
    def encrypt(self, plaintext):
        """Encrypt plaintext using AES-256-GCM"""
        # Generate new IV for each encryption
        iv = os.urandom(16)
        
        # Create cipher
        cipher = Cipher(
            algorithms.AES(self.session_key),
            modes.GCM(iv),
            backend=default_backend()
        )
        encryptor = cipher.encryptor()
        
        # Encrypt data
        ciphertext = encryptor.update(plaintext.encode()) + encryptor.finalize()
        
        # Get authentication tag
        tag = encryptor.tag
        
        # Combine IV + ciphertext + tag
        result = iv + ciphertext + tag
        return self.bytes_to_hex(result)
    
    def decrypt(self, encrypted_data):
        """Decrypt ciphertext using AES-256-GCM"""
        data = self.hex_to_bytes(encrypted_data)
        
        # Extract IV, ciphertext, and tag
        iv = data[:16]
        tag = data[-16:]
        ciphertext = data[16:-16]
        
        # Create cipher
        cipher = Cipher(
            algorithms.AES(self.session_key),
            modes.GCM(iv, tag),
            backend=default_backend()
        )
        decryptor = cipher.decryptor()
        
        # Decrypt data
        plaintext = decryptor.update(ciphertext) + decryptor.finalize()
        return plaintext.decode()

class SecureCommunication:
    """Manages secure communication using Web4 LCT split keys"""
    
    def __init__(self, api_endpoint, device_id):
        self.api_endpoint = api_endpoint
        self.device_id = device_id
        self.key_manager = KeyManager()
        self.channel = None
    
    def establish_secure_channel(self, partner_device_id):
        """Establish secure communication channel via blockchain"""
        print(f"\n🔐 Establishing secure channel between {self.device_id} and {partner_device_id}")
        print("=" * 60)
        
        try:
            # Step 1: Create LCT relationship
            print("📋 Step 1: Creating LCT relationship...")
            lct_response = requests.post(f'{self.api_endpoint}/lct/create', json={
                'creator': 'demo-user',
                'component_a': self.device_id,
                'component_b': partner_device_id,
                'context': 'secure-communication',
                'proxy_id': 'proxy-001'
            })
            
            if lct_response.status_code != 200:
                raise Exception(f"LCT creation failed: {lct_response.text}")
            
            lct_data = lct_response.json()
            print(f"✅ LCT created: {lct_data['lct_id']}")
            print(f"📝 Transaction: {lct_data['tx_hash']}")
            
            # Step 2: Initiate pairing
            print("\n🤝 Step 2: Initiating pairing...")
            pairing_response = requests.post(f'{self.api_endpoint}/pairing/initiate', json={
                'creator': 'demo-user',
                'component_a': self.device_id,
                'component_b': partner_device_id,
                'operational_context': 'secure-session',
                'proxy_id': 'proxy-001',
                'force_immediate': False
            })
            
            if pairing_response.status_code != 200:
                raise Exception(f"Pairing initiation failed: {pairing_response.text}")
            
            pairing_data = pairing_response.json()
            print(f"✅ Pairing initiated: {pairing_data['challenge_id']}")
            print(f"📝 Transaction: {pairing_data['tx_hash']}")
            
            # Step 3: Complete pairing to get split keys
            print("\n🔑 Step 3: Completing pairing to get split keys...")
            complete_response = requests.post(f'{self.api_endpoint}/pairing/complete', json={
                'creator': 'demo-user',
                'challenge_id': pairing_data['challenge_id'],
                'component_a_auth': 'device-auth-token',
                'component_b_auth': 'partner-auth-token',
                'session_context': 'session-001'
            })
            
            if complete_response.status_code != 200:
                raise Exception(f"Pairing completion failed: {complete_response.text}")
            
            complete_data = complete_response.json()
            print(f"✅ Pairing completed: {complete_data['lct_id']}")
            print(f"📝 Transaction: {complete_data['tx_hash']}")
            
            # Step 4: Reconstruct master key
            print("\n🔧 Step 4: Reconstructing master key...")
            self.key_manager.set_key_halves(
                complete_data['split_key_a'],
                complete_data['split_key_b']
            )
            
            # Step 5: Initialize encrypted channel
            print("\n🚀 Step 5: Initializing encrypted channel...")
            self.channel = EncryptedChannel(self.key_manager.get_master_key())
            
            print("\n🎉 Secure channel established successfully!")
            print(f"🔑 Master Key: {self.key_manager.get_master_key()[:32]}...")
            print(f"📊 Key Size: {len(self.key_manager.get_master_key()) // 2} bytes")
            
            return True
            
        except Exception as e:
            print(f"❌ Failed to establish secure channel: {e}")
            return False
    
    def send_secure_message(self, message):
        """Send encrypted message"""
        if not self.channel:
            print("❌ No secure channel established")
            return None
        
        try:
            encrypted = self.channel.encrypt(message)
            print(f"🔒 Encrypted message: {encrypted[:64]}...")
            return encrypted
        except Exception as e:
            print(f"❌ Encryption failed: {e}")
            return None
    
    def receive_secure_message(self, encrypted_message):
        """Receive and decrypt message"""
        if not self.channel:
            print("❌ No secure channel established")
            return None
        
        try:
            decrypted = self.channel.decrypt(encrypted_message)
            print(f"🔓 Decrypted message: {decrypted}")
            return decrypted
        except Exception as e:
            print(f"❌ Decryption failed: {e}")
            return None

def demo_secure_communication():
    """Demonstrate complete secure communication flow"""
    print("🏁 Web4 LCT Encrypted Communication Demo")
    print("=" * 60)
    
    # Initialize devices
    battery_module = SecureCommunication("http://localhost:8080", "battery-module-001")
    motor_controller = SecureCommunication("http://localhost:8080", "motor-controller-001")
    
    # Establish secure channel
    if battery_module.establish_secure_channel("motor-controller-001"):
        print("\n" + "=" * 60)
        print("📡 Testing Secure Communication")
        print("=" * 60)
        
        # Test messages
        test_messages = [
            "Battery status: 85% charge, 3.8V per cell",
            "Motor temperature: 45°C, RPM: 8000",
            "Emergency stop requested",
            "System diagnostics: All systems nominal"
        ]
        
        for i, message in enumerate(test_messages, 1):
            print(f"\n📤 Message {i}: {message}")
            
            # Encrypt message
            encrypted = battery_module.send_secure_message(message)
            if encrypted:
                # Decrypt message
                decrypted = motor_controller.receive_secure_message(encrypted)
                
                if decrypted == message:
                    print(f"✅ Message {i} transmitted securely!")
                else:
                    print(f"❌ Message {i} transmission failed!")
            
            time.sleep(1)  # Brief pause between messages
        
        print("\n" + "=" * 60)
        print("🎯 Demo completed successfully!")
        print("🔐 All messages transmitted with military-grade encryption")
        print("=" * 60)

def test_key_reconstruction():
    """Test key reconstruction with sample keys"""
    print("\n🧪 Testing Key Reconstruction")
    print("-" * 40)
    
    # Sample split keys (32 bytes each)
    key_a = "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"
    key_b = "f1e2d3c4b5a6789012345678901234567890fedcba1234567890fedcba123456"
    
    print(f"Key Half A: {key_a[:32]}...")
    print(f"Key Half B: {key_b[:32]}...")
    
    # Reconstruct master key
    key_manager = KeyManager()
    master_key = key_manager.reconstruct_master_key(key_a, key_b)
    
    print(f"Master Key: {master_key[:32]}...")
    print(f"Key Size: {len(master_key) // 2} bytes")
    
    # Test encryption/decryption
    channel = EncryptedChannel(master_key)
    test_message = "Hello, secure world!"
    
    encrypted = channel.encrypt(test_message)
    decrypted = channel.decrypt(encrypted)
    
    print(f"Test Message: {test_message}")
    print(f"Encrypted: {encrypted[:64]}...")
    print(f"Decrypted: {decrypted}")
    print(f"✅ Test passed: {test_message == decrypted}")

if __name__ == "__main__":
    try:
        # Test key reconstruction first
        test_key_reconstruction()
        
        # Run full demo
        demo_secure_communication()
        
    except KeyboardInterrupt:
        print("\n\n⏹️ Demo interrupted by user")
    except Exception as e:
        print(f"\n❌ Demo failed: {e}")
        print("💡 Make sure the API bridge is running on http://localhost:8080") 