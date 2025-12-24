package componentregistry

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	componentregistrysimulation "racecar-web/x/componentregistry/simulation"
	"racecar-web/x/componentregistry/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	componentregistryGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&componentregistryGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgRegisterComponent          = "op_weight_msg_componentregistry"
		defaultWeightMsgRegisterComponent int = 100
	)

	var weightMsgRegisterComponent int
	simState.AppParams.GetOrGenerate(opWeightMsgRegisterComponent, &weightMsgRegisterComponent, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterComponent = defaultWeightMsgRegisterComponent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRegisterComponent,
		componentregistrysimulation.SimulateMsgRegisterComponent(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAuthorization          = "op_weight_msg_componentregistry"
		defaultWeightMsgUpdateAuthorization int = 100
	)

	var weightMsgUpdateAuthorization int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAuthorization, &weightMsgUpdateAuthorization, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAuthorization = defaultWeightMsgUpdateAuthorization
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAuthorization,
		componentregistrysimulation.SimulateMsgUpdateAuthorization(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgVerifyComponent          = "op_weight_msg_componentregistry"
		defaultWeightMsgVerifyComponent int = 100
	)

	var weightMsgVerifyComponent int
	simState.AppParams.GetOrGenerate(opWeightMsgVerifyComponent, &weightMsgVerifyComponent, nil,
		func(_ *rand.Rand) {
			weightMsgVerifyComponent = defaultWeightMsgVerifyComponent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVerifyComponent,
		componentregistrysimulation.SimulateMsgVerifyComponent(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
