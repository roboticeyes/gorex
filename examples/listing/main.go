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

	// Print user information
	fmt.Println(controller.GetUser())

	// Get the first 10 projects
	projects, err := controller.GetProjects(10, 0)
	if err != nil {
		fmt.Println("Cannot get project list: ", err)
	}

	fmt.Printf("Found %d projects ...\n\n", len(projects))

	for _, p := range projects {
		fmt.Println(p.Name)
	}
}
