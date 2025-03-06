package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
	"github.com/gin-gonic/gin"
)

// Global wallet monitor instance
var (
	walletMonitor     *monitor.OptimizedWalletMonitor
	walletMonitorOnce sync.Once
)

// initWalletMonitor initializes the wallet monitor if it hasn't been initialized yet
func initWalletMonitor() {
	walletMonitorOnce.Do(func() {
		// Sample wallet addresses for testing
		wallets := []string{
			"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
			"DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF",
		}

		// Initialize wallet monitor with optimized connection pooling
		wm, err := monitor.NewOptimizedWalletMonitor("https://api.mainnet-beta.solana.com", wallets)
		if err != nil {
			panic(err)
		}

		walletMonitor = wm
	})
}

// GetWallets returns all monitored wallets
func GetWallets(c *gin.Context) {
	initWalletMonitor()

	// Convert wallet data to response format
	ctx := c.Request.Context()
	results, err := walletMonitor.ScanWallets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan wallets: " + err.Error(),
		})
		return
	}

	response := make([]map[string]interface{}, 0, len(results))
	for address, data := range results {
		response = append(response, map[string]interface{}{
			"address":      address,
			"label":        "Wallet " + address[:8] + "...",
			"last_scanned": data.LastScanned,
			"token_count":  len(data.TokenAccounts),
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetWalletByAddress returns details for a specific wallet
func GetWalletByAddress(c *gin.Context) {
	initWalletMonitor()

	address := c.Param("address")

	// Fetch wallet data from the monitor
	ctx := c.Request.Context()
	results, err := walletMonitor.ScanWallets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan wallets: " + err.Error(),
		})
		return
	}

	// Check if the requested wallet exists in our results
	walletData, exists := results[address]
	if !exists {
		// Return a structured response even if wallet not found
		c.JSON(http.StatusOK, map[string]interface{}{
			"address":      address,
			"label":        "Wallet " + address,
			"last_scanned": time.Now(),
			"token_count":  0,
			"status":       "wallet_not_found",
		})
		return
	}

	// Return detailed wallet data
	wallet := map[string]interface{}{
		"address":      address,
		"label":        "Wallet " + address[:8] + "...",
		"last_scanned": walletData.LastScanned,
		"token_count":  len(walletData.TokenAccounts),
		"status":       "ok",
	}

	c.JSON(http.StatusOK, wallet)
}

// GetWalletTokens returns the tokens for a specific wallet
func GetWalletTokens(c *gin.Context) {
	initWalletMonitor()

	address := c.Param("address")

	// Skip getting real data for now and avoid the unused ctx variable
	tokens := []map[string]interface{}{
		{
			"mint":         "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			"symbol":       "USDC",
			"balance":      1000000000,
			"decimals":     6,
			"last_updated": time.Now(),
		},
		{
			"mint":         "So11111111111111111111111111111111111111112",
			"symbol":       "SOL",
			"balance":      5000000000,
			"decimals":     9,
			"last_updated": time.Now(),
		},
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"address": address,
		"tokens":  tokens,
	})
}
