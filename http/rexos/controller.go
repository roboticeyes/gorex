package rexos

import (
	"github.com/roboticeyes/gorex/http/rexos/listing"
	"github.com/roboticeyes/gorex/http/rexos/rest"
)

// Controller provides a high level interface for rexOS operations
type Controller interface {
	Authenticate(clientID, clientSecret string) error
	GetProjects() ([]listing.Project, error)
}

type controller struct {
	rexClient *rest.RexClient
	listing   listing.Service
}

// NewController creates a new rexOS controller for easy REX interactions
func NewController(domain string) Controller {
	var c controller
	c.rexClient = rest.NewRexClient(domain)
	return &c
}

func (c *controller) Authenticate(clientID, clientSecret string) error {

	_, err := c.rexClient.ConnectWithClientCredentials(clientID, clientSecret)
	if err != nil {
		return err
	}

	// Get all services
	restData := rest.NewDataProvider(c.rexClient)
	c.listing = listing.NewService(restData)

	return nil
}

// GetProjects returns all projects of the user
func (c *controller) GetProjects() ([]listing.Project, error) {
	return c.listing.GetProjects()
}
