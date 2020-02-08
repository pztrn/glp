package structs

import (
	// stdlib
	"strings"
)

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

// FormatSourcePaths tries to create templates which will be used for
// paths formatting. E.g. when generating path to license file.
// This is required because for some repositories github.com (and
// probably gitlab.com too) might not return go-source element in
// page's <head> tag.
func (vd *VCSData) FormatSourcePaths() {
	// Do nothing if templates was filled (e.g. when parsing HTML page
	// for repository with "?go-get=1" parameter).
	if vd.SourceURLDirTemplate != "" && vd.SourceURLFileTemplate != "" {
		return
	}

	// If no URL templates was provided by github and we know that
	// dependency is using it as VCS storage - generate proper
	// template URLs.
	if vd.VCS == "git" && vd.VCSPath != "" && strings.Contains(vd.VCSPath, "github.com") {
		repoPathSplitted := strings.Split(vd.VCSPath, ".")
		vd.SourceURLDirTemplate = strings.Join(repoPathSplitted[:len(repoPathSplitted)-1], ".") + "/blob/" + vd.Branch + "{/dir}"
		vd.SourceURLFileTemplate = strings.Join(repoPathSplitted[:len(repoPathSplitted)-1], ".") + "/blob/" + vd.Branch + "{/dir}/{file}#L{line}"
	}
}
