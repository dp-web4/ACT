# Web4 Race Car Demo - Complete API Reference

**Version**: 1.0  
**Chain ID**: `racecarweb`  
**Base URL**: `https://your-domain.com/api`  
**Documentation Format**: Confluence/Jira Compatible

---

## Table of Contents

1. [Authentication](#authentication)
2. [Component Registry API](#component-registry-api)
3. [Pairing Queue API](#pairing-queue-api)
4. [LCT Manager API](#lct-manager-api)
5. [Pairing API](#pairing-api)
6. [Trust Tensor API](#trust-tensor-api)
7. [Energy Cycle API](#energy-cycle-api)
8. [WebSocket Events](#websocket-events)
9. [Error Codes](#error-codes)
10. [Rate Limits](#rate-limits)

---

## Authentication

### Overview
Web4 uses Cosmos SDK standard authentication with account addresses and transaction signing.

### Demo Account
```
Address: cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k
Mnemonic: Available in development environment
```

### Transaction Format
All state-changing operations use the standard Cosmos transaction format:

```json
{
  "tx": {
    "body": {
      "messages": [/* Message objects */],
      "memo": "",
      "timeout_height": "0"
    }
  },
  "mode": "BROADCAST_MODE_SYNC"
}
```

---

## Component Registry API

### Register Component

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Register a new ModBatt component in the Web4 system

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.componentregistry.v1.MsgRegisterComponent",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "component_id": "MODBATT-MOD-RC001-001",
        "component_type": "module",
        "manufacturer_data": "ModBatt Pro Series MB-4680-LFP"
      }]
    }
  },
  "mode": "BROADCAST_MODE_SYNC"
}
```

**Response**:
```json
{
  "tx_response": {
    "code": 0,
    "txhash": "ABC123DEF456...",
    "height": "12345",
    "events": [
      {
        "type": "component_registered",
        "attributes": [
          {"key": "component_id", "value": "MODBATT-MOD-RC001-001"},
          {"key": "component_type", "value": "module"},
          {"key": "status", "value": "registered"}
        ]
      }
    ]
  }
}
```

**Parameters**:
- `creator`: Cosmos account address
- `component_id`: Unique component identifier (manufacturer format)
- `component_type`: One of: `module`, `pack`, `host`
- `manufacturer_data`: Component specifications and metadata

---

### Get Component

**Endpoint**: `GET /racecarweb/componentregistry/v1/component/{component_id}`

**Description**: Retrieve component information

**Path Parameters**:
- `component_id`: Component identifier

**Response**:
```json
{
  "component_identity": {
    "component_id": "MODBATT-MOD-RC001-001",
    "component_type": "module",
    "manufacturer_data": "ModBatt Pro Series MB-4680-LFP",
    "device_key_half": "encrypted_key_material",
    "status": "active",
    "last_seen": 1704067200
  }
}
```

---

### Check Pairing Authorization

**Endpoint**: `GET /racecarweb/componentregistry/v1/pairing-auth`

**Description**: Verify bidirectional pairing authorization between components

**Query Parameters**:
- `component_a`: First component ID
- `component_b`: Second component ID

**Response**:
```json
{
  "a_can_pair_b": true,
  "b_can_pair_a": true,
  "reason": "bidirectional pairing authorized"
}
```

---

### Update Component Authorization

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.componentregistry.v1.MsgUpdateAuthorization",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "component_id": "MODBATT-MOD-RC001-001",
        "auth_rules": "{\"allowed_partner_types\": [\"pack\"], \"max_connections\": 2}"
      }]
    }
  }
}
```

---

## Pairing Queue API

### Queue Pairing Request

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Queue a pairing request for offline processing

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.pairingqueue.v1.MsgQueuePairingRequest",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "initiator_id": "MODBATT-MOD-RC001-001",
        "target_id": "MODBATT-PACK-RC001-A",
        "request_type": "bidirectional_pairing",
        "proxy_id": "MODBATT-HOST-RC001"
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "pairing_request_queued",
        "attributes": [
          {"key": "request_id", "value": "req-12345-67890"},
          {"key": "status", "value": "queued"}
        ]
      }
    ]
  }
}
```

---

### Get Queued Requests

**Endpoint**: `GET /racecarweb/pairingqueue/v1/component/{component_id}/requests`

**Description**: Retrieve all queued requests for a component

**Response**:
```json
{
  "pairing_requests": [
    {
      "request_id": "req-12345-67890",
      "initiator_id": "MODBATT-MOD-RC001-001",
      "target_id": "MODBATT-PACK-RC001-A",
      "status": "queued",
      "created_at": 1704067200
    }
  ]
}
```

---

### Process Offline Queue

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.pairingqueue.v1.MsgProcessOfflineQueue",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "component_id": "MODBATT-MOD-RC001-001"
      }]
    }
  }
}
```

---

## LCT Manager API

### Create LCT Relationship

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Create a new LCT relationship between two components

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.lctmanager.v1.MsgCreateLctRelationship",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "component_a": "MODBATT-MOD-RC001-001",
        "component_b": "MODBATT-PACK-RC001-A",
        "context": "race_car_operation",
        "proxy_id": "MODBATT-HOST-RC001"
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "lct_relationship_created",
        "attributes": [
          {"key": "lct_id", "value": "lct-MODBATT-MOD-001-PACK-A-1704067200"},
          {"key": "key_exchange", "value": "device_key_half_encrypted"},
          {"key": "status", "value": "active"}
        ]
      }
    ]
  }
}
```

---

### Get LCT Relationship

**Endpoint**: `GET /racecarweb/lctmanager/v1/lct/{lct_id}`

**Description**: Retrieve LCT relationship information

**Response**:
```json
{
  "linked_context_token": {
    "lct_id": "lct-MODBATT-MOD-001-PACK-A-1704067200",
    "component_a_id": "MODBATT-MOD-RC001-001",
    "component_b_id": "MODBATT-PACK-RC001-A",
    "lct_key_half": "encrypted_lct_key_material",
    "pairing_status": "active",
    "created_at": 1704067200,
    "updated_at": 1704067200,
    "last_contact_at": 1704067200,
    "trust_anchor": "trust-anchor-MOD-001-PACK-A",
    "operational_context": "race_car_operation",
    "proxy_component_id": "MODBATT-HOST-RC001",
    "authorization_rules": ""
  }
}
```

---

### Get Component Relationships

**Endpoint**: `GET /racecarweb/lctmanager/v1/component/{component_id}/relationships`

**Description**: Get all LCT relationships for a component (many-to-many support)

**Response**:
```json
{
  "component_relationships": [
    {
      "lct_id": "lct-MODBATT-MOD-001-PACK-A-1704067200",
      "component_a_id": "MODBATT-MOD-RC001-001",
      "component_b_id": "MODBATT-PACK-RC001-A",
      "pairing_status": "active",
      "operational_context": "race_car_operation"
    }
  ],
  "lct_count": 1
}
```

---

### Terminate LCT Relationship

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.lctmanager.v1.MsgTerminateLctRelationship",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "lct_id": "lct-MODBATT-MOD-001-PACK-A-1704067200",
        "reason": "maintenance_required",
        "notify_offline": true
      }]
    }
  }
}
```

---

## Pairing API

### Initiate Bidirectional Pairing

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Start bidirectional pairing process between components

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.pairing.v1.MsgInitiateBidirectionalPairing",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "component_a": "MODBATT-MOD-RC001-001",
        "component_b": "MODBATT-PACK-RC001-A",
        "operational_context": "race_car_battery_system",
        "proxy_id": "MODBATT-HOST-RC001",
        "force_immediate": true
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "bidirectional_pairing_initiated",
        "attributes": [
          {"key": "challenge_id", "value": "challenge-12345-67890"},
          {"key": "lct_id", "value": "lct-MODBATT-MOD-001-PACK-A-1704067200"},
          {"key": "status", "value": "active"},
          {"key": "queue_id", "value": ""}
        ]
      }
    ]
  }
}
```

---

### Validate Bidirectional Authorization

**Endpoint**: `GET /racecarweb/pairing/v1/validate-auth`

**Query Parameters**:
- `component_a`: First component ID
- `component_b`: Second component ID
- `context`: Operational context

**Response**:
```json
{
  "a_can_pair_b": true,
  "b_can_pair_a": true,
  "required_conditions": "trust_threshold_met,proximity_verified"
}
```

---

### Get Pairing Status

**Endpoint**: `GET /racecarweb/pairing/v1/challenge/{challenge_id}`

**Response**:
```json
{
  "pairing_challenge": {
    "challenge_id": "challenge-12345-67890",
    "requester_component": "MODBATT-MOD-RC001-001",
    "target_component": "MODBATT-PACK-RC001-A",
    "challenge_data": "encrypted_challenge_material",
    "response_data": "",
    "expires_at": 1704070800,
    "status": "pending",
    "proxy_component": "MODBATT-HOST-RC001"
  }
}
```

---

### List Active Pairings

**Endpoint**: `GET /racecarweb/pairing/v1/component/{component_id}/active`

**Response**:
```json
{
  "active_lcts": [
    "lct-MODBATT-MOD-001-PACK-A-1704067200",
    "lct-MODBATT-MOD-001-PACK-B-1704067300"
  ],
  "pairing_count": 2
}
```

---

## Trust Tensor API

### Create Relationship Tensor

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Create a trust tensor for an LCT relationship

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.trusttensor.v1.MsgCreateRelationshipTensor",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "lct_id": "lct-MODBATT-MOD-001-PACK-A-1704067200",
        "tensor_type": "T3",
        "context": "race_car_operation"
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "relationship_tensor_created",
        "attributes": [
          {"key": "tensor_id", "value": "tensor-T3-12345-67890"}
        ]
      }
    ]
  }
}
```

---

### Update Tensor Score

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Update a trust tensor score dimension

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.trusttensor.v1.MsgUpdateTensorScore",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "tensor_id": "tensor-T3-12345-67890",
        "dimension": "training",
        "value": "0.85",
        "context": "performance_improvement",
        "witness_data": "race_data_validation_hash"
      }]
    }
  }
}
```

---

### Calculate Relationship Trust

**Endpoint**: `GET /racecarweb/trusttensor/v1/relationship-trust`

**Query Parameters**:
- `lct_id`: LCT relationship identifier
- `context`: Calculation context

**Response**:
```json
{
  "trust_score": "0.750",
  "factors": "T3_composite:talent(0.7)+training(0.85)+temperament(0.65)"
}
```

---

### Get Relationship Tensor

**Endpoint**: `GET /racecarweb/trusttensor/v1/relationship-tensor`

**Query Parameters**:
- `lct_id`: LCT relationship identifier
- `tensor_type`: Tensor type (`T3` or `V3`)

**Response**:
```json
{
  "relationship_trust_tensor": {
    "tensor_id": "tensor-T3-12345-67890",
    "lct_id": "lct-MODBATT-MOD-001-PACK-A-1704067200",
    "tensor_type": "T3",
    "talent_score": "0.70",
    "training_score": "0.85",
    "temperament_score": "0.65",
    "context": "race_car_operation",
    "created_at": 1704067200,
    "updated_at": 1704067800,
    "version": 3
  }
}
```

---

## Energy Cycle API

### Create Relationship Energy Operation

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Create energy operation between LCT relationships (generates ATP tokens)

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.energycycle.v1.MsgCreateRelationshipEnergyOperation",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "source_lct": "lct-MODBATT-PACK-A-HOST-1704067200",
        "target_lct": "lct-MODBATT-MOD-001-PACK-A-1704067200",
        "energy_amount": "50.0",
        "operation_type": "discharge"
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "energy_operation_created",
        "attributes": [
          {"key": "operation_id", "value": "energy-op-discharge-12345"},
          {"key": "atp_tokens", "value": "ATP-001,ATP-002,ATP-003"},
          {"key": "trust_validated", "value": "true"}
        ]
      }
    ]
  }
}
```

---

### Execute Energy Transfer

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.energycycle.v1.MsgExecuteEnergyTransfer",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "operation_id": "energy-op-discharge-12345",
        "transfer_data": "race_car_energy_transfer_validated"
      }]
    }
  }
}
```

---

### Validate Relationship Value

**Endpoint**: `POST /cosmos/tx/v1beta1/txs`

**Description**: Validate energy operation value (converts ATP to ADP tokens)

**Request Body**:
```json
{
  "tx": {
    "body": {
      "messages": [{
        "@type": "/racecarweb.energycycle.v1.MsgValidateRelationshipValue",
        "creator": "cosmos1k0kju6a63dmfz9dhjx0ralxp60fewtf29wwd9k",
        "operation_id": "energy-op-discharge-12345",
        "recipient_validation": "confirmed",
        "utility_rating": "0.90",
        "trust_context": "race_performance_validated"
      }]
    }
  }
}
```

**Response**:
```json
{
  "tx_response": {
    "events": [
      {
        "type": "energy_value_validated",
        "attributes": [
          {"key": "v3_score", "value": "V3-tensor-67890"},
          {"key": "adp_tokens", "value": "ADP-001,ADP-002,ADP-003"}
        ]
      }
    ]
  }
}
```

---

### Get Relationship Energy Balance

**Endpoint**: `GET /racecarweb/energycycle/v1/relationship/{lct_id}/balance`

**Response**:
```json
{
  "atp_balance": "125.500",
  "adp_balance": "75.200",
  "total_energy": "200.700",
  "trust_weighted_balance": "150.525"
}
```

---

### Get Energy Flow History

**Endpoint**: `GET /racecarweb/energycycle/v1/energy-flow-history`

**Query Parameters**:
- `lct_id`: LCT relationship identifier

**Response**:
```json
{
  "energy_operations": [
    {
      "operation_id": "energy-op-discharge-12345",
      "source_lct": "lct-MODBATT-PACK-A-HOST-1704067200",
      "target_lct": "lct-MODBATT-MOD-001-PACK-A-1704067200",
      "energy_amount": "50.0",
      "operation_type": "discharge",
      "status": "validated",
      "timestamp": 1704067200,
      "trust_score": "0.85"
    }
  ]
}
```

---

## WebSocket Events

### Connection

**Endpoint**: `wss://your-domain.com/websocket`

**Connection Example**:
```javascript
const ws = new WebSocket('wss://your-domain.com/websocket');

ws.onopen = function() {
  // Subscribe to events
  ws.send(JSON.stringify({
    "jsonrpc": "2.0",
    "method": "subscribe",
    "id": 1,
    "params": {
      "query": "tm.event='Tx' AND lct_relationship_created.lct_id EXISTS"
    }
  }));
};
```

### Available Event Types

#### LCT Relationship Events
```javascript
// Subscribe to LCT relationship creation
{
  "query": "tm.event='Tx' AND lct_relationship_created.lct_id EXISTS"
}

// Subscribe to LCT termination
{
  "query": "tm.event='Tx' AND lct_relationship_terminated.lct_id EXISTS"
}
```

#### Pairing Events
```javascript
// Subscribe to pairing initiation
{
  "query": "tm.event='Tx' AND bidirectional_pairing_initiated.challenge_id EXISTS"
}

// Subscribe to pairing completion
{
  "query": "tm.event='Tx' AND pairing_completed.lct_id EXISTS"
}
```

#### Energy Cycle Events
```javascript
// Subscribe to energy operations
{
  "query": "tm.event='Tx' AND energy_operation_created.operation_id EXISTS"
}

// Subscribe to ATP/ADP conversions
{
  "query": "tm.event='Tx' AND energy_value_validated.v3_score EXISTS"
}
```

#### Trust Tensor Events
```javascript
// Subscribe to trust updates
{
  "query": "tm.event='Tx' AND relationship_tensor_created.tensor_id EXISTS"
}
```

### Event Response Format

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "query": "tm.event='Tx'",
    "data": {
      "type": "tendermint/event/Tx",
      "value": {
        "TxResult": {
          "height": "12345",
          "tx": "base64_encoded_tx",
          "result": {
            "events": [
              {
                "type": "lct_relationship_created",
                "attributes": [
                  {"key": "lct_id", "value": "lct-12345"},
                  {"key": "status", "value": "active"}
                ]
              }
            ]
          }
        }
      }
    }
  }
}
```

---

## Error Codes

### HTTP Status Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | Request successful |
| 400 | Bad Request | Invalid request format |
| 401 | Unauthorized | Invalid or missing authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 502 | Bad Gateway | Blockchain node unavailable |
| 503 | Service Unavailable | System maintenance |

### Web4 Specific Error Codes

| Code | Message | Description | Resolution |
|------|---------|-------------|------------|
| 5001 | Component not found | Component ID not in registry | Verify component ID and register if needed |
| 5002 | Pairing not authorized | Components cannot pair with each other | Check bidirectional authorization rules |
| 5003 | LCT relationship exists | Relationship already established | Use existing LCT or terminate first |
| 5004 | Bidirectional auth failed | One-way authorization only | Ensure both components authorize pairing |
| 5005 | Component offline | Target component unavailable | Use proxy component or queue operation |
| 5006 | Trust threshold not met | Insufficient trust score for operation | Improve component trust score |
| 5007 | Invalid energy amount | Energy amount outside valid range | Use value between 0.1 and 1000.0 kWh |
| 5008 | ATP tokens not found | No active ATP tokens for operation | Create energy operation first |
| 5009 | Tensor not found | Trust tensor does not exist | Create relationship tensor first |
| 5010 | Operation not completed | Energy operation not in completed state | Execute energy transfer first |

### Error Response Format

```json
{
  "error": {
    "code": 5002,
    "message": "Pairing not authorized",
    "details": "Component MODBATT-MOD-RC001-001 cannot pair with MODBATT-HOST-RC001. Direct module-host pairing forbidden.",
    "suggestion": "Pair module with pack controller first"
  }
}
```

---

## Rate Limits

### Default Limits (per IP address)

| Endpoint Category | Rate Limit | Burst Limit | Window |
|------------------|------------|-------------|--------|
| Component Queries | 200 req/min | 20 req | 1 minute |
| General API | 100 req/min | 10 req | 1 minute |
| Transaction Broadcasting | 10 req/min | 5 req | 1 minute |
| WebSocket Connections | 5 conn/IP | 2 conn | 1 minute |
| Trust Calculations | 50 req/min | 10 req | 1 minute |
| Energy Operations | 20 req/min | 5 req | 1 minute |

### Rate Limit Headers

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1704067200
X-RateLimit-Retry-After: 60
```

### Rate Limit Exceeded Response

```json
{
  "error": {
    "code": 429,
    "message": "Rate limit exceeded",
    "details": "Too many requests. Limit: 100 req/min",
    "retry_after": 60
  }
}
```

---

## API Versioning

### Current Version
- **API Version**: v1
- **Chain Version**: racecarweb-1.0
- **Compatibility**: Cosmos SDK v0.53.0

### Version Header
```http
Accept: application/json; version=1
```

### Deprecation Policy
- 90 days notice for breaking changes
- 180 days support for deprecated endpoints
- Version-specific documentation maintained

---

## SDK Integration

### Python SDK
```python
from web4_client import Web4Client

client = Web4Client(
    rest_endpoint="https://your-domain.com/api",
    chain_id="racecarweb"
)

# Register component
result = client.register_component(
    "MODBATT-MOD-RC001-001", 
    "module", 
    "ModBatt Pro Series"
)

# Create LCT relationship
lct = client.create_lct_relationship(
    "MODBATT-MOD-RC001-001",
    "MODBATT-PACK-RC001-A"
)
```

### JavaScript/TypeScript SDK
```javascript
import { Web4Client } from '@web4/client';

const client = new Web4Client({
  restEndpoint: 'https://your-domain.com/api',
  chainId: 'racecarweb'
});

// Subscribe to events
client.onLctCreated((event) => {
  console.log('New LCT:', event.lct_id);
});
```

---

## Testing

### API Testing Tools

**cURL Examples**:
```bash
# Health check
curl https://your-domain.com/api/cosmos/base/tendermint/v1beta1/node_info

# Get component
curl https://your-domain.com/api/racecarweb/componentregistry/v1/component/MODBATT-MOD-RC001-001

# Check pairing auth
curl "https://your-domain.com/api/racecarweb/componentregistry/v1/pairing-auth?component_a=MODBATT-MOD-RC001-001&component_b=MODBATT-PACK-RC001-A"
```

**Postman Collection**: Available in repository `/docs/postman/`

**Integration Tests**: Available in repository `/tests/integration/`

---

## Support

### Documentation Links
- **GitHub Repository**: https://github.com/your-org/web4-racecar-demo
- **Cosmos SDK Docs**: https://docs.cosmos.network/
- **Web4 Architecture**: See project documentation

### Community Support
- **Discord**: Web4 Development Community
- **GitHub Issues**: Bug reports and feature requests
- **Stack Overflow**: Tag questions with `web4` and `cosmos-sdk`

---

**Document Version**: 1.0  
**Last Updated**: December 2024  
**Compatibility**: Cosmos SDK v0.53.0, Ignite CLI v29  
**Format**: Confluence/Jira Compatible Markdown