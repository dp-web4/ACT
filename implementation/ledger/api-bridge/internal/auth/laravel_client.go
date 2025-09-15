package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type LaravelClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     zerolog.Logger
}

type LaravelKeyInfo struct {
	Valid bool `json:"valid"`
	User  struct {
		UserID        int      `json:"user_id"`
		Username      string   `json:"username"`
		Email         string   `json:"email"`
		ComponentID   string   `json:"component_id"`
		DeviceType    string   `json:"device_type"`
		Permissions   []string `json:"permissions"`
		RateLimit     int      `json:"rate_limit"`
		MaxConcurrent int      `json:"max_concurrent"`
		Roles         []string `json:"roles"`
	} `json:"user"`
}

type LaravelAuthCheck struct {
	Authorized bool   `json:"authorized"`
	ComponentA string `json:"component_a"`
	ComponentB string `json:"component_b"`
}

type APIUsage struct {
	APIKeyID       int    `json:"api_key_id"`
	Endpoint       string `json:"endpoint"`
	Method         string `json:"method"`
	IPAddress      string `json:"ip_address"`
	UserAgent      string `json:"user_agent"`
	ResponseStatus int    `json:"response_status"`
	ResponseTimeMs int    `json:"response_time_ms"`
}

func NewLaravelClient(baseURL, apiKey string, logger zerolog.Logger) *LaravelClient {
	return &LaravelClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// Validate API key with Laravel backend
func (l *LaravelClient) ValidateAPIKey(ctx context.Context, apiKey string) (*LaravelKeyInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/api/auth/validate-key", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+l.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Laravel validation failed: %d", resp.StatusCode)
	}

	var keyInfo LaravelKeyInfo
	if err := json.NewDecoder(resp.Body).Decode(&keyInfo); err != nil {
		return nil, err
	}

	return &keyInfo, nil
}

// Check component authorization with Laravel
func (l *LaravelClient) CheckComponentAuthorization(ctx context.Context, componentA, componentB string) (bool, error) {
	payload := map[string]string{
		"component_a": componentA,
		"component_b": componentB,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/api/auth/check-authorization", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+l.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Laravel authorization check failed: %d", resp.StatusCode)
	}

	var authCheck LaravelAuthCheck
	if err := json.NewDecoder(resp.Body).Decode(&authCheck); err != nil {
		return false, err
	}

	return authCheck.Authorized, nil
}

// Log API usage to Laravel for monitoring
func (l *LaravelClient) LogAPIUsage(ctx context.Context, usage APIUsage) error {
	jsonPayload, err := json.Marshal(usage)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/api/auth/log-usage", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+l.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Fire and forget - don't block on logging
	go func() {
		resp, err := l.httpClient.Do(req)
		if err != nil {
			l.logger.Error().Err(err).Msg("Failed to log API usage")
			return
		}
		defer resp.Body.Close()
	}()

	return nil
}
