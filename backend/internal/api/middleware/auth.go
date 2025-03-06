package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey     string
	TokenExpiry   time.Duration
	TokenIssuer   string
	TokenAudience string
}

// Default JWT configuration
var defaultJWTConfig = JWTConfig{
	SecretKey:     "your-secret-key", // In production, use a proper secret key from environment variables
	TokenExpiry:   time.Hour * 24,
	TokenIssuer:   "insider-monitor",
	TokenAudience: "insider-monitor-api",
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Authenticate middleware checks for valid authentication
func Authenticate() gin.HandlerFunc {
	return authenticateWithConfig(defaultJWTConfig)
}

// AuthenticateWithConfig allows custom JWT configuration
func AuthenticateWithConfig(config JWTConfig) gin.HandlerFunc {
	return authenticateWithConfig(config)
}

// authenticateWithConfig is the internal implementation
func authenticateWithConfig(config JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if Authorization header exists and has the correct format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.SecretKey), nil
		})

		// Handle parsing errors
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			return
		}

		// Validate claims
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Set user information in the context
			c.Set("user", claims.Username)
			c.Set("role", claims.Role)

			// Continue to the next handler
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}
	}
}

// Rate limiting middleware
func RateLimit(maxRequests int, duration time.Duration) gin.HandlerFunc {
	// Create a map to store request counts by IP
	type client struct {
		count    int
		lastSeen time.Time
	}

	clients := make(map[string]*client)

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Check if client exists
		cl, exists := clients[ip]
		if !exists {
			// Create new client
			clients[ip] = &client{
				count:    1,
				lastSeen: time.Now(),
			}
			c.Next()
			return
		}

		// Check if duration has passed
		if time.Since(cl.lastSeen) > duration {
			// Reset count
			cl.count = 1
			cl.lastSeen = time.Now()
			c.Next()
			return
		}

		// Check if count exceeds max requests
		if cl.count >= maxRequests {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		// Increment count
		cl.count++
		cl.lastSeen = time.Now()

		c.Next()
	}
}

// GenerateToken creates a new JWT token
func GenerateToken(username, role string) (string, error) {
	return GenerateTokenWithConfig(username, role, defaultJWTConfig)
}

// GenerateTokenWithConfig creates a new JWT token with custom configuration
func GenerateTokenWithConfig(username, role string, config JWTConfig) (string, error) {
	// Create the claims
	claims := JWTClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.TokenIssuer,
			Subject:   username,
			Audience:  []string{config.TokenAudience},
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString([]byte(config.SecretKey))
}
