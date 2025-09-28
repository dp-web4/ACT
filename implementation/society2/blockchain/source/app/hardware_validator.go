package app

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "cosmossdk.io/log"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// HardwareValidator handles hardware binding validation
type HardwareValidator struct {
    logger            log.Logger
    genesisBinding    *HardwareBinding
    validationEnabled bool
    lastCheckHeight   int64
}

// HardwareBinding represents the hardware configuration
type HardwareBinding struct {
    Platform     string            `json:"platform"`
    HardwareHash string            `json:"hardware_hash"`
    Components   HardwareComponents `json:"components"`
    Timestamp    int64             `json:"timestamp"`
}

// HardwareComponents contains hardware identifiers
type HardwareComponents struct {
    WindowsUUID string `json:"windows_uuid"`
    HyperVUUID  string `json:"hyperv_uuid"`
    WSLBootID   string `json:"wsl_boot_id"`
    CPUInfo     string `json:"cpu_info"`
    MemoryKB    int64  `json:"memory_kb"`
}

// NewHardwareValidator creates a new hardware validator
func NewHardwareValidator(logger log.Logger) *HardwareValidator {
    return &HardwareValidator{
        logger:            logger,
        validationEnabled: false, // Disabled by default for now
    }
}

// LoadGenesisBinding loads the hardware binding from genesis
func (hv *HardwareValidator) LoadGenesisBinding() error {
    // Try to load from hardware_binding.json file
    bindingPath := os.ExpandEnv("$HOME/.society4chain/hardware_binding.json")
    data, err := os.ReadFile(bindingPath)
    if err != nil {
        // If file doesn't exist, extract current hardware as genesis
        hv.logger.Info("No genesis hardware binding found, extracting current hardware")
        binding, err := ExtractCurrentHardware()
        if err != nil {
            return fmt.Errorf("failed to extract hardware: %w", err)
        }
        hv.genesisBinding = binding

        // Save it for future reference
        if jsonData, err := json.MarshalIndent(map[string]interface{}{
            "hardware_binding": binding,
        }, "", "    "); err == nil {
            os.WriteFile(bindingPath, jsonData, 0644)
        }

        return nil
    }

    // Parse the JSON
    var wrapper struct {
        HardwareBinding HardwareBinding `json:"hardware_binding"`
    }
    if err := json.Unmarshal(data, &wrapper); err != nil {
        return fmt.Errorf("failed to parse hardware binding: %w", err)
    }

    hv.genesisBinding = &wrapper.HardwareBinding
    hv.logger.Info("Loaded genesis hardware binding",
        "hardware_hash", hv.genesisBinding.HardwareHash[:16]+"...",
        "platform", hv.genesisBinding.Platform)

    return nil
}

// ValidateHardware checks if current hardware matches genesis
func (hv *HardwareValidator) ValidateHardware(ctx sdk.Context) error {
    if !hv.validationEnabled {
        return nil
    }

    // Only check every 100 blocks to reduce overhead
    if ctx.BlockHeight()-hv.lastCheckHeight < 100 {
        return nil
    }
    hv.lastCheckHeight = ctx.BlockHeight()

    if hv.genesisBinding == nil {
        return fmt.Errorf("genesis hardware binding not loaded")
    }

    // Extract current hardware
    current, err := ExtractCurrentHardware()
    if err != nil {
        return fmt.Errorf("failed to extract current hardware: %w", err)
    }

    // Compare persistent identifiers (allow boot ID to change)
    if current.Components.WindowsUUID != hv.genesisBinding.Components.WindowsUUID {
        return fmt.Errorf("hardware mismatch: Windows UUID changed")
    }

    if current.Components.CPUInfo != hv.genesisBinding.Components.CPUInfo {
        return fmt.Errorf("hardware mismatch: CPU changed")
    }

    hv.logger.Info("Hardware validation successful",
        "height", ctx.BlockHeight(),
        "hardware_hash", current.HardwareHash[:16]+"...")

    return nil
}

// ExtractCurrentHardware gets current machine's hardware identifiers
func ExtractCurrentHardware() (*HardwareBinding, error) {
    // Try to execute inline hardware extraction
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

// EnableValidation enables hardware validation
func (hv *HardwareValidator) EnableValidation(enable bool) {
    hv.validationEnabled = enable
    hv.logger.Info("Hardware validation state changed",
        "enabled", enable)
}

// BeginBlock performs hardware validation at block start
func (hv *HardwareValidator) BeginBlock(ctx sdk.Context) {
    if err := hv.ValidateHardware(ctx); err != nil {
        hv.logger.Error("Hardware validation failed",
            "error", err,
            "height", ctx.BlockHeight())

        // In production with full integration, this would panic to halt the chain
        // For now, just log the error
        // panic(fmt.Sprintf("CRITICAL: Hardware validation failed: %v", err))
    }
}