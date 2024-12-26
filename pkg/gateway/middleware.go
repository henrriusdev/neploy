package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status    int
	size      int64
	buf       *bytes.Buffer
	tee       io.Writer
	committed bool
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.tee.Write(b)
	w.size += int64(n)
	return n, err
}

func (w *responseWriter) WriteHeader(status int) {
	if w.committed {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
	w.committed = true
}

// LoggingMiddleware wraps an http.Handler and logs request/response details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a buffer to store the response
		buf := &bytes.Buffer{}

		// Create a response writer that captures the response
		rw := &responseWriter{
			ResponseWriter: w,
			buf:            buf,
			tee:            io.MultiWriter(w, buf),
			status:         http.StatusOK, // Default status
		}

		// Read and store the request body
		var reqBody []byte
		if r.Body != nil {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Format request headers
		headers := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		// Create log entry
		logEntry := struct {
			Timestamp  string            `json:"timestamp"`
			Method     string            `json:"method"`
			Path       string            `json:"path"`
			Query      string            `json:"query"`
			Status     int               `json:"status"`
			Size       int64             `json:"size"`
			Duration   string            `json:"duration"`
			RemoteAddr string            `json:"remote_addr"`
			UserAgent  string            `json:"user_agent"`
			Headers    map[string]string `json:"headers"`
			ReqBody    string            `json:"request_body,omitempty"`
			RespBody   string            `json:"response_body,omitempty"`
		}{
			Timestamp:  time.Now().Format(time.RFC3339),
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			Status:     rw.status,
			Size:       rw.size,
			Duration:   duration.String(),
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Headers:    headers,
		}

		// Add request body if present and not too large
		if len(reqBody) > 0 && len(reqBody) < 10000 { // Limit to 10KB
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, reqBody, "", "  "); err == nil {
				logEntry.ReqBody = prettyJSON.String()
			} else {
				logEntry.ReqBody = string(reqBody)
			}
		}

		// Add response body if present and not too large
		if rw.buf.Len() > 0 && rw.buf.Len() < 10000 { // Limit to 10KB
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, rw.buf.Bytes(), "", "  "); err == nil {
				logEntry.RespBody = prettyJSON.String()
			} else {
				logEntry.RespBody = rw.buf.String()
			}
		}

		// Convert log entry to JSON
		logJSON, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Error marshaling log entry: %v", err)
			return
		}

		// Log based on status code
		switch {
		case rw.status >= 500:
			log.Printf("ERROR: %s", string(logJSON))
		case rw.status >= 400:
			log.Printf("WARN: %s", string(logJSON))
		default:
			log.Printf("INFO: %s", string(logJSON))
		}
	})
}

// RateLimitMiddleware implements rate limiting for gateway routes
func RateLimitMiddleware(next http.Handler, rateLimit int) http.Handler {
	// TODO: Implement rate limiting using a token bucket or similar algorithm
	return next
}

// AuthMiddleware implements authentication checking for gateway routes
func AuthMiddleware(next http.Handler, requiresAuth bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requiresAuth {
			next.ServeHTTP(w, r)
			return
		}

		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		// TODO: Implement token validation logic here
		// This should validate the token format, expiration, and signature

		next.ServeHTTP(w, r)
	})
}
