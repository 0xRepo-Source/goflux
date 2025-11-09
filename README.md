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
- `--token <token>` - Authentication token (or use `GOFLUX_TOKEN` env var)
- `--version` - Print version

### Authentication

**Enable authentication on server:**
```bash
.\bin\goflux-server.exe --storage ./data --tokens tokens.json
```

**Manage tokens with goflux-admin:**
```bash
# Create a token
.\bin\goflux-admin.exe create --user alice --permissions upload,download,list --days 30

# List tokens
.\bin\goflux-admin.exe list

# Revoke a token
.\bin\goflux-admin.exe revoke tok_abc123def456
```

**Use tokens with client:**
```bash
# Via flag
.\bin\goflux.exe --token <your-token> put file.txt /file.txt

# Via environment variable
$env:GOFLUX_TOKEN = "<your-token>"
.\bin\goflux.exe put file.txt /file.txt
```

**Permissions:**
- `upload` - Upload files
- `download` - Download files
- `list` - List files
- `*` - All permissions

## Features

✅ **Implemented (v0.1.0):**
- HTTP transport for file transfer
- Chunked uploads with integrity verification (SHA-256)
- Automatic chunk reassembly on server
- Local filesystem storage backend
- Simple put/get/ls commands
- Web UI with drag-and-drop upload and file browser
- **Token-based authentication with permission control**
- **Admin CLI tool for token management**
- **Token revocation support**

🚧 **Planned:**
- Resume support (track partial uploads)
- QUIC transport
- SSH transport
- Parallel chunk uploads
- Progress indicators
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
