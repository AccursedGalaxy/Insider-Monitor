#!/bin/bash

# Script to run performance tests against the Insider-Monitor backend
# This script:
# 1. Builds the performance testing tool
# 2. Starts the backend server
# 3. Runs performance tests
# 4. Shuts down the backend server
# 5. Opens the performance report

# Build the performance testing tool
echo "Building performance testing tool..."
cd "$(dirname "$0")"
go build -o perftest

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Failed to build performance testing tool!"
    exit 1
fi

# Set default configuration
BACKEND_PATH="../../cmd/server"
TEST_MODE="all"
CONCURRENCY=20
DURATION="30s"
RAMPUP="5s"
OUTPUT_FILE="perftest_report.html"
VERBOSE=false
CONFIG="perftest_config.json"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        --backend)
            BACKEND_PATH="$2"
            shift
            shift
            ;;
        --mode)
            TEST_MODE="$2"
            shift
            shift
            ;;
        --concurrency)
            CONCURRENCY="$2"
            shift
            shift
            ;;
        --duration)
            DURATION="$2"
            shift
            shift
            ;;
        --rampup)
            RAMPUP="$2"
            shift
            shift
            ;;
        --output)
            OUTPUT_FILE="$2"
            shift
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --config)
            CONFIG="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--backend path] [--mode api|websocket|all] [--concurrency N] [--duration Xs] [--rampup Xs] [--output file.html] [--verbose] [--config file.json]"
            exit 1
            ;;
    esac
done

# Build the backend if needed
echo "Building backend..."
cd "$BACKEND_PATH"
go build -o server
if [ $? -ne 0 ]; then
    echo "Failed to build backend server!"
    exit 1
fi

# Start the backend server
echo "Starting backend server..."
./server &
SERVER_PID=$!

# Give the server time to start
echo "Waiting for server to start..."
sleep 5

# Run performance tests
echo "Running performance tests..."
cd - > /dev/null
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-verbose"
fi

./perftest -mode="$TEST_MODE" -concurrency="$CONCURRENCY" -duration="$DURATION" -ramp-up="$RAMPUP" -output="$OUTPUT_FILE" -config="$CONFIG" $VERBOSE_FLAG

# Check if tests were successful
if [ $? -ne 0 ]; then
    echo "Performance tests failed!"
    # Clean up
    kill $SERVER_PID
    exit 1
fi

# Shutdown the server
echo "Shutting down backend server..."
kill $SERVER_PID

# Wait for server to shutdown
sleep 2

# Display results
echo "Performance test completed successfully!"
echo "Test report saved to: $OUTPUT_FILE"

# Open the report in the default browser if on Linux
if command -v xdg-open > /dev/null; then
    echo "Opening test report..."
    xdg-open "$OUTPUT_FILE"
elif command -v open > /dev/null; then
    echo "Opening test report..."
    open "$OUTPUT_FILE"
else
    echo "To view the report, open $OUTPUT_FILE in a web browser"
fi

echo "Done!"
exit 0
