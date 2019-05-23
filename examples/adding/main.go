package main

// This example shows how to list all project of the user

import (
	"fmt"

	"github.com/roboticeyes/gorex/http/creator"
)

var (
	domain       = "rex.robotic-eyes.com"
	clientID     = ""
	clientSecret = ""
)

func main() {

	controller := creator.NewController(domain)

	// Authenticate
	err := controller.Authenticate(clientID, clientSecret)
	if err != nil {
		fmt.Println("Cannot authenticate: ", err)
		panic(err)
	}

	// Create a new project
	project, err := controller.CreateProject("First Project 123")
	if err != nil {
		fmt.Println("Cannot create project: ", err)
		panic(err)
	}

	fmt.Println("Successfully created project ", project.Name)
}
