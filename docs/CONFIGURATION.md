# Configuration File Guide

goflux now supports configuration files for easy management of server and client settings.

## Quick Start

### 1. Generate Default Config

When you run the server or client for the first time, it will create a default `goflux.json`:

```bash
.\bin\goflux-server.exe --config goflux.json
# Creates goflux.json with defaults if it doesn't exist
```

### 2. Edit Configuration

Edit `goflux.json` to match your setup:

```json
{
  "server": {
    "address": "0.0.0.0:80",
    "storage_dir": "./data",
    "webui_dir": "./web",
    "meta_dir": "./.goflux-meta",
    "tokens_file": "",
    "tls_cert": "",
    "tls_key": ""
  },
  "client": {
    "server_url": "http://95.145.216.175",
    "chunk_size": 1048576,
    "token": ""
  }
}
```

### 3. Use the Config

**Server:**
```bash
.\bin\goflux-server.exe --config goflux.json
```

**Client:**
```bash
.\bin\goflux.exe --config goflux.json put file.txt /file.txt
```

## Configuration Options

### Server Section

| Field | Description | Example |
|-------|-------------|---------|
| `address` | Listen address and port | `"0.0.0.0:80"` or `":9000"` |
| `storage_dir` | Directory to store uploaded files | `"./data"` |
| `webui_dir` | Web UI directory (empty to disable) | `"./web"` or `""` |
| `meta_dir` | Metadata directory for resume sessions | `"./.goflux-meta"` |
| `tokens_file` | Path to tokens file (empty to disable auth) | `"tokens.json"` or `""` |
| `tls_cert` | TLS certificate file (for HTTPS) | `"cert.pem"` or `""` |
| `tls_key` | TLS private key file (for HTTPS) | `"key.pem"` or `""` |

### Client Section

| Field | Description | Example |
|-------|-------------|---------|
| `server_url` | Server URL to connect to | `"http://95.145.216.175"` |
| `chunk_size` | Chunk size in bytes | `1048576` (1MB) |
| `token` | Authentication token | `"your-token-here"` or `""` |

## Multiple Configurations

You can maintain different configs for different scenarios:

**Local Testing (`goflux-local.json`):**
```json
{
  "client": {
    "server_url": "http://localhost",
    ...
  }
}
```

**External Access (`goflux-external.json`):**
```json
{
  "client": {
    "server_url": "http://95.145.216.175",
    ...
  }
}
```

**Production (`goflux-prod.json`):**
```json
{
  "server": {
    "address": "0.0.0.0:443",
    "tokens_file": "tokens.json",
    "tls_cert": "cert.pem",
    "tls_key": "key.pem"
  },
  "client": {
    "server_url": "https://yourdomain.com",
    "token": "your-secure-token"
  }
}
```

Then use them:
```bash
.\bin\goflux.exe --config goflux-local.json put file.txt /file.txt
.\bin\goflux.exe --config goflux-external.json put file.txt /file.txt  
.\bin\goflux.exe --config goflux-prod.json put file.txt /file.txt
```

## Environment Configuration Pattern

**Best practice:** Keep environment-specific configs separate

**goflux.json** (default, checked into git):
```json
{
  "server": {
    "address": "0.0.0.0:80",
    "storage_dir": "./data"
  },
  "client": {
    "server_url": "http://localhost",
    "chunk_size": 1048576
  }
}
```

**goflux-prod.json** (production, in .gitignore):
```json
{
  "client": {
    "server_url": "https://myserver.com",
    "token": "tok_production_secret"
  }
}
```

Then:
```bash
# Development (uses goflux.json)
.\bin\goflux.exe ls

# Production
.\bin\goflux.exe --config goflux-prod.json ls
```

## Benefits

✅ **No Hardcoding** - Change IPs/domains/ports without recompiling  
✅ **Multiple Environments** - Easy switching between dev/staging/prod  
✅ **Team Sharing** - Commit default configs to git (except tokens!)  
✅ **Simpler CLI** - No messy flags to remember  
✅ **Portable** - Share config files between team members

## Security Note

⚠️ **Never commit tokens to git!** Either:
- Use empty `token` in config and set via `GOFLUX_TOKEN` env var
- Add `goflux-*-prod.json` to `.gitignore`
- Keep tokens only in environment variables
