package configuration

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// This structure holds whole configuration for glp.
type config struct {
	Log struct {
		Debug bool `yaml:"debug"`
	} `yaml:"log"`
}

// Tries to parse configuration.
func (c *config) initialize() error {
	// Check if file exists.
	if _, err := os.Stat(configurationPath); os.IsNotExist(err) {
		return err
	}

	// Validate configuration file path.
	// First - replace any "~" that might appear.
	if strings.Contains(configurationPath, "~") {
		userDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configurationPath = strings.Replace(configurationPath, "~", userDir, -1)
	}

	// Then - make relative paths to be absolute.
	absPath, err1 := filepath.Abs(configurationPath)
	if err1 != nil {
		return err1
	}

	configurationPath = absPath

	log.Println("Trying to load configuration file data from '" + configurationPath + "'")

	// Read file into memory.
	fileData, err2 := ioutil.ReadFile(configurationPath)
	if err2 != nil {
		return err2
	}

	// Try to parse loaded data.
	err3 := yaml.Unmarshal(fileData, c)
	if err3 != nil {
		return err3
	}

	if c.Log.Debug {
		log.Printf("Configuration parsed: %+v\n", c)
	}

	return nil
}
