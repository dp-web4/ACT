package types

func NewMsgUpdateLctStatus(creator string, lctId string, newStatus string, reason string) *MsgUpdateLctStatus {
	return &MsgUpdateLctStatus{
		Creator:   creator,
		LctId:     lctId,
		NewStatus: newStatus,
		Reason:    reason,
	}
}
