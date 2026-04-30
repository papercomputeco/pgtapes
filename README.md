# `pgtapes`

The Postgres image ready made for `tapes` with a CloudNativePG base + `pg_duckdb`.

## Quickstart

```bash
# build the image locally - dagger exports the image to your local docker images
make build

# run the image - you must provide postgres a password
docker run -it --rm -e POSTGRES_PASSWORD=password postgres:17.7-pgduckdb-1.1.1
```

## Build and publish

Build the local-platform Postgres image:

```sh
dagger call build-postgres-image
```

Publish the multi-platform images:

```sh
dagger call build-push-postgres-images \
  --registry "$REGISTRY" \
  --tags "17.7-pgduckdb-1.1.1"
```

CloudNativePG requires Postgres image tags start with the PostgreSQL version
so it can detect the major version.
We conform with this with these images by using the following version schema:

```
{postgres-version}-pgduckdb-{pgduckdb-version}
```
