package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSystemStatus returns the current system status
func GetSystemStatus(c *gin.Context) {
	// TODO: Implement system status retrieval from various services
	// This is a placeholder implementation
	status := map[string]interface{}{
		"status":            "running",
		"uptime":            3600, // seconds
		"version":           "1.0.0",
		"backend_version":   "1.0.0",
		"frontend_version":  "1.0.0",
		"last_scan":         time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		"next_scan":         time.Now().Add(55 * time.Second).Format(time.RFC3339),
		"connected_clients": 2,
		"memory_usage":      "64MB",
		"cpu_usage":         "2%",
	}

	c.JSON(http.StatusOK, status)
}

// GetScanStatus returns the status of wallet scanning
func GetScanStatus(c *gin.Context) {
	// TODO: Implement scan status retrieval from the monitor service
	// This is a placeholder implementation
	status := map[string]interface{}{
		"is_scanning":       false,
		"scan_interval":     "1m",
		"last_scan":         time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		"next_scan":         time.Now().Add(55 * time.Second).Format(time.RFC3339),
		"monitored_wallets": 2,
		"scanned_tokens":    15,
		"scan_history": []map[string]interface{}{
			{
				"timestamp":  time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
				"duration":   "1.5s",
				"successful": true,
			},
			{
				"timestamp":  time.Now().Add(-6 * time.Minute).Format(time.RFC3339),
				"duration":   "1.2s",
				"successful": true,
			},
		},
	}

	c.JSON(http.StatusOK, status)
}
