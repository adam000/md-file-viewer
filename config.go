package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adam000/goutils/config"
)

const configFileName = "md-file-viewer.conf"

type configuration struct {
	RootDir     string
	StyleDir    string
	TemplateDir string
	Port        int
}

// Validate checks that the configuration is valid (directories exist; port is a valid port).
func (c configuration) Validate() error {
	def := defaults()
	if strings.TrimSpace(c.RootDir) == "" {
		c.RootDir = def.RootDir
		log.Printf("Configuration has no value for root directory, assuming default of '%s'", c.RootDir)
	}
	if _, err := os.Stat(c.RootDir); err != nil {
		return fmt.Errorf("Could not find root directory '%s': %s", c.RootDir, err)
	}

	if strings.TrimSpace(c.StyleDir) == "" {
		c.StyleDir = def.StyleDir
		log.Printf("Configuration has no value for style directory, assuming default of '%s'", c.StyleDir)
	}
	if _, err := os.Stat(c.StyleDir); err != nil {
		return fmt.Errorf("Could not find style directory '%s': %s", c.StyleDir, err)
	}

	if strings.TrimSpace(c.TemplateDir) == "" {
		c.TemplateDir = def.TemplateDir
		log.Printf("Configuration has no value for template directory, assuming default of '%s'", c.TemplateDir)
	}
	if _, err := os.Stat(c.TemplateDir); err != nil {
		return fmt.Errorf("Could not find template directory '%s': %s", c.TemplateDir, err)
	}

	if c.Port == 0 {
		c.Port = def.Port
		log.Printf("Configuration port unset; using default %d", def.Port)
	}

	return nil
}

func defaults() configuration {
	return configuration{
		RootDir:     ".",
		StyleDir:    "css",
		TemplateDir: "templates",
		Port:        6060,
	}
}

// Create a new configion file and return the default configuration.
func createNewConfigFile() configuration {
	c := defaults()

	//TODO finish this
	//xdgdir.Config.Create(configFileName)

	return c
}

func loadConfiguration() (configuration, error) {
	cfg := defaults()

	var bytes []byte
	log.Printf("Attempting to read file at %s", configFileName)
	file, err := config.Open(configFileName)
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
