---
title: pgtapes
description: The PostgreSQL image the tapes server runs on — a CloudNativePG base with pg_duckdb and pgvector built in.
sidebar:
  order: 1
---

`pgtapes` builds the PostgreSQL container image that
[tapes](https://tapes.dev) runs on: a
[CloudNativePG](https://cloudnative-pg.io/) base image with
[`pg_duckdb`](https://github.com/duckdb/pg_duckdb) and
[`pgvector`](https://github.com/pgvector/pgvector) built in.

The published image:

```
public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Built for `linux/amd64` and `linux/arm64`. The registry also carries older
tags you should not use — see [Image tags](./tags.md).

## Why this image exists

Tapes is a session-capture server. It stores captured agent sessions in
PostgreSQL and uses `pgvector` for semantic search over span embeddings, so it
needs a PostgreSQL image with that extension present and loadable rather than
a stock one. `pg_duckdb` is preloaded alongside it to give an analytical query
path over the same data.

Because the base is a CloudNativePG image, one build serves both deployment
shapes: `docker run` (or Compose) for a local stack, and the CloudNativePG
operator in Kubernetes, without being rebuilt. The two paths configure
themselves differently — that split is the main thing worth understanding
about this image, and [Run the image](./running.md) walks both.

This repository contains no server code and no database schema. It is a
Dockerfile, an init script, and a build module; the tapes server owns its own
migrations. To run tapes itself, see
<https://tapes.dev/docs/installation/> — `tapes local up` provisions this
published image for a local stack.

## Where to go

- [Run the image](./running.md) — `docker run` and CloudNativePG, and what
  each path does and does not configure for you.
- [What is in the image](./image.md) — base, extension versions, and the
  mechanisms the build adds.
- [Image tags](./tags.md) — the tag scheme, and which tags on the registry to
  avoid.
- [Build and publish](./building.md) — building the image from this
  repository.
