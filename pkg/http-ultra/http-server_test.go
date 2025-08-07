package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewFastHTTPServer(t *testing.T) {
	// Test with default config
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if server == nil {
		t.Fatal("Expected server to be created")
	}

	// Test with custom config
	config := &Config{
		Host:              "localhost",
		Port:              "9090",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		MaxHeaderBytes:    512 * 1024,
		EnableKeepAlives:  true,
		EnableCompression: true,
		EnableCORS:        true,
	}

	server, err = New(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if server == nil {
		t.Fatal("Expected server to be created")
	}
}

func TestFastHTTPServer_AddRoute(t *testing.T) {
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Add a simple route
	server.AddRoute("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Test the route using httptest
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Get the underlying server's handler and serve the request
	handler := server.GetServer().Handler
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "Hello, World!"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestFastHTTPServer_RESTMethods(t *testing.T) {
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	fastServer := server.(*HttpServer)

	// Test all REST methods
	fastServer.GET("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET"))
	})
	fastServer.POST("/post", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POST"))
	})
	fastServer.PUT("/put", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PUT"))
	})
	fastServer.DELETE("/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DELETE"))
	})
	fastServer.PATCH("/patch", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PATCH"))
	})

	handler := server.GetServer().Handler

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{"GET", "/get", "GET"},
		{"POST", "/post", "POST"},
		{"PUT", "/put", "PUT"},
		{"DELETE", "/delete", "DELETE"},
		{"PATCH", "/patch", "PATCH"},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for %s %s, got %d", http.StatusOK, test.method, test.path, w.Code)
		}

		if w.Body.String() != test.expected {
			t.Errorf("Expected body %q for %s %s, got %q", test.expected, test.method, test.path, w.Body.String())
		}
	}
}

func TestFastHTTPServer_NotFound(t *testing.T) {
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	handler := server.GetServer().Handler
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestFastHTTPServer_Middleware(t *testing.T) {
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Add middleware that adds a custom header
	server.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	})

	server.AddRoute("GET", "/middleware-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("middleware"))
	})

	req := httptest.NewRequest("GET", "/middleware-test", nil)
	w := httptest.NewRecorder()

	handler := server.GetServer().Handler
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("X-Test-Middleware") != "true" {
		t.Error("Expected middleware header to be set")
	}
}

func TestFastHTTPServer_JSON(t *testing.T) {
	server, err := New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	type Response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	server.AddRoute("GET", "/json", func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Message: "Hello, JSON!",
			Status:  "success",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	req := httptest.NewRequest("GET", "/json", nil)
	w := httptest.NewRecorder()

	handler := server.GetServer().Handler
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Message != "Hello, JSON!" {
		t.Errorf("Expected message 'Hello, JSON!', got %s", response.Message)
	}
}

func TestServerLifecycle(t *testing.T) {
	config := &Config{
		Host: "localhost",
		Port: "0", // Use port 0 to get a random available port
	}

	server, err := New(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	fastServer := server.(*HttpServer)

	// Test that server is not running initially
	if fastServer.IsRunning() {
		t.Error("Expected server to not be running initially")
	}

	// Start server in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		server.StartWithContext(ctx)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is running
	if !fastServer.IsRunning() {
		t.Error("Expected server to be running after start")
	}

	// Stop server
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := server.Stop(stopCtx); err != nil {
		t.Errorf("Expected no error stopping server, got %v", err)
	}

	// Test that server is not running after stop
	if fastServer.IsRunning() {
		t.Error("Expected server to not be running after stop")
	}
}

func BenchmarkFastHTTPServer_SimpleRoute(b *testing.B) {
	server, err := New(nil)
	if err != nil {
		b.Fatalf("Expected no error, got %v", err)
	}

	server.AddRoute("GET", "/benchmark", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := server.GetServer().Handler

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/benchmark", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		}
	})
}
