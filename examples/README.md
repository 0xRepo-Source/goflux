# Examples

This folder contains example usage scenarios and demo scripts for goflux.

## Demo Script (PowerShell)

Run the demo to see goflux in action:

```powershell
cd examples
.\demo.ps1
```

This demonstrates:
- Starting a goflux server
- Creating test files
- Uploading files with different chunk sizes
- Downloading files back
- Verifying file integrity

## Example Usage

### Basic File Transfer

```bash
# Terminal 1: Start server
../bin/goflux-server --storage ./demo-data

# Terminal 2: Upload a file
../bin/goflux put myfile.txt /files/myfile.txt

# Download it back
../bin/goflux get /files/myfile.txt downloaded.txt
```

### Chunked Transfer

```bash
# Upload large file with 32KB chunks
../bin/goflux --chunk-size 32768 put largefile.bin /backups/largefile.bin
```

### List Files

```bash
../bin/goflux ls /files
../bin/goflux ls /backups
```

## Ideas for More Examples

- Resume after network interruption (planned feature)
- Batch upload multiple files
- Integration with cloud storage backends
- Custom authentication setup
