---
title: What is in the image
description: The base image, extension versions, and the entrypoint and configuration mechanisms the build adds on top of it.
sidebar:
  order: 3
---

| | |
| --- | --- |
| Base | `ghcr.io/cloudnative-pg/postgresql:17.7-standard-bookworm` |
| PostgreSQL | 17.7 |
| `pg_duckdb` | 1.1.0, staged from `pgduckdb/pgduckdb:17-v1.1.1` |
| `pgvector` | 0.8.1, from the CloudNativePG standard base |
| Platforms | `linux/amd64`, `linux/arm64` |

The `pg_duckdb` extension version (1.1.0) and the image it is staged from
(v1.1.1) genuinely differ — the tag names the upstream image release, the
extension inside it reports 1.1.0.

## What the build adds

The Dockerfile adds two things to the CloudNativePG base:

- **`pg_duckdb`.** Its runtime artifacts (`libduckdb.so`, `pg_duckdb.so`,
  bitcode, and extension SQL) are copied from the official `pg_duckdb` image,
  and `pg_duckdb` is appended to `shared_preload_libraries` in
  `postgresql.conf.sample`. PostgreSQL copies that sample when it initializes
  a data directory, which is how the standalone path gets the preload for
  free — and why the CloudNativePG path, which generates its own
  configuration, does not.
- **The standard PostgreSQL entrypoint.** `docker-entrypoint.sh`, the initdb
  helpers, and `gosu` are staged in so the image also runs under a plain
  `docker run`, not only under the CloudNativePG operator.

One init script, `initdb.d/0000-install-extensions.sql`, is installed into
`/docker-entrypoint-initdb.d/`. The standard entrypoint runs it the first time
a data directory is initialized, which issues:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_duckdb;
```

Both mechanisms — the sample config and the init directory — belong to the
stock entrypoint, so they apply to the standalone path only. Under
CloudNativePG the `Cluster` manifest carries the equivalent settings; see
[Run the image](./running.md).

## Checking a running container

Against a standalone container started as in [Run the image](./running.md):

```bash
psql "postgres://postgres:password@localhost:5432/postgres" \
  -c "show shared_preload_libraries;"
```

```
 shared_preload_libraries
--------------------------
 pg_duckdb
```

Both extensions answer queries:

```bash
psql "postgres://postgres:password@localhost:5432/postgres" \
  -c "select '[1,2,3]'::vector <-> '[3,2,1]'::vector as l2;" \
  -c "select * from duckdb.query('select 42 as answer');"
```

```
         l2
--------------------
 2.8284271247461903
(1 row)

 answer
--------
     42
(1 row)
```
