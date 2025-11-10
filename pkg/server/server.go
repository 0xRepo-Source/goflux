package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/0xRepo-Source/goflux/pkg/auth"
	"github.com/0xRepo-Source/goflux/pkg/resume"
	"github.com/0xRepo-Source/goflux/pkg/storage"
	"github.com/0xRepo-Source/goflux/pkg/transport"
)

// Server is a goflux server instance.
type Server struct {
	storage      storage.Storage
	chunksDir    string               // directory for temporary chunk storage
	sessionStore *resume.SessionStore // tracks upload sessions for resume
	mu           sync.Mutex
	authMiddle   *auth.Middleware // nil if auth disabled
}

// New creates a new Server.
func New(store storage.Storage, metaDir string) (*Server, error) {
	sessionStore, err := resume.NewSessionStore(metaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create session store: %w", err)
	}

	// Create chunks directory for temporary storage
	chunksDir := filepath.Join(metaDir, "chunks")
	if err := os.MkdirAll(chunksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chunks directory: %w", err)
	}

	return &Server{
		storage:      store,
		chunksDir:    chunksDir,
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

	// Create session-specific chunks directory using path hash
	sessionHash := fmt.Sprintf("%x", []byte(chunkData.Path))
	sessionChunksDir := filepath.Join(s.chunksDir, sessionHash[:16])
	if err := os.MkdirAll(sessionChunksDir, 0755); err != nil {
		http.Error(w, fmt.Sprintf("failed to create session chunks dir: %v", err), http.StatusInternalServerError)
		return
	}

	// Write chunk to disk
	chunkPath := filepath.Join(sessionChunksDir, fmt.Sprintf("chunk_%06d.dat", chunkData.ChunkID))
	if err := os.WriteFile(chunkPath, chunkData.Data, 0644); err != nil {
		http.Error(w, fmt.Sprintf("failed to write chunk: %v", err), http.StatusInternalServerError)
		return
	}

	// Mark chunk as received in session
	if err := s.sessionStore.MarkChunkReceived(chunkData.Path, chunkData.ChunkID); err != nil {
		http.Error(w, fmt.Sprintf("failed to mark chunk: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if upload is complete
	if session.Completed {
		// Reassemble file from disk chunks
		if err := s.reassembleFromDisk(sessionChunksDir, chunkData.Path, chunkData.Total); err != nil {
			http.Error(w, fmt.Sprintf("reassembly failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Clean up chunks directory and session
		os.RemoveAll(sessionChunksDir)
		if err := s.sessionStore.DeleteSession(chunkData.Path); err != nil {
			fmt.Printf("Warning: failed to delete session metadata: %v\n", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "chunk %d/%d received", chunkData.ChunkID+1, chunkData.Total)
}

// reassembleFromDisk reads chunks from disk and assembles the final file
func (s *Server) reassembleFromDisk(chunksDir, remotePath string, totalChunks int) error {
	// Open output file for writing
	tempPath := filepath.Join(s.chunksDir, "temp_"+filepath.Base(remotePath))
	outFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer outFile.Close()

	// Read and write each chunk in order
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(chunksDir, fmt.Sprintf("chunk_%06d.dat", i))
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to read chunk %d: %w", i, err)
		}

		if _, err := outFile.Write(chunkData); err != nil {
			return fmt.Errorf("failed to write chunk %d: %w", i, err)
		}
	}

	outFile.Close()

	// Read the assembled file and put into storage
	finalData, err := os.ReadFile(tempPath)
	if err != nil {
		return fmt.Errorf("failed to read assembled file: %w", err)
	}

	if err := s.storage.Put(remotePath, finalData); err != nil {
		return fmt.Errorf("storage failed: %w", err)
	}

	// Clean up temp file
	os.Remove(tempPath)

	fmt.Printf("File saved: %s (%d bytes)\n", remotePath, len(finalData))
	return nil
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
