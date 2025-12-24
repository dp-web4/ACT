package types

func NewMsgCreateRelationshipEnergyOperation(creator string, sourceLct string, targetLct string, energyAmount string, operationType string) *MsgCreateRelationshipEnergyOperation {
	return &MsgCreateRelationshipEnergyOperation{
		Creator:       creator,
		SourceLct:     sourceLct,
		TargetLct:     targetLct,
		EnergyAmount:  energyAmount,
		OperationType: operationType,
	}
}
