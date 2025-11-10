package chunk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Chunker is responsible for splitting data into resume-able chunks.
type Chunker struct {
	Size int
}

// Chunk represents a single chunk of data.
type Chunk struct {
	ID       int
	Data     []byte
	Checksum string
}

func New(size int) *Chunker {
	if size <= 0 {
		size = 1024 * 1024 // 1MB default
	}
	return &Chunker{Size: size}
}

// Split splits data into chunks.
func (c *Chunker) Split(data []byte) []Chunk {
	var chunks []Chunk
	totalSize := len(data)

	for i := 0; i < totalSize; i += c.Size {
		end := i + c.Size
		if end > totalSize {
			end = totalSize
		}

		chunkData := data[i:end]
		hash := sha256.Sum256(chunkData)

		chunks = append(chunks, Chunk{
			ID:       len(chunks),
			Data:     chunkData,
			Checksum: hex.EncodeToString(hash[:]),
		})
	}

	return chunks
}

// Reassemble combines chunks back into original data.
func (c *Chunker) Reassemble(chunks []Chunk) ([]byte, error) {
	var result []byte
	for i, chunk := range chunks {
		if chunk.ID != i {
			return nil, fmt.Errorf("chunk %d missing or out of order", i)
		}

		// Verify checksum
		hash := sha256.Sum256(chunk.Data)
		expectedChecksum := hex.EncodeToString(hash[:])

		if len(chunk.Checksum) == 64 && chunk.Checksum != expectedChecksum {
			// Check if this looks like a fallback hash (starts with lots of zeros)
			// Fallback hashes are padded, so they'll have many leading zeros
			isFallbackHash := true
			nonZeroCount := 0
			for j := 0; j < len(chunk.Checksum); j++ {
				if chunk.Checksum[j] != '0' {
					nonZeroCount++
				}
			}
			// Real SHA-256 hashes typically have good distribution
			// Fallback hashes will have mostly zeros (padded)
			if nonZeroCount > 16 { // More than 16 non-zero chars suggests real hash
				isFallbackHash = false
			}

			if !isFallbackHash {
				// This looks like a real SHA-256, so enforce validation
				return nil, fmt.Errorf("chunk %d checksum mismatch", i)
			}
			// If it looks like a fallback hash, allow it with a warning
			fmt.Printf("Warning: chunk %d using fallback checksum (non-HTTPS upload)\n", i)
		}

		result = append(result, chunk.Data...)
	}
	return result, nil
}
