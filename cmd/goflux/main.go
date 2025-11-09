package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/example/goflux/pkg/chunk"
	"github.com/example/goflux/pkg/transport"
)

func main() {
	// Global flags
	serverAddr := flag.String("server", "http://localhost:8080", "server address")
	chunkSize := flag.Int("chunk-size", 1024*1024, "chunk size in bytes")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux client version: 0.1.0")
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	client := transport.NewHTTPClient(*serverAddr)
	chunker := chunk.New(*chunkSize)

	command := args[0]
	switch command {
	case "put":
		if len(args) < 3 {
			fmt.Println("Usage: goflux put <local-file> <remote-path>")
			os.Exit(1)
		}
		if err := doPut(client, chunker, args[1], args[2]); err != nil {
			log.Fatalf("Upload failed: %v", err)
		}
	case "get":
		if len(args) < 3 {
			fmt.Println("Usage: goflux get <remote-path> <local-file>")
			os.Exit(1)
		}
		if err := doGet(client, args[1], args[2]); err != nil {
			log.Fatalf("Download failed: %v", err)
		}
	case "ls":
		path := "/"
		if len(args) > 1 {
			path = args[1]
		}
		if err := doList(client, path); err != nil {
			log.Fatalf("List failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func doPut(client *transport.HTTPClient, chunker *chunk.Chunker, localPath, remotePath string) error {
	// Read local file
	data, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Split into chunks
	chunks := chunker.Split(data)
	fmt.Printf("Uploading %s (%d bytes, %d chunks)...\n", localPath, len(data), len(chunks))

	// Upload each chunk
	for i, c := range chunks {
		chunkData := transport.ChunkData{
			Path:     remotePath,
			ChunkID:  c.ID,
			Data:     c.Data,
			Checksum: c.Checksum,
			Total:    len(chunks),
		}

		if err := client.UploadChunk(chunkData); err != nil {
			return fmt.Errorf("failed to upload chunk %d: %w", i, err)
		}
		fmt.Printf("  Uploaded chunk %d/%d\n", i+1, len(chunks))
	}

	fmt.Printf("✓ Upload complete: %s → %s\n", localPath, remotePath)
	return nil
}

func doGet(client *transport.HTTPClient, remotePath, localPath string) error {
	fmt.Printf("Downloading %s...\n", remotePath)

	data, err := client.Download(remotePath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✓ Download complete: %s → %s (%d bytes)\n", remotePath, localPath, len(data))
	return nil
}

func doList(client *transport.HTTPClient, path string) error {
	files, err := client.List(path)
	if err != nil {
		return err
	}

	fmt.Printf("Files in %s:\n", path)
	for _, f := range files {
		fmt.Printf("  %s\n", f)
	}
	return nil
}

func printUsage() {
	fmt.Println("goflux - Fast, resumable file transfer")
	fmt.Println("\nUsage:")
	fmt.Println("  goflux put <local-file> <remote-path>   Upload a file")
	fmt.Println("  goflux get <remote-path> <local-file>   Download a file")
	fmt.Println("  goflux ls [path]                        List files")
	fmt.Println("\nFlags:")
	fmt.Println("  --server <url>        Server address (default: http://localhost:8080)")
	fmt.Println("  --chunk-size <bytes>  Chunk size (default: 1048576)")
	fmt.Println("  --version             Print version")
}
