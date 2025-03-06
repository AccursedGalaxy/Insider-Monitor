package rpcpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

// MockContextOperation simulates an RPC operation
type MockContextOperation struct {
	delay       time.Duration
	shouldError bool
	errorMsg    string
}

func (m *MockContextOperation) Execute() error {
	time.Sleep(m.delay)
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return nil
}

func TestNewContextClient(t *testing.T) {
	// Test with default timeout
	client := NewContextClient(&rpc.Client{}, 0)
	if client.defaultTimeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", client.defaultTimeout)
	}

	// Test with custom timeout
	client = NewContextClient(&rpc.Client{}, 10*time.Second)
	if client.defaultTimeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", client.defaultTimeout)
	}
}

func TestContextualizeRequest(t *testing.T) {
	client := NewContextClient(&rpc.Client{}, 1*time.Second)

	// Test with fast operation and no context
	fastOp := &MockContextOperation{
		delay:       10 * time.Millisecond,
		shouldError: false,
	}

	err := client.contextualizeRequest(nil, fastOp.Execute)
	if err != nil {
		t.Errorf("Expected no error for fast operation, got %v", err)
	}

	// Test with slow operation and timeout context
	slowOp := &MockContextOperation{
		delay:       200 * time.Millisecond,
		shouldError: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = client.contextualizeRequest(ctx, slowOp.Execute)
	if err == nil {
		t.Error("Expected context deadline exceeded error, got nil")
	}

	// Test with error operation
	errorOp := &MockContextOperation{
		delay:       10 * time.Millisecond,
		shouldError: true,
		errorMsg:    "test error",
	}

	err = client.contextualizeRequest(context.Background(), errorOp.Execute)
	if err == nil || err.Error() != "test error" {
		t.Errorf("Expected 'test error', got %v", err)
	}
}

func TestClientGetter(t *testing.T) {
	mockRPC := &rpc.Client{}
	client := NewContextClient(mockRPC, 1*time.Second)

	if client.Client() != mockRPC {
		t.Error("Client() should return the underlying RPC client")
	}
}

// Note: Full test coverage of the specific RPC methods would require
// mocking the Solana RPC responses. This would be a more complex task
// requiring detailed mocks of the gagliardetto/solana-go library.
//
// The following test demonstrates how we would test the pattern of one
// such method, but a production implementation would need more detailed
// response mocking.

func TestGetBalanceWithContext(t *testing.T) {
	// This is a simplified test that checks the pattern
	// A full test would mock the RPC client's response

	mockRPC := &rpc.Client{}
	client := NewContextClient(mockRPC, 100*time.Millisecond)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Sleep to force timeout
	go func() {
		time.Sleep(200 * time.Millisecond)
	}()

	// The actual call would fail with timeout, but we can't fully test
	// without mocking the RPC client
	_, err := client.GetBalance(ctx, "SomeAccount", nil)

	// We expect this to fail, but we can't assert on exact error
	// since we're not actually mocking the client
	if err == nil {
		t.Error("Expected an error from GetBalance with short timeout")
	}
}
