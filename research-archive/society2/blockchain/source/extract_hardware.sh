#!/bin/bash

# Extract hardware identifiers for WSL2 environment
# These create a unique fingerprint for this specific WSL2 instance

extract_windows_uuid() {
    # Get Windows machine UUID from registry via PowerShell
    powershell.exe -Command "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID" 2>/dev/null | tr -d '\r\n'
}

extract_hyperv_uuid() {
    # Get Hyper-V VM UUID (WSL2 runs in Hyper-V)
    cat /sys/class/dmi/id/product_uuid 2>/dev/null || echo "HYPERV_UNAVAILABLE"
}

extract_wsl_boot_id() {
    # Get WSL boot ID (changes on WSL restart but stable during session)
    cat /proc/sys/kernel/random/boot_id 2>/dev/null || echo "BOOT_ID_UNAVAILABLE"
}

extract_cpu_info() {
    # Get CPU model for additional verification
    grep "model name" /proc/cpuinfo | head -1 | cut -d: -f2 | xargs
}

extract_memory_size() {
    # Get total memory allocated to WSL2
    grep MemTotal /proc/meminfo | awk '{print $2}'
}

# Generate composite hardware fingerprint
generate_hardware_hash() {
    local windows_uuid=$(extract_windows_uuid)
    local hyperv_uuid=$(extract_hyperv_uuid)
    local wsl_boot_id=$(extract_wsl_boot_id)
    local cpu_info=$(extract_cpu_info)
    local memory_size=$(extract_memory_size)

    # Combine all identifiers
    local composite="${windows_uuid}|${hyperv_uuid}|${wsl_boot_id}|${cpu_info}|${memory_size}"

    # Generate SHA256 hash
    echo -n "$composite" | sha256sum | cut -d' ' -f1
}

# Output JSON format
output_json() {
    local windows_uuid=$(extract_windows_uuid)
    local hyperv_uuid=$(extract_hyperv_uuid)
    local wsl_boot_id=$(extract_wsl_boot_id)
    local cpu_info=$(extract_cpu_info)
    local memory_size=$(extract_memory_size)
    local hardware_hash=$(generate_hardware_hash)

    cat <<EOF
{
    "hardware_binding": {
        "platform": "wsl2",
        "hardware_hash": "${hardware_hash}",
        "components": {
            "windows_uuid": "${windows_uuid}",
            "hyperv_uuid": "${hyperv_uuid}",
            "wsl_boot_id": "${wsl_boot_id}",
            "cpu_info": "${cpu_info}",
            "memory_kb": ${memory_size}
        },
        "timestamp": $(date +%s)
    }
}
EOF
}

# Main execution
case "${1:-}" in
    json)
        output_json
        ;;
    hash)
        generate_hardware_hash
        ;;
    *)
        echo "Usage: $0 {json|hash}"
        echo "  json - Output full hardware binding information"
        echo "  hash - Output only the hardware hash"
        exit 1
        ;;
esac