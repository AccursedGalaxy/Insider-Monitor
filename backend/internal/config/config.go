package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Config represents the application configuration
type Config struct {
	NetworkURL    string                  `json:"network_url"`
	Wallets       []string                `json:"wallets"`
	ScanInterval  string                  `json:"scan_interval"`
	Scan          ScanConfig              `json:"scan"`
	WalletConfigs map[string]WalletConfig `json:"wallet_configs,omitempty"`
	Alerts        AlertConfig             `json:"alerts"`
	Discord       DiscordConfig           `json:"discord"`
	mu            sync.RWMutex            `json:"-"` // Exclude from JSON
	filepath      string                  `json:"-"` // Exclude from JSON
}

// AlertConfig holds alert settings
type AlertConfig struct {
	MinimumBalance    float64  `json:"minimum_balance"`    // Minimum balance to trigger alerts
	SignificantChange float64  `json:"significant_change"` // e.g., 0.20 for 20% change
	IgnoreTokens      []string `json:"ignore_tokens"`      // Tokens to ignore
}

// ScanConfig holds scan settings
type ScanConfig struct {
	IncludeTokens []string `json:"include_tokens"` // Specific tokens to include (if empty, include all)
	ExcludeTokens []string `json:"exclude_tokens"` // Specific tokens to exclude
	ScanMode      string   `json:"scan_mode"`      // "all", "whitelist", or "blacklist"
}

// DiscordConfig holds Discord integration settings
type DiscordConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	ChannelID  string `json:"channel_id"`
}

// WalletConfig holds wallet-specific configuration
type WalletConfig struct {
	Scan *ScanConfig `json:"scan,omitempty"` // Wallet-specific scan configuration
}

// UpdateRequest represents a request to update the configuration
type UpdateRequest struct {
	NetworkURL    *string                 `json:"network_url,omitempty"`
	Wallets       []string                `json:"wallets,omitempty"`
	ScanInterval  *string                 `json:"scan_interval,omitempty"`
	Scan          *ScanConfig             `json:"scan,omitempty"`
	WalletConfigs map[string]WalletConfig `json:"wallet_configs,omitempty"`
	Alerts        *AlertConfig            `json:"alerts,omitempty"`
	Discord       *DiscordConfig          `json:"discord,omitempty"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.NetworkURL == "" {
		return errors.New("network URL is required")
	}
	if len(c.Wallets) == 0 {
		return errors.New("at least one wallet address is required")
	}

	// Validate scan mode
	if c.Scan.ScanMode != "" && c.Scan.ScanMode != "all" && c.Scan.ScanMode != "whitelist" && c.Scan.ScanMode != "blacklist" {
		return fmt.Errorf("invalid scan mode: %s (must be 'all', 'whitelist', or 'blacklist')", c.Scan.ScanMode)
	}

	// Default to "all" if not set
	if c.Scan.ScanMode == "" {
		c.Scan.ScanMode = "all"
	}

	// Validate per-wallet configs
	for wallet, walletCfg := range c.WalletConfigs {
		if walletCfg.Scan != nil {
			mode := walletCfg.Scan.ScanMode
			if mode != "" && mode != "all" && mode != "whitelist" && mode != "blacklist" {
				return fmt.Errorf("invalid scan mode for wallet %s: %s (must be 'all', 'whitelist', or 'blacklist')", wallet, mode)
			}

			// Default to "all" if not set
			if walletCfg.Scan.ScanMode == "" {
				walletCfg.Scan.ScanMode = "all"
			}
		}
	}

	return nil
}

// GetTestConfig returns a test configuration
func GetTestConfig() *Config {
	return &Config{
		NetworkURL:   "https://api.mainnet-beta.solana.com",
		Wallets:      []string{"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"},
		ScanInterval: "1m",
		Scan: ScanConfig{
			ScanMode:      "all",
			IncludeTokens: []string{},
			ExcludeTokens: []string{},
		},
		WalletConfigs: map[string]WalletConfig{},
		Alerts: AlertConfig{
			MinimumBalance:    1000,
			SignificantChange: 0.05,
			IgnoreTokens:      []string{},
		},
		Discord: DiscordConfig{
			Enabled:    false,
			WebhookURL: "",
			ChannelID:  "",
		},
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If not, create a default config
		config := GetTestConfig()
		config.filepath = path

		// Ensure the directory exists
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		// Save the default config
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}

		return config, nil
	}

	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set the filepath
	config.filepath = path

	// Validate the config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// Update updates the configuration from an update request
func (c *Config) Update(update UpdateRequest) error {
	c.mu.Lock()

	// Update fields if provided
	if update.NetworkURL != nil {
		c.NetworkURL = *update.NetworkURL
	}
	if update.Wallets != nil {
		c.Wallets = update.Wallets
	}
	if update.ScanInterval != nil {
		c.ScanInterval = *update.ScanInterval
	}
	if update.Scan != nil {
		c.Scan = *update.Scan
	}
	if update.WalletConfigs != nil {
		// Merge with existing wallet configs
		if c.WalletConfigs == nil {
			c.WalletConfigs = make(map[string]WalletConfig)
		}
		for wallet, cfg := range update.WalletConfigs {
			c.WalletConfigs[wallet] = cfg
		}
	}
	if update.Alerts != nil {
		c.Alerts = *update.Alerts
	}
	if update.Discord != nil {
		c.Discord = *update.Discord
	}

	// Validate the updated config
	var err error
	if err = c.Validate(); err != nil {
		c.mu.Unlock()
		return err
	}

	// Release the lock before calling Save() to avoid deadlock
	c.mu.Unlock()

	// Save the changes
	if err = c.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Save saves the configuration to a file
func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure we have a filepath
	if c.filepath == "" {
		return errors.New("no filepath specified for saving config")
	}

	// Ensure the directory exists
	dir := filepath.Dir(c.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	log.Printf("Configuration saved to %s", c.filepath)
	return nil
}

// GetScanConfigForWallet returns the scan configuration for a specific wallet
func (c *Config) GetScanConfigForWallet(walletAddress string) ScanConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check if there's a wallet-specific config
	if wallet, exists := c.WalletConfigs[walletAddress]; exists && wallet.Scan != nil {
		return *wallet.Scan
	}

	// Otherwise, return the global scan config
	return c.Scan
}
