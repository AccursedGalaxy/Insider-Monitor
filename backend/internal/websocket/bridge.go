package websocket

import (
	"log"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
)

// ConnectMonitorEvents connects monitor events to WebSocket broadcasts
func ConnectMonitorEvents(wsManager *Manager) {
	// Get the event emitter
	emitter := monitor.GetEventEmitter()

	// Subscribe to WalletScanStarted events
	emitter.Subscribe(monitor.WalletScanStarted, func(event monitor.Event) {
		log.Printf("Received WalletScanStarted event")
		wsManager.BroadcastToAll(StatusUpdateMsg, event.Payload)
	})

	// Subscribe to WalletScanComplete events
	emitter.Subscribe(monitor.WalletScanComplete, func(event monitor.Event) {
		log.Printf("Received WalletScanComplete event")
		wsManager.BroadcastToAll(StatusUpdateMsg, event.Payload)
	})

	// Subscribe to WalletScanError events
	emitter.Subscribe(monitor.WalletScanError, func(event monitor.Event) {
		log.Printf("Received WalletScanError event")
		wsManager.BroadcastToAll(StatusUpdateMsg, event.Payload)
	})

	// Subscribe to WalletDataUpdated events
	emitter.Subscribe(monitor.WalletDataUpdated, func(event monitor.Event) {
		log.Printf("Received WalletDataUpdated event")
		if walletData, ok := event.Payload["wallet"].(map[string]interface{}); ok {
			if address, ok := walletData["address"].(string); ok {
				// Broadcast to the specific wallet topic
				wsManager.Broadcast(WalletUpdateMsg, "wallet/"+address, walletData)
			}
		}
	})

	// Subscribe to TokenChange events
	emitter.Subscribe(monitor.TokenChange, func(event monitor.Event) {
		log.Printf("Received TokenChange event")
		if walletAddress, ok := event.Payload["wallet_address"].(string); ok {
			// Broadcast to the specific wallet topic
			wsManager.Broadcast(WalletUpdateMsg, "wallet/"+walletAddress, event.Payload)

			// Also broadcast to the token topic if we have a token mint
			if tokenMint, ok := event.Payload["token_mint"].(string); ok {
				wsManager.Broadcast(WalletUpdateMsg, "token/"+tokenMint, event.Payload)
			}
		}
	})

	log.Println("Monitor events connected to WebSocket")
}

// InitializeWebSocketBridge sets up the WebSocket bridge with the monitor
func InitializeWebSocketBridge(wsManager *Manager) {
	// Connect monitor events
	ConnectMonitorEvents(wsManager)

	log.Println("WebSocket bridge initialized")
}
