package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthTest() (*gin.Engine, *JWTConfig) {
	// Use test mode to avoid debug logs
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.New()

	// Create a test JWT config with a known secret
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenExpiry:   time.Hour,
		TokenIssuer:   "test-issuer",
		TokenAudience: "test-audience",
	}

	return r, &config
}

func TestGenerateToken(t *testing.T) {
	_, config := setupAuthTest()

	// Generate a token with the test config
	token, err := GenerateTokenWithConfig("testuser", "admin", *config)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthenticateWithValidToken(t *testing.T) {
	r, config := setupAuthTest()

	// Generate a valid token
	token, err := GenerateTokenWithConfig("testuser", "admin", *config)
	assert.NoError(t, err)

	// Create the auth middleware with our config
	authMiddleware := AuthenticateWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", authMiddleware, func(c *gin.Context) {
		// Get the username from the context
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		c.JSON(http.StatusOK, gin.H{
			"message":  "protected resource accessed",
			"username": username,
			"role":     role,
		})
	})

	// Create a test request with the token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "protected resource accessed")
}

func TestAuthenticateWithInvalidToken(t *testing.T) {
	r, config := setupAuthTest()

	// Create the auth middleware with our config
	authMiddleware := AuthenticateWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Create a test request with an invalid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthenticateWithNoToken(t *testing.T) {
	r, config := setupAuthTest()

	// Create the auth middleware with our config
	authMiddleware := AuthenticateWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Create a test request with no token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authentication required")
}

func TestExpiredToken(t *testing.T) {
	r, config := setupAuthTest()

	// Create a config with a very short expiry
	expiredConfig := *config
	expiredConfig.TokenExpiry = time.Millisecond * 10

	// Generate a token that will expire quickly
	token, err := GenerateTokenWithConfig("testuser", "admin", expiredConfig)
	assert.NoError(t, err)

	// Wait for the token to expire
	time.Sleep(time.Millisecond * 20)

	// Create the auth middleware with our config
	authMiddleware := AuthenticateWithConfig(*config)

	// Add a protected endpoint
	r.GET("/protected", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "protected resource accessed",
		})
	})

	// Create a test request with the expired token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token is expired")
}
