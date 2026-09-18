SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
VERSION ?= $(shell cat VERSION)
GO ?= go
export GOTOOLCHAIN ?= go$(shell cat .go-version)
export GOCACHE ?= $(CURDIR)/.cache/go-build
export STATICCHECK_CACHE ?= $(CURDIR)/.cache/staticcheck
TOOLS_DIR := $(CURDIR)/.tools
export GOBIN := $(TOOLS_DIR)
LDFLAGS := -s -w -X main.version=$(VERSION)
STATICCHECK_VERSION := 2026.2.1
GOVULNCHECK_VERSION := v1.8.0
ACTIONLINT_VERSION := v1.7.12

.PHONY: build test vet fmt fmt-check lint vuln workflows tools check pty fixtures release install
build:
	mkdir -p bin
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o bin/scrcpy-tui ./cmd/scrcpy-tui
test:
	$(GO) test -race -count=1 -timeout=120s -coverprofile=coverage.out ./...
vet:
	$(GO) vet ./...
fmt:
	gofmt -w cmd internal scripts/render
fmt-check:
	@test -z "$$(gofmt -l cmd internal scripts/render)" || { gofmt -l cmd internal scripts/render; exit 1; }
tools:
	mkdir -p $(TOOLS_DIR)
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	$(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	$(GO) install github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)
lint:
	$(TOOLS_DIR)/staticcheck ./...
vuln:
	$(TOOLS_DIR)/govulncheck ./...
workflows:
	$(TOOLS_DIR)/actionlint
pty: build
	python3 scripts/pty_test.py bin/scrcpy-tui
fixtures:
	$(GO) run ./scripts/render
check: fmt-check vet test lint workflows pty
release:
	python3 scripts/release.py --version '$(VERSION)'
install: build
	install -Dm755 bin/scrcpy-tui "$${PREFIX:-$$HOME/.local}/bin/scrcpy-tui"

.PHONY: notices
notices:
	python3 scripts/notices.py
