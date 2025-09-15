package energycycle

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	energycyclesimulation "racecar-web/x/energycycle/simulation"
	"racecar-web/x/energycycle/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	energycycleGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&energycycleGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateRelationshipEnergyOperation          = "op_weight_msg_energycycle"
		defaultWeightMsgCreateRelationshipEnergyOperation int = 100
	)

	var weightMsgCreateRelationshipEnergyOperation int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateRelationshipEnergyOperation, &weightMsgCreateRelationshipEnergyOperation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRelationshipEnergyOperation = defaultWeightMsgCreateRelationshipEnergyOperation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateRelationshipEnergyOperation,
		energycyclesimulation.SimulateMsgCreateRelationshipEnergyOperation(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgExecuteEnergyTransfer          = "op_weight_msg_energycycle"
		defaultWeightMsgExecuteEnergyTransfer int = 100
	)

	var weightMsgExecuteEnergyTransfer int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteEnergyTransfer, &weightMsgExecuteEnergyTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteEnergyTransfer = defaultWeightMsgExecuteEnergyTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteEnergyTransfer,
		energycyclesimulation.SimulateMsgExecuteEnergyTransfer(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgValidateRelationshipValue          = "op_weight_msg_energycycle"
		defaultWeightMsgValidateRelationshipValue int = 100
	)

	var weightMsgValidateRelationshipValue int
	simState.AppParams.GetOrGenerate(opWeightMsgValidateRelationshipValue, &weightMsgValidateRelationshipValue, nil,
		func(_ *rand.Rand) {
			weightMsgValidateRelationshipValue = defaultWeightMsgValidateRelationshipValue
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgValidateRelationshipValue,
		energycyclesimulation.SimulateMsgValidateRelationshipValue(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
