# Demo Society Implementation

## Overview

The Demo Society is a minimal but complete implementation of a Web4 society that provides the infrastructure for ACT to operate. It includes LCT registry, law oracle, ATP pool management, witness network, and immutable ledger.

## Components

### 1. LCT Registry
Manages the lifecycle of Linked Context Tokens:
- Birth certificate issuance
- LCT registration and lookup
- Revocation management
- Trust tensor tracking

### 2. Law Oracle
Enforces society rules:
- Permission validation
- ATP charging rules
- Dispute resolution
- Governance execution

### 3. ATP Pool
Economic system management:
- Token minting (ADP state)
- Charging authorization (ADP → ATP)
- Discharge tracking (ATP → ADP)
- Anti-hoarding enforcement

### 4. Witness Network
Trust building infrastructure:
- Event witnessing
- Attestation signing
- Trust calculation
- MRH graph updates

### 5. Immutable Ledger
Permanent record keeping:
- Transaction history
- Trust evolution
- Audit trails
- Proof generation

## Quick Start

```bash
# Install dependencies
npm install

# Initialize database
npm run db:init

# Start all services
npm run start

# Or start individually:
npm run registry:start
npm run oracle:start
npm run pool:start
npm run witness:start
npm run ledger:start
```

## Configuration

### Society Parameters
```yaml
# config/society.yaml
society:
  name: "ACT Demo Society"
  id: "soc:web4:demo:act"
  version: "1.0.0"
  
governance:
  type: "democratic"
  proposal_threshold: 3
  voting_period: 86400  # 24 hours
  quorum: 0.51

economy:
  initial_pool: 1000000  # ADP tokens
  charge_rate: 1.0       # 1 ADP = 1 ATP
  discharge_fee: 0.01    # 1% transaction fee
  max_accumulation: 1000 # Per entity

trust:
  initial_score: 0.5
  witness_weight: 0.1
  decay_rate: 0.001
  minimum_witnesses: 3
```

## API Endpoints

### Registry Service
```typescript
// POST /registry/lct/create
{
  "entity_type": "human",
  "public_key": "...",
  "witnesses": ["lct:1", "lct:2", "lct:3"],
  "metadata": {}
}

// GET /registry/lct/:id
// Returns LCT details and current trust tensor

// POST /registry/lct/:id/revoke
{
  "reason": "compromised",
  "signature": "..."
}
```

### Law Oracle Service
```typescript
// POST /oracle/validate
{
  "action": "transfer",
  "actor": "lct:123",
  "parameters": {},
  "context": {}
}

// GET /oracle/rules
// Returns current society rules

// POST /oracle/proposal
{
  "type": "rule_change",
  "description": "...",
  "changes": {}
}
```

### ATP Pool Service
```typescript
// POST /pool/charge
{
  "producer": "lct:producer:1",
  "amount": 100,
  "proof_of_value": "..."
}

// POST /pool/discharge
{
  "from": "lct:123",
  "to": "lct:456",
  "amount": 10,
  "purpose": "service_payment"
}

// GET /pool/balance/:lct_id
// Returns ATP/ADP balance
```

### Witness Service
```typescript
// POST /witness/attest
{
  "event": "pairing",
  "subject": "lct:123",
  "object": "lct:456",
  "signature": "..."
}

// GET /witness/events/:lct_id
// Returns witnessed events

// POST /witness/subscribe
{
  "lct_id": "lct:123",
  "event_types": ["pairing", "transfer"]
}
```

### Ledger Service
```typescript
// GET /ledger/block/:height
// Returns block at height

// GET /ledger/tx/:hash
// Returns transaction details

// GET /ledger/history/:lct_id
// Returns transaction history for LCT
```

## Database Schema

### SQLite Tables
```sql
-- LCT Registry
CREATE TABLE lcts (
    id TEXT PRIMARY KEY,
    entity_type TEXT NOT NULL,
    public_key TEXT UNIQUE NOT NULL,
    created_at INTEGER NOT NULL,
    revoked_at INTEGER,
    trust_competence REAL DEFAULT 0.5,
    trust_reliability REAL DEFAULT 0.5,
    trust_transparency REAL DEFAULT 1.0
);

-- ATP Pool
CREATE TABLE balances (
    lct_id TEXT PRIMARY KEY,
    atp_balance INTEGER DEFAULT 0,
    adp_balance INTEGER DEFAULT 0,
    last_updated INTEGER NOT NULL,
    FOREIGN KEY (lct_id) REFERENCES lcts(id)
);

-- Witness Network
CREATE TABLE attestations (
    id TEXT PRIMARY KEY,
    witness_lct TEXT NOT NULL,
    subject_lct TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_data TEXT,
    signature TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (witness_lct) REFERENCES lcts(id),
    FOREIGN KEY (subject_lct) REFERENCES lcts(id)
);

-- Immutable Ledger
CREATE TABLE blocks (
    height INTEGER PRIMARY KEY,
    hash TEXT UNIQUE NOT NULL,
    previous_hash TEXT NOT NULL,
    merkle_root TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    transactions TEXT NOT NULL
);

CREATE TABLE transactions (
    hash TEXT PRIMARY KEY,
    block_height INTEGER NOT NULL,
    from_lct TEXT NOT NULL,
    to_lct TEXT,
    type TEXT NOT NULL,
    data TEXT NOT NULL,
    signature TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    FOREIGN KEY (block_height) REFERENCES blocks(height)
);
```

## Security Considerations

### Key Management
- Private keys never stored
- Only public keys in database
- Hardware security module support
- Key rotation protocols

### Access Control
- API key authentication
- Rate limiting per endpoint
- IP whitelisting option
- Audit logging

### Data Integrity
- Merkle tree verification
- Signature validation
- Hash chain integrity
- Backup and recovery

## Monitoring

### Metrics
```yaml
System Metrics:
  - LCT creation rate
  - Transaction throughput
  - Witness participation
  - ATP circulation

Performance Metrics:
  - API response times
  - Database query times
  - Block creation time
  - Network latency

Health Checks:
  - Service availability
  - Database connectivity
  - Disk usage
  - Memory usage
```

### Logging
```typescript
// Structured logging
logger.info('LCT created', {
  lct_id: 'lct:123',
  entity_type: 'human',
  witnesses: 3,
  timestamp: Date.now()
});

// Log levels
// ERROR: System errors
// WARN: Anomalies
// INFO: Important events
// DEBUG: Detailed debugging
```

## Testing

### Unit Tests
```bash
npm run test:unit
```

### Integration Tests
```bash
npm run test:integration
```

### Load Testing
```bash
npm run test:load
```

### Test Coverage
```bash
npm run test:coverage
```

## Deployment

### Docker
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 3000-3005
CMD ["npm", "start"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  registry:
    build: .
    command: npm run registry:start
    ports:
      - "3001:3001"
  
  oracle:
    build: .
    command: npm run oracle:start
    ports:
      - "3002:3002"
  
  pool:
    build: .
    command: npm run pool:start
    ports:
      - "3003:3003"
  
  witness:
    build: .
    command: npm run witness:start
    ports:
      - "3004:3004"
  
  ledger:
    build: .
    command: npm run ledger:start
    ports:
      - "3005:3005"
  
  db:
    image: postgres:14
    environment:
      POSTGRES_DB: demo_society
      POSTGRES_PASSWORD: secure_password
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
```

## Maintenance

### Backup
```bash
# Backup database
npm run db:backup

# Restore database
npm run db:restore backup_file.sql
```

### Updates
```bash
# Check for updates
npm outdated

# Update dependencies
npm update

# Run migrations
npm run db:migrate
```

### Monitoring
```bash
# View logs
npm run logs

# Check health
npm run health

# View metrics
npm run metrics
```

---

*"The Demo Society is Web4 in miniature—fully functional, completely transparent, and ready to demonstrate the future of trust-native computing."*