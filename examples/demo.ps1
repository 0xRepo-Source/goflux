# goflux Demo Script
# This script demonstrates basic goflux functionality

Write-Host "=== goflux Demo ===" -ForegroundColor Cyan

# 1. Start server in background
Write-Host "`n1. Starting server..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "..\bin\gofluxd.exe" -ArgumentList "--storage","./demo-data"
Start-Sleep -Seconds 2

# 2. Create test files
Write-Host "`n2. Creating test files..." -ForegroundColor Yellow
"Hello from goflux!" | Out-File -FilePath small.txt
1..10000 | ForEach-Object { "Line $_ - goflux chunking demo" } | Out-File -FilePath large.txt

$smallSize = (Get-Item small.txt).Length
$largeSize = (Get-Item large.txt).Length
Write-Host "  small.txt: $smallSize bytes"
Write-Host "  large.txt: $largeSize bytes"

# 3. Upload small file
Write-Host "`n3. Uploading small file..." -ForegroundColor Yellow
..\bin\goflux.exe put small.txt /files/small.txt

# 4. Upload large file with custom chunk size
Write-Host "`n4. Uploading large file (32KB chunks)..." -ForegroundColor Yellow
..\bin\goflux.exe put large.txt /files/large.txt --chunk-size 32768

# 5. Download files
Write-Host "`n5. Downloading files..." -ForegroundColor Yellow
..\bin\goflux.exe get /files/small.txt downloaded_small.txt
..\bin\goflux.exe get /files/large.txt downloaded_large.txt

# 6. Verify integrity
Write-Host "`n6. Verifying file integrity..." -ForegroundColor Yellow
$smallMatch = (Get-FileHash small.txt).Hash -eq (Get-FileHash downloaded_small.txt).Hash
$largeMatch = (Get-FileHash large.txt).Hash -eq (Get-FileHash downloaded_large.txt).Hash

if ($smallMatch -and $largeMatch) {
    Write-Host "  ✓ All files verified successfully!" -ForegroundColor Green
} else {
    Write-Host "  ✗ File verification failed!" -ForegroundColor Red
}

# 7. List files
Write-Host "`n7. Listing remote files..." -ForegroundColor Yellow
..\bin\goflux.exe ls /files

Write-Host "`n=== Demo Complete ===" -ForegroundColor Cyan
Write-Host "Server is still running. Press Ctrl+C to stop it." -ForegroundColor Gray
