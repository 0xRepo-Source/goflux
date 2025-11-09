# goflux

[![CI](https://github.com/0xRepo-Source/goflux/actions/workflows/ci.yml/badge.svg)](https://github.com/0xRepo-Source/goflux/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xRepo-Source/goflux)](https://goreportcard.com/report/github.com/0xRepo-Source/goflux)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

"Fast, resumable, auth-flexible file transfer over QUIC/SSH, written in Go."

## Quick Start

### Build

```bash
go build -o bin/goflux-server.exe ./cmd/goflux-server
go build -o bin/goflux.exe ./cmd/goflux
```

### Run Server

```bash
# With web UI (default)
.\bin\goflux-server.exe --storage ./data --addr :8080

# Without web UI
.\bin\goflux-server.exe --storage ./data --addr :8080 --web ""

# Then open http://localhost:8080 in your browser
# Domain open http/s:<yourdomain.com> in your browser
```

Options:
- `--storage <dir>` - Directory to store files (default: `./data`)
- `--addr <address>` - Server listen address (default: `:8080`)
- `--web <dir>` - Web UI directory (default: `./web`, use `""` to disable)
- `--version` - Print version

### Use Client

**Upload a file:**
```bash
.\bin\goflux.exe put ./myfile.txt /remote/path/myfile.txt
```

**Download a file:**
```bash
.\bin\goflux.exe get /remote/path/myfile.txt ./downloaded.txt
```

**List files:**
```bash
.\bin\goflux.exe ls /remote/path
```

**Options:**
- `--server <url>` - Server URL (default: `http://localhost:8080`)
- `--chunk-size <bytes>` - Chunk size for uploads (default: 1048576 = 1MB)
- `--version` - Print version

## Features

✅ **Implemented (v0.1.0):**
- HTTP transport for file transfer
- Chunked uploads with integrity verification (SHA-256)
- Automatic chunk reassembly on server
- Local filesystem storage backend
- Simple put/get/ls commands
- **Web UI with drag-and-drop upload and file browser**

🚧 **Planned:**
- Resume support (track partial uploads)
- QUIC transport
- SSH transport
- Parallel chunk uploads
- Progress indicators
- Token authentication
- S3 storage backend
- Capability negotiation

## Architecture

```
goflux/
  cmd/
    goflux-server/  # server binary
    goflux/         # client CLI
  pkg/
    auth/           # SSH, token, JWT (planned)
    proto/          # request/response types
    server/         # session handling, HTTP endpoints
    storage/        # local FS, S3 (planned)
    transport/      # HTTP client, QUIC/SSH (planned)
    chunk/          # chunking, resume, integrity
  internal/
    log/            # logging helpers
    config/         # configuration
```

**📖 See [docs/architecture.md](docs/architecture.md) for detailed architecture diagrams and deployment guides.**

See `coreidea.md` for design philosophy and `docs/overview.md` for more details.
