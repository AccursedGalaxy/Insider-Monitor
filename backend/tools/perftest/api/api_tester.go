package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/config"
)

// RequestResult represents the result of a single API request
type RequestResult struct {
	URL          string
	Method       string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	StatusCode   int
	ResponseSize int64
	Error        error
	ErrorMessage string
}

// APITestResult represents the results of an API endpoint test
type APITestResult struct {
	URL                 string
	Method              string
	TotalRequests       int
	SuccessfulRequests  int
	FailedRequests      int
	MinResponseTime     time.Duration
	MaxResponseTime     time.Duration
	AverageResponseTime time.Duration
	P50ResponseTime     time.Duration // 50th percentile response time
	P90ResponseTime     time.Duration // 90th percentile response time
	P95ResponseTime     time.Duration // 95th percentile response time
	P99ResponseTime     time.Duration // 99th percentile response time
	TotalBytes          int64
	StartTime           time.Time
	EndTime             time.Time
	ErrorCodes          map[int]int     // Map of status code to count
	Errors              []string        // List of error messages
	RawResults          []RequestResult // Individual request results
}

// APILoadTester is responsible for load testing API endpoints
type APILoadTester struct {
	config    *config.Config
	client    *http.Client
	authToken string
}

// NewAPILoadTester creates a new API load tester
func NewAPILoadTester(cfg *config.Config) *APILoadTester {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &APILoadTester{
		config: cfg,
		client: client,
	}
}

// TestEndpoint tests a specific API endpoint
func (t *APILoadTester) TestEndpoint(ctx context.Context, endpoint config.Endpoint) (*APITestResult, error) {
	start := time.Now()
	results := make([]RequestResult, 0)
	var mutex sync.Mutex

	// Create a rate limiter to control request rate during ramp-up
	var requestRate float64
	durationSecs := t.config.Duration.Seconds()
	if t.config.Verbose {
		log.Printf("Duration seconds: %f", durationSecs)
	}
	if durationSecs <= 0 {
		// Default to 1 request per second per user if duration is invalid
		requestRate = float64(t.config.Concurrency)
		if t.config.Verbose {
			log.Printf("Using default request rate: %f requests/sec", requestRate)
		}
	} else {
		requestRate = float64(t.config.Concurrency) / durationSecs
		if t.config.Verbose {
			log.Printf("Calculated request rate: %f requests/sec", requestRate)
		}
	}

	// Ensure we have a reasonable minimum rate
	if requestRate <= 0 {
		requestRate = 1.0 // Default to at least 1 request per second total
		if t.config.Verbose {
			log.Printf("Adjusted to minimum request rate: %f requests/sec", requestRate)
		}
	}

	// Calculate how many requests to make
	totalDuration := t.config.Duration + t.config.RampUpPeriod
	totalRequests := int(requestRate * totalDuration.Seconds())

	if t.config.Verbose {
		log.Printf("Testing endpoint %s %s with %d concurrent users for %s (total requests: ~%d)",
			endpoint.Method, endpoint.URL, t.config.Concurrency, totalDuration, totalRequests)
	}

	// Create a WaitGroup to wait for all requests to complete
	var wg sync.WaitGroup

	// Create worker pool
	for i := 0; i < t.config.Concurrency; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			// Calculate delay for this worker to implement ramp-up
			var startDelay time.Duration
			if t.config.RampUpPeriod > 0 {
				startDelay = time.Duration(float64(workerId) / float64(t.config.Concurrency) * float64(t.config.RampUpPeriod))
				time.Sleep(startDelay)
			}

			// Main worker loop
			var tickInterval time.Duration
			workerRate := requestRate / float64(t.config.Concurrency)
			if t.config.Verbose && workerId == 0 {
				log.Printf("Worker rate: %f, concurrency: %d", workerRate, t.config.Concurrency)
			}
			if workerRate <= 0 {
				// Default to one request per second if rate is invalid
				tickInterval = time.Second
				if t.config.Verbose && workerId == 0 {
					log.Printf("Using default tick interval of 1 second")
				}
			} else {
				// Calculate interval in nanoseconds directly to avoid division by zero
				intervalNanos := float64(time.Second) / workerRate
				tickInterval = time.Duration(intervalNanos)
				if t.config.Verbose && workerId == 0 {
					log.Printf("Calculated tick interval: %v", tickInterval)
				}
			}
			ticker := time.NewTicker(tickInterval)
			defer ticker.Stop()

			endTime := time.Now().Add(t.config.Duration)
			for time.Now().Before(endTime) {
				select {
				case <-ctx.Done():
					return // Context cancelled
				case <-ticker.C:
					// Make a request
					result := t.makeRequest(endpoint)

					// Add result to the results slice
					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
				}
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()

	// Calculate statistics
	endTime := time.Now()
	return t.calculateStats(endpoint, results, start, endTime), nil
}

// makeRequest makes a single HTTP request to the specified endpoint
func (t *APILoadTester) makeRequest(endpoint config.Endpoint) RequestResult {
	result := RequestResult{
		URL:       endpoint.URL,
		Method:    endpoint.Method,
		StartTime: time.Now(),
	}

	// Create URL with base URL
	url := t.config.BaseURL + endpoint.URL

	// Create request
	var reqBody []byte
	var err error
	if endpoint.Body != nil {
		reqBody, err = json.Marshal(endpoint.Body)
		if err != nil {
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			result.Error = err
			result.ErrorMessage = fmt.Sprintf("Failed to marshal request body: %v", err)
			return result
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(endpoint.Method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Error = err
		result.ErrorMessage = fmt.Sprintf("Failed to create request: %v", err)
		return result
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	// Add authentication if required
	if endpoint.RequiresAuth && t.config.AuthEnabled {
		switch t.config.AuthType {
		case "basic":
			req.SetBasicAuth(t.config.Username, t.config.Password)
		case "jwt":
			req.Header.Set("Authorization", "Bearer "+t.config.JWTToken)
		case "api_key":
			req.Header.Set("X-API-Key", t.config.APIKey)
		}
	}

	// Make the request
	resp, err := t.client.Do(req)
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if err != nil {
		result.Error = err
		result.ErrorMessage = fmt.Sprintf("Request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Record result
	result.StatusCode = resp.StatusCode

	// Get response size
	if resp.ContentLength > 0 {
		result.ResponseSize = resp.ContentLength
	} else {
		// If ContentLength is not provided, we could read and count bytes
		// but for performance testing purposes, this might be unnecessary overhead
		result.ResponseSize = 0
	}

	return result
}

// calculateStats calculates statistics from the request results
func (t *APILoadTester) calculateStats(endpoint config.Endpoint, results []RequestResult, startTime, endTime time.Time) *APITestResult {
	if len(results) == 0 {
		return &APITestResult{
			URL:        endpoint.URL,
			Method:     endpoint.Method,
			StartTime:  startTime,
			EndTime:    endTime,
			ErrorCodes: make(map[int]int),
		}
	}

	// Initialize result
	testResult := &APITestResult{
		URL:             endpoint.URL,
		Method:          endpoint.Method,
		TotalRequests:   len(results),
		StartTime:       startTime,
		EndTime:         endTime,
		ErrorCodes:      make(map[int]int),
		RawResults:      results,
		MinResponseTime: results[0].Duration,
	}

	// Calculate statistics
	var totalDuration time.Duration
	var totalBytes int64
	durations := make([]time.Duration, 0, len(results))

	for _, result := range results {
		durations = append(durations, result.Duration)

		if result.Error != nil {
			testResult.FailedRequests++
			testResult.Errors = append(testResult.Errors, result.ErrorMessage)
			continue
		}

		if result.StatusCode >= 200 && result.StatusCode < 400 {
			testResult.SuccessfulRequests++
		} else {
			testResult.FailedRequests++
		}

		// Track status codes
		testResult.ErrorCodes[result.StatusCode]++

		// Update min/max response time
		if result.Duration < testResult.MinResponseTime {
			testResult.MinResponseTime = result.Duration
		}
		if result.Duration > testResult.MaxResponseTime {
			testResult.MaxResponseTime = result.Duration
		}

		// Accumulate total duration and bytes
		totalDuration += result.Duration
		totalBytes += result.ResponseSize
	}

	// Calculate average response time
	if len(results) > 0 {
		testResult.AverageResponseTime = totalDuration / time.Duration(len(results))
	}

	// Calculate percentiles
	if len(durations) > 0 {
		// Sort durations
		rand.Shuffle(len(durations), func(i, j int) {
			durations[i], durations[j] = durations[j], durations[i]
		})

		testResult.P50ResponseTime = percentile(durations, 50)
		testResult.P90ResponseTime = percentile(durations, 90)
		testResult.P95ResponseTime = percentile(durations, 95)
		testResult.P99ResponseTime = percentile(durations, 99)
	}

	testResult.TotalBytes = totalBytes

	return testResult
}

// percentile calculates the p-th percentile of durations
func percentile(durations []time.Duration, p int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Sort durations (assuming they're not sorted)
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	if p <= 0 {
		return durations[0]
	}
	if p >= 100 {
		return durations[len(durations)-1]
	}

	// Calculate index
	idx := (p * len(durations)) / 100
	if idx >= len(durations) {
		idx = len(durations) - 1
	}

	return durations[idx]
}
