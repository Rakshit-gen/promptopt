BINARY      := promptopt
PKG         := github.com/rakshit-gen/promptopt
CMD         := ./cmd/promptopt
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE        := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS     := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE)

GO          ?= go
GOBIN       ?= $(shell $(GO) env GOPATH)/bin

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build the promptopt binary into ./bin
	$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o bin/$(BINARY) $(CMD)

.PHONY: install
install: ## Install promptopt into $(GOBIN)
	$(GO) install -trimpath -ldflags '$(LDFLAGS)' $(CMD)

.PHONY: run
run: ## Run promptopt (pass args with ARGS="analyze prompt.txt")
	$(GO) run $(CMD) $(ARGS)

.PHONY: test
test: ## Run the unit test suite
	$(GO) test ./...

.PHONY: test-race
test-race: ## Run tests with the race detector
	$(GO) test -race ./...

.PHONY: cover
cover: ## Run tests and write a coverage report
	$(GO) test -coverprofile=coverage.txt -covermode=atomic ./...
	$(GO) tool cover -func=coverage.txt | tail -1

.PHONY: integration
integration: ## Run integration tests (needs GROQ_API_KEY)
	PROMPTOPT_INTEGRATION=1 $(GO) test -tags=integration -run Integration ./tests/...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint if installed, else gofmt + vet
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found; running gofmt + vet"; \
		test -z "$$(gofmt -l . | tee /dev/stderr)"; \
		$(GO) vet ./...; \
	fi

.PHONY: fmt
fmt: ## Format all Go source
	gofmt -w .

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofmt-clean
	@test -z "$$(gofmt -l . | tee /dev/stderr)" || (echo "run 'make fmt'"; exit 1)

.PHONY: snapshot
snapshot: ## Build a local snapshot release with GoReleaser
	goreleaser release --snapshot --clean --skip=publish

.PHONY: release
release: ## Tag-driven release (run in CI; needs GITHUB_TOKEN)
	goreleaser release --clean

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin dist coverage.txt coverage.html

.PHONY: web-dev
web-dev: ## Run the website dev server
	cd web && npm run dev

.PHONY: web-build
web-build: ## Build the website
	cd web && npm ci && npm run build

.PHONY: help
help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
