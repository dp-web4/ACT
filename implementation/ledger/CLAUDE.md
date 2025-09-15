# Claude Context for Web4 ModBatt Demo

## Project Context System

**IMPORTANT**: A comprehensive context system exists at `/mnt/c/projects/ai-agents/misc/context-system/`

Quick access:
```bash
# Get overview of this project's role
cd /mnt/c/projects/ai-agents/misc/context-system
python3 query_context.py project web4-modbatt

# See blockchain integration plans
python3 query_context.py search "digital twin"

# Find related projects
cat /mnt/c/projects/ai-agents/misc/context-system/relationships/blockchain-bridge.md
```

## This Project's Role

Web4 ModBatt Demo bridges physical battery systems to blockchain:
- Built on Cosmos SDK with custom modules
- Provides digital twin functionality for battery hierarchy
- Enables decentralized energy markets and attestation

## Key Modules
- `energycycle`: Energy trading and lifecycle
- `trusttensor`: Trust and reputation metrics
- `pairing`: Device pairing and authentication
- `componentregistry`: Component tracking
- `lctmanager`: Lifecycle management

## Key Relationships
- **Bridges**: Physical battery systems to blockchain
- **Integrates With**: modbatt-CAN (documented integration)
- **Future**: Pack-Controller Bluetooth LE could connect directly

## Current Status
Scaffolded but core business logic not fully implemented. Represents the vision of battery systems as participants in decentralized energy networks, embodying distributed intelligence at the economic layer.