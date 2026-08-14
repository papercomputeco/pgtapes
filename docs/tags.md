---
title: Image tags
description: The tag scheme, why CloudNativePG constrains it, and which tags on the registry to use and which to avoid.
sidebar:
  order: 4
---

Published references follow `<registry>/postgres:<tag>`, and tags use:

```
{postgres-version}-pgduckdb-{pgduckdb-version}
```

The current published tag is `17.7-pgduckdb-1.1.1`. The `pgduckdb` component
names the upstream `pg_duckdb` image release the extension was staged from.

## Why the scheme is constrained

CloudNativePG detects the PostgreSQL major version from the image tag, so tags
must begin with the PostgreSQL version. The publish pipeline validates every
tag against that rule and a non-conforming tag — including `latest` — fails
the build rather than producing an image CloudNativePG cannot place.

## What the registry actually serves

The registry predates that validation, and it still carries tags from an
earlier publishing scheme:

| tag | status |
| --- | --- |
| `17.7-pgduckdb-1.1.1` | current — use this |
| `latest` | legacy — an older build, not updated by the current pipeline |
| `v0.5.4`, `v0.5.3` | legacy — the pre-scheme version numbers |

`latest` does not point at the newest image: it and `v0.5.4` name the same
older build, published before the current scheme, and nothing updates them.
They also fail CloudNativePG's major-version detection. Pin the versioned
scheme; treat every other tag on the registry as historical.

## Upgrading

Moving a deployment to a newer image means changing the pinned reference —
the `imageName` in a CloudNativePG `Cluster`, or the image in your `docker
run`/Compose configuration. Within the same PostgreSQL major version that is
an ordinary binary swap over the existing data directory. Across PostgreSQL
major versions it is not: PostgreSQL requires a data migration between
majors, and this image does not automate one.

After moving to an image with a newer extension build, upgrade the installed
extensions in each database that uses them:

```sql
ALTER EXTENSION vector UPDATE;
ALTER EXTENSION pg_duckdb UPDATE;
```
