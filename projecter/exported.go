package projecter

import (
	// stdlib
	"log"
	"strings"
	"sync"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/outputters"
	"go.dev.pztrn.name/glp/structs"
)

var (
	packages []string

	outputFormat string
	outputFile   string

	projects      map[string]*Project
	projectsMutex sync.RWMutex
)

// Initialize initializes package.
func Initialize(pkgs string, outFormat string, outFile string) {
	log.Println("Initializing projects handler...")

	packages = strings.Split(pkgs, ",")
	projects = make(map[string]*Project)

	outputFormat = outFormat
	outputFile = outFile

	log.Println("Packages list that was passed:", packages)
}

// GetProject returns project by it's path.
func GetProject(path string) *Project {
	projectsMutex.RLock()
	defer projectsMutex.RUnlock()

	prj, found := projects[path]
	if !found {
		return nil
	}

	return prj
}

// Parse starts projects parsing.
func Parse() {
	// Create project for every passed package.
	// This is done in main goroutine and therefore no mutex is used.
	for _, pkgPath := range packages {
		prj := NewProject(pkgPath)
		projects[pkgPath] = prj
	}

	if configuration.Cfg.Log.Debug {
		log.Printf("Projects generated: %+v\n", projects)
	}

	// We should start asynchronous projects parsing.
	var wg sync.WaitGroup

	for _, prj := range projects {
		wg.Add(1)
		go func(prj *Project) {
			prj.process()
			wg.Done()
		}(prj)
	}

	// Wait until all projects will be parsed.
	wg.Wait()

	// Collect dependencies list from all parsed projects.
	var deps []*structs.Dependency

	for _, prj := range projects {
		deps = append(deps, prj.GetDeps()...)
	}

	outputters.Write(outputFormat, outputFile, deps)

	log.Println("Parsing done")
}
