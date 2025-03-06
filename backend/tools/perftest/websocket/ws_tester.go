package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/config"
	"github.com/gorilla/websocket"
)

// MessageResult represents the result of a single WebSocket message
type MessageResult struct {
	Type         string
	SentTime     time.Time
	ReceivedTime time.Time
	Duration     time.Duration
	Size         int
	Error        error
	ErrorMessage string
}

// ConnectionResult represents the result of a WebSocket connection
type ConnectionResult struct {
	ConnectionID     int
	ConnectStartTime time.Time
	ConnectEndTime   time.Time
	ConnectDuration  time.Duration
	DisconnectTime   time.Time
	Connected        bool
	ConnectError     error
	MessagesSent     int
	MessagesReceived int
	MessageResults   []MessageResult
}

// WebSocketTestResult represents the results of a WebSocket test
type WebSocketTestResult struct {
	TotalConnections      int
	SuccessfulConnections int
	FailedConnections     int
	MinConnectTime        time.Duration
	MaxConnectTime        time.Duration
	AverageConnectTime    time.Duration
	TotalMessagesSent     int
	TotalMessagesReceived int
	MinMessageTime        time.Duration
	MaxMessageTime        time.Duration
	AverageMessageTime    time.Duration
	StartTime             time.Time
	EndTime               time.Time
	ConnectionResults     []ConnectionResult
}

// WebSocketLoadTester is responsible for load testing WebSocket functionality
type WebSocketLoadTester struct {
	config *config.Config
}

// NewWebSocketLoadTester creates a new WebSocket load tester
func NewWebSocketLoadTester(cfg *config.Config) *WebSocketLoadTester {
	return &WebSocketLoadTester{
		config: cfg,
	}
}

// TestWebSocketPerformance tests WebSocket performance
func (t *WebSocketLoadTester) TestWebSocketPerformance(ctx context.Context) (*WebSocketTestResult, error) {
	startTime := time.Now()
	results := make([]ConnectionResult, 0)
	var mutex sync.Mutex

	if t.config.Verbose {
		log.Printf("Testing WebSocket with %d concurrent connections for %s",
			t.config.Concurrency, t.config.Duration)
	}

	// Create a WaitGroup to wait for all connections to complete
	var wg sync.WaitGroup

	// Create connections
	for i := 0; i < t.config.Concurrency; i++ {
		wg.Add(1)
		go func(connID int) {
			defer wg.Done()

			// Calculate delay for this connection to implement ramp-up
			var startDelay time.Duration
			if t.config.RampUpPeriod > 0 {
				startDelay = time.Duration(float64(connID) / float64(t.config.Concurrency) * float64(t.config.RampUpPeriod))
				time.Sleep(startDelay)
			}

			// Create connection result
			result := ConnectionResult{
				ConnectionID:     connID,
				ConnectStartTime: time.Now(),
			}

			// Fixed WebSocket connection to use gorilla/websocket Dialer
			// Connect to WebSocket
			dialer := websocket.DefaultDialer
			header := http.Header{}
			header.Set("Origin", "http://localhost")
			conn, _, err := dialer.Dial(t.config.WebSocketURL, header)
			result.ConnectEndTime = time.Now()
			result.ConnectDuration = result.ConnectEndTime.Sub(result.ConnectStartTime)

			if err != nil {
				result.ConnectError = err
				result.Connected = false

				// Add result to results
				mutex.Lock()
				results = append(results, result)
				mutex.Unlock()
				return
			}

			result.Connected = true
			defer conn.Close()

			// Create message channel and error channel
			messageCh := make(chan MessageResult, 100)
			errorCh := make(chan error, 1)

			// Set up receiver goroutine
			go t.receiveMessages(conn, messageCh, errorCh)

			// Send messages
			t.sendMessages(ctx, conn, messageCh, &result)

			// Wait for test duration
			testEndTime := time.Now().Add(t.config.Duration)
			select {
			case <-ctx.Done():
				// Context cancelled
			case err := <-errorCh:
				if err != nil && t.config.Verbose {
					log.Printf("WebSocket error for connection %d: %v", connID, err)
				}
			case <-time.After(time.Until(testEndTime)):
				// Test duration completed
			}

			// Disconnect
			result.DisconnectTime = time.Now()

			// Collect message results
			close(messageCh)
			for msg := range messageCh {
				result.MessageResults = append(result.MessageResults, msg)
				if msg.Error == nil {
					result.MessagesReceived++
				}
			}

			// Add result to results
			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(i)
	}

	// Wait for all connections to complete
	wg.Wait()

	// Calculate statistics
	endTime := time.Now()
	return t.calculateStats(results, startTime, endTime), nil
}

// sendMessages sends messages to the WebSocket connection
func (t *WebSocketLoadTester) sendMessages(ctx context.Context, conn *websocket.Conn, messageCh chan<- MessageResult, result *ConnectionResult) {
	// Calculate message interval based on message rate
	messageInterval := time.Duration(1000.0/t.config.MessageRate) * time.Millisecond

	// Create ticker for sending messages
	ticker := time.NewTicker(messageInterval)
	defer ticker.Stop()

	// Send messages until test duration is reached
	endTime := time.Now().Add(t.config.Duration)
	for time.Now().Before(endTime) {
		select {
		case <-ctx.Done():
			return // Context cancelled
		case <-ticker.C:
			// Create message
			msgType := t.getRandomMessageType()
			message := map[string]interface{}{
				"type": msgType,
				"data": map[string]interface{}{
					"timestamp": time.Now().UnixNano(),
					"id":        result.ConnectionID,
				},
			}

			// Marshal message
			msgBytes, err := json.Marshal(message)
			if err != nil {
				if t.config.Verbose {
					log.Printf("Failed to marshal WebSocket message: %v", err)
				}
				continue
			}

			// Send message
			msgResult := MessageResult{
				Type:     msgType,
				SentTime: time.Now(),
			}

			err = conn.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				msgResult.Error = err
				msgResult.ErrorMessage = err.Error()
				messageCh <- msgResult
				return
			}

			result.MessagesSent++
		}
	}
}

// receiveMessages receives messages from the WebSocket connection
func (t *WebSocketLoadTester) receiveMessages(conn *websocket.Conn, messageCh chan<- MessageResult, errorCh chan<- error) {
	for {
		// Read message
		_, msg, err := conn.ReadMessage()
		now := time.Now()

		if err != nil {
			errorCh <- err
			return
		}

		// Parse message
		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			if t.config.Verbose {
				log.Printf("Failed to unmarshal WebSocket message: %v", err)
				log.Printf("Message content: %s", string(msg))

				// Try to handle multiple JSON objects concatenated together
				parts := splitJSONMessages(string(msg))
				if len(parts) > 1 && t.config.Verbose {
					log.Printf("Detected multiple concatenated JSON messages: %d parts", len(parts))

					// Process each valid part
					for _, part := range parts {
						var partMessage map[string]interface{}
						if err := json.Unmarshal([]byte(part), &partMessage); err != nil {
							if t.config.Verbose {
								log.Printf("Failed to unmarshal part: %v", err)
							}
						} else {
							// Successfully parsed this part, process it as a normal message
							msgType, _ := partMessage["type"].(string)
							msgResult := MessageResult{
								Type:         msgType,
								ReceivedTime: now,
								Size:         len(part),
							}

							// Check if message has a timestamp (for calculating round-trip time)
							if data, ok := partMessage["data"].(map[string]interface{}); ok {
								if timestamp, ok := data["timestamp"].(float64); ok {
									sentTime := time.Unix(0, int64(timestamp))
									msgResult.SentTime = sentTime
									msgResult.Duration = now.Sub(sentTime)
								}
							}

							// Add message result
							messageCh <- msgResult
						}
					}
				}
			}
			continue
		}

		// Create message result
		msgType, _ := message["type"].(string)
		msgResult := MessageResult{
			Type:         msgType,
			ReceivedTime: now,
			Size:         len(msg),
		}

		// Check if message has a timestamp (for calculating round-trip time)
		if data, ok := message["data"].(map[string]interface{}); ok {
			if timestamp, ok := data["timestamp"].(float64); ok {
				sentTime := time.Unix(0, int64(timestamp))
				msgResult.SentTime = sentTime
				msgResult.Duration = now.Sub(sentTime)
			}
		}

		// Add message result
		messageCh <- msgResult
	}
}

// splitJSONMessages attempts to split concatenated JSON messages
// Returns an array of potential JSON objects
func splitJSONMessages(data string) []string {
	var result []string
	var openBraces int
	var start int

	for i, char := range data {
		switch char {
		case '{':
			if openBraces == 0 {
				start = i
			}
			openBraces++
		case '}':
			openBraces--
			if openBraces == 0 {
				// Found a complete JSON object
				result = append(result, data[start:i+1])
			}
		}
	}

	return result
}

// getRandomMessageType returns a random message type from the configuration
func (t *WebSocketLoadTester) getRandomMessageType() string {
	if len(t.config.MessageTypes) == 0 {
		return "ping"
	}
	return t.config.MessageTypes[time.Now().UnixNano()%int64(len(t.config.MessageTypes))]
}

// calculateStats calculates statistics from the connection results
func (t *WebSocketLoadTester) calculateStats(results []ConnectionResult, startTime, endTime time.Time) *WebSocketTestResult {
	if len(results) == 0 {
		return &WebSocketTestResult{
			StartTime: startTime,
			EndTime:   endTime,
		}
	}

	// Initialize result
	testResult := &WebSocketTestResult{
		TotalConnections:  len(results),
		StartTime:         startTime,
		EndTime:           endTime,
		ConnectionResults: results,
		MinConnectTime:    results[0].ConnectDuration,
	}

	// Calculate statistics
	var totalConnectTime time.Duration
	var totalMessageTime time.Duration
	var totalMessages int
	var minMessageTime time.Duration
	var maxMessageTime time.Duration
	var hasMessage bool

	for _, result := range results {
		// Connection statistics
		if result.ConnectError != nil {
			testResult.FailedConnections++
			continue
		}

		testResult.SuccessfulConnections++
		testResult.TotalMessagesSent += result.MessagesSent
		testResult.TotalMessagesReceived += result.MessagesReceived

		// Connection time statistics
		if result.ConnectDuration < testResult.MinConnectTime {
			testResult.MinConnectTime = result.ConnectDuration
		}
		if result.ConnectDuration > testResult.MaxConnectTime {
			testResult.MaxConnectTime = result.ConnectDuration
		}
		totalConnectTime += result.ConnectDuration

		// Message statistics
		for _, msg := range result.MessageResults {
			if msg.Error != nil {
				continue
			}

			totalMessages++
			totalMessageTime += msg.Duration

			if !hasMessage || msg.Duration < minMessageTime {
				minMessageTime = msg.Duration
				hasMessage = true
			}
			if msg.Duration > maxMessageTime {
				maxMessageTime = msg.Duration
			}
		}
	}

	// Calculate averages
	if testResult.SuccessfulConnections > 0 {
		testResult.AverageConnectTime = totalConnectTime / time.Duration(testResult.SuccessfulConnections)
	}
	if totalMessages > 0 {
		testResult.AverageMessageTime = totalMessageTime / time.Duration(totalMessages)
		testResult.MinMessageTime = minMessageTime
		testResult.MaxMessageTime = maxMessageTime
	}

	return testResult
}
