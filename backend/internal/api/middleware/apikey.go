package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyConfig holds the API key configuration
type APIKeyConfig struct {
	Header string
	Keys   map[string]string // key -> description mapping
}

// Default API key configuration
var defaultAPIKeyConfig = APIKeyConfig{
	Header: "X-API-Key",
	Keys: map[string]string{
		"insider-monitor-key-123": "default",
		// Add more API keys for different clients if needed
	},
}

// APIKey middleware checks for a valid API key
func APIKey() gin.HandlerFunc {
	return apiKeyWithConfig(defaultAPIKeyConfig)
}

// APIKeyWithConfig allows custom API key configuration
func APIKeyWithConfig(config APIKeyConfig) gin.HandlerFunc {
	return apiKeyWithConfig(config)
}

// apiKeyWithConfig is the internal implementation
func apiKeyWithConfig(config APIKeyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader(config.Header)

		// Allow Bearer token format for consistency with JWT
		if strings.HasPrefix(apiKey, "Bearer ") {
			apiKey = strings.TrimPrefix(apiKey, "Bearer ")
		}

		// Check if API key exists
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			return
		}

		// Check if API key is valid
		clientName, ok := config.Keys[apiKey]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		// Set client information in context
		c.Set("client", clientName)

		// Continue to next handler
		c.Next()
	}
}

// OAuthOrAPIKey allows authentication using either OAuth JWT or API key
func OAuthOrAPIKey() gin.HandlerFunc {
	jwtAuth := authenticateWithConfig(defaultJWTConfig)
	apiKeyAuth := apiKeyWithConfig(defaultAPIKeyConfig)

	return func(c *gin.Context) {
		// Check for API key first
		apiKey := c.GetHeader(defaultAPIKeyConfig.Header)
		if apiKey != "" {
			apiKeyAuth(c)
			return
		}

		// Fall back to JWT authentication
		jwtAuth(c)
	}
}
