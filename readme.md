# ACT — Agentic Context Tool

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
![Status: Experimental](https://img.shields.io/badge/Status-Experimental-orange.svg)

> **Ledger Documentation**: Canonical ledger specifications and reference implementations are maintained in [web4/ledgers/](https://github.com/dp-web4/web4/tree/main/ledgers). This repository contains the Cosmos SDK implementation.

## What ACT Is

ACT is a **deployable blockchain reference build** for the Web4 ontology. Any Web4 entity — an AI agent fleet, a research group, a governance domain — can instantiate its own ACT copy as part of its memory and trust infrastructure.

It implements the Web4 equation on-chain:

```
Web4 = MCP + RDF + LCT + T3/V3*MRH + ATP/ADP
```

Where: `/` = "verified by", `*` = "contextualized by", `+` = "augmented with"

**Design principles**:
- **Offline-first**: Everything works locally. Chain settlement is opportunistic, not required.
- **Settlement on demand**: No periodic polling. Settle when meaningful — session boundary, trust milestone, federation event.
- **Portable**: Each deployment is self-contained. No dependency on a central always-on chain.
- **Rust when possible**: Go core stays (Cosmos SDK), new modules and tooling prefer Rust.

## Current State (March 2026)

Built on Cosmos SDK v0.53.0 with Ignite CLI. The chain runs with functional modules for identity, trust, energy, and governance.

### Modules

| Module | Status | Purpose |
|--------|--------|---------|
| `lctmanager` | Complete | LCT identity lifecycle — creation, pairing, revocation |
| `energycycle` | Complete | ATP/ADP economy — discharge, recharge, conservation |
| `trusttensor` | Complete | T3/V3 trust and value tensor calculations |
| `componentregistry` | Complete | Hardware component tracking and attestation |
| `pairing` | Complete | Device pairing and authentication |
| `pairingqueue` | Complete | Queue management for pairing operations |
| `mrh` | Partial | Markov Relevancy Horizon — context boundaries |
| `societytodo` | Partial | Society governance and law oracle |

### What Works

- Chain builds and runs (`ignite chain build --skip-proto`)
- Tendermint RPC (`:26657`), REST API (`:1317`), Token Faucet (`:4500`)
- Society 4 (monarchic governance) at 7.9/10 Web4 compliance
- Hardware binding with attestation hashes (software-level — TPM/TEE integration planned)
- Test accounts operational (Alice + Bob with genesis allocations)

### What's Dormant

- Society 2 (democratic governance) — partial, needs build fixes
- Federation witnessing — designed but not connected to SAGE fleet
- Cross-society interoperability — architecture exists, not tested end-to-end

## Architecture

```
Web4 Entity (SAGE instance, 4-life society, etc.)
    │
    │  instantiates
    ▼
┌─────────────────────────────────────────────┐
│  ACT Instance (self-contained)              │
│                                             │
│  Cosmos SDK chain with Web4 modules:        │
│  ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │lctmanager│ │energycycl│ │trusttensor │  │
│  └──────────┘ └──────────┘ └────────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ pairing  │ │pairingqu.│ │compreg.    │  │
│  └──────────┘ └──────────┘ └────────────┘  │
│  ┌──────────┐ ┌──────────┐                  │
│  │   mrh    │ │societytod│  (partial)       │
│  └──────────┘ └──────────┘                  │
└─────────────────────────────────────────────┘
```

## In the Ecosystem

ACT is the **ledger of record** in the dp-web4 stack:

- **[Web4](https://github.com/dp-web4/web4)** defines the ontology — ACT implements it on-chain
- **[SAGE](https://github.com/dp-web4/SAGE)** runs cognition at the edge — ACT records trust, identity, and ATP settlement
- **Hardbound** provides oversight — ACT stores coherence attestations and audit anchors

Currently the SAGE fleet runs ACT as a federation-level testbed. As integration matures, other Web4 entities will instantiate their own copies — [4-Life](https://github.com/dp-web4/4-life) is planned as the first portability prototype.

## Security Notice

**This is an experimental implementation. Do not use with real value at stake.**

Known limitations:
- Hardware binding uses static hashes — vulnerable to replay attacks. Production requires TPM/secure enclave integration.
- Identity verification lacks cryptographic proof of hardware possession.
- Federation witnessing is designed but not yet implemented.

Future security architecture includes TPM/TrustZone attestation, zero-knowledge proofs for identity, and federation witnessing across chains. See `PROPOSAL_004_HARDWARE_BINDING_ROOT_IDENTITY.md` for the design.

## Quick Start

```bash
# Prerequisites: Go 1.24+, Ignite CLI v29.4.0-dev

# Clone and build
git clone https://github.com/dp-web4/ACT.git
cd ACT/implementation/ledger

# Build (skip proto gen — protos are pre-compiled)
ignite chain build --skip-proto

# Initialize
racecar-webd init mynode --chain-id racecarweb
racecar-webd keys add alice --keyring-backend test
racecar-webd genesis add-genesis-account alice 1000000000stake --keyring-backend test
racecar-webd genesis gentx alice 100000000stake --keyring-backend test --chain-id racecarweb
racecar-webd genesis collect-gentxs

# Start
racecar-webd start --api.enable --grpc.enable
```

## Attribution & Licensing

### Patent Notice
LCT (Linked Context Token) technology is covered by U.S. Patents 11,477,027 and 12,278,913, owned by Metalinxx Inc. Licensed under AGPL-3.0.

### License
GNU Affero General Public License v3.0 — see [LICENSE](LICENSE)

## Contributing

We welcome contributions. This project is in active R&D:

1. Check Issues for current work
2. Fork & branch from `main`
3. Test locally — `ignite chain build --skip-proto` must pass
4. Submit PR with clear description

Areas where help is needed:
- **Rust modules**: New tooling and integration bridges
- **MRH completion**: Context boundary graph implementation
- **Federation**: Cross-instance identity and trust protocols
- **Testing**: Unit and integration test coverage

## Contact

Dennis Palatov
dp@metalinxx.io
