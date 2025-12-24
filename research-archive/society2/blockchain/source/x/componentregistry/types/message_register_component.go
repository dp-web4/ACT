package types

func NewMsgRegisterComponent(creator string, componentId string, componentType string, manufacturerData string) *MsgRegisterComponent {
	return &MsgRegisterComponent{
		Creator:          creator,
		ComponentId:      componentId,
		ComponentType:    componentType,
		ManufacturerData: manufacturerData,
	}
}
