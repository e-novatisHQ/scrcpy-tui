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
CYCLONEDX_GOMOD_VERSION := v1.12.0
MIN_COVERAGE ?= 70

.PHONY: build test coverage vet fmt fmt-check lint vuln workflows tools tools-quality tools-vuln tools-release check pty installer-check fixtures sbom release install notices
build:
	mkdir -p bin
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(LDFLAGS)' -o bin/scrcpy-tui ./cmd/scrcpy-tui
test:
	$(GO) test -race -count=1 -timeout=120s -coverprofile=coverage.out ./...
coverage:
	python3 scripts/check_coverage.py coverage.out $(MIN_COVERAGE)
vet:
	$(GO) vet ./...
fmt:
	gofmt -w cmd internal scripts/render
fmt-check:
	@test -z "$$(gofmt -l cmd internal scripts/render)" || { gofmt -l cmd internal scripts/render; exit 1; }
tools: tools-quality tools-vuln tools-release
tools-quality:
	mkdir -p $(TOOLS_DIR)
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
	$(GO) install github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)
tools-vuln:
	mkdir -p $(TOOLS_DIR)
	$(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
tools-release:
	mkdir -p $(TOOLS_DIR)
	$(GO) install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@$(CYCLONEDX_GOMOD_VERSION)
lint:
	$(TOOLS_DIR)/staticcheck ./...
vuln:
	$(TOOLS_DIR)/govulncheck ./...
workflows:
	$(TOOLS_DIR)/actionlint
pty: build
	python3 scripts/pty_test.py bin/scrcpy-tui
installer-check:
	sh scripts/install_test.sh
fixtures:
	$(GO) run ./scripts/render
check: fmt-check vet test coverage lint workflows pty installer-check
sbom:
	mkdir -p dist
	for target in linux_amd64 linux_arm64 windows_amd64; do \
		GOOS="$${target%_*}" GOARCH="$${target#*_}" CGO_ENABLED=0 $(TOOLS_DIR)/cyclonedx-gomod app -json -output-version 1.6 -packages -noserial -notimestamp -main cmd/scrcpy-tui -output "dist/scrcpy-tui_$(VERSION)_$$target.cdx.json" .; \
	done
release:
	python3 scripts/release.py --version '$(VERSION)'
install: build
	install -Dm755 bin/scrcpy-tui "$${PREFIX:-$$HOME/.local}/bin/scrcpy-tui"

notices:
	python3 scripts/notices.py
