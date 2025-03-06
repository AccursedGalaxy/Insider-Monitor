package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type WalletMonitor struct {
	rpcPool      *RPCPool
	wallets      []solana.PublicKey
	networkURL   string
	isConnected  bool
	mu           sync.RWMutex
	ticker       *time.Ticker
	scanInterval time.Duration
	walletData   map[string]*WalletData
	config       struct {
		NetworkURL string
	}
	scanConfigs map[string]ScanConfigInfo // Maps wallet address to scan config
}

func NewWalletMonitor(networkURL string, wallets []string, options interface{}) (*WalletMonitor, error) {
	// Create RPC pool
	clientPool := NewRPCPool(networkURL, 10)

	// Convert wallet addresses to PublicKeys
	pubKeys := make([]solana.PublicKey, len(wallets))
	for i, addr := range wallets {
		pubKey, err := solana.PublicKeyFromBase58(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid wallet address %s: %v", addr, err)
		}
		pubKeys[i] = pubKey
	}

	return &WalletMonitor{
		rpcPool:     clientPool,
		wallets:     pubKeys,
		networkURL:  networkURL,
		scanConfigs: make(map[string]ScanConfigInfo),
	}, nil
}

func (w *WalletMonitor) getTokenAccountsWithRetry(ctx context.Context, wallet solana.PublicKey) (*rpc.GetTokenAccountsResult, error) {
	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Get a client from the pool
		client := w.rpcPool.Get()

		// Make sure to return the client to the pool when done
		defer w.rpcPool.Put(client)

		// Use a timeout for the individual request
		requestCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		accounts, err := client.GetTokenAccountsByOwner(
			requestCtx,
			wallet,
			&rpc.GetTokenAccountsConfig{
				ProgramId: solana.TokenProgramID.ToPointer(),
			},
			&rpc.GetTokenAccountsOpts{
				Encoding: solana.EncodingBase64,
			},
		)

		if err == nil {
			return accounts, nil
		}

		lastErr = err
		if strings.Contains(err.Error(), "429") {
			log.Printf("Rate limited on attempt %d for wallet %s, waiting %v before retry",
				attempt+1, wallet.String(), backoff)
			time.Sleep(backoff)

			// Exponential backoff with max
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// If it's not a rate limit error, return immediately
		return nil, err
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
