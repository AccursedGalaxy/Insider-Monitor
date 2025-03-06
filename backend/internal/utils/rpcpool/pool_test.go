package rpcpool

import (
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

// MockRPCClient returns a mock client for testing
func mockClient(endpoint string) *rpc.Client {
	// In a real situation, this would be rpc.New(endpoint), but for testing
	// we just return a non-nil pointer since we won't actually use the client
	return &rpc.Client{}
}

func TestNewClientPool(t *testing.T) {
	// Test with nil options
	pool := NewClientPool("https://test.example.com", nil)
	if pool == nil {
		t.Fatal("Expected non-nil pool with nil options")
	}

	// Test with custom options
	pool = NewClientPool("https://test.example.com", &ClientPoolOptions{
		MinSize:     2,
		MaxSize:     10,
		IdleTimeout: 5 * time.Minute,
		CreateClient: func(endpoint string) *rpc.Client {
			return mockClient(endpoint)
		},
	})

	if pool.minSize != 2 {
		t.Errorf("Expected minSize to be 2, got %d", pool.minSize)
	}

	if pool.maxSize != 10 {
		t.Errorf("Expected maxSize to be 10, got %d", pool.maxSize)
	}

	if pool.idleTimeout != 5*time.Minute {
		t.Errorf("Expected idleTimeout to be 5m, got %v", pool.idleTimeout)
	}
}

func TestClientPoolGetPut(t *testing.T) {
	// Create a pool with low limits for testing
	pool := NewClientPool("https://test.example.com", &ClientPoolOptions{
		MinSize:     1,
		MaxSize:     2,
		IdleTimeout: 1 * time.Minute,
		CreateClient: func(endpoint string) *rpc.Client {
			return mockClient(endpoint)
		},
	})

	// Get a client
	client1, err := pool.Get()
	if err != nil {
		t.Fatalf("Failed to get client1: %v", err)
	}
	if client1 == nil {
		t.Fatal("Expected non-nil client1")
	}

	// Check pool stats
	if pool.InUse() != 1 {
		t.Errorf("Expected 1 client in use, got %d", pool.InUse())
	}

	// Get a second client
	client2, err := pool.Get()
	if err != nil {
		t.Fatalf("Failed to get client2: %v", err)
	}
	if client2 == nil {
		t.Fatal("Expected non-nil client2")
	}

	// Check pool is now exhausted
	if pool.InUse() != 2 {
		t.Errorf("Expected 2 clients in use, got %d", pool.InUse())
	}

	// Try to get a third client which should fail
	_, err = pool.Get()
	if err != ErrPoolExhausted {
		t.Errorf("Expected pool exhausted error, got %v", err)
	}

	// Return a client
	pool.Put(client1)

	// Check stats
	if pool.InUse() != 1 {
		t.Errorf("Expected 1 client in use after Put, got %d", pool.InUse())
	}

	// Should be able to get another client now
	client3, err := pool.Get()
	if err != nil {
		t.Fatalf("Failed to get client after Put: %v", err)
	}
	if client3 == nil {
		t.Fatal("Expected non-nil client after Put")
	}
}

func TestClientPoolClose(t *testing.T) {
	pool := NewClientPool("https://test.example.com", &ClientPoolOptions{
		MinSize:      1,
		MaxSize:      5,
		CreateClient: mockClient,
	})

	// Get a few clients
	client1, _ := pool.Get()
	client2, _ := pool.Get()

	// Close the pool
	pool.Close()

	// Put back clients
	pool.Put(client1)
	pool.Put(client2)

	// Try to get a client from closed pool
	_, err := pool.Get()
	if err != ErrPoolClosed {
		t.Errorf("Expected pool closed error, got %v", err)
	}
}

func TestClientPoolConcurrentAccess(t *testing.T) {
	// Create a pool with reasonable limits
	pool := NewClientPool("https://test.example.com", &ClientPoolOptions{
		MinSize:      2,
		MaxSize:      10,
		CreateClient: mockClient,
	})

	// Run concurrent access test
	var wg sync.WaitGroup
	clientChan := make(chan *rpc.Client, 10)

	// Launch goroutines to get clients
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := pool.Get()
			if err != nil {
				t.Errorf("Failed to get client: %v", err)
				return
			}
			time.Sleep(10 * time.Millisecond) // Simulate work
			clientChan <- client
		}()
	}

	// Wait for all get operations to complete
	wg.Wait()
	close(clientChan)

	// Return all clients
	var returnWg sync.WaitGroup
	returnWg.Add(1)
	go func() {
		defer returnWg.Done()
		for client := range clientChan {
			pool.Put(client)
		}
	}()

	// Wait for all clients to be returned
	returnWg.Wait()

	// Wait a bit to ensure all goroutines complete
	time.Sleep(20 * time.Millisecond)

	// Check pool stats after all operations
	if inUse := pool.InUse(); inUse != 0 {
		t.Errorf("Expected 0 clients in use after test, got %d", inUse)
	}
}

func TestClientPoolLen(t *testing.T) {
	pool := NewClientPool("https://test.example.com", &ClientPoolOptions{
		MinSize:      1,
		MaxSize:      5,
		CreateClient: mockClient,
	})

	// Get and put some clients to increase pool size
	client1, _ := pool.Get()
	client2, _ := pool.Get()
	pool.Put(client1)
	pool.Put(client2)

	// Check total pool size and in-use count
	if pool.Len() != 2 {
		t.Errorf("Expected pool length to be 2, got %d", pool.Len())
	}

	if pool.InUse() != 0 {
		t.Errorf("Expected 0 clients in use, got %d", pool.InUse())
	}

	if pool.Total() != 2 {
		t.Errorf("Expected total pool size to be 2, got %d", pool.Total())
	}
}
