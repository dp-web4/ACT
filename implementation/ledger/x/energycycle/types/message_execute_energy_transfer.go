package types

func NewMsgExecuteEnergyTransfer(creator string, operationId string, transferData string) *MsgExecuteEnergyTransfer {
	return &MsgExecuteEnergyTransfer{
		Creator:      creator,
		OperationId:  operationId,
		TransferData: transferData,
	}
}
