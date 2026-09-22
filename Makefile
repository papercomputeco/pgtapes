REGISTRY ?= public.ecr.aws/g4e5l3z3/papercomputeco
NAME ?= postgres
TAG ?= 17.7-pgduckdb-1.1.1
CONTAINER_TOOL ?= docker
PLATFORMS ?= linux/amd64,linux/arm64
BUILDX_BUILDER ?= pgtapes-builder

.PHONY: check
check: ## Build the image and verify PostgreSQL reports ready through Dagger
	dagger check

.PHONY: build
build: ## Build the postgres image and load it into the local image store
	$(CONTAINER_TOOL) build -t $(NAME):$(TAG) -f Dockerfile .

.PHONY: validate-tag
validate-tag: ## Validate that TAG is compatible with CloudNativePG version detection
	@printf '%s\n' '$(TAG)' | grep -Eq '^[0-9]+(\.[0-9]+)*([._-].+)?$$' || \
		{ echo 'TAG must start with the PostgreSQL major version and must not be latest' >&2; exit 1; }

.PHONY: build-push
build-push: validate-tag ## Build and publish the multi-platform image to the public registry
	@$(CONTAINER_TOOL) buildx inspect $(BUILDX_BUILDER) >/dev/null 2>&1 || \
		($(CONTAINER_TOOL) buildx create --driver docker-container --name $(BUILDX_BUILDER) >/dev/null 2>&1 || \
		 $(CONTAINER_TOOL) buildx inspect $(BUILDX_BUILDER) >/dev/null)
	$(CONTAINER_TOOL) buildx build \
		--builder $(BUILDX_BUILDER) \
		--platform $(PLATFORMS) \
		--tag $(REGISTRY)/postgres:$(TAG) \
		--provenance=false \
		--sbom=false \
		--push \
		-f Dockerfile .

.PHONY: help
.DEFAULT_GOAL := help
help: ## Prints this help message
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
