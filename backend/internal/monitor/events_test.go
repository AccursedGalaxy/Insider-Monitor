package monitor

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// MockEventEmitter extends EventEmitter to capture events for testing
type MockEventEmitter struct {
	*EventEmitter
	CapturedEvents []Event
	mutex          sync.Mutex
}

// NewMockEventEmitter creates a new mock event emitter
func NewMockEventEmitter() *MockEventEmitter {
	mock := &MockEventEmitter{
		EventEmitter:   NewEventEmitter(),
		CapturedEvents: make([]Event, 0),
	}

	return mock
}

// Emit captures an event and notifies subscribers
func (m *MockEventEmitter) Emit(eventType EventType, payload map[string]interface{}) {
	// Call the underlying implementation
	m.EventEmitter.Emit(eventType, payload)

	// Additionally capture the event for testing
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	m.mutex.Lock()
	m.CapturedEvents = append(m.CapturedEvents, event)
	m.mutex.Unlock()
}

// Test helper to directly inject and capture events
func captureEvents(t *testing.T, fn func()) []*Event {
	// Save original emitter
	origEmitter := globalEventEmitter

	// Create a channel to receive events
	eventChan := make(chan *Event, 10)

	// Create a new emitter for testing
	testEmitter := NewEventEmitter()

	// Hook up all event types that we're interested in
	eventTypes := []EventType{
		WalletScanStarted,
		WalletScanComplete,
		WalletScanError,
		TokenChange,
		WalletDataUpdated,
		NewToken,
		TokenRemoved,
	}

	for _, et := range eventTypes {
		// Capture each event type with a stable variable
		eventType := et
		testEmitter.Subscribe(eventType, func(event Event) {
			eventCopy := event // Make a copy to avoid data races
			eventChan <- &eventCopy
			fmt.Printf("Event captured: %s\n", eventType)
		})
	}

	// Replace the global emitter
	globalEventEmitter = testEmitter

	// Run the function that should emit events
	fn()

	// Allow some time for events to be processed
	time.Sleep(100 * time.Millisecond)

	// Restore the global emitter
	globalEventEmitter = origEmitter

	// Collect all events from channel
	var events []*Event
	close(eventChan)
	for ev := range eventChan {
		events = append(events, ev)
	}

	return events
}

func TestGetEventEmitter(t *testing.T) {
	// Reset the global event emitter
	originalEmitter := globalEventEmitter
	globalEventEmitter = nil
	defer func() { globalEventEmitter = originalEmitter }()

	// Get the event emitter
	emitter := GetEventEmitter()
	if emitter == nil {
		t.Fatal("Expected emitter to be non-nil")
	}

	// Get it again - should be the same instance
	emitter2 := GetEventEmitter()
	if emitter != emitter2 {
		t.Error("Expected GetEventEmitter to return the same instance on subsequent calls")
	}
}

func TestEventEmitterSubscribeAndEmit(t *testing.T) {
	// Create a new event emitter
	emitter := NewEventEmitter()

	// Create a channel to receive events
	received := make(chan Event, 1)

	// Subscribe to events
	emitter.Subscribe(WalletScanStarted, func(event Event) {
		received <- event
	})

	// Emit an event
	payload := map[string]interface{}{
		"test": "value",
	}
	emitter.Emit(WalletScanStarted, payload)

	// Wait for the event to be received
	select {
	case event := <-received:
		if event.Type != WalletScanStarted {
			t.Errorf("Expected event type %s, got %s", WalletScanStarted, event.Type)
		}

		if value, ok := event.Payload["test"]; !ok || value != "value" {
			t.Errorf("Expected payload to contain test=value, got %v", event.Payload)
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for event")
	}
}

func TestPublishWalletScanStarted(t *testing.T) {
	events := captureEvents(t, func() {
		PublishWalletScanStarted()
	})

	// Check if we captured the event
	foundEvent := false
	for _, event := range events {
		if event.Type == WalletScanStarted {
			foundEvent = true
			if _, ok := event.Payload["message"]; !ok {
				t.Error("Expected message in payload")
			}
			if _, ok := event.Payload["time"]; !ok {
				t.Error("Expected time in payload")
			}
		}
	}

	if !foundEvent {
		t.Error("Expected WalletScanStarted event to be published")
	}
}

func TestPublishWalletScanComplete(t *testing.T) {
	// Create mock wallet data
	walletData := map[string]*WalletData{
		"test-wallet": {
			WalletAddress: "test-wallet",
			TokenAccounts: map[string]TokenAccountInfo{
				"token1": {
					Balance:  1000,
					Symbol:   "TKN1",
					Decimals: 6,
				},
			},
			LastScanned: time.Now(),
		},
	}

	events := captureEvents(t, func() {
		PublishWalletScanComplete(walletData)
	})

	// Check if we captured the event
	foundEvent := false
	for _, event := range events {
		if event.Type == WalletScanComplete {
			foundEvent = true
			if _, ok := event.Payload["message"]; !ok {
				t.Error("Expected message in payload")
			}
			if _, ok := event.Payload["time"]; !ok {
				t.Error("Expected time in payload")
			}
			if wallets, ok := event.Payload["wallets"].(map[string]interface{}); !ok {
				t.Error("Expected wallets in payload")
			} else if _, ok := wallets["test-wallet"]; !ok {
				t.Error("Expected test-wallet in wallets")
			}
		}
	}

	if !foundEvent {
		t.Error("Expected WalletScanComplete event to be published")
	}
}

func TestPublishWalletScanError(t *testing.T) {
	// Create a test error
	testError := errors.New("test error")

	events := captureEvents(t, func() {
		PublishWalletScanError(testError)
	})

	// Check if we captured the event
	foundEvent := false
	for _, event := range events {
		if event.Type == WalletScanError {
			foundEvent = true
			if _, ok := event.Payload["message"]; !ok {
				t.Error("Expected message in payload")
			}
			if _, ok := event.Payload["time"]; !ok {
				t.Error("Expected time in payload")
			}
			if errorMsg, ok := event.Payload["error"].(string); !ok {
				t.Error("Expected error in payload")
			} else if errorMsg != "test error" {
				t.Errorf("Expected error to be 'test error', got '%s'", errorMsg)
			}
		}
	}

	if !foundEvent {
		t.Error("Expected WalletScanError event to be published")
	}
}

func TestPublishTokenChange(t *testing.T) {
	// Create test data
	wallet := "test-wallet"
	token := "test-token"
	oldBalance := uint64(1000)
	newBalance := uint64(2000)
	tokenInfo := TokenAccountInfo{
		Balance:  newBalance,
		Symbol:   "TKN",
		Decimals: 6,
	}

	events := captureEvents(t, func() {
		PublishTokenChange(wallet, token, oldBalance, newBalance, tokenInfo)
	})

	// Check if we captured the event
	foundEvent := false
	for _, event := range events {
		if event.Type == TokenChange {
			foundEvent = true
			if walletAddr, ok := event.Payload["wallet_address"].(string); !ok {
				t.Error("Expected wallet_address in payload")
			} else if walletAddr != wallet {
				t.Errorf("Expected wallet_address to be '%s', got '%s'", wallet, walletAddr)
			}

			if tokenMint, ok := event.Payload["token_mint"].(string); !ok {
				t.Error("Expected token_mint in payload")
			} else if tokenMint != token {
				t.Errorf("Expected token_mint to be '%s', got '%s'", token, tokenMint)
			}

			if old, ok := event.Payload["old_balance"].(uint64); !ok {
				t.Error("Expected old_balance in payload")
			} else if old != oldBalance {
				t.Errorf("Expected old_balance to be %d, got %d", oldBalance, old)
			}

			if new, ok := event.Payload["new_balance"].(uint64); !ok {
				t.Error("Expected new_balance in payload")
			} else if new != newBalance {
				t.Errorf("Expected new_balance to be %d, got %d", newBalance, new)
			}
		}
	}

	if !foundEvent {
		t.Error("Expected TokenChange event to be published")
	}
}
