package hardwarebinding

import (
    "context"
    "encoding/json"
    "fmt"

    "cosmossdk.io/core/appmodule"
    "cosmossdk.io/core/registry"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"

    "society4chain/x/hardwarebinding/keeper"
    "society4chain/x/hardwarebinding/types"
)

var (
    _ module.AppModule      = AppModule{}
    _ module.HasGenesis     = AppModule{}
    _ module.HasServices    = AppModule{}
    _ module.HasConsensusVersion = AppModule{}
    _ appmodule.AppModule   = AppModule{}
)

// ConsensusVersion defines the current module consensus version
const ConsensusVersion = 1

type AppModule struct {
    cdc    codec.Codec
    keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
    return AppModule{
        cdc:    cdc,
        keeper: keeper,
    }
}

// Name returns the module's name
func (AppModule) Name() string {
    return types.ModuleName
}

// RegisterLegacyAminoCodec registers the module's types on the LegacyAmino codec
func (AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
    types.RegisterLegacyAminoCodec(cdc)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
    // No gRPC gateway routes for this module
}

// RegisterInterfaces registers interfaces and implementations
func (AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
    types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
    return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation
func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
    var data types.GenesisState
    if err := cdc.UnmarshalJSON(bz, &data); err != nil {
        return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
    }
    return types.ValidateGenesis(data)
}

// InitGenesis performs genesis initialization
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
    var genesisState types.GenesisState
    cdc.MustUnmarshalJSON(data, &genesisState)
    am.keeper.InitGenesis(ctx, genesisState)
}

// ExportGenesis returns the exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
    gs := am.keeper.ExportGenesis(ctx)
    return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements HasConsensusVersion
func (AppModule) ConsensusVersion() uint64 {
    return ConsensusVersion
}

// RegisterServices registers module services
func (am AppModule) RegisterServices(registrar registry.ServiceRegistrar) error {
    types.RegisterMsgServer(registrar.MsgServiceRouter(), keeper.NewMsgServerImpl(am.keeper))
    types.RegisterQueryServer(registrar.QueryServiceRouter(), am.keeper)
    return nil
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface
func (am AppModule) IsAppModule() {}

// BeginBlock performs begin block functionality
func (am AppModule) BeginBlock(ctx context.Context) error {
    sdkCtx := sdk.UnwrapSDKContext(ctx)
    return am.keeper.BeginBlocker(sdkCtx)
}

// EndBlock performs end block functionality
func (am AppModule) EndBlock(ctx context.Context) error {
    sdkCtx := sdk.UnwrapSDKContext(ctx)
    return am.keeper.EndBlocker(sdkCtx)
}