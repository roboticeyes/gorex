package main

import (
	"fmt"
)

type controller struct {
	Config *Configuration
}

// NewController creates a new controller
func NewController(config *Configuration) Controller {
	return &controller{
		Config: config,
	}
}

func (c *controller) Connect() error {
	return fmt.Errorf(c.Config.Default)
}
