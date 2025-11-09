# Contributing to goflux

Thank you for your interest in contributing to goflux!

## Development Setup

1. **Prerequisites**
   - Go 1.20 or later
   - Git

2. **Clone and build**
   ```bash
   git clone https://github.com/0xRepo-Source/goflux.git
   cd goflux
   go mod download
   go build -o bin/goflux-server ./cmd/goflux-server
   go build -o bin/goflux ./cmd/goflux
   ```

3. **Run tests**
   ```bash
   go test -v ./...
   ```

## Project Structure

- `cmd/` - Binary entry points (server and client)
- `pkg/` - Public libraries
- `internal/` - Private implementation details
- `docs/` - Documentation
- `examples/` - Usage examples

## How to Contribute

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clear, commented code
   - Follow Go best practices
   - Add tests for new features

4. **Test your changes**
   ```bash
   go test ./...
   go build ./...
   ```

5. **Commit with clear messages**
   ```bash
   git commit -m "Add feature: description"
   ```

6. **Push and create a Pull Request**

## Areas We Need Help

- [ ] QUIC transport implementation
- [ ] SSH transport implementation
- [ ] Resume/partial upload tracking
- [ ] Parallel chunk uploads
- [ ] Authentication (token, JWT, OAuth)
- [ ] S3/cloud storage backends
- [ ] Web UI
- [ ] Tests and benchmarks
- [ ] Documentation improvements

## Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

## Questions?

Open an issue for discussion before starting major changes.
