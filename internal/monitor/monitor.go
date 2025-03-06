package monitor

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	bin "github.com/gagliardetto/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

type WalletMonitor struct {
	client       *rpc.Client
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
}

func NewWalletMonitor(networkURL string, wallets []string, options interface{}) (*WalletMonitor, error) {
	client := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		networkURL,
		rate.Every(time.Second/4),
		1,
	))

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
		client:     client,
		wallets:    pubKeys,
		networkURL: networkURL,
	}, nil
}

// Simplified TokenAccountInfo
type TokenAccountInfo struct {
	Balance     uint64    `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
	Symbol      string    `json:"symbol"`
	Decimals    uint8     `json:"decimals"`
}

// Simplified WalletData
type WalletData struct {
	WalletAddress string                      `json:"wallet_address"`
	TokenAccounts map[string]TokenAccountInfo `json:"token_accounts"` // mint -> info
	LastScanned   time.Time                   `json:"last_scanned"`
}

// Add these constants for retry configuration
const (
	maxRetries     = 5
	initialBackoff = 5 * time.Second
	maxBackoff     = 30 * time.Second
)

func (w *WalletMonitor) getTokenAccountsWithRetry(wallet solana.PublicKey) (*rpc.GetTokenAccountsResult, error) {
	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		accounts, err := w.client.GetTokenAccountsByOwner(
			context.Background(),
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

func (w *WalletMonitor) GetWalletData(wallet solana.PublicKey) (*WalletData, error) {
	walletData := &WalletData{
		WalletAddress: wallet.String(),
		TokenAccounts: make(map[string]TokenAccountInfo),
		LastScanned:   time.Now(),
	}

	// Use the retry version instead
	accounts, err := w.getTokenAccountsWithRetry(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	// Process token accounts
	for _, acc := range accounts.Value {
		var tokenAccount token.Account
		err = bin.NewBinDecoder(acc.Account.Data.GetBinary()).Decode(&tokenAccount)
		if err != nil {
			log.Printf("warning: failed to decode token account: %v", err)
			continue
		}

		// Only include accounts with positive balance
		if tokenAccount.Amount > 0 {
			mint := tokenAccount.Mint.String()
			walletData.TokenAccounts[mint] = TokenAccountInfo{
				Balance:     tokenAccount.Amount,
				LastUpdated: time.Now(),
				Symbol:      mint[:8] + "...",
				Decimals:    9,
			}
		}
	}

	log.Printf("Wallet %s: found %d token accounts", wallet.String(), len(walletData.TokenAccounts))
	return walletData, nil
}

// Add these type definitions
type Change struct {
	WalletAddress string
	TokenMint     string
	TokenSymbol   string // Add symbol
	TokenDecimals uint8  // Add decimals
	ChangeType    string
	OldBalance    uint64
	NewBalance    uint64
	ChangePercent float64
	TokenBalances map[string]uint64 `json:",omitempty"`
}

func calculatePercentageChange(old, new uint64) float64 {
	if old == 0 {
		return 100.0 // Return 100% for new additions
	}

	// Convert to float64 before division to maintain precision
	oldFloat := float64(old)
	newFloat := float64(new)

	// Calculate percentage change
	change := ((newFloat - oldFloat) / oldFloat) * 100.0

	// Round to 2 decimal places to avoid floating point precision issues
	change = float64(int64(change*100)) / 100

	return change
}

// Utility function for absolute values
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// formatTokenAmount formats a token amount based on its decimals and adds appropriate suffix
func formatTokenAmount(amount uint64, decimals uint8) string {
	if decimals == 0 {
		return fmt.Sprintf("%d", amount)
	}

	// Convert to float with proper decimal places
	divisor := math.Pow10(int(decimals))
	value := float64(amount) / divisor

	// Format based on size
	if value >= 1000000 {
		return fmt.Sprintf("%.2fM", value/1000000)
	} else if value >= 1000 {
		return fmt.Sprintf("%.2fK", value/1000)
	}

	// Regular formatting for smaller numbers (show 4 decimal places)
	return fmt.Sprintf("%.4f", value)
}

func DetectChanges(oldData, newData map[string]*WalletData, significantChange float64) []Change {
	var changes []Change

	// Check for changes in existing wallets
	for walletAddr, newWalletData := range newData {
		oldWalletData, existed := oldData[walletAddr]

		if !existed {
			continue // Skip new wallet detection for now
		}

		// Check for changes in existing wallet
		for mint, newInfo := range newWalletData.TokenAccounts {
			oldInfo, existed := oldWalletData.TokenAccounts[mint]

			if !existed {
				// New token detected
				changes = append(changes, Change{
					WalletAddress: walletAddr,
					TokenMint:     mint,
					TokenSymbol:   newInfo.Symbol,
					TokenDecimals: newInfo.Decimals,
					ChangeType:    "new_token",
					NewBalance:    newInfo.Balance,
				})
				continue
			}

			// Check for significant balance changes
			pctChange := calculatePercentageChange(oldInfo.Balance, newInfo.Balance)
			absChange := abs(pctChange)

			if absChange >= significantChange {
				changes = append(changes, Change{
					WalletAddress: walletAddr,
					TokenMint:     mint,
					TokenSymbol:   newInfo.Symbol,
					TokenDecimals: newInfo.Decimals,
					ChangeType:    "balance_change",
					OldBalance:    oldInfo.Balance,
					NewBalance:    newInfo.Balance,
					ChangePercent: pctChange,
				})
			}
		}
	}

	return changes
}

// UpdateScanInterval updates the scan interval
func (m *WalletMonitor) UpdateScanInterval(interval time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update the interval
	m.scanInterval = interval

	// If you have a ticker, reset it
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = time.NewTicker(interval)
	}
}

// AddWallet adds a wallet to monitor
func (m *WalletMonitor) AddWallet(address string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Convert address to public key
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return fmt.Errorf("invalid wallet address: %w", err)
	}

	// Check if wallet already exists
	for _, wallet := range m.wallets {
		if wallet.String() == address {
			return nil // Already exists, no error
		}
	}

	// Add to wallets list
	m.wallets = append(m.wallets, pubKey)

	// Initialize wallet data if needed
	if m.walletData == nil {
		m.walletData = make(map[string]*WalletData)
	}

	// Add wallet data entry
	m.walletData[address] = &WalletData{
		WalletAddress: address,
		TokenAccounts: make(map[string]TokenAccountInfo),
		LastScanned:   time.Time{},
	}

	return nil
}

// RemoveWallet removes a wallet from monitoring
func (m *WalletMonitor) RemoveWallet(address string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create a new slice without the removed wallet
	var newWallets []solana.PublicKey
	for _, wallet := range m.wallets {
		if wallet.String() != address {
			newWallets = append(newWallets, wallet)
		}
	}
	m.wallets = newWallets

	// Remove from wallet data
	if m.walletData != nil {
		delete(m.walletData, address)
	}
}

func (w *WalletMonitor) checkConnection() error {
	// Try to get slot number as a simple connection test
	_, err := w.client.GetSlot(context.Background(), rpc.CommitmentFinalized)
	w.isConnected = err == nil
	return err
}

// ScanAllWallets scans all wallets and returns their data
func (w *WalletMonitor) ScanAllWallets() (map[string]*WalletData, error) {
	// Check connection first
	if err := w.checkConnection(); err != nil {
		return nil, fmt.Errorf("connection check failed: %w", err)
	}

	results := make(map[string]*WalletData)
	batchSize := 2

	for i := 0; i < len(w.wallets); i += batchSize {
		end := i + batchSize
		if end > len(w.wallets) {
			end = len(w.wallets)
		}

		log.Printf("Processing wallets %d-%d of %d", i+1, end, len(w.wallets))

		// Process batch
		for _, wallet := range w.wallets[i:end] {
			data, err := w.GetWalletData(wallet)
			if err != nil {
				log.Printf("error scanning wallet %s: %v", wallet.String(), err)
				continue
			}
			results[wallet.String()] = data
		}

		// Larger wait between batches
		if end < len(w.wallets) {
			waitTime := 3 * time.Second
			log.Printf("Waiting %v before next batch...", waitTime)
			time.Sleep(waitTime)
		}
	}

	return results, nil
}
