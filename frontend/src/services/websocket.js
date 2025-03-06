/**
 * WebSocket client for real-time communication with the backend
 */

const WS_URL = window.location.origin.replace(/^http/, 'ws') + '/ws';

class WebSocketClient {
  constructor() {
    this.socket = null;
    this.connected = false;
    this.listeners = {};
    this.reconnectTimeout = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000; // Initial delay of 1 second
  }

  /**
   * Initialize WebSocket connection
   */
  connect() {
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      console.log('WebSocket connection already exists');
      return;
    }

    console.log('Connecting to WebSocket server...');
    this.socket = new WebSocket(WS_URL);

    this.socket.onopen = this._onOpen.bind(this);
    this.socket.onclose = this._onClose.bind(this);
    this.socket.onmessage = this._onMessage.bind(this);
    this.socket.onerror = this._onError.bind(this);
  }

  /**
   * Close WebSocket connection
   */
  disconnect() {
    if (!this.socket) return;

    // Remove event listeners
    this.socket.onopen = null;
    this.socket.onclose = null;
    this.socket.onmessage = null;
    this.socket.onerror = null;

    // Clear any pending reconnect
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }

    // Close the connection
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.close();
    }

    this.socket = null;
    this.connected = false;
    console.log('WebSocket disconnected');
  }

  /**
   * Subscribe to a topic
   */
  subscribe(topic) {
    if (!this.connected) {
      console.warn('Cannot subscribe - WebSocket not connected');
      return;
    }

    this.send({
      type: 'subscribe',
      topic: topic
    });

    console.log(`Subscribed to topic: ${topic}`);
  }

  /**
   * Unsubscribe from a topic
   */
  unsubscribe(topic) {
    if (!this.connected) {
      console.warn('Cannot unsubscribe - WebSocket not connected');
      return;
    }

    this.send({
      type: 'unsubscribe',
      topic: topic
    });

    console.log(`Unsubscribed from topic: ${topic}`);
  }

  /**
   * Send a message to the server
   */
  send(data) {
    if (!this.connected) {
      console.warn('Cannot send message - WebSocket not connected');
      return;
    }

    this.socket.send(JSON.stringify(data));
  }

  /**
   * Add an event listener
   */
  addEventListener(type, callback) {
    if (!this.listeners[type]) {
      this.listeners[type] = [];
    }

    this.listeners[type].push(callback);
    return this;
  }

  /**
   * Remove an event listener
   */
  removeEventListener(type, callback) {
    if (!this.listeners[type]) return this;

    this.listeners[type] = this.listeners[type].filter(cb => cb !== callback);
    return this;
  }

  /**
   * Handle WebSocket open event
   */
  _onOpen() {
    this.connected = true;
    this.reconnectAttempts = 0;
    console.log('WebSocket connected');

    // Send a ping to keep the connection alive
    this._startPingInterval();

    // Notify event listeners
    this._dispatchEvent('open', {});
  }

  /**
   * Handle WebSocket close event
   */
  _onClose(event) {
    this.connected = false;
    console.log(`WebSocket closed: ${event.code} ${event.reason}`);

    this._stopPingInterval();

    // Attempt to reconnect
    this._reconnect();

    // Notify event listeners
    this._dispatchEvent('close', event);
  }

  /**
   * Handle WebSocket message event
   */
  _onMessage(event) {
    let data;
    try {
      data = JSON.parse(event.data);
    } catch (e) {
      console.error('Invalid WebSocket message format:', e);
      return;
    }

    // Handle system messages
    if (data.type === 'pong') {
      // Server responded to ping
      return;
    }

    // Dispatch event to listeners
    this._dispatchEvent('message', data);

    // Also dispatch based on message type
    if (data.type) {
      this._dispatchEvent(data.type, data);
    }
  }

  /**
   * Handle WebSocket error event
   */
  _onError(error) {
    console.error('WebSocket error:', error);

    // Notify event listeners
    this._dispatchEvent('error', error);
  }

  /**
   * Dispatch an event to listeners
   */
  _dispatchEvent(type, data) {
    if (!this.listeners[type]) return;

    this.listeners[type].forEach(callback => {
      try {
        callback(data);
      } catch (e) {
        console.error(`Error in WebSocket '${type}' event handler:`, e);
      }
    });
  }

  /**
   * Start ping interval to keep connection alive
   */
  _startPingInterval() {
    this.pingInterval = setInterval(() => {
      if (this.connected) {
        this.send({ type: 'ping' });
      }
    }, 30000); // Ping every 30 seconds
  }

  /**
   * Stop ping interval
   */
  _stopPingInterval() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  /**
   * Attempt to reconnect with exponential backoff
   */
  _reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Maximum reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(30000, this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1));

    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);

    this.reconnectTimeout = setTimeout(() => {
      this.connect();
    }, delay);
  }
}

// Create singleton instance
const wsClient = new WebSocketClient();

export default wsClient;
