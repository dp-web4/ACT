# ACT - Agentic Context Tool

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
![Status: Experimental](https://img.shields.io/badge/Status-Experimental-orange.svg)
![Progress: 65%](https://img.shields.io/badge/Progress-65%25-yellow.svg)

## 🚧 Development Status

**This repository is now PUBLIC but remains in EXPERIMENTAL/DEVELOPMENT stage.** The Web4 reference implementation is approximately 65% complete. Core protobuf definitions are finalized, keeper implementations are functional, but proto generation and module wiring are still in progress. Expect breaking changes as we iterate toward production readiness.

## Overview

ACT is the human interface to Web4 - a complete implementation of the Agentic Context Protocol (ACP) that enables humans to interact with MCP servers through their Linked Context Tokens (LCTs). ACT bridges the gap between human intent and autonomous agent execution in the Web4 ecosystem.

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
Minimal viable society providing:
- LCT issuance and management
- Law Oracle with basic rules
- ATP/ADP token pool
- Witness network
- Immutable ledger

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