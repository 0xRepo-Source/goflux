package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/0xRepo-Source/goflux/pkg/auth"
	"github.com/0xRepo-Source/goflux/pkg/config"
	"github.com/0xRepo-Source/goflux/pkg/server"
	"github.com/0xRepo-Source/goflux/pkg/storage"
)

func main() {
	configFile := flag.String("config", "goflux.json", "configuration file path")
	addr := flag.String("addr", "", "server address (overrides config)")
	storageDir := flag.String("storage", "", "storage directory (overrides config)")
	webUI := flag.String("web", "", "web UI directory (overrides config)")
	tokenFile := flag.String("tokens", "", "tokens file for authentication (overrides config)")
	metaDir := flag.String("meta", "", "metadata directory for resume sessions (overrides config)")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux-server version: 0.3.0")
		return
	}

	// Load or create configuration
	cfg, err := config.LoadOrCreateConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Command-line flags override config file
	if *addr != "" {
		cfg.Server.Address = *addr
	}
	if *storageDir != "" {
		cfg.Server.StorageDir = *storageDir
	}
	if *webUI != "" {
		cfg.Server.WebUIDir = *webUI
	}
	if *tokenFile != "" {
		cfg.Server.TokensFile = *tokenFile
	}
	if *metaDir != "" {
		cfg.Server.MetaDir = *metaDir
	}

	// Create storage backend
	store, err := storage.NewLocal(cfg.Server.StorageDir)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create server
	srv, err := server.New(store, cfg.Server.MetaDir)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Enable authentication if token file provided
	if cfg.Server.TokensFile != "" {
		tokenStore, err := auth.NewTokenStore(cfg.Server.TokensFile)
		if err != nil {
			log.Fatalf("Failed to load tokens: %v", err)
		}
		srv.EnableAuth(tokenStore)
		fmt.Printf("Loaded authentication from: %s\n", cfg.Server.TokensFile)
	}

	fmt.Printf("Starting goflux-server on %s\n", cfg.Server.Address)
	fmt.Printf("Storage directory: %s\n", cfg.Server.StorageDir)
	fmt.Printf("Configuration file: %s\n", *configFile)

	if err := srv.Start(cfg.Server.Address, cfg.Server.WebUIDir); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
