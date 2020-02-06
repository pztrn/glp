package golang

import (
	// stdlib
	"log"
	"os"
	"path/filepath"
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
