package routes

import (
	"log"
	"sync"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/api/middleware"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/websocket"
	"github.com/gin-gonic/gin"
)

var (
	wsManager     *websocket.Manager
	wsManagerOnce sync.Once
)

// GetWSManager returns a singleton WebSocket manager
func GetWSManager() *websocket.Manager {
	wsManagerOnce.Do(func() {
		wsManager = websocket.NewManager()
		// Start the manager in a goroutine
		go wsManager.Start()
		log.Println("WebSocket manager started")
	})
	return wsManager
}

// SetupWebSocketRoutes configures the WebSocket routes
func SetupWebSocketRoutes(router *gin.Engine) {
	// Get singleton manager
	manager := GetWSManager()

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		// Upgrade to WebSocket
		manager.HandleWebSocket(c.Writer, c.Request)
	})

	// Protected WebSocket endpoint for admin-level operations
	router.GET("/ws/admin", middleware.Authenticate(), func(c *gin.Context) {
		// This endpoint could be used for admin-specific WebSocket communications
		// For now, it's just using the same handler but with authentication
		manager.HandleWebSocket(c.Writer, c.Request)
	})
}
