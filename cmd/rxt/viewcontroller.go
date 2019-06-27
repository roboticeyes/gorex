package main

import (
	"github.com/roboticeyes/gorex/http/creator"
	"github.com/roboticeyes/gorex/http/creator/listing"
)

// ViewController handles all requests to the rexOS interface
type ViewController struct {
	config     *Configuration
	controller *creator.Controller
	user       listing.User // cached user information
}

// NewViewController creates new view controller
func NewViewController(config *Configuration) *ViewController {
	return &ViewController{
		config:     config,
		controller: creator.NewController(config.Active.Domain),
	}
}

// Connect to rexOS
func (c *ViewController) Connect() (string, error) {
	err := c.controller.Authenticate(c.config.Active.ClientID, c.config.Active.ClientSecret)
	if err != nil {
		return "", err
	}
	c.user = c.controller.GetUser()
	return c.user.Username, nil
}

// GetConfiguration gets the configuration
func (c *ViewController) GetConfiguration() *Configuration {
	return c.config
}

// GetUserName returns the user name of the authenticated user
func (c *ViewController) GetUserName() string {
	return c.user.Username
}

// GetUserID returns the user ID of the authenticated user
func (c *ViewController) GetUserID() string {
	return c.user.UserID
}

// GetProjects delivers a list of all projects related to the user
func (c *ViewController) GetProjects() ([]listing.Project, error) {
	return c.controller.GetProjects(100, 0)
}

// GetProjectFiles delivers the list of project files for a given project
func (c *ViewController) GetProjectFiles() ([]listing.ProjectFile, error) {
	var dummy []listing.ProjectFile

	dummy = append(dummy, listing.ProjectFile{Name: "File 1", Type: "", FileSize: 1234422})
	dummy = append(dummy, listing.ProjectFile{Name: "File 2", Type: "rex", FileSize: 34422})
	return dummy, nil
}
