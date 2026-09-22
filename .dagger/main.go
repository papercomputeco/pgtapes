// Package main provides pgtapes' Dagger checks.
package main

import "dagger/pgtapes/internal/dagger"

// Pgtapes provides repository checks.
type Pgtapes struct {
	// +private
	Source *dagger.Directory
}

// New creates the pgtapes Dagger module.
func New(
	// +defaultPath="/"
	// +ignore=[".git", ".direnv", ".devenv", "build", "tmp"]
	source *dagger.Directory,
) *Pgtapes {
	return &Pgtapes{Source: source}
}
