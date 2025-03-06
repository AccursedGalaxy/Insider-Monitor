package models

import (
	"time"
)

// WalletData represents data for a monitored wallet
type WalletData struct {
	Address       string                  `json:"address"`
	Label         string                  `json:"label,omitempty"`
	LastScanned   time.Time               `json:"last_scanned"`
	TokenCount    int                     `json:"token_count"`
	TokenAccounts map[string]TokenAccount `json:"token_accounts,omitempty"`
}

// TokenAccount represents data for a token account
type TokenAccount struct {
	Mint        string    `json:"mint"`
	Symbol      string    `json:"symbol"`
	Balance     uint64    `json:"balance"`
	Decimals    uint8     `json:"decimals"`
	LastUpdated time.Time `json:"last_updated"`
}

// WalletConfig represents configuration for a specific wallet
type WalletConfig struct {
	ScanConfig *ScanConfig `json:"scan,omitempty"`
}

// ScanConfig represents scanning configuration
type ScanConfig struct {
	ScanMode      string   `json:"scan_mode"`      // "all", "whitelist", or "blacklist"
	IncludeTokens []string `json:"include_tokens"` // Specific tokens to include (if using whitelist)
	ExcludeTokens []string `json:"exclude_tokens"` // Specific tokens to exclude (if using blacklist)
}
