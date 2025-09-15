package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog"
)

type GRPCAuthInterceptor struct {
	authService *AuthService
	logger      zerolog.Logger
}

func NewGRPCAuthInterceptor(authService *AuthService, logger zerolog.Logger) *GRPCAuthInterceptor {
	return &GRPCAuthInterceptor{
		authService: authService,
		logger:      logger,
	}
}

func (g *GRPCAuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Skip authentication for health checks
	if strings.HasSuffix(info.FullMethod, "/Health") {
		return handler(ctx, req)
	}

	// Extract API key from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata required")
	}

	var apiKey string

	// Check for API key in metadata
	if keys := md["x-api-key"]; len(keys) > 0 {
		apiKey = keys[0]
	} else if auth := md["authorization"]; len(auth) > 0 {
		if strings.HasPrefix(auth[0], "ApiKey ") {
			apiKey = strings.TrimPrefix(auth[0], "ApiKey ")
		}
	}

	if apiKey == "" {
		return nil, status.Errorf(codes.Unauthenticated, "API key required")
	}

	// Validate API key
	keyInfo, err := g.authService.ValidateAPIKey(ctx, apiKey)
	if err != nil {
		g.logger.Error().Err(err).Msg("gRPC API key validation failed")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API key")
	}

	// Add authentication info to context
	ctx = context.WithValue(ctx, "authenticated", true)
	ctx = context.WithValue(ctx, "user_id", keyInfo.UserID)
	ctx = context.WithValue(ctx, "username", keyInfo.Username)
	ctx = context.WithValue(ctx, "component_id", keyInfo.ComponentID)
	ctx = context.WithValue(ctx, "device_type", keyInfo.DeviceType)
	ctx = context.WithValue(ctx, "permissions", keyInfo.Permissions)
	ctx = context.WithValue(ctx, "roles", keyInfo.Roles)

	return handler(ctx, req)
}
