package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gagliardetto/solana-go/rpc"
)

type Config struct {
	NetworkURL   string   `json:"network_url"`
	Wallets      []string `json:"wallets"`
	ScanInterval string   `json:"scan_interval"`
	Alerts       struct {
		MinimumBalance    float64  `json:"minimum_balance"`
		SignificantChange float64  `json:"significant_change"`
		IgnoreTokens      []string `json:"ignore_tokens"`
	} `json:"alerts"`
	Discord struct {
		Enabled    bool   `json:"enabled"`
		WebhookURL string `json:"webhook_url"`
		ChannelID  string `json:"channel_id"`
	} `json:"discord"`
	mu       sync.RWMutex `json:"-"` // Exclude from JSON
	filepath string       `json:"-"` // Exclude from JSON
}

type AlertConfig struct {
	MinimumBalance    uint64   `json:"minimum_balance"`    // Minimum balance to trigger alerts
	SignificantChange float64  `json:"significant_change"` // e.g., 0.20 for 20% change
	IgnoreTokens      []string `json:"ignore_tokens"`      // Tokens to ignore
}

type ScanConfig struct {
	IncludeTokens []string `json:"include_tokens"` // Specific tokens to include (if empty, include all)
	ExcludeTokens []string `json:"exclude_tokens"` // Specific tokens to exclude
	ScanMode      string   `json:"scan_mode"`      // "all", "whitelist", or "blacklist"
}

type DiscordConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	ChannelID  string `json:"channel_id"`
}

func (c *Config) Validate() error {
	if c.NetworkURL == "" {
		return errors.New("network URL is required")
	}
	if len(c.Wallets) == 0 {
		return errors.New("at least one wallet address is required")
	}
	return nil
}

var TestWallets = []string{
	"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", // Known test wallet
	"DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF", // Another test wallet
}

// Test configuration
func GetTestConfig() *Config {
	return &Config{
		NetworkURL:   rpc.DevNet_RPC,
		Wallets:      TestWallets,
		ScanInterval: "5s",
		Alerts: struct {
			MinimumBalance    float64  `json:"minimum_balance"`
			SignificantChange float64  `json:"significant_change"`
			IgnoreTokens      []string `json:"ignore_tokens"`
		}{
			MinimumBalance:    1000,
			SignificantChange: 0.05,
			IgnoreTokens:      []string{},
		},
		Discord: struct {
			Enabled    bool   `json:"enabled"`
			WebhookURL string `json:"webhook_url"`
			ChannelID  string `json:"channel_id"`
		}{
			Enabled:    false,
			WebhookURL: "",
			ChannelID:  "",
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Store the path
	cfg.filepath = path

	log.Printf("Loaded config: NetworkURL=%s, Wallets=%d, ScanInterval=%s",
		cfg.NetworkURL, len(cfg.Wallets), cfg.ScanInterval)

	return &cfg, nil
}

// Update updates the configuration with the provided update request
func (c *Config) Update(update UpdateRequest) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Apply the update
	if err := update.Apply(c); err != nil {
		return err
	}

	// Save the updated configuration
	return c.Save()
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	log.Printf("DEBUG: Save method called")

	// Use absolute path to config file
	path := "config.json"
	if c.filepath != "" {
		path = c.filepath
		log.Printf("DEBUG: Using stored filepath: %s", path)
	} else {
		// Get absolute path of the working directory
		wd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(wd, "config.json")
			log.Printf("DEBUG: Using working directory path: %s", path)
		}
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("ERROR: Failed to marshal config: %v", err)
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Print part of the JSON data
	preview := string(data)
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}
	log.Printf("DEBUG: About to write config data: %s", preview)

	// Write file with explicit error handling
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write config file to %s: %v", path, err)
		return fmt.Errorf("error writing config file: %w", err)
	}

	log.Printf("SUCCESS: Saved configuration to %s", path)

	// Verify the file exists after writing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("ERROR: File does not exist after writing: %s", path)
	} else {
		log.Printf("DEBUG: Confirmed file exists: %s", path)
	}

	return nil
}
