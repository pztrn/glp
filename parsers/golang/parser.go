package golang

import (
	// stdlib
	"sync"

	// local
	"go.dev.pztrn.name/glp/structs"
)

const (
	// Package managers names. Used in Detect() for flavor returning.
	packageManagerGoMod = "go mod"
	packageManagerDep   = "dep"
)

// This structure responsible for parsing projects that written in Go.
type golangParser struct{}

// Detect detects if passed project path can be parsed with this parser
// and additionally detect package manager used.
func (gp *golangParser) Detect(pkgPath string) (bool, string) {
	// Go projects usually using go modules or dep for dependencies
	// management.
	isModules := gp.detectModulesUsage(pkgPath)
	if isModules {
		return true, packageManagerGoMod
	}

	isDep := gp.detectDepUsage(pkgPath)
	if isDep {
		return true, packageManagerDep
	}

	return false, ""
}

// GetDependencies extracts dependencies from project.
func (gp *golangParser) GetDependencies(flavor string, pkgPath string) []*structs.Dependency {
	var deps []*structs.Dependency

	switch flavor {
	case packageManagerDep:
		deps = gp.getDependenciesFromDep(pkgPath)
	case packageManagerGoMod:
		deps = gp.getDependenciesFromModules(pkgPath)
	}

	// Return early if no dependencies was found.
	if len(deps) == 0 {
		return nil
	}

	// For every dependency we should get additional data - go-import
	// and go-source. Asynchronously.
	var wg sync.WaitGroup

	for _, dep := range deps {
		wg.Add(1)
		go func(dep *structs.Dependency) {
			getGoData(dep)
			wg.Done()
		}(dep)
	}

	wg.Wait()

	return deps
}
