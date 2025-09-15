# Event System Implementation

## Overview

The API Bridge includes an optional event emission system that can notify external systems (like SQL databases) when blockchain transactions occur. This system is designed to be:

- **Flexible**: Configurable webhook endpoints per event type
- **Reliable**: Built-in retry logic with exponential backoff
- **Performant**: In-memory queue with configurable size
- **Safe**: Disabled by default, only active when explicitly configured

## Architecture

```
┌─────────────────┐    HTTP Request    ┌─────────────────┐
│   API Bridge    │ ─────────────────► │   Webhook       │
│                 │                    │   Receiver      │
│ ┌─────────────┐ │                    │                 │
│ │ Event Queue │ │                    │ ┌─────────────┐ │
│ │             │ │                    │ │ SQL Database│ │
│ │ - Retry     │ │                    │ │             │ │
│ │ - Backoff   │ │                    │ │ - Audit Log │ │
│ │ - Multiple  │ │                    │ │ - Analytics │ │
│ │   Endpoints │ │                    │ │ - Reports   │ │
│ └─────────────┘ │                    │ └─────────────┘ │
└─────────────────┘                    └─────────────────┘
```

## Implementation Details

### Core Components

#### 1. Event Queue (`internal/events/events.go`)
- **In-memory queue** with configurable size
- **Worker goroutine** for processing events
- **Retry logic** with exponential backoff
- **Multiple endpoints** per event type
- **Graceful shutdown** support

#### 2. Configuration (`internal/config/config.go`)
- **EventsConfig** struct for settings
- **Default values** (disabled by default)
- **Endpoint mapping** per event type
- **Retry and queue settings**

#### 3. Handler Integration (`internal/handlers/handlers.go`)
- **Event emission** after successful blockchain transactions
- **Conditional emission** (only if enabled)
- **Rich event data** including transaction hashes
- **Graceful shutdown** integration

### Event Types

| Event Type | Triggered When | Data Fields |
|------------|----------------|-------------|
| `component_registered` | Component registration succeeds | `component_id`, `creator`, `component_data`, `context`, `timestamp`, `tx_hash` |
| `component_verified` | Component verification succeeds | `component_id`, `verifier`, `context`, `timestamp`, `tx_hash` |
| `pairing_initiated` | Pairing initiation succeeds | `challenge_id`, `creator`, `component_a`, `component_b`, `operational_context`, `timestamp`, `tx_hash` |
| `pairing_completed` | Pairing completion succeeds | `challenge_id`, `creator`, `session_context`, `lct_id`, `timestamp`, `tx_hash` |
| `lct_created` | LCT creation succeeds | `lct_id`, `creator`, `component_a`, `component_b`, `context`, `timestamp`, `tx_hash` |
| `trust_tensor_created` | Trust tensor creation succeeds | `tensor_id`, `creator`, `component_a`, `component_b`, `context`, `initial_score`, `timestamp`, `tx_hash` |
| `energy_transfer` | Energy transfer succeeds | `operation_id`, `creator`, `amount`, `context`, `timestamp`, `tx_hash` |

### Event Data Structure

```json
{
  "event_type": "component_registered",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "component_id": "COMP-alice-1705312200",
    "creator": "alice",
    "component_data": "battery_pack_v1",
    "context": "test",
    "timestamp": 1705312200,
    "tx_hash": "ABC123DEF456..."
  }
}
```

## Configuration

### Basic Configuration

```yaml
events:
  enabled: true
  max_retries: 3
  retry_delay: 5  # seconds
  queue_size: 1000
  endpoints:
    component_registered:
      - "http://localhost:3000/webhooks/component-registered"
    energy_transfer:
      - "http://sql-audit-service:8080/events"
```

### Advanced Configuration

```yaml
events:
  enabled: true
  max_retries: 5
  retry_delay: 2
  queue_size: 5000
  endpoints:
    component_registered:
      - "http://audit-service:8080/events"
      - "http://monitoring-service:8080/webhooks"
      - "http://sql-database:8080/api/events"
    component_verified:
      - "http://audit-service:8080/events"
    pairing_initiated:
      - "http://audit-service:8080/events"
      - "http://real-time-dashboard:8080/events"
    pairing_completed:
      - "http://audit-service:8080/events"
      - "http://analytics-service:8080/events"
    lct_created:
      - "http://audit-service:8080/events"
    trust_tensor_created:
      - "http://audit-service:8080/events"
    energy_transfer:
      - "http://audit-service:8080/events"
      - "http://energy-monitoring:8080/events"
      - "http://billing-service:8080/events"
```

## Usage Examples

### 1. SQL Database Integration

Set up a webhook receiver that stores events in a SQL database:

```python
# webhook_receiver.py
from flask import Flask, request
import sqlite3
import json

app = Flask(__name__)

@app.route('/events', methods=['POST'])
def receive_event():
    event = request.json
    
    # Store in SQL database
    conn = sqlite3.connect('audit.db')
    cursor = conn.cursor()
    
    cursor.execute('''
        INSERT INTO events (event_type, timestamp, data, tx_hash)
        VALUES (?, ?, ?, ?)
    ''', (
        event['event_type'],
        event['timestamp'],
        json.dumps(event['data']),
        event['data'].get('tx_hash', '')
    ))
    
    conn.commit()
    conn.close()
    
    return {'status': 'received'}, 200

if __name__ == '__main__':
    app.run(port=3000)
```

### 2. Real-time Monitoring

Create a webhook receiver for real-time monitoring:

```javascript
// monitoring_receiver.js
const express = require('express');
const app = express();

app.use(express.json());

app.post('/webhooks/component-registered', (req, res) => {
    const event = req.body;
    
    // Send to monitoring dashboard
    console.log(`🔔 Component registered: ${event.data.component_id}`);
    console.log(`   Creator: ${event.data.creator}`);
    console.log(`   TX Hash: ${event.data.tx_hash}`);
    
    // Could also send to WebSocket clients, Slack, etc.
    
    res.json({ status: 'received' });
});

app.listen(3000, () => {
    console.log('Monitoring receiver running on port 3000');
});
```

### 3. Audit Logging

Set up comprehensive audit logging:

```yaml
# config.yaml
events:
  enabled: true
  max_retries: 3
  retry_delay: 5
  queue_size: 1000
  endpoints:
    component_registered:
      - "http://audit-service:8080/events"
    component_verified:
      - "http://audit-service:8080/events"
    pairing_initiated:
      - "http://audit-service:8080/events"
    pairing_completed:
      - "http://audit-service:8080/events"
    lct_created:
      - "http://audit-service:8080/events"
    trust_tensor_created:
      - "http://audit-service:8080/events"
    energy_transfer:
      - "http://audit-service:8080/events"
```

## Testing

### Test Scripts

**Windows:**
```cmd
test_event_system.bat
```

**Linux/macOS:**
```bash
./test_event_system.sh
```

### Manual Testing

1. **Enable events** in configuration
2. **Start webhook receiver** (e.g., using netcat or a simple HTTP server)
3. **Make API calls** that trigger events
4. **Check webhook receiver** for received events

### Example Test Commands

```bash
# Start simple webhook receiver
python -m http.server 3000

# Test component registration
curl -X POST http://localhost:8080/api/v1/components/register \
  -H "Content-Type: application/json" \
  -d '{"creator": "test-user", "component_data": "test-battery", "context": "test"}'

# Check webhook receiver logs for event
```

## Security Considerations

### Current Implementation (Demo)
- **No authentication** on webhook endpoints
- **No encryption** of event data
- **No rate limiting** on webhook calls
- **No validation** of webhook responses

### Production Recommendations
- **HTTPS only** for webhook endpoints
- **Authentication** (API keys, JWT tokens)
- **Event signing** for integrity verification
- **Rate limiting** to prevent abuse
- **Response validation** and monitoring
- **Dead letter queues** for failed events

## Performance Considerations

### Current Implementation
- **In-memory queue** (events lost on restart)
- **Synchronous webhook calls** (can block processing)
- **No batching** of events
- **No compression** of event data

### Production Enhancements
- **Persistent queue** (Redis, RabbitMQ)
- **Asynchronous processing** with worker pools
- **Event batching** for high-volume scenarios
- **Compression** for large event payloads
- **Monitoring and metrics** for queue health

## Future Enhancements

### Planned Features
- **Redis queue backend** for persistence
- **Event filtering** and transformation
- **Event schema validation**
- **Webhook authentication** support
- **Event replay** capabilities
- **Metrics and monitoring** dashboard

### Integration Possibilities
- **Apache Kafka** for high-throughput event streaming
- **AWS SQS/SNS** for cloud-native event processing
- **Elasticsearch** for event search and analytics
- **Grafana** for event visualization
- **Prometheus** for event metrics

## Troubleshooting

### Common Issues

#### Events Not Being Sent
1. **Check configuration**: Ensure `events.enabled: true`
2. **Verify endpoints**: Check webhook URLs are accessible
3. **Check logs**: Look for event queue errors
4. **Test connectivity**: Verify webhook endpoints respond

#### Webhook Failures
1. **Check retry settings**: Increase `max_retries` if needed
2. **Verify endpoint health**: Ensure webhook receivers are running
3. **Check network**: Verify connectivity to webhook endpoints
4. **Review logs**: Look for specific error messages

#### Performance Issues
1. **Increase queue size**: Adjust `queue_size` setting
2. **Optimize webhook receivers**: Ensure they respond quickly
3. **Add monitoring**: Track event processing metrics
4. **Consider batching**: Group multiple events together

### Debug Commands

```bash
# Check event queue status
curl http://localhost:8080/health

# Test webhook endpoint
curl -X POST http://localhost:3000/webhooks/test \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'

# Monitor API bridge logs
./api-bridge.exe --log-level debug
```

## Conclusion

The event system provides a flexible foundation for integrating the API Bridge with external systems. While the current implementation is suitable for demos and development, production deployments should consider the security and performance enhancements outlined above.

The system is designed to be:
- **Non-intrusive**: Disabled by default, doesn't affect core functionality
- **Extensible**: Easy to add new event types and endpoints
- **Reliable**: Built-in retry logic handles temporary failures
- **Observable**: Comprehensive logging for debugging and monitoring 