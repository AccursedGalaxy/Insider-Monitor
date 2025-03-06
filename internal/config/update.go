package config

import (
	"log"
	"time"
)

// UpdateRequest represents a request to update the configuration
type UpdateRequest struct {
	NetworkURL    *string                        `json:"network_url,omitempty"`
	Wallets       *[]string                      `json:"wallets,omitempty"`
	ScanInterval  *string                        `json:"scan_interval,omitempty"`
	Scan          *ScanConfigUpdate              `json:"scan,omitempty"`
	WalletConfigs *map[string]WalletConfigUpdate `json:"wallet_configs,omitempty"`
	Alerts        *AlertsUpdate                  `json:"alerts,omitempty"`
	Discord       *DiscordUpdate                 `json:"discord,omitempty"`
}

// ScanConfigUpdate represents an update to scan configuration
type ScanConfigUpdate struct {
	IncludeTokens *[]string `json:"include_tokens,omitempty"`
	ExcludeTokens *[]string `json:"exclude_tokens,omitempty"`
	ScanMode      *string   `json:"scan_mode,omitempty"`
}

// WalletConfigUpdate represents an update to wallet-specific configuration
type WalletConfigUpdate struct {
	Scan *ScanConfigUpdate `json:"scan,omitempty"`
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

	// Apply scan configuration update
	if u.Scan != nil {
		log.Printf("Updating global scan configuration")

		if u.Scan.ScanMode != nil {
			// Validate scan mode
			mode := *u.Scan.ScanMode
			if mode != "all" && mode != "whitelist" && mode != "blacklist" {
				log.Printf("Invalid scan mode: %s, defaulting to 'all'", mode)
				mode = "all"
			}
			cfg.Scan.ScanMode = mode
		}

		if u.Scan.IncludeTokens != nil {
			cfg.Scan.IncludeTokens = *u.Scan.IncludeTokens
		}

		if u.Scan.ExcludeTokens != nil {
			cfg.Scan.ExcludeTokens = *u.Scan.ExcludeTokens
		}
	}

	// Apply wallet-specific configurations
	if u.WalletConfigs != nil {
		log.Printf("Updating wallet-specific configurations")

		// Initialize the map if it doesn't exist
		if cfg.WalletConfigs == nil {
			cfg.WalletConfigs = make(map[string]WalletConfig)
		}

		// Process each wallet config update
		for walletAddr, walletCfgUpdate := range *u.WalletConfigs {
			// Get or create wallet config
			walletCfg, exists := cfg.WalletConfigs[walletAddr]
			if !exists {
				walletCfg = WalletConfig{}
			}

			// Apply scan config update if provided
			if walletCfgUpdate.Scan != nil {
				// Create new scan config if it doesn't exist
				if walletCfg.Scan == nil {
					walletCfg.Scan = &ScanConfig{}
				}

				// Update scan mode if provided
				if walletCfgUpdate.Scan.ScanMode != nil {
					// Validate scan mode
					mode := *walletCfgUpdate.Scan.ScanMode
					if mode != "all" && mode != "whitelist" && mode != "blacklist" {
						log.Printf("Invalid scan mode for wallet %s: %s, defaulting to 'all'", walletAddr, mode)
						mode = "all"
					}
					walletCfg.Scan.ScanMode = mode
				}

				// Update include tokens if provided
				if walletCfgUpdate.Scan.IncludeTokens != nil {
					walletCfg.Scan.IncludeTokens = *walletCfgUpdate.Scan.IncludeTokens
				}

				// Update exclude tokens if provided
				if walletCfgUpdate.Scan.ExcludeTokens != nil {
					walletCfg.Scan.ExcludeTokens = *walletCfgUpdate.Scan.ExcludeTokens
				}
			}

			// Update the wallet config in the main config
			cfg.WalletConfigs[walletAddr] = walletCfg
		}
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
