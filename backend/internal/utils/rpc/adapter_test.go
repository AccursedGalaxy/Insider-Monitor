package rpc

import (
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

func TestNew(t *testing.T) {
	// Test creating a new RPC client
	client := New("https://api.mainnet-beta.solana.com")

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	// Verify it's the correct type
	_, ok := interface{}(client).(*rpc.Client)
	if !ok {
		t.Error("Expected client to be of type *rpc.Client")
	}
}

func TestNewWithLimiter(t *testing.T) {
	// Test creating a new RPC client with a rate limiter
	client := NewWithLimiter("https://api.mainnet-beta.solana.com", rate.Limit(10), 1)

	if client == nil {
		t.Fatal("Expected non-nil client with rate limiter")
	}

	// Verify it's the correct type
	_, ok := interface{}(client).(*rpc.Client)
	if !ok {
		t.Error("Expected client to be of type *rpc.Client")
	}
}

func TestNewWithCustomRPCClient(t *testing.T) {
	// Create a mock client
	mockClient := &rpc.Client{}

	// Test using the custom client
	client := NewWithCustomRPCClient(mockClient)

	if client == nil {
		t.Fatal("Expected non-nil custom client")
	}

	// Verify it's the same client instance
	if client != mockClient {
		t.Error("Expected custom client to be the same instance as the input")
	}
}

func TestTypeAliases(t *testing.T) {
	// Test that our type aliases match the expected types from solana-go

	// GetBalanceOpts is our custom type, so we can directly instantiate it
	balanceOpts := GetBalanceOpts{}

	// Verify we can create instances of the aliased types
	var (
		_ GetAccountInfoResult          = rpc.GetAccountInfoResult{}
		_ GetAccountInfoOpts            = rpc.GetAccountInfoOpts{}
		_ GetBalanceResult              = rpc.GetBalanceResult{}
		_ GetTokenAccountsResult        = rpc.GetTokenAccountsResult{}
		_ GetTokenAccountsConfig        = rpc.GetTokenAccountsConfig{}
		_ GetTokenAccountsOpts          = rpc.GetTokenAccountsOpts{}
		_ GetProgramAccountsResult      = rpc.GetProgramAccountsResult{}
		_ GetProgramAccountsOpts        = rpc.GetProgramAccountsOpts{}
		_ GetTransactionResult          = rpc.GetTransactionResult{}
		_ GetTransactionOpts            = rpc.GetTransactionOpts{}
		_ GetSignaturesForAddressResult = GetSignaturesForAddressResult{}
		_ GetSignaturesForAddressOpts   = GetSignaturesForAddressOpts{}
	)

	// This is primarily a compile-time check, so as long as we can instantiate
	// the variables above, the test passes. We add this assertion to avoid
	// unused variable warnings.
	if balanceOpts != (GetBalanceOpts{}) {
		t.Error("Unexpected value for GetBalanceOpts")
	}
}
