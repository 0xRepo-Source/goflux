package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Token represents an authentication token
type Token struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	User        string    `json:"user"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Revoked     bool      `json:"revoked"`
}

// TokenStore holds all tokens with thread-safe access
type TokenStore struct {
	mu       sync.RWMutex
	tokens   map[string]*Token // key is token hash
	filename string
}

// TokenStoreFile represents the JSON file format
type TokenStoreFile struct {
	Tokens []Token `json:"tokens"`
}

// NewTokenStore creates a new token store
func NewTokenStore(filename string) (*TokenStore, error) {
	ts := &TokenStore{
		tokens:   make(map[string]*Token),
		filename: filename,
	}

	if err := ts.Load(); err != nil {
		return nil, err
	}

	return ts, nil
}

// Load reads tokens from file
func (ts *TokenStore) Load() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	data, err := os.ReadFile(ts.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return fmt.Errorf("error reading token file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	var storeFile TokenStoreFile
	if err := json.Unmarshal(data, &storeFile); err != nil {
		return fmt.Errorf("error parsing token file: %w", err)
	}

	// Build token map
	ts.tokens = make(map[string]*Token)
	for i := range storeFile.Tokens {
		token := &storeFile.Tokens[i]
		ts.tokens[token.TokenHash] = token
	}

	return nil
}

// Reload reloads tokens from file
func (ts *TokenStore) Reload() error {
	return ts.Load()
}

// Validate checks if a token is valid and returns the associated user and permissions
func (ts *TokenStore) Validate(tokenStr string) (string, []string, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	token, exists := ts.tokens[tokenHash]
	if !exists {
		return "", nil, fmt.Errorf("invalid token")
	}

	if token.Revoked {
		return "", nil, fmt.Errorf("token has been revoked")
	}

	if time.Now().After(token.ExpiresAt) {
		return "", nil, fmt.Errorf("token has expired")
	}

	return token.User, token.Permissions, nil
}

// HasPermission checks if a user has a specific permission
func HasPermission(permissions []string, required string) bool {
	for _, perm := range permissions {
		if perm == required || perm == "*" {
			return true
		}
	}
	return false
}
