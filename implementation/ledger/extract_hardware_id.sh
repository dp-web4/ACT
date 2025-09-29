#!/bin/bash
# Hardware Identity Extraction for Sprout (Jetson Orin Nano)
# Generates deterministic hardware fingerprint for blockchain binding

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔐 Extracting Jetson Hardware Identity${NC}"
echo "================================================"

# Initialize hardware data object
HARDWARE_JSON='{"platform": "jetson_orin_nano", "tier": 2}'

# 1. Device Tree Serial Number (unique per Jetson module)
if [ -f /proc/device-tree/serial-number ]; then
    SERIAL=$(cat /proc/device-tree/serial-number | tr -d '\0' | tr -d ' ')
    echo -e "${GREEN}✓${NC} Device Serial: ${SERIAL}"
    HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {device_serial: \"$SERIAL\"}")
else
    echo -e "${YELLOW}⚠${NC} Device serial not found"
fi

# 2. SoC Information
if [ -d /sys/devices/soc0 ]; then
    SOC_ID=$(cat /sys/devices/soc0/soc_id 2>/dev/null || echo "unknown")
    SOC_FAMILY=$(cat /sys/devices/soc0/family 2>/dev/null || echo "unknown")
    SOC_REVISION=$(cat /sys/devices/soc0/revision 2>/dev/null || echo "unknown")
    echo -e "${GREEN}✓${NC} SoC: $SOC_FAMILY ID=$SOC_ID Rev=$SOC_REVISION"
    HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {soc: {id: \"$SOC_ID\", family: \"$SOC_FAMILY\", revision: \"$SOC_REVISION\"}}")
fi

# 3. Boot ID (changes each boot, good for session binding)
if [ -f /proc/sys/kernel/random/boot_id ]; then
    BOOT_ID=$(cat /proc/sys/kernel/random/boot_id)
    echo -e "${GREEN}✓${NC} Boot Session: ${BOOT_ID:0:8}..."
    HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {boot_id: \"$BOOT_ID\"}")
fi

# 4. Machine ID (persistent across reboots)
if [ -f /etc/machine-id ]; then
    MACHINE_ID=$(cat /etc/machine-id)
    echo -e "${GREEN}✓${NC} Machine ID: ${MACHINE_ID:0:8}..."
    HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {machine_id: \"$MACHINE_ID\"}")
fi

# 5. Network Interface MACs (for additional entropy)
MACS=$(ip link show | grep -E "link/ether" | awk '{print $2}' | tr '\n' ',' | sed 's/,$//')
if [ -n "$MACS" ]; then
    echo -e "${GREEN}✓${NC} Network MACs: $(echo $MACS | cut -d',' -f1)"
    HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {network_macs: \"$MACS\"}")
fi

# 6. CPU Information
CPU_INFO=$(lscpu | grep "Model name" | cut -d':' -f2 | xargs)
CORES=$(nproc)
echo -e "${GREEN}✓${NC} CPU: $CPU_INFO ($CORES cores)"
HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {cpu: {model: \"$CPU_INFO\", cores: $CORES}}")

# 7. Memory Information
MEM_TOTAL=$(grep MemTotal /proc/meminfo | awk '{print $2}')
echo -e "${GREEN}✓${NC} Memory: $((MEM_TOTAL / 1024 / 1024))GB"
HARDWARE_JSON=$(echo "$HARDWARE_JSON" | jq ". + {memory_kb: $MEM_TOTAL}")

# Generate deterministic hardware hash
echo -e "\n${YELLOW}🔄 Generating Hardware Binding Hash...${NC}"

# Create canonical string for hashing (exclude boot_id for persistence)
CANONICAL=$(echo "$HARDWARE_JSON" | jq -S 'del(.boot_id)' | tr -d ' \n')
HARDWARE_HASH=$(echo -n "$CANONICAL" | sha256sum | cut -d' ' -f1)

echo -e "${GREEN}✓${NC} Hardware Hash: ${HARDWARE_HASH}"

# Add hash to JSON
FINAL_JSON=$(echo "$HARDWARE_JSON" | jq ". + {hardware_hash: \"$HARDWARE_HASH\", timestamp: \"$(date -Iseconds)\"}")

# Output to file
OUTPUT_FILE="/home/sprout/ai-workspace/ACT/implementation/ledger/hardware_identity.json"
echo "$FINAL_JSON" | jq '.' > "$OUTPUT_FILE"

echo -e "\n${GREEN}✅ Hardware identity extracted successfully!${NC}"
echo -e "📁 Saved to: $OUTPUT_FILE"
echo -e "\n${YELLOW}Key Identifiers:${NC}"
echo -e "  Device Serial: ${SERIAL}"
echo -e "  Hardware Hash: ${HARDWARE_HASH:0:16}..."
echo -e "  Platform Tier: 2 (Silicon Identifiers)"
echo -e "\n${YELLOW}Note:${NC} This hardware binding is suitable for Society identity."
echo "The Jetson's device serial provides strong hardware uniqueness."