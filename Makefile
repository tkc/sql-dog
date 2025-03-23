# Linting tasks
lint:
	$(HOME)/go/bin/golangci-lint run --no-config --disable-all --enable=gofmt,govet,errcheck

# Testing tasks
test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-cover:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Build tasks
build:
	go build -o bin/sql-dog-lint ./cmd/lint/main.go
	go build -o bin/sql-dog-clean ./cmd/clean/main.go
	@echo "Binaries built in ./bin/"

# Installation task
install:
	go install ./cmd/lint
	go install ./cmd/clean

# Clean up tasks
clean:
	rm -rf bin/

.PHONY: lint test test-race test-cover build install clean