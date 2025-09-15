package types

import (
	"context"
	lctmanagertypes "racecar-web/x/lctmanager/types"
	pairingqueuetypes "racecar-web/x/pairingqueue/types"

	"cosmossdk.io/core/address"
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

// PairingqueueKeeper defines the expected interface for the Pairing Queue module.
type PairingqueueKeeper interface {
	// Queue management methods that pairing needs
	GetPairingRequest(ctx context.Context, requestId string) (pairingqueuetypes.PairingRequest, bool)
	QueueOfflineOperation(ctx context.Context, componentId, operationType string) (string, error)
	ProcessOfflineQueue(ctx context.Context, componentId string) (int, int, error)
}

// Use the LctmanagerKeeper interface from the lctmanager module
type LctmanagerKeeper = lctmanagertypes.LctmanagerKeeper
