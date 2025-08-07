# HTTP Server Package

A high-performance HTTP server implementation for Go that focuses on speed and simplicity.

## Features

- **Ultra-fast routing**: Custom router optimized for speed with minimal allocations
- **Connection pooling**: Efficient connection reuse and keep-alive support
- **Built-in middleware**: Compression, CORS, recovery, and logging middleware
- **Graceful shutdown**: Context-based server lifecycle management
- **Minimal allocations**: Optimized for low memory usage and garbage collection
- **Thread-safe**: Concurrent-safe operations with proper mutex usage
- **Configurable timeouts**: Customizable read, write, and idle timeouts
- **JSON support**: Easy JSON request/response handling

## Performance Optimizations

1. **Custom Router**: Fast path-based routing without regex overhead
2. **Object Pooling**: Reused gzip writers and other objects
3. **Minimal Middleware Chain**: Lightweight middleware execution
4. **Keep-Alive Connections**: Connection reuse for better performance
5. **Optimized Headers**: Efficient header management
6. **Context-Aware**: Proper context handling for cancellation

## Quick Start

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    httpserver "your-module/pkg/http-server"
)

func main() {
    // Create server with default config
    server, err := httpserver.NewFastHTTPServer(nil)
    if err != nil {
        log.Fatal(err)
    }

    // Add routes
    server.AddRoute("GET", "/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    server.AddRoute("GET", "/json", func(w http.ResponseWriter, r *http.Request) {
        response := map[string]interface{}{
            "message": "Hello, JSON!",
            "timestamp": time.Now(),
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    // Start server
    log.Println("Starting server on :8080")
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Advanced Usage

### Custom Configuration

```go
config := &httpserver.Config{
    Host:               "0.0.0.0",
    Port:               "8080",
    ReadTimeout:        10 * time.Second,
    WriteTimeout:       10 * time.Second,
    IdleTimeout:        120 * time.Second,
    ReadHeaderTimeout:  5 * time.Second,
    MaxHeaderBytes:     1 << 20, // 1 MB
    EnableKeepAlives:   true,
    EnableCompression:  true,
    EnableCORS:         true,
    CORSOrigins:        []string{"*"},
    CORSMethods:        []string{"GET", "POST", "PUT", "DELETE"},
    CORSHeaders:        []string{"Content-Type", "Authorization"},
}

server, err := httpserver.NewFastHTTPServer(config)
```

### REST Methods

```go
// Using convenience methods
fastServer := server.(*httpserver.FastHTTPServer)

fastServer.GET("/users", getUsersHandler)
fastServer.POST("/users", createUserHandler)
fastServer.PUT("/users/:id", updateUserHandler)
fastServer.DELETE("/users/:id", deleteUserHandler)
```

### Custom Middleware

```go
// Add custom middleware
server.AddMiddleware(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic before request
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        // Custom logic after request
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
    })
})
```

### Graceful Shutdown

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Start server with context
go func() {
    if err := server.StartWithContext(ctx); err != nil {
        log.Printf("Server error: %v", err)
    }
}()

// Handle shutdown signal
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c

// Graceful shutdown
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
defer shutdownCancel()

if err := server.Stop(shutdownCtx); err != nil {
    log.Printf("Server shutdown error: %v", err)
}
```

## Benchmark Results

The server is optimized for high performance:

- **Low latency**: Sub-millisecond response times for simple routes
- **High throughput**: Capable of handling thousands of requests per second
- **Low memory usage**: Minimal allocations and efficient garbage collection
- **Connection efficiency**: Keep-alive and connection pooling support

## API Reference

### HTTPServer Interface

```go
type HTTPServer interface {
    Start() error
    StartWithContext(ctx context.Context) error
    Stop(ctx context.Context) error
    AddRoute(method, path string, handler http.HandlerFunc)
    AddMiddleware(middleware ...MiddlewareFunc)
    GetServer() *http.Server
}
```

### Configuration

```go
type Config struct {
    Host               string
    Port               string
    ReadTimeout        time.Duration
    WriteTimeout       time.Duration
    IdleTimeout        time.Duration
    ReadHeaderTimeout  time.Duration
    MaxHeaderBytes     int
    EnableKeepAlives   bool
    EnableCompression  bool
    EnableCORS         bool
    CORSOrigins        []string
    CORSMethods        []string
    CORSHeaders        []string
}
```

### Built-in Middleware

- **CompressionMiddleware()**: Gzip compression for responses
- **CORSMiddleware()**: Cross-Origin Resource Sharing support
- **RecoveryMiddleware()**: Panic recovery and error handling
- **LoggingMiddleware()**: Basic request logging

## Error Handling

The package provides comprehensive error handling:

```go
var (
    ErrServerNotStarted      = errors.New("http server not started")
    ErrServerAlreadyStarted  = errors.New("http server already started")
    ErrInvalidConfig         = errors.New("invalid server configuration")
    ErrInvalidRoute          = errors.New("invalid route configuration")
    ErrServerShutdown        = errors.New("server shutdown")
    ErrContextCanceled       = errors.New("context canceled")
    ErrTimeout               = errors.New("operation timeout")
)
```

## Best Practices

1. **Use connection pooling**: Enable keep-alives for better performance
2. **Configure timeouts**: Set appropriate timeouts for your use case
3. **Use middleware wisely**: Only add necessary middleware to minimize overhead
4. **Handle context properly**: Use context for cancellation and timeouts
5. **Monitor performance**: Use built-in metrics and logging
6. **Graceful shutdown**: Always implement proper shutdown handling

## Testing

Run the test suite:

```bash
go test ./pkg/http-server -v
```

Run benchmarks:

```bash
go test ./pkg/http-server -bench=. -benchmem
```

## Contributing

1. Ensure all tests pass
2. Add benchmarks for performance-critical code
3. Follow Go best practices
4. Update documentation for API changes
