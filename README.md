# `pgtapes`

The PostgreSQL image that [`tapes`](https://github.com/papercomputeco/tapes) runs
on: a [CloudNativePG](https://cloudnative-pg.io/) base image with
[`pg_duckdb`](https://github.com/duckdb/pg_duckdb) and
[`pgvector`](https://github.com/pgvector/pgvector) built in. Run it directly and
both extensions are created on first start; under the CloudNativePG operator the
`Cluster` manifest asks for them.

Published images:

```
public.ecr.aws/g4e5l3z3/papercomputeco/postgres:17.7-pgduckdb-1.1.1
```

Built for `linux/amd64` and `linux/arm64`.

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

## Documentation

Reference documentation lives in [`docs/`](docs/introduction.md): running the
image standalone and under CloudNativePG, what is in it, the tag scheme, and
building it from this repository.

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
