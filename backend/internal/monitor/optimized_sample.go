package monitor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/utils/async"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/utils/batch"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/utils/breaker"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/utils/objpool"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/utils/rpcpool"
)

// OptimizedWalletScanner demonstrates how to use performance optimizations
type OptimizedWalletScanner struct {
	// RPC client pool
	rpcPool *rpcpool.ClientPool

	// Circuit breaker for external service calls
	circuitBreaker *breaker.CircuitBreaker

	// Object pool for token data
	tokenDataPool *objpool.MessagePool[*TokenData]

	// Batch processor for wallet scanning
	walletBatcher *batch.BatchProcessor[string, *WalletScanResult]

	// Worker pool for processing wallet data
	processingWorker *async.Worker[*ProcessedWalletData]

	// Processing pipeline for alert generation
	alertPipeline *async.Pipeline[*AlertData, *ProcessedAlert]

	// Map pool for temporary storage
	mapPool *objpool.MapPool[string, interface{}]

	// Context for shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// TokenData represents token data (simplified for example)
type TokenData struct {
	Mint    string
	Balance uint64
	Symbol  string
}

// WalletScanResult represents a wallet scan result
type WalletScanResult struct {
	Address string
	Tokens  []*TokenData
	Error   error
}

// ProcessedWalletData represents processed wallet data
type ProcessedWalletData struct {
	Address   string
	TokenData map[string]*TokenData
}

// AlertData represents alert data
type AlertData struct {
	WalletAddress string
	TokenMint     string
	OldBalance    uint64
	NewBalance    uint64
}

// ProcessedAlert represents a processed alert
type ProcessedAlert struct {
	AlertID     string
	AlertLevel  string
	Description string
}

// NewOptimizedWalletScanner creates a new optimized wallet scanner
func NewOptimizedWalletScanner(rpcEndpoint string) *OptimizedWalletScanner {
	// Create a context for shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create RPC client pool
	clientPool := rpcpool.NewClientPool(rpcEndpoint, &rpcpool.ClientPoolOptions{
		MinSize:     5,
		MaxSize:     20,
		IdleTimeout: 30 * time.Second,
	})

	// Create circuit breaker for RPC calls
	cb := breaker.New(&breaker.Options{
		Threshold: 5,
		Timeout:   10 * time.Second,
		OnStateChange: func(from, to breaker.CircuitState) {
			log.Printf("Circuit state changed from %v to %v", from, to)
		},
	})

	// Create token data object pool
	tokenPool := objpool.NewMessagePool[*TokenData](
		func() *TokenData {
			return &TokenData{}
		},
		func(td *TokenData) {
			// Reset token data
			td.Mint = ""
			td.Balance = 0
			td.Symbol = ""
		},
	)

	// Create map pool
	mapPool := objpool.NewMapPool[string, interface{}]()

	// Create batch processor for wallet scanning
	walletBatcher := batch.New[string, *WalletScanResult](
		func(ctx context.Context, walletAddresses []string) ([]*WalletScanResult, []error) {
			results := make([]*WalletScanResult, len(walletAddresses))
			errors := make([]error, len(walletAddresses))

			// Process wallet batch (simplified)
			for i, addr := range walletAddresses {
				results[i] = &WalletScanResult{Address: addr}
			}

			return results, errors
		},
		&batch.Options{
			MaxBatchSize: 10,
			MaxWaitTime:  50 * time.Millisecond,
		},
	)

	// Create worker pool for processing
	worker := async.NewWorker[*ProcessedWalletData](10, 100)

	// Create processing pipeline for alerts
	alertPipeline := async.NewPipeline[*AlertData, *ProcessedAlert](
		func(ctx context.Context, data *AlertData) (*ProcessedAlert, error) {
			// Process alert (simplified)
			return &ProcessedAlert{
				AlertID:     "alert-" + data.WalletAddress + "-" + data.TokenMint,
				AlertLevel:  "INFO",
				Description: "Balance changed",
			}, nil
		},
		&async.Options{
			Workers:      5,
			InputBuffer:  100,
			OutputBuffer: 100,
		},
	)

	return &OptimizedWalletScanner{
		rpcPool:          clientPool,
		circuitBreaker:   cb,
		tokenDataPool:    tokenPool,
		walletBatcher:    walletBatcher,
		processingWorker: worker,
		alertPipeline:    alertPipeline,
		mapPool:          mapPool,
		ctx:              ctx,
		cancel:           cancel,
	}
}

// ScanWallet demonstrates how to use the optimizations
func (s *OptimizedWalletScanner) ScanWallet(ctx context.Context, address string) (*WalletScanResult, error) {
	// Use circuit breaker to protect against external service failures
	var result *WalletScanResult
	err := s.circuitBreaker.Execute(func() error {
		// Submit wallet for batch processing
		promise, err := s.walletBatcher.Process(ctx, address)
		if err != nil {
			return err
		}

		// Wait for result
		var batchErr error
		result, batchErr = promise.Wait()
		return batchErr
	})

	return result, err
}

// ProcessWalletData demonstrates async processing
func (s *OptimizedWalletScanner) ProcessWalletData(walletResult *WalletScanResult) {
	// Submit to worker pool
	s.processingWorker.Submit(func(ctx context.Context) (*ProcessedWalletData, error) {
		// Get a map from the pool
		tokenMap := s.mapPool.Get()
		defer s.mapPool.Put(tokenMap)

		// Process tokens
		for _, token := range walletResult.Tokens {
			tokenMap[token.Mint] = token
		}

		return &ProcessedWalletData{
			Address:   walletResult.Address,
			TokenData: nil, // Would convert tokenMap, but simplified for example
		}, nil
	})
}

// GenerateAlert demonstrates pipeline processing
func (s *OptimizedWalletScanner) GenerateAlert(walletAddr, tokenMint string, oldBalance, newBalance uint64) {
	// Create alert data
	alertData := &AlertData{
		WalletAddress: walletAddr,
		TokenMint:     tokenMint,
		OldBalance:    oldBalance,
		NewBalance:    newBalance,
	}

	// Submit to alert pipeline
	s.alertPipeline.Submit(alertData)
}

// Close shuts down the scanner
func (s *OptimizedWalletScanner) Close() {
	// Cancel context
	s.cancel()

	// Close all components
	s.walletBatcher.Close()
	s.processingWorker.Close()
	s.alertPipeline.Close()
	s.rpcPool.Close()
}

// MonitorWallets demonstrates how to use the optimizations together
func MonitorWallets(ctx context.Context, rpcEndpoint string, walletAddresses []string) {
	// Create the optimized scanner
	scanner := NewOptimizedWalletScanner(rpcEndpoint)
	defer scanner.Close()

	// Create wait group for concurrent scanning
	var wg sync.WaitGroup

	// Process results from worker pool in background
	go func() {
		for {
			select {
			case result := <-scanner.processingWorker.Results():
				log.Printf("Processed wallet %s with %d tokens", result.Address, len(result.TokenData))
			case err := <-scanner.processingWorker.Errors():
				log.Printf("Error processing wallet: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Process alerts from pipeline in background
	go func() {
		for {
			select {
			case alert := <-scanner.alertPipeline.Result():
				log.Printf("Alert generated: %s (%s)", alert.AlertID, alert.AlertLevel)
			case err := <-scanner.alertPipeline.Errors():
				log.Printf("Error generating alert: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Scan wallets concurrently
	for _, addr := range walletAddresses {
		wg.Add(1)

		addr := addr // Capture for goroutine
		go func() {
			defer wg.Done()

			result, err := scanner.ScanWallet(ctx, addr)
			if err != nil {
				log.Printf("Error scanning wallet %s: %v", addr, err)
				return
			}

			scanner.ProcessWalletData(result)
		}()
	}

	// Wait for all scans to complete
	wg.Wait()
	log.Println("All wallets scanned")
}
