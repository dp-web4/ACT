package types

import (
	"context"
	"testing"
)

func TestMockMySQLBackend(t *testing.T) {
	backend := NewMockMySQLBackend()
	ctx := context.Background()

	t.Run("Valid Pairing", func(t *testing.T) {
		// Test valid pairing: battery module with battery pack
		allowed, reason, err := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-001", "MODBATT-PACK-RC001-A")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Expected pairing to be allowed, but got: %s", reason)
		}
		if reason == "" {
			t.Error("Expected reason message, got empty string")
		}
	})

	t.Run("Invalid Pairing", func(t *testing.T) {
		// Test invalid pairing: battery module with sensor
		allowed, reason, err := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-001", "TEMP-SENSOR-RC001")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if allowed {
			t.Errorf("Expected pairing to be denied, but got: %s", reason)
		}
		if reason == "" {
			t.Error("Expected reason message, got empty string")
		}
	})

	t.Run("Negative Indicator", func(t *testing.T) {
		// Test negative indicator: faulty component
		allowed, reason, err := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-FAULTY", "MODBATT-PACK-RC001-A")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if allowed {
			t.Errorf("Expected pairing to be denied due to negative indicator, but got: %s", reason)
		}
		if reason == "" {
			t.Error("Expected reason message, got empty string")
		}
	})

	t.Run("Bidirectional Pairing", func(t *testing.T) {
		// Test bidirectional pairing: A->B and B->A should both work
		allowed1, _, err1 := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-001", "MODBATT-PACK-RC001-A")
		allowed2, _, err2 := backend.VerifyComponentPairing(ctx, "MODBATT-PACK-RC001-A", "MODBATT-MOD-RC001-001")

		if err1 != nil || err2 != nil {
			t.Fatalf("Unexpected errors: %v, %v", err1, err2)
		}
		if !allowed1 || !allowed2 {
			t.Error("Expected bidirectional pairing to work")
		}
	})

	t.Run("Component Metadata", func(t *testing.T) {
		// Test component metadata retrieval
		metadata, err := backend.GetComponentMetadata(ctx, "MODBATT-MOD-RC001-001")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if metadata["type"] != "battery_module" {
			t.Errorf("Expected type 'battery_module', got: %v", metadata["type"])
		}
		if metadata["capacity"] != "25.6kWh" {
			t.Errorf("Expected capacity '25.6kWh', got: %v", metadata["capacity"])
		}
	})

	t.Run("Unknown Component", func(t *testing.T) {
		// Test unknown component metadata
		metadata, err := backend.GetComponentMetadata(ctx, "UNKNOWN-COMPONENT")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if metadata["type"] != "unknown" {
			t.Errorf("Expected type 'unknown', got: %v", metadata["type"])
		}
	})

	t.Run("Dynamic Rules", func(t *testing.T) {
		// Test adding dynamic pairing rules
		backend.AddPairingRule("NEW-MODULE", []string{"NEW-PACK"})

		allowed, _, err := backend.VerifyComponentPairing(ctx, "NEW-MODULE", "NEW-PACK")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Error("Expected dynamic pairing rule to work")
		}
	})

	t.Run("Dynamic Negative Indicator", func(t *testing.T) {
		// Test adding dynamic negative indicators
		backend.AddNegativeIndicator("TEST-BLOCKED")

		allowed, reason, err := backend.VerifyComponentPairing(ctx, "MODBATT-MOD-RC001-TEST-BLOCKED", "MODBATT-PACK-RC001-A")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if allowed {
			t.Errorf("Expected pairing to be denied due to dynamic negative indicator, but got: %s", reason)
		}
	})
}
