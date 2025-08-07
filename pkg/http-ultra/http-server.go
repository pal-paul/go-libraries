package httpserver

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

// HttpServer implements an ultra-optimized HTTP server for maximum performance
type HttpServer struct {
	config  *Config
	server  *http.Server
	router  *Router
	running int32 // Use atomic instead of mutex for better performance
}

// New creates a new ultra-fast HTTP server instance
func New(config *Config) (IHttpServer, error) {
	if config == nil {
		config = UltraFastConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	router := NewRouter()

	// Custom transport for optimal performance
	server := &http.Server{
		Addr:              config.Address(),
		Handler:           router,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,

		// Optimizations for latency
		DisableGeneralOptionsHandler: true, // Skip OPTIONS handling overhead
	}

	// Disable keep-alives if configured for pure speed
	server.SetKeepAlivesEnabled(config.EnableKeepAlives)

	httpServer := &HttpServer{
		config: config,
		server: server,
		router: router,
	}

	return httpServer, nil
}

// Start starts the HTTP server with optimized settings
func (s *HttpServer) Start() error {
	return s.StartWithContext(context.Background())
}

// StartWithContext starts the HTTP server with context and TCP optimizations
func (s *HttpServer) StartWithContext(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return ErrServerAlreadyStarted
	}

	// Create optimized listener
	listener, err := s.createOptimizedListener()
	if err != nil {
		atomic.StoreInt32(&s.running, 0)
		return err
	}

	// Start server with context handling
	if ctx != context.Background() {
		go func() {
			<-ctx.Done()
			s.Stop(context.Background())
		}()
	}

	if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		atomic.StoreInt32(&s.running, 0)
		return err
	}

	return nil
}

// createOptimizedListener creates a TCP listener with optimizations
func (s *HttpServer) createOptimizedListener() (net.Listener, error) {
	addr := s.config.Address()

	// Create TCP listener with optimizations
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Apply TCP optimizations if possible
	if tcpListener, ok := listener.(*net.TCPListener); ok {
		// Wrap with optimized TCP listener
		return &OptimizedTCPListener{tcpListener}, nil
	}

	return listener, nil
}

// OptimizedTCPListener wraps TCPListener with performance optimizations
type OptimizedTCPListener struct {
	*net.TCPListener
}

// Accept accepts connections with TCP optimizations
func (l *OptimizedTCPListener) Accept() (net.Conn, error) {
	conn, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}

	// Apply TCP optimizations
	conn.SetNoDelay(true)                     // Disable Nagle algorithm for low latency
	conn.SetKeepAlive(true)                   // Enable keep-alive
	conn.SetKeepAlivePeriod(30 * time.Second) // Set keep-alive period

	return conn, nil
}

// Stop gracefully stops the HTTP server
func (s *HttpServer) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		return ErrServerNotStarted
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return s.server.Shutdown(ctx)
}

// AddRoute adds a route with handler
func (s *HttpServer) AddRoute(method, path string, handler http.HandlerFunc) {
	if method == "" || path == "" || handler == nil {
		return
	}
	s.router.AddRoute(method, path, handler)
}

// AddMiddleware adds middleware to the server (simplified for ultra-fast performance)
func (s *HttpServer) AddMiddleware(middleware ...MiddlewareFunc) {
	// Store middleware in router for minimal overhead
	s.router.AddMiddleware(middleware...)
}

// GetServer returns the underlying http.Server
func (s *HttpServer) GetServer() *http.Server {
	return s.server
}

// Router returns the router for advanced usage
func (s *HttpServer) Router() *Router {
	return s.router
}

// IsRunning returns whether the server is currently running (thread-safe)
func (s *HttpServer) IsRunning() bool {
	return atomic.LoadInt32(&s.running) == 1
}

// GET adds a GET route
func (s *HttpServer) GET(path string, handler http.HandlerFunc) {
	s.router.GET(path, handler)
}

// POST adds a POST route
func (s *HttpServer) POST(path string, handler http.HandlerFunc) {
	s.router.POST(path, handler)
}

// PUT adds a PUT route
func (s *HttpServer) PUT(path string, handler http.HandlerFunc) {
	s.router.PUT(path, handler)
}

// DELETE adds a DELETE route
func (s *HttpServer) DELETE(path string, handler http.HandlerFunc) {
	s.router.DELETE(path, handler)
}

// PATCH adds a PATCH route
func (s *HttpServer) PATCH(path string, handler http.HandlerFunc) {
	s.router.PATCH(path, handler)
}

// OPTIONS adds an OPTIONS route
func (s *HttpServer) OPTIONS(path string, handler http.HandlerFunc) {
	s.router.OPTIONS(path, handler)
}

// HEAD adds a HEAD route
func (s *HttpServer) HEAD(path string, handler http.HandlerFunc) {
	s.router.HEAD(path, handler)
}
