# Sonic JSON Library Excluded

## Security Decision

The `bytedance/sonic` library has been permanently excluded from this build because:

1. **Uses unsafe pointers** - Direct memory manipulation bypasses Go's safety guarantees
2. **Assembly optimizations** - Platform-specific code that breaks with Go runtime changes
3. **Version brittleness** - Breaks with each Go version update
4. **Unnecessary for blockchain** - Standard library JSON is more than fast enough for consensus

## Implementation

We use standard library `encoding/json` instead, which:
- Is memory-safe
- Works with all Go versions
- Is maintained by Go team
- Has predictable performance
- Is sufficient for blockchain consensus with thousands of nodes

## Performance Impact

Negligible. Blockchain consensus is limited by:
- Network latency (milliseconds)
- Disk I/O (milliseconds)
- Cryptographic operations (microseconds)

JSON parsing (nanoseconds) is not the bottleneck.

## To Ensure Sonic Exclusion

Add to go.mod:
```go
replace (
    github.com/bytedance/sonic => encoding/json v0.0.0
    github.com/bytedance/sonic/loader => encoding/json v0.0.0
)
```

Or use build tags:
```bash
go build -tags=nosonic
```

## Bottom Line

**Safety > Speed** for blockchain infrastructure.

Unsafe optimizations in consensus-critical code are an unacceptable risk.