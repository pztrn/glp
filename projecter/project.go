package projecter

import (
	// stdlib
	"log"
	"os"
	"path/filepath"
	"strings"

	// local
	"go.dev.pztrn.name/glp/parsers"
	"go.dev.pztrn.name/glp/structs"
)

// Project represents single project (or package) that was passed via
// -pkgs parameter.
type Project struct {
	packagePath string
	parserName  string
	flavor      string

	deps []*structs.Dependency
}

// NewProject creates new project and returns it.
func NewProject(packagePath string) *Project {
	p := &Project{}
	p.initialize(packagePath)

	return p
}

// GetDeps returns list of dependencies for project.
func (p *Project) GetDeps() []*structs.Dependency {
	return p.deps
}

// Initializes project.
func (p *Project) initialize(packagePath string) {
	p.packagePath = packagePath

	// Prepare package path to be used.
	// First - replace "~" with actual home directory.
	if strings.Contains(p.packagePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("Failed to get user's home directory:", err.Error())
		}

		p.packagePath = strings.Replace(p.packagePath, "~", homeDir, -1)
	}

	// Get absolute path.
	var err error
	p.packagePath, err = filepath.Abs(p.packagePath)
	if err != nil {
		log.Fatalln("Failed to get absolute path for package '"+p.packagePath+":", err.Error())
	}
}

// Starts project parsing.
func (p *Project) process() {
	// We should determine project type.
	p.parserName, p.flavor = parsers.Detect(p.packagePath)

	if p.parserName == "unknown" {
		log.Println("Project", p.packagePath, "cannot be parsed with glp")
		return
	}

	// Lets try to get dependencies, their versions and URLs.
	deps, err := parsers.GetDependencies(p.parserName, p.flavor, p.packagePath)
	if err != nil {
		log.Fatalln("Failed to get dependencies:", err.Error())
	}

	p.deps = deps

	// ToDo: get licensing information.
}
