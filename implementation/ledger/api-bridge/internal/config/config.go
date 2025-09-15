package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	Server     ServerConfig     `mapstructure:"server"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Events     EventsConfig     `mapstructure:"events"`
	Security   SecurityConfig   `mapstructure:"security"` // New security config
}

// BlockchainConfig holds blockchain connection settings
type BlockchainConfig struct {
	RESTEndpoint string `mapstructure:"rest_endpoint"`
	GRPCEndpoint string `mapstructure:"grpc_endpoint"`
	ChainID      string `mapstructure:"chain_id"`
	Timeout      int    `mapstructure:"timeout"`
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// EventsConfig holds event queue settings
type EventsConfig struct {
	Enabled    bool                `mapstructure:"enabled"`
	MaxRetries int                 `mapstructure:"max_retries"`
	RetryDelay int                 `mapstructure:"retry_delay"`
	QueueSize  int                 `mapstructure:"queue_size"`
	Endpoints  map[string][]string `mapstructure:"endpoints"`
}

// SecurityConfig holds security and authentication settings
type SecurityConfig struct {
	Enabled      bool               `mapstructure:"enabled"`
	Laravel      LaravelConfig      `mapstructure:"laravel"`
	Cache        CacheConfig        `mapstructure:"cache"`
	RateLimiting RateLimitingConfig `mapstructure:"rate_limiting"`
}

// LaravelConfig holds Laravel backend integration settings
type LaravelConfig struct {
	BaseURL   string          `mapstructure:"base_url"`
	APIKey    string          `mapstructure:"api_key"`
	Timeout   time.Duration   `mapstructure:"timeout"`
	Endpoints EndpointsConfig `mapstructure:"endpoints"`
}

// EndpointsConfig holds Laravel API endpoint paths
type EndpointsConfig struct {
	ValidateKey string `mapstructure:"validate_key"`
	CheckAuth   string `mapstructure:"check_authorization"`
	LogUsage    string `mapstructure:"log_usage"`
}

// CacheConfig holds authentication cache settings
type CacheConfig struct {
	Duration        time.Duration `mapstructure:"duration"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// RateLimitingConfig holds rate limiting settings
type RateLimitingConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// Load loads configuration from file
func Load(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Read environment variables
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("blockchain.rest_endpoint", "http://0.0.0.0:1317")
	viper.SetDefault("blockchain.grpc_endpoint", "localhost:9090")
	viper.SetDefault("blockchain.chain_id", "racecarweb")
	viper.SetDefault("blockchain.timeout", 30)

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Event queue defaults - disabled by default
	viper.SetDefault("events.enabled", false)
	viper.SetDefault("events.max_retries", 3)
	viper.SetDefault("events.retry_delay", 5)
	viper.SetDefault("events.queue_size", 1000)
	viper.SetDefault("events.endpoints", map[string][]string{
		"component_registered": {},
		"component_verified":   {},
		"pairing_initiated":    {},
		"pairing_completed":    {},
		"lct_created":          {},
		"trust_tensor_created": {},
		"energy_transfer":      {},
	})

	// Security defaults - disabled by default
	viper.SetDefault("security.enabled", false)
	viper.SetDefault("security.laravel.base_url", "http://localhost:8000")
	viper.SetDefault("security.laravel.api_key", "")
	viper.SetDefault("security.laravel.timeout", "10s")
	viper.SetDefault("security.laravel.endpoints.validate_key", "/api/auth/validate-key")
	viper.SetDefault("security.laravel.endpoints.check_authorization", "/api/auth/check-authorization")
	viper.SetDefault("security.laravel.endpoints.log_usage", "/api/auth/log-usage")
	viper.SetDefault("security.cache.duration", "5m")
	viper.SetDefault("security.cache.cleanup_interval", "1m")
	viper.SetDefault("security.rate_limiting.enabled", true)
	viper.SetDefault("security.rate_limiting.cleanup_interval", "1h")
}

// Save saves configuration to file
func (c *Config) Save(configFile string) error {
	viper.Set("blockchain", c.Blockchain)
	viper.Set("server", c.Server)
	viper.Set("logging", c.Logging)
	viper.Set("events", c.Events)
	viper.Set("security", c.Security)

	return viper.WriteConfigAs(configFile)
}
