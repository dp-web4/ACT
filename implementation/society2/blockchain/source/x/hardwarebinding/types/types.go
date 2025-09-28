package types

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os/exec"
    "strings"
)

const (
    // ModuleName defines the module name
    ModuleName = "hardwarebinding"

    // StoreKey defines the primary module store key
    StoreKey = ModuleName

    // RouterKey defines the module's message routing key
    RouterKey = ModuleName

    // MemStoreKey defines the in-memory store key
    MemStoreKey = "mem_hardware"
)

// HardwareBinding represents the binding between chain and hardware
type HardwareBinding struct {
    Platform     string            `json:"platform"`
    HardwareHash string            `json:"hardware_hash"`
    Components   HardwareComponents `json:"components"`
    Timestamp    int64             `json:"timestamp"`
    BlockHeight  int64             `json:"block_height"`
}

// HardwareComponents contains individual hardware identifiers
type HardwareComponents struct {
    WindowsUUID string `json:"windows_uuid"`
    HyperVUUID  string `json:"hyperv_uuid"`
    WSLBootID   string `json:"wsl_boot_id"`
    CPUInfo     string `json:"cpu_info"`
    MemoryKB    int64  `json:"memory_kb"`
}

// ExtractCurrentHardware gets current machine's hardware identifiers
func ExtractCurrentHardware() (*HardwareBinding, error) {
    // Try to execute the hardware extraction script
    cmd := exec.Command("/bin/bash", "-c", `
        extract_windows_uuid() {
            powershell.exe -Command "Get-WmiObject -Class Win32_ComputerSystemProduct | Select-Object -ExpandProperty UUID" 2>/dev/null | tr -d '\r\n'
        }

        extract_hyperv_uuid() {
            cat /sys/class/dmi/id/product_uuid 2>/dev/null || echo "HYPERV_UNAVAILABLE"
        }

        extract_wsl_boot_id() {
            cat /proc/sys/kernel/random/boot_id 2>/dev/null || echo "BOOT_ID_UNAVAILABLE"
        }

        extract_cpu_info() {
            grep "model name" /proc/cpuinfo | head -1 | cut -d: -f2 | xargs
        }

        extract_memory_size() {
            grep MemTotal /proc/meminfo | awk '{print $2}'
        }

        windows_uuid=$(extract_windows_uuid)
        hyperv_uuid=$(extract_hyperv_uuid)
        wsl_boot_id=$(extract_wsl_boot_id)
        cpu_info=$(extract_cpu_info)
        memory_size=$(extract_memory_size)

        echo "$windows_uuid|$hyperv_uuid|$wsl_boot_id|$cpu_info|$memory_size"
    `)

    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to extract hardware: %w", err)
    }

    parts := strings.Split(strings.TrimSpace(string(output)), "|")
    if len(parts) != 5 {
        return nil, fmt.Errorf("unexpected hardware output format")
    }

    var memKB int64
    fmt.Sscanf(parts[4], "%d", &memKB)

    components := HardwareComponents{
        WindowsUUID: parts[0],
        HyperVUUID:  parts[1],
        WSLBootID:   parts[2],
        CPUInfo:     parts[3],
        MemoryKB:    memKB,
    }

    // Calculate hardware hash
    composite := fmt.Sprintf("%s|%s|%s|%s|%d",
        components.WindowsUUID,
        components.HyperVUUID,
        components.WSLBootID,
        components.CPUInfo,
        components.MemoryKB,
    )

    hash := sha256.Sum256([]byte(composite))

    return &HardwareBinding{
        Platform:     "wsl2",
        HardwareHash: hex.EncodeToString(hash[:]),
        Components:   components,
    }, nil
}

// VerifyHardware checks if current hardware matches the binding
func (hb *HardwareBinding) VerifyHardware() error {
    current, err := ExtractCurrentHardware()
    if err != nil {
        return fmt.Errorf("failed to get current hardware: %w", err)
    }

    // Allow boot ID to change (WSL restarts) but other components must match
    if current.Components.WindowsUUID != hb.Components.WindowsUUID {
        return fmt.Errorf("Windows UUID mismatch: expected %s, got %s",
            hb.Components.WindowsUUID, current.Components.WindowsUUID)
    }

    if current.Components.HyperVUUID != hb.Components.HyperVUUID &&
       hb.Components.HyperVUUID != "HYPERV_UNAVAILABLE" {
        return fmt.Errorf("Hyper-V UUID mismatch: expected %s, got %s",
            hb.Components.HyperVUUID, current.Components.HyperVUUID)
    }

    if current.Components.CPUInfo != hb.Components.CPUInfo {
        return fmt.Errorf("CPU info mismatch: expected %s, got %s",
            hb.Components.CPUInfo, current.Components.CPUInfo)
    }

    return nil
}