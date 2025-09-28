package lctmanager

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	lctmanagersimulation "racecar-web/x/lctmanager/simulation"
	"racecar-web/x/lctmanager/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	lctmanagerGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&lctmanagerGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateLctRelationship          = "op_weight_msg_lctmanager"
		defaultWeightMsgCreateLctRelationship int = 100
	)

	var weightMsgCreateLctRelationship int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateLctRelationship, &weightMsgCreateLctRelationship, nil,
		func(_ *rand.Rand) {
			weightMsgCreateLctRelationship = defaultWeightMsgCreateLctRelationship
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateLctRelationship,
		lctmanagersimulation.SimulateMsgCreateLctRelationship(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateLctStatus          = "op_weight_msg_lctmanager"
		defaultWeightMsgUpdateLctStatus int = 100
	)

	var weightMsgUpdateLctStatus int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateLctStatus, &weightMsgUpdateLctStatus, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateLctStatus = defaultWeightMsgUpdateLctStatus
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateLctStatus,
		lctmanagersimulation.SimulateMsgUpdateLctStatus(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgTerminateLctRelationship          = "op_weight_msg_lctmanager"
		defaultWeightMsgTerminateLctRelationship int = 100
	)

	var weightMsgTerminateLctRelationship int
	simState.AppParams.GetOrGenerate(opWeightMsgTerminateLctRelationship, &weightMsgTerminateLctRelationship, nil,
		func(_ *rand.Rand) {
			weightMsgTerminateLctRelationship = defaultWeightMsgTerminateLctRelationship
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTerminateLctRelationship,
		lctmanagersimulation.SimulateMsgTerminateLctRelationship(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
