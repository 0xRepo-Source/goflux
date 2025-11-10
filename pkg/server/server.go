package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/0xRepo-Source/goflux/pkg/auth"
	"github.com/0xRepo-Source/goflux/pkg/chunk"
	"github.com/0xRepo-Source/goflux/pkg/resume"
	"github.com/0xRepo-Source/goflux/pkg/storage"
	"github.com/0xRepo-Source/goflux/pkg/transport"
)

// Server is a goflux server instance.
type Server struct {
	storage      storage.Storage
	chunks       map[string][]chunk.Chunk // path -> chunks being assembled
	sessionStore *resume.SessionStore     // tracks upload sessions for resume
	mu           sync.Mutex
	authMiddle   *auth.Middleware // nil if auth disabled
}

// New creates a new Server.
func New(store storage.Storage, metaDir string) (*Server, error) {
	sessionStore, err := resume.NewSessionStore(metaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create session store: %w", err)
	}

	return &Server{
		storage:      store,
		chunks:       make(map[string][]chunk.Chunk),
		sessionStore: sessionStore,
	}, nil
}

// EnableAuth enables authentication on the server
func (s *Server) EnableAuth(tokenStore *auth.TokenStore) {
	s.authMiddle = auth.NewMiddleware(tokenStore)
}

// Start starts the HTTP server.
func (s *Server) Start(addr string, webRoot string) error {
	// Create a new ServeMux to avoid conflicts with default mux
	mux := http.NewServeMux()

	// Register handlers with authentication if enabled
	if s.authMiddle != nil {
		mux.HandleFunc("/upload", s.authMiddle.RequireAuth("upload", s.handleUpload))
		mux.HandleFunc("/upload/status", s.authMiddle.RequireAuth("upload", s.handleUploadStatus))
		mux.HandleFunc("/download", s.authMiddle.RequireAuth("download", s.handleDownload))
		mux.HandleFunc("/list", s.authMiddle.RequireAuth("list", s.handleList))
		fmt.Println("Authentication enabled")
	} else {
		mux.HandleFunc("/upload", s.handleUpload)
		mux.HandleFunc("/upload/status", s.handleUploadStatus)
		mux.HandleFunc("/download", s.handleDownload)
		mux.HandleFunc("/list", s.handleList)
		fmt.Println("⚠️  Authentication disabled - all endpoints are public!")
	}

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

	// Get or create upload session
	session, err := s.sessionStore.GetOrCreateSession(chunkData.Path, chunkData.Total, len(chunkData.Data))
	if err != nil {
		http.Error(w, fmt.Sprintf("session error: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the chunk
	if s.chunks[chunkData.Path] == nil {
		s.chunks[chunkData.Path] = make([]chunk.Chunk, chunkData.Total)
	}

	s.chunks[chunkData.Path][chunkData.ChunkID] = chunk.Chunk{
		ID:       chunkData.ChunkID,
		Data:     chunkData.Data,
		Checksum: chunkData.Checksum,
	}

	// Mark chunk as received in session
	if err := s.sessionStore.MarkChunkReceived(chunkData.Path, chunkData.ChunkID); err != nil {
		http.Error(w, fmt.Sprintf("failed to mark chunk: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if all chunks received
	allReceived := session.Completed || true // recheck manually
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

		// Clean up chunks and session
		delete(s.chunks, chunkData.Path)
		s.sessionStore.DeleteSession(chunkData.Path)
		fmt.Printf("File saved: %s (%d bytes)\n", chunkData.Path, len(data))
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "chunk %d/%d received", chunkData.ChunkID+1, chunkData.Total)
}

// UploadStatusResponse contains the status of an upload session
type UploadStatusResponse struct {
	Exists        bool   `json:"exists"`         // whether a session exists
	TotalChunks   int    `json:"total_chunks"`   // total chunks expected
	ReceivedMap   []bool `json:"received_map"`   // bitmap of received chunks
	MissingChunks []int  `json:"missing_chunks"` // list of missing chunk IDs
	Completed     bool   `json:"completed"`      // upload completed
}

func (s *Server) handleUploadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}

	session, exists := s.sessionStore.GetSession(path)

	response := UploadStatusResponse{
		Exists: exists,
	}

	if exists {
		missing, err := s.sessionStore.GetMissingChunks(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get missing chunks: %v", err), http.StatusInternalServerError)
			return
		}

		response.TotalChunks = session.TotalChunks
		response.ReceivedMap = session.ReceivedMap
		response.MissingChunks = missing
		response.Completed = session.Completed
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("encode failed: %v", err), http.StatusInternalServerError)
		return
	}
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
