package server

import (
	"context"
	"net/http"
	"time"

	"api-bridge/internal/auth"
	"api-bridge/internal/config"
	grpcServer "api-bridge/internal/grpc"
	"api-bridge/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server represents the API bridge server
type Server struct {
	config         *config.Config
	logger         zerolog.Logger
	router         *gin.Engine
	server         *http.Server
	handler        *handlers.Handler
	grpcServer     *grpcServer.Server
	authMiddleware *auth.AuthMiddleware
	authzService   *auth.AuthorizationService
}

// New creates a new server instance
func New(cfg *config.Config, logger zerolog.Logger) (*Server, error) {
	// Set Gin mode
	if logger.GetLevel() == zerolog.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(logger))

	// Create handler
	handler, err := handlers.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Create gRPC server
	grpcSrv := grpcServer.NewServer(handler.GetBlockchainClient(), cfg)
	grpcSrv.SetLogger(logger)

	// Initialize authentication services if enabled
	var authMiddleware *auth.AuthMiddleware
	var authzService *auth.AuthorizationService

	if cfg.Security.Enabled {
		// Laravel Client
		laravelClient := auth.NewLaravelClient(
			cfg.Security.Laravel.BaseURL,
			cfg.Security.Laravel.APIKey,
			logger,
		)

		// Auth Service
		authService := auth.NewAuthService(laravelClient, logger)

		// Auth Middleware
		authMiddleware = auth.NewAuthMiddleware(authService, logger)

		// Authorization Service
		authzService = auth.NewAuthorizationService(authService, handler.GetBlockchainClient(), logger)

		// Set up gRPC authentication interceptor
		grpcAuthInterceptor := auth.NewGRPCAuthInterceptor(authService, logger)
		grpcSrv.SetAuthInterceptor(grpcAuthInterceptor)
	}

	// Setup routes
	setupRoutes(router, handler, authMiddleware, authzService)

	// Create HTTP server
	server := &http.Server{
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		config:         cfg,
		logger:         logger,
		router:         router,
		server:         server,
		handler:        handler,
		grpcServer:     grpcSrv,
		authMiddleware: authMiddleware,
		authzService:   authzService,
	}, nil
}

// Start starts the REST server
func (s *Server) Start(addr string) error {
	s.server.Addr = addr
	s.logger.Info().Str("addr", addr).Msg("Starting REST server")
	return s.server.ListenAndServe()
}

// StartGRPC starts the gRPC server
func (s *Server) StartGRPC(port int) error {
	return s.grpcServer.Start(port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down server")

	// Shutdown handler (which includes event queue)
	s.handler.Shutdown()

	return s.server.Shutdown(ctx)
}

// setupRoutes configures the API routes with authentication and authorization
func setupRoutes(router *gin.Engine, handler *handlers.Handler, authMiddleware *auth.AuthMiddleware, authzService *auth.AuthorizationService) {
	// Public routes (no authentication required)
	router.GET("/health", handler.HealthCheck)
	router.GET("/blockchain/status", handler.BlockchainStatus)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Apply authentication to all protected routes if enabled
	if authMiddleware != nil {
		v1.Use(authMiddleware.RequireAPIKey())
	}

	{
		// Component Registry endpoints with intelligent authorization
		components := v1.Group("/components")
		{
			// Basic component info - system-level access
			components.GET("/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetComponent)

			// Component registration - system-level access with permission
			components.POST("/register",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("component:register")),
				handler.RegisterComponent)

			components.GET("/:id/identity",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetComponentIdentity)

			components.POST("/:id/verify",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("component:verify")),
				handler.VerifyComponent)

			// Privacy-focused anonymous component endpoints - system-level access
			components.POST("/register-anonymous",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.RegisterAnonymousComponent)

			components.POST("/verify-pairing-hashes",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.VerifyComponentPairingWithHashes)

			components.POST("/authorization-anonymous",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.CreateAnonymousPairingAuthorization)

			components.POST("/revocation-anonymous",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.CreateAnonymousRevocationEvent)

			components.GET("/metadata-anonymous/:hash",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetAnonymousComponentMetadata)
		}

		// Pairing endpoints - requires LCT relationship
		pairing := v1.Group("/pairing")
		{
			pairing.POST("/initiate",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.InitiatePairing)

			pairing.POST("/complete",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.CompletePairing)

			pairing.DELETE("/revoke",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.RevokePairing)

			pairing.GET("/status/:challenge_id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetPairingStatus)
		}

		// LCT Management endpoints - mixed authorization
		lct := v1.Group("/lct")
		{
			// Create LCT - system-level access with permission
			lct.POST("/create",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("lct:create")),
				handler.CreateLCT)

			// Get LCT info - system access or be part of LCT
			lct.GET("/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetLCT)

			// Update LCT status - system access with permission
			lct.PUT("/:id/status",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("lct:write")),
				handler.UpdateLCTStatus)
		}

		// Trust Tensor endpoints - require LCT relationship
		trust := v1.Group("/trust")
		{
			trust.POST("/tensor",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.CreateTrustTensor)

			trust.GET("/tensor/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetTrustTensor)

			trust.PUT("/tensor/:id/score",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.UpdateTrustScore)
		}

		// Energy Cycle endpoints - require LCT relationship + trust score
		energy := v1.Group("/energy")
		{
			energy.POST("/operation",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.CreateEnergyOperation)

			energy.POST("/transfer",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				applyAuthzIfEnabled(authzService, authzService.RequireMinimumTrust(0.7)),
				handler.ExecuteEnergyTransfer)

			energy.GET("/balance/:component_id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetEnergyBalance)
		}

		// Queue Management endpoints - infrastructure access
		queue := v1.Group("/queue")
		{
			queue.POST("/pairing-request",
				applyAuthzIfEnabled(authzService, authzService.RequireInfrastructureAccess()),
				handler.QueuePairingRequest)

			queue.GET("/status/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetQueueStatus)

			queue.POST("/process-offline/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireInfrastructureAccess()),
				handler.ProcessOfflineQueue)

			queue.DELETE("/cancel/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.CancelRequest)

			queue.GET("/requests/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetQueuedRequests)

			queue.GET("/proxy/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireInfrastructureAccess()),
				handler.ListProxyQueue)
		}

		// Authorization Management endpoints - system-level access
		auth := v1.Group("/authorization")
		{
			auth.POST("/pairing",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.CreatePairingAuthorization)

			auth.GET("/component/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetComponentAuthorizations)

			auth.PUT("/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("authorization:write")),
				handler.UpdateAuthorization)

			auth.DELETE("/:id",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("authorization:delete")),
				handler.RevokeAuthorization)

			auth.GET("/check",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.CheckPairingAuthorization)
		}

		// Enhanced Trust Tensor endpoints - require LCT relationship
		trustEnhanced := v1.Group("/trust-enhanced")
		{
			trustEnhanced.POST("/calculate",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.CalculateRelationshipTrust)

			trustEnhanced.GET("/relationship",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetRelationshipTensor)

			trustEnhanced.PUT("/score",
				applyAuthzIfEnabled(authzService, authzService.RequireLCTRelationship()),
				handler.UpdateTensorScore)
		}

		// Account Management endpoints - system-level access
		accounts := v1.Group("/accounts")
		{
			accounts.GET("",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetAccounts)

			accounts.POST("",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				applyAuthIfEnabled(authMiddleware, authMiddleware.RequirePermission("account:create")),
				handler.CreateAccount)

			accounts.GET("/info",
				applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
				handler.GetAccountInfo)
		}

		// Testing endpoints - admin role required
		v1.GET("/test/ignite",
			applyAuthIfEnabled(authMiddleware, authMiddleware.RequireRole("admin")),
			handler.TestIgniteCLI)

		v1.GET("/test/transaction-format",
			applyAuthIfEnabled(authMiddleware, authMiddleware.RequireRole("admin")),
			handler.TestTransactionFormat)

		v1.GET("/test/ignite-help",
			applyAuthIfEnabled(authMiddleware, authMiddleware.RequireRole("admin")),
			handler.GetIgniteHelp)
	}

	// WebSocket endpoint for real-time events - system-level access
	router.GET("/ws",
		applyAuthzIfEnabled(authzService, authzService.RequireSystemAccess()),
		handler.WebSocketHandler)
}

// Helper functions to conditionally apply middleware
func applyAuthIfEnabled(authMiddleware *auth.AuthMiddleware, middleware gin.HandlerFunc) gin.HandlerFunc {
	if authMiddleware != nil {
		return middleware
	}
	return gin.HandlerFunc(func(c *gin.Context) { c.Next() })
}

func applyAuthzIfEnabled(authzService *auth.AuthorizationService, middleware gin.HandlerFunc) gin.HandlerFunc {
	if authzService != nil {
		return middleware
	}
	return gin.HandlerFunc(func(c *gin.Context) { c.Next() })
}

// loggerMiddleware adds logging to Gin requests
func loggerMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Str("user_agent", param.Request.UserAgent()).
			Msg("HTTP request")
		return ""
	})
}
