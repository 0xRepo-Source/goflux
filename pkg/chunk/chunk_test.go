package chunk

import (
	"bytes"
	"testing"
)

func TestChunkerNew(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		expected int
	}{
		{"default size", 0, 1024 * 1024},
		{"custom size", 2048, 2048},
		{"negative size", -100, 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunker := New(tt.size)
			if chunker.Size != tt.expected {
				t.Errorf("New(%d) = %d, want %d", tt.size, chunker.Size, tt.expected)
			}
		})
	}
}

func TestChunkerSplit(t *testing.T) {
	chunker := New(10)
	data := []byte("Hello, World! This is a test.")

	chunks := chunker.Split(data)

	if len(chunks) != 3 {
		t.Errorf("Split() produced %d chunks, want 3", len(chunks))
	}

	// Verify chunk IDs
	for i, chunk := range chunks {
		if chunk.ID != i {
			t.Errorf("Chunk %d has ID %d, want %d", i, chunk.ID, i)
		}
	}

	// Verify checksums exist
	for i, chunk := range chunks {
		if chunk.Checksum == "" {
			t.Errorf("Chunk %d has empty checksum", i)
		}
	}
}

func TestChunkerReassemble(t *testing.T) {
	chunker := New(10)
	original := []byte("Hello, World! This is a test for reassembly.")

	// Split
	chunks := chunker.Split(original)

	// Reassemble
	result, err := chunker.Reassemble(chunks)
	if err != nil {
		t.Fatalf("Reassemble() error = %v", err)
	}

	// Compare
	if !bytes.Equal(result, original) {
		t.Errorf("Reassemble() = %v, want %v", result, original)
	}
}

func TestChunkerReassembleChecksumMismatch(t *testing.T) {
	chunker := New(10)
	data := []byte("Test data for checksum verification")

	chunks := chunker.Split(data)

	// Corrupt a chunk
	if len(chunks) > 0 {
		chunks[0].Data[0] ^= 0xFF // Flip bits
	}

	// Should fail reassembly
	_, err := chunker.Reassemble(chunks)
	if err == nil {
		t.Error("Reassemble() expected error for corrupted chunk, got nil")
	}
}

func TestChunkerReassembleOutOfOrder(t *testing.T) {
	chunker := New(10)
	data := []byte("Test data for order verification")

	chunks := chunker.Split(data)

	// Corrupt chunk ID
	if len(chunks) > 1 {
		chunks[1].ID = 99
	}

	// Should fail reassembly
	_, err := chunker.Reassemble(chunks)
	if err == nil {
		t.Error("Reassemble() expected error for out-of-order chunks, got nil")
	}
}

func TestChunkerLargeData(t *testing.T) {
	chunker := New(1024)
	data := make([]byte, 10*1024) // 10KB
	for i := range data {
		data[i] = byte(i % 256)
	}

	chunks := chunker.Split(data)

	if len(chunks) != 10 {
		t.Errorf("Split() produced %d chunks, want 10", len(chunks))
	}

	result, err := chunker.Reassemble(chunks)
	if err != nil {
		t.Fatalf("Reassemble() error = %v", err)
	}

	if !bytes.Equal(result, data) {
		t.Error("Reassemble() produced different data")
	}
}
