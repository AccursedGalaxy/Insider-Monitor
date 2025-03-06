package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LogLevel represents the severity of the log message
type LogLevel string

const (
	// Debug level for development logging
	Debug LogLevel = "DEBUG"
	// Info level for informational messages
	Info LogLevel = "INFO"
	// Warn level for warning messages
	Warn LogLevel = "WARN"
	// Error level for error messages
	Error LogLevel = "ERROR"
)

// Logger contains the logger instance and options
type Logger struct {
	logger *log.Logger
	level  LogLevel
}

var (
	defaultLogger *Logger
)

// init initializes the default logger
func init() {
	defaultLogger = NewLogger(os.Stdout, Info)
}

// NewLogger creates a new logger with the specified output and level
func NewLogger(output io.Writer, level LogLevel) *Logger {
	return &Logger{
		logger: log.New(output, "", log.LstdFlags),
		level:  level,
	}
}

// SetLevel sets the log level for the logger
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// logMessage logs a message with the specified level
func (l *Logger) logMessage(level LogLevel, message string) {
	// Only log if the message level is at or above the logger level
	if shouldLog(l.level, level) {
		l.logger.Printf("[%s] %s", level, message)
	}
}

// shouldLog determines if a message should be logged based on the logger level
func shouldLog(loggerLevel, messageLevel LogLevel) bool {
	// Order of severity
	levels := map[LogLevel]int{
		Debug: 0,
		Info:  1,
		Warn:  2,
		Error: 3,
	}

	return levels[messageLevel] >= levels[loggerLevel]
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logMessage(Debug, fmt.Sprintf(format, v...))
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.logMessage(Info, fmt.Sprintf(format, v...))
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.logMessage(Warn, fmt.Sprintf(format, v...))
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.logMessage(Error, fmt.Sprintf(format, v...))
}

// RequestLogger middleware logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Generate request ID if not already present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set("X-Request-ID", requestID)
		}

		// Save request ID in context for use in other parts of the application
		c.Set("RequestID", requestID)

		// Read request body if needed
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create buffer for response
		responseBody := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = responseBody

		// Process request
		c.Next()

		// Collect response status
		statusCode := c.Writer.Status()
		statusClass := statusCode / 100

		// Calculate latency
		latency := time.Since(startTime)

		// Determine log level based on status
		var level LogLevel
		switch statusClass {
		case 2:
			level = Info
		case 3:
			level = Info
		case 4:
			level = Warn
		case 5:
			level = Error
		default:
			level = Info
		}

		// Log request details
		defaultLogger.logMessage(
			level,
			fmt.Sprintf(
				"[%s] %s %s %d %s",
				requestID,
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				latency,
			),
		)

		// Log request/response bodies in debug mode
		if shouldLog(defaultLogger.level, Debug) {
			// Truncate large bodies to avoid flooding logs
			const maxBodyLogSize = 1024 // 1KB

			// Log request body if present
			if len(requestBody) > 0 {
				displayBody := string(requestBody)
				if len(displayBody) > maxBodyLogSize {
					displayBody = displayBody[:maxBodyLogSize] + "... [truncated]"
				}
				defaultLogger.Debug("[%s] Request Body: %s", requestID, displayBody)
			}

			// Log response body if present
			if responseBody.body.Len() > 0 {
				displayBody := responseBody.body.String()
				if len(displayBody) > maxBodyLogSize {
					displayBody = displayBody[:maxBodyLogSize] + "... [truncated]"
				}
				defaultLogger.Debug("[%s] Response Body: %s", requestID, displayBody)
			}
		}
	}
}

// responseBodyWriter captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body before writing it
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// WriteString captures the response body as a string before writing it
func (r *responseBodyWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

// RequestMetrics middleware collects request metrics
func RequestMetrics() gin.HandlerFunc {
	// Simple in-memory metrics for now
	metrics := struct {
		totalRequests    int
		requestsByRoute  map[string]int
		requestsByStatus map[int]int
		averageLatency   time.Duration
		maxLatency       time.Duration
		minLatency       time.Duration
		mutex            sync.RWMutex
	}{
		requestsByRoute:  make(map[string]int),
		requestsByStatus: make(map[int]int),
		minLatency:       time.Hour, // Initialize with large value
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Get request details
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		statusCode := c.Writer.Status()
		latency := time.Since(startTime)

		// Update metrics
		metrics.mutex.Lock()
		defer metrics.mutex.Unlock()

		metrics.totalRequests++
		metrics.requestsByRoute[route]++
		metrics.requestsByStatus[statusCode]++

		// Update latency stats
		if latency > metrics.maxLatency {
			metrics.maxLatency = latency
		}
		if latency < metrics.minLatency {
			metrics.minLatency = latency
		}

		// Calculate running average
		currentAvg := metrics.averageLatency
		n := metrics.totalRequests
		metrics.averageLatency = (currentAvg*time.Duration(n-1) + latency) / time.Duration(n)
	}
}
