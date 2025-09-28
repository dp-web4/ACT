package lctmanager

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/lctmanager/types"
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
					RpcMethod:      "GetLct",
					Use:            "get-lct [lct-id]",
					Short:          "Query get-lct",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}},
				},

				{
					RpcMethod:      "GetComponentRelationships",
					Use:            "get-component-relationships [component-id]",
					Short:          "Query get-component-relationships",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},

				{
					RpcMethod:      "ValidateLctAccess",
					Use:            "validate-lct-access [lct-id] [requestor-id]",
					Short:          "Query validate-lct-access",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "requestor_id"}},
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
					RpcMethod:      "CreateLctRelationship",
					Use:            "create-lct-relationship [component-a] [component-b] [context] [proxy-id]",
					Short:          "Send a create-lct-relationship tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_a"}, {ProtoField: "component_b"}, {ProtoField: "context"}, {ProtoField: "proxy_id"}},
				},
				{
					RpcMethod:      "UpdateLctStatus",
					Use:            "update-lct-status [lct-id] [new-status] [reason]",
					Short:          "Send a update-lct-status tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "new_status"}, {ProtoField: "reason"}},
				},
				{
					RpcMethod:      "TerminateLctRelationship",
					Use:            "terminate-lct-relationship [lct-id] [reason] [notify-offline]",
					Short:          "Send a terminate-lct-relationship tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "reason"}, {ProtoField: "notify_offline"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
