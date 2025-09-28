package types

func NewMsgRevokePairing(creator string, lctId string, reason string, notifyOffline bool) *MsgRevokePairing {
	return &MsgRevokePairing{
		Creator:       creator,
		LctId:         lctId,
		Reason:        reason,
		NotifyOffline: notifyOffline,
	}
}
