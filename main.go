package main

import (
	"fmt"
	"github.com/roboticeyes/gorex/gorex"
	"log"
	"os"
)

func main() {
	fmt.Println("Welcome at rexos ...")

	baseURL := os.Getenv("REX_BASEURL")
	clientID := os.Getenv("REX_CLIENT_ID")
	clientSecret := os.Getenv("REX_CLIENT_SECRET")

	if len(baseURL) == 0 {
		log.Fatal("Please set REX_BASEURL")
	}

	// Setup REX client connection
	cli := gorex.NewRexClient(baseURL)

	_, err := cli.ConnectWithClientCredentials(clientID, clientSecret)
	if err != nil {
		log.Fatal("Error during connection", err)
	}

	// Setup services

	projectService := gorex.NewProjectService(cli)
	userService := gorex.NewUserService(cli)

	user, err := userService.GetCurrentUser()
	if err != nil {
		fmt.Println("Cannot get user", err)
	}
	fmt.Println(user)

	name := "test"
	owner := user.UserID
	project, err := projectService.FindByNameAndOwner(name, owner)

	if err != nil {
		fmt.Println("Cannot get project", err)
	}

	fmt.Println(project)
}
