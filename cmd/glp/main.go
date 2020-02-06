package main

import (
	// stdlib
	"flag"
	"log"
	"os"

	// local
	"go.dev.pztrn.name/glp/configuration"
	"go.dev.pztrn.name/glp/httpclient"
	"go.dev.pztrn.name/glp/outputters"
	"go.dev.pztrn.name/glp/parsers"
	"go.dev.pztrn.name/glp/projecter"
)

var (
	configurationPath string
	packagesPaths     string
	outputFormat      string
	outputFile        string
)

func main() {
	log.Println("Starting glp")

	flag.StringVar(&configurationPath, "config", "./.glp.yaml", "Path to configuration file.")
	flag.StringVar(&packagesPaths, "pkgs", "", "Packages that should be analyzed. Use comma to delimit packages.")
	flag.StringVar(&outputFormat, "outformat", "csv", "Output file format. Only 'csv' for now.")
	flag.StringVar(&outputFile, "outfile", "", "File to write licensing information to.")

	flag.Parse()

	if packagesPaths == "" {
		log.Println("Packages paths that should be analyzed should be defined.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if outputFile == "" {
		log.Println("Output file path should be defined.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	configuration.Initialize(configurationPath)
	parsers.Initialize()
	outputters.Initialize()
	httpclient.Initialize()

	projecter.Initialize(packagesPaths, outputFormat, outputFile)
	projecter.Parse()
}
