package pairingqueue

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"

	componentregistrytypes "racecar-web/x/componentregistry/types"
	"racecar-web/x/pairingqueue/keeper"
	"racecar-web/x/pairingqueue/types"
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

	PairingqueueKeeper keeper.Keeper
	Module             appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AuthKeeper,
		in.BankKeeper,
		in.ComponentregistryKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{PairingqueueKeeper: k, Module: m}
}
