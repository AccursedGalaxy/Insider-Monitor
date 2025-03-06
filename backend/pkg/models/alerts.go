package models

import (
	"time"
)

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	// Info level is for informational alerts
	Info AlertLevel = "INFO"
	// Warning level is for potential issues
	Warning AlertLevel = "WARNING"
	// Critical level is for severe issues
	Critical AlertLevel = "CRITICAL"
)

// Alert represents a notification about a wallet event
type Alert struct {
	ID            string                 `json:"id"`
	Timestamp     time.Time              `json:"timestamp"`
	WalletAddress string                 `json:"wallet_address"`
	TokenMint     string                 `json:"token_mint,omitempty"`
	AlertType     string                 `json:"alert_type"`
	Message       string                 `json:"message"`
	Level         AlertLevel             `json:"level"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

// AlertSettings represents configuration for alerting
type AlertSettings struct {
	MinimumBalance    float64       `json:"minimum_balance"`
	SignificantChange float64       `json:"significant_change"`
	IgnoreTokens      []string      `json:"ignore_tokens"`
	Discord           DiscordConfig `json:"discord"`
}

// DiscordConfig represents Discord webhook configuration
type DiscordConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	ChannelID  string `json:"channel_id"`
}
