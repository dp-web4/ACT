package types

func NewMsgInitiateBidirectionalPairing(creator string, componentA string, componentB string, operationalContext string, proxyId string, forceImmediate bool) *MsgInitiateBidirectionalPairing {
	return &MsgInitiateBidirectionalPairing{
		Creator:            creator,
		ComponentA:         componentA,
		ComponentB:         componentB,
		OperationalContext: operationalContext,
		ProxyId:            proxyId,
		ForceImmediate:     forceImmediate,
	}
}
