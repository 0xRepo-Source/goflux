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
func (s *Server) EnableWebUI(webRoot string) error {
	// Check if web directory exists
	if _, err := os.Stat(webRoot); os.IsNotExist(err) {
		return err
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(webRoot, "static")))))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(webRoot, "index.html"))
		} else {
			// Let other handlers handle their routes
			http.DefaultServeMux.ServeHTTP(w, r)
		}
	})

	return nil
}
