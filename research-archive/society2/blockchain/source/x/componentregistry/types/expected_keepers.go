package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// TrusttensorKeeper defines the expected interface for the Trust Tensor module.
type TrusttensorKeeper interface {
	// Trust calculation methods that component registry needs
	CalculateRelationshipTrust(ctx context.Context, lctId string, operationalContext string) (string, string, error)
	CalculateV3CompositeScore(ctx context.Context, operationID string) (math.LegacyDec, error)
}

// LctmanagerKeeper defines the expected interface for the LCT Manager module.
type LctmanagerKeeper interface {
	// LCT relationship methods that component registry needs
	GetLinkedContextToken(ctx context.Context, lctId string) (LinkedContextToken, bool)
}

// LinkedContextToken defines the LCT structure (simplified for interface)
type LinkedContextToken struct {
	LctId         string
	PairingStatus string
	// Add other fields as needed
}

// ComponentregistryKeeper defines the expected interface for the Component Registry module.
type ComponentregistryKeeper interface {
	// Component verification methods that other modules need
	GetComponentIdentity(ctx context.Context, componentId string) (ComponentIdentity, bool)
	VerifyComponentForPairing(ctx context.Context, componentId string) (bool, string)
	CheckBidirectionalPairingAuth(ctx context.Context, componentA, componentB string) (bool, bool, string)
}
