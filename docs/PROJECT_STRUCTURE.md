# Project Structure

This document describes the organization of the goflux project.

## Directory Layout

```
goflux/
├── .github/workflows/     # GitHub Actions CI/CD
│   ├── ci.yml            # Build and test workflow
│   └── release.yml       # Release automation
│
├── bin/                   # Compiled binaries (gitignored)
│   ├── goflux-server.exe
│   ├── goflux.exe
│   └── goflux-admin.exe
│
├── cmd/                   # Command-line applications
│   ├── goflux-server/    # Server binary
│   │   └── main.go
│   ├── goflux/           # Client CLI
│   │   └── main.go
│   └── goflux-admin/     # Admin CLI
│       └── main.go
│
├── pkg/                   # Public libraries
│   ├── auth/             # Authentication
│   │   ├── token.go      # Token storage and validation
│   │   └── middleware.go # HTTP middleware
│   ├── chunk/            # File chunking
│   │   ├── chunk.go      # Chunker implementation
│   │   └── chunk_test.go # Unit tests
│   ├── config/           # Configuration management
│   │   └── config.go     # JSON config loading/saving
│   ├── resume/           # Resume functionality
│   │   └── session.go    # Upload session tracking
│   ├── server/           # HTTP server
│   │   ├── server.go     # Server implementation
│   │   └── web.go        # Web UI serving
│   ├── storage/          # Storage backends
│   │   ├── storage.go    # Storage interface
│   │   └── storage_test.go
│   └── transport/        # Network transport
│       └── transport.go  # HTTP client
│
├── web/                   # Web UI
│   ├── index.html        # Main HTML
│   └── static/
│       ├── style.css     # Styling
│       └── app.js        # Client-side logic
│
├── docs/                  # Documentation
│   ├── architecture.md   # Architecture diagrams
│   ├── AUTHENTICATION.md # Auth guide
│   ├── CONFIGURATION.md  # Config file guide
│   ├── GETTING_STARTED.md # Beginner guide
│   ├── MISSING.md        # Feature roadmap
│   └── coreidea.md       # Design philosophy
│
├── examples/              # Usage examples
│   ├── README.md
│   └── demo.ps1          # PowerShell demo script
│
├── goflux.json            # Default configuration
├── goflux-local.json      # Local testing config (gitignored)
├── goflux-external.json   # External config (gitignored)
├── tokens.json            # Token storage (gitignored)
│
├── .gitignore             # Git ignore rules
├── CONTRIBUTING.md        # Contribution guidelines
├── LICENSE                # MIT License
├── README.md              # Main documentation
├── go.mod                 # Go module definition
└── go.sum                 # Go dependencies checksum
```

## Key Directories

### `/cmd` - Binaries
Contains the main entry points for all executables:
- **goflux-server**: HTTP server with upload/download endpoints
- **goflux**: Client CLI for file operations
- **goflux-admin**: Token management tool

### `/pkg` - Libraries
Reusable packages that implement core functionality:
- **auth**: Token-based authentication and middleware
- **chunk**: File chunking with SHA-256 integrity verification
- **config**: JSON configuration file management
- **resume**: Upload session tracking for resume functionality
- **server**: HTTP server and endpoint handlers
- **storage**: Storage backend interface (currently local filesystem)
- **transport**: HTTP client for file transfer

### `/web` - Web UI
Static files for the browser-based interface:
- Drag-and-drop file uploads
- File browser
- Real-time upload progress

### `/docs` - Documentation
User guides and technical documentation:
- Getting started guides
- Architecture diagrams
- Feature roadmaps
- Authentication setup

### `/examples` - Examples
Sample scripts and usage examples

## Generated/Runtime Files (Gitignored)

These are created during runtime or builds:
- `/bin/` - Compiled executables
- `/data/` - Uploaded files storage
- `/.goflux-meta/` - Resume session metadata
- `tokens.json` - Authentication tokens
- Local configuration overrides

## Configuration Files

- **goflux.json**: Default configuration (committed)
- **goflux-local.json**: Local testing (gitignored)
- **goflux-external.json**: External access (gitignored)

See [CONFIGURATION.md](CONFIGURATION.md) for details.

## Build Artifacts

- Built using standard Go toolchain
- Binaries output to `/bin/`
- GitHub Actions builds for Linux, macOS, Windows
- Release artifacts attached to GitHub releases

## Testing

- Unit tests: `*_test.go` files alongside implementation
- Test data directories are gitignored
- Run tests: `go test ./...`
- CI runs tests on every push

## Dependencies

Managed via Go modules:
- `go.mod` - Module definition and requirements
- `go.sum` - Dependency checksums
- External deps: progressbar/v3

## Version Control

- `.gitignore` excludes binaries, test data, local configs
- Configuration examples are committed
- Documentation is committed
- Test files are excluded
