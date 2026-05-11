# Contributing to Next.go

## Development Setup

```bash
git clone https://github.com/nextgo/nextgo.git
cd nextgo
go mod download
go build -o nextgo .
go test ./...
```

## Project Structure

```
nextgo/
├── cmd/                    # CLI commands (dev, build, start, create)
├── packages/
│   ├── router/             # File-system routing
│   ├── server/             # Gin-based HTTP server
│   ├── apihandler/         # Dev mode API handler compilation
│   ├── build/              # Production build system
│   ├── config/             # YAML configuration
│   ├── middleware/          # Logger, CORS, Security
│   ├── watcher/            # File watcher for hot reload
│   ├── static/             # Static file serving
│   └── api/                # API route registration
├── main.go                 # CLI entry point
└── docs/README.md          # Complete documentation
```

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -m 'Add feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## Code Style

- Follow standard Go conventions
- Run `go fmt` and `go vet` before committing
- Add tests for new features
- Update docs if behavior changes

## Running Tests

```bash
go test ./...
go vet ./packages/... ./cmd/...
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
