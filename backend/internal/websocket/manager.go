package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientID is a unique identifier for connected clients
type ClientID string

// Client represents a connected WebSocket client
type Client struct {
	ID        ClientID
	Conn      *websocket.Conn
	Send      chan []byte
	Manager   *Manager
	mu        sync.Mutex
	connected bool
	topics    map[string]bool
}

// Manager handles WebSocket clients and broadcasts
type Manager struct {
	clients    map[ClientID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	topics     map[string]map[ClientID]bool
	mu         sync.RWMutex
	upgrader   websocket.Upgrader
}

// MessageType identifies the type of message
type MessageType string

const (
	// Message types
	WalletUpdateMsg MessageType = "wallet_update"
	ConfigUpdateMsg MessageType = "config_update"
	AlertMsg        MessageType = "alert"
	StatusUpdateMsg MessageType = "status_update"

	// Client message types
	SubscribeMsg   MessageType = "subscribe"
	UnsubscribeMsg MessageType = "unsubscribe"
	PingMsg        MessageType = "ping"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType            `json:"type"`
	Topic   string                 `json:"topic,omitempty"`
	Payload map[string]interface{} `json:"payload"`
	Time    time.Time              `json:"time"`
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[ClientID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		topics:     make(map[string]map[ClientID]bool),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all connections - in production you should restrict this
				return true
			},
		},
	}
}

// Start begins the WebSocket manager
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.registerClient(client)
		case client := <-m.unregister:
			m.unregisterClient(client)
		case message := <-m.broadcast:
			m.broadcastMessage(message)
		}
	}
}

// registerClient adds a client to the manager
func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.ID] = client
	log.Printf("Client registered: %s", client.ID)
}

// unregisterClient removes a client from the manager
func (m *Manager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client.ID]; ok {
		delete(m.clients, client.ID)
		close(client.Send)

		// Remove client from all topics
		for topic, clients := range m.topics {
			if _, ok := clients[client.ID]; ok {
				delete(m.topics[topic], client.ID)
			}
		}

		log.Printf("Client unregistered: %s", client.ID)
	}
}

// broadcastMessage sends a message to clients subscribed to the topic
func (m *Manager) broadcastMessage(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// If topic is specified, send only to clients subscribed to that topic
	if message.Topic != "" {
		if clients, ok := m.topics[message.Topic]; ok {
			for clientID := range clients {
				if client, ok := m.clients[clientID]; ok {
					select {
					case client.Send <- data:
					default:
						m.unregister <- client
					}
				}
			}
		}
	} else {
		// Otherwise broadcast to all clients
		for _, client := range m.clients {
			select {
			case client.Send <- data:
			default:
				m.unregister <- client
			}
		}
	}
}

// HandleWebSocket handles WebSocket requests from clients
func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientID := ClientID(time.Now().String())
	client := &Client{
		ID:        clientID,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Manager:   m,
		connected: true,
		topics:    make(map[string]bool),
	}

	m.register <- client

	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}

// Subscribe adds a client to a topic
func (m *Manager) Subscribe(clientID ClientID, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.topics[topic]; !ok {
		m.topics[topic] = make(map[ClientID]bool)
	}

	m.topics[topic][clientID] = true

	// Add topic to client's subscribed topics
	if client, ok := m.clients[clientID]; ok {
		client.mu.Lock()
		client.topics[topic] = true
		client.mu.Unlock()
	}

	log.Printf("Client %s subscribed to topic: %s", clientID, topic)
}

// Unsubscribe removes a client from a topic
func (m *Manager) Unsubscribe(clientID ClientID, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clients, ok := m.topics[topic]; ok {
		delete(clients, clientID)
	}

	// Remove topic from client's subscribed topics
	if client, ok := m.clients[clientID]; ok {
		client.mu.Lock()
		delete(client.topics, topic)
		client.mu.Unlock()
	}

	log.Printf("Client %s unsubscribed from topic: %s", clientID, topic)
}

// Broadcast sends a message to all clients or those subscribed to a specific topic
func (m *Manager) Broadcast(messageType MessageType, topic string, payload map[string]interface{}) {
	message := &Message{
		Type:    messageType,
		Topic:   topic,
		Payload: payload,
		Time:    time.Now(),
	}

	m.broadcast <- message
}

// BroadcastToAll sends a message to all connected clients
func (m *Manager) BroadcastToAll(messageType MessageType, payload map[string]interface{}) {
	m.Broadcast(messageType, "", payload)
}
