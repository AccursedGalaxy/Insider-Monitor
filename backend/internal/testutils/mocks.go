package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/alerts"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/config"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
)

// MockTime allows controlling time in tests
type MockTime struct {
	CurrentTime time.Time
}

// Now returns the mock current time
func (mt *MockTime) Now() time.Time {
	return mt.CurrentTime
}

// MockEventEmitter is a mock implementation of event emitter for testing
type MockEventEmitter struct {
	Events     []monitor.Event
	Subscribed map[monitor.EventType][]monitor.EventHandler
}

// NewMockEventEmitter creates a new mock event emitter
func NewMockEventEmitter() *MockEventEmitter {
	return &MockEventEmitter{
		Events:     make([]monitor.Event, 0),
		Subscribed: make(map[monitor.EventType][]monitor.EventHandler),
	}
}

// Subscribe adds a handler for an event type
func (m *MockEventEmitter) Subscribe(eventType monitor.EventType, handler monitor.EventHandler) {
	if _, exists := m.Subscribed[eventType]; !exists {
		m.Subscribed[eventType] = make([]monitor.EventHandler, 0)
	}
	m.Subscribed[eventType] = append(m.Subscribed[eventType], handler)
}

// Emit captures an event and notifies subscribers
func (m *MockEventEmitter) Emit(eventType monitor.EventType, payload map[string]interface{}) {
	event := monitor.Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
	m.Events = append(m.Events, event)

	// Notify subscribers
	if handlers, exists := m.Subscribed[eventType]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// MockAlertDelivery captures alerts for testing
type MockAlertDelivery struct {
	Alerts []alerts.Alert
}

// DeliverAlert captures an alert
func (d *MockAlertDelivery) DeliverAlert(alert alerts.Alert) error {
	d.Alerts = append(d.Alerts, alert)
	return nil
}

// FailingAlertDelivery always fails to deliver alerts
type FailingAlertDelivery struct{}

// DeliverAlert always returns an error
func (d *FailingAlertDelivery) DeliverAlert(alert alerts.Alert) error {
	return fmt.Errorf("mock delivery failure")
}

// SetupTestConfig creates a temporary config file for testing
func SetupTestConfig() (*config.Config, string, error) {
	// Create a temporary file
	file, err := ioutil.TempFile("", "config-*.json")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}

	path := file.Name()
	file.Close()

	// Create a test config by loading from the new file path
	// This will create a default config at the specified path
	testConfig, err := config.LoadConfig(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load test config: %w", err)
	}

	return testConfig, path, nil
}

// CleanupTestConfig removes a temporary config file
func CleanupTestConfig(path string) error {
	return os.Remove(path)
}

// CreateMockWalletData creates mock wallet data for testing
func CreateMockWalletData() map[string]*monitor.WalletData {
	return map[string]*monitor.WalletData{
		"55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr": {
			WalletAddress: "55kBY9yxqQzj2zxZqRkqENYq6R8PkXmn5GKyQN9YeVFr",
			TokenAccounts: map[string]monitor.TokenAccountInfo{
				"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": {
					Balance:     1000000000,
					Decimals:    6,
					Symbol:      "USDC",
					LastUpdated: time.Now(),
				},
				"So11111111111111111111111111111111111111112": {
					Balance:     5000000000,
					Decimals:    9,
					Symbol:      "SOL",
					LastUpdated: time.Now(),
				},
			},
			LastScanned: time.Now(),
		},
		"DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF": {
			WalletAddress: "DWuopnuSqYdBhCXqxfqjqzPGibnhkj6SQqFvgC4jkvjF",
			TokenAccounts: map[string]monitor.TokenAccountInfo{
				"So11111111111111111111111111111111111111112": {
					Balance:     10000000000,
					Decimals:    9,
					Symbol:      "SOL",
					LastUpdated: time.Now(),
				},
			},
			LastScanned: time.Now(),
		},
	}
}
