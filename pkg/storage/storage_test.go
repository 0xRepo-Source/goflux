package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLocal(t *testing.T) {
	tmpDir := t.TempDir()

	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	if store.Root != tmpDir {
		t.Errorf("NewLocal().Root = %s, want %s", store.Root, tmpDir)
	}
}

func TestLocalPutGet(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	testData := []byte("Hello, storage test!")
	testPath := "test/file.txt"

	// Put
	err = store.Put(testPath, testData)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// Get
	result, err := store.Get(testPath)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if string(result) != string(testData) {
		t.Errorf("Get() = %s, want %s", result, testData)
	}
}

func TestLocalExists(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	testPath := "exists/test.txt"

	// Should not exist initially
	if store.Exists(testPath) {
		t.Error("Exists() = true for non-existent file, want false")
	}

	// Create file
	err = store.Put(testPath, []byte("test"))
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// Should exist now
	if !store.Exists(testPath) {
		t.Error("Exists() = false for existing file, want true")
	}
}

func TestLocalList(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	// Create test files
	files := map[string][]byte{
		"dir/file1.txt": []byte("content1"),
		"dir/file2.txt": []byte("content2"),
		"dir/file3.txt": []byte("content3"),
	}

	for path, data := range files {
		err := store.Put(path, data)
		if err != nil {
			t.Fatalf("Put(%s) error = %v", path, err)
		}
	}

	// List directory
	list, err := store.List("dir")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(list) != 3 {
		t.Errorf("List() returned %d files, want 3", len(list))
	}
}

func TestLocalPutCreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	deepPath := "a/b/c/d/file.txt"
	testData := []byte("nested file")

	err = store.Put(deepPath, testData)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// Verify file exists
	fullPath := filepath.Join(tmpDir, deepPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Put() did not create nested directories")
	}

	// Verify content
	result, err := store.Get(deepPath)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if string(result) != string(testData) {
		t.Error("Get() returned incorrect data for nested file")
	}
}

func TestLocalGetNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewLocal(tmpDir)
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}

	_, err = store.Get("nonexistent.txt")
	if err == nil {
		t.Error("Get() for non-existent file should return error")
	}
}
