package handlers

import (
	"net/http"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/alerts"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/config"
	"github.com/gin-gonic/gin"
)

// GetAlerts returns recent alerts
func GetAlerts(c *gin.Context) {
	// Get the alerts service
	alertsService := alerts.GetInstance()

	// Get query parameters
	walletAddress := c.Query("wallet")

	// Get alerts
	var alertsList []alerts.Alert
	if walletAddress != "" {
		// Filter by wallet
		alertsList = alertsService.GetAlertsByWallet(walletAddress)
	} else {
		// Get all alerts
		alertsList = alertsService.GetAlerts()
	}

	c.JSON(http.StatusOK, alertsList)
}

// GetAlertSettings returns alert configuration settings
func GetAlertSettings(c *gin.Context) {
	// Get the config service
	configService := config.GetInstance()

	// Get the current configuration
	cfg := configService.GetConfig()

	// Return just the alerts section and Discord settings
	settings := map[string]interface{}{
		"minimum_balance":    cfg.Alerts.MinimumBalance,
		"significant_change": cfg.Alerts.SignificantChange,
		"ignore_tokens":      cfg.Alerts.IgnoreTokens,
		"discord": map[string]interface{}{
			"enabled":     cfg.Discord.Enabled,
			"webhook_url": cfg.Discord.WebhookURL,
			"channel_id":  cfg.Discord.ChannelID,
		},
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateAlertSettings updates alert configuration settings
func UpdateAlertSettings(c *gin.Context) {
	var updateRequest struct {
		MinimumBalance    *float64 `json:"minimum_balance"`
		SignificantChange *float64 `json:"significant_change"`
		IgnoreTokens      []string `json:"ignore_tokens"`
		Discord           struct {
			Enabled    *bool   `json:"enabled"`
			WebhookURL *string `json:"webhook_url"`
			ChannelID  *string `json:"channel_id"`
		} `json:"discord"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the config service
	configService := config.GetInstance()

	// Get the current configuration for updating
	cfg := configService.GetConfig()

	// Create an update request
	configUpdate := config.UpdateRequest{}

	// Update alerts configuration
	alerts := cfg.Alerts
	if updateRequest.MinimumBalance != nil {
		alerts.MinimumBalance = *updateRequest.MinimumBalance
	}
	if updateRequest.SignificantChange != nil {
		alerts.SignificantChange = *updateRequest.SignificantChange
	}
	if updateRequest.IgnoreTokens != nil {
		alerts.IgnoreTokens = updateRequest.IgnoreTokens
	}
	configUpdate.Alerts = &alerts

	// Update Discord configuration
	discord := cfg.Discord
	if updateRequest.Discord.Enabled != nil {
		discord.Enabled = *updateRequest.Discord.Enabled
	}
	if updateRequest.Discord.WebhookURL != nil {
		discord.WebhookURL = *updateRequest.Discord.WebhookURL
	}
	if updateRequest.Discord.ChannelID != nil {
		discord.ChannelID = *updateRequest.Discord.ChannelID
	}
	configUpdate.Discord = &discord

	// Update the configuration
	if err := configService.UpdateConfig(configUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update alert settings: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Alert settings updated successfully",
		"settings": map[string]interface{}{
			"minimum_balance":    alerts.MinimumBalance,
			"significant_change": alerts.SignificantChange,
			"ignore_tokens":      alerts.IgnoreTokens,
			"discord": map[string]interface{}{
				"enabled":     discord.Enabled,
				"webhook_url": discord.WebhookURL,
				"channel_id":  discord.ChannelID,
			},
		},
	})
}
