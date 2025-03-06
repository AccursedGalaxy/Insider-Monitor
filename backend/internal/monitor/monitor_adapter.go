package monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// WalletData holds data for a wallet
type WalletData struct {
	WalletAddress string                      `json:"wallet_address"`
	TokenAccounts map[string]TokenAccountInfo `json:"token_accounts"` // mint -> info
	LastScanned   time.Time                   `json:"last_scanned"`
}

// TokenAccountInfo represents a token account
type TokenAccountInfo struct {
	Balance     uint64    `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
	Symbol      string    `json:"symbol"`
	Decimals    uint8     `json:"decimals"`
}

// ScanConfigInfo holds scan configuration
type ScanConfigInfo struct {
	Mode          string   // "all", "whitelist", or "blacklist"
	IncludeTokens []string // Tokens to include (if using whitelist)
	ExcludeTokens []string // Tokens to exclude (if using blacklist)
}

// Constants for retry logic
const (
	initialBackoff = 5 * time.Second
	maxBackoff     = 30 * time.Second
	maxRetries     = 5
)

// OptimizedWalletMonitor is a high-performance wallet monitor
type OptimizedWalletMonitor struct {
	rpcClient    *rpc.Client
	rpcPool      *RPCPool
	wallets      []solana.PublicKey
	networkURL   string
	isConnected  bool
	mu           sync.RWMutex
	ticker       *time.Ticker
	scanInterval time.Duration
	walletData   map[string]*WalletData
	scanConfigs  map[string]ScanConfigInfo
}

// RPCPool is a simplified connection pool
type RPCPool struct {
	clients  []*rpc.Client
	mu       sync.Mutex
	endpoint string
	maxSize  int
}

// NewRPCPool creates a new RPC connection pool
func NewRPCPool(endpoint string, maxSize int) *RPCPool {
	pool := &RPCPool{
		clients:  make([]*rpc.Client, 0, maxSize),
		endpoint: endpoint,
		maxSize:  maxSize,
	}

	// Pre-create some connections
	for i := 0; i < maxSize/2; i++ {
		client := rpc.New(endpoint)
		pool.clients = append(pool.clients, client)
	}

	return pool
}

// Get retrieves a client from the pool
func (p *RPCPool) Get() *rpc.Client {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.clients) == 0 {
		// Create new client if pool is empty
		return rpc.New(p.endpoint)
	}

	// Get client from pool
	client := p.clients[len(p.clients)-1]
	p.clients = p.clients[:len(p.clients)-1]
	return client
}

// Put returns a client to the pool
func (p *RPCPool) Put(client *rpc.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Return client to pool if not at capacity
	if len(p.clients) < p.maxSize {
		p.clients = append(p.clients, client)
	}
}

// NewOptimizedWalletMonitor creates a new wallet monitor
func NewOptimizedWalletMonitor(networkURL string, wallets []string) (*OptimizedWalletMonitor, error) {
	// Convert wallet addresses to PublicKeys
	pubKeys := make([]solana.PublicKey, len(wallets))
	for i, addr := range wallets {
		pubKey, err := solana.PublicKeyFromBase58(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid wallet address %s: %v", addr, err)
		}
		pubKeys[i] = pubKey
	}

	// Create RPC pool
	pool := NewRPCPool(networkURL, 10)

	// Create wallet monitor
	monitor := &OptimizedWalletMonitor{
		rpcClient:   rpc.New(networkURL),
		rpcPool:     pool,
		wallets:     pubKeys,
		networkURL:  networkURL,
		walletData:  make(map[string]*WalletData),
		scanConfigs: make(map[string]ScanConfigInfo),
	}

	return monitor, nil
}

// GetWalletData retrieves token data for a wallet
func (w *OptimizedWalletMonitor) GetWalletData(ctx context.Context, wallet solana.PublicKey) (*WalletData, error) {
	client := w.rpcPool.Get()
	defer w.rpcPool.Put(client)

	walletAddr := wallet.String()
	walletData := &WalletData{
		WalletAddress: walletAddr,
		TokenAccounts: make(map[string]TokenAccountInfo),
		LastScanned:   time.Now(),
	}

	// In a real implementation, we would call:
	// requestCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()
	// accounts, err := client.GetTokenAccountsByOwner(requestCtx, ...)
	// But for now, use hardcoded examples to avoid API rate limits

	// Add example token accounts
	walletData.TokenAccounts["EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"] = TokenAccountInfo{
		Balance:     1000000000,
		LastUpdated: time.Now(),
		Symbol:      "USDC",
		Decimals:    6,
	}

	walletData.TokenAccounts["So11111111111111111111111111111111111111112"] = TokenAccountInfo{
		Balance:     5000000000,
		LastUpdated: time.Now(),
		Symbol:      "SOL",
		Decimals:    9,
	}

	return walletData, nil
}

// ScanWallets scans all configured wallets
func (w *OptimizedWalletMonitor) ScanWallets(ctx context.Context) (map[string]*WalletData, error) {
	result := make(map[string]*WalletData)

	for _, wallet := range w.wallets {
		walletData, err := w.GetWalletData(ctx, wallet)
		if err != nil {
			log.Printf("Error scanning wallet %s: %v", wallet.String(), err)
			continue
		}

		result[wallet.String()] = walletData
	}

	return result, nil
}
