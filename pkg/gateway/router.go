package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"neploy.dev/pkg/repository"
)

type Route struct {
	AppID  string
	Port   string
	Domain string
	Path   string
}

type Router struct {
	routes            map[string]*httputil.ReverseProxy
	routeInfo         map[string]Route
	mu                sync.RWMutex
	metrics           map[string]*MetricsCollector
	metricsAggregator *MetricsAggregator
	version           *repository.ApplicationVersion
	conf              *repository.GatewayConfig
	vtrace            *repository.VisitorTrace
}

func NewRouter(appStatRepo *repository.ApplicationStat, version *repository.ApplicationVersion, conf *repository.GatewayConfig, vtrace *repository.VisitorTrace) *Router {
	router := &Router{
		routes:    make(map[string]*httputil.ReverseProxy),
		routeInfo: make(map[string]Route),
		metrics:   make(map[string]*MetricsCollector),
		mu:        sync.RWMutex{},
		version:   version,
		conf:      conf,
		vtrace:    vtrace,
	}

	// Create metrics aggregator without a specific collector
	router.metricsAggregator = NewMetricsAggregator(nil, appStatRepo)
	router.metricsAggregator.Start()

	return router
}

func (r *Router) Close() {
	if r.metricsAggregator != nil {
		r.metricsAggregator.Stop()
	}
}

func (r *Router) AddRoute(route Route) error {
	if err := ValidateRoute(route); err != nil {
		return err
	}

	target, err := url.Parse(fmt.Sprintf("http://localhost:%s", route.Port))
	if err != nil {
		return fmt.Errorf("invalid target URL: %v", err)
	}

	r.mu.Lock()
	if _, exists := r.metrics[route.AppID]; !exists {
		metrics, err := NewMetricsCollector("./data/metrics", route.AppID)
		if err != nil {
			log.Printf("ERROR: Failed to create metrics collector for app %s: %v", route.AppID, err)
		} else {
			r.metrics[route.AppID] = metrics
			if r.metricsAggregator != nil {
				r.metricsAggregator.AddCollector(metrics)
				r.metricsAggregator.aggregateAndSaveMetrics()

			}
		}
	}
	r.mu.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Save the original path before any modifications
		originalPath := req.URL.Path
		log.Printf("DEBUG: Director original path: %s", originalPath)
		
		originalDirector(req)
		req.Host = req.URL.Host

		if route.Path != "" {
			versionPrefix := req.Header.Get("Resolved-Version")
			basePath := "/" + versionPrefix + route.Path

			trimmed := strings.TrimPrefix(req.URL.Path, basePath)

			// Fallback si no coincide completamente
			if trimmed == req.URL.Path {
				trimmed = strings.TrimPrefix(req.URL.Path, route.Path)
			}

			if trimmed == "" {
				trimmed = "/" // fallback para evitar path vacío
			}

			// Asegurar que comienza con /
			if !strings.HasPrefix(trimmed, "/") {
				trimmed = "/" + trimmed
			}

			req.URL.Path = trimmed
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("ERROR: Proxy error for route %s to %s: %v", route.Path, target.String(), err)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Error: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	routeKey := route.Path
	r.routes[routeKey] = proxy
	r.routeInfo[routeKey] = route

	return nil
}

func (r *Router) RemoveRoute(routeKey string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.routes, routeKey)
	delete(r.routeInfo, routeKey)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Store the original URL path at the very beginning
	originalPath := req.URL.Path
	log.Printf("DEBUG: Original request path: %s", originalPath)

	config, err := r.conf.Get(req.Context())
	if err != nil {
		log.Printf("ERROR: Failed to get gateway config: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	resolver := VersionRoutingMiddleware(config, r.version)
	resolver(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.mu.RLock()
		defer r.mu.RUnlock()
		
		// Store the original path in a context value to use for route matching
		ctx := context.WithValue(req.Context(), "originalPath", originalPath)
		req = req.WithContext(ctx)
		
		// Log all available routes for debugging
		log.Printf("DEBUG: Available routes:")
		for routeKey := range r.routes {
			route := r.routeInfo[routeKey]
			log.Printf("DEBUG: Route: path=%s, domain=%s", route.Path, route.Domain)
		}

		if strings.Contains(req.URL.Path, "//v") {
			req.URL.Path = strings.ReplaceAll(req.URL.Path, "//v", "/v")
		}

		// Get the original path from the context
		pathToMatch, _ := req.Context().Value("originalPath").(string)
		if pathToMatch == "" {
			pathToMatch = req.URL.Path
		}
		log.Printf("DEBUG: Using path for route matching: %s", pathToMatch)
		
		for routeKey, proxy := range r.routes {
			route := r.routeInfo[routeKey]
			log.Printf("DEBUG: Checking route %s against path %s", route.Path, pathToMatch)
			if r.matchesRoute(req, route) {
				var handler http.Handler = proxy

				if r.metrics != nil {
					handler = LoggingMiddleware(handler, r.metrics[route.AppID])
				} else {
					log.Printf("WARN: Metrics collector not available")
				}

				handler = CacheMiddleware(handler)
				handler = VisitorTraceMiddleware(r.vtrace)(handler)

				handler.ServeHTTP(w, req)
				return
			}
		}

		if strings.Contains(req.URL.Path, ".well-known") {
			w.WriteHeader(http.StatusContinue)
			return
		}

		origPath, _ := req.Context().Value("originalPath").(string)
		log.Printf("WARN: No matching route found for path: %s, original path: %s, host: %s", req.URL.Path, origPath, req.Host)

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found")
	})).ServeHTTP(w, req)
}

func (r *Router) matchesRoute(req *http.Request, route Route) bool {
	// Get the original path from context
	path := req.URL.Path
	if originalPath, ok := req.Context().Value("originalPath").(string); ok && originalPath != "" {
		path = originalPath
	}
	
	// Also check if this is a versioned path that needs special handling
	if strings.HasPrefix(path, "/v") && len(path) > 2 && (path[2] >= '0' && path[2] <= '9') {
		log.Printf("DEBUG: Detected versioned path: %s", path)
	}

	matches := false
	if route.Path != "" {
		// Check if the route path is a prefix of the request path
		matches = strings.HasPrefix(path, route.Path)
		
		// Special handling for versioned paths
		if !matches && strings.HasPrefix(path, "/v") {
			// Try to extract version and check if the rest matches
			parts := strings.SplitN(path, "/", 3)
			if len(parts) >= 3 {
				versionPart := parts[1] // e.g. "v1.0.0"
				restPath := "/" + parts[2]
				log.Printf("DEBUG: Checking versioned path: version=%s, rest=%s against route=%s", 
					versionPart, restPath, route.Path)
				
				// Check if the route path matches after removing version
				if strings.HasPrefix(restPath, route.Path) {
					matches = true
				}
			}
		}
	} else {
		matches = true
	}

	log.Printf("DEBUG: Route match result for path %s and route %s: %v", path, route.Path, matches)
	return matches
}

func ValidateRoute(route Route) error {
	if route.AppID == "" {
		return fmt.Errorf("appID is required")
	}
	if route.Port == "" {
		return fmt.Errorf("port is required")
	}
	if route.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	return nil
}
