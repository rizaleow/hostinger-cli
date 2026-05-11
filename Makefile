BINARY      := hostinger-cli
PKG         := github.com/rizaleow/hostinger-cli
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE        := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS     := -s -w \
               -X $(PKG)/internal/version.Version=$(VERSION) \
               -X $(PKG)/internal/version.Commit=$(COMMIT) \
               -X $(PKG)/internal/version.Date=$(DATE)

SPEC_URL ?= https://developers.hostinger.com/openapi/openapi.json
SPEC_FILE := api.json

.PHONY: all build install generate update-spec test lint tidy clean run snapshot

all: build

build:
	go build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(BINARY) ./cmd/$(BINARY)

install:
	go install -trimpath -ldflags '$(LDFLAGS)' ./cmd/$(BINARY)

generate:
	go generate ./...

update-spec:
	curl -sSfL $(SPEC_URL) -o $(SPEC_FILE)
	$(MAKE) generate
	@echo "Spec refreshed. Review with: git diff --stat $(SPEC_FILE) internal/api/zz_generated.go"

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

clean:
	rm -rf bin dist coverage.out

run:
	go run ./cmd/$(BINARY) $(ARGS)

snapshot:
	goreleaser release --snapshot --clean
