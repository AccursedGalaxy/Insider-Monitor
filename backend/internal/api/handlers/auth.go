package handlers

import (
	"net/http"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

// Default credentials (for development only)
// In a production environment, use a proper authentication system
var defaultCredentials = map[string]string{
	"admin": "admin123",
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists and password is correct
	// This is a simplified example - in a real application, use a proper user store and password hashing
	password, exists := defaultCredentials[req.Username]
	if !exists || password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Determine user role (simplified - use a proper role system in production)
	role := "user"
	if req.Username == "admin" {
		role = "admin"
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(req.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Set expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Return the token
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Username:  req.Username,
		Role:      role,
	})
}

// Logout is a placeholder for logout functionality
// In a stateless JWT system, actual logout happens on the client side by removing the token
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	// Get user information from the context (set by the authentication middleware)
	username, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok {
		roleStr = "user" // Default role
	}

	// Generate new token
	token, err := middleware.GenerateToken(username.(string), roleStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Set expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Return the token
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Username:  username.(string),
		Role:      roleStr,
	})
}
