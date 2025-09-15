package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"racecar-web/x/energycycle/types"
	lctmanagertypes "racecar-web/x/lctmanager/types"
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

	// Energy cycle storage
	EnergyOperations      collections.Map[string, types.EnergyOperation]
	RelationshipAtpTokens collections.Map[string, types.RelationshipAtpToken]
	RelationshipAdpTokens collections.Map[string, types.RelationshipAdpToken]

	bankKeeper        types.BankKeeper
	lctmanagerKeeper  lctmanagertypes.LctmanagerKeeper
	trusttensorKeeper types.TrusttensorKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
	lctmanagerKeeper lctmanagertypes.LctmanagerKeeper,
	trusttensorKeeper types.TrusttensorKeeper,
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

		bankKeeper:            bankKeeper,
		lctmanagerKeeper:      lctmanagerKeeper,
		trusttensorKeeper:     trusttensorKeeper,
		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		EnergyOperations:      collections.NewMap(sb, types.EnergyOperationKey, "energy_operations", collections.StringKey, codec.CollValue[types.EnergyOperation](cdc)),
		RelationshipAtpTokens: collections.NewMap(sb, types.RelationshipAtpTokenKey, "relationship_atp_tokens", collections.StringKey, codec.CollValue[types.RelationshipAtpToken](cdc)),
		RelationshipAdpTokens: collections.NewMap(sb, types.RelationshipAdpTokenKey, "relationship_adp_tokens", collections.StringKey, codec.CollValue[types.RelationshipAdpToken](cdc)),
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
