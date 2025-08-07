package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/pal-paul/go-libraries/pkg/http-ultra"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Create custom configuration for maximum performance
	config := &httpserver.Config{
		Host:              "0.0.0.0",
		Port:              "8080",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		EnableKeepAlives:  true,
		EnableCompression: true,
		EnableCORS:        true,
		CORSOrigins:       []string{"*"},
		CORSMethods:       []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSHeaders:       []string{"Content-Type", "Authorization"},
	}

	// Create the fast HTTP server
	server, err := httpserver.New(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Cast to IHttpServer for convenience methods
	fastServer := server.(*httpserver.HttpServer)

	// Add custom middleware for request timing
	server.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			log.Printf("%s %s - %v", r.Method, r.URL.Path, duration)
		})
	})

	// Initialize helpers
	responseHelper := httpserver.NewResponseHelper()
	requestHelper := httpserver.NewRequestHelper()

	// Sample data
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}

	// Routes using convenience methods
	fastServer.GET("/", func(w http.ResponseWriter, r *http.Request) {
		responseHelper.JSON(w, http.StatusOK, map[string]interface{}{
			"message": "Welcome to Fast HTTP Server!",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /",
				"GET /health",
				"GET /users",
				"POST /users",
				"GET /ping",
				"GET /slow",
			},
		})
	})

	// Health check endpoint
	fastServer.GET("/health", httpserver.HealthCheckHandler())

	// Fast ping endpoint for benchmarking
	fastServer.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// Users API
	fastServer.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		// Support pagination
		page := requestHelper.GetQueryParamInt(r, "page", 1)
		limit := requestHelper.GetQueryParamInt(r, "limit", 10)

		responseHelper.JSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
			"pagination": map[string]int{
				"page":  page,
				"limit": limit,
				"total": len(users),
			},
		})
	})

	fastServer.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		var newUser User
		if err := requestHelper.ParseJSON(r, &newUser); err != nil {
			responseHelper.Error(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		newUser.ID = len(users) + 1
		users = append(users, newUser)

		responseHelper.Success(w, newUser)
	})

	// Simulate slow endpoint for testing
	fastServer.GET("/slow", func(w http.ResponseWriter, r *http.Request) {
		delay := requestHelper.GetQueryParamInt(r, "delay", 1000)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		responseHelper.JSON(w, http.StatusOK, map[string]interface{}{
			"message":  "Slow response",
			"delay_ms": delay,
		})
	})

	// JSON echo endpoint
	fastServer.POST("/echo", func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		if err := requestHelper.ParseJSON(r, &data); err != nil {
			responseHelper.Error(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		responseHelper.JSON(w, http.StatusOK, map[string]interface{}{
			"echo":       data,
			"client_ip":  requestHelper.GetClientIP(r),
			"user_agent": requestHelper.GetHeader(r, "User-Agent"),
		})
	})

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Fast HTTP Server starting on %s", config.Address())
		log.Printf("📊 Performance optimizations enabled:")
		log.Printf("   - Keep-alive connections: %v", config.EnableKeepAlives)
		log.Printf("   - Gzip compression: %v", config.EnableCompression)
		log.Printf("   - CORS support: %v", config.EnableCORS)
		log.Printf("   - Connection timeouts optimized")
		log.Printf("📡 Available endpoints:")
		log.Printf("   GET  / - API info")
		log.Printf("   GET  /health - Health check")
		log.Printf("   GET  /ping - Fast ping")
		log.Printf("   GET  /users - List users")
		log.Printf("   POST /users - Create user")
		log.Printf("   POST /echo - Echo JSON")
		log.Printf("   GET  /slow?delay=1000 - Slow endpoint")

		if err := server.StartWithContext(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
