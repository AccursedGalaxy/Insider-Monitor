package config

import (
	"log"
	"os"
	"sync"
)

// Service provides configuration management functionality
type Service struct {
	config *Config
	mutex  sync.RWMutex
}

var (
	configService     *Service
	configServiceOnce sync.Once
)

// NewService creates a new configuration service
func NewService(configPath string) (*Service, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		config: config,
	}, nil
}

// GetInstance returns the singleton instance of the config service
func GetInstance() *Service {
	configServiceOnce.Do(func() {
		// Use a standard location for config
		configPath := "config.json"

		// Try to load from environment variable if set
		if envPath := getEnvOrDefault("CONFIG_PATH", ""); envPath != "" {
			configPath = envPath
		}

		service, err := NewService(configPath)
		if err != nil {
			log.Printf("Warning: Failed to load config: %v", err)
			log.Printf("Using default configuration")
			service = &Service{
				config: GetTestConfig(),
			}
		}

		configService = service
	})

	return configService
}

// GetConfig returns the current configuration
func (s *Service) GetConfig() *Config {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a deep copy to prevent concurrent access issues
	return s.config
}

// UpdateConfig updates the configuration
func (s *Service) UpdateConfig(update UpdateRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.config.Update(update)
}

// AddWallet adds a new wallet to the configuration
func (s *Service) AddWallet(wallet string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if wallet already exists
	for _, w := range s.config.Wallets {
		if w == wallet {
			// Already exists, nothing to do
			return nil
		}
	}

	// Add the wallet
	s.config.Wallets = append(s.config.Wallets, wallet)

	// Save the changes
	return s.config.Save()
}

// RemoveWallet removes a wallet from the configuration
func (s *Service) RemoveWallet(wallet string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find the wallet index
	index := -1
	for i, w := range s.config.Wallets {
		if w == wallet {
			index = i
			break
		}
	}

	// If not found, nothing to do
	if index == -1 {
		return nil
	}

	// Remove the wallet from the slice
	s.config.Wallets = append(s.config.Wallets[:index], s.config.Wallets[index+1:]...)

	// Remove any wallet-specific configuration
	delete(s.config.WalletConfigs, wallet)

	// Save the changes
	return s.config.Save()
}

// GetWalletConfig returns the configuration for a specific wallet
func (s *Service) GetWalletConfig(wallet string) *WalletConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if cfg, exists := s.config.WalletConfigs[wallet]; exists {
		return &cfg
	}

	// Return a default wallet config if none exists
	return &WalletConfig{
		Scan: &ScanConfig{
			ScanMode:      "all",
			IncludeTokens: []string{},
			ExcludeTokens: []string{},
		},
	}
}

// UpdateWalletConfig updates the configuration for a specific wallet
func (s *Service) UpdateWalletConfig(wallet string, config WalletConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Initialize wallet configs if nil
	if s.config.WalletConfigs == nil {
		s.config.WalletConfigs = make(map[string]WalletConfig)
	}

	// Update the wallet config
	s.config.WalletConfigs[wallet] = config

	// Save the changes
	return s.config.Save()
}

// Helper function to get environment variable with a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// LookupEnv looks up an environment variable
// This is a function variable to allow for testing
var LookupEnv = func(key string) (string, bool) {
	return LookupEnvFunc(key)
}

// LookupEnvFunc is the actual implementation of LookupEnv
var LookupEnvFunc = func(key string) (string, bool) {
	return os.LookupEnv(key)
}
