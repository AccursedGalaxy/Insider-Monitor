package handlers

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemMetrics holds collected system and application metrics
type SystemMetrics struct {
	// System metrics
	NumGoroutine    int       `json:"num_goroutines"`
	NumCPU          int       `json:"num_cpu"`
	AllocatedMemory uint64    `json:"allocated_memory_bytes"`
	TotalAllocated  uint64    `json:"total_allocated_bytes"`
	SystemMemory    uint64    `json:"system_memory_bytes"`
	NumGC           uint32    `json:"num_gc"`
	Uptime          float64   `json:"uptime_seconds"`
	StartTime       time.Time `json:"start_time"`

	// Request metrics
	TotalRequests     int                    `json:"total_requests"`
	RequestsByRoute   map[string]int         `json:"requests_by_route"`
	RequestsByStatus  map[string]int         `json:"requests_by_status"`
	AverageLatency    string                 `json:"average_latency"`
	MaxLatency        string                 `json:"max_latency"`
	MinLatency        string                 `json:"min_latency"`
	ActiveConnections int                    `json:"active_connections"`
	ErrorRate         float64                `json:"error_rate"`
	RequestsPerMinute float64                `json:"requests_per_minute"`
	CustomMetrics     map[string]interface{} `json:"custom_metrics,omitempty"`
}

// Application start time for calculating uptime
var startTime = time.Now()

// Global metrics store for the demo - in a real app, use a proper metrics system
var requestMetrics = struct {
	totalRequests     int
	requestsByRoute   map[string]int
	requestsByStatus  map[int]int // Status codes as integers
	averageLatency    time.Duration
	maxLatency        time.Duration
	minLatency        time.Duration
	activeConnections int
}{
	requestsByRoute:  make(map[string]int),
	requestsByStatus: make(map[int]int),
	minLatency:       time.Hour, // Initialize with large value
}

// GetMetrics returns system and application metrics
func GetMetrics(c *gin.Context) {
	// Collect runtime metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Convert request status codes from int to string for JSON
	requestsByStatus := make(map[string]int)
	for status, count := range requestMetrics.requestsByStatus {
		statusText := http.StatusText(status)
		if statusText == "" {
			statusText = strconv.Itoa(status)
		}
		requestsByStatus[statusText] = count
	}

	// Calculate error rate
	errorCount := 0
	for status, count := range requestMetrics.requestsByStatus {
		if status >= 400 {
			errorCount += count
		}
	}
	errorRate := 0.0
	if requestMetrics.totalRequests > 0 {
		errorRate = float64(errorCount) / float64(requestMetrics.totalRequests)
	}

	// Calculate requests per minute
	uptime := time.Since(startTime).Seconds()
	requestsPerMinute := 0.0
	if uptime > 0 {
		requestsPerMinute = float64(requestMetrics.totalRequests) / uptime * 60.0
	}

	// Create metrics response
	metrics := SystemMetrics{
		// System metrics
		NumGoroutine:    runtime.NumGoroutine(),
		NumCPU:          runtime.NumCPU(),
		AllocatedMemory: memStats.Alloc,
		TotalAllocated:  memStats.TotalAlloc,
		SystemMemory:    memStats.Sys,
		NumGC:           memStats.NumGC,
		Uptime:          uptime,
		StartTime:       startTime,

		// Request metrics
		TotalRequests:     requestMetrics.totalRequests,
		RequestsByRoute:   requestMetrics.requestsByRoute,
		RequestsByStatus:  requestsByStatus,
		AverageLatency:    requestMetrics.averageLatency.String(),
		MaxLatency:        requestMetrics.maxLatency.String(),
		MinLatency:        requestMetrics.minLatency.String(),
		ActiveConnections: requestMetrics.activeConnections,
		ErrorRate:         errorRate,
		RequestsPerMinute: requestsPerMinute,
		CustomMetrics:     make(map[string]interface{}),
	}

	// Add custom metrics
	metrics.CustomMetrics["websocket_connections"] = 0 // Placeholder - integrate with actual WebSocket metrics

	c.JSON(http.StatusOK, metrics)
}
