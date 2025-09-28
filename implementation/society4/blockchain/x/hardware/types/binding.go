package types

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

// HardwareBinding represents the binding between self-LCT and hardware
type HardwareBinding struct {
    Platform     string            `json:"platform"`
    HardwareHash string            `json:"hardware_hash"`
    Components   HardwareComponents `json:"components"`
    Timestamp    int64             `json:"timestamp"`
    Signature    []byte            `json:"signature"`
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
    cmd := exec.Command("/bin/bash", "hardware/extract_hardware.sh", "json")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to extract hardware: %w", err)
    }

    var result struct {
        HardwareBinding HardwareBinding `json:"hardware_binding"`
    }

    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse hardware info: %w", err)
    }

    return &result.HardwareBinding, nil
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

    if current.Components.HyperVUUID != hb.Components.HyperVUUID {
        return fmt.Errorf("Hyper-V UUID mismatch: expected %s, got %s",
            hb.Components.HyperVUUID, current.Components.HyperVUUID)
    }

    if current.Components.CPUInfo != hb.Components.CPUInfo {
        return fmt.Errorf("CPU info mismatch: expected %s, got %s",
            hb.Components.CPUInfo, current.Components.CPUInfo)
    }

    // Recalculate hash with current boot ID
    currentHash := hb.CalculateHash(current.Components.WSLBootID)
    if currentHash != current.HardwareHash {
        return fmt.Errorf("hardware hash mismatch")
    }

    return nil
}

// CalculateHash generates hardware hash with optional boot ID override
func (hb *HardwareBinding) CalculateHash(bootID string) string {
    if bootID == "" {
        bootID = hb.Components.WSLBootID
    }

    composite := fmt.Sprintf("%s|%s|%s|%s|%d",
        hb.Components.WindowsUUID,
        hb.Components.HyperVUUID,
        bootID,
        hb.Components.CPUInfo,
        hb.Components.MemoryKB,
    )

    hash := sha256.Sum256([]byte(composite))
    return hex.EncodeToString(hash[:])
}

// Hash returns the canonical hash of the binding
func (hb *HardwareBinding) Hash() []byte {
    h := sha256.New()
    h.Write([]byte(hb.Platform))
    h.Write([]byte(hb.HardwareHash))
    h.Write([]byte(hb.Components.WindowsUUID))
    h.Write([]byte(hb.Components.HyperVUUID))
    h.Write([]byte(hb.Components.CPUInfo))
    h.Write([]byte(fmt.Sprintf("%d", hb.Components.MemoryKB)))
    return h.Sum(nil)
}

// ValidateBinding ensures the hardware binding is valid
func (hb *HardwareBinding) ValidateBinding() error {
    if hb.Platform != "wsl2" {
        return fmt.Errorf("unsupported platform: %s", hb.Platform)
    }

    if hb.HardwareHash == "" {
        return fmt.Errorf("hardware hash is empty")
    }

    if hb.Components.WindowsUUID == "" {
        return fmt.Errorf("Windows UUID is required")
    }

    if hb.Components.HyperVUUID == "" && hb.Components.HyperVUUID != "HYPERV_UNAVAILABLE" {
        return fmt.Errorf("Hyper-V UUID is required")
    }

    return nil
}