package rpcpool

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	// Import our local RPC adapter that defines custom types
	localrpc "github.com/accursedgalaxy/insider-monitor/backend/internal/utils/rpc"
)

// ContextClient wraps an RPC client with context support
type ContextClient struct {
	client         *rpc.Client
	defaultTimeout time.Duration
}

// NewContextClient creates a new context-aware client
func NewContextClient(client *rpc.Client, defaultTimeout time.Duration) *ContextClient {
	if defaultTimeout == 0 {
		defaultTimeout = 30 * time.Second
	}

	return &ContextClient{
		client:         client,
		defaultTimeout: defaultTimeout,
	}
}

// ErrContextDeadlineExceeded is returned when a context deadline is exceeded
var ErrContextDeadlineExceeded = errors.New("context deadline exceeded")

// Client returns the underlying RPC client
func (c *ContextClient) Client() *rpc.Client {
	return c.client
}

// contextualizeRequest adds context awareness to an RPC request
func (c *ContextClient) contextualizeRequest(ctx context.Context, fn func() error) error {
	// If no context is provided, use a default timeout
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.defaultTimeout)
		defer cancel()
	}

	// Create a channel to receive the result of the RPC call
	done := make(chan error, 1)

	// Execute the RPC call in a goroutine
	go func() {
		done <- fn()
	}()

	// Wait for either the RPC call to complete or the context to be canceled
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ErrContextDeadlineExceeded
	}
}

// GetAccountInfo gets account information with context
func (c *ContextClient) GetAccountInfo(
	ctx context.Context,
	account string,
	opts *rpc.GetAccountInfoOpts,
) (*rpc.GetAccountInfoResult, error) {
	pubkey, err := solana.PublicKeyFromBase58(account)
	if err != nil {
		return nil, err
	}

	var result *rpc.GetAccountInfoResult
	err = c.contextualizeRequest(ctx, func() error {
		var rpcErr error
		result, rpcErr = c.client.GetAccountInfo(ctx, pubkey)
		return rpcErr
	})
	return result, err
}

// GetBalance gets account balance with context
func (c *ContextClient) GetBalance(
	ctx context.Context,
	account string,
	opts *localrpc.GetBalanceOpts,
) (*rpc.GetBalanceResult, error) {
	pubkey, err := solana.PublicKeyFromBase58(account)
	if err != nil {
		return nil, err
	}

	var result *rpc.GetBalanceResult
	err = c.contextualizeRequest(ctx, func() error {
		var rpcErr error
		// The underlying API might expect a commitment type parameter
		var commitment rpc.CommitmentType
		result, rpcErr = c.client.GetBalance(ctx, pubkey, commitment)
		return rpcErr
	})
	return result, err
}

// GetTokenAccountsByOwner gets token accounts by owner with context
func (c *ContextClient) GetTokenAccountsByOwner(
	ctx context.Context,
	wallet string,
	opts *rpc.GetTokenAccountsConfig,
	args *rpc.GetTokenAccountsOpts,
) (*rpc.GetTokenAccountsResult, error) {
	pubkey, err := solana.PublicKeyFromBase58(wallet)
	if err != nil {
		return nil, err
	}

	var result *rpc.GetTokenAccountsResult
	err = c.contextualizeRequest(ctx, func() error {
		var rpcErr error
		result, rpcErr = c.client.GetTokenAccountsByOwner(ctx, pubkey, opts, args)
		return rpcErr
	})
	return result, err
}

// GetProgramAccounts gets program accounts with context
func (c *ContextClient) GetProgramAccounts(
	ctx context.Context,
	programID string,
	opts *rpc.GetProgramAccountsOpts,
) (*rpc.GetProgramAccountsResult, error) {
	pubkey, err := solana.PublicKeyFromBase58(programID)
	if err != nil {
		return nil, err
	}

	var result rpc.GetProgramAccountsResult
	var rpcErr error
	err = c.contextualizeRequest(ctx, func() error {
		// Fix: The client method returns two values - result and error
		result, rpcErr = c.client.GetProgramAccounts(ctx, pubkey)
		return rpcErr
	})
	return &result, err
}

// GetTransaction gets transaction details with context
func (c *ContextClient) GetTransaction(
	ctx context.Context,
	signature string,
	opts *rpc.GetTransactionOpts,
) (*rpc.GetTransactionResult, error) {
	sig, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return nil, err
	}

	var result *rpc.GetTransactionResult
	err = c.contextualizeRequest(ctx, func() error {
		var rpcErr error
		result, rpcErr = c.client.GetTransaction(ctx, sig, opts)
		return rpcErr
	})
	return result, err
}

// GetSignaturesForAddress gets signatures for an address with context
func (c *ContextClient) GetSignaturesForAddress(
	ctx context.Context,
	address string,
	opts *localrpc.GetSignaturesForAddressOpts,
) (*localrpc.GetSignaturesForAddressResult, error) {
	// Validate the address but don't use it since this is a placeholder
	_, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, err
	}

	// This is a placeholder - we're using our custom types from the adapter
	// You would need to adapt this to the actual signatures API
	var result *localrpc.GetSignaturesForAddressResult
	err = c.contextualizeRequest(ctx, func() error {
		// In a real implementation, we would call:
		// result, rpcErr = c.client.GetSignaturesForAddress(ctx, pubkey, ...)
		// But this is just a placeholder for the custom type
		result = &localrpc.GetSignaturesForAddressResult{}
		return nil
	})
	return result, err
}

// WithHTTPClient sets a custom HTTP client with context support
func WithHTTPClient(client *rpc.Client, httpClient *http.Client) *rpc.Client {
	// Note: This is a simplified approach. In a real implementation,
	// you might need to modify the internal HTTP client of the RPC client,
	// which might not be directly accessible depending on the library's design.

	// This is a placeholder that simulates setting a custom HTTP client.
	// The actual implementation would depend on the RPC library's structure.

	return client
}
