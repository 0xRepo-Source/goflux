# goflux Architecture & Deployment

## Overview

goflux is designed as a client-server file transfer system with support for chunking, resumability, and multiple transport protocols.

## Production Architecture

```mermaid
graph TB
    subgraph Client["Client Machine"]
        CLI[goflux CLI]
    end
    
    subgraph Internet["Internet"]
        HTTPS[HTTPS/TLS Connection]
    end
    
    subgraph Server["Cloud Server / VPS"]
        subgraph Proxy["Reverse Proxy Layer"]
            RP[Nginx/Caddy<br/>Port 443<br/>TLS/SSL + Auth]
        end
        
        subgraph App["Application Layer"]
            GFS[goflux-server<br/>Port 80<br/>HTTP]
        end
        
        subgraph Data["Storage Layer"]
            FS[Local Filesystem<br/>/var/goflux/data]
            S3[S3-Compatible<br/>planned]
        end
    end
    
    CLI -->|Upload/Download| HTTPS
    HTTPS --> RP
    RP -->|Reverse Proxy| GFS
    GFS -->|Read/Write| FS
    GFS -.->|Future| S3
    
    style CLI fill:#4CAF50
    style RP fill:#2196F3
    style GFS fill:#FF9800
    style FS fill:#9C27B0
    style S3 fill:#607D8B
```

## Current Implementation (v0.1.0)

**Local/Development Setup:**

```mermaid
sequenceDiagram
    participant C as goflux Client
    participant S as goflux-server
    participant FS as File Storage
    
    Note over C,S: File Upload Flow
    C->>C: Split file into chunks
    C->>C: Calculate SHA-256 per chunk
    loop For each chunk
        C->>S: POST /upload (chunk data + checksum)
        S->>S: Store chunk in memory
        S->>S: Verify checksum
        S-->>C: Chunk received
    end
    S->>S: Reassemble all chunks
    S->>S: Verify integrity
    S->>FS: Write complete file
    S-->>C: Upload complete
    
    Note over C,S: File Download Flow
    C->>S: GET /download?path=/file.txt
    S->>FS: Read file
    FS-->>S: File data
    S-->>C: Return complete file
```

## Transport Layers

### Current: HTTP
- ✅ Implemented
- Simple REST API
- Works on any network
- No encryption by default

### Planned: QUIC
- Better performance over lossy networks
- Built-in encryption (TLS 1.3)
- Connection migration (mobile-friendly)
- Reduced latency

### Planned: SSH
- Drop-in SFTP replacement
- Familiar authentication (SSH keys)
- Encrypted by default
- Compatible with existing SSH infrastructure

## Deployment Options

### Option 1: Behind Reverse Proxy (Recommended)

```mermaid
graph LR
    A[Internet] -->|HTTPS :443| B[Caddy/Nginx]
    B -->|HTTP :80| C[goflux-server]
    C --> D[Storage]
    
    style B fill:#2196F3
    style C fill:#FF9800
```

**Setup:**
```bash
# 1. Start goflux-server (local only)
./goflux-server --config goflux.json

# 2. Configure Caddy (auto HTTPS)
# Caddyfile:
your-domain.com {
    reverse_proxy localhost:80
}
```

### Option 2: Direct HTTPS (Future)

```mermaid
graph LR
    A[Internet] -->|HTTPS :443| B[goflux-server<br/>with TLS]
    B --> C[Storage]
    
    style B fill:#4CAF50
```

**When implemented:**
```bash
./goflux-server --addr :443 --tls-cert cert.pem --tls-key key.pem --storage /data
```

### Option 3: Local Network Only (Current)

```mermaid
graph LR
    A[Your PC] -->|HTTP :8080| B[goflux-server]
    B --> C[Storage]
    
    style A fill:#4CAF50
    style B fill:#FF9800
```

**Usage:**
```bash
# Server
./goflux-server --addr :8080 --storage ./data

# Client
./goflux --server http://localhost:8080 put file.txt /file.txt
```

## Security Considerations

### ⚠️ Current Version (v0.1.0)

**Not recommended for internet use because:**
- ❌ No encryption (HTTP only)
- ❌ No authentication
- ❌ No access control
- ❌ No rate limiting

**Safe to use:**
- ✅ Localhost (127.0.0.1)
- ✅ Trusted private networks
- ✅ Behind VPN
- ✅ Development/testing

### 🔒 Future Security Features

**Planned for v0.2.0+:**
- 🔐 TLS/HTTPS support
- 🔑 Token-based authentication
- 👤 User management
- 📊 Access logging
- ⚡ Rate limiting
- 🔐 End-to-end encryption option

## Scalability

```mermaid
graph TB
    subgraph LB["Load Balancer"]
        NGINX[Nginx]
    end
    
    subgraph Servers["goflux-server Instances"]
        S1[Server 1]
        S2[Server 2]
        S3[Server 3]
    end
    
    subgraph Storage["Shared Storage"]
        NFS[NFS/S3/MinIO]
    end
    
    NGINX --> S1
    NGINX --> S2
    NGINX --> S3
    S1 --> NFS
    S2 --> NFS
    S3 --> NFS
    
    style NGINX fill:#2196F3
    style S1 fill:#FF9800
    style S2 fill:#FF9800
    style S3 fill:#FF9800
    style NFS fill:#9C27B0
```

**For high availability:**
- Multiple goflux-server instances
- Shared storage backend (S3, NFS, MinIO)
- Load balancer in front
- Session state in external store (Redis - future)

## Network Requirements

### Firewall Rules

**Server-side:**
```bash
# Development (HTTP)
Allow TCP port 8080 from trusted IPs

# Production (with reverse proxy)
Allow TCP port 443 from anywhere
Allow TCP port 8080 from localhost only
```

### Bandwidth Considerations

**Chunk size affects:**
- Memory usage (chunks held in RAM during assembly)
- Network efficiency (smaller chunks = more HTTP overhead)
- Resume granularity (lose at most 1 chunk on disconnect)

**Recommendations:**
- LAN: 1-4 MB chunks
- Internet: 512 KB - 1 MB chunks
- High-latency/lossy: 256-512 KB chunks

## Example Production Setup

**1. Server (Ubuntu 22.04 on DigitalOcean/AWS):**
```bash
# Install goflux
wget https://github.com/0xRepo-Source/goflux/releases/download/v0.1.0/goflux-server-linux-amd64
chmod +x goflux-server-linux-amd64
sudo mv goflux-server-linux-amd64 /usr/local/bin/goflux-server

# Create systemd service
sudo nano /etc/systemd/system/goflux.service
```

```ini
[Unit]
Description=goflux File Transfer Server
After=network.target

[Service]
Type=simple
User=goflux
ExecStart=/usr/local/bin/goflux-server --addr localhost:8080 --storage /var/goflux/data
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

**2. Install Caddy:**
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

**3. Configure Caddy:**
```bash
sudo nano /etc/caddy/Caddyfile
```

```
files.yourdomain.com {
    reverse_proxy localhost:8080
}
```

**4. Start services:**
```bash
sudo systemctl enable --now goflux
sudo systemctl enable --now caddy
```

**5. Use from anywhere:**
```bash
./goflux --server https://files.yourdomain.com put backup.tar.gz /backups/backup.tar.gz
```

## See Also

- [README.md](../README.md) - Quick start guide
- [coreidea.md](../coreidea.md) - Design philosophy
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Development guide
