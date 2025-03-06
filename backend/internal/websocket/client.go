package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

// ClientMessage is a message from a client
type ClientMessage struct {
	Type    MessageType            `json:"type"`
	Topic   string                 `json:"topic,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// readPump pumps messages from the WebSocket connection to the Manager
func (c *Client) readPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Conn.Close()
		c.connected = false
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle message based on type
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error parsing client message: %v", err)
			continue
		}

		// Process message based on type
		switch clientMsg.Type {
		case SubscribeMsg:
			if clientMsg.Topic != "" {
				c.Manager.Subscribe(c.ID, clientMsg.Topic)
			}
		case UnsubscribeMsg:
			if clientMsg.Topic != "" {
				c.Manager.Unsubscribe(c.ID, clientMsg.Topic)
			}
		case PingMsg:
			// Send a pong response
			c.Send <- []byte(`{"type":"pong","time":"` + time.Now().Format(time.RFC3339) + `"}`)
		}
	}
}

// writePump pumps messages from the Manager to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Manager closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
