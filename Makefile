.PHONY: build dev test clean install

# Build the nextgo binary
build:
	go build -o nextgo ./main.go
	@echo "✓ Built nextgo binary"

# Run in development mode
dev:
	go run main.go dev

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f nextgo
	rm -rf .next
	rm -rf dist

# Install globally
install:
	go install ./...
	@echo "✓ Installed nextgo globally"

# Tidy modules
tidy:
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/nextgo-linux-amd64 ./main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/nextgo-darwin-amd64 ./main.go
	GOOS=windows GOARCH=amd64 go build -o dist/nextgo-windows-amd64.exe ./main.go
	@echo "✓ Built for all platforms"

# Create a new project
create:
	go run main.go create my-app

# Run example
example:
	cd examples/api-example && go run main.go
