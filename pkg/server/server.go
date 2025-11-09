package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/0xRepo-Source/goflux/pkg/chunk"
	"github.com/0xRepo-Source/goflux/pkg/storage"
	"github.com/0xRepo-Source/goflux/pkg/transport"
)

// Server is a goflux server instance.
type Server struct {
	storage storage.Storage
	chunks  map[string][]chunk.Chunk // path -> chunks being assembled
	mu      sync.Mutex
}

// New creates a new Server.
func New(store storage.Storage) *Server {
	return &Server{
		storage: store,
		chunks:  make(map[string][]chunk.Chunk),
	}
}

// Start starts the HTTP server.
func (s *Server) Start(addr string, webRoot string) error {
	// Create a new ServeMux to avoid conflicts with default mux
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", s.handleUpload)
	mux.HandleFunc("/download", s.handleDownload)
	mux.HandleFunc("/list", s.handleList)

	// Enable web UI if webRoot provided
	if webRoot != "" {
		if err := s.EnableWebUI(mux, webRoot); err != nil {
			fmt.Printf("Warning: Could not enable web UI: %v\n", err)
		} else {
			fmt.Printf("Web UI enabled at http://%s\n", addr)
		}
	}

	fmt.Printf("goflux server listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var chunkData transport.ChunkData
	if err := json.Unmarshal(body, &chunkData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the chunk
	if s.chunks[chunkData.Path] == nil {
		s.chunks[chunkData.Path] = make([]chunk.Chunk, chunkData.Total)
	}

	s.chunks[chunkData.Path][chunkData.ChunkID] = chunk.Chunk{
		ID:       chunkData.ChunkID,
		Data:     chunkData.Data,
		Checksum: chunkData.Checksum,
	}

	// Check if all chunks received
	allReceived := true
	for _, c := range s.chunks[chunkData.Path] {
		if c.Data == nil {
			allReceived = false
			break
		}
	}

	if allReceived {
		// Reassemble and save
		chunker := chunk.New(0)
		data, err := chunker.Reassemble(s.chunks[chunkData.Path])
		if err != nil {
			http.Error(w, fmt.Sprintf("reassembly failed: %v", err), http.StatusInternalServerError)
			return
		}

		if err := s.storage.Put(chunkData.Path, data); err != nil {
			http.Error(w, fmt.Sprintf("storage failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Clean up chunks
		delete(s.chunks, chunkData.Path)
		fmt.Printf("File saved: %s (%d bytes)\n", chunkData.Path, len(data))
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "chunk %d/%d received", chunkData.ChunkID+1, chunkData.Total)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}

	data, err := s.storage.Get(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(data); err != nil {
		http.Error(w, fmt.Sprintf("write failed: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	files, err := s.storage.List(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, fmt.Sprintf("encode failed: %v", err), http.StatusInternalServerError)
		return
	}
}
