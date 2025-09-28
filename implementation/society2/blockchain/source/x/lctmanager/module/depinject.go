package lctmanager

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	componentregistrytypes "racecar-web/x/componentregistry/types"
	"racecar-web/x/lctmanager/keeper"
	"racecar-web/x/lctmanager/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper              types.AuthKeeper
	BankKeeper              types.BankKeeper
	ComponentregistryKeeper componentregistrytypes.ComponentregistryKeeper
}

type ModuleOutputs struct {
	depinject.Out

	LctmanagerKeeper types.LctmanagerKeeper
	Module           appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority.Bytes(),
		in.BankKeeper,
		in.ComponentregistryKeeper,
		nil, // pairingqueueKeeper - removed to break circular dependency
		nil, // logger - will be set by the module
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{LctmanagerKeeper: k, Module: m}
}
