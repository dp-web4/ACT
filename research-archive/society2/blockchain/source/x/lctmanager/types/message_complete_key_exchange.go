package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgInitiateLCTMediatedPairing{}
var _ sdk.Msg = &MsgCompleteLCTMediatedPairing{}

func NewMsgInitiateLCTMediatedPairing(creator string, initiatorLctId string, targetLctId string, context string, proxyLctId string, expiresAt int64) *MsgInitiateLCTMediatedPairing {
	return &MsgInitiateLCTMediatedPairing{
		Creator:        creator,
		InitiatorLctId: initiatorLctId,
		TargetLctId:    targetLctId,
		Context:        context,
		ProxyLctId:     proxyLctId,
		ExpiresAt:      expiresAt,
	}
}

func NewMsgCompleteLCTMediatedPairing(creator string, pairingId string, initiatorResponse string, targetResponse string, sessionKeyData []byte) *MsgCompleteLCTMediatedPairing {
	return &MsgCompleteLCTMediatedPairing{
		Creator:           creator,
		PairingId:         pairingId,
		InitiatorResponse: initiatorResponse,
		TargetResponse:    targetResponse,
		SessionKeyData:    sessionKeyData,
	}
}
