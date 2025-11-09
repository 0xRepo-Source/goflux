package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/example/goflux/pkg/server"
	"github.com/example/goflux/pkg/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "server address")
	storageDir := flag.String("storage", "./data", "storage directory")
	webUI := flag.String("web", "./web", "web UI directory (empty to disable)")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println("goflux-server version: 0.1.0")
		return
	}

	// Create storage backend
	store, err := storage.NewLocal(*storageDir)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create and start server
	srv := server.New(store)
	fmt.Printf("Starting goflux-server on %s\n", *addr)
	fmt.Printf("Storage directory: %s\n", *storageDir)

	if err := srv.Start(*addr, *webUI); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
