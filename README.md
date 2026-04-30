# pgtapes

The Postgres image ready made for `tapes` with a CloudNativePG base + `pg_duckdb`.

## Build

Build the local-platform Postgres image:

```sh
dagger call build-postgres-image
```

Publish multi-platform images:

```sh
dagger call build-push-postgres-images \
  --registry "$REGISTRY" \
  --tags '17.7-pgduckdb-v1.1.1'
```

CloudNativePG requires Postgres image tags to start with the PostgreSQL version
so it can detect the major version.
We version these images with `{postgres-version}-pgduckdb-{pgduckdb-version}`.
