package outputinterface

import (
	// local
	"go.dev.pztrn.name/glp/structs"
)

// Interface is a generic output writer interface.
type Interface interface {
	Write(deps []*structs.Dependency, outFile string)
}
