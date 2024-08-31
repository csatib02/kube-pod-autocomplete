export PATH := $(abspath bin/):${PATH}

##@ General

# Targets commented with ## will be visible in "make help" info.
# Comments marked with ##@ will be used as categories for a group of targets.

.PHONY: help
default: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development


##@ Build

.PHONY: build
build: ## Build binary
	@mkdir -p build
	go build -race -o build/kube-pod-autocomplete .

.PHONY: container-image
container-image: ## Build container image
	docker build .

##@ Checks

.PHONY: check
check: test lint-go ## Run tests and lint check

.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: lint-go
lint-go:
	$(GOLANGCI_LINT_BIN) run $(if ${CI},--out-format github-actions,)

.PHONY: fmt
fmt: ## Format code
	$(GOLANGCI_LINT_BIN) run --fix

##@ Autogeneration

##@ Dependencies

deps: bin/golangci-lint
deps: ## Install dependencies

# Dependency versions
GOLANGCI_LINT_VERSION = 1.60.3

# Dependency binaries
GOLANGCI_LINT_BIN := golangci-lint

bin/golangci-lint:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_LINT_VERSION}
