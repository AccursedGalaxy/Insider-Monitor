package rpcpool

import (
	"errors"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

var (
	// ErrPoolExhausted indicates all connections are in use
	ErrPoolExhausted = errors.New("connection pool exhausted")

	// ErrPoolClosed indicates the pool has been closed
	ErrPoolClosed = errors.New("connection pool is closed")
)

// ClientPool manages a pool of RPC clients
type ClientPool struct {
	mu           sync.Mutex
	clients      []*rpc.Client
	inUse        map[*rpc.Client]struct{}
	maxSize      int
	minSize      int
	endpoint     string
	idleTimeout  time.Duration
	closed       bool
	createClient func(string) *rpc.Client
}

// ClientPoolOptions configures a client pool
type ClientPoolOptions struct {
	// MinSize is the minimum number of connections to keep open
	MinSize int

	// MaxSize is the maximum number of connections allowed
	MaxSize int

	// IdleTimeout is how long connections can remain idle before being closed
	IdleTimeout time.Duration

	// Custom client creation function
	CreateClient func(string) *rpc.Client
}

// NewClientPool creates a new RPC client pool
func NewClientPool(endpoint string, options *ClientPoolOptions) *ClientPool {
	if options == nil {
		options = &ClientPoolOptions{
			MinSize:     2,
			MaxSize:     10,
			IdleTimeout: 30 * time.Second,
		}
	}

	if options.MaxSize < options.MinSize {
		options.MaxSize = options.MinSize
	}

	createClient := options.CreateClient
	if createClient == nil {
		createClient = func(endpoint string) *rpc.Client {
			return rpc.New(endpoint)
		}
	}

	pool := &ClientPool{
		clients:      make([]*rpc.Client, 0, options.MaxSize),
		inUse:        make(map[*rpc.Client]struct{}),
		maxSize:      options.MaxSize,
		minSize:      options.MinSize,
		endpoint:     endpoint,
		idleTimeout:  options.IdleTimeout,
		createClient: createClient,
	}

	// Pre-create minimum connections
	for i := 0; i < options.MinSize; i++ {
		client := createClient(endpoint)
		pool.clients = append(pool.clients, client)
	}

	// Start the idle connection cleanup
	if options.IdleTimeout > 0 {
		go pool.cleanupIdleConnections()
	}

	return pool
}

// Get returns a client from the pool
func (p *ClientPool) Get() (*rpc.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, ErrPoolClosed
	}

	// If there's an available client, return it
	if len(p.clients) > 0 {
		client := p.clients[len(p.clients)-1]
		p.clients = p.clients[:len(p.clients)-1]
		p.inUse[client] = struct{}{}
		return client, nil
	}

	// If we're at capacity, return error
	if len(p.inUse) >= p.maxSize {
		return nil, ErrPoolExhausted
	}

	// Create a new client
	client := p.createClient(p.endpoint)
	p.inUse[client] = struct{}{}

	return client, nil
}

// Put returns a client to the pool
func (p *ClientPool) Put(client *rpc.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		client.Close()
		return
	}

	// Check if the client is from this pool
	if _, ok := p.inUse[client]; !ok {
		// Not our client, just close it
		client.Close()
		return
	}

	// Remove from in-use set
	delete(p.inUse, client)

	// If we have too many idle clients, close this one
	if len(p.clients) >= p.maxSize {
		client.Close()
		return
	}

	// Otherwise, return it to the pool
	p.clients = append(p.clients, client)
}

// Close closes the pool and all clients
func (p *ClientPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	// Close all available clients
	for _, client := range p.clients {
		client.Close()
	}

	// Clear the pool
	p.clients = nil

	// Note: in-use clients will be closed when Put is called
}

// Len returns the number of clients in the pool
func (p *ClientPool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.clients)
}

// InUse returns the number of clients currently in use
func (p *ClientPool) InUse() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.inUse)
}

// Total returns the total number of clients managed by the pool
func (p *ClientPool) Total() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.clients) + len(p.inUse)
}

// cleanupIdleConnections periodically checks for idle connections to close
func (p *ClientPool) cleanupIdleConnections() {
	ticker := time.NewTicker(p.idleTimeout)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()

		// Stop if pool is closed
		if p.closed {
			p.mu.Unlock()
			return
		}

		// Keep minimum connections, close excess
		excessCount := len(p.clients) - p.minSize
		if excessCount > 0 {
			// Close excess connections
			for i := 0; i < excessCount; i++ {
				if len(p.clients) == 0 {
					break
				}

				client := p.clients[len(p.clients)-1]
				p.clients = p.clients[:len(p.clients)-1]
				client.Close()
			}
		}

		p.mu.Unlock()
	}
}
