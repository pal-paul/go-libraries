package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Configuration
	const (
		url            = "http://localhost:8080/ping"
		totalRequests  = 10000
		concurrentReqs = 100
	)

	fmt.Printf("🔥 Benchmarking Fast HTTP Server\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Concurrent requests: %d\n", concurrentReqs)
	fmt.Printf("Starting benchmark...\n\n")

	// Create HTTP client optimized for performance
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true, // Disable for raw performance test
		},
	}

	// Channels for coordination
	requests := make(chan int, totalRequests)
	results := make(chan time.Duration, totalRequests)

	// Fill requests channel
	for i := 0; i < totalRequests; i++ {
		requests <- i
	}
	close(requests)

	// Start timer
	startTime := time.Now()

	// Start concurrent workers
	var wg sync.WaitGroup
	for i := 0; i < concurrentReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requests {
				reqStart := time.Now()

				resp, err := client.Get(url)
				if err != nil {
					log.Printf("Request failed: %v", err)
					continue
				}

				resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Printf("Unexpected status code: %d", resp.StatusCode)
					continue
				}

				duration := time.Since(reqStart)
				results <- duration
			}
		}()
	}

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	// Calculate statistics
	totalTime := time.Since(startTime)
	var durations []time.Duration
	var totalDuration time.Duration
	minDuration := time.Hour // Start with large value
	var maxDuration time.Duration

	for duration := range results {
		durations = append(durations, duration)
		totalDuration += duration
		if duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
	}

	successfulRequests := len(durations)
	avgDuration := totalDuration / time.Duration(successfulRequests)
	requestsPerSecond := float64(successfulRequests) / totalTime.Seconds()

	// Print results
	fmt.Printf("📊 Benchmark Results:\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Total time:          %v\n", totalTime)
	fmt.Printf("Successful requests: %d/%d\n", successfulRequests, totalRequests)
	fmt.Printf("Requests per second: %.2f\n", requestsPerSecond)
	fmt.Printf("Average latency:     %v\n", avgDuration)
	fmt.Printf("Min latency:         %v\n", minDuration)
	fmt.Printf("Max latency:         %v\n", maxDuration)

	// Calculate percentiles
	if len(durations) > 0 {
		// Sort durations for percentile calculation
		for i := 0; i < len(durations)-1; i++ {
			for j := i + 1; j < len(durations); j++ {
				if durations[i] > durations[j] {
					durations[i], durations[j] = durations[j], durations[i]
				}
			}
		}

		p50 := durations[len(durations)*50/100]
		p90 := durations[len(durations)*90/100]
		p95 := durations[len(durations)*95/100]
		p99 := durations[len(durations)*99/100]

		fmt.Printf("\n📈 Latency Percentiles:\n")
		fmt.Printf("50th percentile:     %v\n", p50)
		fmt.Printf("90th percentile:     %v\n", p90)
		fmt.Printf("95th percentile:     %v\n", p95)
		fmt.Printf("99th percentile:     %v\n", p99)
	}

	fmt.Printf("\n✅ Benchmark completed!\n")

	// Performance rating
	if requestsPerSecond > 5000 {
		fmt.Printf("🚀 Excellent performance!\n")
	} else if requestsPerSecond > 2000 {
		fmt.Printf("⚡ Good performance!\n")
	} else if requestsPerSecond > 1000 {
		fmt.Printf("👍 Decent performance.\n")
	} else {
		fmt.Printf("⚠️  Performance could be improved.\n")
	}
}
