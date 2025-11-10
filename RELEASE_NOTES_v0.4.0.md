# goflux v0.4.0 Release Notes

## 🎉 Major Changes

### Simplified Configuration System
- **BREAKING CHANGE:** Removed command-line flag overrides in favor of JSON config files
- Server now only accepts `--config` and `--version` flags
- Client now only accepts `--config` and `--version` flags
- All settings configured via `goflux.json` file
- Cleaner CLI with no confusing flag precedence
- Better support for multiple environments (dev/staging/prod)

### Port Change
- **Default port changed from 8080 to 80**
- Server: `0.0.0.0:80`
- Client: `http://localhost` (port 80)
- Easier for production deployments

### Web UI Improvements
- **Material Design dark mode** - Clean, flat, professional interface
- Removed all emoji icons - Now uses geometric shapes
- Better color contrast and accessibility
- Roboto font for Material look
- Smoother animations with Material motion curves

### Bug Fixes
- **Fixed crypto.subtle error** for non-HTTPS contexts
- Added fallback hash function for HTTP uploads
- Auto-add `http://` scheme when missing from server URL
- Improved checksum validation logic
- Fixed test failures in chunk validation

## 📦 What's Included

- `goflux.exe` - Client CLI for uploading/downloading files
- `goflux-server.exe` - HTTP server with web UI
- `goflux-admin.exe` - Token management tool
- Configuration files:
  - `goflux.json` - Default configuration
  - `goflux.example.json` - Template for new configs

## 🚀 Quick Start

1. Download the release
2. Extract files
3. Run server: `.\bin\goflux-server.exe`
4. Open web UI: http://localhost
5. Or use CLI: `.\bin\goflux.exe ls`

## 📝 Configuration

Create `goflux.json`:
```json
{
  "server": {
    "address": "0.0.0.0:80",
    "storage_dir": "./data",
    "webui_dir": "./web",
    "tokens_file": ""
  },
  "client": {
    "server_url": "http://localhost",
    "chunk_size": 1048576,
    "token": ""
  }
}
```

## 🔄 Upgrading from v0.3.0

**Breaking Changes:**
- Replace command-line flags with config file settings
- Update server address references from `:8080` to `:80`
- Update client URLs from `http://localhost:8080` to `http://localhost`

**Migration:**
```bash
# Old (v0.3.0)
.\bin\goflux-server.exe --addr :8080 --storage ./data --tokens tokens.json
.\bin\goflux.exe --server http://localhost:8080 --token "xxx" ls

# New (v0.4.0)
# Edit goflux.json:
{
  "server": {"address": "0.0.0.0:80", "storage_dir": "./data", "tokens_file": "tokens.json"},
  "client": {"server_url": "http://localhost", "token": "xxx"}
}

# Then:
.\bin\goflux-server.exe
.\bin\goflux.exe ls
```

## ✅ Features (Unchanged)

- Chunked uploads with SHA-256 verification
- Automatic resume for interrupted uploads
- Real-time progress bars with ETA
- Token-based authentication
- Web UI with drag-and-drop
- Cross-platform (Windows, macOS, Linux)

## 🐛 Known Issues

- HTTPS/TLS support not yet implemented (use reverse proxy)
- QUIC transport planned for future release
- S3 storage backend in development

## 📚 Documentation

- [Getting Started Guide](docs/GETTING_STARTED.md)
- [Configuration Guide](docs/CONFIGURATION.md)
- [Authentication Guide](docs/AUTHENTICATION.md)
- [Architecture Overview](docs/architecture.md)

## 🙏 Acknowledgments

Thanks to all contributors and users for feedback and bug reports!

---

**Full Changelog**: https://github.com/0xRepo-Source/goflux/compare/v0.3.0...v0.4.0
