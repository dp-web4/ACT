package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthMiddleware struct {
	authService *AuthService
	rateLimiter *RateLimiter
	logger      zerolog.Logger
}

func NewAuthMiddleware(authService *AuthService, logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		rateLimiter: NewRateLimiter(logger),
		logger:      logger,
	}
}

func (a *AuthMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Extract API key from header
		apiKey := a.extractAPIKey(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Validate API key
		keyInfo, err := a.authService.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			a.logger.Error().Err(err).Msg("API key validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Check rate limiting
		if !a.rateLimiter.IsAllowed(keyInfo.UserID, keyInfo.RateLimit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Store authentication info in context
		c.Set("authenticated", true)
		c.Set("user_id", keyInfo.UserID)
		c.Set("username", keyInfo.Username)
		c.Set("component_id", keyInfo.ComponentID)
		c.Set("device_type", keyInfo.DeviceType)
		c.Set("permissions", keyInfo.Permissions)
		c.Set("roles", keyInfo.Roles)

		// Continue processing
		c.Next()

		// Log API usage (async)
		go a.logAPIUsage(c, keyInfo, startTime)
	}
}

func (a *AuthMiddleware) extractAPIKey(c *gin.Context) string {
	// Check X-API-Key header first
	if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
		return apiKey
	}

	// Check Authorization header as fallback
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "ApiKey ") {
		return strings.TrimPrefix(authHeader, "ApiKey ")
	}

	return ""
}

func (a *AuthMiddleware) logAPIUsage(c *gin.Context, keyInfo *KeyInfo, startTime time.Time) {
	usage := APIUsage{
		APIKeyID:       keyInfo.UserID, // Use user ID as reference
		Endpoint:       c.Request.URL.Path,
		Method:         c.Request.Method,
		IPAddress:      c.ClientIP(),
		UserAgent:      c.Request.UserAgent(),
		ResponseStatus: c.Writer.Status(),
		ResponseTimeMs: int(time.Since(startTime).Milliseconds()),
	}

	if err := a.authService.LogAPIUsage(c.Request.Context(), usage); err != nil {
		a.logger.Error().Err(err).Msg("Failed to log API usage")
	}
}

// Permission checking middleware
func (a *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		userPermissions, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid permissions format"})
			c.Abort()
			return
		}

		// Check if user has required permission or admin role
		hasPermission := false
		for _, perm := range userPermissions {
			if perm == permission || perm == "admin" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Role checking middleware
func (a *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No roles found"})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid roles format"})
			c.Abort()
			return
		}

		// Check if user has required role or admin role
		hasRole := false
		for _, role := range userRoles {
			if role == requiredRole || role == "admin" {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
