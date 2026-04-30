package main

import (
	"context"

	"dagger/pgtapes/internal/dagger"
)

// BuildPostgresImage builds the local-platform postgres container image.
func (p *Pgtapes) BuildPostgresImage() *dagger.Container {
	return p.buildDockerfileImage()
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

	return p.BuildPushImages(ctx, registry, tags)
}
