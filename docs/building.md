---
title: Build and publish
description: Build the image locally with Docker, and publish multi-platform images with Buildx.
sidebar:
  order: 5
---

The repository builds its image directly from the Dockerfile. The Nix flake
provides `make` and the Dagger CLI used by the automatic boot/readiness check
and external PR metadata checks; local image builds use the Docker-compatible
engine installed on the host:

```bash
nix develop      # or `direnv allow`, which loads the same shell
```

`make help` lists every target:

```
check                          Build the image and verify PostgreSQL reports ready through Dagger
build                          Build the postgres image and load it into the local image store
validate-tag                   Validate that TAG is compatible with CloudNativePG version detection
build-push                     Build and publish the multi-platform image to the public registry
help                           Prints this help message
```

`make check` runs the same native-platform Dagger check discovered by Dagger
Cloud on pull requests. It builds the Dockerfile, starts the resulting image,
and requires `pg_isready` to report that PostgreSQL is accepting connections.
It does not publish an image.

## Build locally

`make build` builds the image for your local platform and loads it into your
local container image store:

```bash
make build
```

That invokes the host engine directly:

```bash
docker build -t postgres:17.7-pgduckdb-1.1.1 -f Dockerfile .
```

`NAME`, `TAG`, and `CONTAINER_TOOL` override the local image reference and
Docker-compatible command:

```bash
make build NAME=pgtapes TAG=dev
```

## Publish

`make build-push` builds the multi-platform image (`linux/amd64` and
`linux/arm64`) and pushes it to a registry. It needs credentials for whatever
registry you point it at:

```bash
make build-push REGISTRY=public.ecr.aws/your-alias TAG=17.7-pgduckdb-1.1.1
```

The target provisions or reuses the repository-scoped `pgtapes-builder`
Buildx builder, so multi-platform builds do not depend on the capabilities of
the engine's default builder. `BUILDX_BUILDER` and `PLATFORMS` are
overridable.

The target validates the tag before building — the constraint and its reason
are in [Image tags](./tags.md). CI builds both image platforms on pushes to
`main` and on pull requests; publishing to the public registry is a manually
triggered workflow.
