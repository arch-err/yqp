# yqp — interactive TUI playground for yq

# Build the binary
build:
    go build -o yqp .

# Run tests
test:
    go test -count=1 -race ./...

# Lint
lint:
    golangci-lint run

# Verify module integrity
verify:
    go mod verify

# Run all checks (verify + lint + test)
check: verify lint test

# Security audit
audit:
    go vet ./...
    go mod verify
    govulncheck ./...

# Install locally
install:
    go install .

# Clean build artifacts
clean:
    rm -f yqp
