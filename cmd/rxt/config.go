package main

import (
	"encoding/json"
	"os"
)

// Configuration is the top level configuration file
type Configuration struct {
	Default      string        `json:"default"`
	Environments []Environment `json:"environments"`
	Active       Environment   `json:"-"`
}

// Environment describes a specific rexOS environment
type Environment struct {
	Name         string `json:"name"`
	Domain       string `json:"domain"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// LoadConfiguration loads the total configuration from a file
func LoadConfiguration(fileName string) (*Configuration, error) {
	config := new(Configuration)
	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	// find and set active environment
	for _, e := range config.Environments {
		if e.Name == config.Default {
			config.Active = e
			break
		}
	}
	return config, nil
}
