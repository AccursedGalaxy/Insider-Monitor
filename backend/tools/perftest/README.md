# Insider-Monitor Performance Testing Tool

This tool is designed to test and benchmark the performance of the Insider-Monitor backend, focusing on:

1. API endpoints performance under load
2. WebSocket connections and message handling performance
3. Generating reports with optimization recommendations

## Usage

```bash
go run main.go [options]
```

### Options

- `-config string`: Path to performance test configuration file (default "perftest_config.json")
- `-mode string`: Test mode: 'api', 'websocket', or 'all' (default "all")
- `-duration duration`: Duration of the performance test (default 30s)
- `-concurrency int`: Number of concurrent users to simulate (default 10)
- `-ramp-up duration`: Ramp-up period to gradually increase load (default 5s)
- `-verbose`: Enable verbose logging
- `-output string`: Output file for the performance report (default "perftest_report.html")

## Configuration

The test configuration is defined in a JSON file (default: `perftest_config.json`). Example:

```json
{
  "base_url": "http://localhost:8080",
  "concurrency": 20,
  "duration": "30s",
  "ramp_up_period": "5s",
  "verbose": true,
  "api_endpoints": [
    {
      "url": "/health",
      "method": "GET",
      "expected_code": 200,
      "weight": 1,
      "requires_auth": false
    },
    ...
  ],
  "websocket_url": "ws://localhost:8080/ws",
  "message_types": ["subscribe", "unsubscribe", "ping"],
  "message_rate": 2.0,
  "auth_enabled": false,
  "auth_type": "none"
}
```

## Reports

After running the tests, an HTML report is generated with:

1. Overall performance metrics
2. Detailed endpoint-by-endpoint statistics
3. WebSocket connection and message handling metrics
4. Optimization recommendations based on test results

## Building

```bash
cd backend/tools/perftest
go build -o perftest
```

## Running Tests

Example for testing only API endpoints with high concurrency:

```bash
./perftest -mode=api -concurrency=50 -duration=60s
```

Example for testing only WebSockets:

```bash
./perftest -mode=websocket -concurrency=100
```

## Optimizing Performance

Based on the test results, the tool will provide specific recommendations for optimization. Common recommendations include:

1. Adding caching for slow endpoints
2. Optimizing database queries
3. Implementing connection pooling
4. Tuning garbage collection
5. Adding compression for large responses

## Integration with CI/CD

You can run this tool as part of your CI/CD pipeline to identify performance regressions:

```bash
./perftest -output=performance-report.html
# Then add threshold checks for test results
if grep -q "error" performance-report.html; then
  echo "Performance test failed - endpoints exceeding threshold"
  exit 1
fi
```
