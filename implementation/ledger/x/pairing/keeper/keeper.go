package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	componentregistrytypes "racecar-web/x/componentregistry/types"
	lctmanagertypes "racecar-web/x/lctmanager/types"
	"racecar-web/x/pairing/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	// Collections for pairing sessions
	PairingSessions collections.Map[string, types.PairingSession]

	bankKeeper              types.BankKeeper
	componentregistryKeeper componentregistrytypes.ComponentregistryKeeper
	pairingqueueKeeper      types.PairingqueueKeeper
	lctmanagerKeeper        lctmanagertypes.LctmanagerKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
	componentregistryKeeper componentregistrytypes.ComponentregistryKeeper,
	pairingqueueKeeper types.PairingqueueKeeper,
	lctmanagerKeeper lctmanagertypes.LctmanagerKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		bankKeeper:              bankKeeper,
		componentregistryKeeper: componentregistryKeeper,
		pairingqueueKeeper:      pairingqueueKeeper,
		lctmanagerKeeper:        lctmanagerKeeper,
		Params:                  collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		PairingSessions:         collections.NewMap(sb, types.PairingSessionPrefix, "pairing_sessions", collections.StringKey, codec.CollValue[types.PairingSession](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return types.Params{}, nil
}
