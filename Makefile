BINARY   := pong-ball
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE     := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all run build install test lint snapshot release clean

all: build

## run — build & launch the game directly
run:
	go run $(LDFLAGS) . play

## build — compile a binary for the current OS/arch
build:
	go build $(LDFLAGS) -o bin/$(BINARY) .

## install — install to $GOPATH/bin
install:
	go install $(LDFLAGS) .

## test — run all unit tests with race detector
test:
	go test -race -count=1 ./...

## cover — run tests and open coverage report
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## lint — run golangci-lint (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...

## snapshot — build snapshot release (no publish, no git tag needed)
snapshot:
	goreleaser release --snapshot --clean

## release — full release (requires GITHUB_TOKEN + git tag)
release:
	goreleaser release --clean

## clean — remove build artifacts
clean:
	rm -rf bin/ dist/ coverage.out

## help — list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'