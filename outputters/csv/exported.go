package csv

import (
	// stdlib
	"log"

	// local
	"go.dev.pztrn.name/glp/outputters/outputinterface"
)

func Initialize() outputinterface.Interface {
	log.Println("Initializing csv outputter...")

	c := &outputter{}
	return outputinterface.Interface(c)
}
