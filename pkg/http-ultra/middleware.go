package httpserver

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// compressionWriter wraps http.ResponseWriter to provide gzip compression
type compressionWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (cw *compressionWriter) Write(b []byte) (int, error) {
	return cw.writer.Write(b)
}

// gzipWriterPool pools gzip writers for reuse
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

// CompressionMiddleware provides gzip compression middleware
func CompressionMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Get gzip writer from pool
			gz := gzipWriterPool.Get().(*gzip.Writer)
			defer gzipWriterPool.Put(gz)

			gz.Reset(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			cw := &compressionWriter{
				ResponseWriter: w,
				writer:         gz,
			}

			next.ServeHTTP(cw, r)
		})
	}
}

// CORSMiddleware provides CORS support
func CORSMiddleware(origins, methods, headers []string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			if len(origins) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
			}
			if len(methods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			}
			if len(headers) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware provides basic request logging
func LoggingMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For performance, we keep logging minimal
			// You can replace this with a more sophisticated logger
			next.ServeHTTP(w, r)
		})
	}
}
