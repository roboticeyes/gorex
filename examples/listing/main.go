package main

// This example shows how to list all project of the user

import (
	"fmt"

	"github.com/roboticeyes/gorex/http/creator"
)

var (
	rexDomain    = "https://rex.robotic-eyes.com"
	clientID     = ""
	clientSecret = ""
)

func main() {

	controller := creator.NewController(rexDomain)

	// Authenticate
	err := controller.Authenticate(clientID, clientSecret)
	if err != nil {
		fmt.Println("Cannot authenticate to rexOS: ", err)
	}

	// Get projects
	projects, err := controller.GetProjects()
	if err != nil {
		fmt.Println("Cannot get project list: ", err)
	}

	fmt.Printf("Found %d projects ...\n\n", len(projects))

	for _, p := range projects {
		fmt.Println(p.Name)
	}
}
