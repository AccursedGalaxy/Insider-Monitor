package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Server is the frontend server that serves static files and proxies API requests
type Server struct {
	StaticDir   string
	APIEndpoint string
	WSEndpoint  string
}

// NewServer creates a new frontend server
func NewServer(staticDir, apiEndpoint, wsEndpoint string) *Server {
	return &Server{
		StaticDir:   staticDir,
		APIEndpoint: apiEndpoint,
		WSEndpoint:  wsEndpoint,
	}
}

// Handler returns an HTTP handler for the frontend server
func (s *Server) Handler() http.Handler {
	// Create file server for static files
	fileServer := http.FileServer(http.Dir(s.StaticDir))

	// Create API reverse proxy
	apiURL, err := url.Parse(s.APIEndpoint)
	if err != nil {
		panic(err)
	}
	apiProxy := httputil.NewSingleHostReverseProxy(apiURL)

	// Create WebSocket reverse proxy
	wsURL, err := url.Parse(s.WSEndpoint)
	if err != nil {
		panic(err)
	}
	wsProxy := httputil.NewSingleHostReverseProxy(wsURL)

	// Configure WebSocket proxy with special headers
	wsProxy.ModifyResponse = func(resp *http.Response) error {
		// Ensure WebSocket upgrade headers are preserved
		if resp.StatusCode == http.StatusSwitchingProtocols {
			resp.Header.Set("Connection", "Upgrade")
			resp.Header.Set("Upgrade", "websocket")
		}
		return nil
	}

	// Create main handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proxy API requests
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiProxy.ServeHTTP(w, r)
			return
		}

		// Proxy WebSocket requests
		if strings.HasPrefix(r.URL.Path, "/ws") {
			// Set necessary headers for WebSocket proxying
			if r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "websocket" {
				wsProxy.ServeHTTP(w, r)
				return
			}
		}

		// Check if file exists
		path := filepath.Join(s.StaticDir, r.URL.Path)
		_, err := os.Stat(path)

		// If path doesn't exist or is a directory, serve index.html for SPA routing
		if os.IsNotExist(err) || (err == nil && strings.HasSuffix(path, "/")) {
			// Check if it's a standard extension that should 404 if missing
			ext := filepath.Ext(path)
			if ext != "" && ext != ".html" {
				// 404 for missing assets like .js, .css, .png, etc.
				http.NotFound(w, r)
				return
			}

			// Otherwise serve index.html for client-side routing
			r.URL.Path = "/"
		}

		// Serve static files
		fileServer.ServeHTTP(w, r)
	})
}
