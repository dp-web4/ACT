package types

func NewMsgCancelRequest(creator string, requestId string, reason string) *MsgCancelRequest {
	return &MsgCancelRequest{
		Creator:   creator,
		RequestId: requestId,
		Reason:    reason,
	}
}
