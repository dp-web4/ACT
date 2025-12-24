package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"

	lctmanagertypes "racecar-web/x/lctmanager/types"
	"racecar-web/x/trusttensor/types"
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

	// Trust tensor storage
	RelationshipTensors collections.Map[string, types.RelationshipTrustTensor]
	ValueTensors        collections.Map[string, types.ValueTensor]

	bankKeeper       types.BankKeeper
	lctmanagerKeeper lctmanagertypes.LctmanagerKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
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

		bankKeeper:          bankKeeper,
		lctmanagerKeeper:    lctmanagerKeeper,
		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		RelationshipTensors: collections.NewMap(sb, types.RelationshipTrustTensorKey, "relationship_tensors", collections.StringKey, codec.CollValue[types.RelationshipTrustTensor](cdc)),
		ValueTensors:        collections.NewMap(sb, types.ValueTensorKey, "value_tensors", collections.StringKey, codec.CollValue[types.ValueTensor](cdc)),
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

// CalculateRelationshipTrust calculates composite trust score for an LCT relationship
func (k Keeper) CalculateRelationshipTrust(ctx context.Context, lctId, operationalContext string) (string, string, error) {
	// Verify LCT relationship exists (if lctmanagerKeeper is available)
	if k.lctmanagerKeeper != nil {
		lct, found := k.lctmanagerKeeper.GetLinkedContextToken(ctx, lctId)
		if !found {
			// Return default trust score if LCT not found
			return "0.5", "default_trust_no_lct_found", nil
		}

		if lct.PairingStatus != "active" {
			return "0.3", "low_trust_inactive_lct", nil
		}
	}

	// Calculate T3 composite score using the new business logic
	t3Score, err := k.CalculateT3CompositeScore(ctx, lctId)
	if err != nil {
		return "0.5", "error_calculating_t3_score", err
	}

	// Apply context modifier
	contextModifier := k.GetContextModifier(ctx, operationalContext)
	finalScore := t3Score.Mul(contextModifier)

	// Ensure bounds [0, 1]
	if finalScore.GT(math.LegacyOneDec()) {
		finalScore = math.LegacyOneDec()
	}
	if finalScore.IsNegative() {
		finalScore = math.LegacyZeroDec()
	}

	factors := fmt.Sprintf("t3_score_%.3f_context_%s_modifier_%.3f",
		t3Score.MustFloat64(), operationalContext, contextModifier.MustFloat64())

	return finalScore.String(), factors, nil
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return types.Params{}, nil
}
