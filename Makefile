SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

.PHONY: help verify fmt fmt-check lint test test-race coverage build smoke integration install clean open-pilot-check docs-check

help:
	@awk 'BEGIN {FS = ":.*##"; printf "Commands:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

verify: fmt-check lint test-race coverage build smoke open-pilot-check docs-check ## Run all checks

fmt: ## Format Go files
	@find cmd internal -type f -name '*.go' -print0 | xargs -0 -r gofmt -w

fmt-check: ## Fail when Go files need formatting
	@test -z "$$(find cmd internal -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"

lint: ## Run Go vet and golangci-lint when available
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else echo 'golangci-lint not installed; go vet completed'; fi

test: ## Run Go tests
	go test ./...

test-race: ## Run tests with the race detector
	go test -race ./...

coverage: ## Generate a coverage profile
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out | tail -n 1
	@total="$$(go tool cover -func=coverage.out | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
		awk -v total="$$total" 'BEGIN {if (total + 0 < 80) {printf "coverage %.1f%% is below 80%%\n", total; exit 1}}'

build: ## Build the CLI
	go build -trimpath -o bin/twitchdl ./cmd/twitchdl

smoke: build ## Check the built CLI without contacting Twitch
	./bin/twitchdl --help >/dev/null
	./bin/twitchdl --version
	./bin/twitchdl doctor

integration: build ## Record and validate a generated local HLS stream
	./tests/integration-local-hls.sh

install: build ## Install for the current user
	./scripts/install.sh

clean: ## Remove generated build and coverage artifacts
	go clean
	@rm -f -- coverage.out bin/twitchdl

docs-check: ## Verify README local documentation links
	@test -f docs/usage.md
	@test -f docs/live-recording.md
	@test -f docs/architecture.md
	@test -f docs/troubleshooting.md

open-pilot-check: ## Verify Open Pilot issue template exists
	@test -f .gitea/ISSUE_TEMPLATE/open-pilot-task.yaml
	@grep -q 'Test Command' .gitea/ISSUE_TEMPLATE/open-pilot-task.yaml
