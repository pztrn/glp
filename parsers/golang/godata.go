package golang

import (
	// stdlib
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/httpclient"
	"go.dev.pztrn.name/glp/structs"
)

// attrValue returns the attribute value for the case-insensitive key
// `name', or the empty string if nothing is found.
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

// Gets go-import and go-source data and fill it in dependency.
func getGoData(dependency *structs.Dependency) {
	// Dependencies are imported using URL which can be called with
	// "?go-get=1" parameter to obtain required VCS data.
	req, _ := http.NewRequest("GET", "http://"+dependency.Name, nil)

	q := req.URL.Query()
	q.Add("go-get", "1")

	req.URL.RawQuery = q.Encode()

	respBody := httpclient.GET(req)
	if respBody == nil {
		return
	}

	// HTML is hard to parse properly statically, so we will go
	// line-by-line for <head> parsing.
	resp := bytes.NewBuffer(respBody)

	// Adopted headers parsing algo from Go itself.
	// See https://github.com/golang/go/blob/95e1ea4598175a3461f40d00ce47a51e5fa6e5ea/src/cmd/go/internal/get/discovery.go

	decoder := xml.NewDecoder(resp)
	decoder.Strict = false

	for {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				log.Fatalln("Failed to parse dependency's go-source and go-import things:", err.Error())
			}

			break
		}

		if e, ok := token.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			break
		}

		if e, ok := token.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			break
		}

		e, ok := token.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}

		// Check if we haven't found "go-import" or "go-source" in token's
		// attributes.
		if attrValue(e.Attr, "name") != "go-import" && attrValue(e.Attr, "name") != "go-source" {
			continue
		}

		// Parse go-import data first.
		if attrValue(e.Attr, "name") == "go-import" {
			if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 3 {
				dependency.VCS.VCS = f[1]
				dependency.VCS.VCSPath = f[2]
			}
		}

		// Then - go-source data.
		if attrValue(e.Attr, "name") == "go-source" {
			if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 4 {
				dependency.VCS.SourceURLDirTemplate = f[2]
				dependency.VCS.SourceURLFileTemplate = f[3]
			}
		}
	}

	if configuration.Cfg.Log.Debug {
		log.Printf("go-import and go-source data parsed: %+v\n", dependency.VCS)
	}
}
