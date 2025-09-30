package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Simple types for standalone tool
type Web4LCT struct {
	LCTID            string            `json:"lct_id"`
	Subject          string            `json:"subject"`
	Binding          Web4Binding       `json:"binding"`
	BirthCertificate *Web4BirthCert    `json:"birth_certificate,omitempty"`
	MRH              Web4MRH           `json:"mrh"`
	Policy           Web4Policy        `json:"policy"`
	Attestations     []Web4Attestation `json:"attestations"`
	Lineage          []Web4Lineage     `json:"lineage"`
	Revocation       *Web4Revocation   `json:"revocation,omitempty"`
}

type Web4Binding struct {
	EntityType     string `json:"entity_type"`
	PublicKey      string `json:"public_key"`
	HardwareAnchor string `json:"hardware_anchor,omitempty"`
	CreatedAt      string `json:"created_at"`
	BindingProof   string `json:"binding_proof"`
}

type Web4BirthCert struct {
	CitizenRole     string   `json:"citizen_role"`
	Context         string   `json:"context"`
	BirthTimestamp  string   `json:"birth_timestamp"`
	ParentEntity    string   `json:"parent_entity,omitempty"`
	BirthWitnesses  []string `json:"birth_witnesses"`
	FoundingPurpose string   `json:"founding_purpose,omitempty"`
}

type Web4MRH struct {
	Bound        []Web4BoundRelation   `json:"bound"`
	Paired       []Web4PairedRelation  `json:"paired"`
	Witnessing   []Web4WitnessRelation `json:"witnessing"`
	HorizonDepth int                   `json:"horizon_depth"`
	LastUpdated  string                `json:"last_updated"`
}

type Web4BoundRelation struct {
	LCTID        string `json:"lct_id"`
	Type         string `json:"type"`
	Timestamp    string `json:"ts"`
	BindingCtx   string `json:"binding_context,omitempty"`
}

type Web4PairedRelation struct {
	LCTID       string `json:"lct_id"`
	PairingType string `json:"pairing_type"`
	Permanent   bool   `json:"permanent"`
	Context     string `json:"context,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	Timestamp   string `json:"ts"`
}

type Web4WitnessRelation struct {
	LCTID           string `json:"lct_id"`
	Role            string `json:"role"`
	LastAttestation string `json:"last_attestation"`
}

type Web4Policy struct {
	Capabilities []string               `json:"capabilities"`
	Constraints  map[string]interface{} `json:"constraints"`
}

type Web4Attestation struct {
	Witness   string `json:"witness"`
	Type      string `json:"type"`
	Signature string `json:"sig"`
	Timestamp string `json:"ts"`
	Evidence  string `json:"evidence,omitempty"`
}

type Web4Lineage struct {
	Parent    string `json:"parent,omitempty"`
	Reason    string `json:"reason"`
	Timestamp string `json:"ts"`
}

type Web4Revocation struct {
	Status    string `json:"status"`
	Timestamp string `json:"ts,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

func extractHardwareHash() (string, error) {
	scriptPath := "./source/extract_hardware.sh"
	_, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("hardware script not found: %v", err)
	}

	// Run the script
	cmd := exec.Command("bash", scriptPath)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to extract hardware: %v", err)
	}

	// Parse output for hardware hash
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Hardware Hash") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[len(parts)-1], nil
			}
		}
	}

	return "", fmt.Errorf("hardware hash not found in output")
}

func createSociety4LCT(hardwareHash string) *Web4LCT {
	now := "2025-10-01T00:00:00Z" // Will be set to actual genesis time

	lct := &Web4LCT{
		LCTID:   "lct:web4:mb32:society4self0001",  // Temporary, will be computed from binding_proof
		Subject: "did:web4:society4:king:claudius",

		Binding: Web4Binding{
			EntityType:     "device",
			PublicKey:      "mb64:pending",  // Will be generated
			HardwareAnchor: fmt.Sprintf("eat:mb64:hw:%s", hardwareHash),
			CreatedAt:      now,
			BindingProof:   "cose:pending",  // Will be signed
		},

		BirthCertificate: nil, // Will be added after federation witnesses

		MRH: Web4MRH{
			Bound: []Web4BoundRelation{
				{
					LCTID:      fmt.Sprintf("lct:web4:hardware:wsl2:%s", hardwareHash[:16]),
					Type:       "parent",
					Timestamp:  now,
					BindingCtx: "wsl2_hardware_sovereignty",
				},
			},
			Paired:       []Web4PairedRelation{},
			Witnessing:   []Web4WitnessRelation{},
			HorizonDepth: 3,
			LastUpdated:  now,
		},

		Policy: Web4Policy{
			Capabilities: []string{
				"pairing:initiate",
				"consensus:participate",
				"hardware:validate",
				"temporal:authenticate",
				"pending:consensus",
			},
			Constraints: map[string]interface{}{
				"hardware_hash":    hardwareHash,
				"network_mobility": true,
				"requires_quorum":  3,
				"atp_allocation":   1000,
			},
		},

		Attestations: []Web4Attestation{},

		Lineage: []Web4Lineage{
			{
				Parent:    "",
				Reason:    "genesis",
				Timestamp: now,
			},
		},

		Revocation: &Web4Revocation{
			Status: "active",
		},
	}

	return lct
}

func main() {
	fmt.Println("Society 4 Genesis LCT Creator")
	fmt.Println("===============================\n")

	// Extract hardware hash
	fmt.Println("Extracting hardware hash...")
	hardwareHash, err := extractHardwareHash()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nUsing known Society 4 hardware hash:")
		hardwareHash = "93e766842ee7882a248e7d55ef3269b95e1735b0be88b94287b18029d1851759"
	}
	fmt.Printf("Hardware Hash: %s\n\n", hardwareHash)

	// Create LCT
	fmt.Println("Creating web4-compliant LCT...")
	lct := createSociety4LCT(hardwareHash)

	// Output as JSON
	output, err := json.MarshalIndent(lct, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Save to file
	outputFile := "./society4_genesis_lct.json"
	err = os.WriteFile(outputFile, output, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("LCT saved to: %s\n\n", outputFile)

	// Print summary
	fmt.Println("LCT Summary:")
	fmt.Printf("  LCT ID:      %s\n", lct.LCTID)
	fmt.Printf("  Subject:     %s\n", lct.Subject)
	fmt.Printf("  Entity Type: %s\n", lct.Binding.EntityType)
	fmt.Printf("  Hardware:    %s...\n", hardwareHash[:40])
	fmt.Printf("  Status:      %s\n", lct.Revocation.Status)
	fmt.Println("\nNext Steps:")
	fmt.Println("  1. Generate Ed25519 keypair")
	fmt.Println("  2. Sign binding with private key (COSE)")
	fmt.Println("  3. Compute final LCT ID from binding_proof")
	fmt.Println("  4. Request birth certificate from ACT Federation")
	fmt.Println("  5. Obtain witness signatures (Genesis, Society2, Sprout)")
	fmt.Println("\n⚠️  Birth certificate still required for full compliance")
}
