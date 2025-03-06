package routes

import (
	"net/http"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/api/handlers"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns a router with all API routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Add custom middlewares
	router.Use(middleware.RequestLogger())
	router.Use(middleware.RequestMetrics())

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add rate limiting middleware
	router.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/refresh", middleware.Authenticate(), handlers.RefreshToken)
		}

		// Wallet routes
		wallets := v1.Group("/wallets")
		{
			wallets.GET("", handlers.GetWallets)
			wallets.GET("/:address", handlers.GetWalletByAddress)
			wallets.GET("/:address/tokens", handlers.GetWalletTokens)
		}

		// Configuration routes - protected by authentication
		config := v1.Group("/config")
		config.Use(middleware.OAuthOrAPIKey())
		{
			config.GET("", handlers.GetConfig)
			config.PUT("", handlers.UpdateConfig)
			config.GET("/wallets", handlers.GetWalletConfigs)
			config.POST("/wallets", handlers.AddWallet)
			config.PUT("/wallets/:address", handlers.UpdateWalletConfig)
			config.DELETE("/wallets/:address", handlers.DeleteWalletConfig)
		}

		// Alert routes
		alerts := v1.Group("/alerts")
		{
			alerts.GET("", handlers.GetAlerts)
			alerts.GET("/settings", middleware.OAuthOrAPIKey(), handlers.GetAlertSettings)
			alerts.PUT("/settings", middleware.OAuthOrAPIKey(), handlers.UpdateAlertSettings)
		}

		// Status routes
		status := v1.Group("/status")
		{
			status.GET("", handlers.GetSystemStatus)
			status.GET("/scan", handlers.GetScanStatus)

			// Add metrics endpoint (protected)
			status.GET("/metrics", middleware.OAuthOrAPIKey(), handlers.GetMetrics)
		}
	}

	// Set up WebSocket routes
	SetupWebSocketRoutes(router)

	return router
}
