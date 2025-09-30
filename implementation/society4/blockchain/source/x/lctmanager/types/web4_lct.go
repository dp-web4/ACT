package types

import (
	"time"
)

// Web4Compliant LCT following web4-lct.md specification
type Web4LCT struct {
	LCTID            string            `json:"lct_id"`              // lct:web4:mb32...
	Subject          string            `json:"subject"`             // did:web4:key:z6Mk...
	Binding          Web4Binding       `json:"binding"`
	BirthCertificate *Web4BirthCert    `json:"birth_certificate,omitempty"`
	MRH              Web4MRH           `json:"mrh"`
	Policy           Web4Policy        `json:"policy"`
	Attestations     []Web4Attestation `json:"attestations,omitempty"`
	Lineage          []Web4Lineage     `json:"lineage,omitempty"`
	Revocation       *Web4Revocation   `json:"revocation,omitempty"`
}

// Web4Binding - Establishes permanent link between LCT and entity
type Web4Binding struct {
	EntityType     string    `json:"entity_type"` // device, ai, human, organization, etc
	PublicKey      string    `json:"public_key"`  // mb64:coseKey format
	HardwareAnchor string    `json:"hardware_anchor,omitempty"` // EAT token (RFC 9334)
	CreatedAt      time.Time `json:"created_at"`
	BindingProof   string    `json:"binding_proof"` // COSE signature
}

// Web4BirthCert - Foundational identity and context
type Web4BirthCert struct {
	CitizenRole     string    `json:"citizen_role"`  // lct:web4:role:citizen:...
	Context         string    `json:"context"`       // federation, nation, platform, etc
	BirthTimestamp  time.Time `json:"birth_timestamp"`
	ParentEntity    string    `json:"parent_entity,omitempty"` // lct:web4:...
	BirthWitnesses  []string  `json:"birth_witnesses"`         // Array of witness LCT IDs
	FoundingPurpose string    `json:"founding_purpose,omitempty"` // Optional purpose statement
}

// Web4MRH - Markov Relevancy Horizon tracking
type Web4MRH struct {
	Bound        []Web4BoundRelation   `json:"bound"`
	Paired       []Web4PairedRelation  `json:"paired"`
	Witnessing   []Web4WitnessRelation `json:"witnessing,omitempty"`
	HorizonDepth int                   `json:"horizon_depth,omitempty"` // default: 3
	LastUpdated  time.Time             `json:"last_updated"`
}

// Web4BoundRelation - Binding relationships (permanent)
type Web4BoundRelation struct {
	LCTID          string    `json:"lct_id"`
	RelationType   string    `json:"type"` // parent, child, sibling
	Timestamp      time.Time `json:"ts"`
	BindingContext string    `json:"binding_context,omitempty"` // Additional context
}

// Web4PairedRelation - Active pairings (can be temporary)
type Web4PairedRelation struct {
	LCTID       string    `json:"lct_id"`
	PairingType string    `json:"pairing_type"` // birth_certificate, role, operational
	Permanent   bool      `json:"permanent"`
	Context     string    `json:"context,omitempty"`     // For non-birth pairings
	SessionID   string    `json:"session_id,omitempty"`  // For operational pairings
	Timestamp   time.Time `json:"ts"`
}

// Web4WitnessRelation - Witness relationships
type Web4WitnessRelation struct {
	LCTID           string    `json:"lct_id"`
	Role            string    `json:"role"` // time, audit, oracle
	LastAttestation time.Time `json:"last_attestation"`
}

// Web4Policy - Capabilities and constraints
type Web4Policy struct {
	Capabilities []string               `json:"capabilities"` // e.g., pairing:initiate, metering:grant
	Constraints  map[string]interface{} `json:"constraints,omitempty"`
}

// Web4Attestation - Witnessing events
type Web4Attestation struct {
	Witness   string    `json:"witness"` // DID of witness
	Type      string    `json:"type"`    // time, audit, oracle, existence, action, state, quality
	Signature string    `json:"sig"`     // COSE signature
	Timestamp time.Time `json:"ts"`
	Evidence  string    `json:"evidence,omitempty"` // Optional evidence data
}

// Web4Lineage - LCT evolution tracking
type Web4Lineage struct {
	Parent    string    `json:"parent,omitempty"` // Previous LCT ID
	Reason    string    `json:"reason"`           // genesis, rotation, fork, upgrade
	Timestamp time.Time `json:"ts"`
}

// Web4Revocation - LCT revocation status
type Web4Revocation struct {
	Status    string    `json:"status"` // active, revoked
	Timestamp time.Time `json:"ts"`
	Reason    string    `json:"reason,omitempty"` // compromise, superseded, expired
}

// Society4Self LCT - The genesis 'self' LCT for Society 4
// This wraps the hardware binding in proper web4 format
type Society4SelfLCT struct {
	LCT          Web4LCT `json:"lct"`
	HardwareHash string  `json:"hardware_hash"` // The actual WSL2 hardware hash
	NetworkState string  `json:"network_state"` // home_federation, work_isolated
}

// NewSociety4SelfLCT creates the genesis self-LCT for Society 4
func NewSociety4SelfLCT(hardwareHash string, publicKey string, hardwareAnchorEAT string) *Society4SelfLCT {
	now := time.Now()

	lct := Web4LCT{
		// Will be computed: lct:web4:mb32(SHA256(binding_proof))
		LCTID:   "lct:web4:society4:self:pending",
		Subject: "did:web4:society4:king:claudius",

		Binding: Web4Binding{
			EntityType:     "device", // Society 4 is hardware-bound
			PublicKey:      publicKey, // mb64:coseKey format
			HardwareAnchor: hardwareAnchorEAT, // EAT token with hardware hash
			CreatedAt:      now,
			BindingProof:   "pending", // Will be COSE signature
		},

		// Birth certificate pending - needs federation witnesses
		BirthCertificate: nil,

		MRH: Web4MRH{
			Bound: []Web4BoundRelation{
				// Hardware binding relationship
				{
					LCTID:          "lct:web4:hardware:wsl2:" + hardwareHash[:16],
					RelationType:   "parent",
					Timestamp:      now,
					BindingContext: "wsl2_hardware_sovereignty",
				},
			},
			Paired: []Web4PairedRelation{
				// Will add birth certificate pairing once obtained
			},
			Witnessing: []Web4WitnessRelation{},
			HorizonDepth: 3,
			LastUpdated:  now,
		},

		Policy: Web4Policy{
			Capabilities: []string{
				"pairing:initiate",
				"consensus:participate",
				"hardware:validate",
				"temporal:authenticate",
			},
			Constraints: map[string]interface{}{
				"hardware_hash":    hardwareHash,
				"network_mobility": true,
				"requires_quorum":  3,
			},
		},

		Attestations: []Web4Attestation{},

		Lineage: []Web4Lineage{
			{
				Parent:    "", // Genesis - no parent
				Reason:    "genesis",
				Timestamp: now,
			},
		},

		Revocation: &Web4Revocation{
			Status: "active",
		},
	}

	return &Society4SelfLCT{
		LCT:          lct,
		HardwareHash: hardwareHash,
		NetworkState: "work_isolated", // Will be updated by temporal auth
	}
}

// AddBirthCertificate updates the LCT with federation birth certificate
func (s *Society4SelfLCT) AddBirthCertificate(
	parentEntity string,
	witnesses []string,
	purpose string,
) error {
	now := time.Now()

	s.LCT.BirthCertificate = &Web4BirthCert{
		CitizenRole:     "lct:web4:federation:act:citizen:society4",
		Context:         "federation",
		BirthTimestamp:  now,
		ParentEntity:    parentEntity, // lct:web4:federation:act
		BirthWitnesses:  witnesses,    // Genesis, Society2, Sprout
		FoundingPurpose: purpose,
	}

	// Add birth certificate pairing to MRH
	s.LCT.MRH.Paired = append([]Web4PairedRelation{
		{
			LCTID:       s.LCT.BirthCertificate.CitizenRole,
			PairingType: "birth_certificate",
			Permanent:   true,
			Context:     "ACT Federation citizenship",
			Timestamp:   now,
		},
	}, s.LCT.MRH.Paired...)

	s.LCT.MRH.LastUpdated = now
	return nil
}

// AddWitness adds a witness attestation to the LCT
func (s *Society4SelfLCT) AddWitness(
	witnessLCTID string,
	witnessDID string,
	witnessType string,
	signature string,
	evidence string,
) {
	now := time.Now()

	// Add to attestations
	s.LCT.Attestations = append(s.LCT.Attestations, Web4Attestation{
		Witness:   witnessDID,
		Type:      witnessType,
		Signature: signature,
		Timestamp: now,
		Evidence:  evidence,
	})

	// Add to MRH witnessing if not already present
	found := false
	for i := range s.LCT.MRH.Witnessing {
		if s.LCT.MRH.Witnessing[i].LCTID == witnessLCTID {
			s.LCT.MRH.Witnessing[i].LastAttestation = now
			found = true
			break
		}
	}

	if !found {
		s.LCT.MRH.Witnessing = append(s.LCT.MRH.Witnessing, Web4WitnessRelation{
			LCTID:           witnessLCTID,
			Role:            witnessType,
			LastAttestation: now,
		})
	}

	s.LCT.MRH.LastUpdated = now
}

// UpdateNetworkState updates the temporal/network context
func (s *Society4SelfLCT) UpdateNetworkState(newState string) {
	s.NetworkState = newState

	// Add as constraint
	if s.LCT.Policy.Constraints == nil {
		s.LCT.Policy.Constraints = make(map[string]interface{})
	}
	s.LCT.Policy.Constraints["current_network"] = newState
	s.LCT.Policy.Constraints["last_network_update"] = time.Now()
}

// ValidateCompliance checks if LCT meets web4 requirements
func (lct *Web4LCT) ValidateCompliance() []string {
	issues := []string{}

	// Required fields
	if lct.LCTID == "" || lct.LCTID == "pending" {
		issues = append(issues, "lct_id must be computed from binding_proof")
	}
	if lct.Subject == "" {
		issues = append(issues, "subject DID is required")
	}
	if lct.Binding.EntityType == "" {
		issues = append(issues, "binding.entity_type is required")
	}
	if lct.Binding.PublicKey == "" {
		issues = append(issues, "binding.public_key is required")
	}
	if lct.Binding.BindingProof == "" || lct.Binding.BindingProof == "pending" {
		issues = append(issues, "binding.binding_proof must be valid COSE signature")
	}

	// Birth certificate required
	if lct.BirthCertificate == nil {
		issues = append(issues, "birth_certificate is required for compliance")
	} else {
		if lct.BirthCertificate.CitizenRole == "" {
			issues = append(issues, "birth_certificate.citizen_role is required")
		}
		if len(lct.BirthCertificate.BirthWitnesses) == 0 {
			issues = append(issues, "birth_certificate.birth_witnesses required (min 1)")
		}
	}

	// MRH must have birth cert pairing as first entry
	if len(lct.MRH.Paired) == 0 {
		issues = append(issues, "mrh.paired must contain birth_certificate pairing")
	} else if lct.MRH.Paired[0].PairingType != "birth_certificate" {
		issues = append(issues, "mrh.paired[0] must be birth_certificate pairing")
	}

	// Policy capabilities required
	if len(lct.Policy.Capabilities) == 0 {
		issues = append(issues, "policy.capabilities cannot be empty")
	}

	// Lineage must have genesis entry
	if len(lct.Lineage) == 0 {
		issues = append(issues, "lineage must contain at least genesis entry")
	}

	return issues
}
