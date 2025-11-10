# Release Notes - v0.4.1

**Release Date:** November 10, 2025

## Overview

v0.4.1 is a critical patch release that fixes severe memory issues when uploading and downloading large files. This release resolves bugs discovered in v0.4.0 where both client and server would load entire files into memory, causing hangs and crashes with multi-gigabyte files.

## Critical Bug Fixes

### Memory Efficiency Improvements

**Client-Side Streaming Upload** (CRITICAL FIX)
- Fixed memory exhaustion when uploading large files (6GB+)
- Client now streams file in chunks instead of loading entire file with `os.ReadFile()`
- Calculates SHA-256 checksums incrementally per chunk
- Memory usage remains constant regardless of file size
- **Impact:** Enables uploads of arbitrarily large files without client crashes

**Server-Side Streaming Storage** (CRITICAL FIX)
- Fixed memory exhaustion when receiving large file uploads
- Server now writes chunks to disk immediately instead of holding in memory
- Reassembles files from disk chunks when upload completes
- Memory usage remains constant regardless of file size
- **Impact:** Server can now handle multiple concurrent large file uploads

### Technical Details

**Client Changes:**
- `cmd/goflux/main.go`: Rewrote `doPut()` function
  - Changed from `os.ReadFile()` to `os.Open()` + streaming
  - Reads file in `chunker.Size` blocks (1MB by default)
  - Computes SHA-256 hash per chunk while reading
  - Uploads each chunk immediately after reading (no buffering)
  - Resume logic adapted for streaming approach

**Server Changes:**
- `pkg/server/server.go`: Rewrote upload handler
  - Removed in-memory `chunks` map
  - Added `chunksDir` for temporary chunk storage
  - Writes each received chunk to disk: `.goflux-meta/chunks/<session>/<chunk_id>.dat`
  - Reassembles file from disk when all chunks received
  - Cleans up temporary chunks after successful reassembly

## Testing

Tested with:
- 6GB ISO file (backbox-9-desktop-amd64.iso): ✅ Success
- Resume functionality with large files: ✅ Verified
- Small files (< 1MB): ✅ No regression
- Concurrent uploads: ✅ Server handles multiple sessions

## Performance

**Before (v0.4.0):**
- Client: Loaded entire file into memory before chunking
- Server: Stored all chunks in memory before reassembly
- 6GB file: ~12GB RAM usage (client + server)
- Result: Hangs/crashes on large files

**After (v0.4.1):**
- Client: Constant ~1MB memory overhead per upload
- Server: Constant ~1MB memory overhead per concurrent upload
- 6GB file: ~2MB total RAM usage
- Result: Handles files of any size efficiently

## Migration Guide

### From v0.4.0 to v0.4.1

**No Configuration Changes Required**

This is a drop-in replacement for v0.4.0. Simply replace binaries:

```powershell
# Stop server
Stop-Process -Name goflux-server

# Replace binaries
cp goflux-v0.4.1.exe goflux.exe
cp goflux-server-v0.4.1.exe goflux-server.exe
cp goflux-admin-v0.4.1.exe goflux-admin.exe

# Restart server
.\goflux-server.exe
```

**Note on Incomplete Uploads:**

If you had incomplete uploads from v0.4.0 (marked as "completed" but files missing):
1. Delete old session metadata: `Remove-Item .goflux-meta\*.json`
2. Restart server
3. Re-upload affected files (will start fresh)

## Known Issues

- Web UI path issue: Server warns about `./web` not found (cosmetic warning, doesn't affect functionality)
- Session metadata from v0.4.0 may show as "completed" without actual file (delete and re-upload)

## Upgrade Priority

**CRITICAL** - All v0.4.0 users should upgrade immediately

- If you upload files larger than available RAM: **URGENT**
- If you run server with limited memory: **URGENT**  
- If you only transfer small files (< 100MB): **RECOMMENDED**

## Files Changed

- `cmd/goflux/main.go` - Streaming upload implementation
- `pkg/server/server.go` - Disk-based chunk storage
- Version bumped to 0.4.1 in all binaries

## Checksums

```
# Binary checksums will be added after build
goflux.exe: TBD
goflux-server.exe: TBD
goflux-admin.exe: TBD
```

## Contributors

- Matthew Galpin (@0xRepo-Source) - Bug discovery and testing
- GitHub Copilot - Implementation

## Next Steps

v0.5.0 (planned):
- HTTPS/TLS support
- QUIC transport option
- S3 storage backend
- Parallel chunk uploads
- Web UI improvements

---

**Full Changelog:** https://github.com/0xRepo-Source/goflux/compare/v0.4.0...v0.4.1
