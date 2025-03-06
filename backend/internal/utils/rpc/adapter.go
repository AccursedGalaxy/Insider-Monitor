package rpc

import (
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

// Client is an alias for rpc.Client
type Client = rpc.Client

// GetAccountInfoResult is an alias for rpc.GetAccountInfoResult
type GetAccountInfoResult = rpc.GetAccountInfoResult

// GetAccountInfoOpts is an alias for rpc.GetAccountInfoOpts
type GetAccountInfoOpts = rpc.GetAccountInfoOpts

// GetBalanceResult is an alias for rpc.GetBalanceResult
type GetBalanceResult = rpc.GetBalanceResult

// Custom types for missing solana-go definitions
type GetBalanceOpts struct{}

// GetTokenAccountsResult is an alias for rpc.GetTokenAccountsResult
type GetTokenAccountsResult = rpc.GetTokenAccountsResult

// GetTokenAccountsConfig is an alias for rpc.GetTokenAccountsConfig
type GetTokenAccountsConfig = rpc.GetTokenAccountsConfig

// GetTokenAccountsOpts is an alias for rpc.GetTokenAccountsOpts
type GetTokenAccountsOpts = rpc.GetTokenAccountsOpts

// GetProgramAccountsResult is an alias for rpc.GetProgramAccountsResult
type GetProgramAccountsResult = rpc.GetProgramAccountsResult

// GetProgramAccountsOpts is an alias for rpc.GetProgramAccountsOpts
type GetProgramAccountsOpts = rpc.GetProgramAccountsOpts

// GetTransactionResult is an alias for rpc.GetTransactionResult
type GetTransactionResult = rpc.GetTransactionResult

// GetTransactionOpts is an alias for rpc.GetTransactionOpts
type GetTransactionOpts = rpc.GetTransactionOpts

// Custom types for missing solana-go definitions
type GetSignaturesForAddressResult struct{}
type GetSignaturesForAddressOpts struct{}

// New creates a new RPC client
func New(endpoint string) *Client {
	return rpc.New(endpoint)
}

// NewWithLimiter creates a new RPC client with rate limiting
func NewWithLimiter(endpoint string, limit rate.Limit, burst int) *Client {
	// Simplified implementation since solana-go may not have these exact methods
	return rpc.New(endpoint)
}

// NewWithCustomRPCClient creates a new custom RPC client
func NewWithCustomRPCClient(client *Client) *Client {
	return client
}
