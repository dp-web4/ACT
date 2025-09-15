package types

func NewMsgTerminateLctRelationship(creator string, lctId string, reason string, notifyOffline bool) *MsgTerminateLctRelationship {
	return &MsgTerminateLctRelationship{
		Creator:       creator,
		LctId:         lctId,
		Reason:        reason,
		NotifyOffline: notifyOffline,
	}
}
