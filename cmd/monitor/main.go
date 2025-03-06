package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/accursedgalaxy/insider-monitor/internal/alerts"
	"github.com/accursedgalaxy/insider-monitor/internal/config"
	"github.com/accursedgalaxy/insider-monitor/internal/monitor"
	"github.com/accursedgalaxy/insider-monitor/internal/storage"
	"github.com/accursedgalaxy/insider-monitor/internal/web"
	"github.com/joho/godotenv"
)

// WalletScanner interface defines the contract for wallet monitoring
type WalletScanner interface {
	ScanAllWallets() (map[string]*monitor.WalletData, error)
}

func main() {
	// Load .env file if it exists
	godotenv.Load()

	testMode := flag.Bool("test", false, "Run in test mode with accelerated scanning")
	configPath := flag.String("config", "config.json", "Path to configuration file")
	webMode := flag.Bool("web", false, "Run with web UI")
	webPort := flag.Int("port", 8080, "Port for web UI (when --web is used)")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *testMode {
		cfg = config.GetTestConfig()
		log.Println("Running in test mode with 5-second scan interval")
	} else {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Initialize scanner
	var scanner WalletScanner
	if *testMode {
		scanner = monitor.NewMockWalletMonitor()
	} else {
		scanner, err = monitor.NewWalletMonitor(cfg.NetworkURL, cfg.Wallets, nil)
		if err != nil {
			log.Fatalf("failed to create wallet monitor: %v", err)
		}

		// Configure scan settings if using the real monitor
		if walletMonitor, ok := scanner.(*monitor.WalletMonitor); ok {
			// Set global scan config
			walletMonitor.SetGlobalScanConfig(
				cfg.Scan.ScanMode,
				cfg.Scan.IncludeTokens,
				cfg.Scan.ExcludeTokens,
			)

			// Set per-wallet scan configs if defined
			for walletAddr, walletCfg := range cfg.WalletConfigs {
				if walletCfg.Scan != nil {
					walletMonitor.SetScanConfig(
						walletAddr,
						walletCfg.Scan.ScanMode,
						walletCfg.Scan.IncludeTokens,
						walletCfg.Scan.ExcludeTokens,
					)
				}
			}
		}
	}

	// Initialize alerter
	var alerter alerts.Alerter
	if cfg.Discord.Enabled {
		alerter = alerts.NewDiscordAlerter(cfg.Discord.WebhookURL, cfg.Discord.ChannelID)
		log.Println("Discord alerts enabled")
	} else {
		alerter = &alerts.ConsoleAlerter{}
		log.Println("Console alerts enabled")
	}

	// Parse the scan interval
	scanInterval, err := time.ParseDuration(cfg.ScanInterval)
	if err != nil {
		log.Printf("invalid scan interval '%s', using default of 1 minute", cfg.ScanInterval)
		scanInterval = time.Minute
	}

	// Initialize storage
	storage := storage.New("./data")

	// If web mode, start the web server
	if *webMode {
		log.Printf("Starting web UI on port %d", *webPort)

		// Start the web server in a goroutine
		go func() {
			walletMonitor, ok := scanner.(*monitor.WalletMonitor)
			if !ok {
				log.Fatalf("Web mode requires the real wallet monitor, not a mock")
			}

			webServer := web.NewServer(cfg, walletMonitor, storage, *webPort)
			if err := webServer.Start(); err != nil {
				log.Fatalf("Failed to start web server: %v", err)
			}
		}()
	}

	// Start the monitor
	runMonitor(scanner, alerter, cfg, scanInterval, storage)
}

func runMonitor(scanner WalletScanner, alerter alerts.Alerter, cfg *config.Config, scanInterval time.Duration, storage *storage.Storage) {
	// Create buffered channels for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Track connection state
	var lastSuccessfulScan time.Time
	var connectionLost bool

	// Define maximum allowed time between successful scans
	maxTimeBetweenScans := scanInterval * 3

	// Initialize previousData from storage at startup
	var previousData map[string]*monitor.WalletData
	if savedData, err := storage.LoadWalletData(); err == nil {
		previousData = savedData
		log.Println("Loaded previous wallet data from storage")
	} else {
		log.Printf("Warning: Could not load previous data: %v. Will initialize after first scan.", err)
		previousData = make(map[string]*monitor.WalletData)
	}

	// Set up timer for regular scans
	ticker := time.NewTicker(scanInterval)

	// Monitor loop
	go func() {
		// Run initial scan
		log.Println("Running initial wallet scan...")
		if latestData, err := runScan(scanner, alerter, previousData, cfg); err == nil {
			previousData = latestData
			if err := storage.SaveWalletData(latestData); err != nil {
				log.Printf("Warning: Failed to save wallet data: %v", err)
			}
			lastSuccessfulScan = time.Now()
			connectionLost = false
		} else {
			log.Printf("Initial scan failed: %v", err)
			connectionLost = true
		}

		// Regular scan loop
		for {
			select {
			case <-ticker.C:
				if latestData, err := runScan(scanner, alerter, previousData, cfg); err == nil {
					// Successful scan
					previousData = latestData
					if err := storage.SaveWalletData(latestData); err != nil {
						log.Printf("Warning: Failed to save wallet data: %v", err)
					}
					lastSuccessfulScan = time.Now()

					// Handle reconnection
					if connectionLost {
						log.Println("Connection restored. Monitoring resumed.")
						connectionLost = false
					}
				} else {
					// Failed scan
					log.Printf("Scan failed: %v", err)

					// Check if connection has been lost for too long
					if !connectionLost && time.Since(lastSuccessfulScan) > maxTimeBetweenScans {
						log.Println("Connection appears to be lost. Will continue retrying.")
						connectionLost = true
					}
				}
			case <-done:
				return
			}
		}
	}()

	// Wait for interrupt signal
	<-interrupt
	log.Println("Shutting down gracefully...")
	ticker.Stop()
	done <- true
	log.Println("Shutdown complete.")
}

func runScan(scanner WalletScanner, alerter alerts.Alerter, previousData map[string]*monitor.WalletData, cfg *config.Config) (map[string]*monitor.WalletData, error) {
	log.Println("Scanning wallets...")
	latestData, err := scanner.ScanAllWallets()
	if err != nil {
		return nil, fmt.Errorf("wallet scan failed: %w", err)
	}

	// Check for changes and send alerts
	for address, walletData := range latestData {
		previousWallet, exists := previousData[address]
		if !exists {
			// First time seeing this wallet, don't alert
			log.Printf("New wallet detected: %s", address)
			continue
		}

		// Compare token accounts
		for mint, tokenAccount := range walletData.TokenAccounts {
			// Skip if token is in the ignore list
			if contains(cfg.Alerts.IgnoreTokens, mint) {
				continue
			}

			// Skip if balance is below minimum threshold
			if float64(tokenAccount.Balance) < cfg.Alerts.MinimumBalance {
				continue
			}

			previousToken, existed := previousWallet.TokenAccounts[mint]

			if !existed {
				// New token detected
				sendTokenAlert(alerter, address, mint, 0, tokenAccount.Balance, tokenAccount.Symbol, "NEW_TOKEN", cfg)
				continue
			}

			// Skip if no change in balance
			if previousToken.Balance == tokenAccount.Balance {
				continue
			}

			// Calculate percentage change
			var percentChange float64
			if previousToken.Balance > 0 {
				percentChange = float64(tokenAccount.Balance-previousToken.Balance) / float64(previousToken.Balance)
			} else if tokenAccount.Balance > 0 {
				percentChange = 1.0 // 100% increase from zero
			} else {
				percentChange = 0.0
			}

			// If change exceeds the significant threshold, send an alert
			if absFloat(percentChange) >= cfg.Alerts.SignificantChange {
				changeType := "INCREASE"
				if percentChange < 0 {
					changeType = "DECREASE"
				}
				sendTokenAlert(alerter, address, mint, previousToken.Balance, tokenAccount.Balance, tokenAccount.Symbol, changeType, cfg)
			}
		}

		// Check for removed tokens (tokens that existed before but not now)
		for mint, previousToken := range previousWallet.TokenAccounts {
			// Skip if token is in the ignore list
			if contains(cfg.Alerts.IgnoreTokens, mint) {
				continue
			}

			// Skip if balance is below minimum threshold
			if float64(previousToken.Balance) < cfg.Alerts.MinimumBalance {
				continue
			}

			if _, exists := walletData.TokenAccounts[mint]; !exists {
				// Token no longer exists in the wallet
				sendTokenAlert(alerter, address, mint, previousToken.Balance, 0, previousToken.Symbol, "REMOVED", cfg)
			}
		}
	}

	return latestData, nil
}

func sendTokenAlert(alerter alerts.Alerter, walletAddress, tokenMint string, oldBalance, newBalance uint64, symbol string, changeType string, cfg *config.Config) {
	var message string
	var level alerts.AlertLevel
	var absChange float64

	data := map[string]interface{}{
		"wallet_address": walletAddress,
		"token_mint":     tokenMint,
		"old_balance":    oldBalance,
		"new_balance":    newBalance,
		"token_symbol":   symbol,
	}

	if changeType == "NEW_TOKEN" {
		message = fmt.Sprintf("New token %s (%s) detected in wallet with balance %d", symbol, tokenMint, newBalance)
		level = alerts.Info
	} else if changeType == "REMOVED" {
		message = fmt.Sprintf("Token %s (%s) removed from wallet (previous balance: %d)", symbol, tokenMint, oldBalance)
		level = alerts.Warning
	} else {
		// Calculate percentage change for increase/decrease
		var percentChange float64
		if oldBalance > 0 {
			percentChange = float64(newBalance-oldBalance) / float64(oldBalance)
			absChange = absFloat(percentChange)
		} else {
			percentChange = 1.0
			absChange = 1.0
		}

		// Set alert level based on magnitude of change
		if absChange >= cfg.Alerts.SignificantChange*5 {
			level = alerts.Critical
		} else if absChange >= cfg.Alerts.SignificantChange*2 {
			level = alerts.Warning
		} else {
			level = alerts.Info
		}

		// Create message
		percentStr := fmt.Sprintf("%.1f%%", percentChange*100)
		if percentChange > 0 {
			percentStr = "+" + percentStr
		}

		message = fmt.Sprintf("Token %s (%s) balance %s by %s from %d to %d", symbol, tokenMint, strings.ToLower(changeType), percentStr, oldBalance, newBalance)
	}

	// Send the alert
	alert := alerts.Alert{
		Timestamp:     time.Now(),
		WalletAddress: walletAddress,
		TokenMint:     tokenMint,
		AlertType:     changeType,
		Message:       message,
		Level:         level,
		Data:          data,
	}

	if err := alerter.SendAlert(alert); err != nil {
		log.Printf("Failed to send alert: %v", err)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
