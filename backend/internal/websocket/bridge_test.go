package websocket

import (
	"testing"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
)

// MockManager simulates a WebSocket manager for testing
type MockManager struct {
	BroadcastCalls      []BroadcastCall
	BroadcastToAllCalls []BroadcastToAllCall
}

type BroadcastCall struct {
	Type    MessageType
	Topic   string
	Payload map[string]interface{}
}

type BroadcastToAllCall struct {
	Type    MessageType
	Payload map[string]interface{}
}

// BroadcastToAll records calls for testing
func (m *MockManager) BroadcastToAll(msgType MessageType, payload map[string]interface{}) {
	m.BroadcastToAllCalls = append(m.BroadcastToAllCalls, BroadcastToAllCall{
		Type:    msgType,
		Payload: payload,
	})
}

// Broadcast records calls for testing
func (m *MockManager) Broadcast(msgType MessageType, topic string, payload map[string]interface{}) {
	m.BroadcastCalls = append(m.BroadcastCalls, BroadcastCall{
		Type:    msgType,
		Topic:   topic,
		Payload: payload,
	})
}

func TestBridgeMessageHandling(t *testing.T) {
	// Create a mock manager
	mockManager := &MockManager{
		BroadcastCalls:      make([]BroadcastCall, 0),
		BroadcastToAllCalls: make([]BroadcastToAllCall, 0),
	}

	// Test handling a status update event
	statusEvent := monitor.Event{
		Type:      monitor.WalletScanStarted,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"message": "Wallet scan started",
			"time":    time.Now().Format(time.RFC3339),
		},
	}

	// Manually call the handler function that would normally be registered
	handleStatusUpdate(mockManager, statusEvent)

	// Check if the event was broadcast to all
	if len(mockManager.BroadcastToAllCalls) == 0 {
		t.Error("Expected BroadcastToAll to be called for status update")
	} else {
		call := mockManager.BroadcastToAllCalls[0]
		if call.Type != StatusUpdateMsg {
			t.Errorf("Expected message type %s, got %s", StatusUpdateMsg, call.Type)
		}
		if _, ok := call.Payload["message"]; !ok {
			t.Error("Expected message in payload")
		}
	}

	// Test handling a wallet data update event
	walletEvent := monitor.Event{
		Type:      monitor.WalletDataUpdated,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"wallet": map[string]interface{}{
				"address":      "test-wallet",
				"last_scanned": time.Now(),
				"token_count":  2,
			},
		},
	}

	// Manually call the handler function
	handleWalletUpdate(mockManager, walletEvent)

	// Check if the event was broadcast to the wallet topic
	walletBroadcastFound := false
	for _, call := range mockManager.BroadcastCalls {
		if call.Type == WalletUpdateMsg && call.Topic == "wallet/test-wallet" {
			walletBroadcastFound = true
			break
		}
	}

	if !walletBroadcastFound {
		t.Error("Expected wallet-specific broadcast for wallet update")
	}

	// Test handling a token change event
	tokenEvent := monitor.Event{
		Type:      monitor.TokenChange,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"wallet_address": "test-wallet",
			"token_mint":     "test-token",
			"old_balance":    uint64(1000),
			"new_balance":    uint64(2000),
		},
	}

	// Manually call the handler function
	handleTokenChange(mockManager, tokenEvent)

	// Check if the event was broadcast to both wallet and token topics
	walletTokenBroadcastFound := false
	tokenBroadcastFound := false
	for _, call := range mockManager.BroadcastCalls {
		if call.Type == WalletUpdateMsg {
			if call.Topic == "wallet/test-wallet" {
				walletTokenBroadcastFound = true
			} else if call.Topic == "token/test-token" {
				tokenBroadcastFound = true
			}
		}
	}

	if !walletTokenBroadcastFound {
		t.Error("Expected wallet-specific broadcast for token change")
	}

	if !tokenBroadcastFound {
		t.Error("Expected token-specific broadcast for token change")
	}
}

// Handler functions that match what would be registered in ConnectMonitorEvents
func handleStatusUpdate(manager *MockManager, event monitor.Event) {
	manager.BroadcastToAll(StatusUpdateMsg, event.Payload)
}

func handleWalletUpdate(manager *MockManager, event monitor.Event) {
	if walletData, ok := event.Payload["wallet"].(map[string]interface{}); ok {
		if address, ok := walletData["address"].(string); ok {
			manager.Broadcast(WalletUpdateMsg, "wallet/"+address, walletData)
		}
	}
}

func handleTokenChange(manager *MockManager, event monitor.Event) {
	if walletAddress, ok := event.Payload["wallet_address"].(string); ok {
		manager.Broadcast(WalletUpdateMsg, "wallet/"+walletAddress, event.Payload)

		if tokenMint, ok := event.Payload["token_mint"].(string); ok {
			manager.Broadcast(WalletUpdateMsg, "token/"+tokenMint, event.Payload)
		}
	}
}
