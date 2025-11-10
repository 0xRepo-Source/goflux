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
.\bin\goflux-server.exe

# Then open http://localhost in your browser
# Or access via your domain/IP: http://yourdomain.com
```

The server uses `goflux.json` for configuration (created automatically if missing).

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

### Configuration

goflux uses JSON configuration files instead of command-line flags for cleaner usage:

**Default config (goflux.json):**
```json
{
  "server": {
    "address": "0.0.0.0:80",
    "storage_dir": "./data",
    "webui_dir": "./web",
    "meta_dir": "./.goflux-meta",
    "tokens_file": ""
  },
  "client": {
    "server_url": "http://localhost",
    "chunk_size": 1048576,
    "token": ""
  }
}
```

**Usage with config:**
```bash
# Uses goflux.json by default
.\bin\goflux.exe ls

# Use a different config file
.\bin\goflux.exe --config prod.json put file.txt /file.txt

# Server also uses config
.\bin\goflux-server.exe --config goflux-production.json
```

**Environment variable for tokens:**
```powershell
$env:GOFLUX_TOKEN = "tok_your_token_here"
.\bin\goflux.exe ls
```

**Config priority:** Config file → Environment variable (tokens only)

### Resume Interrupted Uploads

goflux automatically resumes interrupted uploads! If an upload is interrupted (network failure, client crash, etc.), simply run the same `put` command again:

```bash
# Initial upload (interrupted after 50% )
.\bin\goflux.exe put largefile.zip /largefile.zip
# ... network disconnects ...

# Resume upload (automatically skips already-uploaded chunks)
.\bin\goflux.exe put largefile.zip /largefile.zip
# Output: 🔄 Resuming upload: 127/250 chunks already uploaded
```

**How it works:**
- Server tracks upload sessions in metadata files (`.goflux-meta/`)
- Client queries server before uploading to check for existing sessions
- Only missing chunks are uploaded, saving time and bandwidth
- Sessions are automatically cleaned up after successful uploads

### Authentication

**Enable authentication on server:**

Edit your config file to set `tokens_file`:
```json
{
  "server": {
    "tokens_file": "tokens.json"
  }
}
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

Set token in config file or use environment variable:
```powershell
$env:GOFLUX_TOKEN = "tok_your_token_here"
.\bin\goflux.exe put file.txt /file.txt
```

**Permissions:**
- `upload` - Upload files
- `download` - Download files
- `list` - List files
- `*` - All permissions

## Features

✅ **Implemented (v0.3.0):**
- HTTP transport for file transfer
- Chunked uploads with integrity verification (SHA-256)
- Automatic chunk reassembly on server
- **Resume interrupted uploads automatically**
  - Server tracks upload sessions with persistent metadata
  - Client automatically detects and resumes partial uploads
  - Skip already-uploaded chunks to save bandwidth
  - Session cleanup after successful uploads
- **Real-time progress bars** 🆕
  - Visual upload progress with speed and ETA
  - Color-coded progress indicators
  - Resume progress shows new vs existing chunks
  - Spinner for downloads
- **JSON configuration system** 🆕
  - Simple config file management
  - No messy command-line flags
  - Environment variable support for tokens
- Local filesystem storage backend
- Simple put/get/ls commands
- Web UI with drag-and-drop upload and file browser (Material Design dark mode)
- Token-based authentication with permission control
- Admin CLI tool for token management
- Token revocation support

🚧 **Planned:**
- QUIC transport
- SSH transport
- Parallel chunk uploads
- S3 storage backend
- Capability negotiation

## Architecture

```
goflux/
  cmd/
    goflux-server/    # Server binary
    goflux/           # Client CLI
    goflux-admin/     # Token management CLI
  pkg/
    auth/             # Token-based authentication
    server/           # HTTP server and handlers
    storage/          # Storage backends (local filesystem)
    transport/        # HTTP client
    chunk/            # Chunking and integrity verification
    resume/           # Upload session management
    config/           # Configuration file support
  web/                # Web UI (HTML/CSS/JS)
  docs/               # Documentation
  examples/           # Usage examples
```

**📖 See [docs/architecture.md](docs/architecture.md) for detailed architecture diagrams and deployment guides.**

**📝 See [docs/coreidea.md](docs/coreidea.md) for design philosophy.**
