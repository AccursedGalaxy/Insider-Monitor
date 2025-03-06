package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/internal/alerts"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/api/routes"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/config"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/monitor"
	"github.com/accursedgalaxy/insider-monitor/backend/internal/websocket"
)

func main() {
	// Initialize configuration
	configService := config.GetInstance()
	cfg := configService.GetConfig()

	// Initialize alerts service
	_ = alerts.GetInstance() // Initialization only, no need to use the variable
	log.Printf("Alerts service initialized with Discord integration: %v", cfg.Discord.Enabled)

	// Get port from environment or use default
	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil && p > 0 {
			port = p
		}
	}

	// Initialize router
	router := routes.SetupRouter()

	// Initialize WebSocket manager
	wsManager := routes.GetWSManager()

	// Initialize WebSocket bridge
	websocket.InitializeWebSocketBridge(wsManager)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Initialize the wallet monitor based on config
	walletMonitor, err := monitor.NewOptimizedWalletMonitor(cfg.NetworkURL, cfg.Wallets)
	if err != nil {
		log.Printf("Warning: Failed to initialize wallet monitor: %v", err)
	} else {
		log.Printf("Wallet monitor initialized with %d wallets", len(cfg.Wallets))
	}

	// Start scanning wallets in the background
	if walletMonitor != nil {
		go func() {
			// Determine scan interval
			scanInterval := 1 * time.Minute // Default
			if cfg.ScanInterval != "" {
				if d, err := time.ParseDuration(cfg.ScanInterval); err == nil {
					scanInterval = d
				}
			}

			ticker := time.NewTicker(scanInterval)
			defer ticker.Stop()

			log.Printf("Starting wallet scanning with interval: %s", scanInterval)

			// Run an initial scan
			ctx := context.Background()

			// Publish scan started event
			monitor.PublishWalletScanStarted()

			// Perform the scan
			data, err := walletMonitor.ScanWallets(ctx)
			if err != nil {
				log.Printf("Warning: Initial wallet scan failed: %v", err)
				monitor.PublishWalletScanError(err)
			} else {
				log.Printf("Initial wallet scan completed")
				monitor.PublishWalletScanComplete(data)
			}

			// Continue scanning at intervals
			for range ticker.C {
				ctx := context.Background()

				// Publish scan started event
				monitor.PublishWalletScanStarted()

				// Perform the scan
				data, err := walletMonitor.ScanWallets(ctx)
				if err != nil {
					log.Printf("Warning: Wallet scan failed: %v", err)
					monitor.PublishWalletScanError(err)
				} else {
					log.Printf("Wallet scan completed")
					monitor.PublishWalletScanComplete(data)
				}
			}
		}()
	}

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
