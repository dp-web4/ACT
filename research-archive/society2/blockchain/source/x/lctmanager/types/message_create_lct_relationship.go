package types

func NewMsgCreateLctRelationship(creator string, componentA string, componentB string, context string, proxyId string) *MsgCreateLctRelationship {
	return &MsgCreateLctRelationship{
		Creator:    creator,
		ComponentA: componentA,
		ComponentB: componentB,
		Context:    context,
		ProxyId:    proxyId,
	}
}
