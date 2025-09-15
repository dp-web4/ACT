package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"racecar-web/x/componentregistry/types"
)

func main() {
	fmt.Println("🚗 Race Car Demo: Mock MySQL Component Verification Backend")
	fmt.Println(strings.Repeat("=", 60))

	// Create the mock MySQL backend
	backend := types.NewMockMySQLBackend()
	ctx := context.Background()

	// Demo 1: Valid component pairing
	fmt.Println("\n✅ Demo 1: Valid Component Pairing")
	allowed, reason, err := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-001", "MODBATT-PACK-RC001-A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Pairing: MODBATT-MOD-RC001-001 ↔ MODBATT-PACK-RC001-A\n")
	fmt.Printf("   Result: %t\n", allowed)
	fmt.Printf("   Reason: %s\n", reason)

	// Demo 2: Invalid component pairing
	fmt.Println("\n❌ Demo 2: Invalid Component Pairing")
	allowed, reason, err = backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-001", "TEMP-SENSOR-RC001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Pairing: MODBATT-MOD-RC001-001 ↔ TEMP-SENSOR-RC001\n")
	fmt.Printf("   Result: %t\n", allowed)
	fmt.Printf("   Reason: %s\n", reason)

	// Demo 3: Negative indicator detection
	fmt.Println("\n🚫 Demo 3: Negative Indicator Detection")
	allowed, reason, err = backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-FAULTY", "MODBATT-PACK-RC001-A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Pairing: MODBATT-MOD-RC001-FAULTY ↔ MODBATT-PACK-RC001-A\n")
	fmt.Printf("   Result: %t\n", allowed)
	fmt.Printf("   Reason: %s\n", reason)

	// Demo 4: Component metadata
	fmt.Println("\n📋 Demo 4: Component Metadata")
	metadata, err := backend.GetComponentMetadata(ctx, "MODBATT-MOD-RC001-001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Component: MODBATT-MOD-RC001-001\n")
	fmt.Printf("   Type: %v\n", metadata["type"])
	fmt.Printf("   Capacity: %v\n", metadata["capacity"])
	fmt.Printf("   Voltage: %v\n", metadata["voltage"])
	fmt.Printf("   Manufacturer: %v\n", metadata["manufacturer"])

	// Demo 5: Dynamic rule addition
	fmt.Println("\n🔧 Demo 5: Dynamic Rule Addition")
	backend.AddPairingRule("CUSTOM-MODULE", []string{"CUSTOM-PACK"})
	allowed, reason, err = backend.VerifyComponentPairing(ctx, "CUSTOM-MODULE", "CUSTOM-PACK")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Pairing: CUSTOM-MODULE ↔ CUSTOM-PACK\n")
	fmt.Printf("   Result: %t\n", allowed)
	fmt.Printf("   Reason: %s\n", reason)

	// Demo 6: Dynamic negative indicator
	fmt.Println("\n🚫 Demo 6: Dynamic Negative Indicator")
	backend.AddNegativeIndicator("DEMO-BLOCKED")
	allowed, reason, err = backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-DEMO-BLOCKED", "MODBATT-PACK-RC001-A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Pairing: MODBATT-MOD-RC001-DEMO-BLOCKED ↔ MODBATT-PACK-RC001-A\n")
	fmt.Printf("   Result: %t\n", allowed)
	fmt.Printf("   Reason: %s\n", reason)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎯 Demo Complete! The mock MySQL backend is working perfectly.")
	fmt.Println("   This provides a simple, configurable way to verify component pairings")
	fmt.Println("   without needing a real MySQL database for the race car demo.")
}
