package parsers

import (
	// stdlib
	"errors"
	"log"
	"sync"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/parsers/golang"
	"go.dev.pztrn.name/glp/parsers/parserinterface"
	"go.dev.pztrn.name/glp/structs"
)

var (
	parsers      map[string]parserinterface.Interface
	parsersMutex sync.RWMutex
)

// Initialize initializes package.
func Initialize() {
	log.Println("Initializing parsers...")

	parsers = make(map[string]parserinterface.Interface)

	// Initialize parsers.
	golangIface, golangName := golang.Initialize()
	parsers[golangName] = golangIface
}

// Detect tries to launch parsers for project detection. It returns
// parser name that should be used and optional flavor (e.g. dependencies
// manager name) that might be returned by parser's Detect() function.
func Detect(pkgPath string) (string, string) {
	parsersMutex.RLock()
	defer parsersMutex.RUnlock()

	for parserName, parserIface := range parsers {
		if configuration.Cfg.Log.Debug {
			log.Println("Checking if parser '" + parserName + "' can parse project '" + pkgPath + "'...")
		}

		useThisParser, flavor := parserIface.Detect(pkgPath)
		if useThisParser {
			return parserName, flavor
		}
	}

	return "unknown", ""
}

// GetDependencies asks parser to extract dependencies from project.
func GetDependencies(parserName string, flavor string, pkgPath string) ([]*structs.Dependency, error) {
	parsersMutex.RLock()
	defer parsersMutex.RUnlock()
	parser, found := parsers[parserName]

	if !found {
		return nil, errors.New("parser with such name isn't registered")
	}

	deps := parser.GetDependencies(flavor, pkgPath)

	return deps, nil
}
