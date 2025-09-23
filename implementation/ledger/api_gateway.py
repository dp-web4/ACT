#!/usr/bin/env python3
"""
ACT Society API Gateway
Emergency implementation for Society4 onboarding
"""

import json
import hashlib
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import subprocess
import threading
import time
import os

class SocietyAPIHandler(BaseHTTPRequestHandler):
    
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path == '/api/v1/society/info':
            self.handle_society_info()
        elif path == '/api/v1/society/genesis':
            self.handle_genesis_download()
        elif path == '/api/v1/society/peers':
            self.handle_peers_info()
        elif path == '/api/v1/federation/status':
            self.handle_federation_status()
        elif path == '/health':
            self.handle_health_check()
        else:
            self.send_error(404, "Endpoint not found")
    
    def do_POST(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path == '/api/v1/society/join':
            self.handle_join_request()
        else:
            self.send_error(404, "Endpoint not found")
    
    def handle_society_info(self):
        """Return society information for discovery"""
        try:
            # Get current node info
            status = self.get_node_status()
            peers = self.get_net_info()
            
            info = {
                "society_id": "act-society-genesis",
                "chain_id": "act-web4",
                "moniker": "act-society",
                "version": "0.1.0",
                "network": {
                    "p2p": f"c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656",
                    "rpc": "http://10.0.0.72:26657",
                    "api": "http://10.0.0.72:1317",
                    "api_gateway": "http://10.0.0.72:8080"
                },
                "federation": {
                    "role": "genesis_validator",
                    "voting_power": 100,
                    "commission": 0.0
                },
                "capabilities": [
                    "lct_management",
                    "trust_tensor", 
                    "energy_cycle",
                    "task_delegation",
                    "society_todo",
                    "federation_resilience"
                ],
                "onboarding": {
                    "accepting_members": True,
                    "requirements": {
                        "minimum_stake": "100000stake",
                        "trust_score": 0.0,
                        "energy_commitment": 50
                    }
                },
                "status": {
                    "block_height": status.get("latest_block_height", "unknown"),
                    "block_time": status.get("latest_block_time", "unknown"),
                    "catching_up": status.get("catching_up", False),
                    "peers": peers.get("n_peers", "0")
                }
            }
            
            self.send_json_response(info)
            
        except Exception as e:
            self.send_error(500, f"Failed to get society info: {str(e)}")
    
    def handle_genesis_download(self):
        """Serve the genesis file for new nodes"""
        try:
            genesis_path = "/home/dp/ai-workspace/act/implementation/ledger/society/config/genesis.json"
            
            if os.path.exists(genesis_path):
                with open(genesis_path, 'r') as f:
                    genesis = json.load(f)
                
                self.send_json_response(genesis)
            else:
                self.send_error(404, "Genesis file not found")
                
        except Exception as e:
            self.send_error(500, f"Failed to serve genesis: {str(e)}")
    
    def handle_peers_info(self):
        """Return active peer information"""
        try:
            net_info = self.get_net_info()
            
            peers_response = {
                "persistent_peers": "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656",
                "active_peers": [],
                "total_peers": int(net_info.get("n_peers", 0))
            }
            
            # Add peer details from net_info
            if "peers" in net_info:
                for peer in net_info["peers"]:
                    peer_info = {
                        "id": peer["node_info"]["id"],
                        "moniker": peer["node_info"]["moniker"],
                        "remote_ip": peer["remote_ip"],
                        "status": "connected"
                    }
                    peers_response["active_peers"].append(peer_info)
            
            self.send_json_response(peers_response)
            
        except Exception as e:
            self.send_error(500, f"Failed to get peers info: {str(e)}")
    
    def handle_federation_status(self):
        """Return federation-wide status"""
        try:
            status = self.get_node_status()
            net_info = self.get_net_info()
            
            federation_status = {
                "active_societies": [
                    {
                        "id": "society-1-genesis",
                        "status": "online",
                        "block_height": int(status.get("latest_block_height", 0)),
                        "peers": int(net_info.get("n_peers", 0))
                    }
                ],
                "federation_health": 0.75 if int(net_info.get("n_peers", 0)) > 0 else 0.25,
                "consensus_status": "degraded" if int(net_info.get("n_peers", 0)) < 2 else "active",
                "pending_proposals": 3,  # Our three proposals
                "total_societies": 1 + int(net_info.get("n_peers", 0)),
                "chain_id": "act-web4"
            }
            
            # Add peer societies
            if "peers" in net_info:
                for peer in net_info["peers"]:
                    society = {
                        "id": peer["node_info"]["moniker"],
                        "status": "online",
                        "block_height": "unknown",
                        "peers": 1
                    }
                    federation_status["active_societies"].append(society)
            
            self.send_json_response(federation_status)
            
        except Exception as e:
            self.send_error(500, f"Failed to get federation status: {str(e)}")
    
    def handle_join_request(self):
        """Process society join request"""
        try:
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            request_data = json.loads(post_data.decode('utf-8'))
            
            # Basic validation
            required_fields = ["entity_id", "moniker", "type"]
            for field in required_fields:
                if field not in request_data:
                    self.send_error(400, f"Missing required field: {field}")
                    return
            
            # Generate onboarding response
            response = {
                "status": "approved",
                "onboarding_method": "genesis_sync",
                "genesis_url": "http://10.0.0.72:8080/api/v1/society/genesis",
                "sync_info": {
                    "trust_height": 1,
                    "trust_hash": "genesis",
                    "rpc_servers": [
                        "10.0.0.72:26657"
                    ]
                },
                "config": {
                    "chain_id": "act-web4",
                    "minimum_gas_prices": "0stake",
                    "persistent_peers": "c1a129e14fad4cb7c95f9e2b5e9586013941ebf5@10.0.0.72:26656"
                },
                "next_steps": [
                    "Download genesis file from genesis_url",
                    "Initialize node with provided chain_id",
                    "Configure persistent_peers in config.toml",
                    "Start node and wait for sync",
                    "Request validator status after sync completion"
                ],
                "estimated_sync_time": "5-10 minutes",
                "contact": "This API for support"
            }
            
            self.send_json_response(response)
            
        except json.JSONDecodeError:
            self.send_error(400, "Invalid JSON in request body")
        except Exception as e:
            self.send_error(500, f"Failed to process join request: {str(e)}")
    
    def handle_health_check(self):
        """Simple health check endpoint"""
        health = {
            "status": "healthy",
            "timestamp": int(time.time()),
            "version": "api-gateway-0.1.0"
        }
        self.send_json_response(health)
    
    def get_node_status(self):
        """Get blockchain node status"""
        try:
            result = subprocess.run(
                ["curl", "-s", "localhost:26657/status"],
                capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                data = json.loads(result.stdout)
                return data.get("result", {}).get("sync_info", {})
            else:
                return {}
        except:
            return {}
    
    def get_net_info(self):
        """Get network peer information"""
        try:
            result = subprocess.run(
                ["curl", "-s", "localhost:26657/net_info"],
                capture_output=True, text=True, timeout=5
            )
            
            if result.returncode == 0:
                data = json.loads(result.stdout)
                return data.get("result", {})
            else:
                return {}
        except:
            return {}
    
    def send_json_response(self, data):
        """Send JSON response with CORS headers"""
        response = json.dumps(data, indent=2)
        
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.send_header('Content-Length', str(len(response)))
        self.end_headers()
        
        self.wfile.write(response.encode('utf-8'))
    
    def do_OPTIONS(self):
        """Handle CORS preflight requests"""
        self.send_response(200)
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.end_headers()

def run_api_gateway(port=8080):
    """Start the API gateway server"""
    server_address = ('', port)
    httpd = HTTPServer(server_address, SocietyAPIHandler)
    
    print(f"🚀 ACT Society API Gateway starting on port {port}")
    print(f"📡 Endpoints available:")
    print(f"   GET  /api/v1/society/info     - Society discovery")
    print(f"   GET  /api/v1/society/genesis  - Genesis file download")
    print(f"   GET  /api/v1/society/peers    - Peer information")
    print(f"   GET  /api/v1/federation/status - Federation status")
    print(f"   POST /api/v1/society/join     - Join request")
    print(f"   GET  /health                   - Health check")
    print(f"")
    print(f"🎯 Ready for Society4 onboarding!")
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print(f"\n🛑 API Gateway shutting down...")
        httpd.server_close()

if __name__ == "__main__":
    run_api_gateway()