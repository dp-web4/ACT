package types

func NewMsgUpdateTensorScore(creator string, tensorId string, dimension string, value string, context string, witnessData string) *MsgUpdateTensorScore {
	return &MsgUpdateTensorScore{
		Creator:     creator,
		TensorId:    tensorId,
		Dimension:   dimension,
		Value:       value,
		Context:     context,
		WitnessData: witnessData,
	}
}
