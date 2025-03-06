package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAPIKeyTest() (*gin.Engine, *APIKeyConfig) {
	// Use test mode to avoid debug logs
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.New()

	// Create a test API key config
	config := APIKeyConfig{
		Header: "X-API-Key",
		Keys: map[string]string{
			"test-api-key-123": "test-client",
			"test-api-key-456": "another-client",
		},
	}

	return r, &config
}

func TestAPIKeyWithValidKey(t *testing.T) {
	r, config := setupAPIKeyTest()

	// Create the API key middleware with our config
	apiKeyMiddleware := APIKeyWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", apiKeyMiddleware, func(c *gin.Context) {
		// Get the client from the context
		client, _ := c.Get("client")

		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
			"client":  client,
		})
	})

	// Create a test request with a valid API key
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "test-api-key-123")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "protected resource accessed")
	assert.Contains(t, w.Body.String(), "test-client")
}

func TestAPIKeyWithValidBearerKey(t *testing.T) {
	r, config := setupAPIKeyTest()

	// Create the API key middleware with our config
	apiKeyMiddleware := APIKeyWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", apiKeyMiddleware, func(c *gin.Context) {
		// Get the client from the context
		client, _ := c.Get("client")

		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
			"client":  client,
		})
	})

	// Create a test request with a valid API key in Bearer format
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "Bearer test-api-key-123")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "protected resource accessed")
	assert.Contains(t, w.Body.String(), "test-client")
}

func TestAPIKeyWithInvalidKey(t *testing.T) {
	r, config := setupAPIKeyTest()

	// Create the API key middleware with our config
	apiKeyMiddleware := APIKeyWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", apiKeyMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Create a test request with an invalid API key
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "invalid-api-key")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid API key")
}

func TestAPIKeyWithNoKey(t *testing.T) {
	r, config := setupAPIKeyTest()

	// Create the API key middleware with our config
	apiKeyMiddleware := APIKeyWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", apiKeyMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Create a test request with no API key
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "API key required")
}

func TestOAuthOrAPIKey(t *testing.T) {
	r := gin.New()

	// Set up a test endpoint with the combined middleware
	auth := OAuthOrAPIKey()
	r.GET("/protected", auth, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Test with API key - Use the key directly from the config
	t.Run("Valid API Key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// The API key is the key in the map, not the value
		req.Header.Set("X-API-Key", "insider-monitor-key-123")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test with JWT
	t.Run("Valid JWT", func(t *testing.T) {
		// Generate a valid token
		token, err := GenerateToken("testuser", "admin")
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test with no authentication
	t.Run("No Authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
