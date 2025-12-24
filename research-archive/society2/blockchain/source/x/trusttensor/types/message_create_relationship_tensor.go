package types

func NewMsgCreateRelationshipTensor(creator string, lctId string, tensorType string, context string) *MsgCreateRelationshipTensor {
	return &MsgCreateRelationshipTensor{
		Creator:    creator,
		LctId:      lctId,
		TensorType: tensorType,
		Context:    context,
	}
}
