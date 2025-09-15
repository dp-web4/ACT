package energycycle

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/energycycle/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "GetRelationshipEnergyBalance",
					Use:            "get-relationship-energy-balance [lct-id]",
					Short:          "Query get-relationship-energy-balance",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}},
				},

				{
					RpcMethod:      "CalculateRelationshipV3",
					Use:            "calculate-relationship-v-3 [operation-id]",
					Short:          "Query calculate-relationship-v3",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "operation_id"}},
				},

				{
					RpcMethod:      "GetEnergyFlowHistory",
					Use:            "get-energy-flow-history [lct-id]",
					Short:          "Query get-energy-flow-history",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateRelationshipEnergyOperation",
					Use:            "create-relationship-energy-operation [source-lct] [target-lct] [energy-amount] [operation-type]",
					Short:          "Send a create-relationship-energy-operation tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "source_lct"}, {ProtoField: "target_lct"}, {ProtoField: "energy_amount"}, {ProtoField: "operation_type"}},
				},
				{
					RpcMethod:      "ExecuteEnergyTransfer",
					Use:            "execute-energy-transfer [operation-id] [transfer-data]",
					Short:          "Send a execute-energy-transfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "operation_id"}, {ProtoField: "transfer_data"}},
				},
				{
					RpcMethod:      "ValidateRelationshipValue",
					Use:            "validate-relationship-value [operation-id] [recipient-validation] [utility-rating] [trust-context]",
					Short:          "Send a validate-relationship-value tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "operation_id"}, {ProtoField: "recipient_validation"}, {ProtoField: "utility_rating"}, {ProtoField: "trust_context"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
