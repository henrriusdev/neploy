package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mssola/user_agent"
	"io"
	"log"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
	"net/http"
	"slices"
	"strings"
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
func LoggingMiddleware(next http.Handler, metrics *MetricsCollector) http.Handler {
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

		// Record metrics
		isError := rw.status >= 400
		metrics.RecordRequest(start, isError)

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

// VersionRoutingMiddleware enruta a la versión correcta según el header o la ruta
func VersionRoutingMiddleware(config model.GatewayConfig, appVersionRepo *repository.ApplicationVersion) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			println("Path segments:", pathSegments)
			var resolvedVersion string

			if slices.Contains(pathSegments, ".well-known") {
				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, r)
				return
			}

			if config.DefaultVersioningType == model.VersioningTypeHeader {
				resolvedVersion = r.Header.Get("X-API-Version")
				url := pathSegments[0]
				// Save the original path before modifying it
				if r.Header == nil {
					r.Header = make(http.Header)
				}
				r.Header.Set("X-Original-Path", r.URL.Path)
				r.URL.Path = fmt.Sprintf("/%s/%s/", resolvedVersion, url)
			} else {
				if len(pathSegments) > 1 && strings.HasPrefix(pathSegments[0], "v") {
					resolvedVersion = pathSegments[0]
				}
			}

			// Validar si la versión existe para la app actual
			if len(pathSegments) > 1 {
				appName := pathSegments[1] // se asume /vX/app-name
				exists, err := appVersionRepo.ExistsByName(r.Context(), appName, resolvedVersion)
				if err != nil || !exists {
					http.Error(w, "API version not found", http.StatusNotFound)
					return
				}
			}

			r.Header.Set("Resolved-Version", resolvedVersion)
			next.ServeHTTP(w, r)
		})
	}
}

func VisitorTraceMiddleware(visitorTrace *repository.VisitorTrace) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Parsear datos del user-agent
			ua := user_agent.New(r.UserAgent())
			browser, version := ua.Browser()

			// Continuar con la traza después de la respuesta
			next.ServeHTTP(w, r)
			duration := int(time.Since(start).Milliseconds())

			resolvedVersion := r.Header.Get("Resolved-Version")
			appName := ExtractAppName(r.URL.Path)
			if appName == "" {
				return
			}

			trace := model.VisitorTrace{
				ApplicationID:    appName,
				IpAddress:        r.RemoteAddr,
				Device:           ua.Platform(),
				Os:               ua.OS(),
				Browser:          fmt.Sprintf("%s v%s", browser, version),
				PageVisited:      r.URL.Path,
				VisitDuration:    duration,
				VisitedTimestamp: model.NewDateNow(),
			}

			go func() {
				visitorTrace.Create(context.Background(), trace, resolvedVersion)
			}()
		})
	}
}

func ExtractAppName(path string) string {
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")
	var appName string
	if len(pathSegments) > 1 && strings.HasPrefix(pathSegments[0], "v") {
		appName = pathSegments[1]
	} else if len(pathSegments) > 0 {
		appName = pathSegments[0]
	}

	return appName
}
