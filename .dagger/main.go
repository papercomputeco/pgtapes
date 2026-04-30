// Pgtapes CI/CD
//
// Package main provides reproducible image builds and publishing for the
// pgtapes container image repository.
package main

import "dagger/pgtapes/internal/dagger"

// Pgtapes is the main module for the pgtapes CI/CD pipeline.
type Pgtapes struct {
	// Project source directory.
	//
	// +private
	Source *dagger.Directory
}

// New creates a new pgtapes CI/CD module instance.
func New(
	// Project source directory.
	//
	// +defaultPath="/"
	// +ignore=[".git", ".direnv", ".devenv", "build", "tmp"]
	source *dagger.Directory,
) *Pgtapes {
	return &Pgtapes{Source: source}
}
