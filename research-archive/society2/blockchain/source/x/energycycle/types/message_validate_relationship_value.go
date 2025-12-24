package types

func NewMsgValidateRelationshipValue(creator string, operationId string, recipientValidation string, utilityRating string, trustContext string) *MsgValidateRelationshipValue {
	return &MsgValidateRelationshipValue{
		Creator:             creator,
		OperationId:         operationId,
		RecipientValidation: recipientValidation,
		UtilityRating:       utilityRating,
		TrustContext:        trustContext,
	}
}
