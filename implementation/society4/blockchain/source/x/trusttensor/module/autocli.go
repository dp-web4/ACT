package trusttensor

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/trusttensor/types"
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
					RpcMethod:      "GetRelationshipTensor",
					Use:            "get-relationship-tensor [lct-id] [tensor-type]",
					Short:          "Query get-relationship-tensor",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "tensor_type"}},
				},

				{
					RpcMethod:      "CalculateRelationshipTrust",
					Use:            "calculate-relationship-trust [lct-id] [context]",
					Short:          "Query calculate-relationship-trust",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "context"}},
				},

				{
					RpcMethod:      "GetTensorHistory",
					Use:            "get-tensor-history [tensor-id]",
					Short:          "Query get-tensor-history",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tensor_id"}},
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
					RpcMethod:      "CreateRelationshipTensor",
					Use:            "create-relationship-tensor [lct-id] [tensor-type] [context]",
					Short:          "Send a create-relationship-tensor tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "tensor_type"}, {ProtoField: "context"}},
				},
				{
					RpcMethod:      "UpdateTensorScore",
					Use:            "update-tensor-score [tensor-id] [dimension] [value] [context] [witness-data]",
					Short:          "Send a update-tensor-score tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tensor_id"}, {ProtoField: "dimension"}, {ProtoField: "value"}, {ProtoField: "context"}, {ProtoField: "witness_data"}},
				},
				{
					RpcMethod:      "AddTensorWitness",
					Use:            "add-tensor-witness [tensor-id] [dimension] [witness-lct] [confidence] [evidence-hash]",
					Short:          "Send a add-tensor-witness tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tensor_id"}, {ProtoField: "dimension"}, {ProtoField: "witness_lct"}, {ProtoField: "confidence"}, {ProtoField: "evidence_hash"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
