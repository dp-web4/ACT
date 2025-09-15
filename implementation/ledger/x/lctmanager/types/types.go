package types

type LCT struct {
	LctId              string
	ComponentAId       string
	ComponentBId       string
	PairingStatus      string
	CreatedAt          int64
	UpdatedAt          int64
	TrustAnchor        string
	OperationalContext string
	ProxyComponentId   string
}
