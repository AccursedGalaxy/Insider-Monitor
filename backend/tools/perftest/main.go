package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/api"
	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/config"
	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/report"
	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/websocket"
)

func main() {
	// Create our own flag set
	fs := flag.NewFlagSet("perftest", flag.ExitOnError)

	// Define flags
	configFile := fs.String("config", "perftest_config.json", "Path to performance test configuration file")
	mode := fs.String("mode", "all", "Test mode: 'api', 'websocket', or 'all'")
	duration := fs.Duration("duration", 30*time.Second, "Duration of the performance test")
	concurrency := fs.Int("concurrency", 10, "Number of concurrent users to simulate")
	rampUpPeriod := fs.Duration("ramp-up", 5*time.Second, "Ramp-up period to gradually increase load")
	verbose := fs.Bool("verbose", false, "Enable verbose logging")
	outputFile := fs.String("output", "perftest_report.html", "Output file for the performance report")

	// Define a variable to track which flags were explicitly set
	flagsSet := make(map[string]bool)

	// Parse the flag set
	fs.Parse(os.Args[1:])

	// Visit all flags to see which ones were set by the user
	fs.Visit(func(f *flag.Flag) {
		flagsSet[f.Name] = true
	})

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override config with command line flags if specified
	if flagsSet["concurrency"] {
		cfg.Concurrency = *concurrency
	}
	if flagsSet["duration"] {
		cfg.Duration = *duration
	}
	if flagsSet["ramp-up"] {
		cfg.RampUpPeriod = *rampUpPeriod
	}
	if flagsSet["verbose"] {
		cfg.Verbose = *verbose
	}

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		log.Println("Received shutdown signal, stopping tests...")
		cancel()
	}()

	// Initialize results collector
	results := report.NewResults()

	// Run the selected tests
	switch *mode {
	case "api":
		log.Println("Running API performance tests...")
		runAPITests(ctx, cfg, results)
	case "websocket":
		log.Println("Running WebSocket performance tests...")
		runWebSocketTests(ctx, cfg, results)
	case "all":
		log.Println("Running all performance tests...")
		runAPITests(ctx, cfg, results)
		runWebSocketTests(ctx, cfg, results)
	default:
		log.Fatalf("Unknown test mode: %s", *mode)
	}

	// Generate and save the performance report
	log.Println("Generating performance report...")
	err = report.GenerateHTMLReport(results, *outputFile)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	log.Printf("Performance test completed. Report saved to %s", *outputFile)

	// Print optimization recommendations based on test results
	recommendations := results.GenerateOptimizationRecommendations()
	fmt.Println("\nOptimization Recommendations:")
	for i, rec := range recommendations {
		fmt.Printf("%d. %s\n", i+1, rec)
	}
}

func runAPITests(ctx context.Context, cfg *config.Config, results *report.Results) {
	tester := api.NewAPILoadTester(cfg)

	// Test endpoints defined in the configuration
	for _, endpoint := range cfg.APIEndpoints {
		log.Printf("Testing endpoint: %s %s", endpoint.Method, endpoint.URL)

		// Run the test for this endpoint
		endpointResults, err := tester.TestEndpoint(ctx, endpoint)
		if err != nil {
			log.Printf("Error testing endpoint %s: %v", endpoint.URL, err)
			continue
		}

		// Add results to the collector
		results.AddAPIResult(endpoint.URL, endpoint.Method, endpointResults)
	}
}

func runWebSocketTests(ctx context.Context, cfg *config.Config, results *report.Results) {
	tester := websocket.NewWebSocketLoadTester(cfg)

	// Run WebSocket connection test
	wsResults, err := tester.TestWebSocketPerformance(ctx)
	if err != nil {
		log.Printf("Error testing WebSocket: %v", err)
		return
	}

	// Add results to the collector
	results.AddWebSocketResult(wsResults)
}
