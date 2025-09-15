package pairing

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	pairingsimulation "racecar-web/x/pairing/simulation"
	"racecar-web/x/pairing/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	pairingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pairingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgInitiateBidirectionalPairing          = "op_weight_msg_pairing"
		defaultWeightMsgInitiateBidirectionalPairing int = 100
	)

	var weightMsgInitiateBidirectionalPairing int
	simState.AppParams.GetOrGenerate(opWeightMsgInitiateBidirectionalPairing, &weightMsgInitiateBidirectionalPairing, nil,
		func(_ *rand.Rand) {
			weightMsgInitiateBidirectionalPairing = defaultWeightMsgInitiateBidirectionalPairing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInitiateBidirectionalPairing,
		pairingsimulation.SimulateMsgInitiateBidirectionalPairing(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCompletePairing          = "op_weight_msg_pairing"
		defaultWeightMsgCompletePairing int = 100
	)

	var weightMsgCompletePairing int
	simState.AppParams.GetOrGenerate(opWeightMsgCompletePairing, &weightMsgCompletePairing, nil,
		func(_ *rand.Rand) {
			weightMsgCompletePairing = defaultWeightMsgCompletePairing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCompletePairing,
		pairingsimulation.SimulateMsgCompletePairing(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRevokePairing          = "op_weight_msg_pairing"
		defaultWeightMsgRevokePairing int = 100
	)

	var weightMsgRevokePairing int
	simState.AppParams.GetOrGenerate(opWeightMsgRevokePairing, &weightMsgRevokePairing, nil,
		func(_ *rand.Rand) {
			weightMsgRevokePairing = defaultWeightMsgRevokePairing
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRevokePairing,
		pairingsimulation.SimulateMsgRevokePairing(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
