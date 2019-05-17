package main

import (
	"github.com/roboticeyes/gorex/http/rexos"
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// ViewController handles all requests to the rexOS interface
type ViewController struct {
	config        *Configuration
	rexController rexos.Controller
}

// NewViewController creates new view controller
func NewViewController(config *Configuration) *ViewController {
	return &ViewController{
		config:        config,
		rexController: rexos.NewController(config.APIUrl),
	}
}

// Connect to rexOS
func (c *ViewController) Connect() (string, error) {
	err := c.rexController.Authenticate(c.config.Active.ClientID, c.config.Active.ClientSecret)
	if err != nil {
		return "", err
	}
	return c.rexController.GetUserInformation().Username, nil
}

// GetConfiguration gets the configuration
func (c *ViewController) GetConfiguration() *Configuration {
	return c.config
}

// GetUserName returns the user name of the authenticated user
func (c *ViewController) GetUserName() string {
	return c.rexController.GetUserInformation().Username
}

// GetUserID returns the user ID of the authenticated user
func (c *ViewController) GetUserID() string {
	return c.rexController.GetUserInformation().UserID
}

// GetProjects delivers a list of all projects related to the user
func (c *ViewController) GetProjects() ([]listing.Project, error) {
	return c.rexController.GetProjects()
}

// GetProjectFiles delivers the list of project files for a given project
func (c *ViewController) GetProjectFiles() ([]listing.ProjectFile, error) {
	var dummy []listing.ProjectFile

	dummy = append(dummy, listing.ProjectFile{Name: "File 1", Type: "", FileSize: 1234422})
	dummy = append(dummy, listing.ProjectFile{Name: "File 2", Type: "rex", FileSize: 34422})
	return dummy, nil
}
