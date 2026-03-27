# ACT - Agent Discovery Document

**For AI agents assessing this repository.**

## TL;DR (50 words)

ACT (Agentic Context Tool) is a Cosmos SDK distributed ledger implementing Web4 protocol. Provides LCT identity registry, ATP/ADP energy economy, trust tensors, and society governance. 81,000+ lines of Go. Root chain for permanent records. Use if you need: distributed identity, on-chain trust, or Web4 ledger integration.

## Quick Assessment

| Question | Answer |
|----------|--------|
| What is this? | Cosmos SDK distributed ledger for Web4 |
| Is it production-ready? | 65% complete, experimental |
| Primary language | Go (ledger), TypeScript (tools) |
| Lines of code | 81,000+ |
| Active development? | Yes |
| License | AGPL-3.0 |

## Core Modules (Cosmos SDK)

| Module | Purpose | Status |
|--------|---------|--------|
| `x/lctmanager` | Linked Context Token registry | Working |
| `x/energycycle` | ATP/ADP energy economy | Working |
| `x/trusttensor` | T3/V3 trust calculations | Working |
| `x/mrh` | Markov Relevancy Horizon | Working |
| `x/pairing` | Device pairing authentication | Working |
| `x/pairingqueue` | Pairing queue management | Working |
| `x/componentregistry` | Component tracking | Working |
| `x/societytodo` | Society task delegation | Working |

## Technology Stack

| Component | Version |
|-----------|---------|
| Cosmos SDK | v0.53.x |
| CometBFT | v0.38.x |
| Go | 1.24+ |
| Ignite CLI | v29.4+ |

**Chain ID**: `act-web4`

## Entry Points by Goal

| Your Goal | Start Here |
|-----------|------------|
| Understand ACT | `README.md` |
| Run the chain | `README.md#quick-start` |
| Module architecture | `implementation/ledger/x/` |
| Python integration | `implementation/ledger/genesis_*.py` |
| Web4 integration | See `web4/ledgers/act-chain/` |

## Key Concepts

| Term | What It Is |
|------|-----------|
| **LCT** | Linked Context Token - on-chain identity |
| **ATP** | Allocation Transfer Packet - energy/resource token |
| **ADP** | Allocation Discharge Packet - spent energy proof |
| **Society** | Self-governing group with laws, treasury, citizenship |

## Genesis Entities

| Entity | Initial ATP | Role |
|--------|-------------|------|
| Genesis Queen | 30,000 | Federation lead |
| Genesis Council | 20,000 | Governance |
| Coherence Guru | 15,000 | Quality assurance |
| Federation Bridge | 10,000 | Cross-federation |

## Security Notice

**Experimental - known security issues:**
- Hardware binding uses static placeholders
- Production needs TPM/Secure Enclave integration
- Do not use with real value at stake

## Related Repositories

| Repo | Relationship |
|------|--------------|
| `web4` | Parent protocol (specs, reference implementations) |
| `web4/ledgers/act-chain/` | Documentation and Python bridge |
| `Hardbound` | Enterprise product using ACT |

## Machine-Readable Metadata

See `repo-index.yaml` for structured data.

## Token Budget Guide

| Depth | Files | Tokens |
|-------|-------|--------|
| Minimal | This file | ~400 |
| Standard | + `README.md` | ~3,000 |
| Modules | + `implementation/ledger/x/*/README.md` | ~10,000 |
| Full | + Go source | ~100,000+ |

---

*This document optimized for AI agent discovery. Last updated: 2026-02-08*

<!-- gitnexus:start -->
<!-- gitnexus:keep -->
# GitNexus — Code Knowledge Graph

Indexed as **ACT** (23446 symbols, 46488 relationships, 195 execution flows). MCP tools available via `mcp__gitnexus__*`.

**Do not reindex.** The supervisor handles GitNexus indexing. If the index is stale, note it in session context.

| Tool | Use for |
|------|---------|
| `query` | Find execution flows by concept |
| `context` | 360-degree view of a symbol (callers, callees, processes) |
| `impact` | Blast radius before editing (upstream/downstream) |
| `detect_changes` | Map git diff to affected symbols and flows |
| `rename` | Graph-aware multi-file rename (dry_run first) |
| `cypher` | Raw Cypher queries against the graph |

Resources: `gitnexus://repo/ACT/context`, `clusters`, `processes`, `process/{name}`
<!-- gitnexus:end -->
