package golang

import (
	// stdlib
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/httpclient"
	"go.dev.pztrn.name/glp/structs"
)

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

	var (
		// This flag shows that we're currently parsing <head> from HTML.
		headCurrentlyParsing bool
	)

	for {
		line, err := resp.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalln("Failed to read HTML response line-by-line:", err.Error())
		} else if err != nil && err == io.EOF {
			break
		}

		if headCurrentlyParsing {
			// Check for go-import data.
			if strings.Contains(line, `<meta name="go-import"`) {
				// Get content.
				// Import things are in element #4.
				lineSplitted := strings.Split(line, `"`)

				// Check line length. This is not so good approach, but
				// should work for 99% of dependencies.
				if len(lineSplitted) < 5 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - tag itself appears to be incomplete. Skipping")
					continue
				}

				if len(lineSplitted) > 5 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - line where meta tag is located appears to be too long. Skipping")
					continue
				}

				// Import line contains data like VCS name and VCS URL.
				// They're delimited with whitespace.
				importDataSplitted := strings.Split(lineSplitted[3], " ")

				// Import line should contain at least 3 elements.
				if len(importDataSplitted) < 3 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - import data is too small. Skipping")
					continue
				}

				// Fill dependency data with this data.
				// First element is a module name and we do not actually
				// need it, because it is already filled previously.
				dependency.VCS.VCS = importDataSplitted[1]
				dependency.VCS.VCSPath = importDataSplitted[2]
			}

			// Check for go-source data.
			if strings.Contains(line, `<meta name="go-source"`) {
				// Get content.
				// Import things are in element #4.
				lineSplitted := strings.Split(line, `"`)

				// Check line length. This is not so good approach, but
				// should work for 99% of dependencies.
				if len(lineSplitted) < 5 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - tag itself appears to be incomplete. Skipping")
					continue
				}

				if len(lineSplitted) > 5 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - line where meta tag is located appears to be too long. Skipping")
					continue
				}

				// Source line contains data like VCS paths templates.
				// They're delimited with whitespace.
				sourceDataSplitted := strings.Split(lineSplitted[3], " ")

				// Source data line should contain at least 3 elements.
				if len(sourceDataSplitted) < 4 {
					log.Println("Got line: '" + line + "', but it cannot be parsed. Probably badly formed - source data is too small. Skipping")
					continue
				}

				// Fill dependency data.
				dependency.VCS.SourceURLDirTemplate = sourceDataSplitted[2]
				dependency.VCS.SourceURLFileTemplate = sourceDataSplitted[3]
			}
		}

		if strings.Contains(strings.ToLower(line), "<head>") {
			headCurrentlyParsing = true
		}

		if strings.Contains(strings.ToLower(line), "</head>") {
			headCurrentlyParsing = false
		}
	}

	if configuration.Cfg.Log.Debug {
		log.Printf("go-import and go-source data parsed: %+v\n", dependency.VCS)
	}
}
