package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGetTestConfig(t *testing.T) {
	config := GetTestConfig()

	if config == nil {
		t.Fatal("Expected test config to be non-nil")
	}

	if config.NetworkURL == "" {
		t.Error("Expected NetworkURL to be non-empty")
	}

	if len(config.Wallets) == 0 {
		t.Error("Expected Wallets to have at least one entry")
	}
}

func TestLoadConfig_NonExistent(t *testing.T) {
	// Create a temporary file path
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Get the path and close/remove the file
	path := tempFile.Name()
	tempFile.Close()
	os.Remove(path)

	// Load the config - should create a default one
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	defer os.Remove(path) // Clean up

	if config == nil {
		t.Fatal("Expected config to be non-nil")
	}

	// Check that the file was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(*Config)
		expectError bool
	}{
		{
			name:        "Valid config",
			modify:      func(c *Config) {},
			expectError: false,
		},
		{
			name: "Missing network URL",
			modify: func(c *Config) {
				c.NetworkURL = ""
			},
			expectError: true,
		},
		{
			name: "No wallets",
			modify: func(c *Config) {
				c.Wallets = []string{}
			},
			expectError: true,
		},
		{
			name: "Invalid scan mode",
			modify: func(c *Config) {
				c.Scan.ScanMode = "invalid"
			},
			expectError: true,
		},
		{
			name: "Invalid wallet scan mode",
			modify: func(c *Config) {
				if c.WalletConfigs == nil {
					c.WalletConfigs = make(map[string]WalletConfig)
				}
				c.WalletConfigs["test"] = WalletConfig{
					Scan: &ScanConfig{
						ScanMode: "invalid",
					},
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetTestConfig()
			tt.modify(config)

			err := config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected validation error but got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestConfigUpdate(t *testing.T) {
	// Create a temporary file for the config
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Get the path
	path := tempFile.Name()

	// Write a default config to the file
	defaultConfig := GetTestConfig()
	configData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	_, err = tempFile.Write(configData)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tempFile.Close()

	defer os.Remove(path) // Clean up

	// Load the config
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test updating various fields
	newURL := "https://api.testnet.solana.com"
	newInterval := "5m"
	newWallets := []string{"newwallet1", "newwallet2"}

	update := UpdateRequest{
		NetworkURL:   &newURL,
		ScanInterval: &newInterval,
		Wallets:      newWallets,
		Scan: &ScanConfig{
			ScanMode:      "whitelist",
			IncludeTokens: []string{"token1", "token2"},
		},
		Alerts: &AlertConfig{
			MinimumBalance:    2000,
			SignificantChange: 0.1,
			IgnoreTokens:      []string{"ignoreme"},
		},
		Discord: &DiscordConfig{
			Enabled:    true,
			WebhookURL: "https://discord.com/api/webhooks/test",
			ChannelID:  "testchannel",
		},
	}

	// Update the config
	err = config.Update(update)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify the updates were applied
	if config.NetworkURL != newURL {
		t.Errorf("Expected NetworkURL to be %s, got %s", newURL, config.NetworkURL)
	}

	if config.ScanInterval != newInterval {
		t.Errorf("Expected ScanInterval to be %s, got %s", newInterval, config.ScanInterval)
	}

	if len(config.Wallets) != len(newWallets) {
		t.Errorf("Expected %d wallets, got %d", len(newWallets), len(config.Wallets))
	}

	if config.Scan.ScanMode != "whitelist" {
		t.Errorf("Expected ScanMode to be whitelist, got %s", config.Scan.ScanMode)
	}

	if len(config.Scan.IncludeTokens) != 2 {
		t.Errorf("Expected 2 included tokens, got %d", len(config.Scan.IncludeTokens))
	}

	if config.Alerts.MinimumBalance != 2000 {
		t.Errorf("Expected MinimumBalance to be 2000, got %f", config.Alerts.MinimumBalance)
	}

	if config.Discord.Enabled != true {
		t.Errorf("Expected Discord.Enabled to be true")
	}

	// Test loading the config again to verify it was saved
	loadedConfig, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if loadedConfig.NetworkURL != newURL {
		t.Errorf("Loaded config: Expected NetworkURL to be %s, got %s", newURL, loadedConfig.NetworkURL)
	}
}

func TestGetScanConfigForWallet(t *testing.T) {
	config := GetTestConfig()

	// Add a wallet-specific config
	if config.WalletConfigs == nil {
		config.WalletConfigs = make(map[string]WalletConfig)
	}

	testWallet := "testwallet"
	config.WalletConfigs[testWallet] = WalletConfig{
		Scan: &ScanConfig{
			ScanMode:      "blacklist",
			ExcludeTokens: []string{"exclude1", "exclude2"},
		},
	}

	// Test getting wallet-specific config
	walletConfig := config.GetScanConfigForWallet(testWallet)
	if walletConfig.ScanMode != "blacklist" {
		t.Errorf("Expected wallet-specific ScanMode to be blacklist, got %s", walletConfig.ScanMode)
	}

	if len(walletConfig.ExcludeTokens) != 2 {
		t.Errorf("Expected 2 excluded tokens, got %d", len(walletConfig.ExcludeTokens))
	}

	// Test getting global config for non-existent wallet
	globalConfig := config.GetScanConfigForWallet("nonexistentwallet")
	if globalConfig.ScanMode != config.Scan.ScanMode {
		t.Errorf("Expected global ScanMode, got different value")
	}
}
