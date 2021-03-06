// Copyright 2019 Robotic Eyes. All rights reserved.

// Package creator provides a high-level access to the rexOS resources.
// The low-level core REST implementation where you get the complete
// HAL interface.
//
// Please see the examples folder for some samples.
package creator

import (
	"context"
	"net/http"
	"sync"

	"github.com/roboticeyes/gorex/http/core"
	"github.com/roboticeyes/gorex/http/creator/adding"
	"github.com/roboticeyes/gorex/http/creator/listing"
)

const (
	rexAPIScheme = "https://"
	rexAPIPrefix = "/api/v2"
)

// Controller provides a high level interface for rexOS operations
type Controller struct {
	config  core.RexConfig
	ctx     context.Context
	wg      sync.WaitGroup
	client  *core.Client
	listing listing.Service
	adding  adding.Service
	user    listing.User
}

// NewController creates a new rexOS controller for easy REX interactions.
// The domain should not contain any schema information, should be in form of rex.robotic-eyes.com
func NewController(domain string) *Controller {
	config := core.RexConfig{
		ProjectResourceURL: rexAPIScheme + domain + rexAPIPrefix,
		UserResourceURL:    rexAPIScheme + domain + rexAPIPrefix,
		AuthenticationURL:  rexAPIScheme + domain,
	}
	c := &Controller{
		config: config,
		ctx:    context.Background(),
	}
	// make sure to wait till authentication is done
	c.wg.Add(1)
	return c
}

// Authenticate the user with the given clientID and clientSecret
func (c *Controller) Authenticate(clientID, clientSecret string) error {

	defer c.wg.Done()
	token, err := core.Authenticate(c.config.AuthenticationURL, clientID, clientSecret)
	if err != nil {
		return err
	}

	// Create a context with access token to get user information
	var contextData core.ContextData
	contextData.AccessToken = token
	c.ctx = context.WithValue(c.ctx, core.ContextDataKey, contextData)

	// Get all services
	rexAccessor := core.NewController(c.config)
	c.listing = listing.NewService(rexAccessor)
	c.adding = adding.NewService(rexAccessor)

	// Get and cache user information
	c.user, err = c.listing.GetUser(c.ctx)
	if err != nil {
		return err
	}

	// Set a fresh context with accesstoken and user id
	contextData.UserID = c.user.UserID
	c.ctx = context.WithValue(context.Background(), core.ContextDataKey, contextData)

	return nil
}

// GetProjects returns all projects of the user
func (c *Controller) GetProjects(size, page uint64) ([]listing.Project, error) {
	c.wg.Wait()
	res, ret := c.listing.GetProjects(c.ctx, size, page)
	if ret.Code != http.StatusOK {
		return res, ret
	}
	return res, nil
}

// GetUser returns the user information
func (c *Controller) GetUser() listing.User {
	c.wg.Wait()
	return c.user
}

// CreateProject creates a new CreatorProject
func (c *Controller) CreateProject(name string) (*adding.Project, error) {
	c.wg.Wait()
	return c.adding.CreateProject(c.ctx, name)
}
