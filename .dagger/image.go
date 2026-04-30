package main

import (
	"context"
	"fmt"
	"regexp"

	"dagger/pgtapes/internal/dagger"
)

const (
	postgresImageName  = "postgres"
	postgresDockerfile = "Dockerfile"
)

var (
	imagePlatforms = []dagger.Platform{
		"linux/amd64",
		"linux/arm64",
	}

	// CloudNativePG can auto-detect the PostgreSQL major version only when the
	// tag starts with the PostgreSQL major version, optionally followed by
	// numeric dot-separated version components. Examples: 17, 17.7,
	// 17.7-pgduckdb-v1.1.1, 17.7_20260429, 13.3.2.1-1.
	cloudNativePGTagPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+)*([._-].+)?$`)
)

func validateCloudNativePGTags(tags []string) error {
	for _, tag := range tags {
		if !cloudNativePGTagPattern.MatchString(tag) {
			return fmt.Errorf("tag %q is not CloudNativePG-compatible: tags must start with the PostgreSQL major version and must not be latest", tag)
		}
	}
	return nil
}

func (p *Pgtapes) buildDockerfileImage() *dagger.Container {
	return p.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: postgresDockerfile,
	})
}

func (p *Pgtapes) buildDockerfileImageVariants() []*dagger.Container {
	variants := make([]*dagger.Container, 0, len(imagePlatforms))
	for _, platform := range imagePlatforms {
		variant := p.Source.DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: postgresDockerfile,
			Platform:   platform,
		})
		variants = append(variants, variant)
	}

	return variants
}

// BuildImage builds the local-platform container image by name.
func (p *Pgtapes) BuildImage(
	// Image to build. Currently supported: "postgres".
	// +default="postgres"
	image string,
) (*dagger.Container, error) {
	return p.buildDockerfileImage(), nil
}

// BuildPushImages builds a multi-platform image and publishes it to the
// provided registry with each supplied tag.
//
// Published refs use the naming convention: <registry>/<image>:<tag>.
func (p *Pgtapes) BuildPushImages(
	ctx context.Context,

	// Container registry address, e.g. "123456789.dkr.ecr.us-east-1.amazonaws.com".
	registry string,

	// Image tags to apply. For CloudNativePG image auto-detection, postgres tags
	// must start with a PostgreSQL version, e.g. ["17", "17.7-pgduckdb-v1.1.1"].
	tags []string,
) ([]string, error) {
	if err := validateCloudNativePGTags(tags); err != nil {
		return nil, err
	}

	published := []string{}
	variants := p.buildDockerfileImageVariants()

	for _, tag := range tags {
		ref := fmt.Sprintf("%s/%s:%s", registry, postgresImageName, tag)
		addr, err := dag.Container().Publish(ctx, ref, dagger.ContainerPublishOpts{
			PlatformVariants: variants,
		})
		if err != nil {
			return published, fmt.Errorf("failed to publish %s: %w", ref, err)
		}
		published = append(published, addr)
	}

	return published, nil
}
