BINARY := aaasp
BUILD_DIR := ./bin
CMD := ./cmd/aaasp
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build install release tidy clean

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) $(CMD)

install:
	go install -ldflags="-s -w -X main.version=$(VERSION)" $(CMD)

# Cuts a real release via goreleaser — requires GITHUB_TOKEN and a pushed tag.
release:
	goreleaser release --clean

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR) dist/
