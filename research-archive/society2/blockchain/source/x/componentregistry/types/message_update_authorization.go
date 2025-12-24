package types

func NewMsgUpdateAuthorization(creator string, componentId string, authRules string) *MsgUpdateAuthorization {
	return &MsgUpdateAuthorization{
		Creator:     creator,
		ComponentId: componentId,
		AuthRules:   authRules,
	}
}
