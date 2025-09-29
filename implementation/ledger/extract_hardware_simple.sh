#!/bin/bash
# Simple Hardware Identity Extraction for Sprout (no jq dependency)

set -e

echo "🔐 Extracting Jetson Hardware Identity"
echo "================================================"

# Get hardware values
DEVICE_SERIAL=$(cat /proc/device-tree/serial-number 2>/dev/null | tr -d '\0' | tr -d ' ' || echo "unknown")
SOC_ID=$(cat /sys/devices/soc0/soc_id 2>/dev/null || echo "unknown")
SOC_FAMILY=$(cat /sys/devices/soc0/family 2>/dev/null || echo "unknown")
SOC_REVISION=$(cat /sys/devices/soc0/revision 2>/dev/null || echo "unknown")
BOOT_ID=$(cat /proc/sys/kernel/random/boot_id 2>/dev/null || echo "unknown")
MACHINE_ID=$(cat /etc/machine-id 2>/dev/null || echo "unknown")
MAC=$(ip link show | grep -E "link/ether" | head -1 | awk '{print $2}' || echo "unknown")
CPU_MODEL=$(lscpu | grep "Model name" | cut -d':' -f2 | xargs || echo "unknown")
CORES=$(nproc)
MEM_KB=$(grep MemTotal /proc/meminfo | awk '{print $2}')
MEM_GB=$((MEM_KB / 1024 / 1024))

echo "✓ Device Serial: $DEVICE_SERIAL"
echo "✓ SoC: $SOC_FAMILY ID=$SOC_ID Rev=$SOC_REVISION"
echo "✓ Boot Session: ${BOOT_ID:0:8}..."
echo "✓ Machine ID: ${MACHINE_ID:0:8}..."
echo "✓ Primary MAC: $MAC"
echo "✓ CPU: $CPU_MODEL ($CORES cores)"
echo "✓ Memory: ${MEM_GB}GB"

# Create canonical string for hashing (exclude boot_id)
CANONICAL="device_serial=$DEVICE_SERIAL,soc_id=$SOC_ID,soc_family=$SOC_FAMILY,soc_revision=$SOC_REVISION,machine_id=$MACHINE_ID,mac=$MAC,cpu_model=$CPU_MODEL,cores=$CORES,memory_kb=$MEM_KB"
HARDWARE_HASH=$(echo -n "$CANONICAL" | sha256sum | cut -d' ' -f1)

echo ""
echo "🔑 Hardware Hash: $HARDWARE_HASH"

# Create JSON manually
OUTPUT_FILE="/home/sprout/ai-workspace/ACT/implementation/ledger/hardware_identity.json"
TIMESTAMP=$(date -Iseconds)

cat > "$OUTPUT_FILE" << EOF
{
  "platform": "jetson_orin_nano",
  "tier": 2,
  "device_serial": "$DEVICE_SERIAL",
  "soc": {
    "id": "$SOC_ID",
    "family": "$SOC_FAMILY",
    "revision": "$SOC_REVISION"
  },
  "boot_id": "$BOOT_ID",
  "machine_id": "$MACHINE_ID",
  "primary_mac": "$MAC",
  "cpu": {
    "model": "$CPU_MODEL",
    "cores": $CORES
  },
  "memory_kb": $MEM_KB,
  "hardware_hash": "$HARDWARE_HASH",
  "timestamp": "$TIMESTAMP"
}
EOF

echo ""
echo "✅ Hardware identity extracted successfully!"
echo "📁 Saved to: $OUTPUT_FILE"
echo ""
echo "Key Identifiers:"
echo "  Device Serial: $DEVICE_SERIAL"
echo "  Hardware Hash: ${HARDWARE_HASH:0:16}..."
echo "  Platform Tier: 2 (Silicon Identifiers)"
echo ""
echo "Note: This hardware binding provides strong identity for Sprout society."
echo "The Jetson's device serial ($DEVICE_SERIAL) is unique per module."