---
title: Build and publish
description: Build the image locally with the bundled build module, and publish multi-platform images to a registry.
sidebar:
  order: 5
---

The build is a [Dagger](https://dagger.io/) module; `make` targets wrap its
entrypoints. The repository ships a Nix flake with the Go toolchain, `make`,
and `dagger`:

```bash
nix develop      # or `direnv allow`, which loads the same shell
```

`make help` lists every target:

```
build       Builds the postgres image and exports locally
build-push  Build and publish the multi-arch image to the public registry
help        Prints this help message
```

## Build locally

`make build` builds the image for your local platform and loads it into your
local container image store:

```bash
make build
```

That is a thin wrapper over the module's own entrypoint:

```bash
dagger call build-postgres-image export-image --name postgres:17.7-pgduckdb-1.1.1
```

`NAME` and `TAG` override the exported image reference:

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

Which calls:

```bash
dagger call build-push-postgres-images \
  --registry "$REGISTRY" \
  --tags "17.7-pgduckdb-1.1.1"
```

The module validates every tag before building — the constraint and its
reason are in [Image tags](./tags.md). CI builds the image on pushes to `main`
and on pull requests; publishing to the public registry is a manually
triggered workflow.
