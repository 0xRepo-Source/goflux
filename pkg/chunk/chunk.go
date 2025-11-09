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
		if hex.EncodeToString(hash[:]) != chunk.Checksum {
			return nil, fmt.Errorf("chunk %d checksum mismatch", i)
		}
		result = append(result, chunk.Data...)
	}
	return result, nil
}
