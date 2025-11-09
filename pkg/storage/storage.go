package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// Storage is an interface for storing and retrieving files.
type Storage interface {
	Put(path string, data []byte) error
	Get(path string) ([]byte, error)
	Exists(path string) bool
	List(path string) ([]string, error)
}

// Local is a simple local filesystem storage implementation.
type Local struct {
	Root string
}

// NewLocal creates a new local filesystem storage backend.
func NewLocal(root string) (*Local, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create root dir: %w", err)
	}
	return &Local{Root: root}, nil
}

func (l *Local) Put(path string, data []byte) error {
	fullPath := filepath.Join(l.Root, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(fullPath, data, 0644)
}

func (l *Local) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(l.Root, path)
	return os.ReadFile(fullPath)
}

func (l *Local) Exists(path string) bool {
	fullPath := filepath.Join(l.Root, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (l *Local) List(path string) ([]string, error) {
	fullPath := filepath.Join(l.Root, path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}
