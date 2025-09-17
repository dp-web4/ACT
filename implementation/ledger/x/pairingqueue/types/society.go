package types

import (
    "cosmossdk.io/math"
)

// Society membership for Web4 ACT
type SocietyMembership struct {
    SocietyLCT    string   `json:"society_lct"`
    MemberLCT     string   `json:"member_lct"`
    CitizenRole   string   `json:"citizen_role"`
    Rights        []string `json:"rights"`
    Responsibilities []string `json:"responsibilities"`
    JoinedAt      int64    `json:"joined_at"`
    ATP_Allocated math.Int  `json:"atp_allocated"`
}

// Birth certificate for new entities
type BirthCertificate struct {
    EntityLCT  string   `json:"entity_lct"`
    EntityType string   `json:"entity_type"`
    Society    string   `json:"society"`
    IssuedBy   string   `json:"issued_by"`
    IssuedAt   int64    `json:"issued_at"`
    Witnesses  []string `json:"witnesses"`
}
