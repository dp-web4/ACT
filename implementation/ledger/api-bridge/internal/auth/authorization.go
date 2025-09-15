package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthorizationService struct {
	authService *AuthService
	blockchain  interface{} // Using interface{} to avoid circular dependency
	logger      zerolog.Logger
}

type AuthorizationStrategy string

const (
	// Component-to-component authorization (requires LCT relationship)
	StrategyLCTRelationship AuthorizationStrategy = "lct_relationship"

	// System-level authorization (authenticated access)
	StrategySystemLevel AuthorizationStrategy = "system_level"

	// Infrastructure authorization (proxy/gateway components)
	StrategyInfrastructure AuthorizationStrategy = "infrastructure"

	// Public operations (health checks, discovery)
	StrategyPublic AuthorizationStrategy = "public"
)

func NewAuthorizationService(authService *AuthService, blockchain interface{}, logger zerolog.Logger) *AuthorizationService {
	return &AuthorizationService{
		authService: authService,
		blockchain:  blockchain,
		logger:      logger,
	}
}

// Determine appropriate authorization strategy
func (a *AuthorizationService) GetAuthorizationStrategy(deviceType string, operation string) AuthorizationStrategy {
	switch deviceType {
	case "embedded", "sensor", "controller":
		if a.requiresLCTRelationship(operation) {
			return StrategyLCTRelationship
		}
		return StrategySystemLevel

	case "gateway", "proxy":
		return StrategyInfrastructure

	case "web_app", "mobile_app":
		return StrategySystemLevel

	default:
		return StrategyPublic
	}
}

// Check if operation requires LCT relationship
func (a *AuthorizationService) requiresLCTRelationship(operation string) bool {
	lctRequiredOps := []string{
		"energy_transfer",
		"component_pairing",
		"trust_tensor_update",
		"secure_communication",
	}

	for _, op := range lctRequiredOps {
		if op == operation {
			return true
		}
	}
	return false
}

// Middleware for LCT relationship requirement
func (a *AuthorizationService) RequireLCTRelationship() gin.HandlerFunc {
	return func(c *gin.Context) {
		componentID, exists := c.Get("component_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Component ID required for LCT operations"})
			c.Abort()
			return
		}

		sourceComponent := componentID.(string)
		targetComponent := c.Param("target_component")

		if targetComponent == "" {
			c.Next() // Skip if no target component
			return
		}

		// Check LCT relationship via Laravel
		hasRelationship, err := a.authService.CheckComponentAuthorization(c.Request.Context(), sourceComponent, targetComponent)
		if err != nil {
			a.logger.Error().Err(err).Msg("Failed to check LCT relationship")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization check failed"})
			c.Abort()
			return
		}

		if !hasRelationship {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "No LCT relationship found",
				"source": sourceComponent,
				"target": targetComponent,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Middleware for system-level access
func (a *AuthorizationService) RequireSystemAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("authenticated"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		perms := permissions.([]string)
		if !a.hasSystemPermission(perms) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient system permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Middleware for infrastructure access
func (a *AuthorizationService) RequireInfrastructureAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceType, exists := c.Get("device_type")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Device type required"})
			c.Abort()
			return
		}

		if !a.isInfrastructureDevice(deviceType.(string)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized for infrastructure operations"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Check if user has system-level permissions
func (a *AuthorizationService) hasSystemPermission(perms []string) bool {
	for _, perm := range perms {
		if perm == "system:read" || perm == "system:write" || perm == "admin" {
			return true
		}
	}
	return false
}

// Check if device is infrastructure type
func (a *AuthorizationService) isInfrastructureDevice(deviceType string) bool {
	infraDevices := []string{"gateway", "proxy", "controller"}
	for _, dev := range infraDevices {
		if dev == deviceType {
			return true
		}
	}
	return false
}

// Get trust score from blockchain (placeholder - would need blockchain client integration)
func (a *AuthorizationService) GetTrustScore(ctx context.Context, componentA, componentB string) (float64, error) {
	// This is a placeholder - in a real implementation, you would call the blockchain client
	// For now, return a default trust score
	return 0.8, nil
}

// Middleware for minimum trust score requirement
func (a *AuthorizationService) RequireMinimumTrust(minTrustScore float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		componentID, exists := c.Get("component_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Component ID required"})
			c.Abort()
			return
		}

		sourceComponent := componentID.(string)
		targetComponent := c.Param("target_component")

		if targetComponent == "" {
			c.Next() // Skip trust check if no target component
			return
		}

		// Get trust score from blockchain
		trustScore, err := a.GetTrustScore(c.Request.Context(), sourceComponent, targetComponent)
		if err != nil {
			a.logger.Error().Err(err).Msg("Failed to get trust score")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Trust check failed"})
			c.Abort()
			return
		}

		if trustScore < minTrustScore {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Insufficient trust score: %.2f < %.2f", trustScore, minTrustScore),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
