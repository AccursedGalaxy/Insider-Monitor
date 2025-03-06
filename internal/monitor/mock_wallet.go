package monitor

import (
	"time"
)

type MockWalletMonitor struct {
	startTime   time.Time
	data        map[string]*WalletData
	scanConfigs map[string]ScanConfigInfo // Maps wallet address to scan config
}

func NewMockWalletMonitor() *MockWalletMonitor {
	// Initialize with more realistic test data
	return &MockWalletMonitor{
		startTime: time.Now(),
		data: map[string]*WalletData{
			"TestWallet1": {
				WalletAddress: "TestWallet1",
				TokenAccounts: map[string]TokenAccountInfo{
					"So11111111111111111111111111111111111111112": { // SOL
						Balance:     1000000000, // 1 SOL
						LastUpdated: time.Now(),
					},
					"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": { // USDC
						Balance:     1000000, // 1 USDC
						LastUpdated: time.Now(),
					},
				},
			},
		},
		scanConfigs: make(map[string]ScanConfigInfo),
	}
}

// SetScanConfig sets scan configuration for a specific wallet
func (m *MockWalletMonitor) SetScanConfig(walletAddr string, scanMode string, includeTokens, excludeTokens []string) {
	m.scanConfigs[walletAddr] = ScanConfigInfo{
		Mode:          scanMode,
		IncludeTokens: includeTokens,
		ExcludeTokens: excludeTokens,
	}
}

// SetGlobalScanConfig applies the same scan configuration to all currently monitored wallets
func (m *MockWalletMonitor) SetGlobalScanConfig(scanMode string, includeTokens, excludeTokens []string) {
	// Apply to all wallets in the data map
	for walletAddr := range m.data {
		m.scanConfigs[walletAddr] = ScanConfigInfo{
			Mode:          scanMode,
			IncludeTokens: includeTokens,
			ExcludeTokens: excludeTokens,
		}
	}

	// Also apply to test wallets
	m.scanConfigs["TestWallet1"] = ScanConfigInfo{
		Mode:          scanMode,
		IncludeTokens: includeTokens,
		ExcludeTokens: excludeTokens,
	}
	m.scanConfigs["TestWallet2"] = ScanConfigInfo{
		Mode:          scanMode,
		IncludeTokens: includeTokens,
		ExcludeTokens: excludeTokens,
	}
}

// shouldIncludeToken determines if a token should be included based on the scan configuration
func (m *MockWalletMonitor) shouldIncludeToken(walletAddr, tokenMint string) bool {
	scanConfig, exists := m.scanConfigs[walletAddr]
	if !exists {
		return true // Default to including all tokens if no config exists
	}

	switch scanConfig.Mode {
	case "whitelist":
		// Only include tokens in the include list
		for _, token := range scanConfig.IncludeTokens {
			if token == tokenMint {
				return true
			}
		}
		return false

	case "blacklist":
		// Include all tokens except those in the exclude list
		for _, token := range scanConfig.ExcludeTokens {
			if token == tokenMint {
				return false
			}
		}
		return true

	case "all", "": // Default or explicitly set to "all"
		return true

	default:
		return true
	}
}

func (m *MockWalletMonitor) ScanAllWallets() (map[string]*WalletData, error) {
	results := make(map[string]*WalletData)
	now := time.Now()
	elapsed := now.Sub(m.startTime)

	// Base wallet always present
	baseWallet := &WalletData{
		WalletAddress: "TestWallet1",
		TokenAccounts: make(map[string]TokenAccountInfo),
		LastScanned:   now,
	}

	// Add tokens based on scan configuration
	solToken := "So11111111111111111111111111111111111111112"
	if m.shouldIncludeToken("TestWallet1", solToken) {
		baseWallet.TokenAccounts[solToken] = TokenAccountInfo{
			Balance:     1000000000,
			LastUpdated: now,
			Symbol:      "SOL",
			Decimals:    9,
		}
	}

	usdcToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	if m.shouldIncludeToken("TestWallet1", usdcToken) {
		baseWallet.TokenAccounts[usdcToken] = TokenAccountInfo{
			Balance:     1000000,
			LastUpdated: now,
			Symbol:      "USDC",
			Decimals:    6,
		}
	}

	// Apply changes based on time
	if elapsed >= 5*time.Second {
		// Add BONK token after 5 seconds
		bonkToken := "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"
		if m.shouldIncludeToken("TestWallet1", bonkToken) {
			baseWallet.TokenAccounts[bonkToken] = TokenAccountInfo{
				Balance:     5000000,
				LastUpdated: now,
				Symbol:      "BONK",
				Decimals:    5,
			}
		}
	}

	if elapsed >= 10*time.Second {
		// Increase SOL balance after 10 seconds
		if m.shouldIncludeToken("TestWallet1", solToken) {
			baseWallet.TokenAccounts[solToken] = TokenAccountInfo{
				Balance:     2000000000,
				LastUpdated: now,
				Symbol:      "SOL",
				Decimals:    9,
			}
		}
	}

	results["TestWallet1"] = baseWallet

	if elapsed >= 15*time.Second {
		// Add second wallet after 15 seconds
		wallet2 := &WalletData{
			WalletAddress: "TestWallet2",
			TokenAccounts: make(map[string]TokenAccountInfo),
			LastScanned:   now,
		}

		// Add SOL token based on scan configuration
		if m.shouldIncludeToken("TestWallet2", solToken) {
			wallet2.TokenAccounts[solToken] = TokenAccountInfo{
				Balance:     5000000000,
				LastUpdated: now,
				Symbol:      "SOL",
				Decimals:    9,
			}
		}

		results["TestWallet2"] = wallet2
	}

	return results, nil
}
