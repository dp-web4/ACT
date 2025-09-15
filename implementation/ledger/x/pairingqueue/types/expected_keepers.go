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

// TrustTensorKeeper defines the expected interface for the Trust Tensor module.
type TrustTensorKeeper interface {
	GetRelationshipTensor(ctx context.Context, componentA, componentB string) (RelationshipTensor, bool)
	UpdateTensorScore(ctx context.Context, componentA, componentB string, score math.LegacyDec, evidence string) error
}

// LCTManagerKeeper defines the expected interface for the LCT Manager module.
type LCTManagerKeeper interface {
	GetLCTRelationship(ctx context.Context, lctID string) (LCTRelationship, bool)
	ValidateCryptographicIntegrity(ctx context.Context, lctID string, message []byte, signature []byte) (bool, error)
}

// ComponentRegistryKeeper defines the expected interface for the Component Registry module.
type ComponentRegistryKeeper interface {
	GetComponentIdentity(ctx context.Context, componentID string) (ComponentIdentity, bool)
	CheckPairingAuthorization(ctx context.Context, initiatorID, targetID string) (bool, string, error)
}

// RelationshipTensor represents a trust tensor relationship
type RelationshipTensor struct {
	ComponentA    string         `json:"component_a"`
	ComponentB    string         `json:"component_b"`
	TensorScore   math.LegacyDec `json:"tensor_score"`
	EvidenceCount int64          `json:"evidence_count"`
	LastUpdated   int64          `json:"last_updated"`
}

// LCTRelationship represents an LCT relationship
type LCTRelationship struct {
	LctId       string `json:"lct_id"`
	ComponentA  string `json:"component_a"`
	ComponentB  string `json:"component_b"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	LastUpdated int64  `json:"last_updated"`
}

// ComponentIdentity represents a component identity
type ComponentIdentity struct {
	ComponentId   string `json:"component_id"`
	ComponentType string `json:"component_type"`
	Owner         string `json:"owner"`
	IsVerified    bool   `json:"is_verified"`
	CreatedAt     int64  `json:"created_at"`
}
