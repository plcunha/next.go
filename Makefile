.PHONY: build dev test clean install tidy build-all create example

# Build the nextgo binary
build:
	go build -o nextgo .
	@echo "Built nextgo binary"

# Run in development mode
dev:
	go run . dev

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	if exist nextgo del nextgo
	if exist .next rmdir /s /q .next
	if exist dist rmdir /s /q dist

# Install globally (equivalent to npm install -g next)
# Usage: make install
# After install, run: nextgo --help
install:
	go install .
	@echo "Installed nextgo globally to %GOPATH%\bin"

# Tidy modules
tidy:
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/nextgo-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/nextgo-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o dist/nextgo-windows-amd64.exe .
	@echo "Built for all platforms"

# Create a new project
create:
	go run . create my-app

# Run example
example:
	cd examples/api-example && go run main.go
