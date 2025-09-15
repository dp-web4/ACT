package pairing

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/pairing/types"
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
					RpcMethod:      "ValidateBidirectionalAuth",
					Use:            "validate-bidirectional-auth [component-a] [component-b] [context]",
					Short:          "Query validate-bidirectional-auth",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_a"}, {ProtoField: "component_b"}, {ProtoField: "context"}},
				},

				{
					RpcMethod:      "GetPairingStatus",
					Use:            "get-pairing-status [challenge-id]",
					Short:          "Query get-pairing-status",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}},
				},

				{
					RpcMethod:      "ListActivePairings",
					Use:            "list-active-pairings [component-id]",
					Short:          "Query list-active-pairings",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
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
					RpcMethod:      "InitiateBidirectionalPairing",
					Use:            "initiate-bidirectional-pairing [component-a] [component-b] [operational-context] [proxy-id] [force-immediate]",
					Short:          "Send a initiate-bidirectional-pairing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_a"}, {ProtoField: "component_b"}, {ProtoField: "operational_context"}, {ProtoField: "proxy_id"}, {ProtoField: "force_immediate"}},
				},
				{
					RpcMethod:      "CompletePairing",
					Use:            "complete-pairing [challenge-id] [component-a-auth] [component-b-auth] [session-context]",
					Short:          "Send a complete-pairing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}, {ProtoField: "component_a_auth"}, {ProtoField: "component_b_auth"}, {ProtoField: "session_context"}},
				},
				{
					RpcMethod:      "RevokePairing",
					Use:            "revoke-pairing [lct-id] [reason] [notify-offline]",
					Short:          "Send a revoke-pairing tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "lct_id"}, {ProtoField: "reason"}, {ProtoField: "notify_offline"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
