# `pgtapes`

The PostgreSQL image that [`tapes`](https://github.com/papercomputeco/tapes) runs
on: a [CloudNativePG](https://cloudnative-pg.io/) base image with
[`pg_duckdb`](https://github.com/duckdb/pg_duckdb) and
[`pgvector`](https://github.com/pgvector/pgvector) available and enabled on first
start.

Published images:

```
public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Built for `linux/amd64` and `linux/arm64`.

## What's in the image

| | |
| --- | --- |
| Base | `ghcr.io/cloudnative-pg/postgresql:17.7-standard-bookworm` |
| PostgreSQL | 17.7 |
| `pg_duckdb` | 1.1.0, staged from `pgduckdb/pgduckdb:17-v1.1.1` |
| `pgvector` | 0.8.1, from the CloudNativePG standard base |
| Runs as | uid `26` (`postgres`) |
| Port | 5432 |
| Data volume | `/var/lib/postgresql/data` |

The image adds two things to the CloudNativePG base:

- **`pg_duckdb`.** Its runtime artifacts (`libduckdb.so`, `pg_duckdb.so`, bitcode,
  and extension SQL) are copied from the official `pg_duckdb` image, and
  `pg_duckdb` is appended to `shared_preload_libraries` in the sample config so
  it loads at startup.
- **The standard Postgres entrypoint.** `docker-entrypoint.sh`, the initdb
  helpers, and `gosu` are staged in so the same image works both under the
  CloudNativePG operator and with a plain `docker run`.

`initdb.d/0000-install-extensions.sql` is installed into
`/docker-entrypoint-initdb.d/`. The standard entrypoint runs it the first time a
data directory is initialized, which issues:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_duckdb;
```

## Quickstart

Run the published image directly. PostgreSQL requires a password, and the data
directory is initialized on first start:

```bash
docker run --rm -e POSTGRES_PASSWORD=password -p 5432:5432 \
  public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Then confirm both extensions came up:

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

## Building

The build is a [Dagger](https://dagger.io/) module. `make build` builds the
image for your local platform and loads it into your local container image
store:

```bash
make build
```

That is a thin wrapper over the module's own entrypoint:

```bash
dagger call build-postgres-image export-image --name postgres:17.7-pgduckdb-1.1.1
```

`make help` lists every target. `NAME` and `TAG` override the exported image
reference:

```bash
make build NAME=pgtapes TAG=dev
```

## Publishing

`make build-push` builds the multi-platform image and pushes it to a registry.
It needs credentials for whatever registry you point it at:

```bash
make build-push REGISTRY=public.ecr.aws/your-alias TAG=17.7-pgduckdb-1.1.1
```

Which calls:

```bash
dagger call build-push-postgres-images \
  --registry "$REGISTRY" \
  --tags "17.7-pgduckdb-1.1.1"
```

Published references follow `<registry>/postgres:<tag>`.

## Image tags

CloudNativePG detects the PostgreSQL major version from the image tag, so tags
must begin with the PostgreSQL version. These images use:

```
{postgres-version}-pgduckdb-{pgduckdb-version}
```

Tags are validated at publish time and a non-conforming tag — including
`latest` — fails the build rather than producing an image CloudNativePG cannot
place.

## How `tapes` uses it

[`tapes`](https://github.com/papercomputeco/tapes) is a session-capture server
that stores captured agent sessions in PostgreSQL, and uses `pgvector` for
semantic search over span embeddings — so it needs a PostgreSQL image with that
extension present and loadable rather than a stock one. `pgtapes` is that image:
`tapes local up` provisions this published image for a local stack. `pg_duckdb`
is preloaded alongside it to give an analytical query path over the same data.
Because the base is a CloudNativePG image, the same build also runs under the
CloudNativePG operator in Kubernetes without being rebuilt.

## Development

The repository ships a Nix flake with the Go toolchain, `make`, and `dagger`:

```bash
nix develop      # or `direnv allow`, which loads the same shell
```

## License

Dual-licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE))
- MIT license ([LICENSE-MIT](LICENSE-MIT))

at your option. Unless you explicitly state otherwise, any contribution
intentionally submitted for inclusion in the work by you, as defined in the
Apache-2.0 license, shall be dual licensed as above, without any additional
terms or conditions.

These terms cover this repository's own sources — the Dockerfile, build module,
init scripts, and packaging. The container image they produce bundles
third-party software, including PostgreSQL, `pgvector`, `pg_duckdb`, DuckDB, and
the CloudNativePG base image, each of which remains under its own license.
