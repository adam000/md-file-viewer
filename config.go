package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go4.org/xdgdir"
)

const configFileName = "md-file-viewer.conf"

type config struct {
	RootDir     string
	StyleDir    string
	TemplateDir string
}

// Validate checks that the config is valid (directories exist).
func (c config) Validate() error {
	def := defaults()
	if strings.TrimSpace(c.RootDir) == "" {
		c.RootDir = def.RootDir
		log.Printf("Config has no value for root directory, assuming default of '%s'", c.RootDir)
	}
	if _, err := os.Stat(c.RootDir); err != nil {
		return fmt.Errorf("Could not find root directory '%s': %s", c.RootDir, err)
	}

	if strings.TrimSpace(c.StyleDir) == "" {
		c.StyleDir = def.StyleDir
		log.Printf("Config has no value for style directory, assuming default of '%s'", c.StyleDir)
	}
	if _, err := os.Stat(c.StyleDir); err != nil {
		return fmt.Errorf("Could not find style directory '%s': %s", c.StyleDir, err)
	}

	if strings.TrimSpace(c.TemplateDir) == "" {
		c.TemplateDir = def.TemplateDir
		log.Printf("Config has no value for template directory, assuming default of '%s'", c.TemplateDir)
	}
	if _, err := os.Stat(c.TemplateDir); err != nil {
		return fmt.Errorf("Could not find template directory '%s': %s", c.TemplateDir, err)
	}

	return nil
}

func defaults() config {
	return config{
		RootDir:     ".",
		StyleDir:    "css",
		TemplateDir: "templates",
	}
}

// Create a new config file and return the default config.
func createNewConfigFile() config {
	cfg := defaults()

	//TODO finish this
	//xdgdir.Config.Create(configFileName)

	return cfg
}

func loadConfig() (config, error) {
	cfg := defaults()

	var bytes []byte
	file, err := xdgdir.Config.Open(configFileName)
	if err != nil {
		// Try current directory for config.
		bytes, err = ioutil.ReadFile(configFileName)
		if err != nil {
			return defaults(), fmt.Errorf("Could not find file named %s in $XDG_CONFIG_HOME or the same directory as the executable", configFileName)
		}
	} else {
		// File was opened before, close it and let ReadAll reopen.
		file.Close()
		bytes, err = ioutil.ReadAll(file)
		if err != nil {
			return defaults(), fmt.Errorf("Could not read config file: %s", err)
		}
	}

	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return defaults(), fmt.Errorf("Could not read JSON config: %s", err)
	}

	return cfg, nil
}
