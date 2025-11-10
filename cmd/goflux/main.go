package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0xRepo-Source/goflux/pkg/chunk"
	"github.com/0xRepo-Source/goflux/pkg/config"
	"github.com/0xRepo-Source/goflux/pkg/transport"
	"github.com/schollz/progressbar/v3"
)

func main() {
	// Global flags
	configFile := flag.String("config", "goflux.json", "configuration file path")
	serverAddr := flag.String("server", "", "server address (overrides config)")
	chunkSize := flag.Int("chunk-size", 0, "chunk size in bytes (overrides config)")
	token := flag.String("token", "", "authentication token (overrides config and env var)")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux version: 0.3.0")
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	// Load or create configuration
	cfg, err := config.LoadOrCreateConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Command-line flags override config file
	if *serverAddr != "" {
		cfg.Client.ServerURL = *serverAddr
	}
	if *chunkSize > 0 {
		cfg.Client.ChunkSize = *chunkSize
	}

	// Get token: flag > config > environment
	authToken := *token
	if authToken == "" && cfg.Client.Token != "" {
		authToken = cfg.Client.Token
	}
	if authToken == "" {
		authToken = os.Getenv("GOFLUX_TOKEN")
	}

	client := transport.NewHTTPClient(cfg.Client.ServerURL)
	if authToken != "" {
		client.SetAuthToken(authToken)
	}

	chunker := chunk.New(cfg.Client.ChunkSize)

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

	// Query server for existing upload session
	status, err := client.QueryUploadStatus(remotePath)

	var bar *progressbar.ProgressBar
	var totalToUpload int

	if err != nil {
		fmt.Printf("⚠️  Could not query upload status: %v (starting fresh upload)\n", err)
		totalToUpload = len(chunks)
		bar = progressbar.NewOptions(totalToUpload,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionSetDescription("[cyan]Uploading...[reset]"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	} else if status.Exists && !status.Completed {
		alreadyUploaded := len(chunks) - len(status.MissingChunks)
		fmt.Printf("🔄 Resuming upload: %d/%d chunks already uploaded\n", alreadyUploaded, len(chunks))

		totalToUpload = len(status.MissingChunks)
		bar = progressbar.NewOptions(totalToUpload,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionSetDescription("[cyan]Resuming...[reset]"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[yellow]=[reset]",
				SaucerHead:    "[yellow]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		// Create a map of missing chunks for quick lookup
		missingMap := make(map[int]bool)
		for _, chunkID := range status.MissingChunks {
			missingMap[chunkID] = true
		}

		// Only upload missing chunks
		for i, c := range chunks {
			if !missingMap[c.ID] {
				// Chunk already uploaded, skip it
				continue
			}

			chunkData := transport.ChunkData{
				Path:     remotePath,
				ChunkID:  c.ID,
				Data:     c.Data,
				Checksum: c.Checksum,
				Total:    len(chunks),
			}

			if err := client.UploadChunk(chunkData); err != nil {
				bar.Close()
				return fmt.Errorf("failed to upload chunk %d: %w", i, err)
			}
			bar.Add(1)
		}

		bar.Finish()
		fmt.Printf("\n✓ Resume complete: uploaded %d new chunks\n", totalToUpload)
		return nil
	} else {
		// Fresh upload - no existing session
		totalToUpload = len(chunks)
		bar = progressbar.NewOptions(totalToUpload,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionSetDescription("[cyan]Uploading...[reset]"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	}

	// Upload chunks
	for i, c := range chunks {
		chunkData := transport.ChunkData{
			Path:     remotePath,
			ChunkID:  c.ID,
			Data:     c.Data,
			Checksum: c.Checksum,
			Total:    len(chunks),
		}

		if err := client.UploadChunk(chunkData); err != nil {
			bar.Close()
			return fmt.Errorf("failed to upload chunk %d: %w", i, err)
		}
		bar.Add(1)
	}

	bar.Finish()
	fmt.Printf("\n✓ Upload complete: %s → %s\n", localPath, remotePath)
	return nil
}

func doGet(client *transport.HTTPClient, remotePath, localPath string) error {
	fmt.Printf("Downloading %s...\n", remotePath)

	// Create indeterminate progress bar (we don't know size beforehand)
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("[cyan]Downloading...[reset]"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionEnableColorCodes(true),
	)
	bar.RenderBlank()

	data, err := client.Download(remotePath)
	if err != nil {
		bar.Close()
		return err
	}

	bar.Finish()
	fmt.Printf("\n")

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
