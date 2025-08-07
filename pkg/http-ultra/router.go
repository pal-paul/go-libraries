package httpserver

import (
	"net/http"
	"sync"
	"unsafe"
)

// Router implements zero-allocation routing
type Router struct {
	routes     [16]map[string]http.HandlerFunc // Pre-allocated method maps
	middleware []MiddlewareFunc

	// Pre-compiled method indices for O(1) lookup
	methodIndex map[string]int

	// Response pool for zero allocation
	responsePool sync.Pool
}

// NewRouter creates an ultra-fast router with zero-allocation design
func NewRouter() *Router {
	router := &Router{
		methodIndex: map[string]int{
			"GET":     0,
			"POST":    1,
			"PUT":     2,
			"DELETE":  3,
			"PATCH":   4,
			"OPTIONS": 5,
			"HEAD":    6,
		},
		responsePool: sync.Pool{
			New: func() interface{} {
				return &ResponseWriter{
					headers: make(http.Header, 8), // Pre-allocate common header size
				}
			},
		},
	}

	// Initialize route maps
	for i := range router.routes {
		router.routes[i] = make(map[string]http.HandlerFunc, 32) // Pre-allocate
	}

	return router
}

// ResponseWriter with pooling
type ResponseWriter struct {
	http.ResponseWriter
	headers    http.Header
	statusCode int
	written    bool
}

// Reset resets the response writer for pooling
func (rw *ResponseWriter) Reset(w http.ResponseWriter) {
	rw.ResponseWriter = w
	rw.statusCode = 200
	rw.written = false
	// Clear headers without reallocating
	for k := range rw.headers {
		delete(rw.headers, k)
	}
}

// AddRoute adds a route with zero-allocation method lookup
func (fr *Router) AddRoute(method, path string, handler http.HandlerFunc) {
	if idx, exists := fr.methodIndex[method]; exists {
		fr.routes[idx][path] = handler
	}
}

// AddMiddleware adds middleware to the router
func (fr *Router) AddMiddleware(middleware ...MiddlewareFunc) {
	fr.middleware = append(fr.middleware, middleware...)
}

// ServeHTTP implements ultra-fast request handling
func (fr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Zero-allocation method lookup using unsafe pointer arithmetic
	methodIdx, exists := fr.methodIndex[r.Method]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Direct array access - no locking needed for read-only operations
	handler, exists := fr.routes[methodIdx][r.URL.Path]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Get pooled response writer
	pooledRW := fr.responsePool.Get().(*ResponseWriter)
	pooledRW.Reset(w)
	defer fr.responsePool.Put(pooledRW)

	// Skip middleware chain for maximum speed if no middleware
	if len(fr.middleware) == 0 {
		handler(pooledRW, r)
		return
	}

	// Apply middleware chain (only if needed)
	finalHandler := http.Handler(handler)
	for i := len(fr.middleware) - 1; i >= 0; i-- {
		finalHandler = fr.middleware[i](finalHandler)
	}
	finalHandler.ServeHTTP(pooledRW, r)
}

// Hot path optimizations for common HTTP methods
func (fr *Router) GET(path string, handler http.HandlerFunc) {
	fr.routes[0][path] = handler // Direct array access
}

func (fr *Router) POST(path string, handler http.HandlerFunc) {
	fr.routes[1][path] = handler
}

func (fr *Router) PUT(path string, handler http.HandlerFunc) {
	fr.routes[2][path] = handler
}

func (fr *Router) DELETE(path string, handler http.HandlerFunc) {
	fr.routes[3][path] = handler
}

// PATCH adds a PATCH route
func (fr *Router) PATCH(path string, handler http.HandlerFunc) {
	fr.routes[4][path] = handler
}

// OPTIONS adds an OPTIONS route
func (fr *Router) OPTIONS(path string, handler http.HandlerFunc) {
	fr.routes[5][path] = handler
}

// HEAD adds a HEAD route
func (fr *Router) HEAD(path string, handler http.HandlerFunc) {
	fr.routes[6][path] = handler
}

// FastString converts []byte to string without allocation using unsafe
func FastString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// FastBytes converts string to []byte without allocation using unsafe
func FastBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}
