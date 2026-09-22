package main

import (
	"context"

	"dagger/pgtapes/internal/dagger"
)

const postgresPort = 5432

// Test builds the image, boots PostgreSQL, and waits for it to report ready.
//
// +check
func (p *Pgtapes) Test(ctx context.Context) error {
	image := p.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: "Dockerfile",
	})

	postgres := image.
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithExposedPort(postgresPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})

	_, err := image.
		WithEntrypoint(nil).
		WithServiceBinding("postgres", postgres).
		WithExec([]string{
			"pg_isready",
			"--host=postgres",
			"--port=5432",
			"--username=postgres",
			"--dbname=postgres",
		}).
		Sync(ctx)
	return err
}
