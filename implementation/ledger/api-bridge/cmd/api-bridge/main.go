package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-bridge/internal/config"
	"api-bridge/internal/server"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	configFile string
	restPort   int
	grpcPort   int
	logLevel   string
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "api-bridge",
		Short: "Web4 Race Car Battery Management API Bridge",
		Long: `API Bridge service for Web4 Race Car Battery Management System.
Provides REST and WebSocket APIs to interact with the blockchain.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(logger)
		},
	}

	// Add flags
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Configuration file path")
	rootCmd.Flags().IntVarP(&restPort, "rest-port", "p", 8080, "REST server port")
	rootCmd.Flags().IntVarP(&grpcPort, "grpc-port", "g", 9090, "gRPC server port")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to execute command")
	}
}

func runServer(logger zerolog.Logger) error {
	// Set log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger = logger.Level(level)

	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create server
	srv, err := server.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Start REST server
	go func() {
		logger.Info().Int("port", restPort).Msg("Starting REST API bridge server")
		if err := srv.Start(fmt.Sprintf(":%d", restPort)); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("REST server failed to start")
		}
	}()

	// Start gRPC server
	go func() {
		logger.Info().Int("port", grpcPort).Msg("Starting gRPC API bridge server")
		if err := srv.StartGRPC(grpcPort); err != nil {
			logger.Fatal().Err(err).Msg("gRPC server failed to start")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
		return err
	}

	logger.Info().Msg("Server exited")
	return nil
}
