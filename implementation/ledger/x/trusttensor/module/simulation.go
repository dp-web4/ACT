package trusttensor

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	trusttensorsimulation "racecar-web/x/trusttensor/simulation"
	"racecar-web/x/trusttensor/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	trusttensorGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&trusttensorGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateRelationshipTensor          = "op_weight_msg_trusttensor"
		defaultWeightMsgCreateRelationshipTensor int = 100
	)

	var weightMsgCreateRelationshipTensor int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateRelationshipTensor, &weightMsgCreateRelationshipTensor, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRelationshipTensor = defaultWeightMsgCreateRelationshipTensor
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateRelationshipTensor,
		trusttensorsimulation.SimulateMsgCreateRelationshipTensor(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateTensorScore          = "op_weight_msg_trusttensor"
		defaultWeightMsgUpdateTensorScore int = 100
	)

	var weightMsgUpdateTensorScore int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTensorScore, &weightMsgUpdateTensorScore, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTensorScore = defaultWeightMsgUpdateTensorScore
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTensorScore,
		trusttensorsimulation.SimulateMsgUpdateTensorScore(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddTensorWitness          = "op_weight_msg_trusttensor"
		defaultWeightMsgAddTensorWitness int = 100
	)

	var weightMsgAddTensorWitness int
	simState.AppParams.GetOrGenerate(opWeightMsgAddTensorWitness, &weightMsgAddTensorWitness, nil,
		func(_ *rand.Rand) {
			weightMsgAddTensorWitness = defaultWeightMsgAddTensorWitness
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddTensorWitness,
		trusttensorsimulation.SimulateMsgAddTensorWitness(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
