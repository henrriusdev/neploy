package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type Route struct {
	AppID     string
	Port      string
	Domain    string
	Path      string
	Subdomain string
}

type Router struct {
	routes    map[string]*httputil.ReverseProxy
	routeInfo map[string]Route
	mu        sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		routes:    make(map[string]*httputil.ReverseProxy),
		routeInfo: make(map[string]Route),
	}
}

// AddRoute adds a new route to the router
func (r *Router) AddRoute(route Route) error {
	if err := ValidateRoute(route); err != nil {
		return err
	}

	target, err := url.Parse(fmt.Sprintf("http://localhost:%s", route.Port))
	if err != nil {
		return fmt.Errorf("invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Custom director to handle path rewriting
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = req.URL.Host
		if route.Path != "" {
			// Strip the route path prefix
			req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}
	}

	// Add error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Error: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	routeKey := r.generateRouteKey(route)
	r.routes[routeKey] = proxy
	r.routeInfo[routeKey] = route

	return nil
}

// RemoveRoute removes a route from the router
func (r *Router) RemoveRoute(route Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	routeKey := r.generateRouteKey(route)
	delete(r.routes, routeKey)
	delete(r.routeInfo, routeKey)
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try to find a matching route
	for routeKey, proxy := range r.routes {
		route := r.routeInfo[routeKey]
		if r.matchesRoute(req, route) {
			proxy.ServeHTTP(w, req)
			return
		}
	}

	// No matching route found
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not Found")
}

// matchesRoute checks if the request matches a route
func (r *Router) matchesRoute(req *http.Request, route Route) bool {
	// Check subdomain
	if route.Subdomain != "" {
		hostParts := strings.Split(req.Host, ".")
		if len(hostParts) < 2 || hostParts[0] != route.Subdomain {
			return false
		}
	}

	// Check path
	if route.Path != "" {
		return strings.HasPrefix(req.URL.Path, route.Path)
	}

	return true
}

// generateRouteKey creates a unique key for a route
func (r *Router) generateRouteKey(route Route) string {
	if route.Subdomain != "" {
		return fmt.Sprintf("%s.%s", route.Subdomain, route.Domain)
	}
	return fmt.Sprintf("%s%s", route.Domain, route.Path)
}

// ValidateRoute validates a route configuration
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
	if route.Subdomain == "" && route.Path == "" {
		return fmt.Errorf("either subdomain or path must be specified")
	}
	return nil
}
