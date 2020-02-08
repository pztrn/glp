package golang

import (
	// stdlib
	"log"
	"os"
	"path/filepath"
	"strings"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/structs"

	// other
	"github.com/BurntSushi/toml"
)

var goDepFilesToCheck = []string{"Gopkg.toml", "Gopkg.lock"}

type depLockConfig struct {
	Projects []struct {
		Branch    string
		Digest    string
		Name      string
		Packages  []string
		PruneOpts string
		Revision  string
		Version   string
	}
	SolveMeta struct {
		AnalyzerName    string   `toml:"analyzer-name"`
		AnalyzerVersion int      `toml:"analyzer-version"`
		InputImports    []string `toml:"input-imports"`
		SolverName      string   `toml:"solver-name"`
		SolverVersion   int      `toml:"solver-version"`
	} `toml:"solve-meta"`
}

// Detects if project is using dep for dependencies management.
func (gp *golangParser) detectDepUsage(pkgPath string) bool {
	var goDepFilesFound bool
	for _, fileName := range goDepFilesToCheck {
		pathToCheck := filepath.Join(pkgPath, fileName)
		if _, err := os.Stat(pathToCheck); err == nil {
			goDepFilesFound = true
		}
	}

	if goDepFilesFound {
		log.Println("Project '" + pkgPath + "' is using dep for dependencies management")
	}

	return goDepFilesFound
}

// Gets dependencies data from dep-enabled projects.
func (gp *golangParser) getDependenciesFromDep(pkgPath string) []*structs.Dependency {
	deps := make([]*structs.Dependency, 0)

	// Try to figure out parent package name for all dependencies.
	parent := gp.getParentForDep(pkgPath)

	// All dependencies for project will be taken from Gopkg.lock file.
	lockFile := &depLockConfig{}
	_, err := toml.DecodeFile(filepath.Join(pkgPath, "Gopkg.lock"), lockFile)
	if err != nil {
		log.Fatalln("Failed to parse dep lock file:", err.Error())
	}

	if configuration.Cfg.Log.Debug {
		log.Printf("dep lock file parsed: %+v\n", lockFile)
	}

	// Parse dependencies.
	for _, dep := range lockFile.Projects {
		dependency := &structs.Dependency{
			Name:   dep.Name,
			Parent: parent,
			VCS: structs.VCSData{
				Branch:   dep.Branch,
				Revision: dep.Revision,
			},
			Version: dep.Version,
		}

		// If branch is empty - assume master.
		if dependency.VCS.Branch == "" {
			dependency.VCS.Branch = "master"
		}

		// If version is empty - write branch name here with commit hash.
		if dependency.Version == "" {
			dependency.Version = dependency.VCS.Revision + "@" + dependency.VCS.Branch
		}

		// All dep-controlled dependencies are vendored. We should get
		// it's path.
		dependency.LocalPath = filepath.Join(pkgPath, "vendor", dep.Name)

		deps = append(deps, dependency)

		if configuration.Cfg.Log.Debug {
			log.Printf("Initial dependency structure formed: %+v\n", dependency)
		}
	}

	return deps
}

// Tries to get package name for passed package path.
func (gp *golangParser) getParentForDep(pkgPath string) string {
	// Dep-managed projects are in 99% of cases are placed in GOPATH.
	if strings.Contains(pkgPath, "src") {
		return strings.Split(pkgPath, "src/")[1]
	}

	return ""
}
