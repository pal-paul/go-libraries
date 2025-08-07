package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	httpserver "github.com/pal-paul/go-libraries/pkg/http-ultra"
	"github.com/valyala/fasthttp"
)

type BenchmarkResult struct {
	Server         string  `json:"server"`
	RequestsPerSec float64 `json:"requests_per_sec"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	P95LatencyMs   float64 `json:"p95_latency_ms"`
	SuccessRate    float64 `json:"success_rate"`
}

type ComparisonResults struct {
	Results   []BenchmarkResult `json:"results"`
	Timestamp time.Time         `json:"timestamp"`
	Config    BenchmarkConfig   `json:"config"`
}

type BenchmarkConfig struct {
	Requests    int `json:"requests"`
	Concurrency int `json:"concurrency"`
}

func startHttpUltraServer(port string) (func(), error) {
	config := httpserver.DefaultConfig()
	config.Host = "localhost"
	config.Port = port

	server, err := httpserver.New(config)
	if err != nil {
		return nil, err
	}

	server.AddRoute("GET", "/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ultra HTTP server error: %v\n", err)
		}
	}()

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Stop(ctx)
		wg.Wait()
	}

	return cleanup, nil
}

func startChiServer(port string) (func(), error) {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: r,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Chi server error: %v\n", err)
		}
	}()

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		wg.Wait()
	}

	return cleanup, nil
}

func startFastHttpServer(port string) (func(), error) {
	handler := func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Path()) == "/ping" {
			ctx.SetContentType("text/plain")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.WriteString("pong")
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}

	server := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe("localhost:" + port); err != nil {
			fmt.Printf("FastHTTP server error: %v\n", err)
		}
	}()

	cleanup := func() {
		server.Shutdown()
		wg.Wait()
	}

	return cleanup, nil
}

func runLoadTest(url string, requests, concurrency int) (*BenchmarkResult, error) {
	fmt.Printf("  Testing %s with %d requests, %d concurrency...\n", url, requests, concurrency)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var (
		successCount int
		latencies    []time.Duration
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	semaphore := make(chan struct{}, concurrency)
	start := time.Now()

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			reqStart := time.Now()
			resp, err := client.Get(url)
			reqEnd := time.Now()
			latency := reqEnd.Sub(reqStart)

			mu.Lock()
			latencies = append(latencies, latency)
			if err == nil && resp.StatusCode == 200 {
				successCount++
				resp.Body.Close()
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)

	successRate := float64(successCount) / float64(requests) * 100
	rps := float64(requests) / totalTime.Seconds()

	if len(latencies) == 0 {
		return nil, fmt.Errorf("no successful requests")
	}

	// Sort latencies
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	avgLatency := time.Duration(0)
	for _, lat := range latencies {
		avgLatency += lat
	}
	avgLatency /= time.Duration(len(latencies))

	p95Index := int(float64(len(latencies)) * 0.95)
	if p95Index >= len(latencies) {
		p95Index = len(latencies) - 1
	}

	return &BenchmarkResult{
		RequestsPerSec: rps,
		AvgLatencyMs:   float64(avgLatency.Nanoseconds()) / 1e6,
		P95LatencyMs:   float64(latencies[p95Index].Nanoseconds()) / 1e6,
		SuccessRate:    successRate,
	}, nil
}

func generateHTMLReport(results ComparisonResults) error {
	labels := getLabelsJS(results.Results)
	rpsData := getRPSDataJS(results.Results)
	avgLatencyData := getAvgLatencyDataJS(results.Results)
	p95LatencyData := getP95LatencyDataJS(results.Results)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>HTTP Server Performance Comparison</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        .chart-container { width: 100%%; height: 400px; margin: 30px 0; }
        .stats-table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
        .stats-table th, .stats-table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .stats-table th { background-color: #f8f9fa; font-weight: bold; }
        .config-info { background: #e9ecef; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .metric-card { display: inline-block; margin: 10px; padding: 20px; background: #f8f9fa; border-radius: 8px; min-width: 200px; }
        .metric-title { font-weight: bold; color: #495057; margin-bottom: 10px; }
        .metric-value { font-size: 24px; font-weight: bold; color: #007bff; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 HTTP Server Performance Comparison</h1>
        <h2>http-ultra vs go-chi vs fasthttp</h2>
        
        <div class="config-info">
            <strong>Test Configuration:</strong>
            Requests: %d | Concurrency: %d | Timestamp: %s
        </div>

        <div class="chart-container">
            <canvas id="rpsChart"></canvas>
        </div>

        <div class="chart-container">
            <canvas id="latencyChart"></canvas>
        </div>

        <table class="stats-table">
            <thead>
                <tr>
                    <th>Server</th>
                    <th>Requests/Sec</th>
                    <th>Avg Latency (ms)</th>
                    <th>P95 Latency (ms)</th>
                    <th>Success Rate (%%)</th>
                </tr>
            </thead>
            <tbody>`,
		results.Config.Requests, results.Config.Concurrency, results.Timestamp.Format("2006-01-02 15:04:05"))

	for _, result := range results.Results {
		html += fmt.Sprintf(`
                <tr>
                    <td><strong>%s</strong></td>
                    <td>%.1f</td>
                    <td>%.2f</td>
                    <td>%.2f</td>
                    <td>%.1f%%</td>
                </tr>`,
			result.Server, result.RequestsPerSec, result.AvgLatencyMs,
			result.P95LatencyMs, result.SuccessRate)
	}

	html += fmt.Sprintf(`
            </tbody>
        </table>

        <div class="metric-card">
            <div class="metric-title">🏆 Fastest Server (RPS)</div>
            <div class="metric-value">%s</div>
        </div>

        <div class="metric-card">
            <div class="metric-title">⚡ Lowest Latency</div>
            <div class="metric-value">%s</div>
        </div>
    </div>

    <script>
        const rpsCtx = document.getElementById('rpsChart').getContext('2d');
        new Chart(rpsCtx, {
            type: 'bar',
            data: {
                labels: [%s],
                datasets: [{
                    label: 'Requests per Second',
                    data: [%s],
                    backgroundColor: ['#007bff', '#28a745', '#ffc107'],
                    borderColor: ['#0056b3', '#1e7e34', '#e0a800'],
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { title: { display: true, text: 'Requests per Second' } },
                scales: { y: { beginAtZero: true } }
            }
        });

        const latencyCtx = document.getElementById('latencyChart').getContext('2d');
        new Chart(latencyCtx, {
            type: 'line',
            data: {
                labels: [%s],
                datasets: [{
                    label: 'Average Latency (ms)',
                    data: [%s],
                    borderColor: '#007bff',
                    backgroundColor: 'rgba(0, 123, 255, 0.1)',
                    tension: 0.4
                }, {
                    label: 'P95 Latency (ms)',
                    data: [%s],
                    borderColor: '#ffc107',
                    backgroundColor: 'rgba(255, 193, 7, 0.1)',
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { title: { display: true, text: 'Latency Comparison' } },
                scales: { y: { beginAtZero: true } }
            }
        });
    </script>
</body>
</html>`,
		getFastestServer(results.Results), getLowestLatencyServer(results.Results),
		labels, rpsData, labels, avgLatencyData, p95LatencyData)

	return os.WriteFile("comparison_report.html", []byte(html), 0o644)
}

func getLabelsJS(results []BenchmarkResult) string {
	var labels []string
	for _, r := range results {
		labels = append(labels, fmt.Sprintf("'%s'", r.Server))
	}
	return strings.Join(labels, ", ")
}

func getRPSDataJS(results []BenchmarkResult) string {
	var data []string
	for _, r := range results {
		data = append(data, fmt.Sprintf("%.1f", r.RequestsPerSec))
	}
	return strings.Join(data, ", ")
}

func getAvgLatencyDataJS(results []BenchmarkResult) string {
	var data []string
	for _, r := range results {
		data = append(data, fmt.Sprintf("%.2f", r.AvgLatencyMs))
	}
	return strings.Join(data, ", ")
}

func getP95LatencyDataJS(results []BenchmarkResult) string {
	var data []string
	for _, r := range results {
		data = append(data, fmt.Sprintf("%.2f", r.P95LatencyMs))
	}
	return strings.Join(data, ", ")
}

func getFastestServer(results []BenchmarkResult) string {
	fastest := results[0]
	for _, r := range results {
		if r.RequestsPerSec > fastest.RequestsPerSec {
			fastest = r
		}
	}
	return fastest.Server
}

func getLowestLatencyServer(results []BenchmarkResult) string {
	lowest := results[0]
	for _, r := range results {
		if r.AvgLatencyMs < lowest.AvgLatencyMs {
			lowest = r
		}
	}
	return lowest.Server
}

func main() {
	fmt.Println("🚀 HTTP Server Performance Comparison")
	fmt.Println("http-ultra vs go-chi vs fasthttp")
	fmt.Println("=====================================")

	config := BenchmarkConfig{
		Requests:    1000,
		Concurrency: 50,
	}

	var results []BenchmarkResult

	// Test http-ultra
	fmt.Println("\n🔧 Starting http-ultra server...")
	cleanup1, err := startHttpUltraServer("8080")
	if err != nil {
		fmt.Printf("Failed to start http-ultra server: %v\n", err)
		return
	}
	time.Sleep(2 * time.Second)

	result1, err := runLoadTest("http://localhost:8080/ping", config.Requests, config.Concurrency)
	if err != nil {
		fmt.Printf("Failed to benchmark http-ultra: %v\n", err)
	} else {
		result1.Server = "http-ultra"
		results = append(results, *result1)
	}
	cleanup1()
	time.Sleep(2 * time.Second)

	// Test Chi
	fmt.Println("\n🔧 Starting go-chi server...")
	cleanup2, err := startChiServer("8081")
	if err != nil {
		fmt.Printf("Failed to start go-chi server: %v\n", err)
		return
	}
	time.Sleep(2 * time.Second)

	result2, err := runLoadTest("http://localhost:8081/ping", config.Requests, config.Concurrency)
	if err != nil {
		fmt.Printf("Failed to benchmark go-chi: %v\n", err)
	} else {
		result2.Server = "go-chi"
		results = append(results, *result2)
	}
	cleanup2()
	time.Sleep(2 * time.Second)

	// Test FastHTTP
	fmt.Println("\n🔧 Starting fasthttp server...")
	cleanup3, err := startFastHttpServer("8082")
	if err != nil {
		fmt.Printf("Failed to start fasthttp server: %v\n", err)
		return
	}
	time.Sleep(2 * time.Second)

	result3, err := runLoadTest("http://localhost:8082/ping", config.Requests, config.Concurrency)
	if err != nil {
		fmt.Printf("Failed to benchmark fasthttp: %v\n", err)
	} else {
		result3.Server = "fasthttp"
		results = append(results, *result3)
	}
	cleanup3()

	// Generate report
	comparisonResults := ComparisonResults{
		Results:   results,
		Timestamp: time.Now(),
		Config:    config,
	}

	// Save JSON results
	jsonData, _ := json.MarshalIndent(comparisonResults, "", "  ")
	os.WriteFile("benchmark_results.json", jsonData, 0o644)

	// Generate HTML report
	if err := generateHTMLReport(comparisonResults); err != nil {
		fmt.Printf("Failed to generate HTML report: %v\n", err)
		return
	}

	fmt.Println("\n📊 Benchmark Results:")
	fmt.Println("=====================")
	for _, result := range results {
		fmt.Printf("%-12s: %8.1f RPS, %6.2fms avg, %6.2fms P95, %5.1f%% success\n",
			result.Server, result.RequestsPerSec, result.AvgLatencyMs, result.P95LatencyMs, result.SuccessRate)
	}

	fmt.Println("\n✅ Reports generated:")
	fmt.Println("   📄 comparison_report.html - Interactive charts")
	fmt.Println("   📊 benchmark_results.json - Raw data")
	fmt.Println("\n🌐 Open comparison_report.html in your browser!")

	// Try to open the HTML report automatically
	exec.Command("open", "comparison_report.html").Start()
}
