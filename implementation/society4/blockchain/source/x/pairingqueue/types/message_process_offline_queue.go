package types

func NewMsgProcessOfflineQueue(creator string, componentId string) *MsgProcessOfflineQueue {
	return &MsgProcessOfflineQueue{
		Creator:     creator,
		ComponentId: componentId,
	}
}
