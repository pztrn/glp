package outputters

import (
	// stdlib
	"log"

	// local
	"go.dev.pztrn.name/glp/outputters/csv"
	"go.dev.pztrn.name/glp/outputters/outputinterface"
	"go.dev.pztrn.name/glp/structs"
)

var (
	outputters map[string]outputinterface.Interface
)

func Initialize() {
	log.Println("Initializing output providers")

	outputters = make(map[string]outputinterface.Interface)

	csvIface := csv.Initialize()
	outputters["csv"] = csvIface
}

// Write pushes parsed data into outputter for writing.
func Write(outputter string, filePath string, deps []*structs.Dependency) {
	outputterIface, found := outputters[outputter]
	if !found {
		log.Fatalln("Failed to find outputter '" + outputter + "'!")
	}

	outputterIface.Write(deps, filePath)
}
