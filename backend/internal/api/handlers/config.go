package handlers

import (
	"net/http"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/config"
	"github.com/gin-gonic/gin"
)

// GetConfig returns the current application configuration
func GetConfig(c *gin.Context) {
	// Get the config service
	configService := config.GetInstance()

	// Get the current configuration
	cfg := configService.GetConfig()

	// Return the configuration
	c.JSON(http.StatusOK, cfg)
}

// UpdateConfig updates the application configuration
func UpdateConfig(c *gin.Context) {
	var updateRequest config.UpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the config service
	configService := config.GetInstance()

	// Update the configuration
	if err := configService.UpdateConfig(updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update configuration: " + err.Error(),
		})
		return
	}

	// Get the updated configuration to return
	updatedConfig := configService.GetConfig()

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config":  updatedConfig,
	})
}

// GetWalletConfigs returns the configuration for all wallets
func GetWalletConfigs(c *gin.Context) {
	// Get the config service
	configService := config.GetInstance()

	// Get the current configuration
	cfg := configService.GetConfig()

	// Return the wallet configurations
	c.JSON(http.StatusOK, cfg.WalletConfigs)
}

// UpdateWalletConfig updates the configuration for a specific wallet
func UpdateWalletConfig(c *gin.Context) {
	address := c.Param("address")

	// Validate the address
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address is required"})
		return
	}

	var walletConfig config.WalletConfig
	if err := c.ShouldBindJSON(&walletConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the config service
	configService := config.GetInstance()

	// Update the wallet configuration
	if err := configService.UpdateWalletConfig(address, walletConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update wallet configuration: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet configuration updated successfully",
		"address": address,
		"config":  walletConfig,
	})
}

// DeleteWalletConfig removes the configuration for a specific wallet
func DeleteWalletConfig(c *gin.Context) {
	address := c.Param("address")

	// Validate the address
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address is required"})
		return
	}

	// Get the config service
	configService := config.GetInstance()

	// Remove the wallet
	if err := configService.RemoveWallet(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete wallet configuration: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet configuration deleted successfully",
		"address": address,
	})
}

// AddWallet adds a new wallet to the configuration
func AddWallet(c *gin.Context) {
	var request struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the address
	if request.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet address is required"})
		return
	}

	// Get the config service
	configService := config.GetInstance()

	// Add the wallet
	if err := configService.AddWallet(request.Address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add wallet: " + err.Error(),
		})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet added successfully",
		"address": request.Address,
	})
}
