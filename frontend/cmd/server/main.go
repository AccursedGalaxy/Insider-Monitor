package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/accursedgalaxy/insider-monitor/frontend/internal/server"
)

func main() {
	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Configure server
	port := 3000
	if envPort := os.Getenv("FRONTEND_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	// Get the path to the frontend static files
	frontendDir := filepath.Join("frontend", "public")
	if envDir := os.Getenv("FRONTEND_DIR"); envDir != "" {
		frontendDir = envDir
	}

	// Get API and WebSocket endpoints
	apiEndpoint := "http://localhost:8080"
	if envAPI := os.Getenv("API_ENDPOINT"); envAPI != "" {
		apiEndpoint = envAPI
	}

	wsEndpoint := "ws://localhost:8080"
	if envWS := os.Getenv("WS_ENDPOINT"); envWS != "" {
		wsEndpoint = envWS
	}

	// Create frontend server
	frontendServer := server.NewServer(frontendDir, apiEndpoint, wsEndpoint)

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: frontendServer.Handler(),
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Frontend server starting on port %d serving files from %s", port, frontendDir)
		log.Printf("API endpoint: %s", apiEndpoint)
		log.Printf("WebSocket endpoint: %s", wsEndpoint)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Frontend server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-signalCh
	log.Println("Shutdown signal received, shutting down frontend server gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown the server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Frontend server shutdown error: %v", err)
	}

	log.Println("Frontend server shutdown complete")
}
