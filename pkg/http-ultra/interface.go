package httpserver

import (
	"context"
	"net/http"
	"time"
)

// IHttpServer defines the interface for HTTP server operations
type IHttpServer interface {
	// Start starts the HTTP server
	Start() error

	// StartWithContext starts the HTTP server with context
	StartWithContext(ctx context.Context) error

	// Stop gracefully stops the HTTP server
	Stop(ctx context.Context) error

	// AddRoute adds a route with handler
	AddRoute(method, path string, handler http.HandlerFunc)

	// AddMiddleware adds middleware to the server
	AddMiddleware(middleware ...MiddlewareFunc)

	// GetServer returns the underlying http.Server
	GetServer() *http.Server
}

// MiddlewareFunc defines the middleware function signature
type MiddlewareFunc func(http.Handler) http.Handler

// HandlerFunc defines a custom handler function signature
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Config holds the HTTP server configuration
type Config struct {
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	MaxHeaderBytes    int
	EnableKeepAlives  bool
	EnableCompression bool
	EnableCORS        bool
	CORSOrigins       []string
	CORSMethods       []string
	CORSHeaders       []string
}
