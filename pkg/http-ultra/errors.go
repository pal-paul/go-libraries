package httpserver

import (
	"errors"
)

var (
	// ErrServerNotStarted is returned when server is not started
	ErrServerNotStarted = errors.New("http server not started")

	// ErrServerAlreadyStarted is returned when server is already started
	ErrServerAlreadyStarted = errors.New("http server already started")

	// ErrInvalidConfig is returned when configuration is invalid
	ErrInvalidConfig = errors.New("invalid server configuration")

	// ErrInvalidRoute is returned when route is invalid
	ErrInvalidRoute = errors.New("invalid route configuration")

	// ErrServerShutdown is returned during server shutdown
	ErrServerShutdown = errors.New("server shutdown")

	// ErrContextCanceled is returned when context is canceled
	ErrContextCanceled = errors.New("context canceled")

	// ErrTimeout is returned when operation times out
	ErrTimeout = errors.New("operation timeout")
)
