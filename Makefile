BINARY := aaasp
BUILD_DIR := ./bin
CMD := ./cmd/aaasp
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build install release tidy clean

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) $(CMD)

install:
	go install -ldflags="-s -w -X main.version=$(VERSION)" $(CMD)

GORELEASER ?= $(shell command -v goreleaser 2>/dev/null || echo $(HOME)/go/bin/goreleaser)

# Cuts a real release via goreleaser. Needs a pushed tag and a clean tree.
# GITHUB_TOKEN falls back to the gh CLI's token. Verify afterwards by
# downloading a published binary and running `--version`.
release:
	GITHUB_TOKEN=$${GITHUB_TOKEN:-$$(gh auth token)} $(GORELEASER) release --clean

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR) dist/
