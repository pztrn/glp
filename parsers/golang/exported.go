package golang

import (
	// stdlib
	"log"

	// local
	"go.dev.pztrn.name/glp/parsers/parserinterface"
)

func Initialize() (parserinterface.Interface, string) {
	log.Println("Initializing Golang projects parser")

	goDatas = make(map[string]*godata)

	p := &golangParser{}
	return parserinterface.Interface(p), "golang"
}
