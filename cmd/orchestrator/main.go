package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Service represents a running service process
type Service struct {
	Name    string
	Cmd     *exec.Cmd
	DoneCh  chan error
	Started bool
}

func main() {
	log.Println("Starting Insider Monitor services...")

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Find the binary directory and set as working directory
	binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Failed to determine binary directory: %v", err)
	}

	// Create wait group for services
	var wg sync.WaitGroup
	serviceErrors := make(chan error, 2)

	// Start backend service
	backendService := &Service{
		Name:   "Backend",
		DoneCh: make(chan error, 1),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startBackendService(ctx, backendService, binDir); err != nil {
			serviceErrors <- err
		}
	}()

	// Start frontend service
	frontendService := &Service{
		Name:   "Frontend",
		DoneCh: make(chan error, 1),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startFrontendService(ctx, frontendService, binDir); err != nil {
			serviceErrors <- err
		}
	}()

	// Wait for signal or service error
	select {
	case <-signalCh:
		log.Println("Received shutdown signal")
	case err := <-serviceErrors:
		log.Printf("Service error: %v", err)
	}

	// Cancel context to signal shutdown to all services
	cancel()

	// Give services a chance to shut down gracefully
	shutdownTimeout := time.Second * 10
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Wait for services to shut down with timeout
	shutdownComplete := make(chan struct{})
	go func() {
		wg.Wait()
		close(shutdownComplete)
	}()

	select {
	case <-shutdownComplete:
		log.Println("All services shut down successfully")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timed out, forcing exit")
	}
}

func startBackendService(ctx context.Context, service *Service, binDir string) error {
	log.Println("Starting backend service...")

	// Configure backend service
	service.Cmd = exec.CommandContext(ctx, filepath.Join(binDir, "backend"))

	// Configure environment
	service.Cmd.Env = os.Environ()

	// Redirect output to this process
	service.Cmd.Stdout = os.Stdout
	service.Cmd.Stderr = os.Stderr

	// Start the process
	if err := service.Cmd.Start(); err != nil {
		return err
	}

	service.Started = true
	log.Printf("Backend service started with PID %d", service.Cmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		service.DoneCh <- service.Cmd.Wait()
		log.Println("Backend service process exited")
	}()

	return nil
}

func startFrontendService(ctx context.Context, service *Service, binDir string) error {
	log.Println("Starting frontend service...")

	// Configure frontend service
	service.Cmd = exec.CommandContext(ctx, filepath.Join(binDir, "frontend"))

	// Configure environment
	service.Cmd.Env = os.Environ()

	// Redirect output to this process
	service.Cmd.Stdout = os.Stdout
	service.Cmd.Stderr = os.Stderr

	// Start the process
	if err := service.Cmd.Start(); err != nil {
		return err
	}

	service.Started = true
	log.Printf("Frontend service started with PID %d", service.Cmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		service.DoneCh <- service.Cmd.Wait()
		log.Println("Frontend service process exited")
	}()

	return nil
}
