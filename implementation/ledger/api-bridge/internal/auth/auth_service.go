package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type AuthService struct {
	laravelClient *LaravelClient
	cache         *AuthCache
	logger        zerolog.Logger
}

func NewAuthService(laravelClient *LaravelClient, logger zerolog.Logger) *AuthService {
	return &AuthService{
		laravelClient: laravelClient,
		cache:         NewAuthCache(5 * time.Minute), // 5 minute cache
		logger:        logger,
	}
}

// Validate API key with caching
func (a *AuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*KeyInfo, error) {
	// Check cache first
	if cached := a.cache.Get(apiKey); cached != nil {
		return cached, nil
	}

	// Validate with Laravel
	laravelInfo, err := a.laravelClient.ValidateAPIKey(ctx, apiKey)
	if err != nil {
		a.logger.Error().Err(err).Msg("Laravel API key validation failed")
		return nil, err
	}

	if !laravelInfo.Valid {
		return nil, fmt.Errorf("invalid API key")
	}

	// Convert to internal format
	keyInfo := &KeyInfo{
		UserID:        laravelInfo.User.UserID,
		Username:      laravelInfo.User.Username,
		ComponentID:   laravelInfo.User.ComponentID,
		DeviceType:    laravelInfo.User.DeviceType,
		Permissions:   laravelInfo.User.Permissions,
		RateLimit:     laravelInfo.User.RateLimit,
		MaxConcurrent: laravelInfo.User.MaxConcurrent,
		Roles:         laravelInfo.User.Roles,
	}

	// Cache the result
	a.cache.Set(apiKey, keyInfo)

	return keyInfo, nil
}

// Check component authorization
func (a *AuthService) CheckComponentAuthorization(ctx context.Context, componentA, componentB string) (bool, error) {
	return a.laravelClient.CheckComponentAuthorization(ctx, componentA, componentB)
}

// Log API usage
func (a *AuthService) LogAPIUsage(ctx context.Context, usage APIUsage) error {
	return a.laravelClient.LogAPIUsage(ctx, usage)
}

// Revoke API key from cache
func (a *AuthService) RevokeAPIKey(apiKey string) {
	a.cache.Delete(apiKey)
}
