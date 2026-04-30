# pgtapes

Container image build repo for Tapes Postgres deployments.

## Image

- `postgres`: CloudNativePG-compatible Postgres with `pg_duckdb` and `vector` enabled.

The primary image is built from the top-level `Dockerfile`, following the same simple convention as CloudNativePG's `postgres-containers` repo. Image-specific support files live under `images/postgres/`.

## Building with Dagger

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

CloudNativePG requires Postgres image tags to start with the PostgreSQL version so it can detect the major version. Valid examples include `17`, `17.7`, and `17.7-pgduckdb-v1.1.1`; `latest` is not valid.
