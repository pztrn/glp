package parserinterface

import (
	// local
	"go.dev.pztrn.name/glp/structs"
)

// Interface is a generic parser interface.
type Interface interface {
	// Detect should return true if project should be parsed using
	// this parser and false otherwise. May optionally return package
	// flavor (e.g. dependency management utility name).
	Detect(pkgPath string) (bool, string)
	// GetDependencies parses project for dependencies.
	GetDependencies(flavor string, pkgPath string) []*structs.Dependency
}
