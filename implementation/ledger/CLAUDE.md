# Claude Context for ACT Web4 Blockchain

## What This Is

ACT is a deployable blockchain reference build for Web4. Built on Cosmos SDK v0.53.0, it implements Web4's ontology on-chain: LCT identity, T3/V3 trust tensors, ATP/ADP energy economy, and MRH context boundaries.

**Design**: Offline-first, portable, settlement on demand. Any Web4 entity can instantiate its own copy.

## Modules (in `x/`)

| Module | Status | Key Files |
|--------|--------|-----------|
| `lctmanager` | Complete | keeper/, types/, genesis.go |
| `energycycle` | Complete | keeper/, types/, genesis.go |
| `trusttensor` | Complete | keeper/, types/, genesis.go |
| `componentregistry` | Complete | keeper/, types/, genesis.go |
| `pairing` | Complete | keeper/, types/, genesis.go |
| `pairingqueue` | Complete | keeper/, types/, genesis.go |
| `mrh` | Partial | 6 files, needs graph implementation |
| `societytodo` | Partial | Law oracle, governance rules |

## Build & Run

```bash
# Build (protos are pre-compiled)
ignite chain build --skip-proto

# Initialize and start
racecar-webd init mynode --chain-id racecarweb
racecar-webd keys add alice --keyring-backend test
racecar-webd genesis add-genesis-account alice 1000000000stake --keyring-backend test
racecar-webd genesis gentx alice 100000000stake --keyring-backend test --chain-id racecarweb
racecar-webd genesis collect-gentxs
racecar-webd start --api.enable --grpc.enable
```

Services: Tendermint RPC `:26657`, REST API `:1317`, Faucet `:4500`

## Build Requirements

- Go 1.24.0 (with sonic replace directives in go.mod)
- Ignite CLI v29.4.0-dev
- Cosmos SDK v0.53.0

## Web4 Ontological Context

```
Web4 = MCP + RDF + LCT + T3/V3*MRH + ATP/ADP
```

- `/` = "verified by", `*` = "contextualized by", `+` = "augmented with"
- Web4 is an **ontology**, not architecture or infrastructure
- Canonical T3 dimensions: **Talent / Training / Temperament** (per Web4 spec)
- Canonical V3 dimensions: **Valuation / Veracity / Validity**

## Integration Context

ACT is the **ledger of record** in the dp-web4 stack:
- **Web4** (spec) defines what ACT implements
- **SAGE** (cognition) registers identity and settles ATP on ACT
- **Hardbound** (oversight) records coherence attestations on ACT

Ecosystem integration plan: `private-context/plans/ecosystem-integration-plan.md`

## Key Issues

- T3 dimension names in `act/trusttensor/` proto use Competence/Reliability/Transparency instead of canonical Talent/Training/Temperament. The `racecarweb/` layer uses correct names. Phase 0 alignment pending.
- MRH module is partial — needs typed graph structure (`bound[]`, `paired[]`, `witnessing[]`, `broadcast[]`)
- Society 2 (democratic) blockchain doesn't build yet

## Push Command

```bash
grep GITHUB_PAT /mnt/c/projects/ai-agents/.env | cut -d= -f2 | xargs -I {} git push https://dp-web4:{}@github.com/dp-web4/ACT.git
```
