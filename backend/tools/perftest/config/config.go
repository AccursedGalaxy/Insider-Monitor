package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the performance test configuration
type Config struct {
	// General settings
	BaseURL      string        `json:"base_url"`
	Concurrency  int           `json:"concurrency"`
	Duration     time.Duration `json:"duration"`
	RampUpPeriod time.Duration `json:"ramp_up_period"`
	Verbose      bool          `json:"verbose"`

	// API test settings
	APIEndpoints []Endpoint `json:"api_endpoints"`

	// WebSocket test settings
	WebSocketURL string   `json:"websocket_url"`
	MessageTypes []string `json:"message_types"`
	MessageRate  float64  `json:"message_rate"`

	// Authentication
	AuthEnabled bool   `json:"auth_enabled"`
	AuthType    string `json:"auth_type"` // "basic", "jwt", "api_key"
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
	JWTToken    string `json:"jwt_token,omitempty"`
}

// Endpoint represents an API endpoint to test
type Endpoint struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         interface{}       `json:"body,omitempty"`
	ExpectedCode int               `json:"expected_code"`
	Weight       int               `json:"weight"` // For weighted random selection of endpoints
	RequiresAuth bool              `json:"requires_auth"`
}

// CustomConfig is used for custom JSON unmarshaling
type CustomConfig struct {
	// General settings
	BaseURL      string `json:"base_url"`
	Concurrency  int    `json:"concurrency"`
	Duration     string `json:"duration"`
	RampUpPeriod string `json:"ramp_up_period"`
	Verbose      bool   `json:"verbose"`

	// API test settings
	APIEndpoints []Endpoint `json:"api_endpoints"`

	// WebSocket test settings
	WebSocketURL string   `json:"websocket_url"`
	MessageTypes []string `json:"message_types"`
	MessageRate  float64  `json:"message_rate"`

	// Authentication
	AuthEnabled bool   `json:"auth_enabled"`
	AuthType    string `json:"auth_type"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
	JWTToken    string `json:"jwt_token,omitempty"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL:      "http://localhost:8080",
		Concurrency:  10,
		Duration:     30 * time.Second,
		RampUpPeriod: 5 * time.Second,
		Verbose:      false,
		WebSocketURL: "ws://localhost:8080/ws",
		MessageRate:  1.0, // Messages per second per connection
		APIEndpoints: []Endpoint{
			{
				URL:          "/api/v1/wallets",
				Method:       "GET",
				ExpectedCode: 200,
				Weight:       5,
				RequiresAuth: false,
			},
			{
				URL:          "/api/v1/status",
				Method:       "GET",
				ExpectedCode: 200,
				Weight:       2,
				RequiresAuth: false,
			},
			{
				URL:          "/health",
				Method:       "GET",
				ExpectedCode: 200,
				Weight:       1,
				RequiresAuth: false,
			},
		},
		MessageTypes: []string{"subscribe", "unsubscribe"},
		AuthEnabled:  false,
		AuthType:     "none",
	}
}

// LoadConfig loads a configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	// If file doesn't exist, create default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		return cfg, saveConfig(cfg, path)
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the configuration using custom config
	var customCfg CustomConfig
	if err := json.Unmarshal(data, &customCfg); err != nil {
		return nil, err
	}

	// Create the real config from the custom config
	cfg := &Config{
		BaseURL:      customCfg.BaseURL,
		Concurrency:  customCfg.Concurrency,
		Verbose:      customCfg.Verbose,
		APIEndpoints: customCfg.APIEndpoints,
		WebSocketURL: customCfg.WebSocketURL,
		MessageTypes: customCfg.MessageTypes,
		MessageRate:  customCfg.MessageRate,
		AuthEnabled:  customCfg.AuthEnabled,
		AuthType:     customCfg.AuthType,
		Username:     customCfg.Username,
		Password:     customCfg.Password,
		APIKey:       customCfg.APIKey,
		JWTToken:     customCfg.JWTToken,
	}

	// Parse durations
	if customCfg.Duration != "" {
		duration, err := time.ParseDuration(customCfg.Duration)
		if err != nil {
			return nil, err
		}
		cfg.Duration = duration
	} else {
		cfg.Duration = 30 * time.Second
	}

	if customCfg.RampUpPeriod != "" {
		rampUp, err := time.ParseDuration(customCfg.RampUpPeriod)
		if err != nil {
			return nil, err
		}
		cfg.RampUpPeriod = rampUp
	} else {
		cfg.RampUpPeriod = 5 * time.Second
	}

	return cfg, nil
}

// saveConfig saves the configuration to a JSON file
func saveConfig(cfg *Config, path string) error {
	// Convert Config to CustomConfig for JSON serialization
	customCfg := CustomConfig{
		BaseURL:      cfg.BaseURL,
		Concurrency:  cfg.Concurrency,
		Duration:     cfg.Duration.String(),
		RampUpPeriod: cfg.RampUpPeriod.String(),
		Verbose:      cfg.Verbose,
		APIEndpoints: cfg.APIEndpoints,
		WebSocketURL: cfg.WebSocketURL,
		MessageTypes: cfg.MessageTypes,
		MessageRate:  cfg.MessageRate,
		AuthEnabled:  cfg.AuthEnabled,
		AuthType:     cfg.AuthType,
		Username:     cfg.Username,
		Password:     cfg.Password,
		APIKey:       cfg.APIKey,
		JWTToken:     cfg.JWTToken,
	}

	data, err := json.MarshalIndent(customCfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
