# goflux

"Fast, resumable, auth-flexible file transfer over QUIC/SSH, written in Go."

## Quick Start

### Build

```bash
go build -o bin/gofluxd.exe ./cmd/gofluxd
go build -o bin/goflux.exe ./cmd/goflux
```

### Run Server

```bash
.\bin\gofluxd.exe --storage ./data --addr :8080
```

Options:
- `--storage <dir>` - Directory to store files (default: `./data`)
- `--addr <address>` - Server listen address (default: `:8080`)
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
    gofluxd/        # server binary
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

See `coreidea.md` for design philosophy and `docs/overview.md` for more details.
