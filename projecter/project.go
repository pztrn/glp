package projecter

import (
	// stdlib
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/parsers"
	"go.dev.pztrn.name/glp/structs"

	// other
	"gopkg.in/src-d/go-license-detector.v3/licensedb"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/filer"
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

// Parses license file for copyrights.
func (p *Project) parseLicenseForCopyrights(licencePath string) []string {
	f, err := os.Open(licencePath)
	if err != nil {
		log.Println("Failed to open license file for reading:", err.Error())
		return nil
	}

	var copyrights []string

	// Read file data line by line.
	gosum := bufio.NewScanner(f)
	gosum.Split(bufio.ScanLines)

	for gosum.Scan() {
		line := gosum.Text()

		if strings.HasPrefix(strings.ToLower(line), "copyright ") && !strings.Contains(strings.ToLower(line), "notice") {
			copyrights = append(copyrights, line)
		}
	}

	return copyrights
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

	// Get licensing information for every dependency.
	for _, dep := range p.deps {
		// Prepare dependency's things. For now - only check if
		// file/directory templates defined and, if not, generate
		// them.
		dep.VCS.FormatSourcePaths()

		depDir, err := filer.FromDirectory(dep.LocalPath)
		if err != nil {
			log.Println("Failed to prepare directory path for dependency license scan:", err.Error())
			continue
		}

		licenses, err1 := licensedb.Detect(depDir)
		if err1 != nil {
			log.Println("Failed to detect license for", dep.Name+":", err1.Error())

			dep.License.Name = "Unknown"

			continue
		}

		if configuration.Cfg.Log.Debug {
			log.Printf("Got licenses result for '%s': %+v\n", dep.Name, licenses)
		}

		// Get highest ranked license.
		var (
			licenseFile string
			licenseName string
			licenseRank float32
		)

		for name, result := range licenses {
			if licenseRank < result.Confidence {
				licenseName = name
				licenseRank = result.Confidence

				for fileName, confidence := range result.Files {
					if confidence == licenseRank {
						licenseFile = fileName
						break
					}
				}
			}
		}

		if licenseName == "" {
			dep.License.Name = "Unknown"
			continue
		}

		log.Printf("Got license for '%s': %s", dep.Name, licenseName)

		dep.License.Name = licenseName

		// Generate license URL.
		urlFormatter := strings.NewReplacer("{dir}", "", "{/dir}", "", "{file}", licenseFile, "{/file}", licenseFile, "#L{line}", "")
		dep.License.URL = urlFormatter.Replace(dep.VCS.SourceURLFileTemplate)

		// As we should have dependency locally available we should try
		// to parse license file to get copyrights.
		dep.License.Copyrights = p.parseLicenseForCopyrights(filepath.Join(dep.LocalPath, licenseFile))
	}
}
