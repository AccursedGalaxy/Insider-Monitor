package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Define context key type to avoid collisions
type contextKey string

// Define keys for context values
const (
	ClaimsContextKey contextKey = "claims"
)

// User represents a user of the system
type User struct {
	Username string
	Password string
	Role     string
}

// Auth handles authentication and authorization
type Auth struct {
	users     map[string]User
	jwtSecret []byte
}

// New creates a new Auth instance
func New() *Auth {
	// In a real application, these would come from a database
	// For now, we'll use a simple map with a default admin user
	defaultPassword := os.Getenv("ADMIN_PASSWORD")
	if defaultPassword == "" {
		defaultPassword = "admin" // Default password if not set in environment
	}

	users := map[string]User{
		"admin": {
			Username: "admin",
			Password: defaultPassword,
			Role:     "admin",
		},
	}

	// Use environment variable for JWT secret or default
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		// For development only - in production use a proper secret
		jwtSecret = []byte("insider-monitor-secret-key")
	}

	return &Auth{
		users:     users,
		jwtSecret: jwtSecret,
	}
}

// GenerateToken creates a new JWT token for a user
func (a *Auth) GenerateToken(username, password string) (string, error) {
	user, exists := a.users[username]
	if !exists || subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) != 1 {
		return "", fmt.Errorf("invalid credentials")
	}

	// Create the JWT claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (a *Auth) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// BasicAuthMiddleware provides HTTP Basic Auth protection
func (a *Auth) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, return 401 Unauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if it's Basic auth
		if !strings.HasPrefix(authHeader, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get the credentials
		payload, err := base64.StdEncoding.DecodeString(authHeader[6:])
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split into username:password
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, password := pair[0], pair[1]

		// Check if username exists and password matches
		user, exists := a.users[username]
		if !exists || subtle.ConstantTimeCompare([]byte(user.Password), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Authentication successful, call the next handler
		next.ServeHTTP(w, r)
	})
}

// JWTAuthMiddleware provides JWT token protection
func (a *Auth) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if it's Bearer auth
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := authHeader[7:]

		// Validate the token
		claims, err := a.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in the request context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

		// Call the next handler with the augmented context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
