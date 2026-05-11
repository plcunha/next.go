# Contributing to Next.go

Thank you for your interest in contributing to Next.go!

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/nextgo.git
cd nextgo

# Install dependencies (when Go is installed)
go mod download

# Build
go build -o nextgo .

# Run tests
go test ./...
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add tests for new features
- Update documentation as needed

## Reporting Issues

Please use the GitHub issue tracker to report bugs or suggest features.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
