package server

import (
	"net/http"
	"os"
	"path/filepath"
)

// ServeWeb serves the web UI from filesystem
func (s *Server) ServeWeb(webRoot string) http.Handler {
	return http.FileServer(http.Dir(webRoot))
}

// EnableWebUI adds web UI routes to the server
func (s *Server) EnableWebUI(mux *http.ServeMux, webRoot string) error {
	// Check if web directory exists
	if _, err := os.Stat(webRoot); os.IsNotExist(err) {
		return err
	}

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(webRoot, "static")))))

	// Serve index.html at root only
	// Note: "/" pattern in Go's ServeMux matches all paths, so we need to be specific
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only serve index.html for the exact root path
		// API routes (/upload, /download, /list) are already registered and take precedence
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webRoot, "index.html"))
			return
		}
		// For all other paths that aren't handled by other routes, return 404
		http.NotFound(w, r)
	})

	return nil
}
