package structs

// Dependency represents single dependency data.
type Dependency struct {
	// License is a license name for dependency.
	License License
	// LocalPath is a path to dependency (if vendored or in GOPATH or
	// in module cache).
	LocalPath string
	// Name is a dependency name as it appears in package manager's
	// lock file or in sources if no package manager is used.
	Name string
	// Parent is a path to parent package.
	Parent string
	// VCS is a VCS data obtained for dependency.
	VCS VCSData
	// Version is a dependency version used in project.
	Version string
	// URL is a web URL for that dependency (Github, Gitlab, etc.).
	URL string
}
