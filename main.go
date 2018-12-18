package main

import (
	"fmt"
	"github.com/roboticeyes/gorex/gorex"
	"os"
)

func main() {
	fmt.Println("Welcome at rexos ...")

	baseURL := os.Getenv("REX_BASEURL")
	clientID := os.Getenv("REX_CLIENT_ID")
	clientSecret := os.Getenv("REX_CLIENT_SECRET")

	cli := gorex.NewRexClient(baseURL)

	token, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	if token == nil {
		fmt.Println("Token is not retrieved")
	}
	if err != nil {
		fmt.Println("Error during connection")
	}

	// Project service
	projectService := gorex.NewProjectService(cli)

	name := "test"
	owner := "fb1b3be2-1783-4aa8-9ed4-b0d118b80fac"
	project, err := projectService.FindByNameAndOwner(name, owner)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println(project)
}
