package csv

import (
	// stdlib
	c "encoding/csv"
	"log"
	"os"
	"strconv"

	// local
	"go.dev.pztrn.name/glp/structs"
)

var (
	headers = []string{"Module", "License", "Repository URL", "License URL", "Project", "Copyrights"}
)

// Responsible for pushing passed data into CSV file.
type outputter struct{}

func (o *outputter) Write(deps []*structs.Dependency, outFile string) {
	log.Println("Got", strconv.Itoa(len(deps)), "dependencies to write")

	// Check if file exists and remove it if so.
	if _, err := os.Stat(outFile); !os.IsNotExist(err) || err == nil {
		os.Remove(outFile)
	}

	// Open file and create writer.
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Failed to open '"+outFile+"' for writing:", err.Error())
	}

	writer := c.NewWriter(f)
	writer.Comma = ';'

	// Write header first.
	_ = writer.Write(headers)

	// Write dependencies information.
	for _, dep := range deps {
		_ = writer.Write([]string{dep.Name, dep.License.Name, dep.VCS.VCSPath, dep.License.URL, dep.Parent, dep.License.Copyrights})
	}

	writer.Flush()

	f.Close()
}
