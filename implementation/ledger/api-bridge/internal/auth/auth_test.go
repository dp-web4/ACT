package auth

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthCache(t *testing.T) {
	cache := NewAuthCache(5 * time.Minute)

	// Test setting and getting
	keyInfo := &KeyInfo{
		UserID:      1,
		Username:    "testuser",
		ComponentID: "test-component",
		DeviceType:  "embedded",
		Permissions: []string{"read", "write"},
		RateLimit:   10,
		Roles:       []string{"user"},
	}

	cache.Set("test-key", keyInfo)

	// Test get
	retrieved := cache.Get("test-key")
	assert.NotNil(t, retrieved)
	assert.Equal(t, keyInfo.UserID, retrieved.UserID)
	assert.Equal(t, keyInfo.Username, retrieved.Username)

	// Test non-existent key
	notFound := cache.Get("non-existent")
	assert.Nil(t, notFound)

	// Test delete
	cache.Delete("test-key")
	deleted := cache.Get("test-key")
	assert.Nil(t, deleted)
}

func TestRateLimiter(t *testing.T) {
	logger := zerolog.New(nil)
	limiter := NewRateLimiter(logger)

	// Test first request
	assert.True(t, limiter.IsAllowed(1, 5))

	// Test multiple requests within limit
	for i := 0; i < 4; i++ {
		assert.True(t, limiter.IsAllowed(1, 5))
	}

	// Test rate limit exceeded
	assert.False(t, limiter.IsAllowed(1, 5))

	// Test different user
	assert.True(t, limiter.IsAllowed(2, 3))
}

func TestAuthMiddleware_ExtractAPIKey(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	middleware := NewAuthMiddleware(authService, logger)

	// This is a simple test to ensure the middleware can be created
	assert.NotNil(t, middleware)
}

func TestAuthorizationService_GetAuthorizationStrategy(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	authzService := NewAuthorizationService(authService, nil, logger)

	// Test embedded device with LCT operation
	strategy := authzService.GetAuthorizationStrategy("embedded", "energy_transfer")
	assert.Equal(t, StrategyLCTRelationship, strategy)

	// Test embedded device with non-LCT operation
	strategy = authzService.GetAuthorizationStrategy("embedded", "read_status")
	assert.Equal(t, StrategySystemLevel, strategy)

	// Test gateway device
	strategy = authzService.GetAuthorizationStrategy("gateway", "any_operation")
	assert.Equal(t, StrategyInfrastructure, strategy)

	// Test web app
	strategy = authzService.GetAuthorizationStrategy("web_app", "any_operation")
	assert.Equal(t, StrategySystemLevel, strategy)

	// Test unknown device type
	strategy = authzService.GetAuthorizationStrategy("unknown", "any_operation")
	assert.Equal(t, StrategyPublic, strategy)
}

func TestAuthorizationService_RequiresLCTRelationship(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	authzService := NewAuthorizationService(authService, nil, logger)

	// Test operations that require LCT relationship
	assert.True(t, authzService.requiresLCTRelationship("energy_transfer"))
	assert.True(t, authzService.requiresLCTRelationship("component_pairing"))
	assert.True(t, authzService.requiresLCTRelationship("trust_tensor_update"))
	assert.True(t, authzService.requiresLCTRelationship("secure_communication"))

	// Test operations that don't require LCT relationship
	assert.False(t, authzService.requiresLCTRelationship("read_status"))
	assert.False(t, authzService.requiresLCTRelationship("get_component"))
	assert.False(t, authzService.requiresLCTRelationship("list_accounts"))
}

func TestAuthorizationService_HasSystemPermission(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	authzService := NewAuthorizationService(authService, nil, logger)

	// Test with system permissions
	perms := []string{"system:read", "user:write"}
	assert.True(t, authzService.hasSystemPermission(perms))

	// Test with admin role
	perms = []string{"user:read", "admin"}
	assert.True(t, authzService.hasSystemPermission(perms))

	// Test without system permissions
	perms = []string{"user:read", "user:write"}
	assert.False(t, authzService.hasSystemPermission(perms))
}

func TestAuthorizationService_IsInfrastructureDevice(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	authzService := NewAuthorizationService(authService, nil, logger)

	// Test infrastructure devices
	assert.True(t, authzService.isInfrastructureDevice("gateway"))
	assert.True(t, authzService.isInfrastructureDevice("proxy"))
	assert.True(t, authzService.isInfrastructureDevice("controller"))

	// Test non-infrastructure devices
	assert.False(t, authzService.isInfrastructureDevice("embedded"))
	assert.False(t, authzService.isInfrastructureDevice("sensor"))
	assert.False(t, authzService.isInfrastructureDevice("web_app"))
}

func TestAuthorizationService_GetTrustScore(t *testing.T) {
	logger := zerolog.New(nil)
	authService := &AuthService{} // Mock service
	authzService := NewAuthorizationService(authService, nil, logger)

	// Test trust score retrieval (placeholder implementation)
	score, err := authzService.GetTrustScore(context.Background(), "component-a", "component-b")
	assert.NoError(t, err)
	assert.Equal(t, 0.8, score) // Default trust score from placeholder
}
