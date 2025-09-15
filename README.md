# ACT - Agentic Context Tool

## Overview

ACT is the human interface to Web4 - a complete implementation of the Agentic Context Protocol (ACP) that enables humans to interact with MCP servers through their Linked Context Tokens (LCTs). ACT bridges the gap between human intent and autonomous agent execution in the Web4 ecosystem.

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

## Contributing

This is currently a private repository. Contact the maintainers for access.

## License

Proprietary (during development)
Will transition to AGPL-3.0 for public release

## Contact

Dennis Palatov  
dp@metalinxx.io

---

*"ACT transforms Web4 from a protocol into a product - making trust-native computing accessible to everyone."*