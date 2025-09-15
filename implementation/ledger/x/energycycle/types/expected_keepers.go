package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	SetAccount(context.Context, sdk.AccountI)
	NewAccountWithAddress(context.Context, sdk.AccAddress) sdk.AccountI
}

// AuthKeeper defines the expected interface for the Auth module (alias for AccountKeeper).
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	SetAccount(context.Context, sdk.AccountI)
	NewAccountWithAddress(context.Context, sdk.AccAddress) sdk.AccountI
}

// TrusttensorKeeper defines the expected interface for the Trust Tensor module.
type TrusttensorKeeper interface {
	// Trust calculation methods that energycycle needs
	CalculateRelationshipTrust(ctx context.Context, lctId string, operationalContext string) (string, string, error)
	CalculateV3CompositeScore(ctx context.Context, operationID string) (math.LegacyDec, error)
}
