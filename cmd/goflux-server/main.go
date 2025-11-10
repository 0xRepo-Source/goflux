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
	configFile := flag.String("config", "goflux.json", "path to configuration file")
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
