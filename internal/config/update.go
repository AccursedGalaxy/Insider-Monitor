package config

import (
	"log"
	"time"
)

// UpdateRequest represents a request to update the configuration
type UpdateRequest struct {
	NetworkURL   *string        `json:"network_url,omitempty"`
	Wallets      *[]string      `json:"wallets,omitempty"`
	ScanInterval *string        `json:"scan_interval,omitempty"`
	Alerts       *AlertsUpdate  `json:"alerts,omitempty"`
	Discord      *DiscordUpdate `json:"discord,omitempty"`
}

// AlertsUpdate represents an update to the alerts configuration
type AlertsUpdate struct {
	MinimumBalance    *float64  `json:"minimum_balance,omitempty"`
	SignificantChange *float64  `json:"significant_change,omitempty"`
	IgnoreTokens      *[]string `json:"ignore_tokens,omitempty"`
}

// DiscordUpdate represents an update to the Discord configuration
type DiscordUpdate struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	WebhookURL *string `json:"webhook_url,omitempty"`
	ChannelID  *string `json:"channel_id,omitempty"`
}

// Apply applies the update to the configuration
func (u *UpdateRequest) Apply(cfg *Config) error {
	// Apply network URL update
	if u.NetworkURL != nil {
		cfg.NetworkURL = *u.NetworkURL
	}

	// Apply wallets update
	if u.Wallets != nil {
		cfg.Wallets = *u.Wallets
	}

	// Apply scan interval update
	if u.ScanInterval != nil {
		// Just validate the string can be parsed as duration
		_, err := time.ParseDuration(*u.ScanInterval)
		if err != nil {
			return err
		}
		cfg.ScanInterval = *u.ScanInterval
	}

	// Apply alerts update
	if u.Alerts != nil {
		log.Printf("Updating alerts: %+v", u.Alerts)

		if u.Alerts.MinimumBalance != nil {
			log.Printf("Setting minimum balance to: %f", *u.Alerts.MinimumBalance)
			cfg.Alerts.MinimumBalance = *u.Alerts.MinimumBalance
		}

		if u.Alerts.SignificantChange != nil {
			log.Printf("Setting significant change to: %f", *u.Alerts.SignificantChange)
			cfg.Alerts.SignificantChange = *u.Alerts.SignificantChange
		}

		if u.Alerts.IgnoreTokens != nil {
			cfg.Alerts.IgnoreTokens = *u.Alerts.IgnoreTokens
		}
	}

	// Apply Discord update
	if u.Discord != nil {
		if u.Discord.Enabled != nil {
			cfg.Discord.Enabled = *u.Discord.Enabled
		}
		if u.Discord.WebhookURL != nil {
			cfg.Discord.WebhookURL = *u.Discord.WebhookURL
		}
		if u.Discord.ChannelID != nil {
			cfg.Discord.ChannelID = *u.Discord.ChannelID
		}
	}

	return nil
}
