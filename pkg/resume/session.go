package resume

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UploadSession tracks the state of a partial upload
type UploadSession struct {
	Path         string    `json:"path"`          // destination path
	TotalChunks  int       `json:"total_chunks"`  // expected number of chunks
	ChunkSize    int       `json:"chunk_size"`    // size of each chunk
	FileHash     string    `json:"file_hash"`     // SHA-256 of complete file (optional)
	ReceivedMap  []bool    `json:"received_map"`  // bitmap of received chunks
	CreatedAt    time.Time `json:"created_at"`    // when upload started
	LastModified time.Time `json:"last_modified"` // last chunk received
	Completed    bool      `json:"completed"`     // upload completed
}

// SessionStore manages upload sessions with persistence
type SessionStore struct {
	sessions map[string]*UploadSession // keyed by upload ID (hash of path)
	metaDir  string                    // directory for metadata files
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore(metaDir string) (*SessionStore, error) {
	// Create metadata directory if it doesn't exist
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metadata directory: %w", err)
	}

	store := &SessionStore{
		sessions: make(map[string]*UploadSession),
		metaDir:  metaDir,
	}

	// Load existing sessions
	if err := store.loadSessions(); err != nil {
		return nil, fmt.Errorf("failed to load sessions: %w", err)
	}

	return store, nil
}

// GetOrCreateSession gets an existing session or creates a new one
func (s *SessionStore) GetOrCreateSession(path string, totalChunks, chunkSize int) (*UploadSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := s.makeSessionID(path)

	// Check if session exists
	if session, exists := s.sessions[sessionID]; exists {
		// Validate session matches request
		if session.TotalChunks != totalChunks {
			return nil, fmt.Errorf("chunk count mismatch: session has %d, request has %d", session.TotalChunks, totalChunks)
		}
		return session, nil
	}

	// Create new session
	session := &UploadSession{
		Path:         path,
		TotalChunks:  totalChunks,
		ChunkSize:    chunkSize,
		ReceivedMap:  make([]bool, totalChunks),
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
		Completed:    false,
	}

	s.sessions[sessionID] = session

	// Persist to disk
	if err := s.saveSession(sessionID, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

// MarkChunkReceived marks a chunk as received
func (s *SessionStore) MarkChunkReceived(path string, chunkID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := s.makeSessionID(path)
	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found for path: %s", path)
	}

	if chunkID < 0 || chunkID >= session.TotalChunks {
		return fmt.Errorf("invalid chunk ID: %d (total: %d)", chunkID, session.TotalChunks)
	}

	session.ReceivedMap[chunkID] = true
	session.LastModified = time.Now()

	// Check if all chunks received
	allReceived := true
	for _, received := range session.ReceivedMap {
		if !received {
			allReceived = false
			break
		}
	}
	session.Completed = allReceived

	// Persist to disk
	return s.saveSession(sessionID, session)
}

// GetSession retrieves a session by path
func (s *SessionStore) GetSession(path string) (*UploadSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionID := s.makeSessionID(path)
	session, exists := s.sessions[sessionID]
	return session, exists
}

// DeleteSession removes a completed session
func (s *SessionStore) DeleteSession(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := s.makeSessionID(path)
	delete(s.sessions, sessionID)

	// Delete metadata file
	metaFile := filepath.Join(s.metaDir, sessionID+".json")
	if err := os.Remove(metaFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete session file: %w", err)
	}

	return nil
}

// GetMissingChunks returns a list of chunk IDs that haven't been received
func (s *SessionStore) GetMissingChunks(path string) ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionID := s.makeSessionID(path)
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found for path: %s", path)
	}

	missing := []int{}
	for i, received := range session.ReceivedMap {
		if !received {
			missing = append(missing, i)
		}
	}

	return missing, nil
}

// CleanupOldSessions removes sessions older than the specified duration
func (s *SessionStore) CleanupOldSessions(maxAge time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	toDelete := []string{}

	for sessionID, session := range s.sessions {
		if session.LastModified.Before(cutoff) && !session.Completed {
			toDelete = append(toDelete, sessionID)
		}
	}

	for _, sessionID := range toDelete {
		delete(s.sessions, sessionID)

		metaFile := filepath.Join(s.metaDir, sessionID+".json")
		if err := os.Remove(metaFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete session file: %w", err)
		}
	}

	if len(toDelete) > 0 {
		fmt.Printf("Cleaned up %d old sessions\n", len(toDelete))
	}

	return nil
}

// makeSessionID creates a unique session ID from the path
func (s *SessionStore) makeSessionID(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])[:16] // Use first 16 chars
}

// saveSession persists a session to disk
func (s *SessionStore) saveSession(sessionID string, session *UploadSession) error {
	metaFile := filepath.Join(s.metaDir, sessionID+".json")

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metaFile, data, 0644)
}

// loadSessions loads all sessions from disk
func (s *SessionStore) loadSessions() error {
	files, err := os.ReadDir(s.metaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // metadata directory doesn't exist yet
		}
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		sessionID := file.Name()[:len(file.Name())-5] // remove .json
		metaFile := filepath.Join(s.metaDir, file.Name())

		data, err := os.ReadFile(metaFile)
		if err != nil {
			fmt.Printf("Warning: failed to read session file %s: %v\n", metaFile, err)
			continue
		}

		var session UploadSession
		if err := json.Unmarshal(data, &session); err != nil {
			fmt.Printf("Warning: failed to parse session file %s: %v\n", metaFile, err)
			continue
		}

		s.sessions[sessionID] = &session
	}

	if len(s.sessions) > 0 {
		fmt.Printf("Loaded %d existing upload sessions\n", len(s.sessions))
	}

	return nil
}
