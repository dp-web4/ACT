package pairingqueue

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/pairingqueue/types"
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
					RpcMethod:      "GetQueuedRequests",
					Use:            "get-queued-requests [component-id]",
					Short:          "Query get-queued-requests",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},

				{
					RpcMethod:      "GetRequestStatus",
					Use:            "get-request-status [request-id]",
					Short:          "Query get-request-status",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "request_id"}},
				},

				{
					RpcMethod:      "ListProxyQueue",
					Use:            "list-proxy-queue [proxy-id]",
					Short:          "Query list-proxy-queue",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "proxy_id"}},
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
					RpcMethod:      "QueuePairingRequest",
					Use:            "queue-pairing-request [initiator-id] [target-id] [request-type] [proxy-id]",
					Short:          "Send a queue-pairing-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "initiator_id"}, {ProtoField: "target_id"}, {ProtoField: "request_type"}, {ProtoField: "proxy_id"}},
				},
				{
					RpcMethod:      "ProcessOfflineQueue",
					Use:            "process-offline-queue [component-id]",
					Short:          "Send a process-offline-queue tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},
				{
					RpcMethod:      "CancelRequest",
					Use:            "cancel-request [request-id] [reason]",
					Short:          "Send a cancel-request tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "request_id"}, {ProtoField: "reason"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
