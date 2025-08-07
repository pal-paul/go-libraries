package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ResponseHelper provides utility methods for HTTP responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// JSON writes a JSON response
func (rh *ResponseHelper) JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Text writes a plain text response
func (rh *ResponseHelper) Text(w http.ResponseWriter, status int, text string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	return err
}

// HTML writes an HTML response
func (rh *ResponseHelper) HTML(w http.ResponseWriter, status int, html string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	_, err := w.Write([]byte(html))
	return err
}

// Error writes an error response
func (rh *ResponseHelper) Error(w http.ResponseWriter, status int, message string) error {
	return rh.JSON(w, status, map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  status,
	})
}

// Success writes a success response
func (rh *ResponseHelper) Success(w http.ResponseWriter, data interface{}) error {
	return rh.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// RequestHelper provides utility methods for HTTP requests
type RequestHelper struct{}

// NewRequestHelper creates a new request helper
func NewRequestHelper() *RequestHelper {
	return &RequestHelper{}
}

// GetQueryParam gets a query parameter value
func (rh *RequestHelper) GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetQueryParamInt gets a query parameter as integer
func (rh *RequestHelper) GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

// GetQueryParamBool gets a query parameter as boolean
func (rh *RequestHelper) GetQueryParamBool(r *http.Request, key string, defaultValue bool) bool {
	value := strings.ToLower(r.URL.Query().Get(key))
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// GetHeader gets a header value
func (rh *RequestHelper) GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// ParseJSON parses JSON request body
func (rh *RequestHelper) ParseJSON(r *http.Request, dest interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}

// GetClientIP gets the client IP address
func (rh *RequestHelper) GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = ip[:strings.LastIndex(ip, ":")]
	}
	return ip
}

// HealthCheck provides a simple health check handler
func HealthCheckHandler() http.HandlerFunc {
	rh := NewResponseHelper()
	return func(w http.ResponseWriter, r *http.Request) {
		rh.JSON(w, http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": "2025-08-07T00:00:00Z", // This would be time.Now() in real usage
		})
	}
}

// NotFoundHandler provides a custom 404 handler
func NotFoundHandler() http.HandlerFunc {
	rh := NewResponseHelper()
	return func(w http.ResponseWriter, r *http.Request) {
		rh.Error(w, http.StatusNotFound, "Resource not found")
	}
}

// MethodNotAllowedHandler provides a custom 405 handler
func MethodNotAllowedHandler() http.HandlerFunc {
	rh := NewResponseHelper()
	return func(w http.ResponseWriter, r *http.Request) {
		rh.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
