package config

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
)

func TestNewService(t *testing.T) {
	// Create a temporary config file
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

	// Create a new service
	service, err := NewService(path)
	if err != nil {
		t.Fatalf("Failed to create new service: %v", err)
	}

	if service == nil {
		t.Fatal("Expected service to be non-nil")
	}

	if service.config == nil {
		t.Fatal("Expected service.config to be non-nil")
	}
}

func TestGetInstance(t *testing.T) {
	// Reset the singleton before testing
	configService = nil
	configServiceOnce = sync.Once{}

	// Override the LookupEnv function for testing
	origLookupEnv := LookupEnvFunc
	defer func() { LookupEnvFunc = origLookupEnv }()

	LookupEnvFunc = func(key string) (string, bool) {
		return "", false // Always return not found
	}

	// Get the singleton instance
	service := GetInstance()
	if service == nil {
		t.Fatal("Expected service to be non-nil")
	}

	// Get it again - should be the same instance
	service2 := GetInstance()
	if service != service2 {
		t.Error("Expected GetInstance to return the same instance on subsequent calls")
	}
}

func TestServiceGetConfig(t *testing.T) {
	// Create a service with a test config
	config := GetTestConfig()
	service := &Service{
		config: config,
	}

	// Get the config
	cfg := service.GetConfig()
	if cfg == nil {
		t.Fatal("Expected cfg to be non-nil")
	}

	if cfg != config {
		t.Error("Expected GetConfig to return the same config instance")
	}
}

func TestServiceUpdateConfig(t *testing.T) {
	// Create a temporary config file
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

	// Create a new service
	service, err := NewService(path)
	if err != nil {
		t.Fatalf("Failed to create new service: %v", err)
	}

	// Create an update request
	newURL := "https://api.testnet.solana.com"
	update := UpdateRequest{
		NetworkURL: &newURL,
	}

	// Update the config
	err = service.UpdateConfig(update)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify the update was applied
	if service.config.NetworkURL != newURL {
		t.Errorf("Expected NetworkURL to be %s, got %s", newURL, service.config.NetworkURL)
	}
}

func TestServiceAddWallet(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Get the path
	path := tempFile.Name()
	tempFile.Close()

	// Create a service with a test config that has a filepath
	config := GetTestConfig()

	service := &Service{
		config: config,
	}

	// Set the filepath directly
	service.config.filepath = path

	// Add a new wallet
	newWallet := "newTestWallet"
	err = service.AddWallet(newWallet)
	if err != nil {
		t.Fatalf("Failed to add wallet: %v", err)
	}

	// Check if the wallet was added
	found := false
	for _, wallet := range service.config.Wallets {
		if wallet == newWallet {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected wallet %s to be added", newWallet)
	}

	// Test adding the same wallet again - should not duplicate
	originalCount := len(service.config.Wallets)
	err = service.AddWallet(newWallet)
	if err != nil {
		t.Fatalf("Failed to add duplicate wallet: %v", err)
	}

	if len(service.config.Wallets) != originalCount {
		t.Errorf("Expected wallet count to remain %d, got %d", originalCount, len(service.config.Wallets))
	}
}

func TestServiceRemoveWallet(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Get the path
	path := tempFile.Name()
	tempFile.Close()

	// Create a service with a test config that has a specific wallet
	config := GetTestConfig()
	testWallet := "testWallet"
	config.Wallets = append(config.Wallets, testWallet)

	// Also add a wallet-specific config
	if config.WalletConfigs == nil {
		config.WalletConfigs = make(map[string]WalletConfig)
	}
	config.WalletConfigs[testWallet] = WalletConfig{
		Scan: &ScanConfig{
			ScanMode: "whitelist",
		},
	}

	service := &Service{
		config: config,
	}

	// Set the filepath directly
	service.config.filepath = path

	// Remove the wallet
	err = service.RemoveWallet(testWallet)
	if err != nil {
		t.Fatalf("Failed to remove wallet: %v", err)
	}

	// Check if the wallet was removed
	for _, wallet := range service.config.Wallets {
		if wallet == testWallet {
			t.Errorf("Expected wallet %s to be removed", testWallet)
		}
	}

	// Check if the wallet config was removed
	if _, exists := service.config.WalletConfigs[testWallet]; exists {
		t.Errorf("Expected wallet config for %s to be removed", testWallet)
	}

	// Test removing a non-existent wallet - should not error
	err = service.RemoveWallet("nonexistentwallet")
	if err != nil {
		t.Errorf("Expected removing non-existent wallet to succeed, got error: %v", err)
	}
}

func TestGetWalletConfig(t *testing.T) {
	// Create a service with a test config
	config := GetTestConfig()

	// Add a wallet-specific config
	testWallet := "testWallet"
	if config.WalletConfigs == nil {
		config.WalletConfigs = make(map[string]WalletConfig)
	}
	config.WalletConfigs[testWallet] = WalletConfig{
		Scan: &ScanConfig{
			ScanMode: "blacklist",
		},
	}

	service := &Service{
		config: config,
	}

	// Get the wallet config
	walletConfig := service.GetWalletConfig(testWallet)
	if walletConfig == nil {
		t.Fatal("Expected wallet config to be non-nil")
	}

	if walletConfig.Scan.ScanMode != "blacklist" {
		t.Errorf("Expected wallet ScanMode to be blacklist, got %s", walletConfig.Scan.ScanMode)
	}

	// Test getting config for non-existent wallet - should return default
	defaultConfig := service.GetWalletConfig("nonexistentwallet")
	if defaultConfig == nil {
		t.Fatal("Expected default wallet config to be non-nil")
	}

	if defaultConfig.Scan == nil {
		t.Fatal("Expected default wallet Scan config to be non-nil")
	}

	if defaultConfig.Scan.ScanMode != "all" {
		t.Errorf("Expected default ScanMode to be all, got %s", defaultConfig.Scan.ScanMode)
	}
}

func TestUpdateWalletConfig(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Get the path
	path := tempFile.Name()
	tempFile.Close()

	// Create a service with a test config that has a filepath
	config := GetTestConfig()

	service := &Service{
		config: config,
	}

	// Set the filepath directly
	service.config.filepath = path

	// Update a wallet config
	testWallet := "testWallet"
	walletConfig := WalletConfig{
		Scan: &ScanConfig{
			ScanMode:      "whitelist",
			IncludeTokens: []string{"token1", "token2"},
		},
	}

	err = service.UpdateWalletConfig(testWallet, walletConfig)
	if err != nil {
		t.Fatalf("Failed to update wallet config: %v", err)
	}

	// Check if the wallet config was updated
	if _, exists := service.config.WalletConfigs[testWallet]; !exists {
		t.Fatalf("Expected wallet config for %s to exist", testWallet)
	}

	updatedConfig := service.config.WalletConfigs[testWallet]
	if updatedConfig.Scan.ScanMode != "whitelist" {
		t.Errorf("Expected wallet ScanMode to be whitelist, got %s", updatedConfig.Scan.ScanMode)
	}

	if len(updatedConfig.Scan.IncludeTokens) != 2 {
		t.Errorf("Expected 2 included tokens, got %d", len(updatedConfig.Scan.IncludeTokens))
	}
}
