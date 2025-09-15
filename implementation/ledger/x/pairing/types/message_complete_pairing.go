package types

func NewMsgCompletePairing(creator string, challengeId string, componentAAuth string, componentBAuth string, sessionContext string) *MsgCompletePairing {
	return &MsgCompletePairing{
		Creator:        creator,
		ChallengeId:    challengeId,
		ComponentAAuth: componentAAuth,
		ComponentBAuth: componentBAuth,
		SessionContext: sessionContext,
	}
}
