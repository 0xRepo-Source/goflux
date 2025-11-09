# goflux — overview

goflux is a resumable, auth-flexible file transfer system over QUIC/SSH (planned) and HTTP.

## Current Status (v0.1.0)

This repository contains a working HTTP-based file transfer implementation with chunking support.

**Implemented:**
- HTTP transport layer
- Chunked uploads with SHA-256 verification
- Automatic chunk reassembly
- Local filesystem storage backend
- Simple CLI (put/get/ls commands)
- Server with concurrent request handling

**Planned:**
- Resume support for interrupted transfers
- QUIC transport for better performance
- SSH transport for drop-in SFTP replacement
- Parallel chunk uploads
- Token/JWT authentication
- S3 and cloud storage backends
- Capability negotiation protocol
- Web UI

## Architecture

### Binaries

- `cmd/gofluxd` — Server binary
- `cmd/goflux` — Client CLI

### Core Libraries

- `pkg/auth` — Authentication (placeholder for SSH, token, JWT)
- `pkg/transport` — Transport layer (HTTP implemented, QUIC/SSH planned)
- `pkg/chunk` — Chunking with SHA-256 integrity verification
- `pkg/storage` — Storage backends (local FS implemented, S3 planned)
- `pkg/server` — HTTP server with chunk handling
- `pkg/proto` — Request/response types

### Internal

- `internal/log` — Logging helpers
- `internal/config` — Configuration management

## Getting Started

See the main [README.md](../README.md) for build and usage instructions.

## Design Philosophy

See [coreidea.md](../coreidea.md) for the core design principles and vision.
