package componentregistry

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"racecar-web/x/componentregistry/types"
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
					RpcMethod:      "GetComponent",
					Use:            "get-component [component-id]",
					Short:          "Query get-component",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},
				{
					RpcMethod:      "GetComponentVerification",
					Use:            "get-component-verification [component-id]",
					Short:          "Query component verification status",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},
				{
					RpcMethod:      "CheckPairingAuth",
					Use:            "check-pairing-auth [component-a] [component-b]",
					Short:          "Query check-pairing-auth",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_a"}, {ProtoField: "component_b"}},
				},
				{
					RpcMethod:      "ListAuthorizedPartners",
					Use:            "list-authorized-partners [component-id]",
					Short:          "Query list-authorized-partners",
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
					RpcMethod:      "RegisterComponent",
					Use:            "register-component [component-id] [component-type] [manufacturer-data]",
					Short:          "Send a register-component tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}, {ProtoField: "component_type"}, {ProtoField: "manufacturer_data"}},
				},
				{
					RpcMethod:      "UpdateAuthorization",
					Use:            "update-authorization [component-id] [auth-rules]",
					Short:          "Send a update-authorization tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}, {ProtoField: "auth_rules"}},
				},
				{
					RpcMethod:      "VerifyComponent",
					Use:            "verify-component [component-id]",
					Short:          "Send a verify-component tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "component_id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
