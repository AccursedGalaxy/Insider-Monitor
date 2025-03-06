package report

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"

	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/api"
	"github.com/accursedgalaxy/insider-monitor/backend/tools/perftest/websocket"
)

// Results collects all test results
type Results struct {
	APIResults       map[string]*api.APITestResult  // Map of endpoint URL to API test result
	WebSocketResults *websocket.WebSocketTestResult // WebSocket test result
	StartTime        time.Time                      // Start time of the entire test
	EndTime          time.Time                      // End time of the entire test
}

// NewResults creates a new Results instance
func NewResults() *Results {
	return &Results{
		APIResults: make(map[string]*api.APITestResult),
		StartTime:  time.Now(),
	}
}

// AddAPIResult adds an API test result
func (r *Results) AddAPIResult(url, method string, result *api.APITestResult) {
	key := fmt.Sprintf("%s %s", method, url)
	r.APIResults[key] = result
}

// AddWebSocketResult adds a WebSocket test result
func (r *Results) AddWebSocketResult(result *websocket.WebSocketTestResult) {
	r.WebSocketResults = result
}

// Finalize finalizes the results
func (r *Results) Finalize() {
	r.EndTime = time.Now()
}

// GenerateHTMLReport generates an HTML report of the test results
func GenerateHTMLReport(results *Results, outputFile string) error {
	// Finalize results
	results.Finalize()

	// Create the HTML report template
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Performance Test Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        .header {
            background-color: #f8f9fa;
            padding: 10px 20px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        .summary {
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
            margin-bottom: 30px;
        }
        .summary-box {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            flex: 1;
            min-width: 200px;
        }
        .good {
            color: green;
        }
        .warning {
            color: orange;
        }
        .error {
            color: red;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 30px;
        }
        th, td {
            padding: 12px 15px;
            border-bottom: 1px solid #ddd;
            text-align: left;
        }
        th {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .chart-container {
            margin: 30px 0;
        }
        .recommendations {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 5px;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Performance Test Report</h1>
        <p>Test started: {{.StartTime.Format "Jan 02, 2006 15:04:05"}}</p>
        <p>Test duration: {{formatDuration (timeSince .StartTime .EndTime)}}</p>
    </div>

    <h2>Summary</h2>
    <div class="summary">
        <div class="summary-box">
            <h3>API Endpoints</h3>
            <p>Total endpoints tested: {{len .APIResults}}</p>
            <p>Total requests: {{apiTotalRequests .APIResults}}</p>
            <p>Average response time: {{apiAverageResponseTime .APIResults}}</p>
            <p>Success rate: {{apiSuccessRate .APIResults}}%</p>
        </div>

        {{if .WebSocketResults}}
        <div class="summary-box">
            <h3>WebSocket</h3>
            <p>Total connections: {{.WebSocketResults.TotalConnections}}</p>
            <p>Successful connections: {{.WebSocketResults.SuccessfulConnections}}</p>
            <p>Messages sent: {{.WebSocketResults.TotalMessagesSent}}</p>
            <p>Messages received: {{.WebSocketResults.TotalMessagesReceived}}</p>
            <p>Average connection time: {{formatDuration .WebSocketResults.AverageConnectTime}}</p>
            <p>Average message time: {{formatDuration .WebSocketResults.AverageMessageTime}}</p>
        </div>
        {{end}}
    </div>

    <h2>API Endpoints</h2>
    <table>
        <thead>
            <tr>
                <th>Endpoint</th>
                <th>Requests</th>
                <th>Success</th>
                <th>Avg Time</th>
                <th>P95 Time</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            {{range $key, $result := .APIResults}}
            <tr>
                <td>{{$result.Method}} {{$result.URL}}</td>
                <td>{{$result.TotalRequests}}</td>
                <td>{{percentage $result.SuccessfulRequests $result.TotalRequests}}%</td>
                <td>{{formatDuration $result.AverageResponseTime}}</td>
                <td>{{formatDuration $result.P95ResponseTime}}</td>
                <td class="{{apiEndpointStatus $result}}">
                    {{if eq (apiEndpointStatus $result) "good"}}Good{{end}}
                    {{if eq (apiEndpointStatus $result) "warning"}}Warning{{end}}
                    {{if eq (apiEndpointStatus $result) "error"}}Error{{end}}
                </td>
            </tr>
            {{end}}
        </tbody>
    </table>

    <div class="recommendations">
        <h2>Optimization Recommendations</h2>
        <ul>
            {{range .Recommendations}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`

	// Create template functions
	funcMap := template.FuncMap{
		"formatDuration": func(d time.Duration) string {
			return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000.0)
		},
		"timeSince": func(start, end time.Time) time.Duration {
			return end.Sub(start)
		},
		"apiTotalRequests": func(results map[string]*api.APITestResult) int {
			total := 0
			for _, result := range results {
				total += result.TotalRequests
			}
			return total
		},
		"apiAverageResponseTime": func(results map[string]*api.APITestResult) string {
			var total time.Duration
			var count int
			for _, result := range results {
				total += result.AverageResponseTime * time.Duration(result.TotalRequests)
				count += result.TotalRequests
			}
			if count == 0 {
				return "0 ms"
			}
			avg := total / time.Duration(count)
			return fmt.Sprintf("%.2f ms", float64(avg.Microseconds())/1000.0)
		},
		"apiSuccessRate": func(results map[string]*api.APITestResult) float64 {
			var totalRequests, successfulRequests int
			for _, result := range results {
				totalRequests += result.TotalRequests
				successfulRequests += result.SuccessfulRequests
			}
			if totalRequests == 0 {
				return 100.0
			}
			return float64(successfulRequests) / float64(totalRequests) * 100.0
		},
		"percentage": func(part, total int) float64 {
			if total == 0 {
				return 0.0
			}
			return float64(part) / float64(total) * 100.0
		},
		"apiEndpointStatus": func(result *api.APITestResult) string {
			if result.AverageResponseTime > 500*time.Millisecond {
				return "error"
			}
			if result.AverageResponseTime > 200*time.Millisecond {
				return "warning"
			}
			successRate := float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100.0
			if successRate < 95 {
				return "error"
			}
			if successRate < 99 {
				return "warning"
			}
			return "good"
		},
	}

	// Parse template
	t, err := template.New("report").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	// Generate optimization recommendations
	recommendations := results.GenerateOptimizationRecommendations()

	// Create data for template
	data := struct {
		*Results
		Recommendations []string
	}{
		Results:         results,
		Recommendations: recommendations,
	}

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template
	return t.Execute(file, data)
}

// GenerateOptimizationRecommendations generates optimization recommendations based on test results
func (r *Results) GenerateOptimizationRecommendations() []string {
	var recommendations []string

	// API recommendations
	slowEndpoints := make([]string, 0)
	highErrorRateEndpoints := make([]string, 0)

	for key, result := range r.APIResults {
		// Check for slow endpoints
		if result.AverageResponseTime > 200*time.Millisecond {
			slowEndpoints = append(slowEndpoints, key)
		}

		// Check for endpoints with high error rates
		successRate := float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100.0
		if successRate < 95 {
			highErrorRateEndpoints = append(highErrorRateEndpoints, key)
		}
	}

	// Sort slow endpoints by average response time (slowest first)
	sort.Slice(slowEndpoints, func(i, j int) bool {
		return r.APIResults[slowEndpoints[i]].AverageResponseTime > r.APIResults[slowEndpoints[j]].AverageResponseTime
	})

	// Add recommendations for slow endpoints
	if len(slowEndpoints) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Consider optimizing the following slow endpoints: %v", slowEndpoints))

		// Add specific recommendations for the slowest endpoint
		slowestEndpoint := slowEndpoints[0]
		result := r.APIResults[slowestEndpoint]
		recommendations = append(recommendations, fmt.Sprintf("The slowest endpoint is %s with average response time of %.2f ms. Consider adding caching or optimizing database queries.",
			slowestEndpoint, float64(result.AverageResponseTime.Microseconds())/1000.0))
	}

	// Add recommendations for endpoints with high error rates
	if len(highErrorRateEndpoints) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Fix high error rates in the following endpoints: %v", highErrorRateEndpoints))
	}

	// WebSocket recommendations
	if r.WebSocketResults != nil {
		// Check for WebSocket connection issues
		if r.WebSocketResults.FailedConnections > 0 {
			failRate := float64(r.WebSocketResults.FailedConnections) / float64(r.WebSocketResults.TotalConnections) * 100.0
			if failRate > 5 {
				recommendations = append(recommendations, fmt.Sprintf("WebSocket connections have a high failure rate (%.1f%%). Consider increasing the maximum number of concurrent connections on the server.", failRate))
			}
		}

		// Check for message processing issues
		if r.WebSocketResults.TotalMessagesSent > 0 && r.WebSocketResults.TotalMessagesReceived > 0 {
			messageReceiveRate := float64(r.WebSocketResults.TotalMessagesReceived) / float64(r.WebSocketResults.TotalMessagesSent) * 100.0
			if messageReceiveRate < 95 {
				recommendations = append(recommendations, fmt.Sprintf("WebSocket message receive rate is low (%.1f%%). Consider optimizing message handling and implement retry mechanisms.", messageReceiveRate))
			}
		}

		// Check for slow message processing
		if r.WebSocketResults.AverageMessageTime > 100*time.Millisecond {
			recommendations = append(recommendations, fmt.Sprintf("WebSocket message processing is slow (%.2f ms). Consider optimizing the message handling logic.", float64(r.WebSocketResults.AverageMessageTime.Microseconds())/1000.0))
		}
	}

	// General recommendations
	recommendations = append(recommendations, "Consider implementing connection pooling for database and external service connections.")
	recommendations = append(recommendations, "Add caching for frequently accessed data to reduce response times.")
	recommendations = append(recommendations, "Implement rate limiting to prevent resource exhaustion during traffic spikes.")
	recommendations = append(recommendations, "Monitor memory usage and implement garbage collection tuning if needed.")

	return recommendations
}
