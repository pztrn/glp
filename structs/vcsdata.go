package structs

// VCSData describes structure of go-import and go-source data.
type VCSData struct {
	// Branch is a VCS branch used.
	Branch string
	// Revision is a VCS revision used.
	Revision string
	// SourceURLDirTemplate is a template for sources dirs URLs. E.g.:
	// https://sources.dev.pztrn.name/pztrn/glp/src/branch/master{/dir}
	SourceURLDirTemplate string
	// SourceURLFileTemplate is a template for sources files URLs. E.g.:
	// https://sources.dev.pztrn.name/pztrn/glp/src/branch/master{/dir}/{file}#L{line}
	SourceURLFileTemplate string
	// VCS is a VCS name (e.g. "git").
	VCS string
	// VCSPath is a VCS repository path.
	VCSPath string
}
