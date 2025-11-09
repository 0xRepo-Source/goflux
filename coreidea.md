# Core Ideas

## Multiple transports

**SSH (for drop-in SFTP replacement)**
- Compatible with existing SSH infrastructure
- Familiar authentication model
- Easy adoption for ops teams

**QUIC (for speed, mobile, high-latency, UDP-friendly)**
- Modern transport built on UDP
- Better performance in lossy/high-latency networks
- Connection migration for mobile clients
- Reduced head-of-line blocking

**HTTP(S) control API (for listing, tokens, etc.)**
- RESTful API for file operations
- Token-based authentication
- Easy integration with web services and tooling

## Content-addressed or checksum-aware uploads

Client sends file manifest (size + hash) → server checks if already present → instant dedupe.

- Avoid re-uploading identical content
- Server-side deduplication
- Faster uploads when content exists
- Bandwidth savings

## Chunked + resumable

Files are uploaded in chunks with IDs, so if the connection dies, you resume at chunk N.

- Resilient to network interruptions
- No need to restart large transfers
- Progress tracking per chunk
- Integrity verification per chunk

## Capabilities instead of a fixed SFTP verb set

Client asks: "what do you support?" → server replies: `{upload, resume, list, stat, mkdir, share, tags}`

- Extensible protocol
- Feature negotiation
- Backward compatibility
- Server can advertise custom extensions
- Clients gracefully degrade if features unavailable

## Modern auth

**SSH keys for CLI**
- Standard public key authentication
- Integration with ssh-agent
- Familiar workflow for developers

**Tokens / API keys for service-to-service**
- Long-lived credentials for automation
- Scoped permissions
- Revocable without key rotation

**Optional OAuth/JWT for web UI**
- Modern web authentication flows
- Single sign-on integration
- Short-lived tokens with refresh

## Streaming downloads

Ability to stream a file directly to stdout or another service.

- Pipe files to processing tools
- No intermediate disk writes
- Composable with Unix tools
- Efficient for large files

## Client CLI (what users will love)

**Upload**
```bash
goflux put ./backup.tar /remote/backups/
```

**Resume if broken**
```bash
goflux put --resume ./backup.tar /remote/backups/
```

**Download**
```bash
goflux get /remote/backups/backup.tar .
```

**List**
```bash
goflux ls /remote/backups
```

### Flags

- `--transport ssh|quic|https` — Choose the transport layer
- `--parallel 4` — Number of parallel chunk transfers
- `--checksum sha256` — Checksum algorithm for integrity
- `--show-progress` — Display transfer progress

## Why it's "better than SFTP"

**Resumable by design** (SFTP usually needs client logic to do this)
- Built-in chunk tracking
- Automatic resume state management
- No manual offset calculation

**Parallel chunking for high-bandwidth links**
- Multiple chunks in flight simultaneously
- Better utilization of available bandwidth
- Configurable concurrency

**Pluggable storage (local, S3, MinIO)**
- Not limited to filesystem
- Cloud-native storage backends
- Custom storage implementations

**Modern transports (QUIC)**
- UDP-based for better performance
- Connection migration
- Reduced latency

**Extensible protocol (declare capabilities)**
- Feature negotiation at connection time
- Forward-compatible
- Custom extensions without breaking compatibility
