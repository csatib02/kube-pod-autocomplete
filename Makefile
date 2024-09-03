export PATH := $(abspath bin/):${PATH}

PROJECT_NAME = kube-pod-autocomplete
CONTAINER_IMAGE_REF = ghcr.io/csatib02/${PROJECT_NAME}:dev

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
	${KIND_BIN} create cluster --name ${PROJECT_NAME}

.PHONY: down
down: ## Stop development environment
	${KIND_BIN} delete cluster --name ${PROJECT_NAME}

.PHONY: deploy
deploy: container-image ## Deploy kube-pod-autocomplete to the development environment
	kind load docker-image ${CONTAINER_IMAGE_REF} --name ${PROJECT_NAME}
	kubectl create ns ${PROJECT_NAME}
	${HELM_BIN} upgrade --install ${PROJECT_NAME} deploy/charts/${PROJECT_NAME} --namespace ${PROJECT_NAME} --set image.tag=dev

.PHONY: deploy-testdata
deploy-testdata: ## Deploy testdata to the development environment
	kubectl create ns staging
	kubectl create ns prod
	kubectl apply -f e2e/test/

##@ Build

.PHONY: build
build: ## Build binary
	@mkdir -p build
	go build -race -o build/${PROJECT_NAME} .

.PHONY: artifacts
artifacts: container-image helm-chart
artifacts: ## Build artifacts

.PHONY: container-image
container-image: ## Build container image
	docker build -t ${CONTAINER_IMAGE_REF} .

.PHONY: helm-chart
helm-chart: ## Build Helm chart
	@mkdir -p build
	$(HELM_BIN) package -d build/ deploy/charts/${PROJECT_NAME}

##@ Checks

.PHONY: check
check: test lint ## Run tests and lint check

.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests
	go test -race -v -timeout 900s -tags e2e ./e2e/

.PHONY: test-e2e-local
test-e2e-local: container-image ## Run e2e tests locally
	LOAD_IMAGE=${CONTAINER_IMAGE_REF} VERSION=dev ${MAKE} test-e2e

.PHONY: lint
lint: lint-go lint-helm ## Run linters

.PHONY: lint-go
lint-go:
	$(GOLANGCI_LINT_BIN) run $(if ${CI},--out-format github-actions,)

.PHONY: lint-helm
lint-helm:
	$(HELM_BIN) lint deploy/charts/${PROJECT_NAME}

.PHONY: fmt
fmt: ## Format code
	$(GOLANGCI_LINT_BIN) run --fix

##@ Dependencies

deps: bin/golangci-lint bin/kind
deps: ## Install dependencies

# Dependency versions
GOLANGCI_LINT_VERSION = 1.60.3
KIND_VERSION = 0.24.0

# Dependency binaries
GOLANGCI_LINT_BIN := golangci-lint
KIND_BIN := kind
HELM_BIN := helm

bin/golangci-lint:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- v${GOLANGCI_LINT_VERSION}

bin/kind:
	@mkdir -p bin
	curl -Lo bin/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-$(shell uname -s | tr '[:upper:]' '[:lower:]')-$(shell uname -m | sed -e "s/aarch64/arm64/; s/x86_64/amd64/")
	@chmod +x bin/kind

bin/helm:
	@mkdir -p bin
	curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | USE_SUDO=false HELM_INSTALL_DIR=bin bash
	@chmod +x bin/helm
