package types

// TrustRecord tracks trust for an LCT in a specific role
type TrustRecord struct {
  LctId      string  `json:"lct_id"`
  Role       string  `json:"role"`
  T3Score    float64 `json:"t3_score"`
  V3Score    float64 `json:"v3_score"`
  LastUpdate int64   `json:"last_update"`
}

// T3Tensor - Talent, Training, Temperament
type T3Tensor struct {
  Talent      float64 `json:"talent"`
  Training    float64 `json:"training"`
  Temperament float64 `json:"temperament"`
}

// V3Tensor - Veracity, Validity, Value
type V3Tensor struct {
  Veracity float64 `json:"veracity"`
  Validity float64 `json:"validity"`
  Value    float64 `json:"value"`
}

// Outcome represents the result of an action
type Outcome struct {
  Success        bool    `json:"success"`
  ValueGenerated float64 `json:"value_generated"`
  Witnesses      []string `json:"witnesses"`
}

// GravityRecord tracks trust gravity effects
type GravityRecord struct {
  LctId     string  `json:"lct_id"`
  Gravity   float64 `json:"gravity"`
  Timestamp int64   `json:"timestamp"`
}