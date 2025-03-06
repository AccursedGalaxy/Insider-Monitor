package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	// Test creating a new logger
	buf := new(bytes.Buffer)
	logger := NewLogger(buf, Info)

	// Verify logger is created
	assert.NotNil(t, logger)
}

func TestRequestLogger(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router with our logger middleware
	r := gin.New()
	r.Use(RequestLogger())

	// Add a test route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test successful"})
	})

	// Create a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check response - we can only verify the status code
	// since the logging output goes to the default logger
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogLevels(t *testing.T) {
	// Create a buffer to capture log output
	buf := new(bytes.Buffer)
	logger := NewLogger(buf, Info)

	// Log at different levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")

	// Verify output based on level
	output := buf.String()

	// Debug should not be logged when level is Info
	assert.NotContains(t, output, "DEBUG")
	assert.NotContains(t, output, "debug message")

	// Info and higher should be logged
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "WARN")
	assert.Contains(t, output, "warning message")
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "error message")
}

func TestRequestWithPathParams(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router with our logger middleware
	r := gin.New()
	r.Use(RequestLogger())

	// Add a test route with a path parameter
	r.GET("/users/:id", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond) // simulate some processing time
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	})

	// Create a request
	req := httptest.NewRequest("GET", "/users/123?name=test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)

	// The actual logging output is hard to test because it goes to the default logger
	// So we only check the response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "123", response["id"])
}
