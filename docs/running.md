---
title: Run the image
description: Start the image with docker run or under the CloudNativePG operator, and what each path does and does not configure for you.
sidebar:
  order: 2
---

The image runs two ways, and they configure themselves differently:

- **Standalone** (`docker run`, Compose): the stock PostgreSQL entrypoint runs,
  both extensions are created on first start, and `pg_duckdb` is preloaded with
  no extra configuration.
- **Under the CloudNativePG operator**: the operator supplies its own
  entrypoint and generates its own configuration, so the `Cluster` manifest
  must ask for the extensions explicitly.

## Standalone

PostgreSQL requires a password, and the data directory is initialized on first
start:

```bash
docker run --rm -e POSTGRES_PASSWORD=password -p 5432:5432 \
  public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Confirm both extensions came up:

```bash
psql "postgres://postgres:password@localhost:5432/postgres" \
  -c "select extname, extversion from pg_extension order by extname;"
```

```
  extname  | extversion
-----------+------------
 pg_duckdb | 1.1.0
 plpgsql   | 1.0
 vector    | 0.8.1
```

The entrypoint is the standard PostgreSQL one, so its usual environment
variables apply: `POSTGRES_PASSWORD` is required, and `POSTGRES_USER` and
`POSTGRES_DB` create a different superuser and application database on first
start. The bundled init script runs against that application database, so the
extensions land where the application connects:

```bash
docker run --rm -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=tapes -e POSTGRES_DB=tapes -p 5432:5432 \
  public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Runtime facts, all fixed by the image:

| | |
| --- | --- |
| Port | 5432 |
| Data volume | `/var/lib/postgresql/data` |
| Runs as | uid `26` (`postgres`) |

Extension creation and preloading are entrypoint-and-sample-config mechanisms,
which is why they apply to the standalone path only — the details are in
[What is in the image](./image.md).

## Under CloudNativePG

The CloudNativePG operator uses neither the stock entrypoint nor the sample
configuration. It bootstraps the cluster itself and generates PostgreSQL's
configuration from the `Cluster` spec, so a `Cluster` on this image starts
with no extensions created and `pg_duckdb` not preloaded unless the manifest
asks for them:

```yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: tapes
spec:
  imageName: public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
  postgresql:
    shared_preload_libraries:
      - pg_duckdb # vector needs no preloading
  bootstrap:
    initdb:
      postInitApplicationSQL:
        - CREATE EXTENSION IF NOT EXISTS vector;
        - CREATE EXTENSION IF NOT EXISTS pg_duckdb;
```

`postInitApplicationSQL` runs once, against the application database, when the
cluster is bootstrapped. On a cluster that already exists, issue the same
`CREATE EXTENSION` statements against that database yourself.

CloudNativePG detects the PostgreSQL major version from the image tag, which
constrains what a valid tag looks like — see [Image tags](./tags.md).
