package main

import (
	"fmt"
	"os"
)

var (
	configPath string
	configFile string
)

func init() {
	homeDir, _ := os.UserHomeDir()
	configPath = homeDir + "/.config/rxt"
	configFile = configPath + string(os.PathSeparator) + "config.json"
}

func exit(reason string, code int) {
	fmt.Println(reason)
	os.Exit(code)
}

func main() {

	var config *Configuration

	// check if config file is a parameter
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else if len(os.Args) < 2 {
		if stat, err := os.Stat(configPath); err != nil || !stat.IsDir() {
			err = os.MkdirAll(configPath, os.ModePerm)
			exit("Please setup config file in "+configPath, 1)
		}
	}
	config, err := LoadConfiguration(configFile)
	if err != nil {
		exit(err.Error(), 1)
	}

	NewTui(NewViewController(config)).Run()
}
