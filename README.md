# ACT - Agentic Context Tool

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
![Status: Experimental](https://img.shields.io/badge/Status-Experimental-orange.svg)
![Progress: 65%](https://img.shields.io/badge/Progress-65%25-yellow.svg)

## 🚧 Development Status

**This repository is now PUBLIC but remains in EXPERIMENTAL/DEVELOPMENT stage.** The Web4 reference implementation is approximately 65% complete. Core protobuf definitions are finalized, keeper implementations are functional, but proto generation and module wiring are still in progress. Expect breaking changes as we iterate toward production readiness.

## ⚠️ Security Notice

**IMPORTANT: This is an experimental implementation with known security vulnerabilities.** Current implementations are for research and development purposes only:

### Known Security Issues
- **Hardware Binding**: Currently uses static hashes as placeholders - vulnerable to replay attacks. Production implementation will require TPM/secure enclave integration.
- **Identity Verification**: Lacks proper cryptographic proof of hardware possession. Anyone can copy published hardware hashes.
- **Federation Trust**: Missing piecewise identity factor verification across chains.
- **Key Management**: Some implementations may expose sensitive identifiers in logs or commits.

### Contributing to Security
We actively encourage security researchers and contributors to:
- **Identify vulnerabilities** - Help us find unknown security flaws
- **Propose solutions** - Submit PRs with security improvements
- **Document issues** - Create detailed vulnerability reports
- **Design protocols** - Help design federation witnessing and cross-chain identity verification

### Future Security Architecture
The production implementation will include:
- **TPM/Secure Enclave** integration for unforgeable hardware binding
- **Federation Witnessing** - Piecewise identity factors linked to LCTs across multiple chains
- **Zero-Knowledge Proofs** for identity verification without exposure
- **Hardware Security Modules** where available
- **Web of Trust** for platforms without hardware security features

**Do not use this implementation for production systems or with real value at stake.**

## Overview

ACT is the human interface to Web4 - a complete implementation of the Agentic Context Protocol (ACP) that enables humans to interact with MCP servers through their Linked Context Tokens (LCTs). ACT bridges the gap between human intent and autonomous agent execution in the Web4 ecosystem.

### Web4 Society Integration

ACT implements the foundational Web4 Society concept, where all entities (humans, agents, services) are citizens of self-governing digital societies. Each society maintains:
- **Laws**: Codified governance rules that all citizens must follow
- **Ledger**: Immutable record of all society events and citizen actions
- **Treasury**: Society-managed ATP/ADP token pools (tokens belong to society, not individuals)
- **Citizenship**: Witnessed relationships between entities and their societies

### Web4 as Living Standard

The Web4 standard is not fixed - it's a proposal that co-evolves with implementations. ACT serves as both reference implementation and research platform, with discoveries feeding back into the standard:

- **Bidirectional Alignment**: Standard guides experiments, experiments inform standard
- **Pattern Discovery**: ACT revealed roles as attention partitions, not power structures
- **Readiness Economy**: ~33% resources maintain "idle" readiness for resilience
- **Synthon Consciousness**: Human-AI collaboration creates new entity types

For complete Web4 specifications and philosophy, see: [github.com/dp-web4/web4](https://github.com/dp-web4/web4)

## Attribution & Licensing

### Swarm Architecture
The fractal swarm orchestration system is based on [Claude-Flow](https://github.com/ruvnet/claude-flow) (MIT Licensed), adapted for Web4 compliance and ATP economy tracking.

### Patent Notice
LCT (Linked Context Token) technology is covered by U.S. Patents 11,477,027 and 12,278,913, owned by Metalinxx Inc. The license to use this technology is granted under the terms of the GNU Affero General Public License v3.0 (AGPL-3.0).

## Vision

ACT enables humans to:
- **Own their digital presence** through personal LCTs
- **Delegate authority** to AI agents with precise scope and limits
- **Interact with MCP servers** using trust-based authentication
- **Build reputation** through successful interactions
- **Participate in the ATP economy** by creating and consuming value

## Architecture

```
Human User
    ↓
[ACT Client App]
    ↓
[Agent LCT] ← Paired with → [Human Root LCT]
    ↓
[ACP Engine]
    ↓
[MCP Bridges] → [MCP Servers]
    ↓
[Demo Society] → [LCT Registry + Law Oracle + ATP Pool]
```

## Core Components

### 1. Client Application
User interface for humans to:
- Manage their LCT identity
- Create and approve agent plans
- Monitor agent actions
- Review execution history
- Manage ATP balance

### 2. ACP Implementation
Complete Agentic Context Protocol:
- Agent Plans with triggers
- Intent generation and routing
- Decision collection (human/auto)
- Execution with proof-of-agency
- Result recording and witnessing

### 3. Demo Society
Minimal viable Web4 society implementing:
- **Society LCT**: The society's own identity token
- **Law Oracle**: Starting with "all decisions unanimous"
- **Society Treasury**: ATP/ADP pool owned by society (not individuals)
- **Citizenship Records**: Witnessed relationships on immutable ledger
- **Ledger Type**: Confined (citizens-only access) initially
- **Fractal Ready**: Can become citizen of larger societies

### 4. MCP Bridges
Connectors to existing MCP servers:
- Claude MCP
- OpenAI plugins
- Custom tools
- Web services
- Databases

## Quick Start

```bash
# Install dependencies
npm install

# Start demo society
npm run society:start

# Start ACP engine
npm run acp:start

# Launch client app
npm run client:dev

# Create your first LCT
npm run lct:create
```

## Use Cases

### Personal AI Assistant
```yaml
Plan: Daily Assistant
Triggers:
  - cron: "0 9 * * *"  # 9 AM daily
Actions:
  - Check emails (read-only)
  - Summarize news
  - Review calendar
  - Prepare briefing
Human Approval: Required for actions
ATP Cost: 10 per day
```

### Developer Automation
```yaml
Plan: Code Review Bot
Triggers:
  - event: "pull_request"
Actions:
  - Run tests
  - Check style
  - Analyze security
  - Post comments
Human Approval: Auto if tests pass
ATP Cost: 5 per review
```

### Business Process
```yaml
Plan: Invoice Processor
Triggers:
  - event: "invoice_received"
Actions:
  - Validate invoice
  - Check budget
  - Route for approval
  - Process payment
Human Approval: Required if > $1000
ATP Cost: Variable by amount
```

## Technical Stack

- **Frontend**: React/TypeScript for client app
- **Backend**: Node.js/TypeScript for ACP engine
- **Storage**: SQLite for demo society ledger
- **Crypto**: libsodium for Ed25519/X25519
- **RDF**: N3.js for MRH graphs
- **MCP**: Official MCP SDK for server communication

## Project Structure

```
ACT/
├── tool/                # 🚧 Experimental MCP Direct Interface
│   ├── src/            # Server and discovery logic
│   ├── public/         # Web UI for MCP interaction
│   └── README.md       # Tool documentation
├── core-spec/           # ACT-specific specifications
│   ├── human-lct-binding.md
│   ├── agent-pairing.md
│   ├── permission-model.md
│   └── ui-requirements.md
├── implementation/      # Core ACT implementation
│   ├── acp-engine/     # ACP protocol implementation
│   ├── lct-manager/    # LCT lifecycle management
│   ├── mrh-graph/      # RDF relationship graphs
│   └── atp-wallet/     # ATP/ADP token management
├── demo-society/       # Minimal society implementation
│   ├── registry/       # LCT registry
│   ├── law-oracle/     # Basic law engine
│   ├── witness/        # Witness network
│   └── ledger/         # Immutable record
├── client-app/         # Human interface
│   ├── web/           # Web application
│   ├── mobile/        # Mobile app (future)
│   └── cli/           # Command-line interface
├── mcp-bridges/        # MCP server connectors
│   ├── claude/        # Anthropic Claude
│   ├── openai/        # OpenAI GPT
│   ├── generic/       # Generic MCP
│   └── custom/        # Custom implementations
├── docs/              # Documentation
│   ├── user-guide/    # End-user documentation
│   ├── developer/     # Developer guides
│   └── api/           # API references
└── testing/           # Test suites
    ├── unit/          # Unit tests
    ├── integration/   # Integration tests
    └── e2e/           # End-to-end tests
```

### Experimental Tool

The [`tool/`](tool/) directory contains an **experimental prototype** for direct human interaction with MCP servers. This is a discovery and learning platform that allows exploration of MCP capabilities without AI intermediaries. [Learn more →](tool/README.md)

## Development Roadmap

### Phase 1: Foundation (Weeks 1-2)
- [ ] Basic LCT creation and management
- [ ] Simple ACP engine with manual triggers
- [ ] CLI interface for testing
- [ ] SQLite-based demo society

### Phase 2: Core Features (Weeks 3-4)
- [ ] Web UI with LCT dashboard
- [ ] Agent plan creation interface
- [ ] Basic MCP bridge (echo server)
- [ ] ATP wallet functionality

### Phase 3: Integration (Weeks 5-6)
- [ ] Claude MCP integration
- [ ] Multi-step plan execution
- [ ] Witness network implementation
- [ ] Trust tensor tracking

### Phase 4: Production (Weeks 7-8)
- [ ] Security audit
- [ ] Performance optimization
- [ ] Documentation completion
- [ ] Public demo deployment

## Business Model

ACT can operate as:

### 1. SaaS Platform
- Host demo society
- Provide client apps
- Charge for ATP tokens
- Premium features

### 2. Enterprise Solution
- Private society deployment
- Custom law oracle
- Integration services
- Support contracts

### 3. Open Source with Services
- Core technology free
- Paid hosting
- Professional services
- Training and certification

## Why ACT Matters

ACT is the missing piece that makes Web4 accessible to humans. By providing a user-friendly interface to the trust-native internet, ACT enables:

- **Individual sovereignty**: Own your digital identity
- **AI accountability**: Every action is traceable
- **Economic participation**: Earn ATP through value creation
- **Trust building**: Reputation that matters
- **Interoperability**: Work with any MCP server

## 🔨 What's Needed (Contributors Welcome!)

### Immediate Priorities
1. **Proto Generation**: Run `make proto-gen` and fix compilation errors
2. **Module Wiring**: Complete app.go integration for all 5 Web4 modules
3. **Genesis Configuration**: Define initial state for demo society
4. **Chain Startup**: Get the blockchain running with Web4 modules

### Help Wanted
- **Go Developers**: Cosmos SDK experience helpful but not required
- **Protocol Designers**: Help refine Web4 specifications
- **Documentation**: Improve setup guides and API docs
- **Testing**: Write unit and integration tests
- **Frontend**: Build demo UI for society interactions

## Contributing

We welcome contributions! This project is in active development:

1. **Check Issues**: See what needs work
2. **Fork & Branch**: Create feature branches from `main`
3. **Test Locally**: Ensure changes don't break existing functionality
4. **Submit PR**: Include clear description of changes
5. **Be Patient**: This is experimental - we're figuring it out together!

### Development Setup
```bash
# Clone the repository
git clone https://github.com/dp-web4/ACT.git
cd ACT

# Install Go 1.21+
# Install Node.js 18+

# Set up the ledger
cd implementation/ledger
make install

# Run the swarm monitor
cd ../../swarm-bootstrap
node monitor-swarm.js
```

## License

GNU Affero General Public License v3.0 (AGPL-3.0) - see [LICENSE](LICENSE) file

## Contact

Dennis Palatov  
dp@metalinxx.io

---

*"ACT transforms Web4 from a protocol into a product - making trust-native computing accessible to everyone."*