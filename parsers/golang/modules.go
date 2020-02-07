package golang

import (
	// stdlib

	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	// local
	"go.dev.pztrn.name/glp/structs"
)

var goModulesFilesToCheck = []string{"go.mod", "go.sum"}

// Detects if project is using go modules for dependencies management.
func (gp *golangParser) detectModulesUsage(pkgPath string) bool {
	var goModulesFileFound bool
	for _, fileName := range goModulesFilesToCheck {
		pathToCheck := filepath.Join(pkgPath, fileName)
		if _, err := os.Stat(pathToCheck); err == nil {
			goModulesFileFound = true
		}
	}

	if goModulesFileFound {
		log.Println("Project '" + pkgPath + "' is using Go modules for dependencies management")
	}

	return goModulesFileFound
}

// Gets dependencies from go.mod/go.sum files.
func (gp *golangParser) getDependenciesFromModules(pkgPath string) []*structs.Dependency {
	deps := make([]*structs.Dependency, 0)

	// Try to figure out parent package name for all dependencies.
	parent := gp.getParentForDep(pkgPath)

	// Get GOPATH for future dependency path composing.
	gopath, found := os.LookupEnv("GOPATH")
	if !found {
		log.Fatalln("Go modules project found but no GOPATH environment variable defined. Cannot continue.")
	}

	// To get really all dependencies we should use go.sum file.
	filePath := filepath.Join(pkgPath, "go.sum")

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed to open go.sum file for reading:", err.Error())
		return nil
	}

	// We do not need multiple lines of dependencies in reports which
	// describes same name and version.
	createdDeps := make(map[string]bool)

	// Read file data line by line.
	gosum := bufio.NewScanner(f)
	gosum.Split(bufio.ScanLines)

	for gosum.Scan() {
		depLine := strings.Split(gosum.Text(), " ")

		// Version should be cleared out from possible "/go.mod"
		// substring.
		version := strings.Split(depLine[1], "/")[0]

		// Check if we've already processed that dependency.
		_, processed := createdDeps[depLine[0]+"@"+version]
		if processed {
			continue
		}

		// Go modules present on disk either in vendor or in GOPATH/pkg
		// directory. But vendor here should not be trusted because it
		// might contain old versions.
		dependencyPath := filepath.Join(gopath, "pkg", "mod", depLine[0]+"@"+version)

		// Check if this module exists on disk. Absence means that it
		// isn't actually used and just pollute go.sum.
		if _, err := os.Stat(dependencyPath); err != nil {
			continue
		}

		dependency := &structs.Dependency{
			LocalPath: dependencyPath,
			Name:      depLine[0],
			Parent:    parent,
			Version:   version,
		}

		deps = append(deps, dependency)

		// Mark dependency as processed.
		createdDeps[depLine[0]+"@"+version] = true
	}

	return deps
}
