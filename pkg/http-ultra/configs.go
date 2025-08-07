package httpserver

import (
	"time"
)

// DefaultConfig returns a default configuration optimized for ultra performance
func DefaultConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: "8080",

		// Aggressive timeouts for low latency
		ReadTimeout:       1 * time.Second,        // Very fast read timeout
		WriteTimeout:      1 * time.Second,        // Very fast write timeout
		IdleTimeout:       10 * time.Second,       // Short idle timeout
		ReadHeaderTimeout: 100 * time.Millisecond, // Ultra-fast header reading

		// Optimized buffer sizes
		MaxHeaderBytes: 4096, // Smaller headers = faster parsing

		// Performance optimizations
		EnableKeepAlives:  true,  // Reuse connections
		EnableCompression: false, // Disable compression for speed
		EnableCORS:        false, // Disable CORS overhead

		// Empty CORS config for minimal overhead
		CORSOrigins: nil,
		CORSMethods: nil,
		CORSHeaders: nil,
	}
}

// UltraFastConfig returns a configuration optimized for minimum latency
func UltraFastConfig() *Config {
	return DefaultConfig() // Same as default now
}

// LatencyOptimizedConfig returns config specifically for latency optimization
func LatencyOptimizedConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: "8080",

		// Minimal timeouts
		ReadTimeout:       500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 50 * time.Millisecond,

		// Tiny buffers for speed
		MaxHeaderBytes: 2048,

		// All performance features disabled for pure speed
		EnableKeepAlives:  false, // No keep-alive overhead
		EnableCompression: false,
		EnableCORS:        false,

		CORSOrigins: nil,
		CORSMethods: nil,
		CORSHeaders: nil,
	}
}

// Validate validates the server configuration
func (c *Config) Validate() error {
	if c.Host == "" {
		return ErrInvalidConfig
	}
	if c.Port == "" {
		return ErrInvalidConfig
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 1 * time.Second
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 1 * time.Second
	}
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = 10 * time.Second
	}
	if c.ReadHeaderTimeout <= 0 {
		c.ReadHeaderTimeout = 100 * time.Millisecond
	}
	if c.MaxHeaderBytes <= 0 {
		c.MaxHeaderBytes = 4096
	}
	return nil
}

// Address returns the server address
func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}
