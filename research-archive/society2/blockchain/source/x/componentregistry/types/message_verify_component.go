package types

func NewMsgVerifyComponent(creator string, componentId string) *MsgVerifyComponent {
	return &MsgVerifyComponent{
		Creator:     creator,
		ComponentId: componentId,
	}
}
