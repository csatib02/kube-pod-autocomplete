export PATH := $(abspath bin/):${PATH}

##@ General

# Targets commented with ## will be visible in "make help" info.
# Comments marked with ##@ will be used as categories for a group of targets.

.PHONY: help
default: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: up
up: ## Start development environment
	${KIND_BIN} create cluster --name kube-pod-autocomplete

.PHONY: down
down: ## Stop development environment
	${KIND_BIN} delete cluster --name kube-pod-autocomplete

.PHONY: deploy
deploy: ## Deploy kube-pod-autocomplete to the development environment
	kubectl apply -f deploy/

.PHONY: deploy-testdata
deploy-testdata: ## Deploy testdata to the development environment
	kubectl create ns staging
	kubectl create ns prod
	kubectl apply -f test/testdata/

##@ Build

.PHONY: build
build: ## Build binary
	@mkdir -p build
	go build -race -o build/kube-pod-autocomplete .

.PHONY: artifacts
artifacts: container-image binary-snapshot
artifacts: ## Build artifacts

.PHONY: container-image
container-image: ## Build container image
	docker build .

.PHONY: binary-snapshot
binary-snapshot: ## Build binary snapshot
	VERSION=v${GORELEASER_VERSION} ${GORELEASER_BIN} release --clean --skip=publish --snapshot

##@ Checks

.PHONY: check
check: test lint-go ## Run tests and lint check

.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests
	go test -race -v -timeout 900s -tags e2e ./e2e/

.PHONY: lint-go
lint-go:
	$(GOLANGCI_LINT_BIN) run $(if ${CI},--out-format github-actions,)

.PHONY: fmt
fmt: ## Format code
	$(GOLANGCI_LINT_BIN) run --fix

##@ Dependencies

deps: bin/golangci-lint bin/kind bin/cosign bin/goreleaser
deps: ## Install dependencies

# Dependency versions
GOLANGCI_LINT_VERSION = 1.60.3
KIND_VERSION = 0.24.0
COSIGN_VERSION = 2.4.0
GORELEASER_VERSION = 2.2.0

# Dependency binaries
GOLANGCI_LINT_BIN := golangci-lint
KIND_BIN := kind
GORELEASER_BIN := goreleaser

bin/golangci-lint:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_LINT_VERSION}

bin/kind:
	@mkdir -p bin
	curl -Lo bin/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x bin/kind

# Goreleaser uses cosign for signing binaries
bin/cosign:
	@mkdir -p bin
	@OS=$$(uname -s); \
	case $$OS in \
		"Linux") \
			curl -sSfL https://github.com/sigstore/cosign/releases/download/v${COSIGN_VERSION}/cosign-linux-amd64 -o bin/cosign; \
			;; \
		"Darwin") \
			curl -sSfL https://github.com/sigstore/cosign/releases/download/v${COSIGN_VERSION}/cosign-darwin-arm64 -o bin/cosign; \
			;; \
		*) \
			echo "Unsupported OS: $$OS"; \
			exit 1; \
			;; \
	esac
	@chmod +x bin/cosign

bin/goreleaser:
	@mkdir -p bin
	curl -sfL https://goreleaser.com/static/run -o bin/goreleaser
	@chmod +x bin/goreleaser
