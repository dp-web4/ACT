package types

import (
  "cosmossdk.io/math"
)

// EnergyPool manages ATP/ADP tokens
type EnergyPool struct {
  ID                  string   `json:"id"`
  AtpBalance          math.Int `json:"atp_balance"`
  AdpBalance          math.Int `json:"adp_balance"`
  VelocityRequirement float64  `json:"velocity_requirement"`
  DemurrageRate       float64  `json:"demurrage_rate"`
}

// ATPToken represents charged energy
type ATPToken struct {
  ID           string   `json:"id"`
  Amount       math.Int `json:"amount"`
  RechargedBy  string   `json:"recharged_by"`
  RechargeTime int64    `json:"recharge_time"`
  WorkProof    []byte   `json:"work_proof"`
}

// ADPToken represents discharged energy
type ADPToken struct {
  ID            string    `json:"id"`
  Amount        math.Int  `json:"amount"`
  DischargedBy  string    `json:"discharged_by"`
  DischargeTime int64     `json:"discharge_time"`
  R6Action      *R6Action `json:"r6_action,omitempty"`
}

// R6Action represents an action in the R6 framework
type R6Action struct {
  Rules     string `json:"rules"`
  Roles     string `json:"roles"`
  Request   string `json:"request"`
  Reference string `json:"reference"`
  Resource  string `json:"resource"`
  Result    string `json:"result"`
}