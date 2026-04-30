REGISTRY ?= public.ecr.aws/g4e5l3z3/papercomputeco
NAME ?= postgres
TAG ?= 17.7-pgduckdb-1.1.1

.PHONY:build
build: ## Builds the postgres image and exports locally
	dagger call \
		build-postgres-image \
		export-image --name $(NAME):$(TAG)

.PHONY: build-push
build-push:
	dagger call \
		build-push-postgres-images \
			--registry "$(REGISTRY)" \
			--tags "$(TAG)"
