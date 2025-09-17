package types

import (
  trusttypes "racecar-web/x/trusttensor/types"
    "cosmossdk.io/math"
)

// Web4 Compliant LCT (Linked Context Token)
type LCT struct {
    ID               string        `json:"id"`
    EntityType       string        `json:"entity_type"` 
    Identity         LCTIdentity   `json:"identity"`
    MRH              MRH           `json:"mrh"`
    BirthCertificate *BirthCertificate `json:"birth_certificate,omitempty"`
    CreatedAt        int64         `json:"created_at"`
    UpdatedAt        int64         `json:"updated_at"`

    T3Tensor         *trusttypes.T3Tensor `json:"t3_tensor,omitempty"`
    V3Tensor         *trusttypes.V3Tensor `json:"v3_tensor,omitempty"`
}

// Cryptographic Identity for Web4
type LCTIdentity struct {
    Ed25519PublicKey []byte `json:"ed25519_public_key"`
    X25519PublicKey  []byte `json:"x25519_public_key"`
    BindingSignature []byte `json:"binding_signature"`
}

// Markov Relevancy Horizon
type MRH struct {
    Bound      []string `json:"bound"`
    Paired     []string `json:"paired"`
    Witnessing []string `json:"witnessing"`
    Broadcast  []string `json:"broadcast"`
}

type BirthCertificate struct {
    Society          string   `json:"society"`
    Rights           []string `json:"rights"`
    Responsibilities []string `json:"responsibilities"`
    InitialATP       math.Int  `json:"initial_atp"`
    IssuedAt         int64    `json:"issued_at"`
}
