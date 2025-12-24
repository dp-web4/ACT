package types

func NewMsgQueuePairingRequest(creator string, initiatorId string, targetId string, requestType string, proxyId string) *MsgQueuePairingRequest {
	return &MsgQueuePairingRequest{
		Creator:     creator,
		InitiatorId: initiatorId,
		TargetId:    targetId,
		RequestType: requestType,
		ProxyId:     proxyId,
	}
}
