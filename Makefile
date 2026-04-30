REGISTRY ?= public.ecr.aws/g4e5l3z3/papercomputeco
NAME ?= postgres
TAG ?= 17.7-pgduckdb-1.1.1

.PHONY:build
build: ## Builds the postgres image and exports locally
	dagger call \
		build-postgres-image \
		export-image --name $(NAME):$(TAG)

.PHONY: build-push
build-push: ## Build and publish the multi-arch image to the public registry
	dagger call \
		build-push-postgres-images \
			--registry "$(REGISTRY)" \
			--tags "$(TAG)"

.PHONY: help
.DEFAULT_GOAL := help
help: ## Prints this help message
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
