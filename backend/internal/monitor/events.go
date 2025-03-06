package monitor

import (
	"time"
)

// EventType represents the type of event
type EventType string

// Event types
const (
	WalletDataUpdated  EventType = "wallet_data_updated"
	WalletScanStarted  EventType = "wallet_scan_started"
	WalletScanComplete EventType = "wallet_scan_complete"
	WalletScanError    EventType = "wallet_scan_error"
	TokenChange        EventType = "token_change"
	NewToken           EventType = "new_token"
	TokenRemoved       EventType = "token_removed"
)

// Event represents a monitor event
type Event struct {
	Type      EventType
	Timestamp time.Time
	Payload   map[string]interface{}
}

// EventHandler is a function that handles events
type EventHandler func(event Event)

// EventEmitter manages event publishing and subscription
type EventEmitter struct {
	handlers map[EventType][]EventHandler
}

// Global event emitter instance
var (
	globalEventEmitter *EventEmitter
)

// GetEventEmitter returns the global event emitter instance
func GetEventEmitter() *EventEmitter {
	if globalEventEmitter == nil {
		globalEventEmitter = NewEventEmitter()
	}
	return globalEventEmitter
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a handler for an event type
func (e *EventEmitter) Subscribe(eventType EventType, handler EventHandler) {
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// Emit publishes an event to all registered handlers
func (e *EventEmitter) Emit(eventType EventType, payload map[string]interface{}) {
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	if handlers, exists := e.handlers[eventType]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

// PublishWalletScanStarted publishes a wallet scan started event
func PublishWalletScanStarted() {
	payload := map[string]interface{}{
		"message": "Wallet scan started",
		"time":    time.Now().Format(time.RFC3339),
	}

	GetEventEmitter().Emit(WalletScanStarted, payload)
}

// PublishWalletScanComplete publishes a wallet scan completed event
func PublishWalletScanComplete(walletData map[string]*WalletData) {
	// Create a simplified version of the data for the payload
	simplifiedData := make(map[string]interface{})
	for addr, data := range walletData {
		walletInfo := map[string]interface{}{
			"address":      addr,
			"last_scanned": data.LastScanned,
			"token_count":  len(data.TokenAccounts),
		}
		simplifiedData[addr] = walletInfo
	}

	payload := map[string]interface{}{
		"message": "Wallet scan completed",
		"time":    time.Now().Format(time.RFC3339),
		"wallets": simplifiedData,
	}

	GetEventEmitter().Emit(WalletScanComplete, payload)

	// Also emit individual wallet updates
	for addr, data := range walletData {
		walletPayload := map[string]interface{}{
			"address":       addr,
			"last_scanned":  data.LastScanned,
			"token_count":   len(data.TokenAccounts),
			"token_details": data.TokenAccounts,
		}

		// Create a wallet-specific event
		walletUpdatePayload := map[string]interface{}{
			"wallet": walletPayload,
		}

		GetEventEmitter().Emit(WalletDataUpdated, walletUpdatePayload)
	}
}

// PublishWalletScanError publishes a wallet scan error event
func PublishWalletScanError(err error) {
	payload := map[string]interface{}{
		"message": "Wallet scan failed",
		"error":   err.Error(),
		"time":    time.Now().Format(time.RFC3339),
	}

	GetEventEmitter().Emit(WalletScanError, payload)
}

// PublishTokenChange publishes a token change event
func PublishTokenChange(walletAddress string, tokenMint string, oldBalance uint64, newBalance uint64, tokenInfo TokenAccountInfo) {
	payload := map[string]interface{}{
		"wallet_address": walletAddress,
		"token_mint":     tokenMint,
		"old_balance":    oldBalance,
		"new_balance":    newBalance,
		"decimals":       tokenInfo.Decimals,
		"symbol":         tokenInfo.Symbol,
		"time":           time.Now().Format(time.RFC3339),
	}

	GetEventEmitter().Emit(TokenChange, payload)
}

// GetFormattedTime returns the current time formatted as RFC3339
func GetFormattedTime() string {
	return time.Now().Format(time.RFC3339)
}

// For testability
var (
	nowFunc    = func() time.Time { return time.Now() }
	timeFormat = time.RFC3339
)
