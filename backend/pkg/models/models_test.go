package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWalletDataJSON(t *testing.T) {
	// Create a sample WalletData
	now := time.Now().UTC().Truncate(time.Millisecond)
	walletData := WalletData{
		Address:     "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
		Label:       "Test Wallet",
		LastScanned: now,
		TokenCount:  2,
		TokenAccounts: map[string]TokenAccount{
			"token1": {
				Mint:        "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				Symbol:      "USDC",
				Balance:     1000000,
				Decimals:    6,
				LastUpdated: now,
			},
			"token2": {
				Mint:        "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
				Symbol:      "BONK",
				Balance:     10000000000,
				Decimals:    5,
				LastUpdated: now,
			},
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(walletData)
	if err != nil {
		t.Fatalf("Failed to marshal WalletData: %v", err)
	}

	// Deserialize from JSON
	var decoded WalletData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal WalletData: %v", err)
	}

	// Verify fields
	if decoded.Address != walletData.Address {
		t.Errorf("Expected address %s, got %s", walletData.Address, decoded.Address)
	}

	if decoded.Label != walletData.Label {
		t.Errorf("Expected label %s, got %s", walletData.Label, decoded.Label)
	}

	if !decoded.LastScanned.Equal(walletData.LastScanned) {
		t.Errorf("Expected last scanned %v, got %v", walletData.LastScanned, decoded.LastScanned)
	}

	if decoded.TokenCount != walletData.TokenCount {
		t.Errorf("Expected token count %d, got %d", walletData.TokenCount, decoded.TokenCount)
	}

	// Check token accounts
	if len(decoded.TokenAccounts) != len(walletData.TokenAccounts) {
		t.Errorf("Expected %d token accounts, got %d", len(walletData.TokenAccounts), len(decoded.TokenAccounts))
	}

	for key, original := range walletData.TokenAccounts {
		decodedToken, ok := decoded.TokenAccounts[key]
		if !ok {
			t.Errorf("Missing token account: %s", key)
			continue
		}

		if decodedToken.Mint != original.Mint {
			t.Errorf("Token %s: expected mint %s, got %s", key, original.Mint, decodedToken.Mint)
		}

		if decodedToken.Symbol != original.Symbol {
			t.Errorf("Token %s: expected symbol %s, got %s", key, original.Symbol, decodedToken.Symbol)
		}

		if decodedToken.Balance != original.Balance {
			t.Errorf("Token %s: expected balance %d, got %d", key, original.Balance, decodedToken.Balance)
		}

		if decodedToken.Decimals != original.Decimals {
			t.Errorf("Token %s: expected decimals %d, got %d", key, original.Decimals, decodedToken.Decimals)
		}

		if !decodedToken.LastUpdated.Equal(original.LastUpdated) {
			t.Errorf("Token %s: expected last updated %v, got %v", key, original.LastUpdated, decodedToken.LastUpdated)
		}
	}
}

func TestAlertJSON(t *testing.T) {
	// Create a sample Alert
	now := time.Now().UTC().Truncate(time.Millisecond)
	alert := Alert{
		ID:            "alert-123",
		Timestamp:     now,
		WalletAddress: "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
		TokenMint:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		AlertType:     "token_change",
		Message:       "Token balance has changed significantly",
		Level:         Warning,
		Data: map[string]interface{}{
			"old_balance": 1000000,
			"new_balance": 2000000,
			"change":      1.0,
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(alert)
	if err != nil {
		t.Fatalf("Failed to marshal Alert: %v", err)
	}

	// Deserialize from JSON
	var decoded Alert
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Alert: %v", err)
	}

	// Verify fields
	if decoded.ID != alert.ID {
		t.Errorf("Expected ID %s, got %s", alert.ID, decoded.ID)
	}

	if !decoded.Timestamp.Equal(alert.Timestamp) {
		t.Errorf("Expected timestamp %v, got %v", alert.Timestamp, decoded.Timestamp)
	}

	if decoded.WalletAddress != alert.WalletAddress {
		t.Errorf("Expected wallet address %s, got %s", alert.WalletAddress, decoded.WalletAddress)
	}

	if decoded.TokenMint != alert.TokenMint {
		t.Errorf("Expected token mint %s, got %s", alert.TokenMint, decoded.TokenMint)
	}

	if decoded.AlertType != alert.AlertType {
		t.Errorf("Expected alert type %s, got %s", alert.AlertType, decoded.AlertType)
	}

	if decoded.Message != alert.Message {
		t.Errorf("Expected message %s, got %s", alert.Message, decoded.Message)
	}

	if decoded.Level != alert.Level {
		t.Errorf("Expected level %s, got %s", alert.Level, decoded.Level)
	}

	// Check data map (comparing a few keys)
	oldBalance, ok := decoded.Data["old_balance"]
	if !ok || oldBalance != float64(1000000) {
		t.Errorf("Expected old_balance 1000000, got %v", oldBalance)
	}

	newBalance, ok := decoded.Data["new_balance"]
	if !ok || newBalance != float64(2000000) {
		t.Errorf("Expected new_balance 2000000, got %v", newBalance)
	}
}

func TestConfigJSON(t *testing.T) {
	// Create a sample Config
	config := Config{
		NetworkURL:   "https://api.mainnet-beta.solana.com",
		Wallets:      []string{"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr", "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF"},
		ScanInterval: "5m",
		Scan: ScanConfig{
			ScanMode:      "all",
			IncludeTokens: []string{},
			ExcludeTokens: []string{},
		},
		WalletConfigs: map[string]WalletConfig{
			"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr": {
				ScanConfig: &ScanConfig{
					ScanMode:      "whitelist",
					IncludeTokens: []string{"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"},
					ExcludeTokens: []string{},
				},
			},
		},
		Alerts: AlertSettings{
			MinimumBalance:    1000,
			SignificantChange: 0.1,
			IgnoreTokens:      []string{},
			Discord: DiscordConfig{
				Enabled:    true,
				WebhookURL: "https://discord.com/api/webhooks/test",
				ChannelID:  "test-channel",
			},
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal Config: %v", err)
	}

	// Deserialize from JSON
	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Config: %v", err)
	}

	// Verify fields
	if decoded.NetworkURL != config.NetworkURL {
		t.Errorf("Expected network URL %s, got %s", config.NetworkURL, decoded.NetworkURL)
	}

	if len(decoded.Wallets) != len(config.Wallets) {
		t.Errorf("Expected %d wallets, got %d", len(config.Wallets), len(decoded.Wallets))
	} else {
		for i, wallet := range config.Wallets {
			if decoded.Wallets[i] != wallet {
				t.Errorf("Expected wallet %s at index %d, got %s", wallet, i, decoded.Wallets[i])
			}
		}
	}

	if decoded.ScanInterval != config.ScanInterval {
		t.Errorf("Expected scan interval %s, got %s", config.ScanInterval, decoded.ScanInterval)
	}

	// Check scan config
	if decoded.Scan.ScanMode != config.Scan.ScanMode {
		t.Errorf("Expected scan mode %s, got %s", config.Scan.ScanMode, decoded.Scan.ScanMode)
	}

	// Check wallet configs
	if len(decoded.WalletConfigs) != len(config.WalletConfigs) {
		t.Errorf("Expected %d wallet configs, got %d", len(config.WalletConfigs), len(decoded.WalletConfigs))
	}

	// Check alerts settings
	if decoded.Alerts.MinimumBalance != config.Alerts.MinimumBalance {
		t.Errorf("Expected minimum balance %f, got %f", config.Alerts.MinimumBalance, decoded.Alerts.MinimumBalance)
	}

	if decoded.Alerts.SignificantChange != config.Alerts.SignificantChange {
		t.Errorf("Expected significant change %f, got %f", config.Alerts.SignificantChange, decoded.Alerts.SignificantChange)
	}

	if decoded.Alerts.Discord.Enabled != config.Alerts.Discord.Enabled {
		t.Errorf("Expected Discord enabled %t, got %t", config.Alerts.Discord.Enabled, decoded.Alerts.Discord.Enabled)
	}
}
