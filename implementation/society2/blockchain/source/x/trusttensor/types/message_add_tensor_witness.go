package types

func NewMsgAddTensorWitness(creator string, tensorId string, dimension string, witnessLct string, confidence string, evidenceHash string) *MsgAddTensorWitness {
	return &MsgAddTensorWitness{
		Creator:      creator,
		TensorId:     tensorId,
		Dimension:    dimension,
		WitnessLct:   witnessLct,
		Confidence:   confidence,
		EvidenceHash: evidenceHash,
	}
}
