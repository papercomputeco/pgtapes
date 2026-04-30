package main

import (
	"context"
	"fmt"

	"dagger/pgtapes/internal/dagger"
)

const (
	testPgUser = "tapes"
	testPgPass = "tapes"
	testPgDB   = "tapes"
	testPgPort = 5432
)

// PostgresService provides a ready-to-run Postgres service with the pg_duckdb
// and vector extensions installed for local smoke tests and downstream tests.
func (p *Pgtapes) PostgresService() *dagger.Service {
	ctr := p.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: images[postgresImageName].Dockerfile,
	})

	return ctr.
		WithEnvVariable("POSTGRES_USER", testPgUser).
		WithEnvVariable("POSTGRES_PASSWORD", testPgPass).
		WithEnvVariable("POSTGRES_DB", testPgDB).
		WithExposedPort(testPgPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})
}

// BuildPostgresImage builds the local-platform postgres container image.
func (p *Pgtapes) BuildPostgresImage() *dagger.Container {
	return p.buildDockerfileImage(images[postgresImageName])
}

// BuildPushPostgresImages builds a multi-platform postgres image and publishes
// it to the provided registry.
//
// Published refs use the naming convention: <registry>/postgres:<tag>. Tags
// are validated for CloudNativePG PostgreSQL major-version auto-detection and
// must start with a PostgreSQL version, e.g. "17" or "17.7-pgduckdb-v1.1.1".
func (p *Pgtapes) BuildPushPostgresImages(
	ctx context.Context,

	// Container registry address, e.g. "123456789.dkr.ecr.us-east-1.amazonaws.com".
	registry string,

	// Image tags to apply, e.g. ["17", "17.7-pgduckdb-v1.1.1"].
	tags []string,
) ([]string, error) {
	if err := validateCloudNativePGTags(tags); err != nil {
		return nil, err
	}

	return p.BuildPushImages(ctx, postgresImageName, registry, tags)
}

// PostgresDSN returns the connection string used to reach PostgresService.
func (p *Pgtapes) PostgresDSN() string {
	return fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=%d sslmode=disable", testPgUser, testPgPass, testPgDB, testPgPort)
}
