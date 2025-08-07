# HTTP Server Performance Comparison

A simple comparison tool between **http-ultra**, **go-chi**, and **fasthttp** with interactive graphs.

## 🚀 Quick Start

### Run Performance Comparison with Graphs
```bash
go run benchmark.go
```

This will:
- Start each server individually
- Run load tests (1000 requests, 50 concurrency)
- Generate an interactive HTML report with charts
- Automatically open the results in your browser

### Start Individual Servers
```bash
# Start http-ultra server on port 8080
go run servers.go http-ultra

# Start go-chi server on port 8081  
go run servers.go chi

# Start fasthttp server on port 8082
go run servers.go fasthttp
```

### Test Endpoints
```bash
curl http://localhost:8080/ping  # http-ultra
curl http://localhost:8081/ping  # go-chi
curl http://localhost:8082/ping  # fasthttp
```

## 📊 What Gets Compared

| Server | Port | Key Features |
|--------|------|--------------|
| **http-ultra** | 8080 | Zero-allocation router, TCP tuning, Object pooling |
| **go-chi** | 8081 | Middleware stack, Context support, Route patterns |
| **fasthttp** | 8082 | Zero-copy operations, High concurrency, Async processing |

## 📈 Output

The benchmark generates:
- **comparison_report.html** - Interactive charts showing RPS and latency
- **benchmark_results.json** - Raw performance data

### Metrics Measured
- **Requests per Second (RPS)** - Higher is better
- **Average Latency** - Lower is better  
- **P95 Latency** - 95th percentile response time
- **Success Rate** - Percentage of successful requests

## 🔧 Files

```
examples/comparison/
├── benchmark.go         # Performance comparison tool with graphs
├── servers.go          # Individual server implementations  
├── go.mod              # Go modules
├── go.sum              # Dependency checksums
└── README.md           # This file
```

## 💡 Expected Performance

Typical results (vary by hardware):

1. **fasthttp** - Highest RPS, lowest latency
2. **http-ultra** - Balanced performance with extra features
3. **go-chi** - Good performance with developer-friendly features

## 🎯 Use Cases

- **fasthttp** - Maximum performance, minimal features
- **http-ultra** - High performance with standard library compatibility  
- **go-chi** - Good performance with rich middleware ecosystem

Run `go run benchmark.go` to see how they perform on your system!
