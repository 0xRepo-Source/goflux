package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/0xRepo-Source/goflux/pkg/auth"
	"github.com/0xRepo-Source/goflux/pkg/server"
	"github.com/0xRepo-Source/goflux/pkg/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "server address")
	storageDir := flag.String("storage", "./data", "storage directory")
	webUI := flag.String("web", "./web", "web UI directory (empty to disable)")
	tokenFile := flag.String("tokens", "", "tokens file for authentication (empty to disable auth)")
	metaDir := flag.String("meta", "./.goflux-meta", "metadata directory for resume sessions")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux-server version: 0.2.0")
		return
	}

	// Create storage backend
	store, err := storage.NewLocal(*storageDir)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create server
	srv, err := server.New(store, *metaDir)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Enable authentication if token file provided
	if *tokenFile != "" {
		tokenStore, err := auth.NewTokenStore(*tokenFile)
		if err != nil {
			log.Fatalf("Failed to load tokens: %v", err)
		}
		srv.EnableAuth(tokenStore)
		fmt.Printf("Loaded authentication from: %s\n", *tokenFile)
	}

	fmt.Printf("Starting goflux-server on %s\n", *addr)
	fmt.Printf("Storage directory: %s\n", *storageDir)

	if err := srv.Start(*addr, *webUI); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
