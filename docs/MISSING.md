# goflux - Missing Features & Roadmap

## Current Status: v0.1.0

This document tracks what's **implemented** vs. what's **planned** based on the core vision in `coreidea.md`.

---

## ✅ Implemented (v0.1.0)

### Core Features
- ✅ HTTP transport (basic)
- ✅ Chunked uploads with SHA-256 verification
- ✅ Automatic chunk reassembly on server
- ✅ Local filesystem storage backend
- ✅ Basic CLI commands (`put`, `get`, `ls`)
- ✅ Configurable chunk size
- ✅ Per-chunk integrity verification

### Infrastructure
- ✅ GitHub CI/CD workflows
- ✅ Automated testing (chunk, storage)
- ✅ Multi-platform builds (Linux, Windows, macOS)
- ✅ Documentation with architecture diagrams

---

## ❌ Missing Features

### 🚨 Critical (Security & Production Readiness)

#### 1. **TLS/HTTPS Support**
**Status:** ❌ Not implemented  
**Priority:** 🔴 Critical

**Current:** HTTP only (plain text, insecure)  
**Needed:**
```go
// pkg/server/server.go
func (s *Server) StartTLS(addr, certFile, keyFile string) error {
    return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
```

**Files to create/modify:**
- `pkg/server/server.go` - Add TLS support
- `cmd/goflux-server/main.go` - Add `--tls-cert` and `--tls-key` flags
- `pkg/transport/transport.go` - Update client to use HTTPS

---

#### 2. **Authentication & Authorization**
**Status:** ❌ Not implemented  
**Priority:** 🔴 Critical

**Current:** No auth - anyone can upload/download  
**Needed:**

**Token-based auth:**
```go
// pkg/auth/token.go
type TokenAuth struct {
    tokens map[string]*User
}

func (t *TokenAuth) Validate(token string) (*User, error)
```

**SSH key auth:**
```go
// pkg/auth/ssh.go
type SSHKeyAuth struct {
    authorizedKeys []ssh.PublicKey
}

func (s *SSHKeyAuth) Validate(key ssh.PublicKey) (*User, error)
```

**JWT auth:**
```go
// pkg/auth/jwt.go
type JWTAuth struct {
    secret []byte
}

func (j *JWTAuth) GenerateToken(user *User) (string, error)
func (j *JWTAuth) ValidateToken(token string) (*User, error)
```

**Files to create:**
- `pkg/auth/token.go` - Token authentication
- `pkg/auth/ssh.go` - SSH key authentication  
- `pkg/auth/jwt.go` - JWT authentication
- `pkg/auth/middleware.go` - HTTP auth middleware
- `internal/config/auth.go` - Auth configuration

**Server changes:**
```go
// Add to server handlers
func (s *Server) authenticate(r *http.Request) (*User, error) {
    // Check Authorization header
    // Validate token/JWT
    // Return user or error
}
```

---

#### 3. **Rate Limiting**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Needed:**
```go
// pkg/server/ratelimit.go
type RateLimiter struct {
    requests map[string]*rate.Limiter
}

func (r *RateLimiter) Allow(ip string) bool
```

**Files to create:**
- `pkg/server/ratelimit.go` - Rate limiting logic
- Middleware for HTTP handlers

---

### 🚀 Core Features (From coreidea.md)

#### 4. **Resume Support**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Current:** If upload fails, must restart from beginning  
**Needed:**

```go
// pkg/server/resume.go
type UploadState struct {
    Path         string
    TotalChunks  int
    ReceivedMask []bool  // Track which chunks received
    ExpireAt     time.Time
}

func (s *Server) GetUploadState(path string) (*UploadState, error)
func (s *Server) SaveUploadState(path string, state *UploadState) error
```

**Client changes:**
```bash
# Check server for existing upload state
goflux put --resume ./large-file.iso /backups/file.iso
```

**Files to create:**
- `pkg/server/resume.go` - Upload state tracking
- `pkg/proto/resume.go` - Resume request/response types
- Update `cmd/goflux/main.go` - Add `--resume` flag
- Add persistence (JSON/DB) for upload states

---

#### 5. **QUIC Transport**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Needed:**
```go
// pkg/transport/quic.go
type QUICTransport struct {
    conn quic.Connection
}

func NewQUICClient(addr string) (*QUICTransport, error)
func (q *QUICTransport) Upload(chunk ChunkData) error
func (q *QUICTransport) Download(path string) ([]byte, error)
```

**Dependencies:**
- `github.com/quic-go/quic-go` - QUIC implementation

**Files to create:**
- `pkg/transport/quic.go` - QUIC client
- `pkg/server/quic.go` - QUIC server
- Update CLI: `goflux --transport quic put file.txt /file.txt`

**Benefits:**
- Better performance on lossy networks
- Connection migration (mobile)
- Multiplexing without head-of-line blocking
- Built-in TLS 1.3

---

#### 6. **SSH Transport**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Needed:**
```go
// pkg/transport/ssh.go
type SSHTransport struct {
    client *ssh.Client
}

func NewSSHClient(addr string, config *ssh.ClientConfig) (*SSHTransport, error)
func (s *SSHTransport) Upload(chunk ChunkData) error
func (s *SSHTransport) Download(path string) ([]byte, error)
```

**Dependencies:**
- `golang.org/x/crypto/ssh` - SSH library

**Files to create:**
- `pkg/transport/ssh.go` - SSH client
- `pkg/server/ssh.go` - SSH server
- Update CLI: `goflux --transport ssh put file.txt /file.txt`

**Benefits:**
- Drop-in SFTP replacement
- Familiar SSH key authentication
- Works through firewalls
- Encrypted by default

---

#### 7. **Parallel Chunk Uploads**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Current:** Chunks uploaded sequentially  
**Needed:**

```go
// cmd/goflux/main.go
func doParallelPut(client *transport.HTTPClient, chunks []chunk.Chunk, parallel int) error {
    sem := make(chan struct{}, parallel)
    errChan := make(chan error, len(chunks))
    
    for _, c := range chunks {
        sem <- struct{}{}
        go func(chunk chunk.Chunk) {
            defer func() { <-sem }()
            errChan <- client.UploadChunk(...)
        }(c)
    }
    
    // Wait and collect errors
}
```

**Changes:**
- Update `cmd/goflux/main.go` - Add goroutine pool for parallel uploads
- Add `--parallel N` flag
- Add progress tracking for concurrent uploads

---

#### 8. **Content-Addressed / Deduplication**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Needed:**

```go
// pkg/server/dedupe.go
type DedupeStore struct {
    storage storage.Storage
    hashes  map[string]string  // hash -> path
}

func (d *DedupeStore) CheckExists(hash string) (bool, string)
func (d *DedupeStore) StoreByHash(hash string, data []byte) error
```

**Protocol flow:**
1. Client calculates file hash
2. Client sends manifest: `{path, size, hash}`
3. Server checks if hash exists
4. If exists: server links/copies, no upload needed
5. If not: proceed with chunked upload

**Files to create:**
- `pkg/server/dedupe.go` - Deduplication logic
- `pkg/proto/manifest.go` - Manifest request/response
- Add `/check` endpoint to server

**Benefits:**
- Instant "upload" of duplicate files
- Bandwidth savings
- Storage efficiency

---

#### 9. **Capability Negotiation**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Needed:**

```go
// pkg/proto/capabilities.go
type ServerCapabilities struct {
    Features []string  // "upload", "resume", "dedupe", "streaming", etc.
    MaxChunkSize int
    Transports []string  // "http", "quic", "ssh"
}

func (s *Server) GetCapabilities() *ServerCapabilities
```

**Protocol:**
```bash
# Client queries server
GET /capabilities
→ {"features": ["upload", "resume", "list"], "max_chunk_size": 10485760}

# Client adapts behavior based on response
```

**Files to create:**
- `pkg/proto/capabilities.go` - Capability types
- Add `/capabilities` endpoint to server
- Update client to query capabilities before operations

---

#### 10. **Streaming Downloads**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Current:** Download entire file to disk, then use  
**Needed:**

```bash
# Stream to stdout
goflux get /large-video.mp4 - | ffmpeg -i - output.webm

# Stream to pipe
goflux get /backup.tar.gz - | tar xzf -
```

**Changes:**
- Update `cmd/goflux/main.go` - Detect `-` as stdout
- Stream response body directly to stdout
- Add `--stream` flag for explicit streaming

---

### 📦 Storage Backends

#### 11. **S3-Compatible Storage**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Needed:**

```go
// pkg/storage/s3.go
type S3Storage struct {
    client *s3.Client
    bucket string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucket string) (*S3Storage, error)
func (s *S3Storage) Put(path string, data []byte) error
func (s *S3Storage) Get(path string) ([]byte, error)
```

**Dependencies:**
- `github.com/aws/aws-sdk-go-v2` - AWS SDK

**Files to create:**
- `pkg/storage/s3.go` - S3 storage implementation
- Support for AWS S3, MinIO, DigitalOcean Spaces, etc.

**Server flag:**
```bash
goflux-server --storage s3 --s3-endpoint s3.amazonaws.com --s3-bucket my-bucket
```

---

#### 12. **Multi-Storage Support**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Low

**Needed:**
```go
// pkg/storage/multi.go
type MultiStorage struct {
    backends map[string]storage.Storage  // prefix -> backend
}

// Example: /local/* → local FS, /s3/* → S3
```

---

### 🎨 User Experience

#### 13. **Progress Bars**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Current:** Simple text output  
**Needed:**

```bash
Uploading large-file.iso (4.7 GB)
████████████████████░░░░░░ 68% | 3.2 GB/4.7 GB | 15.3 MB/s | ETA 1m 23s
```

**Dependencies:**
- `github.com/schollz/progressbar/v3` - Progress bar library

**Files to update:**
- `cmd/goflux/main.go` - Add progress tracking

---

#### 14. **Web UI**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Low

**Planned:**
- Simple web interface for browsing files
- Upload/download via browser
- User management interface

**Files to create:**
- `web/index.html` - Web UI
- `web/static/` - CSS, JS
- `pkg/server/web.go` - Serve web UI

---

### 🛠️ Operations & Management

#### 15. **Access Logging**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Needed:**
```go
// internal/log/access.go
func LogRequest(user, method, path string, size int64, duration time.Duration)
```

**Output:**
```
2025-11-09 22:30:15 user@example.com PUT /backups/db.sql 15728640 bytes 2.3s
```

---

#### 16. **Metrics & Monitoring**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Needed:**
- Prometheus metrics endpoint
- Request counters
- Bandwidth tracking
- Error rates
- Chunk processing times

**Dependencies:**
- `github.com/prometheus/client_golang` - Prometheus client

---

#### 17. **User Quotas**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Low

**Needed:**
```go
// pkg/server/quota.go
type Quota struct {
    MaxStorage int64  // bytes
    MaxFiles   int
}

func (s *Server) CheckQuota(user string) error
```

---

#### 18. **Expiring Uploads / TTL**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Low

**Needed:**
```bash
# Upload with 7-day expiration
goflux put --ttl 7d file.zip /temp/file.zip
```

---

### 🧪 Testing & Quality

#### 19. **Integration Tests**
**Status:** ❌ Not implemented  
**Priority:** 🟡 High

**Needed:**
- End-to-end tests (start server, upload, download, verify)
- Network failure simulation
- Concurrent upload tests

**Files to create:**
- `test/integration/` - Integration test suite

---

#### 20. **Benchmarks**
**Status:** ❌ Not implemented  
**Priority:** 🟢 Medium

**Needed:**
```go
// pkg/chunk/chunk_bench_test.go
func BenchmarkChunkerSplit(b *testing.B)
func BenchmarkChunkerReassemble(b *testing.B)
```

---

## 📋 Summary

### By Priority

**🔴 Critical (Security - Required for production):**
- TLS/HTTPS support
- Authentication (token, SSH key, JWT)
- Rate limiting

**🟡 High (Core features):**
- Resume support
- QUIC transport
- SSH transport
- Parallel chunk uploads
- Access logging

**🟢 Medium/Low (Nice-to-have):**
- Content deduplication
- Capability negotiation
- S3 storage backend
- Progress bars
- Web UI
- Metrics
- Streaming downloads

---

## 🗺️ Suggested Roadmap

### v0.2.0 - Security & Production Readiness
- [ ] TLS/HTTPS support
- [ ] Token authentication
- [ ] Rate limiting
- [ ] Access logging
- [ ] Integration tests

### v0.3.0 - Resume & Performance
- [ ] Upload resume support
- [ ] Parallel chunk uploads
- [ ] Progress bars
- [ ] Benchmarks

### v0.4.0 - Advanced Transports
- [ ] QUIC transport
- [ ] SSH transport
- [ ] Multi-transport support

### v0.5.0 - Cloud & Scale
- [ ] S3 storage backend
- [ ] Content deduplication
- [ ] Metrics/monitoring
- [ ] User quotas

### v0.6.0 - Polish
- [ ] Web UI
- [ ] Capability negotiation
- [ ] Streaming downloads
- [ ] Multi-storage routing

---

## 📝 How to Contribute

See [CONTRIBUTING.md](../CONTRIBUTING.md) for development setup.

**Pick a feature from above and:**
1. Open an issue to discuss approach
2. Create a feature branch
3. Implement with tests
4. Submit PR

**Quick wins for first contributors:**
- Progress bars (dependency: progressbar)
- Access logging (use existing log package)
- Integration tests (test/ directory)
- Benchmarks (add _test.go files)
