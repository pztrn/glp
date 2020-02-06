package configuration

import (
	// stdlib
	"flag"
	"log"
	"os"
)

var (
	configurationPath string

	Cfg *config
)

// Initialize initializes package.
func Initialize(cfgpath string) {
	log.Println("Initializing configuration")

	configurationPath = cfgpath

	Cfg = &config{}
	err := Cfg.initialize()
	if err != nil {
		log.Println("Error appeared when loading configuration:", err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}
}
