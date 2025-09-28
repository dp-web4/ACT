#!/bin/bash

# Extract Hardware Fingerprint for CBP Machine
# This creates a unique hardware identity for this specific WSL2 instance

MACHINE_ID_FILE="/mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain/keys/machine_id.json"
HARDWARE_INFO_FILE="/mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain/keys/hardware_info.txt"

echo "Extracting hardware fingerprint for CBP machine..."

# Get WSL2 instance UUID (unique per WSL installation)
WSL_UUID=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "cbp-default-$(date +%s)")

# Get system information
HOSTNAME=$(hostname)
KERNEL=$(uname -r)
ARCH=$(uname -m)

# Get CPU info (consistent in WSL2)
CPU_INFO=$(grep -m1 "model name" /proc/cpuinfo | cut -d: -f2 | xargs)
CPU_CORES=$(nproc)

# Get memory info
MEM_TOTAL=$(grep MemTotal /proc/meminfo | awk '{print $2}')

# Get filesystem UUID (for data persistence)
FS_UUID=$(findmnt -n -o UUID / 2>/dev/null || echo "wsl2-root-fs")

# Create composite hardware ID
HARDWARE_ID=$(echo -n "${HOSTNAME}-${WSL_UUID}-${CPU_CORES}-${MEM_TOTAL}" | sha256sum | cut -d' ' -f1)

# Generate deterministic keypair seed from hardware
SEED=$(echo -n "${HARDWARE_ID}-${USER}-cbp-chain" | sha256sum | cut -d' ' -f1)

# Save hardware info
cat > "$HARDWARE_INFO_FILE" << EOF
CBP Machine Hardware Information
================================
Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
Hostname: ${HOSTNAME}
WSL UUID: ${WSL_UUID}
Kernel: ${KERNEL}
Architecture: ${ARCH}
CPU: ${CPU_INFO}
CPU Cores: ${CPU_CORES}
Memory: ${MEM_TOTAL} KB
Filesystem UUID: ${FS_UUID}
Hardware ID: ${HARDWARE_ID}
Seed: ${SEED}
EOF

# Create machine identity JSON
cat > "$MACHINE_ID_FILE" << EOF
{
  "machine_id": "${HARDWARE_ID}",
  "hostname": "${HOSTNAME}",
  "wsl_uuid": "${WSL_UUID}",
  "created": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "hardware": {
    "cpu": "${CPU_INFO}",
    "cores": ${CPU_CORES},
    "memory_kb": ${MEM_TOTAL},
    "kernel": "${KERNEL}",
    "arch": "${ARCH}"
  },
  "chain_config": {
    "chain_id": "cbp-chain-${HARDWARE_ID:0:8}",
    "moniker": "cbp-validator",
    "home": "/mnt/c/exe/projects/ai-agents/ACT/implementation/cbp-chain"
  },
  "seed": "${SEED}"
}
EOF

echo "Hardware fingerprint extracted successfully!"
echo "Hardware ID: ${HARDWARE_ID}"
echo "Chain ID: cbp-chain-${HARDWARE_ID:0:8}"
echo ""
echo "Files created:"
echo "  - ${MACHINE_ID_FILE}"
echo "  - ${HARDWARE_INFO_FILE}"