package pairingqueue

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	pairingqueuesimulation "racecar-web/x/pairingqueue/simulation"
	"racecar-web/x/pairingqueue/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	pairingqueueGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pairingqueueGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgQueuePairingRequest          = "op_weight_msg_pairingqueue"
		defaultWeightMsgQueuePairingRequest int = 100
	)

	var weightMsgQueuePairingRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgQueuePairingRequest, &weightMsgQueuePairingRequest, nil,
		func(_ *rand.Rand) {
			weightMsgQueuePairingRequest = defaultWeightMsgQueuePairingRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgQueuePairingRequest,
		pairingqueuesimulation.SimulateMsgQueuePairingRequest(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgProcessOfflineQueue          = "op_weight_msg_pairingqueue"
		defaultWeightMsgProcessOfflineQueue int = 100
	)

	var weightMsgProcessOfflineQueue int
	simState.AppParams.GetOrGenerate(opWeightMsgProcessOfflineQueue, &weightMsgProcessOfflineQueue, nil,
		func(_ *rand.Rand) {
			weightMsgProcessOfflineQueue = defaultWeightMsgProcessOfflineQueue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgProcessOfflineQueue,
		pairingqueuesimulation.SimulateMsgProcessOfflineQueue(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCancelRequest          = "op_weight_msg_pairingqueue"
		defaultWeightMsgCancelRequest int = 100
	)

	var weightMsgCancelRequest int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelRequest, &weightMsgCancelRequest, nil,
		func(_ *rand.Rand) {
			weightMsgCancelRequest = defaultWeightMsgCancelRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelRequest,
		pairingqueuesimulation.SimulateMsgCancelRequest(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
