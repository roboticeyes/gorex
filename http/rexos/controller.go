// Copyright 2019 Robotic Eyes. All rights reserved.

// Package rexos provides two different levels of accessing rexOS.
// The low-level REST implementation where you get the complete
// HAL interface, and a high-level controller interface which allows
// to easily access the most important information.
//
// Please see the examples folder for some samples
package rexos

import (
	"github.com/roboticeyes/gorex/http/rexos/listing"
	"github.com/roboticeyes/gorex/http/rexos/rest"
)

// Controller provides a high level interface for rexOS operations
type Controller interface {
	Authenticate(clientID, clientSecret string) error
	GetProjects() ([]listing.Project, error)
	GetUserName() string
}

type controller struct {
	rexClient       *rest.RexClient
	listing         listing.Service
	userInformation listing.User
}

// NewController creates a new rexOS controller for easy REX interactions
func NewController(domain string) Controller {
	return &controller{
		rexClient: rest.NewRexClient(domain),
	}
}

func (c *controller) Authenticate(clientID, clientSecret string) error {

	_, err := c.rexClient.ConnectWithClientCredentials(clientID, clientSecret)
	if err != nil {
		return err
	}

	// Get all services
	restData := rest.NewDataProvider(c.rexClient)
	c.listing = listing.NewService(restData)

	// Get and cache user information
	c.userInformation, err = c.listing.GetUserInformation()
	if err != nil {
		return err
	}

	return nil
}

// GetProjects returns all projects of the user
func (c *controller) GetProjects() ([]listing.Project, error) {
	return c.listing.GetProjects()
}

func (c *controller) GetUserName() string {
	return c.userInformation.Username
}
